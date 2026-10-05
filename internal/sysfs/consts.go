// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package sysfs

const (
	// scratchSize is the initial size of a pooled read buffer. It comfortably
	// holds every file this library reads except /usr/share/hwdata/pci.ids,
	// which is streamed instead.
	scratchSize = 4096

	// maxFileSize caps a single read so a malformed or hostile path cannot
	// make the library allocate without bound.
	maxFileSize = 1 << 20

	// bufferGrowth is the factor a scratch buffer grows by when the file turns
	// out to be larger than it.
	bufferGrowth = 2

	// zero is the empty length and equally the index of a leading byte. It is
	// one constant rather than two because a package may not declare two
	// constants sharing a value.
	zero = 0

	// quoteWidth is the single byte a surrounding quote occupies at each end.
	// A slice shorter than two bytes therefore cannot carry a matching pair.
	quoteWidth = 1

	// commentByte introduces a comment line in a KEY=VALUE file.
	commentByte = '#'

	// separatorByte divides a key from its value.
	separatorByte = '='

	// noValue is the answer for an absent key or an unreadable file.
	noValue = ""

	errReadChunk = "sysfs: read chunk: %w"
)
