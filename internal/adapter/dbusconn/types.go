// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package dbusconn

import (
	"context"
	"sync"

	"github.com/godbus/dbus/v5"
	"github.com/gostafa/linuxdesktop/internal/port"
)

type (
	// Bus implements port.Bus over two lazily established private connections.
	Bus = connectionBus[port.BusKind, *dbus.Conn, *port.Object, *port.PropertyQuery]

	connection interface {
		BusObject() dbus.BusObject
		Object(destination string, path dbus.ObjectPath) dbus.BusObject
	}

	// connectionBus separates bus operations from connection selection and lifetime.
	connectionBus[K any, C connection, O interface {
		Address() (destination, path string)
	}, Q interface {
		PropertyAddress() port.PropertyAddress
	}] struct {
		source interface {
			Get(kind K) (C, error)
			Shutdown() error
		}
	}

	// connectionSource supplies connection acquisition and release independently of bus operations.
	connectionSource[K, C any] struct {
		// Acquire selects or opens a connection.
		Acquire func(K) (C, error)
		// Release closes the connections owned by the source.
		Release func() error
	}

	// connections caches the result of opening each private connection once.
	connections[K, C any] struct {
		open  func(context.Context, K) (C, error)
		errs  [2]error
		conns [2]C
		once  [2]sync.Once
	}
)
