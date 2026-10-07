// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package linuxdesktop

import (
	"errors"
	"time"

	"github.com/gostafa/linuxdesktop/internal/adapter/dbusconn"
	"github.com/gostafa/linuxdesktop/internal/adapter/desktop"
	"github.com/gostafa/linuxdesktop/internal/adapter/drm"
	"github.com/gostafa/linuxdesktop/internal/adapter/env"
	"github.com/gostafa/linuxdesktop/internal/adapter/gl"
	"github.com/gostafa/linuxdesktop/internal/adapter/logind"
	"github.com/gostafa/linuxdesktop/internal/adapter/osinfo"
	"github.com/gostafa/linuxdesktop/internal/adapter/portal"
	"github.com/gostafa/linuxdesktop/internal/adapter/procscan"
	"github.com/gostafa/linuxdesktop/internal/adapter/wayland"
	"github.com/gostafa/linuxdesktop/internal/adapter/x11"
	"github.com/gostafa/linuxdesktop/internal/port"
)

var (
	// Compile-time proof that every adapter still satisfies the port it is wired
	// to. These cost nothing at runtime and catch a broken signature at build time
	// rather than at the injection site.
	_ port.EnvProbe     = env.Probe(nil)
	_ port.OSProbe      = osinfo.Probe(nil)
	_ port.SessionProbe = logind.Probe(nil)
	_ port.X11Probe     = x11.Probe(nil)
	_ port.WaylandProbe = wayland.Probe(nil)
	_ port.DesktopProbe = desktop.Probe(nil)
	_ port.GPUProbe     = drm.Probe(nil)
	_ port.StackProbe   = gl.Probe(nil)
	_ port.PortalProbe  = portal.Probe(nil)
	_ port.ProcessProbe = (*procscan.Probe)(nil)
	_ port.Bus          = (*dbusconn.Bus)(nil)

	// ErrNotLinux is joined into the error from Detect when the program is not
	// running on Linux. The returned Environment is still valid: it is zero-valued
	// and marked headless, so a cross-platform caller can import this package
	// unconditionally and branch on the result rather than on the build tag.
	ErrNotLinux = errors.New("linuxdesktop: not running on Linux")

	// ErrInvalidRetryPolicy reports invalid D-Bus connection retry settings.
	ErrInvalidRetryPolicy = errors.New("linuxdesktop: invalid connection retry policy")

	// compile-time assertion that the exported defaults stay durations.
	_ time.Duration = DefaultTimeout
)
