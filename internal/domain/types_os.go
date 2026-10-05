// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package domain

type (
	// OSInfo describes the operating system and kernel.
	OSInfo struct {
		// Name human-readable name reported by the source.
		Name string `json:"name"`
		// PrettyName human-readable distribution name from os-release.
		PrettyName string `json:"pretty_name"`
		// ID identifier reported by the underlying system.
		ID string `json:"id"`
		// IDLike related distribution identifiers from os-release.
		IDLike string `json:"id_like"`
		// Version version reported by the source, when available.
		Version string `json:"version"`
		// VersionID machine-readable distribution version.
		VersionID string `json:"version_id"`
		// Kernel kernel name.
		Kernel string `json:"kernel"`
		// KernelRelease kernel release string.
		KernelRelease string `json:"kernel_release"`
		// Architecture architecture in uname-style spelling.
		Architecture string `json:"architecture"`
		// Hostname host name reported by the operating system.
		Hostname string `json:"hostname"`
	}
)
