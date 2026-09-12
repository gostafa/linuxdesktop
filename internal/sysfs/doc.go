// Package sysfs is the shared primitive for reading procfs, sysfs and other
// small key=value files.
//
// It exists because os.ReadFile is the wrong tool for pseudo-filesystems: the
// st_size reported by procfs and sysfs is a lie (0, or a flat 4096), so
// ReadFile either over-allocates a 4 KiB buffer for a nine-byte value or sizes
// the buffer at zero and grows it. A detection run touches dozens of these
// files, so this package reads them through a pooled scratch buffer and
// allocates exactly once, for the value that is actually kept.
//
// The parsing helpers are here rather than in each adapter because os-release,
// uevent and the logind session files all share the same KEY=VALUE shape.
package sysfs
