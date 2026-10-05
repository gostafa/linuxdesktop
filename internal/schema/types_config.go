// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package schema

import "time"

type (
	// Config is the shared Config layout, parameterized by its component types.
	Config[SectionsType any] struct {
		Sections     SectionsType
		Timeout      time.Duration
		ProbeTimeout time.Duration
		NativeGL     bool
		ProcessScan  bool
	}
)
