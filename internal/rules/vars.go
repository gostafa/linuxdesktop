// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package rules

import (
	"slices"

	"github.com/gostafa/linuxdesktop/internal/domain"
)

// desktopTokens maps a normalised $XDG_CURRENT_DESKTOP, $XDG_SESSION_DESKTOP
// or $DESKTOP_SESSION token onto a desktop environment. Tokens are matched in
// the order they appear in the variable, so "Budgie:GNOME" resolves to Budgie
// and "ubuntu:GNOME" to GNOME.
func desktopTokens(key string) (value domain.DesktopEnvironment, ok bool) {
	aliases := desktopAliases()

	for desktop := range aliases {
		if slices.Contains(aliases[desktop], key) {
			return desktop, true
		}
	}

	return value, false
}

func desktopAliases() map[domain.DesktopEnvironment][]string {
	return map[domain.DesktopEnvironment][]string{
		domain.Unknown:         nil,
		domain.DesktopGNOME:    gnomeTokens(),
		domain.DesktopKDE:      plasmaTokens(),
		domain.DesktopXFCE:     {"xfce", "xfce4", "xubuntu"},
		domain.DesktopCinnamon: {"xcinnamon", "cinnamon"},
		domain.DesktopMATE:     {"mate"},
		domain.DesktopLXQt:     {"lxqt"},
		domain.DesktopLXDE:     {"lxde"},
		domain.DesktopBudgie:   {"budgie", "budgiedesktop"},
		domain.DesktopCOSMIC:   {"cosmic"},
		domain.DesktopPantheon: {"pantheon", "elementary"},
	}
}

func gnomeTokens() []string {
	return []string{
		"gnome", "gnomeclassic", "gnomeflashback", "gnomexorg", "gnomewayland",
		nameGnomeshell, "ubuntu", "ubuntuwayland", "ubuntuxorg", "pop", "zorin",
	}
}

func plasmaTokens() []string {
	return []string{"kde", "plasma", "plasma5", "plasma6", "plasmawayland", "plasmax11"}
}

// envFingerprints are the strongest signal available: each of these variables
// is exported by exactly one compositor.
func envFingerprints() []envFingerprint {
	return []envFingerprint{
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
	found := hyprlandFingerprints()

	found = append(found, kdeFingerprints()...)

	return append(found, otherFingerprints()...)
}

func hyprlandFingerprints() []fingerprint {
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
	}
}

func kdeFingerprints() []fingerprint {
	return []fingerprint{
		{iface: "org_kde_plasma_shell", match: match{domain.CompositorKWin, nameKWin}},
		{iface: "org_kde_kwin_blur_manager", match: match{domain.CompositorKWin, nameKWin}},
		{iface: "org_kde_kwin_appmenu_manager", match: match{domain.CompositorKWin, nameKWin}},
		{iface: "kde_output_device_v2", match: match{domain.CompositorKWin, nameKWin}},
		{iface: "zkde_screencast_unstable_v1", match: match{domain.CompositorKWin, nameKWin}},
	}
}

func otherFingerprints() []fingerprint {
	return []fingerprint{
		{iface: "gtk_shell1", match: match{domain.CompositorMutter, nameMutter}},
		{iface: "weston_desktop_shell", match: match{domain.CompositorWeston, nameWeston}},
		{iface: "weston_screenshooter", match: match{domain.CompositorWeston, nameWeston}},
		{iface: "zwf_shell_manager_v2", match: match{domain.CompositorWayfire, nameWayfire}},
		{iface: "wayfire_shell", match: match{domain.CompositorWayfire, nameWayfire}},
		{iface: "river_status_manager_v1", match: match{domain.CompositorRiver, nameRiver}},
		{iface: "river_control_v1", match: match{domain.CompositorRiver, nameRiver}},
		{iface: "zcosmic_toplevel_info_v1", match: match{domain.Unknown, nameCosmicComp}},
		{
			iface: "zcosmic_workspace_manager_v1",
			match: match{domain.Unknown, nameCosmicComp},
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
	tables := []map[string]match{
		windowManagerNamesTableA(),
		windowManagerNamesTableB(),
		windowManagerNamesTableC(),
	}

	for i := range tables {
		if found, exists := tables[i][key]; exists {
			return found, true
		}
	}

	return value, false
}

func windowManagerNamesTableA() map[string]match {
	return map[string]match{
		tokenMutter:    {domain.CompositorMutter, nameMutter},
		nameGnomeshell: {domain.CompositorMutter, "GNOME Shell"},
		"muttermuffin": {domain.Unknown, nameMuffin},
		nameMuffin:     {domain.Unknown, nameMuffin},
		tokenKwin:      {domain.CompositorKWin, nameKWin},
		tokenXfwm4:     {domain.CompositorXfwm, nameXfwm4},
		tokenMarco:     {domain.CompositorMarco, nameMarco},
		tokenOpenbox:   {domain.CompositorOpenbox, nameOpenbox},
		nameI3:         {domain.CompositorI3, nameI3},
		nameAwesome:    {domain.CompositorAwesome, nameAwesome},
	}
}

func windowManagerNamesTableB() map[string]match {
	return map[string]match{
		nameSway:   {domain.CompositorSway, nameSway},
		nameLabwc:  {domain.CompositorLabwc, nameLabwc},
		nameWeston: {domain.CompositorWeston, nameWeston},
		nameGala:   {domain.Unknown, nameGala},
		"budgiewm": {domain.Unknown, nameBudgieWM},
		"compiz":   {domain.Unknown, "Compiz"},
		"metacity": {domain.Unknown, "Metacity"},
		"icewm":    {domain.Unknown, "IceWM"},
		"fluxbox":  {domain.Unknown, "Fluxbox"},
		nameBspwm:  {domain.Unknown, nameBspwm},
	}
}

func windowManagerNamesTableC() map[string]match {
	return map[string]match{
		nameDwm:          {domain.Unknown, nameDwm},
		nameQtile:        {domain.Unknown, nameQtile},
		nameXmonad:       {domain.Unknown, nameXmonad},
		nameHerbstluftwm: {domain.Unknown, nameHerbstluftwm},
		nameSpectrwm:     {domain.Unknown, nameSpectrwm},
	}
}

// processNames maps a /proc comm value onto a compositor. This is the weakest
// signal: a matching process may belong to another seat entirely.
func processNames(key string) (value match, ok bool) {
	tables := []map[string]match{processNamesTableA(), processNamesTableB(), processNamesTableC()}

	for i := range tables {
		if found, exists := tables[i][key]; exists {
			return found, true
		}
	}

	return value, false
}

func processNamesTableA() map[string]match {
	return map[string]match{
		"gnome-shell":   {domain.CompositorMutter, nameMutter},
		tokenMutter:     {domain.CompositorMutter, nameMutter},
		"mutter-x11-fr": {domain.CompositorMutter, nameMutter},
		"kwin_wayland":  {domain.CompositorKWin, nameKWin},
		"kwin_x11":      {domain.CompositorKWin, nameKWin},
		tokenKwin:       {domain.CompositorKWin, nameKWin},
		nameSway:        {domain.CompositorSway, nameSway},
		nameHyprland:    {domain.CompositorHyprland, nameHyprland},
		"hyprland":      {domain.CompositorHyprland, nameHyprland},
		nameWayfire:     {domain.CompositorWayfire, nameWayfire},
	}
}

func processNamesTableB() map[string]match {
	return map[string]match{
		nameRiver:     {domain.CompositorRiver, nameRiver},
		nameWeston:    {domain.CompositorWeston, nameWeston},
		nameLabwc:     {domain.CompositorLabwc, nameLabwc},
		tokenXfwm4:    {domain.CompositorXfwm, nameXfwm4},
		tokenMarco:    {domain.CompositorMarco, nameMarco},
		tokenOpenbox:  {domain.CompositorOpenbox, nameOpenbox},
		nameI3:        {domain.CompositorI3, nameI3},
		nameAwesome:   {domain.CompositorAwesome, nameAwesome},
		nameMuffin:    {domain.Unknown, nameMuffin},
		"cosmic-comp": {domain.Unknown, nameCosmicComp},
	}
}

func processNamesTableC() map[string]match {
	return map[string]match{
		nameGala:    {domain.Unknown, nameGala},
		"budgie-wm": {domain.Unknown, nameBudgieWM},
		nameNiri:    {domain.Unknown, nameNiri},
	}
}

// desktopCompositors is the compositor a desktop environment ships with. This
// is an inference, not an observation, so it never yields better than medium
// confidence.
func desktopCompositors(key domain.DesktopEnvironment) (match, bool) {
	table := desktopCompositorsTable()
	found, ok := table[key]

	return found, ok
}

func desktopCompositorsTable() map[domain.DesktopEnvironment]match {
	return map[domain.DesktopEnvironment]match{
		domain.Unknown:         {kind: domain.Unknown, name: noValue},
		domain.DesktopLXQt:     {kind: domain.Unknown, name: noValue},
		domain.DesktopGNOME:    {domain.CompositorMutter, nameMutter},
		domain.DesktopKDE:      {domain.CompositorKWin, nameKWin},
		domain.DesktopXFCE:     {domain.CompositorXfwm, nameXfwm4},
		domain.DesktopMATE:     {domain.CompositorMarco, nameMarco},
		domain.DesktopLXDE:     {domain.CompositorOpenbox, nameOpenbox},
		domain.DesktopCinnamon: {domain.Unknown, nameMuffin},
		domain.DesktopCOSMIC:   {domain.Unknown, nameCosmicComp},
		domain.DesktopPantheon: {domain.Unknown, nameGala},
		domain.DesktopBudgie:   {domain.Unknown, nameBudgieWM},
	}
}
