package drm

import (
	"bufio"
	"context"
	"os"
	"path/filepath"
	"sort"
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
		return nil, "", nil
	}
	sort.Strings(entries)

	renderNodes := renderNodesByAddress(entries)
	gpus := make([]domain.GPUInfo, 0, 2)
	primary := ""

	for _, entry := range entries {
		if !isCard(entry) {
			continue
		}
		base := filepath.Join(classDir, entry)
		gpu := domain.GPUInfo{
			ID:         entry,
			VendorID:   sysfs.Trimmed(filepath.Join(base, attrVendor)),
			DeviceID:   sysfs.Trimmed(filepath.Join(base, attrDevice)),
			PCIAddress: sysfs.LinkBase(filepath.Join(base, linkDevice)),
			DRMDevice:  filepath.Join(deviceDir, entry),
		}

		if data, err := sysfs.Bytes(filepath.Join(base, attrUevent)); err == nil {
			gpu.Driver = sysfs.Field(data, keyDriver)
			if slot := sysfs.Field(data, keySlot); slot != "" {
				gpu.PCIAddress = slot
			}
		}
		gpu.RenderDevice = renderNodes[gpu.PCIAddress]
		gpu.Vendor = vendorNames[gpu.VendorID]
		classify(&gpu)

		if sysfs.Trimmed(filepath.Join(base, attrBootVGA)) == "1" {
			primary = gpu.ID
		}
		gpus = append(gpus, gpu)
	}

	if primary == "" && len(gpus) > 0 {
		primary = gpus[0].ID
	}
	resolveModels(gpus)
	return gpus, primary, nil
}

// classify decides whether a GPU is integrated or discrete. Both flags stay
// false when the device is not on a PCI bus at all, which is the honest answer
// for an SoC display engine.
func classify(gpu *domain.GPUInfo) {
	if gpu.PCIAddress == "" || !strings.Contains(gpu.PCIAddress, ":") {
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
	nodes := make(map[string]string, 2)
	for _, entry := range entries {
		if !strings.HasPrefix(entry, renderPrefix) {
			continue
		}
		address := sysfs.LinkBase(filepath.Join(classDir, entry, linkDevice))
		if address == "" {
			continue
		}
		nodes[address] = filepath.Join(deviceDir, entry)
	}
	return nodes
}

// resolveModels fills in Model for every GPU, preferring NVIDIA's own listing
// and otherwise streaming pci.ids exactly once.
func resolveModels(gpus []domain.GPUInfo) {
	pending := make(map[string]map[string][]int, 2)
	for i := range gpus {
		if model := nvidiaModelFor(&gpus[i]); model != "" {
			gpus[i].Model = model
			continue
		}
		vendor, device := trimHex(gpus[i].VendorID), trimHex(gpus[i].DeviceID)
		if vendor == "" || device == "" {
			continue
		}
		gpus[i].Model = gpus[i].DeviceID // replaced below if pci.ids has a name
		if pending[vendor] == nil {
			pending[vendor] = make(map[string][]int, 1)
		}
		pending[vendor][device] = append(pending[vendor][device], i)
	}
	if len(pending) == 0 {
		return
	}

	file := openPCIIDs()
	if file == nil {
		return
	}
	defer file.Close()

	remaining := 0
	for _, devices := range pending {
		remaining += len(devices)
	}

	scanner := bufio.NewScanner(file)
	var devices map[string][]int

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) < idFieldWidth || line[0] == commentByte {
			continue
		}
		if line[0] != indentByte {
			devices = pending[string(line[:idFieldWidth])]
			continue
		}
		if devices == nil || line[1] == indentByte {
			continue // no interest in this vendor, or a subsystem line
		}
		entry := line[1:]
		if len(entry) < idFieldWidth {
			continue
		}
		indices, ok := devices[string(entry[:idFieldWidth])]
		if !ok {
			continue
		}
		name := strings.TrimSpace(string(entry[idFieldWidth:]))
		for _, i := range indices {
			gpus[i].Model = name
		}
		delete(devices, string(entry[:idFieldWidth]))
		if remaining -= len(indices); remaining <= 0 {
			return
		}
	}
}

// nvidiaModelFor reads the marketing name the proprietary driver publishes.
func nvidiaModelFor(gpu *domain.GPUInfo) string {
	if gpu.Driver != driverNvidia && gpu.VendorID != vendorIDNvidia {
		return ""
	}
	if gpu.PCIAddress == "" {
		return ""
	}
	data, err := sysfs.Bytes(filepath.Join(nvidiaGPUDir, gpu.PCIAddress, nvidiaInfo))
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		if rest, found := strings.CutPrefix(line, nvidiaModel); found {
			return strings.TrimSpace(rest)
		}
	}
	return ""
}

func openPCIIDs() *os.File {
	for _, path := range pciIDsPaths {
		if f, err := os.Open(path); err == nil {
			return f
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
func trimHex(s string) string {
	return strings.TrimPrefix(s, "0x")
}
