// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package port

const (
	// SystemBus selects the system D-Bus instance; SessionBus selects the user session instance.
	SystemBus BusKind = iota
	// SessionBus selects the user session D-Bus instance.
	SessionBus
)
