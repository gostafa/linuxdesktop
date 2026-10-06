// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package env

import (
	"testing"
)

func TestSnapshot(t *testing.T) {
	t.Setenv(keyUser, "alice")
	t.Setenv(keyLogname, "bob")
	t.Setenv(keySessionID, "7")
	t.Setenv(keyWaylandDisplay, "wayland-1")
	t.Setenv(keySSHConnection, "remote")
	got := New().Snapshot()
	if got.User != "alice" || got.SessionID != "7" || got.WaylandDisplay != "wayland-1" ||
		got.SSHConnection != "remote" {
		t.Fatalf("snapshot: %+v", got)
	}
	if firstNonEmpty("", "bob") != "bob" || firstNonEmpty("", "") != "" {
		t.Fatal("fallback selection failed")
	}
}
