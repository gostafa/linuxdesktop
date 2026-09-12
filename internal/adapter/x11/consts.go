package x11

// EWMH and ICCCM atom names.
const (
	atomSupportingWMCheck = "_NET_SUPPORTING_WM_CHECK"
	atomNetWMName         = "_NET_WM_NAME"
	atomUTF8String        = "UTF8_STRING"
)

// Local display sockets. "unix" is the historical spelling for a local
// display in $DISPLAY and means the same thing as an empty host.
const (
	unixSocketDir = "/tmp/.X11-unix"
	localHost     = "unix"
)

// maxNameWords bounds the GetProperty read for a window manager name, in
// 32-bit words. 64 words is 256 bytes, far more than any real name.
const maxNameWords = 64
