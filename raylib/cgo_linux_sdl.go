//go:build linux && sdl && !rgfw && !drm && !android
// +build linux,sdl,!rgfw,!drm,!android

package rl

/*
#cgo linux,!es2 LDFLAGS: -lm
#cgo linux CFLAGS: -Wno-stringop-overflow -Dgbutton=cbutton -DSDL_GetJoysticks=rl_sdl2_GetJoysticks -include sdl_shim.h
#cgo linux,sdl CFLAGS: -DPLATFORM_DESKTOP_SDL -DUSING_SDL2_PROJECT
#cgo linux,sdl pkg-config: sdl2

#cgo linux,!es2,!es3 LDFLAGS: -lGL

#cgo linux,opengl11,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_11
#cgo linux,opengl21,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_21
#cgo linux,opengl43,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_43
#cgo linux,!opengl11,!opengl21,!opengl43,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_33
#cgo linux,es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_ES2
#cgo linux,es3,!es2 CFLAGS: -DGRAPHICS_API_OPENGL_ES3

#ifdef USING_SDL2_PROJECT
#include "SDL2/SDL.h"
SDL_JoystickID* rl_sdl2_GetJoysticks(int* count) {
    int n = SDL_NumJoysticks();
    if (count) *count = n;
    SDL_JoystickID* ids = (SDL_JoystickID*)SDL_malloc(sizeof(SDL_JoystickID) * (n > 0 ? n : 1));
    if (!ids) return NULL;
    for (int i = 0; i < n; ++i) {
        if (!SDL_IsGameController(i)) { ids[i] = -1; continue; }
        SDL_GameController* gc = SDL_GameControllerOpen(i);
        if (!gc) { ids[i] = -1; continue; }
        ids[i] = SDL_JoystickInstanceID(SDL_GameControllerGetJoystick(gc));
        SDL_GameControllerClose(gc);
    }
    return ids;
}
#endif
*/
import "C"
