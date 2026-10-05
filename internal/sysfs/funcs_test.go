// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package sysfs

import (
	"errors"
	"io"
	"strings"
	"testing"
)

type partialReader struct{ err error }

func (reader partialReader) Read(buf []byte) (int, error) {
	return copy(buf, "partial"), reader.err
}

func TestDrainWrappedReadErrors(t *testing.T) {
	failure := errors.New("read failed")
	for _, test := range []struct {
		name   string
		reader io.Reader
		want   string
		err    error
	}{
		{"EOF", strings.NewReader("complete"), "complete", nil},
		{"partial EOF", partialReader{err: io.EOF}, "partial", nil},
		{"partial failure", partialReader{err: failure}, "partial", failure},
	} {
		t.Run(test.name, func(t *testing.T) {
			buf := make([]byte, 16)
			n, err := drain(test.reader, &buf)
			if string(buf[:n]) != test.want || !errors.Is(err, test.err) {
				t.Fatalf("drain = %q, %v; want %q, %v", buf[:n], err, test.want, test.err)
			}
		})
	}
}
