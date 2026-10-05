// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package x11

import (
	"github.com/gostafa/linuxdesktop/internal/domain"
	"github.com/gostafa/linuxdesktop/internal/probe"
	"github.com/jezek/xgb/xproto"
)

type (
	// Probe implements port.X11Probe. It is stateless.
	Probe = probe.X11Func[domain.Env, *domain.X11Info]

	// atomCookies holds the three InternAtom requests while they are in flight.
	atomCookies struct {
		check xproto.InternAtomCookie
		name  xproto.InternAtomCookie
		utf8  xproto.InternAtomCookie
	}

	// atomSet is those three atoms once the server has answered. Any of them
	// may be zero, meaning this server does not know that name.
	atomSet struct {
		check xproto.Atom
		name  xproto.Atom
		utf8  xproto.Atom
	}

	// property is the name and type of one X property to read. A zero name
	// means there is no property to ask for.
	property struct {
		name xproto.Atom
		typ  xproto.Atom
	}
)
