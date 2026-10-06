// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package env

import (
	"github.com/gostafa/linuxdesktop/internal/domain"
)

type (
	// Source is the contract consumed by the detection engine.
	Source interface {
		Snapshot() domain.Env
	}

	// Func delegates the probe to its configured function.
	Func[T any] func() T

	// Probe implements port.EnvProbe. It is stateless, so the zero value is ready
	// to use and copies are free.
	Probe = Func[domain.Env]
)
