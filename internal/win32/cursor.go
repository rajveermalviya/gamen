//go:build windows

package win32

import (
	"github.com/rajveermalviya/gamen/cursors"
	"github.com/rajveermalviya/gamen/internal/win32/procs"
)

// TODO: we should embed missing cursors using https://github.com/tc-hib/go-winres
//
// firefox: https://github.com/mozilla/gecko-dev/tree/master/widget/windows/res
// chromium: https://chromium.googlesource.com/chromium/src.git/+/refs/heads/main/ui/resources/cursors/

func toWin32Cursor(icon cursors.Icon) uintptr {
	switch icon {
	case cursors.Default, cursors.Pointer:
		return procs.IDC_ARROW

	case cursors.Help:
		return procs.IDC_HELP

	case cursors.Progress:
		return procs.IDC_APPSTARTING

	case cursors.Wait:
		return procs.IDC_WAIT

	case cursors.Crosshair:
		return procs.IDC_CROSS

	case cursors.Text, cursors.VerticalText:
		return procs.IDC_IBEAM

	case cursors.Move:
		return procs.IDC_SIZEALL

	case cursors.NoDrop, cursors.NotAllowed:
		return procs.IDC_NO

	case cursors.AllScroll:
		return procs.IDC_SIZEALL

	case cursors.ColResize,
		cursors.EResize,
		cursors.WResize,
		cursors.EWResize:
		return procs.IDC_SIZEWE

	case cursors.RowResize,
		cursors.NResize,
		cursors.SResize,
		cursors.NSResize:
		return procs.IDC_SIZENS

	case cursors.NEResize,
		cursors.SWResize,
		cursors.NESWResize:
		return procs.IDC_SIZENESW

	case cursors.NWResize,
		cursors.SEResize,
		cursors.NWSEResize:
		return procs.IDC_SIZENWSE

	case cursors.ContextMenu,
		cursors.Cell,
		cursors.Alias,
		cursors.Copy,
		cursors.Grab,
		cursors.Grabbing,
		cursors.ZoomIn,
		cursors.ZoomOut:
		return procs.IDC_ARROW

	default:
		return procs.IDC_ARROW
	}
}
