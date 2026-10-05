// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package portal

import (
	"context"
	"encoding/xml"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gostafa/linuxdesktop/internal/domain"
	"github.com/gostafa/linuxdesktop/internal/port"
	"github.com/gostafa/linuxdesktop/internal/sysfs"
)

// New returns a portal probe. bus may be nil, which reports the portal as
// unavailable without attempting to connect.
func New(bus port.Bus) *Probe { return &Probe{bus: bus} }

// Portal reports whether xdg-desktop-portal is running and which of its
// interfaces the active backend exports.
func (probe *Probe) Portal(ctx context.Context, env *domain.Env) (domain.PortalInfo, error) {
	info := domain.PortalInfo{Backend: backend(env)}
	if probe.bus == nil {
		return info, nil
	}

	owned, err := probe.bus.HasOwner(ctx, port.SessionBus, busName)
	info.Available = owned

	if err != nil || !owned {
		return info, err
	}

	return probe.interfaces(ctx, info)
}

// interfaces asks the running portal what it exports, which is the only way to
// learn what the active backend actually implements.
func (probe *Probe) interfaces(
	ctx context.Context,
	info domain.PortalInfo,
) (domain.PortalInfo, error) {
	raw, err := probe.bus.Introspect(ctx, port.SessionBus, busName, objectPath)
	if err != nil {
		return info, err
	}

	var root node

	err = xml.Unmarshal([]byte(raw), &root)
	if err != nil {
		return info, err
	}

	applyInterfaces(&info, root.Interfaces)

	return info, nil
}

// applyInterfaces marks every portal interface the backend exports.
func applyInterfaces(info *domain.PortalInfo, ifaces []iface) {
	flags := interfaceFlags()
	for i := range ifaces {
		record(info, flags, ifaces[i].Name)
	}
}

// record marks the one interface a fully qualified name refers to. A name
// outside the portal namespace is some other interface on the same object.
func record(info *domain.PortalInfo, flags map[string]flagSetter, full string) {
	name, found := strings.CutPrefix(full, ifacePrefix)
	if !found {
		return
	}

	info.DesktopPortal = true

	set, ok := flags[name]
	if ok {
		set(info)
	}
}

// interfaceFlags maps each portal interface this library reports on onto the
// field that records it.
func interfaceFlags() map[string]flagSetter {
	return map[string]flagSetter{
		ifaceScreenCast:    func(info *domain.PortalInfo) { info.ScreenCast = true },
		ifaceScreenshot:    func(info *domain.PortalInfo) { info.Screenshot = true },
		ifaceFileChooser:   func(info *domain.PortalInfo) { info.FileChooser = true },
		ifaceOpenURI:       func(info *domain.PortalInfo) { info.OpenURI = true },
		ifaceRemoteDesktop: func(info *domain.PortalInfo) { info.RemoteDesktop = true },
		ifaceInhibit:       func(info *domain.PortalInfo) { info.Inhibit = true },
		ifaceNotification:  func(info *domain.PortalInfo) { info.Notification = true },
	}
}

// backend names the portal implementation that should be servicing this
// desktop, preferring explicit configuration over the legacy UseIn matching.
func backend(env *domain.Env) string {
	paths := configPaths(env)
	for i := range paths {
		if name := preferredDefault(paths[i]); name != noValue {
			return name
		}
	}

	return backendFromPortalFiles(env)
}

// configPaths lists portals.conf locations from most to least specific.
func configPaths(env *domain.Env) []string {
	paths := slices.Grow([]string(nil), expectedPaths)

	paths = append(paths, userConfig(env)...)

	tokens := desktopTokens(env)
	for i := range tokens {
		paths = append(paths, filepath.Join(shareDir, strings.ToLower(tokens[i])+configSuffix))
	}

	return append(paths, filepath.Join(shareDir, configName))
}

// userConfig is the per-user portals.conf, when the environment says where the
// user's configuration lives.
func userConfig(env *domain.Env) []string {
	if env.ConfigHome != noValue {
		return []string{filepath.Join(env.ConfigHome, portalDir, configName)}
	}

	if env.Home != noValue {
		return []string{filepath.Join(env.Home, userConfigDir, portalDir, configName)}
	}

	return nil
}

// preferredDefault reads the [preferred] default= entry, which is a
// semicolon-separated preference list.
func preferredDefault(path string) string {
	data, err := sysfs.Bytes(path)
	if err != nil {
		return noValue
	}

	for candidate := range strings.SplitSeq(sysfs.Field(data, keyDefault), listSeparator) {
		name := strings.TrimSpace(candidate)
		if named(name) {
			return name
		}
	}

	return noValue
}

// named reports whether a preference entry names a backend, rather than
// deferring the choice or refusing one.
func named(candidate string) bool {
	return candidate != noValue && candidate != anyBackend && candidate != noBackend
}

// backendFromPortalFiles matches the UseIn field of each installed .portal file
// against $XDG_CURRENT_DESKTOP. A lone installed backend wins by default.
func backendFromPortalFiles(env *domain.Env) string {
	names, err := sysfs.DirNames(portalsDir)
	if err != nil {
		return noValue
	}

	found := installed(names)
	wanted := desktopTokens(env)

	for i := range found {
		if servesAny(found[i].uses, wanted) {
			return found[i].name
		}
	}

	return onlyOne(found)
}

// installed reads every .portal file in the portals directory, skipping the
// ones that are unreadable or not portal files at all.
func installed(names []string) []backendFile {
	found := make([]backendFile, zero, len(names))

	for i := range names {
		file, ok := readPortalFile(names[i])
		if !ok {
			continue
		}

		found = append(found, file)
	}

	return found
}

// readPortalFile describes one .portal file.
func readPortalFile(name string) (backendFile, bool) {
	if !strings.HasSuffix(name, portalSuffix) {
		return backendFile{}, false
	}

	data, err := sysfs.Bytes(filepath.Join(portalsDir, name))
	if err != nil {
		return backendFile{}, false
	}

	file := backendFile{
		name: implName(strings.TrimSuffix(name, portalSuffix), data),
		uses: uses(data),
	}

	return file, true
}

// implName is the last element of the DBusName field, falling back to what the
// filename already said.
func implName(fallback string, data []byte) string {
	dbusName := sysfs.Field(data, keyDBusName)
	if dbusName == noValue {
		return fallback
	}

	last := afterLastDot(dbusName)
	if last == noValue {
		return fallback
	}

	return last
}

// afterLastDot is the final element of a dotted D-Bus name, or "" when there is
// no dot or nothing after it.
func afterLastDot(name string) string {
	dot := strings.LastIndexByte(name, dotByte)
	if dot < zero {
		return noValue
	}

	return name[dot+dotWidth:]
}

// uses is the desktop list a .portal file declares itself for.
func uses(data []byte) []string {
	var list []string

	for use := range strings.SplitSeq(sysfs.Field(data, keyUseIn), listSeparator) {
		trimmed := strings.TrimSpace(use)
		if trimmed != noValue {
			list = append(list, trimmed)
		}
	}

	return list
}

// servesAny reports whether a backend declares itself for any of the desktops
// currently running.
func servesAny(declared, wanted []string) bool {
	for i := range declared {
		if matchesAny(declared[i], wanted) {
			return true
		}
	}

	return false
}

func matchesAny(use string, wanted []string) bool {
	for i := range wanted {
		if strings.EqualFold(use, wanted[i]) {
			return true
		}
	}

	return false
}

// onlyOne names the single installed backend, which wins by default when none
// of them declared itself for this desktop.
func onlyOne(found []backendFile) string {
	name := noValue

	for i := range found {
		if name != noValue {
			return noValue // more than one is installed, so none is the default
		}

		name = found[i].name
	}

	return name
}

func desktopTokens(env *domain.Env) []string {
	if env.CurrentDesktop == noValue {
		return nil
	}

	return strings.Split(env.CurrentDesktop, desktopSeparator)
}
