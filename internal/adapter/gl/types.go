package gl

// Probe implements port.StackProbe.
//
// native selects the dlopen path, which asks the drivers directly at the cost
// of initialising them. When it is false the probe reads manifests only.
type Probe struct {
	native bool
}

// manifest is the shape shared by Vulkan ICD files and glvnd EGL vendor files.
type manifest struct {
	ICD struct {
		LibraryPath string `json:"library_path"`
		APIVersion  string `json:"api_version"`
	} `json:"ICD"`
}
