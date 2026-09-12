package osinfo

// Source files.
const (
	pathOSRelease    = "/etc/os-release"
	pathOSReleaseAlt = "/usr/lib/os-release"
	pathKernelType   = "/proc/sys/kernel/ostype"
	pathKernelRelase = "/proc/sys/kernel/osrelease"
	pathKernelVer    = "/proc/sys/kernel/version"
)

// os-release keys.
const (
	keyName       = "NAME"
	keyPrettyName = "PRETTY_NAME"
	keyID         = "ID"
	keyIDLike     = "ID_LIKE"
	keyVersion    = "VERSION"
	keyVersionID  = "VERSION_ID"
)
