// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package linuxdesktop

import (
	"time"
)

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

	// zero is an empty selection or the initial length of an accumulating slice.
	zero = 0

	// Unknown indicates an unidentified session, display, desktop, compositor,
	// confidence, or detection method. It is assignable to each identity type.
	Unknown = "unknown"

	// SessionTypeWayland and the other session constants mirror logind's TYPE=
	// and $XDG_SESSION_TYPE.
	SessionTypeWayland SessionType = "wayland"
	SessionTypeX11     SessionType = "x11"
	SessionTypeTTY     SessionType = "tty"
	SessionTypeMir     SessionType = "mir"

	// DesktopGNOME and the other desktop constants identify recognized environments.
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

	// CompositorMutter and the other compositor constants identify recognized
	// compositors or window managers. For Unknown, CompositorInfo.Name preserves
	// the unrecognized raw name.
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

	// ConfidenceLow and the other confidence constants grade detection evidence.
	ConfidenceLow    DetectionConfidence = "low"
	ConfidenceMedium DetectionConfidence = "medium"
	ConfidenceHigh   DetectionConfidence = "high"

	// DetectedEnvironment and the other detection methods are roughly ordered
	// from cheapest to most invasive.
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

	// goosLinux is the operating system supported by detection.
	goosLinux = "linux"

	// DefaultTimeout bounds a whole detection run and
	// DefaultProbeTimeout bounds each probe within it, so one unresponsive server
	// cannot consume the entire budget.
	DefaultTimeout = 2 * time.Second
	// DefaultProbeTimeout bounds each individual detection probe.
	DefaultProbeTimeout = 500 * time.Millisecond
)
