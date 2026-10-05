// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package osinfo

// archNames maps Go's GOARCH spelling onto the machine name uname(2) reports,
// which is what users expect to see and what package managers key off.
func archNames(key string) (value string, ok bool) {
	tables := []map[string]string{archNamesTableA(), archNamesTableB()}

	for i := range tables {
		if found, exists := tables[i][key]; exists {
			return found, true
		}
	}

	return value, false
}

func archNamesTableA() map[string]string {
	return map[string]string{
		"386":      "i686",
		"amd64":    "x86_64",
		"arm":      "armv7l",
		"arm64":    "aarch64",
		"loong64":  "loongarch64",
		archMips:   archMips,
		archMips64: archMips64,
		"mips64le": "mips64el",
		"mipsle":   "mipsel",
		archPpc64:  archPpc64,
	}
}

func archNamesTableB() map[string]string {
	return map[string]string{
		archPpc64le: archPpc64le,
		archRiscv64: archRiscv64,
		archS390x:   archS390x,
	}
}
