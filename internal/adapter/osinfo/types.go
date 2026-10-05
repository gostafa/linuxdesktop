// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package osinfo

import "github.com/gostafa/linuxdesktop/internal/domain"

type (
	// Probe implements port.OSProbe. It is stateless.
	Probe struct{}

	// releaseSetter fills one field of an OSInfo from one os-release value.
	releaseSetter func(*domain.OSInfo, string)
)
