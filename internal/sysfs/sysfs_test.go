// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package sysfs

import (
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
)

type emptyReader struct{}

func (emptyReader) Read([]byte) (int, error) { return 0, nil }

func TestKeyValueFiles(t *testing.T) {
	data := []byte(
		"\n# comment=x\ninvalid\nNAME='quoted'\nEMPTY=\nDOUBLE=\"two\"\nRAW=\"mismatch'\n",
	)
	got := map[string]string{}
	Each(data, func(key, value string) bool { got[key] = value; return true })
	want := map[string]string{"NAME": "quoted", "EMPTY": "", "DOUBLE": "two", "RAW": "\"mismatch'"}
	if !reflect.DeepEqual(got, want) {
		t.Fatal(got)
	}
	count := 0
	Each(data, func(string, string) bool { count++; return false })
	if count != 1 {
		t.Fatal(count)
	}
	if Field(data, "NAME") != "quoted" || Field(data, "absent") != "" ||
		unquote([]byte("raw")) != "raw" {
		t.Fatal("fields")
	}
}

func TestBufferLimitsAndFreshBuffer(t *testing.T) {
	buf := make([]byte, 2)
	read, err := drain(strings.NewReader("growing buffer"), &buf)
	if err != nil || string(buf[:read]) != "growing buffer" {
		t.Fatal(read, err)
	}
	if read, err = drain(emptyReader{}, &buf); err != nil || read != 0 {
		t.Fatal(read, err)
	}
	buf = make([]byte, maxFileSize)
	if read, err = step(
		strings.NewReader("ignored"),
		&buf,
		maxFileSize,
	); !errors.Is(err, io.EOF) ||
		read != 0 {
		t.Fatal(read, err)
	}
	if got := take(); len(*got) != scratchSize {
		t.Fatal("fresh buffer size")
	}
}

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
