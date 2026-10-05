// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package domain

type (
	// DesktopEnvironment is a recognized desktop environment.
	DesktopEnvironment string

	// DesktopInfo identifies the desktop environment.
	DesktopInfo struct {
		Environment     DesktopEnvironment `json:"environment"`
		Name            string             `json:"name"`
		Version         string             `json:"version,omitempty"`
		CurrentDesktop  string             `json:"current_desktop"`
		SessionDesktop  string             `json:"session_desktop"`
		DesktopSession  string             `json:"desktop_session"`
		CurrentDesktops []string           `json:"current_desktops"`
	}
)
