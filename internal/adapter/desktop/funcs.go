// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package desktop

import (
	"context"

	"github.com/gostafa/linuxdesktop/internal/domain"
	"github.com/gostafa/linuxdesktop/internal/port"
	"github.com/gostafa/linuxdesktop/internal/rules"
)

// New returns a desktop probe. bus may be nil to skip the version lookup.
func New(bus port.Bus) Probe {
	return func(ctx context.Context, env *domain.Env) (domain.DesktopInfo, error) { return desktop(ctx, env, bus) }
}

// Desktop identifies the desktop environment and, where one is cheaply
// available, its version.
func desktop(ctx context.Context, env *domain.Env, bus port.Bus) (domain.DesktopInfo, error) {
	info := rules.Desktop(env)

	info.Version = env.KDESessionVersion

	info.Version = version(ctx, &info, bus)

	return info, nil
}

func version(
	ctx context.Context,
	info *domain.DesktopInfo,
	bus port.Bus,
) string {
	if info.Environment == domain.DesktopGNOME {
		return gnomeVersion(ctx, bus)
	}

	if info.Environment == domain.DesktopKDE {
		return info.Version
	}

	return ""
}

func gnomeVersion(ctx context.Context, bus port.Bus) string {
	if bus == nil {
		return ""
	}

	query := port.PropertyQuery{
		Object:    port.Object{Destination: shellName, Path: shellPath},
		Interface: shellName, Name: shellVersion,
	}

	value, err := bus.Property(ctx, port.SessionBus, &query)
	if err != nil {
		return ""
	}

	s, ok := value.(string)
	if !ok {
		return ""
	}

	return s
}
