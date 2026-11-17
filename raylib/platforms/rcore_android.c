/**********************************************************************************************
*
*   rcore_android - Functions to manage window, graphics device and inputs
*
*   PLATFORM: ANDROID
*       - Android (ARM, ARM64)
*
*   LIMITATIONS:
*       - Limitation 01
*       - Limitation 02
*
*   POSSIBLE IMPROVEMENTS:
*       - Improvement 01
*       - Improvement 02
*
*   CONFIGURATION:
*       #define RCORE_PLATFORM_CUSTOM_FLAG
*           Custom flag for rcore on target platform -not used-
*
*   DEPENDENCIES:
*       - Android NDK: Provides C API to access Android functionality
*       - gestures: Gestures system for touch-ready devices (or simulated from mouse inputs)
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

#include <android_native_app_glue.h>    // Required for: android_app struct and activity management
#include <android/window.h>             // Required for: AWINDOW_FLAG_FULLSCREEN definition and others
#include <android/log.h>                // Required for: Android log system: __android_log_vprint()
#include <android/asset_manager.h>      // Required for: AAssetManager
//#include <android/sensor.h>           // Required for: Android sensors functions (accelerometer, gyroscope, light...)

#include <errno.h>                      // Required for: error types
#include <stdatomic.h>                  // Required for: atomic_int, atomic_bool, atomic_uint
#include <sys/eventfd.h>                // Required for: eventfd()
#include <unistd.h>                     // Required for: read(), write()
#include <jni.h>                        // Required for: JNIEnv and JavaVM [Used in OpenURL() and GetCurrentMonitor()]

#include <EGL/egl.h>                    // Native platform windowing system interface

//----------------------------------------------------------------------------------
// Types and Structures Definition
//----------------------------------------------------------------------------------
typedef struct {
    // Application data
    struct android_app *app;            // Android activity
    struct android_poll_source *source; // Android events polling source
    bool appEnabled;                    // Flag to detect if app is active ** = true
    bool contextRebindRequired;         // Used to know context rebind required

    // Display data
    EGLDisplay device;                  // Native display device (physical screen connection)
    EGLSurface surface;                 // Surface to draw on, framebuffers (connected to context)
    EGLContext context;                 // Graphic context, mode in which drawing can be done
    EGLConfig config;                   // Graphic config
    int requestedWidth;                 // Screen width requested on window initialization
    int requestedHeight;                // Screen height requested on window initialization
    double displayCheckEnd;             // Time to stop checking the display size after a configuration change

    // Input data
    char *clipboardText;                // Clipboard text, returned by GetClipboardText()
    int deadChar;                       // Pending dead key accent, combined with the next character
    int gamepadDeviceId[MAX_GAMEPADS];  // Input device id of each gamepad, -1 if not connected
    double gamepadCheckTime;            // Next time to check for connected gamepads
    int mouseButtonState;               // Mouse buttons down (AMOTION_EVENT_BUTTON_*)
    bool keyPressedInFrame[MAX_KEYBOARD_KEYS];      // Keys pressed in current frame
    bool keyReleasePending[MAX_KEYBOARD_KEYS];      // Keys released in the frame they were pressed, release applied on next frame
    bool buttonPressedInFrame[MAX_MOUSE_BUTTONS];   // Mouse buttons and touch pressed in current frame
    bool buttonReleasePending[MAX_MOUSE_BUTTONS];   // Mouse buttons and touch released in the frame they were pressed
} PlatformData;

typedef struct {
    // Store data for both Hover and Touch events
    // Used to ignore Hover events which are interpreted as Touch events
    int32_t pointCount;                             // Number of touch points active
    int32_t pointId[MAX_TOUCH_POINTS];              // Point identifiers
    Vector2 position[MAX_TOUCH_POINTS];             // Touch position on screen

    int32_t hoverPoints[MAX_TOUCH_POINTS];          // Hover Points
} TouchRaw;

// Requests applied on the main (UI) thread, views can only be changed there
#define UI_REQUEST_CURSOR           1       // Update pointer icon
#define UI_REQUEST_CAPTURE          2       // Update pointer capture
#define UI_REQUEST_SYSTEM_BARS      4       // Update system bars visibility
#define UI_REQUEST_SOFT_KEYBOARD    8       // Update on-screen keyboard visibility

#define POINTER_ICON_NULL           0       // PointerIcon.TYPE_NULL
#define SYSTEM_BARS_TRANSIENT       2       // WindowInsetsController.BEHAVIOR_SHOW_TRANSIENT_BARS_BY_SWIPE
#define SYSTEM_UI_IMMERSIVE         0x1706  // View.SYSTEM_UI_FLAG_IMMERSIVE_STICKY, LAYOUT_STABLE, LAYOUT_HIDE_NAVIGATION, LAYOUT_FULLSCREEN, FULLSCREEN, HIDE_NAVIGATION

typedef struct {
    ALooper *looper;                    // Main (UI) thread looper
    int eventFd;                        // Event file descriptor waking the main thread looper
    atomic_uint requests;               // Pending requests (UI_REQUEST_*)
    atomic_int cursor;                  // Requested mouse cursor (MOUSE_CURSOR_*), -1 for hidden
    atomic_bool pointerCapture;         // Requested pointer capture state
    atomic_bool immersive;              // Requested system bars state, hidden if immersive
    atomic_bool softKeyboard;           // Requested on-screen keyboard visibility
} UiThreadData;

//----------------------------------------------------------------------------------
// Global Variables Definition
//----------------------------------------------------------------------------------
extern CoreData CORE;                   // Global CORE state context
static PlatformData platform = { 0 };   // Platform specific data

//----------------------------------------------------------------------------------
// Local Variables Definition
//----------------------------------------------------------------------------------
#define KEYCODE_MAP_SIZE 162
static const KeyboardKey mapKeycode[KEYCODE_MAP_SIZE] = {
    KEY_NULL,           // AKEYCODE_UNKNOWN
    0,                  // AKEYCODE_SOFT_LEFT
    0,                  // AKEYCODE_SOFT_RIGHT
    0,                  // AKEYCODE_HOME
    KEY_BACK,           // AKEYCODE_BACK
    0,                  // AKEYCODE_CALL
    0,                  // AKEYCODE_ENDCALL
    KEY_ZERO,           // AKEYCODE_0
    KEY_ONE,            // AKEYCODE_1
    KEY_TWO,            // AKEYCODE_2
    KEY_THREE,          // AKEYCODE_3
    KEY_FOUR,           // AKEYCODE_4
    KEY_FIVE,           // AKEYCODE_5
    KEY_SIX,            // AKEYCODE_6
    KEY_SEVEN,          // AKEYCODE_7
    KEY_EIGHT,          // AKEYCODE_8
    KEY_NINE,           // AKEYCODE_9
    0,                  // AKEYCODE_STAR
    0,                  // AKEYCODE_POUND
    KEY_UP,             // AKEYCODE_DPAD_UP
    KEY_DOWN,           // AKEYCODE_DPAD_DOWN
    KEY_LEFT,           // AKEYCODE_DPAD_LEFT
    KEY_RIGHT,          // AKEYCODE_DPAD_RIGHT
    0,                  // AKEYCODE_DPAD_CENTER
    KEY_VOLUME_UP,      // AKEYCODE_VOLUME_UP
    KEY_VOLUME_DOWN,    // AKEYCODE_VOLUME_DOWN
    0,                  // AKEYCODE_POWER
    0,                  // AKEYCODE_CAMERA
    0,                  // AKEYCODE_CLEAR
    KEY_A,              // AKEYCODE_A
    KEY_B,              // AKEYCODE_B
    KEY_C,              // AKEYCODE_C
    KEY_D,              // AKEYCODE_D
    KEY_E,              // AKEYCODE_E
    KEY_F,              // AKEYCODE_F
    KEY_G,              // AKEYCODE_G
    KEY_H,              // AKEYCODE_H
    KEY_I,              // AKEYCODE_I
    KEY_J,              // AKEYCODE_J
    KEY_K,              // AKEYCODE_K
    KEY_L,              // AKEYCODE_L
    KEY_M,              // AKEYCODE_M
    KEY_N,              // AKEYCODE_N
    KEY_O,              // AKEYCODE_O
    KEY_P,              // AKEYCODE_P
    KEY_Q,              // AKEYCODE_Q
    KEY_R,              // AKEYCODE_R
    KEY_S,              // AKEYCODE_S
    KEY_T,              // AKEYCODE_T
    KEY_U,              // AKEYCODE_U
    KEY_V,              // AKEYCODE_V
    KEY_W,              // AKEYCODE_W
    KEY_X,              // AKEYCODE_X
    KEY_Y,              // AKEYCODE_Y
    KEY_Z,              // AKEYCODE_Z
    KEY_COMMA,          // AKEYCODE_COMMA
    KEY_PERIOD,         // AKEYCODE_PERIOD
    KEY_LEFT_ALT,       // AKEYCODE_ALT_LEFT
    KEY_RIGHT_ALT,      // AKEYCODE_ALT_RIGHT
    KEY_LEFT_SHIFT,     // AKEYCODE_SHIFT_LEFT
    KEY_RIGHT_SHIFT,    // AKEYCODE_SHIFT_RIGHT
    KEY_TAB,            // AKEYCODE_TAB
    KEY_SPACE,          // AKEYCODE_SPACE
    0,                  // AKEYCODE_SYM
    0,                  // AKEYCODE_EXPLORER
    0,                  // AKEYCODE_ENVELOPE
    KEY_ENTER,          // AKEYCODE_ENTER
    KEY_BACKSPACE,      // AKEYCODE_DEL
    KEY_GRAVE,          // AKEYCODE_GRAVE
    KEY_MINUS,          // AKEYCODE_MINUS
    KEY_EQUAL,          // AKEYCODE_EQUALS
    KEY_LEFT_BRACKET,   // AKEYCODE_LEFT_BRACKET
    KEY_RIGHT_BRACKET,  // AKEYCODE_RIGHT_BRACKET
    KEY_BACKSLASH,      // AKEYCODE_BACKSLASH
    KEY_SEMICOLON,      // AKEYCODE_SEMICOLON
    KEY_APOSTROPHE,     // AKEYCODE_APOSTROPHE
    KEY_SLASH,          // AKEYCODE_SLASH
    0,                  // AKEYCODE_AT
    0,                  // AKEYCODE_NUM
    0,                  // AKEYCODE_HEADSETHOOK
    0,                  // AKEYCODE_FOCUS
    0,                  // AKEYCODE_PLUS
    KEY_MENU,           // AKEYCODE_MENU
    0,                  // AKEYCODE_NOTIFICATION
    0,                  // AKEYCODE_SEARCH
    0,                  // AKEYCODE_MEDIA_PLAY_PAUSE
    0,                  // AKEYCODE_MEDIA_STOP
    0,                  // AKEYCODE_MEDIA_NEXT
    0,                  // AKEYCODE_MEDIA_PREVIOUS
    0,                  // AKEYCODE_MEDIA_REWIND
    0,                  // AKEYCODE_MEDIA_FAST_FORWARD
    0,                  // AKEYCODE_MUTE
    KEY_PAGE_UP,        // AKEYCODE_PAGE_UP
    KEY_PAGE_DOWN,      // AKEYCODE_PAGE_DOWN
    0,                  // AKEYCODE_PICTSYMBOLS
    0,                  // AKEYCODE_SWITCH_CHARSET
    0,                  // AKEYCODE_BUTTON_A
    0,                  // AKEYCODE_BUTTON_B
    0,                  // AKEYCODE_BUTTON_C
    0,                  // AKEYCODE_BUTTON_X
    0,                  // AKEYCODE_BUTTON_Y
    0,                  // AKEYCODE_BUTTON_Z
    0,                  // AKEYCODE_BUTTON_L1
    0,                  // AKEYCODE_BUTTON_R1
    0,                  // AKEYCODE_BUTTON_L2
    0,                  // AKEYCODE_BUTTON_R2
    0,                  // AKEYCODE_BUTTON_THUMBL
    0,                  // AKEYCODE_BUTTON_THUMBR
    0,                  // AKEYCODE_BUTTON_START
    0,                  // AKEYCODE_BUTTON_SELECT
    0,                  // AKEYCODE_BUTTON_MODE
    KEY_ESCAPE,         // AKEYCODE_ESCAPE
    KEY_DELETE,         // AKEYCODE_FORWARD_DELL
    KEY_LEFT_CONTROL,   // AKEYCODE_CTRL_LEFT
    KEY_RIGHT_CONTROL,  // AKEYCODE_CTRL_RIGHT
    KEY_CAPS_LOCK,      // AKEYCODE_CAPS_LOCK
    KEY_SCROLL_LOCK,    // AKEYCODE_SCROLL_LOCK
    KEY_LEFT_SUPER,     // AKEYCODE_META_LEFT
    KEY_RIGHT_SUPER,    // AKEYCODE_META_RIGHT
    0,                  // AKEYCODE_FUNCTION
    KEY_PRINT_SCREEN,   // AKEYCODE_SYSRQ
    KEY_PAUSE,          // AKEYCODE_BREAK
    KEY_HOME,           // AKEYCODE_MOVE_HOME
    KEY_END,            // AKEYCODE_MOVE_END
    KEY_INSERT,         // AKEYCODE_INSERT
    0,                  // AKEYCODE_FORWARD
    0,                  // AKEYCODE_MEDIA_PLAY
    0,                  // AKEYCODE_MEDIA_PAUSE
    0,                  // AKEYCODE_MEDIA_CLOSE
    0,                  // AKEYCODE_MEDIA_EJECT
    0,                  // AKEYCODE_MEDIA_RECORD
    KEY_F1,             // AKEYCODE_F1
    KEY_F2,             // AKEYCODE_F2
    KEY_F3,             // AKEYCODE_F3
    KEY_F4,             // AKEYCODE_F4
    KEY_F5,             // AKEYCODE_F5
    KEY_F6,             // AKEYCODE_F6
    KEY_F7,             // AKEYCODE_F7
    KEY_F8,             // AKEYCODE_F8
    KEY_F9,             // AKEYCODE_F9
    KEY_F10,            // AKEYCODE_F10
    KEY_F11,            // AKEYCODE_F11
    KEY_F12,            // AKEYCODE_F12
    KEY_NUM_LOCK,       // AKEYCODE_NUM_LOCK
    KEY_KP_0,           // AKEYCODE_NUMPAD_0
    KEY_KP_1,           // AKEYCODE_NUMPAD_1
    KEY_KP_2,           // AKEYCODE_NUMPAD_2
    KEY_KP_3,           // AKEYCODE_NUMPAD_3
    KEY_KP_4,           // AKEYCODE_NUMPAD_4
    KEY_KP_5,           // AKEYCODE_NUMPAD_5
    KEY_KP_6,           // AKEYCODE_NUMPAD_6
    KEY_KP_7,           // AKEYCODE_NUMPAD_7
    KEY_KP_8,           // AKEYCODE_NUMPAD_8
    KEY_KP_9,           // AKEYCODE_NUMPAD_9
    KEY_KP_DIVIDE,      // AKEYCODE_NUMPAD_DIVIDE
    KEY_KP_MULTIPLY,    // AKEYCODE_NUMPAD_MULTIPLY
    KEY_KP_SUBTRACT,    // AKEYCODE_NUMPAD_SUBTRACT
    KEY_KP_ADD,         // AKEYCODE_NUMPAD_ADD
    KEY_KP_DECIMAL,     // AKEYCODE_NUMPAD_DOT
    0,                  // AKEYCODE_NUMPAD_COMMA
    KEY_KP_ENTER,       // AKEYCODE_NUMPAD_ENTER
    KEY_KP_EQUAL        // AKEYCODE_NUMPAD_EQUALS
};

// Android pointer icon types (PointerIcon.TYPE_*) for raylib mouse cursors
static const int mapCursor[MOUSE_CURSOR_NOT_ALLOWED + 1] = {
    1000,               // MOUSE_CURSOR_DEFAULT: TYPE_ARROW
    1000,               // MOUSE_CURSOR_ARROW: TYPE_ARROW
    1008,               // MOUSE_CURSOR_IBEAM: TYPE_TEXT
    1007,               // MOUSE_CURSOR_CROSSHAIR: TYPE_CROSSHAIR
    1002,               // MOUSE_CURSOR_POINTING_HAND: TYPE_HAND
    1014,               // MOUSE_CURSOR_RESIZE_EW: TYPE_HORIZONTAL_DOUBLE_ARROW
    1015,               // MOUSE_CURSOR_RESIZE_NS: TYPE_VERTICAL_DOUBLE_ARROW
    1017,               // MOUSE_CURSOR_RESIZE_NWSE: TYPE_TOP_LEFT_DIAGONAL_DOUBLE_ARROW
    1016,               // MOUSE_CURSOR_RESIZE_NESW: TYPE_TOP_RIGHT_DIAGONAL_DOUBLE_ARROW
    1013,               // MOUSE_CURSOR_RESIZE_ALL: TYPE_ALL_SCROLL
    1012                // MOUSE_CURSOR_NOT_ALLOWED: TYPE_NO_DROP
};

static TouchRaw touchRaw = { 0 };
static UiThreadData ui = { .eventFd = -1 };     // NOTE: Kept for the process, the main thread looper does not change

//----------------------------------------------------------------------------------
// Module Internal Functions Declaration
//----------------------------------------------------------------------------------
int InitPlatform(void);          // Initialize platform (graphics, inputs and more)
void ClosePlatform(void);        // Close platform

static void AndroidCommandCallback(struct android_app *app, int32_t cmd);           // Process Android activity lifecycle commands
static int32_t AndroidInputCallback(struct android_app *app, AInputEvent *event);   // Process Android inputs
static GamepadButton AndroidTranslateGamepadButton(int button);                     // Map Android gamepad button to raylib gamepad button

static void SetupFramebuffer(int width, int height); // Setup main framebuffer (required by InitPlatform())
static void UpdateDisplaySize(void);                // Update framebuffer and viewport if the display size changed (rotation)

static JNIEnv *AttachJni(void);                     // Attach current thread to the Java VM, NULL on failure
static void DetachJni(void);                        // Detach current thread from the Java VM
static bool JniFailed(JNIEnv *env);                 // Check and clear a pending Java exception
static jobject GetDefaultDisplay(JNIEnv *env);      // Get default display object (android.view.Display)
static int GetKeyCharacter(int deviceId, int keycode, int metaState); // Get character produced by a key, considering dead keys
static int GetGamepadIndex(int deviceId, bool connect); // Get gamepad index for an input device, connecting it if required
static void ConnectGamepad(int gamepad, int deviceId); // Register a connected gamepad
static void DisconnectGamepad(int gamepad);         // Unregister a disconnected gamepad
static void CheckGamepads(void);                    // Check connected and disconnected gamepads
static void HandleKey(KeyboardKey key, bool pressed); // Register key state change
static void HandleMouseButton(int button, bool pressed); // Register mouse button or touch state change
static Vector2 GetScreenPosition(float x, float y); // Convert display position to screen position
static void RequestUiUpdate(unsigned int request);  // Request an update applied on the main (UI) thread
static void UpdateSystemBars(void);                 // Request system bars visibility from window flags
static int UiThreadCallback(int fd, int events, void *data); // Apply pending requests on the main (UI) thread

static int android_read(void *cookie, char *buf, int size);
static int android_write(void *cookie, const char *buf, int size);
static fpos_t android_seek(void *cookie, fpos_t offset, int whence);
static int android_close(void *cookie);

// WARNING: fopen() calls are intercepted via linker flag -Wl,--wrap=fopen: the linker renames
// the original fopen -> __real_fopen and redirects all call sites to __wrap_fopen
// The flag MUST be applied at every final link step that needs wrapping;
// it has no effect when only building a static archive (.a)
//
// STATIC library (.a) — wrapping deferred to consumer's final link step:
//   both raylib and consumer fopen calls are wrapped together in one link
//       CMake: handled automatically — the PUBLIC flag propagates as INTERFACE_LINK_OPTIONS
//              to the consumer's final link via target_link_libraries
//        Make: pass -Wl,--wrap=fopen to the linker command producing the final artifact
//      custom: pass -Wl,--wrap=fopen to the linker command producing the final artifact
//
// SHARED library (.so) — wrapping is self-contained:
//   only fopen calls linked into the .so are wrapped; the consumer's own fopen calls
//   are NOT wrapped unless the consumer also links with -Wl,--wrap=fopen independently
//       CMake: handled automatically — CMakeLists.txt sets target_link_options(raylib PUBLIC
//              -Wl,--wrap=fopen) which applies the flag to the .so link;
//              only raylib internals are wrapped, app code requires a separate flag
//        Make: handled automatically — src/Makefile sets LDFLAGS += -Wl,--wrap=fopen;
//              only raylib internals are wrapped, app code requires a separate flag
//      custom: apply -Wl,--wrap=fopen to the linker command producing the .so
FILE *__real_fopen(const char *fileName, const char *mode); // Real fopen, provided by the linker (--wrap=fopen)
FILE *__wrap_fopen(const char *fileName, const char *mode); // Replacement for fopen()

FILE *funopen(const void *cookie, int (*readfn)(void *, char *, int), int (*writefn)(void *, const char *, int),
              fpos_t (*seekfn)(void *, fpos_t, int), int (*closefn)(void *));

//----------------------------------------------------------------------------------
// Module Functions Declaration
//----------------------------------------------------------------------------------
// NOTE: Functions declaration is provided by raylib.h

//----------------------------------------------------------------------------------
// Module Functions Definition: Application
//----------------------------------------------------------------------------------

// To allow easier porting to android, allow the user to define a
// custom main function which is called from android_main
extern void android_run(void);

// Android main function
void android_main(struct android_app *app)
{
    platform.app = app;

    android_run();

    // Request to end the native activity
    ANativeActivity_finish(app->activity);

    // Android ALooper_pollOnce() variables
    int pollResult = 0;
    int pollEvents = 0;

    // Waiting for application events before complete finishing
    while (!app->destroyRequested)
    {
        // Poll all events until return value TIMEOUT is reached, meaning no events left to process
        while ((pollResult = ALooper_pollOnce(0, NULL, &pollEvents, (void **)&platform.source)) > ALOOPER_POLL_TIMEOUT)
        {
            if (platform.source != NULL) platform.source->process(app, platform.source);
        }
    }

    platform.app = NULL;
}

// NOTE: Add this to header (if apps really need it)
struct android_app *GetAndroidApp(void)
{
    return platform.app;
}

// Activity entry point of android_native_app_glue, renamed when compiling it (-DANativeActivity_onCreate=AndroidGlueOnCreate)
void AndroidGlueOnCreate(ANativeActivity *activity, void *savedState, size_t savedStateSize);

// Activity entry point, called on the main (UI) thread
JNIEXPORT void ANativeActivity_onCreate(ANativeActivity *activity, void *savedState, size_t savedStateSize)
{
    if (ui.looper == NULL)
    {
        ui.looper = ALooper_forThread();
        ui.eventFd = eventfd(0, EFD_NONBLOCK | EFD_CLOEXEC);

        if ((ui.looper != NULL) && (ui.eventFd >= 0))
        {
            ALooper_acquire(ui.looper);
            ALooper_addFd(ui.looper, ui.eventFd, ALOOPER_POLL_CALLBACK, ALOOPER_EVENT_INPUT, UiThreadCallback, NULL);
        }
    }

    AndroidGlueOnCreate(activity, savedState, savedStateSize);
}

// Show on-screen keyboard
// NOTE: Characters typed are received as key events
void ShowSoftKeyboard(void)
{
    atomic_store(&ui.softKeyboard, true);
    RequestUiUpdate(UI_REQUEST_SOFT_KEYBOARD);
}

// Hide on-screen keyboard
void HideSoftKeyboard(void)
{
    atomic_store(&ui.softKeyboard, false);
    RequestUiUpdate(UI_REQUEST_SOFT_KEYBOARD);
}

//----------------------------------------------------------------------------------
// Module Functions Definition: Window and Graphics Device
//----------------------------------------------------------------------------------

// Check if application should close
bool WindowShouldClose(void)
{
    if (CORE.Window.ready) return CORE.Window.shouldClose;
    else return true;
}

// Toggle fullscreen mode
// NOTE: Fullscreen hides the system bars (immersive mode)
void ToggleFullscreen(void)
{
    if (FLAG_IS_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE)) FLAG_CLEAR(CORE.Window.flags, FLAG_FULLSCREEN_MODE);
    else
    {
        FLAG_CLEAR(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE);
        FLAG_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE);
    }

    UpdateSystemBars();
}

// Toggle borderless windowed mode
// NOTE: Same as fullscreen, the display mode is never changed
void ToggleBorderlessWindowed(void)
{
    if (FLAG_IS_SET(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE)) FLAG_CLEAR(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE);
    else
    {
        FLAG_CLEAR(CORE.Window.flags, FLAG_FULLSCREEN_MODE);
        FLAG_SET(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE);
    }

    UpdateSystemBars();
}

// Set window state: maximized, if resizable
void MaximizeWindow(void)
{
    TRACELOG(LOG_WARNING, "MaximizeWindow() not available on target platform");
}

// Set window state: minimized
void MinimizeWindow(void)
{
    TRACELOG(LOG_WARNING, "MinimizeWindow() not available on target platform");
}

// Restore window from being minimized/maximized
void RestoreWindow(void)
{
    TRACELOG(LOG_WARNING, "RestoreWindow() not available on target platform");
}

// Set window configuration state using flags
void SetWindowState(unsigned int flags)
{
    if (!CORE.Window.ready) TRACELOG(LOG_WARNING, "WINDOW: SetWindowState does nothing before window initialization, Use \"SetConfigFlags\" instead");

    // State change: FLAG_WINDOW_ALWAYS_RUN
    if (FLAG_IS_SET(flags, FLAG_WINDOW_ALWAYS_RUN)) FLAG_SET(CORE.Window.flags, FLAG_WINDOW_ALWAYS_RUN);

    // State change: FLAG_FULLSCREEN_MODE
    if (FLAG_IS_SET(flags, FLAG_FULLSCREEN_MODE) && !FLAG_IS_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE)) ToggleFullscreen();

    // State change: FLAG_BORDERLESS_WINDOWED_MODE
    if (FLAG_IS_SET(flags, FLAG_BORDERLESS_WINDOWED_MODE) && !FLAG_IS_SET(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE)) ToggleBorderlessWindowed();
}

// Clear window configuration state flags
void ClearWindowState(unsigned int flags)
{
    // State change: FLAG_WINDOW_ALWAYS_RUN
    if (FLAG_IS_SET(flags, FLAG_WINDOW_ALWAYS_RUN)) FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_ALWAYS_RUN);

    // State change: FLAG_FULLSCREEN_MODE
    if (FLAG_IS_SET(flags, FLAG_FULLSCREEN_MODE) && FLAG_IS_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE)) ToggleFullscreen();

    // State change: FLAG_BORDERLESS_WINDOWED_MODE
    if (FLAG_IS_SET(flags, FLAG_BORDERLESS_WINDOWED_MODE) && FLAG_IS_SET(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE)) ToggleBorderlessWindowed();
}

// Set icon for window
void SetWindowIcon(Image image)
{
    TRACELOG(LOG_WARNING, "SetWindowIcon() not available on target platform");
}

// Set icon for window
void SetWindowIcons(Image *images, int count)
{
    TRACELOG(LOG_WARNING, "SetWindowIcons() not available on target platform");
}

// Set title for window
void SetWindowTitle(const char *title)
{
    CORE.Window.title = title;
}

// Set window position on screen (windowed mode)
void SetWindowPosition(int x, int y)
{
    TRACELOG(LOG_WARNING, "SetWindowPosition() not available on target platform");
}

// Set monitor for the current window
void SetWindowMonitor(int monitor)
{
    TRACELOG(LOG_WARNING, "SetWindowMonitor() not available on target platform");
}

// Set window minimum dimensions (FLAG_WINDOW_RESIZABLE)
void SetWindowMinSize(int width, int height)
{
    CORE.Window.screenMin.width = width;
    CORE.Window.screenMin.height = height;
}

// Set window maximum dimensions (FLAG_WINDOW_RESIZABLE)
void SetWindowMaxSize(int width, int height)
{
    CORE.Window.screenMax.width = width;
    CORE.Window.screenMax.height = height;
}

// Set window dimensions
void SetWindowSize(int width, int height)
{
    TRACELOG(LOG_WARNING, "SetWindowSize() not available on target platform");
}

// Set window opacity, value opacity is between 0.0 and 1.0
void SetWindowOpacity(float opacity)
{
    TRACELOG(LOG_WARNING, "SetWindowOpacity() not available on target platform");
}

// Set window focused
void SetWindowFocused(void)
{
    TRACELOG(LOG_WARNING, "SetWindowFocused() not available on target platform");
}

// Get native window handle
void *GetWindowHandle(void)
{
    return (void *)platform.app->window; // Type: ANativeWindow*
}

// Get number of monitors
// NOTE: Only the display the activity is on is available
int GetMonitorCount(void)
{
    return 1;
}

// Get current monitor where window is placed
int GetCurrentMonitor(void)
{
    return 0;
}

// Get selected monitor position
Vector2 GetMonitorPosition(int monitor)
{
    return (Vector2){ 0, 0 };
}

// Get selected monitor width (currently used by monitor)
int GetMonitorWidth(int monitor)
{
    int width = 0;

    JNIEnv *env = AttachJni();
    if (env == NULL) return width;

    jobject display = GetDefaultDisplay(env);
    if (display != NULL)
    {
        // Display size in pixels, in the current orientation
        jclass pointClass = (*env)->FindClass(env, "android/graphics/Point");
        jobject size = (*env)->NewObject(env, pointClass, (*env)->GetMethodID(env, pointClass, "<init>", "()V"));
        jclass displayClass = (*env)->GetObjectClass(env, display);
        (*env)->CallVoidMethod(env, display, (*env)->GetMethodID(env, displayClass, "getRealSize", "(Landroid/graphics/Point;)V"), size);
        if (!JniFailed(env)) width = (*env)->GetIntField(env, size, (*env)->GetFieldID(env, pointClass, "x", "I"));
    }

    DetachJni();

    return width;
}

// Get selected monitor height (currently used by monitor)
int GetMonitorHeight(int monitor)
{
    int height = 0;

    JNIEnv *env = AttachJni();
    if (env == NULL) return height;

    jobject display = GetDefaultDisplay(env);
    if (display != NULL)
    {
        // Display size in pixels, in the current orientation
        jclass pointClass = (*env)->FindClass(env, "android/graphics/Point");
        jobject size = (*env)->NewObject(env, pointClass, (*env)->GetMethodID(env, pointClass, "<init>", "()V"));
        jclass displayClass = (*env)->GetObjectClass(env, display);
        (*env)->CallVoidMethod(env, display, (*env)->GetMethodID(env, displayClass, "getRealSize", "(Landroid/graphics/Point;)V"), size);
        if (!JniFailed(env)) height = (*env)->GetIntField(env, size, (*env)->GetFieldID(env, pointClass, "y", "I"));
    }

    DetachJni();

    return height;
}

// Get selected monitor physical width in millimetres
int GetMonitorPhysicalWidth(int monitor)
{
    int width = 0;

    JNIEnv *env = AttachJni();
    if (env == NULL) return width;

    jobject display = GetDefaultDisplay(env);
    if (display != NULL)
    {
        // NOTE: Real pixel density of the display, AConfiguration density is a density bucket
        jclass metricsClass = (*env)->FindClass(env, "android/util/DisplayMetrics");
        jobject metrics = (*env)->NewObject(env, metricsClass, (*env)->GetMethodID(env, metricsClass, "<init>", "()V"));
        jclass displayClass = (*env)->GetObjectClass(env, display);
        (*env)->CallVoidMethod(env, display, (*env)->GetMethodID(env, displayClass, "getRealMetrics", "(Landroid/util/DisplayMetrics;)V"), metrics);
        if (!JniFailed(env))
        {
            int pixels = (*env)->GetIntField(env, metrics, (*env)->GetFieldID(env, metricsClass, "widthPixels", "I"));
            float dpi = (*env)->GetFloatField(env, metrics, (*env)->GetFieldID(env, metricsClass, "xdpi", "F"));
            if (dpi > 0.0f) width = (int)roundf((float)pixels/dpi*25.4f);
        }
    }

    DetachJni();

    return width;
}

// Get selected monitor physical height in millimetres
int GetMonitorPhysicalHeight(int monitor)
{
    int height = 0;

    JNIEnv *env = AttachJni();
    if (env == NULL) return height;

    jobject display = GetDefaultDisplay(env);
    if (display != NULL)
    {
        // NOTE: Real pixel density of the display, AConfiguration density is a density bucket
        jclass metricsClass = (*env)->FindClass(env, "android/util/DisplayMetrics");
        jobject metrics = (*env)->NewObject(env, metricsClass, (*env)->GetMethodID(env, metricsClass, "<init>", "()V"));
        jclass displayClass = (*env)->GetObjectClass(env, display);
        (*env)->CallVoidMethod(env, display, (*env)->GetMethodID(env, displayClass, "getRealMetrics", "(Landroid/util/DisplayMetrics;)V"), metrics);
        if (!JniFailed(env))
        {
            int pixels = (*env)->GetIntField(env, metrics, (*env)->GetFieldID(env, metricsClass, "heightPixels", "I"));
            float dpi = (*env)->GetFloatField(env, metrics, (*env)->GetFieldID(env, metricsClass, "ydpi", "F"));
            if (dpi > 0.0f) height = (int)roundf((float)pixels/dpi*25.4f);
        }
    }

    DetachJni();

    return height;
}

// Get selected monitor refresh rate
int GetMonitorRefreshRate(int monitor)
{
    int refreshRate = 0;

    JNIEnv *env = AttachJni();
    if (env == NULL) return refreshRate;

    jobject display = GetDefaultDisplay(env);
    if (display != NULL)
    {
        jclass displayClass = (*env)->GetObjectClass(env, display);
        float rate = (*env)->CallFloatMethod(env, display, (*env)->GetMethodID(env, displayClass, "getRefreshRate", "()F"));
        if (!JniFailed(env)) refreshRate = (int)roundf(rate);
    }

    DetachJni();

    return refreshRate;
}

// Get the human-readable, UTF-8 encoded name of the selected monitor
const char *GetMonitorName(int monitor)
{
    static char name[128] = { 0 };
    name[0] = '\0';

    JNIEnv *env = AttachJni();
    if (env == NULL) return name;

    jobject display = GetDefaultDisplay(env);
    if (display != NULL)
    {
        jclass displayClass = (*env)->GetObjectClass(env, display);
        jstring displayName = (jstring)(*env)->CallObjectMethod(env, display, (*env)->GetMethodID(env, displayClass, "getName", "()Ljava/lang/String;"));
        if (!JniFailed(env) && (displayName != NULL))
        {
            const char *text = (*env)->GetStringUTFChars(env, displayName, NULL);
            snprintf(name, sizeof(name), "%s", text);
            (*env)->ReleaseStringUTFChars(env, displayName, text);
        }
    }

    DetachJni();

    return name;
}

// Get window position XY on monitor
// NOTE: The window covers the display
Vector2 GetWindowPosition(void)
{
    return (Vector2){ 0, 0 };
}

// Get window scale DPI factor for current monitor
Vector2 GetWindowScaleDPI(void)
{
    int density = AConfiguration_getDensity(platform.app->config);
    float scale = (float)density/160;
    return (Vector2){ scale, scale };
}

// Set clipboard text content
// NOTE: Implementation based on https://stackoverflow.com/questions/76759222/how-to-copy-to-clipboard-on-android-via-jni
void SetClipboardText(const char *text)
{
    JNIEnv *env = NULL;
    JavaVM *vm = platform.app->activity->vm;
    jint result = (*vm)->AttachCurrentThread(vm, &env, NULL);
    if (result != JNI_OK) {
        TRACELOG(LOG_WARNING, "Failed to allocate a thread to set the clipboard text");
        return;
    }
    jobject activity = platform.app->activity->clazz;
    if (!activity) {
        TRACELOG(LOG_WARNING, "Failed to get Android activity class reference");
        (*vm)->DetachCurrentThread(vm);
        return;
    }

    // ClipboardManager
    jclass activityClass = (*env)->GetObjectClass(env, activity);
    jmethodID getSystemService = (*env)->GetMethodID(env, activityClass, "getSystemService", "(Ljava/lang/String;)Ljava/lang/Object;");
    jstring serviceName = (*env)->NewStringUTF(env, "clipboard");
    jobject clipboard = (*env)->CallObjectMethod(env, activity, getSystemService, serviceName);

    // ClipData
    jclass clipDataClass = (*env)->FindClass(env, "android/content/ClipData");
    jmethodID newPlainText = (*env)->GetStaticMethodID(env, clipDataClass, "newPlainText", "(Ljava/lang/CharSequence;Ljava/lang/CharSequence;)Landroid/content/ClipData;");
    const char *windowTitle = (CORE.Window.title != NULL) ? CORE.Window.title : "raylib";
    jstring label = (*env)->NewStringUTF(env, windowTitle);
    jstring textString = (*env)->NewStringUTF(env, text);
    jobject clip = (*env)->CallStaticObjectMethod(env, clipDataClass, newPlainText, label, textString);

    // cm.setPrimaryClip(clip)
    jclass clipboardClass = (*env)->FindClass(env, "android/content/ClipboardManager");
    jmethodID setPrimaryClip = (*env)->GetMethodID(env, clipboardClass, "setPrimaryClip", "(Landroid/content/ClipData;)V");
    (*env)->CallVoidMethod(env, clipboard, setPrimaryClip, clip);

    (*vm)->DetachCurrentThread(vm);
}

// Get clipboard text content
// NOTE: Returned string is kept until the next call
const char *GetClipboardText(void)
{
    RL_FREE(platform.clipboardText);
    platform.clipboardText = NULL;

    JNIEnv *env = AttachJni();
    if (env == NULL) return NULL;

    jobject activity = platform.app->activity->clazz;
    jclass activityClass = (*env)->GetObjectClass(env, activity);
    jmethodID getSystemService = (*env)->GetMethodID(env, activityClass, "getSystemService", "(Ljava/lang/String;)Ljava/lang/Object;");
    jstring serviceName = (*env)->NewStringUTF(env, "clipboard");
    jobject clipboard = (*env)->CallObjectMethod(env, activity, getSystemService, serviceName);

    if (!JniFailed(env) && (clipboard != NULL))
    {
        jclass clipboardClass = (*env)->FindClass(env, "android/content/ClipboardManager");
        jobject clip = (*env)->CallObjectMethod(env, clipboard, (*env)->GetMethodID(env, clipboardClass, "getPrimaryClip", "()Landroid/content/ClipData;"));

        if (!JniFailed(env) && (clip != NULL))
        {
            jclass clipDataClass = (*env)->FindClass(env, "android/content/ClipData");
            int count = (*env)->CallIntMethod(env, clip, (*env)->GetMethodID(env, clipDataClass, "getItemCount", "()I"));

            if (!JniFailed(env) && (count > 0))
            {
                // NOTE: Item converted to text, also for URIs and intents
                jobject item = (*env)->CallObjectMethod(env, clip, (*env)->GetMethodID(env, clipDataClass, "getItemAt", "(I)Landroid/content/ClipData$Item;"), 0);
                jclass itemClass = (*env)->FindClass(env, "android/content/ClipData$Item");
                jobject sequence = (*env)->CallObjectMethod(env, item, (*env)->GetMethodID(env, itemClass, "coerceToText", "(Landroid/content/Context;)Ljava/lang/CharSequence;"), activity);

                if (!JniFailed(env) && (sequence != NULL))
                {
                    jclass sequenceClass = (*env)->FindClass(env, "java/lang/CharSequence");
                    jstring text = (jstring)(*env)->CallObjectMethod(env, sequence, (*env)->GetMethodID(env, sequenceClass, "toString", "()Ljava/lang/String;"));

                    if (!JniFailed(env) && (text != NULL))
                    {
                        const char *utf8 = (*env)->GetStringUTFChars(env, text, NULL);
                        int length = (int)strlen(utf8);
                        platform.clipboardText = (char *)RL_CALLOC(length + 1, sizeof(char));
                        memcpy(platform.clipboardText, utf8, length);
                        (*env)->ReleaseStringUTFChars(env, text, utf8);
                    }
                }
            }
        }
    }

    DetachJni();

    return platform.clipboardText;
}

// Get clipboard image
Image GetClipboardImage(void)
{
    Image image = { 0 };

    TRACELOG(LOG_WARNING, "GetClipboardImage() not implemented on target platform");

    return image;
}

// Show mouse cursor
void ShowCursor(void)
{
    CORE.Input.Mouse.cursorHidden = false;
    atomic_store(&ui.cursor, CORE.Input.Mouse.cursor);
    RequestUiUpdate(UI_REQUEST_CURSOR);
}

// Hide mouse cursor
void HideCursor(void)
{
    CORE.Input.Mouse.cursorHidden = true;
    atomic_store(&ui.cursor, -1);
    RequestUiUpdate(UI_REQUEST_CURSOR);
}

// Enable cursor (unlock cursor)
void EnableCursor(void)
{
    // Set cursor position in the middle
    SetMousePosition(CORE.Window.screen.width/2, CORE.Window.screen.height/2);

    CORE.Input.Mouse.cursorLocked = false;
    atomic_store(&ui.pointerCapture, false);
    RequestUiUpdate(UI_REQUEST_CAPTURE);
}

// Disables cursor (lock cursor)
// NOTE: Pointer capture, mouse movement is received as relative motion
void DisableCursor(void)
{
    // Set cursor position in the middle
    SetMousePosition(CORE.Window.screen.width/2, CORE.Window.screen.height/2);

    CORE.Input.Mouse.cursorLocked = true;
    atomic_store(&ui.pointerCapture, true);
    RequestUiUpdate(UI_REQUEST_CAPTURE);
}

// Swap back buffer with front buffer (screen drawing)
void SwapScreenBuffer(void)
{
    if (platform.surface != EGL_NO_SURFACE) eglSwapBuffers(platform.device, platform.surface);
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
        JNIEnv *env = NULL;
        JavaVM *vm = platform.app->activity->vm;
        (*vm)->AttachCurrentThread(vm, &env, NULL);

        jstring urlString = (*env)->NewStringUTF(env, url);
        jclass uriClass = (*env)->FindClass(env, "android/net/Uri");
        jmethodID uriParse = (*env)->GetStaticMethodID(env, uriClass, "parse", "(Ljava/lang/String;)Landroid/net/Uri;");
        jobject uri = (*env)->CallStaticObjectMethod(env, uriClass, uriParse, urlString);

        jclass intentClass = (*env)->FindClass(env, "android/content/Intent");
        jfieldID actionViewId = (*env)->GetStaticFieldID(env, intentClass, "ACTION_VIEW", "Ljava/lang/String;");
        jobject actionView = (*env)->GetStaticObjectField(env, intentClass, actionViewId);
        jmethodID newIntent = (*env)->GetMethodID(env, intentClass, "<init>", "(Ljava/lang/String;Landroid/net/Uri;)V");
        jobject intent = (*env)->AllocObject(env, intentClass);

        (*env)->CallVoidMethod(env, intent, newIntent, actionView, uri);
        jclass activityClass = (*env)->FindClass(env, "android/app/Activity");
        jmethodID startActivity = (*env)->GetMethodID(env, activityClass, "startActivity", "(Landroid/content/Intent;)V");
        (*env)->CallVoidMethod(env, platform.app->activity->clazz, startActivity, intent);

        (*vm)->DetachCurrentThread(vm);
    }
}

//----------------------------------------------------------------------------------
// Module Functions Definition: Inputs
//----------------------------------------------------------------------------------

// Set internal gamepad mappings
int SetGamepadMappings(const char *mappings)
{
    TRACELOG(LOG_WARNING, "SetGamepadMappings() not implemented on target platform");
    return 0;
}

// Set gamepad vibration
// NOTE: Gamepads with two motors are driven separately on Android 12 and later, one motor uses the strongest value
void SetGamepadVibration(int gamepad, float leftMotor, float rightMotor, float duration)
{
    if ((gamepad < 0) || (gamepad >= MAX_GAMEPADS) || !CORE.Input.Gamepad.ready[gamepad] || (duration <= 0.0f)) return;

    if (leftMotor < 0.0f) leftMotor = 0.0f;
    if (leftMotor > 1.0f) leftMotor = 1.0f;
    if (rightMotor < 0.0f) rightMotor = 0.0f;
    if (rightMotor > 1.0f) rightMotor = 1.0f;
    if (duration > MAX_GAMEPAD_VIBRATION_TIME) duration = MAX_GAMEPAD_VIBRATION_TIME;

    JNIEnv *env = AttachJni();
    if (env == NULL) return;

    int sdkVersion = AConfiguration_getSdkVersion(platform.app->config);
    jlong milliseconds = (jlong)(duration*1000.0f);
    int leftAmplitude = (int)roundf(leftMotor*255.0f);
    int rightAmplitude = (int)roundf(rightMotor*255.0f);
    bool vibrating = false;

    jclass deviceClass = (*env)->FindClass(env, "android/view/InputDevice");
    jobject device = (*env)->CallStaticObjectMethod(env, deviceClass, (*env)->GetStaticMethodID(env, deviceClass, "getDevice", "(I)Landroid/view/InputDevice;"), platform.gamepadDeviceId[gamepad]);

    if (!JniFailed(env) && (device != NULL) && (sdkVersion >= 31))
    {
        // Both motors are driven by one combined vibration
        jclass managerClass = (*env)->FindClass(env, "android/os/VibratorManager");
        jobject manager = (*env)->CallObjectMethod(env, device, (*env)->GetMethodID(env, deviceClass, "getVibratorManager", "()Landroid/os/VibratorManager;"));
        jintArray ids = (!JniFailed(env) && (manager != NULL))? (jintArray)(*env)->CallObjectMethod(env, manager, (*env)->GetMethodID(env, managerClass, "getVibratorIds", "()[I")) : NULL;
        int idCount = (!JniFailed(env) && (ids != NULL))? (*env)->GetArrayLength(env, ids) : 0;

        if (idCount >= 2)
        {
            jint idValues[2] = { 0 };
            (*env)->GetIntArrayRegion(env, ids, 0, 2, idValues);

            if ((leftAmplitude == 0) && (rightAmplitude == 0)) (*env)->CallVoidMethod(env, manager, (*env)->GetMethodID(env, managerClass, "cancel", "()V"));
            else
            {
                jclass effectClass = (*env)->FindClass(env, "android/os/VibrationEffect");
                jmethodID createOneShot = (*env)->GetStaticMethodID(env, effectClass, "createOneShot", "(JI)Landroid/os/VibrationEffect;");
                jclass combinedClass = (*env)->FindClass(env, "android/os/CombinedVibration");
                jclass parallelClass = (*env)->FindClass(env, "android/os/CombinedVibration$ParallelCombination");
                jmethodID addVibrator = (*env)->GetMethodID(env, parallelClass, "addVibrator", "(ILandroid/os/VibrationEffect;)Landroid/os/CombinedVibration$ParallelCombination;");
                jobject parallel = (*env)->CallStaticObjectMethod(env, combinedClass, (*env)->GetStaticMethodID(env, combinedClass, "startParallel", "()Landroid/os/CombinedVibration$ParallelCombination;"));
                int amplitudes[2] = { leftAmplitude, rightAmplitude };

                for (int i = 0; (i < 2) && !JniFailed(env); i++)
                {
                    if (amplitudes[i] == 0) continue;
                    jobject effect = (*env)->CallStaticObjectMethod(env, effectClass, createOneShot, milliseconds, (jint)amplitudes[i]);
                    if (!JniFailed(env)) (*env)->CallObjectMethod(env, parallel, addVibrator, idValues[i], effect);
                }

                jobject combined = JniFailed(env)? NULL : (*env)->CallObjectMethod(env, parallel, (*env)->GetMethodID(env, parallelClass, "combine", "()Landroid/os/CombinedVibration;"));
                if (!JniFailed(env) && (combined != NULL)) (*env)->CallVoidMethod(env, manager, (*env)->GetMethodID(env, managerClass, "vibrate", "(Landroid/os/CombinedVibration;)V"), combined);
            }

            vibrating = !JniFailed(env);
        }
    }

    if (!vibrating && (device != NULL))
    {
        jclass vibratorClass = (*env)->FindClass(env, "android/os/Vibrator");
        jobject vibrator = (*env)->CallObjectMethod(env, device, (*env)->GetMethodID(env, deviceClass, "getVibrator", "()Landroid/os/Vibrator;"));
        int amplitude = (leftAmplitude > rightAmplitude)? leftAmplitude : rightAmplitude;

        if (!JniFailed(env) && (vibrator != NULL))
        {
            if (amplitude == 0) (*env)->CallVoidMethod(env, vibrator, (*env)->GetMethodID(env, vibratorClass, "cancel", "()V"));
            else if (sdkVersion >= 26)
            {
                jclass effectClass = (*env)->FindClass(env, "android/os/VibrationEffect");
                jobject effect = (*env)->CallStaticObjectMethod(env, effectClass, (*env)->GetStaticMethodID(env, effectClass, "createOneShot", "(JI)Landroid/os/VibrationEffect;"), milliseconds, (jint)amplitude);
                if (!JniFailed(env)) (*env)->CallVoidMethod(env, vibrator, (*env)->GetMethodID(env, vibratorClass, "vibrate", "(Landroid/os/VibrationEffect;)V"), effect);
            }
            else (*env)->CallVoidMethod(env, vibrator, (*env)->GetMethodID(env, vibratorClass, "vibrate", "(J)V"), milliseconds);

            JniFailed(env);
        }
    }

    DetachJni();
}

// Set mouse position XY
void SetMousePosition(int x, int y)
{
    CORE.Input.Mouse.currentPosition = (Vector2){ (float)x, (float)y };
    CORE.Input.Mouse.previousPosition = CORE.Input.Mouse.currentPosition;
}

// Set mouse cursor
void SetMouseCursor(int cursor)
{
    if ((cursor < MOUSE_CURSOR_DEFAULT) || (cursor > MOUSE_CURSOR_NOT_ALLOWED)) return;

    CORE.Input.Mouse.cursor = cursor;
    if (!CORE.Input.Mouse.cursorHidden)
    {
        atomic_store(&ui.cursor, cursor);
        RequestUiUpdate(UI_REQUEST_CURSOR);
    }
}

// Get physical key name
// NOTE: Character produced by the key with the virtual keyboard layout, only for keys producing characters
const char *GetKeyName(int key)
{
    static char name[8] = { 0 };

    if ((key <= KEY_NULL) || (key == KEY_SPACE)) return NULL;

    int keycode = 0;
    for (int i = 1; i < KEYCODE_MAP_SIZE; i++)
    {
        if (mapKeycode[i] == key)
        {
            keycode = i;
            break;
        }
    }

    if (keycode == 0) return NULL;

    int deadChar = platform.deadChar;
    platform.deadChar = 0;
    int codepoint = GetKeyCharacter(-1, keycode, 0);
    platform.deadChar = deadChar;
    if (codepoint < 32) return NULL;

    int size = 0;
    const char *utf8 = CodepointToUTF8(codepoint, &size);
    memcpy(name, utf8, size);
    name[size] = '\0';

    return name;
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
    // Reset key repeats
    for (int i = 0; i < MAX_KEYBOARD_KEYS; i++) CORE.Input.Keyboard.keyRepeatInFrame[i] = 0;

    // Reset last gamepad button/axis registered state
    CORE.Input.Gamepad.lastButtonPressed = 0; // GAMEPAD_BUTTON_UNKNOWN
    //CORE.Input.Gamepad.axisCount = 0;

    for (int i = 0; i < MAX_GAMEPADS; i++)
    {
        if (CORE.Input.Gamepad.ready[i]) // Check if gamepad is available
        {
            // Register previous gamepad states
            for (int k = 0; k < MAX_GAMEPAD_BUTTONS; k++)
                CORE.Input.Gamepad.previousButtonState[i][k] = CORE.Input.Gamepad.currentButtonState[i][k];
        }
    }

    // Register previous touch states
    for (int i = 0; i < MAX_TOUCH_POINTS; i++) CORE.Input.Touch.previousTouchState[i] = CORE.Input.Touch.currentTouchState[i];

    // Reset touch positions
    //for (int i = 0; i < MAX_TOUCH_POINTS; i++) CORE.Input.Touch.position[i] = (Vector2){ 0, 0 };

    // Register previous keys states
    for (int i = 0; i < MAX_KEYBOARD_KEYS; i++)
    {
        CORE.Input.Keyboard.previousKeyState[i] = CORE.Input.Keyboard.currentKeyState[i];
        CORE.Input.Keyboard.keyRepeatInFrame[i] = 0;
    }

    // Register previous mouse states
    CORE.Input.Mouse.previousPosition = CORE.Input.Mouse.currentPosition;
    for (int i = 0; i < MAX_MOUSE_BUTTONS; i++) CORE.Input.Mouse.previousButtonState[i] = CORE.Input.Mouse.currentButtonState[i];
    CORE.Input.Mouse.previousWheelMove = CORE.Input.Mouse.currentWheelMove;
    CORE.Input.Mouse.currentWheelMove = (Vector2){ 0.0f, 0.0f };

    // Apply releases of keys, mouse buttons and touches that were pressed and released in the previous frame
    bool releasePending = false;
    for (int i = 0; i < MAX_KEYBOARD_KEYS; i++)
    {
        if (platform.keyReleasePending[i])
        {
            CORE.Input.Keyboard.currentKeyState[i] = 0;
            platform.keyReleasePending[i] = false;
            releasePending = true;
        }

        platform.keyPressedInFrame[i] = false;
    }

    for (int i = 0; i < MAX_MOUSE_BUTTONS; i++)
    {
        if (platform.buttonReleasePending[i])
        {
            CORE.Input.Mouse.currentButtonState[i] = 0;
            CORE.Input.Touch.currentTouchState[i] = 0;
            platform.buttonReleasePending[i] = false;
            releasePending = true;
        }

        platform.buttonPressedInFrame[i] = false;
    }

    CORE.Window.resizedLastFrame = false;

    // Display size changes after a configuration change (rotation) are applied when the window is resized
    if (platform.displayCheckEnd > 0.0)
    {
        UpdateDisplaySize();
        if (GetTime() > platform.displayCheckEnd) platform.displayCheckEnd = 0.0;
    }

    CheckGamepads();

    // Android ALooper_pollOnce() variables
    int pollResult = 0;
    int pollEvents = 0;

    // Poll Events (registered events) until TIMEOUT is reached which indicates there are no events left to poll
    // NOTE: Activity is paused if not enabled (platform.appEnabled) and always run flag is not set (FLAG_WINDOW_ALWAYS_RUN),
    // events are waited for if required, a release applied in this frame or a pending display check is an event already
    int timeout = 0;
    if (!platform.appEnabled && !FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_ALWAYS_RUN)) timeout = -1;
    else if (CORE.Window.eventWaiting && !releasePending && (platform.displayCheckEnd == 0.0)) timeout = 1000;

    while ((pollResult = ALooper_pollOnce(timeout, NULL, &pollEvents, ((void **)&platform.source)) > ALOOPER_POLL_TIMEOUT))
    {
        timeout = 0;

        // Process this event
        if (platform.source != NULL) platform.source->process(platform.app, platform.source);

        // NOTE: Allow closing the window in case a configuration change happened
        // The android_main function should be allowed to return to its caller in order for the
        // Android OS to relaunch the activity
        if (platform.app->destroyRequested != 0)
        {
            CORE.Window.shouldClose = true;
            break;
        }
    }
}

//----------------------------------------------------------------------------------
// Module Internal Functions Definition
//----------------------------------------------------------------------------------

// Initialize platform: graphics, inputs and more
int InitPlatform(void)
{
    // Initialize display basic configuration
    //----------------------------------------------------------------------------
    CORE.Window.currentFbo.width = CORE.Window.screen.width;
    CORE.Window.currentFbo.height = CORE.Window.screen.height;
    platform.requestedWidth = CORE.Window.screen.width;
    platform.requestedHeight = CORE.Window.screen.height;
    platform.displayCheckEnd = 0.0;
    platform.deadChar = 0;
    platform.mouseButtonState = 0;
    platform.gamepadCheckTime = 0.0;
    for (int i = 0; i < MAX_GAMEPADS; i++) platform.gamepadDeviceId[i] = -1;

    // Set desired windows flags before initializing anything
    ANativeActivity_setWindowFlags(platform.app->activity, AWINDOW_FLAG_FULLSCREEN, 0);  //AWINDOW_FLAG_SCALED, AWINDOW_FLAG_DITHER

    int orientation = AConfiguration_getOrientation(platform.app->config);

    if (orientation == ACONFIGURATION_ORIENTATION_PORT) TRACELOG(LOG_INFO, "ANDROID: Window orientation set as portrait");
    else if (orientation == ACONFIGURATION_ORIENTATION_LAND) TRACELOG(LOG_INFO, "ANDROID: Window orientation set as landscape");

    // TODO: Automatic orientation doesn't seem to work
    if (CORE.Window.screen.width <= CORE.Window.screen.height)
    {
        AConfiguration_setOrientation(platform.app->config, ACONFIGURATION_ORIENTATION_PORT);
        TRACELOG(LOG_WARNING, "ANDROID: Window orientation changed to portrait");
    }
    else
    {
        AConfiguration_setOrientation(platform.app->config, ACONFIGURATION_ORIENTATION_LAND);
        TRACELOG(LOG_WARNING, "ANDROID: Window orientation changed to landscape");
    }

    //AConfiguration_getDensity(platform.app->config);
    //AConfiguration_getKeyboard(platform.app->config);
    //AConfiguration_getScreenSize(platform.app->config);
    //AConfiguration_getScreenLong(platform.app->config);

    // Set some default window flags
    FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_HIDDEN);       // false
    FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_MINIMIZED);    // false
    FLAG_SET(CORE.Window.flags, FLAG_WINDOW_MAXIMIZED);      // true
    FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_UNFOCUSED);    // false
    //----------------------------------------------------------------------------

    // Initialize App command system
    // NOTE: On APP_CMD_INIT_WINDOW -> InitGraphicsDevice(), InitTimer(), LoadFontDefault()...
    //----------------------------------------------------------------------------
    platform.app->onAppCmd = AndroidCommandCallback;
    //----------------------------------------------------------------------------

    // Initialize input events system
    //----------------------------------------------------------------------------
    platform.app->onInputEvent = AndroidInputCallback;
    //----------------------------------------------------------------------------

    // Initialize storage system
    //----------------------------------------------------------------------------
    CORE.Storage.basePath = platform.app->activity->internalDataPath;   // Define base path for storage
    //----------------------------------------------------------------------------

    TRACELOG(LOG_INFO, "PLATFORM: ANDROID: Initialized successfully");

    // Android ALooper_pollOnce() variables
    int pollResult = 0;
    int pollEvents = 0;

    // Wait for window to be initialized (display and context)
    while (!CORE.Window.ready)
    {
        // Process events until TIMEOUT is reached, which indicates no more events queued
        while ((pollResult = ALooper_pollOnce(0, NULL, &pollEvents, ((void **)&platform.source)) > ALOOPER_POLL_TIMEOUT))
        {
            // Process this event
            if (platform.source != NULL) platform.source->process(platform.app, platform.source);

            // NOTE: It's highly likely destroyRequested will never be non-zero at the start of the activity lifecycle
            //if (platform.app->destroyRequested != 0) CORE.Window.shouldClose = true;
        }
    }

    for (int i = 0; i < MAX_TOUCH_POINTS; i++) touchRaw.hoverPoints[i] = -1;

    // Gamepads connected at start are available in the first frame
    CheckGamepads();

    if (FLAG_IS_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE) || FLAG_IS_SET(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE)) UpdateSystemBars();

    return 0;
}

// Close platform
void ClosePlatform(void)
{
    // Close surface, context and display
    if (platform.device != EGL_NO_DISPLAY)
    {
        eglMakeCurrent(platform.device, EGL_NO_SURFACE, EGL_NO_SURFACE, EGL_NO_CONTEXT);

        if (platform.surface != EGL_NO_SURFACE)
        {
            eglDestroySurface(platform.device, platform.surface);
            platform.surface = EGL_NO_SURFACE;
        }

        if (platform.context != EGL_NO_CONTEXT)
        {
            eglDestroyContext(platform.device, platform.context);
            platform.context = EGL_NO_CONTEXT;
        }

        eglTerminate(platform.device);
        platform.device = EGL_NO_DISPLAY;
    }

    RL_FREE(platform.clipboardText);
    platform.clipboardText = NULL;

    // NOTE: Reset global state in case the activity is being relaunched
    if (platform.app->destroyRequested != 0)
    {
        CORE = (CoreData){ 0 };
        platform = (PlatformData){ 0 };
    }
}

// Initialize display device and framebuffer
// NOTE: width and height represent the screen (framebuffer) desired size, not actual display size
// If width or height are 0, default display size will be used for framebuffer size
// NOTE: returns false in case graphic device could not be created
static int InitGraphicsDevice(void)
{
    EGLint samples = 0;
    EGLint sampleBuffer = 0;
    if (FLAG_IS_SET(CORE.Window.flags, FLAG_MSAA_4X_HINT))
    {
        samples = 4;
        sampleBuffer = 1;
        TRACELOG(LOG_INFO, "DISPLAY: Trying to enable MSAA x4");
    }

    const EGLint framebufferAttribs[] = {
        EGL_RENDERABLE_TYPE, (rlGetVersion() == RL_OPENGL_ES_30)? EGL_OPENGL_ES3_BIT : EGL_OPENGL_ES2_BIT, // Type of context support
        EGL_RED_SIZE, 8,            // RED color bit depth (alternative: 5)
        EGL_GREEN_SIZE, 8,          // GREEN color bit depth (alternative: 6)
        EGL_BLUE_SIZE, 8,           // BLUE color bit depth (alternative: 5)
        //EGL_TRANSPARENT_TYPE, EGL_NONE, // Request transparent framebuffer (EGL_TRANSPARENT_RGB does not work on RPI)
        EGL_DEPTH_SIZE, 24,         // Depth buffer size (Required to use Depth testing!)
        //EGL_STENCIL_SIZE, 8,      // Stencil buffer size
        EGL_SAMPLE_BUFFERS, sampleBuffer, // Activate MSAA
        EGL_SAMPLES, samples,       // 4x Antialiasing if activated (Free on MALI GPUs)
        EGL_NONE
    };

    const EGLint contextAttribs[] = {
        EGL_CONTEXT_CLIENT_VERSION, 2,
        EGL_NONE
    };

    EGLint numConfigs = 0;

    // Get an EGL device connection
    platform.device = eglGetDisplay(EGL_DEFAULT_DISPLAY);

    if (platform.device == EGL_NO_DISPLAY)
    {
        TRACELOG(LOG_WARNING, "DISPLAY: Failed to initialize EGL device");
        return -1;
    }

    // Initialize the EGL device connection
    if (eglInitialize(platform.device, NULL, NULL) == EGL_FALSE)
    {
        // If all of the calls to eglInitialize returned EGL_FALSE then an error has occurred
        TRACELOG(LOG_WARNING, "DISPLAY: Failed to initialize EGL device");
        return -1;
    }

    // Get an appropriate EGL framebuffer configuration
    eglChooseConfig(platform.device, framebufferAttribs, &platform.config, 1, &numConfigs);

    // Set rendering API
    eglBindAPI(EGL_OPENGL_ES_API);

    // Create an EGL rendering context
    platform.context = eglCreateContext(platform.device, platform.config, EGL_NO_CONTEXT, contextAttribs);
    if (platform.context == EGL_NO_CONTEXT)
    {
        TRACELOG(LOG_WARNING, "DISPLAY: Failed to create EGL context");
        return -1;
    }

    // Create an EGL window surface
    //---------------------------------------------------------------------------------
    EGLint displayFormat = 0;

    // EGL_NATIVE_VISUAL_ID is an attribute of the EGLConfig that is guaranteed to be accepted by ANativeWindow_setBuffersGeometry()
    // As soon as an EGLConfig is picked, it's safe to reconfigure the ANativeWindow buffers to match, using EGL_NATIVE_VISUAL_ID
    eglGetConfigAttrib(platform.device, platform.config, EGL_NATIVE_VISUAL_ID, &displayFormat);

    // At this point render size vs screen size needs to be managed
    // NOTE: This function use and modify global module variables:
    //  -> CORE.Window.screen.width/CORE.Window.screen.height
    //  -> CORE.Window.render.width/CORE.Window.render.height
    //  -> CORE.Window.screenScale
    SetupFramebuffer(CORE.Window.display.width, CORE.Window.display.height);

    ANativeWindow_setBuffersGeometry(platform.app->window, CORE.Window.render.width, CORE.Window.render.height, displayFormat);
    //ANativeWindow_setBuffersGeometry(platform.app->window, 0, 0, displayFormat);       // Force use of native display size

    platform.surface = eglCreateWindowSurface(platform.device, platform.config, platform.app->window, NULL);

    // There must be at least one frame displayed before the buffers are swapped
    //eglSwapInterval(platform.device, 1);

    if (eglMakeCurrent(platform.device, platform.surface, platform.surface, platform.context) == EGL_FALSE)
    {
        TRACELOG(LOG_WARNING, "DISPLAY: Failed to attach EGL rendering context to EGL surface");
        return -1;
    }
    else
    {
        CORE.Window.render.width = CORE.Window.screen.width;
        CORE.Window.render.height = CORE.Window.screen.height;
        CORE.Window.currentFbo.width = CORE.Window.render.width;
        CORE.Window.currentFbo.height = CORE.Window.render.height;

        TRACELOG(LOG_INFO, "DISPLAY: Device initialized successfully");
        TRACELOG(LOG_INFO, "    > Display size: %i x %i", CORE.Window.display.width, CORE.Window.display.height);
        TRACELOG(LOG_INFO, "    > Screen size:  %i x %i", CORE.Window.screen.width, CORE.Window.screen.height);
        TRACELOG(LOG_INFO, "    > Render size:  %i x %i", CORE.Window.render.width, CORE.Window.render.height);
        TRACELOG(LOG_INFO, "    > Viewport offsets: %i, %i", CORE.Window.renderOffset.x, CORE.Window.renderOffset.y);
    }

    // Load OpenGL extensions
    // NOTE: GL procedures address loader is required to load extensions
    rlLoadExtensions(eglGetProcAddress);

    CORE.Window.ready = true;

    if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_MINIMIZED)) MinimizeWindow();

    return 0;
}

// ANDROID: Process activity lifecycle commands
static void AndroidCommandCallback(struct android_app *app, int32_t cmd)
{
    switch (cmd)
    {
        case APP_CMD_START: FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_MINIMIZED); break;
        case APP_CMD_RESUME: break;
        case APP_CMD_INIT_WINDOW:
        {
            if (app->window != NULL)
            {
                if (platform.contextRebindRequired)
                {
                    // Reset screen scaling to full display size
                    EGLint displayFormat = 0;
                    eglGetConfigAttrib(platform.device, platform.config, EGL_NATIVE_VISUAL_ID, &displayFormat);

                    // Adding renderOffset here feels rather hackish, but the viewport scaling is wrong after the
                    // context rebinding if the screen is scaled unless offsets are added. There's probably a more
                    // appropriate way to fix this
                    ANativeWindow_setBuffersGeometry(app->window,
                        CORE.Window.render.width + CORE.Window.renderOffset.x,
                        CORE.Window.render.height + CORE.Window.renderOffset.y,
                        displayFormat);

                    // Recreate display surface and re-attach OpenGL context
                    platform.surface = eglCreateWindowSurface(platform.device, platform.config, app->window, NULL);
                    eglMakeCurrent(platform.device, platform.surface, platform.surface, platform.context);

                    platform.contextRebindRequired = false;

                    // NOTE: Display can be rotated while the window was destroyed
                    UpdateDisplaySize();
                }
                else
                {
                    CORE.Window.display.width = ANativeWindow_getWidth(platform.app->window);
                    CORE.Window.display.height = ANativeWindow_getHeight(platform.app->window);

                    // Initialize graphics device (display device and OpenGL context)
                    InitGraphicsDevice();

                    // Initialize OpenGL context (states and resources)
                    // NOTE: CORE.Window.currentFbo.width and CORE.Window.currentFbo.height not used, stored as globals in rlgl
                    rlglInit(CORE.Window.currentFbo.width, CORE.Window.currentFbo.height);

                    // Setup default viewport
                    // NOTE: It updated CORE.Window.render.width and CORE.Window.render.height
                    SetupViewport(CORE.Window.currentFbo.width, CORE.Window.currentFbo.height);

                    // Initialize hi-res timer
                    InitTimer();

                #if SUPPORT_MODULE_RTEXT
                    // Load default font
                    // WARNING: External function: Module required: rtext
                    LoadFontDefault();
                    #if SUPPORT_MODULE_RSHAPES
                    // Set font white rectangle for shapes drawing, so shapes and text can be batched together
                    // WARNING: rshapes module is required, if not available, default internal white rectangle is used
                    Rectangle rec = GetFontDefault().recs[95];
                    if (FLAG_IS_SET(CORE.Window.flags, FLAG_MSAA_4X_HINT))
                    {
                        // NOTE: Trying to maxime rec padding to avoid pixel bleeding on MSAA filtering
                        SetShapesTexture(GetFontDefault().texture, (Rectangle){ rec.x + 2, rec.y + 2, 1, 1 });
                    }
                    else
                    {
                        // NOTE: Setting up a 1px padding on char rectangle to avoid pixel bleeding
                        SetShapesTexture(GetFontDefault().texture, (Rectangle){ rec.x + 1, rec.y + 1, rec.width - 2, rec.height - 2 });
                    }
                    #endif
                #else
                    #if SUPPORT_MODULE_RSHAPES
                    // Set default texture and rectangle to be used for shapes drawing
                    // NOTE: rlgl default texture is a 1x1 pixel UNCOMPRESSED_R8G8B8A8
                    Texture2D texture = { rlGetTextureIdDefault(), 1, 1, 1, PIXELFORMAT_UNCOMPRESSED_R8G8B8A8 };
                    SetShapesTexture(texture, (Rectangle){ 0.0f, 0.0f, 1.0f, 1.0f });    // WARNING: Module required: rshapes
                    #endif
                #endif

                    // Initialize random seed
                    SetRandomSeed((unsigned int)time(NULL));
                }

                FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_HIDDEN);
            }
        } break;
        case APP_CMD_GAINED_FOCUS:
        {
            platform.appEnabled = true;
            FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_UNFOCUSED);

            // NOTE: Pointer capture is lost with the focus, system bars can be shown again
            if (CORE.Input.Mouse.cursorLocked) RequestUiUpdate(UI_REQUEST_CAPTURE);
            if (FLAG_IS_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE) || FLAG_IS_SET(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE)) UpdateSystemBars();
            //ResumeMusicStream();
        } break;
        case APP_CMD_PAUSE: break;
        case APP_CMD_LOST_FOCUS:
        {
            platform.appEnabled = false;
            FLAG_SET(CORE.Window.flags, FLAG_WINDOW_UNFOCUSED);
            //PauseMusicStream();
        } break;
        case APP_CMD_TERM_WINDOW:
        {
            // Detach OpenGL context and destroy display surface
            // NOTE 1: This case is used when the user exits the app without closing it, context is detached to ensure everything is recoverable upon resuming
            // NOTE 2: Detaching context before destroying display surface avoids losing our resources (textures, shaders, VBOs...)
            // NOTE 3: In some cases (too many context loaded), OS could unload context automatically... :(
            if (platform.device != EGL_NO_DISPLAY)
            {
                eglMakeCurrent(platform.device, EGL_NO_SURFACE, EGL_NO_SURFACE, EGL_NO_CONTEXT);

                if (platform.surface != EGL_NO_SURFACE)
                {
                    eglDestroySurface(platform.device, platform.surface);
                    platform.surface = EGL_NO_SURFACE;
                }

                platform.contextRebindRequired = true;
            }
            // If 'platform.device' is already set to 'EGL_NO_DISPLAY'
            // this means that the user has already called 'CloseWindow()'

            FLAG_SET(CORE.Window.flags, FLAG_WINDOW_HIDDEN);
        } break;
        case APP_CMD_SAVE_STATE: break;
        case APP_CMD_STOP: FLAG_SET(CORE.Window.flags, FLAG_WINDOW_MINIMIZED); break;
        case APP_CMD_DESTROY: break;
        case APP_CMD_CONFIG_CHANGED:
        {
            // NOTE: Window is resized after the configuration change (rotation), its size is checked for a while
            platform.displayCheckEnd = GetTime() + 2.0;
        } break;
        default: break;
    }
}

// ANDROID: Map Android gamepad button to raylib gamepad button
static GamepadButton AndroidTranslateGamepadButton(int button)
{
    switch (button)
    {
        case AKEYCODE_BUTTON_A: return GAMEPAD_BUTTON_RIGHT_FACE_DOWN;
        case AKEYCODE_BUTTON_B: return GAMEPAD_BUTTON_RIGHT_FACE_RIGHT;
        case AKEYCODE_BUTTON_X: return GAMEPAD_BUTTON_RIGHT_FACE_LEFT;
        case AKEYCODE_BUTTON_Y: return GAMEPAD_BUTTON_RIGHT_FACE_UP;
        case AKEYCODE_BUTTON_L1: return GAMEPAD_BUTTON_LEFT_TRIGGER_1;
        case AKEYCODE_BUTTON_R1: return GAMEPAD_BUTTON_RIGHT_TRIGGER_1;
        case AKEYCODE_BUTTON_L2: return GAMEPAD_BUTTON_LEFT_TRIGGER_2;
        case AKEYCODE_BUTTON_R2: return GAMEPAD_BUTTON_RIGHT_TRIGGER_2;
        case AKEYCODE_BUTTON_THUMBL: return GAMEPAD_BUTTON_LEFT_THUMB;
        case AKEYCODE_BUTTON_THUMBR: return GAMEPAD_BUTTON_RIGHT_THUMB;
        case AKEYCODE_BUTTON_START: return GAMEPAD_BUTTON_MIDDLE_RIGHT;
        case AKEYCODE_BUTTON_SELECT: return GAMEPAD_BUTTON_MIDDLE_LEFT;
        case AKEYCODE_BUTTON_MODE: return GAMEPAD_BUTTON_MIDDLE;
        // On some (most?) gamepads dpad events are reported as axis motion instead
        case AKEYCODE_DPAD_DOWN: return GAMEPAD_BUTTON_LEFT_FACE_DOWN;
        case AKEYCODE_DPAD_RIGHT: return GAMEPAD_BUTTON_LEFT_FACE_RIGHT;
        case AKEYCODE_DPAD_LEFT: return GAMEPAD_BUTTON_LEFT_FACE_LEFT;
        case AKEYCODE_DPAD_UP: return GAMEPAD_BUTTON_LEFT_FACE_UP;
        default: return GAMEPAD_BUTTON_UNKNOWN;
    }
}

// ANDROID: Get input events
static int32_t AndroidInputCallback(struct android_app *app, AInputEvent *event)
{
    // If additional inputs are required check:
    // REF: https://developer.android.com/ndk/reference/group/input
    // REF: https://developer.android.com/training/game-controllers/controller-input

    int type = AInputEvent_getType(event);
    int source = AInputEvent_getSource(event);

    if (type == AINPUT_EVENT_TYPE_MOTION)
    {
        if (FLAG_IS_SET(source, AINPUT_SOURCE_JOYSTICK) ||
            FLAG_IS_SET(source, AINPUT_SOURCE_GAMEPAD))
        {
            int gamepad = GetGamepadIndex(AInputEvent_getDeviceId(event), true);
            if (gamepad < 0) return 1;

            CORE.Input.Gamepad.axisState[gamepad][GAMEPAD_AXIS_LEFT_X] = AMotionEvent_getAxisValue(
                    event, AMOTION_EVENT_AXIS_X, 0);
            CORE.Input.Gamepad.axisState[gamepad][GAMEPAD_AXIS_LEFT_Y] = AMotionEvent_getAxisValue(
                    event, AMOTION_EVENT_AXIS_Y, 0);
            CORE.Input.Gamepad.axisState[gamepad][GAMEPAD_AXIS_RIGHT_X] = AMotionEvent_getAxisValue(
                    event, AMOTION_EVENT_AXIS_Z, 0);
            CORE.Input.Gamepad.axisState[gamepad][GAMEPAD_AXIS_RIGHT_Y] = AMotionEvent_getAxisValue(
                    event, AMOTION_EVENT_AXIS_RZ, 0);
            CORE.Input.Gamepad.axisState[gamepad][GAMEPAD_AXIS_LEFT_TRIGGER] = AMotionEvent_getAxisValue(
                    event, AMOTION_EVENT_AXIS_BRAKE, 0)*2.0f - 1.0f;
            CORE.Input.Gamepad.axisState[gamepad][GAMEPAD_AXIS_RIGHT_TRIGGER] = AMotionEvent_getAxisValue(
                    event, AMOTION_EVENT_AXIS_GAS, 0)*2.0f - 1.0f;

            // dpad is reported as an axis on android
            float dpadX = AMotionEvent_getAxisValue(event, AMOTION_EVENT_AXIS_HAT_X, 0);
            float dpadY = AMotionEvent_getAxisValue(event, AMOTION_EVENT_AXIS_HAT_Y, 0);

            if (dpadX == 1.0f)
            {
                CORE.Input.Gamepad.currentButtonState[gamepad][GAMEPAD_BUTTON_LEFT_FACE_RIGHT] = 1;
                CORE.Input.Gamepad.currentButtonState[gamepad][GAMEPAD_BUTTON_LEFT_FACE_LEFT] = 0;
            }
            else if (dpadX == -1.0f)
            {
                CORE.Input.Gamepad.currentButtonState[gamepad][GAMEPAD_BUTTON_LEFT_FACE_RIGHT] = 0;
                CORE.Input.Gamepad.currentButtonState[gamepad][GAMEPAD_BUTTON_LEFT_FACE_LEFT] = 1;
            }
            else
            {
                CORE.Input.Gamepad.currentButtonState[gamepad][GAMEPAD_BUTTON_LEFT_FACE_RIGHT] = 0;
                CORE.Input.Gamepad.currentButtonState[gamepad][GAMEPAD_BUTTON_LEFT_FACE_LEFT] = 0;
            }

            if (dpadY == 1.0f)
            {
                CORE.Input.Gamepad.currentButtonState[gamepad][GAMEPAD_BUTTON_LEFT_FACE_DOWN] = 1;
                CORE.Input.Gamepad.currentButtonState[gamepad][GAMEPAD_BUTTON_LEFT_FACE_UP] = 0;
            }
            else if (dpadY == -1.0f)
            {
                CORE.Input.Gamepad.currentButtonState[gamepad][GAMEPAD_BUTTON_LEFT_FACE_DOWN] = 0;
                CORE.Input.Gamepad.currentButtonState[gamepad][GAMEPAD_BUTTON_LEFT_FACE_UP] = 1;
            }
            else
            {
                CORE.Input.Gamepad.currentButtonState[gamepad][GAMEPAD_BUTTON_LEFT_FACE_DOWN] = 0;
                CORE.Input.Gamepad.currentButtonState[gamepad][GAMEPAD_BUTTON_LEFT_FACE_UP] = 0;
            }

            return 1; // Handled gamepad axis motion
        }

        bool relative = FLAG_IS_SET(source, AINPUT_SOURCE_MOUSE_RELATIVE);

        if (relative || (FLAG_IS_SET(source, AINPUT_SOURCE_MOUSE) && (AMotionEvent_getToolType(event, 0) == AMOTION_EVENT_TOOL_TYPE_MOUSE)))
        {
            int32_t action = AMotionEvent_getAction(event) & AMOTION_EVENT_ACTION_MASK;

            // Mouse position, buttons and wheel
            // NOTE: Captured pointer (disabled cursor) reports relative motion
            if (relative)
            {
                Vector2 scale = { (float)(CORE.Window.screen.width + CORE.Window.renderOffset.x)/(float)CORE.Window.display.width, (float)(CORE.Window.screen.height + CORE.Window.renderOffset.y)/(float)CORE.Window.display.height };
                CORE.Input.Mouse.currentPosition.x += AMotionEvent_getX(event, 0)*scale.x;
                CORE.Input.Mouse.currentPosition.y += AMotionEvent_getY(event, 0)*scale.y;
                CORE.Input.Mouse.cursorOnScreen = true;
            }
            else if (!CORE.Input.Mouse.cursorLocked)
            {
                CORE.Input.Mouse.currentPosition = GetScreenPosition(AMotionEvent_getX(event, 0), AMotionEvent_getY(event, 0));
                CORE.Input.Mouse.cursorOnScreen = (action != AMOTION_EVENT_ACTION_HOVER_EXIT);
            }

            CORE.Input.Touch.position[0] = CORE.Input.Mouse.currentPosition;

            if (action == AMOTION_EVENT_ACTION_SCROLL)
            {
                // NOTE: Horizontal wheel movement is inverted to match the desktop platforms
                CORE.Input.Mouse.currentWheelMove.x -= AMotionEvent_getAxisValue(event, AMOTION_EVENT_AXIS_HSCROLL, 0);
                CORE.Input.Mouse.currentWheelMove.y += AMotionEvent_getAxisValue(event, AMOTION_EVENT_AXIS_VSCROLL, 0);
            }

            static const struct { int32_t state; int button; } buttons[] = {
                { AMOTION_EVENT_BUTTON_PRIMARY, MOUSE_BUTTON_LEFT },
                { AMOTION_EVENT_BUTTON_SECONDARY, MOUSE_BUTTON_RIGHT },
                { AMOTION_EVENT_BUTTON_TERTIARY, MOUSE_BUTTON_MIDDLE },
                { AMOTION_EVENT_BUTTON_BACK, MOUSE_BUTTON_BACK },
                { AMOTION_EVENT_BUTTON_FORWARD, MOUSE_BUTTON_FORWARD }
            };

            // NOTE: Button state is not reported by some events, a primary button is down between down and up actions
            int32_t buttonState = AMotionEvent_getButtonState(event);
            if ((action == AMOTION_EVENT_ACTION_DOWN) && (buttonState == 0)) buttonState = AMOTION_EVENT_BUTTON_PRIMARY;
            else if ((action == AMOTION_EVENT_ACTION_UP) || (action == AMOTION_EVENT_ACTION_CANCEL)) buttonState = 0;
            else if ((action == AMOTION_EVENT_ACTION_MOVE) && (buttonState == 0)) buttonState = platform.mouseButtonState;

            for (int i = 0; i < 5; i++)
            {
                bool down = FLAG_IS_SET(buttonState, buttons[i].state);
                if (down != FLAG_IS_SET(platform.mouseButtonState, buttons[i].state)) HandleMouseButton(buttons[i].button, down);
            }

            platform.mouseButtonState = buttonState;

            // NOTE: Hover events are also passed to the views, the pointer icon is resolved there
            if ((action == AMOTION_EVENT_ACTION_HOVER_MOVE) || (action == AMOTION_EVENT_ACTION_HOVER_ENTER) || (action == AMOTION_EVENT_ACTION_HOVER_EXIT)) return 0;

            return 1; // Handled mouse event
        }
    }
    else if (type == AINPUT_EVENT_TYPE_KEY)
    {
        int32_t keycode = AKeyEvent_getKeyCode(event);
        //int32_t AKeyEvent_getMetaState(event);

        // Handle gamepad button presses and releases. AOSP stamps the
        // KEYBOARD source bit on every key event from a gamepad, so
        // discriminate on the keycode rather than gating on source bits.
        if (FLAG_IS_SET(source, AINPUT_SOURCE_JOYSTICK) ||
            FLAG_IS_SET(source, AINPUT_SOURCE_GAMEPAD))
        {
            GamepadButton button = AndroidTranslateGamepadButton(keycode);
            int gamepad = (button != GAMEPAD_BUTTON_UNKNOWN)? GetGamepadIndex(AInputEvent_getDeviceId(event), true) : -1;

            if (gamepad >= 0)
            {
                if (AKeyEvent_getAction(event) == AKEY_EVENT_ACTION_DOWN)
                {
                    CORE.Input.Gamepad.currentButtonState[gamepad][button] = 1;
                    CORE.Input.Gamepad.lastButtonPressed = button;
                }
                else CORE.Input.Gamepad.currentButtonState[gamepad][button] = 0;  // Key up

                return 1; // Handled gamepad button
            }
            // Unknown keycode: fall through to the keyboard handler below.
        }

        KeyboardKey key = ((keycode > 0) && (keycode < KEYCODE_MAP_SIZE))? mapKeycode[keycode] : KEY_NULL;
        int32_t action = AKeyEvent_getAction(event);

        if (key != KEY_NULL)
        {
            // Save current key and its state
            // NOTE: Repeated key presses are reported as down actions with a repeat count
            if (action == AKEY_EVENT_ACTION_DOWN)
            {
                if (AKeyEvent_getRepeatCount(event) > 0) CORE.Input.Keyboard.keyRepeatInFrame[key] = 1;
                else
                {
                    HandleKey(key, true);
                    if (key == CORE.Input.Keyboard.exitKey) CORE.Window.shouldClose = true;
                }
            }
            else if (action == AKEY_EVENT_ACTION_MULTIPLE) CORE.Input.Keyboard.keyRepeatInFrame[key] = 1;
            else HandleKey(key, false);
        }

        // Register character input, also for repeated key presses
        if ((action == AKEY_EVENT_ACTION_DOWN) && (keycode != AKEYCODE_BACK) && (keycode != AKEYCODE_MENU))
        {
            int codepoint = GetKeyCharacter(AInputEvent_getDeviceId(event), keycode, AKeyEvent_getMetaState(event));

            // Check if there is space available in the queue, control characters are not registered
            if ((codepoint >= 32) && (codepoint != 127) && (CORE.Input.Keyboard.charPressedQueueCount < MAX_CHAR_PRESSED_QUEUE))
            {
                // Add character to the queue
                CORE.Input.Keyboard.charPressedQueue[CORE.Input.Keyboard.charPressedQueueCount] = codepoint;
                CORE.Input.Keyboard.charPressedQueueCount++;
            }
        }

        if (keycode == AKEYCODE_POWER)
        {
            // Let the OS handle input to avoid app stuck. Behaviour: CMD_PAUSE -> CMD_SAVE_STATE -> CMD_STOP -> CMD_CONFIG_CHANGED -> CMD_LOST_FOCUS
            // Resuming Behaviour: CMD_START -> CMD_RESUME -> CMD_CONFIG_CHANGED -> CMD_CONFIG_CHANGED -> CMD_GAINED_FOCUS
            // It seems like locking mobile, screen size (CMD_CONFIG_CHANGED) is affected
            // NOTE: AndroidManifest.xml must have <activity android:configChanges="orientation|keyboardHidden|screenSize" >
            // Before that change, activity was calling CMD_TERM_WINDOW and CMD_DESTROY when locking mobile, so that was not a normal behaviour
            return 0;
        }
        else if ((keycode == AKEYCODE_BACK) || (keycode == AKEYCODE_MENU))
        {
            // Eat BACK_BUTTON and AKEYCODE_MENU, do nothing... and don't let to be handled by OS!
            return 1;
        }
        else if ((keycode == AKEYCODE_VOLUME_UP) || (keycode == AKEYCODE_VOLUME_DOWN))
        {
            // Set default OS behaviour
            return 0;
        }

        return 0;
    }

    // Register touch points count
    touchRaw.pointCount = AMotionEvent_getPointerCount(event);

    for (int i = 0; (i < touchRaw.pointCount) && (i < MAX_TOUCH_POINTS); i++)
    {
        // Register touch points id
        touchRaw.pointId[i] = AMotionEvent_getPointerId(event, i);

        // Register touch points position, normalized for CORE.Window.screen.width and CORE.Window.screen.height
        touchRaw.position[i] = GetScreenPosition(AMotionEvent_getX(event, i), AMotionEvent_getY(event, i));
    }

    int32_t action = AMotionEvent_getAction(event);
    unsigned int flags = action & AMOTION_EVENT_ACTION_MASK;
    int32_t pointerIndex = (action & AMOTION_EVENT_ACTION_POINTER_INDEX_MASK) >> AMOTION_EVENT_ACTION_POINTER_INDEX_SHIFT;

    if (flags == AMOTION_EVENT_ACTION_HOVER_ENTER)
    {
        // The new pointer is hover
        // So add it to hoverPoints
        for (int i = 0; i < MAX_TOUCH_POINTS; i++)
        {
            if (touchRaw.hoverPoints[i] == -1)
            {
                touchRaw.hoverPoints[i] = touchRaw.pointId[pointerIndex];
                break;
            }
        }
    }

#if SUPPORT_GESTURES_SYSTEM
    GestureEvent gestureEvent = { 0 };

    gestureEvent.pointCount = 0;

    // Register touch actions
    if (flags == AMOTION_EVENT_ACTION_DOWN) gestureEvent.touchAction = TOUCH_ACTION_DOWN;
    else if (flags == AMOTION_EVENT_ACTION_UP) gestureEvent.touchAction = TOUCH_ACTION_UP;
    else if (flags == AMOTION_EVENT_ACTION_MOVE) gestureEvent.touchAction = TOUCH_ACTION_MOVE;
    else if (flags == AMOTION_EVENT_ACTION_CANCEL) gestureEvent.touchAction = TOUCH_ACTION_CANCEL;

    for (int i = 0; (i < touchRaw.pointCount) && (i < MAX_TOUCH_POINTS); i++)
    {
        // If the touchPoint is hover, Ignore it
        bool hover = false;
        for (int j = 0; j < MAX_TOUCH_POINTS; j++)
        {
            // Check if the touchPoint is in hoverPointers
            if (touchRaw.hoverPoints[j] == touchRaw.pointId[i])
            {
                hover = true;
                break;
            }
        }
        if (hover) continue;

        gestureEvent.pointId[gestureEvent.pointCount] = touchRaw.pointId[i];
        gestureEvent.position[gestureEvent.pointCount] = touchRaw.position[i];
        gestureEvent.position[gestureEvent.pointCount].x /= (float)GetScreenWidth();
        gestureEvent.position[gestureEvent.pointCount].y /= (float)GetScreenHeight();
        gestureEvent.pointCount++;
    }

    // Gesture data is sent to gestures system for processing
    ProcessGestureEvent(gestureEvent);
#endif

    if (flags == AMOTION_EVENT_ACTION_HOVER_EXIT)
    {
        // Hover exited. So, remove it from hoverPoints
        for (int i = 0; i < MAX_TOUCH_POINTS; i++)
        {
            if (touchRaw.hoverPoints[i] == touchRaw.pointId[pointerIndex])
            {
                touchRaw.hoverPoints[i] = -1;
                break;
            }
        }
    }

    if ((flags == AMOTION_EVENT_ACTION_POINTER_UP) || (flags == AMOTION_EVENT_ACTION_UP))
    {
        // One of the touchpoints is released, remove it from touch point arrays
        for (int i = pointerIndex; (i < touchRaw.pointCount - 1) && (i < MAX_TOUCH_POINTS - 1); i++)
        {
            touchRaw.pointId[i] = touchRaw.pointId[i+1];
            touchRaw.position[i] = touchRaw.position[i+1];
        }
        touchRaw.pointCount--;
    }

    CORE.Input.Touch.pointCount = 0;
    for (int i = 0; (i < touchRaw.pointCount) && (i < MAX_TOUCH_POINTS); i++)
    {
        // If the touchPoint is hover, Ignore it
        bool hover = false;
        for (int j = 0; j < MAX_TOUCH_POINTS; j++)
        {
            // Check if the touchPoint is in hoverPointers
            if (touchRaw.hoverPoints[j] == touchRaw.pointId[i])
            {
                hover = true;
                break;
            }
        }
        if (hover) continue;

        CORE.Input.Touch.pointId[CORE.Input.Touch.pointCount] = touchRaw.pointId[i];
        CORE.Input.Touch.position[CORE.Input.Touch.pointCount] = touchRaw.position[i];
        CORE.Input.Touch.pointCount++;
    }

    // When all touchpoints are tapped and released really quickly, this event is generated
    if (flags == AMOTION_EVENT_ACTION_CANCEL) CORE.Input.Touch.pointCount = 0;

    bool touchDown = (CORE.Input.Touch.pointCount > 0);
    if (touchDown != (CORE.Input.Touch.currentTouchState[MOUSE_BUTTON_LEFT] == 1) || (!touchDown && platform.buttonReleasePending[MOUSE_BUTTON_LEFT])) HandleMouseButton(MOUSE_BUTTON_LEFT, touchDown);

    // Touch down and up do not move the mouse, the delta is only calculated while touch[0] is moving
    if (flags != AMOTION_EVENT_ACTION_MOVE) CORE.Input.Mouse.previousPosition = CORE.Input.Touch.position[0];

    // Map touch[0] as mouse input for convenience
    CORE.Input.Mouse.currentPosition = CORE.Input.Touch.position[0];
    CORE.Input.Mouse.currentWheelMove = (Vector2){ 0.0f, 0.0f };

    return 0;
}

// Compute framebuffer size relative to screen size and display size
// NOTE: Global variables CORE.Window.render.width/CORE.Window.render.height and CORE.Window.renderOffset.x/CORE.Window.renderOffset.y can be modified
static void SetupFramebuffer(int width, int height)
{
    // Calculate CORE.Window.render.width and CORE.Window.render.height, having the display size (input params) and the desired screen size (global var)
    if ((CORE.Window.screen.width > CORE.Window.display.width) || (CORE.Window.screen.height > CORE.Window.display.height))
    {
        TRACELOG(LOG_WARNING, "DISPLAY: Downscaling required: Screen size (%ix%i) is bigger than display size (%ix%i)", CORE.Window.screen.width, CORE.Window.screen.height, CORE.Window.display.width, CORE.Window.display.height);

        // Downscaling to fit display with border-bars
        float widthRatio = (float)CORE.Window.display.width/(float)CORE.Window.screen.width;
        float heightRatio = (float)CORE.Window.display.height/(float)CORE.Window.screen.height;

        if (widthRatio <= heightRatio)
        {
            CORE.Window.render.width = CORE.Window.display.width;
            CORE.Window.render.height = (int)round((float)CORE.Window.screen.height*widthRatio);
            CORE.Window.renderOffset.x = 0;
            CORE.Window.renderOffset.y = (CORE.Window.display.height - CORE.Window.render.height);
        }
        else
        {
            CORE.Window.render.width = (int)round((float)CORE.Window.screen.width*heightRatio);
            CORE.Window.render.height = CORE.Window.display.height;
            CORE.Window.renderOffset.x = (CORE.Window.display.width - CORE.Window.render.width);
            CORE.Window.renderOffset.y = 0;
        }

        // Screen scaling required
        float scaleRatio = (float)CORE.Window.render.width/(float)CORE.Window.screen.width;
        CORE.Window.screenScale = MatrixScale(scaleRatio, scaleRatio, 1.0f);

        // NOTE: Rendering to full display resolution
        // Above parameters need to be calculate for downscale matrix and offsets
        CORE.Window.render.width = CORE.Window.display.width;
        CORE.Window.render.height = CORE.Window.display.height;

        TRACELOG(LOG_WARNING, "DISPLAY: Downscale matrix generated, content will be rendered at (%ix%i)", CORE.Window.render.width, CORE.Window.render.height);
    }
    else if ((CORE.Window.screen.width < CORE.Window.display.width) || (CORE.Window.screen.height < CORE.Window.display.height))
    {
        // Required screen size is smaller than display size
        TRACELOG(LOG_INFO, "DISPLAY: Upscaling required: Screen size (%ix%i) smaller than display size (%ix%i)", CORE.Window.screen.width, CORE.Window.screen.height, CORE.Window.display.width, CORE.Window.display.height);

        if ((CORE.Window.screen.width == 0) || (CORE.Window.screen.height == 0))
        {
            CORE.Window.screen.width = CORE.Window.display.width;
            CORE.Window.screen.height = CORE.Window.display.height;
        }

        // Upscaling to fit display with border-bars
        float displayRatio = (float)CORE.Window.display.width/(float)CORE.Window.display.height;
        float screenRatio = (float)CORE.Window.screen.width/(float)CORE.Window.screen.height;

        if (displayRatio <= screenRatio)
        {
            CORE.Window.render.width = CORE.Window.screen.width;
            CORE.Window.render.height = (int)round((float)CORE.Window.screen.width/displayRatio);
            CORE.Window.renderOffset.x = 0;
            CORE.Window.renderOffset.y = (CORE.Window.render.height - CORE.Window.screen.height);
        }
        else
        {
            CORE.Window.render.width = (int)round((float)CORE.Window.screen.height*displayRatio);
            CORE.Window.render.height = CORE.Window.screen.height;
            CORE.Window.renderOffset.x = (CORE.Window.render.width - CORE.Window.screen.width);
            CORE.Window.renderOffset.y = 0;
        }
    }
    else
    {
        CORE.Window.render.width = CORE.Window.screen.width;
        CORE.Window.render.height = CORE.Window.screen.height;
        CORE.Window.renderOffset.x = 0;
        CORE.Window.renderOffset.y = 0;
    }
}

// Update framebuffer and viewport if the display size changed (rotation)
static void UpdateDisplaySize(void)
{
    if ((platform.app->window == NULL) || (platform.surface == EGL_NO_SURFACE)) return;

    int width = ANativeWindow_getWidth(platform.app->window);
    int height = ANativeWindow_getHeight(platform.app->window);
    if ((width <= 0) || (height <= 0) || ((width == CORE.Window.display.width) && (height == CORE.Window.display.height))) return;

    TRACELOG(LOG_INFO, "DISPLAY: Display size changed from %i x %i to %i x %i", CORE.Window.display.width, CORE.Window.display.height, width, height);

    CORE.Window.display.width = width;
    CORE.Window.display.height = height;

    // Screen size requested by the application is kept, it is scaled to the new display size
    CORE.Window.screen.width = platform.requestedWidth;
    CORE.Window.screen.height = platform.requestedHeight;
    SetupFramebuffer(width, height);

    EGLint displayFormat = 0;
    eglGetConfigAttrib(platform.device, platform.config, EGL_NATIVE_VISUAL_ID, &displayFormat);
    ANativeWindow_setBuffersGeometry(platform.app->window, CORE.Window.render.width, CORE.Window.render.height, displayFormat);

    CORE.Window.render.width = CORE.Window.screen.width;
    CORE.Window.render.height = CORE.Window.screen.height;
    CORE.Window.currentFbo.width = CORE.Window.render.width;
    CORE.Window.currentFbo.height = CORE.Window.render.height;
    SetupViewport(CORE.Window.currentFbo.width, CORE.Window.currentFbo.height);

    CORE.Window.resizedLastFrame = true;
}

// Attach current thread to the Java VM, NULL on failure
static JNIEnv *AttachJni(void)
{
    JNIEnv *env = NULL;
    JavaVM *vm = platform.app->activity->vm;

    if ((*vm)->AttachCurrentThread(vm, &env, NULL) != JNI_OK)
    {
        TRACELOG(LOG_WARNING, "ANDROID: Failed to attach thread to the Java VM");
        return NULL;
    }

    return env;
}

// Detach current thread from the Java VM
static void DetachJni(void)
{
    JavaVM *vm = platform.app->activity->vm;
    (*vm)->DetachCurrentThread(vm);
}

// Check and clear a pending Java exception
static bool JniFailed(JNIEnv *env)
{
    if (!(*env)->ExceptionCheck(env)) return false;

    (*env)->ExceptionClear(env);
    TRACELOG(LOG_WARNING, "ANDROID: Java exception in a platform call");

    return true;
}

// Get default display object (android.view.Display)
static jobject GetDefaultDisplay(JNIEnv *env)
{
    jobject activity = platform.app->activity->clazz;
    jclass activityClass = (*env)->GetObjectClass(env, activity);
    jobject windowManager = (*env)->CallObjectMethod(env, activity, (*env)->GetMethodID(env, activityClass, "getWindowManager", "()Landroid/view/WindowManager;"));
    if (JniFailed(env) || (windowManager == NULL)) return NULL;

    jclass windowManagerClass = (*env)->FindClass(env, "android/view/WindowManager");
    jobject display = (*env)->CallObjectMethod(env, windowManager, (*env)->GetMethodID(env, windowManagerClass, "getDefaultDisplay", "()Landroid/view/Display;"));
    if (JniFailed(env)) return NULL;

    return display;
}

// Get character produced by a key, considering dead keys
// NOTE: Returns 0 if the key does not produce a character or it is a dead key waiting for the next one
static int GetKeyCharacter(int deviceId, int keycode, int metaState)
{
    int codepoint = 0;

    JNIEnv *env = AttachJni();
    if (env == NULL) return codepoint;

    jclass mapClass = (*env)->FindClass(env, "android/view/KeyCharacterMap");
    jobject map = (*env)->CallStaticObjectMethod(env, mapClass, (*env)->GetStaticMethodID(env, mapClass, "load", "(I)Landroid/view/KeyCharacterMap;"), deviceId);

    if (!JniFailed(env) && (map != NULL))
    {
        int value = (*env)->CallIntMethod(env, map, (*env)->GetMethodID(env, mapClass, "get", "(II)I"), keycode, metaState);

        if (!JniFailed(env))
        {
            const int combiningAccent = (int)0x80000000;    // KeyCharacterMap.COMBINING_ACCENT
            const int combiningAccentMask = 0x7fffffff;     // KeyCharacterMap.COMBINING_ACCENT_MASK

            if (value & combiningAccent) platform.deadChar = value & combiningAccentMask;
            else if (value != 0)
            {
                codepoint = value;

                if (platform.deadChar != 0)
                {
                    int combined = (*env)->CallStaticIntMethod(env, mapClass, (*env)->GetStaticMethodID(env, mapClass, "getDeadChar", "(II)I"), platform.deadChar, value);
                    if (!JniFailed(env) && (combined != 0)) codepoint = combined;
                    platform.deadChar = 0;
                }
            }
        }
    }

    DetachJni();

    return codepoint;
}

// Get gamepad index for an input device, connecting it if required
static int GetGamepadIndex(int deviceId, bool connect)
{
    for (int i = 0; i < MAX_GAMEPADS; i++)
    {
        if (platform.gamepadDeviceId[i] == deviceId) return i;
    }

    if (!connect) return -1;

    for (int i = 0; i < MAX_GAMEPADS; i++)
    {
        if (platform.gamepadDeviceId[i] == -1)
        {
            ConnectGamepad(i, deviceId);
            return i;
        }
    }

    return -1;
}

// Register a connected gamepad
static void ConnectGamepad(int gamepad, int deviceId)
{
    platform.gamepadDeviceId[gamepad] = deviceId;
    CORE.Input.Gamepad.ready[gamepad] = true;
    CORE.Input.Gamepad.axisCount[gamepad] = 6;
    CORE.Input.Gamepad.axisState[gamepad][GAMEPAD_AXIS_LEFT_TRIGGER] = -1.0f;
    CORE.Input.Gamepad.axisState[gamepad][GAMEPAD_AXIS_RIGHT_TRIGGER] = -1.0f;
    snprintf(CORE.Input.Gamepad.name[gamepad], MAX_GAMEPAD_NAME_LENGTH, "Gamepad");

    JNIEnv *env = AttachJni();
    if (env == NULL) return;

    jclass deviceClass = (*env)->FindClass(env, "android/view/InputDevice");
    jobject device = (*env)->CallStaticObjectMethod(env, deviceClass, (*env)->GetStaticMethodID(env, deviceClass, "getDevice", "(I)Landroid/view/InputDevice;"), deviceId);

    if (!JniFailed(env) && (device != NULL))
    {
        jstring name = (jstring)(*env)->CallObjectMethod(env, device, (*env)->GetMethodID(env, deviceClass, "getName", "()Ljava/lang/String;"));

        if (!JniFailed(env) && (name != NULL))
        {
            const char *text = (*env)->GetStringUTFChars(env, name, NULL);
            snprintf(CORE.Input.Gamepad.name[gamepad], MAX_GAMEPAD_NAME_LENGTH, "%s", text);
            (*env)->ReleaseStringUTFChars(env, name, text);
        }
    }

    DetachJni();

    TRACELOG(LOG_INFO, "INPUT: Gamepad %i connected: %s", gamepad, CORE.Input.Gamepad.name[gamepad]);
}

// Unregister a disconnected gamepad
static void DisconnectGamepad(int gamepad)
{
    TRACELOG(LOG_INFO, "INPUT: Gamepad %i disconnected: %s", gamepad, CORE.Input.Gamepad.name[gamepad]);

    platform.gamepadDeviceId[gamepad] = -1;
    CORE.Input.Gamepad.ready[gamepad] = false;
    CORE.Input.Gamepad.axisCount[gamepad] = 0;
    memset(CORE.Input.Gamepad.name[gamepad], 0, MAX_GAMEPAD_NAME_LENGTH);
    memset(CORE.Input.Gamepad.currentButtonState[gamepad], 0, sizeof(CORE.Input.Gamepad.currentButtonState[gamepad]));
    memset(CORE.Input.Gamepad.previousButtonState[gamepad], 0, sizeof(CORE.Input.Gamepad.previousButtonState[gamepad]));
    memset(CORE.Input.Gamepad.axisState[gamepad], 0, sizeof(CORE.Input.Gamepad.axisState[gamepad]));
}

// Check connected and disconnected gamepads
// NOTE: Input devices are checked once per second
static void CheckGamepads(void)
{
    double time = GetTime();
    if (time < platform.gamepadCheckTime) return;
    platform.gamepadCheckTime = time + 1.0;

    JNIEnv *env = AttachJni();
    if (env == NULL) return;

    bool connected[MAX_GAMEPADS] = { 0 };
    int newDeviceIds[MAX_GAMEPADS] = { 0 };
    int newCount = 0;

    jclass deviceClass = (*env)->FindClass(env, "android/view/InputDevice");
    jintArray ids = (jintArray)(*env)->CallStaticObjectMethod(env, deviceClass, (*env)->GetStaticMethodID(env, deviceClass, "getDeviceIds", "()[I"));

    if (JniFailed(env) || (ids == NULL))
    {
        DetachJni();
        return;
    }

    int count = (*env)->GetArrayLength(env, ids);
    jint *values = (*env)->GetIntArrayElements(env, ids, NULL);
    jmethodID getDevice = (*env)->GetStaticMethodID(env, deviceClass, "getDevice", "(I)Landroid/view/InputDevice;");
    jmethodID getSources = (*env)->GetMethodID(env, deviceClass, "getSources", "()I");

    for (int i = 0; i < count; i++)
    {
        int index = GetGamepadIndex(values[i], false);
        if (index >= 0)
        {
            connected[index] = true;
            continue;
        }

        jobject device = (*env)->CallStaticObjectMethod(env, deviceClass, getDevice, values[i]);
        if (JniFailed(env) || (device == NULL)) continue;

        int sources = (*env)->CallIntMethod(env, device, getSources);
        (*env)->DeleteLocalRef(env, device);
        if (JniFailed(env)) continue;

        // NOTE: Gamepads report both sources, other devices can report only one of them (keyboards with a joystick source)
        if (FLAG_IS_SET(sources, AINPUT_SOURCE_GAMEPAD) && FLAG_IS_SET(sources, AINPUT_SOURCE_JOYSTICK) && (newCount < MAX_GAMEPADS)) newDeviceIds[newCount++] = values[i];
    }

    (*env)->ReleaseIntArrayElements(env, ids, values, JNI_ABORT);
    DetachJni();

    for (int i = 0; i < MAX_GAMEPADS; i++)
    {
        if ((platform.gamepadDeviceId[i] != -1) && !connected[i]) DisconnectGamepad(i);
    }

    for (int i = 0; i < newCount; i++) GetGamepadIndex(newDeviceIds[i], true);
}

// Register key state change
// NOTE: A key released in the frame it was pressed keeps the press for that frame
static void HandleKey(KeyboardKey key, bool pressed)
{
    if (pressed)
    {
        CORE.Input.Keyboard.currentKeyState[key] = 1;
        platform.keyPressedInFrame[key] = true;
        platform.keyReleasePending[key] = false;

        // Add key to the queue
        if (CORE.Input.Keyboard.keyPressedQueueCount < MAX_KEY_PRESSED_QUEUE)
        {
            CORE.Input.Keyboard.keyPressedQueue[CORE.Input.Keyboard.keyPressedQueueCount] = key;
            CORE.Input.Keyboard.keyPressedQueueCount++;
        }
    }
    else if (platform.keyPressedInFrame[key]) platform.keyReleasePending[key] = true;
    else CORE.Input.Keyboard.currentKeyState[key] = 0;
}

// Register mouse button or touch state change
// NOTE: A button released in the frame it was pressed keeps the press for that frame
static void HandleMouseButton(int button, bool pressed)
{
    if (pressed)
    {
        CORE.Input.Mouse.currentButtonState[button] = 1;
        CORE.Input.Touch.currentTouchState[button] = 1;
        platform.buttonPressedInFrame[button] = true;
        platform.buttonReleasePending[button] = false;
    }
    else if (platform.buttonPressedInFrame[button]) platform.buttonReleasePending[button] = true;
    else
    {
        CORE.Input.Mouse.currentButtonState[button] = 0;
        CORE.Input.Touch.currentTouchState[button] = 0;
    }
}

// Convert display position to screen position
static Vector2 GetScreenPosition(float x, float y)
{
    float widthRatio = (float)(CORE.Window.screen.width + CORE.Window.renderOffset.x)/(float)CORE.Window.display.width;
    float heightRatio = (float)(CORE.Window.screen.height + CORE.Window.renderOffset.y)/(float)CORE.Window.display.height;

    return (Vector2){ x*widthRatio - (float)CORE.Window.renderOffset.x/2, y*heightRatio - (float)CORE.Window.renderOffset.y/2 };
}

// Request an update applied on the main (UI) thread
// NOTE: Never waits for the main thread
static void RequestUiUpdate(unsigned int request)
{
    if (ui.eventFd < 0) return;

    atomic_fetch_or(&ui.requests, request);

    uint64_t value = 1;
    write(ui.eventFd, &value, sizeof(value));
}

// Request system bars visibility from window flags
static void UpdateSystemBars(void)
{
    atomic_store(&ui.immersive, FLAG_IS_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE) || FLAG_IS_SET(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE));
    RequestUiUpdate(UI_REQUEST_SYSTEM_BARS);

    // Window size can change with the system bars
    platform.displayCheckEnd = GetTime() + 2.0;
}

// Apply pending requests on the main (UI) thread
static int UiThreadCallback(int fd, int events, void *data)
{
    uint64_t value = 0;
    read(fd, &value, sizeof(value));

    unsigned int requests = atomic_exchange(&ui.requests, 0);
    struct android_app *app = platform.app;
    if ((requests == 0) || (app == NULL) || (app->activity == NULL)) return 1;

    JNIEnv *env = app->activity->env;
    jobject activity = app->activity->clazz;
    int sdkVersion = app->activity->sdkVersion;

    if ((*env)->PushLocalFrame(env, 32) != JNI_OK) return 1;

    jclass activityClass = (*env)->GetObjectClass(env, activity);
    jobject window = (*env)->CallObjectMethod(env, activity, (*env)->GetMethodID(env, activityClass, "getWindow", "()Landroid/view/Window;"));
    jclass windowClass = (*env)->FindClass(env, "android/view/Window");
    jobject view = (!JniFailed(env) && (window != NULL))? (*env)->CallObjectMethod(env, window, (*env)->GetMethodID(env, windowClass, "getDecorView", "()Landroid/view/View;")) : NULL;
    jclass viewClass = (*env)->FindClass(env, "android/view/View");

    if (!JniFailed(env) && (view != NULL))
    {
        if (FLAG_IS_SET(requests, UI_REQUEST_CURSOR) && (sdkVersion >= 24))
        {
            int cursor = atomic_load(&ui.cursor);
            int iconType = ((cursor >= 0) && (cursor <= MOUSE_CURSOR_NOT_ALLOWED))? mapCursor[cursor] : POINTER_ICON_NULL;

            jclass iconClass = (*env)->FindClass(env, "android/view/PointerIcon");
            jobject icon = (*env)->CallStaticObjectMethod(env, iconClass, (*env)->GetStaticMethodID(env, iconClass, "getSystemIcon", "(Landroid/content/Context;I)Landroid/view/PointerIcon;"), activity, (jint)iconType);
            if (!JniFailed(env)) (*env)->CallVoidMethod(env, view, (*env)->GetMethodID(env, viewClass, "setPointerIcon", "(Landroid/view/PointerIcon;)V"), icon);
            JniFailed(env);
        }

        if (FLAG_IS_SET(requests, UI_REQUEST_CAPTURE) && (sdkVersion >= 26))
        {
            bool capture = atomic_load(&ui.pointerCapture);
            (*env)->CallVoidMethod(env, view, (*env)->GetMethodID(env, viewClass, capture? "requestPointerCapture" : "releasePointerCapture", "()V"));
            JniFailed(env);
        }

        if (FLAG_IS_SET(requests, UI_REQUEST_SOFT_KEYBOARD))
        {
            // NOTE: Input method serves the decor view, the native content view never gets the focus
            jstring serviceName = (*env)->NewStringUTF(env, "input_method");
            jobject manager = (*env)->CallObjectMethod(env, activity, (*env)->GetMethodID(env, activityClass, "getSystemService", "(Ljava/lang/String;)Ljava/lang/Object;"), serviceName);
            jclass managerClass = (*env)->FindClass(env, "android/view/inputmethod/InputMethodManager");

            if (!JniFailed(env) && (manager != NULL))
            {
                if (atomic_load(&ui.softKeyboard)) (*env)->CallBooleanMethod(env, manager, (*env)->GetMethodID(env, managerClass, "showSoftInput", "(Landroid/view/View;I)Z"), view, 0);
                else
                {
                    jobject token = (*env)->CallObjectMethod(env, view, (*env)->GetMethodID(env, viewClass, "getWindowToken", "()Landroid/os/IBinder;"));
                    if (!JniFailed(env)) (*env)->CallBooleanMethod(env, manager, (*env)->GetMethodID(env, managerClass, "hideSoftInputFromWindow", "(Landroid/os/IBinder;I)Z"), token, 0);
                }

                JniFailed(env);
            }
        }

        if (FLAG_IS_SET(requests, UI_REQUEST_SYSTEM_BARS))
        {
            bool immersive = atomic_load(&ui.immersive);

            if (sdkVersion >= 30)
            {
                // Bars shown by a swipe are hidden again after a while (BEHAVIOR_SHOW_TRANSIENT_BARS_BY_SWIPE)
                jobject controller = (*env)->CallObjectMethod(env, window, (*env)->GetMethodID(env, windowClass, "getInsetsController", "()Landroid/view/WindowInsetsController;"));
                jclass typeClass = (*env)->FindClass(env, "android/view/WindowInsets$Type");
                jint types = JniFailed(env)? 0 : (*env)->CallStaticIntMethod(env, typeClass, (*env)->GetStaticMethodID(env, typeClass, "systemBars", "()I"));

                if (!JniFailed(env) && (controller != NULL))
                {
                    jclass controllerClass = (*env)->FindClass(env, "android/view/WindowInsetsController");
                    (*env)->CallVoidMethod(env, controller, (*env)->GetMethodID(env, controllerClass, "setSystemBarsBehavior", "(I)V"), (jint)SYSTEM_BARS_TRANSIENT);
                    (*env)->CallVoidMethod(env, controller, (*env)->GetMethodID(env, controllerClass, immersive? "hide" : "show", "(I)V"), types);
                    JniFailed(env);
                }
            }
            else
            {
                jint visibility = immersive? SYSTEM_UI_IMMERSIVE : 0;
                (*env)->CallVoidMethod(env, view, (*env)->GetMethodID(env, viewClass, "setSystemUiVisibility", "(I)V"), visibility);
                JniFailed(env);
            }
        }
    }

    (*env)->PopLocalFrame(env, NULL);

    return 1;
}

// Replacement for fopen(), used as linker wrap entry point (-Wl,--wrap=fopen)
// REF: https://developer.android.com/ndk/reference/group/asset
__attribute__((visibility("default"))) FILE *__wrap_fopen(const char *fileName, const char *mode)
{
    FILE *file = NULL;

    // NOTE: AAsset provides access to read-only asset, write operations use regular fopen
    if (mode[0] == 'w')
    {
        file = __real_fopen(TextFormat("%s/%s", platform.app->activity->internalDataPath, fileName), mode);
        if (file == NULL) file = __real_fopen(fileName, mode);
    }
    else
    {
        AAsset *asset = AAssetManager_open(platform.app->activity->assetManager, fileName, AASSET_MODE_UNKNOWN);

        if (asset != NULL)
        {
            // Get pointer to file in the assets
            file = funopen(asset, android_read, android_write, android_seek, android_close);
        }
        else
        {
            // Do a regular open if file is not found in the assets
            file = __real_fopen(TextFormat("%s/%s", platform.app->activity->internalDataPath, fileName), mode);
            if (file == NULL) file = __real_fopen(fileName, mode);
        }
    }

    return file;
}

static int android_read(void *cookie, char *data, int dataSize)
{
    return AAsset_read((AAsset *)cookie, data, dataSize);
}

static int android_write(void *cookie, const char *data, int dataSize)
{
    TRACELOG(LOG_WARNING, "ANDROID: Failed to provide write access to APK");

    return EACCES;
}

static fpos_t android_seek(void *cookie, fpos_t offset, int whence)
{
    return AAsset_seek((AAsset *)cookie, offset, whence);
}

static int android_close(void *cookie)
{
    AAsset_close((AAsset *)cookie);
    return 0;
}

// EOF
