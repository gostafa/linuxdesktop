//go:build linux && (amd64 || arm64 || 386 || arm || riscv64 || loong64 || ppc64le)

package native

type (
	// library isolates dynamic linking from the graphics protocol.
	library interface {
		open(name string, flags int) (uintptr, bool)
		symbol(handle uintptr, name string) uintptr
		register(target any, address uintptr)
	}

	// systemLibrary binds native entry points through purego.
	systemLibrary struct{}

	// eglFuncs holds the entry points the OpenGL probe binds.
	eglFuncs struct {
		bindString     func(uintptr) func(uint32) string
		getDisplay     func(uintptr) uintptr
		initialize     func(uintptr, *int32, *int32) uint32
		queryString    func(uintptr, int32) string
		bindAPI        func(uint32) uint32
		chooseConfig   func(uintptr, *int32, *uintptr, int32, *int32) uint32
		createContext  func(uintptr, uintptr, uintptr, *int32) uintptr
		makeCurrent    func(uintptr, uintptr, uintptr, uintptr) uint32
		destroyContext func(uintptr, uintptr) uint32
		terminate      func(uintptr) uint32
		getProcAddress func(string) uintptr
	}

	// clientAPI is one of the two rendering APIs the probe will try, in order.
	clientAPI struct {
		ctxAttrs   []int32
		enum       uint32
		renderable int32
	}
)
