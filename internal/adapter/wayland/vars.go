package wayland

import "errors"

// ErrProtocol reports a frame that does not conform to the Wayland wire
// format. It means the peer is not a Wayland compositor, not that the
// compositor misbehaved.
var ErrProtocol = errors.New("linuxdesktop: malformed wayland message")
