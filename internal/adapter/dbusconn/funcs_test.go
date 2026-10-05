// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package dbusconn

import (
	"context"
	"errors"
	"io"
	"testing"

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

func (conn fakeConnection) BusObject() dbus.BusObject                     { return conn.object }
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
		bus := connectionBus[port.BusKind, fakeConnection]{
			connect: func(port.BusKind) (fakeConnection, error) {
				return fakeConnection{
					fakeObject{call: &dbus.Call{Body: tc.body, Err: tc.err}},
				}, tc.connectErr
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
			if valueErr != nil || value != "value" {
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

func (transport) Read([]byte) (int, error)       { return 0, io.EOF }
func (transport) Write(data []byte) (int, error) { return len(data), nil }
func (conn transport) Close() error              { return conn.closed }

func TestConnectionLifetime(t *testing.T) {
	t.Parallel()
	failure := errors.New("close failure")
	conn, err := dbus.NewConn(transport{closed: failure})
	if err != nil {
		t.Fatal(err)
	}
	calls := 0
	cache := connections{
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
	bus := Bus{release: func() error { return failure }}
	if err = bus.Close(); !errors.Is(err, failure) {
		t.Fatal("close error was lost", err)
	}
}

func TestHostConnectionFailures(t *testing.T) {
	t.Setenv("DBUS_SYSTEM_BUS_ADDRESS", "unix:path=/nonexistent-linuxdesktop-test-bus")
	t.Setenv("DBUS_SESSION_BUS_ADDRESS", "unix:path=/nonexistent-linuxdesktop-test-bus")
	bus := New(nil)
	for _, kind := range []port.BusKind{port.SessionBus, port.SystemBus} {
		if _, err := bus.HasOwner(t.Context(), kind, "fixture"); err == nil {
			t.Fatal("connected to missing bus")
		}
	}
	if err := bus.Close(); err != nil {
		t.Fatal(err)
	}
}
