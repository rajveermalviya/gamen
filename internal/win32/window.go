//go:build windows

package win32

import (
	"math"
	"sync"
	"unicode/utf16"
	"unsafe"

	"github.com/rajveermalviya/gamen/cursors"
	"github.com/rajveermalviya/gamen/dpi"
	"github.com/rajveermalviya/gamen/events"
	"github.com/rajveermalviya/gamen/internal/utils"
	"github.com/rajveermalviya/gamen/internal/win32/procs"
	"golang.org/x/sys/windows"
)

type Window struct {
	d *Display
	// window handle
	hwnd uintptr
	// we allow destroy function to be called multiple
	// times, but in reality we run it once
	destroyOnce sync.Once
	mu          sync.Mutex

	// state
	size               dpi.PhysicalSize[uint32]
	minSize            dpi.PhysicalSize[uint32]
	maxSize            dpi.PhysicalSize[uint32]
	cursorIsOutside    bool
	cursorPos          dpi.PhysicalPosition[int16]
	cursorCaptureCount int
	modifiers          events.ModifiersState
	currentCursor      cursors.Icon
	maximized          bool
	fullscreen         bool
	wpPrev             procs.WINDOWPLACEMENT

	// callbacks
	resizedCb           events.WindowResizedCallback
	closeRequestedCb    events.WindowCloseRequestedCallback
	focusedCb           events.WindowFocusedCallback
	unfocusedCb         events.WindowUnfocusedCallback
	cursorEnteredCb     events.WindowCursorEnteredCallback
	cursorLeftCb        events.WindowCursorLeftCallback
	cursorMovedCb       events.WindowCursorMovedCallback
	mouseWheelCb        events.WindowMouseWheelCallback
	mouseInputCb        events.WindowMouseInputCallback
	modifiersChangedCb  events.WindowModifiersChangedCallback
	keyboardInputCb     events.WindowKeyboardInputCallback
	receivedCharacterCb events.WindowReceivedCharacterCallback
}

var windowClassName = must(windows.UTF16PtrFromString("Window Class"))

func NewWindow(d *Display) (*Window, error) {
	class := procs.WNDCLASSEXW{
		CbSize:        uint32(unsafe.Sizeof(procs.WNDCLASSEXW{})),
		Style:         procs.CS_HREDRAW | procs.CS_VREDRAW,
		LpfnWndProc:   windowProcCb,
		CbClsExtra:    0,
		CbWndExtra:    0,
		HInstance:     procs.GetModuleHandleW(),
		HIcon:         0,
		HCursor:       0,
		HbrBackground: 0,
		LpszMenuName:  nil,
		LpszClassName: windowClassName,
		HIconSm:       0,
	}
	procs.RegisterClassExW(uintptr(unsafe.Pointer(&class)))

	w := &Window{
		currentCursor: cursors.Default,
		minSize:       dpi.PhysicalSize[uint32]{},
		maxSize: dpi.PhysicalSize[uint32]{
			Width:  math.MaxInt16,
			Height: math.MaxInt16,
		},
	}

	hwnd := procs.CreateWindowExW(
		procs.WS_EX_LEFT|
			procs.WS_EX_WINDOWEDGE|
			procs.WS_EX_APPWINDOW,
		uintptr(unsafe.Pointer(windowClassName)),
		0,
		procs.WS_OVERLAPPED| // show title & border
			procs.WS_SIZEBOX| // sizing border
			procs.WS_MAXIMIZEBOX| // maximize button
			procs.WS_CAPTION| // title bar
			procs.WS_MINIMIZEBOX| // minimize button
			procs.WS_BORDER| // thin border
			procs.WS_VISIBLE| // visible
			procs.WS_CLIPSIBLINGS| // clip window behind
			procs.WS_CLIPCHILDREN| // clip window behind
			procs.WS_SYSMENU, // window menu
		procs.CW_USEDEFAULT,
		procs.CW_USEDEFAULT,
		procs.CW_USEDEFAULT,
		procs.CW_USEDEFAULT,
		0,
		0,
		procs.GetModuleHandleW(),
		uintptr(unsafe.Pointer(w)),
	)

	if hwnd == 0 || w.hwnd == 0 {
		return nil, windows.GetLastError()
	}

	d.windows[hwnd] = w

	return w, nil
}

func (w *Window) Win32Hinstance() uintptr { return procs.GetModuleHandleW() }
func (w *Window) Win32Hwnd() uintptr      { return w.hwnd }

func (w *Window) Destroy() {
	w.destroyOnce.Do(func() {
		w.mu.Lock()
		w.resizedCb = nil
		w.closeRequestedCb = nil
		w.focusedCb = nil
		w.unfocusedCb = nil
		w.cursorEnteredCb = nil
		w.cursorLeftCb = nil
		w.cursorMovedCb = nil
		w.mouseWheelCb = nil
		w.mouseInputCb = nil
		w.modifiersChangedCb = nil
		w.keyboardInputCb = nil
		w.receivedCharacterCb = nil
		w.mu.Unlock()

		procs.DestroyWindow(w.hwnd)
	})
}

func (w *Window) SetTitle(title string) {
	titlePtr := windows.StringToUTF16Ptr(title)
	procs.SetWindowTextW(w.hwnd, uintptr(unsafe.Pointer(titlePtr)))
}

func (w *Window) InnerSize() dpi.PhysicalSize[uint32] {
	var rect procs.RECT
	if !procs.GetClientRect(w.hwnd, uintptr(unsafe.Pointer(&rect))) {
		panic("GetClientRect failed")
	}
	return dpi.PhysicalSize[uint32]{
		Width:  uint32(rect.Right),
		Height: uint32(rect.Bottom),
	}
}

func (w *Window) SetInnerSize(size dpi.Size[uint32]) {
	physicalSize := size.ToPhysical(1)

	rect := &procs.RECT{
		Top:    0,
		Left:   0,
		Bottom: int32(physicalSize.Height),
		Right:  int32(physicalSize.Width),
	}

	style := procs.GetWindowLong(w.hwnd, procs.GWL_STYLE)
	styleEx := procs.GetWindowLong(w.hwnd, procs.GWL_EXSTYLE)
	if !procs.AdjustWindowRectEx(w.hwnd, style, styleEx, uintptr(unsafe.Pointer(rect))) {
		panic("AdjustWindowRectEx failed")
	}

	outerX := utils.Abs(rect.Right - rect.Left)
	outerY := utils.Abs(rect.Top - rect.Bottom)

	procs.SetWindowPos(
		w.hwnd,
		0,
		0, 0,
		uintptr(outerX),
		uintptr(outerY),
		procs.SWP_ASYNCWINDOWPOS|
			procs.SWP_NOZORDER|
			procs.SWP_NOREPOSITION|
			procs.SWP_NOMOVE|
			procs.SWP_NOACTIVATE,
	)

	procs.InvalidateRgn(w.hwnd, 0, 0)
}

func (w *Window) SetMinInnerSize(size dpi.Size[uint32]) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.minSize = size.ToPhysical(1)
}

func (w *Window) SetMaxInnerSize(size dpi.Size[uint32]) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.maxSize = size.ToPhysical(1)
}

func (w *Window) Maximized() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.maximized
}

func (w *Window) SetMinimized() {
	procs.ShowWindow(w.hwnd, procs.SW_MINIMIZE)
}

func (w *Window) SetMaximized(maximized bool) {
	if maximized {
		procs.ShowWindow(w.hwnd, procs.SW_MAXIMIZE)
	} else {
		procs.ShowWindow(w.hwnd, procs.SW_RESTORE)
	}
}

func (w *Window) SetCursorIcon(icon cursors.Icon) {
	w.mu.Lock()
	w.currentCursor = icon
	w.mu.Unlock()

	procs.SetCursor(procs.LoadCursorW(0, toWin32Cursor(icon)))
}

func (w *Window) SetCursorVisible(visible bool) {
	if visible {
		procs.ShowCursor(1)
	} else {
		procs.ShowCursor(0)
	}
}

func (w *Window) SetFullscreen(fullscreen bool) {
	dwStyle := procs.GetWindowLong(w.hwnd, procs.GWL_STYLE)

	if fullscreen {
		mi := procs.MONITORINFO{
			CbSize: uint32(unsafe.Sizeof(procs.MONITORINFO{})),
		}
		monitor := procs.MonitorFromWindow(w.hwnd, procs.MONITOR_DEFAULTTOPRIMARY)

		var wp procs.WINDOWPLACEMENT
		if procs.GetWindowPlacement(w.hwnd, uintptr(unsafe.Pointer(&wp))) &&
			procs.GetMonitorInfoW(monitor, uintptr(unsafe.Pointer(&mi))) {
			w.mu.Lock()
			w.wpPrev = wp
			w.mu.Unlock()

			procs.SetWindowLong(w.hwnd, procs.GWL_STYLE, dwStyle&^procs.WS_OVERLAPPEDWINDOW)
			procs.SetWindowPos(
				w.hwnd, procs.HWND_TOP,
				uintptr(mi.RcMonitor.Left), uintptr(mi.RcMonitor.Top),
				uintptr(mi.RcMonitor.Right-mi.RcMonitor.Left),
				uintptr(mi.RcMonitor.Bottom-mi.RcMonitor.Top),
				procs.SWP_NOOWNERZORDER|procs.SWP_FRAMECHANGED,
			)
		}
	} else {
		w.mu.Lock()
		wp := w.wpPrev
		w.mu.Unlock()

		procs.SetWindowLong(w.hwnd, procs.GWL_STYLE, dwStyle|procs.WS_OVERLAPPEDWINDOW)
		procs.SetWindowPlacement(w.hwnd, uintptr(unsafe.Pointer(&wp)))
		procs.SetWindowPos(
			w.hwnd, 0,
			0, 0,
			0, 0,
			procs.SWP_NOMOVE|procs.SWP_NOSIZE|
				procs.SWP_NOZORDER|procs.SWP_NOOWNERZORDER|
				procs.SWP_FRAMECHANGED,
		)
	}
}

func (w *Window) Fullscreen() bool {
	dwStyle := procs.GetWindowLong(w.hwnd, procs.GWL_STYLE)
	if dwStyle&procs.WS_OVERLAPPEDWINDOW != 0 {
		return false
	}
	return true
}

func (w *Window) SetCloseRequestedCallback(cb events.WindowCloseRequestedCallback) {
	w.mu.Lock()
	w.closeRequestedCb = cb
	w.mu.Unlock()
}
func (w *Window) SetResizedCallback(cb events.WindowResizedCallback) {
	w.mu.Lock()
	w.resizedCb = cb
	w.mu.Unlock()
}
func (w *Window) SetFocusedCallback(cb events.WindowFocusedCallback) {
	w.mu.Lock()
	w.focusedCb = cb
	w.mu.Unlock()
}
func (w *Window) SetUnfocusedCallback(cb events.WindowUnfocusedCallback) {
	w.mu.Lock()
	w.unfocusedCb = cb
	w.mu.Unlock()
}
func (w *Window) SetCursorEnteredCallback(cb events.WindowCursorEnteredCallback) {
	w.mu.Lock()
	w.cursorEnteredCb = cb
	w.mu.Unlock()
}
func (w *Window) SetCursorLeftCallback(cb events.WindowCursorLeftCallback) {
	w.mu.Lock()
	w.cursorLeftCb = cb
	w.mu.Unlock()
}
func (w *Window) SetCursorMovedCallback(cb events.WindowCursorMovedCallback) {
	w.mu.Lock()
	w.cursorMovedCb = cb
	w.mu.Unlock()
}
func (w *Window) SetMouseWheelCallback(cb events.WindowMouseWheelCallback) {
	w.mu.Lock()
	w.mouseWheelCb = cb
	w.mu.Unlock()
}
func (w *Window) SetMouseInputCallback(cb events.WindowMouseInputCallback) {
	w.mu.Lock()
	w.mouseInputCb = cb
	w.mu.Unlock()
}
func (w *Window) SetTouchInputCallback(cb events.WindowTouchInputCallback) {
	// TODO:
}
func (w *Window) SetModifiersChangedCallback(cb events.WindowModifiersChangedCallback) {
	w.mu.Lock()
	w.modifiersChangedCb = cb
	w.mu.Unlock()
}
func (w *Window) SetKeyboardInputCallback(cb events.WindowKeyboardInputCallback) {
	w.mu.Lock()
	w.keyboardInputCb = cb
	w.mu.Unlock()
}
func (w *Window) SetReceivedCharacterCallback(cb events.WindowReceivedCharacterCallback) {
	w.mu.Lock()
	w.receivedCharacterCb = cb
	w.mu.Unlock()
}

var windowProcCb = windows.NewCallback(windowProc)

func windowProc(window, msg, wparam, lparam uintptr) uintptr {
	userData := procs.GetWindowLong(window, procs.GWL_USERDATA)

	if userData == 0 {
		if msg == procs.WM_NCCREATE {
			createStruct := (*procs.CREATESTRUCTW)(unsafe.Pointer(lparam))
			w := (*Window)(unsafe.Pointer(createStruct.LpCreateParams))
			w.hwnd = window
			procs.SetWindowLong(window, procs.GWL_USERDATA, uintptr(unsafe.Pointer(w)))
		}
		return procs.DefWindowProcW(window, msg, wparam, lparam)
	}

	w := (*Window)(unsafe.Pointer(userData))

	switch msg {
	case procs.WM_CLOSE:

		w.mu.Lock()
		var closeRequestedCb events.WindowCloseRequestedCallback
		if w.closeRequestedCb != nil {
			closeRequestedCb = w.closeRequestedCb
		}
		w.mu.Unlock()

		if closeRequestedCb != nil {
			closeRequestedCb()
		}
		return 0

	case procs.WM_SIZE:

		size := dpi.PhysicalSize[uint32]{
			Width:  uint32(loword(uint32(lparam))),
			Height: uint32(hiword(uint32(lparam))),
		}

		w.mu.Lock()
		w.maximized = false
		if wparam == procs.SIZE_MAXIMIZED {
			w.maximized = true
		}
		w.size = size
		var resizedCb events.WindowResizedCallback
		if w.resizedCb != nil {
			resizedCb = w.resizedCb
		}
		w.mu.Unlock()

		if resizedCb != nil {
			resizedCb(size.Width, size.Height, 1)
		}
		return 0

	case procs.WM_GETMINMAXINFO:
		mmi := (*procs.MINMAXINFO)(unsafe.Pointer(lparam))

		w.mu.Lock()
		minSize := w.minSize
		w.mu.Unlock()

		var zero dpi.PhysicalSize[uint32]

		if minSize != zero {
			rect := &procs.RECT{
				Top:    0,
				Left:   0,
				Bottom: int32(minSize.Height),
				Right:  int32(minSize.Width),
			}

			style := procs.GetWindowLong(w.hwnd, procs.GWL_STYLE)
			styleEx := procs.GetWindowLong(w.hwnd, procs.GWL_EXSTYLE)
			if !procs.AdjustWindowRectEx(w.hwnd, style, styleEx, uintptr(unsafe.Pointer(rect))) {
				panic("AdjustWindowRectEx failed")
			}

			mmi.PtMinTrackSize = procs.POINT{
				X: rect.Right - rect.Left,
				Y: rect.Bottom - rect.Top,
			}
		}

		w.mu.Lock()
		maxSize := w.maxSize
		w.mu.Unlock()

		if maxSize != zero {
			rect := &procs.RECT{
				Top:    0,
				Left:   0,
				Bottom: int32(maxSize.Height),
				Right:  int32(maxSize.Width),
			}

			style := procs.GetWindowLong(w.hwnd, procs.GWL_STYLE)
			styleEx := procs.GetWindowLong(w.hwnd, procs.GWL_EXSTYLE)
			if !procs.AdjustWindowRectEx(w.hwnd, style, styleEx, uintptr(unsafe.Pointer(rect))) {
				panic("AdjustWindowRectEx failed")
			}

			mmi.PtMaxTrackSize = procs.POINT{
				X: rect.Right - rect.Left,
				Y: rect.Bottom - rect.Top,
			}
		}

		return 0

	case procs.WM_SETCURSOR:
		if loword(uint32(lparam)) == procs.HTCLIENT {
			procs.SetCursor(procs.LoadCursorW(0, toWin32Cursor(w.currentCursor)))
		}

	case procs.WM_MOUSEMOVE:
		if w.cursorIsOutside {
			w.cursorIsOutside = false

			w.mu.Lock()
			var cursorEnteredCb events.WindowCursorEnteredCallback
			if w.cursorEnteredCb != nil {
				cursorEnteredCb = w.cursorEnteredCb
			}
			w.mu.Unlock()

			if cursorEnteredCb != nil {
				cursorEnteredCb()
			}

			procs.TrackMouseEvent(uintptr(unsafe.Pointer(&procs.TRACKMOUSEEVENT{
				CbSize:      uint32(unsafe.Sizeof(procs.TRACKMOUSEEVENT{})),
				DwFlags:     procs.TME_LEAVE,
				HwndTrack:   window,
				DwHoverTime: procs.HOVER_DEFAULT,
			})))
		}

		pos := dpi.PhysicalPosition[int16]{
			X: int16(loword(uint32(lparam))),
			Y: int16(hiword(uint32(lparam))),
		}

		if w.cursorPos != pos {
			w.cursorPos = pos

			w.mu.Lock()
			var cursorMovedCb events.WindowCursorMovedCallback
			if w.cursorMovedCb != nil {
				cursorMovedCb = w.cursorMovedCb
			}
			w.mu.Unlock()

			if cursorMovedCb != nil {
				cursorMovedCb(float64(pos.X), float64(pos.Y))
			}
		}
		return 0

	case procs.WM_MOUSELEAVE:
		w.cursorIsOutside = true

		w.mu.Lock()
		var cursorLeftCb events.WindowCursorLeftCallback
		if w.cursorLeftCb != nil {
			cursorLeftCb = w.cursorLeftCb
		}
		w.mu.Unlock()

		if cursorLeftCb != nil {
			cursorLeftCb()
		}
		return 0

	case procs.WM_MOUSEWHEEL:
		value := float64(int16(wparam>>16)) / procs.WHEEL_DELTA

		w.mu.Lock()
		var mouseWheelCb events.WindowMouseWheelCallback
		if w.mouseWheelCb != nil {
			mouseWheelCb = w.mouseWheelCb
		}
		w.mu.Unlock()

		if mouseWheelCb != nil {
			mouseWheelCb(
				events.MouseScrollDeltaLine,
				events.MouseScrollAxisVertical,
				value,
			)
		}
		return 0

	case procs.WM_MOUSEHWHEEL:
		value := -float64(int16(wparam>>16)) / procs.WHEEL_DELTA

		w.mu.Lock()
		var mouseWheelCb events.WindowMouseWheelCallback
		if w.mouseWheelCb != nil {
			mouseWheelCb = w.mouseWheelCb
		}
		w.mu.Unlock()

		if mouseWheelCb != nil {
			mouseWheelCb(
				events.MouseScrollDeltaLine,
				events.MouseScrollAxisHorizontal,
				value,
			)
		}
		return 0

	case procs.WM_LBUTTONDOWN:
		procs.SetCapture(window)
		w.cursorCaptureCount++

		w.mu.Lock()
		var mouseInputCb events.WindowMouseInputCallback
		if w.mouseInputCb != nil {
			mouseInputCb = w.mouseInputCb
		}
		w.mu.Unlock()

		if mouseInputCb != nil {
			mouseInputCb(
				events.ButtonStatePressed,
				events.MouseButtonLeft,
			)
		}
		return 0

	case procs.WM_LBUTTONUP:
		w.cursorCaptureCount = utils.Max(0, w.cursorCaptureCount-1)
		if w.cursorCaptureCount == 0 {
			procs.ReleaseCapture()
		}

		w.mu.Lock()
		var mouseInputCb events.WindowMouseInputCallback
		if w.mouseInputCb != nil {
			mouseInputCb = w.mouseInputCb
		}
		w.mu.Unlock()

		if mouseInputCb != nil {
			mouseInputCb(
				events.ButtonStateReleased,
				events.MouseButtonLeft,
			)
		}
		return 0

	case procs.WM_RBUTTONDOWN:
		procs.SetCapture(window)
		w.cursorCaptureCount++

		w.mu.Lock()
		var mouseInputCb events.WindowMouseInputCallback
		if w.mouseInputCb != nil {
			mouseInputCb = w.mouseInputCb
		}
		w.mu.Unlock()

		if mouseInputCb != nil {
			mouseInputCb(
				events.ButtonStatePressed,
				events.MouseButtonRight,
			)
		}
		return 0

	case procs.WM_RBUTTONUP:
		w.cursorCaptureCount = utils.Max(0, w.cursorCaptureCount-1)
		if w.cursorCaptureCount == 0 {
			procs.ReleaseCapture()
		}

		w.mu.Lock()
		var mouseInputCb events.WindowMouseInputCallback
		if w.mouseInputCb != nil {
			mouseInputCb = w.mouseInputCb
		}
		w.mu.Unlock()

		if mouseInputCb != nil {
			mouseInputCb(
				events.ButtonStateReleased,
				events.MouseButtonRight,
			)
		}
		return 0

	case procs.WM_MBUTTONDOWN:
		procs.SetCapture(window)
		w.cursorCaptureCount++

		w.mu.Lock()
		var mouseInputCb events.WindowMouseInputCallback
		if w.mouseInputCb != nil {
			mouseInputCb = w.mouseInputCb
		}
		w.mu.Unlock()

		if mouseInputCb != nil {
			mouseInputCb(
				events.ButtonStatePressed,
				events.MouseButtonMiddle,
			)
		}
		return 0

	case procs.WM_MBUTTONUP:
		w.cursorCaptureCount = utils.Max(0, w.cursorCaptureCount-1)
		if w.cursorCaptureCount == 0 {
			procs.ReleaseCapture()
		}

		w.mu.Lock()
		var mouseInputCb events.WindowMouseInputCallback
		if w.mouseInputCb != nil {
			mouseInputCb = w.mouseInputCb
		}
		w.mu.Unlock()

		if mouseInputCb != nil {
			mouseInputCb(
				events.ButtonStateReleased,
				events.MouseButtonMiddle,
			)
		}
		return 0

	case procs.WM_XBUTTONDOWN:
		procs.SetCapture(window)
		w.cursorCaptureCount++

		w.mu.Lock()
		var mouseInputCb events.WindowMouseInputCallback
		if w.mouseInputCb != nil {
			mouseInputCb = w.mouseInputCb
		}
		w.mu.Unlock()

		if mouseInputCb != nil {
			mouseInputCb(
				events.ButtonStatePressed,
				events.MouseButton(loword(uint32(wparam))),
			)
		}
		return 0

	case procs.WM_XBUTTONUP:
		w.cursorCaptureCount = utils.Max(0, w.cursorCaptureCount-1)
		if w.cursorCaptureCount == 0 {
			procs.ReleaseCapture()
		}

		w.mu.Lock()
		var mouseInputCb events.WindowMouseInputCallback
		if w.mouseInputCb != nil {
			mouseInputCb = w.mouseInputCb
		}
		w.mu.Unlock()

		if mouseInputCb != nil {
			mouseInputCb(
				events.ButtonStateReleased,
				events.MouseButton(loword(uint32(wparam))),
			)
		}
		return 0

	case procs.WM_CAPTURECHANGED:
		if lparam != window {
			w.cursorCaptureCount = 0
		}
		return 0

	case procs.WM_SETFOCUS:
		w.mu.Lock()
		var focusedCb events.WindowFocusedCallback
		if w.focusedCb != nil {
			focusedCb = w.focusedCb
		}
		w.mu.Unlock()

		if focusedCb != nil {
			focusedCb()
		}

		if m := getModifiersState(); m != 0 {
			w.modifiers = m

			w.mu.Lock()
			var modifiersChangedCb events.WindowModifiersChangedCallback
			if w.modifiersChangedCb != nil {
				modifiersChangedCb = w.modifiersChangedCb
			}
			w.mu.Unlock()

			if modifiersChangedCb != nil {
				modifiersChangedCb(m)
			}
		}
		return 0

	case procs.WM_KILLFOCUS:
		if w.modifiers != 0 {
			w.modifiers = 0

			w.mu.Lock()
			var modifiersChangedCb events.WindowModifiersChangedCallback
			if w.modifiersChangedCb != nil {
				modifiersChangedCb = w.modifiersChangedCb
			}
			w.mu.Unlock()

			if modifiersChangedCb != nil {
				modifiersChangedCb(0)
			}
		}

		w.mu.Lock()
		var unfocusedCb events.WindowUnfocusedCallback
		if w.unfocusedCb != nil {
			unfocusedCb = w.unfocusedCb
		}
		w.mu.Unlock()

		if unfocusedCb != nil {
			unfocusedCb()
		}
		return 0

	case procs.WM_KEYDOWN, procs.WM_SYSKEYDOWN:
		if msg == procs.WM_SYSKEYDOWN && wparam == procs.VK_F4 {
			return procs.DefWindowProcW(window, msg, wparam, lparam)
		}

		w.mu.Lock()
		var keyboardInputCb events.WindowKeyboardInputCallback
		if w.keyboardInputCb != nil {
			keyboardInputCb = w.keyboardInputCb
		}
		w.mu.Unlock()

		if keyboardInputCb != nil {
			scancode := ((lparam >> 16) & 0xff)
			extended := (lparam & 0x01000000) != 0
			if extended {
				scancode |= 0xE000
			} else {
				scancode |= 0x0000
			}

			keyboardInputCb(
				events.ButtonStatePressed,
				events.ScanCode(scancode),
				mapVK(wparam, scancode, extended),
			)
		}

		m := getModifiersState()
		if w.modifiers != m {
			w.modifiers = m

			w.mu.Lock()
			var modifiersChangedCb events.WindowModifiersChangedCallback
			if w.modifiersChangedCb != nil {
				modifiersChangedCb = w.modifiersChangedCb
			}
			w.mu.Unlock()

			if modifiersChangedCb != nil {
				modifiersChangedCb(m)
			}
		}

		w.mu.Lock()
		var receivedCharacterCb events.WindowReceivedCharacterCallback
		if w.receivedCharacterCb != nil {
			receivedCharacterCb = w.receivedCharacterCb
		}
		w.mu.Unlock()

		// win32 doesn't send WM_CHAR message for delete key
		if receivedCharacterCb != nil && wparam == procs.VK_DELETE {
			receivedCharacterCb('\u007F')
		}
		return 0

	case procs.WM_KEYUP, procs.WM_SYSKEYUP:

		w.mu.Lock()
		var keyboardInputCb events.WindowKeyboardInputCallback
		if w.keyboardInputCb != nil {
			keyboardInputCb = w.keyboardInputCb
		}
		w.mu.Unlock()

		if keyboardInputCb != nil {
			scancode := ((lparam >> 16) & 0xff)
			extended := (lparam & 0x01000000) != 0
			if extended {
				scancode |= 0xE000
			} else {
				scancode |= 0x0000
			}

			keyboardInputCb(
				events.ButtonStateReleased,
				events.ScanCode(scancode),
				mapVK(wparam, scancode, extended),
			)
		}

		m := getModifiersState()
		if w.modifiers != m {
			w.modifiers = m

			w.mu.Lock()
			var modifiersChangedCb events.WindowModifiersChangedCallback
			if w.modifiersChangedCb != nil {
				modifiersChangedCb = w.modifiersChangedCb
			}
			w.mu.Unlock()

			if modifiersChangedCb != nil {
				modifiersChangedCb(m)
			}
		}
		return 0

	case procs.WM_CHAR, procs.WM_SYSCHAR:

		w.mu.Lock()
		var receivedCharacterCb events.WindowReceivedCharacterCallback
		if w.receivedCharacterCb != nil {
			receivedCharacterCb = w.receivedCharacterCb
		}
		w.mu.Unlock()

		if receivedCharacterCb != nil {
			for _, v := range utf16.Decode([]uint16{uint16(wparam)}) {
				receivedCharacterCb(v)
			}
		}
		return 0
	}

	return procs.DefWindowProcW(window, msg, wparam, lparam)
}

func getModifiersState() (m events.ModifiersState) {
	var state [256]byte
	procs.GetKeyboardState(uintptr(unsafe.Pointer(&state)))

	for i, v := range state {
		if v&(1<<7) != 0 { // if pressed
			switch i {
			case procs.VK_SHIFT:
				m |= events.ModifiersStateShift
			case procs.VK_CONTROL:
				m |= events.ModifiersStateCtrl
			case procs.VK_MENU:
				m |= events.ModifiersStateAlt
			case procs.VK_LWIN, procs.VK_RWIN:
				m |= events.ModifiersStateLogo
			}
		}
	}

	return
}
