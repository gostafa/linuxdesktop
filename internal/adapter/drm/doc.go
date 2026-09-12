// Package drm enumerates GPUs from /sys/class/drm.
//
// Every card node points at its underlying device, whose PCI identity, driver
// and boot-VGA flag are all plain files. Render nodes are matched back to
// their card by resolving both device symlinks to the same address, which is
// how a caller gets the /dev/dri/renderD* path it actually needs for
// headless GPU work.
//
// Vendor names come from a built-in table rather than from
// /usr/share/hwdata/pci.ids. The table covers every GPU vendor that ships a
// Linux driver and costs nothing, where the hwdata file is around 1.6 MB.
//
// Model names do come from pci.ids, because there is no other source for them,
// but the file is streamed in one pass that stops as soon as every GPU has
// been resolved — for the usual one or two cards that is a few hundred
// kilobytes, not the whole file. NVIDIA is short-circuited entirely via
// /proc/driver/nvidia, which states the marketing name directly. If pci.ids is
// not installed, Model falls back to the raw device id.
//
// Integrated versus discrete is a heuristic, and is documented as such: a
// device on PCI root bus 0000:00 is treated as integrated, anything further
// out as discrete, with NVIDIA forced discrete because its integrated parts
// are not PCI devices at all.
package drm
