// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package core

import (
	"testing"

	"github.com/gostafa/linuxdesktop/internal/domain"
)

func BenchmarkDetect(b *testing.B) {
	for _, bench := range []struct {
		name string
		cfg  Config
	}{
		{"all", Config{Sections: domain.SectionAll, ProcessScan: true}},
		{"session", Config{Sections: domain.SectionSession}},
	} {
		b.Run(bench.name, func(b *testing.B) {
			f := &fixture{calls: make(map[string]int)}
			deps := fixtureDeps(f)
			engine := New(&deps, &bench.cfg)
			b.ReportAllocs()
			for b.Loop() {
				if _, err := engine.Detect(b.Context()); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
