// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package osinfo

// archNames maps Go's GOARCH spelling onto the machine name uname(2) reports,
// which is what users expect to see and what package managers key off.
func archNames(key string) (value string, ok bool) {
	switch key {
	case "386":
		return "i686", true
	case "amd64":
		return "x86_64", true
	case "arm":
		return "armv7l", true
	case "arm64":
		return "aarch64", true
	case "loong64":
		return "loongarch64", true
	case archMips:
		return archMips, true
	case archMips64:
		return archMips64, true
	case "mips64le":
		return "mips64el", true
	case "mipsle":
		return "mipsel", true
	case archPpc64:
		return archPpc64, true
	case archPpc64le:
		return archPpc64le, true
	case archRiscv64:
		return archRiscv64, true
	case archS390x:
		return archS390x, true
	default:
		return value, false
	}
}
