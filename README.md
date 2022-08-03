# gamen

[![Go Reference](https://pkg.go.dev/badge/github.com/rajveermalviya/gamen.svg)](https://pkg.go.dev/github.com/rajveermalviya/gamen)
[![gamen Matrix](https://img.shields.io/static/v1?label&message=%23gamen&color=blueviolet&logo=matrix)](https://matrix.to/#gamen:matrix.org)

`gamen` is cross-platform windowing library in Go. It natively supports Windows, Linux, Android and Web. on Linux both X11 (via xcb) and Wayland are supported.

`gamen` provides api for creating and handling windows. It also lets you handle events generated by platform window via callbacks. It is fairly low level, it provides native handles for graphics APIs like OpenGL and Vulkan to initialize from, `gamen` by itself doesn't provide you a drawing primitive you have to use a library like [`go-webgpu`](https://github.com/rajveermalviya/go-webgpu) or similar.

`gamen` has a callback based api for handling events i.e it doesn't do queueing<sup>*</sup> by itself. Because most of the backends already do internal queueing of events, doing it again inside the library can introduce unnecessary latency. Also this keep api flexible, a separate library can introduce event queue on top of `gamen` for easier developer experience.
##### * web backend uses channels.

## usage

```go
package main

import (
	"runtime"

	"github.com/rajveermalviya/gamen/display"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	d, err := display.NewDisplay()
	if err != nil {
		panic(err)
	}
	defer d.Destroy()

	w, err := display.NewWindow(d)
	if err != nil {
		panic(err)
	}
	defer w.Destroy()

	w.SetCloseRequestedCallback(func() { d.Destroy() })

	for {
		// render here

		if !d.Wait() {
			break
		}
	}
}
```

check out more examples under [`examples/`](./examples/) dir

examples `wgpu_poll` and `wgpu_wait` shows how to use the event loop for a Game and a GUI respectively. Though an ideal client will switch between `Poll` and `Wait`. (GUI temporarily showing an animation)

## dependencies

### windows

windows (win32) backend **does not** use CGO, i.e **does not** require a C toolchain, only Go compiler is enough

### linux

resulting binaries shouldn't require any dependency to be installed by the users. but developers will need some `devel` packages.

#### fedora

```shell
sudo dnf install wayland-devel libX11-devel libXcursor-devel libxkbcommon-x11-devel xcb-util-image-devel xcb-util-wm-devel
```

#### ubuntu

```shell
sudo apt install libwayland-dev libxkbcommon-x11-dev libx11-xcb-dev libxcb-randr0-dev libxcb-xinput-dev libxcb-icccm4-dev libxcursor-dev libxcb-image0-dev
```

<!-- TODO: other distros -->

### android

android backend uses [`game-activity`](https://developer.android.com/games/agdk/game-activity).

[`tsukuru`](https://github.com/rajveermalviya/tsukuru) should be used to help with dependency resolution and building the android app.

```shell
# make sure you have android sdk installed
# connect your device and setup adb connection / run android emulator

go install github.com/rajveermalviya/tsukuru@latest

tsukuru run apk ./example/hello
```

## features

an incomplete list of features in no particular order that are supported or that we want to support but aren't currently.

| feature                      | win32              | xcb                |  wayland           | android            | web                |
| ---------------------------- | ------------------ | ------------------ | ------------------ | ------------------ | ------------------ |
| window initialization        | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: |
| handles for OpenGL init      | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: |
| handles for Vulkan WSI       | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: |
| show window with decorations | :heavy_check_mark: | :heavy_check_mark: | :exclamation: [#2] | :heavy_check_mark: | **N/A**            |
| window decorations toggle    | :x:                | :x:                | :x:                | :x:                | **N/A**            |
| window resizing events       | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | **N/A**            |
| resizing window manually     | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | **N/A**            | **N/A**            |
| window transparency          | :x:                | :x:                | :x:                | :x:                | :x:                |
| window maximization toggle   | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | **N/A**            | **N/A**            |
| window minimization          | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | **N/A**            | **N/A**            |
| fullscreen toggle            | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: |
| HiDPI support                | :x:                | :x:                | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: |
| popups                       | :x:                | :x:                | :x:                | **N/A**            | **N/A**            |
| monitor list                 | :x:                | :x:                | :x:                | :x:                | :x:                |
| video mode query             | :x:                | :x:                | :x:                | :x:                | :x:                |
| mouse events                 | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | **N/A**            | :heavy_check_mark: |
| cursor locking               | :x:                | :x:                | :x:                | :x:                | :x:                |
| cursor confining             | :x:                | :x:                | :x:                | :x:                | :x:                |
| cursor icon                  | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | **N/A**            | :heavy_check_mark: |
| cursor hittest               | :x:                | :x:                | :x:                | **N/A**            | :x:                |
| touch events                 | :x:                | :x:                | :x:                | :heavy_check_mark: | :x:                |
| keyboard events              | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: | :heavy_check_mark: |
| drag window with cursor      | :x:                | :x:                | :x:                | **N/A**            | **N/A**            |
| drag & drop                  | :x:                | :x:                | :x:                | **N/A**            | :x:                |
| raw device events            | :x:                | :x:                | :x:                | :x:                | :x:                |
| gamepad/joystick events      | :x:                | :x:                | :x:                | :x:                | :x:                |
| ime                          | :x:                | :x:                | :x:                | :x:                | :x:                |
| clipboard                    | :x:                | :x:                | :x:                | :x:                | :x:                |
| theme change events          | :x:                | :x:                | :x:                | :x:                | :x:                |

[#2]: https://github.com/rajveermalviya/gamen/issues/2

(as you can see there are many :x:s so any help will be greatly appreciated)

Note: for macos/ios support see [#1](https://github.com/rajveermalviya/gamen/issues/1)

## contact

join [matrix room](https://matrix.to/#gamen:matrix.org)

## thanks

`gamen`'s api is a mix of glfw and winit and some terminologies from Wayland.

- glfw - https://github.com/glfw/glfw
- winit - https://github.com/rust-windowing/winit
