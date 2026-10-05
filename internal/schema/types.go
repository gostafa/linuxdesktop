// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package schema

type (
	// Environment is the shared Environment layout, parameterized by its component types.
	Environment[
		OSType,
		SessionType,
		DisplayType,
		DesktopType,
		CompositorType,
		GraphicsType,
		PortalType any,
	] struct {
		// OS operating system and kernel identity.
		OS OSType `json:"os"`
		// Session seat and login session details.
		Session SessionType `json:"session"`
		// Display reachable display servers and selected protocol.
		Display DisplayType `json:"display"`
		// Desktop desktop environment identity.
		Desktop DesktopType `json:"desktop"`
		// Compositor window manager or compositor identity and detection evidence.
		Compositor CompositorType `json:"compositor"`
		// Graphics gPU hardware and client graphics implementations.
		Graphics GraphicsType `json:"graphics"`
		// Portal availability and interfaces of the desktop portal.
		Portal PortalType `json:"portal"`
		// Headless whether neither display protocol is reachable.
		Headless bool `json:"headless"`
	}

	// SessionInfo is the shared SessionInfo layout, parameterized by its component types.
	SessionInfo[SessionType any] struct {
		// ID identifier reported by the underlying system.
		ID string `json:"id"`
		// Seat logind seat identifier, such as seat0.
		Seat string `json:"seat"`
		// Type windowing protocol associated with the session.
		Type SessionType `json:"type"`
		// Desktop desktop identifier reported by logind or the environment.
		Desktop string `json:"desktop"`
		// Name human-readable name reported by the source.
		Name string `json:"name"`
		// User user associated with the login session.
		User string `json:"user"`
		// XDGSessionType unmodified XDG_SESSION_TYPE value.
		XDGSessionType string `json:"xdg_session_type"`
		// XDGSessionDesktop unmodified XDG_SESSION_DESKTOP value.
		XDGSessionDesktop string `json:"xdg_session_desktop"`
		// VTNumber virtual terminal number; zero when unavailable.
		VTNumber int `json:"vt_number"`
		// Remote whether the session is remote, including SSH sessions.
		Remote bool `json:"remote"`
		// Active whether the session is currently active.
		Active bool `json:"active"`
	}

	// DisplayInfo is the shared DisplayInfo layout, parameterized by its component types.
	DisplayInfo[XServerType, WaylandServerType, ProtocolType any] struct {
		// X11 details of a reachable X server; nil when unavailable.
		X11 *XServerType `json:"x11,omitempty"`
		// Wayland details of a reachable Wayland compositor; nil when unavailable.
		Wayland *WaylandServerType `json:"wayland,omitempty"`
		// Protocol preferred display protocol for a client.
		Protocol ProtocolType `json:"protocol"`
		// WaylandDisplay wayland socket name from WAYLAND_DISPLAY.
		WaylandDisplay string `json:"wayland_display"`
		// X11Display x display address from DISPLAY.
		X11Display string `json:"x11_display"`
		// X11Available whether the X server probe reached a server.
		X11Available bool `json:"x11_available"`
		// WaylandAvailable whether the Wayland probe reached a compositor.
		WaylandAvailable bool `json:"wayland_available"`
		// XWayland whether the X server runs inside a Wayland session.
		XWayland bool `json:"xwayland"`
	}

	// WaylandInfo is the shared WaylandInfo layout, parameterized by its component types.
	WaylandInfo[GlobalType any] struct {
		// Display wayland socket name used for the connection.
		Display string `json:"display"`
		// Globals global interfaces advertised by the Wayland registry.
		Globals []GlobalType `json:"globals,omitempty"`
		// SocketFD inherited WAYLAND_SOCKET descriptor, when present.
		SocketFD int `json:"socket_fd,omitempty"`
	}
)
