// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package port

import (
	"context"

	"github.com/gostafa/linuxdesktop/internal/adapter/drm"
	"github.com/gostafa/linuxdesktop/internal/adapter/env"
	"github.com/gostafa/linuxdesktop/internal/adapter/osinfo"
	"github.com/gostafa/linuxdesktop/internal/adapter/procscan"
	"github.com/gostafa/linuxdesktop/internal/adapter/wayland"
	"github.com/gostafa/linuxdesktop/internal/adapter/x11"
	"github.com/gostafa/linuxdesktop/internal/domain"
)

type (
	// PropertyValue carries the heterogeneous payload of a D-Bus property.
	PropertyValue = struct {
		// Value is the decoded property payload.
		Value any
	}

	// BusKind selects one of the two D-Bus instances a desktop session has.
	BusKind uint8

	// Object identifies a destination and object path on a selected bus.
	Object struct {
		// Destination is the bus name that owns the object.
		Destination string
		// Path is the D-Bus object path.
		Path string
	}

	// PropertyAddress is the structural address of a property request.
	PropertyAddress = struct {
		// Destination is the bus name that owns the object.
		Destination string
		// Path is the D-Bus object path.
		Path string
		// Interface identifies the interface containing the property.
		Interface string
		// Name identifies the property; empty when reading the entire interface.
		Name string
	}

	// PropertyQuery identifies an interface and optional property on a bus object.
	PropertyQuery struct {
		// Object identifies the destination and object path.
		Object Object
		// Interface is the D-Bus interface containing the property.
		Interface string
		// Name is the property name; empty when reading all properties.
		Name string
	}

	// Bus is the D-Bus port. It is deliberately narrow: the library only ever
	// needs to ask whether a name is claimed, list the interfaces on an object,
	// and read a property.
	Bus interface {
		// HasOwner reports whether a well-known name currently has an owner.
		HasOwner(ctx context.Context, k BusKind, name string) (bool, error)
		// Introspect returns the raw introspection XML for an object path.
		Introspect(ctx context.Context, k BusKind, object *Object) (string, error)
		// Property reads a single property off an interface.
		Property(ctx context.Context, k BusKind, query *PropertyQuery) (PropertyValue, error)
		// Properties reads every property on an interface in one round trip.
		Properties(ctx context.Context, k BusKind, query *PropertyQuery) (map[string]any, error)
		// Close releases any connections that were opened.
		Close() error
	}

	// X11Probe connects to the X server and reports what it found. A nil result
	// with a nil error means there was no X server to talk to.
	X11Probe = x11.Source

	// WaylandProbe connects to the compositor and enumerates its global registry.
	// A nil result with a nil error means there was no compositor to talk to.
	WaylandProbe = wayland.Source

	// DesktopProbe identifies the desktop environment and, where it is cheap to
	// obtain, its version.
	DesktopProbe interface {
		Desktop(ctx context.Context, environment *domain.Env) (domain.DesktopInfo, error)
	}

	// GPUProbe enumerates the machine's GPUs and names the primary one. It is
	// separate from StackProbe because the two answer genuinely different
	// questions from unrelated sources: what hardware is present, and what the
	// client-side driver stack can do with it.
	GPUProbe = drm.Source

	// StackProbe reports the OpenGL and Vulkan implementations.
	StackProbe interface {
		Stack(ctx context.Context) (domain.OpenGLInfo, domain.VulkanInfo, error)
	}

	// PortalProbe inspects xdg-desktop-portal.
	PortalProbe interface {
		Portal(ctx context.Context, environment *domain.Env) (domain.PortalInfo, error)
	}

	// EnvProbe snapshots the process environment. It is the one probe that never
	// fails and never blocks, so it runs synchronously before everything else.
	EnvProbe = env.Source

	// OSProbe reads the distribution and kernel identity.
	OSProbe = osinfo.Source

	// SessionProbe reads the logind seat session.
	SessionProbe interface {
		Session(ctx context.Context, environment *domain.Env) (domain.SessionInfo, error)
	}

	// ProcessProbe lists the command names of the calling user's processes. It is
	// the last-resort signal for compositor detection and is only invoked when
	// explicitly enabled.
	ProcessProbe = procscan.Source
)
