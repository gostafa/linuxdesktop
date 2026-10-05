// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package linuxdesktop

type (
	// SessionType is the kind of seat session the process is attached to.
	SessionType string

	// SessionInfo describes the logind seat session.
	SessionInfo = SessionInfoData[SessionType]
)
