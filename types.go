// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package linuxdesktop

import "time"

type (
	// Section selects which parts of an Environment to populate. Sections are
	// a bitmask so a caller can pay only for what it reads.
	Section uint16

	// Environment is the complete picture of the desktop the process is
	// running in.
	Environment struct {
		OS         OSInfo         `json:"os"`
		Session    SessionInfo    `json:"session"`
		Display    DisplayInfo    `json:"display"`
		Desktop    DesktopInfo    `json:"desktop"`
		Compositor CompositorInfo `json:"compositor"`
		Graphics   GraphicsInfo   `json:"graphics"`
		Portal     PortalInfo     `json:"portal"`
		Headless   bool           `json:"headless"`
	}
)

type (
	// SessionType is the kind of seat session the process is attached to.
	SessionType string

	// SessionInfo describes the logind seat session.
	SessionInfo struct {
		ID                string      `json:"id"`
		Seat              string      `json:"seat"`
		Type              SessionType `json:"type"`
		Desktop           string      `json:"desktop"`
		Name              string      `json:"name"`
		User              string      `json:"user"`
		XDGSessionType    string      `json:"xdg_session_type"`
		XDGSessionDesktop string      `json:"xdg_session_desktop"`
		VTNumber          int         `json:"vt_number"`
		Remote            bool        `json:"remote"`
		Active            bool        `json:"active"`
	}
)

type (
	// DisplayProtocol is the windowing protocol the process should speak.
	DisplayProtocol string

	// DisplayInfo describes which display servers are reachable.
	DisplayInfo struct {
		X11              *X11Info        `json:"x11,omitempty"`
		Wayland          *WaylandInfo    `json:"wayland,omitempty"`
		Protocol         DisplayProtocol `json:"protocol"`
		WaylandDisplay   string          `json:"wayland_display"`
		X11Display       string          `json:"x11_display"`
		X11Available     bool            `json:"x11_available"`
		WaylandAvailable bool            `json:"wayland_available"`
		XWayland         bool            `json:"xwayland"`
	}

	// X11Info is the result of a real connection to an X server.
	X11Info struct {
		Display       string   `json:"display"`
		Vendor        string   `json:"vendor"`
		WindowManager string   `json:"window_manager"`
		Extensions    []string `json:"extensions,omitempty"`
		Screen        int      `json:"screen"`
		ProtocolMajor int      `json:"protocol_major"`
		ProtocolMinor int      `json:"protocol_minor"`
	}

	// WaylandInfo is the result of a real connection to a Wayland compositor.
	WaylandInfo struct {
		Display  string          `json:"display"`
		Globals  []WaylandGlobal `json:"globals,omitempty"`
		SocketFD int             `json:"socket_fd,omitempty"`
	}

	// WaylandGlobal is one entry from the compositor's global registry.
	WaylandGlobal struct {
		Interface string `json:"interface"`
		Name      uint32 `json:"name"`
		Version   uint32 `json:"version"`
	}
)

type (
	// DesktopEnvironment is a recognized desktop environment.
	DesktopEnvironment string

	// DesktopInfo identifies the desktop environment.
	DesktopInfo struct {
		Environment     DesktopEnvironment `json:"environment"`
		Name            string             `json:"name"`
		Version         string             `json:"version,omitempty"`
		CurrentDesktop  string             `json:"current_desktop"`
		SessionDesktop  string             `json:"session_desktop"`
		DesktopSession  string             `json:"desktop_session"`
		CurrentDesktops []string           `json:"current_desktops"`
	}
)

type (
	// CompositorKind is a recognized window manager or Wayland compositor.
	CompositorKind string

	// DetectionConfidence grades how much a detection result can be trusted.
	DetectionConfidence string

	// DetectionMethod records which signal produced a detection result.
	DetectionMethod string

	// CompositorInfo identifies the compositor or window manager, and says how
	// certain that identification is.
	CompositorInfo struct {
		Kind       CompositorKind      `json:"kind"`
		Name       string              `json:"name"`
		Version    string              `json:"version,omitempty"`
		Confidence DetectionConfidence `json:"confidence"`
		DetectedBy DetectionMethod     `json:"detected_by"`
		Wayland    bool                `json:"wayland"`
		X11        bool                `json:"x11"`
	}
)

type (
	// GraphicsInfo describes the GPUs and the client-side graphics stack.
	GraphicsInfo struct {
		OpenGL     OpenGLInfo `json:"opengl"`
		Vulkan     VulkanInfo `json:"vulkan"`
		PrimaryGPU string     `json:"primary_gpu,omitempty"`
		GPUs       []GPUInfo  `json:"gpus"`
	}

	// GPUInfo is one DRM device and its PCI identity.
	GPUInfo struct {
		ID           string `json:"id"`
		Vendor       string `json:"vendor"`
		Model        string `json:"model"`
		VendorID     string `json:"vendor_id"`
		DeviceID     string `json:"device_id"`
		PCIAddress   string `json:"pci_address"`
		Driver       string `json:"driver"`
		DRMDevice    string `json:"drm_device"`
		RenderDevice string `json:"render_device"`
		Integrated   bool   `json:"integrated"`
		Discrete     bool   `json:"discrete"`
	}

	// OpenGLInfo describes the OpenGL implementation.
	OpenGLInfo struct {
		Vendor    string `json:"vendor,omitempty"`
		Renderer  string `json:"renderer,omitempty"`
		Version   string `json:"version,omitempty"`
		Available bool   `json:"available"`
	}

	// VulkanInfo describes the Vulkan loader and its ICDs.
	VulkanInfo struct {
		Version   string `json:"version,omitempty"`
		Available bool   `json:"available"`
	}
)

type (
	// OSInfo describes the operating system and kernel.
	OSInfo struct {
		Name          string `json:"name"`
		PrettyName    string `json:"pretty_name"`
		ID            string `json:"id"`
		IDLike        string `json:"id_like"`
		Version       string `json:"version"`
		VersionID     string `json:"version_id"`
		Kernel        string `json:"kernel"`
		KernelRelease string `json:"kernel_release"`
		Architecture  string `json:"architecture"`
		Hostname      string `json:"hostname"`
	}
)

type (
	// PortalInfo describes xdg-desktop-portal and which of its interfaces the
	// running backend actually exports.
	PortalInfo struct {
		Backend       string `json:"backend,omitempty"`
		Available     bool   `json:"available"`
		DesktopPortal bool   `json:"desktop_portal"`
		ScreenCast    bool   `json:"screencast"`
		Screenshot    bool   `json:"screenshot"`
		FileChooser   bool   `json:"file_chooser"`
		OpenURI       bool   `json:"open_uri"`
		RemoteDesktop bool   `json:"remote_desktop"`
		Inhibit       bool   `json:"inhibit"`
		Notification  bool   `json:"notification"`
	}
)

// Config controls a detection run. DetectContext initializes the defaults before
// applying options in order; callers can supply custom Option functions.
type Config struct {
	// Timeout bounds the whole run. Non-positive values leave it unbounded.
	Timeout time.Duration
	// ProbeTimeout bounds each probe. Non-positive values use the run context.
	ProbeTimeout time.Duration
	// Sections selects the result sections. Zero selects SectionAll.
	Sections Section
	// NativeGL enables direct querying of the graphics drivers.
	NativeGL bool
	// ProcessScan enables the /proc fallback for compositor detection.
	ProcessScan bool
}

// Option configures a detection run by modifying its public configuration.
type Option func(*Config)
