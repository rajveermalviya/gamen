//go:build darwin

package main

import (
	"github.com/rajveermalviya/gamen/display"
	"github.com/rajveermalviya/go-webgpu/wgpu"
	wgpuglfw "github.com/rajveermalviya/go-webgpu/wgpuext/glfw"
)

func getSurfaceDescriptor(w display.Window) *wgpu.SurfaceDescriptor {
	switch w := w.(type) {
	case display.GlfwWindow:
		return wgpuglfw.GetSurfaceDescriptor(w.GlfwWindow())
	default:
		panic("unsupported window")
	}
}
