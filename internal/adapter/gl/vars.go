package gl

// vulkanICDDirs are where the Vulkan loader looks for ICD manifests. Each
// manifest names a driver and, usefully, the API version it implements.
var vulkanICDDirs = []string{
	"/usr/share/vulkan/icd.d",
	"/usr/local/share/vulkan/icd.d",
	"/etc/vulkan/icd.d",
}

// eglVendorDirs are where libglvnd looks for EGL vendor manifests. Their
// presence is what makes OpenGL usable through glvnd.
var eglVendorDirs = []string{
	"/usr/share/glvnd/egl_vendor.d",
	"/usr/local/share/glvnd/egl_vendor.d",
	"/etc/glvnd/egl_vendor.d",
}

// driDirs hold Mesa's DRI drivers, the fallback evidence that OpenGL is
// present on a system without glvnd manifests.
var driDirs = []string{
	"/usr/lib/dri",
	"/usr/lib64/dri",
	"/usr/lib/x86_64-linux-gnu/dri",
	"/usr/lib/aarch64-linux-gnu/dri",
	"/usr/lib/i386-linux-gnu/dri",
}

// driverVendors maps a substring of a driver library name onto a vendor.
var driverVendors = []struct {
	needle string
	vendor string
}{
	{"nvidia", vendorNVIDIA},
	{"mesa", vendorMesa},
	{"radeon", vendorAMD},
	{"amdgpu", vendorAMD},
	{"amd", vendorAMD},
	{"intel", vendorIntel},
	{"iris", vendorIntel},
}
