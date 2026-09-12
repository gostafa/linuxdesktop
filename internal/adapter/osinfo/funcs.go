package osinfo

import (
	"context"
	"os"
	"runtime"

	"github.com/gostafa/linuxdesktop/internal/domain"
	"github.com/gostafa/linuxdesktop/internal/sysfs"
)

// New returns an OS probe.
func New() Probe { return Probe{} }

// OS reads the distribution and kernel identity. A missing os-release is not
// an error: minimal containers legitimately have none, and the kernel fields
// are still worth reporting.
func (Probe) OS(_ context.Context) (domain.OSInfo, error) {
	info := domain.OSInfo{
		Kernel:        sysfs.Trimmed(pathKernelType),
		KernelRelease: sysfs.Trimmed(pathKernelRelase),
		Architecture:  architecture(),
	}
	if info.Kernel == "" && runtime.GOOS == "linux" {
		info.Kernel = "Linux"
	}
	if h, err := os.Hostname(); err == nil {
		info.Hostname = h
	}

	data, err := sysfs.Bytes(pathOSRelease)
	if err != nil {
		data, err = sysfs.Bytes(pathOSReleaseAlt)
	}
	if err != nil {
		return info, nil
	}

	sysfs.Each(data, func(key, value string) bool {
		switch key {
		case keyName:
			info.Name = value
		case keyPrettyName:
			info.PrettyName = value
		case keyID:
			info.ID = value
		case keyIDLike:
			info.IDLike = value
		case keyVersion:
			info.Version = value
		case keyVersionID:
			info.VersionID = value
		}
		return true
	})

	if info.PrettyName == "" {
		info.PrettyName = info.Name
	}
	return info, nil
}

// architecture reports the machine name in uname(2) spelling.
func architecture() string {
	if name, ok := archNames[runtime.GOARCH]; ok {
		return name
	}
	return runtime.GOARCH
}
