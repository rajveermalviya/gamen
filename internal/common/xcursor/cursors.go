//go:build linux

package xcursor

import (
	"strconv"

	"github.com/rajveermalviya/gamen/cursors"
)

func ToXcursorName(icon cursors.Icon) []string {
	switch icon {
	case cursors.Default:
		return []string{"default"}
	case cursors.ContextMenu:
		return []string{"context-menu"}
	case cursors.Help:
		return []string{"help", "question_arrow"}
	case cursors.Pointer:
		return []string{"pointer", "hand"}
	case cursors.Progress:
		return []string{"progress", "left_ptr_watch"}
	case cursors.Wait:
		return []string{"wait", "watch"}
	case cursors.Cell:
		return []string{"cell", "crosshair"}
	case cursors.Crosshair:
		return []string{"crosshair", "cross"}
	case cursors.Text:
		return []string{"text", "xterm"}
	case cursors.VerticalText:
		return []string{"vertical-text", "xterm"}
	case cursors.Alias:
		return []string{"alias", "dnd-link"}
	case cursors.Copy:
		return []string{"copy", "dnd-copy"}
	case cursors.Move:
		return []string{"move", "dnd-move"}
	case cursors.NoDrop:
		return []string{"no-drop", "dnd-none"}
	case cursors.NotAllowed:
		return []string{"not-allowed", "crossed_circle"}
	case cursors.Grab:
		return []string{"grab", "hand2"}
	case cursors.Grabbing:
		return []string{"grabbing", "hand2"}
	case cursors.AllScroll:
		return []string{"all-scroll"}
	case cursors.ColResize:
		return []string{"col-resize", "h_double_arrow"}
	case cursors.RowResize:
		return []string{"row-resize", "v_double_arrow"}
	case cursors.NResize:
		return []string{"n-resize", "top_side"}
	case cursors.EResize:
		return []string{"e-resize", "right_side"}
	case cursors.SResize:
		return []string{"s-resize", "bottom_side"}
	case cursors.WResize:
		return []string{"w-resize", "left_side"}
	case cursors.NEResize:
		return []string{"ne-resize", "top_right_corner"}
	case cursors.NWResize:
		return []string{"nw-resize", "top_left_corner"}
	case cursors.SEResize:
		return []string{"se-resize", "bottom_right_corner"}
	case cursors.SWResize:
		return []string{"sw-resize", "bottom_left_corner"}
	case cursors.EWResize:
		return []string{"ew-resize", "h_double_arrow"}
	case cursors.NSResize:
		return []string{"ns-resize", "v_double_arrow"}
	case cursors.NESWResize:
		return []string{"nesw-resize", "fd_double_arrow"}
	case cursors.NWSEResize:
		return []string{"nwse-resize", "bd_double_arrow"}
	case cursors.ZoomIn:
		return []string{"zoom-in"}
	case cursors.ZoomOut:
		return []string{"zoom-out"}
	}

	panic("invalid cursor: " + strconv.FormatUint(uint64(icon), 10))
}
