// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package probe_test

import (
	"context"
	"errors"
	"testing"

	"github.com/gostafa/linuxdesktop/internal/probe"
)

func TestDefaultContext(t *testing.T) {
	t.Parallel()
	var empty context.Context
	empty = probe.Default(empty, func() context.Context { return t.Context() })
	if empty != t.Context() {
		t.Fatal("nil context did not receive the default")
	}
	parent, cancel := context.WithCancel(t.Context())
	ctx := parent
	ctx = probe.Default(
		ctx,
		func() context.Context { t.Fatal("fallback called for non-nil parent"); return nil },
	)
	cancel()
	if ctx != parent || !errors.Is(ctx.Err(), context.Canceled) {
		t.Fatal("parent context or cancellation lost")
	}
}

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
		func() (int, error) { return probe.OSFunc[int](plain).OS(ctx) },
		func() (int, error) { return probe.DetectFunc[int](plain).Detect(ctx) },
		func() (int, error) { return probe.SessionFunc[int, int](withEnv).Session(ctx, &env) },
		func() (int, error) { return probe.DesktopFunc[int, int](withEnv).Desktop(ctx, &env) },
		func() (int, error) { return probe.X11Func[int, int](withEnv).X11(ctx, &env) },
		func() (int, error) { return probe.WaylandFunc[int, int](withEnv).Wayland(ctx, &env) },
		func() (int, error) { return probe.PortalFunc[int, int](withEnv).Portal(ctx, &env) },
	}
	for _, check := range checks {
		if got, err := check(); got != 7 || !errors.Is(err, failure) {
			t.Fatal(got, err)
		}
	}
	failure = nil
	for _, check := range checks {
		if got, err := check(); got != 7 || err != nil {
			t.Fatal("successful probe changed its result", got, err)
		}
	}
	if got := probe.SnapshotFunc[int](func() int { return 7 }).Snapshot(); got != 7 {
		t.Fatal(got)
	}
	failure = errors.New("paired failure")
	gpu := probe.GPUFunc[int](
		func(context.Context) (int, string, error) { return 7, "primary", failure },
	)
	if got, primary, err := gpu.GPUs(
		ctx,
	); got != 7 || primary != "primary" ||
		!errors.Is(err, failure) {
		t.Fatal(got, primary, err)
	}
	stack := probe.StackFunc[int, string](
		func(context.Context) (int, string, error) { return 7, "vulkan", failure },
	)
	if got, version, err := stack.Stack(
		ctx,
	); got != 7 || version != "vulkan" ||
		!errors.Is(err, failure) {
		t.Fatal(got, version, err)
	}
	failure = nil
	if got, primary, err := gpu.GPUs(ctx); got != 7 || primary != "primary" || err != nil {
		t.Fatal(got, primary, err)
	}
	if got, version, err := stack.Stack(ctx); got != 7 || version != "vulkan" || err != nil {
		t.Fatal(got, version, err)
	}
}
