package wayland

// Object ids. wl_display is always 1; the other two are the only ids this
// client ever allocates.
const (
	objDisplay  uint32 = 1
	objRegistry uint32 = 2
	objCallback uint32 = 3
)

// Request opcodes on wl_display.
const (
	opDisplaySync        uint32 = 0
	opDisplayGetRegistry uint32 = 1
)

// Event opcodes.
const (
	evDisplayError   uint32 = 0
	evRegistryGlobal uint32 = 0
	evCallbackDone   uint32 = 0
)

// Wire layout.
const (
	// headerSize is object id (4 bytes) plus the packed size/opcode word.
	headerSize = 8
	// requestSize is a header plus one new_id argument.
	requestSize = headerSize + 4
	// handshakeSize covers get_registry and sync written together.
	handshakeSize = requestSize * 2
	// minGlobalBody is name, string length and version with an empty string.
	minGlobalBody = 12
)

// defaultDisplay is libwayland's fallback when $WAYLAND_DISPLAY is unset. It
// is only used when a socket of that name actually exists, so an X11 session
// is never mistaken for a Wayland one.
const defaultDisplay = "wayland-0"

// Read limits. A registry is a few kilobytes at most; these bounds exist so a
// malformed or hostile peer cannot make the client allocate without end.
const (
	initialReadBuf  = 8192
	maxMessageBytes = 1 << 20
)
