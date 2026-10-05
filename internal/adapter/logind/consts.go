// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package logind

const (
	// Filesystem sources.
	dirSessions = "/run/systemd/sessions"
	pathCgroup  = "/proc/self/cgroup"

	// Keys in /run/systemd/sessions/<id>.
	keyUID     = "UID"
	keyUser    = "USER"
	keyActive  = "ACTIVE"
	keyRemote  = "REMOTE"
	keySeat    = "SEAT"
	keyVTNR    = "VTNR"
	keyType    = "TYPE"
	keyDesktop = "DESKTOP"
	keyService = "SERVICE"
	keyState   = "STATE"

	// stateActive is the one STATE value that changes what is reported.
	stateActive = "active"

	// D-Bus fallback. The "auto" object path is logind's own alias for the
	// caller's session, so no GetSessionByPID call is needed.
	busName      = "org.freedesktop.login1"
	sessionPath  = "/org/freedesktop/login1/session/auto"
	sessionIface = "org.freedesktop.login1.Session"

	// Property names on org.freedesktop.login1.Session.
	propID      = "Id"
	propType    = "Type"
	propDesktop = "Desktop"
	propSeat    = "Seat"
	propVTNr    = "VTNr"
	propRemote  = "Remote"
	propActive  = "Active"
	propService = "Service"
	propName    = "Name"

	// cgroupPrefix and cgroupSuffix bracket the session id inside a cgroup path
	// such as 0::/user.slice/user-1000.slice/session-3.scope.
	cgroupPrefix = "session-"
	cgroupSuffix = ".scope"

	// zero is the empty length, the index of a leading field, and the virtual
	// terminal a session that is not on one reports. It is one constant rather
	// than three because a package may not declare two constants sharing a
	// value.
	zero = 0

	// noValue is an unset variable or a field logind did not supply.
	noValue = ""
)
