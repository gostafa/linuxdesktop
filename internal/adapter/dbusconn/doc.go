// Package dbusconn is the only package in the library that imports a D-Bus
// implementation.
//
// It adapts github.com/godbus/dbus/v5 to the narrow port.Bus interface, so the
// engine and the other adapters express their needs as "does this name have an
// owner", "list the interfaces on this object" and "read this property"
// without ever seeing a *dbus.Conn.
//
// Both buses are connected lazily and at most once per detection run, behind a
// sync.Once, and are shared by the session, portal and desktop probes. The
// connections are private rather than the process-wide ones returned by
// dbus.SessionBus, so Close actually releases them and the library never
// interferes with a connection the host application owns.
package dbusconn
