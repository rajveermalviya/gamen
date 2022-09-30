//go:build js

package web

import (
	"syscall/js"
	"time"

	"github.com/rajveermalviya/gamen/internal/common/atomicx"
)

type Display struct {
	destroyRequested atomicx.Bool
	destroyed        atomicx.Bool

	windows map[uint64]*Window

	eventCallbacksChan chan func()
}

func NewDisplay() (*Display, error) {
	return &Display{
		windows:            map[uint64]*Window{},
		eventCallbacksChan: make(chan func()),
	}, nil
}

func (d *Display) Poll() bool {
	wait := make(chan struct{})
	cb := js.FuncOf(func(this js.Value, args []js.Value) any {
		if d.destroyed.Load() {
			return nil
		}

		wait <- struct{}{}

		return nil
	})
	defer cb.Release()
	js.Global().Call("requestAnimationFrame", cb)

outerloop:
	for {
		select {
		case cb := <-d.eventCallbacksChan:
			cb()

		innerloop:
			for {
				select {
				case cb := <-d.eventCallbacksChan:
					cb()
				default:
					break innerloop
				}
			}

		case <-wait:
			break outerloop
		}
	}

	if d.destroyRequested.Load() && !d.destroyed.Load() {
		d.destroy()
		return false
	}

	return !d.destroyed.Load()
}

func (d *Display) Wait() bool {
	// wait for first event
	cb := <-d.eventCallbacksChan
	cb()

	// then poll all pending events
	return d.Poll()
}

func (d *Display) WaitTimeout(timeout time.Duration) bool {
	timer := time.NewTimer(timeout)

	select {
	case <-timer.C:
		return !d.destroyed.Load()

	case cb := <-d.eventCallbacksChan:
		if !timer.Stop() {
			<-timer.C
		}

		cb()

		// then poll all pending events
		return d.Poll()
	}
}

func (d *Display) Destroy() {
	d.destroyRequested.Store(true)
}

func (d *Display) destroy() {
	for id, w := range d.windows {
		w.Destroy()

		d.windows[id] = nil
		delete(d.windows, id)
	}

	close(d.eventCallbacksChan)

	d.destroyed.Store(true)
}
