//go:build windows

package win32

import (
	"math"
	"sync"
	"unsafe"

	"github.com/rajveermalviya/gamen/cursors"
	"github.com/rajveermalviya/gamen/dpi"
	"github.com/rajveermalviya/gamen/events"
	"github.com/rajveermalviya/gamen/internal/common/atomicx"
	"github.com/rajveermalviya/gamen/internal/common/mathx"
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

	minSize dpi.PhysicalSize[uint32] // shared mutex
	maxSize dpi.PhysicalSize[uint32] // shared mutex
	wpPrev  procs.WINDOWPLACEMENT    // shared mutex

	currentCursor atomicx.Uint[cursors.Icon] // shared atomic
	maximized     atomicx.Bool               // shared atomic

	cursorIsOutside    bool                        // non-shared
	cursorPos          dpi.PhysicalPosition[int16] // non-shared
	cursorCaptureCount int                         // non-shared
	modifiers          events.ModifiersState       // non-shared

	highSurrogated rune // non-shared

	// callbacks
	resizedCb           atomicx.Pointer[events.WindowResizedCallback]
	closeRequestedCb    atomicx.Pointer[events.WindowCloseRequestedCallback]
	focusedCb           atomicx.Pointer[events.WindowFocusedCallback]
	unfocusedCb         atomicx.Pointer[events.WindowUnfocusedCallback]
	cursorEnteredCb     atomicx.Pointer[events.WindowCursorEnteredCallback]
	cursorLeftCb        atomicx.Pointer[events.WindowCursorLeftCallback]
	cursorMovedCb       atomicx.Pointer[events.WindowCursorMovedCallback]
	mouseWheelCb        atomicx.Pointer[events.WindowMouseScrollCallback]
	mouseInputCb        atomicx.Pointer[events.WindowMouseInputCallback]
	modifiersChangedCb  atomicx.Pointer[events.WindowModifiersChangedCallback]
	keyboardInputCb     atomicx.Pointer[events.WindowKeyboardInputCallback]
	receivedCharacterCb atomicx.Pointer[events.WindowReceivedCharacterCallback]
}

var windowClassName = must(windows.UTF16PtrFromString("Window Class"))

const decoratedWindowStyles = procs.WS_OVERLAPPED |
	procs.WS_SYSMENU |
	procs.WS_CAPTION |
	procs.WS_SIZEBOX

const decoratedWindowExStyles = procs.WS_EX_WINDOWEDGE

const defaultStyles = procs.WS_VISIBLE | // visible
	procs.WS_CLIPSIBLINGS | // clip window behind
	procs.WS_CLIPCHILDREN | // clip window behind
	procs.WS_MAXIMIZEBOX |
	procs.WS_MINIMIZEBOX |
	decoratedWindowStyles

const defaultExStyles = procs.WS_EX_LEFT |
	procs.WS_EX_APPWINDOW |
	decoratedWindowExStyles

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
		minSize: dpi.PhysicalSize[uint32]{},
		maxSize: dpi.PhysicalSize[uint32]{
			Width:  math.MaxInt16,
			Height: math.MaxInt16,
		},
	}
	w.currentCursor.Store(cursors.Default)

	hwnd := procs.CreateWindowExW(
		defaultExStyles,
		uintptr(unsafe.Pointer(windowClassName)),
		0,
		defaultStyles,
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
		w.resizedCb.Store(nil)
		w.closeRequestedCb.Store(nil)
		w.focusedCb.Store(nil)
		w.unfocusedCb.Store(nil)
		w.cursorEnteredCb.Store(nil)
		w.cursorLeftCb.Store(nil)
		w.cursorMovedCb.Store(nil)
		w.mouseWheelCb.Store(nil)
		w.mouseInputCb.Store(nil)
		w.modifiersChangedCb.Store(nil)
		w.keyboardInputCb.Store(nil)
		w.receivedCharacterCb.Store(nil)

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

	outerX := mathx.Abs(rect.Right - rect.Left)
	outerY := mathx.Abs(rect.Top - rect.Bottom)

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
	return w.maximized.Load()
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
	w.currentCursor.Store(icon)
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
	dwExStyle := procs.GetWindowLong(w.hwnd, procs.GWL_EXSTYLE)

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

			procs.SetWindowLong(w.hwnd, procs.GWL_STYLE, dwStyle&^decoratedWindowStyles)
			procs.SetWindowLong(w.hwnd, procs.GWL_EXSTYLE, dwExStyle&^decoratedWindowExStyles)
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

		procs.SetWindowLong(w.hwnd, procs.GWL_STYLE, dwStyle|decoratedWindowStyles)
		procs.SetWindowLong(w.hwnd, procs.GWL_EXSTYLE, dwExStyle|decoratedWindowExStyles)
		procs.SetWindowPlacement(w.hwnd, uintptr(unsafe.Pointer(&wp)))
		procs.SetWindowPos(
			w.hwnd, 0,
			0, 0,
			0, 0,
			procs.SWP_NOMOVE|
				procs.SWP_NOSIZE|
				procs.SWP_NOZORDER|
				procs.SWP_NOOWNERZORDER|
				procs.SWP_FRAMECHANGED,
		)
	}
}

func (w *Window) Fullscreen() bool {
	windowSize := w.InnerSize()

	monitor := procs.MonitorFromWindow(w.hwnd, procs.MONITOR_DEFAULTTOPRIMARY)
	mi := procs.MONITORINFO{CbSize: uint32(unsafe.Sizeof(procs.MONITORINFO{}))}
	procs.GetMonitorInfoW(monitor, uintptr(unsafe.Pointer(&mi)))

	if int32(windowSize.Width) != mathx.Abs(mi.RcMonitor.Right-mi.RcMonitor.Left) ||
		int32(windowSize.Height) != mathx.Abs(mi.RcMonitor.Bottom-mi.RcMonitor.Top) {
		return false
	}
	if w.Decorated() {
		return false
	}
	return true
}

func (w *Window) DragWindow() {
	var pos procs.POINT
	procs.GetCursorPos(uintptr(unsafe.Pointer(&pos)))

	points := procs.POINTS{
		X: int16(pos.X),
		Y: int16(pos.Y),
	}
	procs.ReleaseCapture()
	procs.PostMessageW(
		w.hwnd,
		procs.WM_NCLBUTTONDOWN,
		procs.HTCAPTION,
		uintptr(unsafe.Pointer(&points)),
	)
}

func (w *Window) SetDecorations(decorate bool) {
	dwStyle := procs.GetWindowLong(w.hwnd, procs.GWL_STYLE)
	dwExStyle := procs.GetWindowLong(w.hwnd, procs.GWL_EXSTYLE)

	if decorate {
		dwStyle |= decoratedWindowStyles
		dwExStyle |= decoratedWindowExStyles
	} else {
		dwStyle &^= decoratedWindowStyles
		dwExStyle &^= decoratedWindowExStyles
	}

	procs.SetWindowLong(w.hwnd, procs.GWL_STYLE, dwStyle)
	procs.SetWindowLong(w.hwnd, procs.GWL_EXSTYLE, dwExStyle)
	procs.SetWindowPos(
		w.hwnd, 0,
		0, 0,
		0, 0,
		procs.SWP_NOMOVE|
			procs.SWP_NOSIZE|
			procs.SWP_NOZORDER|
			procs.SWP_NOOWNERZORDER|
			procs.SWP_FRAMECHANGED,
	)
}

func (w *Window) Decorated() bool {
	dwStyle := procs.GetWindowLong(w.hwnd, procs.GWL_STYLE)
	dwExStyle := procs.GetWindowLong(w.hwnd, procs.GWL_EXSTYLE)

	if dwStyle&decoratedWindowStyles != 0 &&
		dwExStyle&decoratedWindowExStyles != 0 {
		return true
	}
	return false
}

func (w *Window) SetCloseRequestedCallback(cb events.WindowCloseRequestedCallback) {
	w.closeRequestedCb.Store(&cb)
}
func (w *Window) SetResizedCallback(cb events.WindowResizedCallback) {
	w.resizedCb.Store(&cb)
}
func (w *Window) SetFocusedCallback(cb events.WindowFocusedCallback) {
	w.focusedCb.Store(&cb)
}
func (w *Window) SetUnfocusedCallback(cb events.WindowUnfocusedCallback) {
	w.unfocusedCb.Store(&cb)
}
func (w *Window) SetCursorEnteredCallback(cb events.WindowCursorEnteredCallback) {
	w.cursorEnteredCb.Store(&cb)
}
func (w *Window) SetCursorLeftCallback(cb events.WindowCursorLeftCallback) {
	w.cursorLeftCb.Store(&cb)
}
func (w *Window) SetCursorMovedCallback(cb events.WindowCursorMovedCallback) {
	w.cursorMovedCb.Store(&cb)
}
func (w *Window) SetMouseScrollCallback(cb events.WindowMouseScrollCallback) {
	w.mouseWheelCb.Store(&cb)
}
func (w *Window) SetMouseInputCallback(cb events.WindowMouseInputCallback) {
	w.mouseInputCb.Store(&cb)
}
func (w *Window) SetTouchInputCallback(cb events.WindowTouchInputCallback) {
	// TODO:
}
func (w *Window) SetModifiersChangedCallback(cb events.WindowModifiersChangedCallback) {
	w.modifiersChangedCb.Store(&cb)
}
func (w *Window) SetKeyboardInputCallback(cb events.WindowKeyboardInputCallback) {
	w.keyboardInputCb.Store(&cb)
}
func (w *Window) SetReceivedCharacterCallback(cb events.WindowReceivedCharacterCallback) {
	w.receivedCharacterCb.Store(&cb)
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

		if cb := w.closeRequestedCb.Load(); cb != nil {
			if cb := (*cb); cb != nil {
				cb()
			}
		}
		return 0

	case procs.WM_SIZE:
		size := dpi.PhysicalSize[uint32]{
			Width:  uint32(loword(uint32(lparam))),
			Height: uint32(hiword(uint32(lparam))),
		}

		if wparam == procs.SIZE_MAXIMIZED {
			w.maximized.Store(true)
		} else {
			w.maximized.Store(false)
		}

		if size.Width != 0 && size.Height != 0 {
			if cb := w.resizedCb.Load(); cb != nil {
				if cb := (*cb); cb != nil {
					cb(size.Width, size.Height, 1)
				}
			}
		}
		return 0

	case procs.WM_GETMINMAXINFO:
		mmi := (*procs.MINMAXINFO)(unsafe.Pointer(lparam))

		w.mu.Lock()
		minSize := w.minSize
		w.mu.Unlock()

		if minSize.Width != 0 && minSize.Height != 0 {
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

		if maxSize.Width != 0 && maxSize.Height != 0 {
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
			procs.SetCursor(procs.LoadCursorW(0, toWin32Cursor(w.currentCursor.Load())))
		}

	case procs.WM_MOUSEMOVE:
		if w.cursorIsOutside {
			w.cursorIsOutside = false

			if cb := w.cursorEnteredCb.Load(); cb != nil {
				if cb := (*cb); cb != nil {
					cb()
				}
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

			if cb := w.cursorMovedCb.Load(); cb != nil {
				if cb := (*cb); cb != nil {
					cb(float64(pos.X), float64(pos.Y))
				}
			}
		}
		return 0

	case procs.WM_MOUSELEAVE:
		w.cursorIsOutside = true

		if cb := w.cursorLeftCb.Load(); cb != nil {
			if cb := (*cb); cb != nil {
				cb()
			}
		}
		return 0

	case procs.WM_MOUSEWHEEL:
		value := float64(int16(wparam>>16)) / procs.WHEEL_DELTA

		if cb := w.mouseWheelCb.Load(); cb != nil {
			if cb := (*cb); cb != nil {
				cb(
					events.MouseScrollDeltaLine,
					events.MouseScrollAxisVertical,
					value,
				)
			}
		}
		return 0

	case procs.WM_MOUSEHWHEEL:
		value := -float64(int16(wparam>>16)) / procs.WHEEL_DELTA

		if cb := w.mouseWheelCb.Load(); cb != nil {
			if cb := (*cb); cb != nil {
				cb(
					events.MouseScrollDeltaLine,
					events.MouseScrollAxisHorizontal,
					value,
				)
			}
		}
		return 0

	case procs.WM_LBUTTONDOWN:
		procs.SetCapture(window)
		w.cursorCaptureCount++

		if cb := w.mouseInputCb.Load(); cb != nil {
			if cb := (*cb); cb != nil {
				cb(
					events.ButtonStatePressed,
					events.MouseButtonLeft,
				)
			}
		}
		return 0

	case procs.WM_LBUTTONUP:
		w.cursorCaptureCount = mathx.Max(0, w.cursorCaptureCount-1)
		if w.cursorCaptureCount == 0 {
			procs.ReleaseCapture()
		}

		if cb := w.mouseInputCb.Load(); cb != nil {
			if cb := (*cb); cb != nil {
				cb(
					events.ButtonStateReleased,
					events.MouseButtonLeft,
				)
			}
		}
		return 0

	case procs.WM_RBUTTONDOWN:
		procs.SetCapture(window)
		w.cursorCaptureCount++

		if cb := w.mouseInputCb.Load(); cb != nil {
			if cb := (*cb); cb != nil {
				cb(
					events.ButtonStatePressed,
					events.MouseButtonRight,
				)
			}
		}
		return 0

	case procs.WM_RBUTTONUP:
		w.cursorCaptureCount = mathx.Max(0, w.cursorCaptureCount-1)
		if w.cursorCaptureCount == 0 {
			procs.ReleaseCapture()
		}

		if cb := w.mouseInputCb.Load(); cb != nil {
			if cb := (*cb); cb != nil {
				cb(
					events.ButtonStateReleased,
					events.MouseButtonRight,
				)
			}
		}
		return 0

	case procs.WM_MBUTTONDOWN:
		procs.SetCapture(window)
		w.cursorCaptureCount++

		if cb := w.mouseInputCb.Load(); cb != nil {
			if cb := (*cb); cb != nil {
				cb(
					events.ButtonStatePressed,
					events.MouseButtonMiddle,
				)
			}
		}
		return 0

	case procs.WM_MBUTTONUP:
		w.cursorCaptureCount = mathx.Max(0, w.cursorCaptureCount-1)
		if w.cursorCaptureCount == 0 {
			procs.ReleaseCapture()
		}

		if cb := w.mouseInputCb.Load(); cb != nil {
			if cb := (*cb); cb != nil {
				cb(
					events.ButtonStateReleased,
					events.MouseButtonMiddle,
				)
			}
		}
		return 0

	case procs.WM_XBUTTONDOWN:
		procs.SetCapture(window)
		w.cursorCaptureCount++

		if cb := w.mouseInputCb.Load(); cb != nil {
			if cb := (*cb); cb != nil {
				cb(
					events.ButtonStatePressed,
					events.MouseButton(loword(uint32(wparam))),
				)
			}
		}
		return 0

	case procs.WM_XBUTTONUP:
		w.cursorCaptureCount = mathx.Max(0, w.cursorCaptureCount-1)
		if w.cursorCaptureCount == 0 {
			procs.ReleaseCapture()
		}

		if cb := w.mouseInputCb.Load(); cb != nil {
			if cb := (*cb); cb != nil {
				cb(
					events.ButtonStateReleased,
					events.MouseButton(loword(uint32(wparam))),
				)
			}
		}
		return 0

	case procs.WM_CAPTURECHANGED:
		if lparam != window {
			w.cursorCaptureCount = 0
		}
		return 0

	case procs.WM_SETFOCUS:
		if cb := w.focusedCb.Load(); cb != nil {
			if cb := (*cb); cb != nil {
				cb()
			}
		}

		if m := getModifiersState(); m != 0 {
			w.modifiers = m

			if cb := w.modifiersChangedCb.Load(); cb != nil {
				if cb := (*cb); cb != nil {
					cb(m)
				}
			}
		}
		return 0

	case procs.WM_KILLFOCUS:
		if w.modifiers != 0 {
			w.modifiers = 0

			if cb := w.modifiersChangedCb.Load(); cb != nil {
				if cb := (*cb); cb != nil {
					cb(0)
				}
			}
		}

		if cb := w.unfocusedCb.Load(); cb != nil {
			if cb := (*cb); cb != nil {
				cb()
			}
		}
		return 0

	case procs.WM_KEYDOWN, procs.WM_SYSKEYDOWN:
		if msg == procs.WM_SYSKEYDOWN && wparam == procs.VK_F4 {
			return procs.DefWindowProcW(window, msg, wparam, lparam)
		}

		if cb := w.keyboardInputCb.Load(); cb != nil {
			if cb := (*cb); cb != nil {
				scancode := ((lparam >> 16) & 0xff)
				extended := (lparam & 0x01000000) != 0
				if extended {
					scancode |= 0xE000
				} else {
					scancode |= 0x0000
				}

				cb(
					events.ButtonStatePressed,
					events.ScanCode(scancode),
					mapVK(wparam, scancode, extended),
				)
			}
		}

		m := getModifiersState()
		if w.modifiers != m {
			w.modifiers = m

			if cb := w.modifiersChangedCb.Load(); cb != nil {
				if cb := (*cb); cb != nil {
					cb(m)
				}
			}
		}

		if cb := w.receivedCharacterCb.Load(); cb != nil {
			if cb := (*cb); cb != nil {
				// win32 doesn't send WM_CHAR message for delete key
				if wparam == procs.VK_DELETE {
					cb('\u007F')
				}
			}
		}
		return 0

	case procs.WM_KEYUP, procs.WM_SYSKEYUP:
		if cb := w.keyboardInputCb.Load(); cb != nil {
			if cb := (*cb); cb != nil {
				scancode := ((lparam >> 16) & 0xff)
				extended := (lparam & 0x01000000) != 0
				if extended {
					scancode |= 0xE000
				} else {
					scancode |= 0x0000
				}

				cb(
					events.ButtonStateReleased,
					events.ScanCode(scancode),
					mapVK(wparam, scancode, extended),
				)
			}
		}

		m := getModifiersState()
		if w.modifiers != m {
			w.modifiers = m

			if cb := w.modifiersChangedCb.Load(); cb != nil {
				if cb := (*cb); cb != nil {
					cb(m)
				}
			}
		}
		return 0
	case procs.WM_CHAR, procs.WM_SYSCHAR:
		// Most UTF16 without surrogates can simply be considered rune.
		ch := rune(wparam)
		// The surrogated UTF16 character is POSTed as two consecutive WM_CHARs: high surrogate and low surrogate.
		if isSurrogatedCharacter(ch) {
			if w.highSurrogated == 0 {
				w.highSurrogated = ch
				return 0
			}
			ch = surrogatedUtf16toRune(w.highSurrogated, ch)
			w.highSurrogated = 0
		}
		if cb := w.receivedCharacterCb.Load(); cb != nil {
			if cb := (*cb); cb != nil {
				cb(ch)
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
