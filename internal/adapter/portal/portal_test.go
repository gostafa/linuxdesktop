// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package portal

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/gostafa/linuxdesktop/internal/domain"
	"github.com/gostafa/linuxdesktop/internal/port"
	"github.com/gostafa/linuxdesktop/internal/sysfs"
)

type fakeBus struct {
	port.Bus
	owned              bool
	ownerErr, errorXML error
	xml                string
}

func (bus fakeBus) HasOwner(context.Context, port.BusKind, string) (bool, error) {
	return bus.owned, bus.ownerErr
}

func (bus fakeBus) Introspect(context.Context, port.BusKind, *port.Object) (string, error) {
	return bus.xml, bus.errorXML
}

func write(t *testing.T, files string, path, data string) {
	t.Helper()
	path = sysfs.Path(files, path)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestPortalInterfaces(t *testing.T) {
	t.Parallel()
	files := t.TempDir()
	env := new(domain.Env)
	if info, err := New(nil).Portal(t.Context(), env); err != nil || info.Available {
		t.Fatal(info, err)
	}
	failure := errors.New("bus unavailable")
	for _, bus := range []fakeBus{{}, {ownerErr: failure}, {owned: true, errorXML: failure}, {owned: true, xml: "<"}} {
		info, err := portalAt(files, bus).Portal(t.Context(), env)
		if bus.ownerErr != nil || bus.errorXML != nil {
			if !errors.Is(err, failure) {
				t.Fatal(err)
			}
		} else if bus.owned {
			if err == nil {
				t.Fatal("malformed XML accepted")
			}
		} else if err != nil || info.Available {
			t.Fatal(info, err)
		}
	}
	xml := "<node><interface name='unrelated'/><interface name='org.freedesktop.portal.Unknown'/>"
	for name := range interfaceFlags() {
		xml += "<interface name='" + ifacePrefix + name + "'/>"
	}
	xml += "</node>"
	info, err := portalAt(files, fakeBus{owned: true, xml: xml}).Portal(t.Context(), env)
	if err != nil || !info.Available || !info.DesktopPortal || !info.ScreenCast ||
		!info.Screenshot ||
		!info.FileChooser ||
		!info.OpenURI ||
		!info.RemoteDesktop ||
		!info.Inhibit ||
		!info.Notification {
		t.Fatal(info, err)
	}
}

func TestBackendPreferenceAndUseIn(t *testing.T) {
	t.Parallel()
	files := t.TempDir()
	env := domain.Env{CurrentDesktop: "GNOME:Custom", ConfigHome: "/user/config"}
	write(
		t,
		files,
		filepath.Join(env.ConfigHome, portalDir, configName),
		"[preferred]\ndefault=*;none; gtk;\n",
	)
	if got := backend(files, &env); got != "gtk" {
		t.Fatal(got)
	}
	if len(userConfig(&domain.Env{Home: "/home/alice"})) != 1 || userConfig(&domain.Env{}) != nil {
		t.Fatal("user paths")
	}
	write(t, files, filepath.Join(env.ConfigHome, portalDir, configName), "default=*;none;\n")
	write(
		t,
		files,
		filepath.Join(portalsDir, "gtk.portal"),
		"DBusName=org.freedesktop.impl.portal.desktop.gtk\nUseIn=gnome;\n",
	)
	write(t, files, filepath.Join(portalsDir, "ignored.txt"), "")
	write(t, files, filepath.Join(portalsDir, "other.portal"), "UseIn=Other;\n")
	if got := backend(files, &env); got != "gtk" {
		t.Fatal(got)
	}
	if got := backendFromPortalFiles(files, &domain.Env{}); got != "" {
		t.Fatal(got)
	}
	if _, ok := readPortalFile(files, "missing.portal"); ok {
		t.Fatal("missing portal")
	}
	if implName("fallback", []byte("DBusName=nodots")) != "fallback" ||
		implName("fallback", []byte("DBusName=name.")) != "fallback" {
		t.Fatal("name fallback")
	}
	if onlyOne([]backendFile{{name: "only"}}) != "only" || onlyOne(nil) != "" {
		t.Fatal("single backend")
	}
}
