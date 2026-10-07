// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package rules_test

import (
	"fmt"
	"testing"

	"github.com/gostafa/linuxdesktop/internal/domain"
	"github.com/gostafa/linuxdesktop/internal/rules"
)

func benchmarkGlobals(count int) []domain.WaylandGlobal {
	globals := make([]domain.WaylandGlobal, 0, count)
	for index := range count {
		globals = append(globals, domain.WaylandGlobal{
			Interface: fmt.Sprintf("zwp_fixture_v%d", index),
			Name:      uint32(index + 1),
			Version:   1,
		})
	}

	return globals
}

func BenchmarkCompositor(b *testing.B) {
	processes := make([]string, 0, 200)
	for index := range 199 {
		processes = append(processes, fmt.Sprintf("process-%d", index))
	}
	processes = append(processes, "sway")
	for _, bench := range []struct {
		name string
		sig  domain.Signals
	}{
		{"env", domain.Signals{Env: domain.Env{HyprlandSignature: "x"}}},
		{"wayland-globals", domain.Signals{
			WaylandGlobals: append(benchmarkGlobals(64), domain.WaylandGlobal{Interface: "gtk_shell1"}),
		}},
		{"x11-window-manager", domain.Signals{X11WindowManager: "KWin"}},
		{"desktop", domain.Signals{Desktop: domain.DesktopGNOME}},
		{"processes", domain.Signals{Processes: processes}},
		{"unknown", domain.Signals{WaylandGlobals: benchmarkGlobals(64)}},
	} {
		b.Run(bench.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = rules.Compositor(&bench.sig)
			}
		})
	}
}

func BenchmarkDesktop(b *testing.B) {
	for _, bench := range []struct {
		name string
		env  domain.Env
	}{
		{"gnome", domain.Env{CurrentDesktop: "ubuntu:GNOME"}},
		{"kde", domain.Env{CurrentDesktop: "KDE", SessionDesktop: "plasmawayland"}},
		{"fallback", domain.Env{CurrentDesktop: "unrecognized:custom", KDEFullSession: "true"}},
		{"empty", domain.Env{}},
	} {
		b.Run(bench.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = rules.Desktop(&bench.env)
			}
		})
	}
}

func BenchmarkIsCompositorProcess(b *testing.B) {
	names := []string{"sway", "gnome-shell", "kwin_wayland", "sleep", "firefox", "Hyprland"}
	b.ReportAllocs()
	for b.Loop() {
		for _, name := range names {
			_ = rules.IsCompositorProcess(name)
		}
	}
}
