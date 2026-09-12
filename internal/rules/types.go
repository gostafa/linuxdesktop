package rules

import "github.com/gostafa/linuxdesktop/internal/domain"

// match is one identified compositor: the enum value plus the name actually
// observed, which matters for the compositors that have no enum value.
type match struct {
	kind domain.CompositorKind
	name string
}

// fingerprint ties a Wayland interface name to the compositor that is the only
// one to advertise it.
type fingerprint struct {
	iface string
	match
}
