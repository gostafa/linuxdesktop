// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package procscan

type (
	// Probe implements port.ProcessProbe. filter selects the command names
	// worth confirming ownership of.
	Probe struct {
		filter func(string) bool
	}

	// scanner walks /proc once, collecting the command names of this user's
	// processes that pass the filter. seen suppresses the duplicates a
	// multi-process desktop produces in quantity.
	scanner struct {
		filter  func(string) bool
		seen    map[string]bool
		matches []string
		uid     string
	}
)
