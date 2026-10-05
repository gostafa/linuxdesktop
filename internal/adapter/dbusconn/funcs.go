// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package dbusconn

import (
	"context"
	"errors"
	"fmt"

	"github.com/godbus/dbus/v5"
	"github.com/gostafa/linuxdesktop/internal/port"
)

// New returns a Bus whose connections live no longer than base.
func New(base context.Context) *Bus {
	if base == nil {
		base = context.Background()
	}

	cache := &connections{open: openConnection}

	return &Bus{
		connect: func(kind port.BusKind) (*dbus.Conn, error) { return connect(base, cache, kind) },
		close:   func() error { return closeConnections(cache) },
	}
}

// HasOwner reports whether a well-known name currently has an owner.
func (bus *connectionBus[K, C]) HasOwner(ctx context.Context, kind K, name string) (bool, error) {
	conn, err := bus.connect(kind)
	if err != nil {
		return false, fmt.Errorf(errCheckOwner, err)
	}

	var has bool

	err = conn.BusObject().
		CallWithContext(ctx, "org.freedesktop.DBus.NameHasOwner", 0, name).
		Store(&has)
	if err != nil {
		return has, fmt.Errorf(errCheckOwner, err)
	}

	return has, nil
}

// Introspect returns the raw introspection XML for an object path.
func (bus *connectionBus[K, C]) Introspect(
	ctx context.Context,
	kind K,
	object *port.Object,
) (string, error) {
	conn, err := bus.connect(kind)
	if err != nil {
		return "", fmt.Errorf(errIntrospect, err)
	}

	var xml string

	err = conn.Object(object.Destination, dbus.ObjectPath(object.Path)).
		CallWithContext(ctx, "org.freedesktop.DBus.Introspectable.Introspect", 0).
		Store(&xml)
	if err != nil {
		return xml, fmt.Errorf(errIntrospect, err)
	}

	return xml, nil
}

// Property reads a single property off an interface.
func (bus *connectionBus[K, C]) Property(
	ctx context.Context,
	kind K,
	query *port.PropertyQuery,
) (any, error) {
	conn, err := bus.connect(kind)
	if err != nil {
		return nil, fmt.Errorf(errReadProperty, err)
	}

	var variant dbus.Variant

	err = conn.Object(query.Object.Destination, dbus.ObjectPath(query.Object.Path)).
		CallWithContext(ctx, "org.freedesktop.DBus.Properties.Get", 0, query.Interface, query.Name).
		Store(&variant)
	if err != nil {
		return nil, fmt.Errorf(errReadProperty, err)
	}

	return variant.Value(), nil
}

// Properties reads every property on an interface in one round trip, which is
// what makes reading a logind session cost one message instead of nine.
func (bus *connectionBus[K, C]) Properties(
	ctx context.Context,
	kind K,
	query *port.PropertyQuery,
) (map[string]any, error) {
	conn, err := bus.connect(kind)
	if err != nil {
		return nil, fmt.Errorf(errReadProperties, err)
	}

	var raw map[string]dbus.Variant

	err = conn.Object(query.Object.Destination, dbus.ObjectPath(query.Object.Path)).
		CallWithContext(ctx, "org.freedesktop.DBus.Properties.GetAll", 0, query.Interface).
		Store(&raw)
	if err != nil {
		return nil, fmt.Errorf(errReadProperties, err)
	}

	return variantValues(raw), nil
}

func variantValues(raw map[string]dbus.Variant) map[string]any {
	out := make(map[string]any, len(raw))
	for kind := range raw {
		variant := raw[kind]

		out[kind] = variant.Value()
	}

	return out
}

// Close releases whichever connections were actually opened.
func closeConnections(cache *connections) error {
	var errs []error

	for i := range cache.conns {
		if cache.conns[i] == nil {
			continue
		}

		errs = append(errs, cache.conns[i].Close())
		cache.conns[i] = nil
	}

	return errors.Join(errs...)
}

// conn establishes a bus on first use and caches the outcome, including the
// failure. A machine without a session bus should pay for exactly one failed
// connect attempt, not one per probe.
func connect(base context.Context, cache *connections, kind port.BusKind) (*dbus.Conn, error) {
	i := int(kind)
	if i >= len(cache.conns) {
		return nil, errors.New("linuxdesktop: unknown bus kind")
	}

	cache.once[i].Do(func() {
		cache.conns[i], cache.errs[i] = cache.open(base, kind)
	})

	return cache.conns[i], cache.errs[i]
}

func openConnection(base context.Context, kind port.BusKind) (*dbus.Conn, error) {
	opt := dbus.WithContext(base)

	if kind == port.SystemBus {
		conn, err := dbus.ConnectSystemBus(opt)
		if err != nil {
			return conn, fmt.Errorf("dbusconn: connect system bus: %w", err)
		}

		return conn, nil
	}

	conn, err := dbus.ConnectSessionBus(opt)
	if err != nil {
		return conn, fmt.Errorf("dbusconn: connect session bus: %w", err)
	}

	return conn, nil
}

// Close releases the private connections owned by this bus.
func (bus *connectionBus[K, C]) Close() error {
	err := bus.close()
	if err != nil {
		return fmt.Errorf("dbusconn: close connections: %w", err)
	}

	return nil
}
