// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package procscan

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func process(t *testing.T, root, pid, comm, uid string) {
	t.Helper()
	dir := filepath.Join(root, pid)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	for name, data := range map[string]string{commFile: comm, statusFile: "Name: test\nUid:\t" + uid + " 0 0 0\n"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(data), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

func TestProcesses(t *testing.T) {
	t.Parallel()
	probe := New(nil)
	probe.root = t.TempDir()
	uid := strconv.Itoa(os.Getuid())
	process(t, probe.root, "1", "sway", uid)
	process(t, probe.root, "2", "sway", uid)
	process(t, probe.root, "3", "other", "different")
	process(t, probe.root, "4", "", uid)
	process(t, probe.root, "not-pid", "ignored", uid)
	names, err := probe.Processes(t.Context())
	if err != nil || len(names) != 1 || names[0] != "sway" {
		t.Fatal(names, err)
	}
	probe.filter = func(string) bool { return false }
	if names, err = probe.Processes(t.Context()); err != nil || len(names) != 0 {
		t.Fatal(names, err)
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if _, err = probe.Processes(ctx); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	probe.root = filepath.Join(probe.root, "missing")
	if _, err = probe.Processes(t.Context()); !errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	}
}

func TestWalkLimitAndParsing(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	scan := newScanner(func(string) bool { return true }, root)
	entries := []string{"bad", "99"}
	for i := range maxMatches + 1 {
		pid := strconv.Itoa(i)
		process(t, root, pid, "name"+pid, scan.uid)
		entries = append(entries, pid)
	}
	names, err := walk(t.Context(), scan, entries)
	if err != nil || len(names) != maxMatches {
		t.Fatal(names, err)
	}
	if scan.ownedBy("99", scan.uid) || isPID("") || isPID("12x") || !isPID("123") {
		t.Fatal("invalid PID/ownership")
	}
	if realUID([]byte("Name: only\n")) != "" || firstField(" \t") != "" {
		t.Fatal("missing UID")
	}
}
