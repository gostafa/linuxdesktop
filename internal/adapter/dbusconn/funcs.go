// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package dbusconn

import (
	"context"
	"errors"

	"github.com/godbus/dbus/v5"
	"github.com/gostafa/linuxdesktop/internal/port"
)

// New returns a Bus whose connections live no longer than base.
func New(base context.Context) *Bus {
	if base == nil {
		base = context.Background()
	}

	return &Bus{base: base}
}

// HasOwner reports whether a well-known name currently has an owner.
func (b *Bus) HasOwner(ctx context.Context, k port.BusKind, name string) (bool, error) {
	c, err := b.conn(k)
	if err != nil {
		return false, err
	}

	var has bool

	err = c.BusObject().
		CallWithContext(ctx, "org.freedesktop.DBus.NameHasOwner", 0, name).
		Store(&has)

	return has, err
}

// Introspect returns the raw introspection XML for an object path.
func (b *Bus) Introspect(ctx context.Context, k port.BusKind, dest, path string) (string, error) {
	c, err := b.conn(k)
	if err != nil {
		return "", err
	}

	var xml string

	err = c.Object(dest, dbus.ObjectPath(path)).
		CallWithContext(ctx, "org.freedesktop.DBus.Introspectable.Introspect", 0).
		Store(&xml)

	return xml, err
}

// Property reads a single property off an interface.
func (b *Bus) Property(
	ctx context.Context,
	k port.BusKind,
	dest, path, iface, prop string,
) (any, error) {
	c, err := b.conn(k)
	if err != nil {
		return nil, err
	}

	var v dbus.Variant

	err = c.Object(dest, dbus.ObjectPath(path)).
		CallWithContext(ctx, "org.freedesktop.DBus.Properties.Get", 0, iface, prop).
		Store(&v)
	if err != nil {
		return nil, err
	}

	return v.Value(), nil
}

// Properties reads every property on an interface in one round trip, which is
// what makes reading a logind session cost one message instead of nine.
func (b *Bus) Properties(
	ctx context.Context,
	k port.BusKind,
	dest, path, iface string,
) (map[string]any, error) {
	c, err := b.conn(k)
	if err != nil {
		return nil, err
	}

	var raw map[string]dbus.Variant

	err = c.Object(dest, dbus.ObjectPath(path)).
		CallWithContext(ctx, "org.freedesktop.DBus.Properties.GetAll", 0, iface).
		Store(&raw)
	if err != nil {
		return nil, err
	}

	out := make(map[string]any, len(raw))
	for k, v := range raw {
		out[k] = v.Value()
	}

	return out, nil
}

// Close releases whichever connections were actually opened.
func (b *Bus) Close() error {
	var errs []error

	for i := range b.conns {
		if c := b.conns[i]; c != nil {
			err := c.Close()
			if err != nil {
				errs = append(errs, err)
			}

			b.conns[i] = nil
		}
	}

	return errors.Join(errs...)
}

// conn establishes a bus on first use and caches the outcome, including the
// failure. A machine without a session bus should pay for exactly one failed
// connect attempt, not one per probe.
func (b *Bus) conn(k port.BusKind) (*dbus.Conn, error) {
	i := int(k)
	if i < 0 || i >= len(b.conns) {
		return nil, errors.New("linuxdesktop: unknown bus kind")
	}

	b.once[i].Do(func() {
		opt := dbus.WithContext(b.base)

		if k == port.SystemBus {
			b.conns[i], b.errs[i] = dbus.ConnectSystemBus(opt)

			return
		}

		b.conns[i], b.errs[i] = dbus.ConnectSessionBus(opt)
	})

	return b.conns[i], b.errs[i]
}
