// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package domain

import (
	"github.com/gostafa/linuxdesktop/internal/schema"
)

type (
	// Section selects which parts of an Environment to populate. Sections are
	// a bitmask so a caller can pay only for what it reads.
	Section uint16

	// Environment is the complete picture of the desktop the process is
	// running in.
	Environment = schema.Environment[
		OSInfo,
		SessionInfo,
		DisplayInfo,
		DesktopInfo,
		CompositorInfo,
		GraphicsInfo,
		PortalInfo,
	]
)
