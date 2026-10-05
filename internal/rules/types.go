// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package rules

import "github.com/gostafa/linuxdesktop/internal/domain"

type (
	// match is one identified compositor: the enum value plus the name actually
	// observed, which matters for the compositors that have no enum value.
	match struct {
		kind domain.CompositorKind
		name string
	}

	// fingerprint ties a Wayland interface name to the compositor that is the
	// only one to advertise it.
	fingerprint struct {
		iface string
		match
	}

	// verdict is one identified compositor together with the method that found
	// it and how far that method can be trusted.
	verdict struct {
		match

		confidence domain.DetectionConfidence
		method     domain.DetectionMethod
	}

	// rung is one step of the compositor detection ladder.
	rung func(*domain.Signals) (verdict, bool)
)
