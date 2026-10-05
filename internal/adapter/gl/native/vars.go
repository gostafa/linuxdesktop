//go:build !(linux && (amd64 || arm64 || 386 || arm || riscv64 || loong64 || ppc64le))

package native

import (
	"github.com/gostafa/linuxdesktop/internal/domain"
)

// This file is the fallback for platforms purego does not support — Windows,
// which has no dlopen, and the architectures whose assembly support is
// incomplete. It carries the same two entry points as funcs.go under the
// inverse build constraint, so the package always compiles and callers never
// need a build tag of their own. Both report failure, which the caller reads
// as "use the filesystem answer".
//
// The functions live in vars.go only because this library constrains every
// package to the doc/funcs/consts/vars/types file names, and funcs.go is
// already taken by the real implementation.

// Vulkan reports nothing on an unsupported platform.
func Vulkan() (string, bool) { return "", false }

// OpenGL reports nothing on an unsupported platform.
func OpenGL() (domain.OpenGLInfo, bool) {
	var empty domain.OpenGLInfo

	return empty, false
}
