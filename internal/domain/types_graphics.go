// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package domain

type (
	// GraphicsInfo describes the GPUs and the client-side graphics stack.
	GraphicsInfo struct {
		OpenGL     OpenGLInfo `json:"opengl"`
		Vulkan     VulkanInfo `json:"vulkan"`
		PrimaryGPU string     `json:"primary_gpu,omitempty"`
		GPUs       []GPUInfo  `json:"gpus"`
	}

	// GPUInfo is one DRM device and its PCI identity.
	GPUInfo struct {
		ID           string `json:"id"`
		Vendor       string `json:"vendor"`
		Model        string `json:"model"`
		VendorID     string `json:"vendor_id"`
		DeviceID     string `json:"device_id"`
		PCIAddress   string `json:"pci_address"`
		Driver       string `json:"driver"`
		DRMDevice    string `json:"drm_device"`
		RenderDevice string `json:"render_device"`
		Integrated   bool   `json:"integrated"`
		Discrete     bool   `json:"discrete"`
	}

	// OpenGLInfo describes the OpenGL implementation.
	OpenGLInfo struct {
		Vendor    string `json:"vendor,omitempty"`
		Renderer  string `json:"renderer,omitempty"`
		Version   string `json:"version,omitempty"`
		Available bool   `json:"available"`
	}

	// VulkanInfo describes the Vulkan loader and its ICDs.
	VulkanInfo struct {
		Version   string `json:"version,omitempty"`
		Available bool   `json:"available"`
	}
)
