package desktop

import (
	"context"

	"github.com/gostafa/linuxdesktop/internal/domain"
	"github.com/gostafa/linuxdesktop/internal/port"
	"github.com/gostafa/linuxdesktop/internal/rules"
)

// New returns a desktop probe. bus may be nil to skip the version lookup.
func New(bus port.Bus) *Probe { return &Probe{bus: bus} }

// Desktop identifies the desktop environment and, where one is cheaply
// available, its version.
func (p *Probe) Desktop(ctx context.Context, env *domain.Env) (domain.DesktopInfo, error) {
	info := rules.Desktop(env)
	info.Version = p.version(ctx, env, info.Environment)
	return info, nil
}

func (p *Probe) version(ctx context.Context, env *domain.Env, de domain.DesktopEnvironment) string {
	switch de {
	case domain.DesktopGNOME:
		if p.bus == nil {
			return ""
		}
		v, err := p.bus.Property(ctx, port.SessionBus, shellName, shellPath, shellIface, shellVersion)
		if err != nil {
			return ""
		}
		s, _ := v.(string)
		return s
	case domain.DesktopKDE:
		// Major version only; Plasma exposes nothing finer without a subprocess.
		return env.KDESessionVersion
	}
	return ""
}
