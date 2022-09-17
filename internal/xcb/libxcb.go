//go:build linux && !android

package xcb

/*

#cgo CFLAGS: -I${SRCDIR}/include
#cgo LDFLAGS: -ldl

#include <dlfcn.h>
#include <X11/Xlib-xcb.h>
#include <xcb/randr.h>
#include <xcb/xinput.h>
#include <xcb/xcb_icccm.h>
#include <xcb/xcb_image.h>
#include <X11/Xcursor/Xcursor.h>
#include <xcb/xkb.h>

// libX11
Status gamen_XInitThreads(void *fp) {
	typedef Status (*proc_XInitThreads)(void);
	return ((proc_XInitThreads)fp)();
}
Display *gamen_XOpenDisplay(void *fp, const char* display_name) {
	typedef Display *(*proc_XOpenDisplay)(const char* display_name);
	return ((proc_XOpenDisplay)fp)(display_name);
}
int gamen_XCloseDisplay(void *fp, Display* dpy) {
	typedef int (*proc_XCloseDisplay)(Display* dpy);
	return ((proc_XCloseDisplay)fp)(dpy);
}

// libX11-xcb
xcb_connection_t *gamen_XGetXCBConnection(void *fp, Display *dpy) {
	typedef xcb_connection_t *(*proc_XGetXCBConnection)(Display *dpy);
	return ((proc_XGetXCBConnection)fp)(dpy);
}
void gamen_XSetEventQueueOwner(void *fp, Display *dpy, enum XEventQueueOwner owner) {
	typedef void (*proc_XSetEventQueueOwner)(Display *dpy, enum XEventQueueOwner owner);
	((proc_XSetEventQueueOwner)fp)(dpy, owner);
}

// libXcursor
Cursor gamen_XcursorLibraryLoadCursor(void *fp, Display *dpy, const char *file) {
	typedef Cursor (*proc_XcursorLibraryLoadCursor)(Display *dpy, const char *file);
	return ((proc_XcursorLibraryLoadCursor)fp)(dpy, file);
}

// libxcb
const struct xcb_setup_t *gamen_xcb_get_setup(void *fp, xcb_connection_t *c) {
	typedef const struct xcb_setup_t *(*proc_xcb_get_setup)(xcb_connection_t *c);
	return ((proc_xcb_get_setup)fp)(c);
}
xcb_screen_iterator_t gamen_xcb_setup_roots_iterator(void *fp, const xcb_setup_t *R) {
	typedef xcb_screen_iterator_t (*proc_xcb_setup_roots_iterator)(const xcb_setup_t *R);
	return ((proc_xcb_setup_roots_iterator)fp)(R);
}
void gamen_xcb_screen_next(void *fp, xcb_screen_iterator_t *i) {
	typedef void (*proc_xcb_screen_next)(xcb_screen_iterator_t *i);
	((proc_xcb_screen_next)fp)(i);
}
const struct xcb_query_extension_reply_t *gamen_xcb_get_extension_data(void *fp, xcb_connection_t *c, xcb_extension_t *ext) {
	typedef const struct xcb_query_extension_reply_t *(*proc_xcb_get_extension_data)(xcb_connection_t *c, xcb_extension_t *ext);
	return ((proc_xcb_get_extension_data)fp)(c, ext);
}
xcb_intern_atom_cookie_t gamen_xcb_intern_atom(void *fp, xcb_connection_t *c, uint8_t only_if_exists, uint16_t name_len, const char *name) {
	typedef xcb_intern_atom_cookie_t (*proc_xcb_intern_atom)(xcb_connection_t *c, uint8_t only_if_exists, uint16_t name_len, const char *name);
	return ((proc_xcb_intern_atom)fp)(c, only_if_exists, name_len, name);
}
xcb_intern_atom_reply_t *gamen_xcb_intern_atom_reply(void *fp, xcb_connection_t *c, xcb_intern_atom_cookie_t cookie, xcb_generic_error_t **e) {
	typedef xcb_intern_atom_reply_t *(*proc_xcb_intern_atom_reply)(xcb_connection_t *c, xcb_intern_atom_cookie_t cookie, xcb_generic_error_t **e);
	return ((proc_xcb_intern_atom_reply)fp)(c, cookie, e);
}
xcb_generic_event_t *gamen_xcb_poll_for_event(void *fp, xcb_connection_t *c) {
	typedef xcb_generic_event_t *(*proc_xcb_poll_for_event)(xcb_connection_t *c);
	return ((proc_xcb_poll_for_event)fp)(c);
}
int gamen_xcb_get_file_descriptor(void *fp, xcb_connection_t *c) {
	typedef int (*proc_xcb_get_file_descriptor)(xcb_connection_t *c);
	return ((proc_xcb_get_file_descriptor)fp)(c);
}
uint32_t gamen_xcb_generate_id(void *fp, xcb_connection_t *c) {
	typedef uint32_t (*proc_xcb_generate_id)(xcb_connection_t *c);
	return ((proc_xcb_generate_id)fp)(c);
}
xcb_void_cookie_t gamen_xcb_create_window_checked(void *fp, xcb_connection_t *c, uint8_t depth, xcb_window_t wid, xcb_window_t parent, int16_t x, int16_t y, uint16_t width, uint16_t height, uint16_t border_width, uint16_t _class, xcb_visualid_t visual, uint32_t value_mask, const void *value_list) {
	typedef xcb_void_cookie_t (*proc_gamen_xcb_create_window_checked)(xcb_connection_t *c, uint8_t depth, xcb_window_t wid, xcb_window_t parent, int16_t x, int16_t y, uint16_t width, uint16_t height, uint16_t border_width, uint16_t _class, xcb_visualid_t visual, uint32_t value_mask, const void *value_list);
	return ((proc_gamen_xcb_create_window_checked)fp)(c, depth, wid, parent, x, y, width, height, border_width, _class, visual, value_mask, value_list);
}
xcb_generic_error_t *gamen_xcb_request_check(void *fp, xcb_connection_t *c, xcb_void_cookie_t cookie) {
	typedef xcb_generic_error_t *(*proc_xcb_request_check)(xcb_connection_t *c, xcb_void_cookie_t cookie);
	return ((proc_xcb_request_check)fp)(c, cookie);
}
xcb_void_cookie_t gamen_xcb_change_property(void *fp, xcb_connection_t *c, uint8_t mode, xcb_window_t window, xcb_atom_t property, xcb_atom_t type, uint8_t format, uint32_t data_len, const void *data) {
	typedef xcb_void_cookie_t (*proc_xcb_change_property)(xcb_connection_t *c, uint8_t mode, xcb_window_t window, xcb_atom_t property, xcb_atom_t type, uint8_t format, uint32_t data_len, const void *data);
	return ((proc_xcb_change_property)fp)(c, mode, window, property, type, format, data_len, data);
}
xcb_void_cookie_t gamen_xcb_map_window_checked(void *fp, xcb_connection_t *c, xcb_window_t window) {
	typedef xcb_void_cookie_t (*proc_xcb_map_window_checked)(xcb_connection_t *c, xcb_window_t window);
	return ((proc_xcb_map_window_checked)fp)(c, window);
}
xcb_void_cookie_t gamen_xcb_destroy_window(void *fp, xcb_connection_t *c, xcb_window_t window) {
	typedef xcb_void_cookie_t (*proc_xcb_destroy_window)(xcb_connection_t *c, xcb_window_t window);
	return ((proc_xcb_destroy_window)fp)(c, window);
}
int gamen_xcb_flush(void *fp, xcb_connection_t *c) {
	typedef int (*proc_xcb_flush)(xcb_connection_t *c);
	return ((proc_xcb_flush)fp)(c);
}
xcb_get_geometry_reply_t *gamen_xcb_get_geometry_reply(void *fp1, void *fp2, xcb_connection_t *c, xcb_drawable_t drawable) {
	typedef xcb_get_geometry_cookie_t (*proc_xcb_get_geometry_unchecked)(xcb_connection_t *c, xcb_drawable_t drawable);
	typedef xcb_get_geometry_reply_t *(*proc_xcb_get_geometry_reply)(xcb_connection_t *c, xcb_get_geometry_cookie_t cookie, xcb_generic_error_t **e);
	return ((proc_xcb_get_geometry_reply)fp2)(c, ((proc_xcb_get_geometry_unchecked)fp1)(c, drawable), NULL);
}
xcb_void_cookie_t gamen_xcb_configure_window(void *fp, xcb_connection_t *c, xcb_window_t window, uint16_t value_mask, const void *value_list) {
	typedef xcb_void_cookie_t (*proc_xcb_configure_window)(xcb_connection_t *c, xcb_window_t window, uint16_t value_mask, const void *value_list);
	return ((proc_xcb_configure_window)fp)(c, window, value_mask, value_list);
}
xcb_get_property_reply_t *gamen_xcb_get_property_reply(void *fp1, void *fp2, xcb_connection_t *c, uint8_t _delete, xcb_window_t window, xcb_atom_t property, xcb_atom_t type, uint32_t long_offset, uint32_t long_length) {
	typedef xcb_get_property_cookie_t (*proc_xcb_get_property_unchecked)(xcb_connection_t *c, uint8_t _delete, xcb_window_t window, xcb_atom_t property, xcb_atom_t type, uint32_t long_offset, uint32_t long_length);
	typedef xcb_get_property_reply_t *(*proc_xcb_get_property_reply)(xcb_connection_t *c, xcb_get_property_cookie_t cookie, xcb_generic_error_t **e);
	return ((proc_xcb_get_property_reply)fp2)(c, ((proc_xcb_get_property_unchecked)fp1)(c, _delete, window, property, type, long_offset, long_length), NULL);
}
void *gamen_xcb_get_property_value(void *fp, const xcb_get_property_reply_t *R) {
	typedef void *(*proc_xcb_get_property_value)(const xcb_get_property_reply_t *R);
	return ((proc_xcb_get_property_value)fp)(R);
}
int gamen_xcb_get_property_value_length(void *fp, const xcb_get_property_reply_t *R) {
	typedef int (*proc_xcb_get_property_value_length)(const xcb_get_property_reply_t *R);
	return ((proc_xcb_get_property_value_length)fp)(R);
}
xcb_void_cookie_t gamen_xcb_send_event(void *fp, xcb_connection_t *c, uint8_t propagate, xcb_window_t destination, uint32_t event_mask, const char *event) {
	typedef xcb_void_cookie_t (*proc_xcb_send_event)(xcb_connection_t *c, uint8_t propagate, xcb_window_t destination, uint32_t event_mask, const char *event);
	return ((proc_xcb_send_event)fp)(c, propagate, destination, event_mask, event);
}
xcb_void_cookie_t gamen_xcb_change_window_attributes(void *fp, xcb_connection_t *c, xcb_window_t window, uint32_t value_mask, const void *value_list) {
	typedef xcb_void_cookie_t (*proc_xcb_change_window_attributes)(xcb_connection_t *c, xcb_window_t window, uint32_t value_mask, const void *value_list);
	return ((proc_xcb_change_window_attributes)fp)(c, window, value_mask, value_list);
}
xcb_translate_coordinates_reply_t *gamen_xcb_translate_coordinates_reply(void *fp1, void *fp2, xcb_connection_t *c, xcb_window_t src_window, xcb_window_t dst_window, int16_t src_x, int16_t src_y) {
	typedef xcb_translate_coordinates_cookie_t (*proc_xcb_translate_coordinates_unchecked)(xcb_connection_t *c, xcb_window_t src_window, xcb_window_t dst_window, int16_t src_x, int16_t src_y);
	typedef xcb_translate_coordinates_reply_t *(*proc_xcb_translate_coordinates_reply)(xcb_connection_t *c, xcb_translate_coordinates_cookie_t cookie, xcb_generic_error_t **e);
	return ((proc_xcb_translate_coordinates_reply)fp2)(c, ((proc_xcb_translate_coordinates_unchecked)fp1)(c, src_window, dst_window, src_x, src_y), NULL);
}
xcb_void_cookie_t gamen_xcb_ungrab_pointer(void *fp, xcb_connection_t *c, xcb_timestamp_t time) {
	typedef xcb_void_cookie_t (*proc_xcb_ungrab_pointer)(xcb_connection_t *c, xcb_timestamp_t time);
	return ((proc_xcb_ungrab_pointer)fp)(c, time);
}
xcb_void_cookie_t gamen_xcb_free_pixmap(void *fp, xcb_connection_t *c, xcb_pixmap_t pixmap) {
	typedef xcb_void_cookie_t (*proc_xcb_free_pixmap)(xcb_connection_t *c, xcb_pixmap_t pixmap);
	return ((proc_xcb_free_pixmap)fp)(c, pixmap);
}
xcb_void_cookie_t gamen_xcb_create_cursor(void *fp, xcb_connection_t *c, xcb_cursor_t cid, xcb_pixmap_t source, xcb_pixmap_t mask, uint16_t fore_red, uint16_t fore_green, uint16_t fore_blue, uint16_t back_red, uint16_t back_green, uint16_t back_blue, uint16_t x, uint16_t y) {
	typedef xcb_void_cookie_t (*proc_xcb_create_cursor)(xcb_connection_t *c, xcb_cursor_t cid, xcb_pixmap_t source, xcb_pixmap_t mask, uint16_t fore_red, uint16_t fore_green, uint16_t fore_blue, uint16_t back_red, uint16_t back_green, uint16_t back_blue, uint16_t x, uint16_t y);
	return ((proc_xcb_create_cursor)fp)(c, cid, source, mask, fore_red, fore_green, fore_blue, back_red, back_green, back_blue, x, y);
}
xcb_void_cookie_t gamen_xcb_free_cursor(void *fp, xcb_connection_t *c, xcb_cursor_t cursor) {
	typedef xcb_void_cookie_t (*proc_xcb_free_cursor)(xcb_connection_t *c, xcb_cursor_t cursor);
	return ((proc_xcb_free_cursor)fp)(c, cursor);
}

// libxcb-randr
xcb_randr_query_version_reply_t *gamen_xcb_randr_query_version_reply(void *fp1, void *fp2, xcb_connection_t *c, uint32_t major_version, uint32_t minor_version) {
	typedef xcb_randr_query_version_cookie_t (*proc_xcb_randr_query_version_unchecked)(xcb_connection_t *c, uint32_t major_version, uint32_t minor_version);
	typedef xcb_randr_query_version_reply_t *(*proc_xcb_randr_query_version_reply)(xcb_connection_t *c, xcb_randr_query_version_cookie_t cookie, xcb_generic_error_t **e);
	return ((proc_xcb_randr_query_version_reply)fp2)(c, ((proc_xcb_randr_query_version_unchecked)fp1)(c, major_version, minor_version), NULL);
}
xcb_void_cookie_t gamen_xcb_randr_select_input(void *fp, xcb_connection_t *c, xcb_window_t window, uint16_t enable) {
	typedef xcb_void_cookie_t (*proc_xcb_randr_select_input)(xcb_connection_t *c, xcb_window_t window, uint16_t enable);
	return ((proc_xcb_randr_select_input)fp)(c, window, enable);
}

// libxcb-xinput
xcb_input_xi_query_version_reply_t *gamen_xcb_input_xi_query_version_reply(void *fp1, void *fp2, xcb_connection_t *c, uint16_t major_version, uint16_t minor_version) {
	typedef xcb_input_xi_query_version_cookie_t (*proc_xcb_input_xi_query_version_unchecked)(xcb_connection_t *c, uint16_t major_version, uint16_t minor_version);
	typedef xcb_input_xi_query_version_reply_t *(*proc_xcb_input_xi_query_version_reply)(xcb_connection_t *c, xcb_input_xi_query_version_cookie_t cookie, xcb_generic_error_t **e);
	return ((proc_xcb_input_xi_query_version_reply)fp2)(c, ((proc_xcb_input_xi_query_version_unchecked)fp1)(c, major_version, minor_version), NULL);
}
xcb_input_xi_query_device_reply_t *gamen_xcb_input_xi_query_device_reply(void *fp1, void *fp2, xcb_connection_t *c, xcb_input_device_id_t deviceid) {
	typedef xcb_input_xi_query_device_cookie_t (*proc_xcb_input_xi_query_device_unchecked)(xcb_connection_t *c, xcb_input_device_id_t deviceid);
	typedef xcb_input_xi_query_device_reply_t *(*proc_xcb_input_xi_query_device_reply)(xcb_connection_t *c, xcb_input_xi_query_device_cookie_t cookie, xcb_generic_error_t **e);
	return ((proc_xcb_input_xi_query_device_reply)fp2)(c, ((proc_xcb_input_xi_query_device_unchecked)fp1)(c, deviceid), NULL);
}
xcb_input_xi_device_info_iterator_t gamen_xcb_input_xi_query_device_infos_iterator(void *fp, const xcb_input_xi_query_device_reply_t *R) {
	typedef xcb_input_xi_device_info_iterator_t (*proc_xcb_input_xi_query_device_infos_iterator)(const xcb_input_xi_query_device_reply_t *R);
	return ((proc_xcb_input_xi_query_device_infos_iterator)fp)(R);
}
void gamen_xcb_input_xi_device_info_next(void *fp, xcb_input_xi_device_info_iterator_t *i) {
	typedef void (*proc_xcb_input_xi_device_info_next)(xcb_input_xi_device_info_iterator_t *i);
	return ((proc_xcb_input_xi_device_info_next)fp)(i);
}
xcb_input_device_class_iterator_t gamen_xcb_input_xi_device_info_classes_iterator(void *fp, const xcb_input_xi_device_info_t *R) {
	typedef xcb_input_device_class_iterator_t (*proc_xcb_input_xi_device_info_classes_iterator)(const xcb_input_xi_device_info_t *R);
	return ((proc_xcb_input_xi_device_info_classes_iterator)fp)(R);
}
void gamen_xcb_input_device_class_next(void *fp, xcb_input_device_class_iterator_t *i) {
	typedef void (*proc_xcb_input_device_class_next)(xcb_input_device_class_iterator_t *i);
	return ((proc_xcb_input_device_class_next)fp)(i);
}
xcb_void_cookie_t gamen_xcb_input_xi_select_events(void *fp, xcb_connection_t *c, xcb_window_t window, uint16_t num_mask, const xcb_input_event_mask_t *masks) {
	typedef xcb_void_cookie_t (*proc_xcb_input_xi_select_events)(xcb_connection_t *c, xcb_window_t window, uint16_t num_mask, const xcb_input_event_mask_t *masks);
	return ((proc_xcb_input_xi_select_events)fp)(c, window, num_mask, masks);
}
int gamen_xcb_input_button_press_valuator_mask_length(void *fp, const xcb_input_button_press_event_t *R) {
	typedef int (*proc_xcb_input_button_press_valuator_mask_length)(const xcb_input_button_press_event_t *R);
	return ((proc_xcb_input_button_press_valuator_mask_length)fp)(R);
}
uint32_t *gamen_xcb_input_button_press_valuator_mask(void *fp, const xcb_input_button_press_event_t *R) {
	typedef uint32_t *(*proc_xcb_input_button_press_valuator_mask)(const xcb_input_button_press_event_t *R);
	return ((proc_xcb_input_button_press_valuator_mask)fp)(R);
}
xcb_input_fp3232_t *gamen_xcb_input_button_press_axisvalues(void *fp, const xcb_input_button_press_event_t *R) {
	typedef xcb_input_fp3232_t *(*proc_xcb_input_button_press_axisvalues)(const xcb_input_button_press_event_t *R);
	return ((proc_xcb_input_button_press_axisvalues)fp)(R);
}
int gamen_xcb_input_button_press_axisvalues_length(void *fp, const xcb_input_button_press_event_t *R) {
	typedef int (*proc_xcb_input_button_press_axisvalues_length)(const xcb_input_button_press_event_t *R);
	return ((proc_xcb_input_button_press_axisvalues_length)fp)(R);
}

// libxcb-icccm
uint8_t gamen_xcb_icccm_get_wm_normal_hints_reply(void *fp1, void *fp2, xcb_connection_t *c, xcb_window_t window, xcb_size_hints_t *hints) {
	typedef xcb_get_property_cookie_t (*proc_xcb_icccm_get_wm_normal_hints_unchecked)(xcb_connection_t *c, xcb_window_t window);
	typedef uint8_t (*proc_xcb_icccm_get_wm_normal_hints_reply)(xcb_connection_t *c, xcb_get_property_cookie_t cookie, xcb_size_hints_t *hints, xcb_generic_error_t **e);
	return ((proc_xcb_icccm_get_wm_normal_hints_reply)fp2)(c, ((proc_xcb_icccm_get_wm_normal_hints_unchecked)fp1)(c, window), hints, NULL);
}
void gamen_xcb_icccm_size_hints_set_min_size(void *fp, xcb_size_hints_t *hints, int32_t min_width, int32_t min_height) {
	typedef void (*proc_xcb_icccm_size_hints_set_min_size)(xcb_size_hints_t *hints, int32_t min_width, int32_t min_height);
	return ((proc_xcb_icccm_size_hints_set_min_size)fp)(hints, min_width, min_height);
}
void gamen_xcb_icccm_size_hints_set_max_size(void *fp, xcb_size_hints_t *hints, int32_t max_width, int32_t max_height) {
	typedef void (*proc_xcb_icccm_size_hints_set_max_size)(xcb_size_hints_t *hints, int32_t max_width, int32_t max_height);
	return ((proc_xcb_icccm_size_hints_set_max_size)fp)(hints, max_width, max_height);
}
xcb_void_cookie_t gamen_xcb_icccm_set_wm_normal_hints(void *fp, xcb_connection_t *c, xcb_window_t window, xcb_size_hints_t *hints) {
	typedef xcb_void_cookie_t (*proc_xcb_icccm_set_wm_normal_hints)(xcb_connection_t *c, xcb_window_t window, xcb_size_hints_t *hints);
	return ((proc_xcb_icccm_set_wm_normal_hints)fp)(c, window, hints);
}
uint8_t gamen_xcb_icccm_get_wm_hints_reply(void *fp1, void *fp2, xcb_connection_t *c, xcb_window_t window, xcb_icccm_wm_hints_t *hints) {
	typedef xcb_get_property_cookie_t (*proc_xcb_icccm_get_wm_hints_unchecked)(xcb_connection_t *c, xcb_window_t window);
	typedef uint8_t (*proc_xcb_icccm_get_wm_hints_reply)(xcb_connection_t *c, xcb_get_property_cookie_t cookie, xcb_icccm_wm_hints_t *hints, xcb_generic_error_t **e);
	return ((proc_xcb_icccm_get_wm_hints_reply)fp2)(c, ((proc_xcb_icccm_get_wm_hints_unchecked)fp1)(c, window), hints, NULL);
}
void gamen_xcb_icccm_wm_hints_set_iconic(void *fp, xcb_icccm_wm_hints_t *hints) {
	typedef void (*proc_xcb_icccm_wm_hints_set_iconic)(xcb_icccm_wm_hints_t *hints);
	return ((proc_xcb_icccm_wm_hints_set_iconic)fp)(hints);
}
xcb_void_cookie_t gamen_xcb_icccm_set_wm_hints(void *fp, xcb_connection_t *c, xcb_window_t window, xcb_icccm_wm_hints_t *hints) {
	typedef xcb_void_cookie_t (*proc_xcb_icccm_set_wm_hints)(xcb_connection_t *c, xcb_window_t window, xcb_icccm_wm_hints_t *hints);
	return ((proc_xcb_icccm_set_wm_hints)fp)(c, window, hints);
}

// libxcb-image
xcb_pixmap_t gamen_xcb_create_pixmap_from_bitmap_data(void *fp, xcb_connection_t *display, xcb_drawable_t d, uint8_t *data, uint32_t width, uint32_t height, uint32_t depth, uint32_t fg, uint32_t bg, xcb_gcontext_t *gcp) {
	typedef xcb_pixmap_t (*proc_xcb_create_pixmap_from_bitmap_data)(xcb_connection_t *display, xcb_drawable_t d, uint8_t *data, uint32_t width, uint32_t height, uint32_t depth, uint32_t fg, uint32_t bg, xcb_gcontext_t *gcp);
	return ((proc_xcb_create_pixmap_from_bitmap_data)fp)(display, d, data, width, height, depth, fg, bg, gcp);
}

// libxcb-xkb
xcb_xkb_use_extension_reply_t *gamen_xcb_xkb_use_extension_reply(void *fp1, void *fp2, xcb_connection_t *c, uint16_t wantedMajor, uint16_t wantedMinor) {
	typedef xcb_xkb_use_extension_cookie_t (*proc_xcb_xkb_use_extension_unchecked)(xcb_connection_t *c, uint16_t wantedMajor, uint16_t wantedMinor);
	typedef xcb_xkb_use_extension_reply_t *(*proc_xcb_xkb_use_extension_reply)(xcb_connection_t *c, xcb_xkb_use_extension_cookie_t cookie, xcb_generic_error_t **e);
	return ((proc_xcb_xkb_use_extension_reply)fp2)(c, ((proc_xcb_xkb_use_extension_unchecked)fp1)(c, wantedMajor, wantedMinor), NULL);
}

*/
import "C"
import (
	"errors"
	"unsafe"
)

type xcb_library struct {
	libX11Handle       unsafe.Pointer
	libX11xcbHandle    unsafe.Pointer
	libXcursor         unsafe.Pointer
	libxcbHandle       unsafe.Pointer
	libxcbRandrHandle  unsafe.Pointer
	libxcbXinputHandle unsafe.Pointer
	libxcbIcccmHandle  unsafe.Pointer
	libxcbImageHandle  unsafe.Pointer
	libxcbXkbHandle    unsafe.Pointer

	// libX11
	XInitThreads_handle  unsafe.Pointer
	XOpenDisplay_handle  unsafe.Pointer
	XCloseDisplay_handle unsafe.Pointer

	// libX11-xcb
	XGetXCBConnection_handle   unsafe.Pointer
	XSetEventQueueOwner_handle unsafe.Pointer

	// libXcursor
	XcursorLibraryLoadCursor_handle unsafe.Pointer

	// libxcb
	xcb_get_setup_handle                       unsafe.Pointer
	xcb_setup_roots_iterator_handle            unsafe.Pointer
	xcb_screen_next_handle                     unsafe.Pointer
	xcb_get_extension_data_handle              unsafe.Pointer
	xcb_intern_atom_handle                     unsafe.Pointer
	xcb_intern_atom_reply_handle               unsafe.Pointer
	xcb_poll_for_event_handle                  unsafe.Pointer
	xcb_get_file_descriptor_handle             unsafe.Pointer
	xcb_generate_id_handle                     unsafe.Pointer
	xcb_create_window_checked_handle           unsafe.Pointer
	xcb_request_check_handle                   unsafe.Pointer
	xcb_change_property_handle                 unsafe.Pointer
	xcb_map_window_checked_handle              unsafe.Pointer
	xcb_destroy_window_handle                  unsafe.Pointer
	xcb_flush_handle                           unsafe.Pointer
	xcb_get_geometry_unchecked_handle          unsafe.Pointer
	xcb_get_geometry_reply_handle              unsafe.Pointer
	xcb_configure_window_handle                unsafe.Pointer
	xcb_get_property_unchecked_handle          unsafe.Pointer
	xcb_get_property_reply_handle              unsafe.Pointer
	xcb_get_property_value_handle              unsafe.Pointer
	xcb_get_property_value_length_handle       unsafe.Pointer
	xcb_send_event_handle                      unsafe.Pointer
	xcb_change_window_attributes_handle        unsafe.Pointer
	xcb_translate_coordinates_unchecked_handle unsafe.Pointer
	xcb_translate_coordinates_reply_handle     unsafe.Pointer
	xcb_ungrab_pointer_handle                  unsafe.Pointer
	xcb_free_pixmap_handle                     unsafe.Pointer
	xcb_create_cursor_handle                   unsafe.Pointer
	xcb_free_cursor_handle                     unsafe.Pointer

	// libxcb-randr
	xcb_randr_id                             unsafe.Pointer
	xcb_randr_query_version_unchecked_handle unsafe.Pointer
	xcb_randr_query_version_reply_handle     unsafe.Pointer
	xcb_randr_select_input_handle            unsafe.Pointer

	// libxcb-xinput
	xcb_input_id                                       unsafe.Pointer
	xcb_input_xi_query_version_unchecked_handle        unsafe.Pointer
	xcb_input_xi_query_version_reply_handle            unsafe.Pointer
	xcb_input_xi_query_device_unchecked_handle         unsafe.Pointer
	xcb_input_xi_query_device_reply_handle             unsafe.Pointer
	xcb_input_xi_query_device_infos_iterator_handle    unsafe.Pointer
	xcb_input_xi_device_info_next_handle               unsafe.Pointer
	xcb_input_xi_device_info_classes_iterator_handle   unsafe.Pointer
	xcb_input_device_class_next_handle                 unsafe.Pointer
	xcb_input_xi_select_events_handle                  unsafe.Pointer
	xcb_input_button_press_valuator_mask_length_handle unsafe.Pointer
	xcb_input_button_press_valuator_mask_handle        unsafe.Pointer
	xcb_input_button_press_axisvalues_handle           unsafe.Pointer
	xcb_input_button_press_axisvalues_length_handle    unsafe.Pointer

	// libxcb-icccm
	xcb_icccm_get_wm_normal_hints_unchecked_handle unsafe.Pointer
	xcb_icccm_get_wm_normal_hints_reply_handle     unsafe.Pointer
	xcb_icccm_size_hints_set_min_size_handle       unsafe.Pointer
	xcb_icccm_size_hints_set_max_size_handle       unsafe.Pointer
	xcb_icccm_set_wm_normal_hints_handle           unsafe.Pointer
	xcb_icccm_get_wm_hints_unchecked_handle        unsafe.Pointer
	xcb_icccm_get_wm_hints_reply_handle            unsafe.Pointer
	xcb_icccm_wm_hints_set_iconic_handle           unsafe.Pointer
	xcb_icccm_set_wm_hints_handle                  unsafe.Pointer

	// libxcb-image
	xcb_create_pixmap_from_bitmap_data_handle unsafe.Pointer

	// libxcb-xkb
	xcb_xkb_id                             unsafe.Pointer
	xcb_xkb_use_extension_unchecked_handle unsafe.Pointer
	xcb_xkb_use_extension_reply_handle     unsafe.Pointer
}

func loadSym(handle unsafe.Pointer, symbol *C.char) (unsafe.Pointer, error) {
	C.dlerror()
	fp := C.dlsym(handle, symbol)
	if fp == nil {
		err := C.dlerror()
		if err != nil {
			return nil, errors.New(C.GoString(err))
		}
	}
	return fp, nil
}

func open_xcb_library() (*xcb_library, error) {
	C.dlerror()

	l := &xcb_library{}

	l.libX11Handle = C.dlopen((*C.char)(unsafe.Pointer(&([]byte("libX11.so.6\x00"))[0])), C.RTLD_LAZY)
	if l.libX11Handle == nil {
		err := C.dlerror()
		if err != nil {
			l.close()
			return nil, errors.New(C.GoString(err))
		}
	}
	l.libXcursor = C.dlopen((*C.char)(unsafe.Pointer(&([]byte("libXcursor.so.1\x00"))[0])), C.RTLD_LAZY)
	if l.libX11xcbHandle == nil {
		err := C.dlerror()
		if err != nil {
			l.close()
			return nil, errors.New(C.GoString(err))
		}
	}
	l.libX11xcbHandle = C.dlopen((*C.char)(unsafe.Pointer(&([]byte("libX11-xcb.so.1\x00"))[0])), C.RTLD_LAZY)
	if l.libX11xcbHandle == nil {
		err := C.dlerror()
		if err != nil {
			l.close()
			return nil, errors.New(C.GoString(err))
		}
	}
	l.libxcbHandle = C.dlopen((*C.char)(unsafe.Pointer(&([]byte("libxcb.so.1\x00"))[0])), C.RTLD_LAZY)
	if l.libxcbHandle == nil {
		err := C.dlerror()
		if err != nil {
			l.close()
			return nil, errors.New(C.GoString(err))
		}
	}
	l.libxcbRandrHandle = C.dlopen((*C.char)(unsafe.Pointer(&([]byte("libxcb-randr.so.0\x00"))[0])), C.RTLD_LAZY)
	if l.libxcbRandrHandle == nil {
		err := C.dlerror()
		if err != nil {
			l.close()
			return nil, errors.New(C.GoString(err))
		}
	}
	l.libxcbXinputHandle = C.dlopen((*C.char)(unsafe.Pointer(&([]byte("libxcb-xinput.so.0\x00"))[0])), C.RTLD_LAZY)
	if l.libxcbXinputHandle == nil {
		err := C.dlerror()
		if err != nil {
			l.close()
			return nil, errors.New(C.GoString(err))
		}
	}
	l.libxcbIcccmHandle = C.dlopen((*C.char)(unsafe.Pointer(&([]byte("libxcb-icccm.so.4\x00"))[0])), C.RTLD_LAZY)
	if l.libxcbIcccmHandle == nil {
		err := C.dlerror()
		if err != nil {
			l.close()
			return nil, errors.New(C.GoString(err))
		}
	}
	l.libxcbImageHandle = C.dlopen((*C.char)(unsafe.Pointer(&([]byte("libxcb-image.so.0\x00"))[0])), C.RTLD_LAZY)
	if l.libxcbImageHandle == nil {
		err := C.dlerror()
		if err != nil {
			l.close()
			return nil, errors.New(C.GoString(err))
		}
	}
	l.libxcbXkbHandle = C.dlopen((*C.char)(unsafe.Pointer(&([]byte("libxcb-xkb.so.1\x00"))[0])), C.RTLD_LAZY)
	if l.libxcbXkbHandle == nil {
		err := C.dlerror()
		if err != nil {
			l.close()
			return nil, errors.New(C.GoString(err))
		}
	}

	var err error

	// libX11
	l.XInitThreads_handle, err = loadSym(l.libX11Handle, (*C.char)(unsafe.Pointer(&([]byte("XInitThreads\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.XOpenDisplay_handle, err = loadSym(l.libX11Handle, (*C.char)(unsafe.Pointer(&([]byte("XOpenDisplay\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.XCloseDisplay_handle, err = loadSym(l.libX11Handle, (*C.char)(unsafe.Pointer(&([]byte("XCloseDisplay\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}

	// libX11-xcb
	l.XGetXCBConnection_handle, err = loadSym(l.libX11xcbHandle, (*C.char)(unsafe.Pointer(&([]byte("XGetXCBConnection\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.XSetEventQueueOwner_handle, err = loadSym(l.libX11xcbHandle, (*C.char)(unsafe.Pointer(&([]byte("XSetEventQueueOwner\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}

	// libXcursor
	l.XcursorLibraryLoadCursor_handle, err = loadSym(l.libXcursor, (*C.char)(unsafe.Pointer(&([]byte("XcursorLibraryLoadCursor\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}

	// libxcb
	l.xcb_get_setup_handle, err = loadSym(l.libxcbHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_get_setup\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_setup_roots_iterator_handle, err = loadSym(l.libxcbHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_setup_roots_iterator\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_screen_next_handle, err = loadSym(l.libxcbHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_screen_next\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_get_extension_data_handle, err = loadSym(l.libxcbHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_get_extension_data\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_intern_atom_handle, err = loadSym(l.libxcbHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_intern_atom\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_intern_atom_reply_handle, err = loadSym(l.libxcbHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_intern_atom_reply\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_poll_for_event_handle, err = loadSym(l.libxcbHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_poll_for_event\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_get_file_descriptor_handle, err = loadSym(l.libxcbHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_get_file_descriptor\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_generate_id_handle, err = loadSym(l.libxcbHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_generate_id\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_create_window_checked_handle, err = loadSym(l.libxcbHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_create_window_checked\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_request_check_handle, err = loadSym(l.libxcbHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_request_check\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_change_property_handle, err = loadSym(l.libxcbHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_change_property\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_map_window_checked_handle, err = loadSym(l.libxcbHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_map_window_checked\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_destroy_window_handle, err = loadSym(l.libxcbHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_destroy_window\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_flush_handle, err = loadSym(l.libxcbHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_flush\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_get_geometry_unchecked_handle, err = loadSym(l.libxcbHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_get_geometry_unchecked\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_get_geometry_reply_handle, err = loadSym(l.libxcbHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_get_geometry_reply\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_configure_window_handle, err = loadSym(l.libxcbHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_configure_window\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_get_property_unchecked_handle, err = loadSym(l.libxcbHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_get_property_unchecked\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_get_property_reply_handle, err = loadSym(l.libxcbHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_get_property_reply\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_get_property_value_handle, err = loadSym(l.libxcbHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_get_property_value\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_get_property_value_length_handle, err = loadSym(l.libxcbHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_get_property_value_length\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_send_event_handle, err = loadSym(l.libxcbHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_send_event\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_change_window_attributes_handle, err = loadSym(l.libxcbHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_change_window_attributes\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_translate_coordinates_unchecked_handle, err = loadSym(l.libxcbHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_translate_coordinates_unchecked\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_translate_coordinates_reply_handle, err = loadSym(l.libxcbHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_translate_coordinates_reply\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_ungrab_pointer_handle, err = loadSym(l.libxcbHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_ungrab_pointer\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_free_pixmap_handle, err = loadSym(l.libxcbHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_free_pixmap\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_create_cursor_handle, err = loadSym(l.libxcbHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_create_cursor\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_free_cursor_handle, err = loadSym(l.libxcbHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_free_cursor\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}

	// libxcb-randr
	l.xcb_randr_id, err = loadSym(l.libxcbRandrHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_randr_id\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_randr_query_version_unchecked_handle, err = loadSym(l.libxcbRandrHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_randr_query_version_unchecked\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_randr_query_version_reply_handle, err = loadSym(l.libxcbRandrHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_randr_query_version_reply\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_randr_select_input_handle, err = loadSym(l.libxcbRandrHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_randr_select_input\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}

	// libxcb-xinput
	l.xcb_input_id, err = loadSym(l.libxcbXinputHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_input_id\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_input_xi_query_version_unchecked_handle, err = loadSym(l.libxcbXinputHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_input_xi_query_version_unchecked\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_input_xi_query_version_reply_handle, err = loadSym(l.libxcbXinputHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_input_xi_query_version_reply\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_input_xi_query_device_unchecked_handle, err = loadSym(l.libxcbXinputHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_input_xi_query_device_unchecked\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_input_xi_query_device_reply_handle, err = loadSym(l.libxcbXinputHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_input_xi_query_device_reply\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_input_xi_query_device_infos_iterator_handle, err = loadSym(l.libxcbXinputHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_input_xi_query_device_infos_iterator\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_input_xi_device_info_next_handle, err = loadSym(l.libxcbXinputHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_input_xi_device_info_next\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_input_xi_device_info_classes_iterator_handle, err = loadSym(l.libxcbXinputHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_input_xi_device_info_classes_iterator\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_input_device_class_next_handle, err = loadSym(l.libxcbXinputHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_input_device_class_next\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_input_xi_select_events_handle, err = loadSym(l.libxcbXinputHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_input_xi_select_events\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_input_button_press_valuator_mask_length_handle, err = loadSym(l.libxcbXinputHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_input_button_press_valuator_mask_length\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_input_button_press_valuator_mask_handle, err = loadSym(l.libxcbXinputHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_input_button_press_valuator_mask\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_input_button_press_axisvalues_handle, err = loadSym(l.libxcbXinputHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_input_button_press_valuator_mask\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_input_button_press_axisvalues_length_handle, err = loadSym(l.libxcbXinputHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_input_button_press_axisvalues_length\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}

	// libxcb-icccm
	l.xcb_icccm_get_wm_normal_hints_unchecked_handle, err = loadSym(l.libxcbIcccmHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_icccm_get_wm_normal_hints_unchecked\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_icccm_get_wm_normal_hints_reply_handle, err = loadSym(l.libxcbIcccmHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_icccm_get_wm_normal_hints_reply\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_icccm_size_hints_set_min_size_handle, err = loadSym(l.libxcbIcccmHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_icccm_size_hints_set_min_size\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_icccm_size_hints_set_max_size_handle, err = loadSym(l.libxcbIcccmHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_icccm_size_hints_set_max_size\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_icccm_set_wm_normal_hints_handle, err = loadSym(l.libxcbIcccmHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_icccm_set_wm_normal_hints\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_icccm_get_wm_hints_unchecked_handle, err = loadSym(l.libxcbIcccmHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_icccm_get_wm_hints_unchecked\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_icccm_get_wm_hints_reply_handle, err = loadSym(l.libxcbIcccmHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_icccm_get_wm_hints_reply\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_icccm_wm_hints_set_iconic_handle, err = loadSym(l.libxcbIcccmHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_icccm_wm_hints_set_iconic\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_icccm_set_wm_hints_handle, err = loadSym(l.libxcbIcccmHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_icccm_set_wm_hints\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}

	// libxcb-image
	l.xcb_create_pixmap_from_bitmap_data_handle, err = loadSym(l.libxcbImageHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_create_pixmap_from_bitmap_data\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}

	// libxcb-xkb
	l.xcb_xkb_id, err = loadSym(l.libxcbXkbHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_xkb_id\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_xkb_use_extension_unchecked_handle, err = loadSym(l.libxcbXkbHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_xkb_use_extension_unchecked\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xcb_xkb_use_extension_reply_handle, err = loadSym(l.libxcbXkbHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_xkb_use_extension_reply\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}

	return l, nil
}

func (l *xcb_library) close() {
	if l.libxcbXkbHandle != nil {
		C.dlclose(l.libxcbXkbHandle)
		l.libxcbXkbHandle = nil
	}
	if l.libxcbImageHandle != nil {
		C.dlclose(l.libxcbImageHandle)
		l.libxcbImageHandle = nil
	}
	if l.libxcbIcccmHandle != nil {
		C.dlclose(l.libxcbIcccmHandle)
		l.libxcbIcccmHandle = nil
	}
	if l.libxcbXinputHandle != nil {
		C.dlclose(l.libxcbXinputHandle)
		l.libxcbXinputHandle = nil
	}
	if l.libxcbRandrHandle != nil {
		C.dlclose(l.libxcbRandrHandle)
		l.libxcbRandrHandle = nil
	}
	if l.libxcbHandle != nil {
		C.dlclose(l.libxcbHandle)
		l.libxcbHandle = nil
	}
	if l.libXcursor != nil {
		C.dlclose(l.libXcursor)
		l.libXcursor = nil
	}
	if l.libX11xcbHandle != nil {
		C.dlclose(l.libX11xcbHandle)
		l.libX11xcbHandle = nil
	}
	if l.libX11Handle != nil {
		C.dlclose(l.libX11Handle)
		l.libX11Handle = nil
	}
}

// libX11
func (l *xcb_library) XInitThreads() C.Status {
	return C.gamen_XInitThreads(l.XInitThreads_handle)
}
func (l *xcb_library) XOpenDisplay(display_name *C.char) *C.Display {
	return C.gamen_XOpenDisplay(l.XOpenDisplay_handle, display_name)
}
func (l *xcb_library) XCloseDisplay(dpy *C.Display) C.int {
	return C.gamen_XCloseDisplay(l.XCloseDisplay_handle, dpy)
}

// libX11-xcb
func (l *xcb_library) XGetXCBConnection(dpy *C.Display) *C.xcb_connection_t {
	return C.gamen_XGetXCBConnection(l.XGetXCBConnection_handle, dpy)
}
func (l *xcb_library) XSetEventQueueOwner(dpy *C.Display, owner C.enum_XEventQueueOwner) {
	C.gamen_XSetEventQueueOwner(l.XSetEventQueueOwner_handle, dpy, owner)
}

// libXcursor
func (l *xcb_library) XcursorLibraryLoadCursor(dpy *C.Display, file *C.char) C.Cursor {
	return C.gamen_XcursorLibraryLoadCursor(l.XcursorLibraryLoadCursor_handle, dpy, file)
}

// libxcb
func (l *xcb_library) xcb_get_setup(c *C.xcb_connection_t) *C.struct_xcb_setup_t {
	return C.gamen_xcb_get_setup(l.xcb_get_setup_handle, c)
}
func (l *xcb_library) xcb_setup_roots_iterator(R *C.xcb_setup_t) C.xcb_screen_iterator_t {
	return C.gamen_xcb_setup_roots_iterator(l.xcb_setup_roots_iterator_handle, R)
}
func (l *xcb_library) xcb_screen_next(i *C.xcb_screen_iterator_t) {
	C.gamen_xcb_screen_next(l.xcb_screen_next_handle, i)
}
func (l *xcb_library) xcb_get_extension_data(c *C.xcb_connection_t, ext *C.xcb_extension_t) *C.struct_xcb_query_extension_reply_t {
	return C.gamen_xcb_get_extension_data(l.xcb_get_extension_data_handle, c, ext)
}
func (l *xcb_library) xcb_intern_atom(c *C.xcb_connection_t, only_if_exists C.uint8_t, name_len C.uint16_t, name *C.char) C.xcb_intern_atom_cookie_t {
	return C.gamen_xcb_intern_atom(l.xcb_intern_atom_handle, c, only_if_exists, name_len, name)
}
func (l *xcb_library) xcb_intern_atom_reply(c *C.xcb_connection_t, cookie C.xcb_intern_atom_cookie_t, e **C.xcb_generic_error_t) *C.xcb_intern_atom_reply_t {
	return C.gamen_xcb_intern_atom_reply(l.xcb_intern_atom_reply_handle, c, cookie, e)
}
func (l *xcb_library) xcb_poll_for_event(c *C.xcb_connection_t) *C.xcb_generic_event_t {
	return C.gamen_xcb_poll_for_event(l.xcb_poll_for_event_handle, c)
}
func (l *xcb_library) xcb_get_file_descriptor(c *C.xcb_connection_t) C.int {
	return C.gamen_xcb_get_file_descriptor(l.xcb_get_file_descriptor_handle, c)
}
func (l *xcb_library) xcb_generate_id(c *C.xcb_connection_t) C.uint32_t {
	return C.gamen_xcb_generate_id(l.xcb_generate_id_handle, c)
}
func (l *xcb_library) xcb_create_window_checked(c *C.xcb_connection_t, depth C.uint8_t, wid C.xcb_window_t, parent C.xcb_window_t, x C.int16_t, y C.int16_t, width C.uint16_t, height C.uint16_t, border_width C.uint16_t, _class C.uint16_t, visual C.xcb_visualid_t, value_mask C.uint32_t, value_list unsafe.Pointer) C.xcb_void_cookie_t {
	return C.gamen_xcb_create_window_checked(l.xcb_create_window_checked_handle, c, depth, wid, parent, x, y, width, height, border_width, _class, visual, value_mask, value_list)
}
func (l *xcb_library) xcb_request_check(c *C.xcb_connection_t, cookie C.xcb_void_cookie_t) *C.xcb_generic_error_t {
	return C.gamen_xcb_request_check(l.xcb_request_check_handle, c, cookie)
}
func (l *xcb_library) xcb_change_property(c *C.xcb_connection_t, mode C.uint8_t, window C.xcb_window_t, property C.xcb_atom_t, type_ C.xcb_atom_t, format C.uint8_t, data_len C.uint32_t, data unsafe.Pointer) C.xcb_void_cookie_t {
	return C.gamen_xcb_change_property(l.xcb_change_property_handle, c, mode, window, property, type_, format, data_len, data)
}
func (l *xcb_library) xcb_map_window_checked(c *C.xcb_connection_t, window C.xcb_window_t) C.xcb_void_cookie_t {
	return C.gamen_xcb_map_window_checked(l.xcb_map_window_checked_handle, c, window)
}
func (l *xcb_library) xcb_destroy_window(c *C.xcb_connection_t, window C.xcb_window_t) C.xcb_void_cookie_t {
	return C.gamen_xcb_destroy_window(l.xcb_destroy_window_handle, c, window)
}
func (l *xcb_library) xcb_flush(c *C.xcb_connection_t) C.int {
	return C.gamen_xcb_flush(l.xcb_flush_handle, c)
}
func (l *xcb_library) xcb_get_geometry_reply(c *C.xcb_connection_t, drawable C.xcb_drawable_t) *C.xcb_get_geometry_reply_t {
	return C.gamen_xcb_get_geometry_reply(l.xcb_get_geometry_unchecked_handle, l.xcb_get_geometry_reply_handle, c, drawable)
}
func (l *xcb_library) xcb_configure_window(c *C.xcb_connection_t, window C.xcb_window_t, value_mask C.uint16_t, value_list unsafe.Pointer) C.xcb_void_cookie_t {
	return C.gamen_xcb_configure_window(l.xcb_configure_window_handle, c, window, value_mask, value_list)
}
func (l *xcb_library) xcb_get_property_reply(c *C.xcb_connection_t, _delete C.uint8_t, window C.xcb_window_t, property C.xcb_atom_t, type_ C.xcb_atom_t, long_offset C.uint32_t, long_length C.uint32_t) *C.xcb_get_property_reply_t {
	return C.gamen_xcb_get_property_reply(l.xcb_get_property_unchecked_handle, l.xcb_get_property_reply_handle, c, _delete, window, property, type_, long_offset, long_length)
}
func (l *xcb_library) xcb_get_property_value(R *C.xcb_get_property_reply_t) unsafe.Pointer {
	return C.gamen_xcb_get_property_value(l.xcb_get_property_value_handle, R)
}
func (l *xcb_library) xcb_get_property_value_length(R *C.xcb_get_property_reply_t) C.int {
	return C.gamen_xcb_get_property_value_length(l.xcb_get_property_value_length_handle, R)
}
func (l *xcb_library) xcb_send_event(c *C.xcb_connection_t, propagate C.uint8_t, destination C.xcb_window_t, event_mask C.uint32_t, event *C.char) C.xcb_void_cookie_t {
	return C.gamen_xcb_send_event(l.xcb_send_event_handle, c, propagate, destination, event_mask, event)
}
func (l *xcb_library) xcb_change_window_attributes(c *C.xcb_connection_t, window C.xcb_window_t, value_mask C.uint32_t, value_list unsafe.Pointer) C.xcb_void_cookie_t {
	return C.gamen_xcb_change_window_attributes(l.xcb_change_window_attributes_handle, c, window, value_mask, value_list)
}
func (l *xcb_library) xcb_translate_coordinates_reply(c *C.xcb_connection_t, src_window C.xcb_window_t, dst_window C.xcb_window_t, src_x C.int16_t, src_y C.int16_t) *C.xcb_translate_coordinates_reply_t {
	return C.gamen_xcb_translate_coordinates_reply(l.xcb_translate_coordinates_unchecked_handle, l.xcb_translate_coordinates_reply_handle, c, src_window, dst_window, src_x, src_y)
}
func (l *xcb_library) xcb_ungrab_pointer(c *C.xcb_connection_t, time C.xcb_timestamp_t) C.xcb_void_cookie_t {
	return C.gamen_xcb_ungrab_pointer(l.xcb_ungrab_pointer_handle, c, time)
}
func (l *xcb_library) xcb_free_pixmap(c *C.xcb_connection_t, pixmap C.xcb_pixmap_t) C.xcb_void_cookie_t {
	return C.gamen_xcb_free_pixmap(l.xcb_free_pixmap_handle, c, pixmap)
}
func (l *xcb_library) xcb_create_cursor(c *C.xcb_connection_t, cid C.xcb_cursor_t, source C.xcb_pixmap_t, mask C.xcb_pixmap_t, fore_red C.uint16_t, fore_green C.uint16_t, fore_blue C.uint16_t, back_red C.uint16_t, back_green C.uint16_t, back_blue C.uint16_t, x C.uint16_t, y C.uint16_t) C.xcb_void_cookie_t {
	return C.gamen_xcb_create_cursor(l.xcb_create_cursor_handle, c, cid, source, mask, fore_red, fore_green, fore_blue, back_red, back_green, back_blue, x, y)
}
func (l *xcb_library) xcb_free_cursor(c *C.xcb_connection_t, cursor C.xcb_cursor_t) C.xcb_void_cookie_t {
	return C.gamen_xcb_free_cursor(l.xcb_free_cursor_handle, c, cursor)
}

// libxcb-randr
func (l *xcb_library) xcb_randr_query_version_reply(c *C.xcb_connection_t, major_version C.uint32_t, minor_version C.uint32_t) *C.xcb_randr_query_version_reply_t {
	return C.gamen_xcb_randr_query_version_reply(l.xcb_randr_query_version_unchecked_handle, l.xcb_randr_query_version_reply_handle, c, major_version, minor_version)
}
func (l *xcb_library) xcb_randr_select_input(c *C.xcb_connection_t, window C.xcb_window_t, enable C.uint16_t) C.xcb_void_cookie_t {
	return C.gamen_xcb_randr_select_input(l.xcb_randr_select_input_handle, c, window, enable)
}

// libxcb-xinput
func (l *xcb_library) xcb_input_xi_query_version_reply(c *C.xcb_connection_t, major_version C.uint16_t, minor_version C.uint16_t) *C.xcb_input_xi_query_version_reply_t {
	return C.gamen_xcb_input_xi_query_version_reply(l.xcb_input_xi_query_version_unchecked_handle, l.xcb_input_xi_query_version_reply_handle, c, major_version, minor_version)
}
func (l *xcb_library) xcb_input_xi_query_device_reply(c *C.xcb_connection_t, deviceid C.xcb_input_device_id_t) *C.xcb_input_xi_query_device_reply_t {
	return C.gamen_xcb_input_xi_query_device_reply(l.xcb_input_xi_query_device_unchecked_handle, l.xcb_input_xi_query_device_reply_handle, c, deviceid)
}
func (l *xcb_library) xcb_input_xi_query_device_infos_iterator(R *C.xcb_input_xi_query_device_reply_t) C.xcb_input_xi_device_info_iterator_t {
	return C.gamen_xcb_input_xi_query_device_infos_iterator(l.xcb_input_xi_query_device_infos_iterator_handle, R)
}
func (l *xcb_library) xcb_input_xi_device_info_next(i *C.xcb_input_xi_device_info_iterator_t) {
	C.gamen_xcb_input_xi_device_info_next(l.xcb_input_xi_device_info_next_handle, i)
}
func (l *xcb_library) xcb_input_xi_device_info_classes_iterator(R *C.xcb_input_xi_device_info_t) C.xcb_input_device_class_iterator_t {
	return C.gamen_xcb_input_xi_device_info_classes_iterator(l.xcb_input_xi_device_info_classes_iterator_handle, R)
}
func (l *xcb_library) xcb_input_device_class_next(i *C.xcb_input_device_class_iterator_t) {
	C.gamen_xcb_input_device_class_next(l.xcb_input_device_class_next_handle, i)
}
func (l *xcb_library) xcb_input_xi_select_events(c *C.xcb_connection_t, window C.xcb_window_t, num_mask C.uint16_t, masks *C.xcb_input_event_mask_t) C.xcb_void_cookie_t {
	return C.gamen_xcb_input_xi_select_events(l.xcb_input_xi_select_events_handle, c, window, num_mask, masks)
}
func (l *xcb_library) xcb_input_button_press_valuator_mask_length(R *C.xcb_input_button_press_event_t) C.int {
	return C.gamen_xcb_input_button_press_valuator_mask_length(l.xcb_input_button_press_valuator_mask_length_handle, R)
}
func (l *xcb_library) xcb_input_button_press_valuator_mask(R *C.xcb_input_button_press_event_t) *C.uint32_t {
	return C.gamen_xcb_input_button_press_valuator_mask(l.xcb_input_button_press_valuator_mask_handle, R)
}
func (l *xcb_library) xcb_input_button_press_axisvalues(R *C.xcb_input_button_press_event_t) *C.xcb_input_fp3232_t {
	return C.gamen_xcb_input_button_press_axisvalues(l.xcb_input_button_press_axisvalues_handle, R)
}
func (l *xcb_library) xcb_input_button_press_axisvalues_length(R *C.xcb_input_button_press_event_t) C.int {
	return C.gamen_xcb_input_button_press_axisvalues_length(l.xcb_input_button_press_axisvalues_length_handle, R)
}

// libxcb-icccm
func (l *xcb_library) xcb_icccm_get_wm_normal_hints_reply(c *C.xcb_connection_t, window C.xcb_window_t, hints *C.xcb_size_hints_t) C.uint8_t {
	return C.gamen_xcb_icccm_get_wm_normal_hints_reply(l.xcb_icccm_get_wm_normal_hints_unchecked_handle, l.xcb_icccm_get_wm_normal_hints_reply_handle, c, window, hints)
}
func (l *xcb_library) xcb_icccm_size_hints_set_min_size(hints *C.xcb_size_hints_t, min_width C.int32_t, min_height C.int32_t) {
	C.gamen_xcb_icccm_size_hints_set_min_size(l.xcb_icccm_size_hints_set_min_size_handle, hints, min_width, min_height)
}
func (l *xcb_library) xcb_icccm_size_hints_set_max_size(hints *C.xcb_size_hints_t, max_width C.int32_t, max_height C.int32_t) {
	C.gamen_xcb_icccm_size_hints_set_max_size(l.xcb_icccm_size_hints_set_max_size_handle, hints, max_width, max_height)
}
func (l *xcb_library) xcb_icccm_set_wm_normal_hints(c *C.xcb_connection_t, window C.xcb_window_t, hints *C.xcb_size_hints_t) C.xcb_void_cookie_t {
	return C.gamen_xcb_icccm_set_wm_normal_hints(l.xcb_icccm_set_wm_normal_hints_handle, c, window, hints)
}
func (l *xcb_library) xcb_icccm_get_wm_hints_reply(c *C.xcb_connection_t, window C.xcb_window_t, hints *C.xcb_icccm_wm_hints_t) C.uint8_t {
	return C.gamen_xcb_icccm_get_wm_hints_reply(l.xcb_icccm_get_wm_hints_unchecked_handle, l.xcb_icccm_get_wm_hints_reply_handle, c, window, hints)
}
func (l *xcb_library) xcb_icccm_wm_hints_set_iconic(hints *C.xcb_icccm_wm_hints_t) {
	C.gamen_xcb_icccm_wm_hints_set_iconic(l.xcb_icccm_wm_hints_set_iconic_handle, hints)
}
func (l *xcb_library) xcb_icccm_set_wm_hints(c *C.xcb_connection_t, window C.xcb_window_t, hints *C.xcb_icccm_wm_hints_t) C.xcb_void_cookie_t {
	return C.gamen_xcb_icccm_set_wm_hints(l.xcb_icccm_set_wm_hints_handle, c, window, hints)
}

// libxcb-image
func (l *xcb_library) xcb_create_pixmap_from_bitmap_data(display *C.xcb_connection_t, d C.xcb_drawable_t, data *C.uint8_t, width C.uint32_t, height C.uint32_t, depth C.uint32_t, fg C.uint32_t, bg C.uint32_t, gcp *C.xcb_gcontext_t) C.xcb_pixmap_t {
	return C.gamen_xcb_create_pixmap_from_bitmap_data(l.xcb_create_pixmap_from_bitmap_data_handle, display, d, data, width, height, depth, fg, bg, gcp)
}

// libxcb-xkb
func (l *xcb_library) xcb_xkb_use_extension_reply(c *C.xcb_connection_t, wantedMajor C.uint16_t, wantedMinor C.uint16_t) *C.xcb_xkb_use_extension_reply_t {
	return C.gamen_xcb_xkb_use_extension_reply(l.xcb_xkb_use_extension_unchecked_handle, l.xcb_xkb_use_extension_reply_handle, c, wantedMajor, wantedMinor)
}
