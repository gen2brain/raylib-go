/**********************************************************************************************
*
*   rcore_desktop_linux - Functions to manage window, graphics device and inputs
*
*   PLATFORM: DESKTOP: LINUX
*       - Linux (Wayland compositors, X11 servers, DRM/KMS without a display server)
*
*   LIMITATIONS:
*       - Limitations of the Wayland, X11 and DRM platforms
*
*   CONFIGURATION:
*       RAYLIB_PLATFORM environment variable
*           wayland, x11, drm: use that platform (no fallback),
*           not set: Wayland if WAYLAND_DISPLAY is set, X11 if DISPLAY is set, DRM otherwise,
*           a platform that fails to initialize falls back to the next one (Wayland, X11, DRM)
*
*   DEPENDENCIES:
*       - rcore_drm: DRM platform, compiled in this unit
*       - rcore_desktop_wayland: Wayland platform, compiled in this unit
*       - rcore_desktop_x11: X11 platform, compiled in this unit
*
*
*   LICENSE: zlib/libpng
*
*   Copyright (c) 2013-2026 Ramon Santamaria (@raysan5) and contributors
*
*   This software is provided "as-is", without any express or implied warranty. In no event
*   will the authors be held liable for any damages arising from the use of this software.
*
*   Permission is granted to anyone to use this software for any purpose, including commercial
*   applications, and to alter it and redistribute it freely, subject to the following restrictions:
*
*     1. The origin of this software must not be misrepresented; you must not claim that you
*     wrote the original software. If you use this software in a product, an acknowledgment
*     in the product documentation would be appreciated but is not required.
*
*     2. Altered source versions must be plainly marked as such, and must not be misrepresented
*     as being the original software.
*
*     3. This notice may not be removed or altered from any source distribution.
*
**********************************************************************************************/

#include <stdlib.h>                 // Required for: getenv()
#include <string.h>                 // Required for: strcmp()

//----------------------------------------------------------------------------------
// Platforms
//----------------------------------------------------------------------------------
// NOTE: All platforms are compiled in this unit, platform functions and names defined by more than one platform are renamed per platform

// Declare platform functions renamed with a suffix
#define DECLARE_PLATFORM_FUNCTIONS(suffix) \
    bool WindowShouldClose##suffix(void); \
    void ToggleFullscreen##suffix(void); \
    void ToggleBorderlessWindowed##suffix(void); \
    void MaximizeWindow##suffix(void); \
    void MinimizeWindow##suffix(void); \
    void RestoreWindow##suffix(void); \
    void SetWindowState##suffix(unsigned int flags); \
    void ClearWindowState##suffix(unsigned int flags); \
    void SetWindowIcon##suffix(Image image); \
    void SetWindowIcons##suffix(Image *images, int count); \
    void SetWindowTitle##suffix(const char *title); \
    void SetWindowPosition##suffix(int x, int y); \
    void SetWindowMonitor##suffix(int monitor); \
    void SetWindowMinSize##suffix(int width, int height); \
    void SetWindowMaxSize##suffix(int width, int height); \
    void SetWindowSize##suffix(int width, int height); \
    void SetWindowOpacity##suffix(float opacity); \
    void SetWindowFocused##suffix(void); \
    void *GetWindowHandle##suffix(void); \
    int GetMonitorCount##suffix(void); \
    int GetCurrentMonitor##suffix(void); \
    Vector2 GetMonitorPosition##suffix(int monitor); \
    int GetMonitorWidth##suffix(int monitor); \
    int GetMonitorHeight##suffix(int monitor); \
    int GetMonitorPhysicalWidth##suffix(int monitor); \
    int GetMonitorPhysicalHeight##suffix(int monitor); \
    int GetMonitorRefreshRate##suffix(int monitor); \
    const char *GetMonitorName##suffix(int monitor); \
    Vector2 GetWindowPosition##suffix(void); \
    Vector2 GetWindowScaleDPI##suffix(void); \
    void SetClipboardText##suffix(const char *text); \
    const char *GetClipboardText##suffix(void); \
    Image GetClipboardImage##suffix(void); \
    void ShowCursor##suffix(void); \
    void HideCursor##suffix(void); \
    void EnableCursor##suffix(void); \
    void DisableCursor##suffix(void); \
    void SwapScreenBuffer##suffix(void); \
    double GetTime##suffix(void); \
    void OpenURL##suffix(const char *url); \
    int SetGamepadMappings##suffix(const char *mappings); \
    void SetGamepadVibration##suffix(int gamepad, float leftMotor, float rightMotor, float duration); \
    void SetMousePosition##suffix(int x, int y); \
    void SetMouseCursor##suffix(int cursor); \
    const char *GetKeyName##suffix(int key); \
    void PollInputEvents##suffix(void); \
    int InitPlatform##suffix(void); \
    void ClosePlatform##suffix(void);

DECLARE_PLATFORM_FUNCTIONS(Drm)
DECLARE_PLATFORM_FUNCTIONS(Wayland)
DECLARE_PLATFORM_FUNCTIONS(X11)

// DRM platform
#define WindowShouldClose WindowShouldCloseDrm
#define ToggleFullscreen ToggleFullscreenDrm
#define ToggleBorderlessWindowed ToggleBorderlessWindowedDrm
#define MaximizeWindow MaximizeWindowDrm
#define MinimizeWindow MinimizeWindowDrm
#define RestoreWindow RestoreWindowDrm
#define SetWindowState SetWindowStateDrm
#define ClearWindowState ClearWindowStateDrm
#define SetWindowIcon SetWindowIconDrm
#define SetWindowIcons SetWindowIconsDrm
#define SetWindowTitle SetWindowTitleDrm
#define SetWindowPosition SetWindowPositionDrm
#define SetWindowMonitor SetWindowMonitorDrm
#define SetWindowMinSize SetWindowMinSizeDrm
#define SetWindowMaxSize SetWindowMaxSizeDrm
#define SetWindowSize SetWindowSizeDrm
#define SetWindowOpacity SetWindowOpacityDrm
#define SetWindowFocused SetWindowFocusedDrm
#define GetWindowHandle GetWindowHandleDrm
#define GetMonitorCount GetMonitorCountDrm
#define GetCurrentMonitor GetCurrentMonitorDrm
#define GetMonitorPosition GetMonitorPositionDrm
#define GetMonitorWidth GetMonitorWidthDrm
#define GetMonitorHeight GetMonitorHeightDrm
#define GetMonitorPhysicalWidth GetMonitorPhysicalWidthDrm
#define GetMonitorPhysicalHeight GetMonitorPhysicalHeightDrm
#define GetMonitorRefreshRate GetMonitorRefreshRateDrm
#define GetMonitorName GetMonitorNameDrm
#define GetWindowPosition GetWindowPositionDrm
#define GetWindowScaleDPI GetWindowScaleDPIDrm
#define SetClipboardText SetClipboardTextDrm
#define GetClipboardText GetClipboardTextDrm
#define GetClipboardImage GetClipboardImageDrm
#define ShowCursor ShowCursorDrm
#define HideCursor HideCursorDrm
#define EnableCursor EnableCursorDrm
#define DisableCursor DisableCursorDrm
#define SwapScreenBuffer SwapScreenBufferDrm
#define GetTime GetTimeDrm
#define OpenURL OpenURLDrm
#define SetGamepadMappings SetGamepadMappingsDrm
#define SetGamepadVibration SetGamepadVibrationDrm
#define SetMousePosition SetMousePositionDrm
#define SetMouseCursor SetMouseCursorDrm
#define GetKeyName GetKeyNameDrm
#define PollInputEvents PollInputEventsDrm
#define InitPlatform InitPlatformDrm
#define ClosePlatform ClosePlatformDrm
#define CloseLibraries CloseLibrariesDrm
#define InitKeyboard InitKeyboardDrm
#define LoadLibraries LoadLibrariesDrm
#define PlatformData PlatformDataDrm
#define linuxToRaylibMap linuxToRaylibMapDrm
#define platform platformDrm

#include "rcore_drm.c"

#undef WindowShouldClose
#undef ToggleFullscreen
#undef ToggleBorderlessWindowed
#undef MaximizeWindow
#undef MinimizeWindow
#undef RestoreWindow
#undef SetWindowState
#undef ClearWindowState
#undef SetWindowIcon
#undef SetWindowIcons
#undef SetWindowTitle
#undef SetWindowPosition
#undef SetWindowMonitor
#undef SetWindowMinSize
#undef SetWindowMaxSize
#undef SetWindowSize
#undef SetWindowOpacity
#undef SetWindowFocused
#undef GetWindowHandle
#undef GetMonitorCount
#undef GetCurrentMonitor
#undef GetMonitorPosition
#undef GetMonitorWidth
#undef GetMonitorHeight
#undef GetMonitorPhysicalWidth
#undef GetMonitorPhysicalHeight
#undef GetMonitorRefreshRate
#undef GetMonitorName
#undef GetWindowPosition
#undef GetWindowScaleDPI
#undef SetClipboardText
#undef GetClipboardText
#undef GetClipboardImage
#undef ShowCursor
#undef HideCursor
#undef EnableCursor
#undef DisableCursor
#undef SwapScreenBuffer
#undef GetTime
#undef OpenURL
#undef SetGamepadMappings
#undef SetGamepadVibration
#undef SetMousePosition
#undef SetMouseCursor
#undef GetKeyName
#undef PollInputEvents
#undef InitPlatform
#undef ClosePlatform
#undef CloseLibraries
#undef InitKeyboard
#undef LoadLibraries
#undef PlatformData
#undef linuxToRaylibMap
#undef platform
#undef KEYMAP_SIZE
#undef LIBEGL_FUNCTIONS
#undef LIBRARY_FUNCTION
#undef LOAD_FUNCTION
#undef eglBindAPI
#undef eglChooseConfig
#undef eglCreateContext
#undef eglCreateWindowSurface
#undef eglDestroyContext
#undef eglDestroySurface
#undef eglGetConfigAttrib
#undef eglGetDisplay
#undef eglGetProcAddress
#undef eglInitialize
#undef eglMakeCurrent
#undef eglSwapBuffers
#undef eglTerminate

// Wayland platform
#define WindowShouldClose WindowShouldCloseWayland
#define ToggleFullscreen ToggleFullscreenWayland
#define ToggleBorderlessWindowed ToggleBorderlessWindowedWayland
#define MaximizeWindow MaximizeWindowWayland
#define MinimizeWindow MinimizeWindowWayland
#define RestoreWindow RestoreWindowWayland
#define SetWindowState SetWindowStateWayland
#define ClearWindowState ClearWindowStateWayland
#define SetWindowIcon SetWindowIconWayland
#define SetWindowIcons SetWindowIconsWayland
#define SetWindowTitle SetWindowTitleWayland
#define SetWindowPosition SetWindowPositionWayland
#define SetWindowMonitor SetWindowMonitorWayland
#define SetWindowMinSize SetWindowMinSizeWayland
#define SetWindowMaxSize SetWindowMaxSizeWayland
#define SetWindowSize SetWindowSizeWayland
#define SetWindowOpacity SetWindowOpacityWayland
#define SetWindowFocused SetWindowFocusedWayland
#define GetWindowHandle GetWindowHandleWayland
#define GetMonitorCount GetMonitorCountWayland
#define GetCurrentMonitor GetCurrentMonitorWayland
#define GetMonitorPosition GetMonitorPositionWayland
#define GetMonitorWidth GetMonitorWidthWayland
#define GetMonitorHeight GetMonitorHeightWayland
#define GetMonitorPhysicalWidth GetMonitorPhysicalWidthWayland
#define GetMonitorPhysicalHeight GetMonitorPhysicalHeightWayland
#define GetMonitorRefreshRate GetMonitorRefreshRateWayland
#define GetMonitorName GetMonitorNameWayland
#define GetWindowPosition GetWindowPositionWayland
#define GetWindowScaleDPI GetWindowScaleDPIWayland
#define SetClipboardText SetClipboardTextWayland
#define GetClipboardText GetClipboardTextWayland
#define GetClipboardImage GetClipboardImageWayland
#define ShowCursor ShowCursorWayland
#define HideCursor HideCursorWayland
#define EnableCursor EnableCursorWayland
#define DisableCursor DisableCursorWayland
#define SwapScreenBuffer SwapScreenBufferWayland
#define GetTime GetTimeWayland
#define OpenURL OpenURLWayland
#define SetGamepadMappings SetGamepadMappingsWayland
#define SetGamepadVibration SetGamepadVibrationWayland
#define SetMousePosition SetMousePositionWayland
#define SetMouseCursor SetMouseCursorWayland
#define GetKeyName GetKeyNameWayland
#define PollInputEvents PollInputEventsWayland
#define InitPlatform InitPlatformWayland
#define ClosePlatform ClosePlatformWayland
#define CloseLibraries CloseLibrariesWayland
#define HandleKey HandleKeyWayland
#define HandleKeyChar HandleKeyCharWayland
#define HandleMouseButton HandleMouseButtonWayland
#define HandleMouseMove HandleMouseMoveWayland
#define HandleTouch HandleTouchWayland
#define InitGraphicsContext InitGraphicsContextWayland
#define LoadLibraries LoadLibrariesWayland
#define LockPointer LockPointerWayland
#define MonitorInfo MonitorInfoWayland
#define PlatformData PlatformDataWayland
#define ProcessDroppedUris ProcessDroppedUrisWayland
#define RequestFullscreen RequestFullscreenWayland
#define UpdateCursor UpdateCursorWayland
#define UpdateDecorations UpdateDecorationsWayland
#define UpdateWindowBuffers UpdateWindowBuffersWayland
#define UpdateWindowScale UpdateWindowScaleWayland
#define linuxToRaylibMap linuxToRaylibMapWayland
#define mgAxisConvertTable mgAxisConvertTableWayland
#define mgButtonConvertTable mgButtonConvertTableWayland
#define platform platformWayland

#include "rcore_desktop_wayland.c"

#undef WindowShouldClose
#undef ToggleFullscreen
#undef ToggleBorderlessWindowed
#undef MaximizeWindow
#undef MinimizeWindow
#undef RestoreWindow
#undef SetWindowState
#undef ClearWindowState
#undef SetWindowIcon
#undef SetWindowIcons
#undef SetWindowTitle
#undef SetWindowPosition
#undef SetWindowMonitor
#undef SetWindowMinSize
#undef SetWindowMaxSize
#undef SetWindowSize
#undef SetWindowOpacity
#undef SetWindowFocused
#undef GetWindowHandle
#undef GetMonitorCount
#undef GetCurrentMonitor
#undef GetMonitorPosition
#undef GetMonitorWidth
#undef GetMonitorHeight
#undef GetMonitorPhysicalWidth
#undef GetMonitorPhysicalHeight
#undef GetMonitorRefreshRate
#undef GetMonitorName
#undef GetWindowPosition
#undef GetWindowScaleDPI
#undef SetClipboardText
#undef GetClipboardText
#undef GetClipboardImage
#undef ShowCursor
#undef HideCursor
#undef EnableCursor
#undef DisableCursor
#undef SwapScreenBuffer
#undef GetTime
#undef OpenURL
#undef SetGamepadMappings
#undef SetGamepadVibration
#undef SetMousePosition
#undef SetMouseCursor
#undef GetKeyName
#undef PollInputEvents
#undef InitPlatform
#undef ClosePlatform
#undef CloseLibraries
#undef HandleKey
#undef HandleKeyChar
#undef HandleMouseButton
#undef HandleMouseMove
#undef HandleTouch
#undef InitGraphicsContext
#undef LoadLibraries
#undef LockPointer
#undef MonitorInfo
#undef PlatformData
#undef ProcessDroppedUris
#undef RequestFullscreen
#undef UpdateCursor
#undef UpdateDecorations
#undef UpdateWindowBuffers
#undef UpdateWindowScale
#undef linuxToRaylibMap
#undef mgAxisConvertTable
#undef mgButtonConvertTable
#undef platform
#undef KEYMAP_SIZE
#undef LIBEGL_FUNCTIONS
#undef LIBRARY_FUNCTION
#undef LIBXKBCOMMON_FUNCTIONS
#undef MAX_MONITORS
#undef MAX_WINDOW_ICONS
#undef eglBindAPI
#undef eglChooseConfig
#undef eglCreateContext
#undef eglCreateWindowSurface
#undef eglDestroyContext
#undef eglDestroySurface
#undef eglGetConfigAttrib
#undef eglGetDisplay
#undef eglGetProcAddress
#undef eglInitialize
#undef eglMakeCurrent
#undef eglSwapBuffers
#undef eglSwapInterval
#undef eglTerminate
#undef xkb_keysym_to_utf32
#undef xkb_keysym_to_utf8

// X11 platform
#define WindowShouldClose WindowShouldCloseX11
#define ToggleFullscreen ToggleFullscreenX11
#define ToggleBorderlessWindowed ToggleBorderlessWindowedX11
#define MaximizeWindow MaximizeWindowX11
#define MinimizeWindow MinimizeWindowX11
#define RestoreWindow RestoreWindowX11
#define SetWindowState SetWindowStateX11
#define ClearWindowState ClearWindowStateX11
#define SetWindowIcon SetWindowIconX11
#define SetWindowIcons SetWindowIconsX11
#define SetWindowTitle SetWindowTitleX11
#define SetWindowPosition SetWindowPositionX11
#define SetWindowMonitor SetWindowMonitorX11
#define SetWindowMinSize SetWindowMinSizeX11
#define SetWindowMaxSize SetWindowMaxSizeX11
#define SetWindowSize SetWindowSizeX11
#define SetWindowOpacity SetWindowOpacityX11
#define SetWindowFocused SetWindowFocusedX11
#define GetWindowHandle GetWindowHandleX11
#define GetMonitorCount GetMonitorCountX11
#define GetCurrentMonitor GetCurrentMonitorX11
#define GetMonitorPosition GetMonitorPositionX11
#define GetMonitorWidth GetMonitorWidthX11
#define GetMonitorHeight GetMonitorHeightX11
#define GetMonitorPhysicalWidth GetMonitorPhysicalWidthX11
#define GetMonitorPhysicalHeight GetMonitorPhysicalHeightX11
#define GetMonitorRefreshRate GetMonitorRefreshRateX11
#define GetMonitorName GetMonitorNameX11
#define GetWindowPosition GetWindowPositionX11
#define GetWindowScaleDPI GetWindowScaleDPIX11
#define SetClipboardText SetClipboardTextX11
#define GetClipboardText GetClipboardTextX11
#define GetClipboardImage GetClipboardImageX11
#define ShowCursor ShowCursorX11
#define HideCursor HideCursorX11
#define EnableCursor EnableCursorX11
#define DisableCursor DisableCursorX11
#define SwapScreenBuffer SwapScreenBufferX11
#define GetTime GetTimeX11
#define OpenURL OpenURLX11
#define SetGamepadMappings SetGamepadMappingsX11
#define SetGamepadVibration SetGamepadVibrationX11
#define SetMousePosition SetMousePositionX11
#define SetMouseCursor SetMouseCursorX11
#define GetKeyName GetKeyNameX11
#define PollInputEvents PollInputEventsX11
#define InitPlatform InitPlatformX11
#define ClosePlatform ClosePlatformX11
#define CloseLibraries CloseLibrariesX11
#define HandleKey HandleKeyX11
#define HandleKeyChar HandleKeyCharX11
#define HandleMouseButton HandleMouseButtonX11
#define HandleMouseMove HandleMouseMoveX11
#define HandleTouch HandleTouchX11
#define InitGraphicsContext InitGraphicsContextX11
#define InitKeyboard InitKeyboardX11
#define LoadLibraries LoadLibrariesX11
#define LockPointer LockPointerX11
#define MonitorInfo MonitorInfoX11
#define PlatformData PlatformDataX11
#define ProcessDroppedUris ProcessDroppedUrisX11
#define RequestFullscreen RequestFullscreenX11
#define UpdateCursor UpdateCursorX11
#define UpdateDecorations UpdateDecorationsX11
#define UpdateWindowBuffers UpdateWindowBuffersX11
#define UpdateWindowScale UpdateWindowScaleX11
#define mgAxisConvertTable mgAxisConvertTableX11
#define mgButtonConvertTable mgButtonConvertTableX11
#define platform platformX11

#include "rcore_desktop_x11.c"

#undef WindowShouldClose
#undef ToggleFullscreen
#undef ToggleBorderlessWindowed
#undef MaximizeWindow
#undef MinimizeWindow
#undef RestoreWindow
#undef SetWindowState
#undef ClearWindowState
#undef SetWindowIcon
#undef SetWindowIcons
#undef SetWindowTitle
#undef SetWindowPosition
#undef SetWindowMonitor
#undef SetWindowMinSize
#undef SetWindowMaxSize
#undef SetWindowSize
#undef SetWindowOpacity
#undef SetWindowFocused
#undef GetWindowHandle
#undef GetMonitorCount
#undef GetCurrentMonitor
#undef GetMonitorPosition
#undef GetMonitorWidth
#undef GetMonitorHeight
#undef GetMonitorPhysicalWidth
#undef GetMonitorPhysicalHeight
#undef GetMonitorRefreshRate
#undef GetMonitorName
#undef GetWindowPosition
#undef GetWindowScaleDPI
#undef SetClipboardText
#undef GetClipboardText
#undef GetClipboardImage
#undef ShowCursor
#undef HideCursor
#undef EnableCursor
#undef DisableCursor
#undef SwapScreenBuffer
#undef GetTime
#undef OpenURL
#undef SetGamepadMappings
#undef SetGamepadVibration
#undef SetMousePosition
#undef SetMouseCursor
#undef GetKeyName
#undef PollInputEvents
#undef InitPlatform
#undef ClosePlatform
#undef CloseLibraries
#undef HandleKey
#undef HandleKeyChar
#undef HandleMouseButton
#undef HandleMouseMove
#undef HandleTouch
#undef InitGraphicsContext
#undef InitKeyboard
#undef LoadLibraries
#undef LockPointer
#undef MonitorInfo
#undef PlatformData
#undef ProcessDroppedUris
#undef RequestFullscreen
#undef UpdateCursor
#undef UpdateDecorations
#undef UpdateWindowBuffers
#undef UpdateWindowScale
#undef mgAxisConvertTable
#undef mgButtonConvertTable
#undef platform
#undef KEYMAP_SIZE
#undef LIBEGL_FUNCTIONS
#undef LIBRARY_FUNCTION
#undef LIBXKBCOMMON_FUNCTIONS
#undef LOAD_FUNCTION
#undef MAX_MONITORS
#undef MAX_WINDOW_ICONS
#undef eglBindAPI
#undef eglChooseConfig
#undef eglCreateContext
#undef eglCreateWindowSurface
#undef eglDestroyContext
#undef eglDestroySurface
#undef eglGetConfigAttrib
#undef eglGetDisplay
#undef eglGetProcAddress
#undef eglInitialize
#undef eglMakeCurrent
#undef eglSwapBuffers
#undef eglSwapInterval
#undef eglTerminate
#undef xkb_keysym_to_utf32
#undef xkb_keysym_to_utf8

//----------------------------------------------------------------------------------
// Types and Structures Definition
//----------------------------------------------------------------------------------
typedef enum {
    LINUX_PLATFORM_NONE = 0,        // Platform not selected yet
    LINUX_PLATFORM_WAYLAND,         // Wayland platform
    LINUX_PLATFORM_X11,             // X11 platform
    LINUX_PLATFORM_DRM              // DRM platform
} LinuxPlatform;

//----------------------------------------------------------------------------------
// Global Variables Definition
//----------------------------------------------------------------------------------
static LinuxPlatform linuxPlatform = LINUX_PLATFORM_NONE;  // Platform in use

//----------------------------------------------------------------------------------
// Module Internal Functions Declaration
//----------------------------------------------------------------------------------
int InitPlatform(void);                                     // Initialize platform (graphics, inputs and more)
void ClosePlatform(void);                                   // Close platform

static LinuxPlatform GetRequestedPlatform(bool *forced);    // Get platform requested by RAYLIB_PLATFORM or the session
static LinuxPlatform GetPlatform(void);                     // Get platform in use, selecting the platform if required
static int InitLinuxPlatform(LinuxPlatform platform);       // Initialize a platform

//----------------------------------------------------------------------------------
// Module Functions Definition
//----------------------------------------------------------------------------------

// Check if application should close
bool WindowShouldClose(void)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) return WindowShouldCloseWayland();
    else if (platform == LINUX_PLATFORM_X11) return WindowShouldCloseX11();
    else return WindowShouldCloseDrm();
}

// Toggle fullscreen mode
void ToggleFullscreen(void)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) ToggleFullscreenWayland();
    else if (platform == LINUX_PLATFORM_X11) ToggleFullscreenX11();
    else ToggleFullscreenDrm();
}

// Toggle borderless windowed mode
void ToggleBorderlessWindowed(void)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) ToggleBorderlessWindowedWayland();
    else if (platform == LINUX_PLATFORM_X11) ToggleBorderlessWindowedX11();
    else ToggleBorderlessWindowedDrm();
}

// Set window state: maximized, if resizable
void MaximizeWindow(void)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) MaximizeWindowWayland();
    else if (platform == LINUX_PLATFORM_X11) MaximizeWindowX11();
    else MaximizeWindowDrm();
}

// Set window state: minimized
void MinimizeWindow(void)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) MinimizeWindowWayland();
    else if (platform == LINUX_PLATFORM_X11) MinimizeWindowX11();
    else MinimizeWindowDrm();
}

// Restore window from being minimized/maximized
void RestoreWindow(void)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) RestoreWindowWayland();
    else if (platform == LINUX_PLATFORM_X11) RestoreWindowX11();
    else RestoreWindowDrm();
}

// Set window configuration state using flags
void SetWindowState(unsigned int flags)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) SetWindowStateWayland(flags);
    else if (platform == LINUX_PLATFORM_X11) SetWindowStateX11(flags);
    else SetWindowStateDrm(flags);
}

// Clear window configuration state flags
void ClearWindowState(unsigned int flags)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) ClearWindowStateWayland(flags);
    else if (platform == LINUX_PLATFORM_X11) ClearWindowStateX11(flags);
    else ClearWindowStateDrm(flags);
}

// Set icon for window
void SetWindowIcon(Image image)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) SetWindowIconWayland(image);
    else if (platform == LINUX_PLATFORM_X11) SetWindowIconX11(image);
    else SetWindowIconDrm(image);
}

// Set icon for window, multiple images
void SetWindowIcons(Image *images, int count)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) SetWindowIconsWayland(images, count);
    else if (platform == LINUX_PLATFORM_X11) SetWindowIconsX11(images, count);
    else SetWindowIconsDrm(images, count);
}

// Set title for window
void SetWindowTitle(const char *title)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) SetWindowTitleWayland(title);
    else if (platform == LINUX_PLATFORM_X11) SetWindowTitleX11(title);
    else SetWindowTitleDrm(title);
}

// Set window position on screen (windowed mode)
void SetWindowPosition(int x, int y)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) SetWindowPositionWayland(x, y);
    else if (platform == LINUX_PLATFORM_X11) SetWindowPositionX11(x, y);
    else SetWindowPositionDrm(x, y);
}

// Set monitor for the current window
void SetWindowMonitor(int monitor)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) SetWindowMonitorWayland(monitor);
    else if (platform == LINUX_PLATFORM_X11) SetWindowMonitorX11(monitor);
    else SetWindowMonitorDrm(monitor);
}

// Set window minimum dimensions (FLAG_WINDOW_RESIZABLE)
void SetWindowMinSize(int width, int height)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) SetWindowMinSizeWayland(width, height);
    else if (platform == LINUX_PLATFORM_X11) SetWindowMinSizeX11(width, height);
    else SetWindowMinSizeDrm(width, height);
}

// Set window maximum dimensions (FLAG_WINDOW_RESIZABLE)
void SetWindowMaxSize(int width, int height)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) SetWindowMaxSizeWayland(width, height);
    else if (platform == LINUX_PLATFORM_X11) SetWindowMaxSizeX11(width, height);
    else SetWindowMaxSizeDrm(width, height);
}

// Set window dimensions
void SetWindowSize(int width, int height)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) SetWindowSizeWayland(width, height);
    else if (platform == LINUX_PLATFORM_X11) SetWindowSizeX11(width, height);
    else SetWindowSizeDrm(width, height);
}

// Set window opacity, value opacity is between 0.0 and 1.0
void SetWindowOpacity(float opacity)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) SetWindowOpacityWayland(opacity);
    else if (platform == LINUX_PLATFORM_X11) SetWindowOpacityX11(opacity);
    else SetWindowOpacityDrm(opacity);
}

// Set window focused
void SetWindowFocused(void)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) SetWindowFocusedWayland();
    else if (platform == LINUX_PLATFORM_X11) SetWindowFocusedX11();
    else SetWindowFocusedDrm();
}

// Get native window handle
void *GetWindowHandle(void)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) return GetWindowHandleWayland();
    else if (platform == LINUX_PLATFORM_X11) return GetWindowHandleX11();
    else return GetWindowHandleDrm();
}

// Get number of monitors
int GetMonitorCount(void)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) return GetMonitorCountWayland();
    else if (platform == LINUX_PLATFORM_X11) return GetMonitorCountX11();
    else return GetMonitorCountDrm();
}

// Get current monitor where window is placed
int GetCurrentMonitor(void)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) return GetCurrentMonitorWayland();
    else if (platform == LINUX_PLATFORM_X11) return GetCurrentMonitorX11();
    else return GetCurrentMonitorDrm();
}

// Get selected monitor position
Vector2 GetMonitorPosition(int monitor)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) return GetMonitorPositionWayland(monitor);
    else if (platform == LINUX_PLATFORM_X11) return GetMonitorPositionX11(monitor);
    else return GetMonitorPositionDrm(monitor);
}

// Get selected monitor width (currently used by monitor)
int GetMonitorWidth(int monitor)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) return GetMonitorWidthWayland(monitor);
    else if (platform == LINUX_PLATFORM_X11) return GetMonitorWidthX11(monitor);
    else return GetMonitorWidthDrm(monitor);
}

// Get selected monitor height (currently used by monitor)
int GetMonitorHeight(int monitor)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) return GetMonitorHeightWayland(monitor);
    else if (platform == LINUX_PLATFORM_X11) return GetMonitorHeightX11(monitor);
    else return GetMonitorHeightDrm(monitor);
}

// Get selected monitor physical width in millimetres
int GetMonitorPhysicalWidth(int monitor)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) return GetMonitorPhysicalWidthWayland(monitor);
    else if (platform == LINUX_PLATFORM_X11) return GetMonitorPhysicalWidthX11(monitor);
    else return GetMonitorPhysicalWidthDrm(monitor);
}

// Get selected monitor physical height in millimetres
int GetMonitorPhysicalHeight(int monitor)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) return GetMonitorPhysicalHeightWayland(monitor);
    else if (platform == LINUX_PLATFORM_X11) return GetMonitorPhysicalHeightX11(monitor);
    else return GetMonitorPhysicalHeightDrm(monitor);
}

// Get selected monitor refresh rate
int GetMonitorRefreshRate(int monitor)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) return GetMonitorRefreshRateWayland(monitor);
    else if (platform == LINUX_PLATFORM_X11) return GetMonitorRefreshRateX11(monitor);
    else return GetMonitorRefreshRateDrm(monitor);
}

// Get the human-readable, UTF-8 encoded name of the selected monitor
const char *GetMonitorName(int monitor)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) return GetMonitorNameWayland(monitor);
    else if (platform == LINUX_PLATFORM_X11) return GetMonitorNameX11(monitor);
    else return GetMonitorNameDrm(monitor);
}

// Get window position XY on monitor
Vector2 GetWindowPosition(void)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) return GetWindowPositionWayland();
    else if (platform == LINUX_PLATFORM_X11) return GetWindowPositionX11();
    else return GetWindowPositionDrm();
}

// Get window scale DPI factor for current monitor
Vector2 GetWindowScaleDPI(void)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) return GetWindowScaleDPIWayland();
    else if (platform == LINUX_PLATFORM_X11) return GetWindowScaleDPIX11();
    else return GetWindowScaleDPIDrm();
}

// Set clipboard text content
void SetClipboardText(const char *text)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) SetClipboardTextWayland(text);
    else if (platform == LINUX_PLATFORM_X11) SetClipboardTextX11(text);
    else SetClipboardTextDrm(text);
}

// Get clipboard text content
const char *GetClipboardText(void)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) return GetClipboardTextWayland();
    else if (platform == LINUX_PLATFORM_X11) return GetClipboardTextX11();
    else return GetClipboardTextDrm();
}

// Get clipboard image
Image GetClipboardImage(void)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) return GetClipboardImageWayland();
    else if (platform == LINUX_PLATFORM_X11) return GetClipboardImageX11();
    else return GetClipboardImageDrm();
}

// Show mouse cursor
void ShowCursor(void)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) ShowCursorWayland();
    else if (platform == LINUX_PLATFORM_X11) ShowCursorX11();
    else ShowCursorDrm();
}

// Hide mouse cursor
void HideCursor(void)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) HideCursorWayland();
    else if (platform == LINUX_PLATFORM_X11) HideCursorX11();
    else HideCursorDrm();
}

// Enable cursor (unlock cursor)
void EnableCursor(void)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) EnableCursorWayland();
    else if (platform == LINUX_PLATFORM_X11) EnableCursorX11();
    else EnableCursorDrm();
}

// Disable cursor (lock cursor)
void DisableCursor(void)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) DisableCursorWayland();
    else if (platform == LINUX_PLATFORM_X11) DisableCursorX11();
    else DisableCursorDrm();
}

// Swap back buffer with front buffer (screen drawing)
void SwapScreenBuffer(void)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) SwapScreenBufferWayland();
    else if (platform == LINUX_PLATFORM_X11) SwapScreenBufferX11();
    else SwapScreenBufferDrm();
}

// Get elapsed time measure in seconds since InitTimer()
double GetTime(void)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) return GetTimeWayland();
    else if (platform == LINUX_PLATFORM_X11) return GetTimeX11();
    else return GetTimeDrm();
}

// Open URL with default system browser (if available)
void OpenURL(const char *url)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) OpenURLWayland(url);
    else if (platform == LINUX_PLATFORM_X11) OpenURLX11(url);
    else OpenURLDrm(url);
}

// Set internal gamepad mappings
int SetGamepadMappings(const char *mappings)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) return SetGamepadMappingsWayland(mappings);
    else if (platform == LINUX_PLATFORM_X11) return SetGamepadMappingsX11(mappings);
    else return SetGamepadMappingsDrm(mappings);
}

// Set gamepad vibration
void SetGamepadVibration(int gamepad, float leftMotor, float rightMotor, float duration)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) SetGamepadVibrationWayland(gamepad, leftMotor, rightMotor, duration);
    else if (platform == LINUX_PLATFORM_X11) SetGamepadVibrationX11(gamepad, leftMotor, rightMotor, duration);
    else SetGamepadVibrationDrm(gamepad, leftMotor, rightMotor, duration);
}

// Set mouse position XY
void SetMousePosition(int x, int y)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) SetMousePositionWayland(x, y);
    else if (platform == LINUX_PLATFORM_X11) SetMousePositionX11(x, y);
    else SetMousePositionDrm(x, y);
}

// Set mouse cursor
void SetMouseCursor(int cursor)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) SetMouseCursorWayland(cursor);
    else if (platform == LINUX_PLATFORM_X11) SetMouseCursorX11(cursor);
    else SetMouseCursorDrm(cursor);
}

// Get physical key name
const char *GetKeyName(int key)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) return GetKeyNameWayland(key);
    else if (platform == LINUX_PLATFORM_X11) return GetKeyNameX11(key);
    else return GetKeyNameDrm(key);
}

// Register all input events
void PollInputEvents(void)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) PollInputEventsWayland();
    else if (platform == LINUX_PLATFORM_X11) PollInputEventsX11();
    else PollInputEventsDrm();
}

// Initialize platform: graphics, inputs and more
// NOTE: A platform that fails to initialize falls back to the next one, unless requested by RAYLIB_PLATFORM
int InitPlatform(void)
{
    bool forced = false;
    LinuxPlatform requested = GetRequestedPlatform(&forced);

    for (int platform = requested; platform <= LINUX_PLATFORM_DRM; platform++)
    {
        linuxPlatform = (LinuxPlatform)platform;
        if (InitLinuxPlatform(linuxPlatform) == 0) return 0;
        if (forced) break;
    }

    linuxPlatform = LINUX_PLATFORM_NONE;

    return -1;
}

// Close platform
void ClosePlatform(void)
{
    LinuxPlatform platform = GetPlatform();

    if (platform == LINUX_PLATFORM_WAYLAND) ClosePlatformWayland();
    else if (platform == LINUX_PLATFORM_X11) ClosePlatformX11();
    else ClosePlatformDrm();

    linuxPlatform = LINUX_PLATFORM_NONE;
}

//----------------------------------------------------------------------------------
// Module Internal Functions Definition
//----------------------------------------------------------------------------------

// Get platform requested by RAYLIB_PLATFORM or the session
static LinuxPlatform GetRequestedPlatform(bool *forced)
{
    const char *requested = getenv("RAYLIB_PLATFORM");
    *forced = true;

    if ((requested != NULL) && (strcmp(requested, "wayland") == 0)) return LINUX_PLATFORM_WAYLAND;
    if ((requested != NULL) && (strcmp(requested, "x11") == 0)) return LINUX_PLATFORM_X11;
    if ((requested != NULL) && (strcmp(requested, "drm") == 0)) return LINUX_PLATFORM_DRM;

    *forced = false;

    if (getenv("WAYLAND_DISPLAY") != NULL) return LINUX_PLATFORM_WAYLAND;
    if (getenv("DISPLAY") != NULL) return LINUX_PLATFORM_X11;

    return LINUX_PLATFORM_DRM;
}

// Get platform in use, selecting the platform if required
// NOTE: Functions called before InitWindow() use the requested platform
static LinuxPlatform GetPlatform(void)
{
    if (linuxPlatform == LINUX_PLATFORM_NONE)
    {
        bool forced = false;
        linuxPlatform = GetRequestedPlatform(&forced);
    }

    return linuxPlatform;
}

// Initialize a platform
static int InitLinuxPlatform(LinuxPlatform platform)
{
    int result = -1;

    if (platform == LINUX_PLATFORM_WAYLAND) result = InitPlatformWayland();
    else if (platform == LINUX_PLATFORM_X11) result = InitPlatformX11();
    else if (platform == LINUX_PLATFORM_DRM) result = InitPlatformDrm();

    if ((result != 0) && (platform != LINUX_PLATFORM_DRM)) TRACELOG(LOG_WARNING, "PLATFORM: %s not available", (platform == LINUX_PLATFORM_WAYLAND)? "Wayland" : "X11");

    return result;
}
