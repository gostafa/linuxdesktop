// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package core

import (
	"context"
	"errors"
	"time"

	"github.com/gostafa/linuxdesktop/internal/domain"
	"github.com/gostafa/linuxdesktop/internal/rules"
)

// NewConfig applies opts to the defaults. It is separate from New because the
// caller needs the effective configuration in order to decide which adapters
// to construct in the first place.
func NewConfig(opts ...Option) Config {
	cfg := DefaultConfig
	for i := range opts {
		if opts[i] != nil {
			opts[i](&cfg)
		}
	}

	return withDefaults(cfg)
}

// New returns an Engine wired to deps.
func New(deps *Deps, cfg Config) *Engine {
	return &Engine{deps: *deps, cfg: withDefaults(cfg)}
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
func (eng *Engine) Detect(ctx context.Context) (*domain.Environment, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if eng.cfg.Timeout > zero {
		var cancel context.CancelFunc

		ctx, cancel = context.WithTimeout(ctx, eng.cfg.Timeout)

		defer cancel()
	}

	gather := newCollector(eng.snapshot())

	eng.run(ctx, gather)
	eng.finish(gather)

	return gather.out, errors.Join(gather.failed...)
}

// finish is stage two: everything below is pure, with every fact collected.
func (eng *Engine) finish(gather *collector) {
	display := gather.display()

	gather.summarize(&display)

	eng.publish(gather, &display)
}

// probeContext derives the per-probe deadline. The returned cancel must always
// be called, which is why probes are spawned through a helper rather than
// inline.
func (eng *Engine) probeContext(ctx context.Context) (context.Context, context.CancelFunc) {
	if eng.cfg.ProbeTimeout <= zero {
		return context.WithCancel(ctx)
	}

	return context.WithTimeout(ctx, eng.cfg.ProbeTimeout)
}

func (eng *Engine) probeDesktop(ctx context.Context, gather *collector) error {
	info, err := eng.deps.Desktop.Desktop(ctx, &gather.sig.Env)

	gather.under(func() { gather.out.Desktop = info })

	return err
}

func (eng *Engine) probeGPUs(ctx context.Context, gather *collector) error {
	gpus, primary, err := eng.deps.GPUs.GPUs(ctx)

	gather.under(func() {
		gather.out.Graphics.GPUs = gpus
		gather.out.Graphics.PrimaryGPU = primary
	})

	return err
}

func (eng *Engine) probeOS(ctx context.Context, gather *collector) error {
	info, err := eng.deps.OS.OS(ctx)

	gather.under(func() { gather.out.OS = info })

	return err
}

func (eng *Engine) probePortal(ctx context.Context, gather *collector) error {
	info, err := eng.deps.Portal.Portal(ctx, &gather.sig.Env)

	gather.under(func() { gather.out.Portal = info })

	return err
}

func (eng *Engine) probeProcesses(ctx context.Context, gather *collector) error {
	names, err := eng.deps.Process.Processes(ctx)

	gather.under(func() { gather.sig.Processes = names })

	return err
}

func (eng *Engine) probeSession(ctx context.Context, gather *collector) error {
	info, err := eng.deps.Session.Session(ctx, &gather.sig.Env)

	gather.under(func() { gather.out.Session = info })

	return err
}

func (eng *Engine) probeStack(ctx context.Context, gather *collector) error {
	opengl, vulkan, err := eng.deps.Stack.Stack(ctx)

	gather.under(func() {
		gather.out.Graphics.OpenGL = opengl
		gather.out.Graphics.Vulkan = vulkan
	})

	return err
}

func (eng *Engine) probeWayland(ctx context.Context, gather *collector) error {
	info, err := eng.deps.Wayland.Wayland(ctx, &gather.sig.Env)

	gather.under(func() { gather.wayland = info })

	return err
}

func (eng *Engine) probeX11(ctx context.Context, gather *collector) error {
	info, err := eng.deps.X11.X11(ctx, &gather.sig.Env)

	gather.under(func() { gather.x11 = info })

	return err
}

// processSection is when the /proc fallback runs, which is never unless the
// caller asked for it: an unselected section is one no probe fires for.
func (eng *Engine) processSection() domain.Section {
	if !eng.cfg.ProcessScan {
		return zero
	}

	return domain.SectionCompositor
}

// publish writes the finished sections the caller asked for.
func (eng *Engine) publish(gather *collector, display *domain.DisplayInfo) {
	if eng.wants(sectionDisplayProbes) {
		// Headless is only meaningful once the display probes have run;
		// asserting it from a run that skipped them would be a lie.
		gather.out.Headless = rules.Headless(display.WaylandAvailable, display.X11Available)
	}

	if eng.wants(domain.SectionDisplay) {
		gather.out.Display = *display
	}

	if eng.wants(domain.SectionCompositor) {
		gather.out.Compositor = rules.Compositor(gather.sig)
	}
}

// run starts every wired probe whose section is wanted, then waits for all of
// them.
func (eng *Engine) run(ctx context.Context, gather *collector) {
	steps := eng.steps()
	for i := range steps {
		if steps[i].wired && eng.wants(steps[i].section) {
			eng.spawn(ctx, gather, steps[i].run)
		}
	}

	gather.wait.Wait()
}

// snapshot reads the environment once, before any probe runs, so that no probe
// pays for a repeated lookup.
func (eng *Engine) snapshot() *domain.Signals {
	sig := &domain.Signals{}
	if eng.deps.Env != nil {
		sig.Env = eng.deps.Env.Snapshot()
	}

	return sig
}

// spawn runs one probe on its own goroutine under its own deadline, so one
// wedged server cannot consume the whole budget.
func (eng *Engine) spawn(ctx context.Context, gather *collector, run probeFunc) {
	gather.wait.Go(func() {
		probeCtx, cancel := eng.probeContext(ctx)
		defer cancel()

		gather.record(run(probeCtx, gather))
	})
}

// steps is every probe the engine knows how to run, with the sections it
// contributes to and whether its adapter was supplied.
func (eng *Engine) steps() []step {
	deps := &eng.deps

	return []step{
		{eng.probeOS, domain.SectionOS, deps.OS != nil},
		{eng.probeSession, domain.SectionSession, deps.Session != nil},
		{eng.probeX11, sectionDisplayProbes, deps.X11 != nil},
		{eng.probeWayland, sectionDisplayProbes, deps.Wayland != nil},
		{eng.probeDesktop, sectionDesktopProbes, deps.Desktop != nil},
		{eng.probeGPUs, domain.SectionGraphics, deps.GPUs != nil},
		{eng.probeStack, domain.SectionGraphics, deps.Stack != nil},
		{eng.probePortal, domain.SectionPortal, deps.Portal != nil},
		{eng.probeProcesses, eng.processSection(), deps.Process != nil},
	}
}

func (eng *Engine) wants(s domain.Section) bool { return eng.cfg.Sections&s != zero }

func newCollector(sig *domain.Signals) *collector {
	return &collector{out: &domain.Environment{}, sig: sig}
}

// display assembles DisplayInfo from the environment and whatever the two
// protocol probes managed to reach.
func (gather *collector) display() domain.DisplayInfo {
	env := &gather.sig.Env
	info := domain.DisplayInfo{
		Protocol:         domain.DisplayProtocolUnknown,
		WaylandDisplay:   env.WaylandDisplay,
		X11Display:       env.Display,
		X11Available:     gather.x11 != nil,
		WaylandAvailable: gather.wayland != nil,
		X11:              gather.x11,
		Wayland:          gather.wayland,
	}
	// An Xwayland server says so itself; otherwise a reachable X server inside
	// a Wayland session can only be Xwayland.
	info.XWayland = hasExtension(gather.x11, domain.ExtensionXWayland) ||
		(info.X11Available && info.WaylandAvailable)

	return info
}

// record keeps a probe failure. Failures are collected rather than returned,
// since one unreachable server should not hide what the others found.
func (gather *collector) record(err error) {
	if err == nil {
		return
	}

	gather.lock.Lock()
	defer gather.lock.Unlock()

	gather.failed = append(gather.failed, err)
}

// recordProtocols copies what the protocol probes saw into the signals the
// classification rules read.
func (gather *collector) recordProtocols() {
	if gather.x11 != nil {
		gather.sig.X11WindowManager = gather.x11.WindowManager
		gather.sig.X11Extensions = gather.x11.Extensions
	}

	if gather.wayland != nil {
		gather.sig.WaylandGlobals = gather.wayland.Globals
	}
}

// summarize records everything the rules reason over, now that every probe has
// answered, and settles which protocol a client should actually speak.
func (gather *collector) summarize(display *domain.DisplayInfo) {
	sig := gather.sig

	sig.Desktop = gather.out.Desktop.Environment
	sig.X11Reachable = display.X11Available
	sig.WaylandReachable = display.WaylandAvailable

	gather.recordProtocols()

	display.Protocol = rules.Protocol(
		gather.out.Session.Type,
		display.WaylandAvailable,
		display.X11Available,
	)
}

// under runs fn while holding the lock. Every probe writes into the same
// Environment from its own goroutine, so all of them go through here.
func (gather *collector) under(fn func()) {
	gather.lock.Lock()
	defer gather.lock.Unlock()

	fn()
}

// withDefaults fills in a configuration the caller left blank.
func withDefaults(cfg Config) Config {
	if cfg.Sections == zero {
		cfg.Sections = domain.SectionAll
	}

	return cfg
}

func hasExtension(info *domain.X11Info, name string) bool {
	if info == nil {
		return false
	}

	for i := range info.Extensions {
		if info.Extensions[i] == name {
			return true
		}
	}

	return false
}
