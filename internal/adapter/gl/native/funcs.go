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
	// purego panics rather than erroring on an unexpected symbol signature.
	defer func() {
		if recover() != nil {
			version, ok = "", false
		}
	}()

	lib, err := purego.Dlopen(libVulkan, purego.RTLD_NOW|purego.RTLD_LOCAL)
	if err != nil || lib == 0 {
		return "", false
	}
	sym, err := purego.Dlsym(lib, vkEnumerateVersionName)
	if err != nil || sym == 0 {
		return vkBaseVersion, true
	}

	var enumerate func(*uint32) int32
	purego.RegisterFunc(&enumerate, sym)

	var packed uint32
	if enumerate(&packed) != vkSuccess {
		return "", false
	}
	return decodeVersion(packed), true
}

// OpenGL initialises EGL and, if it can, a surfaceless context, so that
// glGetString can report the renderer the caller actually cares about.
func OpenGL() (info domain.OpenGLInfo, ok bool) {
	defer func() {
		if recover() != nil {
			ok = false
		}
	}()

	// eglMakeCurrent binds to a thread, not to a goroutine, so the goroutine
	// must not be allowed to migrate while a context is current.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	egl, loaded := loadEGL()
	if !loaded {
		return info, false
	}

	display := egl.getDisplay(eglDefaultDisplay)
	if display == eglNoDisplay {
		return info, false
	}
	var major, minor int32
	if egl.initialize(display, &major, &minor) == 0 {
		return info, false
	}
	defer egl.terminate(display)

	info.Available = true
	info.Vendor = egl.queryString(display, eglVendor)
	info.Version = egl.queryString(display, eglVersion)

	if live, got := currentGL(egl, display); got {
		info.Renderer = live.Renderer
		if live.Vendor != "" {
			info.Vendor = live.Vendor
		}
		if live.Version != "" {
			info.Version = live.Version
		}
	}
	return info, true
}

// currentGL creates a surfaceless context and reads the GL strings, trying
// desktop OpenGL first and OpenGL ES second.
func currentGL(egl eglFuncs, display uintptr) (domain.OpenGLInfo, bool) {
	apis := []clientAPI{
		{eglOpenGLAPI, eglOpenGLBit, []int32{eglNone}},
		{eglOpenGLESAPI, eglOpenGLES2Bit, []int32{eglContextVersion, 2, eglNone}},
	}

	for _, api := range apis {
		if egl.bindAPI(api.enum) == 0 {
			continue
		}
		attrs := []int32{
			eglSurfaceType, eglPbufferBit,
			eglRenderableType, api.renderable,
			eglNone,
		}
		var config uintptr
		var count int32
		if egl.chooseConfig(display, &attrs[0], &config, 1, &count) == 0 || count == 0 {
			continue
		}
		ctx := egl.createContext(display, config, eglNoContext, &api.ctxAttrs[0])
		if ctx == eglNoContext {
			continue
		}

		var info domain.OpenGLInfo
		var ok bool
		if egl.makeCurrent(display, eglNoSurface, eglNoSurface, ctx) != 0 {
			if getString, bound := bindGetString(egl); bound {
				info = domain.OpenGLInfo{
					Available: true,
					Vendor:    getString(glVendor),
					Renderer:  getString(glRenderer),
					Version:   getString(glVersion),
				}
				ok = info.Renderer != ""
			}
			egl.makeCurrent(display, eglNoSurface, eglNoSurface, eglNoContext)
		}
		egl.destroyContext(display, ctx)

		if ok {
			return info, true
		}
	}
	return domain.OpenGLInfo{}, false
}

func bindGetString(egl eglFuncs) (fn func(uint32) string, ok bool) {
	defer func() {
		if recover() != nil {
			fn, ok = nil, false
		}
	}()
	addr := egl.getProcAddress(getStringSymbol)
	if addr == 0 {
		return nil, false
	}
	purego.RegisterFunc(&fn, addr)
	return fn, true
}

func loadEGL() (egl eglFuncs, ok bool) {
	defer func() {
		if recover() != nil {
			ok = false
		}
	}()

	lib, err := purego.Dlopen(libEGL, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil || lib == 0 {
		return egl, false
	}
	purego.RegisterLibFunc(&egl.getDisplay, lib, "eglGetDisplay")
	purego.RegisterLibFunc(&egl.initialize, lib, "eglInitialize")
	purego.RegisterLibFunc(&egl.queryString, lib, "eglQueryString")
	purego.RegisterLibFunc(&egl.bindAPI, lib, "eglBindAPI")
	purego.RegisterLibFunc(&egl.chooseConfig, lib, "eglChooseConfig")
	purego.RegisterLibFunc(&egl.createContext, lib, "eglCreateContext")
	purego.RegisterLibFunc(&egl.makeCurrent, lib, "eglMakeCurrent")
	purego.RegisterLibFunc(&egl.destroyContext, lib, "eglDestroyContext")
	purego.RegisterLibFunc(&egl.terminate, lib, "eglTerminate")
	purego.RegisterLibFunc(&egl.getProcAddress, lib, "eglGetProcAddress")
	return egl, true
}

// decodeVersion unpacks a VK_MAKE_API_VERSION value.
func decodeVersion(packed uint32) string {
	major := (packed >> vkMajorShift) & vkMajorMask
	minor := (packed >> vkMinorShift) & vkMinorMask
	patch := packed & vkPatchMask

	version := strconv.FormatUint(uint64(major), 10) + "." +
		strconv.FormatUint(uint64(minor), 10) + "." +
		strconv.FormatUint(uint64(patch), 10)

	if variant := packed >> vkVariantShift; variant != 0 {
		return version + " (variant " + strconv.FormatUint(uint64(variant), 10) + ")"
	}
	return version
}
