package glfw

import (
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
)

type Display struct {
}

func NewDisplay() (*Display, error) {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.ClientAPI, glfw.NoAPI)
	return &Display{}, nil
}

func (d *Display) NewWindow() (*Window, error) {
	w, err := glfw.CreateWindow(640, 480, "go-webgpu with glfw", nil, nil)
	if err != nil {
		return nil, err
	}
	return &Window{w}, nil
}

func (d *Display) Destroy() {
	println("glfw.Terminate()")
	glfw.Terminate()
}

func (d *Display) Wait() bool {
	glfw.WaitEvents()
	return true
}

func (d *Display) Poll() bool {
	glfw.PollEvents()
	return true
}

func (d *Display) WaitTimeout(timeout time.Duration) bool {
	var seconds float64 = float64(timeout) / float64(time.Second)
	glfw.WaitEventsTimeout(seconds)
	return false
}
