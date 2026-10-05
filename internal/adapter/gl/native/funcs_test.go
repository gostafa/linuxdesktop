//go:build linux && (amd64 || arm64 || 386 || arm || riscv64 || loong64 || ppc64le)

// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package native

import (
	"reflect"
	"testing"
)

func TestDecodeVersion(t *testing.T) {
	t.Parallel()
	packed := uint32(1<<vkMajorShift | 3<<vkMinorShift | 283)
	if got := decodeVersion(packed); got != "1.3.283" {
		t.Fatalf("version = %q", got)
	}
	if got := decodeVersion(packed | 2<<vkVariantShift); got != "1.3.283 (variant 2)" {
		t.Fatalf("variant version = %q", got)
	}
}

func TestChooseConfig(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name   string
		status uint32
		count  int32
		want   bool
	}{
		{"success", 1, 1, true},
		{"no configs", 1, 0, false},
		{"failed", 0, 1, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			egl := eglFuncs{
				chooseConfig: func(display uintptr, attrs *int32, config *uintptr, size int32, count *int32) uint32 {
					if display != 7 || *attrs != eglSurfaceType || size != 1 {
						t.Fatal("invalid configuration request")
					}
					*config, *count = 9, test.count
					return test.status
				},
			}
			config, ok := chooseConfig(&egl, 7, eglOpenGLBit)
			if config != 9 || ok != test.want {
				t.Fatalf("config = %d, ok = %v", config, ok)
			}
		})
	}
}

func TestCurrentGLTriesBothAPIsAndReleasesContexts(t *testing.T) {
	t.Parallel()
	var apis []uint32
	var destroyed []uintptr
	egl := eglFuncs{
		bindAPI: func(api uint32) uint32 { apis = append(apis, api); return 1 },
		chooseConfig: func(_ uintptr, _ *int32, config *uintptr, _ int32, count *int32) uint32 {
			*config, *count = 9, 1
			return 1
		},
		createContext: func(_ uintptr, _ uintptr, _ uintptr, attrs *int32) uintptr {
			if len(apis) == 1 && *attrs != eglNone {
				t.Fatal("desktop context attributes changed")
			}
			if len(apis) == 2 && *attrs != eglContextVersion {
				t.Fatal("ES context attributes changed")
			}
			return uintptr(len(apis))
		},
		makeCurrent:    func(_, _, _, _ uintptr) uint32 { return 0 },
		destroyContext: func(_ uintptr, ctx uintptr) uint32 { destroyed = append(destroyed, ctx); return 1 },
	}
	if info, ok := currentGL(&egl, 7); ok || info.Available {
		t.Fatalf("failed contexts returned a result: %+v", info)
	}
	if !reflect.DeepEqual(apis, []uint32{eglOpenGLAPI, eglOpenGLESAPI}) {
		t.Fatalf("API order = %v", apis)
	}
	if !reflect.DeepEqual(destroyed, []uintptr{1, 2}) {
		t.Fatalf("destroyed contexts = %v", destroyed)
	}
}

func TestInitializeDisplaySkipsAbsentDisplay(t *testing.T) {
	t.Parallel()
	var calls []string
	egl := eglFuncs{
		getDisplay: func(uintptr) uintptr { return 7 },
		initialize: func(_ uintptr, _, _ *int32) uint32 { return 1 },
	}
	display, ok := initializeDisplay(&egl)
	if !ok || display != 7 {
		t.Fatalf("display = %d, ok = %v", display, ok)
	}
	// No display must never reach initialization.
	egl.getDisplay = func(uintptr) uintptr { return eglNoDisplay }
	egl.initialize = func(_ uintptr, _, _ *int32) uint32 { calls = append(calls, "initialize"); return 1 }
	if _, ok := initializeDisplay(&egl); ok || len(calls) != 0 {
		t.Fatal("initialized an absent display")
	}
}
