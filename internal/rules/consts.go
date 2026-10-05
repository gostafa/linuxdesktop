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
)
