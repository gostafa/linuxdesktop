// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package domain

import (
	"github.com/gostafa/linuxdesktop/internal/schema"
)

type (
	// Section selects which parts of an Environment to populate. Sections are
	// a bitmask so a caller can pay only for what it reads.
	Section uint16

	// Environment is the complete picture of the desktop the process is
	// running in.
	Environment = schema.Environment[
		OSInfo,
		SessionInfo,
		DisplayInfo,
		DesktopInfo,
		CompositorInfo,
		GraphicsInfo,
		PortalInfo,
	]

	// CompositorKind is a recognized window manager or Wayland compositor.
	CompositorKind string

	// DetectionConfidence grades how much a detection result can be trusted.
	DetectionConfidence string

	// DetectionMethod records which signal produced a detection result.
	DetectionMethod string

	// CompositorInfo identifies the compositor or window manager, and says how
	// certain that identification is.
	CompositorInfo = schema.CompositorInfo[CompositorKind, DetectionConfidence, DetectionMethod]

	// DesktopEnvironment is a recognized desktop environment.
	DesktopEnvironment string

	// DesktopInfo identifies the desktop environment.
	DesktopInfo = schema.DesktopInfo[DesktopEnvironment]

	// DisplayProtocol is the windowing protocol the process should speak.
	DisplayProtocol string

	// DisplayInfo describes which display servers are reachable.
	DisplayInfo = schema.DisplayInfo[X11Info, WaylandInfo, DisplayProtocol]

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
	WaylandInfo = schema.WaylandInfo[WaylandGlobal]

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
	GraphicsInfo = schema.GraphicsInfo[OpenGLInfo, VulkanInfo, GPUInfo]

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

	// SessionType is the kind of seat session the process is attached to.
	SessionType string

	// SessionInfo describes the logind seat session.
	SessionInfo = schema.SessionInfo[SessionType]

	// Env is every environment variable the library will ever consult, read
	// once at the start of a detection run so no probe pays for a repeated
	// lookup.
	Env struct {
		// SessionID unmodified XDG_SESSION_ID value.
		SessionID string
		// SessionType unmodified XDG_SESSION_TYPE value.
		SessionType string
		// SessionDesktop unmodified XDG_SESSION_DESKTOP value.
		SessionDesktop string
		// SessionClass unmodified XDG_SESSION_CLASS value.
		SessionClass string
		// CurrentDesktop unmodified XDG_CURRENT_DESKTOP value.
		CurrentDesktop string
		// DesktopSession unmodified DESKTOP_SESSION value.
		DesktopSession string
		// WaylandDisplay wayland socket name from WAYLAND_DISPLAY.
		WaylandDisplay string
		// WaylandSocket unmodified WAYLAND_SOCKET descriptor value.
		WaylandSocket string
		// Display unmodified DISPLAY value.
		Display string
		// RuntimeDir unmodified XDG_RUNTIME_DIR value.
		RuntimeDir string
		// ConfigHome unmodified XDG_CONFIG_HOME value.
		ConfigHome string
		// Home unmodified HOME value.
		Home string
		// User uSER value, falling back to LOGNAME.
		User string
		// Seat unmodified XDG_SEAT value.
		Seat string
		// VTNR unmodified XDG_VTNR value.
		VTNR string

		// HyprlandSignature unmodified HYPRLAND_INSTANCE_SIGNATURE value.
		HyprlandSignature string
		// SwaySock unmodified SWAYSOCK value.
		SwaySock string
		// WayfireSocket unmodified WAYFIRE_SOCKET value.
		WayfireSocket string
		// I3Sock unmodified I3SOCK value.
		I3Sock string
		// KDEFullSession unmodified KDE_FULL_SESSION marker.
		KDEFullSession string
		// KDESessionVersion unmodified KDE_SESSION_VERSION value.
		KDESessionVersion string
		// GNOMESessionID unmodified GNOME_DESKTOP_SESSION_ID marker.
		GNOMESessionID string
		// GNOMESetupDisplay unmodified GNOME_SETUP_DISPLAY value.
		GNOMESetupDisplay string

		// SSHConnection unmodified SSH_CONNECTION value.
		SSHConnection string
		// SSHTTY unmodified SSH_TTY value.
		SSHTTY string
		// SSHClient unmodified SSH_CLIENT value.
		SSHClient string
	}

	// Signals carries the raw evidence classification rules reason over. It is
	// filled by the protocol adapters in stage one and consumed by package
	// rules in stage two.
	Signals = evidence[Env, DesktopEnvironment, WaylandGlobal]

	// evidence combines environment, desktop, and registry signals for classification.
	evidence[E, D, G any] struct {
		// E environment snapshot taken before concurrent probes run.
		Env E
		// Desktop desktop environment identified from collected evidence.
		Desktop D
		// X11WindowManager window manager name observed through the X11 probe.
		X11WindowManager string
		// WaylandGlobals registry entries observed through the Wayland probe.
		WaylandGlobals []G
		// X11Extensions extension names observed through the X11 probe.
		X11Extensions []string
		// Processes filtered command names observed for the current user.
		Processes []string
		// WaylandReachable whether the Wayland probe connected successfully.
		WaylandReachable bool
		// X11Reachable whether the X11 probe connected successfully.
		X11Reachable bool
	}
)
