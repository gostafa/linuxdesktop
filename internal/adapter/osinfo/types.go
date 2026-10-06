// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package osinfo

import (
	"context"

	"github.com/gostafa/linuxdesktop/internal/domain"
)

type (
	// Source is the contract consumed by the detection engine.
	Source interface {
		OS(ctx context.Context) (domain.OSInfo, error)
	}

	// Func delegates the probe to its configured function.
	Func[T any] func(ctx context.Context) (T, error)

	// Probe implements port.OSProbe. It is stateless.
	Probe = Func[domain.OSInfo]

	// releaseSetter fills one field of an OSInfo from one os-release value.
	releaseSetter func(*domain.OSInfo, string)
)
