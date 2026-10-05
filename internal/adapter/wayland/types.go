// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package wayland

import (
	"github.com/gostafa/linuxdesktop/internal/domain"
	"github.com/gostafa/linuxdesktop/internal/probe"
)

type (
	// wireEvent is one decoded protocol message.
	wireEvent struct {
		body   []byte
		object uint32
		opcode uint32
	}

	// Probe implements port.WaylandProbe. It is stateless.
	Probe = probe.WaylandFunc[domain.Env, *domain.WaylandInfo]

	// registry accumulates the compositor's global listing across however many
	// reads it takes to receive it.
	registry = registryBuffer[domain.WaylandGlobal]

	// registryBuffer retains decoded globals and incomplete protocol frames.
	registryBuffer[G any] struct {
		globals []G
		buf     []byte
		filled  int
	}
)
