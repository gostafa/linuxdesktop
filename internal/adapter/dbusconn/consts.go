// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package dbusconn

const (
	// noEntries is the initial length of an accumulating result slice.
	noEntries = 0

	errCheckOwner     = "dbus: check name owner: %w"
	errIntrospect     = "dbus: introspect object: %w"
	errReadProperty   = "dbus: read property: %w"
	errReadProperties = "dbus: read properties: %w"
)
