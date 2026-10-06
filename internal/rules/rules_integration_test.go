// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package rules_test

import (
	"testing"

	"github.com/gostafa/linuxdesktop/internal/domain"
	"github.com/gostafa/linuxdesktop/internal/rules"
)

func TestDesktopFallbacks(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		env  domain.Env
		want domain.DesktopEnvironment
	}{
		{domain.Env{}, domain.Unknown},
		{domain.Env{CurrentDesktop: "unrecognized:GNOME"}, domain.DesktopGNOME},
		{domain.Env{CurrentDesktop: "custom", KDEFullSession: "true"}, domain.DesktopKDE},
		{domain.Env{KDEFullSession: "true"}, domain.DesktopKDE},
		{domain.Env{GNOMESessionID: "old"}, domain.DesktopGNOME},
	} {
		if got := rules.Desktop(&tc.env); got.Environment != tc.want {
			t.Fatalf("%+v: %+v", tc.env, got)
		}
	}
}

func TestProtocolFallbacks(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		kind  domain.SessionType
		wl, x bool
		want  domain.DisplayProtocol
	}{
		{domain.Unknown, true, true, domain.DisplayProtocol(domain.SessionTypeWayland)},
		{domain.Unknown, false, true, domain.DisplayProtocol(domain.SessionTypeX11)},
		{domain.SessionTypeWayland, false, false, domain.DisplayProtocol(domain.SessionTypeWayland)},
		{domain.SessionTypeX11, false, false, domain.DisplayProtocol(domain.SessionTypeX11)},
		{domain.SessionTypeTTY, false, false, domain.Unknown},
	} {
		if got := rules.Protocol(tc.kind, tc.wl, tc.x); got != tc.want {
			t.Fatal(got)
		}
	}
	if !rules.Headless(false, false) || rules.Headless(true, false) {
		t.Fatal("headless")
	}
}

// Desktops without a default compositor must still reach the weaker rungs.
func TestCompositorWithoutDesktopDefault(t *testing.T) {
	t.Parallel()

	for _, desktop := range []domain.DesktopEnvironment{
		domain.Unknown,
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
