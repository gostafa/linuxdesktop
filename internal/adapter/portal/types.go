// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package portal

import (
	"github.com/gostafa/linuxdesktop/internal/domain"
	"github.com/gostafa/linuxdesktop/internal/port"
)

type (
	// Probe implements port.PortalProbe.
	Probe struct {
		bus port.Bus
	}

	// node is the subset of D-Bus introspection XML this package needs: the
	// list of interfaces exported by an object.
	node struct {
		Interfaces []iface `xml:"interface"`
	}

	iface struct {
		Name string `xml:"name,attr"`
	}

	// backendFile is one installed .portal file: the implementation it names
	// and the desktops it declares itself for.
	backendFile struct {
		name string
		uses []string
	}

	// flagSetter records that the running portal exports one interface.
	flagSetter func(*domain.PortalInfo)
)
