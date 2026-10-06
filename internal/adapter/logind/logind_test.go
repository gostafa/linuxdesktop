// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package logind

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/gostafa/linuxdesktop/internal/domain"
	"github.com/gostafa/linuxdesktop/internal/port"
	"github.com/gostafa/linuxdesktop/internal/sysfs"
)

type fakeBus struct {
	port.Bus
	props map[string]any
	err   error
}

func (bus fakeBus) Properties(
	context.Context,
	port.BusKind,
	*port.PropertyQuery,
) (map[string]any, error) {
	return bus.props, bus.err
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

func TestSessionSources(t *testing.T) {
	t.Parallel()
	files := t.TempDir()
	if _, err := sessionAt(
		files,
		nil,
	).Session(t.Context(), &domain.Env{}); !errors.Is(
		err,
		ErrNoSession,
	) {
		t.Fatal(err)
	}
	write(
		t,
		files,
		filepath.Join(dirSessions, "7"),
		"UID=0\nUSER=alice\nSEAT=seat0\nTYPE=wayland\nDESKTOP=GNOME\nSERVICE=gdm\nACTIVE=yes\nREMOTE=true\nSTATE=active\nVTNR=2\nIGNORED=x\n",
	)
	info, err := sessionAt(files, nil).Session(t.Context(), &domain.Env{SessionID: "7"})
	if err != nil || info.ID != "7" || info.User != "alice" ||
		info.Type != domain.SessionTypeWayland ||
		!info.Active ||
		!info.Remote ||
		info.VTNumber != 2 ||
		info.Seat != "seat0" ||
		info.Name != "gdm" ||
		info.Desktop != "GNOME" {
		t.Fatal(info, err)
	}
	write(t, files, pathCgroup, "0::/user.slice/session-7.scope\n")
	if sessionID(files, &domain.Env{}) != "7" {
		t.Fatal("cgroup session")
	}
	write(t, files, pathCgroup, "0::/other\n")
	write(
		t,
		files,
		filepath.Join(dirSessions, "8"),
		"UID="+strconv.Itoa(os.Getuid())+"\nSTATE=active\n",
	)
	if sessionID(files, &domain.Env{}) != "8" {
		t.Fatal("scanned session")
	}
	write(t, files, filepath.Join(dirSessions, "9"), "UID=42\nSTATE=online\n")
	if pickSession(files, []string{"missing", "7.ref", "9"}, "42") != "9" {
		t.Fatal("inactive fallback")
	}
	if state, ok := sessionState(files, "9", "43"); ok || state != "" {
		t.Fatal(state, ok)
	}
	if _, err = New(nil).Session(t.Context(), &domain.Env{SessionID: "known"}); err != nil {
		t.Fatal(err)
	}
}

func TestBusAndEnvironment(t *testing.T) {
	t.Parallel()
	files := t.TempDir()
	props := map[string]any{
		propID:      "bus-id",
		propDesktop: "KDE",
		propService: "sddm",
		propName:    "alice",
		propType:    "x11",
		propRemote:  true,
		propActive:  true,
		propVTNr:    uint32(3),
		propSeat:    []any{"seat0", "/seat"},
	}
	info, err := sessionAt(
		files,
		fakeBus{props: props},
	).Session(t.Context(), &domain.Env{SessionType: "tty", VTNR: "bad"})
	if err != nil || info.ID != "bus-id" || info.Seat != "seat0" ||
		info.Type != domain.SessionTypeX11 ||
		info.VTNumber != 3 ||
		!info.Remote ||
		!info.Active {
		t.Fatal(info, err)
	}
	failure := errors.New("offline")
	if err = applyBus(t.Context(), &info, fakeBus{err: failure}); !errors.Is(err, failure) {
		t.Fatal(err)
	}
	info = domain.SessionInfo{Type: domain.Unknown}
	applyEnv(
		&info,
		&domain.Env{
			SessionID:      "env-id",
			SessionType:    "tty",
			SessionDesktop: "custom",
			User:           "bob",
			Seat:           "seat1",
			VTNR:           "bad",
			SSHClient:      "remote",
		},
	)
	if info.ID != "env-id" || info.Type != domain.SessionTypeTTY || !info.Remote ||
		info.VTNumber != 0 {
		t.Fatal(info)
	}
	setState(&info, "online")
	setVTNumber(&info, "4")
	if info.VTNumber != 4 {
		t.Fatal(info)
	}
	for _, raw := range []string{"mir", "unknown"} {
		if normalizeType(raw) == domain.SessionTypeX11 {
			t.Fatal(raw)
		}
	}
	for _, value := range []any{nil, []any{}, []any{42}, []any{"seat"}} {
		_ = structID(value)
	}
	if structID([]any{"seat"}) != "seat" ||
		between([]byte("session-7"), cgroupPrefix, cgroupSuffix) != "" ||
		between([]byte("none"), cgroupPrefix, cgroupSuffix) != "" {
		t.Fatal("pair parsing")
	}
}
