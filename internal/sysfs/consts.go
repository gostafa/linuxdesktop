package sysfs

// scratchSize is the initial size of a pooled read buffer. It comfortably
// holds every file this library reads except /usr/share/hwdata/pci.ids, which
// is streamed instead.
const scratchSize = 4096

// maxFileSize caps a single read so a malformed or hostile path cannot make
// the library allocate without bound.
const maxFileSize = 1 << 20
