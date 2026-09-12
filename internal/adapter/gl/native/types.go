//go:build linux && (amd64 || arm64 || 386 || arm || riscv64 || loong64 || ppc64le)

package native

// eglFuncs holds the entry points the OpenGL probe binds.
type eglFuncs struct {
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
type clientAPI struct {
	enum       uint32
	renderable int32
	ctxAttrs   []int32
}
