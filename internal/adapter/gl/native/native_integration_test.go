//go:build !(linux && (amd64 || arm64 || 386 || arm || riscv64 || loong64 || ppc64le))

// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package native_test

import (
	"github.com/gostafa/linuxdesktop/internal/adapter/gl/native"
	"testing"
)

func TestUnsupportedPlatform(t *testing.T) {
	t.Parallel()
	if version, ok := native.Vulkan(); ok || version != "" {
		t.Fatal("unexpected Vulkan support")
	}
	if info, ok := native.OpenGL(); ok || info.Available {
		t.Fatal("unexpected OpenGL support")
	}
}
