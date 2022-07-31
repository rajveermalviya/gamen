//go:build windows

package win32

import (
	"sync"
	"time"
	"unsafe"

	"github.com/rajveermalviya/gamen/internal/win32/procs"
)

type Display struct {
	windows     map[uintptr]*Window
	destroyed   bool
	destroyOnce sync.Once
}

func NewDisplay() (*Display, error) {
	return &Display{
		windows: map[uintptr]*Window{},
	}, nil
}

func (d *Display) Destroy() {
	d.destroyOnce.Do(func() {
		d.destroyed = true

		for hwnd, w := range d.windows {
			w.Destroy()

			d.windows[hwnd] = nil
			delete(d.windows, hwnd)
		}
	})
}

func (d *Display) Poll() bool {
	var msg procs.MSG

	for procs.PeekMessageW(uintptr(unsafe.Pointer(&msg)), 0, 0, 0, procs.PM_REMOVE) {
		procs.TranslateMessage(uintptr(unsafe.Pointer(&msg)))
		procs.DispatchMessageW(uintptr(unsafe.Pointer(&msg)))
	}

	return !d.destroyed
}

func (d *Display) Wait() bool {
	procs.WaitMessage()
	return d.Poll()
}

func (d *Display) WaitTimeout(timeout time.Duration) bool {
	procs.MsgWaitForMultipleObjects(0, 0, 0, uintptr(timeout.Milliseconds()), procs.QS_ALLEVENTS)
	return d.Poll()
}
