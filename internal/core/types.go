package core

import (
	"time"

	"github.com/gostafa/linuxdesktop/internal/domain"
	"github.com/gostafa/linuxdesktop/internal/port"
)

// Config controls a detection run.
type Config struct {
	// Timeout bounds the whole run. Zero means no overall bound.
	Timeout time.Duration
	// ProbeTimeout bounds each individual probe, so one wedged server cannot
	// consume the entire budget. Zero means each probe may use all of it.
	ProbeTimeout time.Duration
	// Sections selects which parts of the Environment to populate.
	Sections domain.Section
	// NativeGL enables the dlopen path in the graphics stack probe.
	NativeGL bool
	// ProcessScan enables the /proc fallback for compositor detection.
	ProcessScan bool
}

// Option mutates a Config.
type Option func(*Config)

// Deps are the adapters an Engine runs. Every field may be nil, in which case
// the corresponding data is simply left unset.
type Deps struct {
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
type Engine struct {
	deps Deps
	cfg  Config
}
