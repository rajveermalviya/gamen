//go:build android

package android

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/rajveermalviya/gamen/dpi"
	"github.com/rajveermalviya/gamen/events"
)

/*

#cgo CXXFLAGS: -isystem ${SRCDIR}/game-activity/include/
#cgo CFLAGS: -isystem ${SRCDIR}/game-activity/include/
#cgo LDFLAGS: -static-libstdc++ -landroid -llog

#include <game-activity/native_app_glue/android_native_app_glue.h>

extern void display_poll(int timeoutMillis);
extern void display_set_handler(struct android_app* app);

*/
import "C"

var androidApp unsafe.Pointer

//go:linkname main_main main.main
func main_main()

//export android_main
func android_main(app *C.struct_android_app) {
	atomic.StorePointer(&androidApp, unsafe.Pointer(app))

	C.display_set_handler(app)

	main_main()
}

type Display struct{}

func NewDisplay() (*Display, error) {
	return &Display{}, nil
}

func (d *Display) Destroy() {
	app := (*C.struct_android_app)(atomic.LoadPointer(&androidApp))
	if app != nil {
		C.GameActivity_finish(app.activity)
	}
}
func (d *Display) Wait() bool {
	C.display_poll(-1)

	app := (*C.struct_android_app)(atomic.LoadPointer(&androidApp))
	if app != nil && app.destroyRequested != 0 {
		return false
	}

	handleInput()
	return true
}
func (d *Display) Poll() bool {
	C.display_poll(0)

	app := (*C.struct_android_app)(atomic.LoadPointer(&androidApp))
	if app != nil && app.destroyRequested != 0 {
		return false
	}

	handleInput()
	return true
}
func (d *Display) WaitTimeout(t time.Duration) bool {
	C.display_poll(C.int(t.Milliseconds()))

	app := (*C.struct_android_app)(atomic.LoadPointer(&androidApp))
	if app != nil && app.destroyRequested != 0 {
		return false
	}

	handleInput()
	return true
}

var (
	cbMut                           sync.Mutex
	windowSurfaceCreatedCb          events.WindowSurfaceCreatedCallback
	windowSurfaceDestroyedCb        events.WindowSurfaceDestroyedCallback
	windowResizedCallback           events.WindowResizedCallback
	windowFocusedCb                 events.WindowFocusedCallback
	windowUnfocusedCb               events.WindowUnfocusedCallback
	windowTouchInputCb              events.WindowTouchInputCallback
	windowKeyboardInputCb           events.WindowKeyboardInputCallback
	windowReceivedCharacterCallback events.WindowReceivedCharacterCallback
)

func runResizedCallback() {
	cbMut.Lock()
	resizedCb := windowResizedCallback
	cbMut.Unlock()

	app := (*C.struct_android_app)(atomic.LoadPointer(&androidApp))
	var window *C.ANativeWindow
	var config *C.AConfiguration
	if app != nil {
		window = app.window
		config = app.config
	}

	if resizedCb != nil && window != nil {
		newWidth := C.ANativeWindow_getWidth(window)
		newHeight := C.ANativeWindow_getHeight(window)
		scaleFactor := float64(1)
		if config != nil {
			density := C.AConfiguration_getDensity(config)
			scaleFactor = float64(density) / float64(160)
		}

		resizedCb(uint32(newWidth), uint32(newHeight), scaleFactor)
	}
}

//export display_handle_command
func display_handle_command(app *C.struct_android_app, cmd C.int32_t) {
	switch cmd {
	case C.APP_CMD_INIT_WINDOW:
		println("init_window")

		cbMut.Lock()
		surfaceCreatedCb := windowSurfaceCreatedCb
		cbMut.Unlock()

		if surfaceCreatedCb != nil {
			surfaceCreatedCb()
		}

	case C.APP_CMD_TERM_WINDOW:
		println("term_window")

		cbMut.Lock()
		surfaceDestroyedCb := windowSurfaceDestroyedCb
		cbMut.Unlock()

		if surfaceDestroyedCb != nil {
			surfaceDestroyedCb()
		}

	case C.APP_CMD_WINDOW_RESIZED:
		println("window_resized")
		runResizedCallback()

	case C.APP_CMD_CONFIG_CHANGED:
		println("config_changed")
		runResizedCallback()

	case C.APP_CMD_WINDOW_INSETS_CHANGED:
		println("window_insets_changed")
		runResizedCallback()

	case C.APP_CMD_CONTENT_RECT_CHANGED:
		println("content_rect_changed")
		runResizedCallback()

	case C.APP_CMD_LOW_MEMORY:
		println("low_memory")
		runtime.GC()

	case C.APP_CMD_GAINED_FOCUS:
		println("gained_focus")

		cbMut.Lock()
		focusedCb := windowFocusedCb
		cbMut.Unlock()

		if focusedCb != nil {
			focusedCb()
		}

	case C.APP_CMD_LOST_FOCUS:
		println("lost_focus")

		cbMut.Lock()
		unfocusedCb := windowUnfocusedCb
		cbMut.Unlock()

		if unfocusedCb != nil {
			unfocusedCb()
		}

	case C.APP_CMD_WINDOW_REDRAW_NEEDED:
		println("window_redraw_needed")
	case C.APP_CMD_START:
		println("start")
	case C.APP_CMD_RESUME:
		println("resume")
	case C.APP_CMD_SAVE_STATE:
		println("save_state")
	case C.APP_CMD_PAUSE:
		println("pause")
	case C.APP_CMD_STOP:
		println("stop")
	case C.APP_CMD_DESTROY:
		println("destroy")
	}
}

func handleInput() {
	ib := C.android_app_swap_input_buffers((*C.struct_android_app)(atomic.LoadPointer(&androidApp)))
	if ib == nil {
		return
	}

	if ib.motionEventsCount > 0 {
		touchEvents := accumulateTouchEvents(ib)

		cbMut.Lock()
		touchInputCb := windowTouchInputCb
		cbMut.Unlock()

		if touchInputCb != nil {
			for _, ev := range touchEvents {
				touchInputCb(ev.phase, ev.pos, ev.id)
			}
		}

		C.android_app_clear_motion_events(ib)
	}

	if ib.keyEventsCount > 0 {
		handleKeyEvents(ib)
		C.android_app_clear_key_events(ib)
	}
}

type touchInputEvent struct {
	phase events.TouchPhase
	pos   dpi.PhysicalPosition[float64]
	id    events.TouchPointerID
}

var oldPosXs, oldPosYs [C.GAMEACTIVITY_MAX_NUM_POINTERS_IN_MOTION_EVENT]C.float
var oldPosMut sync.Mutex

func accumulateTouchEvents(ib *C.struct_android_input_buffer) []touchInputEvent {
	oldPosMut.Lock()
	defer oldPosMut.Unlock()

	touchEvents := make([]touchInputEvent, 0, ib.motionEventsCount)

	for i := C.uint64_t(0); i < ib.motionEventsCount; i++ {
		event := ib.motionEvents[i]

		mask := event.action & C.AMOTION_EVENT_ACTION_MASK
		var ptrIdx C.int32_t = C.GAMEACTIVITY_MAX_NUM_POINTERS_IN_MOTION_EVENT
		var phase events.TouchPhase

		switch mask {
		case C.AMOTION_EVENT_ACTION_POINTER_DOWN:
			ptrIdx = (event.action & C.AMOTION_EVENT_ACTION_POINTER_INDEX_MASK) >> C.AMOTION_EVENT_ACTION_POINTER_INDEX_SHIFT
			phase = events.TouchPhaseStarted

		case C.AMOTION_EVENT_ACTION_POINTER_UP:
			ptrIdx = (event.action & C.AMOTION_EVENT_ACTION_POINTER_INDEX_MASK) >> C.AMOTION_EVENT_ACTION_POINTER_INDEX_SHIFT
			phase = events.TouchPhaseEnded

		case C.AMOTION_EVENT_ACTION_DOWN:
			ptrIdx = 0
			phase = events.TouchPhaseStarted

		case C.AMOTION_EVENT_ACTION_UP:
			ptrIdx = 0
			phase = events.TouchPhaseEnded

		case C.AMOTION_EVENT_ACTION_MOVE:
			for ptrIdx, ptr := range event.pointers {
				if ptr.rawX == 0 && ptr.rawY == 0 {
					continue
				}

				oldX := oldPosXs[ptrIdx]
				oldY := oldPosYs[ptrIdx]
				newX := C.GameActivityPointerAxes_getAxisValue(&ptr, C.AMOTION_EVENT_AXIS_X)
				newY := C.GameActivityPointerAxes_getAxisValue(&ptr, C.AMOTION_EVENT_AXIS_Y)

				if oldX != newX && oldY != newY {
					oldPosXs[ptrIdx] = newX
					oldPosYs[ptrIdx] = newY

					touchEvents = append(touchEvents, touchInputEvent{
						phase: events.TouchPhaseMoved,
						pos: dpi.PhysicalPosition[float64]{
							X: float64(newX),
							Y: float64(newY),
						},
						id: events.TouchPointerID(ptr.id),
					})
				}
			}
		}

		if ptrIdx != C.GAMEACTIVITY_MAX_NUM_POINTERS_IN_MOTION_EVENT {
			ptr := event.pointers[ptrIdx]
			x := C.GameActivityPointerAxes_getAxisValue(&ptr, C.AMOTION_EVENT_AXIS_X)
			y := C.GameActivityPointerAxes_getAxisValue(&ptr, C.AMOTION_EVENT_AXIS_Y)

			if phase == events.TouchPhaseEnded {
				oldPosXs[ptrIdx] = 0
				oldPosYs[ptrIdx] = 0
			} else {
				oldPosXs[ptrIdx] = x
				oldPosYs[ptrIdx] = y
			}

			touchEvents = append(touchEvents, touchInputEvent{
				phase: phase,
				pos: dpi.PhysicalPosition[float64]{
					X: float64(x),
					Y: float64(y),
				},
				id: events.TouchPointerID(ptr.id),
			})
		}
	}

	return touchEvents
}

func handleKeyEvents(ib *C.struct_android_input_buffer) {
	for i := C.uint64_t(0); i < ib.keyEventsCount; i++ {
		event := ib.keyEvents[i]

		var state events.ButtonState

		switch event.action {
		case C.AKEY_EVENT_ACTION_DOWN,
			C.AKEY_EVENT_ACTION_MULTIPLE:
			state = events.ButtonStatePressed
		case C.AKEY_EVENT_ACTION_UP:
			state = events.ButtonStateReleased

		default:
			panic("unreachable")
		}

		cbMut.Lock()
		keyboardInputCb := windowKeyboardInputCb
		cbMut.Unlock()

		if keyboardInputCb != nil {
			keyboardInputCb(
				state,
				// TODO: current release of GameActivity doesn't expose scancode,
				// but it's available in master, we'll have to wait for next release
				events.ScanCode(0),
				mapKeycode(event.keyCode),
			)
		}

		cbMut.Lock()
		receivedCharacterCallback := windowReceivedCharacterCallback
		cbMut.Unlock()

		if event.unicodeChar != 0 &&
			state == events.ButtonStatePressed &&
			receivedCharacterCallback != nil {
			receivedCharacterCallback(rune(event.unicodeChar))
		}
	}
}
