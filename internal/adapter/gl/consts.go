// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package gl

const (
	// Manifest layout.
	jsonSuffix = ".json"
	driSuffix  = "_dri.so"

	// Vendor names inferred from a driver library filename.
	vendorMesa   = "Mesa"
	vendorNVIDIA = "NVIDIA"
	vendorAMD    = "AMD"
	vendorIntel  = "Intel"

	// versionSeparator divides the fields of a dotted API version.
	versionSeparator = "."

	// zero is the empty length a slice is built at, an unreadable version
	// field, and the point two versions compare equal at.
	zero = 0

	// noValue is an unnamed vendor or an unreported version.
	noValue = ""
)
