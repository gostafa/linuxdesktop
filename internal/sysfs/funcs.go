// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package sysfs

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
)

// String reads a small file and returns its contents with surrounding
// whitespace removed. It allocates once, for the returned string.
func String(path string) (string, error) {
	buf := take()

	read, err := readInto(path, buf)
	if err != nil {
		scratch.Put(buf)

		return noValue, err
	}

	out := string(bytes.TrimSpace((*buf)[:read]))
	scratch.Put(buf)

	return out, nil
}

// Trimmed is String with the error discarded. An unreadable file yields "",
// which is what every caller in this library wants: an absent sysfs attribute
// is missing information, not a failure.
func Trimmed(path string) string {
	out, _ := String(path)

	return out
}

// Bytes reads a file and returns an owned copy of its contents.
func Bytes(path string) ([]byte, error) {
	buf := take()

	read, err := readInto(path, buf)
	if err != nil {
		scratch.Put(buf)

		return nil, err
	}

	out := make([]byte, read)
	copy(out, (*buf)[:read])
	scratch.Put(buf)

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
	dir, err := os.Open(path)
	if err != nil {
		return nil, err
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

// Each walks a KEY=VALUE file, calling fn for every pair. Returning false from
// fn stops the walk. Blank lines and # comments are skipped, and matching
// surrounding quotes are stripped from values.
func Each(data []byte, fn func(key, value string) bool) {
	for line := range bytes.Lines(data) {
		key, value, ok := pair(line)
		if !ok {
			continue
		}

		if !fn(key, value) {
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

// take borrows a scratch buffer, falling back to a fresh one on the impossible
// day the pool hands back something else.
func take() *[]byte {
	if buf, ok := scratch.Get().(*[]byte); ok {
		return buf
	}

	fresh := make([]byte, scratchSize)

	return &fresh
}

// readInto fills *buf with the contents of path and reports how many bytes it
// holds.
func readInto(path string, buf *[]byte) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return zero, err
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
	read := zero

	for {
		got, err := step(file, buf, read)

		read += got

		if err != nil {
			return read, sansEOF(err)
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
	if read < len(*buf) {
		return file.Read((*buf)[read:])
	}

	if read >= maxFileSize {
		return zero, io.EOF
	}

	grown := make([]byte, len(*buf)*bufferGrowth)
	copy(grown, *buf)

	*buf = grown

	return file.Read((*buf)[read:])
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
