// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package desktop

import (
	"context"

	"github.com/gostafa/linuxdesktop/internal/domain"
)

type (
	// Func delegates the probe to its configured function.
	Func[E, T any] func(ctx context.Context, env *E) (T, error)

	// Probe implements port.DesktopProbe. bus may be nil, which simply leaves the
	// GNOME version unset.
	Probe = Func[domain.Env, domain.DesktopInfo]
)
