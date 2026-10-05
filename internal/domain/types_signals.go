// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package domain

type (
	// Env is every environment variable the library will ever consult, read
	// once at the start of a detection run so no probe pays for a repeated
	// lookup.
	Env struct {
		SessionID      string
		SessionType    string
		SessionDesktop string
		SessionClass   string
		CurrentDesktop string
		DesktopSession string
		WaylandDisplay string
		WaylandSocket  string
		Display        string
		RuntimeDir     string
		ConfigHome     string
		Home           string
		User           string
		Seat           string
		VTNR           string

		HyprlandSignature string
		SwaySock          string
		WayfireSocket     string
		I3Sock            string
		KDEFullSession    string
		KDESessionVersion string
		GNOMESessionID    string
		GNOMESetupDisplay string

		SSHConnection string
		SSHTTY        string
		SSHClient     string
	}

	// Signals carries the raw evidence classification rules reason over. It is
	// filled by the protocol adapters in stage one and consumed by package
	// rules in stage two.
	Signals struct {
		Env              Env
		Desktop          DesktopEnvironment
		X11WindowManager string
		WaylandGlobals   []WaylandGlobal
		X11Extensions    []string
		Processes        []string
		WaylandReachable bool
		X11Reachable     bool
	}
)
