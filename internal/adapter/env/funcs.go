// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package env

import (
	"os"

	"github.com/gostafa/linuxdesktop/internal/domain"
)

// New returns an environment probe.
func New() Probe { return Probe{} }

// Snapshot reads the environment once.
func (Probe) Snapshot() domain.Env {
	return domain.Env{
		SessionID:      os.Getenv(keySessionID),
		SessionType:    os.Getenv(keySessionType),
		SessionDesktop: os.Getenv(keySessionDesktop),
		SessionClass:   os.Getenv(keySessionClass),
		CurrentDesktop: os.Getenv(keyCurrentDesktop),
		DesktopSession: os.Getenv(keyDesktopSession),
		WaylandDisplay: os.Getenv(keyWaylandDisplay),
		WaylandSocket:  os.Getenv(keyWaylandSocket),
		Display:        os.Getenv(keyDisplay),
		RuntimeDir:     os.Getenv(keyRuntimeDir),
		ConfigHome:     os.Getenv(keyConfigHome),
		Home:           os.Getenv(keyHome),
		User:           firstNonEmpty(os.Getenv(keyUser), os.Getenv(keyLogname)),
		Seat:           os.Getenv(keySeat),
		VTNR:           os.Getenv(keyVTNR),

		HyprlandSignature: os.Getenv(keyHyprlandSignature),
		SwaySock:          os.Getenv(keySwaySock),
		WayfireSocket:     os.Getenv(keyWayfireSocket),
		I3Sock:            os.Getenv(keyI3Sock),
		KDEFullSession:    os.Getenv(keyKDEFullSession),
		KDESessionVersion: os.Getenv(keyKDESessionVersion),
		GNOMESessionID:    os.Getenv(keyGNOMESessionID),
		GNOMESetupDisplay: os.Getenv(keyGNOMESetupDisplay),

		SSHConnection: os.Getenv(keySSHConnection),
		SSHTTY:        os.Getenv(keySSHTTY),
		SSHClient:     os.Getenv(keySSHClient),
	}
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}

	return ""
}
