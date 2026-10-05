// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package env

import (
	"os"

	"github.com/gostafa/linuxdesktop/internal/domain"
)

// New returns an environment probe.
func New() Probe { return snapshot }

// Snapshot reads the environment once.
func snapshot() domain.Env {
	var snapshot domain.Env

	readSession(&snapshot)
	readPaths(&snapshot)
	readDesktop(&snapshot)
	readRemote(&snapshot)

	return snapshot
}

func readSession(snapshot *domain.Env) {
	snapshot.SessionID = os.Getenv(keySessionID)
	snapshot.SessionType = os.Getenv(keySessionType)
	snapshot.SessionDesktop = os.Getenv(keySessionDesktop)
	snapshot.SessionClass = os.Getenv(keySessionClass)
	snapshot.CurrentDesktop = os.Getenv(keyCurrentDesktop)
	snapshot.DesktopSession = os.Getenv(keyDesktopSession)
	snapshot.User = firstNonEmpty(os.Getenv(keyUser), os.Getenv(keyLogname))
	snapshot.Seat = os.Getenv(keySeat)
	snapshot.VTNR = os.Getenv(keyVTNR)
}

func readPaths(snapshot *domain.Env) {
	snapshot.WaylandDisplay = os.Getenv(keyWaylandDisplay)
	snapshot.WaylandSocket = os.Getenv(keyWaylandSocket)
	snapshot.Display = os.Getenv(keyDisplay)
	snapshot.RuntimeDir = os.Getenv(keyRuntimeDir)
	snapshot.ConfigHome = os.Getenv(keyConfigHome)
	snapshot.Home = os.Getenv(keyHome)
}

func readDesktop(snapshot *domain.Env) {
	snapshot.HyprlandSignature = os.Getenv(keyHyprlandSignature)
	snapshot.SwaySock = os.Getenv(keySwaySock)
	snapshot.WayfireSocket = os.Getenv(keyWayfireSocket)
	snapshot.I3Sock = os.Getenv(keyI3Sock)
	snapshot.KDEFullSession = os.Getenv(keyKDEFullSession)
	snapshot.KDESessionVersion = os.Getenv(keyKDESessionVersion)
	snapshot.GNOMESessionID = os.Getenv(keyGNOMESessionID)
	snapshot.GNOMESetupDisplay = os.Getenv(keyGNOMESetupDisplay)
}

func readRemote(snapshot *domain.Env) {
	snapshot.SSHConnection = os.Getenv(keySSHConnection)
	snapshot.SSHTTY = os.Getenv(keySSHTTY)
	snapshot.SSHClient = os.Getenv(keySSHClient)
}

func firstNonEmpty(vals ...string) string {
	for i := range vals {
		if vals[i] != "" {
			return vals[i]
		}
	}

	return ""
}
