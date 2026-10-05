// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package rules

const (
	// Names of compositors that have no CompositorKind constant. They are
	// reported as CompositorUnknown with one of these in CompositorInfo.Name.
	nameMuffin     = "muffin"
	nameCosmicComp = "cosmic-comp"
	nameGala       = "gala"
	nameBudgieWM   = "budgie-wm"
	nameWlroots    = "wlroots"

	// desktopSeparator splits the colon-delimited $XDG_CURRENT_DESKTOP list.
	desktopSeparator = ":"

	// noValue is an unset environment variable or an unidentified name.
	noValue = ""

	// zero is the empty length a slice is built at.
	zero = 0

	// extraTokens is how many tokens beyond $XDG_CURRENT_DESKTOP can name a
	// desktop: $XDG_SESSION_DESKTOP and $DESKTOP_SESSION.
	extraTokens = 2

	nameGnomeshell   = "gnomeshell"
	nameHyprland     = "Hyprland"
	nameSway         = "sway"
	nameWayfire      = "wayfire"
	nameI3           = "i3"
	nameKWin         = "KWin"
	nameMutter       = "Mutter"
	nameWeston       = "weston"
	nameRiver        = "river"
	tokenMutter      = "mutter"
	tokenKwin        = "kwin"
	tokenXfwm4       = "xfwm4"
	nameXfwm4        = "Xfwm4"
	tokenMarco       = "marco"
	nameMarco        = "Marco"
	tokenOpenbox     = "openbox"
	nameOpenbox      = "Openbox"
	nameAwesome      = "awesome"
	nameLabwc        = "labwc"
	nameBspwm        = "bspwm"
	nameDwm          = "dwm"
	nameQtile        = "qtile"
	nameXmonad       = "xmonad"
	nameHerbstluftwm = "herbstluftwm"
	nameSpectrwm     = "spectrwm"
	nameNiri         = "niri"
)
