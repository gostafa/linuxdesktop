// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package rules

import (
	"testing"

	"github.com/gostafa/linuxdesktop/internal/domain"
)

func TestDesktopFallbacks(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		env  domain.Env
		want domain.DesktopEnvironment
	}{
		{domain.Env{}, domain.DesktopUnknown},
		{domain.Env{CurrentDesktop: "unrecognized:GNOME"}, domain.DesktopGNOME},
		{domain.Env{CurrentDesktop: "custom", KDEFullSession: "true"}, domain.DesktopKDE},
		{domain.Env{KDEFullSession: "true"}, domain.DesktopKDE},
		{domain.Env{GNOMESessionID: "old"}, domain.DesktopGNOME},
	} {
		if got := Desktop(&tc.env); got.Environment != tc.want {
			t.Fatalf("%+v: %+v", tc.env, got)
		}
	}
}

func TestCompositorLadder(t *testing.T) {
	wlrootsMarkers := wlrootsMarkers()

	t.Parallel()
	for _, tc := range []struct {
		sig    domain.Signals
		name   string
		method domain.DetectionMethod
	}{
		{domain.Signals{}, "", domain.DetectedUnknown},
		{domain.Signals{Env: domain.Env{HyprlandSignature: "x"}}, nameHyprland, domain.DetectedEnvironment},
		{domain.Signals{Env: domain.Env{SwaySock: "x"}}, nameSway, domain.DetectedEnvironment},
		{domain.Signals{Env: domain.Env{WayfireSocket: "x"}}, nameWayfire, domain.DetectedEnvironment},
		{domain.Signals{Env: domain.Env{I3Sock: "x"}}, nameI3, domain.DetectedEnvironment},
		{domain.Signals{WaylandGlobals: []domain.WaylandGlobal{{Interface: "gtk_shell1"}}}, nameMutter, domain.DetectedWayland},
		{domain.Signals{X11WindowManager: "KWin"}, nameKWin, domain.DetectedX11EWMH},
		{domain.Signals{X11WindowManager: "custom 2"}, "custom 2", domain.DetectedX11EWMH},
		{domain.Signals{Desktop: domain.DesktopGNOME}, nameMutter, domain.DetectedEnvironment},
		{domain.Signals{WaylandGlobals: []domain.WaylandGlobal{{Interface: wlrootsMarkers[0]}}}, nameWlroots, domain.DetectedWayland},
		{domain.Signals{Processes: []string{"irrelevant", "sway"}}, nameSway, domain.DetectedProcess},
	} {
		got := Compositor(&tc.sig)
		if got.Name != tc.name || got.DetectedBy != tc.method {
			t.Fatalf("%+v: %+v", tc.sig, got)
		}
	}
	if !IsCompositorProcess("sway") || IsCompositorProcess("sleep") {
		t.Fatal("process filter")
	}
	if normalize("Test_123!") != "test123" {
		t.Fatal("normalization")
	}
}

func TestProtocolFallbacks(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		kind  domain.SessionType
		wl, x bool
		want  domain.DisplayProtocol
	}{
		{domain.SessionTypeUnknown, true, true, domain.DisplayProtocolWayland},
		{domain.SessionTypeUnknown, false, true, domain.DisplayProtocolX11},
		{domain.SessionTypeWayland, false, false, domain.DisplayProtocolWayland},
		{domain.SessionTypeX11, false, false, domain.DisplayProtocolX11},
		{domain.SessionTypeTTY, false, false, domain.DisplayProtocolUnknown},
	} {
		if got := Protocol(tc.kind, tc.wl, tc.x); got != tc.want {
			t.Fatal(got)
		}
	}
	if !Headless(false, false) || Headless(true, false) {
		t.Fatal("headless")
	}
}
