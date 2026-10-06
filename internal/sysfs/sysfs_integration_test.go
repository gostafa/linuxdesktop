// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package sysfs_test

import (
	"errors"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/gostafa/linuxdesktop/internal/sysfs"
)

func TestFilesAndLinks(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "value")
	if err := os.WriteFile(path, []byte("  content\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if got, err := sysfs.String(path); err != nil || got != "content" {
		t.Fatal(got, err)
	}
	if got := sysfs.Trimmed(path); got != "content" {
		t.Fatal(got)
	}
	if got, err := sysfs.Bytes(path); err != nil || string(got) != "  content\n" {
		t.Fatal(got, err)
	}
	missing := filepath.Join(dir, "missing")
	if _, err := sysfs.String(missing); !errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	}
	if _, err := sysfs.Bytes(missing); !errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	}
	if sysfs.Trimmed(missing) != "" || sysfs.IsSocket(path) || sysfs.IsSocket(missing) {
		t.Fatal("missing/non-socket path")
	}
	if _, err := sysfs.Bytes(dir); err == nil {
		t.Fatal("directory read succeeded")
	}
	if _, err := sysfs.DirNames(missing); !errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	}
	if _, err := sysfs.DirNames(path); err == nil {
		t.Fatal("file listed as directory")
	}
	if names, err := sysfs.DirNames(dir); err != nil || len(names) != 1 || names[0] != "value" {
		t.Fatal(names, err)
	}
	link := filepath.Join(dir, "link")
	if err := os.Symlink("target/base", link); err != nil {
		t.Fatal(err)
	}
	if sysfs.LinkBase(link) != "base" || sysfs.LinkBase(missing) != "" {
		t.Fatal("symlink resolution")
	}
	socket := filepath.Join(dir, "socket")
	listener, err := net.Listen("unix", socket)
	if err != nil {
		t.Fatal(err)
	}
	if !sysfs.IsSocket(socket) {
		t.Fatal("socket missed")
	}
	if err = listener.Close(); err != nil {
		t.Fatal(err)
	}
	if sysfs.Path(dir, "/etc/file") != filepath.Join(dir, "etc/file") {
		t.Fatal("root resolution")
	}
}
