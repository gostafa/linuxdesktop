// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package domain

type (
	// OSInfo describes the operating system and kernel.
	OSInfo struct {
		Name          string `json:"name"`
		PrettyName    string `json:"pretty_name"`
		ID            string `json:"id"`
		IDLike        string `json:"id_like"`
		Version       string `json:"version"`
		VersionID     string `json:"version_id"`
		Kernel        string `json:"kernel"`
		KernelRelease string `json:"kernel_release"`
		Architecture  string `json:"architecture"`
		Hostname      string `json:"hostname"`
	}
)
