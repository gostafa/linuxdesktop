package port

import (
	"context"

	"github.com/gostafa/linuxdesktop/internal/domain"
)

// BusKind selects one of the two D-Bus instances a desktop session has.
type BusKind uint8

// The two buses.
const (
	SystemBus BusKind = iota
	SessionBus
)

// Bus is the D-Bus port. It is deliberately narrow: the library only ever
// needs to ask whether a name is claimed, list the interfaces on an object,
// and read a property.
type Bus interface {
	// HasOwner reports whether a well-known name currently has an owner.
	HasOwner(ctx context.Context, k BusKind, name string) (bool, error)
	// Introspect returns the raw introspection XML for an object path.
	Introspect(ctx context.Context, k BusKind, dest, path string) (string, error)
	// Property reads a single property off an interface.
	Property(ctx context.Context, k BusKind, dest, path, iface, prop string) (any, error)
	// Properties reads every property on an interface in one round trip.
	Properties(ctx context.Context, k BusKind, dest, path, iface string) (map[string]any, error)
	// Close releases any connections that were opened.
	Close() error
}

// EnvProbe snapshots the process environment. It is the one probe that never
// fails and never blocks, so it runs synchronously before everything else.
type EnvProbe interface {
	Snapshot() domain.Env
}

// OSProbe reads the distribution and kernel identity.
type OSProbe interface {
	OS(ctx context.Context) (domain.OSInfo, error)
}

// SessionProbe reads the logind seat session.
type SessionProbe interface {
	Session(ctx context.Context, env *domain.Env) (domain.SessionInfo, error)
}

// X11Probe connects to the X server and reports what it found. A nil result
// with a nil error means there was no X server to talk to.
type X11Probe interface {
	X11(ctx context.Context, env *domain.Env) (*domain.X11Info, error)
}

// WaylandProbe connects to the compositor and enumerates its global registry.
// A nil result with a nil error means there was no compositor to talk to.
type WaylandProbe interface {
	Wayland(ctx context.Context, env *domain.Env) (*domain.WaylandInfo, error)
}

// DesktopProbe identifies the desktop environment and, where it is cheap to
// obtain, its version.
type DesktopProbe interface {
	Desktop(ctx context.Context, env *domain.Env) (domain.DesktopInfo, error)
}

// GPUProbe enumerates the machine's GPUs and names the primary one. It is
// separate from StackProbe because the two answer genuinely different
// questions from unrelated sources: what hardware is present, and what the
// client-side driver stack can do with it.
type GPUProbe interface {
	GPUs(ctx context.Context) (gpus []domain.GPUInfo, primary string, err error)
}

// StackProbe reports the OpenGL and Vulkan implementations.
type StackProbe interface {
	Stack(ctx context.Context) (domain.OpenGLInfo, domain.VulkanInfo, error)
}

// PortalProbe inspects xdg-desktop-portal.
type PortalProbe interface {
	Portal(ctx context.Context, env *domain.Env) (domain.PortalInfo, error)
}

// ProcessProbe lists the command names of the calling user's processes. It is
// the last-resort signal for compositor detection and is only invoked when
// explicitly enabled.
type ProcessProbe interface {
	Processes(ctx context.Context) ([]string, error)
}
