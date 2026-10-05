// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package drm

const (
	// Sysfs and devfs layout.
	classDir  = "/sys/class/drm"
	deviceDir = "/dev/dri"

	cardPrefix   = "card"
	renderPrefix = "renderD"

	attrVendor  = "device/vendor"
	attrDevice  = "device/device"
	attrUevent  = "device/uevent"
	attrBootVGA = "device/boot_vga"
	linkDevice  = "device"

	// uevent keys.
	keyDriver = "DRIVER"
	keySlot   = "PCI_SLOT_NAME"

	// rootBusPrefix marks a device attached directly to the PCI root complex,
	// which in practice means an integrated GPU.
	rootBusPrefix = "0000:00:"

	// NVIDIA's proprietary driver publishes marketing names here, which saves
	// streaming pci.ids for the most common discrete vendor.
	nvidiaGPUDir   = "/proc/driver/nvidia/gpus"
	nvidiaInfo     = "information"
	nvidiaModel    = "Model:"
	driverNvidia   = "nvidia"
	vendorIDNvidia = "0x10de"

	// pci.ids parsing. A vendor line starts at column zero and a device line is
	// indented by one tab; a second tab marks a subsystem line, which is of no
	// interest here.
	idFieldWidth = 4
	commentByte  = '#'
	indentByte   = '\t'
	indentWidth  = 1

	// expectedGPUs is the capacity every per-card collection is built at. Two
	// covers the laptop that pairs an integrated GPU with a discrete one.
	expectedGPUs = 2

	// zero is the empty length, the index of a leading byte, and an exhausted
	// count. It is one constant rather than three because a package may not
	// declare two constants sharing a value.
	zero = 0

	// noValue is an unreadable sysfs attribute or an unnamed device.
	noValue = ""
)
