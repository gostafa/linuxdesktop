// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package rules

import (
	"github.com/gostafa/linuxdesktop/internal/domain"
)

type (
	// envFingerprint binds an environment value to its compositor identity.
	envFingerprint = environmentFingerprint[domain.Env, match]

	// environmentFingerprint describes evidence supplied by an environment snapshot.
	environmentFingerprint[E, M any] struct {
		value func(*E) string
		match M
	}

	// match is one identified compositor: the enum value plus the name actually
	// observed, which matters for the compositors that have no enum value.
	match = identity[domain.CompositorKind]

	// identity associates a classification key with the observed name.
	identity[K any] struct {
		kind K
		name string
	}

	// fingerprint ties a Wayland interface name to the compositor that is the
	// only one to advertise it.
	fingerprint = interfaceFingerprint[match]

	// interfaceFingerprint maps a protocol interface to an identified implementation.
	interfaceFingerprint[M any] struct {
		match M
		iface string
	}

	// verdict is one identified compositor together with the method that found
	// it and how far that method can be trusted.
	verdict = classified[match, domain.DetectionConfidence, domain.DetectionMethod]

	// classified combines an identity with its confidence and evidence source.
	classified[M, C, D any] struct {
		match M

		confidence C
		method     D
	}

	// rung is one step of the compositor detection ladder.
	rung func(*domain.Signals) (verdict, bool)
)
