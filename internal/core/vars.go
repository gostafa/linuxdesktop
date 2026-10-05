// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package core

import (
	"time"

	"github.com/gostafa/linuxdesktop/internal/domain"
)

const (
	// DefaultTimeout bounds a run. The overall budget is generous enough that a healthy
	// desktop never approaches it, and the per-probe budget is short enough
	// that a single unresponsive server costs a fraction of a second rather
	// than the lot.
	DefaultTimeout = 2 * time.Second
	// DefaultProbeTimeout bounds each individual detection probe.
	DefaultProbeTimeout = 500 * time.Millisecond

	// sectionDisplayProbes is when the X11 and Wayland probes run. They feed
	// compositor detection as well as the display section, so either one
	// being wanted is reason enough.
	sectionDisplayProbes = domain.SectionDisplay | domain.SectionCompositor

	// sectionDesktopProbes is the same arrangement for the desktop probe.
	sectionDesktopProbes = domain.SectionDesktop | domain.SectionCompositor

	// zero is an unset section mask and an unbounded duration. It is one
	// constant rather than two because a package may not declare two constants
	// sharing a value.
	zero = 0
)

// DefaultConfig is the configuration used when the caller supplies no options.
func DefaultConfig() Config {
	return Config{
		Timeout:      DefaultTimeout,
		ProbeTimeout: DefaultProbeTimeout,
		Sections:     domain.SectionAll,
		NativeGL:     false,
		ProcessScan:  false,
	}
}
