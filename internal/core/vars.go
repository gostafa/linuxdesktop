package core

import (
	"time"

	"github.com/gostafa/linuxdesktop/internal/domain"
)

// Defaults for a run. The overall budget is generous enough that a healthy
// desktop never approaches it, and the per-probe budget is short enough that a
// single unresponsive server costs a fraction of a second rather than the lot.
const (
	DefaultTimeout      = 2 * time.Second
	DefaultProbeTimeout = 500 * time.Millisecond
)

// DefaultConfig is the configuration used when the caller supplies no options.
var DefaultConfig = Config{
	Timeout:      DefaultTimeout,
	ProbeTimeout: DefaultProbeTimeout,
	Sections:     domain.SectionAll,
}
