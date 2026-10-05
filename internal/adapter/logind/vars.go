// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package logind

import (
	"errors"
)

// ErrNoSession is reported when neither logind nor the environment knows of a
// session. That is normal inside a container or a bare service unit, so the
// engine treats it as information rather than failure.
var ErrNoSession = errors.New("linuxdesktop: no logind session for this process")
