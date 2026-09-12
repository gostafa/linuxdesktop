// Package portal inspects xdg-desktop-portal.
//
// Which portal interfaces exist is answered with a single
// org.freedesktop.DBus.Introspectable.Introspect call. The alternative — a
// Properties.Get("version") per interface — costs seven round trips to learn
// the same thing, because an interface that is not exported simply is not in
// the introspection XML.
//
// Knowing the portal is running matters more than it looks: on Wayland it is
// the only sanctioned route to screenshots, screen capture and remote input,
// so a caller that finds ScreenCast false has learned it must take a different
// path entirely.
//
// The backend name is resolved from configuration rather than from D-Bus,
// since the running implementation does not advertise itself. Newer
// portals.conf files are consulted first, then the UseIn field of the
// installed .portal files.
package portal
