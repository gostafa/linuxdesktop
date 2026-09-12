// Package rules turns collected evidence into conclusions.
//
// Everything here is a pure function over domain.Signals. No file is read, no
// socket is opened and no clock is consulted, which is what lets the detection
// ladders be stated as data rather than control flow.
//
// The ladders are ordered by how much a signal can be trusted, and each rung
// records both the method that produced the answer and a confidence level, so
// a caller can tell an observed fact from a reasonable guess:
//
//	environment variable   high    only one compositor sets HYPRLAND_INSTANCE_SIGNATURE
//	Wayland registry       high    the advertised interfaces are a fingerprint
//	EWMH _NET_WM_NAME      high    the window manager names itself
//	XDG_CURRENT_DESKTOP    medium  the desktop implies its usual compositor
//	process scan           low     a matching process may not be serving us
//
// A compositor outside the CompositorKind enum, such as muffin or cosmic-comp,
// resolves to CompositorUnknown with its real name preserved in Name. Nothing
// is discarded, and promoting one to a constant is a one-line change here.
package rules
