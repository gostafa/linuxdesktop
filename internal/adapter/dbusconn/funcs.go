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
		base = context.Background() //nolint:contextcheck // Nil context permits an unbounded connection lifetime.
	}

	cache := new(connections[port.BusKind, *dbus.Conn])

	cache.open = openConnection

	return &Bus{source: connectionSource[port.BusKind, *dbus.Conn]{
		Acquire: func(kind port.BusKind) (*dbus.Conn, error) { return connect(base, cache, kind) },
		Release: func() error { return closeConnections(cache) },
	}}
}

// HasOwner reports whether a well-known name currently has an owner.
func (bus *connectionBus[K, C, O, Q]) HasOwner(
	ctx context.Context,
	kind K,
	name string,
) (bool, error) {
	result, err := checkOwner(ctx, selectConnection(bus.source, kind), name)
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
	result, err := readIntrospection(ctx, selectConnection(bus.source, kind), object)

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
//
//nolint:ireturn // D-Bus properties have heterogeneous values required by the Bus port.
func (bus *connectionBus[K, C, O, Q]) Property(
	ctx context.Context,
	kind K,
	query Q,
) (any, error) {
	result, err := readProperty(ctx, selectConnection(bus.source, kind), query)
	if err != nil {
		return nil, fmt.Errorf(errReadProperty, err)
	}

	return result.Value(), nil
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
	result, err := readProperties(ctx, selectConnection(bus.source, kind), query)
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

// Close releases whichever connections were actually opened.
func closeConnections(cache *connections[port.BusKind, *dbus.Conn]) error {
	errs := make([]error, noEntries, len(cache.conns))

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
func connect(
	base context.Context,
	cache *connections[port.BusKind, *dbus.Conn],
	kind port.BusKind,
) (*dbus.Conn, error) {
	i := int(kind)
	if i >= len(cache.conns) {
		return nil, ErrBusKind
	}

	cache.once[i].Do(func() {
		cache.conns[i], cache.errs[i] = cache.open(base, kind)
	})

	if cache.errs[i] != nil {
		return nil, fmt.Errorf("dbusconn: cached connection: %w", cache.errs[i])
	}

	return cache.conns[i], nil
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
func selectConnection[K, C any](
	source interface{ Get(kind K) (C, error) },
	kind K,
) func() (C, error) {
	return func() (C, error) { return source.Get(kind) }
}

// Get selects or opens the connection for kind.
func (source connectionSource[K, C]) Get(kind K) (C, error) {
	result, err := source.Acquire(kind)

	return result, errors.Join(err)
}

// Shutdown releases all connections owned by the source.
func (source connectionSource[K, C]) Shutdown() error {
	return errors.Join(source.Release())
}
