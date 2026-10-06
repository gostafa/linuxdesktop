// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package domain

// Detection flags and identities describe the selected sections and recognized desktop components.
const (
	// SectionOS selects operating system information. Combine the section flags
	// with bitwise OR to select multiple parts of an Environment.
	SectionOS Section = 1 << iota
	// SectionSession selects login session information.
	SectionSession
	// SectionDisplay selects display server information.
	SectionDisplay
	// SectionDesktop selects desktop environment information.
	SectionDesktop
	// SectionCompositor selects compositor identification.
	SectionCompositor
	// SectionGraphics selects GPU and graphics stack information.
	SectionGraphics
	// SectionPortal selects desktop portal information.
	SectionPortal

	// SectionAll selects every section and is the default.
	SectionAll Section = SectionOS | SectionSession | SectionDisplay |
		SectionDesktop | SectionCompositor | SectionGraphics | SectionPortal

	// SessionTypeUnknown indicates an unidentified session type. The session
	// constants mirror logind's TYPE= and $XDG_SESSION_TYPE.
	SessionTypeUnknown SessionType = "unknown"
	SessionTypeWayland SessionType = "wayland"
	SessionTypeX11     SessionType = "x11"
	SessionTypeTTY     SessionType = "tty"
	SessionTypeMir     SessionType = "mir"

	// DisplayProtocolUnknown indicates an unidentified display protocol.
	//nolint:goconst // Distinct enum types intentionally share their wire values.
	DisplayProtocolUnknown DisplayProtocol = DisplayProtocol(SessionTypeUnknown)
	//nolint:goconst // Distinct enum types intentionally share their wire values.
	DisplayProtocolWayland DisplayProtocol = DisplayProtocol(SessionTypeWayland)
	//nolint:goconst // Distinct enum types intentionally share their wire values.
	DisplayProtocolX11 DisplayProtocol = DisplayProtocol(SessionTypeX11)

	// DesktopUnknown indicates an unidentified desktop environment.
	//nolint:goconst // Distinct enum types intentionally share their wire values.
	DesktopUnknown  DesktopEnvironment = DesktopEnvironment(SessionTypeUnknown)
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

	// CompositorUnknown indicates an unrecognized compositor or window manager.
	// Its raw name is preserved in
	// CompositorInfo.Name, so no information is lost.
	//nolint:goconst // Distinct enum types intentionally share their wire values.
	CompositorUnknown  CompositorKind = CompositorKind(SessionTypeUnknown)
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

	// ConfidenceUnknown indicates that detection confidence is unavailable.
	//nolint:goconst // Distinct enum types intentionally share their wire values.
	ConfidenceUnknown DetectionConfidence = DetectionConfidence(SessionTypeUnknown)
	ConfidenceLow     DetectionConfidence = "low"
	ConfidenceMedium  DetectionConfidence = "medium"
	ConfidenceHigh    DetectionConfidence = "high"

	// DetectedUnknown indicates an unidentified detection method. The methods
	// are roughly ordered from cheapest to most invasive.
	//nolint:goconst // Distinct enum types intentionally share their wire values.
	DetectedUnknown     DetectionMethod = DetectionMethod(SessionTypeUnknown)
	DetectedEnvironment DetectionMethod = "environment"
	DetectedLogind      DetectionMethod = "logind"
	DetectedDBus        DetectionMethod = "dbus"
	DetectedX11EWMH     DetectionMethod = "x11_ewmh"
	DetectedWayland     DetectionMethod = "wayland_protocol"
	DetectedSocket      DetectionMethod = "socket"
	DetectedProcess     DetectionMethod = "process"

	// ExtensionXWayland is the X11 extension an Xwayland server advertises and a
	// native X.Org server does not. Its presence in X11Info.Extensions is the only
	// unambiguous way to tell the two apart.
	ExtensionXWayland = "XWAYLAND"
)
