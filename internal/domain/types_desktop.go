// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package domain

import "github.com/gostafa/linuxdesktop/internal/schema"

type (
	// DesktopEnvironment is a recognized desktop environment.
	DesktopEnvironment string

	// DesktopInfo identifies the desktop environment.
	DesktopInfo = schema.DesktopInfo[DesktopEnvironment]
)
