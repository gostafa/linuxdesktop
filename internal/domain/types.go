// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package domain

type (
	// Section selects which parts of an Environment to populate. Sections are
	// a bitmask so a caller can pay only for what it reads.
	Section uint16

	// Environment is the complete picture of the desktop the process is
	// running in.
	Environment struct {
		OS         OSInfo         `json:"os"`
		Session    SessionInfo    `json:"session"`
		Display    DisplayInfo    `json:"display"`
		Desktop    DesktopInfo    `json:"desktop"`
		Compositor CompositorInfo `json:"compositor"`
		Graphics   GraphicsInfo   `json:"graphics"`
		Portal     PortalInfo     `json:"portal"`
		Headless   bool           `json:"headless"`
	}
)
