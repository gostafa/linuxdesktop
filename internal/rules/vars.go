// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package rules

import (
	"github.com/gostafa/linuxdesktop/internal/domain"
)

// desktopTokens maps a normalised $XDG_CURRENT_DESKTOP, $XDG_SESSION_DESKTOP
// or $DESKTOP_SESSION token onto a desktop environment. Tokens are matched in
// the order they appear in the variable, so "Budgie:GNOME" resolves to Budgie
// and "ubuntu:GNOME" to GNOME.
func desktopTokens(key string) (value domain.DesktopEnvironment, ok bool) {
	switch key {
	case "gnome",
		"gnomeclassic",
		"gnomeflashback",
		"gnomexorg",
		"gnomewayland",
		nameGnomeshell,
		"ubuntu",
		"ubuntuwayland",
		"ubuntuxorg",
		"pop",
		"zorin":
		return domain.DesktopGNOME, true
	case "kde", "plasma", "plasma5", "plasma6", "plasmawayland", "plasmax11":
		return domain.DesktopKDE, true
	case "xfce", "xfce4", "xubuntu":
		return domain.DesktopXFCE, true
	case "xcinnamon", "cinnamon":
		return domain.DesktopCinnamon, true
	case "mate":
		return domain.DesktopMATE, true
	case "lxqt":
		return domain.DesktopLXQt, true
	case "lxde":
		return domain.DesktopLXDE, true
	case "budgie", "budgiedesktop":
		return domain.DesktopBudgie, true
	case "cosmic":
		return domain.DesktopCOSMIC, true
	case "pantheon", "elementary":
		return domain.DesktopPantheon, true
	default:
		return value, false
	}
}

// envFingerprints are the strongest signal available: each of these variables
// is exported by exactly one compositor.
func envFingerprints() []struct {
	value func(*domain.Env) string
	match
} {
	return []struct {
		value func(*domain.Env) string
		match
	}{
		{
			func(e *domain.Env) string { return e.HyprlandSignature },
			match{domain.CompositorHyprland, nameHyprland},
		},
		{func(e *domain.Env) string { return e.SwaySock }, match{domain.CompositorSway, nameSway}},
		{
			func(e *domain.Env) string { return e.WayfireSocket },
			match{domain.CompositorWayfire, nameWayfire},
		},
		{func(e *domain.Env) string { return e.I3Sock }, match{domain.CompositorI3, nameI3}},
	}
}

// waylandFingerprints identify a compositor from its advertised globals. Order
// matters: the specific interfaces come before the shared wlroots ones.
func waylandFingerprints() []fingerprint {
	return []fingerprint{
		{
			iface: "hyprland_toplevel_export_manager_v1",
			match: match{domain.CompositorHyprland, nameHyprland},
		},
		{
			iface: "hyprland_global_shortcuts_manager_v1",
			match: match{domain.CompositorHyprland, nameHyprland},
		},
		{
			iface: "hyprland_ctm_control_manager_v1",
			match: match{domain.CompositorHyprland, nameHyprland},
		},
		{
			iface: "hyprland_focus_grab_manager_v1",
			match: match{domain.CompositorHyprland, nameHyprland},
		},

		{iface: "org_kde_plasma_shell", match: match{domain.CompositorKWin, nameKWin}},
		{iface: "org_kde_kwin_blur_manager", match: match{domain.CompositorKWin, nameKWin}},
		{iface: "org_kde_kwin_appmenu_manager", match: match{domain.CompositorKWin, nameKWin}},
		{iface: "kde_output_device_v2", match: match{domain.CompositorKWin, nameKWin}},
		{iface: "zkde_screencast_unstable_v1", match: match{domain.CompositorKWin, nameKWin}},

		{iface: "gtk_shell1", match: match{domain.CompositorMutter, nameMutter}},

		{iface: "weston_desktop_shell", match: match{domain.CompositorWeston, nameWeston}},
		{iface: "weston_screenshooter", match: match{domain.CompositorWeston, nameWeston}},

		{iface: "zwf_shell_manager_v2", match: match{domain.CompositorWayfire, nameWayfire}},
		{iface: "wayfire_shell", match: match{domain.CompositorWayfire, nameWayfire}},

		{iface: "river_status_manager_v1", match: match{domain.CompositorRiver, nameRiver}},
		{iface: "river_control_v1", match: match{domain.CompositorRiver, nameRiver}},

		{iface: "zcosmic_toplevel_info_v1", match: match{domain.CompositorUnknown, nameCosmicComp}},
		{
			iface: "zcosmic_workspace_manager_v1",
			match: match{domain.CompositorUnknown, nameCosmicComp},
		},
	}
}

// wlrootsMarkers are shared by every wlroots-based compositor, so they narrow
// the field without naming a member of it.
func wlrootsMarkers() []string {
	return []string{
		"zwlr_layer_shell_v1",
		"zwlr_output_manager_v1",
		"zwlr_foreign_toplevel_manager_v1",
	}
}

// windowManagerNames maps a normalised EWMH _NET_WM_NAME onto a compositor.
func windowManagerNames(key string) (value match, ok bool) {
	switch key {
	case tokenMutter:
		return match{domain.CompositorMutter, nameMutter}, true
	case nameGnomeshell:
		return match{domain.CompositorMutter, "GNOME Shell"}, true
	case "muttermuffin", nameMuffin:
		return match{domain.CompositorUnknown, nameMuffin}, true
	case tokenKwin:
		return match{domain.CompositorKWin, nameKWin}, true
	case tokenXfwm4:
		return match{domain.CompositorXfwm, nameXfwm4}, true
	case tokenMarco:
		return match{domain.CompositorMarco, nameMarco}, true
	case tokenOpenbox:
		return match{domain.CompositorOpenbox, nameOpenbox}, true
	case nameI3:
		return match{domain.CompositorI3, nameI3}, true
	case nameAwesome:
		return match{domain.CompositorAwesome, nameAwesome}, true
	case nameSway:
		return match{domain.CompositorSway, nameSway}, true
	case nameLabwc:
		return match{domain.CompositorLabwc, nameLabwc}, true
	case nameWeston:
		return match{domain.CompositorWeston, nameWeston}, true
	case nameGala:
		return match{domain.CompositorUnknown, nameGala}, true
	case "budgiewm":
		return match{domain.CompositorUnknown, nameBudgieWM}, true
	case "compiz":
		return match{domain.CompositorUnknown, "Compiz"}, true
	case "metacity":
		return match{domain.CompositorUnknown, "Metacity"}, true
	case "icewm":
		return match{domain.CompositorUnknown, "IceWM"}, true
	case "fluxbox":
		return match{domain.CompositorUnknown, "Fluxbox"}, true
	case nameBspwm:
		return match{domain.CompositorUnknown, nameBspwm}, true
	case nameDwm:
		return match{domain.CompositorUnknown, nameDwm}, true
	case nameQtile:
		return match{domain.CompositorUnknown, nameQtile}, true
	case nameXmonad:
		return match{domain.CompositorUnknown, nameXmonad}, true
	case nameHerbstluftwm:
		return match{domain.CompositorUnknown, nameHerbstluftwm}, true
	case nameSpectrwm:
		return match{domain.CompositorUnknown, nameSpectrwm}, true
	default:
		return value, false
	}
}

// processNames maps a /proc comm value onto a compositor. This is the weakest
// signal: a matching process may belong to another seat entirely.
func processNames(key string) (value match, ok bool) {
	switch key {
	case "gnome-shell", tokenMutter, "mutter-x11-fr":
		return match{domain.CompositorMutter, nameMutter}, true
	case "kwin_wayland", "kwin_x11", tokenKwin:
		return match{domain.CompositorKWin, nameKWin}, true
	case nameSway:
		return match{domain.CompositorSway, nameSway}, true
	case nameHyprland, "hyprland":
		return match{domain.CompositorHyprland, nameHyprland}, true
	case nameWayfire:
		return match{domain.CompositorWayfire, nameWayfire}, true
	case nameRiver:
		return match{domain.CompositorRiver, nameRiver}, true
	case nameWeston:
		return match{domain.CompositorWeston, nameWeston}, true
	case nameLabwc:
		return match{domain.CompositorLabwc, nameLabwc}, true
	case tokenXfwm4:
		return match{domain.CompositorXfwm, nameXfwm4}, true
	case tokenMarco:
		return match{domain.CompositorMarco, nameMarco}, true
	case tokenOpenbox:
		return match{domain.CompositorOpenbox, nameOpenbox}, true
	case nameI3:
		return match{domain.CompositorI3, nameI3}, true
	case nameAwesome:
		return match{domain.CompositorAwesome, nameAwesome}, true
	case nameMuffin:
		return match{domain.CompositorUnknown, nameMuffin}, true
	case "cosmic-comp":
		return match{domain.CompositorUnknown, nameCosmicComp}, true
	case nameGala:
		return match{domain.CompositorUnknown, nameGala}, true
	case "budgie-wm":
		return match{domain.CompositorUnknown, nameBudgieWM}, true
	case nameNiri:
		return match{domain.CompositorUnknown, nameNiri}, true
	default:
		return value, false
	}
}

// desktopCompositors is the compositor a desktop environment ships with. This
// is an inference, not an observation, so it never yields better than medium
// confidence.
func desktopCompositors(key domain.DesktopEnvironment) (value match, ok bool) {
	switch key {
	case domain.DesktopUnknown, domain.DesktopLXQt:
		return match{kind: domain.CompositorUnknown, name: noValue}, true
	case domain.DesktopGNOME:
		return match{domain.CompositorMutter, nameMutter}, true
	case domain.DesktopKDE:
		return match{domain.CompositorKWin, nameKWin}, true
	case domain.DesktopXFCE:
		return match{domain.CompositorXfwm, nameXfwm4}, true
	case domain.DesktopMATE:
		return match{domain.CompositorMarco, nameMarco}, true
	case domain.DesktopLXDE:
		return match{domain.CompositorOpenbox, nameOpenbox}, true
	case domain.DesktopCinnamon:
		return match{domain.CompositorUnknown, nameMuffin}, true
	case domain.DesktopCOSMIC:
		return match{domain.CompositorUnknown, nameCosmicComp}, true
	case domain.DesktopPantheon:
		return match{domain.CompositorUnknown, nameGala}, true
	case domain.DesktopBudgie:
		return match{domain.CompositorUnknown, nameBudgieWM}, true
	default:
		return value, false
	}
}
