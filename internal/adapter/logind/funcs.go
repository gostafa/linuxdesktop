// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package logind

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gostafa/linuxdesktop/internal/domain"
	"github.com/gostafa/linuxdesktop/internal/port"
	"github.com/gostafa/linuxdesktop/internal/sysfs"
)

// New returns a session probe. bus may be nil to disable the D-Bus fallback.
func New(bus port.Bus) Probe { return sessionAt("", bus) }

// Session reads the seat session, preferring the filesystem mirror and only
// falling back to D-Bus when it is unreadable.
func sessionAt(files string, bus port.Bus) Probe {
	return func(ctx context.Context, env *domain.Env) (domain.SessionInfo, error) {
		var info domain.SessionInfo

		info.Type = domain.SessionTypeUnknown

		found := fromMirror(files, &info, sessionID(files, env)) || locate(ctx, &info, bus)

		applyEnv(&info, env)

		if found || known(&info) {
			return info, nil
		}

		return info, ErrNoSession
	}
}

// applyBus fills info from logind's own session object in one round trip.
func applyBus(ctx context.Context, info *domain.SessionInfo, bus port.Bus) error {
	props, err := bus.Properties(ctx, port.SystemBus, &port.PropertyQuery{
		Object:    port.Object{Destination: busName, Path: sessionPath},
		Interface: sessionIface,
		Name:      noValue,
	})
	if err != nil {
		return fmt.Errorf("logind: read session properties: %w", err)
	}

	readStrings(info, props)
	readFlags(info, props)

	info.Seat = structID(props[propSeat])

	return nil
}

// locate fills info from logind's session object when the mirror was unavailable.
func locate(
	ctx context.Context,
	info *domain.SessionInfo,
	bus port.Bus,
) bool {
	if bus == nil {
		return false
	}

	return applyBus(ctx, info, bus) == nil
}

// known reports whether anything at all was learned about the session.
func known(info *domain.SessionInfo) bool {
	return info.Type != domain.SessionTypeUnknown || info.ID != noValue
}

// sessionID resolves this process's session id: from the environment when it
// says, from the cgroup when it does not, and by scanning as a last resort.
func sessionID(files string, env *domain.Env) string {
	if env.SessionID != noValue {
		return env.SessionID
	}

	id := idFromCgroup(files)
	if id != noValue {
		return id
	}

	return idFromScan(files, os.Getuid())
}

// fromMirror fills info from a /run/systemd/sessions/<id> mirror, reporting
// whether there was one to read.
func fromMirror(files string, info *domain.SessionInfo, id string) bool {
	if id == noValue {
		return false
	}

	data, err := sysfs.Bytes(sysfs.Path(files, filepath.Join(dirSessions, id)))
	if err != nil {
		return false
	}

	info.ID = id
	applyFile(info, data)

	return true
}

// applyFile fills info from the KEY=VALUE body of a session mirror.
func applyFile(info *domain.SessionInfo, data []byte) {
	setters := fileSetters()

	sysfs.Each(data, func(key, value string) bool {
		set, ok := setters[key]
		if ok {
			set(info, value)
		}

		return true
	})
}

// fileSetters maps each key of a session mirror onto the field it fills. A key
// that is absent from the table is one this library has no use for.
func fileSetters() map[string]fileSetter {
	return map[string]fileSetter{
		keyUser:    setUser,
		keySeat:    setSeat,
		keyType:    setType,
		keyDesktop: setDesktop,
		keyService: setService,
		keyActive:  setActive,
		keyRemote:  setRemote,
		keyState:   setState,
		keyVTNR:    setVTNumber,
	}
}

func setUser(info *domain.SessionInfo, value string) { info.User = value }

func setSeat(info *domain.SessionInfo, value string) { info.Seat = value }

func setDesktop(info *domain.SessionInfo, value string) { info.Desktop = value }

func setService(info *domain.SessionInfo, value string) { info.Name = value }

func setType(info *domain.SessionInfo, value string) { info.Type = normalizeType(value) }

func setActive(info *domain.SessionInfo, value string) { info.Active = truthy(value) }

func setRemote(info *domain.SessionInfo, value string) { info.Remote = truthy(value) }

// setState records the session state, of which only "active" changes anything.
func setState(info *domain.SessionInfo, value string) {
	if value == stateActive {
		info.Active = true
	}
}

// setVTNumber records the virtual terminal. A value that is not a number means
// the session is not on a VT at all, which is what the zero already says.
func setVTNumber(info *domain.SessionInfo, value string) {
	number, err := strconv.Atoi(value)
	if err != nil {
		return
	}

	info.VTNumber = number
}

// readStrings copies the string-valued properties off the bus reply.
func readStrings(info *domain.SessionInfo, props map[string]any) {
	assignText(&info.ID, props[propID])
	assignText(&info.Desktop, props[propDesktop])
	assignText(&info.Name, props[propService])
	assignText(&info.User, props[propName])

	kind, ok := props[propType].(string)
	if ok {
		info.Type = normalizeType(kind)
	}
}

// readFlags copies the boolean and numeric properties off the same reply.
func readFlags(info *domain.SessionInfo, props map[string]any) {
	assignFlag(&info.Remote, props[propRemote])
	assignFlag(&info.Active, props[propActive])

	number, ok := props[propVTNr].(uint32)
	if ok {
		info.VTNumber = int(number)
	}
}

// assignText writes a string property into target, leaving it alone when
// logind sent nothing, or sent something that is not a string.
func assignText(target *string, value any) {
	text, ok := value.(string)
	if ok {
		*target = text
	}
}

// assignFlag writes a boolean property into target on the same terms.
func assignFlag(target *bool, value any) {
	flag, ok := value.(bool)
	if ok {
		*target = flag
	}
}

// applyEnv backfills anything logind did not supply and always records the raw
// XDG variables, which callers frequently want to compare against.
func applyEnv(info *domain.SessionInfo, env *domain.Env) {
	backfill(info, env)
	markRemote(info, env)

	info.XDGSessionType = env.SessionType
	info.XDGSessionDesktop = env.SessionDesktop
}

// backfill supplies from the environment whatever logind left empty.
func backfill(info *domain.SessionInfo, env *domain.Env) {
	if info.Type == domain.SessionTypeUnknown {
		info.Type = normalizeType(env.SessionType)
	}

	info.ID = orEnv(info.ID, env.SessionID)
	info.Desktop = orEnv(info.Desktop, env.SessionDesktop)
	info.Seat = orEnv(info.Seat, env.Seat)
	info.User = orEnv(info.User, env.User)

	if info.VTNumber == zero {
		setVTNumber(info, env.VTNR)
	}
}

// orEnv keeps what logind said, falling back to the environment when it said
// nothing.
func orEnv(known, fallback string) string {
	if known != noValue {
		return known
	}

	return fallback
}

// markRemote records an SSH login, which logind does not always flag itself.
func markRemote(info *domain.SessionInfo, env *domain.Env) {
	ssh := env.SSHConnection != noValue ||
		env.SSHTTY != noValue ||
		env.SSHClient != noValue

	info.Remote = info.Remote || ssh
}

// idFromCgroup extracts the id from a session-<id>.scope cgroup path, which is
// how libsystemd resolves a pid to a session without talking to logind.
func idFromCgroup(files string) string {
	data, err := sysfs.Bytes(sysfs.Path(files, pathCgroup))
	if err != nil {
		return noValue
	}

	return between(data, cgroupPrefix, cgroupSuffix)
}

// between lifts out the text bracketed by two markers, yielding "" when either
// of them is missing.
func between(data []byte, prefix, suffix string) string {
	start := bytes.Index(data, []byte(prefix))
	if start < zero {
		return noValue
	}

	rest := data[start+len(prefix):]

	parts := bytes.SplitN(rest, []byte(suffix), bracketParts)
	if len(parts) != bracketParts {
		return noValue
	}

	return string(parts[zero])
}

// idFromScan looks for a session owned by uid, preferring an active one.
func idFromScan(files string, uid int) string {
	names, err := sysfs.DirNames(sysfs.Path(files, dirSessions))
	if err != nil {
		return noValue
	}

	return pickSession(files, names, strconv.Itoa(uid))
}

// pickSession is the first active session owned by uid, or failing that the
// first session owned by uid at all.
func pickSession(files string, names []string, uid string) string {
	fallback := noValue

	for i := range names {
		state, ok := sessionState(files, names[i], uid)
		if !ok {
			continue
		}

		if state == stateActive {
			return names[i]
		}

		fallback = orEnv(fallback, names[i])
	}

	return fallback
}

// sessionState reads the STATE of one session mirror, reporting false when the
// entry is not a mirror or does not belong to uid.
func sessionState(files, name, uid string) (string, bool) {
	if strings.ContainsRune(name, '.') {
		return noValue, false // .ref FIFOs, not session mirrors
	}

	data, err := sysfs.Bytes(sysfs.Path(files, filepath.Join(dirSessions, name)))
	if err != nil || sysfs.Field(data, keyUID) != uid {
		return noValue, false
	}

	return sysfs.Field(data, keyState), true
}

// normalizeType maps logind's TYPE= and $XDG_SESSION_TYPE onto the enum.
func normalizeType(raw string) domain.SessionType {
	switch strings.ToLower(raw) {
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
func structID(value any) string {
	fields, ok := value.([]any)
	if !ok || len(fields) == zero {
		return noValue
	}

	text, ok := fields[zero].(string)
	if !ok {
		return noValue
	}

	return text
}

func truthy(raw string) bool {
	return raw == "1" || strings.EqualFold(raw, "yes") || strings.EqualFold(raw, "true")
}
