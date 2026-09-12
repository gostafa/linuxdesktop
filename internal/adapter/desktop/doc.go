// Package desktop identifies the desktop environment.
//
// The identity itself is decided by package rules from the XDG variables, at
// no I/O cost. This adapter exists for the one field that cannot be derived
// from the environment: the version.
//
// Only two desktops report a version cheaply. GNOME Shell exposes a
// ShellVersion property on the session bus, which is a full version string.
// Plasma exports $KDE_SESSION_VERSION, which is a major version only ("6"),
// and is reported as such. Every other desktop would need a subprocess
// (xfce4-session --version and friends), which this library does not do, so
// their Version is left empty rather than guessed at.
package desktop
