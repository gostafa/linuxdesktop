package portal

// The portal service on the session bus.
const (
	busName     = "org.freedesktop.portal.Desktop"
	objectPath  = "/org/freedesktop/portal/desktop"
	ifacePrefix = "org.freedesktop.portal."
)

// Interface suffixes this library reports on.
const (
	ifaceScreenCast    = "ScreenCast"
	ifaceScreenshot    = "Screenshot"
	ifaceFileChooser   = "FileChooser"
	ifaceOpenURI       = "OpenURI"
	ifaceRemoteDesktop = "RemoteDesktop"
	ifaceInhibit       = "Inhibit"
	ifaceNotification  = "Notification"
)

// Backend discovery.
const (
	portalsDir   = "/usr/share/xdg-desktop-portal/portals"
	shareDir     = "/usr/share/xdg-desktop-portal"
	configName   = "portals.conf"
	configSuffix = "-portals.conf"
	portalSuffix = ".portal"

	keyDefault  = "default"
	keyUseIn    = "UseIn"
	keyDBusName = "DBusName"

	listSeparator = ";"
)
