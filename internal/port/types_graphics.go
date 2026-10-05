// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package port

import (
	"context"

	"github.com/gostafa/linuxdesktop/internal/domain"
)

type (
	// GPUProbe enumerates the machine's GPUs and names the primary one. It is
	// separate from StackProbe because the two answer genuinely different
	// questions from unrelated sources: what hardware is present, and what the
	// client-side driver stack can do with it.
	GPUProbe interface {
		GPUs(ctx context.Context) (gpus []domain.GPUInfo, primary string, err error)
	}

	// StackProbe reports the OpenGL and Vulkan implementations.
	StackProbe interface {
		Stack(ctx context.Context) (domain.OpenGLInfo, domain.VulkanInfo, error)
	}

	// PortalProbe inspects xdg-desktop-portal.
	PortalProbe interface {
		Portal(ctx context.Context, env *domain.Env) (domain.PortalInfo, error)
	}
)
