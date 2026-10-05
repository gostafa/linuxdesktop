// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package drm

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gostafa/linuxdesktop/internal/domain"
	"github.com/gostafa/linuxdesktop/internal/sysfs"
)

func write(t *testing.T, files string, path, data string) {
	t.Helper()
	path = sysfs.Path(files, path)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatal(err)
	}
}

func link(t *testing.T, files string, path, target string) {
	t.Helper()
	path = sysfs.Path(files, path)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, path); err != nil {
		t.Fatal(err)
	}
}

func TestGPUEnumeration(t *testing.T) {
	pciIDsPaths := pciIDsPaths()

	t.Parallel()
	files := t.TempDir()
	if gpus, primary, err := gpus(
		files,
		t.Context(),
	); err != nil || len(gpus) != 0 ||
		primary != "" {
		t.Fatal(gpus, primary, err)
	}
	for _, entry := range []string{"card0", "card1"} {
		write(t, files, filepath.Join(classDir, entry, attrVendor), "0x8086")
		write(t, files, filepath.Join(classDir, entry, attrDevice), "0x1234")
	}
	write(
		t,
		files,
		filepath.Join(classDir, "card0", attrUevent),
		"DRIVER=i915\nPCI_SLOT_NAME=0000:00:02.0\n",
	)
	write(t, files, filepath.Join(classDir, "card0", attrBootVGA), "1")
	link(t, files, filepath.Join(classDir, "renderD128", linkDevice), "../../../0000:00:02.0")
	write(t, files, filepath.Join(classDir, "renderD129", "unused"), "")
	write(t, files, filepath.Join(classDir, "card0-HDMI-A-1", "unused"), "")
	write(t, files, pciIDsPaths[0], "# comment\n8086 Intel\n\t1234 Fixture GPU\n")
	gpus, primary, err := gpus(files, t.Context())
	if err != nil || len(gpus) != 2 || primary != "card0" || gpus[0].Model != "Fixture GPU" ||
		gpus[0].RenderDevice != "/dev/dri/renderD128" ||
		gpus[0].Driver != "i915" ||
		!gpus[0].Integrated {
		t.Fatalf("%+v %s %v", gpus, primary, err)
	}
	if orFirst("", gpus) != "card0" || orFirst("", nil) != "" {
		t.Fatal("primary fallback")
	}
	if _, _, err = New().GPUs(t.Context()); err != nil {
		t.Fatal(err)
	}
}

func TestGPUModelsAndClassification(t *testing.T) {
	pciIDsPaths := pciIDsPaths()

	t.Parallel()
	files := t.TempDir()
	gpu := domain.GPUInfo{VendorID: vendorIDNvidia, DeviceID: "0x1234", PCIAddress: "0000:01:00.0"}
	if nvidiaModelFor(files, &gpu) != "" {
		t.Fatal("unexpected model")
	}
	write(
		t,
		files,
		filepath.Join(nvidiaGPUDir, gpu.PCIAddress, nvidiaInfo),
		"IRQ: 1\nModel: Fixture NVIDIA\n",
	)
	if _, _, ok := modelLookup(files, &gpu); ok || gpu.Model != "Fixture NVIDIA" {
		t.Fatal(gpu)
	}
	for _, address := range []string{"", "platform", "0000:00:02.0", "0000:02:00.0"} {
		card := domain.GPUInfo{PCIAddress: address}
		classify(&card)
		if !card.Integrated && !card.Discrete {
			t.Fatal(card)
		}
	}
	classify(&gpu)
	if !gpu.Discrete {
		t.Fatal(gpu)
	}
	if modelLine([]byte("IRQ: 1")) != "" {
		t.Fatal("absent model")
	}
	if err := resolveModels(files, nil); err != nil {
		t.Fatal(err)
	}
	unresolved := []domain.GPUInfo{{}, {VendorID: "0x8086", DeviceID: "0x1234"}}
	if err := resolveModels(files, unresolved); err != nil || unresolved[1].Model != "0x1234" {
		t.Fatal(unresolved, err)
	}
	write(t, files, pciIDsPaths[1], "8086 Intel\n\t1234 Alternate GPU\n")
	if err := resolveModels(
		files,
		unresolved,
	); err != nil ||
		unresolved[1].Model != "Alternate GPU" {
		t.Fatal(unresolved, err)
	}
	write(t, files, filepath.Join(classDir, "card2", attrUevent), "DRIVER=fixture\n")
	var empty domain.GPUInfo
	applyUevent(files, &empty, filepath.Join(classDir, "card2"))
	if empty.Driver != "fixture" || empty.PCIAddress != "" {
		t.Fatal(empty)
	}
}

type failingReader struct{ err error }

func (reader failingReader) Read([]byte) (int, error) { return 0, reader.err }

func TestPCIScan(t *testing.T) {
	t.Parallel()
	gpus := make([]domain.GPUInfo, 2)
	newScan := func() *pciScan {
		return newPCIScan(map[string]map[string][]int{"8086": {"1234": {0}, "5678": {1}}}, gpus)
	}
	scan := newScan()
	listing := "\n# ignored\n9999 other\n\t0001 ignored\n8086 Intel\n\t\t1234 subsystem\n\t0000 unmatched\n\t1234 First\n\t5678 Second\n"
	if err := run(
		scan,
		strings.NewReader(listing),
	); err != nil || gpus[0].Model != "First" ||
		gpus[1].Model != "Second" {
		t.Fatal(gpus, err)
	}
	if err := run(newScan(), strings.NewReader("8086 Intel\n")); err != nil {
		t.Fatal(err)
	}
	failure := io.ErrUnexpectedEOF
	if err := run(newScan(), failingReader{failure}); !errors.Is(err, failure) {
		t.Fatal(err)
	}
}
