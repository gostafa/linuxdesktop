// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package portal

import (
	"context"

	"github.com/gostafa/linuxdesktop/internal/domain"
)

type (
	// Func delegates the probe to its configured function.
	Func[E, T any] func(ctx context.Context, env *E) (T, error)

	// Probe implements port.PortalProbe.
	Probe = Func[domain.Env, domain.PortalInfo]

	// node is the subset of D-Bus introspection XML this package needs: the
	// list of interfaces exported by an object.
	node = interfaceNode[iface]

	// interfaceNode decodes an XML interface listing with a caller-selected entry type.
	interfaceNode[I any] struct {
		// Interfaces interfaces exported by the introspected D-Bus object.
		Interfaces []I `xml:"interface"`
	}

	iface struct {
		// Name fully qualified D-Bus interface name.
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
