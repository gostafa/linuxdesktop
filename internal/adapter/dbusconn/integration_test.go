// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package dbusconn

import (
	"bufio"
	"encoding/binary"
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
