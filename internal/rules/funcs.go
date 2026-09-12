package rules

import (
	"strings"

	"github.com/gostafa/linuxdesktop/internal/domain"
)

// Desktop identifies the desktop environment from the XDG variables.
func Desktop(env *domain.Env) domain.DesktopInfo {
	info := domain.DesktopInfo{
		Environment:    domain.DesktopUnknown,
		CurrentDesktop: env.CurrentDesktop,
		SessionDesktop: env.SessionDesktop,
		DesktopSession: env.DesktopSession,
	}
	if env.CurrentDesktop != "" {
		info.CurrentDesktops = strings.Split(env.CurrentDesktop, desktopSeparator)
	}

	tokens := make([]string, 0, len(info.CurrentDesktops)+2)
	tokens = append(tokens, info.CurrentDesktops...)
	tokens = append(tokens, env.SessionDesktop, env.DesktopSession)

	for _, token := range tokens {
		if token == "" {
			continue
		}
		if de, ok := desktopTokens[normalize(token)]; ok {
			info.Environment = de
			info.Name = token
			return info
		}
		if info.Name == "" {
			info.Name = token
		}
	}

	// Pre-XDG markers, still exported by both projects.
	switch {
	case env.KDEFullSession != "":
		info.Environment = domain.DesktopKDE
		if info.Name == "" {
			info.Name = "KDE"
		}
	case env.GNOMESessionID != "" || env.GNOMESetupDisplay != "":
		info.Environment = domain.DesktopGNOME
		if info.Name == "" {
			info.Name = "GNOME"
		}
	}
	return info
}

// Compositor walks the detection ladder and returns the first answer it
// reaches, along with the method and confidence that produced it.
func Compositor(sig *domain.Signals) domain.CompositorInfo {
	info := domain.CompositorInfo{
		Kind:       domain.CompositorUnknown,
		Confidence: domain.ConfidenceUnknown,
		DetectedBy: domain.DetectedUnknown,
		Wayland:    sig.WaylandReachable,
		X11:        sig.X11Reachable,
	}

	if m, ok := fromEnv(&sig.Env); ok {
		return apply(info, m, domain.ConfidenceHigh, domain.DetectedEnvironment)
	}
	if m, ok := fromGlobals(sig.WaylandGlobals); ok {
		return apply(info, m, domain.ConfidenceHigh, domain.DetectedWayland)
	}
	if name := sig.X11WindowManager; name != "" {
		if m, ok := windowManagerNames[normalize(name)]; ok {
			return apply(info, m, domain.ConfidenceHigh, domain.DetectedX11EWMH)
		}
		// Unrecognised, but the manager did name itself, so keep what it said.
		return apply(info, match{domain.CompositorUnknown, name},
			domain.ConfidenceMedium, domain.DetectedX11EWMH)
	}
	if m, ok := desktopCompositors[sig.Desktop]; ok {
		return apply(info, m, domain.ConfidenceMedium, domain.DetectedEnvironment)
	}
	if hasWlroots(sig.WaylandGlobals) {
		return apply(info, match{domain.CompositorUnknown, nameWlroots},
			domain.ConfidenceLow, domain.DetectedWayland)
	}
	if m, ok := fromProcesses(sig.Processes); ok {
		return apply(info, m, domain.ConfidenceLow, domain.DetectedProcess)
	}
	return info
}

// Protocol reports which display protocol a client should actually use. A
// reachable server beats what $XDG_SESSION_TYPE claims, because a session type
// the client cannot connect to is of no use to it.
func Protocol(sessionType domain.SessionType, waylandAvailable, x11Available bool) domain.DisplayProtocol {
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

func fromEnv(env *domain.Env) (match, bool) {
	for _, f := range envFingerprints {
		if f.value(env) != "" {
			return f.match, true
		}
	}
	return match{}, false
}

func fromGlobals(globals []domain.WaylandGlobal) (match, bool) {
	if len(globals) == 0 {
		return match{}, false
	}
	for _, f := range waylandFingerprints {
		for i := range globals {
			if globals[i].Interface == f.iface {
				return f.match, true
			}
		}
	}
	return match{}, false
}

func hasWlroots(globals []domain.WaylandGlobal) bool {
	for _, marker := range wlrootsMarkers {
		for i := range globals {
			if globals[i].Interface == marker {
				return true
			}
		}
	}
	return false
}

// IsCompositorProcess reports whether a /proc comm value could name a
// compositor. The process scanner uses it to discard the several hundred
// uninteresting processes before paying to verify ownership of the rest.
func IsCompositorProcess(name string) bool {
	_, ok := processNames[name]
	return ok
}

func fromProcesses(names []string) (match, bool) {
	for _, name := range names {
		if m, ok := processNames[name]; ok {
			return m, true
		}
	}
	return match{}, false
}

func apply(info domain.CompositorInfo, m match, c domain.DetectionConfidence, d domain.DetectionMethod) domain.CompositorInfo {
	info.Kind = m.kind
	info.Name = m.name
	info.Confidence = c
	info.DetectedBy = d
	return info
}

// normalize lowercases and drops every non-alphanumeric byte, so that
// "X-Cinnamon", "Mutter (Muffin)" and "plasma-wayland" all reduce to a stable
// lookup key.
func normalize(s string) string {
	b := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case c >= 'A' && c <= 'Z':
			b = append(b, c+('a'-'A'))
		case c >= 'a' && c <= 'z', c >= '0' && c <= '9':
			b = append(b, c)
		}
	}
	return string(b)
}
