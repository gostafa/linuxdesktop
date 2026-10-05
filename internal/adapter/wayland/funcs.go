// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package wayland

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"path/filepath"
	"slices"
	"strconv"

	"github.com/gostafa/linuxdesktop/internal/domain"
	"github.com/gostafa/linuxdesktop/internal/sysfs"
)

// New returns a Wayland probe.
func New() Probe { return Probe{} }

// SocketPath resolves the compositor socket, returning "" when there is none.
// It only reports a path that is actually a socket, so callers can use it as a
// cheap availability test without connecting.
func SocketPath(env *domain.Env) string {
	name := displayName(env)
	if filepath.IsAbs(name) {
		return socketAt(name)
	}

	if env.RuntimeDir == noValue {
		return noValue
	}

	return socketAt(filepath.Join(env.RuntimeDir, name))
}

// Has reports whether a named interface is present in a registry listing.
func Has(globals []domain.WaylandGlobal, iface string) bool {
	for i := range globals {
		if globals[i].Interface == iface {
			return true
		}
	}

	return false
}

// Wayland connects to the compositor and enumerates its global registry. A nil
// result with a nil error means there was no compositor to talk to.
func (Probe) Wayland(ctx context.Context, env *domain.Env) (*domain.WaylandInfo, error) {
	path := SocketPath(env)
	if path == noValue {
		return inherited(env), nil
	}

	globals, err := listen(ctx, path)
	if err != nil {
		return nil, err
	}

	info := &domain.WaylandInfo{Display: displayName(env), Globals: globals}
	if fd, ok := inheritedFD(env); ok {
		info.SocketFD = fd
	}

	return info, nil
}

// newRegistry prepares a listing sized for what a real compositor advertises,
// so the common case costs one allocation apiece for the buffer and the slice.
func newRegistry() *registry {
	return &registry{
		globals: slices.Grow([]domain.WaylandGlobal(nil), expectedGlobals),
		buf:     make([]byte, initialReadBuf),
	}
}

// compact moves the bytes of a partly received message to the front, so the
// next read continues it rather than overwriting it.
func (reg *registry) compact(used int) {
	copy(reg.buf, reg.buf[used:reg.filled])

	reg.filled -= used
}

// dispatch routes one decoded message to the object it names. Anything else is
// a message this client never asked for.
func (reg *registry) dispatch(object, opcode uint32, body []byte) error {
	switch object {
	case objRegistry:
		reg.record(opcode, body)
	case objCallback:
		return complete(opcode)
	case objDisplay:
		return failure(opcode, body)
	default:
		// An object this client never allocated; nothing to do.
	}

	return nil
}

// grow doubles the buffer when it is full, refusing to pass the cap that keeps
// a hostile peer from making the client allocate without end.
func (reg *registry) grow() error {
	if reg.filled < len(reg.buf) {
		return nil
	}

	if reg.filled >= maxMessageBytes {
		return ErrProtocol
	}

	grown := make([]byte, len(reg.buf)+len(reg.buf))
	copy(grown, reg.buf)

	reg.buf = grown

	return nil
}

// message handles the one message starting at offset and reports its size. The
// false result means the buffer does not hold all of it yet.
func (reg *registry) message(offset int) (int, bool, error) {
	order := binary.NativeEndian
	object := order.Uint32(reg.buf[offset:])
	word := order.Uint32(reg.buf[offset+wordBytes:])
	size := int(word >> opcodeBits)

	if size < headerSize || size > maxMessageBytes {
		return size, false, ErrProtocol
	}

	if reg.filled-offset < size {
		return size, false, nil
	}

	body := reg.buf[offset+headerSize : offset+size]

	return size, true, reg.dispatch(object, word&opcodeMask, body)
}

// parse handles every message the buffer now holds in full, then compacts what
// is left of a partly received one to the front.
func (reg *registry) parse() error {
	var used int

	for reg.filled-used >= headerSize {
		size, ready, err := reg.message(used)
		if err != nil {
			return err
		}

		if !ready {
			break
		}

		used += size
	}

	reg.compact(used)

	return nil
}

// record keeps a wl_registry.global event and ignores every other one.
func (reg *registry) record(opcode uint32, body []byte) {
	if opcode != evRegistryGlobal {
		return
	}

	if global, ok := parseGlobal(body); ok {
		reg.globals = append(reg.globals, global)
	}
}

// step reads one chunk from the compositor and handles every message it
// completes. Parsing happens before the read error is acted on: a compositor is
// entitled to send the whole registry and close, delivering the data and the
// EOF in one read, and that data is exactly what was asked for.
func (reg *registry) step(reader io.Reader) error {
	err := reg.grow()
	if err != nil {
		return err
	}

	before := reg.filled
	got, rerr := reader.Read(reg.buf[reg.filled:])

	reg.filled += got

	err = reg.parse()
	if err != nil {
		return err
	}

	return stalled(reg.filled == before, rerr)
}

// listen dials the compositor and returns the registry it advertises.
func listen(ctx context.Context, path string) ([]domain.WaylandGlobal, error) {
	var dialer net.Dialer

	conn, err := dialer.DialContext(ctx, unixNetwork, path)
	if err != nil {
		return nil, err
	}

	globals, err := converse(ctx, conn)

	closed := conn.Close()
	if err == nil {
		err = closed
	}

	return globals, err
}

// converse bounds the connection, sends the handshake and reads the reply.
func converse(ctx context.Context, conn net.Conn) ([]domain.WaylandGlobal, error) {
	err := bound(ctx, conn)
	if err != nil {
		return nil, err
	}

	err = handshake(conn)
	if err != nil {
		return nil, err
	}

	return readGlobals(conn)
}

// bound applies the caller's deadline to the connection, if it set one.
func bound(ctx context.Context, conn net.Conn) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return nil
	}

	return conn.SetDeadline(deadline)
}

// handshake writes wl_display.get_registry and wl_display.sync as one frame, so
// the entire exchange costs a single round trip.
func handshake(writer io.Writer) error {
	var buf [handshakeSize]byte

	putRequest(buf[:requestSize], opDisplayGetRegistry, objRegistry)
	putRequest(buf[requestSize:], opDisplaySync, objCallback)

	_, err := writer.Write(buf[:])

	return err
}

// putRequest writes one wl_display request: the object id, the packed
// size/opcode word, and the new_id the request allocates.
func putRequest(frame []byte, opcode, newID uint32) {
	order := binary.NativeEndian

	order.PutUint32(frame, objDisplay)
	order.PutUint32(frame[wordBytes:], requestSize<<opcodeBits|opcode)
	order.PutUint32(frame[headerSize:], newID)
}

// readGlobals collects wl_registry.global events until the sync callback fires,
// which is the compositor's guarantee that the listing is complete.
func readGlobals(reader io.Reader) ([]domain.WaylandGlobal, error) {
	reg := newRegistry()

	for {
		err := reg.step(reader)
		if err != nil {
			return reg.globals, finished(err)
		}
	}
}

// finished turns the sentinel that ends a complete listing back into success.
func finished(err error) error {
	if errors.Is(err, errComplete) {
		return nil
	}

	return err
}

// stalled reports the read error, or the failure to make any progress that
// would otherwise spin the loop forever against a silent peer.
func stalled(stuck bool, err error) error {
	if err != nil {
		return err
	}

	if stuck {
		return io.ErrNoProgress
	}

	return nil
}

// complete reports the sync callback that ends the listing.
func complete(opcode uint32) error {
	if opcode != evCallbackDone {
		return nil
	}

	return errComplete
}

// failure turns a wl_display.error event into the error it describes.
func failure(opcode uint32, body []byte) error {
	if opcode != evDisplayError {
		return nil
	}

	return parseError(body)
}

// parseGlobal decodes name, interface and version. Wayland strings are a
// length-prefixed, NUL-terminated blob padded to a word boundary.
func parseGlobal(body []byte) (domain.WaylandGlobal, bool) {
	if len(body) < fixedBody {
		return domain.WaylandGlobal{}, false
	}

	order := binary.NativeEndian
	length := int(order.Uint32(body[wordBytes:]))
	rest := body[wordBytes+wordBytes:]
	padded := align(length)

	// A padded length below one word means an empty interface name, which no
	// real global has.
	if padded < wordBytes || padded+wordBytes > len(rest) {
		return domain.WaylandGlobal{}, false
	}

	return domain.WaylandGlobal{
		Name:      order.Uint32(body),
		Interface: text(rest, length),
		Version:   order.Uint32(rest[padded:]),
	}, true
}

// parseError turns a wl_display.error event into a Go error.
func parseError(body []byte) error {
	if len(body) < fixedBody {
		return ErrProtocol
	}

	order := binary.NativeEndian
	code := order.Uint32(body[wordBytes:])
	length := int(order.Uint32(body[wordBytes+wordBytes:]))

	return fmt.Errorf(
		"linuxdesktop: wayland protocol error %d: %s",
		code, text(body[fixedBody:], length),
	)
}

// text lifts the NUL-terminated string of the declared length off the front of
// rest, yielding "" when that length does not fit.
func text(rest []byte, length int) string {
	if length > len(rest) {
		return noValue
	}

	return string(bytes.TrimRight(rest[:length], nul))
}

// align rounds a Wayland string length up to the word boundary the wire format
// pads it to.
func align(length int) int {
	return length + (wordBytes-length%wordBytes)%wordBytes
}

// socketAt returns path when a unix socket lives there, and "" otherwise.
func socketAt(path string) string {
	if sysfs.IsSocket(path) {
		return path
	}

	return noValue
}

func displayName(env *domain.Env) string {
	if env.WaylandDisplay != noValue {
		return env.WaylandDisplay
	}

	return defaultDisplay
}

// inherited reports the connected socket a compositor handed down when there is
// no path to dial, and nothing when there is none.
func inherited(env *domain.Env) *domain.WaylandInfo {
	fd, ok := inheritedFD(env)
	if !ok {
		return nil
	}

	return &domain.WaylandInfo{Display: env.WaylandDisplay, SocketFD: fd}
}

// inheritedFD reads $WAYLAND_SOCKET. The descriptor is only ever reported,
// never consumed: it belongs to the host application.
func inheritedFD(env *domain.Env) (int, bool) {
	if env.WaylandSocket == noValue {
		return noFD, false
	}

	fd, err := strconv.Atoi(env.WaylandSocket)
	if err != nil || fd < noFD {
		return noFD, false
	}

	return fd, true
}
