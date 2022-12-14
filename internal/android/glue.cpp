//go:build android

#include <stdbool.h>
#include <game-activity/GameActivity.cpp>
#include <game-text-input/gametextinput.cpp>

extern "C" {
  #include "_cgo_export.h"
  #include <game-activity/native_app_glue/android_native_app_glue.c>
}

extern "C" {

void display_poll(int timeoutMillis) {
  int events;
  struct android_poll_source* source;

  while ((ALooper_pollAll(timeoutMillis, NULL, &events, (void **) &source)) >= 0) {
    if (source != NULL) {
      source->process(source->app, source);
    }
  }
}

void display_set_handler(struct android_app* app) {
  app->onAppCmd = display_handle_command;

  android_app_set_key_event_filter(app, (android_key_event_filter)display_handle_key_event);
  android_app_set_motion_event_filter(app, (android_motion_event_filter)display_handle_motion_event);
}

}
