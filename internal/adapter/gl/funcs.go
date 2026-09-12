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
func (p Probe) Stack(_ context.Context) (domain.OpenGLInfo, domain.VulkanInfo, error) {
	opengl := openglFromManifests()
	vulkan := vulkanFromManifests()
	if !p.native {
		return opengl, vulkan, nil
	}
	if version, ok := native.Vulkan(); ok {
		vulkan.Available = true
		vulkan.Version = version
	}
	if live, ok := native.OpenGL(); ok {
		merge(&opengl, live)
	}
	return opengl, vulkan, nil
}

// openglFromManifests establishes that OpenGL is installed, and infers the
// vendor from whichever driver library the manifests name.
func openglFromManifests() domain.OpenGLInfo {
	var info domain.OpenGLInfo
	for _, dir := range eglVendorDirs {
		for _, m := range readManifests(dir) {
			info.Available = true
			if info.Vendor == "" {
				info.Vendor = vendorOf(m.ICD.LibraryPath)
			}
		}
	}
	if info.Available {
		return info
	}
	// No glvnd manifests: fall back to the presence of Mesa's DRI drivers.
	for _, dir := range driDirs {
		names, err := sysfs.DirNames(dir)
		if err != nil {
			continue
		}
		for _, name := range names {
			if strings.HasSuffix(name, driSuffix) {
				return domain.OpenGLInfo{Available: true, Vendor: vendorMesa}
			}
		}
	}
	return info
}

// vulkanFromManifests reports the highest API version any installed ICD claims.
func vulkanFromManifests() domain.VulkanInfo {
	var info domain.VulkanInfo
	for _, dir := range vulkanICDDirs {
		for _, m := range readManifests(dir) {
			info.Available = true
			if higher(m.ICD.APIVersion, info.Version) {
				info.Version = m.ICD.APIVersion
			}
		}
	}
	return info
}

func readManifests(dir string) []manifest {
	names, err := sysfs.DirNames(dir)
	if err != nil {
		return nil
	}
	out := make([]manifest, 0, len(names))
	for _, name := range names {
		if !strings.HasSuffix(name, jsonSuffix) {
			continue
		}
		data, err := sysfs.Bytes(filepath.Join(dir, name))
		if err != nil {
			continue
		}
		var m manifest
		if json.Unmarshal(data, &m) != nil {
			continue
		}
		out = append(out, m)
	}
	return out
}

// merge lets a live answer override a manifest-derived one field by field, so
// a partial native result still improves the outcome.
func merge(into *domain.OpenGLInfo, live domain.OpenGLInfo) {
	into.Available = into.Available || live.Available
	if live.Vendor != "" {
		into.Vendor = live.Vendor
	}
	if live.Renderer != "" {
		into.Renderer = live.Renderer
	}
	if live.Version != "" {
		into.Version = live.Version
	}
}

func vendorOf(library string) string {
	lower := strings.ToLower(filepath.Base(library))
	for _, v := range driverVendors {
		if strings.Contains(lower, v.needle) {
			return v.vendor
		}
	}
	return ""
}

// higher compares two dotted version strings numerically, so that 1.3.283
// beats 1.10 rather than losing to it lexically.
func higher(candidate, current string) bool {
	if candidate == "" {
		return false
	}
	if current == "" {
		return true
	}
	a := strings.Split(candidate, ".")
	b := strings.Split(current, ".")
	for i := 0; i < len(a) || i < len(b); i++ {
		x, y := part(a, i), part(b, i)
		if x != y {
			return x > y
		}
	}
	return false
}

func part(fields []string, i int) int {
	if i >= len(fields) {
		return 0
	}
	n, _ := strconv.Atoi(fields[i])
	return n
}
