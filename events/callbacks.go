package events

import (
	"strconv"

	"github.com/rajveermalviya/gamen/dpi"
)

type MouseScrollDelta uint8

const (
	MouseScrollDeltaLine MouseScrollDelta = iota
	MouseScrollDeltaPixel
)

func (v MouseScrollDelta) String() string {
	switch v {
	case MouseScrollDeltaLine:
		return "Line"
	case MouseScrollDeltaPixel:
		return "Pixel"
	default:
		return ""
	}
}

type MouseScrollAxis uint8

const (
	MouseScrollAxisVertical MouseScrollAxis = iota
	MouseScrollAxisHorizontal
)

func (v MouseScrollAxis) String() string {
	switch v {
	case MouseScrollAxisVertical:
		return "Vertical"
	case MouseScrollAxisHorizontal:
		return "Horizontal"
	default:
		return ""
	}
}

type ModifiersState uint8

const (
	ModifiersStateShift ModifiersState = 1 << iota
	ModifiersStateCtrl
	ModifiersStateAlt
	ModifiersStateLogo
)

func (v ModifiersState) String() (s string) {
	if v&ModifiersStateShift != 0 {
		if s == "" {
			s += "Shift"
		} else {
			s += " Shift"
		}
	}
	if v&ModifiersStateCtrl != 0 {
		if s == "" {
			s += "Ctrl"
		} else {
			s += " Ctrl"
		}
	}
	if v&ModifiersStateAlt != 0 {
		if s == "" {
			s += "Alt"
		} else {
			s += " Alt"
		}
	}
	if v&ModifiersStateLogo != 0 {
		if s == "" {
			s += "Logo"
		} else {
			s += " Logo"
		}
	}
	return
}

type ButtonState uint8

const (
	ButtonStatePressed ButtonState = iota
	ButtonStateReleased
)

func (v ButtonState) String() string {
	switch v {
	case ButtonStatePressed:
		return "Pressed"
	case ButtonStateReleased:
		return "Released"
	default:
		return ""
	}
}

type MouseButton uint32

const (
	MouseButtonLeft MouseButton = iota
	MouseButtonRight
	MouseButtonMiddle
)

func (v MouseButton) String() string {
	switch v {
	case MouseButtonLeft:
		return "Left"
	case MouseButtonRight:
		return "Right"
	case MouseButtonMiddle:
		return "Middle"
	default:
		return "0x" + strconv.FormatUint(uint64(v), 8)
	}
}

type TouchPointerID uint8
type TouchPhase uint32

const (
	TouchPhaseStarted TouchPhase = iota
	TouchPhaseMoved
	TouchPhaseEnded
	TouchPhaseCancelled
)

func (v TouchPhase) String() string {
	switch v {
	case TouchPhaseStarted:
		return "Started"
	case TouchPhaseMoved:
		return "Moved"
	case TouchPhaseEnded:
		return "Ended"
	case TouchPhaseCancelled:
		return "Cancelled"
	default:
		return ""
	}
}

type ScanCode uint32

type (
	WindowSurfaceCreatedCallback    func()
	WindowSurfaceDestroyedCallback  func()
	WindowCloseRequestedCallback    func()
	WindowResizedCallback           func(physicalWidth, physicalHeight uint32, scaleFactor float64)
	WindowCursorEnteredCallback     func()
	WindowCursorLeftCallback        func()
	WindowCursorMovedCallback       func(physicalX, physicalY float64)
	WindowMouseWheelCallback        func(delta MouseScrollDelta, axis MouseScrollAxis, value float64)
	WindowMouseInputCallback        func(state ButtonState, button MouseButton)
	WindowFocusedCallback           func()
	WindowUnfocusedCallback         func()
	WindowModifiersChangedCallback  func(state ModifiersState)
	WindowKeyboardInputCallback     func(state ButtonState, scanCode ScanCode, virtualKeyCode VirtualKey)
	WindowReceivedCharacterCallback func(char rune)
	WindowTouchInputCallback        func(phase TouchPhase, location dpi.PhysicalPosition[float64], id TouchPointerID)
)
