// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package probe

import (
	"context"
	"errors"
	"testing"
)

func TestDelegates(t *testing.T) {
	t.Parallel()
	failure := errors.New("probe failure")
	ctx := t.Context()
	plain := func(got context.Context) (int, error) {
		if got != ctx {
			t.Fatal("context lost")
		}
		return 7, failure
	}
	env := 9
	withEnv := func(got context.Context, value *int) (int, error) {
		if value != &env {
			t.Fatal("environment lost")
		}
		return plain(got)
	}
	checks := []func() (int, error){
		func() (int, error) { return OSFunc[int](plain).OS(ctx) },
		func() (int, error) { return DetectFunc[int](plain).Detect(ctx) },
		func() (int, error) { return SessionFunc[int, int](withEnv).Session(ctx, &env) },
		func() (int, error) { return DesktopFunc[int, int](withEnv).Desktop(ctx, &env) },
		func() (int, error) { return X11Func[int, int](withEnv).X11(ctx, &env) },
		func() (int, error) { return WaylandFunc[int, int](withEnv).Wayland(ctx, &env) },
		func() (int, error) { return PortalFunc[int, int](withEnv).Portal(ctx, &env) },
	}
	for _, check := range checks {
		if got, err := check(); got != 7 || !errors.Is(err, failure) {
			t.Fatal(got, err)
		}
	}
	if got := SnapshotFunc[int](func() int { return 7 }).Snapshot(); got != 7 {
		t.Fatal(got)
	}
	gpu := GPUFunc[int](func(context.Context) (int, string, error) { return 7, "primary", failure })
	if got, primary, err := gpu.GPUs(
		ctx,
	); got != 7 || primary != "primary" ||
		!errors.Is(err, failure) {
		t.Fatal(got, primary, err)
	}
	stack := StackFunc[int, string](
		func(context.Context) (int, string, error) { return 7, "vulkan", failure },
	)
	if got, version, err := stack.Stack(
		ctx,
	); got != 7 || version != "vulkan" ||
		!errors.Is(err, failure) {
		t.Fatal(got, version, err)
	}
}
