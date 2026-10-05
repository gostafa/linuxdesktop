// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package wayland

import "github.com/gostafa/linuxdesktop/internal/domain"

type (
	// Probe implements port.WaylandProbe. It is stateless.
	Probe struct{}

	// registry accumulates the compositor's global listing across however many
	// reads it takes to receive it.
	registry struct {
		globals []domain.WaylandGlobal
		buf     []byte
		filled  int
	}
)
