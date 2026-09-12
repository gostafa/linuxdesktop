package core

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/gostafa/linuxdesktop/internal/domain"
	"github.com/gostafa/linuxdesktop/internal/rules"
)

// NewConfig applies opts to the defaults. It is separate from New because the
// caller needs the effective configuration in order to decide which adapters
// to construct in the first place.
func NewConfig(opts ...Option) Config {
	cfg := DefaultConfig
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	if cfg.Sections == 0 {
		cfg.Sections = domain.SectionAll
	}
	return cfg
}

// New returns an Engine wired to deps.
func New(deps Deps, cfg Config) *Engine {
	if cfg.Sections == 0 {
		cfg.Sections = domain.SectionAll
	}
	return &Engine{deps: deps, cfg: cfg}
}

// WithTimeout bounds the whole detection run.
func WithTimeout(d time.Duration) Option {
	return func(c *Config) { c.Timeout = d }
}

// WithProbeTimeout bounds each individual probe.
func WithProbeTimeout(d time.Duration) Option {
	return func(c *Config) { c.ProbeTimeout = d }
}

// WithSections selects which parts of the Environment to populate.
func WithSections(s domain.Section) Option {
	return func(c *Config) { c.Sections = s }
}

// WithNativeGL enables the dlopen graphics probe.
func WithNativeGL() Option {
	return func(c *Config) { c.NativeGL = true }
}

// WithProcessScan enables the /proc compositor fallback.
func WithProcessScan() Option {
	return func(c *Config) { c.ProcessScan = true }
}

// Detect runs the probes and merges their answers. The returned Environment is
// never nil; the error reports non-fatal probe failures and can be ignored by
// callers that only want the data.
func (e *Engine) Detect(ctx context.Context) (*domain.Environment, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if e.cfg.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, e.cfg.Timeout)
		defer cancel()
	}

	out := &domain.Environment{}
	sig := domain.Signals{}
	if e.deps.Env != nil {
		sig.Env = e.deps.Env.Snapshot()
	}

	var (
		mu      sync.Mutex
		failed  []error
		wg      sync.WaitGroup
		x11Info *domain.X11Info
		wlInfo  *domain.WaylandInfo
	)
	record := func(err error) {
		if err == nil {
			return
		}
		mu.Lock()
		failed = append(failed, err)
		mu.Unlock()
	}
	spawn := func(fn func(context.Context) error) {
		wg.Go(func() {
			probeCtx, cancel := e.probeContext(ctx)
			defer cancel()
			record(fn(probeCtx))
		})
	}

	// The display probes feed compositor detection, so they run whenever
	// either section is wanted.
	wantDisplay := e.wants(domain.SectionDisplay) || e.wants(domain.SectionCompositor)

	if e.wants(domain.SectionOS) && e.deps.OS != nil {
		spawn(func(ctx context.Context) error {
			info, err := e.deps.OS.OS(ctx)
			mu.Lock()
			out.OS = info
			mu.Unlock()
			return err
		})
	}

	if e.wants(domain.SectionSession) && e.deps.Session != nil {
		spawn(func(ctx context.Context) error {
			info, err := e.deps.Session.Session(ctx, &sig.Env)
			mu.Lock()
			out.Session = info
			mu.Unlock()
			return err
		})
	}

	if wantDisplay && e.deps.X11 != nil {
		spawn(func(ctx context.Context) error {
			info, err := e.deps.X11.X11(ctx, &sig.Env)
			mu.Lock()
			x11Info = info
			mu.Unlock()
			return err
		})
	}

	if wantDisplay && e.deps.Wayland != nil {
		spawn(func(ctx context.Context) error {
			info, err := e.deps.Wayland.Wayland(ctx, &sig.Env)
			mu.Lock()
			wlInfo = info
			mu.Unlock()
			return err
		})
	}

	if (e.wants(domain.SectionDesktop) || e.wants(domain.SectionCompositor)) && e.deps.Desktop != nil {
		spawn(func(ctx context.Context) error {
			info, err := e.deps.Desktop.Desktop(ctx, &sig.Env)
			mu.Lock()
			out.Desktop = info
			mu.Unlock()
			return err
		})
	}

	if e.wants(domain.SectionGraphics) && e.deps.GPUs != nil {
		spawn(func(ctx context.Context) error {
			gpus, primary, err := e.deps.GPUs.GPUs(ctx)
			mu.Lock()
			out.Graphics.GPUs = gpus
			out.Graphics.PrimaryGPU = primary
			mu.Unlock()
			return err
		})
	}

	if e.wants(domain.SectionGraphics) && e.deps.Stack != nil {
		spawn(func(ctx context.Context) error {
			opengl, vulkan, err := e.deps.Stack.Stack(ctx)
			mu.Lock()
			out.Graphics.OpenGL = opengl
			out.Graphics.Vulkan = vulkan
			mu.Unlock()
			return err
		})
	}

	if e.wants(domain.SectionPortal) && e.deps.Portal != nil {
		spawn(func(ctx context.Context) error {
			info, err := e.deps.Portal.Portal(ctx, &sig.Env)
			mu.Lock()
			out.Portal = info
			mu.Unlock()
			return err
		})
	}

	if e.cfg.ProcessScan && e.wants(domain.SectionCompositor) && e.deps.Process != nil {
		spawn(func(ctx context.Context) error {
			names, err := e.deps.Process.Processes(ctx)
			mu.Lock()
			sig.Processes = names
			mu.Unlock()
			return err
		})
	}

	wg.Wait()

	// Stage two: everything below is pure, with every fact already collected.
	display := buildDisplay(&sig.Env, x11Info, wlInfo)
	sig.Desktop = out.Desktop.Environment
	sig.X11Reachable = display.X11Available
	sig.WaylandReachable = display.WaylandAvailable
	if x11Info != nil {
		sig.X11WindowManager = x11Info.WindowManager
		sig.X11Extensions = x11Info.Extensions
	}
	if wlInfo != nil {
		sig.WaylandGlobals = wlInfo.Globals
	}

	display.Protocol = rules.Protocol(out.Session.Type, display.WaylandAvailable, display.X11Available)
	if wantDisplay {
		// Headless is only meaningful once the display probes have run;
		// asserting it from a run that skipped them would be a lie.
		out.Headless = rules.Headless(display.WaylandAvailable, display.X11Available)
	}
	if e.wants(domain.SectionDisplay) {
		out.Display = display
	}
	if e.wants(domain.SectionCompositor) {
		out.Compositor = rules.Compositor(&sig)
	}

	return out, errors.Join(failed...)
}

func (e *Engine) wants(s domain.Section) bool { return e.cfg.Sections&s != 0 }

// probeContext derives the per-probe deadline. The returned cancel must always
// be called, which is why probes are spawned through a helper rather than
// inline.
func (e *Engine) probeContext(ctx context.Context) (context.Context, context.CancelFunc) {
	if e.cfg.ProbeTimeout <= 0 {
		return context.WithCancel(ctx)
	}
	return context.WithTimeout(ctx, e.cfg.ProbeTimeout)
}

// buildDisplay assembles DisplayInfo from the environment and whatever the two
// protocol probes managed to reach.
func buildDisplay(env *domain.Env, x11Info *domain.X11Info, wlInfo *domain.WaylandInfo) domain.DisplayInfo {
	info := domain.DisplayInfo{
		Protocol:         domain.DisplayProtocolUnknown,
		WaylandDisplay:   env.WaylandDisplay,
		X11Display:       env.Display,
		X11Available:     x11Info != nil,
		WaylandAvailable: wlInfo != nil,
		X11:              x11Info,
		Wayland:          wlInfo,
	}
	// An Xwayland server says so itself; otherwise a reachable X server inside
	// a Wayland session can only be Xwayland.
	info.XWayland = hasExtension(x11Info, domain.ExtensionXWayland) ||
		(info.X11Available && info.WaylandAvailable)
	return info
}

func hasExtension(info *domain.X11Info, name string) bool {
	if info == nil {
		return false
	}
	for _, e := range info.Extensions {
		if e == name {
			return true
		}
	}
	return false
}
