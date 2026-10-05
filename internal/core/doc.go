// Package core is the center of the hexagon: it runs the probes and merges
// their answers into an Environment.
//
// It depends on package port and package rules, and on no adapter. That is
// what makes the wiring replaceable — the root package decides which
// implementations to inject, and this package never learns whether the X11
// answer came from a real server or a stub.
//
// Detection happens in two stages. Stage one fans every independent probe out
// across goroutines, because the cost of a run is dominated by I/O that has no
// reason to be serialized: two socket handshakes, a D-Bus round trip and a few
// dozen small file reads. Stage two is a pure merge over the collected
// signals, so all the reasoning happens with every fact already in hand and
// none of it has to be redone.
//
// No probe failure is fatal. Detect always returns a populated Environment and
// reports what went wrong as a joined error, because on a real system some
// probes are expected to fail: there is no session bus in a container, no X
// server on a Wayland-only session, and no compositor at all on a TTY.
package core
