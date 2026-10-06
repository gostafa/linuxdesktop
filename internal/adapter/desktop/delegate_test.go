// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package desktop

import (
	"context"
	"errors"
	"testing"
)

func TestFuncError(t *testing.T) {
	t.Parallel()
	failure := errors.New("probe failure")
	got, err := Func[int, int](
		func(context.Context, *int) (int, error) { return 7, failure },
	).Desktop(t.Context(), new(int))
	if got != 7 || !errors.Is(err, failure) {
		t.Fatal("partial result or error lost", got, err)
	}
}
