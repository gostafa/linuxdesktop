// Package wayland enumerates a compositor's global registry.
//
// The whole protocol surface the library needs is three events, so this is a
// hand-rolled client rather than a dependency. Pulling in a generated binding
// would add a large module and an event-loop model in exchange for code we do
// not use, and would cost more at runtime than the 24 bytes this package
// writes.
//
// The exchange is one round trip. Both wl_display.get_registry and
// wl_display.sync are written in a single 24-byte frame; the compositor
// answers with every wl_registry.global event it has, followed by the
// wl_callback.done that proves the list is complete. Reading stops there.
//
// The registry is also the single best compositor fingerprint available: the
// set of advertised interfaces identifies wlroots, KWin, Mutter and COSMIC
// unambiguously, without a process scan or a compositor-specific IPC protocol.
package wayland
