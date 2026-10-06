// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package osinfo

import (
	"context"
	"errors"
	"testing"
)

func TestFuncError(t *testing.T) {
	t.Parallel()
	failure := errors.New("probe failure")
	got, err := Func[int](func(context.Context) (int, error) { return 7, failure }).OS(t.Context())
	if got != 7 || !errors.Is(err, failure) {
		t.Fatal("partial result or error lost", got, err)
	}
}
