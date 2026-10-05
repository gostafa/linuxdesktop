// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package domain

import "github.com/gostafa/linuxdesktop/internal/schema"

type (
	// CompositorKind is a recognized window manager or Wayland compositor.
	CompositorKind string

	// DetectionConfidence grades how much a detection result can be trusted.
	DetectionConfidence string

	// DetectionMethod records which signal produced a detection result.
	DetectionMethod string

	// CompositorInfo identifies the compositor or window manager, and says how
	// certain that identification is.
	CompositorInfo = schema.CompositorInfo[CompositorKind, DetectionConfidence, DetectionMethod]
)
