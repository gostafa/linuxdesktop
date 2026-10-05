// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package linuxdesktop

import (
	"errors"
	"runtime"
	"testing"

	"github.com/gostafa/linuxdesktop/internal/core"
)

func TestPublicSections(t *testing.T) {
	t.Parallel()
	if result, err := Detect(); result == nil ||
		(runtime.GOOS != goosLinux && !errors.Is(err, ErrNotLinux)) {
		t.Fatal(result, err)
	}
	for _, call := range []func() error{
		func() error { _, err := OS(); return err }, func() error { _, err := Session(); return err },
		func() error { _, err := Display(); return err }, func() error { _, err := Desktop(); return err },
		func() error { _, err := Compositor(); return err }, func() error { _, err := Graphics(); return err },
		func() error { _, err := Portal(); return err },
	} {
		if err := call(); runtime.GOOS != goosLinux && !errors.Is(err, ErrNotLinux) {
			t.Fatal("section helper lost platform error", err)
		}
	}
	cfg := newConfig(WithSections(SectionOS|SectionGraphics), WithProcessScan(), WithOpenGL())
	result, err := detectOn(t.Context(), cfg, goosLinux)
	if result == nil || err != nil {
		t.Fatal(result, err)
	}
	if _, err = detectOn(t.Context(), cfg, "other"); !errors.Is(err, ErrNotLinux) {
		t.Fatal(err)
	}
	if onLinux(
		"other",
		func() bool { t.Fatal("probe ran on unsupported platform"); return true },
	) ||
		!onLinux(goosLinux, func() bool { return true }) {
		t.Fatal("platform guard")
	}
	_ = waylandAvailable()
	_ = x11Available()
	if graphicsProbe(&core.Config{}) == nil {
		t.Fatal("default graphics probe missing")
	}
}
