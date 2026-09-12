// Package osinfo reads the distribution and kernel identity.
//
// Distribution fields come from /etc/os-release, falling back to
// /usr/lib/os-release as the freedesktop specification requires.
//
// Kernel fields come from /proc/sys/kernel rather than from uname(2). That is
// a deliberate choice: the utsname struct has arch-dependent field types
// ([65]int8 on some platforms, [65]byte on others), so calling uname would
// force build-tagged source files. Reading three procfs entries is portable
// Go, needs no golang.org/x/sys dependency, and is marginally faster than the
// syscall wrapper.
//
// OSInfo.Kernel is the kernel's name ("Linux") and OSInfo.KernelRelease is the
// release string ("6.8.0-45-generic"), matching `uname -s` and `uname -r`.
package osinfo
