// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package gl

type (
	// Probe implements port.StackProbe.
	//
	// native selects the dlopen path, which asks the drivers directly at the
	// cost of initializing them. When it is false the probe reads manifests
	// only.
	Probe struct {
		native bool
	}

	// icd is the driver entry a manifest describes: where the library lives
	// and, for Vulkan, which API version it implements.
	icd struct {
		LibraryPath string `json:"library_path"`
		APIVersion  string `json:"api_version"`
	}

	// manifest is the shape shared by Vulkan ICD files and glvnd EGL vendor
	// files.
	manifest struct {
		ICD icd `json:"ICD"`
	}
)
