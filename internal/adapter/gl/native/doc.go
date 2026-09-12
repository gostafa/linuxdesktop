// Package native asks the graphics drivers themselves what they are.
//
// It dlopens libEGL and libvulkan through purego, so there is no cgo and
// cross-compilation still works. EGL is initialised for the vendor and version
// strings, then a surfaceless context is created so glGetString can name the
// real renderer, preferring desktop OpenGL and falling back to OpenGL ES.
//
// This package exists separately from its parent for one reason: purego does
// not build everywhere Go does. It has no dlopen on Windows and incomplete
// assembly support on some architectures, and this library must still compile
// for a cross-platform application that imports it unconditionally. Confining
// the dependency here means the build constraint is confined here too — the
// parent package, and every other package in the library, stays free of build
// tags.
//
// So the implementation in funcs.go is constrained to the platforms purego
// supports, and vars.go carries the fallback for everywhere else, reporting
// nothing rather than failing to compile. The list is a positive allowlist
// rather than a list of exclusions, so a new Go port is unsupported until it
// has been tried, instead of breaking the build.
//
// Both entry points report success as a second return value. A false there
// means "ask the filesystem instead", never "there is no OpenGL".
package native
