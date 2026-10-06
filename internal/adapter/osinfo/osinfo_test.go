// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package osinfo

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/gostafa/linuxdesktop/internal/sysfs"
)

func write(t *testing.T, root, path, data string) {
	t.Helper()
	path = sysfs.Path(root, path)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestOSRelease(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	write(
		t,
		root,
		pathOSRelease,
		"NAME=Fixture\nPRETTY_NAME=\nID=fixture\nID_LIKE=base\nVERSION=1\nVERSION_ID=1\nUNUSED=yes\n",
	)
	info, err := operatingSystem(root).OS(t.Context())
	if err != nil || info.Name != "Fixture" || info.PrettyName != "Fixture" ||
		info.ID != "fixture" ||
		info.IDLike != "base" ||
		info.Version != "1" ||
		info.VersionID != "1" {
		t.Fatalf("%+v, %v", info, err)
	}
	if err := os.Remove(sysfs.Path(root, pathOSRelease)); err != nil {
		t.Fatal(err)
	}
	write(t, root, pathOSReleaseAlt, "NAME=Fallback\nPRETTY_NAME='Fallback OS'\n")
	info, err = operatingSystem(root).OS(t.Context())
	if err != nil || info.PrettyName != "Fallback OS" {
		t.Fatal(info, err)
	}
	if err := os.Remove(sysfs.Path(root, pathOSReleaseAlt)); err != nil {
		t.Fatal(err)
	}
	if info, err = operatingSystem(root).OS(t.Context()); err != nil || info.Name != "" {
		t.Fatal(info, err)
	}
	if _, err = New().OS(t.Context()); err != nil {
		t.Fatal(err)
	}
}

func TestKernelIdentity(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	if info := kernelInfo(
		root,
		goosLinux,
		"amd64",
	); info.Kernel != "Linux" ||
		info.Architecture != "x86_64" {
		t.Fatal(info)
	}
	write(t, root, pathKernelType, "FixtureKernel")
	write(t, root, pathKernelRelase, "1.2")
	if info := kernelInfo(
		root,
		goosLinux,
		"future",
	); info.Kernel != "FixtureKernel" || info.KernelRelease != "1.2" ||
		info.Architecture != "future" {
		t.Fatal(info)
	}
	if hostname(func() (string, error) { return "fixture", nil }) != "fixture" ||
		hostname(func() (string, error) { return "", errors.New("unavailable") }) != "" {
		t.Fatal("hostname fallback")
	}
}
