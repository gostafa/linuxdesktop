// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package env

import (
	"github.com/gostafa/linuxdesktop/internal/domain"
	"github.com/gostafa/linuxdesktop/internal/probe"
)

type (
	// Probe implements port.EnvProbe. It is stateless, so the zero value is ready
	// to use and copies are free.
	Probe = probe.SnapshotFunc[domain.Env]
)
