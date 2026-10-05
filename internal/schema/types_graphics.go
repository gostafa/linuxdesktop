// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package schema

type (
	// GraphicsInfo is the shared GraphicsInfo layout, parameterized by its component types.
	GraphicsInfo[OpenGLType, VulkanType, GPUType any] struct {
		// OpenGL openGL implementation reported by drivers or manifests.
		OpenGL OpenGLType `json:"opengl"`
		// Vulkan vulkan loader and driver availability.
		Vulkan VulkanType `json:"vulkan"`
		// PrimaryGPU identifier of the firmware-posted GPU, or the first detected GPU.
		PrimaryGPU string `json:"primary_gpu,omitempty"`
		// GPUs detected DRM devices with their PCI identity.
		GPUs []GPUType `json:"gpus"`
	}
)
