// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package drm

import (
	"bufio"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gostafa/linuxdesktop/internal/domain"
	"github.com/gostafa/linuxdesktop/internal/sysfs"
)

// New returns a DRM probe.
func New() Probe { return Probe{} }

// GPUs enumerates the DRM cards and returns them along with the id of the
// primary one, which is the card the firmware posted. A machine with no DRM
// devices is not an error: that is what a headless server looks like.
func (Probe) GPUs(_ context.Context) ([]domain.GPUInfo, string, error) {
	entries, err := sysfs.DirNames(classDir)
	if err != nil {
		return nil, noValue, nil
	}

	slices.Sort(entries)

	gpus, primary := collect(entries)

	err = resolveModels(gpus)

	return gpus, orFirst(primary, gpus), err
}

// newPCIScan prepares a scan for the names the NVIDIA listing did not supply.
func newPCIScan(pending map[string]map[string][]int, gpus []domain.GPUInfo) *pciScan {
	return &pciScan{pending: pending, gpus: gpus, left: countPending(pending)}
}

// device handles one indented device line under the vendor currently selected,
// reporting true when nothing is left to name.
func (scan *pciScan) device(entry []byte) bool {
	if scan.devices == nil || len(entry) < idFieldWidth || entry[zero] == indentByte {
		return false // no interest in this vendor, or a subsystem line
	}

	id := string(entry[:idFieldWidth])

	indices, ok := scan.devices[id]
	if !ok {
		return false
	}

	scan.name(indices, strings.TrimSpace(string(entry[idFieldWidth:])))
	delete(scan.devices, id)

	return scan.left <= zero
}

// line handles one line of the listing, reporting true when nothing is left to
// name. A line flush against the margin selects a vendor; an indented one
// names a device under it.
func (scan *pciScan) line(line []byte) bool {
	if len(line) < idFieldWidth || line[zero] == commentByte {
		return false
	}

	if line[zero] != indentByte {
		scan.devices = scan.pending[string(line[:idFieldWidth])]

		return false
	}

	return scan.device(line[indentWidth:])
}

// name assigns a model to every GPU that was waiting on one device id.
func (scan *pciScan) name(indices []int, model string) {
	for i := range indices {
		scan.gpus[indices[i]].Model = model
	}

	scan.left -= len(indices)
}

// run reads the listing, or as much of it as it takes to name every GPU that is
// still waiting.
func (scan *pciScan) run(file io.Reader) error {
	lines := bufio.NewScanner(file)
	for lines.Scan() {
		if scan.line(lines.Bytes()) {
			return nil
		}
	}

	return lines.Err()
}

// collect builds one GPUInfo per DRM card and names the one the firmware
// posted, which is the only card that carries a boot_vga marker.
func collect(entries []string) ([]domain.GPUInfo, string) {
	nodes := renderNodesByAddress(entries)
	gpus := slices.Grow([]domain.GPUInfo(nil), expectedGPUs)
	primary := noValue

	found := cards(entries)
	for i := range found {
		gpu, boot := describe(found[i], nodes)

		gpus = append(gpus, gpu)

		if boot {
			primary = gpu.ID
		}
	}

	return gpus, primary
}

// cards filters the DRM class directory down to the card devices, leaving out
// the connectors and render nodes that share it.
func cards(entries []string) []string {
	found := make([]string, zero, len(entries))

	for i := range entries {
		if isCard(entries[i]) {
			found = append(found, entries[i])
		}
	}

	return found
}

// describe reads everything sysfs knows about one card, and reports whether the
// firmware posted it.
func describe(entry string, nodes map[string]string) (domain.GPUInfo, bool) {
	base := filepath.Join(classDir, entry)
	gpu := domain.GPUInfo{
		ID:         entry,
		VendorID:   sysfs.Trimmed(filepath.Join(base, attrVendor)),
		DeviceID:   sysfs.Trimmed(filepath.Join(base, attrDevice)),
		PCIAddress: sysfs.LinkBase(filepath.Join(base, linkDevice)),
		DRMDevice:  filepath.Join(deviceDir, entry),
	}

	applyUevent(&gpu, base)

	gpu.RenderDevice = nodes[gpu.PCIAddress]
	gpu.Vendor = vendorNames[gpu.VendorID]

	classify(&gpu)

	return gpu, sysfs.Trimmed(filepath.Join(base, attrBootVGA)) == "1"
}

// applyUevent overlays the driver name and the PCI slot the kernel reports,
// which is more reliable than the device symlink when both are present.
func applyUevent(gpu *domain.GPUInfo, base string) {
	data, err := sysfs.Bytes(filepath.Join(base, attrUevent))
	if err != nil {
		return
	}

	gpu.Driver = sysfs.Field(data, keyDriver)

	if slot := sysfs.Field(data, keySlot); slot != noValue {
		gpu.PCIAddress = slot
	}
}

// orFirst falls back to the first card when the firmware left no boot_vga
// marker, since a machine with a GPU at all has a primary one.
func orFirst(primary string, gpus []domain.GPUInfo) string {
	if primary != noValue || len(gpus) == zero {
		return primary
	}

	return gpus[zero].ID
}

// classify decides whether a GPU is integrated or discrete. Both flags stay
// false when the device is not on a PCI bus at all, which is the honest answer
// for an SoC display engine.
func classify(gpu *domain.GPUInfo) {
	if gpu.PCIAddress == noValue || !strings.Contains(gpu.PCIAddress, ":") {
		gpu.Integrated = true

		return
	}

	if alwaysDiscrete[gpu.VendorID] {
		gpu.Discrete = true

		return
	}

	if strings.HasPrefix(gpu.PCIAddress, rootBusPrefix) {
		gpu.Integrated = true

		return
	}

	gpu.Discrete = true
}

// renderNodesByAddress maps a PCI address to its /dev/dri/renderD* path, by
// resolving each render node's device link to the same device the card points
// at.
func renderNodesByAddress(entries []string) map[string]string {
	nodes := make(map[string]string, expectedGPUs)

	for i := range entries {
		if !strings.HasPrefix(entries[i], renderPrefix) {
			continue
		}

		address := sysfs.LinkBase(filepath.Join(classDir, entries[i], linkDevice))
		if address == noValue {
			continue
		}

		nodes[address] = filepath.Join(deviceDir, entries[i])
	}

	return nodes
}

// resolveModels fills in Model for every GPU, preferring NVIDIA's own listing
// and otherwise streaming pci.ids exactly once.
func resolveModels(gpus []domain.GPUInfo) error {
	pending := pendingModels(gpus)
	if len(pending) == zero {
		return nil
	}

	file := openPCIIDs()
	if file == nil {
		return nil
	}

	scan := newPCIScan(pending, gpus)

	return errors.Join(scan.run(file), file.Close())
}

// pendingModels names every GPU it can from NVIDIA's own listing and returns
// the pci.ids lookups still wanted, keyed by vendor and then by device id.
func pendingModels(gpus []domain.GPUInfo) map[string]map[string][]int {
	pending := make(map[string]map[string][]int, expectedGPUs)

	for i := range gpus {
		vendor, device, ok := modelLookup(&gpus[i])
		if !ok {
			continue
		}

		if pending[vendor] == nil {
			pending[vendor] = make(map[string][]int)
		}

		pending[vendor][device] = append(pending[vendor][device], i)
	}

	return pending
}

// modelLookup settles where a GPU's model comes from. NVIDIA's proprietary
// driver publishes a marketing name, which wins outright; otherwise the device
// id stands in until pci.ids supplies a real one.
func modelLookup(gpu *domain.GPUInfo) (vendor, device string, ok bool) {
	model := nvidiaModelFor(gpu)
	if model != noValue {
		gpu.Model = model

		return noValue, noValue, false
	}

	vendor, device = trimHex(gpu.VendorID), trimHex(gpu.DeviceID)
	if vendor == noValue || device == noValue {
		return noValue, noValue, false
	}

	gpu.Model = gpu.DeviceID // replaced by the scan if pci.ids has a name

	return vendor, device, true
}

// countPending is how many names the scan is looking for in total.
func countPending(pending map[string]map[string][]int) int {
	left := zero

	for vendor := range pending {
		left += len(pending[vendor])
	}

	return left
}

// nvidiaModelFor reads the marketing name the proprietary driver publishes.
func nvidiaModelFor(gpu *domain.GPUInfo) string {
	if !nvidia(gpu) {
		return noValue
	}

	data, err := sysfs.Bytes(filepath.Join(nvidiaGPUDir, gpu.PCIAddress, nvidiaInfo))
	if err != nil {
		return noValue
	}

	return modelLine(data)
}

// nvidia reports whether the proprietary driver could have published a name for
// this card.
func nvidia(gpu *domain.GPUInfo) bool {
	theirs := gpu.Driver == driverNvidia || gpu.VendorID == vendorIDNvidia

	return theirs && gpu.PCIAddress != noValue
}

// modelLine finds the Model: entry in a /proc/driver/nvidia information file.
func modelLine(data []byte) string {
	for line := range strings.SplitSeq(string(data), "\n") {
		if rest, found := strings.CutPrefix(line, nvidiaModel); found {
			return strings.TrimSpace(rest)
		}
	}

	return noValue
}

func openPCIIDs() *os.File {
	for i := range pciIDsPaths {
		if file, err := os.Open(pciIDsPaths[i]); err == nil {
			return file
		}
	}

	return nil
}

func isCard(name string) bool {
	if !strings.HasPrefix(name, cardPrefix) {
		return false
	}

	// card0 is a device; card0-HDMI-A-1 is a connector.
	return !strings.ContainsRune(name[len(cardPrefix):], '-')
}

// trimHex converts a sysfs "0x8086" into the bare "8086" that pci.ids uses.
func trimHex(raw string) string {
	return strings.TrimPrefix(raw, "0x")
}
