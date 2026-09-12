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
	bp := scratch.Get().(*[]byte)
	n, err := readInto(path, bp)
	if err != nil {
		scratch.Put(bp)
		return "", err
	}
	s := string(bytes.TrimSpace((*bp)[:n]))
	scratch.Put(bp)
	return s, nil
}

// Trimmed is String with the error discarded. An unreadable file yields "",
// which is what every caller in this library wants: an absent sysfs attribute
// is missing information, not a failure.
func Trimmed(path string) string {
	s, _ := String(path)
	return s
}

// Bytes reads a file and returns an owned copy of its contents.
func Bytes(path string) ([]byte, error) {
	bp := scratch.Get().(*[]byte)
	n, err := readInto(path, bp)
	if err != nil {
		scratch.Put(bp)
		return nil, err
	}
	out := make([]byte, n)
	copy(out, (*bp)[:n])
	scratch.Put(bp)
	return out, nil
}

// IsSocket reports whether path exists and is a unix socket. This is the
// cheapest possible "is there a display server here" check.
func IsSocket(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.Mode()&os.ModeSocket != 0
}

// DirNames lists the entry names of a directory. It uses Readdirnames so the
// kernel is never asked to stat entries the caller may not care about.
func DirNames(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	names, err := f.Readdirnames(-1)
	f.Close()
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
		return ""
	}
	return filepath.Base(target)
}

// Each walks a KEY=VALUE file, calling fn for every pair. Returning false from
// fn stops the walk. Blank lines and # comments are skipped, and matching
// surrounding quotes are stripped from values.
func Each(data []byte, fn func(key, value string) bool) {
	for len(data) > 0 {
		var line []byte
		line, data = cutLine(data)
		line = bytes.TrimSpace(line)
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		eq := bytes.IndexByte(line, '=')
		if eq < 0 {
			continue
		}
		if !fn(string(bytes.TrimSpace(line[:eq])), unquote(bytes.TrimSpace(line[eq+1:]))) {
			return
		}
	}
}

// Field returns the value of key in a KEY=VALUE file, or "" when it is absent.
// Nothing is allocated until the key matches.
func Field(data []byte, key string) string {
	for len(data) > 0 {
		var line []byte
		line, data = cutLine(data)
		line = bytes.TrimSpace(line)
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		eq := bytes.IndexByte(line, '=')
		if eq < 0 {
			continue
		}
		if string(bytes.TrimSpace(line[:eq])) != key {
			continue
		}
		return unquote(bytes.TrimSpace(line[eq+1:]))
	}
	return ""
}

// readInto fills *bp with the contents of path, growing the buffer only when
// the file turns out to be larger than the pooled scratch size.
func readInto(path string, bp *[]byte) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	buf := *bp
	n := 0
	for {
		if n == len(buf) {
			if n >= maxFileSize {
				break
			}
			grown := make([]byte, len(buf)*2)
			copy(grown, buf)
			buf = grown
		}
		r, rerr := f.Read(buf[n:])
		n += r
		if rerr != nil {
			*bp = buf
			if errors.Is(rerr, io.EOF) {
				return n, nil
			}
			return n, rerr
		}
		if r == 0 {
			break
		}
	}
	*bp = buf
	return n, nil
}

// cutLine splits off the first line, returning it and the remainder.
func cutLine(data []byte) (line, rest []byte) {
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		return data[:i], data[i+1:]
	}
	return data, nil
}

// unquote strips one layer of matching single or double quotes.
func unquote(b []byte) string {
	if len(b) >= 2 {
		if c := b[0]; (c == '"' || c == '\'') && b[len(b)-1] == c {
			b = b[1 : len(b)-1]
		}
	}
	return string(b)
}
