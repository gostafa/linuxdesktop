// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package probe

import "context"

type (
	// X11Func implements X11 by calling the supplied function.
	X11Func[E, T any] func(ctx context.Context, env *E) (T, error)

	// WaylandFunc implements Wayland by calling the supplied function.
	WaylandFunc[E, T any] func(ctx context.Context, env *E) (T, error)
)
