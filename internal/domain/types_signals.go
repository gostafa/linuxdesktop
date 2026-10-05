// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package domain

type (
	// Env is every environment variable the library will ever consult, read
	// once at the start of a detection run so no probe pays for a repeated
	// lookup.
	Env struct {
		// SessionID unmodified XDG_SESSION_ID value.
		SessionID string
		// SessionType unmodified XDG_SESSION_TYPE value.
		SessionType string
		// SessionDesktop unmodified XDG_SESSION_DESKTOP value.
		SessionDesktop string
		// SessionClass unmodified XDG_SESSION_CLASS value.
		SessionClass string
		// CurrentDesktop unmodified XDG_CURRENT_DESKTOP value.
		CurrentDesktop string
		// DesktopSession unmodified DESKTOP_SESSION value.
		DesktopSession string
		// WaylandDisplay wayland socket name from WAYLAND_DISPLAY.
		WaylandDisplay string
		// WaylandSocket unmodified WAYLAND_SOCKET descriptor value.
		WaylandSocket string
		// Display unmodified DISPLAY value.
		Display string
		// RuntimeDir unmodified XDG_RUNTIME_DIR value.
		RuntimeDir string
		// ConfigHome unmodified XDG_CONFIG_HOME value.
		ConfigHome string
		// Home unmodified HOME value.
		Home string
		// User uSER value, falling back to LOGNAME.
		User string
		// Seat unmodified XDG_SEAT value.
		Seat string
		// VTNR unmodified XDG_VTNR value.
		VTNR string

		// HyprlandSignature unmodified HYPRLAND_INSTANCE_SIGNATURE value.
		HyprlandSignature string
		// SwaySock unmodified SWAYSOCK value.
		SwaySock string
		// WayfireSocket unmodified WAYFIRE_SOCKET value.
		WayfireSocket string
		// I3Sock unmodified I3SOCK value.
		I3Sock string
		// KDEFullSession unmodified KDE_FULL_SESSION marker.
		KDEFullSession string
		// KDESessionVersion unmodified KDE_SESSION_VERSION value.
		KDESessionVersion string
		// GNOMESessionID unmodified GNOME_DESKTOP_SESSION_ID marker.
		GNOMESessionID string
		// GNOMESetupDisplay unmodified GNOME_SETUP_DISPLAY value.
		GNOMESetupDisplay string

		// SSHConnection unmodified SSH_CONNECTION value.
		SSHConnection string
		// SSHTTY unmodified SSH_TTY value.
		SSHTTY string
		// SSHClient unmodified SSH_CLIENT value.
		SSHClient string
	}

	// Signals carries the raw evidence classification rules reason over. It is
	// filled by the protocol adapters in stage one and consumed by package
	// rules in stage two.
	Signals = evidence[Env, DesktopEnvironment, WaylandGlobal]

	// evidence combines environment, desktop, and registry signals for classification.
	evidence[E, D, G any] struct {
		// E environment snapshot taken before concurrent probes run.
		Env E
		// Desktop desktop environment identified from collected evidence.
		Desktop D
		// X11WindowManager window manager name observed through the X11 probe.
		X11WindowManager string
		// WaylandGlobals registry entries observed through the Wayland probe.
		WaylandGlobals []G
		// X11Extensions extension names observed through the X11 probe.
		X11Extensions []string
		// Processes filtered command names observed for the current user.
		Processes []string
		// WaylandReachable whether the Wayland probe connected successfully.
		WaylandReachable bool
		// X11Reachable whether the X11 probe connected successfully.
		X11Reachable bool
	}
)
