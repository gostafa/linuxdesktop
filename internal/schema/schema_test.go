// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package schema

import (
	"testing"
	"time"
)

func TestRetryPolicy(t *testing.T) {
	t.Parallel()
	if !ValidRetryPolicy(DefaultRetryPolicy()) {
		t.Fatal("invalid defaults")
	}
	for _, policy := range []RetryPolicy{{}, {MaxAttempts: 1}, {MaxAttempts: 1, InitialInterval: time.Second}} {
		if ValidRetryPolicy(policy) {
			t.Fatal(policy)
		}
	}
}
