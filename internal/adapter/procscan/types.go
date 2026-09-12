package procscan

// Probe implements port.ProcessProbe.
//
// filter decides which command names are worth confirming. It is injected
// rather than hard-coded so the compositor table lives in exactly one place.
type Probe struct {
	filter func(string) bool
}
