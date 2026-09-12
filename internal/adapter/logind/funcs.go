package logind

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gostafa/linuxdesktop/internal/domain"
	"github.com/gostafa/linuxdesktop/internal/port"
	"github.com/gostafa/linuxdesktop/internal/sysfs"
)

// New returns a session probe. bus may be nil to disable the D-Bus fallback.
func New(bus port.Bus) *Probe { return &Probe{bus: bus} }

// Session reads the seat session, preferring the filesystem mirror and only
// falling back to D-Bus when it is unreadable.
func (p *Probe) Session(ctx context.Context, env *domain.Env) (domain.SessionInfo, error) {
	info := domain.SessionInfo{Type: domain.SessionTypeUnknown}

	id := env.SessionID
	if id == "" {
		id = idFromCgroup()
	}
	if id == "" {
		id = idFromScan(os.Getuid())
	}

	found := false
	if id != "" {
		if data, err := sysfs.Bytes(filepath.Join(dirSessions, id)); err == nil {
			info.ID = id
			applyFile(&info, data)
			found = true
		}
	}
	if !found && p.bus != nil {
		found = p.applyBus(ctx, &info) == nil
	}

	applyEnv(&info, env)

	if !found && info.Type == domain.SessionTypeUnknown && info.ID == "" {
		return info, ErrNoSession
	}
	return info, nil
}

// applyFile fills info from a /run/systemd/sessions/<id> mirror.
func applyFile(info *domain.SessionInfo, data []byte) {
	sysfs.Each(data, func(key, value string) bool {
		switch key {
		case keyUser:
			info.User = value
		case keySeat:
			info.Seat = value
		case keyType:
			info.Type = normalizeType(value)
		case keyDesktop:
			info.Desktop = value
		case keyService:
			info.Name = value
		case keyActive:
			info.Active = truthy(value)
		case keyRemote:
			info.Remote = truthy(value)
		case keyState:
			if value == "active" {
				info.Active = true
			}
		case keyVTNR:
			info.VTNumber, _ = strconv.Atoi(value)
		}
		return true
	})
}

// applyBus fills info from logind's own session object in one round trip.
func (p *Probe) applyBus(ctx context.Context, info *domain.SessionInfo) error {
	props, err := p.bus.Properties(ctx, port.SystemBus, busName, sessionPath, sessionIface)
	if err != nil {
		return err
	}
	if s, ok := props[propID].(string); ok {
		info.ID = s
	}
	if s, ok := props[propType].(string); ok {
		info.Type = normalizeType(s)
	}
	if s, ok := props[propDesktop].(string); ok {
		info.Desktop = s
	}
	if s, ok := props[propService].(string); ok {
		info.Name = s
	}
	if s, ok := props[propName].(string); ok {
		info.User = s
	}
	if b, ok := props[propRemote].(bool); ok {
		info.Remote = b
	}
	if b, ok := props[propActive].(bool); ok {
		info.Active = b
	}
	if n, ok := props[propVTNr].(uint32); ok {
		info.VTNumber = int(n)
	}
	info.Seat = structID(props[propSeat])
	return nil
}

// applyEnv backfills anything logind did not supply and always records the raw
// XDG variables, which callers frequently want to compare against.
func applyEnv(info *domain.SessionInfo, env *domain.Env) {
	if info.Type == domain.SessionTypeUnknown {
		info.Type = normalizeType(env.SessionType)
	}
	if info.ID == "" {
		info.ID = env.SessionID
	}
	if info.Desktop == "" {
		info.Desktop = env.SessionDesktop
	}
	if info.Seat == "" {
		info.Seat = env.Seat
	}
	if info.User == "" {
		info.User = env.User
	}
	if info.VTNumber == 0 && env.VTNR != "" {
		info.VTNumber, _ = strconv.Atoi(env.VTNR)
	}
	if env.SSHConnection != "" || env.SSHTTY != "" || env.SSHClient != "" {
		info.Remote = true
	}
	info.XDGSessionType = env.SessionType
	info.XDGSessionDesktop = env.SessionDesktop
}

// idFromCgroup extracts the id from a session-<id>.scope cgroup path, which is
// how libsystemd resolves a pid to a session without talking to logind.
func idFromCgroup() string {
	data, err := sysfs.Bytes(pathCgroup)
	if err != nil {
		return ""
	}
	i := bytes.Index(data, []byte(cgroupPrefix))
	if i < 0 {
		return ""
	}
	rest := data[i+len(cgroupPrefix):]
	j := bytes.Index(rest, []byte(cgroupSuffix))
	if j < 0 {
		return ""
	}
	return string(rest[:j])
}

// idFromScan looks for a session owned by uid, preferring an active one.
func idFromScan(uid int) string {
	names, err := sysfs.DirNames(dirSessions)
	if err != nil {
		return ""
	}
	want := strconv.Itoa(uid)
	fallback := ""
	for _, name := range names {
		if strings.ContainsRune(name, '.') {
			continue // .ref FIFOs, not session mirrors
		}
		data, err := sysfs.Bytes(filepath.Join(dirSessions, name))
		if err != nil || sysfs.Field(data, keyUID) != want {
			continue
		}
		if sysfs.Field(data, keyState) == "active" {
			return name
		}
		if fallback == "" {
			fallback = name
		}
	}
	return fallback
}

// normalizeType maps logind's TYPE= and $XDG_SESSION_TYPE onto the enum.
func normalizeType(s string) domain.SessionType {
	switch strings.ToLower(s) {
	case "wayland":
		return domain.SessionTypeWayland
	case "x11":
		return domain.SessionTypeX11
	case "tty":
		return domain.SessionTypeTTY
	case "mir":
		return domain.SessionTypeMir
	}
	return domain.SessionTypeUnknown
}

// structID pulls the leading string out of a D-Bus (string, objectpath) pair,
// the shape logind uses for Seat and User.
func structID(v any) string {
	fields, ok := v.([]any)
	if !ok || len(fields) == 0 {
		return ""
	}
	s, _ := fields[0].(string)
	return s
}

func truthy(s string) bool {
	return s == "1" || strings.EqualFold(s, "yes") || strings.EqualFold(s, "true")
}
