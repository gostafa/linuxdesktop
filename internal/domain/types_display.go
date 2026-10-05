// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package domain

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
