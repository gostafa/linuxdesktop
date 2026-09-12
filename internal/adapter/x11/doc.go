// Package x11 connects to the X server and reports what it finds.
//
// One connection answers three questions at once: the setup reply carries the
// vendor string and protocol version, ListExtensions says whether this is a
// real X server or Xwayland, and the EWMH _NET_SUPPORTING_WM_CHECK chain names
// the running window manager. Requests are pipelined, so the whole probe costs
// three round trips rather than one per question.
//
// github.com/jezek/xgb is used rather than a hand-rolled client because it
// already parses ~/.Xauthority and performs MIT-MAGIC-COOKIE-1 authentication,
// which is not optional on a real desktop.
//
// xgb's Reply() has no deadline of its own, so the probe runs a watchdog that
// closes the connection when the context is cancelled. That turns a wedged X
// server into a prompt read error instead of a hung detection run.
package x11
