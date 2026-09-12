// Package env snapshots every environment variable the library consults.
//
// It runs once, synchronously, before any other probe. Taking the snapshot up
// front means no probe pays for a repeated os.Getenv, and more importantly it
// means the whole detection run reasons over one consistent view of the
// environment even though the probes run concurrently.
package env
