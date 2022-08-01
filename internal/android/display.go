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

	clearInputBuffers()
	return true
}
func (d *Display) Poll() bool {
	C.display_poll(0)

	app := (*C.struct_android_app)(atomic.LoadPointer(&androidApp))
	if app != nil && app.destroyRequested != 0 {
		return false
	}

	clearInputBuffers()
	return true
}
func (d *Display) WaitTimeout(t time.Duration) bool {
	C.display_poll(C.int(t.Milliseconds()))

	app := (*C.struct_android_app)(atomic.LoadPointer(&androidApp))
	if app != nil && app.destroyRequested != 0 {
		return false
	}

	clearInputBuffers()
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
		cbMut.Lock()
		surfaceCreatedCb := windowSurfaceCreatedCb
		cbMut.Unlock()

		if surfaceCreatedCb != nil {
			surfaceCreatedCb()
		}

	case C.APP_CMD_TERM_WINDOW:
		cbMut.Lock()
		surfaceDestroyedCb := windowSurfaceDestroyedCb
		cbMut.Unlock()

		if surfaceDestroyedCb != nil {
			surfaceDestroyedCb()
		}

	case C.APP_CMD_WINDOW_RESIZED,
		C.APP_CMD_CONFIG_CHANGED,
		C.APP_CMD_WINDOW_INSETS_CHANGED,
		C.APP_CMD_CONTENT_RECT_CHANGED:
		runResizedCallback()

	case C.APP_CMD_LOW_MEMORY:
		runtime.GC()

	case C.APP_CMD_GAINED_FOCUS:
		cbMut.Lock()
		focusedCb := windowFocusedCb
		cbMut.Unlock()

		if focusedCb != nil {
			focusedCb()
		}

	case C.APP_CMD_LOST_FOCUS:
		cbMut.Lock()
		unfocusedCb := windowUnfocusedCb
		cbMut.Unlock()

		if unfocusedCb != nil {
			unfocusedCb()
		}

	case C.APP_CMD_WINDOW_REDRAW_NEEDED:
	case C.APP_CMD_START:
	case C.APP_CMD_RESUME:
	case C.APP_CMD_SAVE_STATE:
	case C.APP_CMD_PAUSE:
	case C.APP_CMD_STOP:
	case C.APP_CMD_DESTROY:
	}
}

//export display_handle_key_event
func display_handle_key_event(event *C.GameActivityKeyEvent) C.bool {
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

	return true
}

var oldPosXs, oldPosYs [C.GAMEACTIVITY_MAX_NUM_POINTERS_IN_MOTION_EVENT]C.float
var oldPosMut sync.Mutex

//export display_handle_motion_event
func display_handle_motion_event(event *C.GameActivityMotionEvent) C.bool {
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

				cbMut.Lock()
				touchInputCb := windowTouchInputCb
				cbMut.Unlock()

				if touchInputCb != nil {
					touchInputCb(
						events.TouchPhaseMoved,
						dpi.PhysicalPosition[float64]{
							X: float64(newX),
							Y: float64(newY),
						},
						events.TouchPointerID(ptr.id),
					)
				}
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

		cbMut.Lock()
		touchInputCb := windowTouchInputCb
		cbMut.Unlock()

		if touchInputCb != nil {
			touchInputCb(
				phase,
				dpi.PhysicalPosition[float64]{
					X: float64(x),
					Y: float64(y),
				},
				events.TouchPointerID(ptr.id),
			)
		}
	}

	return true
}

func clearInputBuffers() {
	ib := C.android_app_swap_input_buffers((*C.struct_android_app)(atomic.LoadPointer(&androidApp)))
	if ib == nil {
		return
	}
	if ib.motionEventsCount > 0 {
		C.android_app_clear_motion_events(ib)
	}
	if ib.keyEventsCount > 0 {
		C.android_app_clear_key_events(ib)
	}
}
