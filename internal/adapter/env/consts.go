// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package env

// Standard XDG and display-server variables.
const (
	keySessionID      = "XDG_SESSION_ID"
	keySessionType    = "XDG_SESSION_TYPE"
	keySessionDesktop = "XDG_SESSION_DESKTOP"
	keySessionClass   = "XDG_SESSION_CLASS"
	keyCurrentDesktop = "XDG_CURRENT_DESKTOP"
	keyDesktopSession = "DESKTOP_SESSION"
	keyRuntimeDir     = "XDG_RUNTIME_DIR"
	keyConfigHome     = "XDG_CONFIG_HOME"
	keyHome           = "HOME"
	keySeat           = "XDG_SEAT"
	keyVTNR           = "XDG_VTNR"
	keyWaylandDisplay = "WAYLAND_DISPLAY"
	keyWaylandSocket  = "WAYLAND_SOCKET"
	keyDisplay        = "DISPLAY"
	keyUser           = "USER"
	keyLogname        = "LOGNAME"
)

// Compositor-specific variables. Each of these is set by exactly one
// compositor, which makes them the strongest and cheapest detection signal.
const (
	keyHyprlandSignature = "HYPRLAND_INSTANCE_SIGNATURE"
	keySwaySock          = "SWAYSOCK"
	keyWayfireSocket     = "WAYFIRE_SOCKET"
	keyI3Sock            = "I3SOCK"
	keyKDEFullSession    = "KDE_FULL_SESSION"
	keyKDESessionVersion = "KDE_SESSION_VERSION"
	keyGNOMESessionID    = "GNOME_DESKTOP_SESSION_ID"
	keyGNOMESetupDisplay = "GNOME_SETUP_DISPLAY"
)

// Remote-session markers.
const (
	keySSHConnection = "SSH_CONNECTION"
	keySSHTTY        = "SSH_TTY"
	keySSHClient     = "SSH_CLIENT"
)
