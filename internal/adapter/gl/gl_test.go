// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package gl

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gostafa/linuxdesktop/internal/domain"
	"github.com/gostafa/linuxdesktop/internal/sysfs"
)

func write(t *testing.T, files string, path, data string) {
	t.Helper()
	path = sysfs.Path(files, path)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestManifestStack(t *testing.T) {
	vulkanICDDirs := vulkanICDDirs()
	eglVendorDirs := eglVendorDirs()
	driDirs := driDirs()

	t.Parallel()
	files := t.TempDir()
	if opengl, vulkan, err := detectStack(t.Context(), files); err != nil || opengl.Available ||
		vulkan.Available {
		t.Fatal(opengl, vulkan, err)
	}
	write(t, files, filepath.Join(driDirs[0], "irrelevant"), "")
	if hasDRIDriver(files, driDirs[0]) {
		t.Fatal("non-driver accepted")
	}
	write(t, files, filepath.Join(driDirs[0], "fixture_dri.so"), "")
	if got := openglFromManifests(files); !got.Available || got.Vendor != vendorMesa {
		t.Fatal(got)
	}
	write(
		t,
		files,
		filepath.Join(eglVendorDirs[0], "vendor.json"),
		`{"ICD":{"library_path":"libEGL_nvidia.so"}}`,
	)
	write(t, files, filepath.Join(eglVendorDirs[0], "ignored.txt"), "")
	write(t, files, filepath.Join(eglVendorDirs[0], "broken.json"), "{")
	write(
		t,
		files,
		filepath.Join(vulkanICDDirs[0], "vendor.json"),
		`{"ICD":{"api_version":"1.3.9"}}`,
	)
	opengl, vulkan, err := detectStack(t.Context(), files)
	if err != nil || !opengl.Available || opengl.Vendor != "NVIDIA" || !vulkan.Available ||
		vulkan.Version != "1.3.9" {
		t.Fatal(opengl, vulkan, err)
	}
	if _, ok := readManifest(files, eglVendorDirs[0], "missing.json"); ok {
		t.Fatal("missing manifest")
	}
	if _, _, err = New().Stack(t.Context()); err != nil {
		t.Fatal(err)
	}
	if _, _, err = NewNative().Stack(t.Context()); err != nil {
		t.Fatal(err)
	}
}

func TestVersionsAndLiveResults(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ current, candidate, want string }{
		{"", "", ""},
		{"", "1.3", "1.3"},
		{"1.3", "", "1.3"},
		{"1.3", "1.3.0", "1.3"},
		{"1.3", "1.10", "1.10"},
		{"1.10", "1.3", "1.10"},
		{"1.3", "1.bad", "1.3"},
	} {
		if got := newer(tc.current, tc.candidate); got != tc.want {
			t.Fatal(tc, got)
		}
	}
	if vendorOf("unknown") != "" || orVendor("known", "new") != "known" {
		t.Fatal("vendor fallback")
	}
	opengl := domain.OpenGLInfo{Vendor: "manifest", Renderer: "retained"}
	var vulkan domain.VulkanInfo
	augmentVulkan(&vulkan, func() (string, bool) { return "1.4", true })
	augmentGL(&opengl, func() (domain.OpenGLInfo, bool) {
		return domain.OpenGLInfo{Available: true, Vendor: "live", Version: "4.6"}, true
	})
	if !vulkan.Available || vulkan.Version != "1.4" || !opengl.Available ||
		opengl.Vendor != "live" ||
		opengl.Renderer != "retained" ||
		opengl.Version != "4.6" {
		t.Fatal(opengl, vulkan)
	}
	augmentVulkan(&vulkan, func() (string, bool) { return "", false })
	augmentGL(&opengl, func() (domain.OpenGLInfo, bool) { return domain.OpenGLInfo{}, false })
	if vulkan.Version != "1.4" || opengl.Version != "4.6" {
		t.Fatal("unavailable drivers overwrote manifests")
	}
}
