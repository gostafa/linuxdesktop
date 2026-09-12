package portal

import "github.com/gostafa/linuxdesktop/internal/port"

// Probe implements port.PortalProbe.
type Probe struct {
	bus port.Bus
}

// node is the subset of D-Bus introspection XML this package needs: the list
// of interfaces exported by an object.
type node struct {
	Interfaces []iface `xml:"interface"`
}

type iface struct {
	Name string `xml:"name,attr"`
}
