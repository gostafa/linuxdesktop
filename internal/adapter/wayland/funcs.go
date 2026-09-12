package wayland

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"path/filepath"
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
	name := env.WaylandDisplay
	if name == "" {
		name = defaultDisplay
	}
	if filepath.IsAbs(name) {
		if sysfs.IsSocket(name) {
			return name
		}
		return ""
	}
	if env.RuntimeDir == "" {
		return ""
	}
	if path := filepath.Join(env.RuntimeDir, name); sysfs.IsSocket(path) {
		return path
	}
	return ""
}

// Wayland connects to the compositor and enumerates its global registry. A nil
// result with a nil error means there was no compositor to talk to.
func (Probe) Wayland(ctx context.Context, env *domain.Env) (*domain.WaylandInfo, error) {
	path := SocketPath(env)
	if path == "" {
		// A process launched by a compositor may inherit a connected socket
		// instead of a path. Report it, but never consume it: the fd belongs
		// to the host application.
		if fd, ok := inheritedFD(env); ok {
			return &domain.WaylandInfo{Display: env.WaylandDisplay, SocketFD: fd}, nil
		}
		return nil, nil
	}

	var dialer net.Dialer
	conn, err := dialer.DialContext(ctx, "unix", path)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if deadline, ok := ctx.Deadline(); ok {
		_ = conn.SetDeadline(deadline)
	}

	if err := handshake(conn); err != nil {
		return nil, err
	}
	globals, err := readGlobals(conn)
	if err != nil {
		return nil, err
	}

	info := &domain.WaylandInfo{Display: displayName(env), Globals: globals}
	if fd, ok := inheritedFD(env); ok {
		info.SocketFD = fd
	}
	return info, nil
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

// handshake writes wl_display.get_registry and wl_display.sync as one frame,
// so the entire exchange costs a single round trip.
func handshake(w io.Writer) error {
	var buf [handshakeSize]byte
	e := binary.NativeEndian

	e.PutUint32(buf[0:], objDisplay)
	e.PutUint32(buf[4:], requestSize<<16|opDisplayGetRegistry)
	e.PutUint32(buf[8:], objRegistry)

	e.PutUint32(buf[12:], objDisplay)
	e.PutUint32(buf[16:], requestSize<<16|opDisplaySync)
	e.PutUint32(buf[20:], objCallback)

	_, err := w.Write(buf[:])
	return err
}

// readGlobals collects wl_registry.global events until the sync callback
// fires, which is the compositor's guarantee that the listing is complete.
func readGlobals(r io.Reader) ([]domain.WaylandGlobal, error) {
	e := binary.NativeEndian
	buf := make([]byte, initialReadBuf)
	globals := make([]domain.WaylandGlobal, 0, 48)
	filled := 0

	for {
		if filled == len(buf) {
			if filled >= maxMessageBytes {
				return globals, ErrProtocol
			}
			grown := make([]byte, len(buf)*2)
			copy(grown, buf)
			buf = grown
		}

		n, rerr := r.Read(buf[filled:])
		filled += n
		if n == 0 && rerr == nil {
			rerr = io.ErrNoProgress
		}

		// Parse before acting on the error: a compositor is entitled to send
		// the whole registry and close, delivering the data and the EOF in one
		// read, and that data is exactly what was asked for.
		used := 0
		for filled-used >= headerSize {
			object := e.Uint32(buf[used:])
			word := e.Uint32(buf[used+4:])
			size := int(word >> 16)
			opcode := word & 0xffff

			if size < headerSize || size > maxMessageBytes {
				return globals, ErrProtocol
			}
			if filled-used < size {
				break
			}
			body := buf[used+headerSize : used+size]
			used += size

			switch object {
			case objRegistry:
				if opcode == evRegistryGlobal {
					if g, ok := parseGlobal(body); ok {
						globals = append(globals, g)
					}
				}
			case objCallback:
				if opcode == evCallbackDone {
					return globals, nil
				}
			case objDisplay:
				if opcode == evDisplayError {
					return globals, parseError(body)
				}
			}
		}

		copy(buf, buf[used:filled])
		filled -= used

		if rerr != nil {
			return globals, rerr
		}
	}
}

// parseGlobal decodes name, interface and version. Wayland strings are a
// length-prefixed, NUL-terminated blob padded to a 4-byte boundary.
func parseGlobal(body []byte) (domain.WaylandGlobal, bool) {
	e := binary.NativeEndian
	if len(body) < minGlobalBody {
		return domain.WaylandGlobal{}, false
	}
	length := int(e.Uint32(body[4:]))
	if length < 1 {
		return domain.WaylandGlobal{}, false
	}
	padded := (length + 3) &^ 3
	if 8+padded+4 > len(body) {
		return domain.WaylandGlobal{}, false
	}
	return domain.WaylandGlobal{
		Name:      e.Uint32(body),
		Interface: string(body[8 : 8+length-1]),
		Version:   e.Uint32(body[8+padded:]),
	}, true
}

// parseError turns a wl_display.error event into a Go error.
func parseError(body []byte) error {
	e := binary.NativeEndian
	if len(body) < 12 {
		return ErrProtocol
	}
	code := e.Uint32(body[4:])
	length := int(e.Uint32(body[8:]))
	message := ""
	if length > 1 && 12+length <= len(body) {
		message = string(body[12 : 12+length-1])
	}
	return fmt.Errorf("linuxdesktop: wayland protocol error %d: %s", code, message)
}

func displayName(env *domain.Env) string {
	if env.WaylandDisplay != "" {
		return env.WaylandDisplay
	}
	return defaultDisplay
}

func inheritedFD(env *domain.Env) (int, bool) {
	if env.WaylandSocket == "" {
		return 0, false
	}
	fd, err := strconv.Atoi(env.WaylandSocket)
	if err != nil || fd < 0 {
		return 0, false
	}
	return fd, true
}
