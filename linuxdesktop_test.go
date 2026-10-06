// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package linuxdesktop

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/gostafa/linuxdesktop/internal/core"
	"github.com/gostafa/linuxdesktop/internal/domain"
)

func TestPublicSections(t *testing.T) {
	t.Parallel()
	if result, err := Detect(); result == nil ||
		(runtime.GOOS != goosLinux && !errors.Is(err, ErrNotLinux)) {
		t.Fatal(result, err)
	}
	for _, call := range []func() error{
		func() error { _, err := OS(); return err }, func() error { _, err := Session(); return err },
		func() error { _, err := Display(); return err }, func() error { _, err := Desktop(); return err },
		func() error { _, err := Compositor(); return err }, func() error { _, err := Graphics(); return err },
		func() error { _, err := Portal(); return err },
	} {
		if err := call(); runtime.GOOS != goosLinux && !errors.Is(err, ErrNotLinux) {
			t.Fatal("section helper lost platform error", err)
		}
	}
	cfg := newConfig(WithSections(SectionOS|SectionGraphics), WithProcessScan(), WithOpenGL())
	result, err := detectOn(t.Context(), cfg, goosLinux)
	if result == nil || err != nil {
		t.Fatal(result, err)
	}
	if _, err = detectOn(t.Context(), cfg, "other"); !errors.Is(err, ErrNotLinux) {
		t.Fatal(err)
	}
	if onLinux(
		"other",
		func() bool { t.Fatal("probe ran on unsupported platform"); return true },
	) ||
		!onLinux(goosLinux, func() bool { return true }) {
		t.Fatal("platform guard")
	}
	_ = waylandAvailable()
	_ = x11Available()
	_ = IsWayland()
	_ = IsX11()
	_ = IsHeadless()
	if graphicsProbe(&core.Config{}) == nil {
		t.Fatal("default graphics probe missing")
	}
}

func TestConfig(t *testing.T) {
	defaults := core.NewConfig()
	if got := newConfig(nil); *got != defaults {
		t.Fatalf("defaults: got %+v, want %+v", got, defaults)
	}
	for _, duration := range []time.Duration{0, -time.Second, time.Second} {
		t.Run(duration.String(), func(t *testing.T) {
			got := newConfig(WithTimeout(duration), WithProbeTimeout(duration), WithSections(0))
			want := core.NewConfig(
				core.WithTimeout(duration),
				core.WithProbeTimeout(duration),
				core.WithSections(0),
			)
			if *got != want {
				t.Fatalf("got %+v, want %+v", got, want)
			}
		})
	}
	got := newConfig(
		WithTimeout(time.Second), nil,
		func(c *Config) {
			if c.Timeout != time.Second || c.ProbeTimeout != DefaultProbeTimeout ||
				c.Sections != SectionAll {
				t.Fatalf("custom option received unexpected configuration: %+v", c)
			}
			c.Timeout = 3 * time.Second
		},
		WithTimeout(4*time.Second), WithProbeTimeout(time.Second),
		WithSections(SectionOS|SectionGraphics), WithOpenGL(), WithProcessScan(),
	)
	want := core.Config{
		Timeout: 4 * time.Second, ProbeTimeout: time.Second,
		Sections: domain.SectionOS | domain.SectionGraphics, NativeGL: true, ProcessScan: true,
	}
	if *got != want {
		t.Fatalf("ordered options: got %+v, want %+v", got, want)
	}
}

func TestDetectionError(t *testing.T) {
	t.Parallel()
	if err := detectionError(nil, "context"); err != nil {
		t.Fatalf("successful detection acquired an error: %v", err)
	}
	failure := errors.New("probe failed")
	err := detectionError(failure, "context")
	if !errors.Is(err, failure) || err.Error() != "linuxdesktop: detect context: probe failed" {
		t.Fatalf("probe error identity or operation lost: %v", err)
	}
}

// Populate every field so an omitted conversion cannot hide behind a zero value.
func populate(v reflect.Value) {
	switch v.Kind() {
	case reflect.Struct:
		for i := range v.NumField() {
			populate(v.Field(i))
		}
	case reflect.Pointer:
		v.Set(reflect.New(v.Type().Elem()))
		populate(v.Elem())
	case reflect.Slice:
		v.Set(reflect.MakeSlice(v.Type(), 2, 2))
		for i := range v.Len() {
			populate(v.Index(i))
		}
	case reflect.String:
		v.SetString(v.Type().String() + "-value")
	case reflect.Bool:
		v.SetBool(true)
	case reflect.Int, reflect.Int32, reflect.Int64:
		v.SetInt(7)
	case reflect.Uint16, reflect.Uint32:
		v.SetUint(11)
	default:
		panic("unhandled fixture type: " + v.Type().String())
	}
}

func assertJSONParity(t *testing.T, in *domain.Environment) *Environment {
	t.Helper()
	out := publicEnvironment(in)
	want, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	got, err := json.Marshal(out)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("conversion changed JSON:\ngot  %s\nwant %s", got, want)
	}
	return out
}

func TestPartialResultConversion(t *testing.T) {
	failure := errors.New("probe failed")
	deps := core.Deps{OS: failingOS{err: failure}}
	in, err := core.New(&deps, newConfig(WithSections(SectionOS))).Detect(context.Background())
	out := assertJSONParity(t, in)
	if out.OS.ID != "partial" || !errors.Is(err, failure) {
		t.Fatalf("partial data or error lost: %+v, %v", out, err)
	}
}

type failingOS struct{ err error }

func (p failingOS) OS(context.Context) (domain.OSInfo, error) {
	return domain.OSInfo{ID: "partial"}, p.err
}

func TestEnvironmentConversion(t *testing.T) {
	in := &domain.Environment{}
	populate(reflect.ValueOf(in).Elem())
	out := assertJSONParity(t, in)
	out.Display.X11.Extensions[0] = "changed"
	out.Display.Wayland.Globals[0].Interface = "changed"
	out.Desktop.CurrentDesktops[0] = "changed"
	out.Graphics.GPUs[0].ID = "changed"
	if in.Display.X11.Extensions[0] == "changed" ||
		in.Display.Wayland.Globals[0].Interface == "changed" ||
		in.Desktop.CurrentDesktops[0] == "changed" ||
		in.Graphics.GPUs[0].ID == "changed" {
		t.Fatal("conversion shared mutable data")
	}
}

func TestEnvironmentConversionOptionalValues(t *testing.T) {
	out := assertJSONParity(t, &domain.Environment{})
	if out.Display.X11 != nil || out.Display.Wayland != nil || out.Graphics.GPUs != nil ||
		out.Desktop.CurrentDesktops != nil {
		t.Fatal("nil values were not preserved")
	}
	for _, empty := range []bool{false, true} {
		in := &domain.Environment{
			Display: domain.DisplayInfo{X11: &domain.X11Info{}, Wayland: &domain.WaylandInfo{}},
		}
		if empty {
			in.Display.X11.Extensions = []string{}
			in.Display.Wayland.Globals = []domain.WaylandGlobal{}
			in.Desktop.CurrentDesktops = []string{}
			in.Graphics.GPUs = []domain.GPUInfo{}
		}
		out = assertJSONParity(t, in)
		if (out.Display.X11.Extensions != nil) != empty ||
			(out.Display.Wayland.Globals != nil) != empty ||
			(out.Desktop.CurrentDesktops != nil) != empty ||
			(out.Graphics.GPUs != nil) != empty {
			t.Fatal("nil versus empty slices were not preserved")
		}
	}
}
