// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package dbusconn

import (
	"context"
	"sync"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/gostafa/linuxdesktop/internal/port"
	"github.com/gostafa/linuxdesktop/internal/schema"
	"github.com/gostafa/singleton"
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
			Get(ctx context.Context, kind K) (C, error)
			Shutdown() error
		}
	}

	// connectionSource supplies connection acquisition and release independently of bus operations.
	connectionSource[K, C any] struct {
		// Acquire selects or opens a connection.
		Acquire func(context.Context, K) (C, error)
		// Release closes the connections owned by the source.
		Release func() error
	}

	// connections caches the result of opening each private connection once.
	connections = connectionCache[port.BusKind, *dbus.Conn]

	// connectionCache binds lazy initialization to an explicitly owned lifetime.
	//nolint:reusability // The owning cache combines generic resources with concrete synchronization and retry policy.
	connectionCache[K, C any] struct {
		//nolint:containedctx // Private connections share the owning run's lifetime.
		base      context.Context
		providers [2]*singleton.Provider[C]
		conns     [2]C
		open      func(context.Context, K) (C, error)
		cancel    context.CancelFunc
		policy    schema.RetryPolicy
		timeout   time.Duration
		lock      sync.Mutex
		closed    bool
	}
)
