// Code generated by internal/wayland/wl/gen; DO NOT EDIT.
// XML file : ./protocols/xdg-decoration-unstable-v1.xml

//go:build linux && !android

package wayland

/*

#include "xdg-decoration-unstable-v1-client-protocol.h"

*/
import "C"
import "unsafe"

// destroy the decoration manager object
//
// Destroy the decoration manager. This doesn't destroy objects created
// with the manager.
func (l *wl_library) zxdg_decoration_manager_v1_destroy(zxdg_decoration_manager_v1 *C.struct_zxdg_decoration_manager_v1) {
	C.gamen_zxdg_decoration_manager_v1_destroy(l.wl_proxy_marshal_flags, l.wl_proxy_get_version, zxdg_decoration_manager_v1)
}

// create a new toplevel decoration object
//
// Create a new decoration object associated with the given toplevel.
//
// Creating an xdg_toplevel_decoration from an xdg_toplevel which has a
// buffer attached or committed is a client error, and any attempts by a
// client to attach or manipulate a buffer prior to the first
// xdg_toplevel_decoration.configure event must also be treated as
// errors.
func (l *wl_library) zxdg_decoration_manager_v1_get_toplevel_decoration(zxdg_decoration_manager_v1 *C.struct_zxdg_decoration_manager_v1, toplevel *C.struct_xdg_toplevel) *C.struct_zxdg_toplevel_decoration_v1 {
	return C.gamen_zxdg_decoration_manager_v1_get_toplevel_decoration(l.wl_proxy_marshal_flags, l.wl_proxy_get_version, zxdg_decoration_manager_v1, toplevel)
}

type zxdg_toplevel_decoration_v1_error C.uint32_t

const (
	// xdg_toplevel has a buffer attached before configure
	ZXDG_TOPLEVEL_DECORATION_V_1_ERROR_UNCONFIGURED_BUFFER zxdg_toplevel_decoration_v1_error = 0
	// xdg_toplevel already has a decoration object
	ZXDG_TOPLEVEL_DECORATION_V_1_ERROR_ALREADY_CONSTRUCTED zxdg_toplevel_decoration_v1_error = 1
	// xdg_toplevel destroyed before the decoration object
	ZXDG_TOPLEVEL_DECORATION_V_1_ERROR_ORPHANED zxdg_toplevel_decoration_v1_error = 2
)

// These values describe window decoration modes.
type zxdg_toplevel_decoration_v1_mode C.uint32_t

const (
	// no server-side window decoration
	ZXDG_TOPLEVEL_DECORATION_V_1_MODE_CLIENT_SIDE zxdg_toplevel_decoration_v1_mode = 1
	// server-side window decoration
	ZXDG_TOPLEVEL_DECORATION_V_1_MODE_SERVER_SIDE zxdg_toplevel_decoration_v1_mode = 2
)

func (l *wl_library) zxdg_toplevel_decoration_v1_add_listener(zxdg_toplevel_decoration_v1 *C.struct_zxdg_toplevel_decoration_v1, listener *C.struct_zxdg_toplevel_decoration_v1_listener, data unsafe.Pointer) C.int {
	return C.gamen_zxdg_toplevel_decoration_v1_add_listener(l.wl_proxy_add_listener_handle, zxdg_toplevel_decoration_v1, listener, data)
}

// destroy the decoration object
//
// Switch back to a mode without any server-side decorations at the next
// commit.
func (l *wl_library) zxdg_toplevel_decoration_v1_destroy(zxdg_toplevel_decoration_v1 *C.struct_zxdg_toplevel_decoration_v1) {
	C.gamen_zxdg_toplevel_decoration_v1_destroy(l.wl_proxy_marshal_flags, l.wl_proxy_get_version, zxdg_toplevel_decoration_v1)
}

// set the decoration mode
//
// Set the toplevel surface decoration mode. This informs the compositor
// that the client prefers the provided decoration mode.
//
// After requesting a decoration mode, the compositor will respond by
// emitting an xdg_surface.configure event. The client should then update
// its content, drawing it without decorations if the received mode is
// server-side decorations. The client must also acknowledge the configure
// when committing the new content (see xdg_surface.ack_configure).
//
// The compositor can decide not to use the client's mode and enforce a
// different mode instead.
//
// Clients whose decoration mode depend on the xdg_toplevel state may send
// a set_mode request in response to an xdg_surface.configure event and wait
// for the next xdg_surface.configure event to prevent unwanted state.
// Such clients are responsible for preventing configure loops and must
// make sure not to send multiple successive set_mode requests with the
// same decoration mode.
func (l *wl_library) zxdg_toplevel_decoration_v1_set_mode(zxdg_toplevel_decoration_v1 *C.struct_zxdg_toplevel_decoration_v1, mode C.uint32_t) {
	C.gamen_zxdg_toplevel_decoration_v1_set_mode(l.wl_proxy_marshal_flags, l.wl_proxy_get_version, zxdg_toplevel_decoration_v1, mode)
}

// unset the decoration mode
//
// Unset the toplevel surface decoration mode. This informs the compositor
// that the client doesn't prefer a particular decoration mode.
//
// This request has the same semantics as set_mode.
func (l *wl_library) zxdg_toplevel_decoration_v1_unset_mode(zxdg_toplevel_decoration_v1 *C.struct_zxdg_toplevel_decoration_v1) {
	C.gamen_zxdg_toplevel_decoration_v1_unset_mode(l.wl_proxy_marshal_flags, l.wl_proxy_get_version, zxdg_toplevel_decoration_v1)
}
