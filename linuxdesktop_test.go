// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package linuxdesktop

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"go/importer"
	"go/token"
	"go/types"
	"io"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/gostafa/linuxdesktop/internal/core"
	"github.com/gostafa/linuxdesktop/internal/domain"
)

func TestPublicAPIUsage(t *testing.T) {
	var session SessionType = SessionTypeWayland
	var protocol DisplayProtocol = DisplayProtocolWayland
	var sections Section = SectionOS | SectionDisplay
	result := Environment{
		Session: SessionInfo{Type: session},
		Display: DisplayInfo{Protocol: protocol, Wayland: &WaylandInfo{
			Globals: []WaylandGlobal{{Interface: "wl_compositor"}},
		}},
	}
	if reflect.TypeOf(result).PkgPath() != "github.com/gostafa/linuxdesktop" {
		t.Fatal("public result has the wrong package identity")
	}
	opts := []Option{
		WithTimeout(time.Second), WithProbeTimeout(time.Millisecond),
		WithSections(sections), WithOpenGL(), WithProcessScan(),
		func(c *Config) { c.Sections = SectionOS },
	}
	var cfg Config
	for _, opt := range opts {
		opt(&cfg)
	}
	if cfg.Timeout != time.Second || cfg.ProbeTimeout != time.Millisecond ||
		cfg.Sections != SectionOS ||
		!cfg.NativeGL ||
		!cfg.ProcessScan {
		t.Fatalf("public options did not configure Config: %+v", cfg)
	}
}

func TestNonLinuxDetection(t *testing.T) {
	if runtime.GOOS == "linux" {
		t.Skip("non-Linux behavior")
	}
	called := false
	out, err := DetectContext(nil, nil, func(c *Config) { called = true })
	if !called || !errors.Is(err, ErrNotLinux) || out == nil || !out.Headless ||
		out.Session.Type != SessionTypeUnknown || out.Display.Protocol != DisplayProtocolUnknown ||
		out.Desktop.Environment != DesktopUnknown {
		t.Fatalf("unexpected non-Linux result: %+v, %v (option called: %v)", out, err, called)
	}
	if IsWayland() || IsX11() || !IsHeadless() {
		t.Fatal("unexpected non-Linux display availability")
	}
	if _, err = OS(); !errors.Is(err, ErrNotLinux) {
		t.Fatalf("section helper lost error: %v", err)
	}
}

// Use compiled export data to inspect every exported declaration, including
// constants and types that a manually maintained list could miss.
func TestPublicAPIBoundary(t *testing.T) {
	const module = "github.com/gostafa/linuxdesktop"
	cmd := exec.Command("go", "list", "-export", "-deps", "-f", "{{.ImportPath}} {{.Export}}", ".")
	data, err := cmd.Output()
	if err != nil {
		t.Fatalf("load export data: %v", err)
	}
	files := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		fields := strings.SplitN(scanner.Text(), " ", 2)
		if len(fields) == 2 && fields[1] != "" {
			files[fields[0]] = fields[1]
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	loader := importer.ForCompiler(
		token.NewFileSet(),
		"gc",
		func(path string) (io.ReadCloser, error) {
			return os.Open(files[path])
		},
	)
	pkg, err := loader.Import(module)
	if err != nil {
		t.Fatal(err)
	}
	seen := make(map[types.Type]bool)
	var visit func(types.Type)
	visit = func(typ types.Type) {
		if typ == nil || seen[typ] {
			return
		}
		seen[typ] = true
		switch typ := typ.(type) {
		case *types.Alias:
			if p := typ.Obj().Pkg(); p != nil && strings.Contains(p.Path(), "/internal/") {
				t.Errorf("public API exposes internal alias %s", typ)
			}
			visit(types.Unalias(typ))
		case *types.Named:
			if p := typ.Obj().Pkg(); p != nil && strings.Contains(p.Path(), "/internal/") {
				t.Errorf("public API exposes internal type %s", typ)
			}
			for i := range typ.NumMethods() {
				if typ.Method(i).Exported() {
					visit(typ.Method(i).Type())
				}
			}
			visit(typ.Underlying())
		case *types.Pointer:
			visit(typ.Elem())
		case *types.Slice:
			visit(typ.Elem())
		case *types.Array:
			visit(typ.Elem())
		case *types.Map:
			visit(typ.Key())
			visit(typ.Elem())
		case *types.Chan:
			visit(typ.Elem())
		case *types.Struct:
			for i := range typ.NumFields() {
				if typ.Field(i).Exported() {
					visit(typ.Field(i).Type())
				}
			}
		case *types.Signature:
			visit(typ.Params())
			visit(typ.Results())
		case *types.Tuple:
			for i := range typ.Len() {
				visit(typ.At(i).Type())
			}
		case *types.Interface:
			for i := range typ.NumMethods() {
				visit(typ.Method(i).Type())
			}
		}
	}
	for _, name := range pkg.Scope().Names() {
		obj := pkg.Scope().Lookup(name)
		if obj.Exported() {
			visit(obj.Type())
		}
	}

	// The duplicated enum values and defaults must agree across the boundary.
	for _, path := range []string{module + "/internal/domain", module + "/internal/core"} {
		internal, err := loader.Import(path)
		if err != nil {
			t.Fatal(err)
		}
		for _, name := range pkg.Scope().Names() {
			publicConst, ok := pkg.Scope().Lookup(name).(*types.Const)
			if !ok || !publicConst.Exported() {
				continue
			}
			if internalConst, ok := internal.Scope().Lookup(name).(*types.Const); ok {
				if publicConst.Val().ExactString() != internalConst.Val().ExactString() {
					t.Errorf("%s differs from %s", name, path)
				}
			}
		}
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
