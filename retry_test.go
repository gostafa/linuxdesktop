// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package linuxdesktop

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestConnectionRetryConfiguration(t *testing.T) {
	t.Parallel()
	defaults := publicConfig(nil).ConnectionRetry
	if defaults != (RetryPolicy{MaxAttempts: 3, InitialInterval: 25 * time.Millisecond, MaxInterval: 100 * time.Millisecond}) {
		t.Fatal(defaults)
	}
	custom := RetryPolicy{
		MaxAttempts:     1,
		InitialInterval: time.Millisecond,
		MaxInterval:     time.Millisecond,
	}
	if got := publicConfig(
		[]Option{WithConnectionRetry(defaults), WithConnectionRetry(custom)},
	); got.ConnectionRetry != custom {
		t.Fatal(got.ConnectionRetry)
	}
	if _, err := detectOn(
		t.Context(),
		newConfig(WithSections(SectionOS), WithConnectionRetry(custom)),
		goosLinux,
	); err != nil {
		t.Fatal(err)
	}
	for _, timeout := range []time.Duration{0, -time.Second} {
		if _, err := detectOn(
			t.Context(),
			newConfig(WithSections(SectionOS), WithTimeout(timeout)),
			goosLinux,
		); err != nil {
			t.Fatal(err)
		}
	}
}

func TestInvalidConnectionRetryPolicy(t *testing.T) {
	t.Parallel()
	for _, policy := range []RetryPolicy{
		{},
		{InitialInterval: time.Millisecond, MaxInterval: time.Millisecond},
		{MaxAttempts: 1, MaxInterval: time.Millisecond},
		{MaxAttempts: 1, InitialInterval: -time.Millisecond, MaxInterval: time.Millisecond},
		{MaxAttempts: 1, InitialInterval: time.Second, MaxInterval: time.Millisecond},
	} {
		cfg := newConfig(WithConnectionRetry(policy))
		result, err := detectOn(t.Context(), cfg, goosLinux)
		if !errors.Is(err, ErrInvalidRetryPolicy) || result == nil ||
			!reflect.DeepEqual(result, new(Environment)) {
			t.Fatal(policy, result, err)
		}
		result, err = detectOn(t.Context(), cfg, "other")
		if result == nil || !errors.Is(err, ErrNotLinux) || errors.Is(err, ErrInvalidRetryPolicy) {
			t.Fatal(result, err)
		}
	}
}
