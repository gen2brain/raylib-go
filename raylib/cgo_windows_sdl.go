//go:build windows && sdl && !rgfw
// +build windows,sdl,!rgfw

package rl

/*
#cgo windows LDFLAGS: -lgdi32 -lwinmm -lole32
#cgo windows,sdl LDFLAGS: -lSDL2
#cgo windows CFLAGS: -Iexternal -Wno-stringop-overflow -Dgbutton=cbutton -DSDL_GetJoysticks=rl_sdl2_GetJoysticks -include sdl_shim.h
#cgo windows,sdl CFLAGS: -DPLATFORM_DESKTOP_SDL -DUSING_SDL2_PROJECT

#cgo windows,!es2,!es3 LDFLAGS: -lopengl32

#cgo windows,opengl11,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_11
#cgo windows,opengl21,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_21
#cgo windows,opengl43,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_43
#cgo windows,!opengl11,!opengl21,!opengl43,!es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_33
#cgo windows,es2,!es3 CFLAGS: -DGRAPHICS_API_OPENGL_ES2
#cgo windows,es3,!es2 CFLAGS: -DGRAPHICS_API_OPENGL_ES3

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
