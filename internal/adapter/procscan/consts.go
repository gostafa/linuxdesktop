// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package procscan

const (
	procDir    = "/proc"
	commFile   = "comm"
	statusFile = "status"

	// uidPrefix introduces the line in /proc/<pid>/status whose first field is
	// the process's real uid.
	uidPrefix = "Uid:"

	// maxMatches bounds the result. More than a couple of compositor-looking
	// processes means something unusual is going on, and the extra names would
	// not improve a low-confidence guess.
	maxMatches = 8

	// zero is the empty length of a line that carries no fields.
	zero = 0

	// noValue is a process that has gone away between the listing and the read.
	noValue = ""
)
