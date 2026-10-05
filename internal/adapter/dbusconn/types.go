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
	Bus = connectionBus[port.BusKind, *dbus.Conn]

	connection interface {
		BusObject() dbus.BusObject
		Object(destination string, path dbus.ObjectPath) dbus.BusObject
	}

	// connectionBus separates bus operations from connection selection and lifetime.
	connectionBus[K any, C connection] struct {
		connect func(K) (C, error)
		close   func() error
	}

	// connections caches the result of opening each private connection once.
	connections struct {
		open  func(context.Context, port.BusKind) (*dbus.Conn, error)
		errs  [2]error
		conns [2]*dbus.Conn
		once  [2]sync.Once
	}
)
