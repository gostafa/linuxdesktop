// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package desktop

import (
	"github.com/gostafa/linuxdesktop/internal/domain"
	"github.com/gostafa/linuxdesktop/internal/probe"
)

type (
	// Probe implements port.DesktopProbe. bus may be nil, which simply leaves the
	// GNOME version unset.
	Probe = probe.DesktopFunc[domain.Env, domain.DesktopInfo]
)
