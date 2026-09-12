package osinfo

// archNames maps Go's GOARCH spelling onto the machine name uname(2) reports,
// which is what users expect to see and what package managers key off.
var archNames = map[string]string{
	"386":      "i686",
	"amd64":    "x86_64",
	"arm":      "armv7l",
	"arm64":    "aarch64",
	"loong64":  "loongarch64",
	"mips":     "mips",
	"mips64":   "mips64",
	"mips64le": "mips64el",
	"mipsle":   "mipsel",
	"ppc64":    "ppc64",
	"ppc64le":  "ppc64le",
	"riscv64":  "riscv64",
	"s390x":    "s390x",
}
