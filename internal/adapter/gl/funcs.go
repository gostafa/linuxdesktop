// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package gl

import (
	"context"
	"encoding/json"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/gostafa/linuxdesktop/internal/adapter/gl/native"
	"github.com/gostafa/linuxdesktop/internal/domain"
	"github.com/gostafa/linuxdesktop/internal/sysfs"
)

// New returns a graphics stack probe that reads installed manifests.
func New() Probe {
	return func(ctx context.Context) (domain.OpenGLInfo, domain.VulkanInfo, error) { return stack(ctx, "") }
}

// NewNative returns a stack probe that also queries native graphics drivers.
func NewNative() Probe {
	return func(ctx context.Context) (domain.OpenGLInfo, domain.VulkanInfo, error) { return nativeStack(ctx, "") }
}

// stack reports the OpenGL and Vulkan implementations from installed manifests.
func stack(_ context.Context, files string) (domain.OpenGLInfo, domain.VulkanInfo, error) {
	opengl := openglFromManifests(files)
	vulkan := vulkanFromManifests(files)

	return opengl, vulkan, nil
}

// nativeStack augments the manifest result with details from live drivers.
func nativeStack(_ context.Context, files string) (domain.OpenGLInfo, domain.VulkanInfo, error) {
	opengl := openglFromManifests(files)
	vulkan := vulkanFromManifests(files)

	augmentGL(&opengl, native.OpenGL)
	augmentVulkan(&vulkan, native.Vulkan)

	return opengl, vulkan, nil
}

func augmentVulkan(vulkan *domain.VulkanInfo, read func() (string, bool)) {
	if version, ok := read(); ok {
		vulkan.Available = true
		vulkan.Version = version
	}
}

func augmentGL(opengl *domain.OpenGLInfo, read func() (domain.OpenGLInfo, bool)) {
	if live, ok := read(); ok {
		merge(opengl, &live)
	}
}

// openglFromManifests establishes that OpenGL is installed, and infers the
// vendor from whichever driver library the manifests name.
func openglFromManifests(files string) domain.OpenGLInfo {
	info := fromEGLVendors(files)
	if info.Available {
		return info
	}

	// No glvnd manifests: fall back to the presence of Mesa's DRI drivers.
	if hasDRIDrivers(files) {
		return domain.OpenGLInfo{
			Available: true,
			Vendor:    vendorMesa,
			Renderer:  "",
			Version:   "",
		}
	}

	return info
}

// fromEGLVendors reads the glvnd manifests whose presence is what makes OpenGL
// usable through glvnd at all.
func fromEGLVendors(files string) domain.OpenGLInfo {
	eglVendorDirs := eglVendorDirs()

	var info domain.OpenGLInfo

	eachManifest(files, eglVendorDirs, eglVendor(&info))

	return info
}

func eglVendor(info *domain.OpenGLInfo) func(*manifest) {
	return func(found *manifest) {
		info.Available = true
		info.Vendor = orVendor(info.Vendor, found.ICD.LibraryPath)
	}
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
func hasDRIDrivers(files string) bool {
	driDirs := driDirs()

	return slices.ContainsFunc(driDirs, func(dir string) bool { return hasDRIDriver(files, dir) })
}

func hasDRIDriver(files, dir string) bool {
	names, err := sysfs.DirNames(sysfs.Path(files, dir))
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
func vulkanFromManifests(files string) domain.VulkanInfo {
	vulkanICDDirs := vulkanICDDirs()

	var info domain.VulkanInfo

	eachManifest(files, vulkanICDDirs, func(found *manifest) {
		info.Available = true
		info.Version = newer(info.Version, found.ICD.APIVersion)
	})

	return info
}

func eachManifest(files string, dirs []string, visit func(*manifest)) {
	for i := range dirs {
		found := readManifests(files, dirs[i])
		for idx := range found {
			visit(&found[idx])
		}
	}
}

// newer keeps whichever of two dotted versions is the higher one.
func newer(current, candidate string) string {
	if higher(candidate, current) {
		return candidate
	}

	return current
}

func readManifests(files, dir string) []manifest {
	names, err := sysfs.DirNames(sysfs.Path(files, dir))
	if err != nil {
		return nil
	}

	out := make([]manifest, zero, len(names))

	for i := range names {
		found, ok := readManifest(files, dir, names[i])
		if !ok {
			continue
		}

		out = append(out, found)
	}

	return out
}

// readManifest decodes one manifest, skipping anything that is not JSON or does
// not parse.
func readManifest(files, dir, name string) (manifest, bool) {
	var empty manifest

	if !strings.HasSuffix(name, jsonSuffix) {
		return empty, false
	}

	data, err := sysfs.Bytes(sysfs.Path(files, filepath.Join(dir, name)))
	if err != nil {
		return empty, false
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
	driverVendors := driverVendors()

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
