// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package drm

import (
	"github.com/gostafa/linuxdesktop/internal/domain"
	"github.com/gostafa/linuxdesktop/internal/probe"
)

type (
	// Probe enumerates DRM devices. It is stateless.
	Probe = probe.GPUFunc[[]domain.GPUInfo]

	// pciScan streams the hwdata device list, naming GPUs as their lines go
	// past. devices is the device table of the vendor the scan is currently
	// inside, and left counts the names still wanted, so the scan can stop as
	// soon as it has them all.
	pciScan struct {
		pending  map[string]map[string][]int
		devices  map[string][]int
		setModel func(int, string)
		left     int
	}
)
