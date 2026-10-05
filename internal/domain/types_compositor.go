// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package domain

type (
	// CompositorKind is a recognized window manager or Wayland compositor.
	CompositorKind string

	// DetectionConfidence grades how much a detection result can be trusted.
	DetectionConfidence string

	// DetectionMethod records which signal produced a detection result.
	DetectionMethod string

	// CompositorInfo identifies the compositor or window manager, and says how
	// certain that identification is.
	CompositorInfo struct {
		Kind       CompositorKind      `json:"kind"`
		Name       string              `json:"name"`
		Version    string              `json:"version,omitempty"`
		Confidence DetectionConfidence `json:"confidence"`
		DetectedBy DetectionMethod     `json:"detected_by"`
		Wayland    bool                `json:"wayland"`
		X11        bool                `json:"x11"`
	}
)
