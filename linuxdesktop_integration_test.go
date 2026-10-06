// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package linuxdesktop_test

import (
	"bufio"
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

	"github.com/gostafa/linuxdesktop"
)

func TestPublicAPIUsage(t *testing.T) {
	var session linuxdesktop.SessionType = linuxdesktop.SessionTypeWayland
	var protocol linuxdesktop.DisplayProtocol = linuxdesktop.DisplayProtocolWayland
	var sections linuxdesktop.Section = linuxdesktop.SectionOS | linuxdesktop.SectionDisplay
	result := linuxdesktop.Environment{
		Session: linuxdesktop.SessionInfo{Type: session},
		Display: linuxdesktop.DisplayInfo{Protocol: protocol, Wayland: &linuxdesktop.WaylandInfo{
			Globals: []linuxdesktop.WaylandGlobal{{Interface: "wl_compositor"}},
		}},
	}
	if reflect.TypeOf(result).PkgPath() != "github.com/gostafa/linuxdesktop" {
		t.Fatal("public result has the wrong package identity")
	}
	opts := []linuxdesktop.Option{
		linuxdesktop.WithTimeout(time.Second),
		linuxdesktop.WithProbeTimeout(time.Millisecond),
		linuxdesktop.WithSections(
			sections,
		),
		linuxdesktop.WithOpenGL(),
		linuxdesktop.WithProcessScan(),
		func(c *linuxdesktop.Config) { c.Sections = linuxdesktop.SectionOS },
	}
	var cfg linuxdesktop.Config
	for _, opt := range opts {
		opt(&cfg)
	}
	if cfg.Timeout != time.Second || cfg.ProbeTimeout != time.Millisecond ||
		cfg.Sections != linuxdesktop.SectionOS ||
		!cfg.NativeGL ||
		!cfg.ProcessScan {
		t.Fatalf("public options did not configure Config: %+v", cfg)
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

func TestNonLinuxDetection(t *testing.T) {
	if runtime.GOOS == "linux" {
		t.Skip("non-Linux behavior")
	}
	called := false
	out, err := linuxdesktop.DetectContext(nil, nil, func(c *linuxdesktop.Config) { called = true })
	if !called || !errors.Is(err, linuxdesktop.ErrNotLinux) || out == nil || !out.Headless ||
		out.Session.Type != linuxdesktop.SessionTypeUnknown || out.Display.Protocol != linuxdesktop.DisplayProtocolUnknown ||
		out.Desktop.Environment != linuxdesktop.DesktopUnknown {
		t.Fatalf("unexpected non-Linux result: %+v, %v (option called: %v)", out, err, called)
	}
	if linuxdesktop.IsWayland() || linuxdesktop.IsX11() || !linuxdesktop.IsHeadless() {
		t.Fatal("unexpected non-Linux display availability")
	}
	if _, err = linuxdesktop.OS(); !errors.Is(err, linuxdesktop.ErrNotLinux) {
		t.Fatalf("section helper lost error: %v", err)
	}
}

func TestPublicNilContextAndOS(t *testing.T) {
	t.Parallel()
	result, err := linuxdesktop.DetectContext(
		nil,
		linuxdesktop.WithSections(linuxdesktop.SectionOS),
	)
	if result == nil || (runtime.GOOS != "linux" && !errors.Is(err, linuxdesktop.ErrNotLinux)) {
		t.Fatal(result, err)
	}
	osInfo, _ := linuxdesktop.OS()
	if runtime.GOOS == "linux" && osInfo.Kernel == "" {
		t.Fatal("operating system result lost")
	}
}
