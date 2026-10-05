// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package osinfo

const (
	// Source files.
	pathOSRelease    = "/etc/os-release"
	pathOSReleaseAlt = "/usr/lib/os-release"
	pathKernelType   = "/proc/sys/kernel/ostype"
	pathKernelRelase = "/proc/sys/kernel/osrelease"
	pathKernelVer    = "/proc/sys/kernel/version"

	// os-release keys.
	keyName       = "NAME"
	keyPrettyName = "PRETTY_NAME"
	keyID         = "ID"
	keyIDLike     = "ID_LIKE"
	keyVersion    = "VERSION"
	keyVersionID  = "VERSION_ID"

	// goosLinux is what the runtime calls Linux, and kernelLinux is what uname
	// calls it. The second stands in when /proc is not mounted.
	goosLinux   = "linux"
	kernelLinux = "Linux"

	// noValue is an unreadable attribute or a field os-release did not carry.
	noValue = ""
)
