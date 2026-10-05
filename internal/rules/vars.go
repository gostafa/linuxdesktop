// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package rules

import "github.com/gostafa/linuxdesktop/internal/domain"

// desktopTokens maps a normalised $XDG_CURRENT_DESKTOP, $XDG_SESSION_DESKTOP
// or $DESKTOP_SESSION token onto a desktop environment. Tokens are matched in
// the order they appear in the variable, so "Budgie:GNOME" resolves to Budgie
// and "ubuntu:GNOME" to GNOME.
var desktopTokens = map[string]domain.DesktopEnvironment{
	"gnome":          domain.DesktopGNOME,
	"gnomeclassic":   domain.DesktopGNOME,
	"gnomeflashback": domain.DesktopGNOME,
	"gnomexorg":      domain.DesktopGNOME,
	"gnomewayland":   domain.DesktopGNOME,
	"gnomeshell":     domain.DesktopGNOME,
	"ubuntu":         domain.DesktopGNOME,
	"ubuntuwayland":  domain.DesktopGNOME,
	"ubuntuxorg":     domain.DesktopGNOME,
	"pop":            domain.DesktopGNOME,
	"zorin":          domain.DesktopGNOME,
	"kde":            domain.DesktopKDE,
	"plasma":         domain.DesktopKDE,
	"plasma5":        domain.DesktopKDE,
	"plasma6":        domain.DesktopKDE,
	"plasmawayland":  domain.DesktopKDE,
	"plasmax11":      domain.DesktopKDE,
	"xfce":           domain.DesktopXFCE,
	"xfce4":          domain.DesktopXFCE,
	"xubuntu":        domain.DesktopXFCE,
	"xcinnamon":      domain.DesktopCinnamon,
	"cinnamon":       domain.DesktopCinnamon,
	"mate":           domain.DesktopMATE,
	"lxqt":           domain.DesktopLXQt,
	"lxde":           domain.DesktopLXDE,
	"budgie":         domain.DesktopBudgie,
	"budgiedesktop":  domain.DesktopBudgie,
	"cosmic":         domain.DesktopCOSMIC,
	"pantheon":       domain.DesktopPantheon,
	"elementary":     domain.DesktopPantheon,
}

// envFingerprints are the strongest signal available: each of these variables
// is exported by exactly one compositor.
var envFingerprints = []struct {
	value func(*domain.Env) string
	match
}{
	{
		func(e *domain.Env) string { return e.HyprlandSignature },
		match{domain.CompositorHyprland, "Hyprland"},
	},
	{func(e *domain.Env) string { return e.SwaySock }, match{domain.CompositorSway, "sway"}},
	{
		func(e *domain.Env) string { return e.WayfireSocket },
		match{domain.CompositorWayfire, "wayfire"},
	},
	{func(e *domain.Env) string { return e.I3Sock }, match{domain.CompositorI3, "i3"}},
}

// waylandFingerprints identify a compositor from its advertised globals. Order
// matters: the specific interfaces come before the shared wlroots ones.
var waylandFingerprints = []fingerprint{
	{"hyprland_toplevel_export_manager_v1", match{domain.CompositorHyprland, "Hyprland"}},
	{"hyprland_global_shortcuts_manager_v1", match{domain.CompositorHyprland, "Hyprland"}},
	{"hyprland_ctm_control_manager_v1", match{domain.CompositorHyprland, "Hyprland"}},
	{"hyprland_focus_grab_manager_v1", match{domain.CompositorHyprland, "Hyprland"}},

	{"org_kde_plasma_shell", match{domain.CompositorKWin, "KWin"}},
	{"org_kde_kwin_blur_manager", match{domain.CompositorKWin, "KWin"}},
	{"org_kde_kwin_appmenu_manager", match{domain.CompositorKWin, "KWin"}},
	{"kde_output_device_v2", match{domain.CompositorKWin, "KWin"}},
	{"zkde_screencast_unstable_v1", match{domain.CompositorKWin, "KWin"}},

	{"gtk_shell1", match{domain.CompositorMutter, "Mutter"}},

	{"weston_desktop_shell", match{domain.CompositorWeston, "weston"}},
	{"weston_screenshooter", match{domain.CompositorWeston, "weston"}},

	{"zwf_shell_manager_v2", match{domain.CompositorWayfire, "wayfire"}},
	{"wayfire_shell", match{domain.CompositorWayfire, "wayfire"}},

	{"river_status_manager_v1", match{domain.CompositorRiver, "river"}},
	{"river_control_v1", match{domain.CompositorRiver, "river"}},

	{"zcosmic_toplevel_info_v1", match{domain.CompositorUnknown, nameCosmicComp}},
	{"zcosmic_workspace_manager_v1", match{domain.CompositorUnknown, nameCosmicComp}},
}

// wlrootsMarkers are shared by every wlroots-based compositor, so they narrow
// the field without naming a member of it.
var wlrootsMarkers = []string{
	"zwlr_layer_shell_v1",
	"zwlr_output_manager_v1",
	"zwlr_foreign_toplevel_manager_v1",
}

// windowManagerNames maps a normalised EWMH _NET_WM_NAME onto a compositor.
var windowManagerNames = map[string]match{
	"mutter":       {domain.CompositorMutter, "Mutter"},
	"gnomeshell":   {domain.CompositorMutter, "GNOME Shell"},
	"muttermuffin": {domain.CompositorUnknown, nameMuffin},
	"muffin":       {domain.CompositorUnknown, nameMuffin},
	"kwin":         {domain.CompositorKWin, "KWin"},
	"xfwm4":        {domain.CompositorXfwm, "Xfwm4"},
	"marco":        {domain.CompositorMarco, "Marco"},
	"openbox":      {domain.CompositorOpenbox, "Openbox"},
	"i3":           {domain.CompositorI3, "i3"},
	"awesome":      {domain.CompositorAwesome, "awesome"},
	"sway":         {domain.CompositorSway, "sway"},
	"labwc":        {domain.CompositorLabwc, "labwc"},
	"weston":       {domain.CompositorWeston, "weston"},
	"gala":         {domain.CompositorUnknown, nameGala},
	"budgiewm":     {domain.CompositorUnknown, nameBudgieWM},
	"compiz":       {domain.CompositorUnknown, "Compiz"},
	"metacity":     {domain.CompositorUnknown, "Metacity"},
	"icewm":        {domain.CompositorUnknown, "IceWM"},
	"fluxbox":      {domain.CompositorUnknown, "Fluxbox"},
	"bspwm":        {domain.CompositorUnknown, "bspwm"},
	"dwm":          {domain.CompositorUnknown, "dwm"},
	"qtile":        {domain.CompositorUnknown, "qtile"},
	"xmonad":       {domain.CompositorUnknown, "xmonad"},
	"herbstluftwm": {domain.CompositorUnknown, "herbstluftwm"},
	"spectrwm":     {domain.CompositorUnknown, "spectrwm"},
}

// processNames maps a /proc comm value onto a compositor. This is the weakest
// signal: a matching process may belong to another seat entirely.
var processNames = map[string]match{
	"gnome-shell":   {domain.CompositorMutter, "Mutter"},
	"mutter":        {domain.CompositorMutter, "Mutter"},
	"mutter-x11-fr": {domain.CompositorMutter, "Mutter"},
	"kwin_wayland":  {domain.CompositorKWin, "KWin"},
	"kwin_x11":      {domain.CompositorKWin, "KWin"},
	"kwin":          {domain.CompositorKWin, "KWin"},
	"sway":          {domain.CompositorSway, "sway"},
	"Hyprland":      {domain.CompositorHyprland, "Hyprland"},
	"hyprland":      {domain.CompositorHyprland, "Hyprland"},
	"wayfire":       {domain.CompositorWayfire, "wayfire"},
	"river":         {domain.CompositorRiver, "river"},
	"weston":        {domain.CompositorWeston, "weston"},
	"labwc":         {domain.CompositorLabwc, "labwc"},
	"xfwm4":         {domain.CompositorXfwm, "Xfwm4"},
	"marco":         {domain.CompositorMarco, "Marco"},
	"openbox":       {domain.CompositorOpenbox, "Openbox"},
	"i3":            {domain.CompositorI3, "i3"},
	"awesome":       {domain.CompositorAwesome, "awesome"},
	"muffin":        {domain.CompositorUnknown, nameMuffin},
	"cosmic-comp":   {domain.CompositorUnknown, nameCosmicComp},
	"gala":          {domain.CompositorUnknown, nameGala},
	"budgie-wm":     {domain.CompositorUnknown, nameBudgieWM},
	"niri":          {domain.CompositorUnknown, "niri"},
}

// desktopCompositors is the compositor a desktop environment ships with. This
// is an inference, not an observation, so it never yields better than medium
// confidence.
var desktopCompositors = map[domain.DesktopEnvironment]match{
	domain.DesktopGNOME:    {domain.CompositorMutter, "Mutter"},
	domain.DesktopKDE:      {domain.CompositorKWin, "KWin"},
	domain.DesktopXFCE:     {domain.CompositorXfwm, "Xfwm4"},
	domain.DesktopMATE:     {domain.CompositorMarco, "Marco"},
	domain.DesktopLXDE:     {domain.CompositorOpenbox, "Openbox"},
	domain.DesktopCinnamon: {domain.CompositorUnknown, nameMuffin},
	domain.DesktopCOSMIC:   {domain.CompositorUnknown, nameCosmicComp},
	domain.DesktopPantheon: {domain.CompositorUnknown, nameGala},
	domain.DesktopBudgie:   {domain.CompositorUnknown, nameBudgieWM},
}
