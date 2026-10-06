// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package wayland

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"path/filepath"
	"testing"
	"time"

	"github.com/gostafa/linuxdesktop/internal/domain"
)

func frame(object, opcode uint32, body []byte) []byte {
	out := make([]byte, headerSize+len(body))
	binary.NativeEndian.PutUint32(out, object)
	binary.NativeEndian.PutUint32(out[wordBytes:], uint32(len(out))<<opcodeBits|opcode)
	copy(out[headerSize:], body)
	return out
}

func globalBody(name string) []byte {
	length := len(name) + 1
	body := make([]byte, wordBytes*3+align(length))
	binary.NativeEndian.PutUint32(body, 7)
	binary.NativeEndian.PutUint32(body[wordBytes:], uint32(length))
	copy(body[wordBytes*2:], name)
	binary.NativeEndian.PutUint32(body[len(body)-wordBytes:], 3)
	return body
}

func TestRegistryMessages(t *testing.T) {
	t.Parallel()
	body := globalBody("wl_compositor")
	event := frame(objRegistry, 0, body)
	reg := newRegistry()
	copy(reg.buf, event[:headerSize])
	reg.filled = headerSize
	if err := parse(reg); err != nil || reg.filled != headerSize {
		t.Fatal(err)
	}
	copy(reg.buf, event)
	reg.filled = len(event)
	if err := parse(reg); err != nil || len(reg.globals) != 1 || reg.filled != 0 {
		t.Fatal(reg, err)
	}
	if !Has(reg.globals, "wl_compositor") || Has(reg.globals, "missing") {
		t.Fatal("registry lookup")
	}
	if err := dispatch(reg, &wireEvent{object: 99}); err != nil {
		t.Fatal(err)
	}
	if err := dispatch(reg, &wireEvent{object: objRegistry, opcode: 1}); err != nil {
		t.Fatal(err)
	}
	if err := dispatch(reg, &wireEvent{object: objRegistry, body: nil}); err != nil {
		t.Fatal(err)
	}
	if err := dispatch(reg, &wireEvent{object: objCallback, opcode: 1}); err != nil {
		t.Fatal(err)
	}
	if err := dispatch(reg, &wireEvent{object: 1, opcode: 1}); err != nil {
		t.Fatal(err)
	}
	if err := dispatch(reg, &wireEvent{object: 1}); !errors.Is(err, ErrProtocol) {
		t.Fatal(err)
	}
	if err := parseError(body); err == nil {
		t.Fatal("protocol error lost")
	}
	if _, ok := parseGlobal(make([]byte, requestSize)); ok {
		t.Fatal("empty global accepted")
	}
	if text([]byte("x"), 2) != "" {
		t.Fatal("truncated text")
	}
	if err := stalled(1, nil); err != nil {
		t.Fatal(err)
	}
	if err := stalled(0, nil); !errors.Is(err, io.ErrNoProgress) {
		t.Fatal(err)
	}
	if err := step(newRegistry(), bytes.NewReader(event)); err != nil {
		t.Fatalf("a complete global is progress even after compaction: %v", err)
	}
	if socketAt(filepath.Join(t.TempDir(), "missing")) != "" {
		t.Fatal("nonexistent socket accepted")
	}
	bad := frame(99, 0, nil)
	binary.NativeEndian.PutUint32(bad[wordBytes:], wordBytes<<opcodeBits)
	copy(reg.buf, bad)
	reg.filled = len(bad)
	if err := parse(reg); !errors.Is(err, ErrProtocol) {
		t.Fatal(err)
	}
	reg.buf = make([]byte, headerSize)
	reg.filled = headerSize
	if err := grow(reg); err != nil || len(reg.buf) != headerSize*2 {
		t.Fatal(err)
	}
	reg.buf = make([]byte, maxMessageBytes)
	reg.filled = maxMessageBytes
	if err := step(reg, bytes.NewReader(nil)); !errors.Is(err, ErrProtocol) {
		t.Fatal(err)
	}
}

type fakeConn struct {
	net.Conn
	deadlineErr, writeErr error
	read                  io.Reader
}

func (conn fakeConn) SetDeadline(time.Time) error { return conn.deadlineErr }

func (conn fakeConn) Write(data []byte) (int, error) {
	if conn.writeErr != nil {
		return 0, conn.writeErr
	}
	return len(data), nil
}

func (conn fakeConn) Read(data []byte) (int, error) { return conn.read.Read(data) }

func TestConverseFailures(t *testing.T) {
	t.Parallel()
	failure := errors.New("connection failed")
	ctx, cancel := context.WithTimeout(t.Context(), time.Second)
	defer cancel()
	for _, conn := range []fakeConn{{deadlineErr: failure}, {writeErr: failure}, {read: bytes.NewReader(nil)}} {
		if _, err := converse(ctx, conn); err == nil {
			t.Fatal("connection failure lost")
		}
	}
	conn := fakeConn{read: bytes.NewReader(frame(objCallback, 0, nil))}
	if _, err := converse(t.Context(), conn); err != nil {
		t.Fatal(err)
	}
	if _, err := listen(t.Context(), filepath.Join(t.TempDir(), "missing")); err == nil {
		t.Fatal("missing socket")
	}
}

func TestWaylandSocket(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, defaultDisplay)
	listener, err := net.Listen("unix", path)
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	server := make(chan error, 1)
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			server <- err
			return
		}
		defer conn.Close()
		request := make([]byte, handshakeSize)
		if _, err = io.ReadFull(conn, request); err != nil {
			server <- err
			return
		}
		// A global and its completion arrive in separate reads.
		if _, err = conn.Write(frame(objRegistry, 0, globalBody("wl_compositor"))); err != nil {
			server <- err
			return
		}
		_, err = conn.Write(frame(objCallback, 0, nil))
		server <- err
	}()
	ctx, cancel := context.WithTimeout(t.Context(), time.Second)
	defer cancel()
	env := domain.Env{RuntimeDir: dir, WaylandSocket: "9"}
	info, err := New().Wayland(ctx, &env)
	if err != nil || info == nil || info.SocketFD != 9 || len(info.Globals) != 1 ||
		info.Globals[0].Interface != "wl_compositor" {
		t.Fatal(info, err)
	}
	if err = <-server; err != nil {
		t.Fatal(err)
	}
	env.WaylandDisplay = path
	if SocketPath(&env) != path {
		t.Fatal("absolute socket")
	}
	if SocketPath(&domain.Env{}) != "" {
		t.Fatal("unset runtime directory")
	}
	// A listening socket that closes before replying produces a probe error.
	go func() {
		conn, err := listener.Accept()
		if err == nil {
			conn.Close()
		}
	}()
	if _, err = New().Wayland(ctx, &env); err == nil {
		t.Fatal("peer close lost")
	}
	if err = listener.Close(); err != nil {
		t.Fatal(err)
	}
	if info, err = New().Wayland(t.Context(), &domain.Env{}); err != nil || info != nil {
		t.Fatal(info, err)
	}
	for _, raw := range []string{"", "bad", "-1", "0"} {
		got, err := New().Wayland(t.Context(), &domain.Env{WaylandSocket: raw})
		if err != nil || (got != nil) != (raw == "0") {
			t.Fatal(raw, got, err)
		}
	}
}

type handshakeWriter struct {
	n   int
	err error
}

func (writer handshakeWriter) Write([]byte) (int, error) {
	return writer.n, writer.err
}

func TestHandshakeWrites(t *testing.T) {
	failure := errors.New("write failed")
	for _, test := range []struct {
		name   string
		writer handshakeWriter
		want   error
	}{
		{"complete", handshakeWriter{n: handshakeSize}, nil},
		{"short", handshakeWriter{n: handshakeSize - 1}, io.ErrShortWrite},
		{"failed", handshakeWriter{n: 1, err: failure}, failure},
	} {
		t.Run(test.name, func(t *testing.T) {
			if err := handshake(test.writer); !errors.Is(err, test.want) {
				t.Fatalf("handshake error = %v, want %v", err, test.want)
			}
		})
	}
}

func TestReadGlobalsCompletion(t *testing.T) {
	frame := make([]byte, headerSize)
	binary.NativeEndian.PutUint32(frame, objCallback)
	binary.NativeEndian.PutUint32(frame[wordBytes:], headerSize<<opcodeBits|wireZero)
	if _, err := readGlobals(bytes.NewReader(frame)); err != nil {
		t.Fatalf("completed registry returned an error: %v", err)
	}
	if _, err := readGlobals(bytes.NewReader(nil)); !errors.Is(err, io.EOF) {
		t.Fatalf("incomplete registry error = %v, want EOF", err)
	}
}
