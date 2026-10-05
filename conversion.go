// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package linuxdesktop

import (
	"slices"

	"github.com/gostafa/linuxdesktop/internal/core"
	"github.com/gostafa/linuxdesktop/internal/domain"
)

// newConfig applies public options before crossing the internal boundary.
func newConfig(opts ...Option) core.Config {
	cfg := Config{
		Timeout:      DefaultTimeout,
		ProbeTimeout: DefaultProbeTimeout,
		Sections:     SectionAll,
	}

	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}

	if cfg.Sections == 0 {
		cfg.Sections = SectionAll
	}

	return core.Config{
		Timeout:      cfg.Timeout,
		ProbeTimeout: cfg.ProbeTimeout,
		Sections:     domain.Section(cfg.Sections),
		NativeGL:     cfg.NativeGL,
		ProcessScan:  cfg.ProcessScan,
	}
}

// publicEnvironment translates the internal model without sharing mutable data.
func publicEnvironment(in *domain.Environment) *Environment {
	return &Environment{
		OS: OSInfo(in.OS),
		Session: SessionInfo{
			ID:                in.Session.ID,
			Seat:              in.Session.Seat,
			Type:              SessionType(in.Session.Type),
			Desktop:           in.Session.Desktop,
			Name:              in.Session.Name,
			User:              in.Session.User,
			XDGSessionType:    in.Session.XDGSessionType,
			XDGSessionDesktop: in.Session.XDGSessionDesktop,
			VTNumber:          in.Session.VTNumber,
			Remote:            in.Session.Remote,
			Active:            in.Session.Active,
		},
		Display: publicDisplay(in.Display),
		Desktop: DesktopInfo{
			Environment:     DesktopEnvironment(in.Desktop.Environment),
			Name:            in.Desktop.Name,
			Version:         in.Desktop.Version,
			CurrentDesktop:  in.Desktop.CurrentDesktop,
			SessionDesktop:  in.Desktop.SessionDesktop,
			DesktopSession:  in.Desktop.DesktopSession,
			CurrentDesktops: slices.Clone(in.Desktop.CurrentDesktops),
		},
		Compositor: CompositorInfo{
			Kind:       CompositorKind(in.Compositor.Kind),
			Name:       in.Compositor.Name,
			Version:    in.Compositor.Version,
			Confidence: DetectionConfidence(in.Compositor.Confidence),
			DetectedBy: DetectionMethod(in.Compositor.DetectedBy),
			Wayland:    in.Compositor.Wayland,
			X11:        in.Compositor.X11,
		},
		Graphics: GraphicsInfo{
			OpenGL:     OpenGLInfo(in.Graphics.OpenGL),
			Vulkan:     VulkanInfo(in.Graphics.Vulkan),
			PrimaryGPU: in.Graphics.PrimaryGPU,
			GPUs:       publicGPUs(in.Graphics.GPUs),
		},
		Portal:   PortalInfo(in.Portal),
		Headless: in.Headless,
	}
}

func publicDisplay(in domain.DisplayInfo) DisplayInfo {
	out := DisplayInfo{
		Protocol:         DisplayProtocol(in.Protocol),
		WaylandDisplay:   in.WaylandDisplay,
		X11Display:       in.X11Display,
		X11Available:     in.X11Available,
		WaylandAvailable: in.WaylandAvailable,
		XWayland:         in.XWayland,
	}
	if in.X11 != nil {
		info := X11Info(*in.X11)
		info.Extensions = slices.Clone(in.X11.Extensions)
		out.X11 = &info
	}
	if in.Wayland != nil {
		out.Wayland = &WaylandInfo{
			Display:  in.Wayland.Display,
			SocketFD: in.Wayland.SocketFD,
		}
		if in.Wayland.Globals != nil {
			out.Wayland.Globals = make([]WaylandGlobal, len(in.Wayland.Globals))
			for i, global := range in.Wayland.Globals {
				out.Wayland.Globals[i] = WaylandGlobal(global)
			}
		}
	}
	return out
}

func publicGPUs(in []domain.GPUInfo) []GPUInfo {
	if in == nil {
		return nil
	}
	out := make([]GPUInfo, len(in))
	for i, gpu := range in {
		out[i] = GPUInfo(gpu)
	}
	return out
}
