// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package schema

import (
	"time"
)

type (
	// RetryPolicy controls the bounded initialization of a D-Bus connection.
	RetryPolicy struct {
		// MaxAttempts counts the first attempt; one disables retries.
		MaxAttempts uint
		// InitialInterval is the initial delay between failed attempts.
		InitialInterval time.Duration
		// MaxInterval caps the nominal delay before jitter.
		MaxInterval time.Duration
	}

	// Environment is the shared Environment layout, parameterized by its component types.
	Environment[
		OSType,
		SessionType,
		DisplayType,
		DesktopType,
		CompositorType,
		GraphicsType,
		PortalType any,
	] = struct {
		// OS operating system and kernel identity.
		OS OSType `json:"os"`
		// Session seat and login session details.
		Session SessionType `json:"session"`
		// Display reachable display servers and selected protocol.
		Display DisplayType `json:"display"`
		// Desktop desktop environment identity.
		Desktop DesktopType `json:"desktop"`
		// Compositor window manager or compositor identity and detection evidence.
		Compositor CompositorType `json:"compositor"`
		// Graphics gPU hardware and client graphics implementations.
		Graphics GraphicsType `json:"graphics"`
		// Portal availability and interfaces of the desktop portal.
		Portal PortalType `json:"portal"`
		// Headless whether neither display protocol is reachable.
		Headless bool `json:"headless"`
	}

	// SessionInfo is the shared SessionInfo layout, parameterized by its component types.
	SessionInfo[SessionType any] = struct {
		// ID identifier reported by the underlying system.
		ID string `json:"id"`
		// Seat logind seat identifier, such as seat0.
		Seat string `json:"seat"`
		// Type windowing protocol associated with the session.
		Type SessionType `json:"type"`
		// Desktop desktop identifier reported by logind or the environment.
		Desktop string `json:"desktop"`
		// Name human-readable name reported by the source.
		Name string `json:"name"`
		// User user associated with the login session.
		User string `json:"user"`
		// XDGSessionType unmodified XDG_SESSION_TYPE value.
		XDGSessionType string `json:"xdg_session_type"`
		// XDGSessionDesktop unmodified XDG_SESSION_DESKTOP value.
		XDGSessionDesktop string `json:"xdg_session_desktop"`
		// VTNumber virtual terminal number; zero when unavailable.
		VTNumber int `json:"vt_number"`
		// Remote whether the session is remote, including SSH sessions.
		Remote bool `json:"remote"`
		// Active whether the session is currently active.
		Active bool `json:"active"`
	}

	// DisplayInfo is the shared DisplayInfo layout, parameterized by its component types.
	DisplayInfo[XServerType, WaylandServerType, ProtocolType any] = struct {
		// X11 details of a reachable X server; nil when unavailable.
		X11 *XServerType `json:"x11,omitempty"`
		// Wayland details of a reachable Wayland compositor; nil when unavailable.
		Wayland *WaylandServerType `json:"wayland,omitempty"`
		// Protocol preferred display protocol for a client.
		Protocol ProtocolType `json:"protocol"`
		// WaylandDisplay wayland socket name from WAYLAND_DISPLAY.
		WaylandDisplay string `json:"wayland_display"`
		// X11Display x display address from DISPLAY.
		X11Display string `json:"x11_display"`
		// X11Available whether the X server probe reached a server.
		X11Available bool `json:"x11_available"`
		// WaylandAvailable whether the Wayland probe reached a compositor.
		WaylandAvailable bool `json:"wayland_available"`
		// XWayland whether the X server runs inside a Wayland session.
		XWayland bool `json:"xwayland"`
	}

	// WaylandInfo is the shared WaylandInfo layout, parameterized by its component types.
	WaylandInfo[GlobalType any] = struct {
		// Display wayland socket name used for the connection.
		Display string `json:"display"`
		// Globals global interfaces advertised by the Wayland registry.
		Globals []GlobalType `json:"globals,omitempty"`
		// SocketFD inherited WAYLAND_SOCKET descriptor, when present.
		SocketFD int `json:"socket_fd,omitempty"`
	}

	// Config is the shared Config layout, parameterized by its component types.
	Config[SectionsType any] = struct {
		Sections        SectionsType
		ConnectionRetry RetryPolicy
		Timeout         time.Duration
		ProbeTimeout    time.Duration
		NativeGL        bool
		ProcessScan     bool
	}

	// DesktopInfo is the shared DesktopInfo layout, parameterized by its component types.
	DesktopInfo[DesktopType any] = struct {
		// Environment recognized desktop environment.
		Environment DesktopType `json:"environment"`
		// Name human-readable name reported by the source.
		Name string `json:"name"`
		// Version version reported by the source, when available.
		Version string `json:"version,omitempty"`
		// CurrentDesktop unmodified XDG_CURRENT_DESKTOP value.
		CurrentDesktop string `json:"current_desktop"`
		// SessionDesktop unmodified XDG_SESSION_DESKTOP value.
		SessionDesktop string `json:"session_desktop"`
		// DesktopSession unmodified DESKTOP_SESSION value.
		DesktopSession string `json:"desktop_session"`
		// CurrentDesktops ordered desktop tokens from XDG_CURRENT_DESKTOP.
		CurrentDesktops []string `json:"current_desktops"`
	}

	// CompositorInfo is the shared CompositorInfo layout, parameterized by its component types.
	CompositorInfo[KindType, ConfidenceType, MethodType any] = struct {
		// Kind identifies the compositor or window manager.
		Kind KindType `json:"kind"`
		// Confidence expresses the certainty of the identification.
		Confidence ConfidenceType `json:"confidence"`
		// DetectedBy identifies the source of detection evidence.
		DetectedBy MethodType `json:"detected_by"`
		// Name is the human-readable compositor name.
		Name string `json:"name"`
		// Version is the reported compositor version, when available.
		Version string `json:"version,omitempty"`
		// Wayland reports whether the compositor supports Wayland.
		Wayland bool `json:"wayland"`
		// X11 reports whether the window manager supports X11.
		X11 bool `json:"x11"`
	}

	// GraphicsInfo is the shared GraphicsInfo layout, parameterized by its component types.
	GraphicsInfo[OpenGLType, VulkanType, GPUType any] = struct {
		// OpenGL openGL implementation reported by drivers or manifests.
		OpenGL OpenGLType `json:"opengl"`
		// Vulkan vulkan loader and driver availability.
		Vulkan VulkanType `json:"vulkan"`
		// PrimaryGPU identifier of the firmware-posted GPU, or the first detected GPU.
		PrimaryGPU string `json:"primary_gpu,omitempty"`
		// GPUs detected DRM devices with their PCI identity.
		GPUs []GPUType `json:"gpus"`
	}
)
