// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package dbusconn_test

import (
	"testing"

	"github.com/gostafa/linuxdesktop/internal/adapter/dbusconn"
	"github.com/gostafa/linuxdesktop/internal/port"
)

func TestHostConnectionFailures(t *testing.T) {
	t.Setenv("DBUS_SYSTEM_BUS_ADDRESS", "unix:path=/nonexistent-linuxdesktop-test-bus")
	t.Setenv("DBUS_SESSION_BUS_ADDRESS", "unix:path=/nonexistent-linuxdesktop-test-bus")
	bus := dbusconn.New(nil)
	for _, kind := range []port.BusKind{port.SessionBus, port.SystemBus} {
		if _, err := bus.HasOwner(t.Context(), kind, "fixture"); err == nil {
			t.Fatal("connected to missing bus")
		}
	}
	if err := bus.Close(); err != nil {
		t.Fatal(err)
	}
}
