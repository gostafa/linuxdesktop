// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package rules

import (
	"testing"

	"github.com/gostafa/linuxdesktop/internal/domain"
)

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
