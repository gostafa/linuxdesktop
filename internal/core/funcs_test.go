// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package core

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"

	"github.com/gostafa/linuxdesktop/internal/domain"
)

type fixture struct {
	lock  sync.Mutex
	calls map[string]int
	err   error
}

func (f *fixture) called(name string) {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.calls[name]++
}

func (f *fixture) Snapshot() domain.Env {
	f.called("env")
	return domain.Env{Display: ":1", WaylandDisplay: "wayland-1"}
}

func (f *fixture) OS(context.Context) (domain.OSInfo, error) {
	f.called("os")
	return domain.OSInfo{ID: "fixture-os"}, f.err
}

func (f *fixture) Session(context.Context, *domain.Env) (domain.SessionInfo, error) {
	f.called("session")
	return domain.SessionInfo{ID: "fixture-session", Type: domain.SessionTypeWayland}, f.err
}

func (f *fixture) X11(context.Context, *domain.Env) (*domain.X11Info, error) {
	f.called("x11")
	return &domain.X11Info{WindowManager: "fixture-wm"}, f.err
}

func (f *fixture) Wayland(context.Context, *domain.Env) (*domain.WaylandInfo, error) {
	f.called("wayland")
	return &domain.WaylandInfo{
		Globals: []domain.WaylandGlobal{{Interface: "fixture-interface"}},
	}, f.err
}

func (f *fixture) Desktop(context.Context, *domain.Env) (domain.DesktopInfo, error) {
	f.called("desktop")
	return domain.DesktopInfo{Name: "fixture-desktop"}, f.err
}

func (f *fixture) GPUs(context.Context) ([]domain.GPUInfo, string, error) {
	f.called("gpus")
	return []domain.GPUInfo{{ID: "fixture-gpu"}}, "fixture-gpu", f.err
}

func (f *fixture) Stack(context.Context) (domain.OpenGLInfo, domain.VulkanInfo, error) {
	f.called("stack")
	return domain.OpenGLInfo{Vendor: "fixture-gl"}, domain.VulkanInfo{Version: "1.3"}, f.err
}

func (f *fixture) Portal(context.Context, *domain.Env) (domain.PortalInfo, error) {
	f.called("portal")
	return domain.PortalInfo{Backend: "fixture-portal"}, f.err
}

func (f *fixture) Processes(context.Context) ([]string, error) {
	f.called("process")
	return []string{"fixture-process"}, f.err
}

func fixtureDeps(f *fixture) Deps {
	return Deps{
		Env: f, OS: f, Session: f, X11: f, Wayland: f,
		Desktop: f, GPUs: f, Stack: f, Portal: f, Process: f,
	}
}

func TestProbeBindingsPreservePartialResults(t *testing.T) {
	t.Parallel()
	failure := errors.New("fixture failure")
	f := &fixture{calls: make(map[string]int), err: failure}
	deps := fixtureDeps(f)
	out, err := New(
		&deps,
		&Config{Sections: domain.SectionAll, ProcessScan: true},
	).Detect(t.Context())
	if !errors.Is(err, failure) {
		t.Fatalf("probe error lost: %v", err)
	}
	if out.OS.ID != "fixture-os" || out.Session.ID != "fixture-session" ||
		out.Desktop.Name != "fixture-desktop" || out.Portal.Backend != "fixture-portal" ||
		out.Graphics.PrimaryGPU != "fixture-gpu" || len(out.Graphics.GPUs) != 1 ||
		out.Graphics.OpenGL.Vendor != "fixture-gl" || out.Graphics.Vulkan.Version != "1.3" ||
		out.Display.X11 == nil || out.Display.X11.WindowManager != "fixture-wm" ||
		out.Display.Wayland == nil || len(out.Display.Wayland.Globals) != 1 {
		t.Fatalf("partial results mapped incorrectly: %+v", out)
	}
	want := map[string]int{
		"env": 1, "os": 1, "session": 1, "x11": 1, "wayland": 1,
		"desktop": 1, "gpus": 1, "stack": 1, "portal": 1, "process": 1,
	}
	if !reflect.DeepEqual(f.calls, want) {
		t.Fatalf("probe calls: got %v, want %v", f.calls, want)
	}
}

func TestProbeBindingsRespectSections(t *testing.T) {
	t.Parallel()
	f := &fixture{calls: make(map[string]int)}
	deps := fixtureDeps(f)
	out, err := New(&deps, &Config{Sections: domain.SectionSession}).Detect(t.Context())
	if err != nil || out.Session.ID != "fixture-session" {
		t.Fatalf("session result: %+v, %v", out, err)
	}
	want := map[string]int{"env": 1, "session": 1}
	if !reflect.DeepEqual(f.calls, want) {
		t.Fatalf("unselected probe ran: got %v, want %v", f.calls, want)
	}
}

func TestProbeBindingsAllowUnwiredAdapters(t *testing.T) {
	t.Parallel()
	deps := Deps{}
	out, err := New(
		&deps,
		&Config{Sections: domain.SectionAll, ProcessScan: true},
	).Detect(t.Context())
	if out == nil || err != nil || !out.Headless {
		t.Fatalf("unwired run: %+v, %v", out, err)
	}
}
