package linuxdesktop

import (
	"context"
	"runtime"
	"time"

	"github.com/gostafa/linuxdesktop/internal/adapter/dbusconn"
	"github.com/gostafa/linuxdesktop/internal/adapter/desktop"
	"github.com/gostafa/linuxdesktop/internal/adapter/drm"
	"github.com/gostafa/linuxdesktop/internal/adapter/env"
	"github.com/gostafa/linuxdesktop/internal/adapter/gl"
	"github.com/gostafa/linuxdesktop/internal/adapter/logind"
	"github.com/gostafa/linuxdesktop/internal/adapter/osinfo"
	"github.com/gostafa/linuxdesktop/internal/adapter/portal"
	"github.com/gostafa/linuxdesktop/internal/adapter/procscan"
	"github.com/gostafa/linuxdesktop/internal/adapter/wayland"
	"github.com/gostafa/linuxdesktop/internal/adapter/x11"
	"github.com/gostafa/linuxdesktop/internal/core"
	"github.com/gostafa/linuxdesktop/internal/domain"
	"github.com/gostafa/linuxdesktop/internal/port"
	"github.com/gostafa/linuxdesktop/internal/rules"
)

// Detect reports the complete desktop environment using the default options.
//
// The returned *Environment is never nil. The error reports probes that did
// not succeed and is safe to ignore; see the package documentation.
func Detect() (*Environment, error) {
	return DetectContext(context.Background())
}

// DetectContext reports the desktop environment, honouring ctx and opts.
func DetectContext(ctx context.Context, opts ...Option) (*Environment, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	cfg := core.NewConfig(opts...)

	if runtime.GOOS != "linux" {
		return &Environment{
			Display:  DisplayInfo{Protocol: DisplayProtocolUnknown},
			Session:  SessionInfo{Type: SessionTypeUnknown},
			Desktop:  DesktopInfo{Environment: DesktopUnknown},
			Headless: true,
		}, ErrNotLinux
	}

	bus := dbusconn.New(ctx)
	defer bus.Close()

	return core.New(adapters(bus, cfg), cfg).Detect(ctx)
}

// OS reports the distribution and kernel identity.
func OS() (OSInfo, error) {
	e, err := section(SectionOS)
	return e.OS, err
}

// Session reports the logind seat session.
func Session() (SessionInfo, error) {
	e, err := section(SectionSession)
	return e.Session, err
}

// Display reports which display servers are reachable, connecting to each one
// that is.
func Display() (DisplayInfo, error) {
	e, err := section(SectionDisplay)
	return e.Display, err
}

// Desktop reports the desktop environment.
func Desktop() (DesktopInfo, error) {
	e, err := section(SectionDesktop)
	return e.Desktop, err
}

// Compositor reports the compositor or window manager, together with the
// method that identified it and how far that method can be trusted.
func Compositor() (CompositorInfo, error) {
	e, err := section(SectionCompositor)
	return e.Compositor, err
}

// Graphics reports the GPUs and the client-side graphics stack. It uses the
// filesystem-only stack probe; for a live renderer string call DetectContext
// with WithOpenGL.
func Graphics() (GraphicsInfo, error) {
	e, err := section(SectionGraphics)
	return e.Graphics, err
}

// Portal reports whether xdg-desktop-portal is running and what it offers.
func Portal() (PortalInfo, error) {
	e, err := section(SectionPortal)
	return e.Portal, err
}

// IsWayland reports whether a Wayland compositor socket is present. It stats a
// path and nothing more, so it is safe on any hot path; Display connects and
// is therefore the authority.
func IsWayland() bool {
	if runtime.GOOS != "linux" {
		return false
	}
	snapshot := env.New().Snapshot()
	return wayland.SocketPath(&snapshot) != ""
}

// IsX11 reports whether an X server appears to be reachable. Like IsWayland it
// only stats; a display on a remote host is assumed reachable.
func IsX11() bool {
	if runtime.GOOS != "linux" {
		return false
	}
	snapshot := env.New().Snapshot()
	return x11.Available(&snapshot)
}

// IsHeadless reports whether there is no display server to draw on.
func IsHeadless() bool {
	return !IsWayland() && !IsX11()
}

// WithTimeout bounds the whole detection run. The default is DefaultTimeout.
func WithTimeout(d time.Duration) Option { return core.WithTimeout(d) }

// WithProbeTimeout bounds each individual probe, so one unresponsive server
// cannot consume the whole budget. The default is DefaultProbeTimeout.
func WithProbeTimeout(d time.Duration) Option { return core.WithProbeTimeout(d) }

// WithSections selects which parts of the Environment to populate. Probes for
// unselected sections never run.
func WithSections(s Section) Option { return core.WithSections(s) }

// WithOpenGL asks the graphics drivers directly instead of reading their
// manifests, which yields a true vendor, renderer and version at the cost of
// initialising the GPU stack — tens of milliseconds, not microseconds.
func WithOpenGL() Option { return core.WithNativeGL() }

// WithProcessScan enables the /proc fallback for compositor detection. It only
// affects the outcome when every stronger signal has failed, and a compositor
// found this way is reported with ConfidenceLow.
func WithProcessScan() Option { return core.WithProcessScan() }

// section runs a detection limited to one part of the Environment.
func section(s Section) (*Environment, error) {
	return DetectContext(context.Background(), WithSections(s))
}

// adapters is the composition root: the one place that decides which
// implementation satisfies each port.
func adapters(bus port.Bus, cfg core.Config) core.Deps {
	deps := core.Deps{
		Env:     env.New(),
		OS:      osinfo.New(),
		Session: logind.New(bus),
		X11:     x11.New(),
		Wayland: wayland.New(),
		Desktop: desktop.New(bus),
		GPUs:    drm.New(),
		Stack:   gl.New(cfg.NativeGL),
		Portal:  portal.New(bus),
	}
	if cfg.ProcessScan {
		deps.Process = procscan.New(rules.IsCompositorProcess)
	}
	return deps
}

// Compile-time proof that every adapter still satisfies the port it is wired
// to. These cost nothing at runtime and catch a broken signature at build time
// rather than at the injection site.
var (
	_ port.EnvProbe     = env.Probe{}
	_ port.OSProbe      = osinfo.Probe{}
	_ port.SessionProbe = (*logind.Probe)(nil)
	_ port.X11Probe     = x11.Probe{}
	_ port.WaylandProbe = wayland.Probe{}
	_ port.DesktopProbe = (*desktop.Probe)(nil)
	_ port.GPUProbe     = drm.Probe{}
	_ port.StackProbe   = gl.Probe{}
	_ port.PortalProbe  = (*portal.Probe)(nil)
	_ port.ProcessProbe = (*procscan.Probe)(nil)
	_ port.Bus          = (*dbusconn.Bus)(nil)

	_ = domain.SectionAll
)
