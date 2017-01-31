package raylib

/*
#cgo linux,!arm LDFLAGS: -lglfw -lGL -lopenal -lm -pthread -ldl -lX11 -lXrandr -lXinerama -lXi -lXxf86vm -lXcursor
#cgo linux,arm,!android LDFLAGS: -lGLESv2 -lEGL -lpthread -lrt -lm -lbcm_host -lvcos -lvchiq_arm -lopenal
#cgo windows LDFLAGS: -lglfw3 -lopengl32 -lgdi32 -lopenal32 -lwinmm
#cgo darwin LDFLAGS: -lglfw -framework OpenGL -framework OpenAL -framework Cocoa
#cgo android LDFLAGS: -llog -landroid -lEGL -lGLESv2 -lOpenSLES -lopenal -lm -landroid_native_app_glue

#cgo CFLAGS: -std=gnu99 -fgnu89-inline -Wno-missing-braces -Wno-unused-result
#cgo linux,windows,darwin,!android,!arm CFLAGS: -DPLATFORM_DESKTOP -DGRAPHICS_API_OPENGL_33
#cgo linux,arm,!android CFLAGS: -DPLATFORM_RPI -DGRAPHICS_API_OPENGL_ES2
#cgo android CFLAGS: -DPLATFORM_ANDROID -DGRAPHICS_API_OPENGL_ES2
*/
import "C"
