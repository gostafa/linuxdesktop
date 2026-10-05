// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package linuxdesktop

import "github.com/gostafa/linuxdesktop/internal/domain"

// Session types, mirroring logind's TYPE= and $XDG_SESSION_TYPE.
const (
	SessionTypeUnknown = domain.SessionTypeUnknown
	SessionTypeWayland = domain.SessionTypeWayland
	SessionTypeX11     = domain.SessionTypeX11
	SessionTypeTTY     = domain.SessionTypeTTY
	SessionTypeMir     = domain.SessionTypeMir
)

// Display protocols.
const (
	DisplayProtocolUnknown = domain.DisplayProtocolUnknown
	DisplayProtocolWayland = domain.DisplayProtocolWayland
	DisplayProtocolX11     = domain.DisplayProtocolX11
)

// Recognized desktop environments.
const (
	DesktopUnknown  = domain.DesktopUnknown
	DesktopGNOME    = domain.DesktopGNOME
	DesktopKDE      = domain.DesktopKDE
	DesktopXFCE     = domain.DesktopXFCE
	DesktopCinnamon = domain.DesktopCinnamon
	DesktopMATE     = domain.DesktopMATE
	DesktopLXQt     = domain.DesktopLXQt
	DesktopLXDE     = domain.DesktopLXDE
	DesktopBudgie   = domain.DesktopBudgie
	DesktopCOSMIC   = domain.DesktopCOSMIC
	DesktopPantheon = domain.DesktopPantheon
)

// Recognized compositors and window managers. Anything outside this set is
// reported as CompositorUnknown with its real name in CompositorInfo.Name.
const (
	CompositorUnknown  = domain.CompositorUnknown
	CompositorMutter   = domain.CompositorMutter
	CompositorKWin     = domain.CompositorKWin
	CompositorSway     = domain.CompositorSway
	CompositorHyprland = domain.CompositorHyprland
	CompositorWayfire  = domain.CompositorWayfire
	CompositorRiver    = domain.CompositorRiver
	CompositorWeston   = domain.CompositorWeston
	CompositorLabwc    = domain.CompositorLabwc
	CompositorXfwm     = domain.CompositorXfwm
	CompositorMarco    = domain.CompositorMarco
	CompositorOpenbox  = domain.CompositorOpenbox
	CompositorI3       = domain.CompositorI3
	CompositorAwesome  = domain.CompositorAwesome
)

// Detection confidence levels.
const (
	ConfidenceUnknown = domain.ConfidenceUnknown
	ConfidenceLow     = domain.ConfidenceLow
	ConfidenceMedium  = domain.ConfidenceMedium
	ConfidenceHigh    = domain.ConfidenceHigh
)

// Detection methods.
const (
	DetectedUnknown     = domain.DetectedUnknown
	DetectedEnvironment = domain.DetectedEnvironment
	DetectedLogind      = domain.DetectedLogind
	DetectedDBus        = domain.DetectedDBus
	DetectedX11EWMH     = domain.DetectedX11EWMH
	DetectedWayland     = domain.DetectedWayland
	DetectedSocket      = domain.DetectedSocket
	DetectedProcess     = domain.DetectedProcess
)

// Sections of an Environment, for WithSections. Compose them with bitwise OR.
const (
	SectionOS         = domain.SectionOS
	SectionSession    = domain.SectionSession
	SectionDisplay    = domain.SectionDisplay
	SectionDesktop    = domain.SectionDesktop
	SectionCompositor = domain.SectionCompositor
	SectionGraphics   = domain.SectionGraphics
	SectionPortal     = domain.SectionPortal
	SectionAll        = domain.SectionAll
)

// ExtensionXWayland is the X11 extension that distinguishes an Xwayland server
// from a native X.Org one. It appears in X11Info.Extensions.
const ExtensionXWayland = domain.ExtensionXWayland

// goosLinux is what the Go runtime calls the one operating system this library
// has anything to say about.
const goosLinux = "linux"
