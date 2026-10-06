// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package procscan

import (
	"context"
)

type (
	// Source is the contract consumed by the detection engine.
	Source interface {
		Processes(ctx context.Context) ([]string, error)
	}

	// Probe implements port.ProcessProbe. filter selects the command names
	// worth confirming ownership of.
	Probe struct {
		filter func(string) bool
		root   string
	}

	// scanner walks /proc once, collecting the command names of this user's
	// processes that pass the filter. seen suppresses the duplicates a
	// multi-process desktop produces in quantity.
	scanner = struct {
		filter  func(string) bool
		root    string
		seen    map[string]bool
		uid     string
		matches []string
	}
)
