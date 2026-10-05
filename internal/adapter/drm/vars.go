// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package drm

// vendorNames covers every vendor that ships a Linux GPU driver, plus the
// virtual devices a desktop is likely to run under. Using a table instead of
// /usr/share/hwdata/pci.ids keeps vendor resolution free, and keeps working on
// systems where hwdata is not installed.
func vendorNames(key string) string {
	tables := []map[string]string{vendorNamesTableA(), vendorNamesTableB(), vendorNamesTableC()}

	for i := range tables {
		if found, exists := tables[i][key]; exists {
			return found
		}
	}

	return noValue
}

func vendorNamesTableA() map[string]string {
	return map[string]string{
		"0x1002":       vendorAMD,
		"0x1022":       vendorAMD,
		vendorIDNvidia: vendorNvidia,
		"0x12d2":       vendorNvidia,
		"0x8086":       "Intel",
		"0x1af4":       vendorRedHat,
		"0x1b36":       vendorRedHat,
		"0x15ad":       "VMware",
		"0x1234":       "Bochs",
		"0x80ee":       "Oracle",
	}
}

func vendorNamesTableB() map[string]string {
	return map[string]string{
		"0x1414": "Microsoft",
		"0x13b5": "ARM",
		"0x5143": "Qualcomm",
		"0x14e4": "Broadcom",
		"0x1d17": "Zhaoxin",
		"0x1ed5": "Moore Threads",
		"0x1e4b": "Innosilicon",
		"0x0014": "Loongson",
		"0x108e": "Sun",
		"0x1a03": "ASPEED",
	}
}

func vendorNamesTableC() map[string]string {
	return map[string]string{
		"0x102b": "Matrox",
		"0x1013": "Cirrus Logic",
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
