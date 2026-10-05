// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package wayland

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"testing"
)

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
