// Package port declares the ports of the hexagon: the interfaces the
// detection engine depends on, expressed only in terms of the domain model.
//
// Nothing here imports a third-party package. In particular the D-Bus port is
// a narrow three-method interface rather than a *dbus.Conn, which keeps
// github.com/godbus/dbus out of every package except its own adapter.
//
// Each port is implemented by exactly one package under internal/adapter, and
// internal/core is the only consumer.
package port
