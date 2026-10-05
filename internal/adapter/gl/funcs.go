// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package gl

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gostafa/linuxdesktop/internal/adapter/gl/native"
	"github.com/gostafa/linuxdesktop/internal/domain"
	"github.com/gostafa/linuxdesktop/internal/sysfs"
)

// New returns a stack probe. Pass native to enable the dlopen path.
func New(useNative bool) Probe { return Probe{native: useNative} }

// Stack reports the OpenGL and Vulkan implementations. The manifest scan
// always runs, so the native path only ever adds detail, never removes it.
func (probe Probe) Stack(_ context.Context) (domain.OpenGLInfo, domain.VulkanInfo, error) {
	opengl := openglFromManifests()
	vulkan := vulkanFromManifests()

	if !probe.native {
		return opengl, vulkan, nil
	}

	if version, ok := native.Vulkan(); ok {
		vulkan.Available = true
		vulkan.Version = version
	}

	if live, ok := native.OpenGL(); ok {
		merge(&opengl, &live)
	}

	return opengl, vulkan, nil
}

// openglFromManifests establishes that OpenGL is installed, and infers the
// vendor from whichever driver library the manifests name.
func openglFromManifests() domain.OpenGLInfo {
	info := fromEGLVendors()
	if info.Available {
		return info
	}

	// No glvnd manifests: fall back to the presence of Mesa's DRI drivers.
	if hasDRIDrivers() {
		return domain.OpenGLInfo{Available: true, Vendor: vendorMesa}
	}

	return info
}

// fromEGLVendors reads the glvnd manifests whose presence is what makes OpenGL
// usable through glvnd at all.
func fromEGLVendors() domain.OpenGLInfo {
	var info domain.OpenGLInfo

	for i := range eglVendorDirs {
		found := readManifests(eglVendorDirs[i])
		for idx := range found {
			info.Available = true
			info.Vendor = orVendor(info.Vendor, found[idx].ICD.LibraryPath)
		}
	}

	return info
}

// orVendor keeps the first vendor a manifest named, since the first glvnd
// manifest is the one glvnd itself would load.
func orVendor(known, library string) string {
	if known != noValue {
		return known
	}

	return vendorOf(library)
}

// hasDRIDrivers reports whether Mesa's DRI drivers are installed, which is the
// fallback evidence that OpenGL is present on a system without glvnd.
func hasDRIDrivers() bool {
	for i := range driDirs {
		if hasDRIDriver(driDirs[i]) {
			return true
		}
	}

	return false
}

func hasDRIDriver(dir string) bool {
	names, err := sysfs.DirNames(dir)
	if err != nil {
		return false
	}

	for i := range names {
		if strings.HasSuffix(names[i], driSuffix) {
			return true
		}
	}

	return false
}

// vulkanFromManifests reports the highest API version any installed ICD claims.
func vulkanFromManifests() domain.VulkanInfo {
	var info domain.VulkanInfo

	for i := range vulkanICDDirs {
		found := readManifests(vulkanICDDirs[i])
		for idx := range found {
			info.Available = true
			info.Version = newer(info.Version, found[idx].ICD.APIVersion)
		}
	}

	return info
}

// newer keeps whichever of two dotted versions is the higher one.
func newer(current, candidate string) string {
	if higher(candidate, current) {
		return candidate
	}

	return current
}

func readManifests(dir string) []manifest {
	names, err := sysfs.DirNames(dir)
	if err != nil {
		return nil
	}

	out := make([]manifest, zero, len(names))

	for i := range names {
		found, ok := readManifest(dir, names[i])
		if !ok {
			continue
		}

		out = append(out, found)
	}

	return out
}

// readManifest decodes one manifest, skipping anything that is not JSON or does
// not parse.
func readManifest(dir, name string) (manifest, bool) {
	if !strings.HasSuffix(name, jsonSuffix) {
		return manifest{}, false
	}

	data, err := sysfs.Bytes(filepath.Join(dir, name))
	if err != nil {
		return manifest{}, false
	}

	var found manifest

	err = json.Unmarshal(data, &found)

	return found, err == nil
}

// merge lets a live answer override a manifest-derived one field by field, so
// a partial native result still improves the outcome.
func merge(into, live *domain.OpenGLInfo) {
	into.Available = into.Available || live.Available
	into.Vendor = orKeep(into.Vendor, live.Vendor)
	into.Renderer = orKeep(into.Renderer, live.Renderer)
	into.Version = orKeep(into.Version, live.Version)
}

// orKeep prefers what the driver said, keeping the manifest answer when the
// driver supplied nothing.
func orKeep(known, live string) string {
	if live != noValue {
		return live
	}

	return known
}

func vendorOf(library string) string {
	lower := strings.ToLower(filepath.Base(library))

	for i := range driverVendors {
		if strings.Contains(lower, driverVendors[i].needle) {
			return driverVendors[i].vendor
		}
	}

	return noValue
}

// higher compares two dotted version strings numerically, so that 1.3.283
// beats 1.10 rather than losing to it lexically.
func higher(candidate, current string) bool {
	if candidate == noValue {
		return false
	}

	if current == noValue {
		return true
	}

	left := strings.Split(candidate, versionSeparator)
	right := strings.Split(current, versionSeparator)

	return compare(left, right) > zero
}

// compare orders two split version strings field by field, returning a positive
// number when the first is the higher one.
func compare(left, right []string) int {
	for i := range max(len(left), len(right)) {
		diff := part(left, i) - part(right, i)
		if diff != zero {
			return diff
		}
	}

	return zero
}

// part is one numeric field of a version, treating a missing or unparseable
// field as zero so that "1.3" and "1.3.0" compare equal.
func part(fields []string, index int) int {
	if index >= len(fields) {
		return zero
	}

	number, err := strconv.Atoi(fields[index])
	if err != nil {
		return zero
	}

	return number
}
