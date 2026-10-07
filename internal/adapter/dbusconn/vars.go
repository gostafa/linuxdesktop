// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package dbusconn

import (
	"errors"
)

// ErrBusKind reports a bus selector outside the supported range.
var (
	ErrBusKind = errors.New("linuxdesktop: unknown bus kind")
	// ErrClosed reports acquisition after the bus has been shut down.
	ErrClosed       = errors.New("linuxdesktop: bus closed")
	errNoConnection = errors.New("dbusconn: opener returned no connection")
)
