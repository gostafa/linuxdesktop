package drm

// Sysfs and devfs layout.
const (
	classDir  = "/sys/class/drm"
	deviceDir = "/dev/dri"

	cardPrefix   = "card"
	renderPrefix = "renderD"

	attrVendor  = "device/vendor"
	attrDevice  = "device/device"
	attrUevent  = "device/uevent"
	attrBootVGA = "device/boot_vga"
	linkDevice  = "device"
)

// uevent keys.
const (
	keyDriver = "DRIVER"
	keySlot   = "PCI_SLOT_NAME"
)

// rootBusPrefix marks a device attached directly to the PCI root complex,
// which in practice means an integrated GPU.
const rootBusPrefix = "0000:00:"

// NVIDIA's proprietary driver publishes marketing names here, which saves
// streaming pci.ids for the most common discrete vendor.
const (
	nvidiaGPUDir   = "/proc/driver/nvidia/gpus"
	nvidiaInfo     = "information"
	nvidiaModel    = "Model:"
	driverNvidia   = "nvidia"
	vendorIDNvidia = "0x10de"
)

// pci.ids parsing.
const (
	idFieldWidth = 4
	commentByte  = '#'
	indentByte   = '\t'
)
