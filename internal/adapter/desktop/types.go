package desktop

import "github.com/gostafa/linuxdesktop/internal/port"

// Probe implements port.DesktopProbe. bus may be nil, which simply leaves the
// GNOME version unset.
type Probe struct {
	bus port.Bus
}
