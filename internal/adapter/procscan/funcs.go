// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

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
func (probe *Probe) Processes(ctx context.Context) ([]string, error) {
	entries, err := sysfs.DirNames(procDir)
	if err != nil {
		return nil, err
	}

	return newScanner(probe.filter).walk(ctx, entries)
}

func newScanner(filter func(string) bool) *scanner {
	return &scanner{
		filter: filter,
		seen:   make(map[string]bool, maxMatches),
		uid:    strconv.Itoa(os.Getuid()),
	}
}

// accepts reports whether a command name is a new one the filter wants.
func (scan *scanner) accepts(name string) bool {
	return name != noValue && !scan.seen[name] && scan.filter(name)
}

// consider records one /proc entry if it names a process worth keeping, and
// reports whether the result is now full.
func (scan *scanner) consider(entry string) bool {
	name, ok := scan.wanted(entry)
	if !ok {
		return false
	}

	scan.seen[name] = true
	scan.matches = append(scan.matches, name)

	return len(scan.matches) == maxMatches
}

// walk visits every process directory, stopping when the result is full or the
// caller gives up.
func (scan *scanner) walk(ctx context.Context, entries []string) ([]string, error) {
	for i := range entries {
		if ctx.Err() != nil {
			return scan.matches, ctx.Err()
		}

		if scan.consider(entries[i]) {
			break
		}
	}

	return scan.matches, nil
}

// wanted reads a process's command name and reports whether it is a new one
// belonging to this user that the filter accepts. Ownership is confirmed last,
// since it costs a second file read.
func (scan *scanner) wanted(entry string) (string, bool) {
	if !isPID(entry) {
		return noValue, false
	}

	name := commName(entry)
	if !scan.accepts(name) {
		return noValue, false
	}

	return name, ownedBy(entry, scan.uid)
}

// commName reads a process's command name, or "" when it has gone away.
func commName(entry string) string {
	name, err := sysfs.String(filepath.Join(procDir, entry, commFile))
	if err != nil {
		return noValue
	}

	return name
}

// isPID reports whether a /proc entry is a process directory. Checking the
// name is far cheaper than stat'ing every entry.
func isPID(name string) bool {
	if name == noValue {
		return false
	}

	for i := range len(name) {
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

	return realUID(data) == uid
}

// realUID is the first field of the Uid: line of a status file.
func realUID(data []byte) string {
	for line := range strings.SplitSeq(string(data), "\n") {
		rest, found := strings.CutPrefix(line, uidPrefix)
		if !found {
			continue
		}

		return firstField(rest)
	}

	return noValue
}

// firstField is the leading whitespace-separated field of a line.
func firstField(line string) string {
	fields := strings.Fields(line)
	if len(fields) == zero {
		return noValue
	}

	return fields[zero]
}
