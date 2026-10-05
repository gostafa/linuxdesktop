//go:build linux && (amd64 || arm64 || 386 || arm || riscv64 || loong64 || ppc64le)

package native

import (
	"runtime"
	"strconv"

	"github.com/ebitengine/purego"
	"github.com/gostafa/linuxdesktop/internal/domain"
)

// Vulkan asks the loader for its API version.
func Vulkan() (version string, ok bool) {
	// purego may panic on an unexpected native symbol signature.
	defer func() {
		if recover() != nil {
			version, ok = "", false
		}
	}()
	return vulkanVersion()
}

func vulkanVersion() (string, bool) {
	lib, loaded := openLibrary(libVulkan, purego.RTLD_NOW|purego.RTLD_LOCAL)
	if !loaded {
		return "", false
	}
	sym, err := purego.Dlsym(lib, vkEnumerateVersionName)
	if err != nil || sym == nullAddress {
		return vkBaseVersion, true
	}
	return enumerateVersion(sym)
}

func enumerateVersion(symbol uintptr) (string, bool) {
	var enumerate func(*uint32) int32
	purego.RegisterFunc(&enumerate, symbol)
	var packed uint32
	if enumerate(&packed) != vkSuccess {
		return "", false
	}
	return decodeVersion(packed), true
}

// OpenGL initializes EGL and reads renderer details through a surfaceless context.
func OpenGL() (info domain.OpenGLInfo, ok bool) {
	defer func() {
		if recover() != nil {
			ok = false
		}
	}()
	// EGL contexts belong to the current OS thread.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	ok = populateGL(&info)
	return info, ok
}

func populateGL(info *domain.OpenGLInfo) bool {
	egl, loaded := loadEGL()
	if !loaded {
		return false
	}
	display, initialized := initializeDisplay(&egl)
	if !initialized {
		return false
	}
	defer egl.terminate(display)
	*info = displayInfo(&egl, display)
	mergeCurrent(info, &egl, display)
	return true
}

func initializeDisplay(egl *eglFuncs) (uintptr, bool) {
	display := egl.getDisplay(eglDefaultDisplay)
	if display == eglNoDisplay {
		return display, false
	}
	var major, minor int32
	return display, egl.initialize(display, &major, &minor) != eglFalse
}

func displayInfo(egl *eglFuncs, display uintptr) domain.OpenGLInfo {
	return domain.OpenGLInfo{
		Available: true,
		Vendor:    egl.queryString(display, eglVendor),
		Version:   egl.queryString(display, eglVersion),
		Renderer:  "",
	}
}

func mergeCurrent(info *domain.OpenGLInfo, egl *eglFuncs, display uintptr) {
	live, got := currentGL(egl, display)
	if !got {
		return
	}
	info.Renderer = live.Renderer
	info.Vendor = preferLive(info.Vendor, live.Vendor)
	info.Version = preferLive(info.Version, live.Version)
}

func preferLive(known, live string) string {
	if live != "" {
		return live
	}
	return known
}

// currentGL tries desktop OpenGL before OpenGL ES.
func currentGL(egl *eglFuncs, display uintptr) (domain.OpenGLInfo, bool) {
	apis := []clientAPI{
		{eglOpenGLAPI, eglOpenGLBit, []int32{eglNone}},
		{eglOpenGLESAPI, eglOpenGLES2Bit, []int32{eglContextVersion, eglESVersion, eglNone}},
	}
	for i := range apis {
		if info, ok := tryAPI(egl, display, &apis[i]); ok {
			return info, true
		}
	}
	var empty domain.OpenGLInfo
	return empty, false
}

func tryAPI(egl *eglFuncs, display uintptr, api *clientAPI) (domain.OpenGLInfo, bool) {
	var empty domain.OpenGLInfo
	ctx := createGLContext(egl, display, api)
	if ctx == eglNoContext {
		return empty, false
	}
	defer egl.destroyContext(display, ctx)
	return contextInfo(egl, display, ctx)
}

func createGLContext(egl *eglFuncs, display uintptr, api *clientAPI) uintptr {
	if egl.bindAPI(api.enum) == eglFalse {
		return eglNoContext
	}
	config, chosen := chooseConfig(egl, display, api.renderable)
	if !chosen {
		return eglNoContext
	}
	return egl.createContext(display, config, eglNoContext, &api.ctxAttrs[0])
}

func chooseConfig(egl *eglFuncs, display uintptr, renderable int32) (uintptr, bool) {
	attrs := []int32{eglSurfaceType, eglPbufferBit, eglRenderableType, renderable, eglNone}
	var config uintptr
	var count int32
	chosen := egl.chooseConfig(display, &attrs[0], &config, eglConfigCount, &count)
	return config, chosen != eglFalse && count != eglNoConfigs
}

func contextInfo(egl *eglFuncs, display, ctx uintptr) (domain.OpenGLInfo, bool) {
	var empty domain.OpenGLInfo
	if egl.makeCurrent(display, eglNoSurface, eglNoSurface, ctx) == eglFalse {
		return empty, false
	}
	defer egl.makeCurrent(display, eglNoSurface, eglNoSurface, eglNoContext)
	return readGLStrings(egl)
}

func readGLStrings(egl *eglFuncs) (domain.OpenGLInfo, bool) {
	var empty domain.OpenGLInfo
	getString, bound := bindGetString(egl)
	if !bound {
		return empty, false
	}
	info := domain.OpenGLInfo{
		Available: true,
		Vendor:    getString(glVendor),
		Renderer:  getString(glRenderer),
		Version:   getString(glVersion),
	}
	return info, info.Renderer != ""
}

func bindGetString(egl *eglFuncs) (getString func(uint32) string, ok bool) {
	defer func() {
		if recover() != nil {
			getString, ok = nil, false
		}
	}()
	addr := egl.getProcAddress(getStringSymbol)
	if addr == nullAddress {
		return nil, false
	}
	purego.RegisterFunc(&getString, addr)
	return getString, true
}

func loadEGL() (egl eglFuncs, ok bool) {
	defer func() {
		if recover() != nil {
			ok = false
		}
	}()
	return eglBindings()
}

func eglBindings() (egl eglFuncs, ok bool) {
	lib, loaded := openLibrary(libEGL, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if !loaded {
		return egl, false
	}
	bindDisplayFunctions(&egl, lib)
	bindContextFunctions(&egl, lib)
	return egl, true
}

func openLibrary(name string, flags int) (uintptr, bool) {
	lib, err := purego.Dlopen(name, flags)
	return lib, err == nil && lib != nullAddress
}

func bindDisplayFunctions(egl *eglFuncs, lib uintptr) {
	purego.RegisterLibFunc(&egl.getDisplay, lib, "eglGetDisplay")
	purego.RegisterLibFunc(&egl.initialize, lib, "eglInitialize")
	purego.RegisterLibFunc(&egl.queryString, lib, "eglQueryString")
	purego.RegisterLibFunc(&egl.bindAPI, lib, "eglBindAPI")
	purego.RegisterLibFunc(&egl.chooseConfig, lib, "eglChooseConfig")
}

func bindContextFunctions(egl *eglFuncs, lib uintptr) {
	purego.RegisterLibFunc(&egl.createContext, lib, "eglCreateContext")
	purego.RegisterLibFunc(&egl.makeCurrent, lib, "eglMakeCurrent")
	purego.RegisterLibFunc(&egl.destroyContext, lib, "eglDestroyContext")
	purego.RegisterLibFunc(&egl.terminate, lib, "eglTerminate")
	purego.RegisterLibFunc(&egl.getProcAddress, lib, "eglGetProcAddress")
}

// decodeVersion unpacks a VK_MAKE_API_VERSION value.
func decodeVersion(packed uint32) string {
	major := (packed >> vkMajorShift) & vkMajorMask
	minor := (packed >> vkMinorShift) & vkMinorMask
	patch := packed & vkPatchMask
	version := strconv.FormatUint(uint64(major), decimalBase) + "." +
		strconv.FormatUint(uint64(minor), decimalBase) + "." +
		strconv.FormatUint(uint64(patch), decimalBase)
	if variant := packed >> vkVariantShift; variant != vkNoVariant {
		return version + " (variant " + strconv.FormatUint(uint64(variant), decimalBase) + ")"
	}
	return version
}
