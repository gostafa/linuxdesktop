// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package sysfs

import (
	"errors"
	"io"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

type emptyReader struct{}

func (emptyReader) Read([]byte) (int, error) { return 0, nil }

func TestFilesAndLinks(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "value")
	if err := os.WriteFile(path, []byte("  content\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if got, err := String(path); err != nil || got != "content" {
		t.Fatal(got, err)
	}
	if got := Trimmed(path); got != "content" {
		t.Fatal(got)
	}
	if got, err := Bytes(path); err != nil || string(got) != "  content\n" {
		t.Fatal(got, err)
	}
	missing := filepath.Join(dir, "missing")
	if _, err := String(missing); !errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	}
	if _, err := Bytes(missing); !errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	}
	if Trimmed(missing) != "" || IsSocket(path) || IsSocket(missing) {
		t.Fatal("missing/non-socket path")
	}
	if _, err := Bytes(dir); err == nil {
		t.Fatal("directory read succeeded")
	}
	if _, err := DirNames(missing); !errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	}
	if _, err := DirNames(path); err == nil {
		t.Fatal("file listed as directory")
	}
	if names, err := DirNames(dir); err != nil || len(names) != 1 || names[0] != "value" {
		t.Fatal(names, err)
	}
	link := filepath.Join(dir, "link")
	if err := os.Symlink("target/base", link); err != nil {
		t.Fatal(err)
	}
	if LinkBase(link) != "base" || LinkBase(missing) != "" {
		t.Fatal("symlink resolution")
	}
	socket := filepath.Join(dir, "socket")
	listener, err := net.Listen("unix", socket)
	if err != nil {
		t.Fatal(err)
	}
	if !IsSocket(socket) {
		t.Fatal("socket missed")
	}
	if err = listener.Close(); err != nil {
		t.Fatal(err)
	}
	if Path(dir, "/etc/file") != filepath.Join(dir, "etc/file") {
		t.Fatal("root resolution")
	}
}

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
