//go:build linux && (amd64 || arm64 || 386 || arm || riscv64 || loong64 || ppc64le)

package native

// Shared library sonames.
const (
	// Native calls return a null address or EGL_FALSE to indicate failure.
	nullAddress    uintptr = 0
	eglFalse       uint32  = 0
	eglNoConfigs   int32   = 0
	eglConfigCount int32   = 1
	eglESVersion   int32   = 2
	decimalBase            = 10
	vkNoVariant    uint32  = 0
	libEGL                 = "libEGL.so.1"
	libVulkan              = "libvulkan.so.1"

	// EGL enumerants.
	eglDefaultDisplay uintptr = 0
	eglNoDisplay      uintptr = 0
	eglNoSurface      uintptr = 0
	eglNoContext      uintptr = 0

	eglVendor  int32 = 0x3053
	eglVersion int32 = 0x3054

	eglOpenGLAPI   uint32 = 0x30A2
	eglOpenGLESAPI uint32 = 0x30A0

	eglSurfaceType    int32 = 0x3033
	eglRenderableType int32 = 0x3040
	eglPbufferBit     int32 = 0x0001
	eglOpenGLBit      int32 = 0x0008
	eglOpenGLES2Bit   int32 = 0x0004
	eglContextVersion int32 = 0x3098
	eglNone           int32 = 0x3038

	// OpenGL enumerants.
	glVendor   uint32 = 0x1F00
	glRenderer uint32 = 0x1F01
	glVersion  uint32 = 0x1F02

	// Vulkan. vkEnumerateInstanceVersion is absent from a 1.0 loader, where the
	// answer is 1.0.0 by definition.
	vkSuccess              int32 = 0
	vkEnumerateVersionName       = "vkEnumerateInstanceVersion"
	vkBaseVersion                = "1.0.0"

	// These decode a packed VK_MAKE_API_VERSION value.
	vkVariantShift = 29
	vkMajorShift   = 22
	vkMajorMask    = 0x7F
	vkMinorShift   = 12
	vkMinorMask    = 0x3FF
	vkPatchMask    = 0xFFF
)

// getStringSymbol is resolved through eglGetProcAddress, which works for both
// the desktop and the ES libraries without knowing which one is loaded.
const getStringSymbol = "glGetString"
