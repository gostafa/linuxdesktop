// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package port

import (
	"context"

	"github.com/gostafa/linuxdesktop/internal/domain"
)

type (
	// EnvProbe snapshots the process environment. It is the one probe that never
	// fails and never blocks, so it runs synchronously before everything else.
	EnvProbe interface {
		Snapshot() domain.Env
	}

	// OSProbe reads the distribution and kernel identity.
	OSProbe interface {
		OS(ctx context.Context) (domain.OSInfo, error)
	}

	// SessionProbe reads the logind seat session.
	SessionProbe interface {
		Session(ctx context.Context, env *domain.Env) (domain.SessionInfo, error)
	}

	// ProcessProbe lists the command names of the calling user's processes. It is
	// the last-resort signal for compositor detection and is only invoked when
	// explicitly enabled.
	ProcessProbe interface {
		Processes(ctx context.Context) ([]string, error)
	}
)
