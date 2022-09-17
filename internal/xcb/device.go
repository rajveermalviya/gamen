//go:build linux && !android

package xcb

/*

#include <stdlib.h>
#include <xcb/xcb.h>
#include <xcb/xinput.h>

*/
import "C"

import (
	"errors"
	"unsafe"
)

type scrollAxis struct {
	index     C.uint16_t
	increment float64
	position  float64
}

type scrollingDevice struct {
	horizontalScroll scrollAxis
	verticalScroll   scrollAxis
}

func (d *Display) xiSetupScrollingDevices(id C.xcb_input_device_id_t) error {
	reply := d.l.xcb_input_xi_query_device_reply(d.xcbConn, id)
	if reply == nil {
		return errors.New("failed to setup xinput scrolling devices")
	}
	defer C.free(unsafe.Pointer(reply))

	for it := d.l.xcb_input_xi_query_device_infos_iterator(reply); it.rem > 0; d.l.xcb_input_xi_device_info_next(&it) {
		deviceInfo := it.data

		if deviceInfo._type != C.XCB_INPUT_DEVICE_TYPE_SLAVE_POINTER {
			continue
		}

		var dev scrollingDevice

		for it := d.l.xcb_input_xi_device_info_classes_iterator(deviceInfo); it.rem != 0; d.l.xcb_input_device_class_next(&it) {
			classInfo := it.data
			switch classInfo._type {
			case C.XCB_INPUT_DEVICE_CLASS_TYPE_SCROLL:
				scrollClassInfo := (*C.xcb_input_scroll_class_t)(unsafe.Pointer(classInfo))

				switch scrollClassInfo.scroll_type {
				case C.XCB_INPUT_SCROLL_TYPE_VERTICAL:
					dev.verticalScroll.index = scrollClassInfo.number
					dev.verticalScroll.increment = fixed3232ToFloat64(scrollClassInfo.increment)

				case C.XCB_INPUT_SCROLL_TYPE_HORIZONTAL:
					dev.horizontalScroll.index = scrollClassInfo.number
					dev.horizontalScroll.increment = fixed3232ToFloat64(scrollClassInfo.increment)
				}

			case C.XCB_INPUT_DEVICE_CLASS_TYPE_VALUATOR:
				vci := (*C.xcb_input_valuator_class_t)(unsafe.Pointer(classInfo))

				switch vci.label {
				case d.relHorizScroll, d.relHorizWheel:
					dev.horizontalScroll.position = fixed3232ToFloat64(vci.value)
				case d.relVertScroll, d.relVertWheel:
					dev.verticalScroll.position = fixed3232ToFloat64(vci.value)
				}
			}
		}

		d.scrollingDevices[deviceInfo.deviceid] = dev
	}

	return nil
}

func (d *Display) resetScrollPosition(id C.xcb_input_device_id_t) {
	reply := d.l.xcb_input_xi_query_device_reply(d.xcbConn, id)
	if reply == nil {
		return
	}
	defer C.free(unsafe.Pointer(reply))

	for it := d.l.xcb_input_xi_query_device_infos_iterator(reply); it.rem > 0; d.l.xcb_input_xi_device_info_next(&it) {
		deviceInfo := it.data

		if deviceInfo._type != C.XCB_INPUT_DEVICE_TYPE_SLAVE_POINTER {
			continue
		}

		dev, ok := d.scrollingDevices[deviceInfo.deviceid]
		if !ok {
			continue
		}

		for it := d.l.xcb_input_xi_device_info_classes_iterator(deviceInfo); it.rem != 0; d.l.xcb_input_device_class_next(&it) {
			classInfo := it.data
			if classInfo._type != C.XCB_INPUT_DEVICE_CLASS_TYPE_VALUATOR {
				continue
			}

			vci := (*C.xcb_input_valuator_class_t)(unsafe.Pointer(classInfo))
			switch vci.label {
			case d.relHorizScroll, d.relHorizWheel:
				dev.horizontalScroll.position = fixed3232ToFloat64(vci.value)

			case d.relVertScroll, d.relVertWheel:
				dev.verticalScroll.position = fixed3232ToFloat64(vci.value)
			}
		}

		d.scrollingDevices[deviceInfo.deviceid] = dev
	}
}
