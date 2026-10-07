// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package sysfs

import (
	"testing"
)

func osReleaseFixture() []byte {
	return []byte(`# os-release fixture
PRETTY_NAME="Ubuntu 24.04.1 LTS"
NAME="Ubuntu"
VERSION_ID="24.04"
VERSION="24.04.1 LTS (Noble Numbat)"
VERSION_CODENAME=noble
ID=ubuntu
ID_LIKE=debian
HOME_URL="https://www.ubuntu.com/"
SUPPORT_URL="https://help.ubuntu.com/"
BUG_REPORT_URL="https://bugs.launchpad.net/ubuntu/"
PRIVACY_POLICY_URL="https://www.ubuntu.com/legal/terms-and-policies/privacy-policy"
UBUNTU_CODENAME=noble
LOGO=ubuntu-logo
`)
}

func BenchmarkEach(b *testing.B) {
	data := osReleaseFixture()
	b.ReportAllocs()
	b.SetBytes(int64(len(data)))
	for b.Loop() {
		count := 0
		Each(data, func(string, string) bool { count++; return true })
		if count == 0 {
			b.Fatal("no fields parsed")
		}
	}
}

func BenchmarkField(b *testing.B) {
	data := osReleaseFixture()
	for _, key := range []string{"PRETTY_NAME", "LOGO", "ABSENT"} {
		b.Run(key, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = Field(data, key)
			}
		})
	}
}
