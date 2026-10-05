// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package domain

type (
	// PortalInfo describes xdg-desktop-portal and which of its interfaces the
	// running backend actually exports.
	PortalInfo struct {
		// Backend preferred or matching portal implementation.
		Backend string `json:"backend,omitempty"`
		// Available whether the implementation or service was detected.
		Available bool `json:"available"`
		// DesktopPortal whether any desktop portal interface is exported.
		DesktopPortal bool `json:"desktop_portal"`
		// ScreenCast whether the ScreenCast interface is exported.
		ScreenCast bool `json:"screencast"`
		// Screenshot whether the Screenshot interface is exported.
		Screenshot bool `json:"screenshot"`
		// FileChooser whether the FileChooser interface is exported.
		FileChooser bool `json:"file_chooser"`
		// OpenURI whether the OpenURI interface is exported.
		OpenURI bool `json:"open_uri"`
		// RemoteDesktop whether the RemoteDesktop interface is exported.
		RemoteDesktop bool `json:"remote_desktop"`
		// Inhibit whether the Inhibit interface is exported.
		Inhibit bool `json:"inhibit"`
		// Notification whether the Notification interface is exported.
		Notification bool `json:"notification"`
	}
)
