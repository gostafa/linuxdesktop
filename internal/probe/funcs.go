// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package probe

// Default supplies a fallback only when the caller supplied the zero value.
func Default[T comparable](value T, fallback func() T) T {
	var zero T

	if value == zero {
		return fallback()
	}

	return value
}
