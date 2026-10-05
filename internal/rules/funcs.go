// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package rules

import (
	"strings"

	"github.com/gostafa/linuxdesktop/internal/domain"
)

// Desktop identifies the desktop environment from the XDG variables.
func Desktop(env *domain.Env) domain.DesktopInfo {
	info := domain.DesktopInfo{
		Environment:     domain.DesktopUnknown,
		CurrentDesktop:  env.CurrentDesktop,
		SessionDesktop:  env.SessionDesktop,
		DesktopSession:  env.DesktopSession,
		CurrentDesktops: nil,
		Name:            "",
		Version:         "",
	}
	if env.CurrentDesktop != noValue {
		info.CurrentDesktops = strings.Split(env.CurrentDesktop, desktopSeparator)
	}

	if fromTokens(&info, desktopTokenList(&info, env)) {
		return info
	}

	fromLegacy(&info, env)

	return info
}

// Compositor walks the detection ladder and returns the first answer it
// reaches, along with the method and confidence that produced it.
func Compositor(sig *domain.Signals) domain.CompositorInfo {
	info := unidentified(sig)

	rungs := ladder()
	for i := range rungs {
		if found, ok := rungs[i](sig); ok {
			apply(&info, &found)

			return info
		}
	}

	return info
}

// Protocol reports which display protocol a client should actually use. A
// reachable server beats what $XDG_SESSION_TYPE claims, because a session type
// the client cannot connect to is of no use to it.
func Protocol(
	sessionType domain.SessionType,
	waylandAvailable, x11Available bool,
) domain.DisplayProtocol {
	switch {
	case waylandAvailable:
		return domain.DisplayProtocolWayland
	case x11Available:
		return domain.DisplayProtocolX11
	case sessionType == domain.SessionTypeWayland:
		return domain.DisplayProtocolWayland
	case sessionType == domain.SessionTypeX11:
		return domain.DisplayProtocolX11
	}

	return domain.DisplayProtocolUnknown
}

// Headless reports whether there is no display server to draw on.
func Headless(waylandAvailable, x11Available bool) bool {
	return !waylandAvailable && !x11Available
}

// IsCompositorProcess reports whether a /proc comm value could name a
// compositor. The process scanner uses it to discard the several hundred
// uninteresting processes before paying to verify ownership of the rest.
func IsCompositorProcess(name string) bool {
	found, ok := processNames(name)

	return ok && found.name != noValue
}

// desktopTokenList is every token that could name a desktop, in the order the
// XDG variables are meant to be consulted.
func desktopTokenList(info *domain.DesktopInfo, env *domain.Env) []string {
	tokens := make([]string, zero, len(info.CurrentDesktops)+extraTokens)

	tokens = append(tokens, info.CurrentDesktops...)

	return append(tokens, env.SessionDesktop, env.DesktopSession)
}

// fromTokens records what the tokens say, reporting whether one of them named a
// known desktop.
func fromTokens(info *domain.DesktopInfo, tokens []string) bool {
	for i := range tokens {
		if matchToken(info, tokens[i]) {
			return true
		}
	}

	return false
}

// matchToken records one token's contribution: a recognized token settles the
// environment, and any non-empty token supplies a name if none was set yet.
func matchToken(info *domain.DesktopInfo, token string) bool {
	if token == noValue {
		return false
	}

	if de, ok := desktopTokens(normalize(token)); ok {
		info.Environment = de
		info.Name = token

		return true
	}

	nameIfUnset(info, token)

	return false
}

// fromLegacy falls back to the pre-XDG markers, still exported by both
// projects.
func fromLegacy(info *domain.DesktopInfo, env *domain.Env) {
	switch {
	case env.KDEFullSession != noValue:
		info.Environment = domain.DesktopKDE
		nameIfUnset(info, "KDE")
	case env.GNOMESessionID != noValue || env.GNOMESetupDisplay != noValue:
		info.Environment = domain.DesktopGNOME
		nameIfUnset(info, "GNOME")
	default:
		// No pre-XDG marker either; the desktop stays unknown.
	}
}

func nameIfUnset(info *domain.DesktopInfo, name string) {
	if info.Name == noValue {
		info.Name = name
	}
}

// unidentified is the answer when no rung of the ladder fires: nothing named,
// though the reachability of each protocol is still known.
func unidentified(sig *domain.Signals) domain.CompositorInfo {
	return domain.CompositorInfo{
		Kind:       domain.CompositorUnknown,
		Confidence: domain.ConfidenceUnknown,
		DetectedBy: domain.DetectedUnknown,
		Wayland:    sig.WaylandReachable,
		X11:        sig.X11Reachable,
		Name:       "",
		Version:    "",
	}
}

// ladder is the detection sequence, strongest signal first. Each rung answers
// from a different source, so the first one to fire is also the most
// trustworthy one available.
func ladder() []rung {
	return []rung{
		byEnv,
		byGlobals,
		byWindowManager,
		byDesktop,
		byWlroots,
		byProcesses,
	}
}

// byEnv reads the strongest signal there is: each of these variables is
// exported by exactly one compositor.
func byEnv(sig *domain.Signals) (verdict, bool) {
	envFingerprints := envFingerprints()

	var empty verdict

	for i := range envFingerprints {
		if envFingerprints[i].value(&sig.Env) == noValue {
			continue
		}

		found := envFingerprints[i].match

		return verdict{found, domain.ConfidenceHigh, domain.DetectedEnvironment}, true
	}

	return empty, false
}

// byGlobals identifies a compositor from an interface only it advertises.
func byGlobals(sig *domain.Signals) (verdict, bool) {
	waylandFingerprints := waylandFingerprints()

	var empty verdict

	for i := range waylandFingerprints {
		if !advertises(sig.WaylandGlobals, waylandFingerprints[i].iface) {
			continue
		}

		found := waylandFingerprints[i].match

		return verdict{found, domain.ConfidenceHigh, domain.DetectedWayland}, true
	}

	return empty, false
}

// byWindowManager reads the EWMH name the manager published about itself.
func byWindowManager(sig *domain.Signals) (verdict, bool) {
	var empty verdict

	name := sig.X11WindowManager
	if name == noValue {
		return empty, false
	}

	if found, ok := windowManagerNames(normalize(name)); ok {
		return verdict{found, domain.ConfidenceHigh, domain.DetectedX11EWMH}, true
	}

	// Unrecognized, but the manager did name itself, so keep what it said.
	said := match{domain.CompositorUnknown, name}

	return verdict{said, domain.ConfidenceMedium, domain.DetectedX11EWMH}, true
}

// byDesktop infers the compositor a desktop environment ships with, which is
// an inference rather than an observation and never better than medium.
func byDesktop(sig *domain.Signals) (verdict, bool) {
	var empty verdict

	found, ok := desktopCompositors(sig.Desktop)
	if !ok || found.name == noValue {
		return empty, false
	}

	return verdict{found, domain.ConfidenceMedium, domain.DetectedEnvironment}, true
}

// byWlroots narrows the field to the wlroots family without naming a member of
// it.
func byWlroots(sig *domain.Signals) (verdict, bool) {
	var empty verdict

	if !hasWlroots(sig.WaylandGlobals) {
		return empty, false
	}

	family := match{domain.CompositorUnknown, nameWlroots}

	return verdict{family, domain.ConfidenceLow, domain.DetectedWayland}, true
}

// byProcesses is the weakest signal: a matching process may belong to another
// seat entirely.
func byProcesses(sig *domain.Signals) (verdict, bool) {
	var empty verdict

	for i := range sig.Processes {
		found, ok := processNames(sig.Processes[i])
		if !ok {
			continue
		}

		return verdict{found, domain.ConfidenceLow, domain.DetectedProcess}, true
	}

	return empty, false
}

func hasWlroots(globals []domain.WaylandGlobal) bool {
	wlrootsMarkers := wlrootsMarkers()

	for i := range wlrootsMarkers {
		if advertises(globals, wlrootsMarkers[i]) {
			return true
		}
	}

	return false
}

// advertises reports whether a registry listing carries a named interface.
func advertises(globals []domain.WaylandGlobal, iface string) bool {
	for i := range globals {
		if globals[i].Interface == iface {
			return true
		}
	}

	return false
}

func apply(info *domain.CompositorInfo, found *verdict) {
	info.Kind = found.match.kind
	info.Name = found.match.name
	info.Confidence = found.confidence
	info.DetectedBy = found.method
}

// normalize lowercases and drops every non-alphanumeric byte, so that
// "X-Cinnamon", "Mutter (Muffin)" and "plasma-wayland" all reduce to a stable
// lookup key.
func normalize(raw string) string {
	out := make([]byte, zero, len(raw))
	for i := range len(raw) {
		out = appendKey(out, raw[i])
	}

	return string(out)
}

// appendKey keeps a byte that can appear in a lookup key, lowercasing letters
// on the way through and dropping everything else.
func appendKey(out []byte, char byte) []byte {
	if upper(char) {
		return append(out, char+('a'-'A'))
	}

	if lower(char) || digit(char) {
		return append(out, char)
	}

	return out
}

func upper(char byte) bool { return char >= 'A' && char <= 'Z' }

func lower(char byte) bool { return char >= 'a' && char <= 'z' }

func digit(char byte) bool { return char >= '0' && char <= '9' }
