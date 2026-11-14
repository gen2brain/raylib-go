/**********************************************************************************************
*
*   rcore_desktop_cocoa - Functions to manage window, graphics device and inputs
*
*   PLATFORM: DESKTOP: COCOA
*       - macOS (x64, arm64)
*
*   LIMITATIONS:
*       - OpenGL 4.3 not available (macOS provides OpenGL up to 4.1 core)
*       - SetGamepadVibration() not available (minigamepad has no rumble support)
*       - Frames stop while the window is resized interactively (AppKit live resize loop)
*
*   POSSIBLE IMPROVEMENTS:
*       - Redraw during the AppKit live resize loop
*
*   CONFIGURATION:
*       #define RCORE_PLATFORM_CUSTOM_FLAG
*           Custom flag for rcore on target platform -not used-
*
*   DEPENDENCIES:
*       - AppKit, QuartzCore: Window, events, software framebuffer presentation
*       - OpenGL framework: NSOpenGLContext (deprecated by Apple, still available)
*       - ANGLE (libEGL.dylib, libGLESv2.dylib, loaded at runtime): OpenGL ES on Metal
*       - IOKit: Gamepads (through minigamepad)
*       - gestures: Gestures system for touch-ready devices (or simulated from mouse inputs)
*
*   NOTE: This file contains Objective-C code, rcore.c must be compiled with -x objective-c
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

// Move MacTypes.h symbols to new names to avoid redefining the same names as rcore
#define Point PointMacOS
#define Size SizeMacOS

#import <Cocoa/Cocoa.h>             // AppKit and Foundation
#import <QuartzCore/QuartzCore.h>   // Required for: CALayer, CATransaction

#undef Point                        // rcore symbol collision
#undef Size                         // rcore symbol collision

#include <dlfcn.h>                  // Required for: dlopen(), dlsym(), dlclose()
#include <time.h>                   // Required for: clock_gettime()
#include <unistd.h>                 // Required for: chdir(), usleep()
#include <math.h>                   // Required for: fmod(), round()
#include <float.h>                  // Required for: FLT_MAX

// Gamepad support (IOKit HID)
#define MG_MAX_GAMEPADS MAX_GAMEPADS
#define MG_API static inline
#define MG_IMPLEMENTATION
#include "../external/RGFW/deps/minigamepad.h"

//----------------------------------------------------------------------------------
// Defines and Macros
//----------------------------------------------------------------------------------
#define KEYMAP_SIZE                 128     // Size of the macOS virtual keycodes to raylib keys map
#define OCCLUDED_FRAME_RATE         60.0    // Frame rate used while the window is occluded and V-Sync is requested
#define PRECISE_SCROLL_FACTOR       0.1     // Trackpad scroll deltas scale, wheel and trackpad move similarly

// Device-dependent modifier masks, distinguish left and right modifier keys
#define DEVICE_LCTL_KEY_MASK        0x00000001
#define DEVICE_LSHIFT_KEY_MASK      0x00000002
#define DEVICE_RSHIFT_KEY_MASK      0x00000004
#define DEVICE_LCMD_KEY_MASK        0x00000008
#define DEVICE_RCMD_KEY_MASK        0x00000010
#define DEVICE_LALT_KEY_MASK        0x00000020
#define DEVICE_RALT_KEY_MASK        0x00000040
#define DEVICE_RCTL_KEY_MASK        0x00002000

// OpenGL ES through EGL, provided by ANGLE (libEGL.dylib, libGLESv2.dylib), loaded at runtime
#if defined(GRAPHICS_API_OPENGL_ES2) || defined(GRAPHICS_API_OPENGL_ES3)
    #define COCOA_EGL

    #define EGL_DEFAULT_DISPLAY                     ((void *)0)
    #define EGL_NO_CONTEXT                          ((void *)0)
    #define EGL_NO_SURFACE                          ((void *)0)
    #define EGL_NONE                                0x3038
    #define EGL_ALPHA_SIZE                          0x3021
    #define EGL_BLUE_SIZE                           0x3022
    #define EGL_GREEN_SIZE                          0x3023
    #define EGL_RED_SIZE                            0x3024
    #define EGL_DEPTH_SIZE                          0x3025
    #define EGL_STENCIL_SIZE                        0x3026
    #define EGL_SAMPLES                             0x3031
    #define EGL_SAMPLE_BUFFERS                      0x3032
    #define EGL_SURFACE_TYPE                        0x3033
    #define EGL_RENDERABLE_TYPE                     0x3040
    #define EGL_WINDOW_BIT                          0x0004
    #define EGL_OPENGL_ES2_BIT                      0x0004
    #define EGL_OPENGL_ES3_BIT                      0x0040
    #define EGL_CONTEXT_CLIENT_VERSION              0x3098
    #define EGL_OPENGL_ES_API                       0x30A0
    #define EGL_PLATFORM_ANGLE_ANGLE                0x3202
    #define EGL_PLATFORM_ANGLE_TYPE_ANGLE           0x3203
    #define EGL_PLATFORM_ANGLE_TYPE_METAL_ANGLE     0x3489
#endif

//----------------------------------------------------------------------------------
// Types and Structures Definition
//----------------------------------------------------------------------------------
#if defined(COCOA_EGL)
typedef int EGLint;
typedef unsigned int EGLBoolean;
typedef unsigned int EGLenum;
typedef void *EGLDisplay;
typedef void *EGLConfig;
typedef void *EGLSurface;
typedef void *EGLContext;

// EGL functions, loaded from libEGL.dylib
typedef struct {
    void *(*GetProcAddress)(const char *name);
    EGLDisplay (*GetDisplay)(void *displayId);
    EGLDisplay (*GetPlatformDisplayEXT)(EGLenum platform, void *nativeDisplay, const EGLint *attribs);
    EGLBoolean (*Initialize)(EGLDisplay display, EGLint *major, EGLint *minor);
    EGLBoolean (*Terminate)(EGLDisplay display);
    EGLBoolean (*BindAPI)(EGLenum api);
    EGLBoolean (*ChooseConfig)(EGLDisplay display, const EGLint *attribs, EGLConfig *configs, EGLint size, EGLint *count);
    EGLSurface (*CreateWindowSurface)(EGLDisplay display, EGLConfig config, void *window, const EGLint *attribs);
    EGLBoolean (*DestroySurface)(EGLDisplay display, EGLSurface surface);
    EGLContext (*CreateContext)(EGLDisplay display, EGLConfig config, EGLContext share, const EGLint *attribs);
    EGLBoolean (*DestroyContext)(EGLDisplay display, EGLContext context);
    EGLBoolean (*MakeCurrent)(EGLDisplay display, EGLSurface draw, EGLSurface read, EGLContext context);
    EGLBoolean (*SwapBuffers)(EGLDisplay display, EGLSurface surface);
    EGLBoolean (*SwapInterval)(EGLDisplay display, EGLint interval);
    EGLint (*GetError)(void);
} EGLFunctions;
#endif

@interface RaylibApplicationDelegate : NSObject <NSApplicationDelegate>
@end

@interface RaylibWindow : NSWindow
@end

@interface RaylibWindowDelegate : NSObject <NSWindowDelegate>
@end

@interface RaylibView : NSView <NSTextInputClient>
{
    NSTrackingArea *trackingArea;
    NSMutableAttributedString *markedText;
    NSRange markedSelection;
}
@end

typedef struct {
    RaylibApplicationDelegate *appDelegate; // Application delegate (quit requests)
    RaylibWindow *window;               // Native window
    RaylibWindowDelegate *windowDelegate; // Window delegate (window events)
    RaylibView *view;                   // Window content view (input events, OpenGL drawable)
    id keyUpMonitor;                    // Local events monitor, forwards key up events while Command is down

#if defined(COCOA_EGL)
    void *eglLibrary;                   // EGL library handle (ANGLE)
    EGLFunctions egl;                   // EGL functions
    EGLDisplay eglDisplay;              // EGL display
    EGLSurface eglSurface;              // EGL window surface, on the content view layer
    EGLContext eglContext;              // EGL context (OpenGL ES)
#elif !defined(GRAPHICS_API_OPENGL_SOFTWARE)
    NSOpenGLContext *glContext;         // OpenGL context
    void *glLibrary;                    // OpenGL framework handle, used to load GL procedures
#else
    unsigned char *softwarePixels;      // Software renderer framebuffer copy (BGRA)
    int softwareWidth;                  // Software renderer framebuffer width
    int softwareHeight;                 // Software renderer framebuffer height
    CGColorSpaceRef colorSpace;         // Color space for the framebuffer image
#endif

    bool initialized;                   // Platform initialized, viewport can be setup
    bool closeRequested;                // Window close requested (close button, Cmd+Q, exit key)
    bool occluded;                      // Window is not visible (covered, other space, minimized)
    bool cursorSystemHidden;            // System cursor hidden with [NSCursor hide] (disabled cursor)
    bool focusOnShow;                   // Window gets focus when shown (FLAG_WINDOW_UNFOCUSED not requested)

    Vector2 cursorWarpDelta;            // Cursor warp distance, included in the next mouse event delta
    NSCursor *hiddenCursor;             // Transparent cursor, used to hide the cursor over the view
    NSString *title;                    // Window title copy, restored when the style mask changes

    NSRect windowedContent;             // Window content rectangle before entering fullscreen/borderless mode

    char *clipboardText;                // Clipboard text, returned by GetClipboardText()

    mg_gamepads minigamepad;            // Gamepads state
} PlatformData;

//----------------------------------------------------------------------------------
// Global Variables Definition
//----------------------------------------------------------------------------------
extern CoreData CORE;                   // Global CORE state context

static PlatformData platform = { 0 };   // Platform specific data

// Map macOS virtual keycodes to raylib KeyboardKey values
static const int macosToRaylibMap[KEYMAP_SIZE] = {
    [0x00] = KEY_A, [0x01] = KEY_S, [0x02] = KEY_D, [0x03] = KEY_F, [0x04] = KEY_H, [0x05] = KEY_G,
    [0x06] = KEY_Z, [0x07] = KEY_X, [0x08] = KEY_C, [0x09] = KEY_V, [0x0B] = KEY_B, [0x0C] = KEY_Q,
    [0x0D] = KEY_W, [0x0E] = KEY_E, [0x0F] = KEY_R, [0x10] = KEY_Y, [0x11] = KEY_T, [0x12] = KEY_ONE,
    [0x13] = KEY_TWO, [0x14] = KEY_THREE, [0x15] = KEY_FOUR, [0x16] = KEY_SIX, [0x17] = KEY_FIVE,
    [0x18] = KEY_EQUAL, [0x19] = KEY_NINE, [0x1A] = KEY_SEVEN, [0x1B] = KEY_MINUS, [0x1C] = KEY_EIGHT,
    [0x1D] = KEY_ZERO, [0x1E] = KEY_RIGHT_BRACKET, [0x1F] = KEY_O, [0x20] = KEY_U, [0x21] = KEY_LEFT_BRACKET,
    [0x22] = KEY_I, [0x23] = KEY_P, [0x24] = KEY_ENTER, [0x25] = KEY_L, [0x26] = KEY_J, [0x27] = KEY_APOSTROPHE,
    [0x28] = KEY_K, [0x29] = KEY_SEMICOLON, [0x2A] = KEY_BACKSLASH, [0x2B] = KEY_COMMA, [0x2C] = KEY_SLASH,
    [0x2D] = KEY_N, [0x2E] = KEY_M, [0x2F] = KEY_PERIOD, [0x30] = KEY_TAB, [0x31] = KEY_SPACE,
    [0x32] = KEY_GRAVE, [0x33] = KEY_BACKSPACE, [0x35] = KEY_ESCAPE, [0x36] = KEY_RIGHT_SUPER,
    [0x37] = KEY_LEFT_SUPER, [0x38] = KEY_LEFT_SHIFT, [0x39] = KEY_CAPS_LOCK, [0x3A] = KEY_LEFT_ALT,
    [0x3B] = KEY_LEFT_CONTROL, [0x3C] = KEY_RIGHT_SHIFT, [0x3D] = KEY_RIGHT_ALT, [0x3E] = KEY_RIGHT_CONTROL,
    [0x41] = KEY_KP_DECIMAL, [0x43] = KEY_KP_MULTIPLY, [0x45] = KEY_KP_ADD, [0x47] = KEY_NUM_LOCK,
    [0x4B] = KEY_KP_DIVIDE, [0x4C] = KEY_KP_ENTER, [0x4E] = KEY_KP_SUBTRACT, [0x51] = KEY_KP_EQUAL,
    [0x52] = KEY_KP_0, [0x53] = KEY_KP_1, [0x54] = KEY_KP_2, [0x55] = KEY_KP_3, [0x56] = KEY_KP_4,
    [0x57] = KEY_KP_5, [0x58] = KEY_KP_6, [0x59] = KEY_KP_7, [0x5B] = KEY_KP_8, [0x5C] = KEY_KP_9,
    [0x60] = KEY_F5, [0x61] = KEY_F6, [0x62] = KEY_F7, [0x63] = KEY_F3, [0x64] = KEY_F8, [0x65] = KEY_F9,
    [0x67] = KEY_F11, [0x69] = KEY_PRINT_SCREEN, [0x6D] = KEY_F10, [0x6E] = KEY_KB_MENU, [0x6F] = KEY_F12,
    [0x72] = KEY_INSERT, [0x73] = KEY_HOME, [0x74] = KEY_PAGE_UP, [0x75] = KEY_DELETE, [0x76] = KEY_F4,
    [0x77] = KEY_END, [0x78] = KEY_F2, [0x79] = KEY_PAGE_DOWN, [0x7A] = KEY_F1, [0x7B] = KEY_LEFT,
    [0x7C] = KEY_RIGHT, [0x7D] = KEY_DOWN, [0x7E] = KEY_UP,
};

// Map minigamepad buttons to raylib GamepadButton values
static const int mgButtonConvertTable[] = {
    [MG_BUTTON_NORTH] = GAMEPAD_BUTTON_RIGHT_FACE_UP,
    [MG_BUTTON_EAST] = GAMEPAD_BUTTON_RIGHT_FACE_RIGHT,
    [MG_BUTTON_SOUTH] = GAMEPAD_BUTTON_RIGHT_FACE_DOWN,
    [MG_BUTTON_WEST] = GAMEPAD_BUTTON_RIGHT_FACE_LEFT,
    [MG_BUTTON_LEFT_SHOULDER] = GAMEPAD_BUTTON_LEFT_TRIGGER_1,
    [MG_BUTTON_RIGHT_SHOULDER] = GAMEPAD_BUTTON_RIGHT_TRIGGER_1,
    [MG_BUTTON_LEFT_TRIGGER] = GAMEPAD_BUTTON_LEFT_TRIGGER_2,
    [MG_BUTTON_RIGHT_TRIGGER] = GAMEPAD_BUTTON_RIGHT_TRIGGER_2,
    [MG_BUTTON_BACK] = GAMEPAD_BUTTON_MIDDLE_LEFT,
    [MG_BUTTON_GUIDE] = GAMEPAD_BUTTON_MIDDLE,
    [MG_BUTTON_START] = GAMEPAD_BUTTON_MIDDLE_RIGHT,
    [MG_BUTTON_DPAD_UP] = GAMEPAD_BUTTON_LEFT_FACE_UP,
    [MG_BUTTON_DPAD_RIGHT] = GAMEPAD_BUTTON_LEFT_FACE_RIGHT,
    [MG_BUTTON_DPAD_DOWN] = GAMEPAD_BUTTON_LEFT_FACE_DOWN,
    [MG_BUTTON_DPAD_LEFT] = GAMEPAD_BUTTON_LEFT_FACE_LEFT,
    [MG_BUTTON_LEFT_STICK] = GAMEPAD_BUTTON_LEFT_THUMB,
    [MG_BUTTON_RIGHT_STICK] = GAMEPAD_BUTTON_RIGHT_THUMB,
};

// Map minigamepad axes to raylib GamepadAxis values
static const int mgAxisConvertTable[] = {
    [MG_AXIS_LEFT_X] = GAMEPAD_AXIS_LEFT_X,
    [MG_AXIS_LEFT_Y] = GAMEPAD_AXIS_LEFT_Y,
    [MG_AXIS_RIGHT_X] = GAMEPAD_AXIS_RIGHT_X,
    [MG_AXIS_RIGHT_Y] = GAMEPAD_AXIS_RIGHT_Y,
    [MG_AXIS_LEFT_TRIGGER] = GAMEPAD_AXIS_LEFT_TRIGGER,
    [MG_AXIS_RIGHT_TRIGGER] = GAMEPAD_AXIS_RIGHT_TRIGGER,
    [MG_AXIS_HAT_DPAD_LEFT_RIGHT] = -1,
    [MG_AXIS_HAT_DPAD_UP_DOWN] = -1,
};

//----------------------------------------------------------------------------------
// Module Internal Functions Declaration
//----------------------------------------------------------------------------------
int InitPlatform(void);          // Initialize platform (graphics, inputs and more)
void ClosePlatform(void);        // Close platform

static void InitApplication(void);                              // Initialize shared application, delegate and menu bar
static void CreateMenuBar(void);                                // Create application menu bar (application and window menus)
static bool InitGraphicsContext(void);                          // Initialize OpenGL context or software framebuffer
#if !defined(GRAPHICS_API_OPENGL_SOFTWARE) && !defined(COCOA_EGL)
static NSOpenGLPixelFormat *CreatePixelFormat(NSOpenGLPixelFormatAttribute profile, bool msaa, NSOpenGLPixelFormatAttribute renderer); // Create OpenGL pixel format
#endif
static void UpdateFramebufferSize(void);                        // Update screen and render sizes from the content view
static void UpdateWindowPosition(void);                         // Update window position (content top-left, screen coordinates)
static NSWindowStyleMask GetWindowStyle(void);                  // Get window style mask from window flags
static bool IsWindowZoomed(void);                               // Check if window is zoomed (maximized)
static NSString *GetTitleString(const char *title);             // Get title string, empty for NULL or invalid UTF-8
static void UpdateWindowStyle(void);                            // Apply window style keeping the content size
static void SetWindowFullscreenFrame(bool enabled);             // Cover the current monitor with a borderless window, or restore it
static NSScreen *GetMonitorScreen(int monitor);                 // Get screen for a monitor index
static CGDirectDisplayID GetScreenDisplayID(NSScreen *screen);  // Get display id for a screen
static float GetPrimaryScreenHeight(void);                      // Get primary screen height, required to flip coordinates
static NSCursor *GetMouseCursor(int cursor);                    // Get system cursor for a raylib mouse cursor
static bool IsCursorInContentView(void);                        // Check if mouse cursor is inside the content view
static void UpdateCursor(void);                                 // Set current cursor, if inside the content view
static void WarpCursor(int x, int y);                           // Move the system cursor to a content view position
static void LockCursor(bool lock);                              // Lock/unlock cursor (disabled cursor mode)
static void HandleKey(int key, bool pressed, bool repeat);      // Register key state change
static void HandleModifiers(unsigned short keyCode, NSUInteger flags); // Register modifier key state change
static void HandleChars(NSString *text);                        // Register characters input
static void HandleMouseButton(int button, bool pressed);        // Register mouse button state change
static void HandleMouseMove(NSEvent *event);                    // Register mouse position change
static void HandleDroppedFiles(NSArray *urls);                  // Register dropped files

static void SetSwapInterval(int interval);                      // Set swap interval (V-Sync)
static void *GetProcAddressCocoa(const char *name);             // Get OpenGL procedure address

//----------------------------------------------------------------------------------
// Module Functions Declaration
//----------------------------------------------------------------------------------
// NOTE: Functions declaration is provided by raylib.h

//----------------------------------------------------------------------------------
// Module Functions Definition: Window and Graphics Device
//----------------------------------------------------------------------------------

// Check if application should close
bool WindowShouldClose(void)
{
    return CORE.Window.shouldClose;
}

// Toggle fullscreen mode
// NOTE: Borderless window covering the monitor above the menu bar, like GLFW, no video mode change
void ToggleFullscreen(void)
{
    @autoreleasepool
    {
        if (!FLAG_IS_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE))
        {
            if (FLAG_IS_SET(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE))
            {
                FLAG_CLEAR(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE);
                FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_UNDECORATED);
            }
            else SetWindowFullscreenFrame(true);

            FLAG_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE);
        }
        else
        {
            FLAG_CLEAR(CORE.Window.flags, FLAG_FULLSCREEN_MODE);
            SetWindowFullscreenFrame(false);
        }

        UpdateFramebufferSize();
    }
}

// Toggle borderless windowed mode
void ToggleBorderlessWindowed(void)
{
    @autoreleasepool
    {
        if (!FLAG_IS_SET(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE))
        {
            if (FLAG_IS_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE)) FLAG_CLEAR(CORE.Window.flags, FLAG_FULLSCREEN_MODE);
            else SetWindowFullscreenFrame(true);

            FLAG_SET(CORE.Window.flags, FLAG_WINDOW_UNDECORATED);
            FLAG_SET(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE);
        }
        else
        {
            FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_UNDECORATED);
            FLAG_CLEAR(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE);
            SetWindowFullscreenFrame(false);
        }

        // Refocus window, like GLFW
        SetWindowFocused();

        UpdateFramebufferSize();
    }
}

// Set window state: maximized, if resizable
void MaximizeWindow(void)
{
    @autoreleasepool
    {
        if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_RESIZABLE))
        {
            // NOTE: Fullscreen/borderless windows already cover the monitor, like GLFW only the flag is set
            bool covering = (FLAG_IS_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE) || FLAG_IS_SET(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE));
            if (!covering && !IsWindowZoomed()) [platform.window zoom:nil];
            FLAG_SET(CORE.Window.flags, FLAG_WINDOW_MAXIMIZED);
        }
    }
}

// Set window state: minimized
void MinimizeWindow(void)
{
    @autoreleasepool
    {
        // NOTE: Window delegate sets FLAG_WINDOW_MINIMIZED
        [platform.window miniaturize:nil];
    }
}

// Restore window from being minimized/maximized
void RestoreWindow(void)
{
    @autoreleasepool
    {
        if ([platform.window isMiniaturized]) [platform.window deminiaturize:nil];
        else if (IsWindowZoomed()) [platform.window zoom:nil];

        FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_MINIMIZED);
        FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_MAXIMIZED);
    }
}

// Set window configuration state using flags
void SetWindowState(unsigned int flags)
{
    @autoreleasepool
    {
        if (!CORE.Window.ready) TRACELOG(LOG_WARNING, "WINDOW: SetWindowState does nothing before window initialization, Use \"SetConfigFlags\" instead");

        // State change: FLAG_VSYNC_HINT
        if (!FLAG_IS_SET(CORE.Window.flags, FLAG_VSYNC_HINT) && FLAG_IS_SET(flags, FLAG_VSYNC_HINT))
        {
            SetSwapInterval(1);
            FLAG_SET(CORE.Window.flags, FLAG_VSYNC_HINT);
        }

        // State change: FLAG_BORDERLESS_WINDOWED_MODE
        if (!FLAG_IS_SET(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE) && FLAG_IS_SET(flags, FLAG_BORDERLESS_WINDOWED_MODE))
        {
            ToggleBorderlessWindowed(); // NOTE: Window state flag updated inside function
        }

        // State change: FLAG_FULLSCREEN_MODE
        if (!FLAG_IS_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE) && FLAG_IS_SET(flags, FLAG_FULLSCREEN_MODE))
        {
            ToggleFullscreen(); // NOTE: Window state flag updated inside function
        }

        // State change: FLAG_WINDOW_RESIZABLE
        if (!FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_RESIZABLE) && FLAG_IS_SET(flags, FLAG_WINDOW_RESIZABLE))
        {
            FLAG_SET(CORE.Window.flags, FLAG_WINDOW_RESIZABLE);
            UpdateWindowStyle();
        }

        // State change: FLAG_WINDOW_UNDECORATED
        if (!FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_UNDECORATED) && FLAG_IS_SET(flags, FLAG_WINDOW_UNDECORATED))
        {
            FLAG_SET(CORE.Window.flags, FLAG_WINDOW_UNDECORATED);
            UpdateWindowStyle();
        }

        // State change: FLAG_WINDOW_HIDDEN
        if (!FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_HIDDEN) && FLAG_IS_SET(flags, FLAG_WINDOW_HIDDEN))
        {
            [platform.window orderOut:nil];
            FLAG_SET(CORE.Window.flags, FLAG_WINDOW_HIDDEN);
        }

        // State change: FLAG_WINDOW_MINIMIZED
        if (!FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_MINIMIZED) && FLAG_IS_SET(flags, FLAG_WINDOW_MINIMIZED))
        {
            MinimizeWindow(); // NOTE: Window state flag updated by the window delegate
        }

        // State change: FLAG_WINDOW_MAXIMIZED
        if (!FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_MAXIMIZED) && FLAG_IS_SET(flags, FLAG_WINDOW_MAXIMIZED))
        {
            MaximizeWindow(); // NOTE: Window state flag updated inside function
        }

        // State change: FLAG_WINDOW_UNFOCUSED
        // NOTE: Only affects the window being shown later, like GLFW_FOCUS_ON_SHOW
        if (!FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_UNFOCUSED) && FLAG_IS_SET(flags, FLAG_WINDOW_UNFOCUSED))
        {
            platform.focusOnShow = false;
            FLAG_SET(CORE.Window.flags, FLAG_WINDOW_UNFOCUSED);
        }

        // State change: FLAG_WINDOW_TOPMOST
        if (!FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_TOPMOST) && FLAG_IS_SET(flags, FLAG_WINDOW_TOPMOST))
        {
            FLAG_SET(CORE.Window.flags, FLAG_WINDOW_TOPMOST);
            if (!FLAG_IS_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE) && !FLAG_IS_SET(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE)) [platform.window setLevel:NSFloatingWindowLevel];
        }

        // State change: FLAG_WINDOW_ALWAYS_RUN
        if (!FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_ALWAYS_RUN) && FLAG_IS_SET(flags, FLAG_WINDOW_ALWAYS_RUN))
        {
            FLAG_SET(CORE.Window.flags, FLAG_WINDOW_ALWAYS_RUN);
        }

        // The following states can not be changed after window creation

        // State change: FLAG_WINDOW_TRANSPARENT
        if (!FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_TRANSPARENT) && FLAG_IS_SET(flags, FLAG_WINDOW_TRANSPARENT))
        {
            TRACELOG(LOG_WARNING, "WINDOW: Framebuffer transparency can only be configured before window initialization");
        }

        // State change: FLAG_WINDOW_HIGHDPI
        if (!FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_HIGHDPI) && FLAG_IS_SET(flags, FLAG_WINDOW_HIGHDPI))
        {
            TRACELOG(LOG_WARNING, "WINDOW: High DPI can only be configured before window initialization");
        }

        // State change: FLAG_WINDOW_MOUSE_PASSTHROUGH
        if (!FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_MOUSE_PASSTHROUGH) && FLAG_IS_SET(flags, FLAG_WINDOW_MOUSE_PASSTHROUGH))
        {
            [platform.window setIgnoresMouseEvents:YES];
            FLAG_SET(CORE.Window.flags, FLAG_WINDOW_MOUSE_PASSTHROUGH);
        }

        // State change: FLAG_MSAA_4X_HINT
        if (!FLAG_IS_SET(CORE.Window.flags, FLAG_MSAA_4X_HINT) && FLAG_IS_SET(flags, FLAG_MSAA_4X_HINT))
        {
            TRACELOG(LOG_WARNING, "WINDOW: MSAA can only be configured before window initialization");
        }

        // State change: FLAG_INTERLACED_HINT
        if (!FLAG_IS_SET(CORE.Window.flags, FLAG_INTERLACED_HINT) && FLAG_IS_SET(flags, FLAG_INTERLACED_HINT))
        {
            TRACELOG(LOG_WARNING, "WINDOW: Interlaced mode can only be configured before window initialization");
        }
    }
}

// Clear window configuration state flags
void ClearWindowState(unsigned int flags)
{
    @autoreleasepool
    {
        // State change: FLAG_VSYNC_HINT
        if (FLAG_IS_SET(CORE.Window.flags, FLAG_VSYNC_HINT) && FLAG_IS_SET(flags, FLAG_VSYNC_HINT))
        {
            SetSwapInterval(0);
            FLAG_CLEAR(CORE.Window.flags, FLAG_VSYNC_HINT);
        }

        // State change: FLAG_BORDERLESS_WINDOWED_MODE
        if (FLAG_IS_SET(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE) && FLAG_IS_SET(flags, FLAG_BORDERLESS_WINDOWED_MODE))
        {
            ToggleBorderlessWindowed(); // NOTE: Window state flag updated inside function
        }

        // State change: FLAG_FULLSCREEN_MODE
        if (FLAG_IS_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE) && FLAG_IS_SET(flags, FLAG_FULLSCREEN_MODE))
        {
            ToggleFullscreen(); // NOTE: Window state flag updated inside function
        }

        // State change: FLAG_WINDOW_RESIZABLE
        if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_RESIZABLE) && FLAG_IS_SET(flags, FLAG_WINDOW_RESIZABLE))
        {
            FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_RESIZABLE);
            UpdateWindowStyle();
        }

        // State change: FLAG_WINDOW_HIDDEN
        if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_HIDDEN) && FLAG_IS_SET(flags, FLAG_WINDOW_HIDDEN))
        {
            FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_HIDDEN);
            if (platform.focusOnShow) [platform.window makeKeyAndOrderFront:nil];
            else [platform.window orderFront:nil];
        }

        // State change: FLAG_WINDOW_MINIMIZED
        if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_MINIMIZED) && FLAG_IS_SET(flags, FLAG_WINDOW_MINIMIZED))
        {
            RestoreWindow(); // NOTE: Window state flag updated inside function
        }

        // State change: FLAG_WINDOW_MAXIMIZED
        if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_MAXIMIZED) && FLAG_IS_SET(flags, FLAG_WINDOW_MAXIMIZED))
        {
            RestoreWindow(); // NOTE: Window state flag updated inside function
        }

        // State change: FLAG_WINDOW_UNDECORATED
        if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_UNDECORATED) && FLAG_IS_SET(flags, FLAG_WINDOW_UNDECORATED))
        {
            FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_UNDECORATED);
            UpdateWindowStyle();
        }

        // State change: FLAG_WINDOW_UNFOCUSED
        if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_UNFOCUSED) && FLAG_IS_SET(flags, FLAG_WINDOW_UNFOCUSED))
        {
            platform.focusOnShow = true;
            FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_UNFOCUSED);
        }

        // State change: FLAG_WINDOW_TOPMOST
        if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_TOPMOST) && FLAG_IS_SET(flags, FLAG_WINDOW_TOPMOST))
        {
            FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_TOPMOST);
            if (!FLAG_IS_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE) && !FLAG_IS_SET(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE)) [platform.window setLevel:NSNormalWindowLevel];
        }

        // State change: FLAG_WINDOW_ALWAYS_RUN
        if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_ALWAYS_RUN) && FLAG_IS_SET(flags, FLAG_WINDOW_ALWAYS_RUN))
        {
            FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_ALWAYS_RUN);
        }

        // The following states can not be changed after window creation

        // State change: FLAG_WINDOW_TRANSPARENT
        if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_TRANSPARENT) && FLAG_IS_SET(flags, FLAG_WINDOW_TRANSPARENT))
        {
            TRACELOG(LOG_WARNING, "WINDOW: Framebuffer transparency can only be configured before window initialization");
        }

        // State change: FLAG_WINDOW_HIGHDPI
        if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_HIGHDPI) && FLAG_IS_SET(flags, FLAG_WINDOW_HIGHDPI))
        {
            TRACELOG(LOG_WARNING, "WINDOW: High DPI can only be configured before window initialization");
        }

        // State change: FLAG_WINDOW_MOUSE_PASSTHROUGH
        if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_MOUSE_PASSTHROUGH) && FLAG_IS_SET(flags, FLAG_WINDOW_MOUSE_PASSTHROUGH))
        {
            [platform.window setIgnoresMouseEvents:NO];
            FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_MOUSE_PASSTHROUGH);
        }

        // State change: FLAG_MSAA_4X_HINT
        if (FLAG_IS_SET(CORE.Window.flags, FLAG_MSAA_4X_HINT) && FLAG_IS_SET(flags, FLAG_MSAA_4X_HINT))
        {
            TRACELOG(LOG_WARNING, "WINDOW: MSAA can only be configured before window initialization");
        }

        // State change: FLAG_INTERLACED_HINT
        if (FLAG_IS_SET(CORE.Window.flags, FLAG_INTERLACED_HINT) && FLAG_IS_SET(flags, FLAG_INTERLACED_HINT))
        {
            TRACELOG(LOG_WARNING, "WINDOW: Interlaced mode can only be configured before window initialization");
        }
    }
}

// Set icon for window
// NOTE: macOS has no window icons, the image is used as application icon (Dock)
void SetWindowIcon(Image image)
{
    @autoreleasepool
    {
        if (image.data == NULL)
        {
            [NSApp setApplicationIconImage:nil];
            return;
        }

        if (image.format != PIXELFORMAT_UNCOMPRESSED_R8G8B8A8)
        {
            TRACELOG(LOG_WARNING, "WINDOW: Window icon image must be in R8G8B8A8 pixel format");
            return;
        }

        NSBitmapImageRep *rep = [[NSBitmapImageRep alloc] initWithBitmapDataPlanes:NULL
            pixelsWide:image.width pixelsHigh:image.height bitsPerSample:8 samplesPerPixel:4 hasAlpha:YES isPlanar:NO
            colorSpaceName:NSDeviceRGBColorSpace bitmapFormat:NSBitmapFormatAlphaNonpremultiplied
            bytesPerRow:image.width*4 bitsPerPixel:32];

        if (rep == nil) return;

        memcpy([rep bitmapData], image.data, image.width*image.height*4);

        NSImage *icon = [[NSImage alloc] initWithSize:NSMakeSize(image.width, image.height)];
        [icon addRepresentation:rep];
        [NSApp setApplicationIconImage:icon];

        [icon release];
        [rep release];
    }
}

// Set icon for window, multiple images
// NOTE: The largest image is used, the system scales it for every Dock size
void SetWindowIcons(Image *images, int count)
{
    if ((images == NULL) || (count <= 0))
    {
        SetWindowIcon((Image){ 0 });
        return;
    }

    int largest = 0;
    for (int i = 1; i < count; i++)
    {
        if ((images[i].width*images[i].height) > (images[largest].width*images[largest].height)) largest = i;
    }

    SetWindowIcon(images[largest]);
}

// Set title for window
void SetWindowTitle(const char *title)
{
    CORE.Window.title = title;

    @autoreleasepool
    {
        [platform.title release];
        platform.title = [GetTitleString(title) retain];
        [platform.window setTitle:platform.title];
    }
}

// Set window position on screen (windowed mode)
void SetWindowPosition(int x, int y)
{
    @autoreleasepool
    {
        NSRect contentRect = [platform.window contentRectForFrameRect:[platform.window frame]];
        contentRect.origin.x = x;
        contentRect.origin.y = GetPrimaryScreenHeight() - y - contentRect.size.height;

        [platform.window setFrameOrigin:[platform.window frameRectForContentRect:contentRect].origin];

        UpdateWindowPosition();
    }
}

// Set monitor for the current window
void SetWindowMonitor(int monitor)
{
    @autoreleasepool
    {
        NSScreen *screen = GetMonitorScreen(monitor);

        if (screen != nil)
        {
            TRACELOG(LOG_INFO, "COCOA: Selected monitor: [%i] %s", monitor, GetMonitorName(monitor));

            if (FLAG_IS_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE) || FLAG_IS_SET(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE))
            {
                [platform.window setFrame:[screen frame] display:YES];
            }
            else
            {
                NSRect screenFrame = [screen visibleFrame];
                NSRect contentRect = [platform.window contentRectForFrameRect:[platform.window frame]];
                contentRect.origin.x = screenFrame.origin.x + (screenFrame.size.width - contentRect.size.width)/2;
                contentRect.origin.y = screenFrame.origin.y + (screenFrame.size.height - contentRect.size.height)/2;

                [platform.window setFrameOrigin:[platform.window frameRectForContentRect:contentRect].origin];
            }

            UpdateWindowPosition();
            UpdateFramebufferSize();
        }
        else TRACELOG(LOG_WARNING, "COCOA: Failed to find selected monitor");
    }
}

// Set window minimum dimensions (for FLAG_WINDOW_RESIZABLE)
void SetWindowMinSize(int width, int height)
{
    @autoreleasepool
    {
        CORE.Window.screenMin.width = width;
        CORE.Window.screenMin.height = height;

        [platform.window setContentMinSize:NSMakeSize((width > 0)? width : 1, (height > 0)? height : 1)];
    }
}

// Set window maximum dimensions (for FLAG_WINDOW_RESIZABLE)
void SetWindowMaxSize(int width, int height)
{
    @autoreleasepool
    {
        CORE.Window.screenMax.width = width;
        CORE.Window.screenMax.height = height;

        [platform.window setContentMaxSize:NSMakeSize((width > 0)? width : FLT_MAX, (height > 0)? height : FLT_MAX)];
    }
}

// Set window dimensions
void SetWindowSize(int width, int height)
{
    @autoreleasepool
    {
        if ((CORE.Window.screenMin.width > 0) && (width < (int)CORE.Window.screenMin.width)) width = CORE.Window.screenMin.width;
        if ((CORE.Window.screenMin.height > 0) && (height < (int)CORE.Window.screenMin.height)) height = CORE.Window.screenMin.height;
        if ((CORE.Window.screenMax.width > 0) && (width > (int)CORE.Window.screenMax.width)) width = CORE.Window.screenMax.width;
        if ((CORE.Window.screenMax.height > 0) && (height > (int)CORE.Window.screenMax.height)) height = CORE.Window.screenMax.height;

        // NOTE: Keep the content top-left corner in place, Cocoa origin is the bottom-left corner
        NSRect contentRect = [platform.window contentRectForFrameRect:[platform.window frame]];
        contentRect.origin.y += contentRect.size.height - height;
        contentRect.size = NSMakeSize(width, height);

        [platform.window setFrame:[platform.window frameRectForContentRect:contentRect] display:YES];

        UpdateFramebufferSize();
    }
}

// Set window opacity, value opacity is between 0.0 and 1.0
void SetWindowOpacity(float opacity)
{
    @autoreleasepool
    {
        if (opacity >= 1.0f) opacity = 1.0f;
        else if (opacity <= 0.0f) opacity = 0.0f;

        [platform.window setAlphaValue:opacity];
    }
}

// Set window focused
void SetWindowFocused(void)
{
    @autoreleasepool
    {
        [NSApp activateIgnoringOtherApps:YES];
        [platform.window makeKeyAndOrderFront:nil];
    }
}

// Get native window handle
void *GetWindowHandle(void)
{
    return (void *)platform.window;
}

// Get number of monitors
int GetMonitorCount(void)
{
    @autoreleasepool
    {
        return (int)[[NSScreen screens] count];
    }
}

// Get current monitor where window is placed
int GetCurrentMonitor(void)
{
    @autoreleasepool
    {
        NSScreen *screen = [platform.window screen];
        NSUInteger index = [[NSScreen screens] indexOfObject:screen];

        return (index == NSNotFound)? 0 : (int)index;
    }
}

// Get selected monitor position
Vector2 GetMonitorPosition(int monitor)
{
    @autoreleasepool
    {
        NSScreen *screen = GetMonitorScreen(monitor);

        if (screen != nil)
        {
            NSRect frame = [screen frame];
            return (Vector2){ (float)frame.origin.x, GetPrimaryScreenHeight() - (float)(frame.origin.y + frame.size.height) };
        }
        else TRACELOG(LOG_WARNING, "COCOA: Failed to find selected monitor");

        return (Vector2){ 0, 0 };
    }
}

// Get selected monitor width (currently used by monitor)
int GetMonitorWidth(int monitor)
{
    @autoreleasepool
    {
        int width = 0;
        NSScreen *screen = GetMonitorScreen(monitor);

        if (screen != nil) width = (int)[screen frame].size.width;
        else TRACELOG(LOG_WARNING, "COCOA: Failed to find selected monitor");

        return width;
    }
}

// Get selected monitor height (currently used by monitor)
int GetMonitorHeight(int monitor)
{
    @autoreleasepool
    {
        int height = 0;
        NSScreen *screen = GetMonitorScreen(monitor);

        if (screen != nil) height = (int)[screen frame].size.height;
        else TRACELOG(LOG_WARNING, "COCOA: Failed to find selected monitor");

        return height;
    }
}

// Get selected monitor physical width in millimetres
int GetMonitorPhysicalWidth(int monitor)
{
    @autoreleasepool
    {
        int width = 0;
        NSScreen *screen = GetMonitorScreen(monitor);

        if (screen != nil) width = (int)CGDisplayScreenSize(GetScreenDisplayID(screen)).width;
        else TRACELOG(LOG_WARNING, "COCOA: Failed to find selected monitor");

        return width;
    }
}

// Get selected monitor physical height in millimetres
int GetMonitorPhysicalHeight(int monitor)
{
    @autoreleasepool
    {
        int height = 0;
        NSScreen *screen = GetMonitorScreen(monitor);

        if (screen != nil) height = (int)CGDisplayScreenSize(GetScreenDisplayID(screen)).height;
        else TRACELOG(LOG_WARNING, "COCOA: Failed to find selected monitor");

        return height;
    }
}

// Get selected monitor refresh rate
int GetMonitorRefreshRate(int monitor)
{
    @autoreleasepool
    {
        int refresh = 0;
        NSScreen *screen = GetMonitorScreen(monitor);

        if (screen != nil)
        {
            CGDisplayModeRef mode = CGDisplayCopyDisplayMode(GetScreenDisplayID(screen));
            if (mode != NULL)
            {
                refresh = (int)round(CGDisplayModeGetRefreshRate(mode));
                CGDisplayModeRelease(mode);
            }

            // NOTE: Built-in displays report a refresh rate of 0, maximumFramesPerSecond requires macOS 12
            if ((refresh == 0) && [screen respondsToSelector:@selector(maximumFramesPerSecond)]) refresh = [[screen valueForKey:@"maximumFramesPerSecond"] intValue];
            if (refresh == 0) refresh = 60;
        }
        else TRACELOG(LOG_WARNING, "COCOA: Failed to find selected monitor");

        return refresh;
    }
}

// Get the human-readable, UTF-8 encoded name of the selected monitor
const char *GetMonitorName(int monitor)
{
    @autoreleasepool
    {
        static char name[128] = { 0 };
        NSScreen *screen = GetMonitorScreen(monitor);

        memset(name, 0, sizeof(name));

        if (screen != nil) snprintf(name, sizeof(name), "%s", [[screen localizedName] UTF8String]);
        else TRACELOG(LOG_WARNING, "COCOA: Failed to find selected monitor");

        return name;
    }
}

// Get window position XY on monitor
Vector2 GetWindowPosition(void)
{
    @autoreleasepool
    {
        UpdateWindowPosition();

        return (Vector2){ (float)CORE.Window.position.x, (float)CORE.Window.position.y };
    }
}

// Get window scale DPI factor for current monitor
Vector2 GetWindowScaleDPI(void)
{
    @autoreleasepool
    {
        Vector2 scale = { 1.0f, 1.0f };

        if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_HIGHDPI))
        {
            scale.x = (float)[platform.window backingScaleFactor];
            scale.y = scale.x;
        }

        return scale;
    }
}

// Set clipboard text content
void SetClipboardText(const char *text)
{
    @autoreleasepool
    {
        NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
        [pasteboard clearContents];
        if (text != NULL) [pasteboard setString:[NSString stringWithUTF8String:text] forType:NSPasteboardTypeString];
    }
}

// Get clipboard text content
// NOTE: returned string is allocated and freed by the platform
const char *GetClipboardText(void)
{
    RL_FREE(platform.clipboardText);
    platform.clipboardText = NULL;

    @autoreleasepool
    {
        NSString *text = [[NSPasteboard generalPasteboard] stringForType:NSPasteboardTypeString];

        if (text != nil)
        {
            const char *utf8 = [text UTF8String];
            size_t length = strlen(utf8);
            platform.clipboardText = (char *)RL_MALLOC(length + 1);
            memcpy(platform.clipboardText, utf8, length + 1);
        }
    }

    return platform.clipboardText;
}

// Get clipboard image
Image GetClipboardImage(void)
{
    Image image = { 0 };

#if SUPPORT_CLIPBOARD_IMAGE && SUPPORT_MODULE_RTEXTURES
    @autoreleasepool
    {
        NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
        NSData *data = [pasteboard dataForType:NSPasteboardTypePNG];

        // Convert any other image type to PNG (screenshots are copied as TIFF)
        if (data == nil)
        {
            NSData *tiff = [pasteboard dataForType:NSPasteboardTypeTIFF];
            if (tiff != nil)
            {
                NSBitmapImageRep *rep = [NSBitmapImageRep imageRepWithData:tiff];
                if (rep != nil) data = [rep representationUsingType:NSBitmapImageFileTypePNG properties:@{}];
            }
        }

        if (data != nil) image = LoadImageFromMemory(".png", (const unsigned char *)[data bytes], (int)[data length]);
        else TRACELOG(LOG_WARNING, "Clipboard image: Couldn't get clipboard data.");
    }
#else
    TRACELOG(LOG_WARNING, "Clipboard image: SUPPORT_CLIPBOARD_IMAGE requires SUPPORT_MODULE_RTEXTURES to work properly");
#endif

    return image;
}

// Show mouse cursor
void ShowCursor(void)
{
    @autoreleasepool
    {
        CORE.Input.Mouse.cursorHidden = false;
        UpdateCursor();
    }
}

// Hides mouse cursor
void HideCursor(void)
{
    @autoreleasepool
    {
        CORE.Input.Mouse.cursorHidden = true;
        UpdateCursor();
    }
}

// Enables cursor (unlock cursor)
void EnableCursor(void)
{
    @autoreleasepool
    {
        LockCursor(false);

        CORE.Input.Mouse.cursorHidden = false;
        CORE.Input.Mouse.cursorLocked = false;

        // Set cursor position in the middle
        SetMousePosition(CORE.Window.screen.width/2, CORE.Window.screen.height/2);

        UpdateCursor();
    }
}

// Disables cursor (lock cursor)
void DisableCursor(void)
{
    @autoreleasepool
    {
        // Reset mouse position within the window area before disabling cursor
        SetMousePosition(CORE.Window.screen.width/2, CORE.Window.screen.height/2);

        CORE.Input.Mouse.cursorHidden = true;
        CORE.Input.Mouse.cursorLocked = true;
        LockCursor(true);
    }
}

// Swap back buffer with front buffer (screen drawing)
void SwapScreenBuffer(void)
{
    @autoreleasepool
    {
#if defined(GRAPHICS_API_OPENGL_SOFTWARE)
        if ((platform.softwarePixels == NULL) || (platform.view == nil)) return;

        int width = platform.softwareWidth;
        int height = platform.softwareHeight;
        bool transparent = FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_TRANSPARENT);

        rlCopyFramebuffer(0, 0, width, height, PIXELFORMAT_UNCOMPRESSED_R8G8B8A8, platform.softwarePixels);

        // NOTE: Software renderer framebuffer copy is in BGRA order, premultiply alpha for transparent windows
        if (transparent)
        {
            for (int i = 0; i < width*height; i++)
            {
                unsigned char *p = platform.softwarePixels + i*4;
                p[0] = (unsigned char)(p[0]*p[3]/255);
                p[1] = (unsigned char)(p[1]*p[3]/255);
                p[2] = (unsigned char)(p[2]*p[3]/255);
            }
        }

        CFDataRef data = CFDataCreate(NULL, platform.softwarePixels, width*height*4);
        CGDataProviderRef provider = CGDataProviderCreateWithCFData(data);
        CGBitmapInfo info = kCGBitmapByteOrder32Little | (transparent? kCGImageAlphaPremultipliedFirst : kCGImageAlphaNoneSkipFirst);
        CGImageRef image = CGImageCreate(width, height, 8, 32, width*4, platform.colorSpace, info, provider, NULL, false, kCGRenderingIntentDefault);

        [CATransaction begin];
        [CATransaction setDisableActions:YES];
        [[platform.view layer] setContents:(id)image];
        [CATransaction commit];

        CGImageRelease(image);
        CGDataProviderRelease(provider);
        CFRelease(data);
#else
        // NOTE: Swap interval is not applied to occluded windows, frames are paced to avoid a busy loop
        if (platform.occluded && FLAG_IS_SET(CORE.Window.flags, FLAG_VSYNC_HINT))
        {
            double period = 1.0/OCCLUDED_FRAME_RATE;
            double delay = period - fmod(GetTime(), period);
            usleep((useconds_t)(delay*1e6));
        }

#if defined(COCOA_EGL)
        platform.egl.SwapBuffers(platform.eglDisplay, platform.eglSurface);
#else
        [platform.glContext flushBuffer];
#endif
#endif
    }
}

//----------------------------------------------------------------------------------
// Module Functions Definition: Misc
//----------------------------------------------------------------------------------

// Get elapsed time measure in seconds since InitTimer()
double GetTime(void)
{
    double time = 0.0;
    struct timespec ts = { 0 };
    clock_gettime(CLOCK_MONOTONIC, &ts);
    unsigned long long nanoSeconds = (unsigned long long)ts.tv_sec*1000000000LLU + (unsigned long long)ts.tv_nsec;

    time = (double)(nanoSeconds - CORE.Time.base)*1e-9; // Elapsed time since InitTimer()

    return time;
}

// Open URL with default system browser (if available)
// WARNING: This function is only safe to use if you control the URL given,
// a user could craft a malicious string to perform and undesired action
// NOTE: Some safety checks have been added to mitigate security issues
void OpenURL(const char *url)
{
    // Security check to (partially) avoid malicious code
    if ((strchr(url, '\'') != NULL) || (strchr(url, '\"') != NULL))
    {
        // Filter characters: ' and "
        TRACELOG(LOG_WARNING, "SYSTEM: Provided URL could be potentially malicious, avoid [\'\"] characters");
    }
    else if ((strncmp(url, "http://", 7) != 0) && (strncmp(url, "https://", 8) != 0))
    {
        // Only allow URL starting with "http://" or "https://" protocols
        TRACELOG(LOG_WARNING, "SYSTEM: Provided URL must start with 'http://' or 'https://' protocols");
    }
    else
    {
        @autoreleasepool
        {
            NSURL *nsurl = [NSURL URLWithString:[NSString stringWithUTF8String:url]];
            if ((nsurl == nil) || ![[NSWorkspace sharedWorkspace] openURL:nsurl]) TRACELOG(LOG_WARNING, "OpenURL() failed to open the URL");
        }
    }
}

//----------------------------------------------------------------------------------
// Module Functions Definition: Inputs
//----------------------------------------------------------------------------------

// Set internal gamepad mappings
int SetGamepadMappings(const char *mappings)
{
    return mg_update_gamepad_mappings(&platform.minigamepad, mappings);
}

// Set gamepad vibration
void SetGamepadVibration(int gamepad, float leftMotor, float rightMotor, float duration)
{
    TRACELOG(LOG_WARNING, "SetGamepadVibration() not available on target platform");
}

// Set mouse position XY
void SetMousePosition(int x, int y)
{
    @autoreleasepool
    {
        CORE.Input.Mouse.currentPosition = (Vector2){ (float)x, (float)y };
        CORE.Input.Mouse.previousPosition = CORE.Input.Mouse.currentPosition;

        if (!CORE.Input.Mouse.cursorLocked) WarpCursor(x, y);
    }
}

// Set mouse cursor
void SetMouseCursor(int cursor)
{
    @autoreleasepool
    {
        CORE.Input.Mouse.cursor = cursor;
        UpdateCursor();
    }
}

// Get physical key name
// NOTE: Returns the character the key produces on the current keyboard layout, NULL for non-printable keys (like GLFW)
const char *GetKeyName(int key)
{
    static char keyName[16] = { 0 };

    if ((key != KEY_KP_EQUAL) && ((key < KEY_KP_0) || (key > KEY_KP_ADD)) && ((key < KEY_APOSTROPHE) || (key > KEY_GRAVE))) return NULL;

    int keyCode = -1;
    for (int i = 0; i < KEYMAP_SIZE; i++)
    {
        if (macosToRaylibMap[i] == key)
        {
            keyCode = i;
            break;
        }
    }

    if (keyCode < 0) return NULL;

    // Text Input Sources functions, from HIToolbox (loaded by AppKit)
    CFBundleRef bundle = CFBundleGetBundleWithIdentifier(CFSTR("com.apple.HIToolbox"));
    if (bundle == NULL) return NULL;

    CFStringRef *unicodeKeyLayoutData = (CFStringRef *)CFBundleGetDataPointerForName(bundle, CFSTR("kTISPropertyUnicodeKeyLayoutData"));
    void *(*CopyCurrentKeyboardLayoutInputSource)(void) = (void *(*)(void))CFBundleGetFunctionPointerForName(bundle, CFSTR("TISCopyCurrentKeyboardLayoutInputSource"));
    void *(*GetInputSourceProperty)(void *, CFStringRef) = (void *(*)(void *, CFStringRef))CFBundleGetFunctionPointerForName(bundle, CFSTR("TISGetInputSourceProperty"));
    UInt8 (*GetKbdType)(void) = (UInt8 (*)(void))CFBundleGetFunctionPointerForName(bundle, CFSTR("LMGetKbdType"));

    if ((unicodeKeyLayoutData == NULL) || (CopyCurrentKeyboardLayoutInputSource == NULL) || (GetInputSourceProperty == NULL) || (GetKbdType == NULL)) return NULL;

    void *inputSource = CopyCurrentKeyboardLayoutInputSource();
    if (inputSource == NULL) return NULL;

    const char *result = NULL;
    CFDataRef layoutData = (CFDataRef)GetInputSourceProperty(inputSource, *unicodeKeyLayoutData);

    if (layoutData != NULL)
    {
        UInt32 deadKeyState = 0;
        UniChar characters[4] = { 0 };
        UniCharCount characterCount = 0;

        if ((UCKeyTranslate((const UCKeyboardLayout *)CFDataGetBytePtr(layoutData), (UInt16)keyCode, kUCKeyActionDisplay, 0, GetKbdType(),
            kUCKeyTranslateNoDeadKeysBit, &deadKeyState, 4, &characterCount, characters) == noErr) && (characterCount > 0))
        {
            CFStringRef string = CFStringCreateWithCharacters(kCFAllocatorDefault, characters, characterCount);
            if (CFStringGetCString(string, keyName, sizeof(keyName), kCFStringEncodingUTF8)) result = keyName;
            CFRelease(string);
        }
    }

    CFRelease((CFTypeRef)inputSource);

    return result;
}

// Register all input events
void PollInputEvents(void)
{
#if SUPPORT_GESTURES_SYSTEM
    // NOTE: Gestures update must be called every frame to reset gestures correctly
    // because ProcessGestureEvent() is called on an event, not every frame
    UpdateGestures();
#endif

    // Reset keys/chars pressed registered
    CORE.Input.Keyboard.keyPressedQueueCount = 0;
    CORE.Input.Keyboard.charPressedQueueCount = 0;

    // Reset last gamepad button/axis registered state
    CORE.Input.Gamepad.lastButtonPressed = GAMEPAD_BUTTON_UNKNOWN;

    // Register previous keys states
    for (int i = 0; i < MAX_KEYBOARD_KEYS; i++)
    {
        CORE.Input.Keyboard.previousKeyState[i] = CORE.Input.Keyboard.currentKeyState[i];
        CORE.Input.Keyboard.keyRepeatInFrame[i] = 0;
    }

    // Register previous mouse states
    for (int i = 0; i < MAX_MOUSE_BUTTONS; i++) CORE.Input.Mouse.previousButtonState[i] = CORE.Input.Mouse.currentButtonState[i];

    // Register previous mouse wheel state
    CORE.Input.Mouse.previousWheelMove = CORE.Input.Mouse.currentWheelMove;
    CORE.Input.Mouse.currentWheelMove = (Vector2){ 0.0f, 0.0f };

    // Register previous mouse position
    CORE.Input.Mouse.previousPosition = CORE.Input.Mouse.currentPosition;

    // Register previous touch states
    for (int i = 0; i < MAX_TOUCH_POINTS; i++) CORE.Input.Touch.previousTouchState[i] = CORE.Input.Touch.currentTouchState[i];

    // Map touch position to mouse position for convenience
    CORE.Input.Touch.position[0] = CORE.Input.Mouse.currentPosition;

    // Register previous gamepad states
    for (int i = 0; i < MAX_GAMEPADS; i++)
    {
        if (CORE.Input.Gamepad.ready[i])
        {
            for (int k = 0; k < MAX_GAMEPAD_BUTTONS; k++) CORE.Input.Gamepad.previousButtonState[i][k] = CORE.Input.Gamepad.currentButtonState[i][k];
        }
    }

    CORE.Window.resizedLastFrame = false;

    // Process window and input events
    @autoreleasepool
    {
        bool waitEvents = (CORE.Window.eventWaiting ||
            (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_MINIMIZED) && !FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_ALWAYS_RUN)));

        // NOTE: When waiting, the first event blocks (drawing is paused), the rest are processed without waiting
        NSDate *until = waitEvents? [NSDate distantFuture] : [NSDate distantPast];

        while (true)
        {
            NSEvent *event = [NSApp nextEventMatchingMask:NSEventMaskAny untilDate:until inMode:NSDefaultRunLoopMode dequeue:YES];
            if (event == nil) break;

            [NSApp sendEvent:event];
            until = [NSDate distantPast];
        }

        if (waitEvents) CORE.Time.previous = GetTime();
    }

    // Process gamepad events
    const int buttonTableSize = (int)(sizeof(mgButtonConvertTable)/sizeof(mgButtonConvertTable[0]));
    const int axisTableSize = (int)(sizeof(mgAxisConvertTable)/sizeof(mgAxisConvertTable[0]));
    mg_event gamepadEvent = { 0 };
    while (mg_gamepads_check_event(&platform.minigamepad, &gamepadEvent))
    {
        int gamepadIndex = gamepadEvent.gamepad->index;

        switch (gamepadEvent.type)
        {
            case MG_EVENT_BUTTON_PRESS:
            {
                int button = ((gamepadEvent.button >= 0) && (gamepadEvent.button < buttonTableSize))? mgButtonConvertTable[gamepadEvent.button] : -1;
                if (button >= 0)
                {
                    CORE.Input.Gamepad.currentButtonState[gamepadIndex][button] = 1;
                    CORE.Input.Gamepad.lastButtonPressed = button;
                }
            } break;
            case MG_EVENT_BUTTON_RELEASE:
            {
                int button = ((gamepadEvent.button >= 0) && (gamepadEvent.button < buttonTableSize))? mgButtonConvertTable[gamepadEvent.button] : -1;
                if (button >= 0)
                {
                    CORE.Input.Gamepad.currentButtonState[gamepadIndex][button] = 0;
                    if (CORE.Input.Gamepad.lastButtonPressed == button) CORE.Input.Gamepad.lastButtonPressed = GAMEPAD_BUTTON_UNKNOWN;
                }
            } break;
            case MG_EVENT_AXIS_MOVE:
            {
                int axis = ((gamepadEvent.axis >= 0) && (gamepadEvent.axis < axisTableSize))? mgAxisConvertTable[gamepadEvent.axis] : -1;
                float value = platform.minigamepad.gamepads[gamepadIndex].axes[gamepadEvent.axis].value;

                switch (axis)
                {
                    case GAMEPAD_AXIS_LEFT_X:
                    case GAMEPAD_AXIS_LEFT_Y:
                    case GAMEPAD_AXIS_RIGHT_X:
                    case GAMEPAD_AXIS_RIGHT_Y:
                    {
                        CORE.Input.Gamepad.axisState[gamepadIndex][axis] = value;
                    } break;
                    case GAMEPAD_AXIS_LEFT_TRIGGER:
                    case GAMEPAD_AXIS_RIGHT_TRIGGER:
                    {
                        CORE.Input.Gamepad.axisState[gamepadIndex][axis] = value;

                        // Trigger button press when axis is all the way
                        int button = (axis == GAMEPAD_AXIS_LEFT_TRIGGER)? GAMEPAD_BUTTON_LEFT_TRIGGER_2 : GAMEPAD_BUTTON_RIGHT_TRIGGER_2;
                        int pressed = (value >= 1.0f);

                        CORE.Input.Gamepad.currentButtonState[gamepadIndex][button] = pressed;
                        if (pressed) CORE.Input.Gamepad.lastButtonPressed = button;
                        else if (CORE.Input.Gamepad.lastButtonPressed == button) CORE.Input.Gamepad.lastButtonPressed = GAMEPAD_BUTTON_UNKNOWN;
                    } break;
                    default: break;
                }
            } break;
            case MG_EVENT_GAMEPAD_CONNECT:
            {
                CORE.Input.Gamepad.ready[gamepadIndex] = true;
                CORE.Input.Gamepad.axisState[gamepadIndex][GAMEPAD_AXIS_LEFT_TRIGGER] = -1.0f;
                CORE.Input.Gamepad.axisState[gamepadIndex][GAMEPAD_AXIS_RIGHT_TRIGGER] = -1.0f;

                int axisCount = 0;
                for (int i = 0; i < MG_AXIS_COUNT; i++)
                {
                    if (platform.minigamepad.gamepads[gamepadIndex].axes[i].supported) axisCount++;
                    else break;
                }

                CORE.Input.Gamepad.axisCount[gamepadIndex] = axisCount;
                snprintf(CORE.Input.Gamepad.name[gamepadIndex], MAX_GAMEPAD_NAME_LENGTH, "%s", platform.minigamepad.gamepads[gamepadIndex].name);
            } break;
            case MG_EVENT_GAMEPAD_DISCONNECT:
            {
                CORE.Input.Gamepad.ready[gamepadIndex] = false;
                memset(CORE.Input.Gamepad.name[gamepadIndex], 0, MAX_GAMEPAD_NAME_LENGTH);
            } break;
            default: break;
        }
    }

#if defined(GRAPHICS_API_OPENGL_SOFTWARE)
    // Resize software renderer framebuffer to the window framebuffer size
    if ((platform.softwareWidth != (int)CORE.Window.render.width) || (platform.softwareHeight != (int)CORE.Window.render.height))
    {
        platform.softwareWidth = CORE.Window.render.width;
        platform.softwareHeight = CORE.Window.render.height;
        RL_FREE(platform.softwarePixels);
        platform.softwarePixels = (unsigned char *)RL_MALLOC(platform.softwareWidth*platform.softwareHeight*4);
        rlResizeFramebuffer(platform.softwareWidth, platform.softwareHeight);
    }
#endif

    CORE.Window.shouldClose = platform.closeRequested;

    // Reset close status for next frame
    platform.closeRequested = false;
}

//----------------------------------------------------------------------------------
// Module Internal Functions Definition
//----------------------------------------------------------------------------------
// Initialize platform: graphics, inputs and more
int InitPlatform(void)
{
    @autoreleasepool
    {
        InitApplication();

        // Initialize window size
        NSScreen *mainScreen = [NSScreen mainScreen];
        if (mainScreen == nil)
        {
            TRACELOG(LOG_WARNING, "COCOA: Failed to get main screen");
            return -1;
        }

        CORE.Window.display.width = (int)[mainScreen frame].size.width;
        CORE.Window.display.height = (int)[mainScreen frame].size.height;

        bool fullscreen = FLAG_IS_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE);
        bool borderless = FLAG_IS_SET(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE);
        if ((CORE.Window.screen.width == 0) || (CORE.Window.screen.height == 0)) fullscreen = true;

        // Fullscreen/borderless modes are applied after the window is created
        FLAG_CLEAR(CORE.Window.flags, FLAG_FULLSCREEN_MODE);
        FLAG_CLEAR(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE);

        if ((CORE.Window.screen.width == 0) || (CORE.Window.screen.height == 0))
        {
            // Set some default screen size in case user decides to exit fullscreen mode
            CORE.Window.screen.width = 800;
            CORE.Window.screen.height = 450;
        }

        // Create window and content view
        NSRect contentRect = NSMakeRect(0, 0, CORE.Window.screen.width, CORE.Window.screen.height);

        platform.window = [[RaylibWindow alloc] initWithContentRect:contentRect styleMask:GetWindowStyle()
            backing:NSBackingStoreBuffered defer:NO];

        if (platform.window == nil)
        {
            TRACELOG(LOG_WARNING, "COCOA: Failed to create window");
            return -1;
        }

        platform.title = [GetTitleString(CORE.Window.title) retain];
        platform.windowDelegate = [[RaylibWindowDelegate alloc] init];
        platform.view = [[RaylibView alloc] initWithFrame:contentRect];

        [platform.window setReleasedWhenClosed:NO];
        [platform.window setDelegate:platform.windowDelegate];
        [platform.window setContentView:platform.view];
        [platform.window makeFirstResponder:platform.view];
        [platform.window setAcceptsMouseMovedEvents:YES];
        [platform.window setRestorable:NO];
        [platform.window setTitle:platform.title];
        [platform.window setTabbingMode:NSWindowTabbingModeDisallowed];
        [platform.window setCollectionBehavior:(FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_RESIZABLE)?
            (NSWindowCollectionBehaviorFullScreenPrimary | NSWindowCollectionBehaviorManaged) : NSWindowCollectionBehaviorFullScreenNone)];

        if ((CORE.Window.screenMin.width > 0) || (CORE.Window.screenMin.height > 0)) SetWindowMinSize(CORE.Window.screenMin.width, CORE.Window.screenMin.height);
        if ((CORE.Window.screenMax.width > 0) || (CORE.Window.screenMax.height > 0)) SetWindowMaxSize(CORE.Window.screenMax.width, CORE.Window.screenMax.height);

        if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_TOPMOST)) [platform.window setLevel:NSFloatingWindowLevel];
        if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_MOUSE_PASSTHROUGH)) [platform.window setIgnoresMouseEvents:YES];

        if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_TRANSPARENT))
        {
            [platform.window setOpaque:NO];
            [platform.window setHasShadow:NO];
            [platform.window setBackgroundColor:[NSColor clearColor]];
        }

        // HiDPI framebuffer only with FLAG_WINDOW_HIGHDPI
        [platform.view setWantsBestResolutionOpenGLSurface:FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_HIGHDPI)];

        // Center window on the main screen work area (without menu bar and Dock)
        NSRect screenFrame = [mainScreen visibleFrame];
        contentRect.origin.x = screenFrame.origin.x + (screenFrame.size.width - contentRect.size.width)/2;
        contentRect.origin.y = screenFrame.origin.y + (screenFrame.size.height - contentRect.size.height)/2;
        [platform.window setFrameOrigin:[platform.window frameRectForContentRect:contentRect].origin];

        // Initialize graphics context
        if (!InitGraphicsContext())
        {
            TRACELOG(LOG_FATAL, "PLATFORM: Failed to initialize graphics device");
            ClosePlatform();
            return -1;
        }

        UpdateWindowPosition();
        UpdateFramebufferSize();

        // Show window
        platform.focusOnShow = !FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_UNFOCUSED);
        if (!FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_HIDDEN))
        {
            if (!platform.focusOnShow) [platform.window orderFront:nil];
            else
            {
                [NSApp activateIgnoringOtherApps:YES];
                [platform.window makeKeyAndOrderFront:nil];
            }
        }

        CORE.Input.Mouse.cursorOnScreen = IsCursorInContentView();

        // Forward key up events while Command is down, AppKit does not send them to the key window
        platform.keyUpMonitor = [NSEvent addLocalMonitorForEventsMatchingMask:NSEventMaskKeyUp handler:^NSEvent *(NSEvent *event)
        {
            if ([event modifierFlags] & NSEventModifierFlagCommand) [[NSApp keyWindow] sendEvent:event];
            return event;
        }];

        CORE.Window.ready = true;

        // Apply initial window states
        if (fullscreen) ToggleFullscreen();
        else if (borderless) ToggleBorderlessWindowed();

        if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_MAXIMIZED))
        {
            FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_MAXIMIZED);
            MaximizeWindow();
        }

        if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_MINIMIZED))
        {
            FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_MINIMIZED);
            MinimizeWindow();
        }
    }

    TRACELOG(LOG_INFO, "DISPLAY: Device initialized successfully %s",
        FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_HIGHDPI)? "(HighDPI)" : "");
    TRACELOG(LOG_INFO, "    > Display size: %i x %i", CORE.Window.display.width, CORE.Window.display.height);
    TRACELOG(LOG_INFO, "    > Screen size:  %i x %i", CORE.Window.screen.width, CORE.Window.screen.height);
    TRACELOG(LOG_INFO, "    > Render size:  %i x %i", CORE.Window.render.width, CORE.Window.render.height);
    TRACELOG(LOG_INFO, "    > Viewport offsets: %i, %i", CORE.Window.renderOffset.x, CORE.Window.renderOffset.y);

    // Load OpenGL extensions
    // NOTE: GL procedures address loader is required to load extensions
    //----------------------------------------------------------------------------
#if !defined(GRAPHICS_API_OPENGL_SOFTWARE)
    rlLoadExtensions(GetProcAddressCocoa);
#endif
    //----------------------------------------------------------------------------

    // Initialize gamepads
    //----------------------------------------------------------------------------
    mg_gamepads_init(&platform.minigamepad);
    //----------------------------------------------------------------------------

    // Initialize timing system
    //----------------------------------------------------------------------------
    InitTimer();
    //----------------------------------------------------------------------------

    // Initialize storage system
    //----------------------------------------------------------------------------
    CORE.Storage.basePath = GetWorkingDirectory();
    //----------------------------------------------------------------------------

    platform.initialized = true;

#if defined(GRAPHICS_API_OPENGL_SOFTWARE)
    TRACELOG(LOG_INFO, "PLATFORM: DESKTOP (COCOA - Software): Initialized successfully");
#else
    TRACELOG(LOG_INFO, "PLATFORM: DESKTOP (COCOA): Initialized successfully");
#endif

    return 0;
}

// Close platform
void ClosePlatform(void)
{
    if (platform.minigamepad.src.hidManager != NULL) mg_gamepads_free(&platform.minigamepad);

    @autoreleasepool
    {
        LockCursor(false);

        if (platform.keyUpMonitor != nil) [NSEvent removeMonitor:platform.keyUpMonitor];

#if defined(GRAPHICS_API_OPENGL_SOFTWARE)
        RL_FREE(platform.softwarePixels);
        if (platform.colorSpace != NULL) CGColorSpaceRelease(platform.colorSpace);
#elif defined(COCOA_EGL)
        if (platform.eglDisplay != NULL)
        {
            platform.egl.MakeCurrent(platform.eglDisplay, EGL_NO_SURFACE, EGL_NO_SURFACE, EGL_NO_CONTEXT);
            if (platform.eglSurface != EGL_NO_SURFACE) platform.egl.DestroySurface(platform.eglDisplay, platform.eglSurface);
            if (platform.eglContext != EGL_NO_CONTEXT) platform.egl.DestroyContext(platform.eglDisplay, platform.eglContext);
            platform.egl.Terminate(platform.eglDisplay);
        }

        if (platform.eglLibrary != NULL) dlclose(platform.eglLibrary);
#else
        if (platform.glContext != nil)
        {
            [NSOpenGLContext clearCurrentContext];
            [platform.glContext clearDrawable];
            [platform.glContext release];
        }

        if (platform.glLibrary != NULL) dlclose(platform.glLibrary);
#endif

        if (platform.window != nil)
        {
            [platform.window setDelegate:nil];
            [platform.window orderOut:nil];
            [platform.window close];
            [platform.window release];
        }

        if (platform.view != nil) [platform.view release];
        if (platform.windowDelegate != nil) [platform.windowDelegate release];
        if (platform.hiddenCursor != nil) [platform.hiddenCursor release];
        if (platform.title != nil) [platform.title release];

        // Process pending events of the closed window
        while (true)
        {
            NSEvent *event = [NSApp nextEventMatchingMask:NSEventMaskAny untilDate:[NSDate distantPast] inMode:NSDefaultRunLoopMode dequeue:YES];
            if (event == nil) break;
            [NSApp sendEvent:event];
        }
    }

    RL_FREE(platform.clipboardText);

    // NOTE: Application delegate and menu bar are kept for a later InitWindow()
    RaylibApplicationDelegate *appDelegate = platform.appDelegate;
    memset(&platform, 0, sizeof(platform));
    platform.appDelegate = appDelegate;
}

// Initialize shared application, delegate and menu bar
static void InitApplication(void)
{
    [NSApplication sharedApplication];

    if (platform.appDelegate != nil) return;

    platform.appDelegate = [[RaylibApplicationDelegate alloc] init];
    [NSApp setDelegate:platform.appDelegate];

    // NOTE: Unbundled executables are background processes by default, without Dock icon and menu bar
    [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];

    // Holding a key repeats it instead of opening the accents popup
    [[NSUserDefaults standardUserDefaults] registerDefaults:@{ @"ApplePressAndHoldEnabled": @NO }];

    // Change working directory to the bundle resources directory, like GLFW
    NSBundle *bundle = [NSBundle mainBundle];
    if ([[bundle bundlePath] hasSuffix:@".app"] && ([bundle resourcePath] != nil)) chdir([[bundle resourcePath] fileSystemRepresentation]);

    if ([NSApp mainMenu] == nil) CreateMenuBar();

    [NSApp finishLaunching];
}

// Create application menu bar (application and window menus)
static void CreateMenuBar(void)
{
    NSString *appName = [[NSProcessInfo processInfo] processName];
    NSString *bundleName = [[NSBundle mainBundle] objectForInfoDictionaryKey:@"CFBundleName"];
    if ([bundleName length] > 0) appName = bundleName;

    NSMenu *menuBar = [[NSMenu alloc] init];
    [NSApp setMainMenu:menuBar];

    // Application menu
    NSMenuItem *appMenuItem = [menuBar addItemWithTitle:@"" action:NULL keyEquivalent:@""];
    NSMenu *appMenu = [[NSMenu alloc] init];
    [appMenuItem setSubmenu:appMenu];

    [appMenu addItemWithTitle:[NSString stringWithFormat:@"About %@", appName] action:@selector(orderFrontStandardAboutPanel:) keyEquivalent:@""];
    [appMenu addItem:[NSMenuItem separatorItem]];
    [appMenu addItemWithTitle:[NSString stringWithFormat:@"Hide %@", appName] action:@selector(hide:) keyEquivalent:@"h"];
    [[appMenu addItemWithTitle:@"Hide Others" action:@selector(hideOtherApplications:) keyEquivalent:@"h"]
        setKeyEquivalentModifierMask:(NSEventModifierFlagOption | NSEventModifierFlagCommand)];
    [appMenu addItemWithTitle:@"Show All" action:@selector(unhideAllApplications:) keyEquivalent:@""];
    [appMenu addItem:[NSMenuItem separatorItem]];
    [appMenu addItemWithTitle:[NSString stringWithFormat:@"Quit %@", appName] action:@selector(terminate:) keyEquivalent:@"q"];

    // Window menu
    NSMenuItem *windowMenuItem = [menuBar addItemWithTitle:@"" action:NULL keyEquivalent:@""];
    NSMenu *windowMenu = [[NSMenu alloc] initWithTitle:@"Window"];
    [windowMenuItem setSubmenu:windowMenu];
    [NSApp setWindowsMenu:windowMenu];

    [windowMenu addItemWithTitle:@"Minimize" action:@selector(performMiniaturize:) keyEquivalent:@"m"];
    [windowMenu addItemWithTitle:@"Zoom" action:@selector(performZoom:) keyEquivalent:@""];
    [windowMenu addItem:[NSMenuItem separatorItem]];
    [windowMenu addItemWithTitle:@"Bring All to Front" action:@selector(arrangeInFront:) keyEquivalent:@""];
    [windowMenu addItem:[NSMenuItem separatorItem]];
    [[windowMenu addItemWithTitle:@"Enter Full Screen" action:@selector(toggleFullScreen:) keyEquivalent:@"f"]
        setKeyEquivalentModifierMask:(NSEventModifierFlagControl | NSEventModifierFlagCommand)];

    [windowMenu release];
    [appMenu release];
    [menuBar release];
}

#if !defined(GRAPHICS_API_OPENGL_SOFTWARE) && !defined(COCOA_EGL)
// Create OpenGL pixel format, hardware accelerated or using the given software renderer
static NSOpenGLPixelFormat *CreatePixelFormat(NSOpenGLPixelFormatAttribute profile, bool msaa, NSOpenGLPixelFormatAttribute renderer)
{
    NSOpenGLPixelFormatAttribute attributes[32] = { 0 };
    int count = 0;

    if (renderer == 0)
    {
        attributes[count++] = NSOpenGLPFAAccelerated;
        attributes[count++] = NSOpenGLPFAClosestPolicy;
    }
    else
    {
        attributes[count++] = NSOpenGLPFARendererID;
        attributes[count++] = renderer;
    }

    attributes[count++] = NSOpenGLPFADoubleBuffer;
    attributes[count++] = NSOpenGLPFAOpenGLProfile;
    attributes[count++] = profile;
    attributes[count++] = NSOpenGLPFAColorSize;
    attributes[count++] = 24;
    attributes[count++] = NSOpenGLPFAAlphaSize;
    attributes[count++] = 8;
    attributes[count++] = NSOpenGLPFADepthSize;
    attributes[count++] = 24;
    attributes[count++] = NSOpenGLPFAStencilSize;
    attributes[count++] = 8;

    if (msaa)
    {
        attributes[count++] = NSOpenGLPFASampleBuffers;
        attributes[count++] = 1;
        attributes[count++] = NSOpenGLPFASamples;
        attributes[count++] = 4;
        attributes[count++] = NSOpenGLPFAMultisample;
    }

    attributes[count++] = 0;

    return [[NSOpenGLPixelFormat alloc] initWithAttributes:attributes];
}
#endif

// Initialize OpenGL context or software framebuffer
static bool InitGraphicsContext(void)
{
#if defined(GRAPHICS_API_OPENGL_SOFTWARE)
    [platform.view setWantsLayer:YES];
    [[platform.view layer] setContentsGravity:kCAGravityResize];
    [[platform.view layer] setOpaque:!FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_TRANSPARENT)];
    platform.colorSpace = CGColorSpaceCreateDeviceRGB();

    return true;
#elif defined(COCOA_EGL)
    // NOTE: ANGLE uses the content view layer as EGL native window, it adds its own Metal layer
    [platform.view setWantsLayer:YES];
    [[platform.view layer] setOpaque:!FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_TRANSPARENT)];
    [[platform.view layer] setContentsScale:FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_HIGHDPI)? [platform.window backingScaleFactor] : 1.0];

    // Load EGL library: next to the executable, from the executable rpath or from the library search paths
    const char *libraryNames[] = { "@executable_path/libEGL.dylib", "@rpath/libEGL.dylib", "libEGL.dylib" };
    for (int i = 0; (i < 3) && (platform.eglLibrary == NULL); i++) platform.eglLibrary = dlopen(libraryNames[i], RTLD_LAZY | RTLD_LOCAL);

    if (platform.eglLibrary == NULL)
    {
        TRACELOG(LOG_WARNING, "EGL: Failed to load libEGL.dylib (ANGLE), required for OpenGL ES");
        return false;
    }

    platform.egl.GetProcAddress = (void *(*)(const char *))dlsym(platform.eglLibrary, "eglGetProcAddress");
    platform.egl.GetDisplay = (EGLDisplay (*)(void *))dlsym(platform.eglLibrary, "eglGetDisplay");
    platform.egl.Initialize = (EGLBoolean (*)(EGLDisplay, EGLint *, EGLint *))dlsym(platform.eglLibrary, "eglInitialize");
    platform.egl.Terminate = (EGLBoolean (*)(EGLDisplay))dlsym(platform.eglLibrary, "eglTerminate");
    platform.egl.BindAPI = (EGLBoolean (*)(EGLenum))dlsym(platform.eglLibrary, "eglBindAPI");
    platform.egl.ChooseConfig = (EGLBoolean (*)(EGLDisplay, const EGLint *, EGLConfig *, EGLint, EGLint *))dlsym(platform.eglLibrary, "eglChooseConfig");
    platform.egl.CreateWindowSurface = (EGLSurface (*)(EGLDisplay, EGLConfig, void *, const EGLint *))dlsym(platform.eglLibrary, "eglCreateWindowSurface");
    platform.egl.DestroySurface = (EGLBoolean (*)(EGLDisplay, EGLSurface))dlsym(platform.eglLibrary, "eglDestroySurface");
    platform.egl.CreateContext = (EGLContext (*)(EGLDisplay, EGLConfig, EGLContext, const EGLint *))dlsym(platform.eglLibrary, "eglCreateContext");
    platform.egl.DestroyContext = (EGLBoolean (*)(EGLDisplay, EGLContext))dlsym(platform.eglLibrary, "eglDestroyContext");
    platform.egl.MakeCurrent = (EGLBoolean (*)(EGLDisplay, EGLSurface, EGLSurface, EGLContext))dlsym(platform.eglLibrary, "eglMakeCurrent");
    platform.egl.SwapBuffers = (EGLBoolean (*)(EGLDisplay, EGLSurface))dlsym(platform.eglLibrary, "eglSwapBuffers");
    platform.egl.SwapInterval = (EGLBoolean (*)(EGLDisplay, EGLint))dlsym(platform.eglLibrary, "eglSwapInterval");
    platform.egl.GetError = (EGLint (*)(void))dlsym(platform.eglLibrary, "eglGetError");

    if ((platform.egl.GetProcAddress == NULL) || (platform.egl.GetDisplay == NULL) || (platform.egl.Initialize == NULL) ||
        (platform.egl.Terminate == NULL) || (platform.egl.BindAPI == NULL) || (platform.egl.ChooseConfig == NULL) ||
        (platform.egl.CreateWindowSurface == NULL) || (platform.egl.DestroySurface == NULL) || (platform.egl.CreateContext == NULL) ||
        (platform.egl.DestroyContext == NULL) || (platform.egl.MakeCurrent == NULL) || (platform.egl.SwapBuffers == NULL) ||
        (platform.egl.GetError == NULL))
    {
        TRACELOG(LOG_WARNING, "EGL: Failed to load EGL functions");
        return false;
    }

    // Request ANGLE Metal backend
    platform.egl.GetPlatformDisplayEXT = (EGLDisplay (*)(EGLenum, void *, const EGLint *))platform.egl.GetProcAddress("eglGetPlatformDisplayEXT");
    if (platform.egl.GetPlatformDisplayEXT != NULL)
    {
        const EGLint displayAttribs[] = { EGL_PLATFORM_ANGLE_TYPE_ANGLE, EGL_PLATFORM_ANGLE_TYPE_METAL_ANGLE, EGL_NONE };
        platform.eglDisplay = platform.egl.GetPlatformDisplayEXT(EGL_PLATFORM_ANGLE_ANGLE, EGL_DEFAULT_DISPLAY, displayAttribs);
    }

    if (platform.eglDisplay == NULL) platform.eglDisplay = platform.egl.GetDisplay(EGL_DEFAULT_DISPLAY);

    EGLint major = 0;
    EGLint minor = 0;
    if ((platform.eglDisplay == NULL) || !platform.egl.Initialize(platform.eglDisplay, &major, &minor))
    {
        TRACELOG(LOG_WARNING, "EGL: Failed to initialize EGL display [ERROR: 0x%x]", platform.egl.GetError());
        return false;
    }

    platform.egl.BindAPI(EGL_OPENGL_ES_API);

    bool msaa = FLAG_IS_SET(CORE.Window.flags, FLAG_MSAA_4X_HINT);
    const EGLint configAttribs[] = {
        EGL_SURFACE_TYPE, EGL_WINDOW_BIT,
        EGL_RENDERABLE_TYPE, (rlGetVersion() == RL_OPENGL_ES_30)? EGL_OPENGL_ES3_BIT : EGL_OPENGL_ES2_BIT,
        EGL_RED_SIZE, 8,
        EGL_GREEN_SIZE, 8,
        EGL_BLUE_SIZE, 8,
        EGL_ALPHA_SIZE, 8,
        EGL_DEPTH_SIZE, 24,
        EGL_STENCIL_SIZE, 8,
        EGL_SAMPLE_BUFFERS, msaa? 1 : 0,
        EGL_SAMPLES, msaa? 4 : 0,
        EGL_NONE
    };

    EGLConfig config = NULL;
    EGLint configCount = 0;
    if (!platform.egl.ChooseConfig(platform.eglDisplay, configAttribs, &config, 1, &configCount) || (configCount < 1))
    {
        TRACELOG(LOG_WARNING, "EGL: Failed to choose a suitable framebuffer config [ERROR: 0x%x]", platform.egl.GetError());
        return false;
    }

    platform.eglSurface = platform.egl.CreateWindowSurface(platform.eglDisplay, config, (void *)[platform.view layer], NULL);
    if (platform.eglSurface == EGL_NO_SURFACE)
    {
        TRACELOG(LOG_WARNING, "EGL: Failed to create window surface [ERROR: 0x%x]", platform.egl.GetError());
        return false;
    }

    const EGLint contextAttribs[] = { EGL_CONTEXT_CLIENT_VERSION, (rlGetVersion() == RL_OPENGL_ES_30)? 3 : 2, EGL_NONE };
    platform.eglContext = platform.egl.CreateContext(platform.eglDisplay, config, EGL_NO_CONTEXT, contextAttribs);
    if (platform.eglContext == EGL_NO_CONTEXT)
    {
        TRACELOG(LOG_WARNING, "EGL: Failed to create OpenGL ES context [ERROR: 0x%x]", platform.egl.GetError());
        return false;
    }

    if (!platform.egl.MakeCurrent(platform.eglDisplay, platform.eglSurface, platform.eglSurface, platform.eglContext))
    {
        TRACELOG(LOG_WARNING, "EGL: Failed to make context current [ERROR: 0x%x]", platform.egl.GetError());
        return false;
    }

    TRACELOG(LOG_INFO, "EGL: Initialized EGL %i.%i (ANGLE)", major, minor);

    SetSwapInterval(FLAG_IS_SET(CORE.Window.flags, FLAG_VSYNC_HINT)? 1 : 0);

    return true;
#else
    NSOpenGLPixelFormatAttribute profile = 0;

    // NOTE: macOS provides OpenGL 4.1 core as the highest version, compatible with OpenGL 3.3 core requests
    switch (rlGetVersion())
    {
        case RL_OPENGL_11:
        case RL_OPENGL_21: profile = NSOpenGLProfileVersionLegacy; break;
        case RL_OPENGL_33: profile = NSOpenGLProfileVersion4_1Core; break;
        case RL_OPENGL_43:
        {
            TRACELOG(LOG_WARNING, "GL: OpenGL 4.3 not available on macOS, using OpenGL 4.1 core");
            profile = NSOpenGLProfileVersion4_1Core;
        } break;
        default: return false;
    }

    bool msaa = FLAG_IS_SET(CORE.Window.flags, FLAG_MSAA_4X_HINT);
    NSOpenGLPixelFormat *pixelFormat = CreatePixelFormat(profile, msaa, 0);

    if (pixelFormat == nil)
    {
        const NSOpenGLPixelFormatAttribute renderers[] = { kCGLRendererGenericFloatID, kCGLRendererAppleSWID };

        for (int i = 0; (i < 2) && (pixelFormat == nil); i++)
        {
            pixelFormat = CreatePixelFormat(profile, msaa, renderers[i]);
            if ((pixelFormat == nil) && msaa) pixelFormat = CreatePixelFormat(profile, false, renderers[i]);
        }

        if (pixelFormat != nil) TRACELOG(LOG_WARNING, "GL: Hardware acceleration not available, using software renderer");
    }

    if (pixelFormat == nil)
    {
        TRACELOG(LOG_WARNING, "GL: Failed to find a suitable pixel format");
        return false;
    }

    platform.glContext = [[NSOpenGLContext alloc] initWithFormat:pixelFormat shareContext:nil];
    [pixelFormat release];

    if (platform.glContext == nil)
    {
        TRACELOG(LOG_WARNING, "GL: Failed to create OpenGL context");
        return false;
    }

    if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_TRANSPARENT))
    {
        GLint opacity = 0;
        [platform.glContext setValues:&opacity forParameter:NSOpenGLContextParameterSurfaceOpacity];
    }

    [platform.glContext setView:platform.view];
    [platform.glContext makeCurrentContext];

    SetSwapInterval(FLAG_IS_SET(CORE.Window.flags, FLAG_VSYNC_HINT)? 1 : 0);

    platform.glLibrary = dlopen("/System/Library/Frameworks/OpenGL.framework/OpenGL", RTLD_LAZY | RTLD_LOCAL);
    if (platform.glLibrary == NULL)
    {
        TRACELOG(LOG_WARNING, "GL: Failed to load OpenGL framework");
        return false;
    }

    return true;
#endif
}

// Update screen and render sizes from the content view
// NOTE: Screen size is in points (mouse coordinates), render size is the framebuffer size in pixels
static void UpdateFramebufferSize(void)
{
    if (platform.view == nil) return;

    NSRect bounds = [platform.view bounds];
    NSRect backing = FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_HIGHDPI)? [platform.view convertRectToBacking:bounds] : bounds;

    int width = (int)bounds.size.width;
    int height = (int)bounds.size.height;
    int renderWidth = (int)backing.size.width;
    int renderHeight = (int)backing.size.height;

    // NOTE: Minimized windows keep their size
    if ((width <= 0) || (height <= 0) || (renderWidth <= 0) || (renderHeight <= 0)) return;

#if defined(COCOA_EGL) || defined(GRAPHICS_API_OPENGL_SOFTWARE)
    [[platform.view layer] setContentsScale:(CGFloat)renderWidth/(CGFloat)width];
#else
    [platform.glContext update];
#endif

    if (FLAG_IS_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE) || FLAG_IS_SET(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE))
    {
        CORE.Window.display.width = width;
        CORE.Window.display.height = height;
    }

    bool resized = (((int)CORE.Window.render.width != renderWidth) || ((int)CORE.Window.render.height != renderHeight) ||
        ((int)CORE.Window.screen.width != width) || ((int)CORE.Window.screen.height != height));

    CORE.Window.screen.width = width;
    CORE.Window.screen.height = height;
    CORE.Window.screenScale = MatrixScale((float)renderWidth/(float)width, (float)renderHeight/(float)height, 1.0f);

    // Reset viewport and projection matrix for new size
    // NOTE: Stores current render size: CORE.Window.render
    if (platform.initialized && resized) SetupViewport(renderWidth, renderHeight);

    CORE.Window.render.width = renderWidth;
    CORE.Window.render.height = renderHeight;
    CORE.Window.currentFbo = CORE.Window.render;
    if (platform.initialized && resized) CORE.Window.resizedLastFrame = true;
}

// Update window position (content top-left, screen coordinates)
static void UpdateWindowPosition(void)
{
    NSRect contentRect = [platform.window contentRectForFrameRect:[platform.window frame]];

    CORE.Window.position.x = (int)contentRect.origin.x;
    CORE.Window.position.y = (int)(GetPrimaryScreenHeight() - (contentRect.origin.y + contentRect.size.height));
}

// Get window style mask from window flags
static NSWindowStyleMask GetWindowStyle(void)
{
    NSWindowStyleMask style = NSWindowStyleMaskMiniaturizable;

    if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_UNDECORATED)) style |= NSWindowStyleMaskBorderless;
    else style |= (NSWindowStyleMaskTitled | NSWindowStyleMaskClosable);

    if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_RESIZABLE)) style |= NSWindowStyleMaskResizable;

    return style;
}

// Check if window is zoomed (maximized)
// NOTE: isZoomed reports non-resizable windows as zoomed, they can not be maximized
static bool IsWindowZoomed(void)
{
    return (([platform.window styleMask] & NSWindowStyleMaskResizable) && [platform.window isZoomed]);
}

// Get title string, empty for NULL or invalid UTF-8
static NSString *GetTitleString(const char *title)
{
    NSString *string = (title != NULL)? [NSString stringWithUTF8String:title] : nil;

    return (string != nil)? string : @"";
}

// Apply window style keeping the content size
static void UpdateWindowStyle(void)
{
    [platform.window setCollectionBehavior:(FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_RESIZABLE)?
        (NSWindowCollectionBehaviorFullScreenPrimary | NSWindowCollectionBehaviorManaged) : NSWindowCollectionBehaviorFullScreenNone)];

    // Fullscreen/borderless windows keep the borderless style until restored
    if (FLAG_IS_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE) || FLAG_IS_SET(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE)) return;

    NSRect contentRect = [platform.window contentRectForFrameRect:[platform.window frame]];

    [platform.window setStyleMask:GetWindowStyle()];
    [platform.window setFrame:[platform.window frameRectForContentRect:contentRect] display:YES];

    // Restore title and first responder, cleared by the style mask change
    [platform.window setTitle:platform.title];
    [platform.window makeFirstResponder:platform.view];

    UpdateFramebufferSize();
}

// Cover the current monitor with a borderless window, or restore it
static void SetWindowFullscreenFrame(bool enabled)
{
    if (enabled)
    {
        NSScreen *screen = [platform.window screen];
        if (screen == nil) screen = [NSScreen mainScreen];

        platform.windowedContent = [platform.window contentRectForFrameRect:[platform.window frame]];

        [platform.window setStyleMask:NSWindowStyleMaskBorderless];
        [platform.window setFrame:[screen frame] display:YES];
        [platform.window setLevel:NSMainMenuWindowLevel + 1];
        [platform.window setHasShadow:NO];
    }
    else
    {
        [platform.window setStyleMask:GetWindowStyle()];
        [platform.window setFrame:[platform.window frameRectForContentRect:platform.windowedContent] display:YES];
        [platform.window setLevel:FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_TOPMOST)? NSFloatingWindowLevel : NSNormalWindowLevel];
        [platform.window setHasShadow:!FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_TRANSPARENT)];
        [platform.window setTitle:platform.title];
    }

    // Restore first responder, cleared by the style mask change
    [platform.window makeFirstResponder:platform.view];
    if (!FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_HIDDEN)) [platform.window makeKeyAndOrderFront:nil];

    UpdateWindowPosition();
}

// Get screen for a monitor index
static NSScreen *GetMonitorScreen(int monitor)
{
    NSArray *screens = [NSScreen screens];

    if ((monitor >= 0) && (monitor < (int)[screens count])) return [screens objectAtIndex:monitor];

    return nil;
}

// Get display id for a screen
static CGDirectDisplayID GetScreenDisplayID(NSScreen *screen)
{
    return [[[screen deviceDescription] objectForKey:@"NSScreenNumber"] unsignedIntValue];
}

// Get primary screen height, required to flip coordinates
// NOTE: Cocoa screen coordinates origin is the bottom-left corner of the primary screen
static float GetPrimaryScreenHeight(void)
{
    NSArray *screens = [NSScreen screens];

    return ([screens count] > 0)? (float)[[screens objectAtIndex:0] frame].size.height : 0.0f;
}

// Get system cursor for a raylib mouse cursor
static NSCursor *GetMouseCursor(int cursor)
{
    NSCursor *result = nil;
    SEL selector = NULL;

    switch (cursor)
    {
        case MOUSE_CURSOR_IBEAM: result = [NSCursor IBeamCursor]; break;
        case MOUSE_CURSOR_CROSSHAIR: result = [NSCursor crosshairCursor]; break;
        case MOUSE_CURSOR_POINTING_HAND: result = [NSCursor pointingHandCursor]; break;
        case MOUSE_CURSOR_RESIZE_EW: result = [NSCursor resizeLeftRightCursor]; break;
        case MOUSE_CURSOR_RESIZE_NS: result = [NSCursor resizeUpDownCursor]; break;
        case MOUSE_CURSOR_RESIZE_ALL: result = [NSCursor closedHandCursor]; break;
        case MOUSE_CURSOR_NOT_ALLOWED: result = [NSCursor operationNotAllowedCursor]; break;

        // NOTE: Diagonal resize cursors are not public, used when available like GLFW
        case MOUSE_CURSOR_RESIZE_NWSE: selector = NSSelectorFromString(@"_windowResizeNorthWestSouthEastCursor"); break;
        case MOUSE_CURSOR_RESIZE_NESW: selector = NSSelectorFromString(@"_windowResizeNorthEastSouthWestCursor"); break;
        default: break;
    }

    if ((selector != NULL) && [NSCursor respondsToSelector:selector])
    {
        id object = [NSCursor performSelector:selector];
        if ([object isKindOfClass:[NSCursor class]]) result = object;
    }

    if (result == nil) result = [NSCursor arrowCursor];

    return result;
}

// Check if mouse cursor is inside the content view
// NOTE: Tracking areas do not report the cursor already inside when the window appears
static bool IsCursorInContentView(void)
{
    if (platform.view == nil) return false;

    NSPoint point = [platform.view convertPoint:[platform.window mouseLocationOutsideOfEventStream] fromView:nil];

    return NSMouseInRect(point, [platform.view bounds], [platform.view isFlipped]);
}

// Set current cursor, if inside the content view
static void UpdateCursor(void)
{
    if (!IsCursorInContentView()) return;

    if (CORE.Input.Mouse.cursorHidden)
    {
        if (platform.hiddenCursor == nil)
        {
            NSImage *image = [[NSImage alloc] initWithSize:NSMakeSize(16, 16)];
            platform.hiddenCursor = [[NSCursor alloc] initWithImage:image hotSpot:NSZeroPoint];
            [image release];
        }

        [platform.hiddenCursor set];
    }
    else [GetMouseCursor(CORE.Input.Mouse.cursor) set];
}

// Move the system cursor to a content view position
static void WarpCursor(int x, int y)
{
    if (platform.window == nil) return;

    NSRect contentRect = [platform.window contentRectForFrameRect:[platform.window frame]];
    CGPoint point = CGPointMake(contentRect.origin.x + x, GetPrimaryScreenHeight() - (contentRect.origin.y + contentRect.size.height) + y);

    // Warp distance, subtracted from the next mouse event delta
    NSPoint current = [platform.view convertPoint:[platform.window mouseLocationOutsideOfEventStream] fromView:nil];
    platform.cursorWarpDelta.x += (float)x - (float)current.x;
    platform.cursorWarpDelta.y += (float)y - (float)([platform.view bounds].size.height - current.y);

    CGWarpMouseCursorPosition(point);

    // Reassociate mouse and cursor after warping
    if (!CORE.Input.Mouse.cursorLocked) CGAssociateMouseAndMouseCursorPosition(true);
}

// Lock/unlock cursor (disabled cursor mode)
// NOTE: Locked cursor stays hidden in the content view center, mouse deltas move the virtual cursor position
static void LockCursor(bool lock)
{
    if (lock)
    {
        if (!platform.cursorSystemHidden)
        {
            [NSCursor hide];
            platform.cursorSystemHidden = true;
        }

        NSRect bounds = [platform.view bounds];
        WarpCursor((int)(bounds.size.width/2), (int)(bounds.size.height/2));
        CGAssociateMouseAndMouseCursorPosition(false);
    }
    else
    {
        CGAssociateMouseAndMouseCursorPosition(true);

        if (platform.cursorSystemHidden)
        {
            [NSCursor unhide];
            platform.cursorSystemHidden = false;
        }
    }
}

// Register key state change
static void HandleKey(int key, bool pressed, bool repeat)
{
    if ((key <= KEY_NULL) || (key >= MAX_KEYBOARD_KEYS)) return;

    if (pressed)
    {
        if (repeat) CORE.Input.Keyboard.keyRepeatInFrame[key] = 1;
        else
        {
            CORE.Input.Keyboard.currentKeyState[key] = 1;

            // Check if there is space available in the key queue
            if (CORE.Input.Keyboard.keyPressedQueueCount < MAX_KEY_PRESSED_QUEUE)
            {
                // Add key to the queue
                CORE.Input.Keyboard.keyPressedQueue[CORE.Input.Keyboard.keyPressedQueueCount] = key;
                CORE.Input.Keyboard.keyPressedQueueCount++;
            }

            // Check the exit key to set close window
            if (key == CORE.Input.Keyboard.exitKey) platform.closeRequested = true;
        }
    }
    else CORE.Input.Keyboard.currentKeyState[key] = 0;
}

// Register modifier key state change
// NOTE: CAPS lock key is considered down while the lock is active
static void HandleModifiers(unsigned short keyCode, NSUInteger flags)
{
    int key = (keyCode < KEYMAP_SIZE)? macosToRaylibMap[keyCode] : (int)KEY_NULL;
    NSUInteger mask = 0;            // Device-dependent mask of the key
    NSUInteger sideMasks = 0;       // Device-dependent masks of both keys (left and right)
    NSUInteger modifier = 0;        // Device-independent modifier flag

    switch (key)
    {
        case KEY_LEFT_SHIFT: mask = DEVICE_LSHIFT_KEY_MASK; break;
        case KEY_RIGHT_SHIFT: mask = DEVICE_RSHIFT_KEY_MASK; break;
        case KEY_LEFT_CONTROL: mask = DEVICE_LCTL_KEY_MASK; break;
        case KEY_RIGHT_CONTROL: mask = DEVICE_RCTL_KEY_MASK; break;
        case KEY_LEFT_ALT: mask = DEVICE_LALT_KEY_MASK; break;
        case KEY_RIGHT_ALT: mask = DEVICE_RALT_KEY_MASK; break;
        case KEY_LEFT_SUPER: mask = DEVICE_LCMD_KEY_MASK; break;
        case KEY_RIGHT_SUPER: mask = DEVICE_RCMD_KEY_MASK; break;
        default: break;
    }

    switch (key)
    {
        case KEY_LEFT_SHIFT: case KEY_RIGHT_SHIFT: sideMasks = DEVICE_LSHIFT_KEY_MASK | DEVICE_RSHIFT_KEY_MASK; modifier = NSEventModifierFlagShift; break;
        case KEY_LEFT_CONTROL: case KEY_RIGHT_CONTROL: sideMasks = DEVICE_LCTL_KEY_MASK | DEVICE_RCTL_KEY_MASK; modifier = NSEventModifierFlagControl; break;
        case KEY_LEFT_ALT: case KEY_RIGHT_ALT: sideMasks = DEVICE_LALT_KEY_MASK | DEVICE_RALT_KEY_MASK; modifier = NSEventModifierFlagOption; break;
        case KEY_LEFT_SUPER: case KEY_RIGHT_SUPER: sideMasks = DEVICE_LCMD_KEY_MASK | DEVICE_RCMD_KEY_MASK; modifier = NSEventModifierFlagCommand; break;
        case KEY_CAPS_LOCK: mask = NSEventModifierFlagCapsLock; sideMasks = mask; modifier = mask; break;
        default: break;
    }

    if (modifier != 0)
    {
        // Synthesized events can miss the device-dependent flags
        bool pressed = ((flags & sideMasks) != 0)? ((flags & mask) != 0) : ((flags & modifier) != 0);
        if (pressed != (CORE.Input.Keyboard.currentKeyState[key] != 0)) HandleKey(key, pressed, false);
    }
}

// Register characters input
static void HandleChars(NSString *text)
{
    NSUInteger length = [text length];
    NSUInteger index = 0;

    while (index < length)
    {
        unsigned int codepoint = [text characterAtIndex:index];
        index++;

        // Combine UTF-16 surrogate pairs
        if ((codepoint >= 0xD800) && (codepoint <= 0xDBFF) && (index < length))
        {
            unsigned int low = [text characterAtIndex:index];
            if ((low >= 0xDC00) && (low <= 0xDFFF))
            {
                codepoint = 0x10000 + ((codepoint - 0xD800) << 10) + (low - 0xDC00);
                index++;
            }
        }

        // NOTE: Control characters and function keys (Unicode private range used by AppKit) are not registered
        if ((codepoint < 32) || (codepoint == 127) || ((codepoint >= 0xF700) && (codepoint <= 0xF8FF))) continue;

        // Check if there is space available in the queue
        if (CORE.Input.Keyboard.charPressedQueueCount < MAX_CHAR_PRESSED_QUEUE)
        {
            // Add character to the queue
            CORE.Input.Keyboard.charPressedQueue[CORE.Input.Keyboard.charPressedQueueCount] = codepoint;
            CORE.Input.Keyboard.charPressedQueueCount++;
        }
    }
}

// Register mouse button state change
static void HandleMouseButton(int button, bool pressed)
{
    if ((button < 0) || (button >= MAX_MOUSE_BUTTONS)) return;

    CORE.Input.Mouse.currentButtonState[button] = pressed;
    CORE.Input.Touch.currentTouchState[button] = pressed;

#if SUPPORT_GESTURES_SYSTEM && SUPPORT_MOUSE_GESTURES
    // Process mouse events as touches to be able to use mouse-gestures
    GestureEvent gestureEvent = { 0 };

    // Register touch actions
    if ((CORE.Input.Mouse.currentButtonState[button] == 1) && (CORE.Input.Mouse.previousButtonState[button] == 0)) gestureEvent.touchAction = TOUCH_ACTION_DOWN;
    else if ((CORE.Input.Mouse.currentButtonState[button] == 0) && (CORE.Input.Mouse.previousButtonState[button] == 1)) gestureEvent.touchAction = TOUCH_ACTION_UP;

    // NOTE: TOUCH_ACTION_MOVE event is registered in HandleMouseMove()

    // Assign a pointer ID
    gestureEvent.pointId[0] = 0;

    // Register touch points count
    gestureEvent.pointCount = 1;

    // Register touch points position, only one point registered
    gestureEvent.position[0] = GetMousePosition();

    // Normalize gestureEvent.position[0] for CORE.Window.screen.width and CORE.Window.screen.height
    gestureEvent.position[0].x /= (float)GetScreenWidth();
    gestureEvent.position[0].y /= (float)GetScreenHeight();

    // Gesture data is sent to gestures-system for processing
    ProcessGestureEvent(gestureEvent);
#endif
}

// Register mouse position change
static void HandleMouseMove(NSEvent *event)
{
    if (CORE.Input.Mouse.cursorLocked)
    {
        // NOTE: Locked cursor does not move, deltas move the virtual cursor position
        CORE.Input.Mouse.currentPosition.x += (float)[event deltaX] - platform.cursorWarpDelta.x;
        CORE.Input.Mouse.currentPosition.y += (float)[event deltaY] - platform.cursorWarpDelta.y;
    }
    else
    {
        NSPoint point = [platform.view convertPoint:[event locationInWindow] fromView:nil];
        CORE.Input.Mouse.currentPosition.x = (float)point.x;
        CORE.Input.Mouse.currentPosition.y = (float)([platform.view bounds].size.height - point.y);
    }

    platform.cursorWarpDelta = (Vector2){ 0.0f, 0.0f };
    CORE.Input.Touch.position[0] = CORE.Input.Mouse.currentPosition;

#if SUPPORT_GESTURES_SYSTEM && SUPPORT_MOUSE_GESTURES
    // Process mouse events as touches to be able to use mouse-gestures
    GestureEvent gestureEvent = { 0 };

    gestureEvent.touchAction = TOUCH_ACTION_MOVE;

    // Assign a pointer ID
    gestureEvent.pointId[0] = 0;

    // Register touch points count
    gestureEvent.pointCount = 1;

    // Register touch points position, only one point registered
    gestureEvent.position[0] = CORE.Input.Mouse.currentPosition;

    // Normalize gestureEvent.position[0] for CORE.Window.screen.width and CORE.Window.screen.height
    gestureEvent.position[0].x /= (float)GetScreenWidth();
    gestureEvent.position[0].y /= (float)GetScreenHeight();

    // Gesture data is sent to gestures-system for processing
    ProcessGestureEvent(gestureEvent);
#endif
}

// Register dropped files
static void HandleDroppedFiles(NSArray *urls)
{
    if ([urls count] == 0) return;

    // In case previous dropped filepaths have not been freed, free them
    if (CORE.Window.dropFileCount > 0)
    {
        for (unsigned int i = 0; i < CORE.Window.dropFileCount; i++) RL_FREE(CORE.Window.dropFilepaths[i]);

        RL_FREE(CORE.Window.dropFilepaths);

        CORE.Window.dropFileCount = 0;
        CORE.Window.dropFilepaths = NULL;
    }

    CORE.Window.dropFilepaths = (char **)RL_CALLOC([urls count], sizeof(char *));

    for (NSURL *url in urls)
    {
        const char *path = [url fileSystemRepresentation];
        if (path == NULL) continue;

        CORE.Window.dropFilepaths[CORE.Window.dropFileCount] = (char *)RL_CALLOC(MAX_FILEPATH_LENGTH, sizeof(char));
        snprintf(CORE.Window.dropFilepaths[CORE.Window.dropFileCount], MAX_FILEPATH_LENGTH, "%s", path);
        CORE.Window.dropFileCount++;
    }
}

// Set swap interval (V-Sync)
static void SetSwapInterval(int interval)
{
#if defined(COCOA_EGL)
    if ((platform.eglDisplay != NULL) && (platform.egl.SwapInterval != NULL)) platform.egl.SwapInterval(platform.eglDisplay, interval);
#elif !defined(GRAPHICS_API_OPENGL_SOFTWARE)
    GLint value = interval;
    [platform.glContext setValues:&value forParameter:NSOpenGLContextParameterSwapInterval];
#endif
}

// Get OpenGL procedure address
static void *GetProcAddressCocoa(const char *name)
{
#if defined(COCOA_EGL)
    return platform.egl.GetProcAddress(name);
#elif !defined(GRAPHICS_API_OPENGL_SOFTWARE)
    return dlsym(platform.glLibrary, name);
#else
    return NULL;
#endif
}

//----------------------------------------------------------------------------------
// Module Internal Classes Implementation
//----------------------------------------------------------------------------------
@implementation RaylibApplicationDelegate

// Quit requested (menu, Cmd+Q, Dock), the application closes through CloseWindow()
- (NSApplicationTerminateReply)applicationShouldTerminate:(NSApplication *)sender
{
    platform.closeRequested = true;

    return NSTerminateCancel;
}

@end

@implementation RaylibWindow

// NOTE: Borderless windows can not become key or main windows by default
- (BOOL)canBecomeKeyWindow
{
    return YES;
}

- (BOOL)canBecomeMainWindow
{
    return YES;
}

@end

@implementation RaylibWindowDelegate

- (BOOL)windowShouldClose:(id)sender
{
    platform.closeRequested = true;

    return NO;
}

- (void)windowDidResize:(NSNotification *)notification
{
    if (IsWindowZoomed()) FLAG_SET(CORE.Window.flags, FLAG_WINDOW_MAXIMIZED);
    else if (!FLAG_IS_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE) && !FLAG_IS_SET(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE)) FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_MAXIMIZED);

    UpdateFramebufferSize();
    UpdateWindowPosition();

    if (CORE.Input.Mouse.cursorLocked && [platform.window isKeyWindow]) LockCursor(true);
}

- (void)windowDidMove:(NSNotification *)notification
{
    UpdateWindowPosition();

    if (CORE.Input.Mouse.cursorLocked && [platform.window isKeyWindow]) LockCursor(true);
}

- (void)windowDidChangeBackingProperties:(NSNotification *)notification
{
    UpdateFramebufferSize();
}

- (void)windowDidMiniaturize:(NSNotification *)notification
{
    FLAG_SET(CORE.Window.flags, FLAG_WINDOW_MINIMIZED);
}

- (void)windowDidDeminiaturize:(NSNotification *)notification
{
    FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_MINIMIZED);
}

- (void)windowDidBecomeKey:(NSNotification *)notification
{
    FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_UNFOCUSED);

    if (CORE.Input.Mouse.cursorLocked) LockCursor(true);
}

- (void)windowDidResignKey:(NSNotification *)notification
{
    FLAG_SET(CORE.Window.flags, FLAG_WINDOW_UNFOCUSED);

    // NOTE: Keys and buttons released outside the window are not reported
    memset(CORE.Input.Keyboard.currentKeyState, 0, sizeof(CORE.Input.Keyboard.currentKeyState));
    memset(CORE.Input.Mouse.currentButtonState, 0, sizeof(CORE.Input.Mouse.currentButtonState));
    memset(CORE.Input.Touch.currentTouchState, 0, sizeof(CORE.Input.Touch.currentTouchState));

    // Give the cursor back to the system while the window is not focused
    if (CORE.Input.Mouse.cursorLocked) LockCursor(false);
}

- (void)windowDidChangeOcclusionState:(NSNotification *)notification
{
    platform.occluded = !([platform.window occlusionState] & NSWindowOcclusionStateVisible);
}

@end

@implementation RaylibView

- (instancetype)initWithFrame:(NSRect)frame
{
    self = [super initWithFrame:frame];

    if (self != nil)
    {
        markedText = [[NSMutableAttributedString alloc] init];
        [self updateTrackingAreas];
        [self registerForDraggedTypes:@[ NSPasteboardTypeFileURL ]];
    }

    return self;
}

- (void)dealloc
{
    [trackingArea release];
    [markedText release];
    [super dealloc];
}

- (BOOL)isOpaque
{
    return !FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_TRANSPARENT);
}

- (BOOL)canBecomeKeyView
{
    return YES;
}

- (BOOL)acceptsFirstResponder
{
    return YES;
}

- (BOOL)acceptsFirstMouse:(NSEvent *)event
{
    return YES;
}

- (BOOL)wantsUpdateLayer
{
    return YES;
}

- (void)updateLayer
{
#if !defined(GRAPHICS_API_OPENGL_SOFTWARE) && !defined(COCOA_EGL)
    [platform.glContext update];
#endif
}

- (void)updateTrackingAreas
{
    if (trackingArea != nil)
    {
        [self removeTrackingArea:trackingArea];
        [trackingArea release];
    }

    NSTrackingAreaOptions options = NSTrackingMouseEnteredAndExited | NSTrackingActiveInKeyWindow |
        NSTrackingEnabledDuringMouseDrag | NSTrackingCursorUpdate | NSTrackingInVisibleRect | NSTrackingAssumeInside;

    trackingArea = [[NSTrackingArea alloc] initWithRect:[self bounds] options:options owner:self userInfo:nil];
    [self addTrackingArea:trackingArea];
    [super updateTrackingAreas];
}

- (void)cursorUpdate:(NSEvent *)event
{
    UpdateCursor();
}

- (void)mouseEntered:(NSEvent *)event
{
    CORE.Input.Mouse.cursorOnScreen = true;
    UpdateCursor();
}

- (void)mouseExited:(NSEvent *)event
{
    CORE.Input.Mouse.cursorOnScreen = false;
}

- (void)mouseMoved:(NSEvent *)event
{
    HandleMouseMove(event);
}

- (void)mouseDragged:(NSEvent *)event
{
    HandleMouseMove(event);
}

- (void)rightMouseDragged:(NSEvent *)event
{
    HandleMouseMove(event);
}

- (void)otherMouseDragged:(NSEvent *)event
{
    HandleMouseMove(event);
}

- (void)mouseDown:(NSEvent *)event
{
    HandleMouseButton(MOUSE_BUTTON_LEFT, true);
}

- (void)mouseUp:(NSEvent *)event
{
    HandleMouseButton(MOUSE_BUTTON_LEFT, false);
}

- (void)rightMouseDown:(NSEvent *)event
{
    HandleMouseButton(MOUSE_BUTTON_RIGHT, true);
}

- (void)rightMouseUp:(NSEvent *)event
{
    HandleMouseButton(MOUSE_BUTTON_RIGHT, false);
}

// NOTE: Other buttons: 2 is middle, 3 and 4 are side buttons (MOUSE_BUTTON_SIDE, MOUSE_BUTTON_EXTRA)...
- (void)otherMouseDown:(NSEvent *)event
{
    HandleMouseButton((int)[event buttonNumber], true);
}

- (void)otherMouseUp:(NSEvent *)event
{
    HandleMouseButton((int)[event buttonNumber], false);
}

- (void)scrollWheel:(NSEvent *)event
{
    double deltaX = [event scrollingDeltaX];
    double deltaY = [event scrollingDeltaY];

    if ([event hasPreciseScrollingDeltas])
    {
        deltaX *= PRECISE_SCROLL_FACTOR;
        deltaY *= PRECISE_SCROLL_FACTOR;
    }

    CORE.Input.Mouse.currentWheelMove.x += (float)deltaX;
    CORE.Input.Mouse.currentWheelMove.y += (float)deltaY;
}

- (void)keyDown:(NSEvent *)event
{
    int key = ([event keyCode] < KEYMAP_SIZE)? macosToRaylibMap[[event keyCode]] : (int)KEY_NULL;
    HandleKey(key, true, [event isARepeat]);

    // Characters input, through the input method (dead keys, compose)
    [self interpretKeyEvents:@[ event ]];
}

- (void)keyUp:(NSEvent *)event
{
    int key = ([event keyCode] < KEYMAP_SIZE)? macosToRaylibMap[[event keyCode]] : (int)KEY_NULL;
    HandleKey(key, false, false);
}

- (void)flagsChanged:(NSEvent *)event
{
    HandleModifiers([event keyCode], [event modifierFlags]);
}

// File drops
- (NSDragOperation)draggingEntered:(id<NSDraggingInfo>)sender
{
    return NSDragOperationGeneric;
}

- (BOOL)performDragOperation:(id<NSDraggingInfo>)sender
{
    NSPoint point = [self convertPoint:[sender draggingLocation] fromView:nil];
    CORE.Input.Mouse.currentPosition = (Vector2){ (float)point.x, (float)([self bounds].size.height - point.y) };

    NSArray *urls = [[sender draggingPasteboard] readObjectsForClasses:@[ [NSURL class] ]
        options:@{ NSPasteboardURLReadingFileURLsOnlyKey: @YES }];

    HandleDroppedFiles(urls);

    return YES;
}

// NSTextInputClient protocol, characters input
- (void)insertText:(id)string replacementRange:(NSRange)replacementRange
{
    // NOTE: Characters typed with Command are shortcuts, not text
    if ([[NSApp currentEvent] modifierFlags] & NSEventModifierFlagCommand) return;

    NSString *text = [string isKindOfClass:[NSAttributedString class]]? [string string] : string;
    HandleChars(text);

    [self unmarkText];
}

- (void)doCommandBySelector:(SEL)selector
{
}

- (void)setMarkedText:(id)string selectedRange:(NSRange)selectedRange replacementRange:(NSRange)replacementRange
{
    [markedText release];

    if ([string isKindOfClass:[NSAttributedString class]]) markedText = [[NSMutableAttributedString alloc] initWithAttributedString:string];
    else markedText = [[NSMutableAttributedString alloc] initWithString:string];

    markedSelection = selectedRange;
}

- (void)unmarkText
{
    [[markedText mutableString] setString:@""];
    markedSelection = NSMakeRange(0, 0);
}

- (BOOL)hasMarkedText
{
    return ([markedText length] > 0);
}

- (NSRange)markedRange
{
    if ([markedText length] > 0) return NSMakeRange(0, [markedText length]);

    return NSMakeRange(NSNotFound, 0);
}

// NOTE: The input method converts the marked text, it requires the selection and the text of the marked range
- (NSRange)selectedRange
{
    if ([markedText length] > 0) return markedSelection;

    return NSMakeRange(0, 0);
}

- (NSArray *)validAttributesForMarkedText
{
    return [NSArray array];
}

- (NSAttributedString *)attributedSubstringForProposedRange:(NSRange)range actualRange:(NSRangePointer)actualRange
{
    NSRange available = NSIntersectionRange(range, NSMakeRange(0, [markedText length]));

    if ((range.location == NSNotFound) || (available.length == 0)) return nil;
    if (actualRange != NULL) *actualRange = available;

    return [markedText attributedSubstringFromRange:available];
}

- (NSUInteger)characterIndexForPoint:(NSPoint)point
{
    return 0;
}

- (NSRect)firstRectForCharacterRange:(NSRange)range actualRange:(NSRangePointer)actualRange
{
    // NOTE: Input method candidates window is placed at the mouse position
    NSRect rect = NSMakeRect(CORE.Input.Mouse.currentPosition.x, [self bounds].size.height - CORE.Input.Mouse.currentPosition.y, 0, 0);

    return [[self window] convertRectToScreen:[self convertRect:rect toView:nil]];
}

@end


