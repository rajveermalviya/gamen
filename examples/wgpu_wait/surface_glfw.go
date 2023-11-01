//go:build darwin

package main

import (
	"github.com/rajveermalviya/gamen/display"
	"github.com/rajveermalviya/gamen/internal/glfw"
	"github.com/rajveermalviya/go-webgpu/wgpu"
	wgpuglfw "github.com/rajveermalviya/go-webgpu/wgpuext/glfw"
)

func getSurfaceDescriptor(w display.Window) *wgpu.SurfaceDescriptor {
	switch w := w.(type) {
	case *glfw.Window:
		return wgpuglfw.GetSurfaceDescriptor(w.Window())
	default:
		panic("unsupported window")
	}
}
