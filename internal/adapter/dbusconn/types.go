// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package dbusconn

import (
	"context"
	"sync"

	"github.com/godbus/dbus/v5"
)

// Bus implements port.Bus over two lazily established private connections.
//
// The stored context is the connection lifetime, not a per-call deadline: it
// is the overall detection context, so a wedged bus cannot outlive the run.
// Individual calls carry their own, shorter deadlines.
type Bus struct {
	errs  [2]error
	base  context.Context
	conns [2]*dbus.Conn
	once  [2]sync.Once
}
