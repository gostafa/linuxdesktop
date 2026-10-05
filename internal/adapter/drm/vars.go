// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package drm

// vendorNames covers every vendor that ships a Linux GPU driver, plus the
// virtual devices a desktop is likely to run under. Using a table instead of
// /usr/share/hwdata/pci.ids keeps vendor resolution free, and keeps working on
// systems where hwdata is not installed.
func vendorNames(key string) string {
	switch key {
	case "0x1002", "0x1022":
		return vendorAMD
	case vendorIDNvidia, "0x12d2":
		return vendorNvidia
	case "0x8086":
		return "Intel"
	case "0x1af4", "0x1b36":
		return vendorRedHat
	case "0x15ad":
		return "VMware"
	case "0x1234":
		return "Bochs"
	case "0x80ee":
		return "Oracle"
	case "0x1414":
		return "Microsoft"
	case "0x13b5":
		return "ARM"
	case "0x5143":
		return "Qualcomm"
	case "0x14e4":
		return "Broadcom"
	case "0x1d17":
		return "Zhaoxin"
	case "0x1ed5":
		return "Moore Threads"
	case "0x1e4b":
		return "Innosilicon"
	case "0x0014":
		return "Loongson"
	case "0x108e":
		return "Sun"
	case "0x1a03":
		return "ASPEED"
	case "0x102b":
		return "Matrox"
	case "0x1013":
		return "Cirrus Logic"
	default:
		return noValue
	}
}

// alwaysDiscrete lists vendors whose PCI GPUs are discrete by construction.
// NVIDIA's integrated parts (Tegra) are not PCI devices, so a PCI NVIDIA GPU
// is always a discrete card regardless of which bus it sits on.
func alwaysDiscrete(key string) bool { return key == vendorIDNvidia }

// pciIDsPaths are the places distributions install the hwdata device list.
func pciIDsPaths() []string {
	return []string{
		"/usr/share/hwdata/pci.ids",
		"/usr/share/misc/pci.ids",
		"/usr/share/pci.ids",
	}
}
