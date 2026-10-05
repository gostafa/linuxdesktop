// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package procscan

import (
	"context"
	"fmt"
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

	return &Probe{filter: filter, root: procDir}
}

// Processes returns the command names of this user's processes that pass the
// filter.
func (probe *Probe) Processes(ctx context.Context) ([]string, error) {
	entries, err := sysfs.DirNames(probe.root)
	if err != nil {
		return nil, fmt.Errorf(errListProcesses, err)
	}

	result, callErr := walk(ctx, newScanner(probe.filter, probe.root), entries)
	if callErr != nil {
		return result, fmt.Errorf(errListProcesses, callErr)
	}

	return result, nil
}

func newScanner(filter func(string) bool, root string) *scanner {
	return &scanner{
		filter:  filter,
		root:    root,
		seen:    make(map[string]bool, maxMatches),
		uid:     strconv.Itoa(os.Getuid()),
		matches: nil,
	}
}

// accepts reports whether a command name is a new one the filter wants.
func accepts(scan *scanner, name string) bool {
	return name != noValue && !scan.seen[name] && scan.filter(name)
}

// consider records one /proc entry if it names a process worth keeping, and
// reports whether the result is now full.
func consider(scan *scanner, entry string) bool {
	name, ok := wanted(scan, entry)
	if !ok {
		return false
	}

	scan.seen[name] = true
	scan.matches = append(scan.matches, name)

	return len(scan.matches) == maxMatches
}

// walk visits every process directory, stopping when the result is full or the
// caller gives up.
func walk(ctx context.Context, scan *scanner, entries []string) ([]string, error) {
	for i := range entries {
		err := ctx.Err()
		if err != nil {
			return scan.matches, fmt.Errorf("procscan: scan processes: %w", err)
		}

		if consider(scan, entries[i]) {
			break
		}
	}

	return scan.matches, nil
}

// wanted reads a process's command name and reports whether it is a new one
// belonging to this user that the filter accepts. Ownership is confirmed last,
// since it costs a second file read.
func wanted(scan *scanner, entry string) (string, bool) {
	if !isPID(entry) {
		return noValue, false
	}

	name := scan.commName(entry)
	if !accepts(scan, name) {
		return noValue, false
	}

	return name, scan.ownedBy(entry, scan.uid)
}

// commName reads a process's command name, or "" when it has gone away.
func (scan *scanner) commName(entry string) string {
	name, err := sysfs.String(filepath.Join(scan.root, entry, commFile))
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
func (scan *scanner) ownedBy(pid, uid string) bool {
	data, err := sysfs.Bytes(filepath.Join(scan.root, pid, statusFile))
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
