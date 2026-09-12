package sysfs

import "sync"

// scratch hands out reusable read buffers so that reading a few dozen sysfs
// files costs a few dozen allocations for the retained values, and none for
// the intermediate buffers.
var scratch = sync.Pool{
	New: func() any {
		b := make([]byte, scratchSize)
		return &b
	},
}
