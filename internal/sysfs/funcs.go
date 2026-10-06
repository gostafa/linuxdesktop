// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package sysfs

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Path resolves a system path under root. An empty root reads the current system.
func Path(root, path string) string { return filepath.Join(root, path) }

// String reads a small file and returns its contents with surrounding
// whitespace removed. Its scratch buffer belongs to this read.
func String(path string) (string, error) {
	buf := take()

	read, err := readInto(path, buf)
	if err != nil {
		return noValue, fmt.Errorf("sysfs: read string: %w", err)
	}

	out := string(bytes.TrimSpace((*buf)[:read]))

	return out, nil
}

// Trimmed is String with the error discarded. An unreadable file yields "",
// which is what every caller in this library wants: an absent sysfs attribute
// is missing information, not a failure.
func Trimmed(path string) string {
	out, err := String(path)
	if err != nil {
		return noValue
	}

	return out
}

// Bytes reads a file and returns an owned copy of its contents.
func Bytes(path string) ([]byte, error) {
	buf := take()

	read, err := readInto(path, buf)
	if err != nil {
		return nil, fmt.Errorf("sysfs: read bytes: %w", err)
	}

	out := bytes.Clone((*buf)[:read])

	return out, nil
}

// IsSocket reports whether path exists and is a unix socket. This is the
// cheapest possible "is there a display server here" check.
func IsSocket(path string) bool {
	stat, err := os.Stat(path)

	return err == nil && stat.Mode().Type() == os.ModeSocket
}

// DirNames lists the entry names of a directory. It uses Readdirnames so the
// kernel is never asked to stat entries the caller may not care about.
func DirNames(path string) ([]string, error) {
	//nolint:gosec // Internal adapters supply system paths or explicitly configured fixture roots.
	dir, err := os.Open(
		path,
	)
	if err != nil {
		return nil, fmt.Errorf("sysfs: open directory: %w", err)
	}

	names, err := dir.Readdirnames(-1)
	closed := dir.Close()

	if err == nil {
		err = closed
	}

	if err != nil && !errors.Is(err, io.EOF) {
		return names, err
	}

	return names, nil
}

// LinkBase reads a symlink and returns the last element of its target. Sysfs
// uses this shape everywhere: /sys/class/drm/card0/device points at
// ../../../0000:00:02.0, whose base is the PCI address.
func LinkBase(path string) string {
	target, err := os.Readlink(path)
	if err != nil {
		return noValue
	}

	return filepath.Base(target)
}

// Each walks a KEY=VALUE file, calling visit for every pair. Returning false from
// visit stops the walk. Blank lines and # comments are skipped, and matching
// surrounding quotes are stripped from values.
func Each(data []byte, visit func(key, value string) bool) {
	for line := range bytes.Lines(data) {
		key, value, ok := pair(line)
		if !ok {
			continue
		}

		if !visit(key, value) {
			return
		}
	}
}

// Field returns the value of key in a KEY=VALUE file, or "" when it is absent.
// Nothing is allocated until the key matches.
func Field(data []byte, key string) string {
	for line := range bytes.Lines(data) {
		name, value, ok := pair(line)
		if ok && name == key {
			return value
		}
	}

	return noValue
}

// take creates the buffer owned by one read.
func take() *[]byte {
	var buffer [scratchSize]byte

	data := buffer[:]

	return &data
}

// readInto fills *buf with the contents of path and reports how many bytes it
// holds.
func readInto(path string, buf *[]byte) (int, error) {
	//nolint:gosec // Internal adapters supply system paths or explicitly configured fixture roots.
	file, err := os.Open(
		path,
	)
	if err != nil {
		return zero, fmt.Errorf("sysfs: open file: %w", err)
	}

	read, err := drain(file, buf)
	closed := file.Close()

	if err == nil {
		err = closed
	}

	return read, err
}

// drain reads file to its end, letting step decide when the buffer must grow.
func drain(file io.Reader, buf *[]byte) (int, error) {
	read, err := readChunks(file, buf)

	callErr := readError(err)
	if callErr != nil {
		return read, fmt.Errorf("sysfs: drain file: %w", callErr)
	}

	return read, nil
}

func readChunks(file io.Reader, buf *[]byte) (int, error) {
	read := zero

	for {
		got, err := step(file, buf, read)

		read += got

		if err != nil {
			return read, fmt.Errorf("sysfs: read chunks: %w", err)
		}

		if got == zero {
			return read, nil
		}
	}
}

// step reads one chunk into the free space of *buf, doubling the buffer first
// when it is full. Reaching maxFileSize ends the read rather than growing past
// it, so a hostile path cannot allocate without bound.
func step(file io.Reader, buf *[]byte, read int) (int, error) {
	if read >= len(*buf) && read >= maxFileSize {
		return zero, io.EOF
	}

	growBuffer(buf, read)

	result, callErr := file.Read((*buf)[read:])
	if callErr != nil {
		return result, fmt.Errorf(errReadChunk, callErr)
	}

	return result, nil
}

// sansEOF turns the io.EOF that ends every complete read into success.
func sansEOF(err error) error {
	if errors.Is(err, io.EOF) {
		return nil
	}

	return err
}

// pair lifts one KEY=VALUE pair off a line, reporting false for a blank line, a
// comment, or a line carrying no separator.
func pair(line []byte) (key, value string, ok bool) {
	before, after, cut := bytes.Cut(bytes.TrimSpace(line), []byte{separatorByte})
	if !cut || comment(before) {
		return noValue, noValue, false
	}

	return string(bytes.TrimSpace(before)), unquote(bytes.TrimSpace(after)), true
}

// comment reports whether what parsed as a key is really a comment.
func comment(key []byte) bool {
	return len(key) > zero && key[zero] == commentByte
}

// unquote strips one layer of matching single or double quotes.
func unquote(raw []byte) string {
	if len(raw) <= quoteWidth {
		return string(raw)
	}

	lead := raw[zero]
	if !quote(lead) || raw[len(raw)-quoteWidth] != lead {
		return string(raw)
	}

	return string(raw[quoteWidth : len(raw)-quoteWidth])
}

// quote reports whether char is one of the two quoting bytes.
func quote(char byte) bool {
	return char == '"' || char == '\''
}

func growBuffer(buf *[]byte, read int) {
	if read < len(*buf) {
		return
	}

	grown := append(*buf, make([]byte, len(*buf)*(bufferGrowth-1))...)

	*buf = grown
}

func readError(err error) error {
	err = sansEOF(err)
	if err != nil {
		return fmt.Errorf("sysfs: read file contents: %w", err)
	}

	return nil
}
