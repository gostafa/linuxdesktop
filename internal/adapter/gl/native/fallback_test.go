//go:build !(linux && (amd64 || arm64 || 386 || arm || riscv64 || loong64 || ppc64le))

// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package native

import (
	"testing"
)

func TestUnsupportedPlatform(t *testing.T) {
	t.Parallel()
	if version, ok := Vulkan(); ok || version != "" {
		t.Fatal("unexpected Vulkan support")
	}
	if info, ok := OpenGL(); ok || info.Available {
		t.Fatal("unexpected OpenGL support")
	}
}
