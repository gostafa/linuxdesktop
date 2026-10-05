// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package core

import (
	"context"
	"sync"
	"time"

	"github.com/gostafa/linuxdesktop/internal/domain"
	"github.com/gostafa/linuxdesktop/internal/port"
)

type (
	// Config controls a detection run.
	Config struct {
		// Timeout bounds the whole run. Zero means no overall bound.
		Timeout time.Duration
		// ProbeTimeout bounds each individual probe, so one wedged server
		// cannot consume the entire budget. Zero means each probe may use all
		// of it.
		ProbeTimeout time.Duration
		// Sections selects which parts of the Environment to populate.
		Sections domain.Section
		// NativeGL enables the dlopen path in the graphics stack probe.
		NativeGL bool
		// ProcessScan enables the /proc fallback for compositor detection.
		ProcessScan bool
	}

	// Option mutates a Config.
	Option func(*Config)

	// Deps are the adapters an Engine runs. Every field may be nil, in which
	// case the corresponding data is simply left unset.
	Deps struct {
		Env     port.EnvProbe
		OS      port.OSProbe
		Session port.SessionProbe
		X11     port.X11Probe
		Wayland port.WaylandProbe
		Desktop port.DesktopProbe
		GPUs    port.GPUProbe
		Stack   port.StackProbe
		Portal  port.PortalProbe
		Process port.ProcessProbe
	}

	// Engine runs a configured set of probes.
	Engine struct {
		deps Deps
		cfg  Config
	}

	// probeFunc is one adapter call, already bound to the engine that owns it.
	probeFunc func(context.Context, *collector) error

	// step is one probe the engine may run: the work it does, the sections it
	// contributes to, and whether its adapter was wired in at all.
	step struct {
		run     probeFunc
		section domain.Section
		wired   bool
	}

	// collector gathers what the probes find. Each one runs on its own
	// goroutine, so every field here is written while holding lock.
	collector struct {
		out     *domain.Environment
		sig     *domain.Signals
		x11     *domain.X11Info
		wayland *domain.WaylandInfo
		failed  []error
		wait    sync.WaitGroup
		lock    sync.Mutex
	}
)
