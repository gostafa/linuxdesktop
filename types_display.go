// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package linuxdesktop

type (
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
)
