package logind

// Filesystem sources.
const (
	dirSessions = "/run/systemd/sessions"
	pathCgroup  = "/proc/self/cgroup"
)

// Keys in /run/systemd/sessions/<id>.
const (
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
)

// D-Bus fallback. The "auto" object path is logind's own alias for the
// caller's session, so no GetSessionByPID call is needed.
const (
	busName      = "org.freedesktop.login1"
	sessionPath  = "/org/freedesktop/login1/session/auto"
	sessionIface = "org.freedesktop.login1.Session"
)

// Property names on org.freedesktop.login1.Session.
const (
	propID      = "Id"
	propType    = "Type"
	propDesktop = "Desktop"
	propSeat    = "Seat"
	propVTNr    = "VTNr"
	propRemote  = "Remote"
	propActive  = "Active"
	propService = "Service"
	propName    = "Name"
)

// cgroupMarker brackets the session id inside a cgroup path such as
// 0::/user.slice/user-1000.slice/session-3.scope.
const (
	cgroupPrefix = "session-"
	cgroupSuffix = ".scope"
)
