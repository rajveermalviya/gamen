//go:build js

package web

import (
	"sync"
	"syscall/js"
	"time"
)

type Display struct {
	destroyed   bool
	destroyOnce sync.Once

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
	js.Global().Call("requestAnimationFrame", js.FuncOf(func(this js.Value, args []js.Value) any {
		if d.destroyed {
			return nil
		}

		wait <- struct{}{}

		return nil
	}))

	for {
		select {
		case cb := <-d.eventCallbacksChan:
			cb()

		loop:
			for {
				select {
				case cb := <-d.eventCallbacksChan:
					cb()
				default:
					break loop
				}
			}

		case <-wait:
			return !d.destroyed
		}
	}
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
		return !d.destroyed

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
	d.destroyOnce.Do(func() {
		for id, w := range d.windows {
			w.Destroy()

			d.windows[id] = nil
			delete(d.windows, id)
		}

		// drain the channel
	loop:
		for {
			select {
			case cb := <-d.eventCallbacksChan:
				cb()
			default:
				break loop
			}
		}

		close(d.eventCallbacksChan)
	})
}
