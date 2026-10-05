// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package gl

import (
	"github.com/gostafa/linuxdesktop/internal/domain"
	"github.com/gostafa/linuxdesktop/internal/probe"
)

type (
	// Probe implements port.StackProbe.
	//
	// native selects the dlopen path, which asks the drivers directly at the
	// cost of initializing them. When it is false the probe reads manifests
	// only.
	Probe = probe.StackFunc[domain.OpenGLInfo, domain.VulkanInfo]

	// icd is the driver entry a manifest describes: where the library lives
	// and, for Vulkan, which API version it implements.
	icd struct {
		// LibraryPath driver library named by the ICD manifest.
		LibraryPath string `json:"library_path"`
		// APIVersion vulkan API version claimed by the driver.
		APIVersion string `json:"api_version"`
	}

	// manifest is the shape shared by Vulkan ICD files and glvnd EGL vendor
	// files.
	manifest = driverManifest[icd]

	// driverManifest decodes a driver entry shared by EGL and Vulkan files.
	driverManifest[D any] struct {
		// ICD driver entry shared by EGL and Vulkan manifests.
		ICD D `json:"ICD"`
	}
)
