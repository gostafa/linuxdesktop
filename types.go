// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package linuxdesktop

import (
	"time"
)

type (
	// Section selects which parts of an Environment to populate. Sections are
	// a bitmask so a caller can pay only for what it reads.
	Section uint16

	// Environment is the complete picture of the desktop the process is
	// running in.
	Environment = EnvironmentData[
		OSInfo,
		SessionInfo,
		DisplayInfo,
		DesktopInfo,
		CompositorInfo,
		GraphicsInfo,
		PortalInfo,
	]

	// Config controls a detection run. DetectContext initializes the defaults before
	// applying options in order; callers can supply custom Option functions.
	Config = ConfigData[Section]

	// Option configures a detection run by modifying its public configuration.
	Option func(*Config)

	// CompositorInfo identifies the compositor or window manager, and says how
	// certain that identification is.
	CompositorInfo = CompositorInfoData[CompositorKind, DetectionConfidence, DetectionMethod]

	// ConfigData is the shared ConfigData layout, parameterized by its component types.
	ConfigData[SectionsType any] struct {
		// Sections selects the environment sections to detect.
		Sections SectionsType
		// Timeout bounds the entire detection run.
		Timeout time.Duration
		// ProbeTimeout bounds each individual probe.
		ProbeTimeout time.Duration
		// NativeGL enables native graphics library probes.
		NativeGL bool
		// ProcessScan enables process-based desktop detection.
		ProcessScan bool
	}

	// EnvironmentData is the shared EnvironmentData layout, parameterized by its component types.
	EnvironmentData[
		OSType,
		SessionType,
		DisplayType,
		DesktopType,
		CompositorType,
		GraphicsType,
		PortalType any,
	] struct {
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

	// SessionInfoData is the shared SessionInfoData layout, parameterized by its component types.
	SessionInfoData[SessionType any] struct {
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

	// DisplayInfoData is the shared DisplayInfoData layout, parameterized by its component types.
	DisplayInfoData[XServerType, WaylandServerType, ProtocolType any] struct {
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

	// WaylandInfoData is the shared WaylandInfoData layout, parameterized by its component types.
	WaylandInfoData[GlobalType any] struct {
		// Display wayland socket name used for the connection.
		Display string `json:"display"`
		// Globals global interfaces advertised by the Wayland registry.
		Globals []GlobalType `json:"globals,omitempty"`
		// SocketFD inherited WAYLAND_SOCKET descriptor, when present.
		SocketFD int `json:"socket_fd,omitempty"`
	}

	// DesktopEnvironment is a recognized desktop environment.
	DesktopEnvironment string

	// DesktopInfo identifies the desktop environment.
	DesktopInfo = DesktopInfoData[DesktopEnvironment]

	// CompositorKind is a recognized window manager or Wayland compositor.
	CompositorKind string

	// DetectionConfidence grades how much a detection result can be trusted.
	DetectionConfidence string

	// DetectionMethod records which signal produced a detection result.
	DetectionMethod string

	// DesktopInfoData is the shared DesktopInfoData layout, parameterized by its component types.
	DesktopInfoData[DesktopType any] struct {
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

	// CompositorInfoData is the shared CompositorInfoData layout, parameterized by its component types.
	CompositorInfoData[KindType, ConfidenceType, MethodType any] struct {
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

	// DisplayProtocol is the windowing protocol the process should speak.
	DisplayProtocol string

	// DisplayInfo describes which display servers are reachable.
	DisplayInfo = DisplayInfoData[X11Info, WaylandInfo, DisplayProtocol]

	// X11Info is the result of a real connection to an X server.
	X11Info struct {
		// Display x display address used for the connection.
		Display string `json:"display"`
		// Vendor vendor name reported by the driver or server.
		Vendor string `json:"vendor"`
		// WindowManager window manager name from EWMH properties.
		WindowManager string `json:"window_manager"`
		// Extensions extension names advertised by the X server.
		Extensions []string `json:"extensions,omitempty"`
		// Screen default X screen index.
		Screen int `json:"screen"`
		// ProtocolMajor major version of the X protocol.
		ProtocolMajor int `json:"protocol_major"`
		// ProtocolMinor minor version of the X protocol.
		ProtocolMinor int `json:"protocol_minor"`
	}

	// WaylandInfo is the result of a real connection to a Wayland compositor.
	WaylandInfo = WaylandInfoData[WaylandGlobal]

	// WaylandGlobal is one entry from the compositor's global registry.
	WaylandGlobal struct {
		// Interface fully qualified Wayland interface name.
		Interface string `json:"interface"`
		// Name numeric object name assigned by the Wayland registry.
		Name uint32 `json:"name"`
		// Version version reported by the source, when available.
		Version uint32 `json:"version"`
	}

	// GraphicsInfo describes the GPUs and the client-side graphics stack.
	GraphicsInfo = GraphicsInfoData[OpenGLInfo, VulkanInfo, GPUInfo]

	// GPUInfo is one DRM device and its PCI identity.
	GPUInfo struct {
		// ID identifier reported by the underlying system.
		ID string `json:"id"`
		// Vendor vendor name reported by the driver or server.
		Vendor string `json:"vendor"`
		// Model marketing name from the driver or PCI device database.
		Model string `json:"model"`
		// VendorID pCI vendor identifier as reported by sysfs.
		VendorID string `json:"vendor_id"`
		// DeviceID pCI device identifier as reported by sysfs.
		DeviceID string `json:"device_id"`
		// PCIAddress pCI slot address, when the GPU is on a PCI bus.
		PCIAddress string `json:"pci_address"`
		// Driver kernel driver bound to the device.
		Driver string `json:"driver"`
		// DRMDevice path to the DRM card device.
		DRMDevice string `json:"drm_device"`
		// RenderDevice path to the corresponding DRM render node.
		RenderDevice string `json:"render_device"`
		// Integrated whether the device is classified as integrated.
		Integrated bool `json:"integrated"`
		// Discrete whether the device is classified as discrete.
		Discrete bool `json:"discrete"`
	}

	// OpenGLInfo describes the OpenGL implementation.
	OpenGLInfo struct {
		// Vendor vendor name reported by the driver or server.
		Vendor string `json:"vendor,omitempty"`
		// Renderer openGL renderer string from the live driver.
		Renderer string `json:"renderer,omitempty"`
		// Version version reported by the source, when available.
		Version string `json:"version,omitempty"`
		// Available whether the implementation or service was detected.
		Available bool `json:"available"`
	}

	// VulkanInfo describes the Vulkan loader and its ICDs.
	VulkanInfo struct {
		// Version version reported by the source, when available.
		Version string `json:"version,omitempty"`
		// Available whether the implementation or service was detected.
		Available bool `json:"available"`
	}

	// GraphicsInfoData is the shared GraphicsInfoData layout, parameterized by its component types.
	GraphicsInfoData[OpenGLType, VulkanType, GPUType any] struct {
		// OpenGL openGL implementation reported by drivers or manifests.
		OpenGL OpenGLType `json:"opengl"`
		// Vulkan vulkan loader and driver availability.
		Vulkan VulkanType `json:"vulkan"`
		// PrimaryGPU identifier of the firmware-posted GPU, or the first detected GPU.
		PrimaryGPU string `json:"primary_gpu,omitempty"`
		// GPUs detected DRM devices with their PCI identity.
		GPUs []GPUType `json:"gpus"`
	}

	// SessionType is the kind of seat session the process is attached to.
	SessionType string

	// SessionInfo describes the logind seat session.
	SessionInfo = SessionInfoData[SessionType]

	// OSInfo describes the operating system and kernel.
	OSInfo struct {
		// Name human-readable name reported by the source.
		Name string `json:"name"`
		// PrettyName human-readable distribution name from os-release.
		PrettyName string `json:"pretty_name"`
		// ID identifier reported by the underlying system.
		ID string `json:"id"`
		// IDLike related distribution identifiers from os-release.
		IDLike string `json:"id_like"`
		// Version version reported by the source, when available.
		Version string `json:"version"`
		// VersionID machine-readable distribution version.
		VersionID string `json:"version_id"`
		// Kernel kernel name.
		Kernel string `json:"kernel"`
		// KernelRelease kernel release string.
		KernelRelease string `json:"kernel_release"`
		// Architecture architecture in uname-style spelling.
		Architecture string `json:"architecture"`
		// Hostname host name reported by the operating system.
		Hostname string `json:"hostname"`
	}

	// PortalInfo describes xdg-desktop-portal and which of its interfaces the
	// running backend actually exports.
	PortalInfo struct {
		// Backend preferred or matching portal implementation.
		Backend string `json:"backend,omitempty"`
		// Available whether the implementation or service was detected.
		Available bool `json:"available"`
		// DesktopPortal whether any desktop portal interface is exported.
		DesktopPortal bool `json:"desktop_portal"`
		// ScreenCast whether the ScreenCast interface is exported.
		ScreenCast bool `json:"screencast"`
		// Screenshot whether the Screenshot interface is exported.
		Screenshot bool `json:"screenshot"`
		// FileChooser whether the FileChooser interface is exported.
		FileChooser bool `json:"file_chooser"`
		// OpenURI whether the OpenURI interface is exported.
		OpenURI bool `json:"open_uri"`
		// RemoteDesktop whether the RemoteDesktop interface is exported.
		RemoteDesktop bool `json:"remote_desktop"`
		// Inhibit whether the Inhibit interface is exported.
		Inhibit bool `json:"inhibit"`
		// Notification whether the Notification interface is exported.
		Notification bool `json:"notification"`
	}
)
