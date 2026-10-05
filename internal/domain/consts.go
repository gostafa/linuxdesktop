// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package domain

// Session types, mirroring logind's TYPE= and $XDG_SESSION_TYPE.
const (
	SessionTypeUnknown SessionType = "unknown"
	SessionTypeWayland SessionType = "wayland"
	SessionTypeX11     SessionType = "x11"
	SessionTypeTTY     SessionType = "tty"
	SessionTypeMir     SessionType = "mir"
)

// Display protocols.
const (
	DisplayProtocolUnknown DisplayProtocol = "unknown"
	DisplayProtocolWayland DisplayProtocol = "wayland"
	DisplayProtocolX11     DisplayProtocol = "x11"
)

// Recognized desktop environments.
const (
	DesktopUnknown  DesktopEnvironment = "unknown"
	DesktopGNOME    DesktopEnvironment = "gnome"
	DesktopKDE      DesktopEnvironment = "kde"
	DesktopXFCE     DesktopEnvironment = "xfce"
	DesktopCinnamon DesktopEnvironment = "cinnamon"
	DesktopMATE     DesktopEnvironment = "mate"
	DesktopLXQt     DesktopEnvironment = "lxqt"
	DesktopLXDE     DesktopEnvironment = "lxde"
	DesktopBudgie   DesktopEnvironment = "budgie"
	DesktopCOSMIC   DesktopEnvironment = "cosmic"
	DesktopPantheon DesktopEnvironment = "pantheon"
)

// Recognized compositors and window managers. A compositor outside this set is
// reported as CompositorUnknown with its raw name preserved in
// CompositorInfo.Name, so no information is lost.
const (
	CompositorUnknown  CompositorKind = "unknown"
	CompositorMutter   CompositorKind = "mutter"
	CompositorKWin     CompositorKind = "kwin"
	CompositorSway     CompositorKind = "sway"
	CompositorHyprland CompositorKind = "hyprland"
	CompositorWayfire  CompositorKind = "wayfire"
	CompositorRiver    CompositorKind = "river"
	CompositorWeston   CompositorKind = "weston"
	CompositorLabwc    CompositorKind = "labwc"
	CompositorXfwm     CompositorKind = "xfwm4"
	CompositorMarco    CompositorKind = "marco"
	CompositorOpenbox  CompositorKind = "openbox"
	CompositorI3       CompositorKind = "i3"
	CompositorAwesome  CompositorKind = "awesome"
)

// Detection confidence levels.
const (
	ConfidenceUnknown DetectionConfidence = "unknown"
	ConfidenceLow     DetectionConfidence = "low"
	ConfidenceMedium  DetectionConfidence = "medium"
	ConfidenceHigh    DetectionConfidence = "high"
)

// Detection methods, roughly ordered from cheapest to most invasive.
const (
	DetectedUnknown     DetectionMethod = "unknown"
	DetectedEnvironment DetectionMethod = "environment"
	DetectedLogind      DetectionMethod = "logind"
	DetectedDBus        DetectionMethod = "dbus"
	DetectedX11EWMH     DetectionMethod = "x11_ewmh"
	DetectedWayland     DetectionMethod = "wayland_protocol"
	DetectedSocket      DetectionMethod = "socket"
	DetectedProcess     DetectionMethod = "process"
)

// ExtensionXWayland is the X11 extension an Xwayland server advertises and a
// native X.Org server does not. Its presence in X11Info.Extensions is the only
// unambiguous way to tell the two apart.
const ExtensionXWayland = "XWAYLAND"

// Sections of an Environment. Compose them with bitwise OR.
const (
	SectionOS Section = 1 << iota
	SectionSession
	SectionDisplay
	SectionDesktop
	SectionCompositor
	SectionGraphics
	SectionPortal

	// SectionAll selects every section and is the default.
	SectionAll Section = SectionOS | SectionSession | SectionDisplay |
		SectionDesktop | SectionCompositor | SectionGraphics | SectionPortal
)
