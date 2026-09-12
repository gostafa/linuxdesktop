package procscan

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gostafa/linuxdesktop/internal/sysfs"
)

// New returns a process probe. filter selects the command names worth
// confirming; a nil filter accepts everything.
func New(filter func(string) bool) *Probe {
	if filter == nil {
		filter = func(string) bool { return true }
	}
	return &Probe{filter: filter}
}

// Processes returns the command names of this user's processes that pass the
// filter.
func (p *Probe) Processes(ctx context.Context) ([]string, error) {
	entries, err := sysfs.DirNames(procDir)
	if err != nil {
		return nil, err
	}

	uid := strconv.Itoa(os.Getuid())
	var matches []string
	seen := make(map[string]struct{}, maxMatches)

	for _, entry := range entries {
		if !isPID(entry) {
			continue
		}
		if ctx.Err() != nil {
			return matches, ctx.Err()
		}
		name, err := sysfs.String(filepath.Join(procDir, entry, commFile))
		if err != nil || name == "" || !p.filter(name) {
			continue
		}
		if _, dup := seen[name]; dup {
			continue
		}
		if !ownedBy(entry, uid) {
			continue
		}
		seen[name] = struct{}{}
		matches = append(matches, name)
		if len(matches) == maxMatches {
			break
		}
	}
	return matches, nil
}

// isPID reports whether a /proc entry is a process directory. Checking the
// name is far cheaper than stat'ing every entry.
func isPID(name string) bool {
	if name == "" {
		return false
	}
	for i := 0; i < len(name); i++ {
		if name[i] < '0' || name[i] > '9' {
			return false
		}
	}
	return true
}

// ownedBy confirms a process belongs to uid by reading the Uid: line of its
// status file, whose first field is the real uid.
func ownedBy(pid, uid string) bool {
	data, err := sysfs.Bytes(filepath.Join(procDir, pid, statusFile))
	if err != nil {
		return false
	}
	for _, line := range strings.Split(string(data), "\n") {
		rest, found := strings.CutPrefix(line, uidPrefix)
		if !found {
			continue
		}
		fields := strings.Fields(rest)
		return len(fields) > 0 && fields[0] == uid
	}
	return false
}
