// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package dbusconn

import (
	"time"
)

const (
	defaultInitializationTimeout = 500 * time.Millisecond
	singleAttempt                = 1
	retryMultiplier              = 2
	retryJitter                  = 0.2
	//nolint:goconst // The number of bus instances is independent of the retry multiplier.
	busCount = 2

	// zero is an empty collection or an unbounded initialization duration.
	zero = 0

	errCheckOwner     = "dbus: check name owner: %w"
	errIntrospect     = "dbus: introspect object: %w"
	errReadProperty   = "dbus: read property: %w"
	errReadProperties = "dbus: read properties: %w"
)
