// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package port

import (
	"context"
)

type (
	// BusKind selects one of the two D-Bus instances a desktop session has.
	BusKind uint8

	// Object identifies a destination and object path on a selected bus.
	Object struct {
		Destination string
		Path        string
	}

	// PropertyQuery identifies an interface and optional property on a bus object.
	PropertyQuery struct {
		Object    Object
		Interface string
		Name      string
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
		Property(ctx context.Context, k BusKind, query *PropertyQuery) (any, error)
		// Properties reads every property on an interface in one round trip.
		Properties(ctx context.Context, k BusKind, query *PropertyQuery) (map[string]any, error)
		// Close releases any connections that were opened.
		Close() error
	}
)
