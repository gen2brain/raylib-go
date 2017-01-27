package raylib

/*
#cgo linux,!arm LDFLAGS: -lraylib -lglfw3 -lGL -lopenal -lm -pthread -ldl -lX11 -lXrandr -lXinerama -lXi -lXxf86vm -lXcursor
#cgo linux,arm,!android LDFLAGS: -lraylib -lGLESv2 -lEGL -lpthread -lrt -lm -lbcm_host -lvcos -lvchiq_arm -lopenal
#cgo windows LDFLAGS: -lraylib -lglfw3 -lopengl32 -lgdi32 -lopenal32 -lwinmm
#cgo darwin LDFLAGS: -lraylib -lglfw3 -framework OpenGL -framework OpenAl -framework Cocoa
#cgo android LDFLAGS: -lraylib -llog -landroid -lEGL -lGLESv2 -lOpenSLES -lopenal -lm -landroid_native_app_glue

#cgo linux,windows,darwin,!android,!arm CFLAGS: -DPLATFORM_DESKTOP
#cgo linux,arm,!android CFLAGS: -DPLATFORM_RPI
#cgo android CFLAGS: -DPLATFORM_ANDROID
*/
import "C"
