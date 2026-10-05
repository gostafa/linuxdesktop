// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package linuxdesktop

type (
	// CompositorInfo identifies the compositor or window manager, and says how
	// certain that identification is.
	CompositorInfo = CompositorInfoData[CompositorKind, DetectionConfidence, DetectionMethod]
)
