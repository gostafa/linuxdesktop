// Package linuxdesktop reports what kind of Linux desktop a process is running
// in: the seat session, the display protocol, the desktop environment, the
// compositor, the graphics stack and whether xdg-desktop-portal is usable.
//
// It exists so that an application answers those questions once, properly,
// instead of scattering os.Getenv("XDG_SESSION_TYPE") checks through its code
// and guessing at the rest.
//
// # Getting the whole picture
//
//	env, err := linuxdesktop.Detect()
//	if env.Display.Protocol == linuxdesktop.DisplayProtocolWayland {
//		// take the Wayland path
//	}
//
// # Getting one answer
//
// Each section has a getter that runs only the probes that section needs, so
// asking for the session type costs a single file read rather than a Wayland
// handshake:
//
//	session, _ := linuxdesktop.Session()
//
// For cancellable or partial detection, DetectContext takes both a context and
// a section mask:
//
//	env, _ := linuxdesktop.DetectContext(ctx,
//		linuxdesktop.WithSections(linuxdesktop.SectionDisplay|linuxdesktop.SectionCompositor))
//
// # Errors are information, not failure
//
// Detect never returns a nil *Environment. The error it returns is a joined
// report of probes that did not succeed, and on a healthy system some of them
// will not: a container has no session bus, a Wayland-only session has no X
// server, a TTY has no compositor. Callers that only want the data can ignore
// the error entirely. The one error worth testing for is ErrNotLinux.
//
// # Confidence
//
// CompositorInfo carries both the method that identified the compositor and
// how far that method can be trusted. A Confidence of ConfidenceHigh means the
// compositor was observed — it set its own environment variable, advertised a
// distinctive Wayland interface, or named itself over EWMH. ConfidenceMedium
// means it was inferred from the desktop environment, and ConfidenceLow means
// a matching process was found and nothing better. Code that is about to do
// something compositor-specific should check it.
//
// A compositor with no matching CompositorKind constant, such as muffin or
// cosmic-comp, is reported as CompositorUnknown with its real name in Name.
// Kind is for branching; Name is for displaying and logging.
//
// # Cost
//
// A run performs no subprocesses at all — no loginctl, no xrandr, no glxinfo.
// Probes run concurrently, connections are made once and shared, and the whole
// thing is normally a few milliseconds. The exception is WithOpenGL, which
// initializes the GPU driver to obtain a true renderer string and can take
// tens of milliseconds; it is off by default.
//
// # Portability
//
// The package builds on every platform Go supports. On anything other than
// Linux, Detect returns a zero-value Environment marked headless together with
// ErrNotLinux, so a cross-platform application can import it unconditionally.
package linuxdesktop
