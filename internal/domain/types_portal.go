// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package domain

type (
	// PortalInfo describes xdg-desktop-portal and which of its interfaces the
	// running backend actually exports.
	PortalInfo struct {
		Backend       string `json:"backend,omitempty"`
		Available     bool   `json:"available"`
		DesktopPortal bool   `json:"desktop_portal"`
		ScreenCast    bool   `json:"screencast"`
		Screenshot    bool   `json:"screenshot"`
		FileChooser   bool   `json:"file_chooser"`
		OpenURI       bool   `json:"open_uri"`
		RemoteDesktop bool   `json:"remote_desktop"`
		Inhibit       bool   `json:"inhibit"`
		Notification  bool   `json:"notification"`
	}
)
