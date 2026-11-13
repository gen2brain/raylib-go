/**********************************************************************************************
*
*   rcore_desktop_wayland - Functions to manage window, graphics device and inputs
*
*   PLATFORM: DESKTOP: WAYLAND
*       - Linux (Wayland compositors)
*
*   LIMITATIONS:
*       - Window position can not be queried or set (not exposed by Wayland)
*       - FLAG_WINDOW_TOPMOST is not supported (not exposed by Wayland)
*       - SetMousePosition() only works while the cursor is disabled (locked)
*       - Minimized state is not reported (not exposed by Wayland)
*
*   POSSIBLE IMPROVEMENTS:
*       - Client-side decorations drop shadow
*       - Primary selection (middle-click paste)
*
*   CONFIGURATION:
*       #define WAYLAND_CLIENT_DECORATIONS
*           Request client-side decorations even if the compositor supports server-side decorations,
*           drawn by libdecor when available or by raylib
*       #define WAYLAND_BUILTIN_DECORATIONS
*           Draw client-side decorations with raylib, libdecor is not loaded
*
*   DEPENDENCIES:
*       - libwayland-client, libwayland-egl, libwayland-cursor: Wayland protocol, EGL window and cursor themes
*       - libxkbcommon: Keyboard keymaps and compose sequences
*       - EGL: OpenGL context creation
*       - libdecor (optional, loaded at runtime): Client-side decorations
*       - minigamepad: Gamepad support (evdev)
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

#include <wayland-client.h>         // Wayland client library
#include <wayland-egl.h>            // Wayland EGL window
#include <wayland-cursor.h>         // Wayland cursor themes
#include <xkbcommon/xkbcommon.h>    // Keyboard keymap and state
#include <xkbcommon/xkbcommon-compose.h> // Keyboard compose sequences (dead keys)
// NOTE: glad.h embeds khrplatform.h without KHRONOS_APIENTRY, required by EGL headers
#if !defined(KHRONOS_APIENTRY)
    #define KHRONOS_APIENTRY
#endif
#include <EGL/egl.h>                // EGL library

#include "../external/wayland/xdg-shell-client-protocol.h"
#include "../external/wayland/xdg-decoration-unstable-v1-client-protocol.h"
#include "../external/wayland/fractional-scale-v1-client-protocol.h"
#include "../external/wayland/viewporter-client-protocol.h"
#include "../external/wayland/cursor-shape-v1-client-protocol.h"
#include "../external/wayland/tablet-v2-client-protocol.h"
#include "../external/wayland/pointer-constraints-unstable-v1-client-protocol.h"
#include "../external/wayland/relative-pointer-unstable-v1-client-protocol.h"
#include "../external/wayland/xdg-activation-v1-client-protocol.h"
#include "../external/wayland/xdg-toplevel-icon-v1-client-protocol.h"
#include "../external/wayland/alpha-modifier-v1-client-protocol.h"
#include "../external/wayland/xdg-output-unstable-v1-client-protocol.h"

// Wayland protocols interfaces code
#include "../external/wayland/xdg-shell-client-protocol-code.h"
#include "../external/wayland/xdg-decoration-unstable-v1-client-protocol-code.h"
#include "../external/wayland/fractional-scale-v1-client-protocol-code.h"
#include "../external/wayland/viewporter-client-protocol-code.h"
#include "../external/wayland/cursor-shape-v1-client-protocol-code.h"
#include "../external/wayland/tablet-v2-client-protocol-code.h"
#include "../external/wayland/pointer-constraints-unstable-v1-client-protocol-code.h"
#include "../external/wayland/relative-pointer-unstable-v1-client-protocol-code.h"
#include "../external/wayland/xdg-activation-v1-client-protocol-code.h"
#include "../external/wayland/xdg-toplevel-icon-v1-client-protocol-code.h"
#include "../external/wayland/alpha-modifier-v1-client-protocol-code.h"
#include "../external/wayland/xdg-output-unstable-v1-client-protocol-code.h"

#include <linux/input-event-codes.h> // Required for: BTN_LEFT, BTN_RIGHT...
#include <sys/mman.h>               // Required for: mmap(), munmap()
#include <time.h>                   // Required for: clock_gettime()
#include <stdlib.h>                 // Required for: getenv(), mkstemp(), strtol()
#include <ctype.h>                  // Required for: isxdigit()
#include <poll.h>                   // Required for: poll()
#include <unistd.h>                 // Required for: close(), read(), write(), pipe(), unlink(), ftruncate()
#include <fcntl.h>                  // Required for: fcntl(), FD_CLOEXEC
#include <errno.h>                  // Required for: errno, EINTR, EAGAIN
#include <dlfcn.h>                  // Required for: dlopen(), dlsym(), dlclose()
#include <math.h>                   // Required for: roundf(), ceilf()

// Gamepad support (evdev)
#define MG_MAX_GAMEPADS MAX_GAMEPADS
#define MG_IMPLEMENTATION
#include "../external/RGFW/deps/minigamepad.h"

//----------------------------------------------------------------------------------
// Defines and Macros
//----------------------------------------------------------------------------------
#ifndef EGL_OPENGL_ES3_BIT
    #define EGL_OPENGL_ES3_BIT      0x40
#endif

#define MAX_MONITORS                16      // Maximum number of monitors (wl_output) tracked
#define MAX_ENTERED_MONITORS        16      // Maximum number of monitors the window can be on at the same time
#define MAX_DATA_OFFER_MIME_TYPES   32      // Maximum number of mime types registered per data offer
#define MAX_WINDOW_ICONS            8       // Maximum number of window icon sizes
#define MAX_KEY_REPEATS_PER_FRAME   32      // Maximum number of key repeats processed per frame
#define KEYMAP_SIZE                 256     // Size of the evdev keycodes to raylib keys map

#define XKB_KEYCODE_OFFSET          8       // Offset between evdev scancodes and xkb keycodes
#define KEY_REPEAT_RATE_DEFAULT     25      // Key repeat rate when the compositor does not send it (repeats per second)
#define KEY_REPEAT_DELAY_DEFAULT    600     // Key repeat delay when the compositor does not send it (milliseconds)
#define CURSOR_SIZE_DEFAULT         24      // Cursor theme size when XCURSOR_SIZE is not set (logical pixels)
#define CURSOR_SIZE_MAX             256     // Maximum cursor theme size accepted from XCURSOR_SIZE
#define DATA_OFFER_READ_CHUNK       4096    // Data offer read buffer growth size (bytes)
#define DATA_OFFER_READ_TIMEOUT     1000    // Maximum time waiting for data offer content (milliseconds)

#define DECORATION_TITLEBAR_HEIGHT  32      // Client-side decorations: titlebar height (logical pixels)
#define DECORATION_BORDER_SIZE      6       // Client-side decorations: invisible resize border (logical pixels)
#define DECORATION_CORNER_RADIUS    8       // Client-side decorations: titlebar top corners radius
#define DECORATION_BUTTON_SIZE      24      // Client-side decorations: titlebar buttons size
#define DECORATION_BUTTON_SPACING   8       // Client-side decorations: titlebar buttons spacing
#define DECORATION_TITLE_MARGIN     16      // Client-side decorations: title margin to the buttons
#define DECORATION_ICON_SIZE        20      // Client-side decorations: titlebar window icon size
#define DECORATION_TITLE_FONT_SCALE 1       // Client-side decorations: title default font integer scale (crisp pixel font)
#define DECORATION_CORNER_SIZE      16      // Client-side decorations: resize corners size
#define DECORATION_TOP_EDGE_SIZE    4       // Client-side decorations: titlebar resize edge size
#define DECORATION_DOUBLE_CLICK_TIME 0.4    // Client-side decorations: titlebar double click time (seconds)

#define FRAME_CALLBACK_TIMEOUT      50      // Maximum time waiting for the compositor frame callback on V-Sync (milliseconds)

//----------------------------------------------------------------------------------
// Types and Structures Definition
//----------------------------------------------------------------------------------

// libdecor declarations, the library is loaded at runtime when available
struct libdecor;
struct libdecor_frame;
struct libdecor_state;
struct libdecor_configuration;

enum libdecor_error {
    LIBDECOR_ERROR_COMPOSITOR_INCOMPATIBLE,
    LIBDECOR_ERROR_INVALID_FRAME_CONFIGURATION,
};

enum libdecor_window_state {
    LIBDECOR_WINDOW_STATE_NONE = 0,
    LIBDECOR_WINDOW_STATE_ACTIVE = 1,
    LIBDECOR_WINDOW_STATE_MAXIMIZED = 2,
    LIBDECOR_WINDOW_STATE_FULLSCREEN = 4,
    LIBDECOR_WINDOW_STATE_TILED_LEFT = 8,
    LIBDECOR_WINDOW_STATE_TILED_RIGHT = 16,
    LIBDECOR_WINDOW_STATE_TILED_TOP = 32,
    LIBDECOR_WINDOW_STATE_TILED_BOTTOM = 64,
    LIBDECOR_WINDOW_STATE_SUSPENDED = 128,
};

enum libdecor_capabilities {
    LIBDECOR_ACTION_MOVE = 1,
    LIBDECOR_ACTION_RESIZE = 2,
    LIBDECOR_ACTION_MINIMIZE = 4,
    LIBDECOR_ACTION_FULLSCREEN = 8,
    LIBDECOR_ACTION_CLOSE = 16,
};

struct libdecor_interface {
    void (*error)(struct libdecor *context, enum libdecor_error error, const char *message);
    void (*reserved[10])(void);
};

struct libdecor_frame_interface {
    void (*configure)(struct libdecor_frame *frame, struct libdecor_configuration *configuration, void *userData);
    void (*close)(struct libdecor_frame *frame, void *userData);
    void (*commit)(struct libdecor_frame *frame, void *userData);
    void (*dismiss_popup)(struct libdecor_frame *frame, const char *seatName, void *userData);
    void (*reserved[10])(void);
};

// libdecor functions loaded at runtime
typedef struct {
    void *handle;                   // Library handle
    struct libdecor *context;       // libdecor context
    struct libdecor_frame *frame;   // Window frame, NULL if libdecor is not used
    struct libdecor *(*New)(struct wl_display *display, const struct libdecor_interface *iface);
    void (*Unref)(struct libdecor *context);
    int (*Dispatch)(struct libdecor *context, int timeout);
    struct libdecor_frame *(*Decorate)(struct libdecor *context, struct wl_surface *surface, const struct libdecor_frame_interface *iface, void *userData);
    void (*FrameUnref)(struct libdecor_frame *frame);
    void (*FrameSetAppId)(struct libdecor_frame *frame, const char *appId);
    void (*FrameSetTitle)(struct libdecor_frame *frame, const char *title);
    void (*FrameSetMinimized)(struct libdecor_frame *frame);
    void (*FrameSetFullscreen)(struct libdecor_frame *frame, struct wl_output *output);
    void (*FrameUnsetFullscreen)(struct libdecor_frame *frame);
    void (*FrameMap)(struct libdecor_frame *frame);
    void (*FrameCommit)(struct libdecor_frame *frame, struct libdecor_state *state, struct libdecor_configuration *configuration);
    void (*FrameSetMinContentSize)(struct libdecor_frame *frame, int width, int height);
    void (*FrameSetMaxContentSize)(struct libdecor_frame *frame, int width, int height);
    void (*FrameSetMaximized)(struct libdecor_frame *frame);
    void (*FrameUnsetMaximized)(struct libdecor_frame *frame);
    void (*FrameSetCapabilities)(struct libdecor_frame *frame, enum libdecor_capabilities capabilities);
    void (*FrameUnsetCapabilities)(struct libdecor_frame *frame, enum libdecor_capabilities capabilities);
    void (*FrameSetVisibility)(struct libdecor_frame *frame, bool visible);
    struct xdg_toplevel *(*FrameGetXdgToplevel)(struct libdecor_frame *frame);
    bool (*ConfigurationGetContentSize)(struct libdecor_configuration *configuration, struct libdecor_frame *frame, int *width, int *height);
    bool (*ConfigurationGetWindowState)(struct libdecor_configuration *configuration, enum libdecor_window_state *state);
    struct libdecor_state *(*StateNew)(int width, int height);
    void (*StateFree)(struct libdecor_state *state);
} LibdecorData;

// Decorations mode
typedef enum {
    DECORATION_MODE_NONE = 0,       // No decorations (undecorated or fullscreen)
    DECORATION_MODE_SERVER,         // Server-side decorations (xdg-decoration)
    DECORATION_MODE_LIBDECOR,       // Client-side decorations drawn by libdecor
    DECORATION_MODE_BUILTIN,        // Client-side decorations drawn by raylib
} DecorationMode;

// Client-side decorations titlebar buttons
typedef enum {
    DECORATION_BUTTON_NONE = 0,
    DECORATION_BUTTON_MINIMIZE,
    DECORATION_BUTTON_MAXIMIZE,
    DECORATION_BUTTON_CLOSE,
} DecorationButton;

// Client-side decorations parts (subsurfaces)
typedef enum {
    DECORATION_PART_TITLEBAR = 0,   // Titlebar, above the window content
    DECORATION_PART_LEFT,           // Left resize border
    DECORATION_PART_RIGHT,          // Right resize border
    DECORATION_PART_BOTTOM,         // Bottom resize border
    DECORATION_PART_COUNT
} DecorationPartType;

// Shared memory buffer
typedef struct {
    struct wl_buffer *buffer;       // Wayland buffer
    unsigned int *pixels;           // Pixel data (ARGB8888, premultiplied alpha)
    int width;                      // Buffer width in pixels
    int height;                     // Buffer height in pixels
    size_t size;                    // Mapped memory size in bytes
    bool busy;                      // Buffer is in use by the compositor
} ShmBuffer;

// Client-side decorations part (subsurface)
typedef struct {
    struct wl_surface *surface;     // Decoration part surface
    struct wl_subsurface *subsurface; // Decoration part subsurface (relative to the window surface)
    struct wp_viewport *viewport;   // Decoration part viewport (logical size)
    ShmBuffer buffers[2];           // Double buffered, borders only use the first one
} DecorationPart;

// Monitor info (wl_output)
typedef struct {
    struct wl_output *output;       // Wayland output
    struct zxdg_output_v1 *xdgOutput; // Output logical info
    unsigned int globalName;        // Registry global name
    char name[128];                 // Output name
    int x;                          // Logical position x
    int y;                          // Logical position y
    int logicalWidth;               // Logical width
    int logicalHeight;              // Logical height
    int width;                      // Current mode width in pixels
    int height;                     // Current mode height in pixels
    int physicalWidth;              // Physical width in millimetres
    int physicalHeight;             // Physical height in millimetres
    int refreshRate;                // Current mode refresh rate in mHz
    int scale;                      // Output integer scale
} MonitorInfo;

// Data offer (clipboard selection or drag and drop)
typedef struct {
    struct wl_data_offer *offer;    // Wayland data offer
    char *mimeTypes[MAX_DATA_OFFER_MIME_TYPES]; // Offered mime types
    int mimeTypeCount;              // Offered mime types count
} DataOffer;

typedef struct {
    // Wayland globals
    struct wl_display *display;     // Wayland display connection
    struct wl_registry *registry;   // Wayland globals registry
    struct wl_compositor *compositor; // Surfaces compositor
    struct wl_subcompositor *subcompositor; // Subsurfaces compositor (client-side decorations)
    struct wl_shm *shm;             // Shared memory buffers
    struct wl_seat *seat;           // Input devices seat
    struct wl_data_device_manager *dataDeviceManager; // Clipboard and drag and drop
    struct xdg_wm_base *wmBase;     // Desktop windows (xdg-shell)
    struct zxdg_decoration_manager_v1 *decorationManager; // Server-side decorations (optional)
    struct wp_fractional_scale_manager_v1 *fractionalScaleManager; // Fractional scaling (optional)
    struct wp_viewporter *viewporter; // Surface scaling and cropping (optional)
    struct wp_cursor_shape_manager_v1 *cursorShapeManager; // Cursor shapes (optional)
    struct zwp_pointer_constraints_v1 *pointerConstraints; // Pointer lock (optional)
    struct zwp_relative_pointer_manager_v1 *relativePointerManager; // Relative pointer motion (optional)
    struct xdg_activation_v1 *activation; // Window focus requests (optional)
    struct xdg_toplevel_icon_manager_v1 *iconManager; // Window icons (optional)
    struct wp_alpha_modifier_v1 *alphaModifier; // Window opacity (optional)
    struct zxdg_output_manager_v1 *outputManager; // Outputs logical position and size (optional)
    unsigned int compositorVersion; // Bound wl_compositor version
    unsigned int seatVersion;       // Bound wl_seat version

    // Window objects
    struct wl_surface *surface;     // Window surface
    struct xdg_surface *xdgSurface; // Window xdg surface
    struct xdg_toplevel *xdgToplevel; // Window toplevel
    struct zxdg_toplevel_decoration_v1 *toplevelDecoration; // Window decorations mode negotiation
    struct wp_fractional_scale_v1 *fractionalScale; // Window preferred fractional scale
    struct wp_viewport *viewport;   // Window viewport (buffer to logical size)
    struct wp_alpha_modifier_surface_v1 *alphaModifierSurface; // Window opacity modifier
    struct wl_callback *frameCallback; // Pending frame callback (V-Sync)
    struct wl_event_queue *frameQueue; // Frame callback events queue, dispatched without processing input events
    struct wl_egl_window *eglWindow; // EGL native window
    char appId[256];                // Window application id (executable name)
    char title[256];                // Window title copy (CORE.Window.title is not owned by raylib)

    // EGL data
    EGLDisplay eglDisplay;          // Native display device (physical screen connection)
    EGLContext eglContext;          // Graphic context, mode in which drawing can be done
    EGLSurface eglSurface;          // Surface to draw on, framebuffers (connected to context)
    EGLConfig eglConfig;            // Graphic config

    // Window state
    int width;                      // Window content width (logical pixels)
    int height;                     // Window content height (logical pixels)
    int bufferWidth;                // Window buffer width (pixels)
    int bufferHeight;               // Window buffer height (pixels)
    float scale;                    // Window scale factor
    int preferredBufferScale;       // Integer scale preferred by the compositor
    int fractionalScale120;         // Fractional scale preferred by the compositor (x120)
    int bufferScale;                // Integer buffer scale for the next buffer (without viewporter)
    int appliedBufferScale;         // Integer buffer scale set on the window surface
    int windowedWidth;              // Last floating window content width (restored when leaving fullscreen/maximized)
    int windowedHeight;             // Last floating window content height (restored when leaving fullscreen/maximized)
    int pendingWidth;               // Pending configure width
    int pendingHeight;              // Pending configure height
    bool pendingActivated;          // Pending configure state: activated
    bool pendingMaximized;          // Pending configure state: maximized
    bool pendingFullscreen;         // Pending configure state: fullscreen
    bool configured;                // First configure event received
    bool visible;                   // Window is mapped (FLAG_WINDOW_HIDDEN not set)
    bool closeRequested;            // Window close requested by user/compositor
    bool connectionLost;            // Connection to the compositor lost (protocol or socket error)
    bool activatedOnce;             // Window was activated by the compositor at least once
    bool fullscreenRequested;       // Last fullscreen state requested to the compositor
    bool fullscreenPending;         // Last fullscreen request not yet confirmed by a configure event
    bool eventReceived;             // Window or input event received (ends event waiting)
    bool initialized;               // Platform initialization finished (rlgl ready on next events)
    bool fullscreen;                // Window is fullscreen
    bool maximized;                 // Window is maximized
    int fullscreenMonitor;          // Monitor requested for fullscreen
    ShmBuffer iconBuffers[MAX_WINDOW_ICONS]; // Window icon buffers
    int iconBufferCount;            // Window icon buffers count

    // Monitors
    MonitorInfo monitors[MAX_MONITORS]; // Connected monitors
    int monitorCount;               // Connected monitors count
    struct wl_output *enteredOutputs[MAX_ENTERED_MONITORS]; // Monitors the window is on, last entered is current
    int enteredOutputCount;         // Monitors the window is on count

    // Pointer
    struct wl_pointer *pointer;     // Seat pointer
    struct wl_surface *pointerFocus; // Surface under the pointer (window or decoration part)
    unsigned int pointerEnterSerial; // Pointer enter serial (required to set the cursor)
    double pointerX;                // Pointer position x on the focused surface (logical pixels)
    double pointerY;                // Pointer position y on the focused surface (logical pixels)
    float wheelMoveX;               // Mouse wheel horizontal move accumulated in the current frame
    float wheelMoveY;               // Mouse wheel vertical move accumulated in the current frame
    bool wheelDiscrete;             // Mouse wheel steps received in the current pointer frame
    struct wp_cursor_shape_device_v1 *cursorShapeDevice; // Pointer cursor shape device
    struct wl_cursor_theme *cursorTheme; // Cursor theme (cursor-shape fallback)
    int cursorThemeScale;           // Cursor theme loaded scale
    struct wl_surface *cursorSurface; // Cursor theme image surface
    struct zwp_locked_pointer_v1 *lockedPointer; // Locked pointer (disabled cursor)
    struct zwp_relative_pointer_v1 *relativePointer; // Relative pointer motion (disabled cursor)

    // Keyboard
    struct wl_keyboard *keyboard;   // Seat keyboard
    struct xkb_context *xkbContext; // Keyboard xkb context
    struct xkb_keymap *xkbKeymap;   // Keyboard keymap sent by the compositor
    struct xkb_state *xkbState;     // Keyboard modifiers and layout state
    struct xkb_compose_table *xkbComposeTable; // Compose sequences table (locale)
    struct xkb_compose_state *xkbComposeState; // Compose sequence state (dead keys)
    int repeatRate;                 // Key repeat rate (characters per second)
    int repeatDelay;                // Key repeat delay (milliseconds)
    unsigned int repeatScancode;    // Key being repeated (evdev scancode)
    double repeatTime;              // Next key repeat time
    unsigned int inputSerial;       // Last input event serial (required by some requests)

    // Touch
    struct wl_touch *touch;         // Seat touch device
    int touchIds[MAX_TOUCH_POINTS];         // Touch point ids per slot, -1 if slot not used
    Vector2 touchPositions[MAX_TOUCH_POINTS]; // Touch point positions per slot

    // Clipboard and drag and drop
    struct wl_data_device *dataDevice; // Seat data device
    DataOffer *selectionOffer;      // Current clipboard selection offer
    DataOffer *dragOffer;           // Current drag and drop offer
    struct wl_data_source *dataSource; // Clipboard selection owned by this client
    char *clipboardText;            // Text owned by this client as clipboard selection
    char *clipboardReceived;        // Text received from the clipboard selection

    // Decorations
    DecorationMode decorationMode;  // Current decorations mode
    int decorationModeSent;         // Last xdg-decoration mode requested to the compositor
    bool decorationModeChanged;     // Decorations mode requested, next configure keeps the content size
    int sentMinWidth;               // Last minimum width sent to the compositor
    int sentMinHeight;              // Last minimum height sent to the compositor
    int sentMaxWidth;               // Last maximum width sent to the compositor
    int sentMaxHeight;              // Last maximum height sent to the compositor
    int sizeLimitsResizable;        // Last resizable capability sent to the compositor (-1 if not sent)
    int decorationVisible;          // Last libdecor decorations visibility (-1 if not set)
    bool decorationRequested;       // Decorations requested (FLAG_WINDOW_UNDECORATED not set)
    bool preferBuiltinDecorations;  // Use raylib decorations instead of libdecor
    bool preferClientDecorations;   // Request client-side decorations even if server-side are available
    DecorationPart decorationParts[DECORATION_PART_COUNT]; // Built-in decorations parts (titlebar and borders)
    DecorationButton decorationHover; // Titlebar button under the pointer
    DecorationButton decorationPressed; // Titlebar button pressed
    double decorationLastClickTime; // Titlebar last click time (double-click to maximize)
    bool decorationDirty;           // Titlebar requires redraw
    bool frameCommitted;            // Window surface has content (first frame presented), required before committing other state
    bool decorationTitleReady;      // Titlebar was drawn with the default font loaded
    LibdecorData libdecor;          // libdecor library and window frame

#if defined(GRAPHICS_API_OPENGL_SOFTWARE)
    // Software renderer
    ShmBuffer softwareBuffers[2];   // Double buffered shared memory buffers
    unsigned char *softwarePixels;  // Software renderer framebuffer copy (RGBA)
    int softwareWidth;              // Software renderer framebuffer width
    int softwareHeight;             // Software renderer framebuffer height
#endif

    // Gamepads
    mg_gamepads minigamepad;        // Gamepads (evdev)
} PlatformData;

//----------------------------------------------------------------------------------
// Global Variables Definition
//----------------------------------------------------------------------------------
extern CoreData CORE;                   // Global CORE state context

static PlatformData platform = { 0 };   // Platform specific data

// Builtin titlebar buttons, from right to left
static const DecorationButton titlebarButtons[3] = { DECORATION_BUTTON_CLOSE, DECORATION_BUTTON_MAXIMIZE, DECORATION_BUTTON_MINIMIZE };

// Clipboard text mime types, in order of preference
static const char *textMimeTypes[] = { "text/plain;charset=utf-8", "text/plain", "UTF8_STRING", "STRING", "TEXT" };

// Map evdev keycodes (linux/input-event-codes.h) to raylib KeyboardKey values
// NOTE: raylib key values match GLFW values, mapping based on GLFW Wayland backend
static const short linuxToRaylibMap[KEYMAP_SIZE] = {
      0, 256,  49,  50,  51,  52,  53,  54,
     55,  56,  57,  48,  45,  61, 259, 258,
     81,  87,  69,  82,  84,  89,  85,  73,
     79,  80,  91,  93, 257, 341,  65,  83,
     68,  70,  71,  72,  74,  75,  76,  59,
     39,  96, 340,  92,  90,  88,  67,  86,
     66,  78,  77,  44,  46,  47, 344, 332,
    342,  32, 280, 290, 291, 292, 293, 294,
    295, 296, 297, 298, 299, 282, 281, 327,
    328, 329, 333, 324, 325, 326, 334, 321,
    322, 323, 320, 330,   0,   0, 162, 300,
    301,   0,   0,   0,   0,   0,   0,   0,
    335, 345, 331,   0, 346,   0, 268, 265,
    266, 263, 262, 269, 264, 267, 260, 261,
      0,   0,   0,   0,   0, 336,   0, 284,
      0,   0,   0,   0,   0, 343, 347, 348,
      0,   0,   0,   0,   0,   0,   0,   0,
      0,   0,   0,   0,   0,   0,   0,   0,
      0,   0,   0,   0,   0,   0,   0,   0,
      0,   0,   0,   0,   0,   0,   0,   0,
      0,   0,   0,   0,   0,   0,   0,   0,
      0,   0,   0,   0,   0,   0,   0,   0,
      0,   0,   0,   0,   0,   0,   0, 302,
    303, 304, 305, 306, 307, 308, 309, 310,
    311, 312, 313,   0,   0,   0,   0,   0,
      0,   0,   0,   0,   0,   0,   0,   0,
      0,   0, 283,   0,   0,   0,   0,   0,
      0,   0,   0,   0,   0,   0,   0,   0,
      0,   0,   0,   0,   0,   0,   0,   0,
      0,   0,   0,   0,   0,   0,   0,   0,
      0,   0,   0,   0,   0,   0,   0,   0,
      0,   0,   0,   0,   0,   0,   0,   0
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

static bool InitGraphicsContext(void);                          // Initialize EGL display, context and window surface
static void DispatchEvents(int timeout);                        // Read and dispatch Wayland events, timeout in milliseconds (-1 blocks)
static void WaitFrameCallback(int timeout);                     // Wait for the compositor frame callback, timeout in milliseconds
static bool ReadEvents(int timeout);                            // Read Wayland events after a prepared read, timeout in milliseconds (-1 blocks)
static void HandleConnectionError(void);                        // Register connection to the compositor lost, the window should close
static void CreateShellObjects(void);                           // Create window role objects (xdg_surface, xdg_toplevel, decorations)
static void DestroyShellObjects(void);                          // Destroy window role objects
static void ApplyConfigure(int width, int height, bool activated, bool maximized, bool fullscreen); // Apply a configure event
static void UpdateWindowScale(void);                            // Update window scale from compositor preferences
static void UpdateWindowBuffers(void);                          // Apply window size and scale to the framebuffer
static void UpdateWindowGeometry(void);                         // Update window geometry and size limits
static void RequestFullscreen(bool fullscreen, struct wl_output *output); // Request fullscreen state to the compositor
static void UpdateWindowIcon(void);                             // Set window icon from the icon buffers, default icon if none
static int GetMonitorIndex(struct wl_output *output);           // Get monitor index for an output
static struct wl_output *GetMonitorOutput(int monitor);         // Get output for a monitor index
static bool CreateShmBuffer(ShmBuffer *buffer, int width, int height); // Create shared memory buffer (ARGB8888)
static void DestroyShmBuffer(ShmBuffer *buffer);                // Destroy shared memory buffer
static void ReleasePointer(struct wl_pointer *pointer);         // Release seat pointer (destroy on old compositors)
static void ReleaseKeyboard(struct wl_keyboard *keyboard);      // Release seat keyboard (destroy on old compositors)
static void ReleaseTouch(struct wl_touch *touch);               // Release seat touch device (destroy on old compositors)
static void ReleaseOutput(struct wl_output *output);            // Release output (destroy on old compositors)
static void HandleKey(unsigned int scancode, bool pressed);     // Register key state change
static void HandleKeyChar(unsigned int scancode);               // Register character input for a key, considering keyboard layout and compose sequences
static void ProcessKeyRepeat(void);                             // Generate key repeat events
static void HandleMouseButton(int button, bool pressed);        // Register mouse button state change
static void HandleMouseMove(double x, double y);                // Register mouse position change
static void HandleTouch(int touchAction, int slot, int id, double x, double y); // Register touch event
static void UpdateCursor(void);                                 // Update cursor image for current pointer focus
static void SetCursorShape(int cursor);                         // Set cursor shape (MouseCursor) on pointer focus
static void SetCursorImage(unsigned int shape, const char *name, const char *fallbackName); // Set cursor image, cursor-shape protocol or cursor theme
static void LockPointer(bool lock);                             // Lock or unlock pointer (disabled cursor)
static void DestroyDataOffer(DataOffer *offer);                 // Destroy data offer and its mime types
static bool DataOfferHasMimeType(DataOffer *offer, const char *mimeType); // Check data offer mime type
static char *ReadDataOffer(DataOffer *offer, const char *mimeType, int *size); // Read data offer content, returned data is allocated and zero-terminated
static void ProcessDroppedUris(char *uris);                     // Register dropped files from a text/uri-list
static bool LoadLibdecor(void);                                 // Load libdecor library
static void UnloadLibdecor(void);                               // Unload libdecor library
static void CreateDecorations(void);                            // Create client-side decorations parts (raylib)
static void DestroyDecorations(void);                           // Destroy client-side decorations parts (raylib)
static Rectangle GetDecorationPartBounds(int part);             // Get decorations part bounds (logical pixels, relative to the window content)
static void AttachDecorationBuffer(DecorationPart *part, ShmBuffer *buffer, Rectangle bounds, int bufferScale); // Attach decorations part buffer, considering window scale
static void CommitWindowSurface(void);                          // Commit window surface state, deferred to the next frame until the first frame was presented
static int GetTitlebarHeight(void);                             // Get builtin titlebar height included in the window geometry
static void UpdateDecorations(void);                            // Update client-side decorations size and position
static void BlendDecorationPixel(ShmBuffer *buffer, int x, int y, unsigned int color, float coverage); // Blend a color into a decorations buffer pixel
static void DrawDecorationLine(ShmBuffer *buffer, float x0, float y0, float x1, float y1, float thick, unsigned int color); // Draw an anti-aliased line into a decorations buffer
static void DrawDecorationCircle(ShmBuffer *buffer, float centerX, float centerY, float radius, unsigned int color); // Draw an anti-aliased circle into a decorations buffer
static void DrawDecorationIcon(ShmBuffer *buffer, int posX, int posY, int size); // Draw window icon into a decorations buffer, scaled down from the closest icon size
#if SUPPORT_MODULE_RTEXT
static void DrawDecorationText(ShmBuffer *buffer, Font font, const char *text, int posX, int posY, int textScale, unsigned int color); // Draw decorations title text
#endif
static void DrawDecorations(void);                              // Draw client-side decorations titlebar
static int GetDecorationPart(struct wl_surface *surface);       // Get decoration part for a surface, -1 if none
static DecorationButton GetDecorationButton(double x, double y); // Get titlebar button at position (logical pixels, relative to the titlebar)
static unsigned int GetDecorationEdge(int part, double x, double y); // Get resize edge at decorations position (logical pixels, relative to the part)
static void HandleDecorationButton(int part, unsigned int button, bool pressed); // Process pointer button on decorations

// Window callbacks events
static void RegistryGlobalCallback(void *data, struct wl_registry *registry, uint32_t name, const char *interface, uint32_t version); // Registry: global object announced
static void RegistryGlobalRemoveCallback(void *data, struct wl_registry *registry, uint32_t name); // Registry: global object removed
static void WmBasePingCallback(void *data, struct xdg_wm_base *wmBase, uint32_t serial);        // xdg_wm_base: compositor checks if the client is responsive
static void SurfaceEnterCallback(void *data, struct wl_surface *surface, struct wl_output *output); // Surface: window entered an output
static void SurfaceLeaveCallback(void *data, struct wl_surface *surface, struct wl_output *output); // Surface: window left an output
static void SurfacePreferredBufferScaleCallback(void *data, struct wl_surface *surface, int32_t factor); // Surface: compositor preferred buffer scale (integer)
static void SurfacePreferredBufferTransformCallback(void *data, struct wl_surface *surface, uint32_t transform); // Surface: compositor preferred buffer transform
static void XdgSurfaceConfigureCallback(void *data, struct xdg_surface *xdgSurface, uint32_t serial); // xdg_surface: configure sequence done, apply pending state
static void XdgToplevelConfigureCallback(void *data, struct xdg_toplevel *toplevel, int32_t width, int32_t height,
    struct wl_array *states);                                                           // xdg_toplevel: suggested size and state
static void XdgToplevelCloseCallback(void *data, struct xdg_toplevel *toplevel);                // xdg_toplevel: window close requested
static void XdgToplevelConfigureBoundsCallback(void *data, struct xdg_toplevel *toplevel, int32_t width, int32_t height); // xdg_toplevel: recommended maximum size
static void XdgToplevelWmCapabilitiesCallback(void *data, struct xdg_toplevel *toplevel, struct wl_array *capabilities); // xdg_toplevel: compositor capabilities
static void ToplevelDecorationConfigureCallback(void *data, struct zxdg_toplevel_decoration_v1 *decoration, uint32_t mode); // xdg-decoration: mode selected by the compositor
static void FractionalScalePreferredScaleCallback(void *data, struct wp_fractional_scale_v1 *fractionalScale, uint32_t scale); // Fractional scale: preferred scale (x120)
static void FrameDoneCallback(void *data, struct wl_callback *callback, uint32_t time); // Frame callback: compositor is ready for a new frame
static void OutputGeometryCallback(void *data, struct wl_output *output, int32_t x, int32_t y, int32_t physicalWidth, int32_t physicalHeight,
    int32_t subpixel, const char *make, const char *model, int32_t transform);          // Output: physical properties
static void OutputModeCallback(void *data, struct wl_output *output, uint32_t flags, int32_t width, int32_t height, int32_t refresh); // Output: video mode
static void OutputDoneCallback(void *data, struct wl_output *output);                           // Output: all properties sent
static void OutputScaleCallback(void *data, struct wl_output *output, int32_t factor);          // Output: integer scale factor
static void OutputNameCallback(void *data, struct wl_output *output, const char *name);         // Output: name
static void OutputDescriptionCallback(void *data, struct wl_output *output, const char *description); // Output: description
static void XdgOutputLogicalPositionCallback(void *data, struct zxdg_output_v1 *xdgOutput, int32_t x, int32_t y); // xdg_output: logical position (compositor space)
static void XdgOutputLogicalSizeCallback(void *data, struct zxdg_output_v1 *xdgOutput, int32_t width, int32_t height); // xdg_output: logical size (compositor space)
static void XdgOutputDoneCallback(void *data, struct zxdg_output_v1 *xdgOutput);                // xdg_output: all properties sent
static void XdgOutputNameCallback(void *data, struct zxdg_output_v1 *xdgOutput, const char *name); // xdg_output: name
static void XdgOutputDescriptionCallback(void *data, struct zxdg_output_v1 *xdgOutput, const char *description); // xdg_output: description
static void ActivationTokenDoneCallback(void *data, struct xdg_activation_token_v1 *token, const char *tokenString); // Activation token ready, request window activation
static void BufferReleaseCallback(void *data, struct wl_buffer *buffer);                        // Shared memory buffer released by the compositor
static void LibdecorErrorCallback(struct libdecor *context, enum libdecor_error error, const char *message); // libdecor: error
static void LibdecorFrameConfigureCallback(struct libdecor_frame *frame, struct libdecor_configuration *configuration, void *data); // libdecor: suggested size and state
static void LibdecorFrameCloseCallback(struct libdecor_frame *frame, void *data);               // libdecor: window close requested
static void LibdecorFrameCommitCallback(struct libdecor_frame *frame, void *data);              // libdecor: window surface commit requested
static void LibdecorFrameDismissPopupCallback(struct libdecor_frame *frame, const char *seatName, void *data); // libdecor: popup dismissed

// Input callbacks events
static void SeatCapabilitiesCallback(void *data, struct wl_seat *seat, uint32_t capabilities);  // Seat: input capabilities changed
static void SeatNameCallback(void *data, struct wl_seat *seat, const char *name);               // Seat: name
static void PointerEnterCallback(void *data, struct wl_pointer *pointer, uint32_t serial, struct wl_surface *surface, wl_fixed_t x, wl_fixed_t y); // Pointer: entered a surface
static void PointerLeaveCallback(void *data, struct wl_pointer *pointer, uint32_t serial, struct wl_surface *surface); // Pointer: left a surface
static void PointerMotionCallback(void *data, struct wl_pointer *pointer, uint32_t time, wl_fixed_t x, wl_fixed_t y); // Pointer: moved over a surface
static void PointerButtonCallback(void *data, struct wl_pointer *pointer, uint32_t serial, uint32_t time, uint32_t button, uint32_t state); // Pointer: button pressed or released
static void PointerAxisCallback(void *data, struct wl_pointer *pointer, uint32_t time, uint32_t axis, wl_fixed_t value); // Pointer: scroll (continuous value)
static void PointerFrameCallback(void *data, struct wl_pointer *pointer);                       // Pointer: end of a group of pointer events
static void PointerAxisSourceCallback(void *data, struct wl_pointer *pointer, uint32_t source); // Pointer: scroll source
static void PointerAxisStopCallback(void *data, struct wl_pointer *pointer, uint32_t time, uint32_t axis); // Pointer: scroll stopped
static void PointerAxisDiscreteCallback(void *data, struct wl_pointer *pointer, uint32_t axis, int32_t discrete); // Pointer: discrete scroll steps (deprecated by axis_value120)
static void PointerAxisValue120Callback(void *data, struct wl_pointer *pointer, uint32_t axis, int32_t value120); // Pointer: high-resolution scroll steps, 120 is one wheel step
static void PointerAxisRelativeDirectionCallback(void *data, struct wl_pointer *pointer, uint32_t axis, uint32_t direction); // Pointer: scroll direction relative to the motion
static void RelativePointerMotionCallback(void *data, struct zwp_relative_pointer_v1 *relativePointer, uint32_t timeHigh, uint32_t timeLow,
    wl_fixed_t dx, wl_fixed_t dy, wl_fixed_t dxUnaccel, wl_fixed_t dyUnaccel);          // Relative pointer: unaccelerated motion while the cursor is disabled
static void KeyboardKeymapCallback(void *data, struct wl_keyboard *keyboard, uint32_t format, int32_t fd, uint32_t size); // Keyboard: keymap provided by the compositor
static void KeyboardEnterCallback(void *data, struct wl_keyboard *keyboard, uint32_t serial, struct wl_surface *surface,
    struct wl_array *keys);                                                             // Keyboard: window got keyboard focus
static void KeyboardLeaveCallback(void *data, struct wl_keyboard *keyboard, uint32_t serial, struct wl_surface *surface); // Keyboard: window lost keyboard focus
static void KeyboardKeyCallback(void *data, struct wl_keyboard *keyboard, uint32_t serial, uint32_t time, uint32_t key, uint32_t state); // Keyboard: key pressed or released
static void KeyboardModifiersCallback(void *data, struct wl_keyboard *keyboard, uint32_t serial, uint32_t depressed, uint32_t latched,
    uint32_t locked, uint32_t group);                                                   // Keyboard: modifiers state changed
static void KeyboardRepeatInfoCallback(void *data, struct wl_keyboard *keyboard, int32_t rate, int32_t delay); // Keyboard: key repeat rate and delay
static void TouchDownCallback(void *data, struct wl_touch *touch, uint32_t serial, uint32_t time, struct wl_surface *surface, int32_t id,
    wl_fixed_t x, wl_fixed_t y);                                                        // Touch: new touch point
static void TouchUpCallback(void *data, struct wl_touch *touch, uint32_t serial, uint32_t time, int32_t id); // Touch: touch point released
static void TouchMotionCallback(void *data, struct wl_touch *touch, uint32_t time, int32_t id, wl_fixed_t x, wl_fixed_t y); // Touch: touch point moved
static void TouchFrameCallback(void *data, struct wl_touch *touch);                             // Touch: end of a group of touch events
static void TouchCancelCallback(void *data, struct wl_touch *touch);                            // Touch: touch session cancelled by the compositor
static void TouchShapeCallback(void *data, struct wl_touch *touch, int32_t id, wl_fixed_t major, wl_fixed_t minor); // Touch: touch point shape
static void TouchOrientationCallback(void *data, struct wl_touch *touch, int32_t id, wl_fixed_t orientation); // Touch: touch point orientation
static void DataOfferOfferCallback(void *data, struct wl_data_offer *offer, const char *mimeType); // Data offer: mime type offered
static void DataOfferSourceActionsCallback(void *data, struct wl_data_offer *offer, uint32_t actions); // Data offer: drag and drop actions supported by the source
static void DataOfferActionCallback(void *data, struct wl_data_offer *offer, uint32_t action);  // Data offer: drag and drop action selected by the compositor
static void DataDeviceDataOfferCallback(void *data, struct wl_data_device *device, struct wl_data_offer *offer); // Data device: new data offer, followed by its mime types
static void DataDeviceEnterCallback(void *data, struct wl_data_device *device, uint32_t serial, struct wl_surface *surface,
    wl_fixed_t x, wl_fixed_t y, struct wl_data_offer *offer);                           // Data device: drag and drop entered a surface
static void DataDeviceLeaveCallback(void *data, struct wl_data_device *device);                 // Data device: drag and drop left the surface
static void DataDeviceMotionCallback(void *data, struct wl_data_device *device, uint32_t time, wl_fixed_t x, wl_fixed_t y); // Data device: drag and drop moved over the surface
static void DataDeviceDropCallback(void *data, struct wl_data_device *device);                  // Data device: drag and drop dropped on the surface
static void DataDeviceSelectionCallback(void *data, struct wl_data_device *device, struct wl_data_offer *offer); // Data device: clipboard selection changed
static void DataSourceTargetCallback(void *data, struct wl_data_source *source, const char *mimeType); // Data source: target mime type accepted (drag and drop)
static void DataSourceSendCallback(void *data, struct wl_data_source *source, const char *mimeType, int32_t fd); // Data source: clipboard content requested by another client
static void DataSourceCancelledCallback(void *data, struct wl_data_source *source);             // Data source: clipboard ownership lost
static void DataSourceDndDropPerformedCallback(void *data, struct wl_data_source *source);      // Data source: drag and drop drop performed
static void DataSourceDndFinishedCallback(void *data, struct wl_data_source *source);           // Data source: drag and drop finished
static void DataSourceActionCallback(void *data, struct wl_data_source *source, uint32_t action); // Data source: drag and drop action selected

// Wayland events listeners
static const struct wl_registry_listener registryListener = { RegistryGlobalCallback, RegistryGlobalRemoveCallback };
static const struct xdg_wm_base_listener wmBaseListener = { WmBasePingCallback };
static const struct wl_surface_listener surfaceListener = {
    SurfaceEnterCallback, SurfaceLeaveCallback, SurfacePreferredBufferScaleCallback,
    SurfacePreferredBufferTransformCallback
};
static const struct xdg_surface_listener xdgSurfaceListener = { XdgSurfaceConfigureCallback };
static const struct xdg_toplevel_listener xdgToplevelListener = {
    XdgToplevelConfigureCallback, XdgToplevelCloseCallback, XdgToplevelConfigureBoundsCallback, XdgToplevelWmCapabilitiesCallback
};
static const struct zxdg_toplevel_decoration_v1_listener toplevelDecorationListener = { ToplevelDecorationConfigureCallback };
static const struct wp_fractional_scale_v1_listener fractionalScaleListener = { FractionalScalePreferredScaleCallback };
static const struct wl_callback_listener frameCallbackListener = { FrameDoneCallback };
static const struct wl_output_listener outputListener = {
    OutputGeometryCallback, OutputModeCallback, OutputDoneCallback, OutputScaleCallback, OutputNameCallback,
    OutputDescriptionCallback
};
static const struct zxdg_output_v1_listener xdgOutputListener = {
    XdgOutputLogicalPositionCallback, XdgOutputLogicalSizeCallback, XdgOutputDoneCallback, XdgOutputNameCallback, XdgOutputDescriptionCallback
};
static const struct wl_seat_listener seatListener = { SeatCapabilitiesCallback, SeatNameCallback };
static const struct wl_pointer_listener pointerListener = {
    PointerEnterCallback, PointerLeaveCallback, PointerMotionCallback, PointerButtonCallback, PointerAxisCallback, PointerFrameCallback,
    PointerAxisSourceCallback, PointerAxisStopCallback, PointerAxisDiscreteCallback, PointerAxisValue120Callback, PointerAxisRelativeDirectionCallback
};
static const struct zwp_relative_pointer_v1_listener relativePointerListener = { RelativePointerMotionCallback };
static const struct wl_keyboard_listener keyboardListener = {
    KeyboardKeymapCallback, KeyboardEnterCallback, KeyboardLeaveCallback, KeyboardKeyCallback, KeyboardModifiersCallback, KeyboardRepeatInfoCallback
};
static const struct wl_touch_listener touchListener = {
    TouchDownCallback, TouchUpCallback, TouchMotionCallback, TouchFrameCallback, TouchCancelCallback,
    TouchShapeCallback, TouchOrientationCallback
};
static const struct wl_data_offer_listener dataOfferListener = { DataOfferOfferCallback, DataOfferSourceActionsCallback, DataOfferActionCallback };
static const struct wl_data_device_listener dataDeviceListener = {
    DataDeviceDataOfferCallback, DataDeviceEnterCallback, DataDeviceLeaveCallback, DataDeviceMotionCallback, DataDeviceDropCallback, DataDeviceSelectionCallback
};
static const struct wl_data_source_listener dataSourceListener = {
    DataSourceTargetCallback, DataSourceSendCallback, DataSourceCancelledCallback, DataSourceDndDropPerformedCallback, DataSourceDndFinishedCallback, DataSourceActionCallback
};
static const struct xdg_activation_token_v1_listener activationTokenListener = { ActivationTokenDoneCallback };
static const struct wl_buffer_listener bufferListener = { BufferReleaseCallback };
static struct libdecor_interface libdecorInterface = { LibdecorErrorCallback };
static struct libdecor_frame_interface libdecorFrameInterface = {
    LibdecorFrameConfigureCallback, LibdecorFrameCloseCallback, LibdecorFrameCommitCallback,
    LibdecorFrameDismissPopupCallback
};

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
        struct wl_output *output = GetMonitorOutput(monitor);

        if (output != NULL)
        {
            CORE.Window.display.width = GetMonitorWidth(monitor);
            CORE.Window.display.height = GetMonitorHeight(monitor);
        }

        platform.fullscreenMonitor = monitor;
        FLAG_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE);
        FLAG_SET(CORE.Window.flags, FLAG_WINDOW_UNDECORATED);

        RequestFullscreen(true, output);
    }
    else
    {
        // Restore previous window size
        CORE.Window.position = CORE.Window.previousPosition;
        FLAG_CLEAR(CORE.Window.flags, FLAG_FULLSCREEN_MODE);
        if (platform.decorationRequested) FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_UNDECORATED);

        RequestFullscreen(false, NULL);
    }

    // NOTE: New window size is applied on next configure event
    wl_display_flush(platform.display);
}

// Toggle borderless windowed mode
// NOTE: On Wayland it is a fullscreen window, there are no video modes to change
void ToggleBorderlessWindowed(void)
{
    // Leave fullscreen before attempting to set borderless windowed mode
    if (FLAG_IS_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE)) ToggleFullscreen();

    if (!FLAG_IS_SET(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE))
    {
        // Store screen position and size
        CORE.Window.previousPosition = CORE.Window.position;
        CORE.Window.previousScreen = CORE.Window.screen;

        int monitor = GetCurrentMonitor();
        platform.fullscreenMonitor = monitor;
        FLAG_SET(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE);
        FLAG_SET(CORE.Window.flags, FLAG_WINDOW_UNDECORATED);

        RequestFullscreen(true, GetMonitorOutput(monitor));
    }
    else
    {
        // Restore previous screen values
        CORE.Window.position = CORE.Window.previousPosition;
        FLAG_CLEAR(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE);
        if (platform.decorationRequested) FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_UNDECORATED);

        RequestFullscreen(false, NULL);
    }

    // NOTE: New window size is applied on next configure event
    wl_display_flush(platform.display);
}

// Set window state: maximized, if resizable
void MaximizeWindow(void)
{
    if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_RESIZABLE))
    {
        if (platform.libdecor.frame != NULL) platform.libdecor.FrameSetMaximized(platform.libdecor.frame);
        else if (platform.xdgToplevel != NULL) xdg_toplevel_set_maximized(platform.xdgToplevel);

        FLAG_SET(CORE.Window.flags, FLAG_WINDOW_MAXIMIZED);
        wl_display_flush(platform.display);
    }
}

// Set window state: minimized
// NOTE: Minimized state is not reported back by Wayland compositors
void MinimizeWindow(void)
{
    if (platform.libdecor.frame != NULL) platform.libdecor.FrameSetMinimized(platform.libdecor.frame);
    else if (platform.xdgToplevel != NULL) xdg_toplevel_set_minimized(platform.xdgToplevel);

    wl_display_flush(platform.display);
}

// Restore window from being minimized/maximized
// NOTE: Wayland clients can not restore a minimized window, only the compositor can
void RestoreWindow(void)
{
    if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_RESIZABLE))
    {
        if (platform.libdecor.frame != NULL) platform.libdecor.FrameUnsetMaximized(platform.libdecor.frame);
        else if (platform.xdgToplevel != NULL) xdg_toplevel_unset_maximized(platform.xdgToplevel);

        FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_MAXIMIZED);
        wl_display_flush(platform.display);
    }
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
        UpdateWindowGeometry();
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
        FLAG_SET(CORE.Window.flags, FLAG_WINDOW_UNFOCUSED);
        DestroyShellObjects();
    }

    // State change: FLAG_WINDOW_MINIMIZED
    if ((FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_MINIMIZED) != FLAG_IS_SET(flags, FLAG_WINDOW_MINIMIZED)) && FLAG_IS_SET(flags, FLAG_WINDOW_MINIMIZED))
    {
        MinimizeWindow();       // NOTE: Window state flag updated on configure event
    }

    // State change: FLAG_WINDOW_MAXIMIZED
    if ((FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_MAXIMIZED) != FLAG_IS_SET(flags, FLAG_WINDOW_MAXIMIZED)) && FLAG_IS_SET(flags, FLAG_WINDOW_MAXIMIZED))
    {
        MaximizeWindow();       // NOTE: Window state flag updated inside function
    }

    // State change: FLAG_WINDOW_UNFOCUSED
    if ((FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_UNFOCUSED) != FLAG_IS_SET(flags, FLAG_WINDOW_UNFOCUSED)) && FLAG_IS_SET(flags, FLAG_WINDOW_UNFOCUSED))
    {
        TRACELOG(LOG_WARNING, "SetWindowState(FLAG_WINDOW_UNFOCUSED) not available on target platform");
    }

    // State change: FLAG_WINDOW_TOPMOST
    if ((FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_TOPMOST) != FLAG_IS_SET(flags, FLAG_WINDOW_TOPMOST)) && FLAG_IS_SET(flags, FLAG_WINDOW_TOPMOST))
    {
        TRACELOG(LOG_WARNING, "SetWindowState(FLAG_WINDOW_TOPMOST) not available on target platform");
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
        struct wl_region *region = wl_compositor_create_region(platform.compositor);
        wl_surface_set_input_region(platform.surface, region);
        wl_region_destroy(region);
        CommitWindowSurface();
        FLAG_SET(CORE.Window.flags, FLAG_WINDOW_MOUSE_PASSTHROUGH);
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
}

// Clear window configuration state flags
void ClearWindowState(unsigned int flags)
{
    // Check previous state and requested state to apply required changes
    // NOTE: In most cases the functions already change the flags internally

    // State change: FLAG_VSYNC_HINT
    if ((FLAG_IS_SET(CORE.Window.flags, FLAG_VSYNC_HINT)) && (FLAG_IS_SET(flags, FLAG_VSYNC_HINT)))
    {
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
        UpdateWindowGeometry();
    }

    // State change: FLAG_WINDOW_HIDDEN
    if ((FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_HIDDEN)) && (FLAG_IS_SET(flags, FLAG_WINDOW_HIDDEN)))
    {
        FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_HIDDEN);
        CreateShellObjects();
    }

    // State change: FLAG_WINDOW_MINIMIZED
    if ((FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_MINIMIZED)) && (FLAG_IS_SET(flags, FLAG_WINDOW_MINIMIZED)))
    {
        TRACELOG(LOG_WARNING, "ClearWindowState(FLAG_WINDOW_MINIMIZED) not available on target platform");
    }

    // State change: FLAG_WINDOW_MAXIMIZED
    if ((FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_MAXIMIZED)) && (FLAG_IS_SET(flags, FLAG_WINDOW_MAXIMIZED)))
    {
        RestoreWindow();       // NOTE: Window state flag updated inside function
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
        wl_surface_set_input_region(platform.surface, NULL);
        CommitWindowSurface();
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
}

// Set icon for window
// NOTE 1: Image must be in RGBA format, 8bit per channel
// NOTE 2: Image is scaled by the compositor for all required sizes
void SetWindowIcon(Image image)
{
    if (image.data == NULL) SetWindowIcons(NULL, 0);
    else SetWindowIcons(&image, 1);
}

// Set icon for window, multiple images
// NOTE 1: Images must be in RGBA format, 8bit per channel
// NOTE 2: Images must be square, the multiple images are used depending on provided sizes
void SetWindowIcons(Image *images, int count)
{
    // NOTE: Icon is also drawn by the builtin titlebar, the compositor may not support window icons
    if ((platform.iconManager == NULL) && (platform.decorationMode != DECORATION_MODE_BUILTIN))
    {
        TRACELOG(LOG_WARNING, "WAYLAND: Window icon not supported by the compositor");
    }

    // Release previous icon buffers
    for (int i = 0; i < platform.iconBufferCount; i++) DestroyShmBuffer(&platform.iconBuffers[i]);
    platform.iconBufferCount = 0;

    for (int i = 0; (images != NULL) && (i < count) && (platform.iconBufferCount < MAX_WINDOW_ICONS); i++)
    {
        if (images[i].format != PIXELFORMAT_UNCOMPRESSED_R8G8B8A8)
        {
            TRACELOG(LOG_WARNING, "WAYLAND: Window icon image must be in R8G8B8A8 pixel format");
            continue;
        }

        if (images[i].width != images[i].height)
        {
            TRACELOG(LOG_WARNING, "WAYLAND: Window icon image must be square");
            continue;
        }

        ShmBuffer *buffer = &platform.iconBuffers[platform.iconBufferCount];
        if (!CreateShmBuffer(buffer, images[i].width, images[i].height)) continue;

        // Convert RGBA to premultiplied ARGB8888
        const unsigned char *pixels = (const unsigned char *)images[i].data;
        for (int p = 0; p < images[i].width*images[i].height; p++)
        {
            unsigned int a = pixels[p*4 + 3];
            unsigned int r = pixels[p*4 + 0]*a/255;
            unsigned int g = pixels[p*4 + 1]*a/255;
            unsigned int b = pixels[p*4 + 2]*a/255;
            buffer->pixels[p] = (a << 24) | (r << 16) | (g << 8) | b;
        }

        platform.iconBufferCount++;
    }

    // NOTE: Icon is applied when the window is shown if it is hidden
    UpdateWindowIcon();
    platform.decorationDirty = true;
    wl_display_flush(platform.display);
}

// Set title for window
void SetWindowTitle(const char *title)
{
    CORE.Window.title = title;
    snprintf(platform.title, sizeof(platform.title), "%s", (title != NULL)? title : "");

    if (platform.libdecor.frame != NULL) platform.libdecor.FrameSetTitle(platform.libdecor.frame, platform.title);
    else if (platform.xdgToplevel != NULL) xdg_toplevel_set_title(platform.xdgToplevel, platform.title);

    platform.decorationDirty = true;
    UpdateDecorations();
}

// Set window position on screen (windowed mode)
// NOTE: Not supported, window position is managed by the compositor
void SetWindowPosition(int x, int y)
{
    TRACELOG(LOG_WARNING, "SetWindowPosition() not available on target platform");
}

// Set monitor for the current window
void SetWindowMonitor(int monitor)
{
    if ((monitor >= 0) && (monitor < platform.monitorCount))
    {
        if (FLAG_IS_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE) || FLAG_IS_SET(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE))
        {
            TRACELOG(LOG_INFO, "WAYLAND: Selected fullscreen monitor: [%i] %s", monitor, platform.monitors[monitor].name);

            platform.fullscreenMonitor = monitor;

            RequestFullscreen(true, platform.monitors[monitor].output);

            wl_display_flush(platform.display);
        }
        else TRACELOG(LOG_WARNING, "WAYLAND: Window monitor can only be selected in fullscreen mode");
    }
    else TRACELOG(LOG_WARNING, "WAYLAND: Failed to find selected monitor");
}

// Set window minimum dimensions (FLAG_WINDOW_RESIZABLE)
void SetWindowMinSize(int width, int height)
{
    CORE.Window.screenMin.width = width;
    CORE.Window.screenMin.height = height;

    UpdateWindowGeometry();
    SetWindowSize(platform.windowedWidth, platform.windowedHeight);
}

// Set window maximum dimensions (FLAG_WINDOW_RESIZABLE)
void SetWindowMaxSize(int width, int height)
{
    CORE.Window.screenMax.width = width;
    CORE.Window.screenMax.height = height;

    UpdateWindowGeometry();
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
        platform.width = width;
        platform.height = height;

        UpdateWindowBuffers();
        UpdateWindowGeometry();
        UpdateDecorations();

        if (platform.libdecor.frame != NULL)
        {
            struct libdecor_state *frameState = platform.libdecor.StateNew(platform.width, platform.height);
            platform.libdecor.FrameCommit(platform.libdecor.frame, frameState, NULL);
            platform.libdecor.StateFree(frameState);
        }

        CommitWindowSurface();
    }
}

// Set window opacity, value opacity is between 0.0 and 1.0
void SetWindowOpacity(float opacity)
{
    if (opacity >= 1.0f) opacity = 1.0f;
    else if (opacity <= 0.0f) opacity = 0.0f;

    if (platform.alphaModifier == NULL)
    {
        TRACELOG(LOG_WARNING, "WAYLAND: Window opacity not supported by the compositor");
        return;
    }

    if (platform.alphaModifierSurface == NULL) platform.alphaModifierSurface = wp_alpha_modifier_v1_get_surface(platform.alphaModifier, platform.surface);

    wp_alpha_modifier_surface_v1_set_multiplier(platform.alphaModifierSurface, (uint32_t)((double)opacity*UINT32_MAX));
    CommitWindowSurface();
}

// Set window focused
// NOTE: Requires xdg-activation, the compositor decides if focus is granted
void SetWindowFocused(void)
{
    if (platform.activation == NULL)
    {
        TRACELOG(LOG_WARNING, "WAYLAND: Window activation not supported by the compositor");
        return;
    }

    struct xdg_activation_token_v1 *token = xdg_activation_v1_get_activation_token(platform.activation);
    xdg_activation_token_v1_add_listener(token, &activationTokenListener, NULL);
    if (platform.seat != NULL) xdg_activation_token_v1_set_serial(token, platform.inputSerial, platform.seat);
    xdg_activation_token_v1_set_surface(token, platform.surface);
    xdg_activation_token_v1_commit(token);
    wl_display_flush(platform.display);
}

// Get native window handle
// NOTE: Returns the window struct wl_surface pointer
void *GetWindowHandle(void)
{
    return (void *)platform.surface;
}

// Get number of monitors
int GetMonitorCount(void)
{
    return platform.monitorCount;
}

// Get current monitor where window is placed
// NOTE: Monitor the window entered last, primary (first) monitor if unknown
int GetCurrentMonitor(void)
{
    int index = 0;

    if (FLAG_IS_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE) && (platform.fullscreenMonitor < platform.monitorCount)) index = platform.fullscreenMonitor;
    else if (platform.enteredOutputCount > 0)
    {
        int monitor = GetMonitorIndex(platform.enteredOutputs[platform.enteredOutputCount - 1]);
        if (monitor >= 0) index = monitor;
    }

    return index;
}

// Get selected monitor position
Vector2 GetMonitorPosition(int monitor)
{
    if ((monitor >= 0) && (monitor < platform.monitorCount))
    {
        return (Vector2){ (float)platform.monitors[monitor].x, (float)platform.monitors[monitor].y };
    }
    else TRACELOG(LOG_WARNING, "WAYLAND: Failed to find selected monitor");

    return (Vector2){ 0.0f, 0.0f };
}

// Get selected monitor width (currently used by monitor)
int GetMonitorWidth(int monitor)
{
    int width = 0;

    if ((monitor >= 0) && (monitor < platform.monitorCount)) width = platform.monitors[monitor].width;
    else TRACELOG(LOG_WARNING, "WAYLAND: Failed to find selected monitor");

    return width;
}

// Get selected monitor height (currently used by monitor)
int GetMonitorHeight(int monitor)
{
    int height = 0;

    if ((monitor >= 0) && (monitor < platform.monitorCount)) height = platform.monitors[monitor].height;
    else TRACELOG(LOG_WARNING, "WAYLAND: Failed to find selected monitor");

    return height;
}

// Get selected monitor physical width in millimetres
int GetMonitorPhysicalWidth(int monitor)
{
    int width = 0;

    if ((monitor >= 0) && (monitor < platform.monitorCount)) width = platform.monitors[monitor].physicalWidth;
    else TRACELOG(LOG_WARNING, "WAYLAND: Failed to find selected monitor");

    return width;
}

// Get selected monitor physical height in millimetres
int GetMonitorPhysicalHeight(int monitor)
{
    int height = 0;

    if ((monitor >= 0) && (monitor < platform.monitorCount)) height = platform.monitors[monitor].physicalHeight;
    else TRACELOG(LOG_WARNING, "WAYLAND: Failed to find selected monitor");

    return height;
}

// Get selected monitor refresh rate
int GetMonitorRefreshRate(int monitor)
{
    int refresh = 0;

    if ((monitor >= 0) && (monitor < platform.monitorCount)) refresh = (int)roundf((float)platform.monitors[monitor].refreshRate/1000.0f);
    else TRACELOG(LOG_WARNING, "WAYLAND: Failed to find selected monitor");

    return refresh;
}

// Get the human-readable, UTF-8 encoded name of the selected monitor
const char *GetMonitorName(int monitor)
{
    if ((monitor >= 0) && (monitor < platform.monitorCount))
    {
        return platform.monitors[monitor].name;
    }
    else TRACELOG(LOG_WARNING, "WAYLAND: Failed to find selected monitor");

    return "";
}

// Get window position XY on monitor
// NOTE: Not supported, window position is managed by the compositor
Vector2 GetWindowPosition(void)
{
    return (Vector2){ (float)CORE.Window.position.x, (float)CORE.Window.position.y };
}

// Get window scale DPI factor for current monitor
Vector2 GetWindowScaleDPI(void)
{
    Vector2 scale = { 1.0f, 1.0f };

    if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_HIGHDPI)) scale = (Vector2){ platform.scale, platform.scale };

    return scale;
}

// Set clipboard text content
void SetClipboardText(const char *text)
{
    if ((platform.dataDeviceManager == NULL) || (platform.dataDevice == NULL)) return;

    if (platform.dataSource != NULL)
    {
        wl_data_source_destroy(platform.dataSource);
        platform.dataSource = NULL;
    }

    RL_FREE(platform.clipboardText);
    platform.clipboardText = NULL;

    if (text == NULL) return;

    int length = (int)strlen(text);
    platform.clipboardText = (char *)RL_CALLOC(length + 1, sizeof(char));
    memcpy(platform.clipboardText, text, length);

    platform.dataSource = wl_data_device_manager_create_data_source(platform.dataDeviceManager);
    wl_data_source_add_listener(platform.dataSource, &dataSourceListener, NULL);
    for (int i = 0; i < (int)(sizeof(textMimeTypes)/sizeof(textMimeTypes[0])); i++) wl_data_source_offer(platform.dataSource, textMimeTypes[i]);
    wl_data_device_set_selection(platform.dataDevice, platform.dataSource, platform.inputSerial);
    wl_display_flush(platform.display);
}

// Get clipboard text content
// NOTE: returned string is allocated and freed by the platform
const char *GetClipboardText(void)
{
    if (platform.dataSource != NULL) return platform.clipboardText;

    RL_FREE(platform.clipboardReceived);
    platform.clipboardReceived = NULL;

    DataOffer *offer = platform.selectionOffer;
    if (offer == NULL) return NULL;

    for (int i = 0; i < (int)(sizeof(textMimeTypes)/sizeof(textMimeTypes[0])); i++)
    {
        if (DataOfferHasMimeType(offer, textMimeTypes[i]))
        {
            platform.clipboardReceived = ReadDataOffer(offer, textMimeTypes[i], NULL);
            break;
        }
    }

    return platform.clipboardReceived;
}

// Get clipboard image
Image GetClipboardImage(void)
{
    Image image = { 0 };

#if SUPPORT_CLIPBOARD_IMAGE && SUPPORT_MODULE_RTEXTURES
    DataOffer *offer = platform.selectionOffer;

    if ((platform.dataSource == NULL) && (offer != NULL))
    {
        const char *imageMimeTypes[] = { "image/png", "image/bmp", "image/jpeg" };
        const char *imageFileTypes[] = { ".png", ".bmp", ".jpg" };

        for (int i = 0; i < (int)(sizeof(imageMimeTypes)/sizeof(imageMimeTypes[0])); i++)
        {
            if (DataOfferHasMimeType(offer, imageMimeTypes[i]))
            {
                int size = 0;
                char *data = ReadDataOffer(offer, imageMimeTypes[i], &size);

                if (data != NULL)
                {
                    image = LoadImageFromMemory(imageFileTypes[i], (const unsigned char *)data, size);
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

    // Set cursor position in the middle
    SetMousePosition(CORE.Window.screen.width/2, CORE.Window.screen.height/2);

    CORE.Input.Mouse.cursorHidden = false;
    CORE.Input.Mouse.cursorLocked = false;
    UpdateCursor();
}

// Disable cursor (lock cursor)
void DisableCursor(void)
{
    // Reset mouse position within the window area before disabling cursor
    SetMousePosition(CORE.Window.screen.width/2, CORE.Window.screen.height/2);

    LockPointer(true);

    CORE.Input.Mouse.cursorHidden = true;
    CORE.Input.Mouse.cursorLocked = true;
    UpdateCursor();
}

// Swap back buffer with front buffer (screen drawing)
void SwapScreenBuffer(void)
{
    // NOTE: Hidden windows have no surface role, buffers can not be attached; frames are paced like occluded windows
    if (!platform.visible || platform.connectionLost)
    {
        if (FLAG_IS_SET(CORE.Window.flags, FLAG_VSYNC_HINT)) WaitTime((double)FRAME_CALLBACK_TIMEOUT/1000.0);
        return;
    }

    // Wait for the compositor to request a new frame, frames are limited to screen refresh rate
    // NOTE: Waiting is limited in time, compositors do not send frame callbacks for occluded windows
    if (FLAG_IS_SET(CORE.Window.flags, FLAG_VSYNC_HINT) && (platform.frameCallback != NULL))
    {
        WaitFrameCallback(FRAME_CALLBACK_TIMEOUT);
    }

    if (platform.frameCallback != NULL)
    {
        wl_callback_destroy(platform.frameCallback);
        platform.frameCallback = NULL;
    }

    // Release frame callback events already read into the queue
    if (wl_display_dispatch_queue_pending(platform.display, platform.frameQueue) < 0) HandleConnectionError();

    // Request next frame callback on its own queue, input events are dispatched only in PollInputEvents()
    if (FLAG_IS_SET(CORE.Window.flags, FLAG_VSYNC_HINT))
    {
        struct wl_surface *surfaceWrapper = (struct wl_surface *)wl_proxy_create_wrapper(platform.surface);
        wl_proxy_set_queue((struct wl_proxy *)surfaceWrapper, platform.frameQueue);
        platform.frameCallback = wl_surface_frame(surfaceWrapper);
        wl_proxy_wrapper_destroy(surfaceWrapper);
        wl_callback_add_listener(platform.frameCallback, &frameCallbackListener, NULL);
    }

    if ((platform.viewport == NULL) && (platform.bufferScale > 0) && (platform.bufferScale != platform.appliedBufferScale))
    {
        platform.appliedBufferScale = platform.bufferScale;
        wl_surface_set_buffer_scale(platform.surface, platform.bufferScale);
    }

#if defined(GRAPHICS_API_OPENGL_SOFTWARE)
    // Copy software renderer framebuffer to a free shared memory buffer
    ShmBuffer *buffer = NULL;
    for (int i = 0; i < 2; i++)
    {
        if (!platform.softwareBuffers[i].busy)
        {
            buffer = &platform.softwareBuffers[i];
            break;
        }
    }

    if ((buffer != NULL) && (platform.softwarePixels != NULL))
    {
        int width = platform.softwareWidth;
        int height = platform.softwareHeight;

        if ((buffer->width != width) || (buffer->height != height))
        {
            DestroyShmBuffer(buffer);
            CreateShmBuffer(buffer, width, height);
        }

        if (buffer->buffer != NULL)
        {
            rlCopyFramebuffer(0, 0, width, height, PIXELFORMAT_UNCOMPRESSED_R8G8B8A8, platform.softwarePixels);

            // NOTE: Software renderer framebuffer copy is in BGRA order, premultiply alpha for transparent windows
            bool transparent = FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_TRANSPARENT);
            for (int i = 0; i < width*height; i++)
            {
                const unsigned char *p = platform.softwarePixels + i*4;
                unsigned int a = transparent? p[3] : 255;
                buffer->pixels[i] = (a << 24) | ((p[2]*a/255) << 16) | ((p[1]*a/255) << 8) | (p[0]*a/255);
            }

            wl_surface_attach(platform.surface, buffer->buffer, 0, 0);
            wl_surface_damage_buffer(platform.surface, 0, 0, width, height);
            buffer->busy = true;
        }
    }

    wl_surface_commit(platform.surface);
    wl_display_flush(platform.display);
#else
    eglSwapBuffers(platform.eglDisplay, platform.eglSurface);
#endif

    platform.frameCommitted = true;
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
    return mg_update_gamepad_mappings(&platform.minigamepad, mappings);
}

// Set gamepad vibration
void SetGamepadVibration(int gamepad, float leftMotor, float rightMotor, float duration)
{
    TRACELOG(LOG_WARNING, "SetGamepadVibration() not available on target platform");
}

// Set mouse position XY
// NOTE: Wayland clients can not move the pointer, the position is only a hint to the compositor while the cursor is disabled
void SetMousePosition(int x, int y)
{
    CORE.Input.Mouse.currentPosition = (Vector2){ (float)x, (float)y };

    if (platform.lockedPointer != NULL)
    {
        zwp_locked_pointer_v1_set_cursor_position_hint(platform.lockedPointer, wl_fixed_from_int(x), wl_fixed_from_int(y));
        CommitWindowSurface();
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

    if ((platform.xkbKeymap == NULL) || (platform.xkbState == NULL) || (key <= 0)) return NULL;

    int scancode = -1;
    for (int i = 0; i < KEYMAP_SIZE; i++)
    {
        if (linuxToRaylibMap[i] == key)
        {
            scancode = i;
            break;
        }
    }

    if (scancode < 0) return NULL;

    const xkb_keysym_t *keysyms = NULL;
    xkb_layout_index_t layout = xkb_state_key_get_layout(platform.xkbState, scancode + XKB_KEYCODE_OFFSET);
    if (layout == XKB_LAYOUT_INVALID) return NULL;

    int count = xkb_keymap_key_get_syms_by_level(platform.xkbKeymap, scancode + XKB_KEYCODE_OFFSET, layout, 0, &keysyms);
    if (count != 1) return NULL;

    if (xkb_keysym_to_utf32(keysyms[0]) < 32) return NULL;

    memset(keyName, 0, sizeof(keyName));
    if (xkb_keysym_to_utf8(keysyms[0], keyName, sizeof(keyName)) <= 0) return NULL;

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
        int timeout = -1;
        if (platform.repeatScancode != 0)
        {
            timeout = (int)((platform.repeatTime - GetTime())*1000.0) + 1;
            if (timeout < 0) timeout = 0;
        }

        DispatchEvents(timeout);
        CORE.Time.previous = GetTime();
    }
    else DispatchEvents(0);

    ProcessKeyRepeat();

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
    // Resize software renderer framebuffer to the window buffer size
    if ((platform.softwareWidth != (int)CORE.Window.render.width) || (platform.softwareHeight != (int)CORE.Window.render.height))
    {
        platform.softwareWidth = CORE.Window.render.width;
        platform.softwareHeight = CORE.Window.render.height;
        RL_FREE(platform.softwarePixels);
        platform.softwarePixels = (unsigned char *)RL_MALLOC(platform.softwareWidth*platform.softwareHeight*4);
        rlResizeFramebuffer(platform.softwareWidth, platform.softwareHeight);
    }
#endif

    // Draw titlebar text once the default font is available
    if ((platform.decorationMode == DECORATION_MODE_BUILTIN) && (!platform.decorationTitleReady || platform.decorationDirty)) DrawDecorations();

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
    // Decorations mode preference
#if defined(WAYLAND_CLIENT_DECORATIONS) || defined(WAYLAND_BUILTIN_DECORATIONS)
    platform.preferClientDecorations = true;
#endif
#if defined(WAYLAND_BUILTIN_DECORATIONS)
    platform.preferBuiltinDecorations = true;
#endif

    // Connect to the compositor and get required globals
    platform.display = wl_display_connect(NULL);
    if (platform.display == NULL)
    {
        TRACELOG(LOG_WARNING, "WAYLAND: Failed to connect to Wayland display");
        return -1;
    }

    platform.frameQueue = wl_display_create_queue(platform.display);
    platform.registry = wl_display_get_registry(platform.display);
    wl_registry_add_listener(platform.registry, &registryListener, NULL);
    wl_display_roundtrip(platform.display);     // Get globals
    wl_display_roundtrip(platform.display);     // Get outputs and seat information

    if ((platform.compositor == NULL) || (platform.wmBase == NULL) || (platform.shm == NULL))
    {
        TRACELOG(LOG_WARNING, "WAYLAND: Compositor does not support required protocols (wl_compositor, wl_shm, xdg_wm_base)");
        ClosePlatform();
        return -1;
    }

    // Initialize keyboard context and compose sequences (dead keys)
    platform.xkbContext = xkb_context_new(XKB_CONTEXT_NO_FLAGS);

    const char *locale = getenv("LC_ALL");
    if ((locale == NULL) || (locale[0] == '\0')) locale = getenv("LC_CTYPE");
    if ((locale == NULL) || (locale[0] == '\0')) locale = getenv("LANG");
    if ((locale == NULL) || (locale[0] == '\0')) locale = "C";

    platform.xkbComposeTable = xkb_compose_table_new_from_locale(platform.xkbContext, locale, XKB_COMPOSE_COMPILE_NO_FLAGS);
    if (platform.xkbComposeTable != NULL) platform.xkbComposeState = xkb_compose_state_new(platform.xkbComposeTable, XKB_COMPOSE_STATE_NO_FLAGS);

    platform.repeatRate = KEY_REPEAT_RATE_DEFAULT;
    platform.repeatDelay = KEY_REPEAT_DELAY_DEFAULT;
    platform.scale = 1.0f;
    platform.preferredBufferScale = 1;
    for (int i = 0; i < MAX_TOUCH_POINTS; i++) platform.touchIds[i] = -1;

    // Initialize window size
    if ((CORE.Window.screen.width == 0) || (CORE.Window.screen.height == 0)) FLAG_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE);

    if (platform.monitorCount > 0)
    {
        CORE.Window.display.width = platform.monitors[0].width;
        CORE.Window.display.height = platform.monitors[0].height;
    }

    if (FLAG_IS_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE))
    {
        // Check if user requested some screen size
        if ((CORE.Window.screen.width == 0) || (CORE.Window.screen.height == 0))
        {
            // Set some default screen size in case user decides to exit fullscreen mode
            CORE.Window.previousScreen.width = 800;
            CORE.Window.previousScreen.height = 450;

            // Set screen width/height to the display logical width/height
            if (CORE.Window.screen.width == 0) CORE.Window.screen.width = (platform.monitorCount > 0)? platform.monitors[0].logicalWidth : 800;
            if (CORE.Window.screen.height == 0) CORE.Window.screen.height = (platform.monitorCount > 0)? platform.monitors[0].logicalHeight : 450;
        }
        else CORE.Window.previousScreen = CORE.Window.screen;
    }

    platform.width = CORE.Window.screen.width;
    platform.height = CORE.Window.screen.height;
    platform.windowedWidth = FLAG_IS_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE)? (int)CORE.Window.previousScreen.width : platform.width;
    platform.windowedHeight = FLAG_IS_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE)? (int)CORE.Window.previousScreen.height : platform.height;

    snprintf(platform.title, sizeof(platform.title), "%s", (CORE.Window.title != NULL)? CORE.Window.title : "");

    // Window application id, used by compositors to group windows and to find the desktop entry
    snprintf(platform.appId, sizeof(platform.appId), "raylib");
    char exePath[512] = { 0 };
    ssize_t exePathLength = readlink("/proc/self/exe", exePath, sizeof(exePath) - 1);
    if (exePathLength > 0) snprintf(platform.appId, sizeof(platform.appId), "%s", GetFileName(exePath));

    // Create window surface
    platform.surface = wl_compositor_create_surface(platform.compositor);
    wl_surface_add_listener(platform.surface, &surfaceListener, NULL);

    if ((platform.fractionalScaleManager != NULL) && (platform.viewporter != NULL))
    {
        platform.fractionalScale = wp_fractional_scale_manager_v1_get_fractional_scale(platform.fractionalScaleManager, platform.surface);
        wp_fractional_scale_v1_add_listener(platform.fractionalScale, &fractionalScaleListener, NULL);
    }

    if (platform.viewporter != NULL) platform.viewport = wp_viewporter_get_viewport(platform.viewporter, platform.surface);

    if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_MOUSE_PASSTHROUGH))
    {
        struct wl_region *region = wl_compositor_create_region(platform.compositor);
        wl_surface_set_input_region(platform.surface, region);
        wl_region_destroy(region);
    }

    // Initialize graphics context
    if (!InitGraphicsContext())
    {
        TRACELOG(LOG_FATAL, "PLATFORM: Failed to initialize graphics device");
        ClosePlatform();
        return -1;
    }

    // Create window role objects and wait for the initial configure
    platform.decorationRequested = !FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_UNDECORATED);
    if (!FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_HIDDEN)) CreateShellObjects();
    else UpdateWindowBuffers();

    CORE.Window.ready = true;

    // Apply FLAG_WINDOW_MINIMIZED once the window exists
    if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_MINIMIZED))
    {
        FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_MINIMIZED);
        MinimizeWindow();
    }

    if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_TOPMOST)) TRACELOG(LOG_WARNING, "SetWindowState(FLAG_WINDOW_TOPMOST) not available on target platform");

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

    const char *decorationNames[] = { "none", "server-side", "libdecor", "builtin" };
    TRACELOG(LOG_INFO, "WAYLAND: Window decorations: %s", decorationNames[platform.decorationMode]);
#if defined(GRAPHICS_API_OPENGL_SOFTWARE)
    TRACELOG(LOG_INFO, "PLATFORM: DESKTOP (WAYLAND - Software): Initialized successfully");
#else
    TRACELOG(LOG_INFO, "PLATFORM: DESKTOP (WAYLAND): Initialized successfully");
#endif

    return 0;
}

// Close platform
void ClosePlatform(void)
{
    mg_gamepads_free(&platform.minigamepad);

    LockPointer(false);
    DestroyShellObjects();

#if defined(GRAPHICS_API_OPENGL_SOFTWARE)
    for (int i = 0; i < 2; i++) DestroyShmBuffer(&platform.softwareBuffers[i]);
    RL_FREE(platform.softwarePixels);
    platform.softwarePixels = NULL;
#else
    if (platform.eglDisplay != EGL_NO_DISPLAY)
    {
        eglMakeCurrent(platform.eglDisplay, EGL_NO_SURFACE, EGL_NO_SURFACE, EGL_NO_CONTEXT);
        if (platform.eglSurface != EGL_NO_SURFACE) eglDestroySurface(platform.eglDisplay, platform.eglSurface);
        if (platform.eglContext != EGL_NO_CONTEXT) eglDestroyContext(platform.eglDisplay, platform.eglContext);
        eglTerminate(platform.eglDisplay);
    }
    if (platform.eglWindow != NULL) wl_egl_window_destroy(platform.eglWindow);
#endif

    for (int i = 0; i < platform.iconBufferCount; i++) DestroyShmBuffer(&platform.iconBuffers[i]);
    if (platform.frameCallback != NULL) wl_callback_destroy(platform.frameCallback);
    if (platform.frameQueue != NULL) wl_event_queue_destroy(platform.frameQueue);
    if (platform.alphaModifierSurface != NULL) wp_alpha_modifier_surface_v1_destroy(platform.alphaModifierSurface);
    if (platform.fractionalScale != NULL) wp_fractional_scale_v1_destroy(platform.fractionalScale);
    if (platform.viewport != NULL) wp_viewport_destroy(platform.viewport);
    if (platform.surface != NULL) wl_surface_destroy(platform.surface);

    if (platform.dataSource != NULL) wl_data_source_destroy(platform.dataSource);
    DestroyDataOffer(platform.selectionOffer);
    DestroyDataOffer(platform.dragOffer);
    if (platform.dataDevice != NULL)
    {
        if (wl_data_device_get_version(platform.dataDevice) >= WL_DATA_DEVICE_RELEASE_SINCE_VERSION) wl_data_device_release(platform.dataDevice);
        else wl_data_device_destroy(platform.dataDevice);
    }
    RL_FREE(platform.clipboardText);
    RL_FREE(platform.clipboardReceived);

    if (platform.cursorShapeDevice != NULL) wp_cursor_shape_device_v1_destroy(platform.cursorShapeDevice);
    if (platform.cursorSurface != NULL) wl_surface_destroy(platform.cursorSurface);
    if (platform.cursorTheme != NULL) wl_cursor_theme_destroy(platform.cursorTheme);
    if (platform.pointer != NULL) ReleasePointer(platform.pointer);
    if (platform.keyboard != NULL) ReleaseKeyboard(platform.keyboard);
    if (platform.touch != NULL) ReleaseTouch(platform.touch);

    if (platform.xkbComposeState != NULL) xkb_compose_state_unref(platform.xkbComposeState);
    if (platform.xkbComposeTable != NULL) xkb_compose_table_unref(platform.xkbComposeTable);
    if (platform.xkbState != NULL) xkb_state_unref(platform.xkbState);
    if (platform.xkbKeymap != NULL) xkb_keymap_unref(platform.xkbKeymap);
    if (platform.xkbContext != NULL) xkb_context_unref(platform.xkbContext);

    for (int i = 0; i < platform.monitorCount; i++)
    {
        if (platform.monitors[i].xdgOutput != NULL) zxdg_output_v1_destroy(platform.monitors[i].xdgOutput);
        ReleaseOutput(platform.monitors[i].output);
    }

    UnloadLibdecor();

    // Destroy Wayland globals
    if (platform.seat != NULL)
    {
        if (platform.seatVersion >= WL_SEAT_RELEASE_SINCE_VERSION) wl_seat_release(platform.seat);
        else wl_seat_destroy(platform.seat);
    }
    if (platform.outputManager != NULL) zxdg_output_manager_v1_destroy(platform.outputManager);
    if (platform.alphaModifier != NULL) wp_alpha_modifier_v1_destroy(platform.alphaModifier);
    if (platform.iconManager != NULL) xdg_toplevel_icon_manager_v1_destroy(platform.iconManager);
    if (platform.activation != NULL) xdg_activation_v1_destroy(platform.activation);
    if (platform.relativePointerManager != NULL) zwp_relative_pointer_manager_v1_destroy(platform.relativePointerManager);
    if (platform.pointerConstraints != NULL) zwp_pointer_constraints_v1_destroy(platform.pointerConstraints);
    if (platform.cursorShapeManager != NULL) wp_cursor_shape_manager_v1_destroy(platform.cursorShapeManager);
    if (platform.viewporter != NULL) wp_viewporter_destroy(platform.viewporter);
    if (platform.fractionalScaleManager != NULL) wp_fractional_scale_manager_v1_destroy(platform.fractionalScaleManager);
    if (platform.decorationManager != NULL) zxdg_decoration_manager_v1_destroy(platform.decorationManager);
    if (platform.wmBase != NULL) xdg_wm_base_destroy(platform.wmBase);
    if (platform.dataDeviceManager != NULL) wl_data_device_manager_destroy(platform.dataDeviceManager);
    if (platform.shm != NULL) wl_shm_destroy(platform.shm);
    if (platform.subcompositor != NULL) wl_subcompositor_destroy(platform.subcompositor);
    if (platform.compositor != NULL) wl_compositor_destroy(platform.compositor);
    if (platform.registry != NULL) wl_registry_destroy(platform.registry);

    if (platform.display != NULL)
    {
        wl_display_flush(platform.display);
        wl_display_disconnect(platform.display);
    }

    memset(&platform, 0, sizeof(platform));
}

// Initialize EGL display, context and window surface
static bool InitGraphicsContext(void)
{
    platform.bufferWidth = platform.width;
    platform.bufferHeight = platform.height;

#if defined(GRAPHICS_API_OPENGL_SOFTWARE)
    platform.softwareWidth = platform.bufferWidth;
    platform.softwareHeight = platform.bufferHeight;
    platform.softwarePixels = (unsigned char *)RL_MALLOC(platform.softwareWidth*platform.softwareHeight*4);

    return true;
#else
    platform.eglDisplay = eglGetDisplay((EGLNativeDisplayType)platform.display);
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
    eglBindAPI(openglES? EGL_OPENGL_ES_API : EGL_OPENGL_API);

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

    EGLint alphaSize = FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_TRANSPARENT)? 8 : 0;

    const EGLint framebufferAttribs[] = {
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
        TRACELOG(LOG_WARNING, "DISPLAY: Failed to choose EGL config");
        return false;
    }

    EGLConfig *configs = (EGLConfig *)RL_CALLOC(numConfigs, sizeof(EGLConfig));
    if ((eglChooseConfig(platform.eglDisplay, framebufferAttribs, configs, numConfigs, &numConfigs) == EGL_FALSE) || (numConfigs == 0))
    {
        TRACELOG(LOG_WARNING, "DISPLAY: Failed to choose EGL config");
        RL_FREE(configs);
        return false;
    }

    // NOTE: Configs are sorted by color depth, the first one can be 10 bit per channel (more bandwidth, compositor conversion)
    platform.eglConfig = configs[0];

    for (EGLint i = 0; i < numConfigs; i++)
    {
        EGLint red = 0;
        EGLint green = 0;
        EGLint blue = 0;
        EGLint alpha = 0;

        eglGetConfigAttrib(platform.eglDisplay, configs[i], EGL_RED_SIZE, &red);
        eglGetConfigAttrib(platform.eglDisplay, configs[i], EGL_GREEN_SIZE, &green);
        eglGetConfigAttrib(platform.eglDisplay, configs[i], EGL_BLUE_SIZE, &blue);
        eglGetConfigAttrib(platform.eglDisplay, configs[i], EGL_ALPHA_SIZE, &alpha);

        if ((red == 8) && (green == 8) && (blue == 8) && (alpha == alphaSize))
        {
            platform.eglConfig = configs[i];
            break;
        }
    }

    RL_FREE(configs);

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

    platform.eglWindow = wl_egl_window_create(platform.surface, platform.bufferWidth, platform.bufferHeight);
    platform.eglSurface = eglCreateWindowSurface(platform.eglDisplay, platform.eglConfig, (EGLNativeWindowType)platform.eglWindow, NULL);
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

    // NOTE: Swap interval is always 0, V-Sync waits for the compositor frame callback in SwapScreenBuffer()
    // to avoid blocking indefinitely when the window is not visible
    eglSwapInterval(platform.eglDisplay, 0);
    if (FLAG_IS_SET(CORE.Window.flags, FLAG_VSYNC_HINT)) TRACELOG(LOG_INFO, "DISPLAY: Trying to enable VSYNC");

    return true;
#endif
}

// Read and dispatch Wayland events, timeout in milliseconds (-1 blocks)
static void DispatchEvents(int timeout)
{
    if ((platform.display == NULL) || platform.connectionLost) return;

    platform.eventReceived = false;

    // NOTE: libdecor dispatches the default queue, events received there end the waiting too
    if (platform.libdecor.context != NULL) platform.libdecor.Dispatch(platform.libdecor.context, 0);

    double endTime = GetTime() + (double)timeout/1000.0;

    // NOTE: Waiting ends when a window or input event is received, other events (e.g. EGL buffers synchronization) are not counted
    while (!platform.connectionLost)
    {
        while (wl_display_prepare_read(platform.display) != 0)
        {
            if (wl_display_dispatch_pending(platform.display) < 0)
            {
                HandleConnectionError();
                return;
            }
        }

        bool decorationRedraw = (platform.decorationDirty && (platform.decorationParts[DECORATION_PART_TITLEBAR].surface != NULL));
        int remaining = timeout;
        if (platform.eventReceived || decorationRedraw) remaining = 0;
        else if (timeout > 0)
        {
            remaining = (int)((endTime - GetTime())*1000.0);
            if (remaining < 0) remaining = 0;
        }

        if (!ReadEvents(remaining)) return;

        if (wl_display_dispatch_pending(platform.display) < 0)
        {
            HandleConnectionError();
            return;
        }

        decorationRedraw = (platform.decorationDirty && (platform.decorationParts[DECORATION_PART_TITLEBAR].surface != NULL));
        if (platform.eventReceived || decorationRedraw || (remaining == 0)) break;
    }
}

// Wait for the compositor frame callback, timeout in milliseconds
// NOTE: Only the frame callback queue is dispatched, other events are read and kept queued
static void WaitFrameCallback(int timeout)
{
    double endTime = GetTime() + (double)timeout/1000.0;

    while ((platform.frameCallback != NULL) && !platform.connectionLost)
    {
        if (wl_display_dispatch_queue_pending(platform.display, platform.frameQueue) < 0)
        {
            HandleConnectionError();
            break;
        }

        if (platform.frameCallback == NULL) break;

        int remaining = (int)((endTime - GetTime())*1000.0);
        if (remaining <= 0) break;

        while (wl_display_prepare_read_queue(platform.display, platform.frameQueue) != 0)
        {
            if (wl_display_dispatch_queue_pending(platform.display, platform.frameQueue) < 0)
            {
                HandleConnectionError();
                return;
            }
        }

        if (!ReadEvents(remaining)) break;
    }
}

// Read Wayland events after a prepared read, waiting up to timeout milliseconds (-1 blocks)
// NOTE: Events are queued, not dispatched; returns false if the connection failed
static bool ReadEvents(int timeout)
{
    if ((wl_display_flush(platform.display) < 0) && (errno != EAGAIN))
    {
        wl_display_cancel_read(platform.display);
        HandleConnectionError();
        return false;
    }

    struct pollfd fds[1] = { { wl_display_get_fd(platform.display), POLLIN, 0 } };
    int result = 0;

    do
    {
        result = poll(fds, 1, timeout);
    } while ((result < 0) && (errno == EINTR));

    if ((result > 0) && FLAG_IS_SET(fds[0].revents, POLLIN)) wl_display_read_events(platform.display);
    else wl_display_cancel_read(platform.display);

    return true;
}

// Register connection to the compositor lost, the window should close
static void HandleConnectionError(void)
{
    if (platform.connectionLost) return;

    platform.connectionLost = true;

    if (wl_display_get_error(platform.display) == EPROTO)
    {
        const struct wl_interface *interface = NULL;
        uint32_t id = 0;
        uint32_t code = wl_display_get_protocol_error(platform.display, &interface, &id);

        TRACELOG(LOG_WARNING, "WAYLAND: Protocol error %u on %s@%u", code, (interface != NULL)? interface->name : "unknown", id);
    }
    else TRACELOG(LOG_WARNING, "WAYLAND: Connection to the compositor lost");
}

// Create window role objects (xdg_surface, xdg_toplevel, decorations)
static void CreateShellObjects(void)
{
    if ((platform.xdgSurface != NULL) || (platform.libdecor.frame != NULL)) return;

    platform.fullscreenPending = false;
    platform.sizeLimitsResizable = -1;
    platform.decorationVisible = 1;

    bool fullscreen = (FLAG_IS_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE) || FLAG_IS_SET(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE));
    struct wl_output *fullscreenOutput = GetMonitorOutput(platform.fullscreenMonitor);

    // Client-side decorations drawn by libdecor, required when the compositor does not support server-side decorations
    // NOTE: libdecor is also used for undecorated windows, its decorations are hidden until requested
    bool useLibdecor = false;
    if (!platform.preferBuiltinDecorations && ((platform.decorationManager == NULL) || platform.preferClientDecorations)) useLibdecor = LoadLibdecor();

    if (useLibdecor)
    {
        platform.libdecor.frame = platform.libdecor.Decorate(platform.libdecor.context, platform.surface, &libdecorFrameInterface, NULL);
        platform.libdecor.FrameSetAppId(platform.libdecor.frame, platform.appId);
        platform.libdecor.FrameSetTitle(platform.libdecor.frame, platform.title);

        if (fullscreen) RequestFullscreen(true, fullscreenOutput);
        if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_MAXIMIZED)) platform.libdecor.FrameSetMaximized(platform.libdecor.frame);

        platform.decorationMode = DECORATION_MODE_LIBDECOR;
        UpdateWindowGeometry();
        platform.libdecor.FrameMap(platform.libdecor.frame);
    }
    else
    {
        platform.xdgSurface = xdg_wm_base_get_xdg_surface(platform.wmBase, platform.surface);
        xdg_surface_add_listener(platform.xdgSurface, &xdgSurfaceListener, NULL);

        platform.xdgToplevel = xdg_surface_get_toplevel(platform.xdgSurface);
        xdg_toplevel_add_listener(platform.xdgToplevel, &xdgToplevelListener, NULL);
        xdg_toplevel_set_app_id(platform.xdgToplevel, platform.appId);
        xdg_toplevel_set_title(platform.xdgToplevel, platform.title);

        if (fullscreen) RequestFullscreen(true, fullscreenOutput);
        if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_MAXIMIZED)) xdg_toplevel_set_maximized(platform.xdgToplevel);

        if (platform.decorationManager != NULL)
        {
            platform.toplevelDecoration = zxdg_decoration_manager_v1_get_toplevel_decoration(platform.decorationManager, platform.xdgToplevel);
            zxdg_toplevel_decoration_v1_add_listener(platform.toplevelDecoration, &toplevelDecorationListener, NULL);

            bool serverSide = (platform.decorationRequested && !platform.preferClientDecorations);
            platform.decorationModeSent = serverSide? ZXDG_TOPLEVEL_DECORATION_V1_MODE_SERVER_SIDE : ZXDG_TOPLEVEL_DECORATION_V1_MODE_CLIENT_SIDE;
            zxdg_toplevel_decoration_v1_set_mode(platform.toplevelDecoration, platform.decorationModeSent);

            platform.decorationMode = serverSide? DECORATION_MODE_SERVER : (platform.decorationRequested? DECORATION_MODE_BUILTIN : DECORATION_MODE_NONE);
        }
        else platform.decorationMode = platform.decorationRequested? DECORATION_MODE_BUILTIN : DECORATION_MODE_NONE;

        UpdateWindowGeometry();
    }

    // Wait for the initial configure
    platform.configured = false;
    wl_surface_commit(platform.surface);

    while (!platform.configured)
    {
        if (platform.libdecor.context != NULL)
        {
            if (platform.libdecor.Dispatch(platform.libdecor.context, -1) < 0) break;
        }
        else if (wl_display_dispatch(platform.display) < 0) break;
    }

    if (platform.libdecor.frame != NULL) platform.xdgToplevel = platform.libdecor.FrameGetXdgToplevel(platform.libdecor.frame);
    UpdateWindowIcon();

    platform.visible = true;
    UpdateDecorations();
}

// Destroy window role objects
static void DestroyShellObjects(void)
{
    DestroyDecorations();

    if (platform.toplevelDecoration != NULL) zxdg_toplevel_decoration_v1_destroy(platform.toplevelDecoration);
    platform.toplevelDecoration = NULL;

    if (platform.libdecor.frame != NULL)
    {
        platform.libdecor.FrameUnref(platform.libdecor.frame);
        platform.libdecor.frame = NULL;
        platform.xdgToplevel = NULL;
    }

    if (platform.xdgToplevel != NULL) xdg_toplevel_destroy(platform.xdgToplevel);
    platform.xdgToplevel = NULL;

    if (platform.xdgSurface != NULL) xdg_surface_destroy(platform.xdgSurface);
    platform.xdgSurface = NULL;

    if (platform.surface != NULL)
    {
        wl_surface_attach(platform.surface, NULL, 0, 0);
        wl_surface_commit(platform.surface);
    }

    platform.decorationMode = DECORATION_MODE_NONE;
    platform.sizeLimitsResizable = -1;
    platform.decorationVisible = -1;
    platform.configured = false;
    platform.frameCommitted = false;
    platform.visible = false;
}

// Apply a configure event
static void ApplyConfigure(int width, int height, bool activated, bool maximized, bool fullscreen)
{
    platform.eventReceived = true;
    // Window state flags
    // NOTE: Window is considered focused until the compositor activated it once, like other platforms
    if (activated)
    {
        platform.activatedOnce = true;
        FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_UNFOCUSED);
    }
    else if (platform.activatedOnce) FLAG_SET(CORE.Window.flags, FLAG_WINDOW_UNFOCUSED);

    if (maximized) FLAG_SET(CORE.Window.flags, FLAG_WINDOW_MAXIMIZED);
    else FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_MAXIMIZED);

    // Fullscreen state changed by the compositor
    // NOTE: Configure events sent before the last fullscreen request was applied do not change the flags
    if (platform.fullscreenPending && (fullscreen == platform.fullscreenRequested)) platform.fullscreenPending = false;

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
    }

    // NOTE: Leaving maximized/fullscreen or changing decorations mode keeps the floating size, the compositor size is a suggestion
    bool restored = ((platform.fullscreen || platform.maximized || platform.decorationModeChanged) && !fullscreen && !maximized);
    platform.decorationModeChanged = false;

    platform.fullscreen = fullscreen;
    platform.maximized = maximized;

    // Configure size includes the builtin decorations titlebar (window geometry)
    int titlebarHeight = GetTitlebarHeight();

    int newWidth = platform.width;
    int newHeight = platform.height;

    // NOTE: Floating non-resizable windows keep their size, the compositor size is only a suggestion
    if ((width > 0) && (height > 0) && (fullscreen || maximized || (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_RESIZABLE) && !restored)))
    {
        newWidth = width;
        newHeight = height - titlebarHeight;
    }
    else if (!fullscreen && !maximized)
    {
        newWidth = platform.windowedWidth;
        newHeight = platform.windowedHeight;
    }

    if (newWidth < 1) newWidth = 1;
    if (newHeight < 1) newHeight = 1;

    if (!fullscreen && !maximized)
    {
        platform.windowedWidth = newWidth;
        platform.windowedHeight = newHeight;
    }

    platform.width = newWidth;
    platform.height = newHeight;

    UpdateWindowBuffers();
    UpdateWindowGeometry();

    platform.decorationDirty = true;
    UpdateDecorations();
}

// Update window scale from compositor preferences
static void UpdateWindowScale(void)
{
    if (platform.surface == NULL) return;

    float scale = 1.0f;

    if ((platform.fractionalScale120 > 0) && (platform.viewport != NULL)) scale = (float)platform.fractionalScale120/120.0f;
    else if (platform.compositorVersion >= 6) scale = (float)platform.preferredBufferScale;
    else
    {
        // Use the highest scale of the outputs the window is on
        for (int i = 0; i < platform.enteredOutputCount; i++)
        {
            int monitor = GetMonitorIndex(platform.enteredOutputs[i]);
            if ((monitor >= 0) && ((float)platform.monitors[monitor].scale > scale)) scale = (float)platform.monitors[monitor].scale;
        }
    }

    if (scale <= 0.0f) scale = 1.0f;

    if (scale != platform.scale)
    {
        platform.scale = scale;
        platform.eventReceived = true;

        UpdateWindowBuffers();
        platform.decorationDirty = true;
        UpdateDecorations();
        UpdateCursor();
    }
}

// Apply window size and scale to the framebuffer
static void UpdateWindowBuffers(void)
{
    float scale = FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_HIGHDPI)? platform.scale : 1.0f;

    platform.bufferWidth = (int)roundf((float)platform.width*scale);
    platform.bufferHeight = (int)roundf((float)platform.height*scale);
    if (platform.bufferWidth < 1) platform.bufferWidth = 1;
    if (platform.bufferHeight < 1) platform.bufferHeight = 1;

#if !defined(GRAPHICS_API_OPENGL_SOFTWARE)
    if (platform.eglWindow != NULL) wl_egl_window_resize(platform.eglWindow, platform.bufferWidth, platform.bufferHeight, 0, 0);
#endif

    // NOTE: Buffer scale is set with the next buffer attached, in SwapScreenBuffer()
    if (platform.viewport != NULL) wp_viewport_set_destination(platform.viewport, platform.width, platform.height);
    else platform.bufferScale = (int)scale;

    // Opaque region allows the compositor to skip drawing what is behind the window
    if (!FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_TRANSPARENT))
    {
        struct wl_region *region = wl_compositor_create_region(platform.compositor);
        wl_region_add(region, 0, 0, platform.width, platform.height);
        wl_surface_set_opaque_region(platform.surface, region);
        wl_region_destroy(region);
    }

    if (FLAG_IS_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE) || FLAG_IS_SET(CORE.Window.flags, FLAG_BORDERLESS_WINDOWED_MODE))
    {
        int monitor = GetCurrentMonitor();
        if (monitor < platform.monitorCount)
        {
            CORE.Window.display.width = platform.monitors[monitor].width;
            CORE.Window.display.height = platform.monitors[monitor].height;
        }
    }

    // Set screen size to logical pixel size and render size to framebuffer size
    // NOTE: Mouse coordinates are already in logical space, no mouse scaling required
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

// Update window geometry and size limits
static void UpdateWindowGeometry(void)
{
    int titlebarHeight = GetTitlebarHeight();

    int minWidth = 0;
    int minHeight = 0;
    int maxWidth = 0;
    int maxHeight = 0;

    if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_RESIZABLE))
    {
        minWidth = CORE.Window.screenMin.width;
        minHeight = CORE.Window.screenMin.height;
        maxWidth = CORE.Window.screenMax.width;
        maxHeight = CORE.Window.screenMax.height;
    }
    else
    {
        minWidth = platform.windowedWidth;
        minHeight = platform.windowedHeight;
        maxWidth = platform.windowedWidth;
        maxHeight = platform.windowedHeight;
    }

    // NOTE: Size limits include the requested titlebar also in fullscreen, they are only sent when changed (every change produces a configure)
    bool clientSideRequested = ((platform.toplevelDecoration != NULL) && (platform.decorationModeSent == ZXDG_TOPLEVEL_DECORATION_V1_MODE_CLIENT_SIDE));
    bool builtinTitlebar = (platform.decorationRequested && (platform.libdecor.frame == NULL) &&
        ((platform.decorationMode == DECORATION_MODE_BUILTIN) || clientSideRequested));
    int limitsTitlebarHeight = builtinTitlebar? DECORATION_TITLEBAR_HEIGHT : 0;
    if ((platform.libdecor.frame == NULL) && (minHeight > 0)) minHeight += limitsTitlebarHeight;
    if ((platform.libdecor.frame == NULL) && (maxHeight > 0)) maxHeight += limitsTitlebarHeight;

    int resizable = FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_RESIZABLE)? 1 : 0;
    bool limitsChanged = ((minWidth != platform.sentMinWidth) || (minHeight != platform.sentMinHeight) ||
        (maxWidth != platform.sentMaxWidth) || (maxHeight != platform.sentMaxHeight) || (resizable != platform.sizeLimitsResizable));

    platform.sentMinWidth = minWidth;
    platform.sentMinHeight = minHeight;
    platform.sentMaxWidth = maxWidth;
    platform.sentMaxHeight = maxHeight;
    platform.sizeLimitsResizable = resizable;

    if (platform.libdecor.frame != NULL)
    {
        if (limitsChanged)
        {
            // NOTE: libdecor restores its own saved limits when the resize capability changes, limits are set after it
            if (resizable == 1) platform.libdecor.FrameSetCapabilities(platform.libdecor.frame, LIBDECOR_ACTION_RESIZE | LIBDECOR_ACTION_FULLSCREEN);
            else platform.libdecor.FrameUnsetCapabilities(platform.libdecor.frame, LIBDECOR_ACTION_RESIZE | LIBDECOR_ACTION_FULLSCREEN);

            platform.libdecor.FrameSetMinContentSize(platform.libdecor.frame, minWidth, minHeight);
            platform.libdecor.FrameSetMaxContentSize(platform.libdecor.frame, maxWidth, maxHeight);

            // NOTE: libdecor only applies the limits on frame commit
            if (platform.configured)
            {
                struct libdecor_state *frameState = platform.libdecor.StateNew(platform.width, platform.height);
                platform.libdecor.FrameCommit(platform.libdecor.frame, frameState, NULL);
                platform.libdecor.StateFree(frameState);
            }
        }
    }
    else if (platform.xdgToplevel != NULL)
    {
        if (limitsChanged)
        {
            xdg_toplevel_set_min_size(platform.xdgToplevel, minWidth, minHeight);
            xdg_toplevel_set_max_size(platform.xdgToplevel, maxWidth, maxHeight);
        }

        if (platform.configured) xdg_surface_set_window_geometry(platform.xdgSurface, 0, -titlebarHeight, platform.width, platform.height + titlebarHeight);
    }
}

// Request fullscreen state to the compositor
static void RequestFullscreen(bool fullscreen, struct wl_output *output)
{
    if ((platform.libdecor.frame == NULL) && (platform.xdgToplevel == NULL)) return;

    if (platform.libdecor.frame != NULL)
    {
        if (fullscreen) platform.libdecor.FrameSetFullscreen(platform.libdecor.frame, output);
        else platform.libdecor.FrameUnsetFullscreen(platform.libdecor.frame);
    }
    else if (fullscreen) xdg_toplevel_set_fullscreen(platform.xdgToplevel, output);
    else xdg_toplevel_unset_fullscreen(platform.xdgToplevel);

    platform.fullscreenRequested = fullscreen;
    platform.fullscreenPending = true;
}

// Set window icon from the icon buffers, default icon if none
static void UpdateWindowIcon(void)
{
    if ((platform.iconManager == NULL) || (platform.xdgToplevel == NULL)) return;

    if (platform.iconBufferCount == 0)
    {
        xdg_toplevel_icon_manager_v1_set_icon(platform.iconManager, platform.xdgToplevel, NULL);
        return;
    }

    struct xdg_toplevel_icon_v1 *icon = xdg_toplevel_icon_manager_v1_create_icon(platform.iconManager);
    for (int i = 0; i < platform.iconBufferCount; i++) xdg_toplevel_icon_v1_add_buffer(icon, platform.iconBuffers[i].buffer, 1);
    xdg_toplevel_icon_manager_v1_set_icon(platform.iconManager, platform.xdgToplevel, icon);
    xdg_toplevel_icon_v1_destroy(icon);
}

// Get monitor index for an output
static int GetMonitorIndex(struct wl_output *output)
{
    for (int i = 0; i < platform.monitorCount; i++)
    {
        if (platform.monitors[i].output == output) return i;
    }

    return -1;
}

// Get output for a monitor index
static struct wl_output *GetMonitorOutput(int monitor)
{
    if ((monitor >= 0) && (monitor < platform.monitorCount)) return platform.monitors[monitor].output;

    return NULL;
}

// Create shared memory buffer (ARGB8888)
static bool CreateShmBuffer(ShmBuffer *buffer, int width, int height)
{
    memset(buffer, 0, sizeof(ShmBuffer));

    const char *runtimeDir = getenv("XDG_RUNTIME_DIR");
    if (runtimeDir == NULL) runtimeDir = "/tmp";

    char path[512] = { 0 };
    snprintf(path, sizeof(path), "%s/raylib-shm-XXXXXX", runtimeDir);

    int fd = mkstemp(path);
    if (fd < 0)
    {
        TRACELOG(LOG_WARNING, "WAYLAND: Failed to create shared memory file");
        return false;
    }

    unlink(path);
    fcntl(fd, F_SETFD, FD_CLOEXEC);

    int stride = width*4;
    size_t size = (size_t)stride*height;

    if (ftruncate(fd, (off_t)size) < 0)
    {
        close(fd);
        TRACELOG(LOG_WARNING, "WAYLAND: Failed to allocate shared memory");
        return false;
    }

    void *pixels = mmap(NULL, size, PROT_READ | PROT_WRITE, MAP_SHARED, fd, 0);
    if (pixels == MAP_FAILED)
    {
        close(fd);
        TRACELOG(LOG_WARNING, "WAYLAND: Failed to map shared memory");
        return false;
    }

    struct wl_shm_pool *pool = wl_shm_create_pool(platform.shm, fd, (int32_t)size);
    buffer->buffer = wl_shm_pool_create_buffer(pool, 0, width, height, stride, WL_SHM_FORMAT_ARGB8888);
    wl_shm_pool_destroy(pool);
    close(fd);

    buffer->pixels = (unsigned int *)pixels;
    buffer->width = width;
    buffer->height = height;
    buffer->size = size;
    wl_buffer_add_listener(buffer->buffer, &bufferListener, buffer);

    return true;
}

// Destroy shared memory buffer
static void DestroyShmBuffer(ShmBuffer *buffer)
{
    if (buffer->buffer != NULL) wl_buffer_destroy(buffer->buffer);
    if (buffer->pixels != NULL) munmap(buffer->pixels, buffer->size);

    memset(buffer, 0, sizeof(ShmBuffer));
}

// Release seat pointer (destroy on old compositors)
static void ReleasePointer(struct wl_pointer *pointer)
{
    if (wl_pointer_get_version(pointer) >= WL_POINTER_RELEASE_SINCE_VERSION) wl_pointer_release(pointer);
    else wl_pointer_destroy(pointer);
}

// Release seat keyboard (destroy on old compositors)
static void ReleaseKeyboard(struct wl_keyboard *keyboard)
{
    if (wl_keyboard_get_version(keyboard) >= WL_KEYBOARD_RELEASE_SINCE_VERSION) wl_keyboard_release(keyboard);
    else wl_keyboard_destroy(keyboard);
}

// Release seat touch device (destroy on old compositors)
static void ReleaseTouch(struct wl_touch *touch)
{
    if (wl_touch_get_version(touch) >= WL_TOUCH_RELEASE_SINCE_VERSION) wl_touch_release(touch);
    else wl_touch_destroy(touch);
}

// Release output (destroy on old compositors)
static void ReleaseOutput(struct wl_output *output)
{
    if (wl_output_get_version(output) >= WL_OUTPUT_RELEASE_SINCE_VERSION) wl_output_release(output);
    else wl_output_destroy(output);
}

// Register key state change
static void HandleKey(unsigned int scancode, bool pressed)
{
    platform.eventReceived = true;
    int key = (scancode < KEYMAP_SIZE)? linuxToRaylibMap[scancode] : (int)KEY_NULL;
    if ((key <= KEY_NULL) || (key >= MAX_KEYBOARD_KEYS)) return;

    if (pressed)
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
    else CORE.Input.Keyboard.currentKeyState[key] = 0;

    // NOTE: CAPS/NUM lock keys are considered down while the lock is active
    if (platform.xkbState != NULL)
    {
        bool capsLocked = (xkb_state_mod_name_is_active(platform.xkbState, XKB_MOD_NAME_CAPS, XKB_STATE_MODS_LOCKED) > 0);
        bool numLocked = (xkb_state_mod_name_is_active(platform.xkbState, XKB_MOD_NAME_NUM, XKB_STATE_MODS_LOCKED) > 0);

        if ((key == KEY_CAPS_LOCK) && capsLocked) CORE.Input.Keyboard.currentKeyState[key] = 1;
        if ((key == KEY_NUM_LOCK) && numLocked) CORE.Input.Keyboard.currentKeyState[key] = 1;
    }
}

// Register character input for a key, considering keyboard layout and compose sequences
static void HandleKeyChar(unsigned int scancode)
{
    if (platform.xkbState == NULL) return;

    xkb_keysym_t keysym = xkb_state_key_get_one_sym(platform.xkbState, scancode + XKB_KEYCODE_OFFSET);

    if ((platform.xkbComposeState != NULL) && (xkb_compose_state_feed(platform.xkbComposeState, keysym) == XKB_COMPOSE_FEED_ACCEPTED))
    {
        switch (xkb_compose_state_get_status(platform.xkbComposeState))
        {
            case XKB_COMPOSE_COMPOSING: keysym = XKB_KEY_NoSymbol; break;
            case XKB_COMPOSE_COMPOSED:
            {
                keysym = xkb_compose_state_get_one_sym(platform.xkbComposeState);
                xkb_compose_state_reset(platform.xkbComposeState);
            } break;
            case XKB_COMPOSE_CANCELLED:
            {
                keysym = XKB_KEY_NoSymbol;
                xkb_compose_state_reset(platform.xkbComposeState);
            } break;
            default: break;
        }
    }

    unsigned int codepoint = xkb_keysym_to_utf32(keysym);

    // Check if there is space available in the queue, control characters are not registered
    if ((codepoint >= 32) && (codepoint != 127) && (CORE.Input.Keyboard.charPressedQueueCount < MAX_CHAR_PRESSED_QUEUE))
    {
        // Add character to the queue
        CORE.Input.Keyboard.charPressedQueue[CORE.Input.Keyboard.charPressedQueueCount] = codepoint;
        CORE.Input.Keyboard.charPressedQueueCount++;
    }
}

// Generate key repeat events
static void ProcessKeyRepeat(void)
{
    if ((platform.repeatScancode == 0) || (platform.repeatRate <= 0)) return;

    double time = GetTime();
    int repeats = 0;

    while ((time >= platform.repeatTime) && (repeats < MAX_KEY_REPEATS_PER_FRAME))
    {
        int key = (platform.repeatScancode < KEYMAP_SIZE)? linuxToRaylibMap[platform.repeatScancode] : (int)KEY_NULL;
        if ((key > KEY_NULL) && (key < MAX_KEYBOARD_KEYS)) CORE.Input.Keyboard.keyRepeatInFrame[key] = 1;

        HandleKeyChar(platform.repeatScancode);

        platform.repeatTime += 1.0/(double)platform.repeatRate;
        repeats++;
    }

    if (time >= platform.repeatTime) platform.repeatTime = time + 1.0/(double)platform.repeatRate;
}

// Register mouse button state change
static void HandleMouseButton(int button, bool pressed)
{
    platform.eventReceived = true;
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

// Register touch event
static void HandleTouch(int touchAction, int slot, int id, double x, double y)
{
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

    if ((touchAction == TOUCH_ACTION_UP) || (touchAction == TOUCH_ACTION_CANCEL)) platform.touchIds[slot] = -1;

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
}

// Update cursor image for current pointer focus
static void UpdateCursor(void)
{
    if ((platform.pointer == NULL) || (platform.pointerFocus == NULL)) return;

    int part = GetDecorationPart(platform.pointerFocus);
    if (part >= 0)
    {
        switch (GetDecorationEdge(part, platform.pointerX, platform.pointerY))
        {
            case XDG_TOPLEVEL_RESIZE_EDGE_TOP: SetCursorImage(WP_CURSOR_SHAPE_DEVICE_V1_SHAPE_N_RESIZE, "n-resize", "top_side"); break;
            case XDG_TOPLEVEL_RESIZE_EDGE_BOTTOM: SetCursorImage(WP_CURSOR_SHAPE_DEVICE_V1_SHAPE_S_RESIZE, "s-resize", "bottom_side"); break;
            case XDG_TOPLEVEL_RESIZE_EDGE_LEFT: SetCursorImage(WP_CURSOR_SHAPE_DEVICE_V1_SHAPE_W_RESIZE, "w-resize", "left_side"); break;
            case XDG_TOPLEVEL_RESIZE_EDGE_RIGHT: SetCursorImage(WP_CURSOR_SHAPE_DEVICE_V1_SHAPE_E_RESIZE, "e-resize", "right_side"); break;
            case XDG_TOPLEVEL_RESIZE_EDGE_TOP_LEFT: SetCursorImage(WP_CURSOR_SHAPE_DEVICE_V1_SHAPE_NW_RESIZE, "nw-resize", "top_left_corner"); break;
            case XDG_TOPLEVEL_RESIZE_EDGE_TOP_RIGHT: SetCursorImage(WP_CURSOR_SHAPE_DEVICE_V1_SHAPE_NE_RESIZE, "ne-resize", "top_right_corner"); break;
            case XDG_TOPLEVEL_RESIZE_EDGE_BOTTOM_LEFT: SetCursorImage(WP_CURSOR_SHAPE_DEVICE_V1_SHAPE_SW_RESIZE, "sw-resize", "bottom_left_corner"); break;
            case XDG_TOPLEVEL_RESIZE_EDGE_BOTTOM_RIGHT: SetCursorImage(WP_CURSOR_SHAPE_DEVICE_V1_SHAPE_SE_RESIZE, "se-resize", "bottom_right_corner"); break;
            default: SetCursorImage(WP_CURSOR_SHAPE_DEVICE_V1_SHAPE_DEFAULT, "default", "left_ptr"); break;
        }

        return;
    }

    if (platform.pointerFocus != platform.surface) return;

    if (CORE.Input.Mouse.cursorHidden) wl_pointer_set_cursor(platform.pointer, platform.pointerEnterSerial, NULL, 0, 0);
    else SetCursorShape(CORE.Input.Mouse.cursor);
}

// Set cursor shape (MouseCursor) on pointer focus
static void SetCursorShape(int cursor)
{
    switch (cursor)
    {
        case MOUSE_CURSOR_IBEAM: SetCursorImage(WP_CURSOR_SHAPE_DEVICE_V1_SHAPE_TEXT, "text", "xterm"); break;
        case MOUSE_CURSOR_CROSSHAIR: SetCursorImage(WP_CURSOR_SHAPE_DEVICE_V1_SHAPE_CROSSHAIR, "crosshair", "cross"); break;
        case MOUSE_CURSOR_POINTING_HAND: SetCursorImage(WP_CURSOR_SHAPE_DEVICE_V1_SHAPE_POINTER, "pointer", "hand2"); break;
        case MOUSE_CURSOR_RESIZE_EW: SetCursorImage(WP_CURSOR_SHAPE_DEVICE_V1_SHAPE_EW_RESIZE, "ew-resize", "sb_h_double_arrow"); break;
        case MOUSE_CURSOR_RESIZE_NS: SetCursorImage(WP_CURSOR_SHAPE_DEVICE_V1_SHAPE_NS_RESIZE, "ns-resize", "sb_v_double_arrow"); break;
        case MOUSE_CURSOR_RESIZE_NWSE: SetCursorImage(WP_CURSOR_SHAPE_DEVICE_V1_SHAPE_NWSE_RESIZE, "nwse-resize", "bottom_right_corner"); break;
        case MOUSE_CURSOR_RESIZE_NESW: SetCursorImage(WP_CURSOR_SHAPE_DEVICE_V1_SHAPE_NESW_RESIZE, "nesw-resize", "bottom_left_corner"); break;
        case MOUSE_CURSOR_RESIZE_ALL: SetCursorImage(WP_CURSOR_SHAPE_DEVICE_V1_SHAPE_ALL_SCROLL, "all-scroll", "fleur"); break;
        case MOUSE_CURSOR_NOT_ALLOWED: SetCursorImage(WP_CURSOR_SHAPE_DEVICE_V1_SHAPE_NOT_ALLOWED, "not-allowed", "crossed_circle"); break;
        default: SetCursorImage(WP_CURSOR_SHAPE_DEVICE_V1_SHAPE_DEFAULT, "default", "left_ptr"); break;
    }
}

// Set cursor image, cursor-shape protocol or cursor theme
static void SetCursorImage(unsigned int shape, const char *name, const char *fallbackName)
{
    if (platform.cursorShapeDevice != NULL)
    {
        wp_cursor_shape_device_v1_set_shape(platform.cursorShapeDevice, platform.pointerEnterSerial, shape);
        return;
    }

    int scale = (int)ceilf(platform.scale);
    if (scale < 1) scale = 1;

    if ((platform.cursorTheme == NULL) || (platform.cursorThemeScale != scale))
    {
        if (platform.cursorTheme != NULL) wl_cursor_theme_destroy(platform.cursorTheme);

        int size = CURSOR_SIZE_DEFAULT;
        const char *sizeString = getenv("XCURSOR_SIZE");
        if (sizeString != NULL)
        {
            size = atoi(sizeString);
            if ((size <= 0) || (size > CURSOR_SIZE_MAX)) size = CURSOR_SIZE_DEFAULT;
        }

        platform.cursorTheme = wl_cursor_theme_load(getenv("XCURSOR_THEME"), size*scale, platform.shm);
        platform.cursorThemeScale = scale;
    }

    if (platform.cursorTheme == NULL) return;

    struct wl_cursor *cursor = wl_cursor_theme_get_cursor(platform.cursorTheme, name);
    if (cursor == NULL) cursor = wl_cursor_theme_get_cursor(platform.cursorTheme, fallbackName);
    if (cursor == NULL) cursor = wl_cursor_theme_get_cursor(platform.cursorTheme, "left_ptr");
    if ((cursor == NULL) || (cursor->image_count == 0)) return;

    struct wl_cursor_image *image = cursor->images[0];
    struct wl_buffer *buffer = wl_cursor_image_get_buffer(image);
    if (buffer == NULL) return;

    if (platform.cursorSurface == NULL) platform.cursorSurface = wl_compositor_create_surface(platform.compositor);

    // NOTE: Cursor image size must be a multiple of the buffer scale, themes do not provide every size
    int bufferScale = scale;
    while ((bufferScale > 1) && (((image->width%bufferScale) != 0) || ((image->height%bufferScale) != 0))) bufferScale--;

    wl_surface_set_buffer_scale(platform.cursorSurface, bufferScale);
    wl_surface_attach(platform.cursorSurface, buffer, 0, 0);
    wl_surface_damage_buffer(platform.cursorSurface, 0, 0, image->width, image->height);
    wl_surface_commit(platform.cursorSurface);
    wl_pointer_set_cursor(platform.pointer, platform.pointerEnterSerial, platform.cursorSurface, image->hotspot_x/bufferScale, image->hotspot_y/bufferScale);
}

// Lock or unlock pointer (disabled cursor)
static void LockPointer(bool lock)
{
    if (lock)
    {
        if ((platform.lockedPointer != NULL) || (platform.pointer == NULL) || (platform.pointerConstraints == NULL)) return;

        platform.lockedPointer = zwp_pointer_constraints_v1_lock_pointer(platform.pointerConstraints, platform.surface,
            platform.pointer, NULL, ZWP_POINTER_CONSTRAINTS_V1_LIFETIME_PERSISTENT);

        if (platform.relativePointerManager != NULL)
        {
            platform.relativePointer = zwp_relative_pointer_manager_v1_get_relative_pointer(platform.relativePointerManager, platform.pointer);
            zwp_relative_pointer_v1_add_listener(platform.relativePointer, &relativePointerListener, NULL);
        }
        else TRACELOG(LOG_WARNING, "WAYLAND: Relative pointer motion not supported by the compositor");
    }
    else
    {
        if (platform.relativePointer != NULL) zwp_relative_pointer_v1_destroy(platform.relativePointer);
        platform.relativePointer = NULL;

        if (platform.lockedPointer != NULL) zwp_locked_pointer_v1_destroy(platform.lockedPointer);
        platform.lockedPointer = NULL;
    }
}

// Destroy data offer and its mime types
static void DestroyDataOffer(DataOffer *offer)
{
    if (offer == NULL) return;

    for (int i = 0; i < offer->mimeTypeCount; i++) RL_FREE(offer->mimeTypes[i]);
    if (offer->offer != NULL) wl_data_offer_destroy(offer->offer);
    RL_FREE(offer);
}

// Check data offer mime type
static bool DataOfferHasMimeType(DataOffer *offer, const char *mimeType)
{
    if (offer == NULL) return false;

    for (int i = 0; i < offer->mimeTypeCount; i++)
    {
        if (strcmp(offer->mimeTypes[i], mimeType) == 0) return true;
    }

    return false;
}

// Read data offer content, returned data is allocated and zero-terminated
static char *ReadDataOffer(DataOffer *offer, const char *mimeType, int *size)
{
    int fds[2] = { -1, -1 };
    if (pipe(fds) != 0) return NULL;

    fcntl(fds[0], F_SETFD, FD_CLOEXEC);
    fcntl(fds[1], F_SETFD, FD_CLOEXEC);

    wl_data_offer_receive(offer->offer, mimeType, fds[1]);
    close(fds[1]);
    wl_display_flush(platform.display);

    int capacity = DATA_OFFER_READ_CHUNK;
    int length = 0;
    char *buffer = (char *)RL_MALLOC(capacity);

    while (true)
    {
        // NOTE: Waiting is limited in time, the source client could be unresponsive
        struct pollfd fd = { fds[0], POLLIN, 0 };
        int result = poll(&fd, 1, DATA_OFFER_READ_TIMEOUT);
        if ((result < 0) && (errno == EINTR)) continue;
        if (result <= 0) break;

        if ((capacity - length) < DATA_OFFER_READ_CHUNK)
        {
            capacity *= 2;
            buffer = (char *)RL_REALLOC(buffer, capacity);
        }

        ssize_t count = read(fds[0], buffer + length, capacity - length - 1);
        if ((count < 0) && (errno == EINTR)) continue;
        if (count <= 0) break;

        length += (int)count;
    }

    close(fds[0]);

    buffer[length] = '\0';
    if (size != NULL) *size = length;

    if (length == 0)
    {
        RL_FREE(buffer);
        return NULL;
    }

    return buffer;
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

// Load libdecor library
static bool LoadLibdecor(void)
{
    if (platform.libdecor.context != NULL) return true;

    platform.libdecor.handle = dlopen("libdecor-0.so.0", RTLD_LAZY | RTLD_LOCAL);
    if (platform.libdecor.handle == NULL)
    {
        TRACELOG(LOG_INFO, "WAYLAND: libdecor not available, using builtin decorations");
        return false;
    }

    void *handle = platform.libdecor.handle;
    platform.libdecor.New = dlsym(handle, "libdecor_new");
    platform.libdecor.Unref = dlsym(handle, "libdecor_unref");
    platform.libdecor.Dispatch = dlsym(handle, "libdecor_dispatch");
    platform.libdecor.Decorate = dlsym(handle, "libdecor_decorate");
    platform.libdecor.FrameUnref = dlsym(handle, "libdecor_frame_unref");
    platform.libdecor.FrameSetAppId = dlsym(handle, "libdecor_frame_set_app_id");
    platform.libdecor.FrameSetTitle = dlsym(handle, "libdecor_frame_set_title");
    platform.libdecor.FrameSetMinimized = dlsym(handle, "libdecor_frame_set_minimized");
    platform.libdecor.FrameSetFullscreen = dlsym(handle, "libdecor_frame_set_fullscreen");
    platform.libdecor.FrameUnsetFullscreen = dlsym(handle, "libdecor_frame_unset_fullscreen");
    platform.libdecor.FrameMap = dlsym(handle, "libdecor_frame_map");
    platform.libdecor.FrameCommit = dlsym(handle, "libdecor_frame_commit");
    platform.libdecor.FrameSetMinContentSize = dlsym(handle, "libdecor_frame_set_min_content_size");
    platform.libdecor.FrameSetMaxContentSize = dlsym(handle, "libdecor_frame_set_max_content_size");
    platform.libdecor.FrameSetMaximized = dlsym(handle, "libdecor_frame_set_maximized");
    platform.libdecor.FrameUnsetMaximized = dlsym(handle, "libdecor_frame_unset_maximized");
    platform.libdecor.FrameSetCapabilities = dlsym(handle, "libdecor_frame_set_capabilities");
    platform.libdecor.FrameUnsetCapabilities = dlsym(handle, "libdecor_frame_unset_capabilities");
    platform.libdecor.FrameSetVisibility = dlsym(handle, "libdecor_frame_set_visibility");
    platform.libdecor.FrameGetXdgToplevel = dlsym(handle, "libdecor_frame_get_xdg_toplevel");
    platform.libdecor.ConfigurationGetContentSize = dlsym(handle, "libdecor_configuration_get_content_size");
    platform.libdecor.ConfigurationGetWindowState = dlsym(handle, "libdecor_configuration_get_window_state");
    platform.libdecor.StateNew = dlsym(handle, "libdecor_state_new");
    platform.libdecor.StateFree = dlsym(handle, "libdecor_state_free");

    if ((platform.libdecor.New == NULL) || (platform.libdecor.Unref == NULL) || (platform.libdecor.Dispatch == NULL) ||
        (platform.libdecor.Decorate == NULL) || (platform.libdecor.FrameUnref == NULL) || (platform.libdecor.FrameSetAppId == NULL) ||
        (platform.libdecor.FrameSetTitle == NULL) || (platform.libdecor.FrameSetMinimized == NULL) || (platform.libdecor.FrameSetFullscreen == NULL) ||
        (platform.libdecor.FrameUnsetFullscreen == NULL) || (platform.libdecor.FrameMap == NULL) || (platform.libdecor.FrameCommit == NULL) ||
        (platform.libdecor.FrameSetMinContentSize == NULL) || (platform.libdecor.FrameSetMaxContentSize == NULL) ||
        (platform.libdecor.FrameSetMaximized == NULL) || (platform.libdecor.FrameUnsetMaximized == NULL) ||
        (platform.libdecor.FrameSetCapabilities == NULL) || (platform.libdecor.FrameUnsetCapabilities == NULL) ||
        (platform.libdecor.FrameSetVisibility == NULL) || (platform.libdecor.FrameGetXdgToplevel == NULL) ||
        (platform.libdecor.ConfigurationGetContentSize == NULL) || (platform.libdecor.ConfigurationGetWindowState == NULL) ||
        (platform.libdecor.StateNew == NULL) || (platform.libdecor.StateFree == NULL))
    {
        TRACELOG(LOG_WARNING, "WAYLAND: Failed to load libdecor functions, using builtin decorations");
        UnloadLibdecor();
        return false;
    }

    platform.libdecor.context = platform.libdecor.New(platform.display, &libdecorInterface);
    if (platform.libdecor.context == NULL)
    {
        TRACELOG(LOG_WARNING, "WAYLAND: Failed to initialize libdecor, using builtin decorations");
        UnloadLibdecor();
        return false;
    }

    // NOTE: libdecor plugins get their globals asynchronously, frames can not be created before
    wl_display_roundtrip(platform.display);

    TRACELOG(LOG_INFO, "WAYLAND: libdecor loaded successfully");

    return true;
}

// Unload libdecor library
static void UnloadLibdecor(void)
{
    if (platform.libdecor.frame != NULL) platform.libdecor.FrameUnref(platform.libdecor.frame);
    if (platform.libdecor.context != NULL) platform.libdecor.Unref(platform.libdecor.context);
    if (platform.libdecor.handle != NULL) dlclose(platform.libdecor.handle);

    memset(&platform.libdecor, 0, sizeof(LibdecorData));
}

// Create client-side decorations parts (raylib)
static void CreateDecorations(void)
{
    if ((platform.subcompositor == NULL) || (platform.decorationParts[DECORATION_PART_TITLEBAR].surface != NULL)) return;

    for (int i = 0; i < DECORATION_PART_COUNT; i++)
    {
        DecorationPart *part = &platform.decorationParts[i];

        part->surface = wl_compositor_create_surface(platform.compositor);
        part->subsurface = wl_subcompositor_get_subsurface(platform.subcompositor, part->surface, platform.surface);
        if (platform.viewporter != NULL) part->viewport = wp_viewporter_get_viewport(platform.viewporter, part->surface);
    }

    platform.decorationDirty = true;
}

// Destroy client-side decorations parts (raylib)
static void DestroyDecorations(void)
{
    for (int i = 0; i < DECORATION_PART_COUNT; i++)
    {
        DecorationPart *part = &platform.decorationParts[i];

        if (part->viewport != NULL) wp_viewport_destroy(part->viewport);
        if (part->subsurface != NULL) wl_subsurface_destroy(part->subsurface);
        if (part->surface != NULL) wl_surface_destroy(part->surface);
        DestroyShmBuffer(&part->buffers[0]);
        DestroyShmBuffer(&part->buffers[1]);

        memset(part, 0, sizeof(DecorationPart));
    }

    if (platform.pointerFocus != platform.surface) platform.pointerFocus = NULL;
    platform.decorationHover = DECORATION_BUTTON_NONE;
    platform.decorationPressed = DECORATION_BUTTON_NONE;
}

// Get decorations part bounds (logical pixels, relative to the window content)
static Rectangle GetDecorationPartBounds(int part)
{
    float width = (float)platform.width;
    float height = (float)platform.height;
    float titlebar = (float)DECORATION_TITLEBAR_HEIGHT;
    float border = (float)DECORATION_BORDER_SIZE;

    Rectangle bounds = { 0 };

    switch (part)
    {
        case DECORATION_PART_TITLEBAR: bounds = (Rectangle){ 0.0f, -titlebar, width, titlebar }; break;
        case DECORATION_PART_LEFT: bounds = (Rectangle){ -border, -titlebar, border, height + titlebar }; break;
        case DECORATION_PART_RIGHT: bounds = (Rectangle){ width, -titlebar, border, height + titlebar }; break;
        case DECORATION_PART_BOTTOM: bounds = (Rectangle){ -border, height, width + 2.0f*border, border }; break;
        default: break;
    }

    return bounds;
}

// Attach decorations part buffer, considering window scale
static void AttachDecorationBuffer(DecorationPart *part, ShmBuffer *buffer, Rectangle bounds, int bufferScale)
{
    wl_surface_attach(part->surface, buffer->buffer, 0, 0);

    if (part->viewport != NULL) wp_viewport_set_destination(part->viewport, (int)bounds.width, (int)bounds.height);
    else wl_surface_set_buffer_scale(part->surface, bufferScale);

    wl_surface_damage_buffer(part->surface, 0, 0, buffer->width, buffer->height);
    wl_surface_commit(part->surface);
    buffer->busy = true;
}

// Commit window surface state, deferred to the next frame until the first frame was presented
static void CommitWindowSurface(void)
{
    if ((platform.surface != NULL) && platform.frameCommitted) wl_surface_commit(platform.surface);
}

// Get builtin titlebar height included in the window geometry
static int GetTitlebarHeight(void)
{
    bool shown = ((platform.decorationMode == DECORATION_MODE_BUILTIN) && (platform.subcompositor != NULL) && !platform.fullscreen);

    return shown? DECORATION_TITLEBAR_HEIGHT : 0;
}

// Update client-side decorations size and position
static void UpdateDecorations(void)
{
    if (platform.libdecor.frame != NULL)
    {
        int visible = platform.decorationRequested? 1 : 0;

        if (visible != platform.decorationVisible)
        {
            platform.decorationVisible = visible;
            platform.libdecor.FrameSetVisibility(platform.libdecor.frame, platform.decorationRequested);
        }
        return;
    }

    // Negotiate decorations mode with the compositor when requested state changed
    bool clientSide = true;

    if (platform.toplevelDecoration != NULL)
    {
        int mode = ZXDG_TOPLEVEL_DECORATION_V1_MODE_CLIENT_SIDE;
        if (platform.decorationRequested && !platform.preferClientDecorations) mode = ZXDG_TOPLEVEL_DECORATION_V1_MODE_SERVER_SIDE;

        if (mode != platform.decorationModeSent)
        {
            platform.decorationModeSent = mode;
            zxdg_toplevel_decoration_v1_set_mode(platform.toplevelDecoration, mode);
            platform.decorationModeChanged = true;
        }

        // NOTE: Client-side mode is not confirmed again by the compositor if it does not change
        clientSide = ((platform.decorationMode != DECORATION_MODE_SERVER) && (mode == ZXDG_TOPLEVEL_DECORATION_V1_MODE_CLIENT_SIDE));
    }

    if (clientSide && (platform.xdgToplevel != NULL))
    {
        DecorationMode decorationMode = platform.decorationRequested? DECORATION_MODE_BUILTIN : DECORATION_MODE_NONE;

        if (decorationMode != platform.decorationMode)
        {
            platform.decorationMode = decorationMode;
            platform.decorationModeChanged = true;
            UpdateWindowGeometry();
        }
    }

    bool visible = ((platform.decorationMode == DECORATION_MODE_BUILTIN) && (platform.xdgToplevel != NULL) && (platform.subcompositor != NULL) &&
        platform.visible && !platform.fullscreen);

    if (!visible)
    {
        if (platform.decorationParts[DECORATION_PART_TITLEBAR].surface != NULL)
        {
            DestroyDecorations();
            CommitWindowSurface();
        }
        return;
    }

    CreateDecorations();

    int bufferScale = (int)ceilf(platform.scale);
    float scale = (platform.viewporter != NULL)? platform.scale : (float)bufferScale;

    // Borders are transparent, only used for resizing
    for (int i = DECORATION_PART_LEFT; i < DECORATION_PART_COUNT; i++)
    {
        DecorationPart *part = &platform.decorationParts[i];
        Rectangle bounds = GetDecorationPartBounds(i);
        int width = (int)ceilf(bounds.width*scale);
        int height = (int)ceilf(bounds.height*scale);

        wl_subsurface_set_position(part->subsurface, (int)bounds.x, (int)bounds.y);

        if ((part->buffers[0].width != width) || (part->buffers[0].height != height))
        {
            DestroyShmBuffer(&part->buffers[0]);
            if (CreateShmBuffer(&part->buffers[0], width, height))
            {
                memset(part->buffers[0].pixels, 0, part->buffers[0].size);
                AttachDecorationBuffer(part, &part->buffers[0], bounds, bufferScale);
            }
        }
    }

    Rectangle titlebarBounds = GetDecorationPartBounds(DECORATION_PART_TITLEBAR);
    wl_subsurface_set_position(platform.decorationParts[DECORATION_PART_TITLEBAR].subsurface, (int)titlebarBounds.x, (int)titlebarBounds.y);

    DrawDecorations();
}

// Blend a color into a decorations buffer pixel
static void BlendDecorationPixel(ShmBuffer *buffer, int x, int y, unsigned int color, float coverage)
{
    if ((x < 0) || (y < 0) || (x >= buffer->width) || (y >= buffer->height) || (coverage <= 0.0f)) return;
    if (coverage > 1.0f) coverage = 1.0f;

    unsigned int *pixel = &buffer->pixels[y*buffer->width + x];
    float alpha = (float)((color >> 24) & 0xff)/255.0f*coverage;

    unsigned int dst = *pixel;
    float r = (float)((color >> 16) & 0xff)*alpha + (float)((dst >> 16) & 0xff)*(1.0f - alpha);
    float g = (float)((color >> 8) & 0xff)*alpha + (float)((dst >> 8) & 0xff)*(1.0f - alpha);
    float b = (float)(color & 0xff)*alpha + (float)(dst & 0xff)*(1.0f - alpha);
    float a = alpha*255.0f + (float)((dst >> 24) & 0xff)*(1.0f - alpha);

    *pixel = ((unsigned int)(a + 0.5f) << 24) | ((unsigned int)(r + 0.5f) << 16) | ((unsigned int)(g + 0.5f) << 8) | (unsigned int)(b + 0.5f);
}

// Draw an anti-aliased line into a decorations buffer
static void DrawDecorationLine(ShmBuffer *buffer, float x0, float y0, float x1, float y1, float thick, unsigned int color)
{
    int minX = (int)floorf(fminf(x0, x1) - thick);
    int maxX = (int)ceilf(fmaxf(x0, x1) + thick);
    int minY = (int)floorf(fminf(y0, y1) - thick);
    int maxY = (int)ceilf(fmaxf(y0, y1) + thick);

    float dx = x1 - x0;
    float dy = y1 - y0;
    float lengthSqr = dx*dx + dy*dy;

    for (int y = minY; y <= maxY; y++)
    {
        for (int x = minX; x <= maxX; x++)
        {
            float px = (float)x + 0.5f;
            float py = (float)y + 0.5f;
            float t = (lengthSqr > 0.0f)? ((px - x0)*dx + (py - y0)*dy)/lengthSqr : 0.0f;
            if (t < 0.0f) t = 0.0f;
            else if (t > 1.0f) t = 1.0f;

            float distX = px - (x0 + t*dx);
            float distY = py - (y0 + t*dy);
            float distance = sqrtf(distX*distX + distY*distY);

            BlendDecorationPixel(buffer, x, y, color, thick/2.0f - distance + 0.5f);
        }
    }
}

// Draw an anti-aliased circle into a decorations buffer
static void DrawDecorationCircle(ShmBuffer *buffer, float centerX, float centerY, float radius, unsigned int color)
{
    for (int y = (int)floorf(centerY - radius - 1.0f); y <= (int)ceilf(centerY + radius + 1.0f); y++)
    {
        for (int x = (int)floorf(centerX - radius - 1.0f); x <= (int)ceilf(centerX + radius + 1.0f); x++)
        {
            float distX = (float)x + 0.5f - centerX;
            float distY = (float)y + 0.5f - centerY;

            BlendDecorationPixel(buffer, x, y, color, radius - sqrtf(distX*distX + distY*distY) + 0.5f);
        }
    }
}

#if SUPPORT_MODULE_RTEXT
// Draw decorations title text
static void DrawDecorationText(ShmBuffer *buffer, Font font, const char *text, int posX, int posY, int textScale, unsigned int color)
{
    int length = TextLength(text);
    int x = posX;

    for (int i = 0; i < length;)
    {
        int codepointSize = 0;
        int codepoint = GetCodepointNext(&text[i], &codepointSize);
        int index = GetGlyphIndex(font, codepoint);
        GlyphInfo glyph = font.glyphs[index];
        int advance = (glyph.advanceX > 0)? glyph.advanceX : (int)font.recs[index].width;

        if ((glyph.image.data != NULL) && (glyph.image.format == PIXELFORMAT_UNCOMPRESSED_GRAY_ALPHA))
        {
            const unsigned char *pixels = (const unsigned char *)glyph.image.data;

            for (int gy = 0; gy < glyph.image.height; gy++)
            {
                for (int gx = 0; gx < glyph.image.width; gx++)
                {
                    float coverage = (float)pixels[(gy*glyph.image.width + gx)*2 + 1]/255.0f;
                    if (coverage <= 0.0f) continue;

                    for (int sy = 0; sy < textScale; sy++)
                    {
                        for (int sx = 0; sx < textScale; sx++)
                        {
                            BlendDecorationPixel(buffer, x + (glyph.offsetX + gx)*textScale + sx, posY + (glyph.offsetY + gy)*textScale + sy, color, coverage);
                        }
                    }
                }
            }
        }

        x += (advance + 1)*textScale;
        i += codepointSize;
    }
}
#endif

// Draw window icon into a decorations buffer, scaled down from the closest icon size
static void DrawDecorationIcon(ShmBuffer *buffer, int posX, int posY, int size)
{
    // Smallest icon not smaller than the required size, otherwise the biggest one
    ShmBuffer *icon = NULL;
    for (int i = 0; i < platform.iconBufferCount; i++)
    {
        ShmBuffer *candidate = &platform.iconBuffers[i];

        if (icon == NULL) icon = candidate;
        else if ((icon->width < size) && (candidate->width > icon->width)) icon = candidate;
        else if ((candidate->width >= size) && (candidate->width < icon->width)) icon = candidate;
    }

    if ((icon == NULL) || (icon->pixels == NULL) || (size <= 0)) return;

    // NOTE: Icon pixels are premultiplied ARGB, every target pixel averages its source area
    for (int y = 0; y < size; y++)
    {
        int srcY0 = y*icon->height/size;
        int srcY1 = (y + 1)*icon->height/size;
        if (srcY1 <= srcY0) srcY1 = srcY0 + 1;

        for (int x = 0; x < size; x++)
        {
            int srcX0 = x*icon->width/size;
            int srcX1 = (x + 1)*icon->width/size;
            if (srcX1 <= srcX0) srcX1 = srcX0 + 1;

            unsigned int sumA = 0;
            unsigned int sumR = 0;
            unsigned int sumG = 0;
            unsigned int sumB = 0;
            unsigned int count = 0;

            for (int sy = srcY0; sy < srcY1; sy++)
            {
                for (int sx = srcX0; sx < srcX1; sx++)
                {
                    unsigned int pixel = icon->pixels[sy*icon->width + sx];
                    sumA += (pixel >> 24) & 0xff;
                    sumR += (pixel >> 16) & 0xff;
                    sumG += (pixel >> 8) & 0xff;
                    sumB += pixel & 0xff;
                    count++;
                }
            }

            unsigned int a = sumA/count;
            if (a == 0) continue;

            unsigned int r = sumR*255/sumA;
            unsigned int g = sumG*255/sumA;
            unsigned int b = sumB*255/sumA;

            BlendDecorationPixel(buffer, posX + x, posY + y, (a << 24) | (r << 16) | (g << 8) | b, 1.0f);
        }
    }
}

// Draw client-side decorations titlebar
static void DrawDecorations(void)
{
    DecorationPart *part = &platform.decorationParts[DECORATION_PART_TITLEBAR];
    if (part->surface == NULL) return;

    ShmBuffer *buffer = NULL;
    if (!part->buffers[0].busy) buffer = &part->buffers[0];
    else if (!part->buffers[1].busy) buffer = &part->buffers[1];

    // NOTE: Both buffers in use by the compositor, titlebar is redrawn later
    if (buffer == NULL)
    {
        platform.decorationDirty = true;
        return;
    }

    int bufferScale = (int)ceilf(platform.scale);
    float scale = (platform.viewporter != NULL)? platform.scale : (float)bufferScale;
    Rectangle bounds = GetDecorationPartBounds(DECORATION_PART_TITLEBAR);
    int width = (int)ceilf(bounds.width*scale);
    int height = (int)ceilf(bounds.height*scale);

    if ((buffer->width != width) || (buffer->height != height))
    {
        DestroyShmBuffer(buffer);
        if (!CreateShmBuffer(buffer, width, height)) return;
    }

    bool active = !FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_UNFOCUSED);
    unsigned int backgroundColor = active? 0xffebebeb : 0xfffafafa;
    unsigned int borderColor = active? 0xffcdc7c2 : 0xffd5d0cc;
    unsigned int textColor = active? 0xff2e3436 : 0xff929595;
    unsigned int hoverColor = 0x1a000000;
    unsigned int pressedColor = 0x33000000;

    // Background with rounded top corners
    float radius = (platform.maximized? 0.0f : (float)DECORATION_CORNER_RADIUS*scale);
    for (int y = 0; y < height; y++)
    {
        for (int x = 0; x < width; x++)
        {
            float coverage = 1.0f;

            if (((float)y < radius) && (((float)x < radius) || ((float)x >= (float)width - radius)))
            {
                float centerX = ((float)x < radius)? radius : (float)width - radius;
                float distX = (float)x + 0.5f - centerX;
                float distY = (float)y + 0.5f - radius;
                coverage = radius - sqrtf(distX*distX + distY*distY) + 0.5f;
                if (coverage < 0.0f) coverage = 0.0f;
                else if (coverage > 1.0f) coverage = 1.0f;
            }

            unsigned int alpha = (unsigned int)(coverage*255.0f + 0.5f);
            unsigned int color = (y >= (height - (int)ceilf(scale)))? borderColor : backgroundColor;
            unsigned int r = ((color >> 16) & 0xff)*alpha/255;
            unsigned int g = ((color >> 8) & 0xff)*alpha/255;
            unsigned int b = (color & 0xff)*alpha/255;

            buffer->pixels[y*width + x] = (alpha << 24) | (r << 16) | (g << 8) | b;
        }
    }

    // Titlebar buttons, from right to left: close, maximize, minimize
    float buttonSize = (float)DECORATION_BUTTON_SIZE*scale;
    float buttonX = (float)width - (float)(DECORATION_BUTTON_SPACING + DECORATION_BUTTON_SIZE)*scale;
    float buttonY = ((float)height - buttonSize)/2.0f;
    float iconSize = 10.0f*scale;
    float thick = 1.5f*scale;
    int buttonCount = 0;

    for (int i = 0; i < 3; i++)
    {
        if ((titlebarButtons[i] == DECORATION_BUTTON_MAXIMIZE) && !FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_RESIZABLE)) continue;

        float centerX = buttonX + buttonSize/2.0f;
        float centerY = buttonY + buttonSize/2.0f;
        float half = iconSize/2.0f;

        if (platform.decorationPressed == titlebarButtons[i]) DrawDecorationCircle(buffer, centerX, centerY, buttonSize/2.0f, pressedColor);
        else if (platform.decorationHover == titlebarButtons[i]) DrawDecorationCircle(buffer, centerX, centerY, buttonSize/2.0f, hoverColor);

        switch (titlebarButtons[i])
        {
            case DECORATION_BUTTON_CLOSE:
            {
                DrawDecorationLine(buffer, centerX - half, centerY - half, centerX + half, centerY + half, thick, textColor);
                DrawDecorationLine(buffer, centerX + half, centerY - half, centerX - half, centerY + half, thick, textColor);
            } break;
            case DECORATION_BUTTON_MAXIMIZE:
            {
                float size = platform.maximized? half*0.8f : half;
                float offset = platform.maximized? half*0.25f : 0.0f;

                // Restore icon is drawn as two overlapped squares
                if (platform.maximized)
                {
                    DrawDecorationLine(buffer, centerX - size + 2.0f*offset, centerY - size - offset, centerX + size + 2.0f*offset, centerY - size - offset, thick, textColor);
                    DrawDecorationLine(buffer, centerX + size + 2.0f*offset, centerY - size - offset, centerX + size + 2.0f*offset, centerY + size - offset, thick, textColor);
                }

                float left = centerX - size - offset;
                float top = centerY - size + offset;
                float right = centerX + size - offset;
                float bottom = centerY + size + offset;

                DrawDecorationLine(buffer, left, top, right, top, thick, textColor);
                DrawDecorationLine(buffer, right, top, right, bottom, thick, textColor);
                DrawDecorationLine(buffer, right, bottom, left, bottom, thick, textColor);
                DrawDecorationLine(buffer, left, bottom, left, top, thick, textColor);
            } break;
            case DECORATION_BUTTON_MINIMIZE:
            {
                DrawDecorationLine(buffer, centerX - half, centerY + half*0.6f, centerX + half, centerY + half*0.6f, thick, textColor);
            } break;
            default: break;
        }

        buttonX -= (float)(DECORATION_BUTTON_SIZE + DECORATION_BUTTON_SPACING)*scale;
        buttonCount++;
    }

    // Window icon, left side
    if (platform.iconBufferCount > 0)
    {
        int windowIconSize = (int)roundf((float)DECORATION_ICON_SIZE*scale);
        DrawDecorationIcon(buffer, (int)roundf((float)DECORATION_BUTTON_SPACING*scale), (height - windowIconSize)/2, windowIconSize);
    }

    // Title text, centered and elided in the middle if too long
    platform.decorationTitleReady = true;

#if SUPPORT_MODULE_RTEXT
    Font font = GetFontDefault();
    const char *title = platform.title;

    if ((font.texture.id == 0) || (font.glyphs == NULL)) platform.decorationTitleReady = false;
    else if (title[0] != '\0')
    {
        int textScale = (int)roundf((float)DECORATION_TITLE_FONT_SCALE*scale);
        if (textScale < 1) textScale = 1;

        int maxWidth = width - 2*(int)((float)(DECORATION_BUTTON_SPACING + DECORATION_BUTTON_SIZE)*scale*(float)buttonCount) - (int)((float)DECORATION_TITLE_MARGIN*scale);
        char text[256] = { 0 };
        snprintf(text, sizeof(text), "%s", title);

        float fontSize = (float)(font.baseSize*textScale);
        float spacing = (float)textScale;

        if ((int)MeasureTextEx(font, text, fontSize, spacing).x > maxWidth)
        {
            // Middle ellipsis keeping both the start and the end of the title, cut on codepoint boundaries
            int codepointCount = GetCodepointCount(title);
            int keep = codepointCount - 1;

            while (keep > 0)
            {
                int headCount = (keep + 1)/2;
                int tailStart = codepointCount - keep/2;
                int headSize = 0;
                int tailOffset = 0;
                int offset = 0;

                for (int i = 0; i < codepointCount; i++)
                {
                    int codepointSize = 0;
                    GetCodepointNext(&title[offset], &codepointSize);
                    offset += codepointSize;

                    if ((i + 1) == headCount) headSize = offset;
                    if ((i + 1) == tailStart) tailOffset = offset;
                }

                snprintf(text, sizeof(text), "%.*s...%s", headSize, title, title + tailOffset);
                if ((int)MeasureTextEx(font, text, fontSize, spacing).x <= maxWidth) break;
                keep--;
            }
        }

        int textWidth = (int)MeasureTextEx(font, text, fontSize, spacing).x;
        int textX = (width - textWidth)/2;
        int textY = (height - font.baseSize*textScale)/2;

        if (maxWidth > 0) DrawDecorationText(buffer, font, text, textX, textY, textScale, textColor);
    }
#endif

    AttachDecorationBuffer(part, buffer, bounds, bufferScale);

    // NOTE: Subsurface state is applied on parent surface commit
    CommitWindowSurface();

    platform.decorationDirty = false;
}

// Get decoration part for a surface, -1 if none
static int GetDecorationPart(struct wl_surface *surface)
{
    if (surface == NULL) return -1;

    for (int i = 0; i < DECORATION_PART_COUNT; i++)
    {
        if (platform.decorationParts[i].surface == surface) return i;
    }

    return -1;
}

// Get titlebar button at position (logical pixels, relative to the titlebar)
static DecorationButton GetDecorationButton(double x, double y)
{
    float buttonY = (float)(DECORATION_TITLEBAR_HEIGHT - DECORATION_BUTTON_SIZE)/2.0f;
    if ((y < buttonY) || (y >= (buttonY + DECORATION_BUTTON_SIZE))) return DECORATION_BUTTON_NONE;

    float buttonX = (float)(platform.width - DECORATION_BUTTON_SPACING - DECORATION_BUTTON_SIZE);

    for (int i = 0; i < 3; i++)
    {
        if ((titlebarButtons[i] == DECORATION_BUTTON_MAXIMIZE) && !FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_RESIZABLE)) continue;

        if ((x >= buttonX) && (x < (buttonX + DECORATION_BUTTON_SIZE))) return titlebarButtons[i];
        buttonX -= (float)(DECORATION_BUTTON_SIZE + DECORATION_BUTTON_SPACING);
    }

    return DECORATION_BUTTON_NONE;
}

// Get resize edge at decorations position (logical pixels, relative to the part)
static unsigned int GetDecorationEdge(int part, double x, double y)
{
    if (!FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_RESIZABLE) || platform.maximized || platform.fullscreen) return XDG_TOPLEVEL_RESIZE_EDGE_NONE;

    const double corner = (double)DECORATION_CORNER_SIZE;
    Rectangle bounds = GetDecorationPartBounds(part);
    unsigned int edge = XDG_TOPLEVEL_RESIZE_EDGE_NONE;

    switch (part)
    {
        case DECORATION_PART_TITLEBAR:
        {
            if (y < (double)DECORATION_TOP_EDGE_SIZE)
            {
                if (x < corner) edge = XDG_TOPLEVEL_RESIZE_EDGE_TOP_LEFT;
                else if (x > (bounds.width - corner)) edge = XDG_TOPLEVEL_RESIZE_EDGE_TOP_RIGHT;
                else edge = XDG_TOPLEVEL_RESIZE_EDGE_TOP;
            }
        } break;
        case DECORATION_PART_LEFT:
        {
            if (y < corner) edge = XDG_TOPLEVEL_RESIZE_EDGE_TOP_LEFT;
            else if (y > (bounds.height - corner)) edge = XDG_TOPLEVEL_RESIZE_EDGE_BOTTOM_LEFT;
            else edge = XDG_TOPLEVEL_RESIZE_EDGE_LEFT;
        } break;
        case DECORATION_PART_RIGHT:
        {
            if (y < corner) edge = XDG_TOPLEVEL_RESIZE_EDGE_TOP_RIGHT;
            else if (y > (bounds.height - corner)) edge = XDG_TOPLEVEL_RESIZE_EDGE_BOTTOM_RIGHT;
            else edge = XDG_TOPLEVEL_RESIZE_EDGE_RIGHT;
        } break;
        case DECORATION_PART_BOTTOM:
        {
            if (x < corner) edge = XDG_TOPLEVEL_RESIZE_EDGE_BOTTOM_LEFT;
            else if (x > (bounds.width - corner)) edge = XDG_TOPLEVEL_RESIZE_EDGE_BOTTOM_RIGHT;
            else edge = XDG_TOPLEVEL_RESIZE_EDGE_BOTTOM;
        } break;
        default: break;
    }

    return edge;
}

// Process pointer button on decorations
static void HandleDecorationButton(int part, unsigned int button, bool pressed)
{
    if ((part < 0) || (platform.xdgToplevel == NULL) || (platform.seat == NULL)) return;

    unsigned int edge = GetDecorationEdge(part, platform.pointerX, platform.pointerY);

    if ((button == BTN_LEFT) && pressed)
    {
        if (edge != XDG_TOPLEVEL_RESIZE_EDGE_NONE) xdg_toplevel_resize(platform.xdgToplevel, platform.seat, platform.inputSerial, edge);
        else if (part == DECORATION_PART_TITLEBAR)
        {
            DecorationButton decorationButton = GetDecorationButton(platform.pointerX, platform.pointerY);

            if (decorationButton != DECORATION_BUTTON_NONE)
            {
                platform.decorationPressed = decorationButton;
                platform.decorationDirty = true;
            }
            else
            {
                // Double click toggles maximized state, single click moves the window
                double time = GetTime();

                if ((time - platform.decorationLastClickTime) < DECORATION_DOUBLE_CLICK_TIME)
                {
                    platform.decorationLastClickTime = 0.0;

                    if (platform.maximized) RestoreWindow();
                    else MaximizeWindow();
                }
                else
                {
                    platform.decorationLastClickTime = time;
                    xdg_toplevel_move(platform.xdgToplevel, platform.seat, platform.inputSerial);
                }
            }
        }
    }
    else if ((button == BTN_LEFT) && !pressed && (platform.decorationPressed != DECORATION_BUTTON_NONE))
    {
        DecorationButton decorationButton = (part == DECORATION_PART_TITLEBAR)? GetDecorationButton(platform.pointerX, platform.pointerY) : DECORATION_BUTTON_NONE;

        if (decorationButton == platform.decorationPressed)
        {
            switch (decorationButton)
            {
                case DECORATION_BUTTON_CLOSE: platform.closeRequested = true; break;
                case DECORATION_BUTTON_MAXIMIZE:
                {
                    if (platform.maximized) RestoreWindow();
                    else MaximizeWindow();
                } break;
                case DECORATION_BUTTON_MINIMIZE: MinimizeWindow(); break;
                default: break;
            }
        }

        platform.decorationPressed = DECORATION_BUTTON_NONE;
        platform.decorationDirty = true;
    }
    else if ((button == BTN_RIGHT) && pressed && (part == DECORATION_PART_TITLEBAR))
    {
        xdg_toplevel_show_window_menu(platform.xdgToplevel, platform.seat, platform.inputSerial, (int)platform.pointerX, (int)platform.pointerY - DECORATION_TITLEBAR_HEIGHT);
    }
}

//----------------------------------------------------------------------------------
// Module Internal Functions Definition
// NOTE: Those functions are only required for current platform
//----------------------------------------------------------------------------------

// Registry: global object announced
static void RegistryGlobalCallback(void *data, struct wl_registry *registry, uint32_t name, const char *interface, uint32_t version)
{
    if (strcmp(interface, wl_compositor_interface.name) == 0)
    {
        platform.compositorVersion = (version < 6)? version : 6;
        platform.compositor = wl_registry_bind(registry, name, &wl_compositor_interface, platform.compositorVersion);
    }
    else if (strcmp(interface, wl_subcompositor_interface.name) == 0) platform.subcompositor = wl_registry_bind(registry, name, &wl_subcompositor_interface, 1);
    else if (strcmp(interface, wl_shm_interface.name) == 0) platform.shm = wl_registry_bind(registry, name, &wl_shm_interface, 1);
    else if (strcmp(interface, xdg_wm_base_interface.name) == 0)
    {
        platform.wmBase = wl_registry_bind(registry, name, &xdg_wm_base_interface, (version < 6)? version : 6);
        xdg_wm_base_add_listener(platform.wmBase, &wmBaseListener, NULL);
    }
    else if (strcmp(interface, wl_seat_interface.name) == 0)
    {
        if (platform.seat == NULL)
        {
            platform.seatVersion = (version < 9)? version : 9;
            platform.seat = wl_registry_bind(registry, name, &wl_seat_interface, platform.seatVersion);
            wl_seat_add_listener(platform.seat, &seatListener, NULL);
        }
    }
    else if (strcmp(interface, wl_data_device_manager_interface.name) == 0)
    {
        platform.dataDeviceManager = wl_registry_bind(registry, name, &wl_data_device_manager_interface, (version < 3)? version : 3);
    }
    else if (strcmp(interface, wl_output_interface.name) == 0)
    {
        if (platform.monitorCount < MAX_MONITORS)
        {
            MonitorInfo *monitor = &platform.monitors[platform.monitorCount];
            memset(monitor, 0, sizeof(MonitorInfo));
            monitor->globalName = name;
            monitor->scale = 1;
            monitor->output = wl_registry_bind(registry, name, &wl_output_interface, (version < 4)? version : 4);
            wl_output_add_listener(monitor->output, &outputListener, NULL);

            if (platform.outputManager != NULL)
            {
                monitor->xdgOutput = zxdg_output_manager_v1_get_xdg_output(platform.outputManager, monitor->output);
                zxdg_output_v1_add_listener(monitor->xdgOutput, &xdgOutputListener, monitor->output);
            }

            platform.monitorCount++;
        }
    }
    else if (strcmp(interface, zxdg_output_manager_v1_interface.name) == 0)
    {
        platform.outputManager = wl_registry_bind(registry, name, &zxdg_output_manager_v1_interface, (version < 3)? version : 3);

        for (int i = 0; i < platform.monitorCount; i++)
        {
            if (platform.monitors[i].xdgOutput == NULL)
            {
                platform.monitors[i].xdgOutput = zxdg_output_manager_v1_get_xdg_output(platform.outputManager, platform.monitors[i].output);
                zxdg_output_v1_add_listener(platform.monitors[i].xdgOutput, &xdgOutputListener, platform.monitors[i].output);
            }
        }
    }
    else if (strcmp(interface, zxdg_decoration_manager_v1_interface.name) == 0)
    {
        platform.decorationManager = wl_registry_bind(registry, name, &zxdg_decoration_manager_v1_interface, 1);
    }
    else if (strcmp(interface, wp_fractional_scale_manager_v1_interface.name) == 0)
    {
        platform.fractionalScaleManager = wl_registry_bind(registry, name, &wp_fractional_scale_manager_v1_interface, 1);
    }
    else if (strcmp(interface, wp_viewporter_interface.name) == 0) platform.viewporter = wl_registry_bind(registry, name, &wp_viewporter_interface, 1);
    else if (strcmp(interface, wp_cursor_shape_manager_v1_interface.name) == 0)
    {
        platform.cursorShapeManager = wl_registry_bind(registry, name, &wp_cursor_shape_manager_v1_interface, 1);
    }
    else if (strcmp(interface, zwp_pointer_constraints_v1_interface.name) == 0)
    {
        platform.pointerConstraints = wl_registry_bind(registry, name, &zwp_pointer_constraints_v1_interface, 1);
    }
    else if (strcmp(interface, zwp_relative_pointer_manager_v1_interface.name) == 0)
    {
        platform.relativePointerManager = wl_registry_bind(registry, name, &zwp_relative_pointer_manager_v1_interface, 1);
    }
    else if (strcmp(interface, xdg_activation_v1_interface.name) == 0) platform.activation = wl_registry_bind(registry, name, &xdg_activation_v1_interface, 1);
    else if (strcmp(interface, xdg_toplevel_icon_manager_v1_interface.name) == 0)
    {
        platform.iconManager = wl_registry_bind(registry, name, &xdg_toplevel_icon_manager_v1_interface, 1);
    }
    else if (strcmp(interface, wp_alpha_modifier_v1_interface.name) == 0)
    {
        platform.alphaModifier = wl_registry_bind(registry, name, &wp_alpha_modifier_v1_interface, 1);
    }

    // Data device requires both the seat and the data device manager
    if ((platform.dataDevice == NULL) && (platform.seat != NULL) && (platform.dataDeviceManager != NULL))
    {
        platform.dataDevice = wl_data_device_manager_get_data_device(platform.dataDeviceManager, platform.seat);
        wl_data_device_add_listener(platform.dataDevice, &dataDeviceListener, NULL);
    }
}

// Registry: global object removed
static void RegistryGlobalRemoveCallback(void *data, struct wl_registry *registry, uint32_t name)
{
    platform.eventReceived = true;
    for (int i = 0; i < platform.monitorCount; i++)
    {
        if (platform.monitors[i].globalName == name)
        {
            struct wl_output *output = platform.monitors[i].output;

            // Remove output from the list of outputs the window is on
            for (int k = 0; k < platform.enteredOutputCount; k++)
            {
                if (platform.enteredOutputs[k] == output)
                {
                    for (int j = k; j < platform.enteredOutputCount - 1; j++) platform.enteredOutputs[j] = platform.enteredOutputs[j + 1];
                    platform.enteredOutputCount--;
                    break;
                }
            }

            if (platform.monitors[i].xdgOutput != NULL) zxdg_output_v1_destroy(platform.monitors[i].xdgOutput);
            ReleaseOutput(output);

            for (int k = i; k < platform.monitorCount - 1; k++) platform.monitors[k] = platform.monitors[k + 1];
            platform.monitorCount--;

            UpdateWindowScale();
            break;
        }
    }
}

// xdg_wm_base: compositor checks if the client is responsive
static void WmBasePingCallback(void *data, struct xdg_wm_base *wmBase, uint32_t serial)
{
    xdg_wm_base_pong(wmBase, serial);
}

// Surface: window entered an output
static void SurfaceEnterCallback(void *data, struct wl_surface *surface, struct wl_output *output)
{
    if (GetMonitorIndex(output) < 0) return;

    if (platform.enteredOutputCount < MAX_ENTERED_MONITORS) platform.enteredOutputs[platform.enteredOutputCount++] = output;

    UpdateWindowScale();
}

// Surface: window left an output
static void SurfaceLeaveCallback(void *data, struct wl_surface *surface, struct wl_output *output)
{
    for (int i = 0; i < platform.enteredOutputCount; i++)
    {
        if (platform.enteredOutputs[i] == output)
        {
            for (int k = i; k < platform.enteredOutputCount - 1; k++) platform.enteredOutputs[k] = platform.enteredOutputs[k + 1];
            platform.enteredOutputCount--;
            break;
        }
    }

    UpdateWindowScale();
}

// Surface: compositor preferred buffer scale (integer)
static void SurfacePreferredBufferScaleCallback(void *data, struct wl_surface *surface, int32_t factor)
{
    platform.preferredBufferScale = factor;
    UpdateWindowScale();
}

// Surface: compositor preferred buffer transform
static void SurfacePreferredBufferTransformCallback(void *data, struct wl_surface *surface, uint32_t transform)
{
    // Buffer transform not used, the compositor transforms the window content
}

// xdg_surface: configure sequence done, apply pending state
static void XdgSurfaceConfigureCallback(void *data, struct xdg_surface *xdgSurface, uint32_t serial)
{
    xdg_surface_ack_configure(xdgSurface, serial);
    platform.configured = true;

    ApplyConfigure(platform.pendingWidth, platform.pendingHeight, platform.pendingActivated, platform.pendingMaximized, platform.pendingFullscreen);
}

// xdg_toplevel: suggested size and state
static void XdgToplevelConfigureCallback(void *data, struct xdg_toplevel *toplevel, int32_t width, int32_t height, struct wl_array *states)
{
    platform.pendingWidth = width;
    platform.pendingHeight = height;
    platform.pendingActivated = false;
    platform.pendingMaximized = false;
    platform.pendingFullscreen = false;

    uint32_t *state = NULL;
    for (state = (uint32_t *)states->data; (const char *)state < ((const char *)states->data + states->size); state++)
    {
        switch (*state)
        {
            case XDG_TOPLEVEL_STATE_ACTIVATED: platform.pendingActivated = true; break;
            case XDG_TOPLEVEL_STATE_MAXIMIZED: platform.pendingMaximized = true; break;
            case XDG_TOPLEVEL_STATE_FULLSCREEN: platform.pendingFullscreen = true; break;
            default: break;
        }
    }
}

// xdg_toplevel: window close requested
static void XdgToplevelCloseCallback(void *data, struct xdg_toplevel *toplevel)
{
    platform.eventReceived = true;
    platform.closeRequested = true;
}

// xdg_toplevel: recommended maximum size
static void XdgToplevelConfigureBoundsCallback(void *data, struct xdg_toplevel *toplevel, int32_t width, int32_t height)
{
    // Bounds not used, the window keeps the size requested by the application
}

// xdg_toplevel: compositor capabilities
static void XdgToplevelWmCapabilitiesCallback(void *data, struct xdg_toplevel *toplevel, struct wl_array *capabilities)
{
    // Capabilities not used, all titlebar buttons are always available
}

// xdg-decoration: decoration mode selected by the compositor
static void ToplevelDecorationConfigureCallback(void *data, struct zxdg_toplevel_decoration_v1 *decoration, uint32_t mode)
{
    DecorationMode decorationMode = DECORATION_MODE_NONE;

    if (mode == ZXDG_TOPLEVEL_DECORATION_V1_MODE_SERVER_SIDE) decorationMode = DECORATION_MODE_SERVER;
    else if (platform.decorationRequested) decorationMode = DECORATION_MODE_BUILTIN;

    if (decorationMode != platform.decorationMode)
    {
        if (platform.decorationMode == DECORATION_MODE_BUILTIN) DestroyDecorations();

        platform.decorationMode = decorationMode;
        UpdateWindowGeometry();
        UpdateDecorations();
    }
}

// Fractional scale: preferred scale (x120)
static void FractionalScalePreferredScaleCallback(void *data, struct wp_fractional_scale_v1 *fractionalScale, uint32_t scale)
{
    platform.fractionalScale120 = (int)scale;
    UpdateWindowScale();
}

// Frame callback: compositor is ready for a new frame
static void FrameDoneCallback(void *data, struct wl_callback *callback, uint32_t time)
{
    wl_callback_destroy(callback);
    if (platform.frameCallback == callback) platform.frameCallback = NULL;
}

// Output: physical properties
static void OutputGeometryCallback(void *data, struct wl_output *output, int32_t x, int32_t y, int32_t physicalWidth, int32_t physicalHeight,
    int32_t subpixel, const char *make, const char *model, int32_t transform)
{
    int index = GetMonitorIndex(output);
    if (index < 0) return;

    MonitorInfo *monitor = &platform.monitors[index];
    if (monitor->xdgOutput == NULL)
    {
        monitor->x = x;
        monitor->y = y;
    }
    monitor->physicalWidth = physicalWidth;
    monitor->physicalHeight = physicalHeight;
    if (monitor->name[0] == '\0') snprintf(monitor->name, sizeof(monitor->name), "%s %s", make, model);
}

// Output: video mode
static void OutputModeCallback(void *data, struct wl_output *output, uint32_t flags, int32_t width, int32_t height, int32_t refresh)
{
    int index = GetMonitorIndex(output);
    if ((index < 0) || !FLAG_IS_SET(flags, WL_OUTPUT_MODE_CURRENT)) return;

    MonitorInfo *monitor = &platform.monitors[index];
    monitor->width = width;
    monitor->height = height;
    monitor->refreshRate = refresh;
}

// Output: all properties sent
static void OutputDoneCallback(void *data, struct wl_output *output)
{
    platform.eventReceived = true;
    int index = GetMonitorIndex(output);
    if (index < 0) return;

    MonitorInfo *monitor = &platform.monitors[index];
    if ((monitor->xdgOutput == NULL) && (monitor->scale > 0))
    {
        monitor->logicalWidth = monitor->width/monitor->scale;
        monitor->logicalHeight = monitor->height/monitor->scale;
    }
}

// Output: integer scale factor
static void OutputScaleCallback(void *data, struct wl_output *output, int32_t factor)
{
    int index = GetMonitorIndex(output);
    if (index < 0) return;

    platform.monitors[index].scale = factor;
    UpdateWindowScale();
}

// Output: name
static void OutputNameCallback(void *data, struct wl_output *output, const char *name)
{
    int index = GetMonitorIndex(output);
    if (index >= 0) snprintf(platform.monitors[index].name, sizeof(platform.monitors[index].name), "%s", name);
}

// Output: description
static void OutputDescriptionCallback(void *data, struct wl_output *output, const char *description)
{
    // Description not used, output name is used as monitor name
}

// xdg_output: logical position (compositor space)
static void XdgOutputLogicalPositionCallback(void *data, struct zxdg_output_v1 *xdgOutput, int32_t x, int32_t y)
{
    int index = GetMonitorIndex((struct wl_output *)data);
    if (index < 0) return;

    platform.monitors[index].x = x;
    platform.monitors[index].y = y;
}

// xdg_output: logical size (compositor space)
static void XdgOutputLogicalSizeCallback(void *data, struct zxdg_output_v1 *xdgOutput, int32_t width, int32_t height)
{
    int index = GetMonitorIndex((struct wl_output *)data);
    if (index < 0) return;

    platform.monitors[index].logicalWidth = width;
    platform.monitors[index].logicalHeight = height;
}

// xdg_output: all properties sent
static void XdgOutputDoneCallback(void *data, struct zxdg_output_v1 *xdgOutput)
{
    // Properties already stored
}

// xdg_output: name
static void XdgOutputNameCallback(void *data, struct zxdg_output_v1 *xdgOutput, const char *name)
{
    int index = GetMonitorIndex((struct wl_output *)data);
    if (index >= 0) snprintf(platform.monitors[index].name, sizeof(platform.monitors[index].name), "%s", name);
}

// xdg_output: description
static void XdgOutputDescriptionCallback(void *data, struct zxdg_output_v1 *xdgOutput, const char *description)
{
    // Description not used, output name is used as monitor name
}

// Activation token ready, request window activation
static void ActivationTokenDoneCallback(void *data, struct xdg_activation_token_v1 *token, const char *tokenString)
{
    if ((platform.activation != NULL) && (platform.surface != NULL)) xdg_activation_v1_activate(platform.activation, tokenString, platform.surface);
    xdg_activation_token_v1_destroy(token);
}

// Shared memory buffer released by the compositor
static void BufferReleaseCallback(void *data, struct wl_buffer *buffer)
{
    ((ShmBuffer *)data)->busy = false;
}

// libdecor: error
static void LibdecorErrorCallback(struct libdecor *context, enum libdecor_error error, const char *message)
{
    TRACELOG(LOG_WARNING, "WAYLAND: libdecor error %i: %s", error, message);
}

// libdecor: suggested size and state
static void LibdecorFrameConfigureCallback(struct libdecor_frame *frame, struct libdecor_configuration *configuration, void *data)
{
    int width = 0;
    int height = 0;
    if (!platform.libdecor.ConfigurationGetContentSize(configuration, frame, &width, &height))
    {
        width = 0;
        height = 0;
    }

    enum libdecor_window_state state = LIBDECOR_WINDOW_STATE_NONE;
    if (!platform.libdecor.ConfigurationGetWindowState(configuration, &state)) state = LIBDECOR_WINDOW_STATE_NONE;

    ApplyConfigure(width, height, FLAG_IS_SET(state, LIBDECOR_WINDOW_STATE_ACTIVE), FLAG_IS_SET(state, LIBDECOR_WINDOW_STATE_MAXIMIZED),
        FLAG_IS_SET(state, LIBDECOR_WINDOW_STATE_FULLSCREEN));

    struct libdecor_state *frameState = platform.libdecor.StateNew(platform.width, platform.height);
    platform.libdecor.FrameCommit(frame, frameState, configuration);
    platform.libdecor.StateFree(frameState);

    platform.configured = true;
}

// libdecor: window close requested
static void LibdecorFrameCloseCallback(struct libdecor_frame *frame, void *data)
{
    platform.eventReceived = true;
    platform.closeRequested = true;
}

// libdecor: window surface commit requested
static void LibdecorFrameCommitCallback(struct libdecor_frame *frame, void *data)
{
    CommitWindowSurface();
}

// libdecor: popup dismissed
static void LibdecorFrameDismissPopupCallback(struct libdecor_frame *frame, const char *seatName, void *data)
{
    // No popups used
}

// Seat: input capabilities changed
static void SeatCapabilitiesCallback(void *data, struct wl_seat *seat, uint32_t capabilities)
{
    if (FLAG_IS_SET(capabilities, WL_SEAT_CAPABILITY_POINTER) && (platform.pointer == NULL))
    {
        platform.pointer = wl_seat_get_pointer(seat);
        wl_pointer_add_listener(platform.pointer, &pointerListener, NULL);

        if (platform.cursorShapeManager != NULL) platform.cursorShapeDevice = wp_cursor_shape_manager_v1_get_pointer(platform.cursorShapeManager, platform.pointer);
    }
    else if (!FLAG_IS_SET(capabilities, WL_SEAT_CAPABILITY_POINTER) && (platform.pointer != NULL))
    {
        LockPointer(false);
        if (platform.cursorShapeDevice != NULL) wp_cursor_shape_device_v1_destroy(platform.cursorShapeDevice);
        platform.cursorShapeDevice = NULL;
        ReleasePointer(platform.pointer);
        platform.pointer = NULL;
        platform.pointerFocus = NULL;
    }

    if (FLAG_IS_SET(capabilities, WL_SEAT_CAPABILITY_KEYBOARD) && (platform.keyboard == NULL))
    {
        platform.keyboard = wl_seat_get_keyboard(seat);
        wl_keyboard_add_listener(platform.keyboard, &keyboardListener, NULL);
    }
    else if (!FLAG_IS_SET(capabilities, WL_SEAT_CAPABILITY_KEYBOARD) && (platform.keyboard != NULL))
    {
        ReleaseKeyboard(platform.keyboard);
        platform.keyboard = NULL;
        platform.repeatScancode = 0;
    }

    if (FLAG_IS_SET(capabilities, WL_SEAT_CAPABILITY_TOUCH) && (platform.touch == NULL))
    {
        platform.touch = wl_seat_get_touch(seat);
        wl_touch_add_listener(platform.touch, &touchListener, NULL);
    }
    else if (!FLAG_IS_SET(capabilities, WL_SEAT_CAPABILITY_TOUCH) && (platform.touch != NULL))
    {
        ReleaseTouch(platform.touch);
        platform.touch = NULL;
    }
}

// Seat: name
static void SeatNameCallback(void *data, struct wl_seat *seat, const char *name)
{
    // Seat name not used
}

// Pointer: entered a surface
static void PointerEnterCallback(void *data, struct wl_pointer *pointer, uint32_t serial, struct wl_surface *surface, wl_fixed_t x, wl_fixed_t y)
{
    platform.eventReceived = true;
    platform.pointerEnterSerial = serial;
    platform.pointerX = wl_fixed_to_double(x);
    platform.pointerY = wl_fixed_to_double(y);

    // Ignore surfaces not created by raylib (libdecor decorations)
    if ((surface != platform.surface) && (GetDecorationPart(surface) < 0))
    {
        platform.pointerFocus = NULL;
        return;
    }

    platform.pointerFocus = surface;

    if (surface == platform.surface)
    {
        CORE.Input.Mouse.cursorOnScreen = true;
        if (platform.lockedPointer == NULL) HandleMouseMove(platform.pointerX, platform.pointerY);
    }
    else
    {
        platform.decorationHover = DECORATION_BUTTON_NONE;
        if (GetDecorationPart(surface) == DECORATION_PART_TITLEBAR) platform.decorationHover = GetDecorationButton(platform.pointerX, platform.pointerY);
        platform.decorationDirty = true;
    }

    UpdateCursor();
}

// Pointer: left a surface
static void PointerLeaveCallback(void *data, struct wl_pointer *pointer, uint32_t serial, struct wl_surface *surface)
{
    platform.eventReceived = true;
    if (surface == platform.surface) CORE.Input.Mouse.cursorOnScreen = false;
    else if (GetDecorationPart(surface) >= 0)
    {
        platform.decorationHover = DECORATION_BUTTON_NONE;
        platform.decorationPressed = DECORATION_BUTTON_NONE;
        platform.decorationDirty = true;
    }

    platform.pointerFocus = NULL;
}

// Pointer: moved over a surface
static void PointerMotionCallback(void *data, struct wl_pointer *pointer, uint32_t time, wl_fixed_t x, wl_fixed_t y)
{
    platform.pointerX = wl_fixed_to_double(x);
    platform.pointerY = wl_fixed_to_double(y);

    if (platform.pointerFocus == NULL) return;

    if (platform.pointerFocus == platform.surface)
    {
        if (platform.lockedPointer == NULL) HandleMouseMove(platform.pointerX, platform.pointerY);
    }
    else
    {
        DecorationButton hover = DECORATION_BUTTON_NONE;
        if (GetDecorationPart(platform.pointerFocus) == DECORATION_PART_TITLEBAR) hover = GetDecorationButton(platform.pointerX, platform.pointerY);

        if (hover != platform.decorationHover)
        {
            platform.decorationHover = hover;
            platform.decorationDirty = true;
        }

        UpdateCursor();
    }
}

// Pointer: button pressed or released
static void PointerButtonCallback(void *data, struct wl_pointer *pointer, uint32_t serial, uint32_t time, uint32_t button, uint32_t state)
{
    platform.inputSerial = serial;

    if (platform.pointerFocus == NULL) return;

    bool pressed = (state == WL_POINTER_BUTTON_STATE_PRESSED);

    if (platform.pointerFocus == platform.surface)
    {
        int mouseButton = -1;

        switch (button)
        {
            case BTN_LEFT: mouseButton = MOUSE_BUTTON_LEFT; break;
            case BTN_RIGHT: mouseButton = MOUSE_BUTTON_RIGHT; break;
            case BTN_MIDDLE: mouseButton = MOUSE_BUTTON_MIDDLE; break;
            case BTN_SIDE: mouseButton = MOUSE_BUTTON_SIDE; break;
            case BTN_EXTRA: mouseButton = MOUSE_BUTTON_EXTRA; break;
            case BTN_FORWARD: mouseButton = MOUSE_BUTTON_FORWARD; break;
            case BTN_BACK: mouseButton = MOUSE_BUTTON_BACK; break;
            default: break;
        }

        if (mouseButton >= 0) HandleMouseButton(mouseButton, pressed);
    }
    else HandleDecorationButton(GetDecorationPart(platform.pointerFocus), button, pressed);
}

// Pointer: scroll (continuous value)
static void PointerAxisCallback(void *data, struct wl_pointer *pointer, uint32_t time, uint32_t axis, wl_fixed_t value)
{
    if (platform.pointerFocus != platform.surface) return;

    // NOTE: Discrete wheel steps (axis_value120) take precedence over the continuous value
    if (platform.wheelDiscrete) return;

    float move = -(float)wl_fixed_to_double(value)/10.0f;

    if (axis == WL_POINTER_AXIS_HORIZONTAL_SCROLL) platform.wheelMoveX += move;
    else platform.wheelMoveY += move;

    // NOTE: Seats older than version 5 do not send frame events
    if (platform.seatVersion < 5) PointerFrameCallback(data, pointer);
}

// Pointer: end of a group of pointer events
static void PointerFrameCallback(void *data, struct wl_pointer *pointer)
{
    platform.eventReceived = true;
    if ((platform.wheelMoveX != 0.0f) || (platform.wheelMoveY != 0.0f))
    {
        CORE.Input.Mouse.currentWheelMove.x += platform.wheelMoveX;
        CORE.Input.Mouse.currentWheelMove.y += platform.wheelMoveY;
    }

    platform.wheelMoveX = 0.0f;
    platform.wheelMoveY = 0.0f;
    platform.wheelDiscrete = false;
}

// Pointer: scroll source
static void PointerAxisSourceCallback(void *data, struct wl_pointer *pointer, uint32_t source)
{
    // Scroll source not used
}

// Pointer: scroll stopped
static void PointerAxisStopCallback(void *data, struct wl_pointer *pointer, uint32_t time, uint32_t axis)
{
    // Kinetic scrolling not used
}

// Pointer: discrete scroll steps (deprecated by axis_value120)
static void PointerAxisDiscreteCallback(void *data, struct wl_pointer *pointer, uint32_t axis, int32_t discrete)
{
    // Discrete steps from axis_value120 are used instead
}

// Pointer: high-resolution scroll steps, 120 is one wheel step
static void PointerAxisValue120Callback(void *data, struct wl_pointer *pointer, uint32_t axis, int32_t value120)
{
    if (platform.pointerFocus != platform.surface) return;

    float move = -(float)value120/120.0f;

    if (axis == WL_POINTER_AXIS_HORIZONTAL_SCROLL) platform.wheelMoveX += move;
    else platform.wheelMoveY += move;

    platform.wheelDiscrete = true;
}

// Pointer: scroll direction relative to the physical motion
static void PointerAxisRelativeDirectionCallback(void *data, struct wl_pointer *pointer, uint32_t axis, uint32_t direction)
{
    // Natural scrolling is already applied by the compositor
}

// Relative pointer: unaccelerated motion while the cursor is disabled
static void RelativePointerMotionCallback(void *data, struct zwp_relative_pointer_v1 *relativePointer, uint32_t timeHigh, uint32_t timeLow,
    wl_fixed_t dx, wl_fixed_t dy, wl_fixed_t dxUnaccel, wl_fixed_t dyUnaccel)
{
    platform.eventReceived = true;
    if (platform.lockedPointer == NULL) return;

    HandleMouseMove(CORE.Input.Mouse.currentPosition.x + wl_fixed_to_double(dxUnaccel),
        CORE.Input.Mouse.currentPosition.y + wl_fixed_to_double(dyUnaccel));
}

// Keyboard: keymap provided by the compositor
static void KeyboardKeymapCallback(void *data, struct wl_keyboard *keyboard, uint32_t format, int32_t fd, uint32_t size)
{
    if (format != WL_KEYBOARD_KEYMAP_FORMAT_XKB_V1)
    {
        close(fd);
        return;
    }

    char *keymapString = (char *)mmap(NULL, size, PROT_READ, MAP_PRIVATE, fd, 0);
    if (keymapString == MAP_FAILED)
    {
        close(fd);
        return;
    }

    struct xkb_keymap *keymap = xkb_keymap_new_from_string(platform.xkbContext, keymapString, XKB_KEYMAP_FORMAT_TEXT_V1, XKB_KEYMAP_COMPILE_NO_FLAGS);
    munmap(keymapString, size);
    close(fd);

    if (keymap == NULL)
    {
        TRACELOG(LOG_WARNING, "WAYLAND: Failed to compile keyboard keymap");
        return;
    }

    if (platform.xkbState != NULL) xkb_state_unref(platform.xkbState);
    if (platform.xkbKeymap != NULL) xkb_keymap_unref(platform.xkbKeymap);

    platform.xkbKeymap = keymap;
    platform.xkbState = xkb_state_new(keymap);
}

// Keyboard: window got keyboard focus
static void KeyboardEnterCallback(void *data, struct wl_keyboard *keyboard, uint32_t serial, struct wl_surface *surface, struct wl_array *keys)
{
    platform.eventReceived = true;
    platform.inputSerial = serial;
}

// Keyboard: window lost keyboard focus
static void KeyboardLeaveCallback(void *data, struct wl_keyboard *keyboard, uint32_t serial, struct wl_surface *surface)
{
    platform.eventReceived = true;
    platform.repeatScancode = 0;

    // Release keys and mouse buttons, their release events are sent to the newly focused window
    for (int i = 0; i < MAX_KEYBOARD_KEYS; i++)
    {
        if ((i != KEY_CAPS_LOCK) && (i != KEY_NUM_LOCK)) CORE.Input.Keyboard.currentKeyState[i] = 0;
    }

    for (int i = 0; i < MAX_MOUSE_BUTTONS; i++) CORE.Input.Mouse.currentButtonState[i] = 0;

    if (platform.xkbComposeState != NULL) xkb_compose_state_reset(platform.xkbComposeState);
}

// Keyboard: key pressed or released
static void KeyboardKeyCallback(void *data, struct wl_keyboard *keyboard, uint32_t serial, uint32_t time, uint32_t key, uint32_t state)
{
    platform.inputSerial = serial;

    bool pressed = (state == WL_KEYBOARD_KEY_STATE_PRESSED);
    HandleKey(key, pressed);

    if (pressed)
    {
        HandleKeyChar(key);

        // Start key repeat
        if ((platform.repeatRate > 0) && (platform.xkbKeymap != NULL) && xkb_keymap_key_repeats(platform.xkbKeymap, key + XKB_KEYCODE_OFFSET))
        {
            platform.repeatScancode = key;
            platform.repeatTime = GetTime() + (double)platform.repeatDelay/1000.0;
        }
    }
    else if (platform.repeatScancode == key) platform.repeatScancode = 0;
}

// Keyboard: modifiers state changed
static void KeyboardModifiersCallback(void *data, struct wl_keyboard *keyboard, uint32_t serial, uint32_t depressed, uint32_t latched, uint32_t locked, uint32_t group)
{
    if (platform.xkbState == NULL) return;

    xkb_state_update_mask(platform.xkbState, depressed, latched, locked, 0, 0, group);

    // NOTE: CAPS/NUM lock keys are considered down while the lock is active
    CORE.Input.Keyboard.currentKeyState[KEY_CAPS_LOCK] = (xkb_state_mod_name_is_active(platform.xkbState, XKB_MOD_NAME_CAPS, XKB_STATE_MODS_LOCKED) > 0);
    CORE.Input.Keyboard.currentKeyState[KEY_NUM_LOCK] = (xkb_state_mod_name_is_active(platform.xkbState, XKB_MOD_NAME_NUM, XKB_STATE_MODS_LOCKED) > 0);
}

// Keyboard: key repeat rate and delay
static void KeyboardRepeatInfoCallback(void *data, struct wl_keyboard *keyboard, int32_t rate, int32_t delay)
{
    platform.repeatRate = rate;
    platform.repeatDelay = delay;
}

// Touch: new touch point
static void TouchDownCallback(void *data, struct wl_touch *touch, uint32_t serial, uint32_t time, struct wl_surface *surface, int32_t id, wl_fixed_t x, wl_fixed_t y)
{
    platform.inputSerial = serial;

    if (surface != platform.surface) return;

    for (int i = 0; i < MAX_TOUCH_POINTS; i++)
    {
        if (platform.touchIds[i] == -1)
        {
            HandleTouch(TOUCH_ACTION_DOWN, i, id, wl_fixed_to_double(x), wl_fixed_to_double(y));
            break;
        }
    }
}

// Touch: touch point released
static void TouchUpCallback(void *data, struct wl_touch *touch, uint32_t serial, uint32_t time, int32_t id)
{
    platform.inputSerial = serial;

    for (int i = 0; i < MAX_TOUCH_POINTS; i++)
    {
        if (platform.touchIds[i] == id)
        {
            HandleTouch(TOUCH_ACTION_UP, i, id, platform.touchPositions[i].x, platform.touchPositions[i].y);
            break;
        }
    }
}

// Touch: touch point moved
static void TouchMotionCallback(void *data, struct wl_touch *touch, uint32_t time, int32_t id, wl_fixed_t x, wl_fixed_t y)
{
    for (int i = 0; i < MAX_TOUCH_POINTS; i++)
    {
        if (platform.touchIds[i] == id)
        {
            HandleTouch(TOUCH_ACTION_MOVE, i, id, wl_fixed_to_double(x), wl_fixed_to_double(y));
            break;
        }
    }
}

// Touch: end of a group of touch events
static void TouchFrameCallback(void *data, struct wl_touch *touch)
{
    // Touch events are processed individually
}

// Touch: touch session cancelled by the compositor
static void TouchCancelCallback(void *data, struct wl_touch *touch)
{
    for (int i = 0; i < MAX_TOUCH_POINTS; i++)
    {
        if (platform.touchIds[i] != -1) HandleTouch(TOUCH_ACTION_CANCEL, i, platform.touchIds[i], platform.touchPositions[i].x, platform.touchPositions[i].y);
    }
}

// Touch: touch point shape
static void TouchShapeCallback(void *data, struct wl_touch *touch, int32_t id, wl_fixed_t major, wl_fixed_t minor)
{
    // Touch shape not used
}

// Touch: touch point orientation
static void TouchOrientationCallback(void *data, struct wl_touch *touch, int32_t id, wl_fixed_t orientation)
{
    // Touch orientation not used
}

// Data offer: mime type offered
static void DataOfferOfferCallback(void *data, struct wl_data_offer *offer, const char *mimeType)
{
    DataOffer *dataOffer = (DataOffer *)data;

    if ((dataOffer != NULL) && (dataOffer->mimeTypeCount < MAX_DATA_OFFER_MIME_TYPES))
    {
        int length = (int)strlen(mimeType);
        dataOffer->mimeTypes[dataOffer->mimeTypeCount] = (char *)RL_CALLOC(length + 1, sizeof(char));
        memcpy(dataOffer->mimeTypes[dataOffer->mimeTypeCount], mimeType, length);
        dataOffer->mimeTypeCount++;
    }
}

// Data offer: drag and drop actions supported by the source
static void DataOfferSourceActionsCallback(void *data, struct wl_data_offer *offer, uint32_t actions)
{
    // Only copy action is used
}

// Data offer: drag and drop action selected by the compositor
static void DataOfferActionCallback(void *data, struct wl_data_offer *offer, uint32_t action)
{
    // Only copy action is used
}

// Data device: new data offer, followed by its mime types
static void DataDeviceDataOfferCallback(void *data, struct wl_data_device *device, struct wl_data_offer *offer)
{
    DataOffer *dataOffer = (DataOffer *)RL_CALLOC(1, sizeof(DataOffer));
    dataOffer->offer = offer;
    wl_data_offer_add_listener(offer, &dataOfferListener, dataOffer);
}

// Data device: drag and drop entered a surface
static void DataDeviceEnterCallback(void *data, struct wl_data_device *device, uint32_t serial, struct wl_surface *surface,
    wl_fixed_t x, wl_fixed_t y, struct wl_data_offer *offer)
{
    DestroyDataOffer(platform.dragOffer);
    platform.dragOffer = NULL;

    if (offer == NULL) return;

    platform.dragOffer = (DataOffer *)wl_data_offer_get_user_data(offer);

    if ((surface == platform.surface) && DataOfferHasMimeType(platform.dragOffer, "text/uri-list"))
    {
        wl_data_offer_accept(offer, serial, "text/uri-list");
        if (wl_data_offer_get_version(offer) >= 3) wl_data_offer_set_actions(offer, WL_DATA_DEVICE_MANAGER_DND_ACTION_COPY, WL_DATA_DEVICE_MANAGER_DND_ACTION_COPY);
    }
    else wl_data_offer_accept(offer, serial, NULL);
}

// Data device: drag and drop left the surface
static void DataDeviceLeaveCallback(void *data, struct wl_data_device *device)
{
    DestroyDataOffer(platform.dragOffer);
    platform.dragOffer = NULL;
}

// Data device: drag and drop moved over the surface
static void DataDeviceMotionCallback(void *data, struct wl_data_device *device, uint32_t time, wl_fixed_t x, wl_fixed_t y)
{
    // Drop position not used
}

// Data device: drag and drop dropped on the surface
static void DataDeviceDropCallback(void *data, struct wl_data_device *device)
{
    DataOffer *offer = platform.dragOffer;
    if (offer == NULL) return;

    if (DataOfferHasMimeType(offer, "text/uri-list"))
    {
        char *uris = ReadDataOffer(offer, "text/uri-list", NULL);
        if (uris != NULL)
        {
            ProcessDroppedUris(uris);
            RL_FREE(uris);
        }

        if (wl_data_offer_get_version(offer->offer) >= 3) wl_data_offer_finish(offer->offer);
    }

    DestroyDataOffer(platform.dragOffer);
    platform.dragOffer = NULL;
}

// Data device: clipboard selection changed
static void DataDeviceSelectionCallback(void *data, struct wl_data_device *device, struct wl_data_offer *offer)
{
    DestroyDataOffer(platform.selectionOffer);
    platform.selectionOffer = (offer != NULL)? (DataOffer *)wl_data_offer_get_user_data(offer) : NULL;
}

// Data source: target mime type accepted (drag and drop)
static void DataSourceTargetCallback(void *data, struct wl_data_source *source, const char *mimeType)
{
    // Drag and drop source not used
}

// Data source: clipboard content requested by another client
static void DataSourceSendCallback(void *data, struct wl_data_source *source, const char *mimeType, int32_t fd)
{
    if (platform.clipboardText != NULL)
    {
        const char *text = platform.clipboardText;
        size_t length = strlen(text);

        while (length > 0)
        {
            ssize_t written = write(fd, text, length);

            if (written < 0)
            {
                if (errno == EINTR) continue;
                break;
            }

            text += written;
            length -= written;
        }
    }

    close(fd);
}

// Data source: clipboard ownership lost
static void DataSourceCancelledCallback(void *data, struct wl_data_source *source)
{
    if (platform.dataSource == source) platform.dataSource = NULL;
    wl_data_source_destroy(source);
}

// Data source: drag and drop drop performed
static void DataSourceDndDropPerformedCallback(void *data, struct wl_data_source *source)
{
    // Drag and drop source not used
}

// Data source: drag and drop finished
static void DataSourceDndFinishedCallback(void *data, struct wl_data_source *source)
{
    // Drag and drop source not used
}

// Data source: drag and drop action selected
static void DataSourceActionCallback(void *data, struct wl_data_source *source, uint32_t action)
{
    // Drag and drop source not used
}

// EOF
