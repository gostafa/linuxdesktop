package linuxdesktop

import (
	"github.com/gostafa/linuxdesktop/internal/core"
	"github.com/gostafa/linuxdesktop/internal/domain"
)

// The model. These are aliases rather than distinct definitions so that the
// adapters inside internal/ can speak the same types the public API exposes,
// without an import cycle and without a translation layer.
type (
	// SessionType is the kind of seat session the process is attached to.
	SessionType = domain.SessionType
	// DisplayProtocol is the windowing protocol the process should speak.
	DisplayProtocol = domain.DisplayProtocol
	// DesktopEnvironment is a recognised desktop environment.
	DesktopEnvironment = domain.DesktopEnvironment
	// CompositorKind is a recognised window manager or Wayland compositor.
	CompositorKind = domain.CompositorKind
	// DetectionConfidence grades how much a detection result can be trusted.
	DetectionConfidence = domain.DetectionConfidence
	// DetectionMethod records which signal produced a detection result.
	DetectionMethod = domain.DetectionMethod
	// Section selects which parts of an Environment to populate.
	Section = domain.Section

	// Environment is the complete picture of the desktop.
	Environment = domain.Environment
	// OSInfo describes the operating system and kernel.
	OSInfo = domain.OSInfo
	// SessionInfo describes the logind seat session.
	SessionInfo = domain.SessionInfo
	// DisplayInfo describes which display servers are reachable.
	DisplayInfo = domain.DisplayInfo
	// X11Info is the result of a real connection to an X server.
	X11Info = domain.X11Info
	// WaylandInfo is the result of a real connection to a compositor.
	WaylandInfo = domain.WaylandInfo
	// WaylandGlobal is one entry from the compositor's global registry.
	WaylandGlobal = domain.WaylandGlobal
	// DesktopInfo identifies the desktop environment.
	DesktopInfo = domain.DesktopInfo
	// CompositorInfo identifies the compositor and how certain that is.
	CompositorInfo = domain.CompositorInfo
	// GraphicsInfo describes the GPUs and the client-side graphics stack.
	GraphicsInfo = domain.GraphicsInfo
	// GPUInfo is one DRM device and its PCI identity.
	GPUInfo = domain.GPUInfo
	// OpenGLInfo describes the OpenGL implementation.
	OpenGLInfo = domain.OpenGLInfo
	// VulkanInfo describes the Vulkan loader and its ICDs.
	VulkanInfo = domain.VulkanInfo
	// PortalInfo describes xdg-desktop-portal.
	PortalInfo = domain.PortalInfo

	// Option configures a detection run.
	Option = core.Option
)
