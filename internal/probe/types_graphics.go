// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package probe

import "context"

type (
	// GPUFunc implements GPUs by calling the supplied function.
	GPUFunc[T any] func(ctx context.Context) (T, string, error)

	// StackFunc implements Stack by calling the supplied function.
	StackFunc[A, B any] func(ctx context.Context) (A, B, error)

	// PortalFunc implements Portal by calling the supplied function.
	PortalFunc[E, T any] func(ctx context.Context, env *E) (T, error)
)
