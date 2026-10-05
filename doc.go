// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

// Package linuxdesktop detects the Linux desktop environment, display servers,
// and graphics stack. Detection returns partial results alongside probe errors;
// callers can use the available data even when a probe fails. On other operating
// systems it returns a headless environment and ErrNotLinux.
//
// DetectContext applies options to the default configuration in order. Use the
// With functions or supply a custom option, for example:
//
//	func(c *linuxdesktop.Config) { c.Sections = linuxdesktop.SectionOS }
//
// Public result and configuration types belong to this package. Probe engines
// and their models are implementation details kept under internal/.
package linuxdesktop
