//go:build linux && !android

#include "_cgo_export.h"
#include "xdg-shell-client-protocol.h"
#include <stdio.h>
#include <wayland-client-protocol.h>

const struct wl_registry_listener wl_registry_listener = {
    .global = (void (*)(void *, struct wl_registry *, uint32_t, const char *, uint32_t))registryHandleGlobal,
    .global_remove = registryHandleGlobalRemove,
};

void outputHandleName(void *data, struct wl_output *wl_output, const char *name) {}

void outputHandleDescription(void *data, struct wl_output *wl_output, const char *description) {}

const struct wl_output_listener wl_output_listener = {
    .geometry = (void (*)(void *, struct wl_output *, int32_t, int32_t, int32_t, int32_t, int32_t, const char *, const char *, int32_t))outputHandleGeometry,
    .mode = outputHandleMode,
    .done = outputHandleDone,
    .scale = outputHandleScale,
    .name = outputHandleName,
    .description = outputHandleDescription,
};

void xdgWmBaseHandlePing(void *data, struct xdg_wm_base *xdg_wm_base, uint32_t serial) {
  xdg_wm_base_pong(xdg_wm_base, serial);
}

const struct xdg_wm_base_listener xdg_wm_base_listener = {
    .ping = xdgWmBaseHandlePing,
};

void seatHandleName(void *data, struct wl_seat *seat, const char *name) {}

const struct wl_seat_listener wl_seat_listener = {
    .capabilities = seatHandleCapabilities,
    .name = seatHandleName,
};

void pointerHandleMotion_cgo(void *data, struct wl_pointer *wl_pointer, uint32_t time, wl_fixed_t surface_x, wl_fixed_t surface_y) {
  pointerHandleMotion(data, wl_pointer, time, wl_fixed_to_double(surface_x), wl_fixed_to_double(surface_y));
}

void pointerHandleAxis_cgo(void *data, struct wl_pointer *wl_pointer, uint32_t time, uint32_t axis, wl_fixed_t value) {
  pointerHandleAxis(data, wl_pointer, time, axis, wl_fixed_to_double(value));
}

void pointerHandleAxisSource(void *data, struct wl_pointer *wl_pointer, uint32_t axis_source) {}

void pointerHandleAxisStop(void *data, struct wl_pointer *wl_pointer, uint32_t time, uint32_t axis) {}

const struct wl_pointer_listener wl_pointer_listener = {
    .enter = pointerHandleEnter,
    .leave = pointerHandleLeave,
    .motion = pointerHandleMotion_cgo,
    .button = pointerHandleButton,
    .axis = pointerHandleAxis_cgo,
    .axis_discrete = pointerHandleAxisDiscrete,
    .axis_source = pointerHandleAxisSource,
    .axis_stop = pointerHandleAxisStop,
    .frame = pointerHandleFrame,
};

const struct wl_keyboard_listener wl_keyboard_listener = {
    .keymap = keyboardHandleKeymap,
    .enter = keyboardHandleEnter,
    .leave = keyboardHandleLeave,
    .key = keyboardHandleKey,
    .modifiers = keyboardHandleModifiers,
    .repeat_info = keyboardHandleRepeatInfo,
};

const struct wl_surface_listener window_surface_listener = {
    .enter = windowSurfaceHandleEnter,
    .leave = windowSurfaceHandleLeave,
};

void xdgSurfaceHandleConfigure(void *data, struct xdg_surface *xdg_surface, uint32_t serial) {
  xdg_surface_ack_configure(xdg_surface, serial);
}

const struct xdg_surface_listener xdg_surface_listener = {
    .configure = xdgSurfaceHandleConfigure,
};

void xdgToplevelConfigureBounds(void *data, struct xdg_toplevel *xdg_toplevel, int32_t width, int32_t height) {}

const struct xdg_toplevel_listener xdg_toplevel_listener = {
    .configure = xdgToplevelHandleConfigure,
    .close = xdgToplevelHandleClose,
    .configure_bounds = xdgToplevelConfigureBounds,
};

const struct wl_callback_listener go_wl_callback_listener  = {
    .done = goWlCallbackDone,
};
