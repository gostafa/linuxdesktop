// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package linuxdesktop

import "time"

type (
	// ConfigData is the shared ConfigData layout, parameterized by its component types.
	ConfigData[SectionsType any] struct {
		Sections     SectionsType
		Timeout      time.Duration
		ProbeTimeout time.Duration
		NativeGL     bool
		ProcessScan  bool
	}
)
