// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package rules_test

import (
	"testing"

	"github.com/gostafa/linuxdesktop/internal/domain"
	"github.com/gostafa/linuxdesktop/internal/rules"
)

// Desktops without a default compositor must still reach the weaker rungs.
func TestCompositorWithoutDesktopDefault(t *testing.T) {
	t.Parallel()

	for _, desktop := range []domain.DesktopEnvironment{
		domain.DesktopUnknown,
		domain.DesktopLXQt,
		domain.DesktopEnvironment("custom"),
	} {
		t.Run(string(desktop), func(t *testing.T) {
			t.Parallel()

			var signals domain.Signals

			signals.Desktop = desktop
			signals.WaylandGlobals = []domain.WaylandGlobal{
				{Interface: "zwlr_layer_shell_v1", Name: 1, Version: 1},
			}

			info := rules.Compositor(&signals)
			if info.Name != "wlroots" || info.DetectedBy != domain.DetectedWayland {
				t.Fatalf("desktop %q stopped compositor fallback: %+v", desktop, info)
			}
		})
	}
}
