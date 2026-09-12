package drm

// vendorNames covers every vendor that ships a Linux GPU driver, plus the
// virtual devices a desktop is likely to run under. Using a table instead of
// /usr/share/hwdata/pci.ids keeps vendor resolution free, and keeps working on
// systems where hwdata is not installed.
var vendorNames = map[string]string{
	"0x1002": "AMD",
	"0x1022": "AMD",
	"0x10de": "NVIDIA",
	"0x12d2": "NVIDIA",
	"0x8086": "Intel",
	"0x1af4": "Red Hat",
	"0x1b36": "Red Hat",
	"0x15ad": "VMware",
	"0x1234": "Bochs",
	"0x80ee": "Oracle",
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
	"0x102b": "Matrox",
	"0x1013": "Cirrus Logic",
}

// alwaysDiscrete lists vendors whose PCI GPUs are discrete by construction.
// NVIDIA's integrated parts (Tegra) are not PCI devices, so a PCI NVIDIA GPU
// is always a discrete card regardless of which bus it sits on.
var alwaysDiscrete = map[string]bool{
	"0x10de": true,
}

// pciIDsPaths are the places distributions install the hwdata device list.
var pciIDsPaths = []string{
	"/usr/share/hwdata/pci.ids",
	"/usr/share/misc/pci.ids",
	"/usr/share/pci.ids",
}
