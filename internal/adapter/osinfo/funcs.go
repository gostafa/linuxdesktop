// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

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
	info := kernelInfo()

	applyOSRelease(&info)

	if info.PrettyName == noValue {
		info.PrettyName = info.Name
	}

	return info, nil
}

// kernelInfo is everything the kernel and the runtime know without reading
// os-release.
func kernelInfo() domain.OSInfo {
	info := domain.OSInfo{
		Kernel:        sysfs.Trimmed(pathKernelType),
		KernelRelease: sysfs.Trimmed(pathKernelRelase),
		Architecture:  architecture(),
		Hostname:      hostname(),
	}
	if info.Kernel == noValue && runtime.GOOS == goosLinux {
		info.Kernel = kernelLinux
	}

	return info
}

// applyOSRelease overlays whatever the distribution says about itself.
func applyOSRelease(info *domain.OSInfo) {
	data, err := sysfs.Bytes(pathOSRelease)
	if err != nil {
		data, err = sysfs.Bytes(pathOSReleaseAlt)
	}

	if err != nil {
		return
	}

	setters := releaseSetters()

	sysfs.Each(data, func(key, value string) bool {
		set, ok := setters[key]
		if ok {
			set(info, value)
		}

		return true
	})
}

// releaseSetters maps each os-release key onto the field it fills. A key that
// is absent from the table is one this library has no use for.
func releaseSetters() map[string]releaseSetter {
	return map[string]releaseSetter{
		keyName:       func(info *domain.OSInfo, value string) { info.Name = value },
		keyPrettyName: func(info *domain.OSInfo, value string) { info.PrettyName = value },
		keyID:         func(info *domain.OSInfo, value string) { info.ID = value },
		keyIDLike:     func(info *domain.OSInfo, value string) { info.IDLike = value },
		keyVersion:    func(info *domain.OSInfo, value string) { info.Version = value },
		keyVersionID:  func(info *domain.OSInfo, value string) { info.VersionID = value },
	}
}

func hostname() string {
	name, err := os.Hostname()
	if err != nil {
		return noValue
	}

	return name
}

// architecture reports the machine name in uname(2) spelling.
func architecture() string {
	if name, ok := archNames[runtime.GOARCH]; ok {
		return name
	}

	return runtime.GOARCH
}
