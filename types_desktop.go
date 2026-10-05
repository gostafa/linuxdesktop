// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package linuxdesktop

type (
	// DesktopEnvironment is a recognized desktop environment.
	DesktopEnvironment string

	// DesktopInfo identifies the desktop environment.
	DesktopInfo = DesktopInfoData[DesktopEnvironment]

	// CompositorKind is a recognized window manager or Wayland compositor.
	CompositorKind string

	// DetectionConfidence grades how much a detection result can be trusted.
	DetectionConfidence string

	// DetectionMethod records which signal produced a detection result.
	DetectionMethod string
)
