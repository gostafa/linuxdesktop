// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package port

import (
	"context"

	"github.com/gostafa/linuxdesktop/internal/domain"
)

type (
	// X11Probe connects to the X server and reports what it found. A nil result
	// with a nil error means there was no X server to talk to.
	X11Probe interface {
		X11(ctx context.Context, env *domain.Env) (*domain.X11Info, error)
	}

	// WaylandProbe connects to the compositor and enumerates its global registry.
	// A nil result with a nil error means there was no compositor to talk to.
	WaylandProbe interface {
		Wayland(ctx context.Context, env *domain.Env) (*domain.WaylandInfo, error)
	}

	// DesktopProbe identifies the desktop environment and, where it is cheap to
	// obtain, its version.
	DesktopProbe interface {
		Desktop(ctx context.Context, env *domain.Env) (domain.DesktopInfo, error)
	}
)
