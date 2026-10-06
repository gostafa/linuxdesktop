// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package dbusconn

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/gostafa/linuxdesktop/internal/port"
)

type fakeObject struct {
	dbus.BusObject
	call *dbus.Call
}

func (object fakeObject) CallWithContext(context.Context, string, dbus.Flags, ...any) *dbus.Call {
	return object.call
}

type fakeConnection struct{ object fakeObject }

func (conn fakeConnection) BusObject() dbus.BusObject { return conn.object }

func (conn fakeConnection) Object(string, dbus.ObjectPath) dbus.BusObject { return conn.object }

func TestBusOperations(t *testing.T) {
	t.Parallel()
	failure := errors.New("offline")
	for _, tc := range []struct {
		body            []any
		err, connectErr error
	}{
		{connectErr: failure},
		{err: failure},
		{body: []any{"wrong type"}},
		{body: []any{true}},
		{body: []any{dbus.MakeVariant("value")}},
		{body: []any{map[string]dbus.Variant{"key": dbus.MakeVariant("value")}}},
	} {
		bus := connectionBus[port.BusKind, fakeConnection, *port.Object, *port.PropertyQuery]{
			source: connectionSource[port.BusKind, fakeConnection]{
				Acquire: func(port.BusKind) (fakeConnection, error) {
					return fakeConnection{
						fakeObject{call: &dbus.Call{Body: tc.body, Err: tc.err}},
					}, tc.connectErr
				},
			},
		}
		ctx := t.Context()
		object := &port.Object{Destination: "fixture", Path: "/fixture"}
		query := &port.PropertyQuery{Object: *object, Interface: "fixture", Name: "Value"}
		owner, ownerErr := bus.HasOwner(ctx, port.SessionBus, "fixture")
		text, textErr := bus.Introspect(ctx, port.SessionBus, object)
		value, valueErr := bus.Property(ctx, port.SessionBus, query)
		values, valuesErr := bus.Properties(ctx, port.SessionBus, query)
		if tc.connectErr != nil || tc.err != nil {
			for _, err := range []error{ownerErr, textErr, valueErr, valuesErr} {
				if !errors.Is(err, failure) {
					t.Fatal(err)
				}
			}
			continue
		}
		switch tc.body[0].(type) {
		case bool:
			if ownerErr != nil || !owner {
				t.Fatal(owner, ownerErr)
			}
		case string:
			if textErr != nil || text != "wrong type" {
				t.Fatal(text, textErr)
			}
		case dbus.Variant:
			if valueErr != nil || value.Value != "value" {
				t.Fatal(value, valueErr)
			}
		case map[string]dbus.Variant:
			if valuesErr != nil || values["key"] != "value" {
				t.Fatal(values, valuesErr)
			}
		}
	}
}

type transport struct{ closed error }

func (transport) Read([]byte) (int, error) { return 0, io.EOF }

func (transport) Write(data []byte) (int, error) { return len(data), nil }

func (conn transport) Close() error { return conn.closed }

func TestConnectionLifetime(t *testing.T) {
	t.Parallel()
	failure := errors.New("close failure")
	conn, err := dbus.NewConn(transport{closed: failure})
	if err != nil {
		t.Fatal(err)
	}
	calls := 0
	cache := connections[port.BusKind, *dbus.Conn]{
		open: func(context.Context, port.BusKind) (*dbus.Conn, error) { calls++; return conn, nil },
	}
	for range 2 {
		got, err := connect(t.Context(), &cache, port.SessionBus)
		if err != nil || got != conn {
			t.Fatal(got, err)
		}
	}
	if calls != 1 {
		t.Fatal("connection not cached")
	}
	if _, err = connect(t.Context(), &cache, port.BusKind(255)); err == nil {
		t.Fatal("unknown bus accepted")
	}
	if err = closeConnections(&cache); !errors.Is(err, failure) || cache.conns[0] != nil {
		t.Fatal(err)
	}
	if err = closeConnections(&cache); err != nil {
		t.Fatal(err)
	}
	if err = New(t.Context()).Close(); err != nil {
		t.Fatal(err)
	}
	bus := Bus{
		source: connectionSource[port.BusKind, *dbus.Conn]{
			Release: func() error { return failure },
		},
	}
	if err = bus.Close(); !errors.Is(err, failure) {
		t.Fatal("close error was lost", err)
	}
}

// TestOpenConnection exercises authentication and Hello over a real local socket,
// independently of any desktop session bus installed on the test host.
func TestOpenConnection(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bus")
	listener, err := net.Listen("unix", path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = listener.Close() })
	t.Setenv("DBUS_SESSION_BUS_ADDRESS", "unix:path="+path)
	done := make(chan error, 1)
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			done <- err
			return
		}
		defer conn.Close()
		if err = conn.SetDeadline(time.Now().Add(5 * time.Second)); err == nil {
			err = serveBusHello(conn)
		}
		done <- err
	}()
	conn, err := openConnection(t.Context(), port.SessionBus)
	if err != nil {
		t.Fatal(err)
	}
	if len(conn.Names()) == 0 || conn.Names()[0] != ":1.1" {
		t.Fatal("Hello did not establish the connection identity", conn.Names())
	}
	if err = conn.Close(); err != nil {
		t.Fatal(err)
	}
	if err = <-done; err != nil {
		t.Fatal(err)
	}
}

func serveBusHello(conn net.Conn) error {
	reader := bufio.NewReader(conn)
	if _, err := reader.ReadByte(); err != nil {
		return err
	}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		line = strings.TrimSpace(line)
		if line == "BEGIN" {
			break
		}
		response := "ERROR unsupported\r\n"
		if line == "AUTH" {
			response = "REJECTED EXTERNAL\r\n"
		} else if strings.HasPrefix(line, "AUTH EXTERNAL") {
			response = "OK 0123456789abcdef0123456789abcdef\r\n"
		}
		if _, err = io.WriteString(conn, response); err != nil {
			return err
		}
	}
	request, err := dbus.DecodeMessage(reader)
	if err != nil {
		return err
	}
	if request.Headers[dbus.FieldMember].Value() != "Hello" {
		return fmt.Errorf("unexpected method: %v", request.Headers[dbus.FieldMember])
	}
	reply := dbus.Message{
		Type: dbus.TypeMethodReply,
		Headers: map[dbus.HeaderField]dbus.Variant{
			dbus.FieldReplySerial: dbus.MakeVariant(request.Serial()),
			dbus.FieldSignature:   dbus.MakeVariant(dbus.SignatureOf(":1.1")),
		},
		Body: []any{":1.1"},
	}
	if err = reply.EncodeTo(conn, binary.LittleEndian); err != nil {
		return err
	}
	_, err = io.Copy(io.Discard, reader)
	return err
}
