// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package probe

import "context"

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
)
