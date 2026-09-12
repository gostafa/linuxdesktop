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
	base  context.Context
	once  [2]sync.Once
	conns [2]*dbus.Conn
	errs  [2]error
}
