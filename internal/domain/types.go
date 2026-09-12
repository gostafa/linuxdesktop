package domain

// SessionType is the kind of seat session the process is attached to.
type SessionType string

// DisplayProtocol is the windowing protocol the process should speak.
type DisplayProtocol string

// DesktopEnvironment is a recognised desktop environment.
type DesktopEnvironment string

// CompositorKind is a recognised window manager or Wayland compositor.
type CompositorKind string

// DetectionConfidence grades how much a detection result can be trusted.
type DetectionConfidence string

// DetectionMethod records which signal produced a detection result.
type DetectionMethod string

// Section selects which parts of an Environment to populate. Sections are a
// bitmask so a caller can pay only for what it reads.
type Section uint16

// Environment is the complete picture of the desktop the process is running in.
type Environment struct {
	OS         OSInfo         `json:"os"`
	Session    SessionInfo    `json:"session"`
	Display    DisplayInfo    `json:"display"`
	Desktop    DesktopInfo    `json:"desktop"`
	Compositor CompositorInfo `json:"compositor"`
	Graphics   GraphicsInfo   `json:"graphics"`
	Portal     PortalInfo     `json:"portal"`
	Headless   bool           `json:"headless"`
}

// OSInfo describes the operating system and kernel.
type OSInfo struct {
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

// SessionInfo describes the logind seat session.
type SessionInfo struct {
	ID                string      `json:"id"`
	Seat              string      `json:"seat"`
	Type              SessionType `json:"type"`
	Desktop           string      `json:"desktop"`
	Name              string      `json:"name"`
	User              string      `json:"user"`
	Remote            bool        `json:"remote"`
	Active            bool        `json:"active"`
	VTNumber          int         `json:"vt_number"`
	XDGSessionType    string      `json:"xdg_session_type"`
	XDGSessionDesktop string      `json:"xdg_session_desktop"`
}

// DisplayInfo describes which display servers are reachable.
type DisplayInfo struct {
	Protocol         DisplayProtocol `json:"protocol"`
	WaylandDisplay   string          `json:"wayland_display"`
	X11Display       string          `json:"x11_display"`
	X11Available     bool            `json:"x11_available"`
	WaylandAvailable bool            `json:"wayland_available"`
	XWayland         bool            `json:"xwayland"`
	X11              *X11Info        `json:"x11,omitempty"`
	Wayland          *WaylandInfo    `json:"wayland,omitempty"`
}

// X11Info is the result of a real connection to an X server.
type X11Info struct {
	Display       string   `json:"display"`
	Screen        int      `json:"screen"`
	Vendor        string   `json:"vendor"`
	ProtocolMajor int      `json:"protocol_major"`
	ProtocolMinor int      `json:"protocol_minor"`
	WindowManager string   `json:"window_manager"`
	Extensions    []string `json:"extensions,omitempty"`
}

// WaylandInfo is the result of a real connection to a Wayland compositor.
type WaylandInfo struct {
	Display  string          `json:"display"`
	SocketFD int             `json:"socket_fd,omitempty"`
	Globals  []WaylandGlobal `json:"globals,omitempty"`
}

// WaylandGlobal is one entry from the compositor's global registry.
type WaylandGlobal struct {
	Name      uint32 `json:"name"`
	Interface string `json:"interface"`
	Version   uint32 `json:"version"`
}

// DesktopInfo identifies the desktop environment.
type DesktopInfo struct {
	Environment     DesktopEnvironment `json:"environment"`
	Name            string             `json:"name"`
	Version         string             `json:"version,omitempty"`
	CurrentDesktop  string             `json:"current_desktop"`
	CurrentDesktops []string           `json:"current_desktops"`
	SessionDesktop  string             `json:"session_desktop"`
	DesktopSession  string             `json:"desktop_session"`
}

// CompositorInfo identifies the compositor or window manager, and says how
// certain that identification is.
type CompositorInfo struct {
	Kind       CompositorKind      `json:"kind"`
	Name       string              `json:"name"`
	Version    string              `json:"version,omitempty"`
	Wayland    bool                `json:"wayland"`
	X11        bool                `json:"x11"`
	Confidence DetectionConfidence `json:"confidence"`
	DetectedBy DetectionMethod     `json:"detected_by"`
}

// GraphicsInfo describes the GPUs and the client-side graphics stack.
type GraphicsInfo struct {
	GPUs       []GPUInfo  `json:"gpus"`
	PrimaryGPU string     `json:"primary_gpu,omitempty"`
	OpenGL     OpenGLInfo `json:"opengl"`
	Vulkan     VulkanInfo `json:"vulkan"`
}

// GPUInfo is one DRM device and its PCI identity.
type GPUInfo struct {
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
type OpenGLInfo struct {
	Available bool   `json:"available"`
	Vendor    string `json:"vendor,omitempty"`
	Renderer  string `json:"renderer,omitempty"`
	Version   string `json:"version,omitempty"`
}

// VulkanInfo describes the Vulkan loader and its ICDs.
type VulkanInfo struct {
	Available bool   `json:"available"`
	Version   string `json:"version,omitempty"`
}

// PortalInfo describes xdg-desktop-portal and which of its interfaces the
// running backend actually exports.
type PortalInfo struct {
	Available     bool   `json:"available"`
	DesktopPortal bool   `json:"desktop_portal"`
	Backend       string `json:"backend,omitempty"`
	ScreenCast    bool   `json:"screencast"`
	Screenshot    bool   `json:"screenshot"`
	FileChooser   bool   `json:"file_chooser"`
	OpenURI       bool   `json:"open_uri"`
	RemoteDesktop bool   `json:"remote_desktop"`
	Inhibit       bool   `json:"inhibit"`
	Notification  bool   `json:"notification"`
}

// Env is every environment variable the library will ever consult, read once
// at the start of a detection run so no probe pays for a repeated lookup.
type Env struct {
	SessionID      string
	SessionType    string
	SessionDesktop string
	SessionClass   string
	CurrentDesktop string
	DesktopSession string
	WaylandDisplay string
	WaylandSocket  string
	Display        string
	RuntimeDir     string
	ConfigHome     string
	Home           string
	User           string
	Seat           string
	VTNR           string

	HyprlandSignature string
	SwaySock          string
	WayfireSocket     string
	I3Sock            string
	KDEFullSession    string
	KDESessionVersion string
	GNOMESessionID    string
	GNOMESetupDisplay string

	SSHConnection string
	SSHTTY        string
	SSHClient     string
}

// Signals carries the raw evidence classification rules reason over. It is
// filled by the protocol adapters in stage one and consumed by package rules
// in stage two.
type Signals struct {
	Env              Env
	Desktop          DesktopEnvironment
	WaylandGlobals   []WaylandGlobal
	WaylandReachable bool
	X11Reachable     bool
	X11WindowManager string
	X11Extensions    []string
	Processes        []string
}
