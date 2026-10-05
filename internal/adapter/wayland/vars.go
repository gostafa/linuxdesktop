// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package wayland

import (
	"errors"
)

var (
	// ErrProtocol reports a frame that does not conform to the Wayland wire
	// format. It means the peer is not a Wayland compositor, not that the
	// compositor misbehaved.
	ErrProtocol = errors.New("linuxdesktop: malformed wayland message")

	// errComplete is the sentinel the sync callback raises to say the listing
	// is complete. It never escapes readGlobals.
	errComplete = errors.New("linuxdesktop: wayland registry complete")
)
