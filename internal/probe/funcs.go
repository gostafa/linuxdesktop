// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package probe

import (
	"context"
	"fmt"
)

// Snapshot delegates to the configured function.
func (run SnapshotFunc[T]) Snapshot() T { return run() }

// OS delegates to the configured function.
func (run OSFunc[T]) OS(ctx context.Context) (T, error) {
	result, err := run(ctx)
	if err != nil {
		return result, fmt.Errorf("probe: run: %w", err)
	}
	return result, nil
}

// Session delegates to the configured function.
func (run SessionFunc[E, T]) Session(
	ctx context.Context,
	env *E,
) (T, error) {
	result, err := run(ctx, env)
	if err != nil {
		return result, fmt.Errorf("probe: run: %w", err)
	}
	return result, nil
}

// Desktop delegates to the configured function.
func (run DesktopFunc[E, T]) Desktop(
	ctx context.Context,
	env *E,
) (T, error) {
	result, err := run(ctx, env)
	if err != nil {
		return result, fmt.Errorf("probe: run: %w", err)
	}
	return result, nil
}

// X11 delegates to the configured function.
func (run X11Func[E, T]) X11(ctx context.Context, env *E) (T, error) {
	result, err := run(ctx, env)
	if err != nil {
		return result, fmt.Errorf("probe: run: %w", err)
	}
	return result, nil
}

// Wayland delegates to the configured function.
func (run WaylandFunc[E, T]) Wayland(
	ctx context.Context,
	env *E,
) (T, error) {
	result, err := run(ctx, env)
	if err != nil {
		return result, fmt.Errorf("probe: run: %w", err)
	}
	return result, nil
}

// GPUs delegates to the configured function.
func (run GPUFunc[T]) GPUs(ctx context.Context) (result T, primary string, err error) {
	result, primary, err = run(ctx)
	if err != nil {
		return result, primary, fmt.Errorf("probe: GPUs: %w", err)
	}
	return result, primary, nil
}

// Stack delegates to the configured function.
func (run StackFunc[A, B]) Stack(ctx context.Context) (opengl A, vulkan B, err error) {
	opengl, vulkan, err = run(ctx)
	if err != nil {
		return opengl, vulkan, fmt.Errorf("probe: stack: %w", err)
	}
	return opengl, vulkan, nil
}

// Portal delegates to the configured function.
func (run PortalFunc[E, T]) Portal(ctx context.Context, env *E) (T, error) {
	result, err := run(ctx, env)
	if err != nil {
		return result, fmt.Errorf("probe: run: %w", err)
	}
	return result, nil
}

// Detect delegates to the configured function.
func (run DetectFunc[T]) Detect(ctx context.Context) (T, error) {
	result, err := run(ctx)
	if err != nil {
		return result, fmt.Errorf("probe: run: %w", err)
	}
	return result, nil
}
