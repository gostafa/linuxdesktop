// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package wayland

const (
	// Object ids. wl_display is always 1; the other two are the only ids this
	// client ever allocates.
	objDisplay  uint32 = 1
	objRegistry uint32 = 2
	objCallback uint32 = 3

	// Request opcodes on wl_display.
	opDisplaySync        uint32 = 0
	opDisplayGetRegistry uint32 = 1

	// Event opcodes.
	evDisplayError   uint32 = 0
	evRegistryGlobal uint32 = 0
	evCallbackDone   uint32 = 0

	// wordBytes is one 32-bit wire word: the unit every Wayland field is sized
	// and aligned in.
	wordBytes = 4

	// headerSize is object id (4 bytes) plus the packed size/opcode word.
	headerSize = wordBytes + wordBytes

	// requestSize is a header plus one new_id argument.
	requestSize = headerSize + wordBytes

	// handshakeSize covers get_registry and sync written together.
	handshakeSize = requestSize + requestSize

	// fixedBody is the part of a wl_registry.global or a wl_display.error body
	// that precedes its string: three 32-bit fields.
	fixedBody = wordBytes * 3

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

	// noValue is an unset environment variable or an unreadable string, and
	// noFD is the absence of an inherited socket.
	noValue = ""
	noFD    = 0
)
