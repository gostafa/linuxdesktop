// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package dbusconn

import (
	"errors"
)

// ErrBusKind reports a bus selector outside the supported range.
var ErrBusKind = errors.New("linuxdesktop: unknown bus kind")
