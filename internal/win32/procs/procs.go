//go:build windows

package procs

import (
	"runtime"

	"golang.org/x/sys/windows"
)

var (
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")
	user32   = windows.NewLazySystemDLL("user32.dll")

	_GetModuleHandleW = kernel32.NewProc("GetModuleHandleW")

	_PeekMessageW              = user32.NewProc("PeekMessageW")
	_TranslateMessage          = user32.NewProc("TranslateMessage")
	_DispatchMessageW          = user32.NewProc("DispatchMessageW")
	_WaitMessage               = user32.NewProc("WaitMessage")
	_MsgWaitForMultipleObjects = user32.NewProc("MsgWaitForMultipleObjects")
	_PostMessageW              = user32.NewProc("PostMessageW")

	_RegisterClassExW         = user32.NewProc("RegisterClassExW")
	_CreateWindowExW          = user32.NewProc("CreateWindowExW")
	_DestroyWindow            = user32.NewProc("DestroyWindow")
	_DefWindowProcW           = user32.NewProc("DefWindowProcW")
	_GetWindowLongW           = user32.NewProc("GetWindowLongW")
	_SetWindowLongW           = user32.NewProc("SetWindowLongW")
	_GetWindowLongPtrW        = user32.NewProc("GetWindowLongPtrW")
	_SetWindowLongPtrW        = user32.NewProc("SetWindowLongPtrW")
	_LoadCursorW              = user32.NewProc("LoadCursorW")
	_SetCursor                = user32.NewProc("SetCursor")
	_ShowCursor               = user32.NewProc("ShowCursor")
	_TrackMouseEvent          = user32.NewProc("TrackMouseEvent")
	_SetCapture               = user32.NewProc("SetCapture")
	_ReleaseCapture           = user32.NewProc("ReleaseCapture")
	_GetKeyboardLayout        = user32.NewProc("GetKeyboardLayout")
	_ToUnicodeEx              = user32.NewProc("ToUnicodeEx")
	_MapVirtualKeyW           = user32.NewProc("MapVirtualKeyW")
	_GetKeyboardState         = user32.NewProc("GetKeyboardState")
	_GetMenu                  = user32.NewProc("GetMenu")
	_GetDpiForWindow          = user32.NewProc("GetDpiForWindow")
	_GetClientRect            = user32.NewProc("GetClientRect")
	_AdjustWindowRectExForDpi = user32.NewProc("AdjustWindowRectExForDpi")
	_AdjustWindowRectEx       = user32.NewProc("AdjustWindowRectEx")
	_SetWindowPos             = user32.NewProc("SetWindowPos")
	_InvalidateRgn            = user32.NewProc("InvalidateRgn")
	_ShowWindow               = user32.NewProc("ShowWindow")
	_GetCursorPos             = user32.NewProc("GetCursorPos")
	_SetWindowTextW           = user32.NewProc("SetWindowTextW")
	_GetWindowPlacement       = user32.NewProc("GetWindowPlacement")
	_SetWindowPlacement       = user32.NewProc("SetWindowPlacement")
	_MonitorFromWindow        = user32.NewProc("MonitorFromWindow")
	_GetMonitorInfoW          = user32.NewProc("GetMonitorInfoW")
)

func GetModuleHandleW() (r uintptr) {
	r, _, _ = _GetModuleHandleW.Call(0)
	return
}

type POINT struct {
	X, Y int32
}

type POINTS struct {
	X, Y int16
}

type MSG struct {
	HWND    uintptr
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	PT      POINT
}

const (
	PM_REMOVE = 1
)

func PeekMessageW(lpMsg, hWnd, wMsgFilterMin, wMsgFilterMax, wRemoveMsg uintptr) bool {
	r, _, _ := _PeekMessageW.Call(lpMsg, hWnd, wMsgFilterMin, wMsgFilterMax, wRemoveMsg)
	return r != 0
}

func TranslateMessage(lpMsg uintptr) bool {
	r, _, _ := _TranslateMessage.Call(lpMsg)
	return r != 0
}

func DispatchMessageW(lpMsg uintptr) {
	_DispatchMessageW.Call(lpMsg)
}

func WaitMessage() bool {
	r, _, _ := _WaitMessage.Call()
	return r != 0
}

const (
	QS_ALLEVENTS = 1215
)

func MsgWaitForMultipleObjects(nCount, pHandles, fWaitAll, dwMilliseconds, dwWakeMask uintptr) uintptr {
	r, _, _ := _MsgWaitForMultipleObjects.Call(nCount, pHandles, fWaitAll, dwMilliseconds, dwWakeMask)
	return r
}

func PostMessageW(hWnd, Msg, wParam, lParam uintptr) bool {
	r, _, _ := _PostMessageW.Call(hWnd, Msg, wParam, lParam)
	return r != 0
}

const (
	CS_VREDRAW = 1
	CS_HREDRAW = 2
)

type WNDCLASSEXW struct {
	CbSize        uint32
	Style         uint32
	LpfnWndProc   uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     uintptr
	HIcon         uintptr
	HCursor       uintptr
	HbrBackground uintptr
	LpszMenuName  *uint16
	LpszClassName *uint16
	HIconSm       uintptr
}

func RegisterClassExW(class uintptr) (r uintptr) {
	r, _, _ = _RegisterClassExW.Call(class)
	return
}

const (
	WS_OVERLAPPED       = 0
	WS_SIZEBOX          = 262144
	WS_MAXIMIZEBOX      = 65536
	WS_CAPTION          = 12582912
	WS_MINIMIZEBOX      = 131072
	WS_BORDER           = 8388608
	WS_VISIBLE          = 268435456
	WS_CLIPSIBLINGS     = 67108864
	WS_CLIPCHILDREN     = 33554432
	WS_SYSMENU          = 524288
	WS_OVERLAPPEDWINDOW = 13565952
)

const (
	WS_EX_LEFT       = 0
	WS_EX_WINDOWEDGE = 256
	WS_EX_APPWINDOW  = 262144
)

const (
	CW_USEDEFAULT = ^uintptr(0) - 2147483647 // -2147483648
)

func CreateWindowExW(
	dwExStyle,
	lpClassName,
	lpWindowName,
	dwStyle,
	X,
	Y,
	nWidth,
	nHeight,
	hWndParent,
	hMenu,
	hInstance,
	lpParam uintptr,
) (r uintptr) {
	r, _, _ = _CreateWindowExW.Call(
		dwExStyle,
		lpClassName,
		lpWindowName,
		dwStyle,
		X,
		Y,
		nWidth,
		nHeight,
		hWndParent,
		hMenu,
		hInstance,
		lpParam,
	)
	return
}

type CREATESTRUCTW struct {
	LpCreateParams uintptr
	HInstance      uintptr
	HMenu          uintptr
	HwndParent     uintptr
	Cy             int32
	Cx             int32
	Y              int32
	X              int32
	Style          int32
	LpszName       *uint16
	LpszClass      *uint16
	DwExStyle      uint32
}

const (
	WM_CREATE         = 1
	WM_NCCREATE       = 129
	WM_CLOSE          = 16
	WM_SIZE           = 5
	WM_MOUSEMOVE      = 512
	WM_MOUSELEAVE     = 675
	WM_MOUSEWHEEL     = 522
	WM_MOUSEHWHEEL    = 526
	WM_NCLBUTTONDOWN  = 161
	WM_LBUTTONDOWN    = 513
	WM_LBUTTONUP      = 514
	WM_RBUTTONDOWN    = 516
	WM_RBUTTONUP      = 517
	WM_MBUTTONDOWN    = 519
	WM_MBUTTONUP      = 520
	WM_XBUTTONDOWN    = 523
	WM_XBUTTONUP      = 524
	WM_CAPTURECHANGED = 533
	WM_SETFOCUS       = 7
	WM_KILLFOCUS      = 8
	WM_SETCURSOR      = 32
	WM_GETMINMAXINFO  = 36

	WM_KEYDOWN    = 256
	WM_SYSKEYDOWN = 260
	WM_KEYUP      = 257
	WM_SYSKEYUP   = 261

	WM_CHAR    = 258
	WM_SYSCHAR = 262
)

const (
	HTCAPTION   = 2
	WHEEL_DELTA = 120
)

func DefWindowProcW(hWnd, Msg, wParam, lParam uintptr) (r uintptr) {
	r, _, _ = _DefWindowProcW.Call(hWnd, Msg, wParam, lParam)
	return r
}

const GWL_STYLE = ^uintptr(0) - 15    // -16
const GWL_EXSTYLE = ^uintptr(0) - 19  // -20
const GWL_USERDATA = ^uintptr(0) - 20 // -21

func GetWindowLong(hWnd, nIndex uintptr) uintptr {
	switch runtime.GOARCH {
	case "amd64", "arm64":
		r, _, _ := _GetWindowLongPtrW.Call(hWnd, nIndex)
		return r

	case "386", "arm":
		r, _, _ := _GetWindowLongW.Call(hWnd, nIndex)
		return r
	}

	panic("unsupported GOARCH: " + runtime.GOARCH)
}

func SetWindowLong(hWnd, nIndex, dwNewLong uintptr) uintptr {
	switch runtime.GOARCH {
	case "amd64", "arm64":
		r, _, _ := _SetWindowLongPtrW.Call(hWnd, nIndex, dwNewLong)
		return r

	case "386", "arm":
		r, _, _ := _SetWindowLongW.Call(hWnd, nIndex, dwNewLong)
		return r
	}

	panic("unsupported GOARCH: " + runtime.GOARCH)
}

const (
	HTCLIENT        = 1
	IDC_APPSTARTING = 32650
	IDC_ARROW       = 32512
	IDC_CROSS       = 32515
	IDC_HAND        = 32649
	IDC_HELP        = 32651
	IDC_IBEAM       = 32513
	IDC_ICON        = 32641
	IDC_NO          = 32648
	IDC_SIZE        = 32640
	IDC_SIZEALL     = 32646
	IDC_SIZENESW    = 32643
	IDC_SIZENS      = 32645
	IDC_SIZENWSE    = 32642
	IDC_SIZEWE      = 32644
	IDC_UPARROW     = 32516
	IDC_WAIT        = 32514
)

func LoadCursorW(hInstance, lpCursorName uintptr) (r uintptr) {
	r, _, _ = _LoadCursorW.Call(hInstance, lpCursorName)
	return
}

func SetCursor(hcursor uintptr) (r uintptr) {
	r, _, _ = _SetCursor.Call(hcursor)
	return
}

func ShowCursor(bShow uintptr) uintptr {
	r, _, _ := _ShowCursor.Call(bShow)
	return r
}

const (
	TME_LEAVE     = 2
	HOVER_DEFAULT = 4294967295
)

type TRACKMOUSEEVENT struct {
	CbSize      uint32
	DwFlags     uint32
	HwndTrack   uintptr
	DwHoverTime uint32
}

func TrackMouseEvent(lpEventTrack uintptr) bool {
	r, _, _ := _TrackMouseEvent.Call(lpEventTrack)
	return r != 0
}

func SetCapture(hwnd uintptr) (r uintptr) {
	r, _, _ = _SetCapture.Call(hwnd)
	return
}

func ReleaseCapture() bool {
	r, _, _ := _ReleaseCapture.Call()
	return r != 0
}

const (
	VK_0                               = 48
	VK_1                               = 49
	VK_2                               = 50
	VK_3                               = 51
	VK_4                               = 52
	VK_5                               = 53
	VK_6                               = 54
	VK_7                               = 55
	VK_8                               = 56
	VK_9                               = 57
	VK_A                               = 65
	VK_B                               = 66
	VK_C                               = 67
	VK_D                               = 68
	VK_E                               = 69
	VK_F                               = 70
	VK_G                               = 71
	VK_H                               = 72
	VK_I                               = 73
	VK_J                               = 74
	VK_K                               = 75
	VK_L                               = 76
	VK_M                               = 77
	VK_N                               = 78
	VK_O                               = 79
	VK_P                               = 80
	VK_Q                               = 81
	VK_R                               = 82
	VK_S                               = 83
	VK_T                               = 84
	VK_U                               = 85
	VK_V                               = 86
	VK_W                               = 87
	VK_X                               = 88
	VK_Y                               = 89
	VK_Z                               = 90
	VK_LBUTTON                         = 1
	VK_RBUTTON                         = 2
	VK_CANCEL                          = 3
	VK_MBUTTON                         = 4
	VK_XBUTTON1                        = 5
	VK_XBUTTON2                        = 6
	VK_BACK                            = 8
	VK_TAB                             = 9
	VK_CLEAR                           = 12
	VK_RETURN                          = 13
	VK_SHIFT                           = 16
	VK_CONTROL                         = 17
	VK_MENU                            = 18
	VK_PAUSE                           = 19
	VK_CAPITAL                         = 20
	VK_KANA                            = 21
	VK_HANGEUL                         = 21
	VK_HANGUL                          = 21
	VK_IME_ON                          = 22
	VK_JUNJA                           = 23
	VK_FINAL                           = 24
	VK_HANJA                           = 25
	VK_KANJI                           = 25
	VK_IME_OFF                         = 26
	VK_ESCAPE                          = 27
	VK_CONVERT                         = 28
	VK_NONCONVERT                      = 29
	VK_ACCEPT                          = 30
	VK_MODECHANGE                      = 31
	VK_SPACE                           = 32
	VK_PRIOR                           = 33
	VK_NEXT                            = 34
	VK_END                             = 35
	VK_HOME                            = 36
	VK_LEFT                            = 37
	VK_UP                              = 38
	VK_RIGHT                           = 39
	VK_DOWN                            = 40
	VK_SELECT                          = 41
	VK_PRINT                           = 42
	VK_EXECUTE                         = 43
	VK_SNAPSHOT                        = 44
	VK_INSERT                          = 45
	VK_DELETE                          = 46
	VK_HELP                            = 47
	VK_LWIN                            = 91
	VK_RWIN                            = 92
	VK_APPS                            = 93
	VK_SLEEP                           = 95
	VK_NUMPAD0                         = 96
	VK_NUMPAD1                         = 97
	VK_NUMPAD2                         = 98
	VK_NUMPAD3                         = 99
	VK_NUMPAD4                         = 100
	VK_NUMPAD5                         = 101
	VK_NUMPAD6                         = 102
	VK_NUMPAD7                         = 103
	VK_NUMPAD8                         = 104
	VK_NUMPAD9                         = 105
	VK_MULTIPLY                        = 106
	VK_ADD                             = 107
	VK_SEPARATOR                       = 108
	VK_SUBTRACT                        = 109
	VK_DECIMAL                         = 110
	VK_DIVIDE                          = 111
	VK_F1                              = 112
	VK_F2                              = 113
	VK_F3                              = 114
	VK_F4                              = 115
	VK_F5                              = 116
	VK_F6                              = 117
	VK_F7                              = 118
	VK_F8                              = 119
	VK_F9                              = 120
	VK_F10                             = 121
	VK_F11                             = 122
	VK_F12                             = 123
	VK_F13                             = 124
	VK_F14                             = 125
	VK_F15                             = 126
	VK_F16                             = 127
	VK_F17                             = 128
	VK_F18                             = 129
	VK_F19                             = 130
	VK_F20                             = 131
	VK_F21                             = 132
	VK_F22                             = 133
	VK_F23                             = 134
	VK_F24                             = 135
	VK_NAVIGATION_VIEW                 = 136
	VK_NAVIGATION_MENU                 = 137
	VK_NAVIGATION_UP                   = 138
	VK_NAVIGATION_DOWN                 = 139
	VK_NAVIGATION_LEFT                 = 140
	VK_NAVIGATION_RIGHT                = 141
	VK_NAVIGATION_ACCEPT               = 142
	VK_NAVIGATION_CANCEL               = 143
	VK_NUMLOCK                         = 144
	VK_SCROLL                          = 145
	VK_OEM_NEC_EQUAL                   = 146
	VK_OEM_FJ_JISHO                    = 146
	VK_OEM_FJ_MASSHOU                  = 147
	VK_OEM_FJ_TOUROKU                  = 148
	VK_OEM_FJ_LOYA                     = 149
	VK_OEM_FJ_ROYA                     = 150
	VK_LSHIFT                          = 160
	VK_RSHIFT                          = 161
	VK_LCONTROL                        = 162
	VK_RCONTROL                        = 163
	VK_LMENU                           = 164
	VK_RMENU                           = 165
	VK_BROWSER_BACK                    = 166
	VK_BROWSER_FORWARD                 = 167
	VK_BROWSER_REFRESH                 = 168
	VK_BROWSER_STOP                    = 169
	VK_BROWSER_SEARCH                  = 170
	VK_BROWSER_FAVORITES               = 171
	VK_BROWSER_HOME                    = 172
	VK_VOLUME_MUTE                     = 173
	VK_VOLUME_DOWN                     = 174
	VK_VOLUME_UP                       = 175
	VK_MEDIA_NEXT_TRACK                = 176
	VK_MEDIA_PREV_TRACK                = 177
	VK_MEDIA_STOP                      = 178
	VK_MEDIA_PLAY_PAUSE                = 179
	VK_LAUNCH_MAIL                     = 180
	VK_LAUNCH_MEDIA_SELECT             = 181
	VK_LAUNCH_APP1                     = 182
	VK_LAUNCH_APP2                     = 183
	VK_OEM_1                           = 186
	VK_OEM_PLUS                        = 187
	VK_OEM_COMMA                       = 188
	VK_OEM_MINUS                       = 189
	VK_OEM_PERIOD                      = 190
	VK_OEM_2                           = 191
	VK_OEM_3                           = 192
	VK_GAMEPAD_A                       = 195
	VK_GAMEPAD_B                       = 196
	VK_GAMEPAD_X                       = 197
	VK_GAMEPAD_Y                       = 198
	VK_GAMEPAD_RIGHT_SHOULDER          = 199
	VK_GAMEPAD_LEFT_SHOULDER           = 200
	VK_GAMEPAD_LEFT_TRIGGER            = 201
	VK_GAMEPAD_RIGHT_TRIGGER           = 202
	VK_GAMEPAD_DPAD_UP                 = 203
	VK_GAMEPAD_DPAD_DOWN               = 204
	VK_GAMEPAD_DPAD_LEFT               = 205
	VK_GAMEPAD_DPAD_RIGHT              = 206
	VK_GAMEPAD_MENU                    = 207
	VK_GAMEPAD_VIEW                    = 208
	VK_GAMEPAD_LEFT_THUMBSTICK_BUTTON  = 209
	VK_GAMEPAD_RIGHT_THUMBSTICK_BUTTON = 210
	VK_GAMEPAD_LEFT_THUMBSTICK_UP      = 211
	VK_GAMEPAD_LEFT_THUMBSTICK_DOWN    = 212
	VK_GAMEPAD_LEFT_THUMBSTICK_RIGHT   = 213
	VK_GAMEPAD_LEFT_THUMBSTICK_LEFT    = 214
	VK_GAMEPAD_RIGHT_THUMBSTICK_UP     = 215
	VK_GAMEPAD_RIGHT_THUMBSTICK_DOWN   = 216
	VK_GAMEPAD_RIGHT_THUMBSTICK_RIGHT  = 217
	VK_GAMEPAD_RIGHT_THUMBSTICK_LEFT   = 218
	VK_OEM_4                           = 219
	VK_OEM_5                           = 220
	VK_OEM_6                           = 221
	VK_OEM_7                           = 222
	VK_OEM_8                           = 223
	VK_OEM_AX                          = 225
	VK_OEM_102                         = 226
	VK_ICO_HELP                        = 227
	VK_ICO_00                          = 228
	VK_PROCESSKEY                      = 229
	VK_ICO_CLEAR                       = 230
	VK_PACKET                          = 231
	VK_OEM_RESET                       = 233
	VK_OEM_JUMP                        = 234
	VK_OEM_PA1                         = 235
	VK_OEM_PA2                         = 236
	VK_OEM_PA3                         = 237
	VK_OEM_WSCTRL                      = 238
	VK_OEM_CUSEL                       = 239
	VK_OEM_ATTN                        = 240
	VK_OEM_FINISH                      = 241
	VK_OEM_COPY                        = 242
	VK_OEM_AUTO                        = 243
	VK_OEM_ENLW                        = 244
	VK_OEM_BACKTAB                     = 245
	VK_ATTN                            = 246
	VK_CRSEL                           = 247
	VK_EXSEL                           = 248
	VK_EREOF                           = 249
	VK_PLAY                            = 250
	VK_ZOOM                            = 251
	VK_NONAME                          = 252
	VK_PA1                             = 253
	VK_OEM_CLEAR                       = 254
	VK_ABNT_C1                         = 193
	VK_ABNT_C2                         = 194
	VK_DBE_ALPHANUMERIC                = 240
	VK_DBE_CODEINPUT                   = 250
	VK_DBE_DBCSCHAR                    = 244
	VK_DBE_DETERMINESTRING             = 252
	VK_DBE_ENTERDLGCONVERSIONMODE      = 253
	VK_DBE_ENTERIMECONFIGMODE          = 248
	VK_DBE_ENTERWORDREGISTERMODE       = 247
	VK_DBE_FLUSHSTRING                 = 249
	VK_DBE_HIRAGANA                    = 242
	VK_DBE_KATAKANA                    = 241
	VK_DBE_NOCODEINPUT                 = 251
	VK_DBE_NOROMAN                     = 246
	VK_DBE_ROMAN                       = 245
	VK_DBE_SBCSCHAR                    = 243
)

func GetKeyboardLayout(idThread uintptr) uintptr {
	r, _, _ := _GetKeyboardLayout.Call(idThread)
	return r
}

func ToUnicodeEx(wVirtKey, wScanCode, lpKeyState, pwszBuff, cchBuff, wFlags, dwhkl uintptr) uintptr {
	r, _, _ := _ToUnicodeEx.Call(wVirtKey, wScanCode, lpKeyState, pwszBuff, cchBuff, wFlags, dwhkl)
	return r
}

const MAPVK_VK_TO_CHAR = 2
const MAPVK_VSC_TO_VK_EX = 3

func MapVirtualKeyW(uCode, uMapType uintptr) uintptr {
	r, _, _ := _MapVirtualKeyW.Call(uCode, uMapType)
	return r
}

func GetKeyboardState(lpKeyState uintptr) bool {
	r, _, _ := _GetKeyboardState.Call(lpKeyState)
	return r != 0
}

func DestroyWindow(hwnd uintptr) bool {
	r, _, _ := _DestroyWindow.Call(hwnd)
	return r != 0
}

type RECT struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

func GetClientRect(hWnd, lpRect uintptr) bool {
	r, _, _ := _GetClientRect.Call(hWnd, lpRect)
	return r != 0
}

func AdjustWindowRectEx(hwnd, style, styleEx, rect uintptr) bool {
	menu, _, _ := _GetMenu.Call(hwnd)
	if menu != 0 {
		menu = 1
	}

	if _GetDpiForWindow.Find() == nil && _AdjustWindowRectExForDpi.Find() == nil {
		dpi, _, _ := _GetDpiForWindow.Call(hwnd)
		r, _, _ := _AdjustWindowRectExForDpi.Call(rect, style, menu, styleEx, dpi)
		return r != 0
	} else {
		r, _, _ := _AdjustWindowRectEx.Call(rect, style, menu, styleEx)
		return r != 0
	}
}

const SWP_ASYNCWINDOWPOS = 16384
const SWP_NOZORDER = 4
const SWP_NOREPOSITION = 512
const SWP_NOMOVE = 2
const SWP_NOACTIVATE = 16
const SWP_NOOWNERZORDER = 512
const SWP_FRAMECHANGED = 32
const SWP_NOSIZE = 1

const HWND_TOP = 0

func SetWindowPos(hWnd, hWndInsertAfter, X, Y, cx, cy, uFlags uintptr) bool {
	r, _, _ := _SetWindowPos.Call(hWnd, hWndInsertAfter, X, Y, cx, cy, uFlags)
	return r != 0
}

func InvalidateRgn(hwnd, hrng, berase uintptr) bool {
	r, _, _ := _InvalidateRgn.Call(hwnd, hrng, berase)
	return r != 0
}

type MINMAXINFO struct {
	PtReserved     POINT
	PtMaxSize      POINT
	PtMaxPosition  POINT
	PtMinTrackSize POINT
	PtMaxTrackSize POINT
}

const (
	SW_MAXIMIZE = 3
	SW_MINIMIZE = 6
	SW_RESTORE  = 9
)

const (
	SIZE_MAXIMIZED = 2
	SIZE_RESTORED  = 0
)

func ShowWindow(hWnd, nCmdShow uintptr) bool {
	r, _, _ := _ShowWindow.Call(hWnd, nCmdShow)
	return r != 0
}

func GetCursorPos(point uintptr) bool {
	r, _, _ := _GetCursorPos.Call(point)
	return r != 0
}

func SetWindowTextW(hwnd, lpstring uintptr) bool {
	r, _, _ := _SetWindowTextW.Call(hwnd, lpstring)
	return r != 0
}

type WINDOWPLACEMENT struct {
	length           uint32
	flags            uint32
	showCmd          uint32
	ptMinPosition    POINT
	ptMaxPosition    POINT
	rcNormalPosition RECT
}

func GetWindowPlacement(hWnd, lpwndpl uintptr) bool {
	r, _, _ := _GetWindowPlacement.Call(hWnd, lpwndpl)
	return r != 0
}
func SetWindowPlacement(hWnd, lpwndpl uintptr) bool {
	r, _, _ := _SetWindowPlacement.Call(hWnd, lpwndpl)
	return r != 0
}

func MonitorFromWindow(hwnd, dwFlags uintptr) uintptr {
	r, _, _ := _MonitorFromWindow.Call(hwnd, dwFlags)
	return r
}

type MONITORINFO struct {
	CbSize    uint32
	RcMonitor RECT
	RcWork    RECT
	DwFlags   uint32
}

const (
	MONITOR_DEFAULTTOPRIMARY = 1
)

func GetMonitorInfoW(hMonitor, lpmi uintptr) bool {
	r, _, _ := _GetMonitorInfoW.Call(hMonitor, lpmi)
	return r != 0
}
