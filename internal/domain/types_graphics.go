// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package domain

import (
	"github.com/gostafa/linuxdesktop/internal/schema"
)

type (
	// GraphicsInfo describes the GPUs and the client-side graphics stack.
	GraphicsInfo = schema.GraphicsInfo[OpenGLInfo, VulkanInfo, GPUInfo]

	// GPUInfo is one DRM device and its PCI identity.
	GPUInfo struct {
		// ID identifier reported by the underlying system.
		ID string `json:"id"`
		// Vendor vendor name reported by the driver or server.
		Vendor string `json:"vendor"`
		// Model marketing name from the driver or PCI device database.
		Model string `json:"model"`
		// VendorID pCI vendor identifier as reported by sysfs.
		VendorID string `json:"vendor_id"`
		// DeviceID pCI device identifier as reported by sysfs.
		DeviceID string `json:"device_id"`
		// PCIAddress pCI slot address, when the GPU is on a PCI bus.
		PCIAddress string `json:"pci_address"`
		// Driver kernel driver bound to the device.
		Driver string `json:"driver"`
		// DRMDevice path to the DRM card device.
		DRMDevice string `json:"drm_device"`
		// RenderDevice path to the corresponding DRM render node.
		RenderDevice string `json:"render_device"`
		// Integrated whether the device is classified as integrated.
		Integrated bool `json:"integrated"`
		// Discrete whether the device is classified as discrete.
		Discrete bool `json:"discrete"`
	}

	// OpenGLInfo describes the OpenGL implementation.
	OpenGLInfo struct {
		// Vendor vendor name reported by the driver or server.
		Vendor string `json:"vendor,omitempty"`
		// Renderer openGL renderer string from the live driver.
		Renderer string `json:"renderer,omitempty"`
		// Version version reported by the source, when available.
		Version string `json:"version,omitempty"`
		// Available whether the implementation or service was detected.
		Available bool `json:"available"`
	}

	// VulkanInfo describes the Vulkan loader and its ICDs.
	VulkanInfo struct {
		// Version version reported by the source, when available.
		Version string `json:"version,omitempty"`
		// Available whether the implementation or service was detected.
		Available bool `json:"available"`
	}
)
