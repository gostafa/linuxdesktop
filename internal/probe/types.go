// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package probe

import (
	"context"
)

type (
	// SnapshotFunc implements Snapshot by calling the supplied function.
	SnapshotFunc[T any] func() T

	// OSFunc implements OS by calling the supplied function.
	OSFunc[T any] func(ctx context.Context) (T, error)

	// SessionFunc implements Session by calling the supplied function.
	SessionFunc[E, T any] func(ctx context.Context, env *E) (T, error)

	// DesktopFunc implements Desktop by calling the supplied function.
	DesktopFunc[E, T any] func(ctx context.Context, env *E) (T, error)

	// DetectFunc implements Detect by calling the supplied function.
	DetectFunc[T any] func(ctx context.Context) (T, error)

	// X11Func implements X11 by calling the supplied function.
	X11Func[E, T any] func(ctx context.Context, env *E) (T, error)

	// WaylandFunc implements Wayland by calling the supplied function.
	WaylandFunc[E, T any] func(ctx context.Context, env *E) (T, error)

	// GPUFunc implements GPUs by calling the supplied function.
	GPUFunc[T any] func(ctx context.Context) (T, string, error)

	// StackFunc implements Stack by calling the supplied function.
	StackFunc[A, B any] func(ctx context.Context) (A, B, error)

	// PortalFunc implements Portal by calling the supplied function.
	PortalFunc[E, T any] func(ctx context.Context, env *E) (T, error)
)
