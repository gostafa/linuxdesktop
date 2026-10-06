// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package gl

import (
	"context"
	"errors"
	"testing"
)

func TestFuncError(t *testing.T) {
	t.Parallel()
	failure := errors.New("probe failure")
	got, _, err := Func[int, string](
		func(context.Context) (int, string, error) { return 7, "vulkan", failure },
	).Stack(t.Context())
	if got != 7 || !errors.Is(err, failure) {
		t.Fatal("partial result or error lost", got, err)
	}
}
