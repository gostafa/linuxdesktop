// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package core

import (
	"context"
	"testing"
	"time"

	"github.com/gostafa/linuxdesktop/internal/domain"
)

func TestOptionsAndDeadlines(t *testing.T) {
	t.Parallel()
	cfg := NewConfig(
		nil,
		WithTimeout(time.Second),
		WithProbeTimeout(time.Millisecond),
		WithSections(0),
		WithNativeGL(),
		WithProcessScan(),
	)
	if cfg.Timeout != time.Second || cfg.ProbeTimeout != time.Millisecond ||
		cfg.Sections != domain.SectionAll ||
		!cfg.NativeGL ||
		!cfg.ProcessScan {
		t.Fatal(cfg)
	}
	deps := Deps{}
	if out, err := New(&deps, &cfg).Detect(nil); err != nil || out == nil {
		t.Fatal(out, err)
	}
	eng := detector{cfg: cfg}
	err := timedProbe(&eng, func(ctx context.Context, gather *collector) error {
		if _, ok := ctx.Deadline(); !ok {
			t.Fatal("probe deadline absent")
		}
		return nil
	})(t.Context(), nil)
	if err != nil {
		t.Fatal(err)
	}
}
