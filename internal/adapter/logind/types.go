// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package logind

import (
	"context"

	"github.com/gostafa/linuxdesktop/internal/domain"
)

type (
	// Func delegates the probe to its configured function.
	Func[E, T any] func(ctx context.Context, env *E) (T, error)

	// Probe implements port.SessionProbe.
	//
	// bus may be nil, in which case only the filesystem path is attempted. That
	// is the common case on a normal desktop, where the session file always
	// exists.
	Probe = Func[domain.Env, domain.SessionInfo]

	// fileSetter fills one field of a SessionInfo from one value of a session
	// mirror.
	fileSetter func(*domain.SessionInfo, string)
)
