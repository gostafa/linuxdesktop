// Package gl reports the OpenGL and Vulkan implementations.
//
// It has two modes, because the honest answer costs real time and most callers
// do not need it.
//
// The default mode reads the filesystem only. Vulkan ICD manifests under
// /usr/share/vulkan/icd.d carry a genuine api_version, and the glvnd EGL
// vendor manifests plus the Mesa DRI driver directory establish that OpenGL is
// installed. This costs a directory listing and a few small JSON files, and it
// is the only mode that is safe to run unconditionally: it never touches the
// GPU. Its limitation is stated plainly in the results — Available means
// installed, and Renderer is left empty because no renderer string exists
// until a context has been created.
//
// The native mode, enabled by the caller, dlopens libEGL and libvulkan through
// purego — no cgo, so cross-compilation still works — and asks the drivers
// themselves. It initialises EGL for the vendor and version strings, then
// creates a surfaceless context so glGetString can name the real renderer,
// preferring desktop OpenGL and falling back to OpenGL ES. This is accurate
// and it is slow: eglInitialize can spin up the GPU driver stack and take tens
// of milliseconds.
//
// Because eglMakeCurrent binds to the calling thread rather than the
// goroutine, the native probe locks its OS thread for the whole exchange and
// unbinds before releasing it. Every step is guarded and every failure falls
// back to the filesystem answer, so a broken driver degrades the result rather
// than failing the detection run.
package gl
