package logind

import "github.com/gostafa/linuxdesktop/internal/port"

// Probe implements port.SessionProbe.
//
// bus may be nil, in which case only the filesystem path is attempted. That is
// the common case on a normal desktop, where the session file always exists.
type Probe struct {
	bus port.Bus
}
