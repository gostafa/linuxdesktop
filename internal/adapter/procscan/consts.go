package procscan

const (
	procDir    = "/proc"
	commFile   = "comm"
	statusFile = "status"

	// uidPrefix introduces the line in /proc/<pid>/status whose first field is
	// the process's real uid.
	uidPrefix = "Uid:"

	// maxMatches bounds the result. More than a couple of compositor-looking
	// processes means something unusual is going on, and the extra names would
	// not improve a low-confidence guess.
	maxMatches = 8
)
