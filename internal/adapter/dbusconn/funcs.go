// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package dbusconn

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/gostafa/linuxdesktop/internal/port"
	"github.com/gostafa/linuxdesktop/internal/probe"
	"github.com/gostafa/linuxdesktop/internal/schema"
)

// New returns a Bus whose connections live no longer than base.
func New(base context.Context) *Bus {
	return NewConfigured(base, schema.DefaultRetryPolicy(), defaultInitializationTimeout)
}

// NewConfigured creates run-owned connections with a shared initialization budget.
func NewConfigured(base context.Context, policy schema.RetryPolicy, timeout time.Duration) *Bus {
	cache := newConnections(
		probe.Default(base, context.Background),
		policy,
		timeout,
		openConnection,
	)

	return &Bus{source: connectionSource[port.BusKind, *dbus.Conn]{
		Acquire: func(ctx context.Context, kind port.BusKind) (*dbus.Conn, error) { return connect(ctx, cache, kind) },
		Release: func() error { return closeConnections(cache) },
	}}
}

// HasOwner reports whether a well-known name currently has an owner.
func (bus *connectionBus[K, C, O, Q]) HasOwner(
	ctx context.Context,
	kind K,
	name string,
) (bool, error) {
	result, err := checkOwner(ctx, selectConnection(ctx, bus.source, kind), name)
	if err != nil {
		return result, fmt.Errorf("dbusconn: owner query: %w", err)
	}

	return result, nil
}

func checkOwner[C connection](
	ctx context.Context,
	connect func() (C, error),
	name string,
) (bool, error) {
	conn, err := connect()
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
func (bus *connectionBus[K, C, O, Q]) Introspect(
	ctx context.Context,
	kind K,
	object O,
) (string, error) {
	result, err := readIntrospection(ctx, selectConnection(ctx, bus.source, kind), object)

	return result, errors.Join(err)
}

func readIntrospection[C connection, O interface {
	Address() (destination, path string)
}](ctx context.Context, connect func() (C, error), object O) (string, error) {
	conn, err := connect()
	if err != nil {
		return "", fmt.Errorf(errIntrospect, err)
	}

	destination, path := object.Address()

	var xml string

	err = conn.Object(destination, dbus.ObjectPath(path)).
		CallWithContext(ctx, "org.freedesktop.DBus.Introspectable.Introspect", 0).
		Store(&xml)
	if err != nil {
		return xml, fmt.Errorf(errIntrospect, err)
	}

	return xml, nil
}

// Property reads a single property off an interface.
func (bus *connectionBus[K, C, O, Q]) Property(
	ctx context.Context,
	kind K,
	query Q,
) (struct{ Value any }, error) {
	result, err := readProperty(ctx, selectConnection(ctx, bus.source, kind), query)
	if err != nil {
		return struct{ Value any }{Value: nil}, fmt.Errorf(errReadProperty, err)
	}

	return struct{ Value any }{Value: result.Value()}, nil
}

func readProperty[C connection, Q interface {
	PropertyAddress() port.PropertyAddress
}](ctx context.Context, connect func() (C, error), query Q) (dbus.Variant, error) {
	conn, err := connect()
	if err != nil {
		return dbus.Variant{}, fmt.Errorf(errReadProperty, err)
	}

	address := query.PropertyAddress()

	var variant dbus.Variant

	err = conn.Object(address.Destination, dbus.ObjectPath(address.Path)).
		CallWithContext(ctx, "org.freedesktop.DBus.Properties.Get", 0, address.Interface, address.Name).
		Store(&variant)
	if err != nil {
		return dbus.Variant{}, fmt.Errorf(errReadProperty, err)
	}

	return variant, nil
}

// Properties reads every property on an interface in one round trip, which is
// what makes reading a logind session cost one message instead of nine.
func (bus *connectionBus[K, C, O, Q]) Properties(
	ctx context.Context,
	kind K,
	query Q,
) (map[string]any, error) {
	result, err := readProperties(ctx, selectConnection(ctx, bus.source, kind), query)
	if err != nil {
		return nil, fmt.Errorf(errReadProperties, err)
	}

	return result, nil
}

func readProperties[C connection, Q interface {
	PropertyAddress() port.PropertyAddress
}](ctx context.Context, connect func() (C, error), query Q) (map[string]any, error) {
	conn, err := connect()
	if err != nil {
		return nil, fmt.Errorf(errReadProperties, err)
	}

	address := query.PropertyAddress()

	var raw map[string]dbus.Variant

	err = conn.Object(address.Destination, dbus.ObjectPath(address.Path)).
		CallWithContext(ctx, "org.freedesktop.DBus.Properties.GetAll", 0, address.Interface).
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

func openConnection(base context.Context, kind port.BusKind) (*dbus.Conn, error) {
	opt := dbus.WithContext(base)
	connect := dbus.ConnectSessionBus

	if kind == port.SystemBus {
		connect = dbus.ConnectSystemBus
	}

	conn, err := connect(opt)
	if err != nil {
		return nil, fmt.Errorf("dbusconn: connect bus: %w", err)
	}

	return conn, nil
}

// Close releases the private connections owned by this bus.
func (bus *connectionBus[K, C, O, Q]) Close() error {
	return errors.Join(releaseConnection(bus.source.Shutdown))
}

func releaseConnection(release func() error) error {
	err := release()
	if err != nil {
		return fmt.Errorf("dbusconn: close connections: %w", err)
	}

	return nil
}

// selectConnection binds a bus kind for operations that only need a connection provider.
func selectConnection[K, Value any](
	ctx context.Context,
	source interface {
		Get(ctx context.Context, kind K) (Value, error)
	},
	kind K,
) func() (Value, error) {
	return func() (Value, error) { return source.Get(ctx, kind) }
}

// Get selects or opens the connection for kind.
func (source connectionSource[K, C]) Get(ctx context.Context, kind K) (C, error) {
	result, err := source.Acquire(ctx, kind)

	return result, errors.Join(err)
}

// Shutdown releases all connections owned by the source.
func (source connectionSource[K, C]) Shutdown() error {
	return errors.Join(source.Release())
}
