/**********************************************************************************************
*
*   rcore_desktop_x11 - Functions to manage window, graphics device and inputs
*
*   PLATFORM: DESKTOP: X11
*       - Linux (X11 servers, Xwayland)
*       - FreeBSD, OpenBSD (X11 servers)
*
*   LIMITATIONS:
*       - Fullscreen does not change the video mode, the window covers the monitor at its current mode
*       - One scale factor for all monitors (Xft.dpi), X11 has no per-monitor scale
*       - Window state requests depend on the window manager (EWMH), tiling window managers decide size and position
*       - Gamepads are only supported on Linux
*
*   POSSIBLE IMPROVEMENTS:
*       - Primary selection (middle-click paste)
*       - Input method preedit drawing (on-the-spot)
*       - Smooth scrolling while the cursor is disabled (XInput2 pointer grab)
*
*   CONFIGURATION:
*       #define RCORE_PLATFORM_CUSTOM_FLAG
*           Custom flag for rcore on target platform -not used-
*
*   DEPENDENCIES:
*       - libX11, libXext (loaded at runtime): X11 protocol, shared memory images, input shape and sync counters
*       - libXrandr (loaded at runtime): Monitors
*       - libXi (loaded at runtime): Smooth scrolling, touch and raw mouse motion (XInput2)
*       - libXcursor (loaded at runtime): Cursor themes
*       - libxkbcommon (loaded at runtime): Key symbols to characters
*       - EGL (loaded at runtime): OpenGL context creation
*       - minigamepad: Gamepad support (evdev, Linux)
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

// NOTE: X11 defines a Font type, renamed to avoid the conflict with raylib Font
#define Font X11Font
#include <X11/Xlib.h>               // X11 client library
#include <X11/Xutil.h>              // Window manager hints, visuals
#include <X11/Xatom.h>              // Predefined atoms
#include <X11/Xresource.h>          // Resource manager database (Xft.dpi)
#include <X11/XKBlib.h>             // Keyboard extension
#include <X11/keysym.h>             // Key symbols
#include <X11/cursorfont.h>         // Core cursor shapes
#include <X11/extensions/Xrandr.h>  // Monitors
#include <X11/extensions/XInput2.h> // Smooth scrolling, touch and raw mouse motion
#include <X11/extensions/shape.h>   // Input shape (mouse passthrough)
#include <X11/extensions/sync.h>    // Sync counters (_NET_WM_SYNC_REQUEST)
#if defined(GRAPHICS_API_OPENGL_SOFTWARE)
    #include <X11/extensions/XShm.h> // Shared memory images
    #include <sys/ipc.h>            // Required for: IPC_PRIVATE, IPC_RMID
    #include <sys/shm.h>            // Required for: shmget(), shmat(), shmdt(), shmctl()
#endif
#undef Font

Cursor XcursorLibraryLoadCursor(Display *display, const char *file); // Cursor themes (libXcursor)

#include <xkbcommon/xkbcommon.h>    // Required for: xkb_keysym_to_utf8(), xkb_keysym_to_utf32()

#if !defined(GRAPHICS_API_OPENGL_SOFTWARE)
    // NOTE: glad.h embeds khrplatform.h without KHRONOS_APIENTRY, required by EGL headers
    #if !defined(KHRONOS_APIENTRY)
        #define KHRONOS_APIENTRY
    #endif
    #include <EGL/egl.h>            // EGL library
#endif

#include <time.h>                   // Required for: clock_gettime()
#include <stdlib.h>                 // Required for: getenv(), unsetenv(), strtol(), getprogname()
#include <ctype.h>                  // Required for: isxdigit()
#include <limits.h>                 // Required for: LONG_MAX
#include <locale.h>                 // Required for: setlocale()
#include <poll.h>                   // Required for: poll()
#include <unistd.h>                 // Required for: readlink(), getpid()
#include <errno.h>                  // Required for: errno, EINTR
#include <dlfcn.h>                  // Required for: dlopen(), dlsym(), dlclose()
#include <math.h>                   // Required for: roundf()

// System library functions, loaded at runtime
#define LIBX11_FUNCTIONS(F) \
    F(XAllocSizeHints) F(XChangeProperty) F(XCheckIfEvent) F(XCloseDisplay) F(XCloseIM) F(XConvertSelection) \
    F(XCreateBitmapFromData) F(XCreateColormap) F(XCreateFontCursor) F(XCreateGC) F(XCreateIC) \
    F(XCreateImage) F(XCreatePixmapCursor) F(XCreateWindow) F(XDefineCursor) F(XDeleteProperty) F(XDestroyIC) \
    F(XDestroyWindow) F(XDisplayKeycodes) F(XDisplayName) F(XEventsQueued) F(XExtendedMaxRequestSize) \
    F(XFilterEvent) F(XFlush) F(XFree) F(XFreeColormap) F(XFreeCursor) F(XFreeEventData) F(XFreeGC) \
    F(XFreePixmap) F(XGetAtomName) F(XGetEventData) F(XGetICValues) F(XGetIMValues) F(XGetInputFocus) \
    F(XGetKeyboardMapping) F(XGetSelectionOwner) F(XGetVisualInfo) F(XGetWindowProperty) F(XGrabPointer) \
    F(XIconifyWindow) F(XInitThreads) F(XInternAtoms) F(XkbFreeKeyboard) F(XkbFreeNames) F(XkbGetMap) \
    F(XkbGetNames) F(XkbGetState) F(XkbKeycodeToKeysym) F(XkbKeysymToModifiers) F(XkbQueryExtension) \
    F(XkbRefreshKeyboardMapping) F(XkbSelectEventDetails) F(XkbSelectEvents) F(XkbSetDetectableAutoRepeat) \
    F(XLookupString) F(XMapWindow) F(XMatchVisualInfo) F(XMaxRequestSize) F(XMoveResizeWindow) F(XMoveWindow) \
    F(XNextEvent) F(XOpenDisplay) F(XOpenIM) F(XPeekEvent) F(XPending) F(XPutImage) F(XQueryExtension) \
    F(XQueryKeymap) F(XQueryPointer) F(XRaiseWindow) F(XRefreshKeyboardMapping) \
    F(XRegisterIMInstantiateCallback) F(XResizeWindow) F(XrmDestroyDatabase) F(XrmGetResource) \
    F(XrmGetStringDatabase) F(XrmInitialize) F(XSelectInput) F(XSendEvent) F(XSetClassHint) \
    F(XSetErrorHandler) F(XSetICFocus) F(XSetIMValues) F(XSetInputFocus) F(XSetIOErrorHandler) \
    F(XSetLocaleModifiers) F(XSetSelectionOwner) F(XSetWMClientMachine) F(XSetWMHints) F(XSetWMNormalHints) \
    F(XSetWMProtocols) F(XStringListToTextProperty) F(XSupportsLocale) F(XSync) F(XTranslateCoordinates) \
    F(XUngrabPointer) F(XUnregisterIMInstantiateCallback) F(XUnsetICFocus) F(Xutf8LookupString) \
    F(XWarpPointer) F(XWithdrawWindow)
#define LIBXEXT_FUNCTIONS(F) \
    F(XShapeCombineMask) F(XShapeCombineRectangles) F(XShapeQueryExtension) F(XShapeQueryVersion) \
    F(XSyncCreateCounter) F(XSyncDestroyCounter) F(XSyncInitialize) F(XSyncIntsToValue) F(XSyncIntToValue) \
    F(XSyncQueryExtension) F(XSyncSetCounter)
#define LIBXRANDR_FUNCTIONS(F) \
    F(XRRFreeCrtcInfo) F(XRRFreeMonitors) F(XRRFreeOutputInfo) F(XRRFreeScreenResources) F(XRRGetCrtcInfo) \
    F(XRRGetMonitors) F(XRRGetOutputInfo) F(XRRGetScreenResourcesCurrent) F(XRRQueryExtension) \
    F(XRRQueryVersion) F(XRRSelectInput) F(XRRUpdateConfiguration)
#define LIBXI_FUNCTIONS(F) \
    F(XIFreeDeviceInfo) F(XIQueryDevice) F(XIQueryVersion) F(XISelectEvents)
#define LIBXCURSOR_FUNCTIONS(F) \
    F(XcursorLibraryLoadCursor)
#define LIBXKBCOMMON_FUNCTIONS(F) \
    F(xkb_keysym_to_utf32) F(xkb_keysym_to_utf8)
#if defined(GRAPHICS_API_OPENGL_SOFTWARE)
    #define LIBXSHM_FUNCTIONS(F) \
        F(XShmAttach) F(XShmCreateImage) F(XShmDetach) F(XShmPutImage) F(XShmQueryExtension)
    #define LIBEGL_FUNCTIONS(F)
#else
    #define LIBXSHM_FUNCTIONS(F)
    #define LIBEGL_FUNCTIONS(F) \
        F(eglBindAPI) F(eglChooseConfig) F(eglCreateContext) F(eglCreateWindowSurface) F(eglDestroyContext) \
        F(eglDestroySurface) F(eglGetConfigAttrib) F(eglGetDisplay) F(eglGetProcAddress) F(eglInitialize) \
        F(eglMakeCurrent) F(eglQueryString) F(eglSwapBuffers) F(eglSwapInterval) F(eglTerminate)
#endif
#define LIBRARY_FUNCTION(name) __typeof__(&name) p##name;

// Gamepad support (evdev), Linux only
#if defined(__linux__)
    #define MG_MAX_GAMEPADS MAX_GAMEPADS
    #if !defined(MG_IMPLEMENTATION)
        #define MG_IMPLEMENTATION
        #include "../external/RGFW/deps/minigamepad.h"
    #endif
#endif

// NOTE: linux/input.h (included by minigamepad) defines KEY_<X> macros that collide with raylib KEY_<X> values,
// macros are undefined to use the raylib values, also in rcore.c
#undef KEY_APOSTROPHE
#undef KEY_COMMA
#undef KEY_MINUS
#undef KEY_SLASH
#undef KEY_SEMICOLON
#undef KEY_EQUAL
#undef KEY_A
#undef KEY_B
#undef KEY_C
#undef KEY_D
#undef KEY_E
#undef KEY_F
#undef KEY_G
#undef KEY_H
#undef KEY_I
#undef KEY_J
#undef KEY_K
#undef KEY_L
#undef KEY_M
#undef KEY_N
#undef KEY_O
#undef KEY_P
#undef KEY_Q
#undef KEY_R
#undef KEY_S
#undef KEY_T
#undef KEY_U
#undef KEY_V
#undef KEY_W
#undef KEY_X
#undef KEY_Y
#undef KEY_Z
#undef KEY_BACKSLASH
#undef KEY_GRAVE
#undef KEY_SPACE
#undef KEY_ENTER
#undef KEY_TAB
#undef KEY_BACKSPACE
#undef KEY_INSERT
#undef KEY_DELETE
#undef KEY_RIGHT
#undef KEY_LEFT
#undef KEY_DOWN
#undef KEY_UP
#undef KEY_HOME
#undef KEY_END
#undef KEY_PAUSE
#undef KEY_F1
#undef KEY_F2
#undef KEY_F3
#undef KEY_F4
#undef KEY_F5
#undef KEY_F6
#undef KEY_F7
#undef KEY_F8
#undef KEY_F9
#undef KEY_F10
#undef KEY_F11
#undef KEY_F12
#undef KEY_BACK
#undef KEY_MENU

//----------------------------------------------------------------------------------
// Defines and Macros
//----------------------------------------------------------------------------------
#ifndef EGL_OPENGL_ES3_BIT
    #define EGL_OPENGL_ES3_BIT      0x40
#endif
#ifndef EGL_PLATFORM_X11_EXT
    #define EGL_PLATFORM_X11_EXT    0x31D5
#endif
#ifndef EGL_PLATFORM_X11_SCREEN_EXT
    #define EGL_PLATFORM_X11_SCREEN_EXT 0x31D6
#endif

#define MAX_MONITORS                16      // Maximum number of monitors tracked
#define MAX_WINDOW_ICONS            8       // Maximum number of window icon sizes
#define MAX_INPUT_DEVICES           32      // Maximum number of pointer devices tracked (XInput2)
#define MAX_SCROLL_AXES             4       // Maximum number of scroll axes per pointer device
#define MAX_INCR_TRANSFERS          4       // Maximum number of simultaneous incremental selection transfers
#define KEYMAP_SIZE                 256     // Size of the X11 keycodes to raylib keys map
#define MOUSE_CURSOR_COUNT          11      // Number of MouseCursor shapes

#define DPI_DEFAULT                 96.0f   // Display DPI at scale factor 1.0
#define SELECTION_TIMEOUT           1000    // Maximum time waiting for a selection owner or clipboard manager (milliseconds)
#define INCR_TRANSFER_TIMEOUT       10.0    // Maximum time waiting for the requestor of an incremental selection transfer (seconds)
#define WINDOW_MAP_TIMEOUT          200     // Maximum time waiting for the window manager to map the window (milliseconds)
#define KEY_REPEAT_EPSILON          20      // Maximum time between release and press of a repeated key (milliseconds)
#define FULLSCREEN_REQUEST_TIMEOUT   1.0     // Maximum time waiting for the window manager to apply a fullscreen request (seconds)

#define XDND_VERSION                5       // Drag and drop protocol version supported

#define NET_WM_STATE_REMOVE         0       // _NET_WM_STATE client message action: remove state
#define NET_WM_STATE_ADD            1       // _NET_WM_STATE client message action: add state

#define MWM_HINTS_DECORATIONS       2       // _MOTIF_WM_HINTS flags: decorations field is valid
#define MWM_DECOR_ALL               1       // _MOTIF_WM_HINTS decorations: all decorations

//----------------------------------------------------------------------------------
// Types and Structures Definition
//----------------------------------------------------------------------------------

// Monitor info
typedef struct {
    char name[128];                 // Monitor name
    int x;                          // Position x (pixels)
    int y;                          // Position y (pixels)
    int width;                      // Width (pixels)
    int height;                     // Height (pixels)
    int physicalWidth;              // Physical width in millimetres
    int physicalHeight;             // Physical height in millimetres
    float refreshRate;              // Refresh rate in Hz
    int index;                      // Index in the RandR monitors list (Xinerama screen index)
} MonitorInfo;

// Scroll axis of a pointer device (XInput2 scroll valuator)
typedef struct {
    int number;                     // Valuator number
    bool horizontal;                // Horizontal scroll axis
    double increment;               // Valuator change for one scroll step
    double value;                   // Last valuator value
    bool valid;                     // Last valuator value is known
} ScrollAxis;

// Motion axis of a pointer device (XInput2 valuators 0 and 1)
typedef struct {
    bool relative;                  // Axis reports relative motion, absolute positions otherwise
    double min;                     // Minimum value (absolute axis)
    double max;                     // Maximum value (absolute axis)
    double value;                   // Last value (absolute axis)
    bool valid;                     // Last value is known
} MotionAxis;

// Pointer device (XInput2)
typedef struct {
    int id;                         // Device id
    MotionAxis axes[2];             // Motion axes: x and y
    ScrollAxis scrollAxes[MAX_SCROLL_AXES]; // Scroll axes
    int scrollAxisCount;            // Scroll axes count
} InputDevice;

// Incremental selection transfer to another client (INCR)
typedef struct {
    Window requestor;               // Window receiving the data, None if not used
    Atom property;                  // Property the data is written to
    Atom type;                      // Data type
    char *data;                     // Data copy
    int size;                       // Data size in bytes
    int offset;                     // Data already sent in bytes
    double time;                    // Time of the last chunk sent, transfers of unresponsive requestors are released
} IncrTransfer;

// Event predicate arguments
typedef struct {
    int type;                       // Event type
    Window window;                  // Event window
    Atom atom;                      // Selection (SelectionNotify) or property (PropertyNotify)
} EventFilter;

// X11 atoms
// NOTE: Structure is filled as an array of atoms, in the order of atomNames
typedef struct {
    Atom WM_PROTOCOLS;
    Atom WM_DELETE_WINDOW;
    Atom WM_STATE;
    Atom UTF8_STRING;
    Atom RESOURCE_MANAGER;
    Atom _MOTIF_WM_HINTS;
    Atom _NET_SUPPORTED;
    Atom _NET_SUPPORTING_WM_CHECK;
    Atom _NET_WM_NAME;
    Atom _NET_WM_ICON_NAME;
    Atom _NET_WM_ICON;
    Atom _NET_WM_PID;
    Atom _NET_WM_PING;
    Atom _NET_WM_SYNC_REQUEST;
    Atom _NET_WM_SYNC_REQUEST_COUNTER;
    Atom _NET_WM_USER_TIME;
    Atom _NET_WM_WINDOW_TYPE;
    Atom _NET_WM_WINDOW_TYPE_NORMAL;
    Atom _NET_WM_STATE;
    Atom _NET_WM_STATE_ABOVE;
    Atom _NET_WM_STATE_FULLSCREEN;
    Atom _NET_WM_STATE_MAXIMIZED_VERT;
    Atom _NET_WM_STATE_MAXIMIZED_HORZ;
    Atom _NET_WM_STATE_HIDDEN;
    Atom _NET_WM_FULLSCREEN_MONITORS;
    Atom _NET_WM_WINDOW_OPACITY;
    Atom _NET_ACTIVE_WINDOW;
    Atom _NET_STARTUP_ID;
    Atom _NET_STARTUP_INFO_BEGIN;
    Atom _NET_STARTUP_INFO;
    Atom CLIPBOARD;
    Atom CLIPBOARD_MANAGER;
    Atom SAVE_TARGETS;
    Atom TARGETS;
    Atom MULTIPLE;
    Atom ATOM_PAIR;
    Atom INCR;
    Atom NULL_;
    Atom TEXT;
    Atom TEXT_PLAIN_UTF8;
    Atom TEXT_URI_LIST;
    Atom IMAGE_PNG;
    Atom IMAGE_BMP;
    Atom IMAGE_JPEG;
    Atom RAYLIB_SELECTION;
    Atom RAYLIB_DROP;
    Atom RAYLIB_TIMESTAMP;
    Atom XdndAware;
    Atom XdndEnter;
    Atom XdndPosition;
    Atom XdndStatus;
    Atom XdndLeave;
    Atom XdndDrop;
    Atom XdndFinished;
    Atom XdndSelection;
    Atom XdndTypeList;
    Atom XdndActionCopy;
} Atoms;

typedef struct {
    // X11 connection
    Display *display;               // X11 display connection
    int screen;                     // Default screen
    Window root;                    // Root window
    Atoms atoms;                    // Interned atoms
    XErrorHandler previousErrorHandler; // Error handler replaced on initialization
    XIOErrorHandler previousIOErrorHandler; // Connection error handler replaced on initialization
    bool ioErrorHandlerSet;         // Connection error handler replaced
    int errorCode;                  // Last X11 error code received
    bool connectionLost;            // Connection to the X server lost

    // Extensions
    bool xkbAvailable;              // Keyboard extension available
    int xkbEventBase;               // Keyboard extension first event
    bool xkbDetectableRepeat;       // Key repeats are sent as press events only
    bool randrAvailable;            // RandR extension available (1.2 or later)
    bool randrMonitors;             // RandR monitors available (1.5 or later)
    int randrEventBase;             // RandR extension first event
    bool xiAvailable;               // XInput2 extension available
    int xiOpcode;                   // XInput2 extension major opcode
    int xiMinor;                    // XInput2 minor version supported by the server
    bool shapeAvailable;            // Shape extension available with input shapes (1.1 or later)
    bool syncAvailable;             // Sync extension available

    // Window objects
    Window window;                  // Window
    Window parent;                  // Window parent, root or the window manager frame
    Visual *visual;                 // Window visual
    int depth;                      // Window depth
    Colormap colormap;              // Window colormap
    long eventMask;                 // Window core events selected
    XSyncCounter syncCounter;       // Counter for _NET_WM_SYNC_REQUEST
    XSyncValue syncValue;           // Counter value requested by the window manager
    bool syncPending;               // Counter value has to be set after the next frame
    char appId[256];                // Window application id (executable name)

#if !defined(GRAPHICS_API_OPENGL_SOFTWARE)
    // EGL data
    EGLDisplay eglDisplay;          // Native display device (physical screen connection)
    EGLContext eglContext;          // Graphic context, mode in which drawing can be done
    EGLSurface eglSurface;          // Surface to draw on, framebuffers (connected to context)
    EGLConfig eglConfig;            // Graphic config
    bool eglPlatformDisplay;        // Display obtained with the X11 platform extension
#else
    // Software renderer
    XImage *image;                  // Framebuffer image presented on the window
    XShmSegmentInfo shmInfo;        // Shared memory segment of the image
    bool shmAvailable;              // Shared memory images available
    bool shmAttached;               // Image uses the shared memory segment
    GC gc;                          // Graphics context to present the image
#endif

    // Window state
    int width;                      // Window content width (logical pixels)
    int height;                     // Window content height (logical pixels)
    int bufferWidth;                // Window content width (pixels)
    int bufferHeight;               // Window content height (pixels)
    int x;                          // Window content position x on the root window (pixels)
    int y;                          // Window content position y on the root window (pixels)
    float scale;                    // Scale factor (Xft.dpi)
    int windowedWidth;              // Last floating window content width (logical pixels)
    int windowedHeight;             // Last floating window content height (logical pixels)
    int windowedX;                  // Last floating window content position x (pixels)
    int windowedY;                  // Last floating window content position y (pixels)
    bool mapped;                    // Window is mapped
    bool withdrawn;                 // Window is not managed by the window manager (never shown or hidden)
    bool fullscreen;                // Window is fullscreen
    bool maximized;                 // Window is maximized
    bool focused;                   // Window has the input focus
    bool fullscreenRequested;       // Last fullscreen state requested to the window manager
    bool fullscreenPending;         // Last fullscreen request not yet confirmed by the window manager
    double fullscreenRequestTime;   // Time of the last fullscreen request
    bool fullscreenManual;          // Fullscreen applied by resizing the window (window manager without EWMH)
    bool decorationRequested;       // Decorations requested (FLAG_WINDOW_UNDECORATED not set)
    bool closeRequested;            // Window close requested by user/window manager
    bool eventReceived;             // Window or input event received (ends event waiting)
    bool initialized;               // Platform initialization finished (rlgl ready on next events)
    Window wmCheckWindow;           // Window of the running window manager (_NET_SUPPORTING_WM_CHECK)
    Atom *wmSupported;              // Hints supported by the window manager (_NET_SUPPORTED)
    int wmSupportedCount;           // Hints supported by the window manager count
    Time lastUserTime;              // Last user interaction time (key or button press)
    unsigned long resizeSerial;     // Serial of the last resize request, size events generated before it are outdated

    // Monitors
    MonitorInfo monitors[MAX_MONITORS]; // Connected monitors
    int monitorCount;               // Connected monitors count

    // Pointer
    InputDevice devices[MAX_INPUT_DEVICES]; // Pointer devices (XInput2)
    int deviceCount;                // Pointer devices count
    Cursor cursors[MOUSE_CURSOR_COUNT]; // Cursors loaded per MouseCursor
    Cursor hiddenCursor;            // Transparent cursor
    bool pointerGrabbed;            // Pointer is grabbed and confined to the window (disabled cursor)
    double pointerX;                // Last pointer position x on the window (pixels)
    double pointerY;                // Last pointer position y on the window (pixels)
    double warpX;                   // Last position x the pointer was moved to (pixels)
    double warpY;                   // Last position y the pointer was moved to (pixels)
    bool keyPressedInFrame[MAX_KEYBOARD_KEYS]; // Keys pressed in the current frame
    bool keyReleasePending[MAX_KEYBOARD_KEYS]; // Keys released in the frame they were pressed, release registered on next frame
    bool buttonPressedInFrame[MAX_MOUSE_BUTTONS]; // Mouse buttons pressed in the current frame
    bool buttonReleasePending[MAX_MOUSE_BUTTONS]; // Mouse buttons released in the frame they were pressed, release registered on next frame
    int touchIds[MAX_TOUCH_POINTS]; // Touch point ids per slot, -1 if slot not used
    Vector2 touchPositions[MAX_TOUCH_POINTS]; // Touch point positions per slot

    // Keyboard
    short keymap[KEYMAP_SIZE];      // Map of X11 keycodes to raylib KeyboardKey values
    XIM inputMethod;                // Input method connection
    XIC inputContext;               // Input context of the window
    Time keyPressTimes[KEYMAP_SIZE]; // Last press time per keycode (input methods send key events twice)
    unsigned int xkbGroup;          // Current keyboard layout group
    unsigned int numLockMask;       // Modifier mask of NUM lock
    bool capsLocked;                // CAPS lock is active
    bool numLocked;                 // NUM lock is active
    bool capsLockDown;              // CAPS lock key is held
    bool numLockDown;               // NUM lock key is held

    // Clipboard and drag and drop
    char *clipboardText;            // Text owned by this client as clipboard selection
    char *clipboardReceived;        // Text received from the clipboard selection
    bool clipboardOwned;            // Clipboard selection is owned by this client
    Time clipboardTime;             // Time the clipboard selection was acquired
    IncrTransfer incrTransfers[MAX_INCR_TRANSFERS]; // Incremental selection transfers in progress
    Window dragSource;              // Drag and drop source window
    int dragVersion;                // Drag and drop protocol version of the source
    bool dragAccepted;              // Drag and drop source offers a file list

#if defined(__linux__)
    // Gamepads
    mg_gamepads minigamepad;        // Gamepads (evdev)
#endif

    // System libraries, loaded at runtime
    void *libX11;                   // X11 client library
    void *libXext;                  // X11 extensions library
    void *libXrandr;                // RandR extension library
    void *libXi;                    // XInput2 extension library
    void *libXcursor;               // Cursor themes library
    void *libxkbcommon;             // Keyboard keymap library
    void *libEGL;                   // EGL library
    struct {
        LIBX11_FUNCTIONS(LIBRARY_FUNCTION)
        LIBXEXT_FUNCTIONS(LIBRARY_FUNCTION)
        LIBXSHM_FUNCTIONS(LIBRARY_FUNCTION)
        LIBXRANDR_FUNCTIONS(LIBRARY_FUNCTION)
        LIBXI_FUNCTIONS(LIBRARY_FUNCTION)
        LIBXCURSOR_FUNCTIONS(LIBRARY_FUNCTION)
        LIBXKBCOMMON_FUNCTIONS(LIBRARY_FUNCTION)
        LIBEGL_FUNCTIONS(LIBRARY_FUNCTION)
    } lib;
} PlatformData;

// System library functions, loaded at runtime
#define XAllocSizeHints                  platform.lib.pXAllocSizeHints
#define XChangeProperty                  platform.lib.pXChangeProperty
#define XCheckIfEvent                    platform.lib.pXCheckIfEvent
#define XCloseDisplay                    platform.lib.pXCloseDisplay
#define XCloseIM                         platform.lib.pXCloseIM
#define XConvertSelection                platform.lib.pXConvertSelection
#define XCreateBitmapFromData            platform.lib.pXCreateBitmapFromData
#define XCreateColormap                  platform.lib.pXCreateColormap
#define XCreateFontCursor                platform.lib.pXCreateFontCursor
#define XCreateGC                        platform.lib.pXCreateGC
#define XCreateIC                        platform.lib.pXCreateIC
#define XCreateImage                     platform.lib.pXCreateImage
#define XCreatePixmapCursor              platform.lib.pXCreatePixmapCursor
#define XCreateWindow                    platform.lib.pXCreateWindow
#define XDefineCursor                    platform.lib.pXDefineCursor
#define XDeleteProperty                  platform.lib.pXDeleteProperty
#define XDestroyIC                       platform.lib.pXDestroyIC
#define XDestroyWindow                   platform.lib.pXDestroyWindow
#define XDisplayKeycodes                 platform.lib.pXDisplayKeycodes
#define XDisplayName                     platform.lib.pXDisplayName
#define XEventsQueued                    platform.lib.pXEventsQueued
#define XExtendedMaxRequestSize          platform.lib.pXExtendedMaxRequestSize
#define XFilterEvent                     platform.lib.pXFilterEvent
#define XFlush                           platform.lib.pXFlush
#define XFree                            platform.lib.pXFree
#define XFreeColormap                    platform.lib.pXFreeColormap
#define XFreeCursor                      platform.lib.pXFreeCursor
#define XFreeEventData                   platform.lib.pXFreeEventData
#define XFreeGC                          platform.lib.pXFreeGC
#define XFreePixmap                      platform.lib.pXFreePixmap
#define XGetAtomName                     platform.lib.pXGetAtomName
#define XGetEventData                    platform.lib.pXGetEventData
#define XGetICValues                     platform.lib.pXGetICValues
#define XGetIMValues                     platform.lib.pXGetIMValues
#define XGetInputFocus                   platform.lib.pXGetInputFocus
#define XGetKeyboardMapping              platform.lib.pXGetKeyboardMapping
#define XGetSelectionOwner               platform.lib.pXGetSelectionOwner
#define XGetVisualInfo                   platform.lib.pXGetVisualInfo
#define XGetWindowProperty               platform.lib.pXGetWindowProperty
#define XGrabPointer                     platform.lib.pXGrabPointer
#define XIconifyWindow                   platform.lib.pXIconifyWindow
#define XInitThreads                     platform.lib.pXInitThreads
#define XInternAtoms                     platform.lib.pXInternAtoms
#define XkbFreeKeyboard                  platform.lib.pXkbFreeKeyboard
#define XkbFreeNames                     platform.lib.pXkbFreeNames
#define XkbGetMap                        platform.lib.pXkbGetMap
#define XkbGetNames                      platform.lib.pXkbGetNames
#define XkbGetState                      platform.lib.pXkbGetState
#define XkbKeycodeToKeysym               platform.lib.pXkbKeycodeToKeysym
#define XkbKeysymToModifiers             platform.lib.pXkbKeysymToModifiers
#define XkbQueryExtension                platform.lib.pXkbQueryExtension
#define XkbRefreshKeyboardMapping        platform.lib.pXkbRefreshKeyboardMapping
#define XkbSelectEventDetails            platform.lib.pXkbSelectEventDetails
#define XkbSelectEvents                  platform.lib.pXkbSelectEvents
#define XkbSetDetectableAutoRepeat       platform.lib.pXkbSetDetectableAutoRepeat
#define XLookupString                    platform.lib.pXLookupString
#define XMapWindow                       platform.lib.pXMapWindow
#define XMatchVisualInfo                 platform.lib.pXMatchVisualInfo
#define XMaxRequestSize                  platform.lib.pXMaxRequestSize
#define XMoveResizeWindow                platform.lib.pXMoveResizeWindow
#define XMoveWindow                      platform.lib.pXMoveWindow
#define XNextEvent                       platform.lib.pXNextEvent
#define XOpenDisplay                     platform.lib.pXOpenDisplay
#define XOpenIM                          platform.lib.pXOpenIM
#define XPeekEvent                       platform.lib.pXPeekEvent
#define XPending                         platform.lib.pXPending
#define XPutImage                        platform.lib.pXPutImage
#define XQueryExtension                  platform.lib.pXQueryExtension
#define XQueryKeymap                     platform.lib.pXQueryKeymap
#define XQueryPointer                    platform.lib.pXQueryPointer
#define XRaiseWindow                     platform.lib.pXRaiseWindow
#define XRefreshKeyboardMapping          platform.lib.pXRefreshKeyboardMapping
#define XRegisterIMInstantiateCallback   platform.lib.pXRegisterIMInstantiateCallback
#define XResizeWindow                    platform.lib.pXResizeWindow
#define XrmDestroyDatabase               platform.lib.pXrmDestroyDatabase
#define XrmGetResource                   platform.lib.pXrmGetResource
#define XrmGetStringDatabase             platform.lib.pXrmGetStringDatabase
#define XrmInitialize                    platform.lib.pXrmInitialize
#define XSelectInput                     platform.lib.pXSelectInput
#define XSendEvent                       platform.lib.pXSendEvent
#define XSetClassHint                    platform.lib.pXSetClassHint
#define XSetErrorHandler                 platform.lib.pXSetErrorHandler
#define XSetICFocus                      platform.lib.pXSetICFocus
#define XSetIMValues                     platform.lib.pXSetIMValues
#define XSetInputFocus                   platform.lib.pXSetInputFocus
#define XSetIOErrorHandler               platform.lib.pXSetIOErrorHandler
#define XSetLocaleModifiers              platform.lib.pXSetLocaleModifiers
#define XSetSelectionOwner               platform.lib.pXSetSelectionOwner
#define XSetWMClientMachine              platform.lib.pXSetWMClientMachine
#define XSetWMHints                      platform.lib.pXSetWMHints
#define XSetWMNormalHints                platform.lib.pXSetWMNormalHints
#define XSetWMProtocols                  platform.lib.pXSetWMProtocols
#define XStringListToTextProperty        platform.lib.pXStringListToTextProperty
#define XSupportsLocale                  platform.lib.pXSupportsLocale
#define XSync                            platform.lib.pXSync
#define XTranslateCoordinates            platform.lib.pXTranslateCoordinates
#define XUngrabPointer                   platform.lib.pXUngrabPointer
#define XUnregisterIMInstantiateCallback platform.lib.pXUnregisterIMInstantiateCallback
#define XUnsetICFocus                    platform.lib.pXUnsetICFocus
#define Xutf8LookupString                platform.lib.pXutf8LookupString
#define XWarpPointer                     platform.lib.pXWarpPointer
#define XWithdrawWindow                  platform.lib.pXWithdrawWindow
#define XShapeCombineMask                platform.lib.pXShapeCombineMask
#define XShapeCombineRectangles          platform.lib.pXShapeCombineRectangles
#define XShapeQueryExtension             platform.lib.pXShapeQueryExtension
#define XShapeQueryVersion               platform.lib.pXShapeQueryVersion
#define XSyncCreateCounter               platform.lib.pXSyncCreateCounter
#define XSyncDestroyCounter              platform.lib.pXSyncDestroyCounter
#define XSyncInitialize                  platform.lib.pXSyncInitialize
#define XSyncIntsToValue                 platform.lib.pXSyncIntsToValue
#define XSyncIntToValue                  platform.lib.pXSyncIntToValue
#define XSyncQueryExtension              platform.lib.pXSyncQueryExtension
#define XSyncSetCounter                  platform.lib.pXSyncSetCounter
#define XRRFreeCrtcInfo                  platform.lib.pXRRFreeCrtcInfo
#define XRRFreeMonitors                  platform.lib.pXRRFreeMonitors
#define XRRFreeOutputInfo                platform.lib.pXRRFreeOutputInfo
#define XRRFreeScreenResources           platform.lib.pXRRFreeScreenResources
#define XRRGetCrtcInfo                   platform.lib.pXRRGetCrtcInfo
#define XRRGetMonitors                   platform.lib.pXRRGetMonitors
#define XRRGetOutputInfo                 platform.lib.pXRRGetOutputInfo
#define XRRGetScreenResourcesCurrent     platform.lib.pXRRGetScreenResourcesCurrent
#define XRRQueryExtension                platform.lib.pXRRQueryExtension
#define XRRQueryVersion                  platform.lib.pXRRQueryVersion
#define XRRSelectInput                   platform.lib.pXRRSelectInput
#define XRRUpdateConfiguration           platform.lib.pXRRUpdateConfiguration
#define XIFreeDeviceInfo                 platform.lib.pXIFreeDeviceInfo
#define XIQueryDevice                    platform.lib.pXIQueryDevice
#define XIQueryVersion                   platform.lib.pXIQueryVersion
#define XISelectEvents                   platform.lib.pXISelectEvents
#define XcursorLibraryLoadCursor         platform.lib.pXcursorLibraryLoadCursor
#define xkb_keysym_to_utf32              platform.lib.pxkb_keysym_to_utf32
#define xkb_keysym_to_utf8               platform.lib.pxkb_keysym_to_utf8
#if defined(GRAPHICS_API_OPENGL_SOFTWARE)
    #define XShmAttach                       platform.lib.pXShmAttach
    #define XShmCreateImage                  platform.lib.pXShmCreateImage
    #define XShmDetach                       platform.lib.pXShmDetach
    #define XShmPutImage                     platform.lib.pXShmPutImage
    #define XShmQueryExtension               platform.lib.pXShmQueryExtension
#else
    #define eglBindAPI                       platform.lib.peglBindAPI
    #define eglChooseConfig                  platform.lib.peglChooseConfig
    #define eglCreateContext                 platform.lib.peglCreateContext
    #define eglCreateWindowSurface           platform.lib.peglCreateWindowSurface
    #define eglDestroyContext                platform.lib.peglDestroyContext
    #define eglDestroySurface                platform.lib.peglDestroySurface
    #define eglGetConfigAttrib               platform.lib.peglGetConfigAttrib
    #define eglGetDisplay                    platform.lib.peglGetDisplay
    #define eglGetProcAddress                platform.lib.peglGetProcAddress
    #define eglInitialize                    platform.lib.peglInitialize
    #define eglMakeCurrent                   platform.lib.peglMakeCurrent
    #define eglQueryString                   platform.lib.peglQueryString
    #define eglSwapBuffers                   platform.lib.peglSwapBuffers
    #define eglSwapInterval                  platform.lib.peglSwapInterval
    #define eglTerminate                     platform.lib.peglTerminate
#endif

// X11 key name (XKB) to raylib KeyboardKey value
typedef struct {
    short key;
    const char *name;
} KeyNameMapping;

//----------------------------------------------------------------------------------
// Global Variables Definition
//----------------------------------------------------------------------------------
extern CoreData CORE;                   // Global CORE state context

static PlatformData platform = { 0 };   // Platform specific data

// Map XKB key names to raylib KeyboardKey values
// NOTE: Key names identify physical keys, independently of the keyboard layout
static const KeyNameMapping keyNameMap[] = {
    { KEY_GRAVE, "TLDE" }, { KEY_ONE, "AE01" }, { KEY_TWO, "AE02" }, { KEY_THREE, "AE03" }, { KEY_FOUR, "AE04" },
    { KEY_FIVE, "AE05" }, { KEY_SIX, "AE06" }, { KEY_SEVEN, "AE07" }, { KEY_EIGHT, "AE08" }, { KEY_NINE, "AE09" },
    { KEY_ZERO, "AE10" }, { KEY_MINUS, "AE11" }, { KEY_EQUAL, "AE12" }, { KEY_Q, "AD01" }, { KEY_W, "AD02" },
    { KEY_E, "AD03" }, { KEY_R, "AD04" }, { KEY_T, "AD05" }, { KEY_Y, "AD06" }, { KEY_U, "AD07" },
    { KEY_I, "AD08" }, { KEY_O, "AD09" }, { KEY_P, "AD10" }, { KEY_LEFT_BRACKET, "AD11" }, { KEY_RIGHT_BRACKET, "AD12" },
    { KEY_A, "AC01" }, { KEY_S, "AC02" }, { KEY_D, "AC03" }, { KEY_F, "AC04" }, { KEY_G, "AC05" },
    { KEY_H, "AC06" }, { KEY_J, "AC07" }, { KEY_K, "AC08" }, { KEY_L, "AC09" }, { KEY_SEMICOLON, "AC10" },
    { KEY_APOSTROPHE, "AC11" }, { KEY_Z, "AB01" }, { KEY_X, "AB02" }, { KEY_C, "AB03" }, { KEY_V, "AB04" },
    { KEY_B, "AB05" }, { KEY_N, "AB06" }, { KEY_M, "AB07" }, { KEY_COMMA, "AB08" }, { KEY_PERIOD, "AB09" },
    { KEY_SLASH, "AB10" }, { KEY_BACKSLASH, "BKSL" }, { 162, "LSGT" }, { KEY_SPACE, "SPCE" }, { KEY_ESCAPE, "ESC" },
    { KEY_ENTER, "RTRN" }, { KEY_TAB, "TAB" }, { KEY_BACKSPACE, "BKSP" }, { KEY_INSERT, "INS" }, { KEY_DELETE, "DELE" },
    { KEY_RIGHT, "RGHT" }, { KEY_LEFT, "LEFT" }, { KEY_DOWN, "DOWN" }, { KEY_UP, "UP" }, { KEY_PAGE_UP, "PGUP" },
    { KEY_PAGE_DOWN, "PGDN" }, { KEY_HOME, "HOME" }, { KEY_END, "END" }, { KEY_CAPS_LOCK, "CAPS" }, { KEY_SCROLL_LOCK, "SCLK" },
    { KEY_NUM_LOCK, "NMLK" }, { KEY_PRINT_SCREEN, "PRSC" }, { KEY_PAUSE, "PAUS" }, { KEY_F1, "FK01" }, { KEY_F2, "FK02" },
    { KEY_F3, "FK03" }, { KEY_F4, "FK04" }, { KEY_F5, "FK05" }, { KEY_F6, "FK06" }, { KEY_F7, "FK07" },
    { KEY_F8, "FK08" }, { KEY_F9, "FK09" }, { KEY_F10, "FK10" }, { KEY_F11, "FK11" }, { KEY_F12, "FK12" },
    { 302, "FK13" }, { 303, "FK14" }, { 304, "FK15" }, { 305, "FK16" }, { 306, "FK17" },
    { 307, "FK18" }, { 308, "FK19" }, { 309, "FK20" }, { 310, "FK21" }, { 311, "FK22" },
    { 312, "FK23" }, { 313, "FK24" }, { 314, "FK25" }, { KEY_KP_0, "KP0" }, { KEY_KP_1, "KP1" },
    { KEY_KP_2, "KP2" }, { KEY_KP_3, "KP3" }, { KEY_KP_4, "KP4" }, { KEY_KP_5, "KP5" }, { KEY_KP_6, "KP6" },
    { KEY_KP_7, "KP7" }, { KEY_KP_8, "KP8" }, { KEY_KP_9, "KP9" }, { KEY_KP_DECIMAL, "KPDL" }, { KEY_KP_DIVIDE, "KPDV" },
    { KEY_KP_MULTIPLY, "KPMU" }, { KEY_KP_SUBTRACT, "KPSU" }, { KEY_KP_ADD, "KPAD" }, { KEY_KP_ENTER, "KPEN" }, { KEY_KP_EQUAL, "KPEQ" },
    { KEY_LEFT_SHIFT, "LFSH" }, { KEY_LEFT_CONTROL, "LCTL" }, { KEY_LEFT_ALT, "LALT" }, { KEY_LEFT_SUPER, "LWIN" }, { KEY_RIGHT_SHIFT, "RTSH" },
    { KEY_RIGHT_CONTROL, "RCTL" }, { KEY_RIGHT_ALT, "RALT" }, { KEY_RIGHT_ALT, "LVL3" }, { KEY_RIGHT_ALT, "MDSW" }, { KEY_RIGHT_SUPER, "RWIN" },
    { KEY_KB_MENU, "COMP" }, { KEY_KB_MENU, "MENU" }
};

// X11 atom names, in the order of the Atoms structure
static const char *atomNames[] = {
    "WM_PROTOCOLS", "WM_DELETE_WINDOW", "WM_STATE", "UTF8_STRING", "RESOURCE_MANAGER", "_MOTIF_WM_HINTS",
    "_NET_SUPPORTED", "_NET_SUPPORTING_WM_CHECK", "_NET_WM_NAME", "_NET_WM_ICON_NAME", "_NET_WM_ICON", "_NET_WM_PID",
    "_NET_WM_PING", "_NET_WM_SYNC_REQUEST", "_NET_WM_SYNC_REQUEST_COUNTER", "_NET_WM_USER_TIME",
    "_NET_WM_WINDOW_TYPE", "_NET_WM_WINDOW_TYPE_NORMAL", "_NET_WM_STATE", "_NET_WM_STATE_ABOVE",
    "_NET_WM_STATE_FULLSCREEN", "_NET_WM_STATE_MAXIMIZED_VERT", "_NET_WM_STATE_MAXIMIZED_HORZ",
    "_NET_WM_STATE_HIDDEN", "_NET_WM_FULLSCREEN_MONITORS", "_NET_WM_WINDOW_OPACITY", "_NET_ACTIVE_WINDOW",
    "_NET_STARTUP_ID", "_NET_STARTUP_INFO_BEGIN", "_NET_STARTUP_INFO", "CLIPBOARD", "CLIPBOARD_MANAGER",
    "SAVE_TARGETS", "TARGETS", "MULTIPLE", "ATOM_PAIR", "INCR", "NULL", "TEXT", "text/plain;charset=utf-8",
    "text/uri-list", "image/png", "image/bmp", "image/jpeg", "RAYLIB_SELECTION", "RAYLIB_DROP", "RAYLIB_TIMESTAMP",
    "XdndAware", "XdndEnter", "XdndPosition", "XdndStatus", "XdndLeave", "XdndDrop", "XdndFinished", "XdndSelection",
    "XdndTypeList", "XdndActionCopy"
};

// Cursor theme names per MouseCursor: CSS name and legacy X11 name
static const char *cursorNames[MOUSE_CURSOR_COUNT][2] = {
    { "default", "left_ptr" },
    { "default", "left_ptr" },
    { "text", "xterm" },
    { "crosshair", "cross" },
    { "pointer", "hand2" },
    { "ew-resize", "sb_h_double_arrow" },
    { "ns-resize", "sb_v_double_arrow" },
    { "nwse-resize", "bottom_right_corner" },
    { "nesw-resize", "bottom_left_corner" },
    { "all-scroll", "fleur" },
    { "not-allowed", "crossed_circle" }
};

// Core cursor shapes per MouseCursor, used when the cursor theme does not provide the cursor
static const unsigned int cursorShapes[MOUSE_CURSOR_COUNT] = {
    XC_left_ptr, XC_left_ptr, XC_xterm, XC_crosshair, XC_hand2, XC_sb_h_double_arrow,
    XC_sb_v_double_arrow, XC_bottom_right_corner, XC_bottom_left_corner, XC_fleur, XC_X_cursor
};

#if defined(__linux__)
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
#endif

//----------------------------------------------------------------------------------
// Module Internal Functions Declaration
//----------------------------------------------------------------------------------
int InitPlatform(void);          // Initialize platform (graphics, inputs and more)
void ClosePlatform(void);        // Close platform

static bool LoadLibraries(void);                                // Load system libraries and their functions
static void CloseLibraries(void);                               // Close system libraries
static int ErrorHandler(Display *display, XErrorEvent *event);  // X11 errors handler, errors are not fatal
static int IOErrorHandler(Display *display);                    // X11 connection errors handler
static void IOErrorExitHandler(Display *display, void *data);   // X11 connection errors exit handler, the process is not terminated
static void InitAtoms(void);                                    // Intern required atoms
static void InitExtensions(void);                               // Query X11 extensions
static bool InitGraphicsContext(void);                          // Initialize EGL display and config, choose the window visual
static bool InitGraphicsSurface(void);                          // Initialize EGL context and window surface
static bool CreateWindow(void);                                 // Create window and set its properties
static void ShowWindow(bool wait);                              // Map window, optionally wait for the window manager
static void CompleteStartupNotification(void);                  // Notify the launcher that the application started

static void ProcessEvents(int timeout);                         // Read and process X11 events, timeout in milliseconds (-1 blocks)
static bool WaitForEvents(int timeout);                         // Wait for data on the X11 connection, timeout in milliseconds (-1 blocks)
static Bool MatchEvent(Display *display, XEvent *event, XPointer data); // Event predicate: selection notifications and property changes for a window and atom
static bool WaitForEvent(int type, Window window, Atom atom, XEvent *event, int timeout); // Wait for an event of a given type, window and atom, other events are kept queued
static Time GetServerTime(void);                                // Get current X server time
static void ProcessEvent(XEvent *event);                        // Process X11 event
static void ProcessClientMessage(XClientMessageEvent *event);   // Process window manager and drag and drop messages
static void ProcessPropertyChange(XPropertyEvent *event);       // Process window and root properties changes

static unsigned long GetWindowProperty(Window window, Atom property, Atom type, unsigned char **value); // Get window property, returns items count
static void SendWindowMessage(Atom type, long data0, long data1, long data2, long data3, long data4); // Send client message to the root window (request to the window manager)
static void UpdateWMSupported(void);                            // Update hints supported by the window manager
static bool IsWMHintSupported(Atom hint);                       // Check hint supported by the window manager
static void SetWindowStateHint(Atom state1, Atom state2, bool enabled); // Request window state change (_NET_WM_STATE)
static void UpdateWindowStateFlags(void);                       // Update window flags from window manager states
static void UpdateNormalHints(void);                            // Update window size limits (WM_NORMAL_HINTS)
static void UpdateDecorations(void);                            // Update window decorations (_MOTIF_WM_HINTS)
static void UpdateWindowTitle(const char *title);               // Update window title properties
static void UpdateWindowScale(void);                            // Update scale factor from the resource manager
static void UpdateWindowBuffers(void);                          // Apply window size and scale to the framebuffer
static void ApplyWindowSize(int width, int height);             // Apply window size reported by the server (pixels)
static void RequestFullscreen(bool fullscreen, int monitor);    // Request fullscreen state to the window manager
static void KeepWindowOnMonitors(void);                         // Move floating window to the main monitor if it is out of all monitors
static float GetOutputRefreshRate(XRRScreenResources *resources, RROutput output); // Get refresh rate of the mode of a RandR output, in Hz
static void UpdateMonitors(void);                               // Update connected monitors
static int GetMonitorAt(int x, int y, int width, int height);   // Get monitor with the largest overlap with a rectangle

static void InitKeyboard(void);                                 // Initialize keyboard extension, keymap and input method
static void UpdateKeymap(void);                                 // Update map of X11 keycodes to raylib keys
static int KeysymToKey(KeySym keysym);                          // Get raylib key for a key symbol
static void InputMethodInstantiateCallback(Display *display, XPointer clientData, XPointer callData); // Input method available
static void CreateInputContext(void);                           // Create window input context
static void InputMethodDestroyCallback(XIM inputMethod, XPointer clientData, XPointer callData); // Input method connection closed
static void InputContextDestroyCallback(XIC inputContext, XPointer clientData, XPointer callData); // Input context destroyed
static void HandleKey(XKeyEvent *event, bool pressed, bool filtered); // Register key state change
static void ReleaseKey(int key);                                // Register key release
static void HandleKeyChar(XKeyEvent *event);                    // Register character input for a key press
static void HandleFocus(bool focused);                          // Register window focus change
static void SyncInputState(void);                               // Release keys and mouse buttons not held anymore

static void InitInputDevices(void);                             // Initialize XInput2 and select pointer events
static void UpdateInputDevices(void);                           // Update pointer devices info (scroll axes, motion mode)
static InputDevice *GetInputDevice(int id);                     // Get pointer device info
static void ResetMotionAxes(void);                              // Forget last values of absolute motion axes
static void ProcessGenericEvent(XGenericEventCookie *cookie);   // Process XInput2 event
static void HandleMouseButton(int button, bool pressed);        // Register mouse button state change
static void HandleMouseMove(double x, double y);                // Register mouse position change (logical pixels)
static void HandlePointerButton(int button, bool pressed, Time time); // Register X11 pointer button, wheel buttons included
static void HandlePointerMotion(double x, double y);            // Register X11 pointer position (pixels)
static void HandleTouch(int touchAction, int id, double x, double y); // Register touch event
static void UpdateCursor(void);                                 // Update window cursor
static Cursor GetCursor(int cursor);                            // Get cursor for a MouseCursor shape
static void LockPointer(bool lock);                             // Grab or release pointer (disabled cursor)
static void SelectRawMotion(bool enabled);                      // Select raw motion events on the root window
static void WarpPointer(int x, int y);                          // Move pointer to window position (pixels)

static void HandleSelectionRequest(XSelectionRequestEvent *request); // Send clipboard content to another client
static bool WriteSelection(Window requestor, Atom property, Atom target); // Write clipboard content to a requestor property
static void ContinueIncrTransfer(IncrTransfer *transfer);       // Send next chunk of an incremental selection transfer
static void ReleaseIncrTransfer(IncrTransfer *transfer);        // Release incremental selection transfer
static unsigned char *ReadSelection(Atom selection, Atom target, int *size, Atom *type); // Read selection content, returned data is allocated and zero-terminated
static unsigned char *ReadSelectionProperty(Atom property, int *size, Atom *type); // Read and delete a selection property, incremental transfers included
static bool SelectionHasTarget(Atom selection, Atom target);    // Check selection owner supported target
static Bool MatchSelectionEvent(Display *display, XEvent *event, XPointer data); // Event predicate: selection events
static void SaveClipboard(void);                                // Hand clipboard content over to the clipboard manager
static void HandleDragEnter(XClientMessageEvent *event);        // Drag and drop: source entered the window
static void HandleDragPosition(XClientMessageEvent *event);     // Drag and drop: source moved over the window
static void HandleDragDrop(XClientMessageEvent *event);         // Drag and drop: data dropped
static void HandleDragData(XSelectionEvent *event);             // Drag and drop: file list received from the source
static void FinishDrag(bool accepted);                          // Drag and drop: notify the source that the drop is finished
static void ProcessDroppedUris(char *uris);                     // Register dropped files from a text/uri-list

#if defined(GRAPHICS_API_OPENGL_SOFTWARE)
static bool CreateSoftwareImage(int width, int height);         // Create software renderer framebuffer image
static void DestroySoftwareImage(void);                         // Destroy software renderer framebuffer image
#endif

//----------------------------------------------------------------------------------
// Module Functions Declaration
//----------------------------------------------------------------------------------
// NOTE: Functions declaration is provided by raylib.h

//----------------------------------------------------------------------------------
// Module Functions Definition: Window and Graphics Device
//----------------------------------------------------------------------------------

// Check if application should close
// NOTE: By default, if KEY_ESCAPE pressed or window close icon clicked
bool WindowShouldClose(void)
{
    if (CORE.Window.ready) return CORE.Window.shouldClose;
    else return true;
}

// Toggle fullscreen mode
void ToggleFullscreen(void)
{
    if (!FLAG_IS_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE))
    {
        // Leave borderless windowed mode before entering fullscreen
        if (FLAG_IS_SET(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE)) ToggleBorderlessWindowed();

        // Store previous screen data (in case exiting fullscreen)
        CORE.Window.previousPosition = CORE.Window.position;
        CORE.Window.previousScreen = CORE.Window.screen;

        int monitor = GetCurrentMonitor();
        CORE.Window.display.width = GetMonitorWidth(monitor);
        CORE.Window.display.height = GetMonitorHeight(monitor);

        FLAG_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE);
        FLAG_SET(CORE.Window.flags, FLAG_WINDOW_UNDECORATED);

        RequestFullscreen(true, monitor);
    }
    else
    {
        // Restore previous window position
        CORE.Window.position = CORE.Window.previousPosition;
        FLAG_CLEAR(CORE.Window.flags, FLAG_FULLSCREEN_MODE);
        if (platform.decorationRequested) FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_UNDECORATED);

        RequestFullscreen(false, 0);
    }

    // NOTE: New window size is applied when the window manager resizes the window
}

// Toggle borderless windowed mode
// NOTE: It is a fullscreen window, the video mode is never changed
void ToggleBorderlessWindowed(void)
{
    // Leave fullscreen before attempting to set borderless windowed mode
    if (FLAG_IS_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE)) ToggleFullscreen();

    if (!FLAG_IS_SET(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE))
    {
        // Store screen position and size
        CORE.Window.previousPosition = CORE.Window.position;
        CORE.Window.previousScreen = CORE.Window.screen;

        FLAG_SET(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE);
        FLAG_SET(CORE.Window.flags, FLAG_WINDOW_UNDECORATED);

        RequestFullscreen(true, GetCurrentMonitor());
    }
    else
    {
        // Restore previous screen values
        CORE.Window.position = CORE.Window.previousPosition;
        FLAG_CLEAR(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE);
        if (platform.decorationRequested) FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_UNDECORATED);

        RequestFullscreen(false, 0);
    }
}

// Set window state: maximized, if resizable
void MaximizeWindow(void)
{
    if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_RESIZABLE))
    {
        if (!IsWMHintSupported(platform.atoms._NET_WM_STATE_MAXIMIZED_VERT))
        {
            TRACELOG(LOG_WARNING, "X11: Window maximization not supported by the window manager");
            return;
        }

        // NOTE: Maximized state flag is updated when the window manager changes the window state, the request can be refused
        SetWindowStateHint(platform.atoms._NET_WM_STATE_MAXIMIZED_VERT, platform.atoms._NET_WM_STATE_MAXIMIZED_HORZ, true);
        if (platform.withdrawn) FLAG_SET(CORE.Window.flags, FLAG_WINDOW_MAXIMIZED);
        XFlush(platform.display);
    }
}

// Set window state: minimized
void MinimizeWindow(void)
{
    // NOTE: Minimized state flag is updated when the window manager changes the window state
    XIconifyWindow(platform.display, platform.window, platform.screen);
    XFlush(platform.display);
}

// Restore window from being minimized/maximized
void RestoreWindow(void)
{
    if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_MINIMIZED))
    {
        XMapWindow(platform.display, platform.window);
    }
    else if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_MAXIMIZED))
    {
        SetWindowStateHint(platform.atoms._NET_WM_STATE_MAXIMIZED_VERT, platform.atoms._NET_WM_STATE_MAXIMIZED_HORZ, false);
    }

    XFlush(platform.display);
}

// Set window configuration state using flags
void SetWindowState(unsigned int flags)
{
    if (!CORE.Window.ready) TRACELOG(LOG_WARNING, "WINDOW: SetWindowState does nothing before window initialization, Use \"SetConfigFlags\" instead");

    // Check previous state and requested state to apply required changes
    // NOTE: In most cases the functions already change the flags internally

    // State change: FLAG_VSYNC_HINT
    if ((FLAG_IS_SET(CORE.Window.flags, FLAG_VSYNC_HINT) != FLAG_IS_SET(flags, FLAG_VSYNC_HINT)) && FLAG_IS_SET(flags, FLAG_VSYNC_HINT))
    {
#if !defined(GRAPHICS_API_OPENGL_SOFTWARE)
        eglSwapInterval(platform.eglDisplay, 1);
#endif
        FLAG_SET(CORE.Window.flags, FLAG_VSYNC_HINT);
    }

    // State change: FLAG_BORDERLESS_WINDOWED_MODE
    // NOTE: This must be handled before FLAG_FULLSCREEN_MODE because ToggleBorderlessWindowed() needs to get some fullscreen values if fullscreen is running
    if ((FLAG_IS_SET(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE) != FLAG_IS_SET(flags, FLAG_BORDERLESS_WINDOWED_MODE)) && FLAG_IS_SET(flags, FLAG_BORDERLESS_WINDOWED_MODE))
    {
        ToggleBorderlessWindowed();     // NOTE: Window state flag updated inside function
    }

    // State change: FLAG_FULLSCREEN_MODE
    if ((FLAG_IS_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE) != FLAG_IS_SET(flags, FLAG_FULLSCREEN_MODE)) && FLAG_IS_SET(flags, FLAG_FULLSCREEN_MODE))
    {
        ToggleFullscreen();     // NOTE: Window state flag updated inside function
    }

    // State change: FLAG_WINDOW_RESIZABLE
    if ((FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_RESIZABLE) != FLAG_IS_SET(flags, FLAG_WINDOW_RESIZABLE)) && FLAG_IS_SET(flags, FLAG_WINDOW_RESIZABLE))
    {
        FLAG_SET(CORE.Window.flags, FLAG_WINDOW_RESIZABLE);
        UpdateNormalHints();
    }

    // State change: FLAG_WINDOW_UNDECORATED
    // NOTE: Flag is also set while fullscreen, requested decorations are tracked separately
    if (platform.decorationRequested && FLAG_IS_SET(flags, FLAG_WINDOW_UNDECORATED))
    {
        FLAG_SET(CORE.Window.flags, FLAG_WINDOW_UNDECORATED);
        platform.decorationRequested = false;
        UpdateDecorations();
    }

    // State change: FLAG_WINDOW_HIDDEN
    if ((FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_HIDDEN) != FLAG_IS_SET(flags, FLAG_WINDOW_HIDDEN)) && FLAG_IS_SET(flags, FLAG_WINDOW_HIDDEN))
    {
        FLAG_SET(CORE.Window.flags, FLAG_WINDOW_HIDDEN);
        LockPointer(false);
        XWithdrawWindow(platform.display, platform.window, platform.screen);
        platform.withdrawn = true;
        platform.mapped = false;
    }

    // State change: FLAG_WINDOW_MINIMIZED
    if ((FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_MINIMIZED) != FLAG_IS_SET(flags, FLAG_WINDOW_MINIMIZED)) && FLAG_IS_SET(flags, FLAG_WINDOW_MINIMIZED))
    {
        MinimizeWindow();       // NOTE: Window state flag updated when the window manager changes the window state
    }

    // State change: FLAG_WINDOW_MAXIMIZED
    if ((FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_MAXIMIZED) != FLAG_IS_SET(flags, FLAG_WINDOW_MAXIMIZED)) && FLAG_IS_SET(flags, FLAG_WINDOW_MAXIMIZED))
    {
        MaximizeWindow();       // NOTE: Window state flag updated when the window manager changes the window state
    }

    // State change: FLAG_WINDOW_UNFOCUSED
    if ((FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_UNFOCUSED) != FLAG_IS_SET(flags, FLAG_WINDOW_UNFOCUSED)) && FLAG_IS_SET(flags, FLAG_WINDOW_UNFOCUSED))
    {
        TRACELOG(LOG_WARNING, "SetWindowState(FLAG_WINDOW_UNFOCUSED) not available on target platform");
    }

    // State change: FLAG_WINDOW_TOPMOST
    if ((FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_TOPMOST) != FLAG_IS_SET(flags, FLAG_WINDOW_TOPMOST)) && FLAG_IS_SET(flags, FLAG_WINDOW_TOPMOST))
    {
        SetWindowStateHint(platform.atoms._NET_WM_STATE_ABOVE, None, true);
        FLAG_SET(CORE.Window.flags, FLAG_WINDOW_TOPMOST);
    }

    // State change: FLAG_WINDOW_ALWAYS_RUN
    if ((FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_ALWAYS_RUN) != FLAG_IS_SET(flags, FLAG_WINDOW_ALWAYS_RUN)) && FLAG_IS_SET(flags, FLAG_WINDOW_ALWAYS_RUN))
    {
        FLAG_SET(CORE.Window.flags, FLAG_WINDOW_ALWAYS_RUN);
    }

    // The following states can not be changed after window creation

    // State change: FLAG_WINDOW_TRANSPARENT
    if ((FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_TRANSPARENT) != FLAG_IS_SET(flags, FLAG_WINDOW_TRANSPARENT)) && FLAG_IS_SET(flags, FLAG_WINDOW_TRANSPARENT))
    {
        TRACELOG(LOG_WARNING, "WINDOW: Framebuffer transparency can only be configured before window initialization");
    }

    // State change: FLAG_WINDOW_HIGHDPI
    if ((FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_HIGHDPI) != FLAG_IS_SET(flags, FLAG_WINDOW_HIGHDPI)) && FLAG_IS_SET(flags, FLAG_WINDOW_HIGHDPI))
    {
        TRACELOG(LOG_WARNING, "WINDOW: High DPI can only be configured before window initialization");
    }

    // State change: FLAG_WINDOW_MOUSE_PASSTHROUGH
    if ((FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_MOUSE_PASSTHROUGH) != FLAG_IS_SET(flags, FLAG_WINDOW_MOUSE_PASSTHROUGH)) && FLAG_IS_SET(flags, FLAG_WINDOW_MOUSE_PASSTHROUGH))
    {
        if (platform.shapeAvailable)
        {
            XShapeCombineRectangles(platform.display, platform.window, ShapeInput, 0, 0, NULL, 0, ShapeSet, Unsorted);
            FLAG_SET(CORE.Window.flags, FLAG_WINDOW_MOUSE_PASSTHROUGH);
        }
        else TRACELOG(LOG_WARNING, "X11: Mouse passthrough not supported by the X server");
    }

    // State change: FLAG_MSAA_4X_HINT
    if ((FLAG_IS_SET(CORE.Window.flags, FLAG_MSAA_4X_HINT) != FLAG_IS_SET(flags, FLAG_MSAA_4X_HINT)) && FLAG_IS_SET(flags, FLAG_MSAA_4X_HINT))
    {
        TRACELOG(LOG_WARNING, "WINDOW: MSAA can only be configured before window initialization");
    }

    // State change: FLAG_INTERLACED_HINT
    if ((FLAG_IS_SET(CORE.Window.flags, FLAG_INTERLACED_HINT) != FLAG_IS_SET(flags, FLAG_INTERLACED_HINT)) && FLAG_IS_SET(flags, FLAG_INTERLACED_HINT))
    {
        TRACELOG(LOG_WARNING, "WINDOW: Interlaced mode can only be configured before window initialization");
    }

    XFlush(platform.display);
}

// Clear window configuration state flags
void ClearWindowState(unsigned int flags)
{
    // Check previous state and requested state to apply required changes
    // NOTE: In most cases the functions already change the flags internally

    // State change: FLAG_VSYNC_HINT
    if ((FLAG_IS_SET(CORE.Window.flags, FLAG_VSYNC_HINT)) && (FLAG_IS_SET(flags, FLAG_VSYNC_HINT)))
    {
#if !defined(GRAPHICS_API_OPENGL_SOFTWARE)
        eglSwapInterval(platform.eglDisplay, 0);
#endif
        FLAG_CLEAR(CORE.Window.flags, FLAG_VSYNC_HINT);
    }

    // State change: FLAG_BORDERLESS_WINDOWED_MODE
    // NOTE: This must be handled before FLAG_FULLSCREEN_MODE because ToggleBorderlessWindowed() needs to get some fullscreen values if fullscreen is running
    if ((FLAG_IS_SET(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE)) && (FLAG_IS_SET(flags, FLAG_BORDERLESS_WINDOWED_MODE)))
    {
        ToggleBorderlessWindowed(); // NOTE: Window state flag updated inside function
    }

    // State change: FLAG_FULLSCREEN_MODE
    if ((FLAG_IS_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE)) && (FLAG_IS_SET(flags, FLAG_FULLSCREEN_MODE)))
    {
        ToggleFullscreen(); // NOTE: Window state flag updated inside function
    }

    // State change: FLAG_WINDOW_RESIZABLE
    if ((FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_RESIZABLE)) && (FLAG_IS_SET(flags, FLAG_WINDOW_RESIZABLE)))
    {
        FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_RESIZABLE);
        UpdateNormalHints();
    }

    // State change: FLAG_WINDOW_HIDDEN
    if ((FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_HIDDEN)) && (FLAG_IS_SET(flags, FLAG_WINDOW_HIDDEN)))
    {
        FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_HIDDEN);
        ShowWindow(false);
    }

    // State change: FLAG_WINDOW_MINIMIZED
    if ((FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_MINIMIZED)) && (FLAG_IS_SET(flags, FLAG_WINDOW_MINIMIZED)))
    {
        RestoreWindow();       // NOTE: Window state flag updated when the window manager changes the window state
    }

    // State change: FLAG_WINDOW_MAXIMIZED
    if ((FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_MAXIMIZED)) && (FLAG_IS_SET(flags, FLAG_WINDOW_MAXIMIZED)))
    {
        RestoreWindow();       // NOTE: Window state flag updated when the window manager changes the window state
    }

    // State change: FLAG_WINDOW_UNDECORATED
    // NOTE: Flag stays set while fullscreen, decorations are shown when leaving it
    if (!platform.decorationRequested && FLAG_IS_SET(flags, FLAG_WINDOW_UNDECORATED))
    {
        if (!platform.fullscreen) FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_UNDECORATED);
        platform.decorationRequested = true;
        UpdateDecorations();
    }

    // State change: FLAG_WINDOW_UNFOCUSED
    if ((FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_UNFOCUSED)) && (FLAG_IS_SET(flags, FLAG_WINDOW_UNFOCUSED)))
    {
        SetWindowFocused();
    }

    // State change: FLAG_WINDOW_TOPMOST
    if ((FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_TOPMOST)) && (FLAG_IS_SET(flags, FLAG_WINDOW_TOPMOST)))
    {
        SetWindowStateHint(platform.atoms._NET_WM_STATE_ABOVE, None, false);
        FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_TOPMOST);
    }

    // State change: FLAG_WINDOW_ALWAYS_RUN
    if ((FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_ALWAYS_RUN)) && (FLAG_IS_SET(flags, FLAG_WINDOW_ALWAYS_RUN)))
    {
        FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_ALWAYS_RUN);
    }

    // The following states can not be changed after window creation

    // State change: FLAG_WINDOW_TRANSPARENT
    if ((FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_TRANSPARENT)) && (FLAG_IS_SET(flags, FLAG_WINDOW_TRANSPARENT)))
    {
        TRACELOG(LOG_WARNING, "WINDOW: Framebuffer transparency can only be configured before window initialization");
    }

    // State change: FLAG_WINDOW_HIGHDPI
    if ((FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_HIGHDPI)) && (FLAG_IS_SET(flags, FLAG_WINDOW_HIGHDPI)))
    {
        TRACELOG(LOG_WARNING, "WINDOW: High DPI can only be configured before window initialization");
    }

    // State change: FLAG_WINDOW_MOUSE_PASSTHROUGH
    if ((FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_MOUSE_PASSTHROUGH)) && (FLAG_IS_SET(flags, FLAG_WINDOW_MOUSE_PASSTHROUGH)))
    {
        if (platform.shapeAvailable) XShapeCombineMask(platform.display, platform.window, ShapeInput, 0, 0, None, ShapeSet);
        FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_MOUSE_PASSTHROUGH);
    }

    // State change: FLAG_MSAA_4X_HINT
    if ((FLAG_IS_SET(CORE.Window.flags, FLAG_MSAA_4X_HINT)) && (FLAG_IS_SET(flags, FLAG_MSAA_4X_HINT)))
    {
        TRACELOG(LOG_WARNING, "WINDOW: MSAA can only be configured before window initialization");
    }

    // State change: FLAG_INTERLACED_HINT
    if ((FLAG_IS_SET(CORE.Window.flags, FLAG_INTERLACED_HINT)) && (FLAG_IS_SET(flags, FLAG_INTERLACED_HINT)))
    {
        TRACELOG(LOG_WARNING, "WINDOW: Interlaced mode can only be configured before window initialization");
    }

    XFlush(platform.display);
}

// Set icon for window
// NOTE: Image must be in RGBA format, 8bit per channel
void SetWindowIcon(Image image)
{
    if (image.data == NULL) SetWindowIcons(NULL, 0);
    else SetWindowIcons(&image, 1);
}

// Set icon for window, multiple images
// NOTE: Images must be in RGBA format, 8bit per channel, the window manager uses the most appropriate size
void SetWindowIcons(Image *images, int count)
{
    int dataCount = 0;
    int iconCount = 0;

    for (int i = 0; (images != NULL) && (i < count) && (iconCount < MAX_WINDOW_ICONS); i++)
    {
        if (images[i].format != PIXELFORMAT_UNCOMPRESSED_R8G8B8A8)
        {
            TRACELOG(LOG_WARNING, "X11: Window icon image must be in R8G8B8A8 pixel format");
            continue;
        }

        dataCount += 2 + images[i].width*images[i].height;
        iconCount++;
    }

    if (dataCount == 0)
    {
        XDeleteProperty(platform.display, platform.window, platform.atoms._NET_WM_ICON);
        XFlush(platform.display);
        return;
    }

    // NOTE: Property items of 32 bits are stored as long: width, height and ARGB pixels for every icon
    unsigned long *data = (unsigned long *)RL_CALLOC(dataCount, sizeof(unsigned long));
    unsigned long *target = data;
    iconCount = 0;

    for (int i = 0; (i < count) && (iconCount < MAX_WINDOW_ICONS); i++)
    {
        if (images[i].format != PIXELFORMAT_UNCOMPRESSED_R8G8B8A8) continue;

        *target++ = images[i].width;
        *target++ = images[i].height;

        const unsigned char *pixels = (const unsigned char *)images[i].data;
        for (int p = 0; p < images[i].width*images[i].height; p++)
        {
            *target++ = ((unsigned long)pixels[p*4 + 3] << 24) | ((unsigned long)pixels[p*4 + 0] << 16) |
                ((unsigned long)pixels[p*4 + 1] << 8) | (unsigned long)pixels[p*4 + 2];
        }

        iconCount++;
    }

    XChangeProperty(platform.display, platform.window, platform.atoms._NET_WM_ICON, XA_CARDINAL, 32, PropModeReplace,
        (unsigned char *)data, dataCount);
    XFlush(platform.display);

    RL_FREE(data);
}

// Set title for window
void SetWindowTitle(const char *title)
{
    CORE.Window.title = title;
    UpdateWindowTitle(title);
    XFlush(platform.display);
}

// Set window position on screen (windowed mode)
void SetWindowPosition(int x, int y)
{
    if (platform.fullscreen) return;

    XMoveWindow(platform.display, platform.window, x, y);
    XFlush(platform.display);
}

// Set monitor for the current window
void SetWindowMonitor(int monitor)
{
    if ((monitor >= 0) && (monitor < platform.monitorCount))
    {
        if (FLAG_IS_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE) || FLAG_IS_SET(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE))
        {
            TRACELOG(LOG_INFO, "X11: Selected fullscreen monitor: [%i] %s", monitor, platform.monitors[monitor].name);

            RequestFullscreen(true, monitor);
        }
        else
        {
            TRACELOG(LOG_INFO, "X11: Selected monitor: [%i] %s", monitor, platform.monitors[monitor].name);

            // Center the window on the selected monitor
            const MonitorInfo *info = &platform.monitors[monitor];
            SetWindowPosition(info->x + (info->width - platform.bufferWidth)/2, info->y + (info->height - platform.bufferHeight)/2);
        }
    }
    else TRACELOG(LOG_WARNING, "X11: Failed to find selected monitor");
}

// Set window minimum dimensions (FLAG_WINDOW_RESIZABLE)
void SetWindowMinSize(int width, int height)
{
    CORE.Window.screenMin.width = width;
    CORE.Window.screenMin.height = height;

    UpdateNormalHints();
    SetWindowSize(platform.windowedWidth, platform.windowedHeight);
}

// Set window maximum dimensions (FLAG_WINDOW_RESIZABLE)
void SetWindowMaxSize(int width, int height)
{
    CORE.Window.screenMax.width = width;
    CORE.Window.screenMax.height = height;

    UpdateNormalHints();
    SetWindowSize(platform.windowedWidth, platform.windowedHeight);
}

// Set window dimensions
// NOTE: Size is limited to the window minimum and maximum dimensions
void SetWindowSize(int width, int height)
{
    if ((width <= 0) || (height <= 0)) return;

    if ((CORE.Window.screenMin.width > 0) && (width < (int)CORE.Window.screenMin.width)) width = CORE.Window.screenMin.width;
    if ((CORE.Window.screenMin.height > 0) && (height < (int)CORE.Window.screenMin.height)) height = CORE.Window.screenMin.height;
    if ((CORE.Window.screenMax.width > 0) && (width > (int)CORE.Window.screenMax.width)) width = CORE.Window.screenMax.width;
    if ((CORE.Window.screenMax.height > 0) && (height > (int)CORE.Window.screenMax.height)) height = CORE.Window.screenMax.height;

    if ((width == platform.width) && (height == platform.height) && (width == platform.windowedWidth) && (height == platform.windowedHeight)) return;

    platform.windowedWidth = width;
    platform.windowedHeight = height;

    if (!platform.fullscreen && !platform.maximized)
    {
        float scale = FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_HIGHDPI)? platform.scale : 1.0f;

        platform.width = width;
        platform.height = height;
        platform.bufferWidth = (int)roundf((float)width*scale);
        platform.bufferHeight = (int)roundf((float)height*scale);

        // NOTE: Non-resizable windows have minimum and maximum size hints set to the window size
        UpdateNormalHints();
        UpdateWindowBuffers();

        // NOTE: Size events generated before this request are outdated
        platform.resizeSerial = NextRequest(platform.display);
        XResizeWindow(platform.display, platform.window, platform.bufferWidth, platform.bufferHeight);
        XFlush(platform.display);
    }
}

// Set window opacity, value opacity is between 0.0 and 1.0
// NOTE: Requires a compositing manager
void SetWindowOpacity(float opacity)
{
    if (opacity >= 1.0f) opacity = 1.0f;
    else if (opacity <= 0.0f) opacity = 0.0f;

    if (opacity >= 1.0f) XDeleteProperty(platform.display, platform.window, platform.atoms._NET_WM_WINDOW_OPACITY);
    else
    {
        unsigned long value = (unsigned long)((double)opacity*(double)0xffffffffu);
        XChangeProperty(platform.display, platform.window, platform.atoms._NET_WM_WINDOW_OPACITY, XA_CARDINAL, 32, PropModeReplace,
            (unsigned char *)&value, 1);
    }

    XFlush(platform.display);
}

// Set window focused
// NOTE: The window manager decides if focus is granted
void SetWindowFocused(void)
{
    if (platform.withdrawn) return;

    if (IsWMHintSupported(platform.atoms._NET_ACTIVE_WINDOW))
    {
        // NOTE: Source indication 1 (application), time of the last user interaction
        SendWindowMessage(platform.atoms._NET_ACTIVE_WINDOW, 1, (long)platform.lastUserTime, 0, 0, 0);
    }
    else if (platform.mapped)
    {
        XRaiseWindow(platform.display, platform.window);
        XSetInputFocus(platform.display, platform.window, RevertToParent, CurrentTime);
    }

    XFlush(platform.display);
}

// Get native window handle
// NOTE: Returns the X11 window id
void *GetWindowHandle(void)
{
    return (void *)platform.window;
}

// Get number of monitors
int GetMonitorCount(void)
{
    return platform.monitorCount;
}

// Get current monitor where window is placed
// NOTE: Monitor with the largest part of the window
int GetCurrentMonitor(void)
{
    return GetMonitorAt(platform.x, platform.y, platform.bufferWidth, platform.bufferHeight);
}

// Get selected monitor position
Vector2 GetMonitorPosition(int monitor)
{
    if ((monitor >= 0) && (monitor < platform.monitorCount))
    {
        return (Vector2){ (float)platform.monitors[monitor].x, (float)platform.monitors[monitor].y };
    }
    else TRACELOG(LOG_WARNING, "X11: Failed to find selected monitor");

    return (Vector2){ 0.0f, 0.0f };
}

// Get selected monitor width (currently used by monitor)
int GetMonitorWidth(int monitor)
{
    int width = 0;

    if ((monitor >= 0) && (monitor < platform.monitorCount)) width = platform.monitors[monitor].width;
    else TRACELOG(LOG_WARNING, "X11: Failed to find selected monitor");

    return width;
}

// Get selected monitor height (currently used by monitor)
int GetMonitorHeight(int monitor)
{
    int height = 0;

    if ((monitor >= 0) && (monitor < platform.monitorCount)) height = platform.monitors[monitor].height;
    else TRACELOG(LOG_WARNING, "X11: Failed to find selected monitor");

    return height;
}

// Get selected monitor physical width in millimetres
int GetMonitorPhysicalWidth(int monitor)
{
    int width = 0;

    if ((monitor >= 0) && (monitor < platform.monitorCount)) width = platform.monitors[monitor].physicalWidth;
    else TRACELOG(LOG_WARNING, "X11: Failed to find selected monitor");

    return width;
}

// Get selected monitor physical height in millimetres
int GetMonitorPhysicalHeight(int monitor)
{
    int height = 0;

    if ((monitor >= 0) && (monitor < platform.monitorCount)) height = platform.monitors[monitor].physicalHeight;
    else TRACELOG(LOG_WARNING, "X11: Failed to find selected monitor");

    return height;
}

// Get selected monitor refresh rate
int GetMonitorRefreshRate(int monitor)
{
    int refresh = 0;

    if ((monitor >= 0) && (monitor < platform.monitorCount)) refresh = (int)roundf(platform.monitors[monitor].refreshRate);
    else TRACELOG(LOG_WARNING, "X11: Failed to find selected monitor");

    return refresh;
}

// Get the human-readable, UTF-8 encoded name of the selected monitor
const char *GetMonitorName(int monitor)
{
    if ((monitor >= 0) && (monitor < platform.monitorCount))
    {
        return platform.monitors[monitor].name;
    }
    else TRACELOG(LOG_WARNING, "X11: Failed to find selected monitor");

    return "";
}

// Get window position XY on monitor
Vector2 GetWindowPosition(void)
{
    return (Vector2){ (float)platform.x, (float)platform.y };
}

// Get window scale DPI factor for current monitor
Vector2 GetWindowScaleDPI(void)
{
    return (Vector2){ platform.scale, platform.scale };
}

// Set clipboard text content
void SetClipboardText(const char *text)
{
    RL_FREE(platform.clipboardText);
    platform.clipboardText = NULL;

    if (text == NULL)
    {
        if (platform.clipboardOwned) XSetSelectionOwner(platform.display, platform.atoms.CLIPBOARD, None, GetServerTime());
        platform.clipboardOwned = false;
        XFlush(platform.display);
        return;
    }

    int length = (int)strlen(text);
    platform.clipboardText = (char *)RL_CALLOC(length + 1, sizeof(char));
    memcpy(platform.clipboardText, text, length);

    platform.clipboardTime = GetServerTime();
    XSetSelectionOwner(platform.display, platform.atoms.CLIPBOARD, platform.window, platform.clipboardTime);
    platform.clipboardOwned = (XGetSelectionOwner(platform.display, platform.atoms.CLIPBOARD) == platform.window);

    if (!platform.clipboardOwned) TRACELOG(LOG_WARNING, "X11: Failed to become owner of the clipboard selection");
}

// Get clipboard text content
// NOTE: returned string is allocated and freed by the platform
const char *GetClipboardText(void)
{
    if (platform.clipboardOwned) return platform.clipboardText;

    RL_FREE(platform.clipboardReceived);
    platform.clipboardReceived = NULL;

    if (XGetSelectionOwner(platform.display, platform.atoms.CLIPBOARD) == None) return NULL;

    const Atom targets[] = { platform.atoms.UTF8_STRING, platform.atoms.TEXT_PLAIN_UTF8, XA_STRING };

    for (int i = 0; i < (int)(sizeof(targets)/sizeof(targets[0])); i++)
    {
        int size = 0;
        Atom type = None;
        unsigned char *data = ReadSelection(platform.atoms.CLIPBOARD, targets[i], &size, &type);
        if (data == NULL) continue;

        if (type == XA_STRING)
        {
            // Convert Latin-1 text to UTF-8
            char *text = (char *)RL_CALLOC(size*2 + 1, sizeof(char));
            int length = 0;

            for (int c = 0; c < size; c++)
            {
                if (data[c] < 0x80) text[length++] = (char)data[c];
                else
                {
                    text[length++] = (char)(0xc0 | (data[c] >> 6));
                    text[length++] = (char)(0x80 | (data[c] & 0x3f));
                }
            }

            RL_FREE(data);
            platform.clipboardReceived = text;
        }
        else platform.clipboardReceived = (char *)data;

        break;
    }

    return platform.clipboardReceived;
}

// Get clipboard image
Image GetClipboardImage(void)
{
    Image image = { 0 };

#if SUPPORT_CLIPBOARD_IMAGE && SUPPORT_MODULE_RTEXTURES
    if (!platform.clipboardOwned && (XGetSelectionOwner(platform.display, platform.atoms.CLIPBOARD) != None))
    {
        const Atom imageTargets[] = { platform.atoms.IMAGE_PNG, platform.atoms.IMAGE_BMP, platform.atoms.IMAGE_JPEG };
        const char *imageFileTypes[] = { ".png", ".bmp", ".jpg" };

        for (int i = 0; i < (int)(sizeof(imageTargets)/sizeof(imageTargets[0])); i++)
        {
            if (SelectionHasTarget(platform.atoms.CLIPBOARD, imageTargets[i]))
            {
                int size = 0;
                unsigned char *data = ReadSelection(platform.atoms.CLIPBOARD, imageTargets[i], &size, NULL);

                if (data != NULL)
                {
                    image = LoadImageFromMemory(imageFileTypes[i], data, size);
                    RL_FREE(data);
                }
                break;
            }
        }
    }

    if (image.data == NULL) TRACELOG(LOG_WARNING, "Clipboard image: Couldn't get clipboard data.");
#else
    TRACELOG(LOG_WARNING, "Clipboard image: SUPPORT_CLIPBOARD_IMAGE requires SUPPORT_MODULE_RTEXTURES to work properly");
#endif

    return image;
}

// Show mouse cursor
void ShowCursor(void)
{
    CORE.Input.Mouse.cursorHidden = false;
    UpdateCursor();
}

// Hide mouse cursor
void HideCursor(void)
{
    CORE.Input.Mouse.cursorHidden = true;
    UpdateCursor();
}

// Enable cursor (unlock cursor)
void EnableCursor(void)
{
    LockPointer(false);

    CORE.Input.Mouse.cursorHidden = false;
    CORE.Input.Mouse.cursorLocked = false;

    // Set cursor position in the middle
    SetMousePosition(CORE.Window.screen.width/2, CORE.Window.screen.height/2);

    UpdateCursor();
}

// Disable cursor (lock cursor)
void DisableCursor(void)
{
    // Reset mouse position within the window area before disabling cursor
    SetMousePosition(CORE.Window.screen.width/2, CORE.Window.screen.height/2);

    CORE.Input.Mouse.cursorHidden = true;
    CORE.Input.Mouse.cursorLocked = true;

    UpdateCursor();
    LockPointer(true);
}

// Swap back buffer with front buffer (screen drawing)
void SwapScreenBuffer(void)
{
    if (platform.connectionLost) return;

    // NOTE: Frames of a window that is not mapped are not presented, they are paced like occluded windows
    if (!platform.mapped)
    {
        if (FLAG_IS_SET(CORE.Window.flags, FLAG_VSYNC_HINT)) WaitTime(1.0/60.0);
        return;
    }

#if defined(GRAPHICS_API_OPENGL_SOFTWARE)
    if (platform.image != NULL)
    {
        int width = platform.image->width;
        int height = platform.image->height;

        // NOTE: Software renderer framebuffer copy is in BGRA order, the byte order of 32 bit images on little-endian servers
        rlCopyFramebuffer(0, 0, width, height, PIXELFORMAT_UNCOMPRESSED_R8G8B8A8, platform.image->data);

        if (platform.shmAttached) XShmPutImage(platform.display, platform.window, platform.gc, platform.image, 0, 0, 0, 0, width, height, False);
        else XPutImage(platform.display, platform.window, platform.gc, platform.image, 0, 0, 0, 0, width, height);

        // NOTE: Wait for the server to process the image, frames are not queued
        XSync(platform.display, False);
    }
#else
    eglSwapBuffers(platform.eglDisplay, platform.eglSurface);
#endif

    // Notify the window manager that the frame for the last resize is presented
    if (platform.syncPending)
    {
        XSyncSetCounter(platform.display, platform.syncCounter, platform.syncValue);
        XFlush(platform.display);
        platform.syncPending = false;
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
        char *cmd = (char *)RL_CALLOC(strlen(url) + 32, sizeof(char));
        sprintf(cmd, "xdg-open '%s'", url); // Alternatives: firefox, x-www-browser
        int result = system(cmd);

        if (result == -1) TRACELOG(LOG_WARNING, "OpenURL() child process could not be created");
        RL_FREE(cmd);
    }
}

//----------------------------------------------------------------------------------
// Module Functions Definition: Inputs
//----------------------------------------------------------------------------------

// Set internal gamepad mappings
int SetGamepadMappings(const char *mappings)
{
#if defined(__linux__)
    return mg_update_gamepad_mappings(&platform.minigamepad, mappings);
#else
    TRACELOG(LOG_WARNING, "SetGamepadMappings() not available on target platform");
    return 0;
#endif
}

// Set gamepad vibration
void SetGamepadVibration(int gamepad, float leftMotor, float rightMotor, float duration)
{
    TRACELOG(LOG_WARNING, "SetGamepadVibration() not available on target platform");
}

// Set mouse position XY
// NOTE: While the cursor is disabled only the virtual position is set, the hidden cursor stays in the window
void SetMousePosition(int x, int y)
{
    CORE.Input.Mouse.currentPosition = (Vector2){ (float)x, (float)y };
    CORE.Input.Mouse.previousPosition = CORE.Input.Mouse.currentPosition;

    if (!CORE.Input.Mouse.cursorLocked)
    {
        float scale = FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_HIGHDPI)? platform.scale : 1.0f;
        WarpPointer((int)roundf((float)x*scale), (int)roundf((float)y*scale));
        XFlush(platform.display);
    }
}

// Set mouse cursor
void SetMouseCursor(int cursor)
{
    CORE.Input.Mouse.cursor = cursor;
    UpdateCursor();
}

// Get physical key name
// NOTE: Returns the character the key produces on the current keyboard layout, NULL for non-printable keys
const char *GetKeyName(int key)
{
    static char keyName[16] = { 0 };

    if (!platform.xkbAvailable || (key <= 0)) return NULL;

    int keycode = -1;
    for (int i = 0; i < KEYMAP_SIZE; i++)
    {
        if (platform.keymap[i] == key)
        {
            keycode = i;
            break;
        }
    }

    if (keycode < 0) return NULL;

    KeySym keysym = XkbKeycodeToKeysym(platform.display, (KeyCode)keycode, platform.xkbGroup, 0);

    // NOTE: Keypad keys produce characters on the second level (NUM lock)
    if ((key >= KEY_KP_0) && (key <= KEY_KP_EQUAL)) keysym = XkbKeycodeToKeysym(platform.display, (KeyCode)keycode, platform.xkbGroup, 1);

    if ((keysym == NoSymbol) || (xkb_keysym_to_utf32((xkb_keysym_t)keysym) <= 32)) return NULL;

    memset(keyName, 0, sizeof(keyName));
    if (xkb_keysym_to_utf8((xkb_keysym_t)keysym, keyName, sizeof(keyName)) <= 0) return NULL;

    return keyName;
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

    // Register keys and mouse buttons released in the frame they were pressed
    bool inputReleased = false;
    for (int i = 0; i < MAX_KEYBOARD_KEYS; i++)
    {
        platform.keyPressedInFrame[i] = false;

        if (platform.keyReleasePending[i])
        {
            platform.keyReleasePending[i] = false;
            ReleaseKey(i);
            inputReleased = true;
        }
    }

    for (int i = 0; i < MAX_MOUSE_BUTTONS; i++)
    {
        platform.buttonPressedInFrame[i] = false;

        if (platform.buttonReleasePending[i])
        {
            platform.buttonReleasePending[i] = false;
            HandleMouseButton(i, false);
            inputReleased = true;
        }
    }

    // Register previous mouse wheel state
    CORE.Input.Mouse.previousWheelMove = CORE.Input.Mouse.currentWheelMove;
    CORE.Input.Mouse.currentWheelMove = (Vector2){ 0.0f, 0.0f };

    // Register previous mouse position
    CORE.Input.Mouse.previousPosition = CORE.Input.Mouse.currentPosition;

    // Register previous touch states
    for (int i = 0; i < MAX_TOUCH_POINTS; i++)
    {
        CORE.Input.Touch.previousTouchState[i] = CORE.Input.Touch.currentTouchState[i];
        CORE.Input.Touch.previousPosition[i] = CORE.Input.Touch.position[i];
    }

    // Map touch position to mouse position for convenience
    if (CORE.Input.Touch.pointCount == 0) CORE.Input.Touch.position[0] = CORE.Input.Mouse.currentPosition;

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
    if ((CORE.Window.eventWaiting) ||
        (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_MINIMIZED) && !FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_ALWAYS_RUN)))
    {
        // Wait for input events before continue (drawing is paused)
        // NOTE: A release registered in this frame is an input event, there is no waiting
        ProcessEvents(inputReleased? 0 : -1);
        CORE.Time.previous = GetTime();
    }
    else ProcessEvents(0);

    // NOTE: Without XInput2 raw motion, the hidden cursor is moved back to the window center to keep reporting motion
    if (platform.pointerGrabbed && !platform.xiAvailable &&
        ((platform.pointerX != (double)(platform.bufferWidth/2)) || (platform.pointerY != (double)(platform.bufferHeight/2))))
    {
        WarpPointer(platform.bufferWidth/2, platform.bufferHeight/2);
        XFlush(platform.display);
    }

#if defined(GRAPHICS_API_OPENGL_SOFTWARE)
    // Resize software renderer framebuffer to the window size
    if ((platform.image != NULL) && ((platform.image->width != (int)CORE.Window.render.width) || (platform.image->height != (int)CORE.Window.render.height)))
    {
        DestroySoftwareImage();
        if (CreateSoftwareImage(CORE.Window.render.width, CORE.Window.render.height)) rlResizeFramebuffer(CORE.Window.render.width, CORE.Window.render.height);
    }
#endif

#if defined(__linux__)
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
#endif

    CORE.Window.shouldClose = (platform.closeRequested || platform.connectionLost);

    // Reset close status for next frame
    platform.closeRequested = false;
}

//----------------------------------------------------------------------------------
// Module Internal Functions Definition
//----------------------------------------------------------------------------------
// Initialize platform: graphics, inputs and more
int InitPlatform(void)
{
    if (!LoadLibraries())
    {
        CloseLibraries();
        return -1;
    }

    XInitThreads();

    platform.display = XOpenDisplay(NULL);
    if (platform.display == NULL)
    {
        TRACELOG(LOG_WARNING, "X11: Failed to open display %s", XDisplayName(NULL));
        return -1;
    }

    platform.screen = DefaultScreen(platform.display);
    platform.root = RootWindow(platform.display, platform.screen);
    platform.scale = 1.0f;
    platform.withdrawn = true;
    for (int i = 0; i < MAX_TOUCH_POINTS; i++) platform.touchIds[i] = -1;

    // NOTE: X11 errors terminate the process by default, most of them are not fatal (e.g. requests to windows of other clients)
    platform.previousErrorHandler = XSetErrorHandler(ErrorHandler);

    // NOTE: Connection errors terminate the process, the default handler exits; with libX11 1.7 or later
    // an exit handler can be set to continue running, the window should close
    void (*setIOErrorExitHandler)(Display *, void (*)(Display *, void *), void *) =
        (void (*)(Display *, void (*)(Display *, void *), void *))dlsym(platform.libX11, "XSetIOErrorExitHandler");
    if (setIOErrorExitHandler != NULL)
    {
        platform.previousIOErrorHandler = XSetIOErrorHandler(IOErrorHandler);
        platform.ioErrorHandlerSet = true;
        setIOErrorExitHandler(platform.display, IOErrorExitHandler, NULL);
    }

    InitAtoms();
    InitExtensions();
    UpdateWMSupported();
    UpdateMonitors();
    XrmInitialize();
    UpdateWindowScale();

    // Get notified about resource manager (Xft.dpi) and window manager changes
    XSelectInput(platform.display, platform.root, PropertyChangeMask);
    if (platform.randrAvailable) XRRSelectInput(platform.display, platform.root, RRScreenChangeNotifyMask | RRCrtcChangeNotifyMask | RROutputChangeNotifyMask);

    // Initialize window size
    float scale = FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_HIGHDPI)? platform.scale : 1.0f;

    if ((CORE.Window.screen.width == 0) || (CORE.Window.screen.height == 0)) FLAG_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE);

    CORE.Window.display.width = platform.monitors[0].width;
    CORE.Window.display.height = platform.monitors[0].height;

    if (FLAG_IS_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE))
    {
        // Check if user requested some screen size
        if ((CORE.Window.screen.width == 0) || (CORE.Window.screen.height == 0))
        {
            // Set some default screen size in case user decides to exit fullscreen mode
            CORE.Window.previousScreen.width = 800;
            CORE.Window.previousScreen.height = 450;
        }
        else CORE.Window.previousScreen = CORE.Window.screen;

        // NOTE: Fullscreen window covers the monitor, the video mode is not changed
        platform.bufferWidth = platform.monitors[0].width;
        platform.bufferHeight = platform.monitors[0].height;
        platform.width = (int)roundf((float)platform.bufferWidth/scale);
        platform.height = (int)roundf((float)platform.bufferHeight/scale);
        platform.windowedWidth = CORE.Window.previousScreen.width;
        platform.windowedHeight = CORE.Window.previousScreen.height;
    }
    else
    {
        platform.width = CORE.Window.screen.width;
        platform.height = CORE.Window.screen.height;
        platform.bufferWidth = (int)roundf((float)platform.width*scale);
        platform.bufferHeight = (int)roundf((float)platform.height*scale);
        platform.windowedWidth = platform.width;
        platform.windowedHeight = platform.height;
    }

    // Window application id, used by window managers to group windows and to find the desktop entry
    snprintf(platform.appId, sizeof(platform.appId), "raylib");
#if defined(__linux__)
    char exePath[512] = { 0 };
    ssize_t exePathLength = readlink("/proc/self/exe", exePath, sizeof(exePath) - 1);
    if (exePathLength > 0) snprintf(platform.appId, sizeof(platform.appId), "%s", GetFileName(exePath));
#else
    if (getprogname() != NULL) snprintf(platform.appId, sizeof(platform.appId), "%s", getprogname());
#endif

    platform.decorationRequested = !FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_UNDECORATED);

    // Initialize graphics device and create window
    // NOTE: Window visual is defined by the graphics config
    if (!InitGraphicsContext() || !CreateWindow() || !InitGraphicsSurface())
    {
        TRACELOG(LOG_FATAL, "PLATFORM: Failed to initialize graphics device");
        ClosePlatform();
        return -1;
    }

    InitKeyboard();
    InitInputDevices();
    UpdateWindowBuffers();
    UpdateCursor();

    if (!FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_HIDDEN)) ShowWindow(true);

    // Initial mouse position, the first mouse movement is relative to it
    Window pointerRoot = None;
    Window pointerChild = None;
    int pointerRootX = 0;
    int pointerRootY = 0;
    int pointerX = 0;
    int pointerY = 0;
    unsigned int pointerMask = 0;
    if (XQueryPointer(platform.display, platform.window, &pointerRoot, &pointerChild, &pointerRootX, &pointerRootY, &pointerX, &pointerY, &pointerMask))
    {
        platform.pointerX = (double)pointerX;
        platform.pointerY = (double)pointerY;
        CORE.Input.Mouse.currentPosition = (Vector2){ (float)pointerX/scale, (float)pointerY/scale };
        CORE.Input.Mouse.previousPosition = CORE.Input.Mouse.currentPosition;
    }

    CompleteStartupNotification();

    CORE.Window.ready = true;

    // NOTE: Minimized state flag is set when the window manager changes the window state
    FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_MINIMIZED);

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
    rlLoadExtensions((void *)eglGetProcAddress);
#endif
    //----------------------------------------------------------------------------

    // Initialize gamepads
    //----------------------------------------------------------------------------
#if defined(__linux__)
    mg_gamepads_init(&platform.minigamepad);
#endif
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
    TRACELOG(LOG_INFO, "PLATFORM: DESKTOP (X11 - Software): Initialized successfully");
#else
    TRACELOG(LOG_INFO, "PLATFORM: DESKTOP (X11): Initialized successfully");
#endif

    return 0;
}

// Close platform
void ClosePlatform(void)
{
    if (platform.display == NULL) return;

#if defined(__linux__)
    mg_gamepads_free(&platform.minigamepad);
#endif

    if (!platform.connectionLost)
    {
        SaveClipboard();
        LockPointer(false);
    }

#if defined(GRAPHICS_API_OPENGL_SOFTWARE)
    if (!platform.connectionLost)
    {
        DestroySoftwareImage();
        if (platform.gc != NULL) XFreeGC(platform.display, platform.gc);
    }
#else
    // NOTE: EGL resources are released before closing the display connection they use
    if (platform.eglDisplay != EGL_NO_DISPLAY)
    {
        eglMakeCurrent(platform.eglDisplay, EGL_NO_SURFACE, EGL_NO_SURFACE, EGL_NO_CONTEXT);
        if (platform.eglSurface != EGL_NO_SURFACE) eglDestroySurface(platform.eglDisplay, platform.eglSurface);
        if (platform.eglContext != EGL_NO_CONTEXT) eglDestroyContext(platform.eglDisplay, platform.eglContext);
        eglTerminate(platform.eglDisplay);
    }
#endif

    XUnregisterIMInstantiateCallback(platform.display, NULL, NULL, NULL, InputMethodInstantiateCallback, NULL);

    // NOTE: No requests can be sent after a connection error, server resources are already released
    if (!platform.connectionLost)
    {
        if (platform.inputContext != NULL) XDestroyIC(platform.inputContext);
        if (platform.inputMethod != NULL) XCloseIM(platform.inputMethod);

        for (int i = 0; i < MOUSE_CURSOR_COUNT; i++)
        {
            if (platform.cursors[i] != None) XFreeCursor(platform.display, platform.cursors[i]);
        }
        if (platform.hiddenCursor != None) XFreeCursor(platform.display, platform.hiddenCursor);

        if (platform.syncCounter != None) XSyncDestroyCounter(platform.display, platform.syncCounter);
        if (platform.window != None) XDestroyWindow(platform.display, platform.window);
        if (platform.colormap != None) XFreeColormap(platform.display, platform.colormap);
    }

    for (int i = 0; i < MAX_INCR_TRANSFERS; i++) RL_FREE(platform.incrTransfers[i].data);
    RL_FREE(platform.wmSupported);
    RL_FREE(platform.clipboardText);
    RL_FREE(platform.clipboardReceived);

    XCloseDisplay(platform.display);

    XSetErrorHandler(platform.previousErrorHandler);
    if (platform.ioErrorHandlerSet) XSetIOErrorHandler(platform.previousIOErrorHandler);

    CloseLibraries();

    memset(&platform, 0, sizeof(platform));
}

// Load system libraries and their functions
static bool LoadLibraries(void)
{
    struct { void **handle; const char *names[2]; } libraries[] = {
        { &platform.libX11, { "libX11.so.6", "libX11.so" } },
        { &platform.libXext, { "libXext.so.6", "libXext.so" } },
        { &platform.libXrandr, { "libXrandr.so.2", "libXrandr.so" } },
        { &platform.libXi, { "libXi.so.6", "libXi.so" } },
        { &platform.libXcursor, { "libXcursor.so.1", "libXcursor.so" } },
        { &platform.libxkbcommon, { "libxkbcommon.so.0", "libxkbcommon.so" } },
#if !defined(GRAPHICS_API_OPENGL_SOFTWARE)
        { &platform.libEGL, { "libEGL.so.1", "libEGL.so" } },
#endif
    };

    for (int i = 0; i < (int)(sizeof(libraries)/sizeof(libraries[0])); i++)
    {
        *libraries[i].handle = dlopen(libraries[i].names[0], RTLD_LAZY | RTLD_LOCAL);
        if (*libraries[i].handle == NULL) *libraries[i].handle = dlopen(libraries[i].names[1], RTLD_LAZY | RTLD_LOCAL);
        if (*libraries[i].handle == NULL)
        {
            TRACELOG(LOG_WARNING, "X11: Failed to load %s", libraries[i].names[0]);
            return false;
        }
    }

    bool loaded = true;
    void *library = NULL;

    #define LOAD_FUNCTION(name) \
        platform.lib.p##name = (__typeof__(platform.lib.p##name))dlsym(library, #name); \
        if (platform.lib.p##name == NULL) { TRACELOG(LOG_WARNING, "X11: Failed to load function %s", #name); loaded = false; }

    library = platform.libX11;
    LIBX11_FUNCTIONS(LOAD_FUNCTION)
    library = platform.libXext;
    LIBXEXT_FUNCTIONS(LOAD_FUNCTION)
    LIBXSHM_FUNCTIONS(LOAD_FUNCTION)
    library = platform.libXrandr;
    LIBXRANDR_FUNCTIONS(LOAD_FUNCTION)
    library = platform.libXi;
    LIBXI_FUNCTIONS(LOAD_FUNCTION)
    library = platform.libXcursor;
    LIBXCURSOR_FUNCTIONS(LOAD_FUNCTION)
    library = platform.libxkbcommon;
    LIBXKBCOMMON_FUNCTIONS(LOAD_FUNCTION)
    library = platform.libEGL;
    LIBEGL_FUNCTIONS(LOAD_FUNCTION)

    #undef LOAD_FUNCTION

    return loaded;
}

// Close system libraries
static void CloseLibraries(void)
{
    void **libraries[] = { &platform.libEGL, &platform.libxkbcommon, &platform.libXcursor, &platform.libXi, &platform.libXrandr, &platform.libXext, &platform.libX11 };

    for (int i = 0; i < (int)(sizeof(libraries)/sizeof(libraries[0])); i++)
    {
        if (*libraries[i] != NULL) dlclose(*libraries[i]);
        *libraries[i] = NULL;
    }
}

// X11 errors handler, errors are not fatal
static int ErrorHandler(Display *display, XErrorEvent *event)
{
    platform.errorCode = event->error_code;

    TRACELOG(LOG_DEBUG, "X11: Error %i on request %i.%i", event->error_code, event->request_code, event->minor_code);

    return 0;
}

// X11 connection errors handler
// NOTE: The display can not be used after a connection error, the window should close
static int IOErrorHandler(Display *display)
{
    if (!platform.connectionLost) TRACELOG(LOG_WARNING, "X11: Connection to the X server lost");

    platform.connectionLost = true;

    return 0;
}

// X11 connection errors exit handler, the process is not terminated
static void IOErrorExitHandler(Display *display, void *data)
{
    platform.connectionLost = true;
}

// Intern required atoms
// NOTE: All the atoms are requested at once (one round trip to the server)
static void InitAtoms(void)
{
    XInternAtoms(platform.display, (char **)atomNames, (int)(sizeof(atomNames)/sizeof(atomNames[0])), False, (Atom *)&platform.atoms);
}

// Query X11 extensions
// NOTE: Extensions are optional, features depending on a missing extension are not available
static void InitExtensions(void)
{
    int opcode = 0;
    int eventBase = 0;
    int errorBase = 0;
    int major = 0;
    int minor = 0;

    // Keyboard extension: physical key names, layout groups and detectable key repeat
    major = XkbMajorVersion;
    minor = XkbMinorVersion;
    if (XkbQueryExtension(platform.display, &opcode, &eventBase, &errorBase, &major, &minor))
    {
        platform.xkbAvailable = true;
        platform.xkbEventBase = eventBase;
    }

    // RandR: monitors
    if (XRRQueryExtension(platform.display, &eventBase, &errorBase) && XRRQueryVersion(platform.display, &major, &minor))
    {
        platform.randrAvailable = ((major > 1) || (minor >= 2));
        platform.randrMonitors = ((major > 1) || (minor >= 5));
        platform.randrEventBase = eventBase;
    }

    // XInput2: smooth scrolling (2.1), touch (2.2) and raw mouse motion
    if (XQueryExtension(platform.display, "XInputExtension", &opcode, &eventBase, &errorBase))
    {
        major = 2;
        minor = 2;
        if ((XIQueryVersion(platform.display, &major, &minor) == Success) && (major == 2))
        {
            platform.xiAvailable = true;
            platform.xiOpcode = opcode;
            platform.xiMinor = minor;
        }
    }

    // Shape: input shape for mouse passthrough (1.1)
    if (XShapeQueryExtension(platform.display, &eventBase, &errorBase) && XShapeQueryVersion(platform.display, &major, &minor))
    {
        platform.shapeAvailable = ((major > 1) || (minor >= 1));
    }

    // Sync: frame counter for resizes synchronized with the window manager
    if (XSyncQueryExtension(platform.display, &eventBase, &errorBase) && XSyncInitialize(platform.display, &major, &minor))
    {
        platform.syncAvailable = true;
    }

#if defined(GRAPHICS_API_OPENGL_SOFTWARE)
    // NOTE: Shared memory images use the byte order of the server, the framebuffer is copied in little-endian order
    platform.shmAvailable = ((XShmQueryExtension(platform.display) == True) && (ImageByteOrder(platform.display) == LSBFirst));
#endif
}

// Initialize EGL display and config, choose the window visual
static bool InitGraphicsContext(void)
{
    bool transparent = FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_TRANSPARENT);

    platform.visual = DefaultVisual(platform.display, platform.screen);
    platform.depth = DefaultDepth(platform.display, platform.screen);

#if defined(GRAPHICS_API_OPENGL_SOFTWARE)
    // NOTE: Framebuffer is presented as a 32 bit per pixel image, transparency requires a visual with alpha channel
    XVisualInfo visualInfo = { 0 };
    if (XMatchVisualInfo(platform.display, platform.screen, transparent? 32 : 24, TrueColor, &visualInfo))
    {
        platform.visual = visualInfo.visual;
        platform.depth = visualInfo.depth;
    }
    else if (transparent) TRACELOG(LOG_WARNING, "DISPLAY: Failed to find a visual with alpha channel, window is not transparent");

    return true;
#else
    // NOTE: EGL platform display is requested by extension, eglGetPlatformDisplay() requires EGL 1.5
    const char *clientExtensions = eglQueryString(EGL_NO_DISPLAY, EGL_EXTENSIONS);
    if ((clientExtensions != NULL) && (strstr(clientExtensions, "EGL_EXT_platform_x11") != NULL))
    {
        EGLDisplay (*getPlatformDisplay)(EGLenum, void *, const EGLint *) =
            (EGLDisplay (*)(EGLenum, void *, const EGLint *))eglGetProcAddress("eglGetPlatformDisplayEXT");
        const EGLint displayAttribs[] = { EGL_PLATFORM_X11_SCREEN_EXT, platform.screen, EGL_NONE };
        if (getPlatformDisplay != NULL) platform.eglDisplay = getPlatformDisplay(EGL_PLATFORM_X11_EXT, platform.display, displayAttribs);
        platform.eglPlatformDisplay = (platform.eglDisplay != EGL_NO_DISPLAY);
    }

    if (platform.eglDisplay == EGL_NO_DISPLAY) platform.eglDisplay = eglGetDisplay((EGLNativeDisplayType)platform.display);
    if (platform.eglDisplay == EGL_NO_DISPLAY)
    {
        TRACELOG(LOG_WARNING, "DISPLAY: Failed to get EGL display");
        return false;
    }

    EGLint majorVersion = 0;
    EGLint minorVersion = 0;
    if (eglInitialize(platform.eglDisplay, &majorVersion, &minorVersion) == EGL_FALSE)
    {
        TRACELOG(LOG_WARNING, "DISPLAY: Failed to initialize EGL display");
        return false;
    }

    bool openglES = ((rlGetVersion() == RL_OPENGL_ES_20) || (rlGetVersion() == RL_OPENGL_ES_30));
    if (eglBindAPI(openglES? EGL_OPENGL_ES_API : EGL_OPENGL_API) == EGL_FALSE)
    {
        TRACELOG(LOG_WARNING, "DISPLAY: Failed to bind EGL rendering API");
        return false;
    }

    EGLint renderableType = EGL_OPENGL_BIT;
    if (rlGetVersion() == RL_OPENGL_ES_20) renderableType = EGL_OPENGL_ES2_BIT;
    else if (rlGetVersion() == RL_OPENGL_ES_30) renderableType = EGL_OPENGL_ES3_BIT;

    EGLint samples = 0;
    EGLint sampleBuffer = 0;
    if (FLAG_IS_SET(CORE.Window.flags, FLAG_MSAA_4X_HINT))
    {
        samples = 4;
        sampleBuffer = 1;
        TRACELOG(LOG_INFO, "DISPLAY: Trying to enable MSAA x4");
    }

    EGLint alphaSize = transparent? 8 : 0;

    EGLint framebufferAttribs[] = {
        EGL_SURFACE_TYPE, EGL_WINDOW_BIT,
        EGL_RENDERABLE_TYPE, renderableType,
        EGL_RED_SIZE, 8,
        EGL_GREEN_SIZE, 8,
        EGL_BLUE_SIZE, 8,
        EGL_ALPHA_SIZE, alphaSize,
        EGL_DEPTH_SIZE, 24,
        EGL_SAMPLE_BUFFERS, sampleBuffer,
        EGL_SAMPLES, samples,
        EGL_NONE
    };

    EGLint numConfigs = 0;
    if ((eglChooseConfig(platform.eglDisplay, framebufferAttribs, NULL, 0, &numConfigs) == EGL_FALSE) || (numConfigs == 0))
    {
        if (samples > 0)
        {
            // Try again without multisampling
            TRACELOG(LOG_WARNING, "DISPLAY: MSAA x4 not available");
            framebufferAttribs[15] = 0;
            framebufferAttribs[17] = 0;
            numConfigs = 0;
            eglChooseConfig(platform.eglDisplay, framebufferAttribs, NULL, 0, &numConfigs);
        }

        if (numConfigs == 0)
        {
            TRACELOG(LOG_WARNING, "DISPLAY: Failed to choose EGL config");
            return false;
        }
    }

    EGLConfig *configs = (EGLConfig *)RL_CALLOC(numConfigs, sizeof(EGLConfig));
    if ((eglChooseConfig(platform.eglDisplay, framebufferAttribs, configs, numConfigs, &numConfigs) == EGL_FALSE) || (numConfigs == 0))
    {
        TRACELOG(LOG_WARNING, "DISPLAY: Failed to choose EGL config");
        RL_FREE(configs);
        return false;
    }

    // Choose the config by its X11 visual, the window must be created with the visual of the config
    // NOTE: A visual with alpha channel (depth 32) makes the window translucent on compositing managers,
    // it is only selected for transparent framebuffers; configs can be 10 bit per channel
    int bestScore = -1;
    for (EGLint i = 0; i < numConfigs; i++)
    {
        EGLint red = 0;
        EGLint green = 0;
        EGLint blue = 0;
        EGLint alpha = 0;
        EGLint visualId = 0;

        eglGetConfigAttrib(platform.eglDisplay, configs[i], EGL_RED_SIZE, &red);
        eglGetConfigAttrib(platform.eglDisplay, configs[i], EGL_GREEN_SIZE, &green);
        eglGetConfigAttrib(platform.eglDisplay, configs[i], EGL_BLUE_SIZE, &blue);
        eglGetConfigAttrib(platform.eglDisplay, configs[i], EGL_ALPHA_SIZE, &alpha);
        eglGetConfigAttrib(platform.eglDisplay, configs[i], EGL_NATIVE_VISUAL_ID, &visualId);

        if (visualId == 0) continue;

        XVisualInfo visualTemplate = { 0 };
        visualTemplate.visualid = (VisualID)visualId;
        visualTemplate.screen = platform.screen;
        int visualCount = 0;
        XVisualInfo *visualInfo = XGetVisualInfo(platform.display, VisualIDMask | VisualScreenMask, &visualTemplate, &visualCount);
        if (visualInfo == NULL) continue;

        int score = 0;
        if ((visualInfo->depth == 32) == transparent) score += 4;
        if ((red == 8) && (green == 8) && (blue == 8)) score += 2;
        if (alpha == alphaSize) score += 1;

        if (score > bestScore)
        {
            bestScore = score;
            platform.eglConfig = configs[i];
            platform.visual = visualInfo->visual;
            platform.depth = visualInfo->depth;
        }

        XFree(visualInfo);
    }

    // NOTE: Configs without a native visual use the default visual
    if (bestScore < 0) platform.eglConfig = configs[0];
    else if (transparent && (platform.depth != 32)) TRACELOG(LOG_WARNING, "DISPLAY: Failed to find a visual with alpha channel, window is not transparent");

    RL_FREE(configs);

    return true;
#endif
}

// Initialize EGL context and window surface
static bool InitGraphicsSurface(void)
{
#if defined(GRAPHICS_API_OPENGL_SOFTWARE)
    platform.gc = XCreateGC(platform.display, platform.window, 0, NULL);

    return CreateSoftwareImage(platform.bufferWidth, platform.bufferHeight);
#else
    EGLint contextAttribs[16] = { 0 };
    int attrib = 0;

    if (rlGetVersion() == RL_OPENGL_21)
    {
        contextAttribs[attrib++] = EGL_CONTEXT_MAJOR_VERSION;
        contextAttribs[attrib++] = 2;
        contextAttribs[attrib++] = EGL_CONTEXT_MINOR_VERSION;
        contextAttribs[attrib++] = 1;
    }
    else if (rlGetVersion() == RL_OPENGL_33)
    {
        contextAttribs[attrib++] = EGL_CONTEXT_MAJOR_VERSION;
        contextAttribs[attrib++] = 3;
        contextAttribs[attrib++] = EGL_CONTEXT_MINOR_VERSION;
        contextAttribs[attrib++] = 3;
        contextAttribs[attrib++] = EGL_CONTEXT_OPENGL_PROFILE_MASK;
        contextAttribs[attrib++] = EGL_CONTEXT_OPENGL_CORE_PROFILE_BIT;
    }
    else if (rlGetVersion() == RL_OPENGL_43)
    {
        contextAttribs[attrib++] = EGL_CONTEXT_MAJOR_VERSION;
        contextAttribs[attrib++] = 4;
        contextAttribs[attrib++] = EGL_CONTEXT_MINOR_VERSION;
        contextAttribs[attrib++] = 3;
        contextAttribs[attrib++] = EGL_CONTEXT_OPENGL_PROFILE_MASK;
        contextAttribs[attrib++] = EGL_CONTEXT_OPENGL_CORE_PROFILE_BIT;
#if RLGL_ENABLE_OPENGL_DEBUG_CONTEXT
        contextAttribs[attrib++] = EGL_CONTEXT_OPENGL_DEBUG;
        contextAttribs[attrib++] = EGL_TRUE;
#endif
    }
    else if (rlGetVersion() == RL_OPENGL_ES_20)
    {
        contextAttribs[attrib++] = EGL_CONTEXT_MAJOR_VERSION;
        contextAttribs[attrib++] = 2;
    }
    else if (rlGetVersion() == RL_OPENGL_ES_30)
    {
        contextAttribs[attrib++] = EGL_CONTEXT_MAJOR_VERSION;
        contextAttribs[attrib++] = 3;
    }
    contextAttribs[attrib++] = EGL_NONE;

    platform.eglContext = eglCreateContext(platform.eglDisplay, platform.eglConfig, EGL_NO_CONTEXT, contextAttribs);
    if (platform.eglContext == EGL_NO_CONTEXT)
    {
        TRACELOG(LOG_WARNING, "DISPLAY: Failed to create EGL context");
        return false;
    }

    // NOTE: Platform surface function takes a pointer to the window, legacy function takes the window
    EGLSurface (*createPlatformWindowSurface)(EGLDisplay, EGLConfig, void *, const EGLint *) = NULL;
    if (platform.eglPlatformDisplay)
    {
        createPlatformWindowSurface = (EGLSurface (*)(EGLDisplay, EGLConfig, void *, const EGLint *))eglGetProcAddress("eglCreatePlatformWindowSurfaceEXT");
    }

    if (createPlatformWindowSurface != NULL) platform.eglSurface = createPlatformWindowSurface(platform.eglDisplay, platform.eglConfig, &platform.window, NULL);
    else platform.eglSurface = eglCreateWindowSurface(platform.eglDisplay, platform.eglConfig, (EGLNativeWindowType)platform.window, NULL);
    if (platform.eglSurface == EGL_NO_SURFACE)
    {
        TRACELOG(LOG_WARNING, "DISPLAY: Failed to create EGL window surface");
        return false;
    }

    if (eglMakeCurrent(platform.eglDisplay, platform.eglSurface, platform.eglSurface, platform.eglContext) == EGL_FALSE)
    {
        TRACELOG(LOG_WARNING, "DISPLAY: Failed to attach EGL rendering context to EGL surface");
        return false;
    }

    if (FLAG_IS_SET(CORE.Window.flags, FLAG_VSYNC_HINT))
    {
        eglSwapInterval(platform.eglDisplay, 1);
        TRACELOG(LOG_INFO, "DISPLAY: Trying to enable VSYNC");
    }
    else eglSwapInterval(platform.eglDisplay, 0);

    return true;
#endif
}

// Create window and set its properties
static bool CreateWindow(void)
{
    Display *display = platform.display;
    Atoms *atoms = &platform.atoms;
    float scale = FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_HIGHDPI)? platform.scale : 1.0f;

    // NOTE: Window is created with the floating size, restored by the window manager when leaving fullscreen
    int width = (int)roundf((float)platform.windowedWidth*scale);
    int height = (int)roundf((float)platform.windowedHeight*scale);

    // Center window on the primary monitor
    platform.x = platform.monitors[0].x + (platform.monitors[0].width - width)/2;
    platform.y = platform.monitors[0].y + (platform.monitors[0].height - height)/2;
    if (platform.x < platform.monitors[0].x) platform.x = platform.monitors[0].x;
    if (platform.y < platform.monitors[0].y) platform.y = platform.monitors[0].y;
    platform.windowedX = platform.x;
    platform.windowedY = platform.y;
    CORE.Window.position.x = platform.x;
    CORE.Window.position.y = platform.y;

    platform.colormap = XCreateColormap(display, platform.root, platform.visual, AllocNone);

    platform.eventMask = StructureNotifyMask | KeyPressMask | KeyReleaseMask | PointerMotionMask | ButtonPressMask | ButtonReleaseMask |
        EnterWindowMask | LeaveWindowMask | FocusChangeMask | PropertyChangeMask | ExposureMask;

    XSetWindowAttributes attributes = { 0 };
    attributes.colormap = platform.colormap;
    attributes.border_pixel = 0;
    attributes.event_mask = platform.eventMask;

    platform.errorCode = 0;
    platform.window = XCreateWindow(display, platform.root, platform.x, platform.y, width, height, 0, platform.depth, InputOutput,
        platform.visual, CWBorderPixel | CWColormap | CWEventMask, &attributes);
    platform.parent = platform.root;

    XSync(display, False);
    if ((platform.window == None) || (platform.errorCode != 0))
    {
        platform.window = None;
        TRACELOG(LOG_WARNING, "X11: Failed to create window");
        return false;
    }

    // Window manager protocols
    Atom protocols[3] = { atoms->WM_DELETE_WINDOW, atoms->_NET_WM_PING, None };
    int protocolCount = 2;

    if (platform.syncAvailable)
    {
        XSyncValue value;
        XSyncIntToValue(&value, 0);
        platform.syncCounter = XSyncCreateCounter(display, value);
        XChangeProperty(display, platform.window, atoms->_NET_WM_SYNC_REQUEST_COUNTER, XA_CARDINAL, 32, PropModeReplace,
            (unsigned char *)&platform.syncCounter, 1);
        protocols[protocolCount++] = atoms->_NET_WM_SYNC_REQUEST;
    }

    XSetWMProtocols(display, platform.window, protocols, protocolCount);

    // Application identification
    // NOTE: _NET_WM_PID requires WM_CLIENT_MACHINE, both are used by window managers to end unresponsive applications
    unsigned long pid = (unsigned long)getpid();
    XChangeProperty(display, platform.window, atoms->_NET_WM_PID, XA_CARDINAL, 32, PropModeReplace, (unsigned char *)&pid, 1);

    char hostname[256] = { 0 };
    if (gethostname(hostname, sizeof(hostname) - 1) == 0)
    {
        char *hostnameList[1] = { hostname };
        XTextProperty hostnameProperty = { 0 };
        if (XStringListToTextProperty(hostnameList, 1, &hostnameProperty) != 0)
        {
            XSetWMClientMachine(display, platform.window, &hostnameProperty);
            XFree(hostnameProperty.value);
        }
    }

    const char *resourceName = getenv("RESOURCE_NAME");
    XClassHint classHint = { 0 };
    classHint.res_name = (char *)(((resourceName != NULL) && (resourceName[0] != '\0'))? resourceName : platform.appId);
    classHint.res_class = platform.appId;
    XSetClassHint(display, platform.window, &classHint);

    Atom windowType = atoms->_NET_WM_WINDOW_TYPE_NORMAL;
    XChangeProperty(display, platform.window, atoms->_NET_WM_WINDOW_TYPE, XA_ATOM, 32, PropModeReplace, (unsigned char *)&windowType, 1);

    const char *startupId = getenv("DESKTOP_STARTUP_ID");
    if ((startupId != NULL) && (startupId[0] != '\0'))
    {
        XChangeProperty(display, platform.window, atoms->_NET_STARTUP_ID, atoms->UTF8_STRING, 8, PropModeReplace,
            (const unsigned char *)startupId, (int)strlen(startupId));
    }

    // NOTE: Window manager does not focus a window mapped with user time 0
    if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_UNFOCUSED))
    {
        unsigned long userTime = 0;
        XChangeProperty(display, platform.window, atoms->_NET_WM_USER_TIME, XA_CARDINAL, 32, PropModeReplace, (unsigned char *)&userTime, 1);
    }

    // Drag and drop
    unsigned long dndVersion = XDND_VERSION;
    XChangeProperty(display, platform.window, atoms->XdndAware, XA_ATOM, 32, PropModeReplace, (unsigned char *)&dndVersion, 1);

    UpdateWindowTitle(CORE.Window.title);
    UpdateNormalHints();
    UpdateDecorations();

    if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_MOUSE_PASSTHROUGH))
    {
        if (platform.shapeAvailable) XShapeCombineRectangles(display, platform.window, ShapeInput, 0, 0, NULL, 0, ShapeSet, Unsorted);
        else
        {
            TRACELOG(LOG_WARNING, "X11: Mouse passthrough not supported by the X server");
            FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_MOUSE_PASSTHROUGH);
        }
    }

    // Transparent cursor, used to hide the cursor
    char cursorData[1] = { 0 };
    XColor cursorColor = { 0 };
    Pixmap cursorPixmap = XCreateBitmapFromData(display, platform.window, cursorData, 1, 1);
    platform.hiddenCursor = XCreatePixmapCursor(display, cursorPixmap, cursorPixmap, &cursorColor, &cursorColor, 0, 0);
    XFreePixmap(display, cursorPixmap);

    return true;
}

// Map window, optionally wait for the window manager to apply size and states (before the first frame)
// NOTE: Waiting is limited in time, the window manager could not map the window (e.g. another workspace) or not be running
static void ShowWindow(bool wait)
{
    Atoms *atoms = &platform.atoms;

    // Window states are applied by the window manager when the window is mapped
    // NOTE: States are removed by the window manager when the window is hidden, they are set on every map
    if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_MAXIMIZED))
    {
        SetWindowStateHint(atoms->_NET_WM_STATE_MAXIMIZED_VERT, atoms->_NET_WM_STATE_MAXIMIZED_HORZ, true);
    }

    if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_TOPMOST)) SetWindowStateHint(atoms->_NET_WM_STATE_ABOVE, None, true);

    if (FLAG_IS_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE) || FLAG_IS_SET(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE))
    {
        FLAG_SET(CORE.Window.flags, FLAG_WINDOW_UNDECORATED);
        RequestFullscreen(true, GetCurrentMonitor());
    }

    // NOTE: Initial state is read by the window manager when the window is mapped
    XWMHints hints = { 0 };
    hints.flags = StateHint | InputHint;
    hints.initial_state = FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_MINIMIZED)? IconicState : NormalState;
    hints.input = True;
    XSetWMHints(platform.display, platform.window, &hints);

    XMapWindow(platform.display, platform.window);
    XFlush(platform.display);
    platform.withdrawn = false;

    if (!wait || FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_MINIMIZED)) return;

    double endTime = GetTime() + (double)WINDOW_MAP_TIMEOUT/1000.0;

    // NOTE: Window size and state set by the window manager on map are applied before the first frame
    while (!platform.connectionLost && (!platform.mapped || platform.fullscreenPending))
    {
        int remaining = (int)((endTime - GetTime())*1000.0);
        if (remaining <= 0) break;

        ProcessEvents(remaining);
    }

    ProcessEvents(0);
}

// Notify the launcher that the application started (startup notification)
static void CompleteStartupNotification(void)
{
    const char *startupId = getenv("DESKTOP_STARTUP_ID");
    if ((startupId == NULL) || (startupId[0] == '\0')) return;

    // Message: remove: ID="<id>", quotes and backslashes of the id are escaped
    char message[1024] = { 0 };
    int length = snprintf(message, sizeof(message), "remove: ID=\"");
    for (int i = 0; (startupId[i] != '\0') && (length < (int)sizeof(message) - 4); i++)
    {
        if ((startupId[i] == '"') || (startupId[i] == '\\')) message[length++] = '\\';
        message[length++] = startupId[i];
    }
    message[length++] = '"';
    length++;   // NOTE: Message is sent with its terminating zero

    // Message is sent to the root window in chunks of 20 bytes
    XEvent event = { 0 };
    event.xclient.type = ClientMessage;
    event.xclient.window = platform.window;
    event.xclient.format = 8;
    event.xclient.message_type = platform.atoms._NET_STARTUP_INFO_BEGIN;

    for (int offset = 0; offset < length; offset += 20)
    {
        memset(event.xclient.data.b, 0, 20);
        memcpy(event.xclient.data.b, message + offset, ((length - offset) < 20)? (length - offset) : 20);
        XSendEvent(platform.display, platform.root, False, PropertyChangeMask, &event);
        event.xclient.message_type = platform.atoms._NET_STARTUP_INFO;
    }

    XFlush(platform.display);
    unsetenv("DESKTOP_STARTUP_ID");
}

// Read and process X11 events, timeout in milliseconds (-1 blocks)
// NOTE: Waiting ends when a window or input event is received
static void ProcessEvents(int timeout)
{
    if ((platform.display == NULL) || platform.connectionLost) return;

    platform.eventReceived = false;

    double endTime = GetTime() + (double)timeout/1000.0;

    while (!platform.connectionLost)
    {
        // NOTE: XPending() flushes the requests and reads the events available without blocking
        while (!platform.connectionLost && (XPending(platform.display) > 0))
        {
            XEvent event = { 0 };
            XNextEvent(platform.display, &event);
            ProcessEvent(&event);
        }

        if (platform.eventReceived || (timeout == 0)) break;

        int remaining = -1;
        if (timeout > 0)
        {
            remaining = (int)((endTime - GetTime())*1000.0);
            if (remaining <= 0) break;
        }

        if (!WaitForEvents(remaining)) break;
    }
}

// Wait for data on the X11 connection, timeout in milliseconds (-1 blocks)
// NOTE: Returns false if the waiting timed out or failed
static bool WaitForEvents(int timeout)
{
    struct pollfd fds[1] = { { ConnectionNumber(platform.display), POLLIN, 0 } };
    int result = 0;

    do
    {
        result = poll(fds, 1, timeout);
    } while ((result < 0) && (errno == EINTR) && (timeout < 0));

    return (result > 0);
}

// Event predicate: selection notifications and property changes for a window and atom
static Bool MatchEvent(Display *display, XEvent *event, XPointer data)
{
    const EventFilter *filter = (const EventFilter *)data;

    if (event->type != filter->type) return False;

    if (event->type == SelectionNotify) return ((event->xselection.requestor == filter->window) && (event->xselection.selection == filter->atom));
    if (event->type == PropertyNotify)
    {
        return ((event->xproperty.window == filter->window) && (event->xproperty.atom == filter->atom) && (event->xproperty.state == PropertyNewValue));
    }

    return (event->xany.window == filter->window);
}

// Wait for an event of a given type, window and atom, other events are kept queued
// NOTE: Waiting is limited in time (milliseconds), other clients could be unresponsive
static bool WaitForEvent(int type, Window window, Atom atom, XEvent *event, int timeout)
{
    EventFilter filter = { type, window, atom };
    double endTime = GetTime() + (double)timeout/1000.0;

    while (!platform.connectionLost)
    {
        if (XCheckIfEvent(platform.display, event, MatchEvent, (XPointer)&filter)) return true;

        int remaining = (int)((endTime - GetTime())*1000.0);
        if (remaining <= 0) break;

        XFlush(platform.display);
        WaitForEvents(remaining);
    }

    return false;
}

// Get current X server time
// NOTE: Selection requests require a time not older than the last selection change, events of the window can be older;
// the time is taken from the notification of a property change
static Time GetServerTime(void)
{
    Atom property = platform.atoms.RAYLIB_TIMESTAMP;

    XChangeProperty(platform.display, platform.window, property, XA_INTEGER, 32, PropModeAppend, NULL, 0);

    XEvent event = { 0 };
    if (!WaitForEvent(PropertyNotify, platform.window, property, &event, SELECTION_TIMEOUT)) return CurrentTime;

    return event.xproperty.time;
}

// Process X11 event
static void ProcessEvent(XEvent *event)
{
    // NOTE: Key code is saved before filtering, some input methods clear it
    unsigned int keycode = ((event->type == KeyPress) || (event->type == KeyRelease))? event->xkey.keycode : 0;

    // Input method events (compose sequences, text composition)
    bool filtered = (XFilterEvent(event, None) == True);

    if (platform.xkbAvailable && (event->type == platform.xkbEventBase))
    {
        XkbEvent *xkbEvent = (XkbEvent *)event;

        if (xkbEvent->any.xkb_type == XkbStateNotify)
        {
            platform.xkbGroup = xkbEvent->state.group;
            platform.capsLocked = FLAG_IS_SET(xkbEvent->state.locked_mods, LockMask);
            platform.numLocked = ((platform.numLockMask != 0) && FLAG_IS_SET(xkbEvent->state.locked_mods, platform.numLockMask));

            // NOTE: CAPS/NUM lock keys are considered down while the lock is active
            if (!platform.capsLockDown) CORE.Input.Keyboard.currentKeyState[KEY_CAPS_LOCK] = platform.capsLocked;
            if (!platform.numLockDown) CORE.Input.Keyboard.currentKeyState[KEY_NUM_LOCK] = platform.numLocked;
        }
        else if ((xkbEvent->any.xkb_type == XkbNewKeyboardNotify) || (xkbEvent->any.xkb_type == XkbMapNotify))
        {
            if (xkbEvent->any.xkb_type == XkbMapNotify) XkbRefreshKeyboardMapping(&xkbEvent->map);
            UpdateKeymap();
        }

        return;
    }

    if (platform.randrAvailable && ((event->type == platform.randrEventBase + RRScreenChangeNotify) || (event->type == platform.randrEventBase + RRNotify)))
    {
        XRRUpdateConfiguration(event);
        UpdateMonitors();
        platform.eventReceived = true;
        return;
    }

    if (event->type == GenericEvent)
    {
        ProcessGenericEvent(&event->xcookie);
        return;
    }

    if (filtered && (event->type != KeyPress) && (event->type != KeyRelease)) return;

    switch (event->type)
    {
        case KeyPress:
        {
            event->xkey.keycode = keycode;
            HandleKey(&event->xkey, true, filtered);
        } break;
        case KeyRelease:
        {
            event->xkey.keycode = keycode;
            HandleKey(&event->xkey, false, filtered);
        } break;
        case ButtonPress: HandlePointerButton(event->xbutton.button, true, event->xbutton.time); break;
        case ButtonRelease: HandlePointerButton(event->xbutton.button, false, event->xbutton.time); break;
        case MotionNotify: HandlePointerMotion((double)event->xmotion.x, (double)event->xmotion.y); break;
        case EnterNotify:
        {
            // NOTE: Scroll axes values changed while the pointer was on other windows
            if (platform.xiAvailable) UpdateInputDevices();
            if (event->xcrossing.mode == NotifyNormal) HandlePointerMotion((double)event->xcrossing.x, (double)event->xcrossing.y);

            // NOTE: Crossing events caused by pointer grabs are ignored (window manager shortcuts, menus)
            if ((event->xcrossing.mode != NotifyGrab) && (event->xcrossing.detail != NotifyInferior)) CORE.Input.Mouse.cursorOnScreen = true;
        } break;
        case LeaveNotify:
        {
            if ((event->xcrossing.mode == NotifyNormal) && (event->xcrossing.detail != NotifyInferior)) CORE.Input.Mouse.cursorOnScreen = false;
        } break;
        case FocusIn:
        case FocusOut:
        {
            if ((event->xfocus.detail == NotifyPointer) || (event->xfocus.detail == NotifyInferior)) break;

            // NOTE: Focus changes caused by pointer or keyboard grabs are ignored (window dragging, window manager shortcuts),
            // keys released while another client held the grab are checked
            if ((event->xfocus.mode == NotifyGrab) || (event->xfocus.mode == NotifyUngrab))
            {
                if ((event->type == FocusIn) && platform.focused) SyncInputState();
                break;
            }

            HandleFocus(event->type == FocusIn);
        } break;
        case MapNotify:
        {
            if (event->xmap.window != platform.window) break;

            platform.mapped = true;
            platform.eventReceived = true;
            FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_MINIMIZED);

            // NOTE: Input shape is set again on map, some window managers reset it
            if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_MOUSE_PASSTHROUGH))
            {
                XShapeCombineRectangles(platform.display, platform.window, ShapeInput, 0, 0, NULL, 0, ShapeSet, Unsorted);
            }
        } break;
        case UnmapNotify:
        {
            if (event->xunmap.window != platform.window) break;

            platform.mapped = false;
            platform.eventReceived = true;
            LockPointer(false);
        } break;
        case ReparentNotify:
        {
            if (event->xreparent.window == platform.window) platform.parent = event->xreparent.parent;
        } break;
        case DestroyNotify:
        {
            // Window manager exit
            if (event->xdestroywindow.window == platform.wmCheckWindow) UpdateWMSupported();

            // Requestor of an incremental selection transfer destroyed
            for (int i = 0; i < MAX_INCR_TRANSFERS; i++)
            {
                if ((platform.incrTransfers[i].requestor != None) && (platform.incrTransfers[i].requestor == event->xdestroywindow.window))
                {
                    RL_FREE(platform.incrTransfers[i].data);
                    memset(&platform.incrTransfers[i], 0, sizeof(IncrTransfer));
                }
            }
        } break;
        case ConfigureNotify:
        {
            if (event->xconfigure.window != platform.window) break;

            platform.eventReceived = true;

            // NOTE: Size events generated before the last resize request are outdated
            if ((long)(event->xconfigure.serial - platform.resizeSerial) >= 0) ApplyWindowSize(event->xconfigure.width, event->xconfigure.height);

            // NOTE: Position of events sent by the server is relative to the parent window (window manager frame),
            // events sent by the window manager are relative to the root window
            int x = event->xconfigure.x;
            int y = event->xconfigure.y;

            if (!event->xconfigure.send_event && (platform.parent != platform.root))
            {
                Window child = None;
                int rootX = 0;
                int rootY = 0;
                if (XTranslateCoordinates(platform.display, platform.parent, platform.root, x, y, &rootX, &rootY, &child))
                {
                    x = rootX;
                    y = rootY;
                }
            }

            platform.x = x;
            platform.y = y;
            CORE.Window.position.x = x;
            CORE.Window.position.y = y;

            if (!platform.fullscreen && !platform.maximized)
            {
                platform.windowedX = x;
                platform.windowedY = y;
            }
        } break;
        case PropertyNotify: ProcessPropertyChange(&event->xproperty); break;
        case ClientMessage: ProcessClientMessage(&event->xclient); break;
        case SelectionRequest: HandleSelectionRequest(&event->xselectionrequest); break;
        case SelectionClear:
        {
            // NOTE: Selection lost before it was acquired again is ignored (server time is a 32 bit value that wraps around)
            if ((event->xselectionclear.selection == platform.atoms.CLIPBOARD) && ((platform.clipboardTime == CurrentTime) ||
                (((unsigned int)event->xselectionclear.time - (unsigned int)platform.clipboardTime) < 0x80000000u))) platform.clipboardOwned = false;
        } break;
        case SelectionNotify:
        {
            // Dropped data converted by the drag and drop source
            if (event->xselection.selection == platform.atoms.XdndSelection) HandleDragData(&event->xselection);
        } break;
        case MappingNotify:
        {
            XRefreshKeyboardMapping(&event->xmapping);
            if (event->xmapping.request != MappingPointer) UpdateKeymap();
        } break;
        default: break;
    }
}

// Process window manager and drag and drop messages
static void ProcessClientMessage(XClientMessageEvent *event)
{
    Atoms *atoms = &platform.atoms;

    if (event->message_type == atoms->WM_PROTOCOLS)
    {
        Atom protocol = (Atom)event->data.l[0];

        if (protocol == atoms->WM_DELETE_WINDOW)
        {
            platform.closeRequested = true;
            platform.eventReceived = true;
        }
        else if (protocol == atoms->_NET_WM_PING)
        {
            // Reply to the window manager, the application is responsive
            XEvent reply = { 0 };
            reply.xclient = *event;
            reply.xclient.window = platform.root;
            XSendEvent(platform.display, platform.root, False, SubstructureNotifyMask | SubstructureRedirectMask, &reply);
        }
        else if (protocol == atoms->_NET_WM_SYNC_REQUEST)
        {
            // NOTE: Counter is set after the frame with the new window size is presented
            XSyncIntsToValue(&platform.syncValue, (unsigned int)event->data.l[2], (int)event->data.l[3]);
            platform.syncPending = true;
        }
    }
    else if (event->message_type == atoms->XdndEnter) HandleDragEnter(event);
    else if (event->message_type == atoms->XdndPosition) HandleDragPosition(event);
    else if (event->message_type == atoms->XdndDrop) HandleDragDrop(event);
    else if (event->message_type == atoms->XdndLeave)
    {
        platform.dragSource = None;
        platform.dragAccepted = false;
    }
}

// Process window and root properties changes
static void ProcessPropertyChange(XPropertyEvent *event)
{
    Atoms *atoms = &platform.atoms;

    if (event->window == platform.root)
    {
        if (event->atom == atoms->RESOURCE_MANAGER) UpdateWindowScale();
        else if ((event->atom == atoms->_NET_SUPPORTED) || (event->atom == atoms->_NET_SUPPORTING_WM_CHECK)) UpdateWMSupported();
    }
    else if (event->window == platform.window)
    {
        if ((event->atom == atoms->_NET_WM_STATE) || (event->atom == atoms->WM_STATE)) UpdateWindowStateFlags();
    }

    if (event->state == PropertyDelete)
    {
        // Incremental selection transfer: the requestor is ready for the next chunk
        for (int i = 0; i < MAX_INCR_TRANSFERS; i++)
        {
            IncrTransfer *transfer = &platform.incrTransfers[i];
            if ((transfer->requestor == event->window) && (transfer->property == event->atom))
            {
                ContinueIncrTransfer(transfer);
                break;
            }
        }
    }
}

// Get window property, returns items count
// NOTE: Property value must be released with XFree(), items of 32 bits are stored as long
static unsigned long GetWindowProperty(Window window, Atom property, Atom type, unsigned char **value)
{
    Atom actualType = None;
    int actualFormat = 0;
    unsigned long itemCount = 0;
    unsigned long bytesAfter = 0;

    *value = NULL;
    if (XGetWindowProperty(platform.display, window, property, 0, LONG_MAX, False, type, &actualType, &actualFormat,
        &itemCount, &bytesAfter, value) != Success) return 0;

    // NOTE: Xlib allocates a value also for empty properties, properties other than text have items of 32 bits
    if ((*value != NULL) && ((itemCount == 0) || ((type != XA_STRING) && (actualFormat != 32))))
    {
        XFree(*value);
        *value = NULL;
    }

    if (*value == NULL) return 0;

    return itemCount;
}

// Send client message to the root window (request to the window manager)
static void SendWindowMessage(Atom type, long data0, long data1, long data2, long data3, long data4)
{
    XEvent event = { 0 };
    event.xclient.type = ClientMessage;
    event.xclient.window = platform.window;
    event.xclient.format = 32;
    event.xclient.message_type = type;
    event.xclient.data.l[0] = data0;
    event.xclient.data.l[1] = data1;
    event.xclient.data.l[2] = data2;
    event.xclient.data.l[3] = data3;
    event.xclient.data.l[4] = data4;

    XSendEvent(platform.display, platform.root, False, SubstructureNotifyMask | SubstructureRedirectMask, &event);
}

// Update hints supported by the window manager (_NET_SUPPORTED)
// NOTE: Hints are only valid while the window manager check window exists (EWMH compliant window manager running)
static void UpdateWMSupported(void)
{
    RL_FREE(platform.wmSupported);
    platform.wmSupported = NULL;
    platform.wmSupportedCount = 0;
    platform.wmCheckWindow = None;

    unsigned char *value = NULL;
    if (GetWindowProperty(platform.root, platform.atoms._NET_SUPPORTING_WM_CHECK, XA_WINDOW, &value) == 0) return;

    Window checkWindow = (Window)(*(unsigned long *)value);
    XFree(value);

    // The check window has the same property pointing to itself
    if (GetWindowProperty(checkWindow, platform.atoms._NET_SUPPORTING_WM_CHECK, XA_WINDOW, &value) == 0) return;

    bool valid = ((Window)(*(unsigned long *)value) == checkWindow);
    XFree(value);
    if (!valid) return;

    // Get notified when the window manager exits
    platform.wmCheckWindow = checkWindow;
    XSelectInput(platform.display, checkWindow, StructureNotifyMask);

    unsigned long count = GetWindowProperty(platform.root, platform.atoms._NET_SUPPORTED, XA_ATOM, &value);
    if (count > 0)
    {
        platform.wmSupported = (Atom *)RL_CALLOC(count, sizeof(Atom));
        memcpy(platform.wmSupported, value, count*sizeof(Atom));
        platform.wmSupportedCount = (int)count;
        XFree(value);
    }
}

// Check hint supported by the window manager
static bool IsWMHintSupported(Atom hint)
{
    for (int i = 0; i < platform.wmSupportedCount; i++)
    {
        if (platform.wmSupported[i] == hint) return true;
    }

    return false;
}

// Request window state change (_NET_WM_STATE), one or two states
// NOTE: States of a managed window (shown or minimized) are changed by the window manager on request,
// states of a window not shown are read by the window manager from the property when the window is mapped
static void SetWindowStateHint(Atom state1, Atom state2, bool enabled)
{
    if (!platform.withdrawn)
    {
        // NOTE: Source indication 1 (application)
        SendWindowMessage(platform.atoms._NET_WM_STATE, enabled? NET_WM_STATE_ADD : NET_WM_STATE_REMOVE, (long)state1, (long)state2, 1, 0);
        return;
    }

    unsigned char *value = NULL;
    unsigned long count = GetWindowProperty(platform.window, platform.atoms._NET_WM_STATE, XA_ATOM, &value);
    Atom *current = (Atom *)value;

    Atom *states = (Atom *)RL_CALLOC(count + 2, sizeof(Atom));
    int stateCount = 0;

    for (unsigned long i = 0; i < count; i++)
    {
        if ((current[i] != state1) && (current[i] != state2)) states[stateCount++] = current[i];
    }

    if (enabled)
    {
        states[stateCount++] = state1;
        if (state2 != None) states[stateCount++] = state2;
    }

    XChangeProperty(platform.display, platform.window, platform.atoms._NET_WM_STATE, XA_ATOM, 32, PropModeReplace,
        (unsigned char *)states, stateCount);

    if (value != NULL) XFree(value);
    RL_FREE(states);
}

// Update window flags from window manager states
static void UpdateWindowStateFlags(void)
{
    Atoms *atoms = &platform.atoms;

    // NOTE: States are removed by the window manager while the window is hidden, flags keep the requested states
    if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_HIDDEN)) return;

    platform.eventReceived = true;

    bool floating = (!platform.fullscreen && !platform.maximized);
    bool fullscreen = false;
    bool maximizedVert = false;
    bool maximizedHorz = false;
    bool minimized = false;

    unsigned char *value = NULL;
    unsigned long count = GetWindowProperty(platform.window, atoms->_NET_WM_STATE, XA_ATOM, &value);
    Atom *states = (Atom *)value;

    for (unsigned long i = 0; i < count; i++)
    {
        if (states[i] == atoms->_NET_WM_STATE_FULLSCREEN) fullscreen = true;
        else if (states[i] == atoms->_NET_WM_STATE_MAXIMIZED_VERT) maximizedVert = true;
        else if (states[i] == atoms->_NET_WM_STATE_MAXIMIZED_HORZ) maximizedHorz = true;
        else if (states[i] == atoms->_NET_WM_STATE_HIDDEN) minimized = true;
    }

    if (value != NULL) XFree(value);

    // ICCCM window state: WithdrawnState, NormalState or IconicState
    if (GetWindowProperty(platform.window, atoms->WM_STATE, atoms->WM_STATE, &value) > 0)
    {
        if (*(unsigned long *)value == IconicState) minimized = true;
        XFree(value);
    }

    if (minimized) FLAG_SET(CORE.Window.flags, FLAG_WINDOW_MINIMIZED);
    else FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_MINIMIZED);

    platform.maximized = (maximizedVert && maximizedHorz);
    if (platform.maximized) FLAG_SET(CORE.Window.flags, FLAG_WINDOW_MAXIMIZED);
    else FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_MAXIMIZED);

    if (platform.fullscreenManual) return;

    // Fullscreen state changed by the window manager
    // NOTE: States notified before the last fullscreen request was applied do not change the flags,
    // a request is not waited for longer than FULLSCREEN_REQUEST_TIMEOUT (the window manager can ignore it)
    if (platform.fullscreenPending && ((fullscreen == platform.fullscreenRequested) ||
        ((GetTime() - platform.fullscreenRequestTime) > FULLSCREEN_REQUEST_TIMEOUT))) platform.fullscreenPending = false;

    bool fullscreenFlags = (FLAG_IS_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE) || FLAG_IS_SET(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE));
    if (!platform.fullscreenPending && (fullscreen != fullscreenFlags))
    {
        if (fullscreen)
        {
            FLAG_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE);
            FLAG_SET(CORE.Window.flags, FLAG_WINDOW_UNDECORATED);
        }
        else
        {
            if (platform.decorationRequested) FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_UNDECORATED);
            FLAG_CLEAR(CORE.Window.flags, FLAG_FULLSCREEN_MODE);
            FLAG_CLEAR(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE);
        }

        platform.fullscreenRequested = fullscreen;
        UpdateNormalHints();
    }

    platform.fullscreen = fullscreen;

    // NOTE: Position restored by the window manager can be on a monitor that is not connected anymore
    if (!floating && !platform.fullscreen && !platform.maximized) KeepWindowOnMonitors();
}

// Update window size limits (WM_NORMAL_HINTS)
// NOTE: Non-resizable windows have the same minimum and maximum size, fullscreen windows have no limits
static void UpdateNormalHints(void)
{
    float scale = FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_HIGHDPI)? platform.scale : 1.0f;

    XSizeHints *hints = XAllocSizeHints();

    // NOTE: Window position is the position of the window content, the window manager frame is placed around it
    hints->flags = PPosition | PWinGravity;
    hints->win_gravity = StaticGravity;

    if (!platform.fullscreenRequested)
    {
        if (!FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_RESIZABLE))
        {
            hints->flags |= (PMinSize | PMaxSize);
            hints->min_width = hints->max_width = (int)roundf((float)platform.windowedWidth*scale);
            hints->min_height = hints->max_height = (int)roundf((float)platform.windowedHeight*scale);
        }
        else
        {
            if ((CORE.Window.screenMin.width > 0) && (CORE.Window.screenMin.height > 0))
            {
                hints->flags |= PMinSize;
                hints->min_width = (int)roundf((float)CORE.Window.screenMin.width*scale);
                hints->min_height = (int)roundf((float)CORE.Window.screenMin.height*scale);
            }

            if ((CORE.Window.screenMax.width > 0) && (CORE.Window.screenMax.height > 0))
            {
                hints->flags |= PMaxSize;
                hints->max_width = (int)roundf((float)CORE.Window.screenMax.width*scale);
                hints->max_height = (int)roundf((float)CORE.Window.screenMax.height*scale);
            }
        }
    }

    XSetWMNormalHints(platform.display, platform.window, hints);
    XFree(hints);
}

// Update window decorations (_MOTIF_WM_HINTS)
static void UpdateDecorations(void)
{
    // Motif window manager hints: flags, functions, decorations, input mode and status
    unsigned long hints[5] = { MWM_HINTS_DECORATIONS, 0, platform.decorationRequested? MWM_DECOR_ALL : 0, 0, 0 };

    XChangeProperty(platform.display, platform.window, platform.atoms._MOTIF_WM_HINTS, platform.atoms._MOTIF_WM_HINTS, 32, PropModeReplace,
        (unsigned char *)hints, 5);
}

// Update window title properties
static void UpdateWindowTitle(const char *title)
{
    if (title == NULL) title = "";
    int length = (int)strlen(title);

    // NOTE: Legacy properties (WM_NAME, WM_ICON_NAME) are also set as UTF-8 text, for window managers without EWMH support
    XChangeProperty(platform.display, platform.window, platform.atoms._NET_WM_NAME, platform.atoms.UTF8_STRING, 8, PropModeReplace,
        (const unsigned char *)title, length);
    XChangeProperty(platform.display, platform.window, platform.atoms._NET_WM_ICON_NAME, platform.atoms.UTF8_STRING, 8, PropModeReplace,
        (const unsigned char *)title, length);
    XChangeProperty(platform.display, platform.window, XA_WM_NAME, platform.atoms.UTF8_STRING, 8, PropModeReplace,
        (const unsigned char *)title, length);
    XChangeProperty(platform.display, platform.window, XA_WM_ICON_NAME, platform.atoms.UTF8_STRING, 8, PropModeReplace,
        (const unsigned char *)title, length);
}

// Update scale factor from the resource manager (Xft.dpi)
// NOTE: Resource manager property of the root window is read, it changes when the desktop scale is changed
static void UpdateWindowScale(void)
{
    float scale = 1.0f;

    unsigned char *value = NULL;
    if (GetWindowProperty(platform.root, platform.atoms.RESOURCE_MANAGER, XA_STRING, &value) > 0)
    {
        XrmDatabase database = XrmGetStringDatabase((const char *)value);
        if (database != NULL)
        {
            char *type = NULL;
            XrmValue dpi = { 0 };
            if (XrmGetResource(database, "Xft.dpi", "Xft.Dpi", &type, &dpi) && (type != NULL) && (strcmp(type, "String") == 0) && (dpi.addr != NULL))
            {
                float dpiValue = (float)atof(dpi.addr);
                if (dpiValue > 0.0f) scale = dpiValue/DPI_DEFAULT;
            }

            XrmDestroyDatabase(database);
        }
    }

    if (value != NULL) XFree(value);

    if (scale == platform.scale) return;

    platform.scale = scale;
    platform.eventReceived = true;

    if ((platform.window == None) || !FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_HIGHDPI)) return;

    if (!platform.fullscreen && !platform.maximized)
    {
        // Keep the window logical size
        platform.bufferWidth = (int)roundf((float)platform.width*scale);
        platform.bufferHeight = (int)roundf((float)platform.height*scale);

        UpdateNormalHints();
        platform.resizeSerial = NextRequest(platform.display);
        XResizeWindow(platform.display, platform.window, platform.bufferWidth, platform.bufferHeight);
    }
    else
    {
        platform.width = (int)roundf((float)platform.bufferWidth/scale);
        platform.height = (int)roundf((float)platform.bufferHeight/scale);
    }

    UpdateWindowBuffers();
}

// Apply window size and scale to the framebuffer
static void UpdateWindowBuffers(void)
{
    float scale = FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_HIGHDPI)? platform.scale : 1.0f;

    if (platform.width < 1) platform.width = 1;
    if (platform.height < 1) platform.height = 1;
    if (platform.bufferWidth < 1) platform.bufferWidth = 1;
    if (platform.bufferHeight < 1) platform.bufferHeight = 1;

    if (FLAG_IS_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE) || FLAG_IS_SET(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE))
    {
        int monitor = GetCurrentMonitor();
        CORE.Window.display.width = platform.monitors[monitor].width;
        CORE.Window.display.height = platform.monitors[monitor].height;
    }

    // Set screen size to logical pixel size and render size to framebuffer size
    bool resized = (((int)CORE.Window.render.width != platform.bufferWidth) || ((int)CORE.Window.render.height != platform.bufferHeight) ||
        ((int)CORE.Window.screen.width != platform.width) || ((int)CORE.Window.screen.height != platform.height));

    CORE.Window.screen.width = platform.width;
    CORE.Window.screen.height = platform.height;
    CORE.Window.screenScale = MatrixScale(scale, scale, 1.0f);

    // Reset viewport and projection matrix for new size
    // NOTE: Stores current render size: CORE.Window.render
    if (platform.initialized && resized) SetupViewport(platform.bufferWidth, platform.bufferHeight);

    CORE.Window.render.width = platform.bufferWidth;
    CORE.Window.render.height = platform.bufferHeight;
    CORE.Window.currentFbo = CORE.Window.render;
    if (resized && platform.initialized) CORE.Window.resizedLastFrame = true;
}

// Apply window size reported by the server (pixels)
static void ApplyWindowSize(int width, int height)
{
    if ((width == platform.bufferWidth) && (height == platform.bufferHeight)) return;

    float scale = FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_HIGHDPI)? platform.scale : 1.0f;

    platform.bufferWidth = width;
    platform.bufferHeight = height;
    platform.width = (int)roundf((float)width/scale);
    platform.height = (int)roundf((float)height/scale);

    if (!platform.fullscreen && !platform.maximized && !platform.fullscreenRequested)
    {
        platform.windowedWidth = platform.width;
        platform.windowedHeight = platform.height;
    }

    UpdateWindowBuffers();
}

// Request fullscreen state to the window manager
// NOTE: The window covers the monitor, the video mode is not changed
static void RequestFullscreen(bool fullscreen, int monitor)
{
    Atoms *atoms = &platform.atoms;
    float scale = FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_HIGHDPI)? platform.scale : 1.0f;

    if ((monitor < 0) || (monitor >= platform.monitorCount)) monitor = 0;
    const MonitorInfo *info = &platform.monitors[monitor];

    platform.fullscreenRequested = fullscreen;
    UpdateNormalHints();

    if (fullscreen)
    {
        if (IsWMHintSupported(atoms->_NET_WM_STATE_FULLSCREEN))
        {
            // Select the monitor, the window manager uses the monitor the window is on by default
            if (platform.monitorCount > 1)
            {
                if (!platform.withdrawn && IsWMHintSupported(atoms->_NET_WM_FULLSCREEN_MONITORS))
                {
                    // NOTE: Monitors edges: top, bottom, left, right; source indication 1 (application)
                    SendWindowMessage(atoms->_NET_WM_FULLSCREEN_MONITORS, info->index, info->index, info->index, info->index, 1);
                }
                else if (GetCurrentMonitor() != monitor)
                {
                    // NOTE: A fullscreen window can not be moved, it leaves fullscreen to change the monitor
                    int width = (int)roundf((float)platform.windowedWidth*scale);
                    int height = (int)roundf((float)platform.windowedHeight*scale);

                    if (platform.fullscreen) SetWindowStateHint(atoms->_NET_WM_STATE_FULLSCREEN, None, false);
                    XMoveWindow(platform.display, platform.window, info->x + (info->width - width)/2, info->y + (info->height - height)/2);
                    platform.fullscreen = false;
                }
            }

            SetWindowStateHint(atoms->_NET_WM_STATE_FULLSCREEN, None, true);
            if (!platform.fullscreen) platform.fullscreenPending = true;
            platform.fullscreenRequestTime = GetTime();
        }
        else
        {
            // Window manager without fullscreen support (or not running): cover the monitor with the window
            platform.fullscreenManual = true;
            platform.fullscreen = true;
            XMoveResizeWindow(platform.display, platform.window, info->x, info->y, info->width, info->height);
            XRaiseWindow(platform.display, platform.window);
        }
    }
    else
    {
        if (platform.fullscreenManual)
        {
            platform.fullscreenManual = false;
            platform.fullscreen = false;
            XMoveResizeWindow(platform.display, platform.window, platform.windowedX, platform.windowedY,
                (int)roundf((float)platform.windowedWidth*scale), (int)roundf((float)platform.windowedHeight*scale));
        }
        else
        {
            // NOTE: Window position and size are restored by the window manager
            SetWindowStateHint(atoms->_NET_WM_STATE_FULLSCREEN, None, false);
            XDeleteProperty(platform.display, platform.window, atoms->_NET_WM_FULLSCREEN_MONITORS);
            if (platform.fullscreen) platform.fullscreenPending = true;
            platform.fullscreenRequestTime = GetTime();
        }
    }

    XFlush(platform.display);
}

// Move floating window to the main monitor if it is out of all monitors
static void KeepWindowOnMonitors(void)
{
    if (platform.monitorCount == 0) return;

    float scale = FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_HIGHDPI)? platform.scale : 1.0f;
    int width = (int)roundf((float)platform.windowedWidth*scale);
    int height = (int)roundf((float)platform.windowedHeight*scale);

    for (int i = 0; i < platform.monitorCount; i++)
    {
        const MonitorInfo *info = &platform.monitors[i];

        if ((platform.windowedX < (info->x + info->width)) && ((platform.windowedX + width) > info->x) &&
            (platform.windowedY < (info->y + info->height)) && ((platform.windowedY + height) > info->y)) return;
    }

    const MonitorInfo *info = &platform.monitors[0];
    XMoveWindow(platform.display, platform.window, info->x + (info->width - width)/2, info->y + (info->height - height)/2);
}

// Get refresh rate of the mode of a RandR output, in Hz
static float GetOutputRefreshRate(XRRScreenResources *resources, RROutput output)
{
    float rate = 0.0f;

    XRROutputInfo *outputInfo = XRRGetOutputInfo(platform.display, resources, output);
    if (outputInfo == NULL) return rate;

    if (outputInfo->crtc != None)
    {
        XRRCrtcInfo *crtcInfo = XRRGetCrtcInfo(platform.display, resources, outputInfo->crtc);
        if (crtcInfo != NULL)
        {
            for (int i = 0; i < resources->nmode; i++)
            {
                const XRRModeInfo *mode = &resources->modes[i];
                if (mode->id != crtcInfo->mode) continue;

                double verticalTotal = (double)mode->vTotal;
                if (FLAG_IS_SET(mode->modeFlags, RR_DoubleScan)) verticalTotal *= 2.0;
                if (FLAG_IS_SET(mode->modeFlags, RR_Interlace)) verticalTotal /= 2.0;

                if ((mode->hTotal > 0) && (verticalTotal > 0.0)) rate = (float)((double)mode->dotClock/((double)mode->hTotal*verticalTotal));
                break;
            }

            XRRFreeCrtcInfo(crtcInfo);
        }
    }

    XRRFreeOutputInfo(outputInfo);

    return rate;
}

// Update connected monitors
// NOTE: Primary monitor is the first one, the whole screen is used when RandR monitors are not available
static void UpdateMonitors(void)
{
    platform.monitorCount = 0;

    if (platform.randrMonitors)
    {
        int count = 0;
        XRRMonitorInfo *infos = XRRGetMonitors(platform.display, platform.root, True, &count);
        XRRScreenResources *resources = XRRGetScreenResourcesCurrent(platform.display, platform.root);

        for (int i = 0; (infos != NULL) && (i < count) && (platform.monitorCount < MAX_MONITORS); i++)
        {
            if ((infos[i].width <= 0) || (infos[i].height <= 0)) continue;

            // Outputs showing the same area are one monitor (mirrored outputs)
            bool mirrored = false;
            for (int j = 0; j < platform.monitorCount; j++)
            {
                const MonitorInfo *other = &platform.monitors[j];
                if ((other->x == infos[i].x) && (other->y == infos[i].y) && (other->width == infos[i].width) && (other->height == infos[i].height)) mirrored = true;
            }

            if (mirrored) continue;

            MonitorInfo *info = &platform.monitors[platform.monitorCount];
            memset(info, 0, sizeof(MonitorInfo));

            char *name = XGetAtomName(platform.display, infos[i].name);
            if (name != NULL)
            {
                snprintf(info->name, sizeof(info->name), "%s", name);
                XFree(name);
            }

            info->index = i;
            info->x = infos[i].x;
            info->y = infos[i].y;
            info->width = infos[i].width;
            info->height = infos[i].height;
            info->physicalWidth = infos[i].mwidth;
            info->physicalHeight = infos[i].mheight;
            if ((resources != NULL) && (infos[i].noutput > 0)) info->refreshRate = GetOutputRefreshRate(resources, infos[i].outputs[0]);

            platform.monitorCount++;
        }

        if (resources != NULL) XRRFreeScreenResources(resources);
        if (infos != NULL) XRRFreeMonitors(infos);
    }

    if (platform.monitorCount == 0)
    {
        MonitorInfo *info = &platform.monitors[0];
        memset(info, 0, sizeof(MonitorInfo));

        snprintf(info->name, sizeof(info->name), "Screen %i", platform.screen);
        info->width = DisplayWidth(platform.display, platform.screen);
        info->height = DisplayHeight(platform.display, platform.screen);
        info->physicalWidth = DisplayWidthMM(platform.display, platform.screen);
        info->physicalHeight = DisplayHeightMM(platform.display, platform.screen);

        platform.monitorCount = 1;
    }
}

// Get monitor with the largest overlap with a rectangle, primary (first) monitor if there is no overlap
static int GetMonitorAt(int x, int y, int width, int height)
{
    int index = 0;
    long bestArea = 0;

    for (int i = 0; i < platform.monitorCount; i++)
    {
        const MonitorInfo *info = &platform.monitors[i];

        int left = (x > info->x)? x : info->x;
        int top = (y > info->y)? y : info->y;
        int right = ((x + width) < (info->x + info->width))? (x + width) : (info->x + info->width);
        int bottom = ((y + height) < (info->y + info->height))? (y + height) : (info->y + info->height);

        if ((right <= left) || (bottom <= top)) continue;

        long area = (long)(right - left)*(long)(bottom - top);
        if (area > bestArea)
        {
            bestArea = area;
            index = i;
        }
    }

    return index;
}

// Initialize keyboard extension, keymap and input method
static void InitKeyboard(void)
{
    if (platform.xkbAvailable)
    {
        // NOTE: Key repeats are sent as press events only, without the release events of core X11
        Bool supported = False;
        if (XkbSetDetectableAutoRepeat(platform.display, True, &supported) && supported) platform.xkbDetectableRepeat = true;

        XkbStateRec state = { 0 };
        if (XkbGetState(platform.display, XkbUseCoreKbd, &state) == Success)
        {
            platform.xkbGroup = state.group;
            platform.capsLocked = FLAG_IS_SET(state.locked_mods, LockMask);
        }

        // Keyboard layout group, lock modifiers and keymap changes
        XkbSelectEventDetails(platform.display, XkbUseCoreKbd, XkbStateNotify, XkbGroupStateMask | XkbModifierLockMask, XkbGroupStateMask | XkbModifierLockMask);
        XkbSelectEvents(platform.display, XkbUseCoreKbd, XkbNewKeyboardNotifyMask | XkbMapNotifyMask, XkbNewKeyboardNotifyMask | XkbMapNotifyMask);

        platform.numLockMask = XkbKeysymToModifiers(platform.display, XK_Num_Lock);
        platform.numLocked = ((platform.numLockMask != 0) && FLAG_IS_SET(state.locked_mods, platform.numLockMask));
    }

    UpdateKeymap();

    // Input method: text input considering keyboard layout, compose sequences (dead keys) and text composition
    // NOTE: Input methods require the locale of the environment, the default "C" locale of a program has no compose sequences
    const char *locale = setlocale(LC_CTYPE, NULL);
    if ((locale != NULL) && (strcmp(locale, "C") == 0)) setlocale(LC_CTYPE, "");

    if (XSupportsLocale())
    {
        // NOTE: Xlib built-in input method is used when the one of the environment (XMODIFIERS) is not available
        if (XSetLocaleModifiers("") == NULL) XSetLocaleModifiers("@im=none");

        // NOTE: Callback is called when an input method is available, also when it is started later
        XRegisterIMInstantiateCallback(platform.display, NULL, NULL, NULL, InputMethodInstantiateCallback, NULL);
        if (platform.inputMethod == NULL) InputMethodInstantiateCallback(platform.display, NULL, NULL);
    }
}

// Update map of X11 keycodes to raylib keys
// NOTE: Keys are identified by XKB key names (physical keys), key symbols are used for keys without a known name
static void UpdateKeymap(void)
{
    memset(platform.keymap, 0, sizeof(platform.keymap));

    int keycodeMin = 8;
    int keycodeMax = 255;
    XDisplayKeycodes(platform.display, &keycodeMin, &keycodeMax);
    if (keycodeMax >= KEYMAP_SIZE) keycodeMax = KEYMAP_SIZE - 1;

    const int nameCount = (int)(sizeof(keyNameMap)/sizeof(keyNameMap[0]));

    if (platform.xkbAvailable)
    {
        XkbDescPtr desc = XkbGetMap(platform.display, 0, XkbUseCoreKbd);

        if ((desc != NULL) && (XkbGetNames(platform.display, XkbKeyNamesMask | XkbKeyAliasesMask, desc) == Success) &&
            (desc->names != NULL) && (desc->names->keys != NULL))
        {
            for (int keycode = desc->min_key_code; (keycode <= desc->max_key_code) && (keycode < KEYMAP_SIZE); keycode++)
            {
                const char *name = desc->names->keys[keycode].name;

                for (int i = 0; (i < nameCount) && (platform.keymap[keycode] == 0); i++)
                {
                    if (strncmp(name, keyNameMap[i].name, XkbKeyNameLength) == 0) platform.keymap[keycode] = keyNameMap[i].key;
                }

                // Key name aliases
                for (int a = 0; (desc->names->key_aliases != NULL) && (a < desc->names->num_key_aliases) && (platform.keymap[keycode] == 0); a++)
                {
                    if (strncmp(desc->names->key_aliases[a].real, name, XkbKeyNameLength) != 0) continue;

                    for (int i = 0; (i < nameCount) && (platform.keymap[keycode] == 0); i++)
                    {
                        if (strncmp(desc->names->key_aliases[a].alias, keyNameMap[i].name, XkbKeyNameLength) == 0) platform.keymap[keycode] = keyNameMap[i].key;
                    }
                }
            }
        }

        if (desc != NULL)
        {
            XkbFreeNames(desc, XkbKeyNamesMask | XkbKeyAliasesMask, True);
            XkbFreeKeyboard(desc, 0, True);
        }

        platform.numLockMask = XkbKeysymToModifiers(platform.display, XK_Num_Lock);
    }

    // Keys without a known name are mapped by the key symbol of the first layout
    int symbolsPerKeycode = 0;
    KeySym *symbols = XGetKeyboardMapping(platform.display, (KeyCode)keycodeMin, keycodeMax - keycodeMin + 1, &symbolsPerKeycode);
    if (symbols != NULL)
    {
        for (int keycode = keycodeMin; keycode <= keycodeMax; keycode++)
        {
            if (platform.keymap[keycode] == 0) platform.keymap[keycode] = (short)KeysymToKey(symbols[(keycode - keycodeMin)*symbolsPerKeycode]);
        }

        XFree(symbols);
    }
}

// Get raylib key for a key symbol
static int KeysymToKey(KeySym keysym)
{
    if ((keysym >= XK_a) && (keysym <= XK_z)) return KEY_A + (int)(keysym - XK_a);
    if ((keysym >= XK_A) && (keysym <= XK_Z)) return KEY_A + (int)(keysym - XK_A);
    if ((keysym >= XK_0) && (keysym <= XK_9)) return KEY_ZERO + (int)(keysym - XK_0);
    if ((keysym >= XK_F1) && (keysym <= XK_F25)) return KEY_F1 + (int)(keysym - XK_F1);
    if ((keysym >= XK_KP_0) && (keysym <= XK_KP_9)) return KEY_KP_0 + (int)(keysym - XK_KP_0);

    switch (keysym)
    {
        case XK_space: return KEY_SPACE;
        case XK_apostrophe: return KEY_APOSTROPHE;
        case XK_comma: return KEY_COMMA;
        case XK_minus: return KEY_MINUS;
        case XK_period: return KEY_PERIOD;
        case XK_slash: return KEY_SLASH;
        case XK_semicolon: return KEY_SEMICOLON;
        case XK_equal: return KEY_EQUAL;
        case XK_bracketleft: return KEY_LEFT_BRACKET;
        case XK_backslash: return KEY_BACKSLASH;
        case XK_bracketright: return KEY_RIGHT_BRACKET;
        case XK_grave: return KEY_GRAVE;
        case XK_Escape: return KEY_ESCAPE;
        case XK_Return: return KEY_ENTER;
        case XK_Tab: return KEY_TAB;
        case XK_BackSpace: return KEY_BACKSPACE;
        case XK_Insert: return KEY_INSERT;
        case XK_Delete: return KEY_DELETE;
        case XK_Right: return KEY_RIGHT;
        case XK_Left: return KEY_LEFT;
        case XK_Down: return KEY_DOWN;
        case XK_Up: return KEY_UP;
        case XK_Page_Up: return KEY_PAGE_UP;
        case XK_Page_Down: return KEY_PAGE_DOWN;
        case XK_Home: return KEY_HOME;
        case XK_End: return KEY_END;
        case XK_Caps_Lock: return KEY_CAPS_LOCK;
        case XK_Scroll_Lock: return KEY_SCROLL_LOCK;
        case XK_Num_Lock: return KEY_NUM_LOCK;
        case XK_Print: return KEY_PRINT_SCREEN;
        case XK_Pause: return KEY_PAUSE;
        case XK_Shift_L: return KEY_LEFT_SHIFT;
        case XK_Control_L: return KEY_LEFT_CONTROL;
        case XK_Alt_L:
        case XK_Meta_L: return KEY_LEFT_ALT;
        case XK_Super_L: return KEY_LEFT_SUPER;
        case XK_Shift_R: return KEY_RIGHT_SHIFT;
        case XK_Control_R: return KEY_RIGHT_CONTROL;
        case XK_Alt_R:
        case XK_Meta_R:
        case XK_Mode_switch:
        case XK_ISO_Level3_Shift: return KEY_RIGHT_ALT;
        case XK_Super_R: return KEY_RIGHT_SUPER;
        case XK_Menu: return KEY_KB_MENU;
        case XK_KP_Decimal:
        case XK_KP_Separator:
        case XK_KP_Delete: return KEY_KP_DECIMAL;
        case XK_KP_Divide: return KEY_KP_DIVIDE;
        case XK_KP_Multiply: return KEY_KP_MULTIPLY;
        case XK_KP_Subtract: return KEY_KP_SUBTRACT;
        case XK_KP_Add: return KEY_KP_ADD;
        case XK_KP_Enter: return KEY_KP_ENTER;
        case XK_KP_Equal: return KEY_KP_EQUAL;
        case XK_KP_Insert: return KEY_KP_0;
        case XK_KP_End: return KEY_KP_1;
        case XK_KP_Down: return KEY_KP_2;
        case XK_KP_Page_Down: return KEY_KP_3;
        case XK_KP_Left: return KEY_KP_4;
        case XK_KP_Begin: return KEY_KP_5;
        case XK_KP_Right: return KEY_KP_6;
        case XK_KP_Home: return KEY_KP_7;
        case XK_KP_Up: return KEY_KP_8;
        case XK_KP_Page_Up: return KEY_KP_9;
        default: break;
    }

    return KEY_NULL;
}

// Input method available: open it and create the window input context
static void InputMethodInstantiateCallback(Display *display, XPointer clientData, XPointer callData)
{
    if (platform.inputMethod != NULL) return;

    platform.inputMethod = XOpenIM(platform.display, NULL, NULL, NULL);
    if (platform.inputMethod == NULL) return;

    // Check the input method supports input without composition feedback drawn by the application
    bool usable = false;
    XIMStyles *styles = NULL;
    if ((XGetIMValues(platform.inputMethod, XNQueryInputStyle, &styles, NULL) == NULL) && (styles != NULL))
    {
        for (int i = 0; i < styles->count_styles; i++)
        {
            if (styles->supported_styles[i] == (XIMPreeditNothing | XIMStatusNothing)) usable = true;
        }

        XFree(styles);
    }

    if (!usable)
    {
        XCloseIM(platform.inputMethod);
        platform.inputMethod = NULL;
        return;
    }

    XIMCallback callback = { 0 };
    callback.callback = (XIMProc)InputMethodDestroyCallback;
    XSetIMValues(platform.inputMethod, XNDestroyCallback, &callback, NULL);

    CreateInputContext();
}

// Create window input context
static void CreateInputContext(void)
{
    if ((platform.inputMethod == NULL) || (platform.window == None) || (platform.inputContext != NULL)) return;

    XIMCallback callback = { 0 };
    callback.callback = (XIMProc)InputContextDestroyCallback;

    platform.inputContext = XCreateIC(platform.inputMethod, XNInputStyle, XIMPreeditNothing | XIMStatusNothing,
        XNClientWindow, platform.window, XNFocusWindow, platform.window, XNDestroyCallback, &callback, NULL);
    if (platform.inputContext == NULL) return;

    // Select the events required by the input method
    unsigned long filterEvents = 0;
    if (XGetICValues(platform.inputContext, XNFilterEvents, &filterEvents, NULL) == NULL)
    {
        platform.eventMask |= (long)filterEvents;
        XSelectInput(platform.display, platform.window, platform.eventMask);
    }

    if (platform.focused) XSetICFocus(platform.inputContext);
}

// Input method connection closed (input method server stopped)
static void InputMethodDestroyCallback(XIM inputMethod, XPointer clientData, XPointer callData)
{
    // NOTE: Input method and its input contexts are already destroyed
    platform.inputMethod = NULL;
    platform.inputContext = NULL;
}

// Input context destroyed
static void InputContextDestroyCallback(XIC inputContext, XPointer clientData, XPointer callData)
{
    platform.inputContext = NULL;
}

// Register key state change
static void HandleKey(XKeyEvent *event, bool pressed, bool filtered)
{
    platform.eventReceived = true;
    unsigned int keycode = event->keycode;
    int key = (keycode < KEYMAP_SIZE)? platform.keymap[keycode] : (int)KEY_NULL;
    bool validKey = ((key > KEY_NULL) && (key < MAX_KEYBOARD_KEYS));

    if (pressed)
    {
        // NOTE: Input methods send again the key events they do not consume, with the time of the original event;
        // the key state is registered once, characters are only registered for events not consumed by the input method
        bool duplicated = false;
        if (keycode < KEYMAP_SIZE)
        {
            // NOTE: Server time is a 32 bit value that wraps around, zero is never sent
            unsigned int previous = (unsigned int)platform.keyPressTimes[keycode];
            unsigned int elapsed = (unsigned int)event->time - previous;
            duplicated = ((previous != 0) && ((elapsed == 0) || (elapsed >= 0x80000000u)));
            if (!duplicated) platform.keyPressTimes[keycode] = event->time;
        }

        if (!duplicated)
        {
            platform.lastUserTime = event->time;
            XChangeProperty(platform.display, platform.window, platform.atoms._NET_WM_USER_TIME, XA_CARDINAL, 32, PropModeReplace,
                (unsigned char *)&platform.lastUserTime, 1);

            if (validKey)
            {
                // NOTE: Press of a key already down is a key repeat, CAPS/NUM lock keys are also down while the lock is active
                bool repeated = (CORE.Input.Keyboard.currentKeyState[key] == 1);

                if (key == KEY_CAPS_LOCK)
                {
                    repeated = platform.capsLockDown;
                    platform.capsLockDown = true;
                }
                else if (key == KEY_NUM_LOCK)
                {
                    repeated = platform.numLockDown;
                    platform.numLockDown = true;
                }

                // NOTE: Key pressed again before its release was registered stays down, it is not a key repeat
                if (platform.keyReleasePending[key]) repeated = false;
                platform.keyPressedInFrame[key] = true;
                platform.keyReleasePending[key] = false;

                if (repeated) CORE.Input.Keyboard.keyRepeatInFrame[key] = 1;
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
        }

        if (!filtered) HandleKeyChar(event);
    }
    else
    {
        if (!platform.xkbDetectableRepeat && (XEventsQueued(platform.display, QueuedAfterReading) > 0))
        {
            // NOTE: Without detectable key repeat, repeats are sent as release and press events with the same time
            XEvent next = { 0 };
            XPeekEvent(platform.display, &next);

            if ((next.type == KeyPress) && (next.xkey.window == event->window) && (next.xkey.keycode == keycode) &&
                ((next.xkey.time - event->time) < KEY_REPEAT_EPSILON)) return;
        }

        // NOTE: Press and release can be received in the same frame (e.g. keys sent by other programs),
        // the release is registered on next frame to not miss the key press
        if (validKey)
        {
            if (platform.keyPressedInFrame[key]) platform.keyReleasePending[key] = true;
            else ReleaseKey(key);
        }
    }
}

// Register key release
static void ReleaseKey(int key)
{
    CORE.Input.Keyboard.currentKeyState[key] = 0;

    // NOTE: CAPS/NUM lock keys are considered down while the lock is active
    if (key == KEY_CAPS_LOCK)
    {
        platform.capsLockDown = false;
        if (platform.capsLocked) CORE.Input.Keyboard.currentKeyState[key] = 1;
    }
    else if (key == KEY_NUM_LOCK)
    {
        platform.numLockDown = false;
        if (platform.numLocked) CORE.Input.Keyboard.currentKeyState[key] = 1;
    }
}

// Register character input for a key press
static void HandleKeyChar(XKeyEvent *event)
{
    // NOTE: Characters are not registered while CTRL or ALT are down (key shortcuts), like other desktop platforms
    if (FLAG_IS_SET(event->state, ControlMask) || FLAG_IS_SET(event->state, Mod1Mask)) return;

    if (platform.inputContext != NULL)
    {
        char buffer[128] = { 0 };
        char *text = buffer;
        Status status = 0;

        int length = Xutf8LookupString(platform.inputContext, event, buffer, sizeof(buffer) - 1, NULL, &status);
        if (status == XBufferOverflow)
        {
            text = (char *)RL_CALLOC(length + 1, sizeof(char));
            length = Xutf8LookupString(platform.inputContext, event, text, length, NULL, &status);
        }

        if ((status == XLookupChars) || (status == XLookupBoth))
        {
            // Decode UTF-8 text
            const unsigned char *bytes = (const unsigned char *)text;
            int i = 0;

            while (i < length)
            {
                unsigned int codepoint = bytes[i];
                int size = 1;

                if ((bytes[i] & 0xe0) == 0xc0) { codepoint = bytes[i] & 0x1f; size = 2; }
                else if ((bytes[i] & 0xf0) == 0xe0) { codepoint = bytes[i] & 0x0f; size = 3; }
                else if ((bytes[i] & 0xf8) == 0xf0) { codepoint = bytes[i] & 0x07; size = 4; }

                if ((i + size) > length) break;
                for (int k = 1; k < size; k++) codepoint = (codepoint << 6) | (bytes[i + k] & 0x3f);
                i += size;

                // Check if there is space available in the queue, control characters are not registered
                if ((codepoint >= 32) && ((codepoint < 127) || (codepoint >= 160)) && (CORE.Input.Keyboard.charPressedQueueCount < MAX_CHAR_PRESSED_QUEUE))
                {
                    // Add character to the queue
                    CORE.Input.Keyboard.charPressedQueue[CORE.Input.Keyboard.charPressedQueueCount] = (int)codepoint;
                    CORE.Input.Keyboard.charPressedQueueCount++;
                }
            }
        }

        if (text != buffer) RL_FREE(text);
    }
    else
    {
        // NOTE: Without input method there are no compose sequences, the key symbol considers layout and modifiers
        KeySym keysym = NoSymbol;
        XLookupString(event, NULL, 0, &keysym, NULL);

        unsigned int codepoint = xkb_keysym_to_utf32((xkb_keysym_t)keysym);

        if ((codepoint >= 32) && ((codepoint < 127) || (codepoint >= 160)) && (CORE.Input.Keyboard.charPressedQueueCount < MAX_CHAR_PRESSED_QUEUE))
        {
            CORE.Input.Keyboard.charPressedQueue[CORE.Input.Keyboard.charPressedQueueCount] = (int)codepoint;
            CORE.Input.Keyboard.charPressedQueueCount++;
        }
    }
}

// Register window focus change
static void HandleFocus(bool focused)
{
    platform.eventReceived = true;
    platform.focused = focused;

    if (focused)
    {
        FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_UNFOCUSED);
        if (platform.inputContext != NULL) XSetICFocus(platform.inputContext);

        SyncInputState();
    }
    else
    {
        FLAG_SET(CORE.Window.flags, FLAG_WINDOW_UNFOCUSED);
        if (platform.inputContext != NULL) XUnsetICFocus(platform.inputContext);

        // Release keys and mouse buttons, release events are not received without focus
        for (int i = 0; i < MAX_KEYBOARD_KEYS; i++)
        {
            CORE.Input.Keyboard.currentKeyState[i] = 0;
            platform.keyReleasePending[i] = false;
        }

        CORE.Input.Keyboard.currentKeyState[KEY_CAPS_LOCK] = platform.capsLocked;
        CORE.Input.Keyboard.currentKeyState[KEY_NUM_LOCK] = platform.numLocked;
        platform.capsLockDown = false;
        platform.numLockDown = false;

        for (int i = 0; i < MAX_MOUSE_BUTTONS; i++)
        {
            if (CORE.Input.Mouse.currentButtonState[i] == 1) HandleMouseButton(i, false);
        }

        // NOTE: Cursor stays disabled, the pointer is grabbed again when it is used on the focused window
        LockPointer(false);
    }
}

// Release keys and mouse buttons not held anymore
// NOTE: Release events are not received while another client holds a grab (window manager shortcuts, menus)
static void SyncInputState(void)
{
    char keys[32] = { 0 };
    XQueryKeymap(platform.display, keys);

    bool keyDown[MAX_KEYBOARD_KEYS] = { 0 };
    for (int keycode = 0; keycode < KEYMAP_SIZE; keycode++)
    {
        int key = platform.keymap[keycode];
        if ((key > KEY_NULL) && (key < MAX_KEYBOARD_KEYS) && FLAG_IS_SET(keys[keycode/8], (1 << (keycode%8)))) keyDown[key] = true;
    }

    for (int key = 0; key < MAX_KEYBOARD_KEYS; key++)
    {
        if (!keyDown[key] && (CORE.Input.Keyboard.currentKeyState[key] == 1)) CORE.Input.Keyboard.currentKeyState[key] = 0;
    }

    // NOTE: CAPS/NUM lock keys are considered down while the lock is active, a lock that was already active is not a key press
    platform.capsLockDown = keyDown[KEY_CAPS_LOCK];
    platform.numLockDown = keyDown[KEY_NUM_LOCK];

    if (platform.capsLocked && (CORE.Input.Keyboard.currentKeyState[KEY_CAPS_LOCK] == 0))
    {
        CORE.Input.Keyboard.currentKeyState[KEY_CAPS_LOCK] = 1;
        CORE.Input.Keyboard.previousKeyState[KEY_CAPS_LOCK] = 1;
    }

    if (platform.numLocked && (CORE.Input.Keyboard.currentKeyState[KEY_NUM_LOCK] == 0))
    {
        CORE.Input.Keyboard.currentKeyState[KEY_NUM_LOCK] = 1;
        CORE.Input.Keyboard.previousKeyState[KEY_NUM_LOCK] = 1;
    }

    Window root = None;
    Window child = None;
    int rootX = 0;
    int rootY = 0;
    int windowX = 0;
    int windowY = 0;
    unsigned int mask = 0;

    if (XQueryPointer(platform.display, platform.window, &root, &child, &rootX, &rootY, &windowX, &windowY, &mask))
    {
        if (!FLAG_IS_SET(mask, Button1Mask) && (CORE.Input.Mouse.currentButtonState[MOUSE_BUTTON_LEFT] == 1)) HandleMouseButton(MOUSE_BUTTON_LEFT, false);
        if (!FLAG_IS_SET(mask, Button2Mask) && (CORE.Input.Mouse.currentButtonState[MOUSE_BUTTON_MIDDLE] == 1)) HandleMouseButton(MOUSE_BUTTON_MIDDLE, false);
        if (!FLAG_IS_SET(mask, Button3Mask) && (CORE.Input.Mouse.currentButtonState[MOUSE_BUTTON_RIGHT] == 1)) HandleMouseButton(MOUSE_BUTTON_RIGHT, false);
    }
}

// Initialize XInput2 and select pointer events
// NOTE: XInput2 events replace the core pointer events, they provide smooth scrolling, touch and subpixel positions
static void InitInputDevices(void)
{
    if (!platform.xiAvailable) return;

    unsigned char mask[XIMaskLen(XI_LASTEVENT)] = { 0 };
    XIEventMask eventMask = { 0 };
    eventMask.mask = mask;
    eventMask.mask_len = sizeof(mask);

    // Pointer events of the window
    eventMask.deviceid = XIAllMasterDevices;
    XISetMask(mask, XI_Motion);
    XISetMask(mask, XI_ButtonPress);
    XISetMask(mask, XI_ButtonRelease);
    if (platform.xiMinor >= 2)
    {
        XISetMask(mask, XI_TouchBegin);
        XISetMask(mask, XI_TouchUpdate);
        XISetMask(mask, XI_TouchEnd);
    }

    if (XISelectEvents(platform.display, platform.window, &eventMask, 1) != Success)
    {
        platform.xiAvailable = false;
        return;
    }

    platform.eventMask &= ~(PointerMotionMask | ButtonPressMask | ButtonReleaseMask);
    XSelectInput(platform.display, platform.window, platform.eventMask);

    // Devices added, removed or changed
    memset(mask, 0, sizeof(mask));
    eventMask.deviceid = XIAllDevices;
    XISetMask(mask, XI_HierarchyChanged);
    XISetMask(mask, XI_DeviceChanged);
    XISelectEvents(platform.display, platform.root, &eventMask, 1);

    UpdateInputDevices();
}

// Update pointer devices info (scroll axes, motion mode)
static void UpdateInputDevices(void)
{
    platform.deviceCount = 0;

    int count = 0;
    XIDeviceInfo *infos = XIQueryDevice(platform.display, XIAllDevices, &count);
    if (infos == NULL) return;

    for (int i = 0; (i < count) && (platform.deviceCount < MAX_INPUT_DEVICES); i++)
    {
        if ((infos[i].use != XIMasterPointer) && (infos[i].use != XISlavePointer) && (infos[i].use != XIFloatingSlave)) continue;

        InputDevice *device = &platform.devices[platform.deviceCount];
        memset(device, 0, sizeof(InputDevice));
        device->id = infos[i].deviceid;
        device->axes[0].relative = true;
        device->axes[1].relative = true;

        for (int c = 0; c < infos[i].num_classes; c++)
        {
            const XIAnyClassInfo *classInfo = infos[i].classes[c];

            if (classInfo->type == XIValuatorClass)
            {
                const XIValuatorClassInfo *valuator = (const XIValuatorClassInfo *)classInfo;
                if ((valuator->number == 0) || (valuator->number == 1))
                {
                    device->axes[valuator->number].relative = (valuator->mode == XIModeRelative);
                    device->axes[valuator->number].min = valuator->min;
                    device->axes[valuator->number].max = valuator->max;
                }
            }
#if defined(XIScrollClass)
            else if ((classInfo->type == XIScrollClass) && (device->scrollAxisCount < MAX_SCROLL_AXES))
            {
                const XIScrollClassInfo *scroll = (const XIScrollClassInfo *)classInfo;
                if (scroll->increment != 0.0)
                {
                    ScrollAxis *axis = &device->scrollAxes[device->scrollAxisCount];
                    axis->number = scroll->number;
                    axis->horizontal = (scroll->scroll_type == XIScrollTypeHorizontal);
                    axis->increment = scroll->increment;
                    device->scrollAxisCount++;
                }
            }
#endif
        }

        // Scroll axes report accumulated values, the current value is the reference for the next scroll event
        for (int a = 0; a < device->scrollAxisCount; a++)
        {
            for (int c = 0; c < infos[i].num_classes; c++)
            {
                const XIValuatorClassInfo *valuator = (const XIValuatorClassInfo *)infos[i].classes[c];
                if ((valuator->type == XIValuatorClass) && (valuator->number == device->scrollAxes[a].number))
                {
                    device->scrollAxes[a].value = valuator->value;
                    device->scrollAxes[a].valid = true;
                }
            }
        }

        platform.deviceCount++;
    }

    XIFreeDeviceInfo(infos);
}

// Get pointer device info
static InputDevice *GetInputDevice(int id)
{
    for (int i = 0; i < platform.deviceCount; i++)
    {
        if (platform.devices[i].id == id) return &platform.devices[i];
    }

    return NULL;
}

// Forget last values of absolute motion axes
// NOTE: Values change while the pointer is not grabbed
static void ResetMotionAxes(void)
{
    for (int i = 0; i < platform.deviceCount; i++)
    {
        platform.devices[i].axes[0].valid = false;
        platform.devices[i].axes[1].valid = false;
    }
}

// Process XInput2 event
static void ProcessGenericEvent(XGenericEventCookie *cookie)
{
    if (!platform.xiAvailable || (cookie->extension != platform.xiOpcode)) return;
    if (!XGetEventData(platform.display, cookie)) return;

    float scale = FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_HIGHDPI)? platform.scale : 1.0f;

    switch (cookie->evtype)
    {
        case XI_RawMotion:
        {
            // Relative motion while the cursor is disabled
            // NOTE: Raw values are not accelerated, absolute devices (tablets, virtual machines) report positions
            const XIRawEvent *raw = (const XIRawEvent *)cookie->data;
            if (!platform.pointerGrabbed) break;

            InputDevice *device = GetInputDevice((platform.xiMinor >= 1)? raw->sourceid : raw->deviceid);
            double delta[2] = { 0.0, 0.0 };
            int index = 0;

            for (int bit = 0; bit < raw->valuators.mask_len*8; bit++)
            {
                if (!XIMaskIsSet(raw->valuators.mask, bit)) continue;

                double value = raw->raw_values[index++];
                if (bit > 1) continue;

                if ((device != NULL) && !device->axes[bit].relative)
                {
                    MotionAxis *axis = &device->axes[bit];
                    double range = axis->max - axis->min;
                    double size = (bit == 0)? (double)DisplayWidth(platform.display, platform.screen) : (double)DisplayHeight(platform.display, platform.screen);

                    if (axis->valid) delta[bit] = (value - axis->value)*((range > 0.0)? size/range : 1.0);
                    axis->value = value;
                    axis->valid = true;
                }
                else delta[bit] = value;
            }

            if ((delta[0] != 0.0) || (delta[1] != 0.0))
            {
                HandleMouseMove((double)CORE.Input.Mouse.currentPosition.x + delta[0], (double)CORE.Input.Mouse.currentPosition.y + delta[1]);
            }
        } break;
        case XI_Motion:
        {
            const XIDeviceEvent *deviceEvent = (const XIDeviceEvent *)cookie->data;

            // Smooth scrolling: scroll axes report accumulated values
            InputDevice *device = GetInputDevice(deviceEvent->sourceid);
            if ((device != NULL) && (device->scrollAxisCount > 0))
            {
                int index = 0;
                for (int bit = 0; bit < deviceEvent->valuators.mask_len*8; bit++)
                {
                    if (!XIMaskIsSet(deviceEvent->valuators.mask, bit)) continue;

                    double value = deviceEvent->valuators.values[index++];

                    for (int a = 0; a < device->scrollAxisCount; a++)
                    {
                        ScrollAxis *axis = &device->scrollAxes[a];
                        if (axis->number != bit) continue;

                        if (axis->valid)
                        {
                            float move = -(float)((value - axis->value)/axis->increment);
                            if (axis->horizontal) CORE.Input.Mouse.currentWheelMove.x += move;
                            else CORE.Input.Mouse.currentWheelMove.y += move;
                            platform.eventReceived = true;
                        }

                        axis->value = value;
                        axis->valid = true;
                    }
                }
            }

            HandlePointerMotion(deviceEvent->event_x, deviceEvent->event_y);
        } break;
        case XI_ButtonPress:
        case XI_ButtonRelease:
        {
            const XIDeviceEvent *deviceEvent = (const XIDeviceEvent *)cookie->data;

            // NOTE: Wheel buttons emulated from scroll axes are ignored, the scroll axes values are used
            bool wheelButton = ((deviceEvent->detail >= Button4) && (deviceEvent->detail <= 7));
            if (wheelButton && FLAG_IS_SET(deviceEvent->flags, XIPointerEmulated)) break;

            HandlePointerButton(deviceEvent->detail, (cookie->evtype == XI_ButtonPress), deviceEvent->time);
        } break;
        case XI_TouchBegin:
        case XI_TouchUpdate:
        case XI_TouchEnd:
        {
            const XIDeviceEvent *deviceEvent = (const XIDeviceEvent *)cookie->data;

            int touchAction = TOUCH_ACTION_MOVE;
            if (cookie->evtype == XI_TouchBegin) touchAction = TOUCH_ACTION_DOWN;
            else if (cookie->evtype == XI_TouchEnd) touchAction = TOUCH_ACTION_UP;

            HandleTouch(touchAction, deviceEvent->detail, deviceEvent->event_x/scale, deviceEvent->event_y/scale);
        } break;
        case XI_HierarchyChanged:
        case XI_DeviceChanged: UpdateInputDevices(); break;
        default: break;
    }

    XFreeEventData(platform.display, cookie);
}

// Register mouse button state change
static void HandleMouseButton(int button, bool pressed)
{
    platform.eventReceived = true;

    // NOTE: Press and release can be received in the same frame (e.g. button also used for scrolling, sent on release),
    // the release is registered on next frame to not miss the button press
    if (pressed)
    {
        platform.buttonPressedInFrame[button] = true;
        platform.buttonReleasePending[button] = false;
    }
    else if (platform.buttonPressedInFrame[button])
    {
        platform.buttonReleasePending[button] = true;
        return;
    }

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

// Register mouse position change (logical pixels)
static void HandleMouseMove(double x, double y)
{
    platform.eventReceived = true;
    CORE.Input.Mouse.currentPosition.x = (float)x;
    CORE.Input.Mouse.currentPosition.y = (float)y;
    if (CORE.Input.Touch.pointCount == 0) CORE.Input.Touch.position[0] = CORE.Input.Mouse.currentPosition;

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

// Register X11 pointer button, wheel buttons included
static void HandlePointerButton(int button, bool pressed, Time time)
{
    if (pressed)
    {
        platform.lastUserTime = time;
        XChangeProperty(platform.display, platform.window, platform.atoms._NET_WM_USER_TIME, XA_CARDINAL, 32, PropModeReplace,
            (unsigned char *)&platform.lastUserTime, 1);

        // NOTE: Disabled cursor is grabbed when the pointer is used on the focused window
        if (CORE.Input.Mouse.cursorLocked) LockPointer(true);
    }

    // Mouse wheel steps are reported as buttons: 4 (up), 5 (down), 6 (left) and 7 (right)
    if ((button >= Button4) && (button <= 7))
    {
        if (pressed)
        {
            if (button == Button4) CORE.Input.Mouse.currentWheelMove.y += 1.0f;
            else if (button == Button5) CORE.Input.Mouse.currentWheelMove.y -= 1.0f;
            else if (button == 6) CORE.Input.Mouse.currentWheelMove.x += 1.0f;
            else CORE.Input.Mouse.currentWheelMove.x -= 1.0f;

            platform.eventReceived = true;
        }

        return;
    }

    int mouseButton = -1;
    if (button == Button1) mouseButton = MOUSE_BUTTON_LEFT;
    else if (button == Button2) mouseButton = MOUSE_BUTTON_MIDDLE;
    else if (button == Button3) mouseButton = MOUSE_BUTTON_RIGHT;
    else if (button >= 8) mouseButton = MOUSE_BUTTON_SIDE + (button - 8);

    if ((mouseButton >= 0) && (mouseButton < MAX_MOUSE_BUTTONS)) HandleMouseButton(mouseButton, pressed);
}

// Register X11 pointer position (pixels)
static void HandlePointerMotion(double x, double y)
{
    if (CORE.Input.Mouse.cursorLocked)
    {
        // NOTE: Disabled cursor is grabbed when the pointer is used on the focused window
        if (!platform.pointerGrabbed && (platform.focused || (platform.wmSupportedCount == 0))) LockPointer(true);

        // NOTE: Without XInput2 raw motion, relative motion is the distance to the previous position,
        // the pointer is moved back to the window center in PollInputEvents()
        if (platform.pointerGrabbed && !platform.xiAvailable && ((x != platform.warpX) || (y != platform.warpY)))
        {
            float scale = FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_HIGHDPI)? platform.scale : 1.0f;
            HandleMouseMove((double)CORE.Input.Mouse.currentPosition.x + (x - platform.pointerX)/scale,
                (double)CORE.Input.Mouse.currentPosition.y + (y - platform.pointerY)/scale);
        }

        platform.pointerX = x;
        platform.pointerY = y;
        return;
    }

    if ((x == platform.pointerX) && (y == platform.pointerY)) return;

    float scale = FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_HIGHDPI)? platform.scale : 1.0f;

    platform.pointerX = x;
    platform.pointerY = y;
    HandleMouseMove(x/scale, y/scale);
}

// Register touch event
static void HandleTouch(int touchAction, int id, double x, double y)
{
    int slot = -1;
    for (int i = 0; i < MAX_TOUCH_POINTS; i++)
    {
        if (platform.touchIds[i] == id) slot = i;
    }

    if ((slot < 0) && (touchAction == TOUCH_ACTION_DOWN))
    {
        for (int i = 0; (i < MAX_TOUCH_POINTS) && (slot < 0); i++)
        {
            if (platform.touchIds[i] == -1) slot = i;
        }
    }

    if (slot < 0) return;

    platform.eventReceived = true;
    platform.touchIds[slot] = id;
    platform.touchPositions[slot] = (Vector2){ (float)x, (float)y };

#if SUPPORT_GESTURES_SYSTEM
    GestureEvent gestureEvent = { 0 };
    gestureEvent.touchAction = touchAction;
#endif

    // Register touch points, released points are included in the gesture event
    int pointCount = 0;
    for (int i = 0; i < MAX_TOUCH_POINTS; i++)
    {
        if (platform.touchIds[i] == -1) continue;

    #if SUPPORT_GESTURES_SYSTEM
        gestureEvent.pointId[pointCount] = platform.touchIds[i];
        gestureEvent.position[pointCount].x = platform.touchPositions[i].x/(float)GetScreenWidth();
        gestureEvent.position[pointCount].y = platform.touchPositions[i].y/(float)GetScreenHeight();
    #endif
        pointCount++;
    }

#if SUPPORT_GESTURES_SYSTEM
    gestureEvent.pointCount = pointCount;
    ProcessGestureEvent(gestureEvent);
#endif

    if (touchAction == TOUCH_ACTION_UP) platform.touchIds[slot] = -1;

    // Update touch points state
    CORE.Input.Touch.pointCount = 0;
    for (int i = 0; i < MAX_TOUCH_POINTS; i++)
    {
        if (platform.touchIds[i] == -1) continue;

        CORE.Input.Touch.pointId[CORE.Input.Touch.pointCount] = platform.touchIds[i];
        CORE.Input.Touch.position[CORE.Input.Touch.pointCount] = platform.touchPositions[i];
        CORE.Input.Touch.currentTouchState[CORE.Input.Touch.pointCount] = 1;
        CORE.Input.Touch.pointCount++;
    }

    for (int i = CORE.Input.Touch.pointCount; i < MAX_TOUCH_POINTS; i++) CORE.Input.Touch.currentTouchState[i] = 0;

    // Map primary touch point to mouse position
    if (CORE.Input.Touch.pointCount > 0) CORE.Input.Mouse.currentPosition = CORE.Input.Touch.position[0];
}

// Update window cursor
// NOTE: Disabled cursor is only hidden while the pointer is grabbed
static void UpdateCursor(void)
{
    if (platform.window == None) return;

    bool hidden = (CORE.Input.Mouse.cursorHidden && (!CORE.Input.Mouse.cursorLocked || platform.pointerGrabbed));

    XDefineCursor(platform.display, platform.window, hidden? platform.hiddenCursor : GetCursor(CORE.Input.Mouse.cursor));
    XFlush(platform.display);
}

// Get cursor for a MouseCursor shape
// NOTE: Cursors are loaded from the cursor theme, core cursor shapes are used when the theme does not provide them
static Cursor GetCursor(int cursor)
{
    if ((cursor < 0) || (cursor >= MOUSE_CURSOR_COUNT)) cursor = MOUSE_CURSOR_DEFAULT;

    if (platform.cursors[cursor] == None)
    {
        platform.cursors[cursor] = XcursorLibraryLoadCursor(platform.display, cursorNames[cursor][0]);
        if (platform.cursors[cursor] == None) platform.cursors[cursor] = XcursorLibraryLoadCursor(platform.display, cursorNames[cursor][1]);
        if (platform.cursors[cursor] == None) platform.cursors[cursor] = XCreateFontCursor(platform.display, cursorShapes[cursor]);
    }

    return platform.cursors[cursor];
}

// Grab or release pointer (disabled cursor)
// NOTE: The pointer is hidden and confined to the window, relative motion is received as XInput2 raw motion;
// the grab fails while another client holds the pointer (window manager, menus), it is requested again on next pointer events
static void LockPointer(bool lock)
{
    if (lock)
    {
        if (platform.pointerGrabbed || !platform.mapped) return;

        // NOTE: Pointer is only grabbed by the window with the input focus (focus follows the pointer without window manager)
        Window focusWindow = None;
        int revertTo = 0;
        XGetInputFocus(platform.display, &focusWindow, &revertTo);
        if ((focusWindow != platform.window) && (focusWindow != PointerRoot)) return;

        int result = XGrabPointer(platform.display, platform.window, False, ButtonPressMask | ButtonReleaseMask | PointerMotionMask,
            GrabModeAsync, GrabModeAsync, platform.window, platform.hiddenCursor, CurrentTime);
        if (result != GrabSuccess) return;

        platform.pointerGrabbed = true;
        ResetMotionAxes();
        SelectRawMotion(true);

        // Hidden cursor is placed in the window center
        WarpPointer(platform.bufferWidth/2, platform.bufferHeight/2);
        platform.pointerX = platform.warpX;
        platform.pointerY = platform.warpY;
    }
    else
    {
        if (!platform.pointerGrabbed) return;

        XUngrabPointer(platform.display, CurrentTime);
        SelectRawMotion(false);
        platform.pointerGrabbed = false;
    }

    UpdateCursor();
}

// Select raw motion events on the root window
static void SelectRawMotion(bool enabled)
{
    if (!platform.xiAvailable) return;

    unsigned char mask[XIMaskLen(XI_RawMotion)] = { 0 };
    XIEventMask eventMask = { 0 };
    eventMask.deviceid = XIAllMasterDevices;
    eventMask.mask = mask;
    eventMask.mask_len = sizeof(mask);
    if (enabled) XISetMask(mask, XI_RawMotion);

    XISelectEvents(platform.display, platform.root, &eventMask, 1);
}

// Move pointer to window position (pixels)
static void WarpPointer(int x, int y)
{
    platform.warpX = (double)x;
    platform.warpY = (double)y;
    XWarpPointer(platform.display, None, platform.window, 0, 0, 0, 0, x, y);
}

// Send clipboard content to another client
static void HandleSelectionRequest(XSelectionRequestEvent *request)
{
    Atoms *atoms = &platform.atoms;

    XEvent reply = { 0 };
    reply.xselection.type = SelectionNotify;
    reply.xselection.requestor = request->requestor;
    reply.xselection.selection = request->selection;
    reply.xselection.target = request->target;
    reply.xselection.time = request->time;
    reply.xselection.property = None;

    if ((request->selection == atoms->CLIPBOARD) && platform.clipboardOwned && (platform.clipboardText != NULL))
    {
        // NOTE: Obsolete clients do not set the property, the target is used
        Atom property = (request->property != None)? request->property : request->target;

        if (request->target == atoms->MULTIPLE)
        {
            // Multiple conversions: pairs of target and property
            unsigned char *value = NULL;
            unsigned long count = (request->property != None)? GetWindowProperty(request->requestor, request->property, atoms->ATOM_PAIR, &value) : 0;
            Atom *pairs = (Atom *)value;

            for (unsigned long i = 0; (i + 1) < count; i += 2)
            {
                if (!WriteSelection(request->requestor, pairs[i + 1], pairs[i])) pairs[i + 1] = None;
            }

            if (count > 0)
            {
                XChangeProperty(platform.display, request->requestor, request->property, atoms->ATOM_PAIR, 32, PropModeReplace, value, (int)count);
                reply.xselection.property = request->property;
            }

            if (value != NULL) XFree(value);
        }
        else if (WriteSelection(request->requestor, property, request->target)) reply.xselection.property = property;
    }

    XSendEvent(platform.display, request->requestor, False, NoEventMask, &reply);
    XFlush(platform.display);
}

// Write clipboard content to a requestor property, returns false if the target is not supported
// NOTE: Content larger than the maximum request size is sent in chunks (incremental transfer)
static bool WriteSelection(Window requestor, Atom property, Atom target)
{
    Atoms *atoms = &platform.atoms;

    if (target == atoms->TARGETS)
    {
        Atom targets[] = { atoms->TARGETS, atoms->MULTIPLE, atoms->UTF8_STRING, atoms->TEXT_PLAIN_UTF8, atoms->TEXT };
        XChangeProperty(platform.display, requestor, property, XA_ATOM, 32, PropModeReplace, (unsigned char *)targets, (int)(sizeof(targets)/sizeof(targets[0])));
        return true;
    }

    if (target == atoms->SAVE_TARGETS)
    {
        // NOTE: Conversion has no data, the clipboard manager requests the targets it saves
        XChangeProperty(platform.display, requestor, property, atoms->NULL_, 32, PropModeReplace, NULL, 0);
        return true;
    }

    if ((target != atoms->UTF8_STRING) && (target != atoms->TEXT_PLAIN_UTF8) && (target != atoms->TEXT)) return false;

    Atom type = (target == atoms->TEXT)? atoms->UTF8_STRING : target;
    int size = (int)strlen(platform.clipboardText);

    long maxRequestSize = XExtendedMaxRequestSize(platform.display);
    if (maxRequestSize == 0) maxRequestSize = XMaxRequestSize(platform.display);
    long maxSize = maxRequestSize*4 - 1024;

    if ((long)size <= maxSize)
    {
        XChangeProperty(platform.display, requestor, property, type, 8, PropModeReplace, (const unsigned char *)platform.clipboardText, size);
        return true;
    }

    // Release transfers of unresponsive requestors and the previous transfer to the same property
    for (int i = 0; i < MAX_INCR_TRANSFERS; i++)
    {
        IncrTransfer *transfer = &platform.incrTransfers[i];
        if (transfer->requestor == None) continue;

        if (((transfer->requestor == requestor) && (transfer->property == property)) ||
            ((GetTime() - transfer->time) > INCR_TRANSFER_TIMEOUT)) ReleaseIncrTransfer(transfer);
    }

    for (int i = 0; i < MAX_INCR_TRANSFERS; i++)
    {
        IncrTransfer *transfer = &platform.incrTransfers[i];
        if (transfer->requestor != None) continue;

        transfer->requestor = requestor;
        transfer->property = property;
        transfer->type = type;
        transfer->size = size;
        transfer->offset = 0;
        transfer->time = GetTime();
        transfer->data = (char *)RL_MALLOC(size);
        memcpy(transfer->data, platform.clipboardText, size);

        // NOTE: The requestor deletes the property when it is ready for the next chunk,
        // events of the own window and the root window are already selected
        unsigned long totalSize = (unsigned long)size;
        if ((requestor != platform.window) && (requestor != platform.root)) XSelectInput(platform.display, requestor, PropertyChangeMask | StructureNotifyMask);
        XChangeProperty(platform.display, requestor, property, atoms->INCR, 32, PropModeReplace, (unsigned char *)&totalSize, 1);
        return true;
    }

    return false;
}

// Send next chunk of an incremental selection transfer, an empty chunk ends the transfer
static void ContinueIncrTransfer(IncrTransfer *transfer)
{
    long maxRequestSize = XExtendedMaxRequestSize(platform.display);
    if (maxRequestSize == 0) maxRequestSize = XMaxRequestSize(platform.display);
    long maxSize = maxRequestSize*4 - 1024;

    int chunk = transfer->size - transfer->offset;
    if ((long)chunk > maxSize) chunk = (int)maxSize;

    XChangeProperty(platform.display, transfer->requestor, transfer->property, transfer->type, 8, PropModeReplace,
        (const unsigned char *)transfer->data + transfer->offset, chunk);
    transfer->offset += chunk;
    transfer->time = GetTime();

    if (chunk == 0) ReleaseIncrTransfer(transfer);

    XFlush(platform.display);
}

// Release incremental selection transfer
static void ReleaseIncrTransfer(IncrTransfer *transfer)
{
    if ((transfer->requestor != platform.window) && (transfer->requestor != platform.root))
    {
        XSelectInput(platform.display, transfer->requestor, NoEventMask);
    }

    RL_FREE(transfer->data);
    memset(transfer, 0, sizeof(IncrTransfer));
}

// Read selection content converted to a target, returned data is allocated and zero-terminated
// NOTE: Waiting is limited in time, the selection owner could be unresponsive
static unsigned char *ReadSelection(Atom selection, Atom target, int *size, Atom *type)
{
    if (size != NULL) *size = 0;

    XConvertSelection(platform.display, selection, target, platform.atoms.RAYLIB_SELECTION, platform.window, GetServerTime());

    // NOTE: Notifications of previous requests that timed out are discarded
    XEvent event = { 0 };
    do
    {
        if (!WaitForEvent(SelectionNotify, platform.window, selection, &event, SELECTION_TIMEOUT)) return NULL;
    } while (event.xselection.target != target);

    if (event.xselection.property == None) return NULL;

    return ReadSelectionProperty(platform.atoms.RAYLIB_SELECTION, size, type);
}

// Read and delete a selection property of the window, returned data is allocated and zero-terminated
// NOTE: Large content is received in chunks (incremental transfer), every chunk is sent when the previous one is deleted
static unsigned char *ReadSelectionProperty(Atom property, int *size, Atom *type)
{
    Display *display = platform.display;

    Atom actualType = None;
    int actualFormat = 0;
    unsigned long itemCount = 0;
    unsigned long bytesAfter = 0;
    unsigned char *value = NULL;

    if (size != NULL) *size = 0;
    if (type != NULL) *type = None;

    if (XGetWindowProperty(display, platform.window, property, 0, LONG_MAX, False, AnyPropertyType, &actualType, &actualFormat,
        &itemCount, &bytesAfter, &value) != Success) return NULL;

    unsigned char *data = NULL;
    int length = 0;

    if (actualType != platform.atoms.INCR)
    {
        // NOTE: Property items of 32 bits are stored as long
        length = (int)itemCount*((actualFormat == 32)? (int)sizeof(long) : actualFormat/8);

        if ((value != NULL) && (length > 0))
        {
            data = (unsigned char *)RL_CALLOC(length + 1, 1);
            memcpy(data, value, length);
        }

        if (value != NULL) XFree(value);
        XDeleteProperty(display, platform.window, property);
    }
    else
    {
        if (value != NULL) XFree(value);

        // Property changes notified before the transfer started are outdated
        EventFilter filter = { PropertyNotify, platform.window, property };
        XEvent event = { 0 };
        while (XCheckIfEvent(display, &event, MatchEvent, (XPointer)&filter));

        // Start the transfer
        XDeleteProperty(display, platform.window, property);
        XFlush(display);

        int capacity = 0;
        bool completed = false;

        while (!completed)
        {
            if (!WaitForEvent(PropertyNotify, platform.window, property, &event, SELECTION_TIMEOUT)) break;

            value = NULL;
            if (XGetWindowProperty(display, platform.window, property, 0, LONG_MAX, True, AnyPropertyType, &actualType, &actualFormat,
                &itemCount, &bytesAfter, &value) != Success) break;

            int chunk = (int)itemCount*((actualFormat == 32)? (int)sizeof(long) : actualFormat/8);

            // NOTE: Empty chunk ends the transfer
            if ((actualType != None) && (chunk == 0)) completed = true;
            else if (chunk > 0)
            {
                if ((length + chunk + 1) > capacity)
                {
                    capacity = (length + chunk + 1)*2;
                    data = (unsigned char *)RL_REALLOC(data, capacity);
                }

                memcpy(data + length, value, chunk);
                length += chunk;
                data[length] = '\0';
            }

            if (value != NULL) XFree(value);
        }

        if (!completed)
        {
            RL_FREE(data);
            data = NULL;
            length = 0;
        }
    }

    if (data == NULL) return NULL;

    if (size != NULL) *size = length;
    if (type != NULL) *type = actualType;

    return data;
}

// Check selection owner supported target
static bool SelectionHasTarget(Atom selection, Atom target)
{
    bool result = false;

    int size = 0;
    Atom *targets = (Atom *)ReadSelection(selection, platform.atoms.TARGETS, &size, NULL);
    if (targets == NULL) return false;

    for (int i = 0; i < (int)(size/sizeof(Atom)); i++)
    {
        if (targets[i] == target)
        {
            result = true;
            break;
        }
    }

    RL_FREE(targets);

    return result;
}

// Event predicate: selection events
static Bool MatchSelectionEvent(Display *display, XEvent *event, XPointer data)
{
    return ((event->type == SelectionRequest) || (event->type == SelectionNotify) || (event->type == PropertyNotify));
}

// Hand clipboard content over to the clipboard manager, the content is lost when the window is destroyed
// NOTE: Waiting is limited in time, the clipboard manager could be unresponsive
static void SaveClipboard(void)
{
    Atoms *atoms = &platform.atoms;

    if (!platform.clipboardOwned || (platform.clipboardText == NULL)) return;
    if (XGetSelectionOwner(platform.display, atoms->CLIPBOARD_MANAGER) == None) return;

    XConvertSelection(platform.display, atoms->CLIPBOARD_MANAGER, atoms->SAVE_TARGETS, None, platform.window, GetServerTime());
    XFlush(platform.display);

    double endTime = GetTime() + (double)SELECTION_TIMEOUT/1000.0;

    while (!platform.connectionLost)
    {
        // The clipboard manager requests the content, it notifies when it is saved
        XEvent event = { 0 };
        while (XCheckIfEvent(platform.display, &event, MatchSelectionEvent, NULL))
        {
            if (event.type == SelectionRequest) HandleSelectionRequest(&event.xselectionrequest);
            else if (event.type == PropertyNotify) ProcessPropertyChange(&event.xproperty);
            else if (event.xselection.selection == atoms->CLIPBOARD_MANAGER) return;
        }

        int remaining = (int)((endTime - GetTime())*1000.0);
        if (remaining <= 0) break;

        XFlush(platform.display);
        WaitForEvents(remaining);
    }
}

// Drag and drop: source entered the window, check the offered data types
static void HandleDragEnter(XClientMessageEvent *event)
{
    Atoms *atoms = &platform.atoms;

    platform.dragSource = None;
    platform.dragAccepted = false;

    int version = (int)((unsigned long)event->data.l[1] >> 24);
    if (version > XDND_VERSION) return;

    platform.dragSource = (Window)event->data.l[0];
    platform.dragVersion = version;

    if (FLAG_IS_SET(event->data.l[1], 1))
    {
        // More than three data types: list stored in a property of the source window
        unsigned char *value = NULL;
        unsigned long count = GetWindowProperty(platform.dragSource, atoms->XdndTypeList, XA_ATOM, &value);
        Atom *types = (Atom *)value;

        for (unsigned long i = 0; i < count; i++)
        {
            if (types[i] == atoms->TEXT_URI_LIST) platform.dragAccepted = true;
        }

        if (value != NULL) XFree(value);
    }
    else
    {
        for (int i = 2; i < 5; i++)
        {
            if ((Atom)event->data.l[i] == atoms->TEXT_URI_LIST) platform.dragAccepted = true;
        }
    }
}

// Drag and drop: source moved over the window, reply if the drop is accepted
static void HandleDragPosition(XClientMessageEvent *event)
{
    if ((platform.dragSource == None) || ((Window)event->data.l[0] != platform.dragSource)) return;

    // NOTE: Pointer events are not received while dragging, position is sent in root window coordinates
    int rootX = (short)(((unsigned long)event->data.l[2] >> 16) & 0xffff);
    int rootY = (short)((unsigned long)event->data.l[2] & 0xffff);
    HandlePointerMotion((double)(rootX - platform.x), (double)(rootY - platform.y));

    XEvent reply = { 0 };
    reply.xclient.type = ClientMessage;
    reply.xclient.window = platform.dragSource;
    reply.xclient.message_type = platform.atoms.XdndStatus;
    reply.xclient.format = 32;
    reply.xclient.data.l[0] = (long)platform.window;
    reply.xclient.data.l[1] = platform.dragAccepted? 1 : 0;
    reply.xclient.data.l[4] = platform.dragAccepted? (long)platform.atoms.XdndActionCopy : None;

    XSendEvent(platform.display, platform.dragSource, False, NoEventMask, &reply);
    XFlush(platform.display);
}

// Drag and drop: data dropped, request the file list to the source
static void HandleDragDrop(XClientMessageEvent *event)
{
    if ((platform.dragSource == None) || ((Window)event->data.l[0] != platform.dragSource)) return;

    if (platform.dragAccepted)
    {
        Time time = (platform.dragVersion >= 1)? (Time)event->data.l[2] : CurrentTime;
        XConvertSelection(platform.display, platform.atoms.XdndSelection, platform.atoms.TEXT_URI_LIST, platform.atoms.RAYLIB_DROP,
            platform.window, time);
        XFlush(platform.display);
    }
    else FinishDrag(false);
}

// Drag and drop: file list received from the source
static void HandleDragData(XSelectionEvent *event)
{
    bool accepted = false;

    if (event->property != None)
    {
        int size = 0;
        char *uris = (char *)ReadSelectionProperty(platform.atoms.RAYLIB_DROP, &size, NULL);

        if (uris != NULL)
        {
            ProcessDroppedUris(uris);
            RL_FREE(uris);
            accepted = true;
        }
    }

    FinishDrag(accepted);
}

// Drag and drop: notify the source that the drop is finished
static void FinishDrag(bool accepted)
{
    if ((platform.dragSource != None) && (platform.dragVersion >= 2))
    {
        XEvent reply = { 0 };
        reply.xclient.type = ClientMessage;
        reply.xclient.window = platform.dragSource;
        reply.xclient.message_type = platform.atoms.XdndFinished;
        reply.xclient.format = 32;
        reply.xclient.data.l[0] = (long)platform.window;
        reply.xclient.data.l[1] = accepted? 1 : 0;
        reply.xclient.data.l[2] = accepted? (long)platform.atoms.XdndActionCopy : None;

        XSendEvent(platform.display, platform.dragSource, False, NoEventMask, &reply);
        XFlush(platform.display);
    }

    platform.dragSource = None;
    platform.dragAccepted = false;
}

// Register dropped files from a text/uri-list
static void ProcessDroppedUris(char *uris)
{
    platform.eventReceived = true;
    int count = 0;
    for (char *c = uris; *c != '\0'; c++)
    {
        if ((*c == '\n') || (*c == '\r')) count++;
    }

    char **filepaths = (char **)RL_CALLOC(count + 1, sizeof(char *));
    unsigned int fileCount = 0;

    char *line = strtok(uris, "\r\n");
    while (line != NULL)
    {
        // Only local files are registered: file://[host]/path, comments start with '#'
        if (strncmp(line, "file://", 7) == 0)
        {
            const char *path = strchr(line + 7, '/');

            if (path != NULL)
            {
                char *filepath = (char *)RL_CALLOC(MAX_FILEPATH_LENGTH, sizeof(char));
                int length = 0;

                // Decode percent-encoded characters
                for (int i = 0; (path[i] != '\0') && (length < (MAX_FILEPATH_LENGTH - 1)); i++)
                {
                    if ((path[i] == '%') && isxdigit((unsigned char)path[i + 1]) && isxdigit((unsigned char)path[i + 2]))
                    {
                        char hex[3] = { path[i + 1], path[i + 2], '\0' };
                        char decoded = (char)strtol(hex, NULL, 16);

                        if (decoded != '\0') filepath[length++] = decoded;
                        i += 2;
                    }
                    else filepath[length++] = path[i];
                }

                filepaths[fileCount++] = filepath;
            }
        }

        line = strtok(NULL, "\r\n");
    }

    if (fileCount == 0)
    {
        RL_FREE(filepaths);
        return;
    }

    // In case previous dropped filepaths have not been freed, free them
    if (CORE.Window.dropFileCount > 0)
    {
        for (unsigned int i = 0; i < CORE.Window.dropFileCount; i++) RL_FREE(CORE.Window.dropFilepaths[i]);

        RL_FREE(CORE.Window.dropFilepaths);
    }

    CORE.Window.dropFilepaths = filepaths;
    CORE.Window.dropFileCount = fileCount;
}

#if defined(GRAPHICS_API_OPENGL_SOFTWARE)
// Create software renderer framebuffer image
// NOTE: Image is shared with the X server when possible (local connections), sent through the connection otherwise
static bool CreateSoftwareImage(int width, int height)
{
    if (platform.shmAvailable)
    {
        platform.image = XShmCreateImage(platform.display, platform.visual, platform.depth, ZPixmap, NULL, &platform.shmInfo, width, height);

        if (platform.image != NULL)
        {
            platform.shmInfo.shmid = shmget(IPC_PRIVATE, (size_t)platform.image->bytes_per_line*height, IPC_CREAT | 0600);
            platform.shmInfo.shmaddr = (platform.shmInfo.shmid >= 0)? (char *)shmat(platform.shmInfo.shmid, NULL, 0) : (char *)-1;
            platform.shmInfo.readOnly = False;

            if (platform.shmInfo.shmaddr != (char *)-1)
            {
                platform.image->data = platform.shmInfo.shmaddr;

                platform.errorCode = 0;
                XShmAttach(platform.display, &platform.shmInfo);
                XSync(platform.display, False);
                platform.shmAttached = (platform.errorCode == 0);
            }

            // NOTE: Segment is released when the last process detaches from it
            if (platform.shmInfo.shmid >= 0) shmctl(platform.shmInfo.shmid, IPC_RMID, NULL);

            if (!platform.shmAttached)
            {
                if (platform.shmInfo.shmaddr != (char *)-1) shmdt(platform.shmInfo.shmaddr);
                platform.image->data = NULL;
                XDestroyImage(platform.image);
                platform.image = NULL;
                platform.shmAvailable = false;
            }
        }
    }

    if (platform.image == NULL)
    {
        // NOTE: Image data is released by XDestroyImage()
        char *data = (char *)calloc((size_t)width*height, 4);
        platform.image = XCreateImage(platform.display, platform.visual, platform.depth, ZPixmap, 0, data, width, height, 32, 0);

        if (platform.image == NULL)
        {
            free(data);
            TRACELOG(LOG_WARNING, "DISPLAY: Failed to create framebuffer image");
            return false;
        }

        platform.image->byte_order = LSBFirst;
    }

    if (platform.image->bits_per_pixel != 32)
    {
        TRACELOG(LOG_WARNING, "DISPLAY: Framebuffer image format not supported (%i bits per pixel)", platform.image->bits_per_pixel);
        DestroySoftwareImage();
        return false;
    }

    return true;
}

// Destroy software renderer framebuffer image
static void DestroySoftwareImage(void)
{
    if (platform.image == NULL) return;

    if (platform.shmAttached)
    {
        XShmDetach(platform.display, &platform.shmInfo);
        shmdt(platform.shmInfo.shmaddr);
        platform.image->data = NULL;
        platform.shmAttached = false;
    }

    XDestroyImage(platform.image);
    platform.image = NULL;
}
#endif

// EOF
