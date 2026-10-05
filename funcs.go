// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package linuxdesktop

import (
	"context"
	"fmt"
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
	"github.com/gostafa/linuxdesktop/internal/port"
	"github.com/gostafa/linuxdesktop/internal/rules"
)

var (
	// Compile-time proof that every adapter still satisfies the port it is wired
	// to. These cost nothing at runtime and catch a broken signature at build time
	// rather than at the injection site.
	_ port.EnvProbe     = env.Probe(nil)
	_ port.OSProbe      = osinfo.Probe(nil)
	_ port.SessionProbe = logind.Probe(nil)
	_ port.X11Probe     = x11.Probe(nil)
	_ port.WaylandProbe = wayland.Probe(nil)
	_ port.DesktopProbe = desktop.Probe(nil)
	_ port.GPUProbe     = drm.Probe(nil)
	_ port.StackProbe   = gl.Probe(nil)
	_ port.PortalProbe  = portal.Probe(nil)
	_ port.ProcessProbe = (*procscan.Probe)(nil)
	_ port.Bus          = (*dbusconn.Bus)(nil)
)

// Detect reports the complete desktop environment using the default options.
//
// The returned *Environment is never nil. The error reports probes that did
// not succeed and is safe to ignore; see the package documentation.
func Detect() (*Environment, error) {
	result, callErr := DetectContext(context.Background())

	if callErr != nil {
		return result, fmt.Errorf("linuxdesktop: detect environment: %w", callErr)
	}
	return result, nil
}

// DetectContext reports the desktop environment, honoring ctx and opts.
func DetectContext(ctx context.Context, opts ...Option) (*Environment, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := detectOn(ctx, newConfig(opts...), runtime.GOOS)
	if err != nil {
		return result, fmt.Errorf("linuxdesktop: detect context: %w", err)
	}
	return result, nil
}

func detectOn(ctx context.Context, cfg *core.Config, goos string) (*Environment, error) {
	if goos != goosLinux {
		return elsewhere(), ErrNotLinux
	}

	bus := dbusconn.New(ctx)
	defer bus.Close()

	deps := adapters(bus, cfg)

	result, err := core.New(&deps, cfg).Detect(ctx)

	return publicEnvironment(result), err
}

// elsewhere is the answer on a system that has no Linux desktop to describe: a
// zero-valued Environment marked headless, so a cross-platform caller can
// import this package unconditionally and branch on the result.
func elsewhere() *Environment {
	result := new(Environment)

	result.Display.Protocol = DisplayProtocolUnknown
	result.Session.Type = SessionTypeUnknown
	result.Desktop.Environment = DesktopUnknown
	result.Headless = true

	return result
}

// OS reports the distribution and kernel identity.
func OS() (OSInfo, error) {
	e, err := section(SectionOS)
	if err != nil {
		return e.OS, fmt.Errorf("linuxdesktop: detect operating system: %w", err)
	}
	return e.OS, nil
}

// Session reports the logind seat session.
func Session() (SessionInfo, error) {
	e, err := section(SectionSession)
	if err != nil {
		return e.Session, fmt.Errorf("linuxdesktop: detect session: %w", err)
	}
	return e.Session, nil
}

// Display reports which display servers are reachable, connecting to each one
// that is.
func Display() (DisplayInfo, error) {
	e, err := section(SectionDisplay)
	if err != nil {
		return e.Display, fmt.Errorf("linuxdesktop: detect display: %w", err)
	}
	return e.Display, nil
}

// Desktop reports the desktop environment.
func Desktop() (DesktopInfo, error) {
	e, err := section(SectionDesktop)
	if err != nil {
		return e.Desktop, fmt.Errorf("linuxdesktop: detect desktop: %w", err)
	}
	return e.Desktop, nil
}

// Compositor reports the compositor or window manager, together with the
// method that identified it and how far that method can be trusted.
func Compositor() (CompositorInfo, error) {
	e, err := section(SectionCompositor)
	if err != nil {
		return e.Compositor, fmt.Errorf("linuxdesktop: detect compositor: %w", err)
	}
	return e.Compositor, nil
}

// Graphics reports the GPUs and the client-side graphics stack. It uses the
// filesystem-only stack probe; for a live renderer string call DetectContext
// with WithOpenGL.
func Graphics() (GraphicsInfo, error) {
	e, err := section(SectionGraphics)
	if err != nil {
		return e.Graphics, fmt.Errorf("linuxdesktop: detect graphics: %w", err)
	}
	return e.Graphics, nil
}

// Portal reports whether xdg-desktop-portal is running and what it offers.
func Portal() (PortalInfo, error) {
	e, err := section(SectionPortal)
	if err != nil {
		return e.Portal, fmt.Errorf("linuxdesktop: detect portal: %w", err)
	}
	return e.Portal, nil
}

// IsWayland reports whether a Wayland compositor socket is present. It stats a
// path and nothing more, so it is safe on any hot path; Display connects and
// is therefore the authority.
func IsWayland() bool { return onLinux(runtime.GOOS, waylandAvailable) }

func waylandAvailable() bool {
	snapshot := env.New().Snapshot()

	return wayland.SocketPath(&snapshot) != ""
}

// IsX11 reports whether an X server appears to be reachable. Like IsWayland it
// only stats; a display on a remote host is assumed reachable.
func IsX11() bool { return onLinux(runtime.GOOS, x11Available) }

func x11Available() bool {
	snapshot := env.New().Snapshot()

	return x11.Available(&snapshot)
}

// IsHeadless reports whether there is no display server to draw on.
func IsHeadless() bool {
	return !IsWayland() && !IsX11()
}

// WithTimeout bounds the whole detection run. The default is DefaultTimeout.
func WithTimeout(d time.Duration) Option { return func(c *Config) { c.Timeout = d } }

// WithProbeTimeout bounds each individual probe, so one unresponsive server
// cannot consume the whole budget. The default is DefaultProbeTimeout.
func WithProbeTimeout(d time.Duration) Option { return func(c *Config) { c.ProbeTimeout = d } }

// WithSections selects which parts of the Environment to populate. Probes for
// unselected sections never run.
func WithSections(s Section) Option { return func(c *Config) { c.Sections = s } }

// WithOpenGL asks the graphics drivers directly instead of reading their
// manifests, which yields a true vendor, renderer and version at the cost of
// initializing the GPU stack — tens of milliseconds, not microseconds.
func WithOpenGL() Option { return func(c *Config) { c.NativeGL = true } }

// WithProcessScan enables the /proc fallback for compositor detection. It only
// affects the outcome when every stronger signal has failed, and a compositor
// found this way is reported with ConfidenceLow.
func WithProcessScan() Option { return func(c *Config) { c.ProcessScan = true } }

// section runs a detection limited to one part of the Environment.
func section(s Section) (*Environment, error) {
	result, callErr := DetectContext(context.Background(), WithSections(s))

	if callErr != nil {
		return result, fmt.Errorf("linuxdesktop: detect section: %w", callErr)
	}
	return result, nil
}

// adapters is the composition root: the one place that decides which
// implementation satisfies each port.
func adapters(bus port.Bus, cfg *core.Config) core.Deps {
	deps := core.Deps{
		Env:     env.New(),
		OS:      osinfo.New(),
		Session: logind.New(bus),
		X11:     x11.New(),
		Wayland: wayland.New(),
		Desktop: desktop.New(bus),
		GPUs:    drm.New(),
		Stack:   graphicsProbe(cfg),
		Portal:  portal.New(bus),
		Process: nil,
	}
	if cfg.ProcessScan {
		deps.Process = procscan.New(rules.IsCompositorProcess)
	}

	return deps
}

func graphicsProbe(cfg *core.Config) gl.Probe {
	if cfg.NativeGL {
		return gl.NewNative()
	}

	return gl.New()
}

func onLinux(goos string, available func() bool) bool {
	return goos == goosLinux && available()
}
