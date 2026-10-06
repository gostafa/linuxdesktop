// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package desktop

import (
	"context"
	"errors"
	"testing"

	"github.com/gostafa/linuxdesktop/internal/domain"
	"github.com/gostafa/linuxdesktop/internal/port"
)

func TestDesktopVersions(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		desktop, version string
		bus              port.Bus
		want             string
	}{
		{"KDE", "6", nil, "6"},
		{"other", "6", nil, ""},
		{"GNOME", "", nil, ""},
		{"GNOME", "", fakeBus{value: "48"}, "48"},
		{"GNOME", "", fakeBus{value: 48}, ""},
		{"GNOME", "", fakeBus{err: errors.New("offline")}, ""},
	} {
		info, err := New(
			tc.bus,
		).Desktop(t.Context(), &domain.Env{CurrentDesktop: tc.desktop, KDESessionVersion: tc.version})
		if err != nil || info.Version != tc.want {
			t.Fatalf("%s: %+v, %v", tc.desktop, info, err)
		}
	}
}

type fakeBus struct {
	port.Bus
	value any
	err   error
}

func (bus fakeBus) Property(
	_ context.Context,
	kind port.BusKind,
	query *port.PropertyQuery,
) (port.PropertyValue, error) {
	if kind != port.SessionBus || query.Name != shellVersion {
		return port.PropertyValue{}, errors.New("invalid query")
	}
	return port.PropertyValue{Value: bus.value}, bus.err
}
