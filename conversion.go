// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package linuxdesktop

import (
	"slices"

	"github.com/gostafa/linuxdesktop/internal/core"
	"github.com/gostafa/linuxdesktop/internal/domain"
)

// newConfig applies public options before crossing the internal boundary.
func newConfig(opts ...Option) *core.Config {
	cfg := publicConfig(opts)

	return &core.Config{
		Timeout:      cfg.Timeout,
		ProbeTimeout: cfg.ProbeTimeout,
		Sections:     domain.Section(cfg.Sections),
		NativeGL:     cfg.NativeGL,
		ProcessScan:  cfg.ProcessScan,
	}
}

func publicConfig(opts []Option) Config {
	cfg := Config{
		Timeout:      DefaultTimeout,
		ProbeTimeout: DefaultProbeTimeout,
		Sections:     SectionAll,
		NativeGL:     false,
		ProcessScan:  false,
	}

	for i := range opts {
		if opts[i] != nil {
			opts[i](&cfg)
		}
	}

	if cfg.Sections == noSections {
		cfg.Sections = SectionAll
	}

	return cfg
}

// publicEnvironment translates the internal model without sharing mutable data.
func publicEnvironment(source *domain.Environment) *Environment {
	return &Environment{
		OS:         OSInfo(source.OS),
		Session:    publicSession(&source.Session),
		Display:    publicDisplay(&source.Display),
		Desktop:    publicDesktop(&source.Desktop),
		Compositor: publicCompositor(&source.Compositor),
		Graphics:   publicGraphics(&source.Graphics),
		Portal:     PortalInfo(source.Portal),
		Headless:   source.Headless,
	}
}

func publicSession(source *domain.SessionInfo) SessionInfo {
	return SessionInfo{
		ID:                source.ID,
		Seat:              source.Seat,
		Type:              SessionType(source.Type),
		Desktop:           source.Desktop,
		Name:              source.Name,
		User:              source.User,
		XDGSessionType:    source.XDGSessionType,
		XDGSessionDesktop: source.XDGSessionDesktop,
		VTNumber:          source.VTNumber,
		Remote:            source.Remote,
		Active:            source.Active,
	}
}

func publicDesktop(source *domain.DesktopInfo) DesktopInfo {
	return DesktopInfo{
		Environment:     DesktopEnvironment(source.Environment),
		Name:            source.Name,
		Version:         source.Version,
		CurrentDesktop:  source.CurrentDesktop,
		SessionDesktop:  source.SessionDesktop,
		DesktopSession:  source.DesktopSession,
		CurrentDesktops: slices.Clone(source.CurrentDesktops),
	}
}

func publicCompositor(source *domain.CompositorInfo) CompositorInfo {
	return CompositorInfo{
		Kind:       CompositorKind(source.Kind),
		Name:       source.Name,
		Version:    source.Version,
		Confidence: DetectionConfidence(source.Confidence),
		DetectedBy: DetectionMethod(source.DetectedBy),
		Wayland:    source.Wayland,
		X11:        source.X11,
	}
}

func publicGraphics(source *domain.GraphicsInfo) GraphicsInfo {
	return GraphicsInfo{
		OpenGL:     OpenGLInfo(source.OpenGL),
		Vulkan:     VulkanInfo(source.Vulkan),
		PrimaryGPU: source.PrimaryGPU,
		GPUs:       publicGPUs(source.GPUs),
	}
}

func publicDisplay(source *domain.DisplayInfo) DisplayInfo {
	out := DisplayInfo{
		Protocol:         DisplayProtocol(source.Protocol),
		WaylandDisplay:   source.WaylandDisplay,
		X11Display:       source.X11Display,
		X11Available:     source.X11Available,
		WaylandAvailable: source.WaylandAvailable,
		XWayland:         source.XWayland,
		Wayland:          nil,
		X11:              nil,
	}

	out.X11 = publicX11(source.X11)
	out.Wayland = publicWayland(source.Wayland)

	return out
}

func publicX11(source *domain.X11Info) *X11Info {
	if source == nil {
		return nil
	}

	info := X11Info(*source)

	info.Extensions = slices.Clone(source.Extensions)

	return &info
}

func publicWayland(source *domain.WaylandInfo) *WaylandInfo {
	if source == nil {
		return nil
	}

	out := &WaylandInfo{Display: source.Display, SocketFD: source.SocketFD}
	if source.Globals != nil {
		out.Globals = make([]WaylandGlobal, len(source.Globals))
		for i := range source.Globals {
			out.Globals[i] = WaylandGlobal(source.Globals[i])
		}
	}

	return out
}

func publicGPUs(source []domain.GPUInfo) []GPUInfo {
	if source == nil {
		return nil
	}

	out := make([]GPUInfo, len(source))
	for i := range source {
		out[i] = GPUInfo(source[i])
	}

	return out
}
