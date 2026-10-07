// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package schema

import (
	"time"
)

const (
	defaultAttempts = 3
	defaultInitial  = 25 * time.Millisecond
	defaultMaximum  = 100 * time.Millisecond
	zero            = 0
)

// DefaultRetryPolicy fits a short desktop detection budget.
func DefaultRetryPolicy() RetryPolicy {
	return RetryPolicy{
		MaxAttempts:     defaultAttempts,
		InitialInterval: defaultInitial,
		MaxInterval:     defaultMaximum,
	}
}

// ValidRetryPolicy reports whether the policy can configure exponential backoff.
func ValidRetryPolicy(policy RetryPolicy) bool {
	return policy.MaxAttempts > zero && policy.InitialInterval > zero &&
		policy.MaxInterval >= policy.InitialInterval
}
