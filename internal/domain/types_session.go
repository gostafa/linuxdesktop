// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package domain

type (
	// SessionType is the kind of seat session the process is attached to.
	SessionType string

	// SessionInfo describes the logind seat session.
	SessionInfo struct {
		ID                string      `json:"id"`
		Seat              string      `json:"seat"`
		Type              SessionType `json:"type"`
		Desktop           string      `json:"desktop"`
		Name              string      `json:"name"`
		User              string      `json:"user"`
		XDGSessionType    string      `json:"xdg_session_type"`
		XDGSessionDesktop string      `json:"xdg_session_desktop"`
		VTNumber          int         `json:"vt_number"`
		Remote            bool        `json:"remote"`
		Active            bool        `json:"active"`
	}
)
