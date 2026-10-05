// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package linuxdesktop

type (
	// DesktopInfoData is the shared DesktopInfoData layout, parameterized by its component types.
	DesktopInfoData[DesktopType any] struct {
		// Environment recognized desktop environment.
		Environment DesktopType `json:"environment"`
		// Name human-readable name reported by the source.
		Name string `json:"name"`
		// Version version reported by the source, when available.
		Version string `json:"version,omitempty"`
		// CurrentDesktop unmodified XDG_CURRENT_DESKTOP value.
		CurrentDesktop string `json:"current_desktop"`
		// SessionDesktop unmodified XDG_SESSION_DESKTOP value.
		SessionDesktop string `json:"session_desktop"`
		// DesktopSession unmodified DESKTOP_SESSION value.
		DesktopSession string `json:"desktop_session"`
		// CurrentDesktops ordered desktop tokens from XDG_CURRENT_DESKTOP.
		CurrentDesktops []string `json:"current_desktops"`
	}

	// CompositorInfoData is the shared CompositorInfoData layout, parameterized by its component types.
	CompositorInfoData[KindType, ConfidenceType, MethodType any] struct {
		Kind       KindType       `json:"kind"`
		Confidence ConfidenceType `json:"confidence"`
		DetectedBy MethodType     `json:"detected_by"`
		Name       string         `json:"name"`
		Version    string         `json:"version,omitempty"`
		Wayland    bool           `json:"wayland"`
		X11        bool           `json:"x11"`
	}
)
