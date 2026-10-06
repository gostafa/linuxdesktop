// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package wayland

import (
	"context"

	"github.com/gostafa/linuxdesktop/internal/domain"
)

type (
	// Source is the contract consumed by the detection engine.
	Source interface {
		Wayland(ctx context.Context, env *domain.Env) (*domain.WaylandInfo, error)
	}

	// Func delegates the probe to its configured function.
	Func[E, T any] func(ctx context.Context, env *E) (T, error)

	// wireEvent is one decoded protocol message.
	wireEvent struct {
		body   []byte
		object uint32
		opcode uint32
	}

	// Probe implements port.WaylandProbe. It is stateless.
	Probe = Func[domain.Env, *domain.WaylandInfo]

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
