// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package wayland

const (
	// wireOne is the display object ID and its get_registry request opcode.
	wireOne = 1

	// The registry and callback are the only object IDs this client allocates.
	objRegistry uint32 = 2
	objCallback uint32 = 3

	// wireZero is the sync request and error/global/done event opcode.
	// It also marks the absence of an inherited socket.
	wireZero = 0

	// wordBytes is one 32-bit wire word: the unit every Wayland field is sized
	// and aligned in.
	wordBytes = 4

	// headerSize is object id (4 bytes) plus the packed size/opcode word.
	headerSize = wordBytes + wordBytes

	// requestSize is a header plus one new_id argument. It also covers the
	// three-word prefix before a registry.global or display.error string.
	requestSize = headerSize + wordBytes

	// handshakeSize covers get_registry and sync written together.
	handshakeSize = requestSize + requestSize

	// opcodeBits is how far a message size sits above the opcode in the second
	// header word, and opcodeMask selects the opcode back out of it.
	opcodeBits = 16
	opcodeMask = 0xffff

	// defaultDisplay is libwayland's fallback when $WAYLAND_DISPLAY is unset.
	// It is only used when a socket of that name actually exists, so an X11
	// session is never mistaken for a Wayland one.
	defaultDisplay = "wayland-0"

	// unixNetwork is the only network a compositor is ever reached over.
	unixNetwork = "unix"

	// expectedGlobals is the capacity a listing is built at; a real compositor
	// advertises a few dozen interfaces.
	expectedGlobals = 48

	// Read limits. A registry is a few kilobytes at most; these bounds exist so
	// a malformed or hostile peer cannot make the client allocate without end.
	initialReadBuf  = 8192
	maxMessageBytes = 1 << 20

	// nul terminates every string on the wire.
	nul = "\x00"

	// noValue is an unset environment variable or an unreadable string.
	noValue = ""

	errDispatchEvent = "wayland: dispatch event: %w"
	errReadEvents    = "wayland: read events: %w"
	errQueryRegistry = "wayland: query registry: %w"
)
