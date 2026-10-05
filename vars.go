// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package linuxdesktop

import (
	"errors"
	"time"

	"github.com/gostafa/linuxdesktop/internal/core"
)

// ErrNotLinux is joined into the error from Detect when the program is not
// running on Linux. The returned Environment is still valid: it is zero-valued
// and marked headless, so a cross-platform caller can import this package
// unconditionally and branch on the result rather than on the build tag.
var ErrNotLinux = errors.New("linuxdesktop: not running on Linux")

// Default timeouts. DefaultTimeout bounds a whole detection run and
// DefaultProbeTimeout bounds each probe within it, so one unresponsive server
// cannot consume the entire budget.
const (
	DefaultTimeout      = core.DefaultTimeout
	DefaultProbeTimeout = core.DefaultProbeTimeout
)

// compile-time assertion that the exported defaults stay durations.
var _ time.Duration = DefaultTimeout
