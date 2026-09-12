package x11

import (
	"bytes"
	"context"
	"strings"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"

	"github.com/gostafa/linuxdesktop/internal/domain"
	"github.com/gostafa/linuxdesktop/internal/sysfs"
)

// New returns an X11 probe.
func New() Probe { return Probe{} }

// Available is a cheap test for an X server, doing no more than a stat. A
// display on a remote host is reported as available without verification,
// since confirming it would mean connecting. Callers that need certainty
// should run the probe, which connects.
func Available(env *domain.Env) bool {
	display := env.Display
	if display == "" {
		return false
	}
	colon := strings.IndexByte(display, ':')
	if colon < 0 {
		return false
	}
	if host := display[:colon]; host != "" && host != localHost {
		return true
	}
	number := display[colon+1:]
	if dot := strings.IndexByte(number, '.'); dot >= 0 {
		number = number[:dot]
	}
	if number == "" {
		return false
	}
	return sysfs.IsSocket(unixSocketDir + "/X" + number)
}

// X11 connects to $DISPLAY and reports what the server says about itself. A
// nil result with a nil error means there was no X server to talk to, which on
// a pure Wayland or headless session is the expected outcome.
func (Probe) X11(ctx context.Context, env *domain.Env) (*domain.X11Info, error) {
	if env.Display == "" {
		return nil, nil
	}

	conn, err := xgb.NewConnDisplay(env.Display)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// xgb's Reply blocks indefinitely, so cancellation is expressed by closing
	// the connection out from under it.
	stop := make(chan struct{})
	defer close(stop)
	go func() {
		select {
		case <-ctx.Done():
			conn.Close()
		case <-stop:
		}
	}()

	setup := xproto.Setup(conn)
	info := &domain.X11Info{
		Display:       env.Display,
		Screen:        conn.DefaultScreen,
		Vendor:        setup.Vendor,
		ProtocolMajor: int(setup.ProtocolMajorVersion),
		ProtocolMinor: int(setup.ProtocolMinorVersion),
	}

	// Round trip one: every request that needs no prior answer goes out before
	// any reply is read.
	extCookie := xproto.ListExtensions(conn)
	checkCookie := internAtom(conn, atomSupportingWMCheck)
	nameCookie := internAtom(conn, atomNetWMName)
	utf8Cookie := internAtom(conn, atomUTF8String)

	if reply, err := extCookie.Reply(); err == nil && reply != nil {
		info.Extensions = make([]string, 0, len(reply.Names))
		for _, n := range reply.Names {
			info.Extensions = append(info.Extensions, n.Name)
		}
	}

	checkAtom := atomOf(checkCookie)
	nameAtom := atomOf(nameCookie)
	utf8Atom := atomOf(utf8Cookie)

	if len(setup.Roots) > conn.DefaultScreen && checkAtom != 0 {
		root := setup.Roots[conn.DefaultScreen].Root
		info.WindowManager = windowManager(conn, root, checkAtom, nameAtom, utf8Atom)
	}

	return info, nil
}

// windowManager walks the EWMH _NET_SUPPORTING_WM_CHECK chain: the root window
// points at a window owned by the window manager, and that window carries the
// manager's name.
func windowManager(conn *xgb.Conn, root xproto.Window, check, netName, utf8 xproto.Atom) string {
	reply, err := xproto.GetProperty(conn, false, root, check, xproto.AtomWindow, 0, 1).Reply()
	if err != nil || reply == nil || len(reply.Value) < 4 {
		return ""
	}
	owner := xproto.Window(xgb.Get32(reply.Value))
	if owner == 0 {
		return ""
	}

	if netName != 0 && utf8 != 0 {
		if name := textProperty(conn, owner, netName, utf8); name != "" {
			return name
		}
	}
	// Window managers that predate EWMH only set WM_NAME.
	return textProperty(conn, owner, xproto.AtomWmName, xproto.AtomString)
}

// textProperty reads a string property, tolerating the BadWindow that a stale
// _NET_SUPPORTING_WM_CHECK produces after a window manager crash.
func textProperty(conn *xgb.Conn, win xproto.Window, prop, typ xproto.Atom) string {
	reply, err := xproto.GetProperty(conn, false, win, prop, typ, 0, maxNameWords).Reply()
	if err != nil || reply == nil || len(reply.Value) == 0 {
		return ""
	}
	value := reply.Value
	if i := bytes.IndexByte(value, 0); i >= 0 {
		value = value[:i]
	}
	return string(value)
}

func internAtom(conn *xgb.Conn, name string) xproto.InternAtomCookie {
	return xproto.InternAtom(conn, true, uint16(len(name)), name)
}

// atomOf resolves an InternAtom cookie, yielding 0 when the atom does not
// exist on this server.
func atomOf(cookie xproto.InternAtomCookie) xproto.Atom {
	reply, err := cookie.Reply()
	if err != nil || reply == nil {
		return 0
	}
	return reply.Atom
}
