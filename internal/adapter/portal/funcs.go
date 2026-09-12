package portal

import (
	"context"
	"encoding/xml"
	"path/filepath"
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
func (p *Probe) Portal(ctx context.Context, env *domain.Env) (domain.PortalInfo, error) {
	info := domain.PortalInfo{Backend: backend(env)}
	if p.bus == nil {
		return info, nil
	}

	owned, err := p.bus.HasOwner(ctx, port.SessionBus, busName)
	if err != nil {
		return info, err
	}
	info.Available = owned
	if !owned {
		return info, nil
	}

	raw, err := p.bus.Introspect(ctx, port.SessionBus, busName, objectPath)
	if err != nil {
		return info, err
	}

	var root node
	if err := xml.Unmarshal([]byte(raw), &root); err != nil {
		return info, err
	}
	for _, i := range root.Interfaces {
		name, found := strings.CutPrefix(i.Name, ifacePrefix)
		if !found {
			continue
		}
		info.DesktopPortal = true
		switch name {
		case ifaceScreenCast:
			info.ScreenCast = true
		case ifaceScreenshot:
			info.Screenshot = true
		case ifaceFileChooser:
			info.FileChooser = true
		case ifaceOpenURI:
			info.OpenURI = true
		case ifaceRemoteDesktop:
			info.RemoteDesktop = true
		case ifaceInhibit:
			info.Inhibit = true
		case ifaceNotification:
			info.Notification = true
		}
	}
	return info, nil
}

// backend names the portal implementation that should be servicing this
// desktop, preferring explicit configuration over the legacy UseIn matching.
func backend(env *domain.Env) string {
	for _, path := range configPaths(env) {
		if name := preferredDefault(path); name != "" {
			return name
		}
	}
	return backendFromPortalFiles(env)
}

// configPaths lists portals.conf locations from most to least specific.
func configPaths(env *domain.Env) []string {
	paths := make([]string, 0, 4)
	if env.ConfigHome != "" {
		paths = append(paths, filepath.Join(env.ConfigHome, "xdg-desktop-portal", configName))
	} else if env.Home != "" {
		paths = append(paths, filepath.Join(env.Home, ".config", "xdg-desktop-portal", configName))
	}
	for _, desktop := range desktopTokens(env) {
		paths = append(paths, filepath.Join(shareDir, strings.ToLower(desktop)+configSuffix))
	}
	return append(paths, filepath.Join(shareDir, configName))
}

// preferredDefault reads the [preferred] default= entry, which is a
// semicolon-separated preference list.
func preferredDefault(path string) string {
	data, err := sysfs.Bytes(path)
	if err != nil {
		return ""
	}
	for _, candidate := range strings.Split(sysfs.Field(data, keyDefault), listSeparator) {
		candidate = strings.TrimSpace(candidate)
		if candidate != "" && candidate != "*" && candidate != "none" {
			return candidate
		}
	}
	return ""
}

// backendFromPortalFiles matches the UseIn field of each installed .portal
// file against $XDG_CURRENT_DESKTOP. A lone installed backend wins by default.
func backendFromPortalFiles(env *domain.Env) string {
	names, err := sysfs.DirNames(portalsDir)
	if err != nil {
		return ""
	}
	wanted := desktopTokens(env)
	only := ""
	count := 0

	for _, name := range names {
		if !strings.HasSuffix(name, portalSuffix) {
			continue
		}
		data, err := sysfs.Bytes(filepath.Join(portalsDir, name))
		if err != nil {
			continue
		}
		impl := strings.TrimSuffix(name, portalSuffix)
		if dbusName := sysfs.Field(data, keyDBusName); dbusName != "" {
			if i := strings.LastIndexByte(dbusName, '.'); i >= 0 && i+1 < len(dbusName) {
				impl = dbusName[i+1:]
			}
		}
		count++
		only = impl

		for _, use := range strings.Split(sysfs.Field(data, keyUseIn), listSeparator) {
			use = strings.TrimSpace(use)
			for _, want := range wanted {
				if use != "" && strings.EqualFold(use, want) {
					return impl
				}
			}
		}
	}

	if count == 1 {
		return only
	}
	return ""
}

func desktopTokens(env *domain.Env) []string {
	if env.CurrentDesktop == "" {
		return nil
	}
	return strings.Split(env.CurrentDesktop, ":")
}
