// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package linuxdesktop

type (
	// Section selects which parts of an Environment to populate. Sections are
	// a bitmask so a caller can pay only for what it reads.
	Section uint16

	// Environment is the complete picture of the desktop the process is
	// running in.
	Environment = EnvironmentData[
		OSInfo,
		SessionInfo,
		DisplayInfo,
		DesktopInfo,
		CompositorInfo,
		GraphicsInfo,
		PortalInfo,
	]

	// Config controls a detection run. DetectContext initializes the defaults before
	// applying options in order; callers can supply custom Option functions.
	Config = ConfigData[Section]

	// Option configures a detection run by modifying its public configuration.
	Option func(*Config)
)
