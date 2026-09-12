// Package domain holds the model the whole library is written against.
//
// It is the centre of the hexagon: pure data, no I/O, no third-party imports.
// The root linuxdesktop package re-exports every type here as an alias and
// every constant by value, so the public API is identical to these
// declarations while adapters stay free to import the model without creating
// an import cycle back through the root package.
package domain
