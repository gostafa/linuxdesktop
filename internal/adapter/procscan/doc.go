// Package procscan looks for a compositor process in /proc.
//
// This is the weakest signal the library has and the only one that is opt-in:
// a running compositor process is not proof that it is the compositor serving
// this process, so a result from here never rates better than low confidence.
// It exists for the case where nothing else worked — a bare X11 session with
// no EWMH support, or a compositor that exports no distinguishing global.
//
// The scan is two-phase to keep it cheap. Every /proc/<pid>/comm is read and
// tested against a caller-supplied filter, which rejects all but a handful of
// candidates; only those survivors pay for the /proc/<pid>/status read that
// confirms the process belongs to this user. Ownership is checked by parsing
// status rather than by stat'ing the directory, because reading the real uid
// out of a FileInfo needs syscall.Stat_t, which does not exist on every
// platform this package must still compile for.
package procscan
