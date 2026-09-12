package rules

// Names of compositors that have no CompositorKind constant. They are reported
// as CompositorUnknown with one of these in CompositorInfo.Name.
const (
	nameMuffin     = "muffin"
	nameCosmicComp = "cosmic-comp"
	nameGala       = "gala"
	nameBudgieWM   = "budgie-wm"
	nameWlroots    = "wlroots"
)

// desktopSeparator splits the colon-delimited $XDG_CURRENT_DESKTOP list.
const desktopSeparator = ":"
