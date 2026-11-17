/**********************************************************************************************
*
*   rcore_desktop_win32 - Functions to manage window, graphics device and inputs
*
*   PLATFORM: DESKTOP: WIN32
*       - Windows (Win32, Win64)
*
*   LIMITATIONS:
*       - Initial development stage, lot of functionality missing
*       - No support for MOUSE_BUTTON_FORWARD/MOUSE_BUTTON_BACK
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
*       - Win32 API (windows.h)
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

// Move windows.h symbols to new names to avoid redefining the same names as raylib
#define CloseWindow CloseWindowWin32
#define Rectangle RectangleWin32
#define ShowCursor ShowCursorWin32
#define DrawTextA DrawTextAWin32
#define DrawTextW DrawTextWin32
#define DrawTextExA DrawTextExAWin32
#define DrawTextExW DrawTextExWin32

#define WIN32_LEAN_AND_MEAN
#include <windows.h>

#undef CloseWindow      // raylib symbol collision
#undef Rectangle        // raylib symbol collision
#undef ShowCursor       // raylib symbol collision
#undef LoadImage        // raylib symbol collision
#undef DrawText         // raylib symbol collision
#undef DrawTextA
#undef DrawTextW
#undef DrawTextEx       // raylib symbol collision
#undef DrawTextExA
#undef DrawTextExW

#include <windowsx.h>
#include <imm.h>            // Required for: COMPOSITIONFORM, CANDIDATEFORM (functions loaded at runtime)
#include <shellapi.h>       // Required for: DragAcceptFiles(), DragQueryFileW(), DragFinish()

// Avoid including shellscalingapi.h by defining these two enums directly
#ifndef DPI_ENUMS_DECLARED
typedef enum PROCESS_DPI_AWARENESS {
    PROCESS_DPI_UNAWARE = 0,
    PROCESS_SYSTEM_DPI_AWARE = 1,
    PROCESS_PER_MONITOR_DPI_AWARE = 2
} PROCESS_DPI_AWARENESS;

typedef enum MONITOR_DPI_TYPE {
    MDT_EFFECTIVE_DPI = 0,
    MDT_ANGULAR_DPI = 1,
    MDT_RAW_DPI = 2,
    MDT_DEFAULT = MDT_EFFECTIVE_DPI
} MONITOR_DPI_TYPE;

#define DPI_ENUMS_DECLARED
#endif

#include <malloc.h>          // Required for alloca()

#if !defined(GRAPHICS_API_OPENGL_SOFTWARE)
    #include <GL/gl.h>
#endif

//----------------------------------------------------------------------------------
// Types and Structures Definition
//----------------------------------------------------------------------------------

// NOTE: appScreenSize is the last screen size requested by the app,
// the backend must keep the client area this size (after DPI scaling is applied)
// when the window isn't fullscreen/maximized/minimized
typedef struct {
    HWND hwnd;              // Window handler
    HDC hdc;                // Graphic context handler
    HGLRC glContext;        // OpenGL context handler

    HICON largeIcon;        // Large icon typically used for the taskbar
    HICON smallIcon;        // Small icon typically used for the window's title bar

    // Software renderer variables
    HDC hdcmem;             // Memory graphic context handler
    HBITMAP hbitmap;        // GDI bitmap handler
    unsigned int *pixels;   // Pointer to pixel data buffer (BGRA format)

    int previousRenderWidth;
    int previousRenderHeight;
    int appScreenWidth;
    int appScreenHeight;
    unsigned int desiredFlags;

    RECT coverRestoreRect;          // Window rectangle to restore when leaving fullscreen or borderless windowed mode
    bool coverRestoreUndecorated;   // FLAG_WINDOW_UNDECORATED state to restore when leaving fullscreen or borderless windowed mode
    bool coverRestoreMaximized;     // FLAG_WINDOW_MAXIMIZED state to restore when leaving fullscreen or borderless windowed mode
    unsigned int coverEnterFlags;   // Fullscreen or borderless windowed mode flag being applied by a window style change
    bool coverLeave;                // Fullscreen or borderless windowed mode being left by a window style change
    float opacity;                  // Window opacity

    HCURSOR cursor;                 // Mouse cursor shape
    bool mouseTracked;              // Mouse leave notification requested
    POINT cursorRestorePosition;    // Cursor screen position to restore when the cursor is enabled
    POINT rawAbsolutePosition;      // Last position of an absolute pointing device (remote desktop, tablets)
    bool rawAbsoluteValid;          // Absolute pointing device position is known
    POINT textInputPosition;        // Last clicked position, input method windows are placed next to it
    WCHAR highSurrogate;            // Pending first half of a character outside of the basic multilingual plane

    int touchIds[MAX_TOUCH_POINTS];                 // Touch points ids, -1 if not used
    Vector2 touchPositions[MAX_TOUCH_POINTS];       // Touch points positions

    double gamepadCheckTime;                        // Next time to check for connected gamepads
    double gamepadVibrationEnd[MAX_GAMEPADS];       // Time to stop gamepads vibration, 0 if not vibrating

    bool keyPressedInFrame[MAX_KEYBOARD_KEYS];      // Keys pressed in current frame
    bool keyReleasePending[MAX_KEYBOARD_KEYS];      // Keys released in the frame they were pressed, release applied on next frame
    bool buttonPressedInFrame[MAX_MOUSE_BUTTONS];   // Mouse buttons pressed in current frame
    bool buttonReleasePending[MAX_MOUSE_BUTTONS];   // Mouse buttons released in the frame they were pressed, release applied on next frame

    bool perMonitorDPIAwareV2; // Flag to check DPI-awareness level. This is used mainly to distinguish between per-monitor DPI awareness v1 and v2

    LARGE_INTEGER timerFrequency;

    char *clipboardText;
} PlatformData;

// Define WGL function pointer types (no wglext.h needed)
typedef HGLRC (WINAPI *PFNWGLCREATECONTEXTATTRIBSARBPROC)(HDC, HGLRC, const int *);
typedef BOOL (WINAPI *PFNWGLCHOOSEPIXELFORMATARBPROC)(HDC, const int *, const FLOAT *, UINT, int *, UINT *);
typedef BOOL (WINAPI *PFNWGLSWAPINTERVALEXTPROC)(int);
typedef const char *(WINAPI *PFNWGLGETEXTENSIONSSTRINGARBPROC)(HDC hdc);

// Define Win32 function pointer types
typedef UINT (WINAPI *GetDpiForWindowProc)(HWND hwnd);
typedef BOOL (WINAPI *AdjustWindowRectExForDpiProc)(LPRECT lpRect, DWORD dwStyle, BOOL bMenu, DWORD dwExStyle, UINT dpi);
typedef BOOL (WINAPI *SetProcessDpiAwarenessContextProc)(HANDLE value);
typedef HRESULT (WINAPI *SetProcessDpiAwarenessProc)(PROCESS_DPI_AWARENESS value);
typedef BOOL (WINAPI *SetProcessDPIAwareProc)(void);
typedef HRESULT (WINAPI *GetDpiForMonitorProc)(HMONITOR monitor, MONITOR_DPI_TYPE dpiType, UINT *dpiX, UINT *dpiY);
typedef HANDLE (WINAPI *GetThreadDpiAwarenessContextProc)(void);
typedef BOOL (WINAPI *AreDpiAwarenessContextsEqualProc)(HANDLE a, HANDLE b);
typedef BOOL (WINAPI *GetPointerTypeProc)(UINT32 pointerId, DWORD *pointerType);
typedef HIMC (WINAPI *ImmGetContextProc)(HWND hwnd);
typedef BOOL (WINAPI *ImmReleaseContextProc)(HWND hwnd, HIMC context);
typedef BOOL (WINAPI *ImmSetCompositionWindowProc)(HIMC context, LPCOMPOSITIONFORM form);
typedef BOOL (WINAPI *ImmSetCandidateWindowProc)(HIMC context, LPCANDIDATEFORM form);

// XInput gamepad state and functions, library loaded at runtime (no xinput.h needed)
typedef struct {
    WORD buttons;
    BYTE leftTrigger;
    BYTE rightTrigger;
    SHORT thumbLX;
    SHORT thumbLY;
    SHORT thumbRX;
    SHORT thumbRY;
} XInputGamepad;

typedef struct {
    DWORD packetNumber;
    XInputGamepad gamepad;
} XInputState;

typedef struct {
    WORD leftMotorSpeed;
    WORD rightMotorSpeed;
} XInputVibration;

typedef DWORD (WINAPI *XInputGetStateProc)(DWORD userIndex, XInputState *state);
typedef DWORD (WINAPI *XInputSetStateProc)(DWORD userIndex, XInputVibration *vibration);

//----------------------------------------------------------------------------------
// Global Variables Definition
//----------------------------------------------------------------------------------
extern CoreData CORE;                   // Global CORE state context

static PlatformData platform = { 0 };   // Platform specific data

// Required WGL functions
static PFNWGLCREATECONTEXTATTRIBSARBPROC wglCreateContextAttribsARB = NULL;
static PFNWGLCHOOSEPIXELFORMATARBPROC wglChoosePixelFormatARB = NULL;
static PFNWGLSWAPINTERVALEXTPROC wglSwapIntervalEXT = NULL;
static PFNWGLGETEXTENSIONSSTRINGARBPROC wglGetExtensionsStringARB = NULL;

// Win32 functions that need to be dynamically loaded to support older versions of Windows
// Before calling any of these functions, they should be null-checked
static GetDpiForWindowProc pGetDpiForWindow = NULL;
static AdjustWindowRectExForDpiProc pAdjustWindowRectExForDpi = NULL;
static SetProcessDpiAwarenessContextProc pSetProcessDpiAwarenessContext = NULL;
static SetProcessDpiAwarenessProc pSetProcessDpiAwareness = NULL;
static SetProcessDPIAwareProc pSetProcessDPIAware = NULL;
static GetDpiForMonitorProc pGetDpiForMonitor = NULL;
static GetThreadDpiAwarenessContextProc pGetThreadDpiAwarenessContext = NULL;
static AreDpiAwarenessContextsEqualProc pAreDpiAwarenessContextsEqual = NULL;

// Touch input function (Windows 8)
static GetPointerTypeProc pGetPointerType = NULL;

// Gamepad functions, XInput library is loaded at runtime
static XInputGetStateProc pXInputGetState = NULL;
static XInputSetStateProc pXInputSetState = NULL;

static HMODULE xinput = NULL; // xinput1_4.dll, xinput1_3.dll or xinput9_1_0.dll: gamepads

// Input method functions, imm32.dll is loaded at runtime
static ImmGetContextProc pImmGetContext = NULL;
static ImmReleaseContextProc pImmReleaseContext = NULL;
static ImmSetCompositionWindowProc pImmSetCompositionWindow = NULL;
static ImmSetCandidateWindowProc pImmSetCandidateWindow = NULL;

static HMODULE imm32 = NULL;  // imm32.dll: input methods
static HMODULE shcore = NULL; // shcore.dll (Windows 8.1 and early versions of Windows 10): loaded once, NULL on other versions

// --------------------------------------------------------------------------------
// This part of the file contains pure functions that never access global state
// This distinction helps keep the backend maintainable as the inputs and outputs
// of every function called in this section can be fully derived from the call-site alone
// --------------------------------------------------------------------------------

// Prevent any code in this part of the file from accessing the global CORE state
#define CORE DONT_USE_CORE_HERE

//----------------------------------------------------------------------------------
// Defines and Macros
//----------------------------------------------------------------------------------
#define A_TO_W_ALLOCA(outWstr, inAnsi)   do {                   \
    size_t outLen = AToWLen(inAnsi);                            \
    outWstr = (WCHAR *)alloca(sizeof(WCHAR)*(outLen + 1));      \
    AToWCopy(inAnsi, outWstr, outLen);                          \
    outWstr[outLen] = 0;                                        \
} while (0)

#define STYLE_MASK_ALL          0xffffffff
#define STYLE_MASK_READONLY     (WS_MINIMIZE | WS_MAXIMIZE)
#define STYLE_MASK_WRITABLE     (~STYLE_MASK_READONLY)

#define STYLE_FLAGS_RESIZABLE   (WS_THICKFRAME | WS_MAXIMIZEBOX)

#define STYLE_FLAGS_UNDECORATED_OFF     (WS_OVERLAPPED | WS_CAPTION | WS_SYSMENU | WS_MINIMIZEBOX)
#define STYLE_FLAGS_UNDECORATED_ON      WS_POPUP

#define WINDOW_STYLE_EX         0

#define CLASS_NAME              L"raylibWindow"

#define FLAG_MASK_OPTIONAL      (FLAG_VSYNC_HINT | FLAG_WINDOW_UNFOCUSED)
#define FLAG_MASK_REQUIRED      ~(FLAG_MASK_OPTIONAL)

// Flags that have no operations to perform during an update
#define FLAG_MASK_NO_UPDATE     (FLAG_WINDOW_HIGHDPI | FLAG_MSAA_4X_HINT | FLAG_WINDOW_TRANSPARENT | FLAG_INTERLACED_HINT)

// Flags of the modes covering the monitor with an undecorated window, there is no video mode change
#define FLAG_MASK_COVER         (FLAG_FULLSCREEN_MODE | FLAG_BORDERLESS_WINDOWED_MODE)

#define MAX_MONITORS            16

// Lowest extended-key scancode, start of keyboardKeyFromExtendedScancode[]
#define EXTENDED_KEY_OFFSET     0x1C

#define WGL_DRAW_TO_WINDOW_ARB              0x2001
#define WGL_ACCELERATION_ARB                0x2003
#define WGL_SUPPORT_OPENGL_ARB              0x2010
#define WGL_DOUBLE_BUFFER_ARB               0x2011
#define WGL_PIXEL_TYPE_ARB                  0x2013
#define WGL_COLOR_BITS_ARB                  0x2014
#define WGL_RED_BITS_ARB                    0x2015
#define WGL_RED_SHIFT_ARB                   0x2016
#define WGL_GREEN_BITS_ARB                  0x2017
#define WGL_GREEN_SHIFT_ARB                 0x2018
#define WGL_BLUE_BITS_ARB                   0x2019
#define WGL_BLUE_SHIFT_ARB                  0x201a
#define WGL_ALPHA_BITS_ARB                  0x201b
#define WGL_ALPHA_SHIFT_ARB                 0x201c
#define WGL_DEPTH_BITS_ARB                  0x2022
#define WGL_STENCIL_BITS_ARB                0x2023
#define WGL_TYPE_RGBA_ARB                   0x202b

// Context acceleration types
#define WGL_NO_ACCELERATION_ARB             0x2025      // OpenGL 1.1 GDI software rasterizer
#define WGL_GENERIC_ACCELERATION_ARB        0x2026
#define WGL_FULL_ACCELERATION_ARB           0x2027      // OpenGL hardware-accelerated, using GPU-drivers provided by vendor

// WGL_ARB_multisample extension supported
#define WGL_SAMPLE_BUFFERS_ARB              0x2041      // Multisampling: 1 if multisample buffers are supported
#define WGL_SAMPLES_ARB                     0x2042      // Multisampling: Number of samples per pixel (4, 8, 16)

// WGL_ARB_framebuffer_sRGB extension supported
#define WGL_FRAMEBUFFER_SRGB_CAPABLE_ARB    0x20a9      // GL_TRUE if the framebuffer can do sRGB conversion

#define WGL_NUMBER_PIXEL_FORMATS_ARB        0x2000
#define WGL_CONTEXT_MAJOR_VERSION_ARB       0x2091
#define WGL_CONTEXT_MINOR_VERSION_ARB       0x2092
#define WGL_CONTEXT_PROFILE_MASK_ARB        0x9126
#define WGL_CONTEXT_CORE_PROFILE_BIT_ARB    0x00000001
#define WGL_CONTEXT_COMPATIBILITY_PROFILE_BIT_ARB 0x00000002
#define WGL_CONTEXT_ES_PROFILE_BIT_EXT        0x00000004
#define WGL_CONTEXT_ES2_PROFILE_BIT_EXT        0x00000004

// Polyfill macro definitions and enums for older Windows SDKs
#ifndef WM_DPICHANGED
    #define WM_DPICHANGED                   0x02E0
#endif

#ifndef WM_GETDPISCALEDSIZE
    #define WM_GETDPISCALEDSIZE             0x02E4
#endif

#ifndef WM_MOUSEHWHEEL
    #define WM_MOUSEHWHEEL                  0x020E
#endif

#ifndef WM_POINTERUPDATE
    #define WM_POINTERUPDATE                0x0245
    #define WM_POINTERDOWN                  0x0246
    #define WM_POINTERUP                    0x0247
#endif
#ifndef WM_POINTERCAPTURECHANGED
    #define WM_POINTERCAPTURECHANGED        0x024C
#endif

#define POINTER_TYPE_TOUCH                  2           // Pointer input type: touch (PT_TOUCH)

#define XINPUT_MAX_GAMEPADS                 4           // Maximum number of gamepads supported by XInput
#define XINPUT_GAMEPAD_NAME                 "Xbox Controller"

#ifndef DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE
    #define DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE ((HANDLE)-3)
#endif
#ifndef DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2
    #define DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2 ((HANDLE)-4)
#endif

//----------------------------------------------------------------------------------
// Types and Structures Definition
//----------------------------------------------------------------------------------
// Maximize-minimize request types
typedef enum {
    MIZED_NONE,
    MIZED_MIN,
    MIZED_MAX
} Mized;

// Flag operations
// NOTE: Some ops need to be deferred
typedef struct {
    DWORD set;
    DWORD clear;
} FlagsOp;

// Monitors list, primary monitor first
typedef struct {
    HMONITOR handles[MAX_MONITORS];
    int count;
} MonitorList;

//----------------------------------------------------------------------------------
// Global Const Variables Definition
//----------------------------------------------------------------------------------
// Maps scancodes to Raylib keys
static const unsigned short keyboardKeyFromScancode[] = {
    [0x01] = KEY_ESCAPE,
    [0x02] = KEY_ONE,
    [0x03] = KEY_TWO,
    [0x04] = KEY_THREE,
    [0x05] = KEY_FOUR,
    [0x06] = KEY_FIVE,
    [0x07] = KEY_SIX,
    [0x08] = KEY_SEVEN,
    [0x09] = KEY_EIGHT,
    [0x0A] = KEY_NINE,
    [0x0B] = KEY_ZERO,
    [0x0C] = KEY_MINUS,
    [0x0D] = KEY_EQUAL,
    [0x0E] = KEY_BACKSPACE,
    [0x0F] = KEY_TAB,
    [0x10] = KEY_Q,
    [0x11] = KEY_W,
    [0x12] = KEY_E,
    [0x13] = KEY_R,
    [0x14] = KEY_T,
    [0x15] = KEY_Y,
    [0x16] = KEY_U,
    [0x17] = KEY_I,
    [0x18] = KEY_O,
    [0x19] = KEY_P,
    [0x1A] = KEY_LEFT_BRACKET,
    [0x1B] = KEY_RIGHT_BRACKET,
    [0x1C] = KEY_ENTER,
    [0x1D] = KEY_LEFT_CONTROL,
    [0x1E] = KEY_A,
    [0x1F] = KEY_S,
    [0x20] = KEY_D,
    [0x21] = KEY_F,
    [0x22] = KEY_G,
    [0x23] = KEY_H,
    [0x24] = KEY_J,
    [0x25] = KEY_K,
    [0x26] = KEY_L,
    [0x27] = KEY_SEMICOLON,
    [0x28] = KEY_APOSTROPHE,
    [0x29] = KEY_GRAVE,
    [0x2A] = KEY_LEFT_SHIFT,
    [0x2B] = KEY_BACKSLASH,
    [0x2C] = KEY_Z,
    [0x2D] = KEY_X,
    [0x2E] = KEY_C,
    [0x2F] = KEY_V,
    [0x30] = KEY_B,
    [0x31] = KEY_N,
    [0x32] = KEY_M,
    [0x33] = KEY_COMMA,
    [0x34] = KEY_PERIOD,
    [0x35] = KEY_SLASH,
    [0x36] = KEY_RIGHT_SHIFT,
    [0x37] = KEY_KP_MULTIPLY,
    [0x38] = KEY_LEFT_ALT,
    [0x39] = KEY_SPACE,
    [0x3A] = KEY_CAPS_LOCK,
    [0x3B] = KEY_F1,
    [0x3C] = KEY_F2,
    [0x3D] = KEY_F3,
    [0x3E] = KEY_F4,
    [0x3F] = KEY_F5,
    [0x40] = KEY_F6,
    [0x41] = KEY_F7,
    [0x42] = KEY_F8,
    [0x43] = KEY_F9,
    [0x44] = KEY_F10,
    [0x45] = KEY_PAUSE,
    [0x46] = KEY_SCROLL_LOCK,
    [0x47] = KEY_KP_7,
    [0x48] = KEY_KP_8,
    [0x49] = KEY_KP_9,
    [0x4A] = KEY_KP_SUBTRACT,
    [0x4B] = KEY_KP_4,
    [0x4C] = KEY_KP_5,
    [0x4D] = KEY_KP_6,
    [0x4E] = KEY_KP_ADD,
    [0x4F] = KEY_KP_1,
    [0x50] = KEY_KP_2,
    [0x51] = KEY_KP_3,
    [0x52] = KEY_KP_0,
    [0x53] = KEY_KP_DECIMAL,
    [0x57] = KEY_F11,
    [0x58] = KEY_F12,
    [0x59] = KEY_KP_EQUAL,

    // Duplicate scancodes
    [0x54] = KEY_PRINT_SCREEN,      // Alt+PrtSc sends the SysRq scancode
};

// Maps scancodes with extended-flag set to Raylib keys; lowest key is 0x1C
static const unsigned short keyboardKeyFromExtendedScancode[] = {
    [0x1C - EXTENDED_KEY_OFFSET] = KEY_KP_ENTER,
    [0x1D - EXTENDED_KEY_OFFSET] = KEY_RIGHT_CONTROL,
    [0x35 - EXTENDED_KEY_OFFSET] = KEY_KP_DIVIDE,
    [0x37 - EXTENDED_KEY_OFFSET] = KEY_PRINT_SCREEN,
    [0x38 - EXTENDED_KEY_OFFSET] = KEY_RIGHT_ALT,
    [0x45 - EXTENDED_KEY_OFFSET] = KEY_NUM_LOCK,
    [0x47 - EXTENDED_KEY_OFFSET] = KEY_HOME,
    [0x48 - EXTENDED_KEY_OFFSET] = KEY_UP,
    [0x49 - EXTENDED_KEY_OFFSET] = KEY_PAGE_UP,
    [0x4B - EXTENDED_KEY_OFFSET] = KEY_LEFT,
    [0x4D - EXTENDED_KEY_OFFSET] = KEY_RIGHT,
    [0x4F - EXTENDED_KEY_OFFSET] = KEY_END,
    [0x50 - EXTENDED_KEY_OFFSET] = KEY_DOWN,
    [0x51 - EXTENDED_KEY_OFFSET] = KEY_PAGE_DOWN,
    [0x52 - EXTENDED_KEY_OFFSET] = KEY_INSERT,
    [0x53 - EXTENDED_KEY_OFFSET] = KEY_DELETE,
    [0x5B - EXTENDED_KEY_OFFSET] = KEY_LEFT_SUPER,
    [0x5C - EXTENDED_KEY_OFFSET] = KEY_RIGHT_SUPER,
    [0x5D - EXTENDED_KEY_OFFSET] = KEY_KB_MENU,

    // Duplicate scancodes
    [0x46 - EXTENDED_KEY_OFFSET] = KEY_PAUSE,           // Ctrl+Pause sends the Break scancode
    [0x36 - EXTENDED_KEY_OFFSET] = KEY_RIGHT_SHIFT,     // Some keyboards set the extended bit on Right Shift
};

//----------------------------------------------------------------------------------
// Module Internal Functions Declaration
//----------------------------------------------------------------------------------
// Get ASCII to WCHAR length
static size_t AToWLen(const char *ascii)
{
    int sizeNeeded = MultiByteToWideChar(CP_UTF8, 0, ascii, -1, NULL, 0);

    if (sizeNeeded < 0) TRACELOG(LOG_ERROR, "WIN32: Failed to calculate wide length [ERROR: %u]", GetLastError());

    return sizeNeeded;
}

// Copy ASCII to WCHAR string
static void AToWCopy(const char *ascii, wchar_t *outPtr, size_t outLen)
{
    int size = MultiByteToWideChar(CP_UTF8, 0, ascii, -1, outPtr, (int)outLen);
    if (size != outLen) TRACELOG(LOG_WARNING, "WIN32: Failed to convert %i UTF-8 chars to WCHAR, converted %i chars", outLen, size);
}

// Dynamically load Windows functions that may not be available on older versions
static void LoadWin32FunctionPointers(void)
{
    HMODULE user32 = GetModuleHandleW(L"user32.dll");

    pGetDpiForWindow = (GetDpiForWindowProc)GetProcAddress(user32, "GetDpiForWindow");
    pAdjustWindowRectExForDpi = (AdjustWindowRectExForDpiProc)GetProcAddress(user32, "AdjustWindowRectExForDpi");
    pSetProcessDpiAwarenessContext = (SetProcessDpiAwarenessContextProc)GetProcAddress(user32, "SetProcessDpiAwarenessContext");

    pSetProcessDPIAware = (SetProcessDPIAwareProc)GetProcAddress(user32, "SetProcessDPIAware");
    pGetThreadDpiAwarenessContext = (GetThreadDpiAwarenessContextProc)GetProcAddress(user32, "GetThreadDpiAwarenessContext");
    pAreDpiAwarenessContextsEqual = (AreDpiAwarenessContextsEqualProc)GetProcAddress(user32, "AreDpiAwarenessContextsEqual");

    // SetProcessDpiAwarenessContext (Windows 10 v1703) supersedes SetProcessDpiAwareness (Windows 8.1)
    // We can avoid loading the shcore.dll if the newer DPI functions are available
    if (pSetProcessDpiAwarenessContext == NULL)
    {
        shcore = LoadLibraryW(L"shcore.dll");
        if (shcore != NULL)
        {
            pSetProcessDpiAwareness = (SetProcessDpiAwarenessProc)GetProcAddress(shcore, "SetProcessDpiAwareness");
            pGetDpiForMonitor = (GetDpiForMonitorProc)GetProcAddress(shcore, "GetDpiForMonitor");
        }
    }

    pGetPointerType = (GetPointerTypeProc)GetProcAddress(user32, "GetPointerType");

    const WCHAR *xinputNames[] = { L"xinput1_4.dll", L"xinput1_3.dll", L"xinput9_1_0.dll" };
    for (int i = 0; (i < 3) && (xinput == NULL); i++) xinput = LoadLibraryW(xinputNames[i]);
    if (xinput != NULL)
    {
        pXInputGetState = (XInputGetStateProc)GetProcAddress(xinput, "XInputGetState");
        pXInputSetState = (XInputSetStateProc)GetProcAddress(xinput, "XInputSetState");
    }

    imm32 = LoadLibraryW(L"imm32.dll");
    if (imm32 != NULL)
    {
        pImmGetContext = (ImmGetContextProc)GetProcAddress(imm32, "ImmGetContext");
        pImmReleaseContext = (ImmReleaseContextProc)GetProcAddress(imm32, "ImmReleaseContext");
        pImmSetCompositionWindow = (ImmSetCompositionWindowProc)GetProcAddress(imm32, "ImmSetCompositionWindow");
        pImmSetCandidateWindow = (ImmSetCandidateWindowProc)GetProcAddress(imm32, "ImmSetCandidateWindow");
    }
}

// GetDpiForWindow with a fallback for older versions of Windows
static UINT GetDpiForWindowCompat(HWND hwnd)
{
    UINT result = 0;

    if (pGetDpiForWindow != NULL)
    {
        // This API requires Windows 10 v1607 or newer
        // Returns 0 on failure
        result = pGetDpiForWindow(hwnd);
    }
    else if (pGetDpiForMonitor != NULL)
    {
        // Windows 8.1 and early versions of Windows 10
        // Get the DPI based on which monitor the window is on
        HMONITOR monitor = MonitorFromWindow(hwnd, MONITOR_DEFAULTTONEAREST);
        UINT dpiX = 0;
        UINT dpiY = 0;
        HRESULT hr = pGetDpiForMonitor(monitor, MDT_EFFECTIVE_DPI, &dpiX, &dpiY);
        if (hr == S_OK) result = dpiX;
    }

    if (result == 0)
    {
        // If running on older versions of Windows or if the modern API calls fail: fallback to system-wide DPI
        // We can call GetDC without a corresponding ReleaseDC because the window was created with CS_OWNDC
        result = GetDeviceCaps(GetDC(hwnd), LOGPIXELSX);
    }

    return result;
}

// AdjustWindowRectExForDpi with a fallback for older versions of Windows
static BOOL AdjustWindowRectExForDpiCompat(LPRECT lpRect, DWORD dwStyle, BOOL bMenu, DWORD dwExStyle, UINT dpi)
{
    BOOL result = FALSE;
    if ((pAdjustWindowRectExForDpi != NULL) && platform.perMonitorDPIAwareV2)
    {
        result = pAdjustWindowRectExForDpi(lpRect, dwStyle, bMenu, dwExStyle, dpi);
    }
    else
    {
        result = AdjustWindowRectEx(lpRect, dwStyle, bMenu, dwExStyle);
    }
    return result;
}

// Set process to highest level of dpi awareness on this version of Windows
// DPI awareness can only be set once. This function will fail on subsequent calls (which is harmless)
// Windows 10 v1703+  - per-monitor awareness v2
// Windows 8.1/10     - per-monitor awareness v1
// Windows Vista/7/8  - system-aware
// Windows XP         - No DPI awareness
static void SetProcessDpiAwareCompat(void)
{
    if (pSetProcessDpiAwarenessContext != NULL)
    {
        // Windows 10 v1703+
        BOOL succeeded = pSetProcessDpiAwarenessContext((HANDLE)DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2);

        if (!succeeded)
        {
            // If monitor aware v2 is not available, try v1
            pSetProcessDpiAwarenessContext((HANDLE)DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE);
        }
    }
    else if (pSetProcessDpiAwareness != NULL)
    {
        // Windows 8.1+
        pSetProcessDpiAwareness(PROCESS_PER_MONITOR_DPI_AWARE);
    }
    else if (pSetProcessDPIAware != NULL)
    {
        // Windows Vista+
        pSetProcessDPIAware();
    }

    platform.perMonitorDPIAwareV2 = false;
    if ((pGetThreadDpiAwarenessContext != NULL) && (pAreDpiAwarenessContextsEqual != NULL))
    {
        // There are subtle differences between per-monitor v1 and v2
        platform.perMonitorDPIAwareV2 = pAreDpiAwarenessContextsEqual(pGetThreadDpiAwarenessContext(), (HANDLE)DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2);
    }
}

static bool DecoratedFromStyle(DWORD style)
{
    if (style & STYLE_FLAGS_UNDECORATED_ON)
    {
        if (style & STYLE_FLAGS_UNDECORATED_OFF) TRACELOG(LOG_ERROR, "WIN32: FLAGS: Style 0x%x has both undecorated on/off flags", style);
        return false; // Not decorated
    }

    DWORD masked = (style & STYLE_FLAGS_UNDECORATED_OFF);
    if (STYLE_FLAGS_UNDECORATED_OFF != masked)
    {
        TRACELOG(LOG_ERROR, "WIN32: FLAGS: Style 0x%x is missing flags 0x%x", masked, masked ^ STYLE_FLAGS_UNDECORATED_OFF);
    }

    return true; // Decorated
}

// Get window style from required flags
static DWORD MakeWindowStyle(unsigned flags)
{
    // Flag is not needed because there are no child windows,
    // but supposedly it improves efficiency, plus, windows adds this
    // flag automatically anyway so it keeps flags in sync with the OS
    DWORD style = WS_CLIPSIBLINGS;

    // NOTE: Fullscreen and borderless windowed modes use an undecorated window with no resize frame
    bool covered = ((flags & FLAG_MASK_COVER) != 0);

    style |= (flags & FLAG_WINDOW_HIDDEN)? 0 : WS_VISIBLE;
    style |= ((flags & FLAG_WINDOW_RESIZABLE) && !covered)? STYLE_FLAGS_RESIZABLE : 0;
    style |= ((flags & FLAG_WINDOW_UNDECORATED) || covered)? STYLE_FLAGS_UNDECORATED_ON : STYLE_FLAGS_UNDECORATED_OFF;

    // Minimized takes precedence over maximized
    int mized = MIZED_NONE;
    if (flags & FLAG_WINDOW_MINIMIZED) mized = MIZED_MIN;
    else if (flags & FLAG_WINDOW_MAXIMIZED) mized = MIZED_MAX;

    switch (mized)
    {
        case MIZED_NONE: break;
        case MIZED_MIN: style |= WS_MINIMIZE; break;
        case MIZED_MAX: style |= WS_MAXIMIZE; break;
        default: break;
    }

    return style;
}

// Check flags state, enforces that the actual window/platform state is in sync with raylib's flags
static void CheckFlags(const char *context, HWND hwnd, DWORD flags, DWORD expectedStyle, DWORD styleCheckMask)
{
    DWORD styleFromFlags = MakeWindowStyle(flags);
    if ((styleFromFlags & styleCheckMask) != (expectedStyle & styleCheckMask))
    {
        TRACELOG(LOG_ERROR, "WIN32: FLAGS: %s: window flags (0x%x) produced style 0x%x which != expected 0x%x (diff=0x%x, mask=0x%x)",
            context, flags, styleFromFlags & styleCheckMask, expectedStyle & styleCheckMask,
            (styleFromFlags & styleCheckMask) ^ (expectedStyle & styleCheckMask), styleCheckMask);
    }

    SetLastError(0);
    LONG actualStyle = (LONG)GetWindowLongPtrW(hwnd, GWL_STYLE);
    if ((actualStyle & styleCheckMask) != (expectedStyle & styleCheckMask))
    {
        TRACELOG(LOG_ERROR, "WIN32: FLAGS: %s: expected style 0x%x but got 0x%x (diff=0x%x, mask=0x%x, lasterror=%lu)",
            context, expectedStyle & styleCheckMask, actualStyle & styleCheckMask,
            (expectedStyle & styleCheckMask) ^ (actualStyle & styleCheckMask),
            styleCheckMask, GetLastError());
    }

    if (styleCheckMask & WS_MINIMIZE)
    {
        bool isIconic = IsIconic(hwnd);
        bool styleMinimized = !!(WS_MINIMIZE & actualStyle);
        if (isIconic != styleMinimized) TRACELOG(LOG_ERROR, "WIN32: FLAGS: IsIconic(%d) != WS_MINIMIZED(%d)", isIconic, styleMinimized);
    }

    if (styleCheckMask & WS_MAXIMIZE)
    {
        WINDOWPLACEMENT placement;
        placement.length = sizeof(placement);
        if (!GetWindowPlacement(hwnd, &placement))
        {
            TRACELOG(LOG_ERROR, "WIN32: FLAGS: %s failed, error=%lu", "GetWindowPlacement", GetLastError());
        }
        bool placementMaximized = (placement.showCmd == SW_SHOWMAXIMIZED);
        bool styleMaximized = WS_MAXIMIZE & actualStyle;
        if (placementMaximized != styleMaximized)
        {
            TRACELOG(LOG_ERROR, "WIN32: FLAGS: Maximized state desync, placement maximized=%d (showCmd=%lu) style maximized=%d",
                placementMaximized, placement.showCmd, styleMaximized);
        }
    }
}

// Calculate window size (with borders, title-bar...) from desired client size (framebuffer size)
static SIZE CalcWindowSize(UINT dpi, SIZE clientSize, DWORD style)
{
    RECT rect = { 0, 0, clientSize.cx, clientSize.cy };

    int result = AdjustWindowRectExForDpiCompat(&rect, style, 0, WINDOW_STYLE_EX, dpi);
    if (result == 0) TRACELOG(LOG_ERROR, "WIN32: Failed to adjust window rect [ERROR: %lu]", GetLastError());

    return (SIZE){ rect.right - rect.left, rect.bottom - rect.top };
}

// Update window size if required
// NOTE: Returns true if the window size was updated, false otherwise
static bool UpdateWindowSize(HWND hwnd, int width, int height, unsigned flags)
{
    if (flags & FLAG_WINDOW_MINIMIZED) return false;

    if (flags & FLAG_WINDOW_MAXIMIZED)
    {
        CheckFlags("UpdateWindowSize(maximized)", hwnd, flags, MakeWindowStyle(flags), STYLE_MASK_ALL);
        return false;
    }

    if (flags & FLAG_MASK_COVER)
    {
        MONITORINFO info = { 0 };
        HMONITOR monitor = MonitorFromWindow(hwnd, MONITOR_DEFAULTTONEAREST);
        info.cbSize = sizeof(info);
        if (!GetMonitorInfoW(monitor, &info)) TRACELOG(LOG_ERROR, "WIN32: Failed to get monitor info [ERROR: %lu]", GetLastError());

        RECT windowRect = { 0 };
        if (!GetWindowRect(hwnd, &windowRect)) TRACELOG(LOG_ERROR, "WIN32: Failed to get window rect [ERROR: %lu]", GetLastError());

        if ((windowRect.left == info.rcMonitor.left) &&
            (windowRect.top == info.rcMonitor.top) &&
            ((windowRect.right - windowRect.left) == (info.rcMonitor.right - info.rcMonitor.left)) &&
            ((windowRect.bottom - windowRect.top) == (info.rcMonitor.bottom - info.rcMonitor.top))) return false;

        if (!SetWindowPos(hwnd, HWND_TOP,
            info.rcMonitor.left, info.rcMonitor.top,
            info.rcMonitor.right - info.rcMonitor.left,
            info.rcMonitor.bottom - info.rcMonitor.top,
            SWP_NOOWNERZORDER))
        {
            TRACELOG(LOG_ERROR, "WIN32: Failed to set window position [ERROR: %lu]", GetLastError());
        }

        return true;
    }

    // Get size in pixels from points, considering high-dpi
    UINT dpi = GetDpiForWindowCompat(hwnd);
    float dpiScale = ((float)dpi)/96.0f;
    bool dpiScaling = flags & FLAG_WINDOW_HIGHDPI;
    SIZE desiredSize = {
        .cx = dpiScaling? (int)roundf((float)width*dpiScale) : width,
        .cy = dpiScaling? (int)roundf((float)height*dpiScale) : height
    };

    // Get client size (framebuffer inside the window)
    RECT rect = { 0 };
    GetClientRect(hwnd, &rect);
    SIZE clientSize = { rect.right, rect.bottom };

    // If client size is alread desired size, no need to update
    if ((clientSize.cx == desiredSize.cx) && (clientSize.cy == desiredSize.cy)) return false;

    TRACELOG(LOG_INFO, "WIN32: Restoring client size from [%dx%d] to [%dx%d] (dpi:%lu dpiScaling:%d app:%ix%i)",
        clientSize.cx, clientSize.cy, desiredSize.cx, desiredSize.cy, dpi, dpiScaling, width, height);

    // Calculate window size from desired framebuffer size and window flags
    SIZE windowSize = CalcWindowSize(dpi, desiredSize, MakeWindowStyle(flags));

    // WARNING: This code must be called after swInit() has been called, after InitPlatform() in [rcore]
    SetWindowPos(hwnd, NULL, 0, 0, windowSize.cx, windowSize.cy, SWP_NOMOVE | SWP_NOZORDER);

    return true;
}

// Get OpenGL function pointers
static void *WglGetProcAddress(const char *procname)
{
    void *proc = (void *)wglGetProcAddress(procname);

    if ((proc == NULL) ||
        // NOTE: Some GPU drivers could return following
        // invalid sentinel values instead of NULL
        (proc == (void *)0x1) ||
        (proc == (void *)0x2) ||
        (proc == (void *)0x3) ||
        (proc == (void *)-1))
    {
        // TODO: Keep gl module pointer as global platform data?
        HMODULE glModule = LoadLibraryW(L"opengl32.dll");
        proc = (void *)GetProcAddress(glModule, procname);

        //if (proc == NULL) TRACELOG(LOG_ERROR, "GL: GetProcAddress() failed to get %s [%p], error=%u", procname, proc, GetLastError());
        //else TRACELOG(LOG_INFO, "GL: Found entry point for %s [%p]", procname, proc);
    }

    return proc;
}

// Get cursor name
static LPCWSTR GetCursorName(int cursor)
{
    LPCWSTR name = (LPCWSTR)IDC_ARROW;

    switch (cursor)
    {
        case MOUSE_CURSOR_DEFAULT: name = (LPCWSTR)IDC_ARROW; break;
        case MOUSE_CURSOR_ARROW: name = (LPCWSTR)IDC_ARROW; break;
        case MOUSE_CURSOR_IBEAM: name = (LPCWSTR)IDC_IBEAM; break;
        case MOUSE_CURSOR_CROSSHAIR: name = (LPCWSTR)IDC_CROSS; break;
        case MOUSE_CURSOR_POINTING_HAND: name = (LPCWSTR)IDC_HAND; break;
        case MOUSE_CURSOR_RESIZE_EW: name = (LPCWSTR)IDC_SIZEWE; break;
        case MOUSE_CURSOR_RESIZE_NS: name = (LPCWSTR)IDC_SIZENS; break;
        case MOUSE_CURSOR_RESIZE_NWSE: name = (LPCWSTR)IDC_SIZENWSE; break;
        case MOUSE_CURSOR_RESIZE_NESW: name = (LPCWSTR)IDC_SIZENESW; break;
        case MOUSE_CURSOR_RESIZE_ALL: name = (LPCWSTR)IDC_SIZEALL; break;
        case MOUSE_CURSOR_NOT_ALLOWED: name = (LPCWSTR)IDC_NO; break;
        default: break;
    }

    return name;
}

// Delete platform.smallIcon and platform.largeIcon
static void DeletePlatformIconHandles(void)
{
    HICON smallIcon = platform.smallIcon;
    HICON largeIcon = platform.largeIcon;

    if (smallIcon != NULL) DestroyIcon(smallIcon);

    // Only delete largeIcon if it's not pointing to the same underlying resource as smallIcon
    if ((largeIcon != NULL) && (largeIcon != smallIcon)) DestroyIcon(largeIcon);

    platform.smallIcon = NULL;
    platform.largeIcon = NULL;
}

// Delete existing icons and replace them with the standard window/taskbar icons
static void ResetToDefaultWindowsIcons(void)
{
    HICON smallIcon = (HICON)GetClassLongPtrW(platform.hwnd, GCLP_HICONSM);
    HICON largeIcon = (HICON)GetClassLongPtrW(platform.hwnd, GCLP_HICON);

    SendMessage(platform.hwnd, WM_SETICON, ICON_SMALL, (LPARAM)smallIcon);
    SendMessage(platform.hwnd, WM_SETICON, ICON_BIG, (LPARAM)largeIcon);

    // The built-in Windows icons shouldn't be deleted, so there's no need to store a reference to them
    // Just delete the previous icons and exit
    DeletePlatformIconHandles();
}

// Convert Image to Win32 HICON (must be RGBA)
static HICON ConvertImageToWindowsIcon(Image image)
{
    if (image.format != PIXELFORMAT_UNCOMPRESSED_R8G8B8A8)
    {
        TRACELOG(LOG_WARNING, "WIN32: Window icon image must be in R8G8B8A8 pixel format");
        return NULL;
    }

    HDC hdc = GetDC(platform.hwnd);

    // Create 32-bit BGRA DIB for color
    BITMAPV5HEADER bi = { 0 };
    ZeroMemory(&bi, sizeof(bi));
    bi.bV5Size = sizeof(bi);
    bi.bV5Width = image.width;
    bi.bV5Height = -image.height; // Negative = top-down bitmap
    bi.bV5Planes = 1;
    bi.bV5BitCount = 32;
    bi.bV5Compression = BI_BITFIELDS;
    bi.bV5RedMask = 0x00FF0000;
    bi.bV5GreenMask = 0x0000FF00;
    bi.bV5BlueMask = 0x000000FF;
    bi.bV5AlphaMask = 0xFF000000;

    unsigned char *targetBits = NULL;
    HBITMAP hColorBitmap = CreateDIBSection(hdc, (BITMAPINFO *)&bi, DIB_RGB_COLORS, (void **)&targetBits, NULL, 0);
    if (!hColorBitmap)
    {
        ReleaseDC(platform.hwnd, hdc);
        return NULL;
    }

    // Copy RGBA > BGRA (Win32 expects BGRA)
    for (int y = 0; y < image.height; y++)
    {
        for (int x = 0; x < image.width; x++)
        {
            int i = (y*image.width + x)*4;
            targetBits[i + 0] = ((unsigned char *)image.data)[i + 2]; // B
            targetBits[i + 1] = ((unsigned char *)image.data)[i + 1]; // G
            targetBits[i + 2] = ((unsigned char *)image.data)[i + 0]; // R
            targetBits[i + 3] = ((unsigned char *)image.data)[i + 3]; // A
        }
    }

    // Create mask bitmap (1-bit, all opaque)
    HBITMAP hMaskBitmap = CreateBitmap(image.width, image.height, 1, 1, NULL);

    // Build icon info
    ICONINFO ii = { 0 };
    ZeroMemory(&ii, sizeof(ii));
    ii.fIcon = TRUE;
    ii.hbmMask = hMaskBitmap;
    ii.hbmColor = hColorBitmap;

    HICON result = CreateIconIndirect(&ii);

    // Clean up GDI bitmaps (icon keeps copies internally)
    DeleteObject(hColorBitmap);
    DeleteObject(hMaskBitmap);
    ReleaseDC(platform.hwnd, hdc);

    return result;
}

// Collect monitors process, primary monitor first
// NOTE: Required by GetMonitorHandles()
static BOOL CALLBACK CollectMonitorsProc(HMONITOR handle, HDC hdc, LPRECT rect, LPARAM lparam)
{
    MonitorList *list = (MonitorList *)lparam;

    if (list->count < MAX_MONITORS)
    {
        MONITORINFO info = { 0 };
        info.cbSize = sizeof(info);

        if (GetMonitorInfoW(handle, &info) && (info.dwFlags & MONITORINFOF_PRIMARY) && (list->count > 0))
        {
            for (int i = list->count; i > 0; i--) list->handles[i] = list->handles[i - 1];
            list->handles[0] = handle;
        }
        else list->handles[list->count] = handle;

        list->count++;
    }

    // Always return TRUE to continue the loop, otherwise, the caller
    // can't distinguish between stopping the loop and an error
    return TRUE;
}

// Get connected monitors, primary monitor first
static MonitorList GetMonitorHandles(void)
{
    MonitorList list = { 0 };

    int result = EnumDisplayMonitors(NULL, NULL, CollectMonitorsProc, (LPARAM)&list);
    if (result == 0) TRACELOG(LOG_ERROR, "%s failed, error=%lu", "EnumDisplayMonitors", GetLastError());

    return list;
}

// Get monitor info for a monitor index, returns false if the monitor is not available
static bool GetMonitorInfoFromIndex(int monitor, MONITORINFOEXW *info)
{
    MonitorList list = GetMonitorHandles();

    if ((monitor < 0) || (monitor >= list.count))
    {
        TRACELOG(LOG_WARNING, "WIN32: Failed to find selected monitor");
        return false;
    }

    ZeroMemory(info, sizeof(MONITORINFOEXW));
    info->cbSize = sizeof(MONITORINFOEXW);

    return (GetMonitorInfoW(list.handles[monitor], (MONITORINFO *)info) != 0);
}

// Get style changed required operations flags
// NOTE: Required for deferred operations
static void GetStyleChangeFlagOps(DWORD coreWindowFlags, STYLESTRUCT *style, FlagsOp *deferredFlags)
{
    // Check window resizable flag change
    // NOTE: Fullscreen and borderless windowed modes have no resize frame
    bool resizable = ((coreWindowFlags & FLAG_WINDOW_RESIZABLE) && !(coreWindowFlags & FLAG_MASK_COVER));
    bool resizableOld = ((style->styleOld & STYLE_FLAGS_RESIZABLE) != 0);
    bool resizableNew = ((style->styleNew & STYLE_FLAGS_RESIZABLE) != 0);
    if (resizable != resizableOld) TRACELOG(LOG_ERROR, "WIN32: Expected resizable %u but got %u", resizable, resizableOld);
    if (resizableOld != resizableNew)
    {
        if (resizableNew) deferredFlags->set |= FLAG_WINDOW_RESIZABLE;
        else deferredFlags->clear |= FLAG_WINDOW_RESIZABLE;
    }

    // Check window decorated flag change
    bool decorated = (0 == (coreWindowFlags & FLAG_WINDOW_UNDECORATED));
    bool decoratedOld = DecoratedFromStyle(style->styleOld);
    bool decoratedNew = DecoratedFromStyle(style->styleNew);
    if (decorated != decoratedOld) TRACELOG(LOG_ERROR, "WIN32: Expected decorated %u but got %u", decorated, decoratedOld);
    if (decoratedOld != decoratedNew)
    {
        if (decoratedNew) deferredFlags->clear |= FLAG_WINDOW_UNDECORATED;
        else deferredFlags->set |= FLAG_WINDOW_UNDECORATED;
    }

    // Check window hidden flag change
    bool hidden = (coreWindowFlags & FLAG_WINDOW_HIDDEN);
    bool hiddenOld = ((style->styleOld & WS_VISIBLE) == 0);
    bool hiddenNew = ((style->styleNew & WS_VISIBLE) == 0);
    if (hidden != hiddenOld) TRACELOG(LOG_ERROR, "WIN32: Expected hidden %u but got %u", hidden, hiddenOld);
    if (hiddenOld != hiddenNew)
    {
        if (hiddenNew) deferredFlags->set |= FLAG_WINDOW_HIDDEN;
        else deferredFlags->clear |= FLAG_WINDOW_HIDDEN;
    }
}

// Adopt window resize
// NOTE: Call when the window is rezised, returns true
// if the new window size should update the desired app size
static bool AdoptWindowResize(unsigned flags)
{
    if (flags & FLAG_WINDOW_MINIMIZED) return false;
    if (flags & FLAG_WINDOW_MAXIMIZED) return false;
    if (flags & FLAG_FULLSCREEN_MODE) return false;
    if (flags & FLAG_BORDERLESS_WINDOWED_MODE) return false;
    if (!(flags & FLAG_WINDOW_RESIZABLE)) return false;

    return true;
}

// ---------------------------------------------------------------------------------------------
// Here's the end of the "pure function section", the rest of the file can access global state
// ---------------------------------------------------------------------------------------------

// Unlock the ability to use CORE in the rest of the file
#undef CORE

//----------------------------------------------------------------------------------
// Module Internal Functions Declaration
//----------------------------------------------------------------------------------
int InitPlatform(void);         // Initialize platform (graphics, inputs and more)
void ClosePlatform(void);       // Close platform

// Win32 process messages management function
static LRESULT CALLBACK WndProc(HWND hwnd, UINT msg, WPARAM wparam, LPARAM lparam);

// Win32: Handle inputs functions
static void HandleKey(WPARAM wparam, LPARAM lparam, char state);
static void HandleChar(unsigned int codepoint);                 // Register character input
static void HandleTouch(int touchAction, int id, float x, float y); // Register touch event
static void HandleDroppedFiles(HDROP drop);                     // Register dropped files
static void PollGamepads(void);                                 // Update gamepads state (XInput)
static void UpdateCursorClip(void);                             // Confine or release the cursor (disabled cursor)
static void UpdateInputMethodPosition(HWND hwnd);               // Place input method windows next to the last clicked position
static void HandleMouseButton(int button, char state);
static void HandleRawInput(LPARAM lparam);
static void HandleWindowResize(HWND hwnd, int *width, int *height);

static void UpdateWindowStyle(HWND hwnd, unsigned desiredFlags);
static unsigned SanitizeFlags(unsigned flags);
static void UpdateFlags(HWND hwnd, unsigned desiredFlags, int width, int height); // Update window flags
static void UpdateLayeredStyle(HWND hwnd, bool passthrough, float opacity); // Update window layered style: mouse passthrough and opacity
static void EnterCoverMode(HWND hwnd, unsigned desiredFlags);   // Cover the monitor with the window: fullscreen or borderless windowed mode
static void LeaveCoverMode(HWND hwnd, unsigned desiredFlags);   // Restore the window from fullscreen or borderless windowed mode

// Check if OpenGL extension is available
static bool IsWglExtensionAvailable(HDC hdc, const char *extension);

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
// NOTE: The window covers the monitor, the video mode is not changed
void ToggleFullscreen(void)
{
    if (CORE.Window.flags & FLAG_FULLSCREEN_MODE) ClearWindowState(FLAG_FULLSCREEN_MODE);
    else SetWindowState(FLAG_FULLSCREEN_MODE);
}

// Toggle borderless windowed mode
void ToggleBorderlessWindowed(void)
{
    if (CORE.Window.flags & FLAG_BORDERLESS_WINDOWED_MODE) ClearWindowState(FLAG_BORDERLESS_WINDOWED_MODE);
    else SetWindowState(FLAG_BORDERLESS_WINDOWED_MODE);
}

// Set window state: maximized, if resizable
void MaximizeWindow(void)
{
    SetWindowState(FLAG_WINDOW_MAXIMIZED);
}

// Set window state: minimized
void MinimizeWindow(void)
{
    SetWindowState(FLAG_WINDOW_MINIMIZED);
}

// Restore window from being minimized/maximized
void RestoreWindow(void)
{
    if ((CORE.Window.flags & FLAG_WINDOW_MAXIMIZED) &&
        (CORE.Window.flags & FLAG_WINDOW_MINIMIZED)) ClearWindowState(FLAG_WINDOW_MINIMIZED);
    else ClearWindowState(FLAG_WINDOW_MINIMIZED | FLAG_WINDOW_MAXIMIZED);
}

// Set window configuration state using flags
void SetWindowState(unsigned int flags)
{
    if (flags & FLAG_MSAA_4X_HINT) TRACELOG(LOG_WARNING, "WINDOW: MSAA can only be configured before window initialization");
    if (flags & FLAG_WINDOW_HIGHDPI) TRACELOG(LOG_WARNING, "WINDOW: High DPI can only be configured before window initialization");

    unsigned requestedFlags = CORE.Window.flags | flags;

    // Fullscreen and borderless windowed modes exclude each other, last one requested is kept
    if (flags & FLAG_FULLSCREEN_MODE) requestedFlags &= ~FLAG_BORDERLESS_WINDOWED_MODE;
    else if (flags & FLAG_BORDERLESS_WINDOWED_MODE) requestedFlags &= ~FLAG_FULLSCREEN_MODE;

    // Decorations and maximized state are restored when leaving fullscreen and borderless windowed modes
    if (requestedFlags & FLAG_MASK_COVER)
    {
        if (!(CORE.Window.flags & FLAG_MASK_COVER))
        {
            platform.coverRestoreUndecorated = ((requestedFlags & FLAG_WINDOW_UNDECORATED) != 0);
            platform.coverRestoreMaximized = ((requestedFlags & FLAG_WINDOW_MAXIMIZED) != 0);
        }
        else
        {
            if (flags & FLAG_WINDOW_UNDECORATED) platform.coverRestoreUndecorated = true;
            if (flags & FLAG_WINDOW_MAXIMIZED) platform.coverRestoreMaximized = true;
        }
    }

    platform.desiredFlags = SanitizeFlags(requestedFlags);
    UpdateFlags(platform.hwnd, platform.desiredFlags, platform.appScreenWidth, platform.appScreenHeight);
}

// Clear window configuration state flags
void ClearWindowState(unsigned int flags)
{
    if (flags & FLAG_MSAA_4X_HINT) TRACELOG(LOG_WARNING, "WINDOW: MSAA can only be configured before window initialization");
    if (flags & FLAG_WINDOW_HIGHDPI) TRACELOG(LOG_WARNING, "WINDOW: High DPI can only be configured before window initialization");

    unsigned requestedFlags = CORE.Window.flags & ~flags;

    // Decorations and maximized state are restored when leaving fullscreen and borderless windowed modes
    if (CORE.Window.flags & FLAG_MASK_COVER)
    {
        if (flags & FLAG_WINDOW_UNDECORATED) platform.coverRestoreUndecorated = false;
        if (flags & FLAG_WINDOW_MAXIMIZED) platform.coverRestoreMaximized = false;

        if (!(requestedFlags & FLAG_MASK_COVER))
        {
            if (platform.coverRestoreUndecorated) requestedFlags |= FLAG_WINDOW_UNDECORATED;
            else requestedFlags &= ~FLAG_WINDOW_UNDECORATED;

            if (platform.coverRestoreMaximized && (requestedFlags & FLAG_WINDOW_RESIZABLE)) requestedFlags |= FLAG_WINDOW_MAXIMIZED;
        }
    }

    platform.desiredFlags = SanitizeFlags(requestedFlags);
    UpdateFlags(platform.hwnd, platform.desiredFlags, platform.appScreenWidth, platform.appScreenHeight);
}

// Set icon for window
void SetWindowIcon(Image image)
{
    if (!platform.hwnd) return;

    if ((image.data == NULL) || (image.width <= 0) || (image.height <= 0))
    {
        ResetToDefaultWindowsIcons();
        return;
    }

    HICON hIcon = ConvertImageToWindowsIcon(image);
    if (hIcon)
    {
        SendMessage(platform.hwnd, WM_SETICON, ICON_SMALL, (LPARAM)hIcon);
        SendMessage(platform.hwnd, WM_SETICON, ICON_BIG, (LPARAM)hIcon);
        DeletePlatformIconHandles();
        platform.smallIcon = hIcon;
        platform.largeIcon = hIcon;
    }
}

// Set icon for window (multiple sizes)
void SetWindowIcons(Image *images, int count)
{
    if ((images == NULL) || (count <= 0))
    {
        ResetToDefaultWindowsIcons();
        return;
    }

    // Large icon has a resolution of (SM_CXICON*SM_CXICON) pixels
    int largeIconRenderSize = GetSystemMetrics(SM_CXICON);

    // Attempt to find a best match for the large icon based on each image's width
    int bestLargeIconIndex = -1;
    int bestLargeIconSizeDifference = INT_MAX;
    for (int i = 0; i < count; i++)
    {
        Image img = images[i];
        int difference = abs(img.width - largeIconRenderSize);
        if ((img.data != NULL) && (img.format == PIXELFORMAT_UNCOMPRESSED_R8G8B8A8) && (difference < bestLargeIconSizeDifference))
        {
            bestLargeIconIndex = i;
            bestLargeIconSizeDifference = difference;
        }
    }

    // Small icon has a resolution of (SM_CYSMICON*SM_CYSMICON) pixels
    int smallIconRenderSize = GetSystemMetrics(SM_CYSMICON);

    // Attempt to find a best match for the small icon based on each image's width
    int bestSmallIconIndex = -1;
    int bestSmallIconSizeDifference = INT_MAX;
    for (int i = 0; i < count; i++)
    {
        Image img = images[i];
        int difference = abs(img.width - smallIconRenderSize);
        if ((img.data != NULL) && (img.format == PIXELFORMAT_UNCOMPRESSED_R8G8B8A8) && (difference < bestSmallIconSizeDifference))
        {
            bestSmallIconIndex = i;
            bestSmallIconSizeDifference = difference;
        }
    }

    if ((bestSmallIconIndex != -1) && (bestLargeIconIndex != -1))
    {
        HICON smallIcon = ConvertImageToWindowsIcon(images[bestSmallIconIndex]);
        HICON largeIcon = ConvertImageToWindowsIcon(images[bestLargeIconIndex]);

        SendMessage(platform.hwnd, WM_SETICON, ICON_SMALL, (LPARAM)smallIcon);
        SendMessage(platform.hwnd, WM_SETICON, ICON_BIG, (LPARAM)largeIcon);

        DeletePlatformIconHandles();

        platform.smallIcon = smallIcon;
        platform.largeIcon = largeIcon;
    }
    else TRACELOG(LOG_WARNING, "WIN32: All icons passed to SetWindowIcons were invalid");
}

void SetWindowTitle(const char *title)
{
    CORE.Window.title = title;

    WCHAR *titleWide = NULL;
    A_TO_W_ALLOCA(titleWide, CORE.Window.title);

    int result = SetWindowTextW(platform.hwnd, titleWide);
    if (result == 0) TRACELOG(LOG_WARNING, "WIN32: Failed to set window title [ERROR: %lu]", GetLastError());
}

// Set window position on screen (windowed mode)
void SetWindowPosition(int x, int y)
{
    if (platform.hwnd != NULL)
    {
        // NOTE: Position refers to the window client area, as returned by GetWindowPosition()
        RECT rect = { x, y, x, y };
        DWORD style = (DWORD)GetWindowLongPtrW(platform.hwnd, GWL_STYLE);
        AdjustWindowRectExForDpiCompat(&rect, style, FALSE, WINDOW_STYLE_EX, GetDpiForWindowCompat(platform.hwnd));

        // Move the window to the new position (keeping size and z-order)
        SetWindowPos(platform.hwnd, NULL, rect.left, rect.top, 0, 0, SWP_NOSIZE | SWP_NOZORDER | SWP_NOACTIVATE);
    }
}

// Set monitor for the current window
void SetWindowMonitor(int monitor)
{
    MONITORINFOEXW info = { 0 };
    if (!GetMonitorInfoFromIndex(monitor, &info)) return;

    TRACELOG(LOG_INFO, "WIN32: Selected monitor: [%i]", monitor);

    RECT windowRect = { 0 };
    GetWindowRect(platform.hwnd, &windowRect);
    int width = windowRect.right - windowRect.left;
    int height = windowRect.bottom - windowRect.top;

    if (CORE.Window.flags & FLAG_MASK_COVER)
    {
        // Cover the selected monitor, position to restore is moved to it
        int restoreWidth = platform.coverRestoreRect.right - platform.coverRestoreRect.left;
        int restoreHeight = platform.coverRestoreRect.bottom - platform.coverRestoreRect.top;
        platform.coverRestoreRect.left = info.rcWork.left + ((info.rcWork.right - info.rcWork.left) - restoreWidth)/2;
        platform.coverRestoreRect.top = info.rcWork.top + ((info.rcWork.bottom - info.rcWork.top) - restoreHeight)/2;
        platform.coverRestoreRect.right = platform.coverRestoreRect.left + restoreWidth;
        platform.coverRestoreRect.bottom = platform.coverRestoreRect.top + restoreHeight;

        SetWindowPos(platform.hwnd, NULL, info.rcMonitor.left, info.rcMonitor.top,
            info.rcMonitor.right - info.rcMonitor.left, info.rcMonitor.bottom - info.rcMonitor.top, SWP_NOZORDER | SWP_NOACTIVATE);
    }
    else if (!(CORE.Window.flags & (FLAG_WINDOW_MINIMIZED | FLAG_WINDOW_MAXIMIZED)))
    {
        // Center the window on the selected monitor
        int x = info.rcWork.left + ((info.rcWork.right - info.rcWork.left) - width)/2;
        int y = info.rcWork.top + ((info.rcWork.bottom - info.rcWork.top) - height)/2;

        SetWindowPos(platform.hwnd, NULL, x, y, 0, 0, SWP_NOSIZE | SWP_NOZORDER | SWP_NOACTIVATE);
    }
}

// Set window minimum dimensions (FLAG_WINDOW_RESIZABLE)
void SetWindowMinSize(int width, int height)
{
    if ((width > CORE.Window.screenMax.width) || (height > CORE.Window.screenMax.height))
    {
        TRACELOG(LOG_WARNING, "WIN32: WINDOW: Cannot set minimum screen size higher than the maximum");
        return;
    }

    CORE.Window.screenMin.width = width;
    CORE.Window.screenMin.height = height;

    SetWindowSize(platform.appScreenWidth, platform.appScreenHeight);
}

// Set window maximum dimensions (FLAG_WINDOW_RESIZABLE)
void SetWindowMaxSize(int width, int height)
{
    if ((width < CORE.Window.screenMin.width) || (height < CORE.Window.screenMin.height))
    {
        TRACELOG(LOG_WARNING, "WIN32: WINDOW: Cannot set maximum screen size lower than the minimum");
        return;
    }

    CORE.Window.screenMax.width = width;
    CORE.Window.screenMax.height = height;

    SetWindowSize(platform.appScreenWidth, platform.appScreenHeight);
}

// Set window dimensions
void SetWindowSize(int width, int height)
{
    int screenWidth = fmaxf(CORE.Window.screenMin.width, fminf(CORE.Window.screenMax.width, width));
    int screenHeight = fmaxf(CORE.Window.screenMin.height, fminf(CORE.Window.screenMax.height, height));

    // NOTE: Size is applied when the window is restored if it is currently minimized, maximized or covering the monitor
    platform.appScreenWidth = screenWidth;
    platform.appScreenHeight = screenHeight;

    UpdateWindowSize(platform.hwnd, screenWidth, screenHeight, platform.desiredFlags);
}

// Set window opacity, value opacity is between 0.0 and 1.0
void SetWindowOpacity(float opacity)
{
    if (opacity >= 1.0f) opacity = 1.0f;
    else if (opacity <= 0.0f) opacity = 0.0f;

    platform.opacity = opacity;
    UpdateLayeredStyle(platform.hwnd, ((CORE.Window.flags & FLAG_WINDOW_MOUSE_PASSTHROUGH) != 0), platform.opacity);
}

// Set window focused
void SetWindowFocused(void)
{
    BringWindowToTop(platform.hwnd);
    SetForegroundWindow(platform.hwnd);
    SetFocus(platform.hwnd);

    if (GetForegroundWindow() != platform.hwnd) TRACELOG(LOG_WARNING, "WIN32: WINDOW: Failed to set window focus");
}

// Get native window handle
void *GetWindowHandle(void)
{
    return (void *)platform.hwnd;
}

int GetMonitorCount(void)
{
    MonitorList list = GetMonitorHandles();

    return list.count;
}

// Get current monitor where window is placed
int GetCurrentMonitor(void)
{
    int index = 0;

    HMONITOR monitor = MonitorFromWindow(platform.hwnd, MONITOR_DEFAULTTOPRIMARY);
    if (!monitor) TRACELOG(LOG_ERROR, "%s failed, error=%lu", "MonitorFromWindow", GetLastError());

    MonitorList list = GetMonitorHandles();
    for (int i = 0; i < list.count; i++)
    {
        if (list.handles[i] == monitor)
        {
            index = i;
            break;
        }
    }

    return index;
}

// Get selected monitor position
Vector2 GetMonitorPosition(int monitor)
{
    MONITORINFOEXW info = { 0 };
    if (!GetMonitorInfoFromIndex(monitor, &info)) return (Vector2){ 0, 0 };

    return (Vector2){ (float)info.rcMonitor.left, (float)info.rcMonitor.top };
}

// Get selected monitor width (currently used by monitor)
int GetMonitorWidth(int monitor)
{
    MONITORINFOEXW info = { 0 };
    if (!GetMonitorInfoFromIndex(monitor, &info)) return 0;

    return (int)(info.rcMonitor.right - info.rcMonitor.left);
}

// Get selected monitor height (currently used by monitor)
int GetMonitorHeight(int monitor)
{
    MONITORINFOEXW info = { 0 };
    if (!GetMonitorInfoFromIndex(monitor, &info)) return 0;

    return (int)(info.rcMonitor.bottom - info.rcMonitor.top);
}

// Get selected monitor physical width in millimetres
int GetMonitorPhysicalWidth(int monitor)
{
    int width = 0;

    MONITORINFOEXW info = { 0 };
    if (!GetMonitorInfoFromIndex(monitor, &info)) return width;

    HDC hdc = CreateDCW(L"DISPLAY", info.szDevice, NULL, NULL);
    if (hdc != NULL)
    {
        width = GetDeviceCaps(hdc, HORZSIZE);
        DeleteDC(hdc);
    }

    return width;
}

// Get selected monitor physical height in millimetres
int GetMonitorPhysicalHeight(int monitor)
{
    int height = 0;

    MONITORINFOEXW info = { 0 };
    if (!GetMonitorInfoFromIndex(monitor, &info)) return height;

    HDC hdc = CreateDCW(L"DISPLAY", info.szDevice, NULL, NULL);
    if (hdc != NULL)
    {
        height = GetDeviceCaps(hdc, VERTSIZE);
        DeleteDC(hdc);
    }

    return height;
}

// Get selected monitor refresh rate
int GetMonitorRefreshRate(int monitor)
{
    MONITORINFOEXW info = { 0 };
    if (!GetMonitorInfoFromIndex(monitor, &info)) return 0;

    DEVMODEW mode = { 0 };
    mode.dmSize = sizeof(mode);
    if (!EnumDisplaySettingsW(info.szDevice, ENUM_CURRENT_SETTINGS, &mode)) return 0;

    return (int)mode.dmDisplayFrequency;
}

// Get the human-readable, UTF-8 encoded name of the selected monitor
const char *GetMonitorName(int monitor)
{
    static char name[128] = { 0 };
    name[0] = '\0';

    MONITORINFOEXW info = { 0 };
    if (!GetMonitorInfoFromIndex(monitor, &info)) return name;

    // NOTE: Monitor attached to the display adapter, adapter name is used if not available
    DISPLAY_DEVICEW device = { 0 };
    device.cb = sizeof(device);
    if (EnumDisplayDevicesW(info.szDevice, 0, &device, 0)) WideCharToMultiByte(CP_UTF8, 0, device.DeviceString, -1, name, sizeof(name) - 1, NULL, NULL);
    else WideCharToMultiByte(CP_UTF8, 0, info.szDevice, -1, name, sizeof(name) - 1, NULL, NULL);

    return name;
}

// Get window position XY on monitor
// NOTE: Position of the window client area
Vector2 GetWindowPosition(void)
{
    POINT position = { 0 };
    if (platform.hwnd != NULL) ClientToScreen(platform.hwnd, &position);

    return (Vector2){ (float)position.x, (float)position.y };
}

// Get window scale DPI factor for current monitor
Vector2 GetWindowScaleDPI(void)
{
    float scale = ((float)GetDpiForWindowCompat(platform.hwnd))/96.0f;
    return (Vector2){ scale, scale };
}

// Set clipboard text content
void SetClipboardText(const char *text)
{
    if (OpenClipboard(platform.hwnd))
    {
        if (text == NULL) EmptyClipboard();
        else
        {
            size_t length = strlen(text);

            int clipboardTextLength = MultiByteToWideChar(CP_UTF8, 0, text, (int)(length + 1), NULL, 0);
            if (clipboardTextLength != 0)
            {
                HGLOBAL globalHandle = GlobalAlloc(GMEM_MOVEABLE, clipboardTextLength*sizeof(wchar_t));

                if (globalHandle)
                {
                    LPVOID clipboardText = GlobalLock(globalHandle);
                    int charactersCopied = 0;
                    if (clipboardText != NULL)
                    {
                        // Directly copies text to clipboardText and converts from UTF-8 to UTF-16
                        charactersCopied = MultiByteToWideChar(CP_UTF8, 0, text, (int)(length + 1), clipboardText, clipboardTextLength);
                        GlobalUnlock(globalHandle);
                    }

                    if (charactersCopied != 0)
                    {
                        EmptyClipboard();

                        // We don't need to call GlobalFree when SetClipboardData succeeds; Windows takes control of the memory
                        HANDLE clipboardData = SetClipboardData(CF_UNICODETEXT, globalHandle);
                        if (!clipboardData)
                        {
                            TRACELOG(LOG_WARNING, "Clipboard: Failed to set clipboard text.");
                            GlobalFree(globalHandle);
                        }
                    }
                    else GlobalFree(globalHandle);
                }
            }
        }

        CloseClipboard();
    }
    else TRACELOG(LOG_WARNING, "Clipboard: Failed to write to clipboard. It may be locked by another program.");
}

// Get clipboard text content
const char *GetClipboardText(void)
{
    if (platform.clipboardText != NULL)
    {
        RL_FREE(platform.clipboardText);
        platform.clipboardText = NULL;
    }

    if (OpenClipboard(platform.hwnd))
    {
        HANDLE clipboardDataHandle = GetClipboardData(CF_UNICODETEXT);
        if (clipboardDataHandle)
        {
            wchar_t *clipboardText = GlobalLock(clipboardDataHandle);
            if (clipboardText != NULL)
            {
                // The clipboard data can't be trusted. It may or may not be null-terminated!
                size_t maxLength = GlobalSize(clipboardDataHandle)/sizeof(wchar_t);
                int length = (int)wcsnlen(clipboardText, maxLength);

                int size = WideCharToMultiByte(CP_UTF8, 0, clipboardText, length, NULL, 0, NULL, NULL);
                char *utf8Text = (char *)RL_CALLOC(size + 1, 1);

                if ((size > 0) && (WideCharToMultiByte(CP_UTF8, 0, clipboardText, length, utf8Text, size, NULL, NULL) == 0)) size = 0;
                utf8Text[size] = '\0';

                platform.clipboardText = utf8Text;

                GlobalUnlock(clipboardDataHandle);
            }
        }

        CloseClipboard();
    }
    else
    {
        TRACELOG(LOG_WARNING, "Clipboard: Failed to read clipboard. It may be locked by another program.");
    }

    return platform.clipboardText;
}

#if SUPPORT_CLIPBOARD_IMAGE && SUPPORT_MODULE_RTEXTURES
    #define WIN32_CLIPBOARD_IMPLEMENTATION
    #define WINUSER_ALREADY_INCLUDED
    #define WINBASE_ALREADY_INCLUDED
    #define WINGDI_ALREADY_INCLUDED
    #include "../external/win32_clipboard.h"
#endif

// Get clipboard image
Image GetClipboardImage(void)
{
    Image image = { 0 };

#if SUPPORT_CLIPBOARD_IMAGE && SUPPORT_MODULE_RTEXTURES
    unsigned int dataSize = 0;
    void *fileData = NULL;
    int width = 0;
    int height = 0;

    fileData = (void *)Win32GetClipboardImageData(&width, &height, &dataSize);

    if (fileData == NULL) TRACELOG(LOG_WARNING, "Clipboard image: Couldn't get clipboard data.");
    else
    {
        image = LoadImageFromMemory(".bmp", (const unsigned char*)fileData, (int)dataSize);
        RL_FREE(fileData);
    }
#else
    TRACELOG(LOG_WARNING, "Clipboard image: SUPPORT_CLIPBOARD_IMAGE requires SUPPORT_MODULE_RTEXTURES to work properly");
#endif // SUPPORT_CLIPBOARD_IMAGE

    return image;
}

// Show mouse cursor
void ShowCursor(void)
{
    SetCursor(platform.cursor);
    CORE.Input.Mouse.cursorHidden = false;
}

// Hide mouse cursor
void HideCursor(void)
{
    // NOTE: Using SetCursor() instead of ShowCursor() because
    // it makes it easy to only hide the cursor while it's inside the client area
    SetCursor(NULL);
    CORE.Input.Mouse.cursorHidden = true;
}

// Enable cursor (unlock cursor)
void EnableCursor(void)
{
    if (CORE.Input.Mouse.cursorLocked)
    {
        CORE.Input.Mouse.cursorLocked = false;
        UpdateCursorClip();

        RAWINPUTDEVICE rid = { 0 };
        rid.usUsagePage = 0x01; // HID_USAGE_PAGE_GENERIC
        rid.usUsage = 0x02; // HID_USAGE_GENERIC_MOUSE
        rid.dwFlags = RIDEV_REMOVE;
        rid.hwndTarget = NULL;
        int result = RegisterRawInputDevices(&rid, 1, sizeof(rid));
        if (result == 0) TRACELOG(LOG_WARNING, "WIN32: Failed to register raw input devices [ERROR: %lu]", GetLastError());

        // Restore cursor position, mouse position follows it with no movement registered
        SetCursorPos(platform.cursorRestorePosition.x, platform.cursorRestorePosition.y);

        POINT position = platform.cursorRestorePosition;
        if (ScreenToClient(platform.hwnd, &position))
        {
            CORE.Input.Mouse.currentPosition = (Vector2){ (float)position.x, (float)position.y };
            CORE.Input.Mouse.previousPosition = CORE.Input.Mouse.currentPosition;
        }

        ShowCursor();
    }
}

// Disables cursor (lock cursor)
// NOTE: Cursor is hidden and confined to the window, mouse position is updated from raw mouse motion
void DisableCursor(void)
{
    if (!CORE.Input.Mouse.cursorLocked)
    {
        // NOTE: First movement of an absolute pointing device is relative to the cursor position
        platform.rawAbsoluteValid = (GetCursorPos(&platform.cursorRestorePosition) != 0);
        platform.rawAbsolutePosition = platform.cursorRestorePosition;

        RAWINPUTDEVICE rid = { 0 };
        rid.usUsagePage = 0x01; // HID_USAGE_PAGE_GENERIC
        rid.usUsage = 0x02; // HID_USAGE_GENERIC_MOUSE
        rid.dwFlags = 0;
        rid.hwndTarget = platform.hwnd;
        int result = RegisterRawInputDevices(&rid, 1, sizeof(rid));
        if (result == 0) TRACELOG(LOG_WARNING, "WIN32: Failed to register raw input devices [ERROR: %lu]", GetLastError());

        HideCursor();

        CORE.Input.Mouse.cursorLocked = true;
        UpdateCursorClip();
    }
}

// Swap back buffer with front buffer (screen drawing)
void SwapScreenBuffer(void)
{
    if (!platform.hdc) abort();

#if defined(GRAPHICS_API_OPENGL_SOFTWARE)
    // Update framebuffer
    rlCopyFramebuffer(0, 0, CORE.Window.render.width, CORE.Window.render.height, PIXELFORMAT_UNCOMPRESSED_R8G8B8A8, platform.pixels);

    // Force redraw
    InvalidateRect(platform.hwnd, NULL, FALSE);
    UpdateWindow(platform.hwnd);
#else
    if (!SwapBuffers(platform.hdc)) TRACELOG(LOG_ERROR, "WIN32: Failed to swap buffers [ERROR: %lu]", GetLastError());
    if (!ValidateRect(platform.hwnd, NULL)) TRACELOG(LOG_ERROR, "WIN32: Failed to validate screen rect [ERROR: %lu]", GetLastError());
#endif
}

//----------------------------------------------------------------------------------
// Module Functions Definition: Misc
//----------------------------------------------------------------------------------

// Get elapsed time measure in seconds
double GetTime(void)
{
    double time = 0.0;
    LARGE_INTEGER now = { 0 };
    QueryPerformanceCounter(&now);
    time = (double)(now.QuadPart - CORE.Time.base)/(double)platform.timerFrequency.QuadPart;

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
        int len = strlen(url) + 16;
        char *cmd = (char *)RL_CALLOC(len, sizeof(char));
        snprintf(cmd, len, "explorer \"%s\"", url);

        int result = system(cmd);
        if (result == -1) TRACELOG(LOG_WARNING, "OpenURL() child process could not be created");
        RL_FREE(cmd);
    }
}

//----------------------------------------------------------------------------------
// Module Functions Definition: Inputs
//----------------------------------------------------------------------------------

// Set internal gamepad mappings
// NOTE: Gamepads are read with XInput, it has a fixed buttons and axes layout
int SetGamepadMappings(const char *mappings)
{
    TRACELOG(LOG_WARNING, "WIN32: Gamepad mappings are not used by XInput");

    return -1;
}

// Set gamepad vibration
void SetGamepadVibration(int gamepad, float leftMotor, float rightMotor, float duration)
{
    if ((gamepad < 0) || (gamepad >= XINPUT_MAX_GAMEPADS) || (gamepad >= MAX_GAMEPADS) || !CORE.Input.Gamepad.ready[gamepad]) return;
    if ((pXInputSetState == NULL) || (duration <= 0.0f)) return;

    if (leftMotor < 0.0f) leftMotor = 0.0f;
    if (leftMotor > 1.0f) leftMotor = 1.0f;
    if (rightMotor < 0.0f) rightMotor = 0.0f;
    if (rightMotor > 1.0f) rightMotor = 1.0f;
    if (duration > MAX_GAMEPAD_VIBRATION_TIME) duration = MAX_GAMEPAD_VIBRATION_TIME;

    XInputVibration vibration = { (WORD)(leftMotor*65535.0f), (WORD)(rightMotor*65535.0f) };
    if (pXInputSetState((DWORD)gamepad, &vibration) == ERROR_SUCCESS) platform.gamepadVibrationEnd[gamepad] = GetTime() + (double)duration;
}

// Set mouse position XY
void SetMousePosition(int x, int y)
{
    if (!CORE.Input.Mouse.cursorLocked)
    {
        // Set the cursor position relative to the top-left corner of the window's client area
        POINT topLeft = { 0 };

        // Skip if window is minimized or if ClientToScreen calculation fails
        if (!IsIconic(platform.hwnd) && ClientToScreen(platform.hwnd, &topLeft))
        {
            Vector2 scale = CORE.Input.Mouse.scale;
            int xScaled = (int)roundf(x/scale.x);
            int yScaled = (int)roundf(y/scale.y);
            if (SetCursorPos(topLeft.x + xScaled, topLeft.y + yScaled))
            {
                CORE.Input.Mouse.currentPosition = (Vector2){ (float)xScaled, (float)yScaled };
            }
        }
    }
    else TRACELOG(LOG_WARNING, "INPUT: MOUSE: Unable to set mouse position. The cursor is locked");
}

// Set mouse cursor
void SetMouseCursor(int cursor)
{
    LPCWSTR cursorName = GetCursorName(cursor);
    HCURSOR hcursor = LoadCursorW(NULL, cursorName);
    if (!hcursor) TRACELOG(LOG_ERROR, "WIN32: Failed to load requested cursor [ERROR: %lu]", GetLastError());
    else platform.cursor = hcursor;

    // NOTE: Cursor is set by the window when the mouse moves over it
    if (!CORE.Input.Mouse.cursorHidden && CORE.Input.Mouse.cursorOnScreen) SetCursor(platform.cursor);
}

// Get physical key name
// NOTE: Character produced by the key with the current keyboard layout, only for keys producing characters
const char *GetKeyName(int key)
{
    static char name[8] = { 0 };

    bool printable = (((key >= KEY_APOSTROPHE) && (key <= KEY_GRAVE)) || ((key >= KEY_KP_0) && (key <= KEY_KP_ADD)) || (key == KEY_KP_EQUAL));
    if (!printable || (key == KEY_SPACE)) return NULL;

    UINT scancode = 0;
    UINT vk = 0;

    if ((key >= KEY_KP_0) && (key <= KEY_KP_ADD))
    {
        // NOTE: Keypad keys depend on NUMLOCK state, their virtual keys are used directly
        static const UINT keypadKeys[] = { VK_NUMPAD0, VK_NUMPAD1, VK_NUMPAD2, VK_NUMPAD3, VK_NUMPAD4, VK_NUMPAD5, VK_NUMPAD6, VK_NUMPAD7,
            VK_NUMPAD8, VK_NUMPAD9, VK_DECIMAL, VK_DIVIDE, VK_MULTIPLY, VK_SUBTRACT, VK_ADD };

        vk = keypadKeys[key - KEY_KP_0];
        scancode = MapVirtualKeyW(vk, MAPVK_VK_TO_VSC);
    }
    else
    {
        for (unsigned int i = 0; i < sizeof(keyboardKeyFromScancode)/sizeof(keyboardKeyFromScancode[0]); i++)
        {
            if (keyboardKeyFromScancode[i] == key)
            {
                scancode = i;
                break;
            }
        }

        if (scancode != 0) vk = MapVirtualKeyW(scancode, MAPVK_VSC_TO_VK);
    }

    if (vk == 0) return NULL;

    BYTE state[256] = { 0 };

    WCHAR chars[16] = { 0 };
    int length = ToUnicode(vk, scancode, state, chars, sizeof(chars)/sizeof(chars[0]), 0);

    // NOTE: Dead keys are stored by the keyboard layout, translated again to get the character and clear that state
    if (length == -1) length = ToUnicode(vk, scancode, state, chars, sizeof(chars)/sizeof(chars[0]), 0);
    if (length < 1) return NULL;

    int size = WideCharToMultiByte(CP_UTF8, 0, chars, 1, name, sizeof(name) - 1, NULL, NULL);
    if (size < 1) return NULL;
    name[size] = '\0';

    return name;
}

// Register all input events
void PollInputEvents(void)
{
    // Reset keys/chars pressed registered
    CORE.Input.Keyboard.keyPressedQueueCount = 0;
    CORE.Input.Keyboard.charPressedQueueCount = 0;

    // Reset key repeats
    for (int i = 0; i < MAX_KEYBOARD_KEYS; i++) CORE.Input.Keyboard.keyRepeatInFrame[i] = 0;

    // Register previous mouse states
    for (int i = 0; i < MAX_MOUSE_BUTTONS; i++) CORE.Input.Mouse.previousButtonState[i] = CORE.Input.Mouse.currentButtonState[i];

    // Reset last gamepad button/axis registered state
    CORE.Input.Gamepad.lastButtonPressed = 0; // GAMEPAD_BUTTON_UNKNOWN
    //CORE.Input.Gamepad.axisCount = 0;

    // Register previous touch states
    for (int i = 0; i < MAX_TOUCH_POINTS; i++) CORE.Input.Touch.previousTouchState[i] = CORE.Input.Touch.currentTouchState[i];

    // Reset touch positions
    // TODO: It resets on target platform the mouse position and not filled again until a move-event,
    // so, if mouse is not moved it returns a (0, 0) position... this behaviour should be reviewed!
    //for (int i = 0; i < MAX_TOUCH_POINTS; i++) CORE.Input.Touch.position[i] = (Vector2){ 0, 0 };

    memcpy(CORE.Input.Keyboard.previousKeyState, CORE.Input.Keyboard.currentKeyState, sizeof(CORE.Input.Keyboard.previousKeyState));

    // Apply releases of keys and mouse buttons that were pressed and released in the previous frame
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

    // Register previous mouse wheel state
    CORE.Input.Mouse.previousWheelMove = CORE.Input.Mouse.currentWheelMove;
    CORE.Input.Mouse.currentWheelMove = (Vector2){ 0.0f, 0.0f };

    // Register previous mouse position
    CORE.Input.Mouse.previousPosition = CORE.Input.Mouse.currentPosition;

    CORE.Window.resizedLastFrame = false;
    CORE.Window.shouldClose = false;

    // Wait for events if required, also while minimized, a release applied in this frame is an event already
    bool waitEvents = (CORE.Window.eventWaiting || ((CORE.Window.flags & FLAG_WINDOW_MINIMIZED) && !(CORE.Window.flags & FLAG_WINDOW_ALWAYS_RUN)));
    if (waitEvents && !releasePending)
    {
        MSG next = { 0 };
        if (!PeekMessageW(&next, NULL, 0, 0, PM_NOREMOVE)) WaitMessage();
    }

    // Process windows messages
    MSG msg = { 0 };
    while (PeekMessageW(&msg, NULL, 0, 0, PM_REMOVE))
    {
        TranslateMessage(&msg);
        DispatchMessageW(&msg);
    }

    if ((platform.previousRenderWidth != CORE.Window.render.width) || (platform.previousRenderHeight != CORE.Window.render.height))
    {
        SetupViewport(CORE.Window.render.width, CORE.Window.render.height);
#if defined(GRAPHICS_API_OPENGL_SOFTWARE)
        if (rlGetVersion() == RL_OPENGL_SOFTWARE)
        {
            BITMAPINFO bmi = { 0 };
            ZeroMemory(&bmi, sizeof(bmi));
            bmi.bmiHeader.biSize = sizeof(BITMAPINFOHEADER);
            bmi.bmiHeader.biWidth = (LONG)CORE.Window.render.width;
            bmi.bmiHeader.biHeight = -(LONG)(CORE.Window.render.height); // Top-down bitmap
            bmi.bmiHeader.biPlanes = 1;
            bmi.bmiHeader.biBitCount  = 32;         // 32-bit BGRA
            bmi.bmiHeader.biCompression = BI_RGB;

            platform.hbitmap = CreateDIBSection(
                platform.hdcmem, &bmi, DIB_RGB_COLORS,
                (void **)&platform.pixels, NULL, 0);

            HGDIOBJ oldBitmap = SelectObject(platform.hdcmem, platform.hbitmap);
            DeleteObject(oldBitmap);

            swResize(CORE.Window.render.width, CORE.Window.render.height);
        }
#endif
        platform.previousRenderWidth = CORE.Window.render.width;
        platform.previousRenderHeight = CORE.Window.render.height;
    }

    PollGamepads();

    // Windows will miss key up events for specific keys; check them manually
    // NOTE: A key pressed in this frame is kept down for the frame
    if (((GetKeyState(VK_LSHIFT) & 0x8000) == 0) && !platform.keyPressedInFrame[KEY_LEFT_SHIFT]) CORE.Input.Keyboard.currentKeyState[KEY_LEFT_SHIFT] = 0;
    if (((GetKeyState(VK_RSHIFT) & 0x8000) == 0) && !platform.keyPressedInFrame[KEY_RIGHT_SHIFT]) CORE.Input.Keyboard.currentKeyState[KEY_RIGHT_SHIFT] = 0;
    if (((GetKeyState(VK_LWIN) & 0x8000) == 0) && !platform.keyPressedInFrame[KEY_LEFT_SUPER]) CORE.Input.Keyboard.currentKeyState[KEY_LEFT_SUPER] = 0;
    if (((GetKeyState(VK_RWIN) & 0x8000) == 0) && !platform.keyPressedInFrame[KEY_RIGHT_SUPER]) CORE.Input.Keyboard.currentKeyState[KEY_RIGHT_SUPER] = 0;

    // Windows won't send a WM_KEYDOWN event for print screen; check it manually
    bool isPrintScreenDown = ((GetAsyncKeyState(VK_SNAPSHOT) & 0x8000) != 0) && (GetForegroundWindow() == platform.hwnd);
    CORE.Input.Keyboard.currentKeyState[KEY_PRINT_SCREEN] = isPrintScreenDown;
    if (IsKeyPressed(KEY_PRINT_SCREEN) && (CORE.Input.Keyboard.exitKey == KEY_PRINT_SCREEN)) CORE.Window.shouldClose = true;
}

//----------------------------------------------------------------------------------
// Module Internal Functions Definition
//----------------------------------------------------------------------------------
// Initialize modern OpenGL context
// NOTE: Creating a dummy context first to query required extensions
HGLRC InitOpenGL(HWND hwnd, HDC hdc)
{
    // Create an invisible window for initial dummy OpenGL context
    WNDCLASSEXW dummyWindowClass = {
        .cbSize = sizeof(WNDCLASSEXW),
        .style = CS_HREDRAW | CS_VREDRAW | CS_OWNDC,
        .lpfnWndProc = DefWindowProcW,
        .hInstance = GetModuleHandleW(NULL),
        .lpszClassName = L"wglDummy"
    };
    RegisterClassExW(&dummyWindowClass);

    HWND dummyWindow = CreateWindowExW(
        0, dummyWindowClass.lpszClassName, L"dummyWindow",
        0, 0, 0, 1, 1, NULL, NULL, NULL, NULL);

    HDC dummyHDC = GetDC(dummyWindow);

    // Create a dummy OpenGL context to get WGL extensions
    PIXELFORMATDESCRIPTOR pixelFormatDesc = {
        .nSize = sizeof(PIXELFORMATDESCRIPTOR),
        .nVersion = 1,
        .dwFlags = PFD_DRAW_TO_WINDOW | PFD_SUPPORT_OPENGL | PFD_DOUBLEBUFFER,
        .iPixelType = PFD_TYPE_RGBA,
        .cColorBits = 32,
        .cAlphaBits = 8,
        .cDepthBits = 24,
        .cStencilBits = 8,
        .iLayerType = PFD_MAIN_PLANE
    };

    int dummyPixelFormat = ChoosePixelFormat(dummyHDC, &pixelFormatDesc);
    SetPixelFormat(dummyHDC, dummyPixelFormat, &pixelFormatDesc);
    //int pixelFormat = ChoosePixelFormat(platform.hdc, &pixelFormatDesc);
    //if (!pixelFormat) { TRACELOG(LOG_ERROR, "%s failed, error=%lu", "ChoosePixelFormat", GetLastError()); return -1; }
    //if (!SetPixelFormat(platform.hdc, pixelFormat, &pixelFormatDesc)) { TRACELOG(LOG_ERROR, "%s failed, error=%lu", "SetPixelFormat", GetLastError()); return -1; }

    HGLRC tempContext = wglCreateContext(dummyHDC);
    //if (!tempContext) { TRACELOG(LOG_ERROR, "%s failed, error=%lu", "wglCreateContext", GetLastError()); return -1; }
    if (!wglMakeCurrent(dummyHDC, tempContext)) TRACELOG(LOG_WARNING, "WIN32: Failed to make the dummy OpenGL context current [ERROR: %lu]", GetLastError());

    // Load WGL extension entry points
    wglCreateContextAttribsARB = (PFNWGLCREATECONTEXTATTRIBSARBPROC)wglGetProcAddress("wglCreateContextAttribsARB");
    wglChoosePixelFormatARB = (PFNWGLCHOOSEPIXELFORMATARBPROC)wglGetProcAddress("wglChoosePixelFormatARB");
    wglSwapIntervalEXT = (PFNWGLSWAPINTERVALEXTPROC)wglGetProcAddress("wglSwapIntervalEXT");
    wglGetExtensionsStringARB = (PFNWGLGETEXTENSIONSSTRINGARBPROC)wglGetProcAddress("wglGetExtensionsStringARB");

    // Setup modern pixel format if extension is available
    BOOL setPixelFormatSucceeded = FALSE;
    if (wglChoosePixelFormatARB)
    {
        bool useMSAA = FLAG_IS_SET(CORE.Window.flags, FLAG_MSAA_4X_HINT);

        int pixelFormatAttribs[] = {
            WGL_ACCELERATION_ARB, WGL_FULL_ACCELERATION_ARB,
            WGL_DRAW_TO_WINDOW_ARB, GL_TRUE,
            WGL_SUPPORT_OPENGL_ARB, GL_TRUE,
            WGL_DOUBLE_BUFFER_ARB, GL_TRUE,
            WGL_PIXEL_TYPE_ARB, WGL_TYPE_RGBA_ARB,
            WGL_COLOR_BITS_ARB, 32,
            WGL_RED_BITS_ARB, 8,
            WGL_GREEN_BITS_ARB, 8,
            WGL_BLUE_BITS_ARB, 8,
            WGL_ALPHA_BITS_ARB, 8,
            WGL_DEPTH_BITS_ARB, 24,
            WGL_STENCIL_BITS_ARB, 8,
            WGL_SAMPLES_ARB, (useMSAA)? 4 : 0,
            WGL_SAMPLE_BUFFERS_ARB, (int)useMSAA,
            0 // Terminator
        };

        int format = 0;
        UINT numFormats = 0;
        if (wglChoosePixelFormatARB(hdc, pixelFormatAttribs, NULL, 1, &format, &numFormats) && (numFormats > 0))
        {
            PIXELFORMATDESCRIPTOR newPixelFormatDescriptor = { 0 };
            DescribePixelFormat(hdc, format, sizeof(newPixelFormatDescriptor), &newPixelFormatDescriptor);
            setPixelFormatSucceeded = SetPixelFormat(hdc, format, &newPixelFormatDescriptor);
        }
    }

    if (!setPixelFormatSucceeded)
    {
        // If modern pixel format selection failed, fall back to old SetPixelFormat
        int pixelFormat = ChoosePixelFormat(hdc, &pixelFormatDesc);
        setPixelFormatSucceeded = SetPixelFormat(hdc, pixelFormat, &pixelFormatDesc);

        if (setPixelFormatSucceeded) TRACELOG(LOG_WARNING, "WGL: Using legacy pixel format, MSAA not available");
        else TRACELOG(LOG_ERROR, "WGL: Unable to find a suitable pixel format with this graphics driver");
    }

    // Create real modern OpenGL context (3.3 core)
    HGLRC realContext = NULL;
    if (wglCreateContextAttribsARB)
    {
        int glContextVersionMajor = 1;
        int glContextVersionMinor = 1;
        int glContextProfile = WGL_CONTEXT_CORE_PROFILE_BIT_ARB;

        if (rlGetVersion() == RL_OPENGL_21)         // Request OpenGL 2.1 context
        {
            glContextVersionMajor = 2;
            glContextVersionMinor = 1;
        }
        else if (rlGetVersion() == RL_OPENGL_33)    // Request OpenGL 3.3 context
        {
            glContextVersionMajor = 3;
            glContextVersionMinor = 3;
        }
        else if (rlGetVersion() == RL_OPENGL_43)    // Request OpenGL 4.3 context
        {
            glContextVersionMajor = 4;
            glContextVersionMinor = 3;
        }
        else if (rlGetVersion() == RL_OPENGL_ES_20) // Request OpenGL ES 2.0 context
        {
            if (IsWglExtensionAvailable(platform.hdc, "WGL_EXT_create_context_es_profile") ||
                IsWglExtensionAvailable(platform.hdc, "WGL_EXT_create_context_es2_profile"))
            {
                glContextVersionMajor = 2;
                glContextVersionMinor = 0;
                glContextProfile = WGL_CONTEXT_ES_PROFILE_BIT_EXT;
            }
            else TRACELOG(LOG_WARNING, "GL: OpenGL ES context not supported by GPU");
        }
        else if (rlGetVersion() == RL_OPENGL_ES_30) // Request OpenGL ES 3.0 context
        {
            if (IsWglExtensionAvailable(platform.hdc, "WGL_EXT_create_context_es_profile") ||
                IsWglExtensionAvailable(platform.hdc, "WGL_EXT_create_context_es2_profile"))
            {
                glContextVersionMajor = 3;
                glContextVersionMinor = 0;
                glContextProfile = WGL_CONTEXT_ES_PROFILE_BIT_EXT;
            }
            else TRACELOG(LOG_WARNING, "GL: OpenGL ES context not supported by GPU");
        }

        int contextAttribs[] = {
            WGL_CONTEXT_MAJOR_VERSION_ARB, glContextVersionMajor,
            WGL_CONTEXT_MINOR_VERSION_ARB, glContextVersionMinor,
            WGL_CONTEXT_PROFILE_MASK_ARB, glContextProfile, // WGL_CONTEXT_COMPATIBILITY_PROFILE_BIT_ARB, WGL_CONTEXT_ES_PROFILE_BIT_EXT (if supported)
            //WGL_CONTEXT_FLAGS_ARB, WGL_CONTEXT_FORWARD_COMPATIBLE_BIT_ARB | WGL_CONTEXT_DEBUG_BIT_ARB [glDebugMessageCallback()]
            0 // Terminator
        };

        // NOTE: Not sharing context resources so, second parameters is NULL
        realContext = wglCreateContextAttribsARB(hdc, NULL, contextAttribs);

        // Check for error context creation errors
        // ERROR_INVALID_VERSION_ARB (0x2095)
        // ERROR_INVALID_PROFILE_ARB (0x2096)
        if (realContext == NULL) TRACELOG(LOG_ERROR, "GL: Error creating requested context: %lu", GetLastError());
    }
    else
    {
        // If modern context extension isn't available, fall back to classic context creation
        realContext = wglCreateContext(hdc);
    }

    // Cleanup dummy temp context and window
    wglMakeCurrent(NULL, NULL);
    wglDeleteContext(tempContext);
    ReleaseDC(dummyWindow, dummyHDC);
    DestroyWindow(dummyWindow);
    UnregisterClassW(dummyWindowClass.lpszClassName, dummyWindowClass.hInstance);

    // Activate real context
    if (realContext) wglMakeCurrent(hdc, realContext);

    // Once a real modern OpenGL context is created,
    // required extensions can be loaded (function pointers)
    rlLoadExtensions(WglGetProcAddress);

    return realContext;
}

// Initialize platform: graphics, inputs and more
int InitPlatform(void)
{
    int result = 0;

    platform.appScreenWidth = CORE.Window.screen.width;
    platform.appScreenHeight = CORE.Window.screen.height;
    platform.coverRestoreUndecorated = ((CORE.Window.flags & FLAG_WINDOW_UNDECORATED) != 0);
    platform.coverRestoreMaximized = ((CORE.Window.flags & FLAG_WINDOW_MAXIMIZED) != 0);
    platform.coverEnterFlags = 0;
    platform.coverLeave = false;
    platform.opacity = 1.0f;
    platform.cursor = LoadCursorW(NULL, (LPCWSTR)IDC_ARROW);
    platform.mouseTracked = false;
    platform.rawAbsoluteValid = false;
    platform.highSurrogate = 0;
    platform.gamepadCheckTime = 0.0;
    for (int i = 0; i < MAX_TOUCH_POINTS; i++) platform.touchIds[i] = -1;
    for (int i = 0; i < MAX_GAMEPADS; i++) platform.gamepadVibrationEnd[i] = 0.0;
    memset(platform.keyPressedInFrame, 0, sizeof(platform.keyPressedInFrame));
    memset(platform.keyReleasePending, 0, sizeof(platform.keyReleasePending));
    memset(platform.buttonPressedInFrame, 0, sizeof(platform.buttonPressedInFrame));
    memset(platform.buttonReleasePending, 0, sizeof(platform.buttonReleasePending));
    platform.desiredFlags = SanitizeFlags(CORE.Window.flags);

    // NOTE: From this point CORE.Window.flags should always reflect the actual state of the window
    CORE.Window.flags = FLAG_WINDOW_HIDDEN | (platform.desiredFlags & FLAG_MASK_NO_UPDATE);
    CORE.Window.screenMax.width = 9999;
    CORE.Window.screenMax.height = 9999;

    LoadWin32FunctionPointers();

    SetProcessDpiAwareCompat();

    HINSTANCE hInstance = GetModuleHandleW(0);

    // Define window class
    WNDCLASSEXW windowClass = {
        .cbSize = sizeof(WNDCLASSEXW),
        .style = CS_HREDRAW | CS_VREDRAW | CS_OWNDC,
        .lpfnWndProc = WndProc,                         // Custom procedure assigned
        .cbWndExtra = sizeof(LONG_PTR),                 // extra space for the Tuple object ptr
        .hInstance = hInstance,
        .hCursor = LoadCursorW(NULL, (LPCWSTR)IDC_ARROW), // TODO: Check if this is really required, since WM_SETCURSOR event is processed
        .lpszClassName = CLASS_NAME                     // Class name: L"raylibWindow"
    };

    // Load user-provided icon if available
    // NOTE: raylib resource file defaults to GLFW_ICON id, so looking for same identifier
    windowClass.hIcon = LoadImageW(hInstance, L"GLFW_ICON", IMAGE_ICON, 0, 0, LR_DEFAULTSIZE | LR_SHARED);
    if (!windowClass.hIcon) windowClass.hIcon = LoadImageW(NULL, (LPCWSTR)IDI_APPLICATION, IMAGE_ICON, 0, 0, LR_DEFAULTSIZE | LR_SHARED);

    // Register window class
    result = (int)RegisterClassExW(&windowClass);
    if (result == 0) TRACELOG(LOG_ERROR, "WIN32: Failed to register window class [ERROR: %lu]", GetLastError());

    // Get primary monitor info
    POINT primaryTopLeft = { 0 };
    HMONITOR monitor = MonitorFromPoint(primaryTopLeft, MONITOR_DEFAULTTOPRIMARY);
    if (monitor != NULL)
    {
        MONITORINFO info = { 0 };
        info.cbSize = sizeof(info);
        result = (int)GetMonitorInfoW(monitor, &info);

        if (result == 0) TRACELOG(LOG_WARNING, "WIN32: DISPLAY: Failed to get monitor info [ERROR: %u]", GetLastError());
        else
        {
            CORE.Window.display.width = info.rcMonitor.right - info.rcMonitor.left;
            CORE.Window.display.height = info.rcMonitor.bottom - info.rcMonitor.top;
        }
    }
    else TRACELOG(LOG_WARNING, "WIN32: DISPLAY: Failed to get primary monitor from point [ERROR: %u]", GetLastError());

    // Create window
    // NOTE: Title string needs to be converted to WCHAR
    WCHAR *titleWide = NULL;
    A_TO_W_ALLOCA(titleWide, CORE.Window.title);

    // Create window and get handle
    platform.hwnd = CreateWindowExW(
        WINDOW_STYLE_EX,
        CLASS_NAME,
        titleWide,
        MakeWindowStyle(CORE.Window.flags | FLAG_WINDOW_HIDDEN),
        CW_USEDEFAULT, CW_USEDEFAULT,
        CW_USEDEFAULT, CW_USEDEFAULT, // Window size/position will be overridden after creation
        NULL, NULL,
        GetModuleHandleW(NULL), NULL);

    if (!platform.hwnd)
    {
        TRACELOG(LOG_ERROR, "WIN32: WINDOW: Failed to create window [ERROR: %lu]", GetLastError());
        return -1;
    }

    // Get handle to device drawing context
    // NOTE: Windows GDI object that represents a drawing surface
    platform.hdc = GetDC(platform.hwnd);

    // Accept files dropped on the window
    DragAcceptFiles(platform.hwnd, TRUE);

    bool highDPIWindowFlag = platform.desiredFlags & FLAG_WINDOW_HIGHDPI;
    UINT dpi = GetDpiForWindowCompat(platform.hwnd);
    float dpiScale = dpi/96.0f;

    if (rlGetVersion() != RL_OPENGL_SOFTWARE)
    {
        // Init hardware-accelerated OpenGL modern context
        platform.glContext = InitOpenGL(platform.hwnd, platform.hdc);
    }

    CORE.Window.ready = true;

    // Window client size, considering high-dpi
    int clientWidth = platform.appScreenWidth;
    int clientHeight = platform.appScreenHeight;
    if (highDPIWindowFlag)
    {
        clientWidth = (int)roundf(platform.appScreenWidth*dpiScale);
        clientHeight = (int)roundf(platform.appScreenHeight*dpiScale);
    }

    // Window size with the requested decorations, fullscreen and borderless windowed modes are applied after the window is shown
    unsigned floatingFlags = platform.desiredFlags & ~(FLAG_MASK_COVER | FLAG_WINDOW_UNDECORATED | FLAG_WINDOW_MINIMIZED | FLAG_WINDOW_MAXIMIZED);
    if (platform.coverRestoreUndecorated) floatingFlags |= FLAG_WINDOW_UNDECORATED;

    RECT rect = { 0, 0, clientWidth, clientHeight };
    AdjustWindowRectExForDpiCompat(&rect, MakeWindowStyle(floatingFlags), FALSE, WINDOW_STYLE_EX, dpi);

    int windowWidth = rect.right - rect.left;
    int windowHeight = rect.bottom - rect.top;

    HMONITOR currentMonitor = MonitorFromWindow(platform.hwnd, MONITOR_DEFAULTTONEAREST);
    MONITORINFO monitorInfo = { 0 };
    monitorInfo.cbSize = sizeof(monitorInfo);
    if (GetMonitorInfoW(currentMonitor, &monitorInfo))
    {
        // rcWork is the monitor area, not including the taskbar
        RECT rcWork = monitorInfo.rcWork;
        int monitorWidth = rcWork.right - rcWork.left;
        int monitorHeight = rcWork.bottom - rcWork.top;
        int x = rcWork.left + (monitorWidth - windowWidth)/2;
        int y = rcWork.top + (monitorHeight - windowHeight)/2;

        // Ensure the window doesn't go off-screen; otherwise, the user might not be able to move it
        if (x < rcWork.left) x = rcWork.left;
        if (y < rcWork.top) y = rcWork.top;

        // Set the client area position and DPI-adjusted size with the current style, it is kept when the decorations are applied
        RECT windowRect = { x - rect.left, y - rect.top, x - rect.left + clientWidth, y - rect.top + clientHeight };
        AdjustWindowRectExForDpiCompat(&windowRect, MakeWindowStyle(CORE.Window.flags | FLAG_WINDOW_HIDDEN), FALSE, WINDOW_STYLE_EX, dpi);

        SetWindowPos(platform.hwnd, NULL, windowRect.left, windowRect.top, windowRect.right - windowRect.left,
            windowRect.bottom - windowRect.top, SWP_NOZORDER | SWP_NOACTIVATE);
    }
    else TRACELOG(LOG_WARNING, "DISPLAY: Unable to retrieve information about monitor");

    if ((platform.desiredFlags & FLAG_WINDOW_HIDDEN) == 0)
    {
        // Show window and taskbar icon, activate unless FLAG_WINDOW_UNFOCUSED is set
        bool isUnfocused = ((platform.desiredFlags & FLAG_WINDOW_UNFOCUSED) != 0);
        if (isUnfocused) CORE.Window.flags |= FLAG_WINDOW_UNFOCUSED;
        int nCmdShow = (isUnfocused)? SW_SHOWNA : SW_SHOWDEFAULT;
        ShowWindow(platform.hwnd, nCmdShow);
    }

    // Update flags (in case of deferred state change required)
    UpdateFlags(platform.hwnd, platform.desiredFlags, platform.appScreenWidth, platform.appScreenHeight);

    // Initial mouse position, first mouse movement is relative to it
    POINT cursorPosition = { 0 };
    if (GetCursorPos(&cursorPosition) && ScreenToClient(platform.hwnd, &cursorPosition))
    {
        CORE.Input.Mouse.currentPosition = (Vector2){ (float)cursorPosition.x, (float)cursorPosition.y };
        CORE.Input.Mouse.previousPosition = CORE.Input.Mouse.currentPosition;
        platform.textInputPosition = cursorPosition;
    }

    CORE.Window.resizedLastFrame = false;
    TRACELOG(LOG_INFO, "DISPLAY: Device initialized successfully %s",
        FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_HIGHDPI)? "(HighDPI)" : "");
    TRACELOG(LOG_INFO, "    > Display size: %i x %i", CORE.Window.display.width, CORE.Window.display.height);
    TRACELOG(LOG_INFO, "    > Screen size:  %i x %i", CORE.Window.screen.width, CORE.Window.screen.height);
    TRACELOG(LOG_INFO, "    > Render size:  %i x %i", CORE.Window.render.width, CORE.Window.render.height);
    TRACELOG(LOG_INFO, "    > Viewport offsets: %i, %i", CORE.Window.renderOffset.x, CORE.Window.renderOffset.y);

    if (rlGetVersion() == RL_OPENGL_SOFTWARE) // Using software renderer
    {
        TRACELOG(LOG_INFO, "GL: OpenGL device information:");
        TRACELOG(LOG_INFO, "    > Vendor:   %s", glGetString(GL_VENDOR));
        TRACELOG(LOG_INFO, "    > Renderer: %s", glGetString(GL_RENDERER));
        TRACELOG(LOG_INFO, "    > Version:  %s", glGetString(GL_VERSION));
        TRACELOG(LOG_INFO, "    > GLSL:     %s", "NOT SUPPORTED");

        // Initialize software framebuffer
        BITMAPINFO bmi = { 0 };
        ZeroMemory(&bmi, sizeof(bmi));
        bmi.bmiHeader.biSize = sizeof(BITMAPINFOHEADER);
        bmi.bmiHeader.biWidth = (LONG)CORE.Window.render.width;
        bmi.bmiHeader.biHeight = -(LONG)(CORE.Window.render.height); // Top-down bitmap
        bmi.bmiHeader.biPlanes = 1;
        bmi.bmiHeader.biBitCount  = 32;         // 32-bit BGRA
        bmi.bmiHeader.biCompression = BI_RGB;

        platform.hdcmem = CreateCompatibleDC(platform.hdc);

        platform.hbitmap = CreateDIBSection(
            platform.hdcmem, &bmi, DIB_RGB_COLORS,
            (void **)&platform.pixels, NULL, 0);

        SelectObject(platform.hdcmem, platform.hbitmap);
        //ReleaseDC(platform.hwnd, platform.hdc); // Not required because the window class uses CS_OWNDC
    }

    // Initialize timing system
    //----------------------------------------------------------------------------
    LARGE_INTEGER time = { 0 };
    QueryPerformanceCounter(&time);
    QueryPerformanceFrequency(&platform.timerFrequency);
    CORE.Time.base = time.QuadPart;

    InitTimer();
    //----------------------------------------------------------------------------

    // Gamepads connected at start are available in the first frame
    PollGamepads();

    // Initialize storage system
    //----------------------------------------------------------------------------
    CORE.Storage.basePath = GetWorkingDirectory();
    //----------------------------------------------------------------------------

    TRACELOG(LOG_INFO, "PLATFORM: DESKTOP: WIN32: Initialized successfully");

    return 0;
}

// Close platform
void ClosePlatform(void)
{
    // Clear required platform data
    if (platform.hwnd)
    {
        // Disable cursor lock (if enabled) and unregister raw input devices
        EnableCursor();

        BOOL result = DestroyWindow(platform.hwnd);
        if (!result) TRACELOG(LOG_WARNING, "WIN32: WINDOW: Failed on window destroy [ERROR: %u]", GetLastError());

        HINSTANCE hInstance = GetModuleHandleW(0);
        result = UnregisterClassW(CLASS_NAME, hInstance);
        if (!result) TRACELOG(LOG_WARNING, "WIN32: WINDOW: Failed to unregister window class [ERROR: %u]", GetLastError());

        DeletePlatformIconHandles();

        platform.hwnd = NULL;
    }

    // Clear clipboard
    if (platform.clipboardText != NULL)
    {
        RL_FREE(platform.clipboardText);
        platform.clipboardText = NULL;
    }

    // Stop gamepads vibration and unload XInput library
    if (xinput != NULL)
    {
        for (int i = 0; (i < XINPUT_MAX_GAMEPADS) && (i < MAX_GAMEPADS); i++)
        {
            if ((platform.gamepadVibrationEnd[i] > 0.0) && (pXInputSetState != NULL))
            {
                XInputVibration vibration = { 0 };
                pXInputSetState((DWORD)i, &vibration);
            }
        }

        FreeLibrary(xinput);
        xinput = NULL;
        pXInputGetState = NULL;
        pXInputSetState = NULL;
    }

    // Unload imm32.dll
    if (imm32 != NULL)
    {
        FreeLibrary(imm32);
        imm32 = NULL;
        pImmGetContext = NULL;
        pImmReleaseContext = NULL;
        pImmSetCompositionWindow = NULL;
        pImmSetCandidateWindow = NULL;
    }

    // Unload shcore.dll
    if (shcore != NULL)
    {
        FreeLibrary(shcore);
        shcore = NULL;
    }

#if SUPPORT_WINMM_HIGHRES_TIMER && !SUPPORT_BUSY_WAIT_LOOP
    timeEndPeriod(1); // Restore time period
#endif
}

// Window procedure, message processing callback
// NOTE: All window event messages are processed here
static LRESULT CALLBACK WndProc(HWND hwnd, UINT msg, WPARAM wparam, LPARAM lparam)
{
    LRESULT result = 0;

    // Sanity check
    DWORD mask = STYLE_MASK_ALL;
    if (platform.hwnd == hwnd)
    {
        if (msg == WM_WINDOWPOSCHANGING) mask &= ~(WS_MINIMIZE | WS_MAXIMIZE);
        CheckFlags("WndProc", hwnd, CORE.Window.flags, MakeWindowStyle(CORE.Window.flags), mask);
    }

    FlagsOp flagsOp = { 0 };
    FlagsOp *deferredFlags = &flagsOp;

    // Message processing
    //------------------------------------------------------------------------------------
    switch (msg)
    {
        case WM_CREATE:
        {
            // WARNING: Not recommended to do OpenGL intialization at this point

        } break;
        //case WM_ACTIVATE
        case WM_DESTROY:
        {
            // Clean up for window destruction
            if (rlGetVersion() == RL_OPENGL_SOFTWARE) // Using software renderer
            {
                if (platform.hdcmem)
                {
                    DeleteDC(platform.hdcmem);
                    platform.hdcmem = NULL;
                }

                if (platform.hbitmap)
                {
                    DeleteObject(platform.hbitmap); // Clears platform.pixels data
                    platform.hbitmap = NULL;
                    platform.pixels = NULL; // NOTE: Pointer invalid after DeleteObject()
                }
            }
            else // OpenGL hardware renderer
            {
                wglMakeCurrent(platform.hdc, NULL);
                if (platform.glContext)
                {
                    if (!wglDeleteContext(platform.glContext)) abort();
                    platform.glContext = NULL;
                }
            }

            if (platform.hdc)
            {
                if (!ReleaseDC(hwnd, platform.hdc)) abort();
                platform.hdc = NULL;
            }

            PostQuitMessage(0);

        } break;
        case WM_CLOSE: CORE.Window.shouldClose = true; break; // Window close button [x], ALT+F4
        //case WM_QUIT: // Application closing, not related to window
        case WM_KILLFOCUS:
        {
            CORE.Window.flags |= FLAG_WINDOW_UNFOCUSED;
            memset(CORE.Input.Keyboard.currentKeyState, 0, sizeof(CORE.Input.Keyboard.currentKeyState));
            memset(CORE.Input.Mouse.currentButtonState, 0, sizeof(CORE.Input.Mouse.currentButtonState));
            memset(CORE.Input.Touch.currentTouchState, 0, sizeof(CORE.Input.Touch.currentTouchState));
            memset(platform.keyReleasePending, 0, sizeof(platform.keyReleasePending));
            memset(platform.buttonReleasePending, 0, sizeof(platform.buttonReleasePending));
            for (int i = 0; i < MAX_TOUCH_POINTS; i++) platform.touchIds[i] = -1;
            CORE.Input.Touch.pointCount = 0;
            platform.highSurrogate = 0;
            UpdateCursorClip();
        } break;
        case WM_SETFOCUS:
        {
            CORE.Window.flags &= ~FLAG_WINDOW_UNFOCUSED;
            UpdateCursorClip();
        } break;
        case WM_SIZING: // Sent to a window that the user is resizing
        {
            if (CORE.Window.flags & FLAG_WINDOW_RESIZABLE)
            {
                //HandleWindowResize(hwnd, &platform.appScreenWidth, &platform.appScreenHeight);
            }

            result = TRUE;
        } break;
        case WM_SIZE:
        {
            // WARNING: Don't trust the docs, they say this message can not be obtained if not calling DefWindowProc()
            // in response to WM_WINDOWPOSCHANGED but looks like when a window is created,
            // this message can be obtained without getting WM_WINDOWPOSCHANGED
            HandleWindowResize(hwnd, &platform.appScreenWidth, &platform.appScreenHeight);
            result = 0; // If an application processes WM_SIZE message, it should return zero
        } break;
        case WM_GETMINMAXINFO:
        {
            // NOTE: Window covering the monitor is not limited by the app size limits
            if ((CORE.Window.flags | platform.desiredFlags) & FLAG_MASK_COVER)
            {
                result = DefWindowProcW(hwnd, msg, wparam, lparam);
                break;
            }

            UINT dpi = GetDpiForWindowCompat(hwnd);
            DWORD style = MakeWindowStyle(platform.desiredFlags);
            float clientScale = (CORE.Window.flags & FLAG_WINDOW_HIGHDPI)? ((float)dpi)/96.0f : 1.0f;
            SIZE maxClientSize = { (LONG)roundf(CORE.Window.screenMax.width*clientScale), (LONG)roundf(CORE.Window.screenMax.height*clientScale) };
            SIZE maxWindowSize = CalcWindowSize(dpi, maxClientSize, style);
            SIZE minClientSize = { (LONG)roundf(CORE.Window.screenMin.width*clientScale), (LONG)roundf(CORE.Window.screenMin.height*clientScale) };
            SIZE minWindowSize = CalcWindowSize(dpi, minClientSize, style);

            // NOTE: System maximized size is kept if smaller, it is adjusted by the system for the monitor the window is maximized on
            LPMINMAXINFO lpmmi = (LPMINMAXINFO) lparam;
            if (maxWindowSize.cx < lpmmi->ptMaxSize.x) lpmmi->ptMaxSize.x = maxWindowSize.cx;
            if (maxWindowSize.cy < lpmmi->ptMaxSize.y) lpmmi->ptMaxSize.y = maxWindowSize.cy;
            if (maxWindowSize.cx < lpmmi->ptMaxTrackSize.x) lpmmi->ptMaxTrackSize.x = maxWindowSize.cx;
            if (maxWindowSize.cy < lpmmi->ptMaxTrackSize.y) lpmmi->ptMaxTrackSize.y = maxWindowSize.cy;
            lpmmi->ptMinTrackSize.x = minWindowSize.cx;
            lpmmi->ptMinTrackSize.y = minWindowSize.cy;
        } break;
        case WM_STYLECHANGING:
        {
            if (wparam == GWL_STYLE)
            {
                STYLESTRUCT *ss = (STYLESTRUCT *)lparam;

                if ((platform.coverEnterFlags != 0) || platform.coverLeave)
                {
                    // Fullscreen or borderless windowed mode change, resizable flag is kept
                    if (platform.coverEnterFlags != 0) deferredFlags->set |= platform.coverEnterFlags;
                    else deferredFlags->clear |= FLAG_MASK_COVER;

                    if (DecoratedFromStyle(ss->styleNew)) deferredFlags->clear |= FLAG_WINDOW_UNDECORATED;
                    else deferredFlags->set |= FLAG_WINDOW_UNDECORATED;
                }
                else GetStyleChangeFlagOps(CORE.Window.flags, ss, deferredFlags);

                UINT dpi = GetDpiForWindowCompat(hwnd);
                // Get client size (framebuffer inside the window)
                RECT rect = { 0 };
                GetClientRect(hwnd, &rect);
                SIZE clientSize = { rect.right, rect.bottom };
                SIZE oldSize = CalcWindowSize(dpi, clientSize, ss->styleOld);
                SIZE newSize = CalcWindowSize(dpi, clientSize, ss->styleNew);

                if (oldSize.cx != newSize.cx || oldSize.cy != newSize.cy)
                {
                    TRACELOG(LOG_INFO, "WIN32: WINDOW: Resize from style change [%dx%d] to [%dx%d]", oldSize.cx, oldSize.cy, newSize.cx, newSize.cy);

                    if (CORE.Window.flags & FLAG_WINDOW_MAXIMIZED)
                    {
                        // looks like windows will automatically "unminimize" a window
                        // if a style changes modifies it's size
                        TRACELOG(LOG_INFO, "WIN32: WINDOW: Style change modified window size, removing maximized flag");
                        deferredFlags->clear |= FLAG_WINDOW_MAXIMIZED;
                    }
                }
            }
        } break;
        case WM_WINDOWPOSCHANGING:
        {
            WINDOWPOS *pos = (WINDOWPOS *)lparam;
            if (pos->flags & SWP_SHOWWINDOW) deferredFlags->clear |= FLAG_WINDOW_HIDDEN;
            else if (pos->flags & SWP_HIDEWINDOW) deferredFlags->set |= FLAG_WINDOW_HIDDEN;

            Mized mized = MIZED_NONE;
            bool isIconic = IsIconic(hwnd);
            bool styleMinimized = !!(WS_MINIMIZE & GetWindowLongPtrW(hwnd, GWL_STYLE));
            if (isIconic != styleMinimized) TRACELOG(LOG_WARNING, "WIN32: IsIconic state different from WS_MINIMIZED state");

            if (isIconic) mized = MIZED_MIN;
            else
            {
                WINDOWPLACEMENT placement;
                placement.length = sizeof(placement);
                if (!GetWindowPlacement(hwnd, &placement)) TRACELOG(LOG_ERROR, "WIN32: WINDOW: FAiled to get monitor placement [ERROR: %lu]", GetLastError());

                if (placement.showCmd == SW_SHOWMAXIMIZED) mized = MIZED_MAX;
            }

            switch (mized)
            {
                case MIZED_NONE:
                {
                    deferredFlags->clear |= (FLAG_WINDOW_MINIMIZED | FLAG_WINDOW_MAXIMIZED);
                } break;
                case MIZED_MIN:
                {
                    // !!! NOTE !!! Do not update the maximized/borderless
                    // flags because when hwnd is minimized it temporarily overrides
                    // the maximized state/flag which gets restored on SW_RESTORE
                    deferredFlags->set |= FLAG_WINDOW_MINIMIZED;
                } break;
                case MIZED_MAX:
                {
                    deferredFlags->clear |= FLAG_WINDOW_MINIMIZED;
                    deferredFlags->set |= FLAG_WINDOW_MAXIMIZED;
                } break;
                default: break;
            }
        } break;
        //case WM_MOVE: break;
        case WM_WINDOWPOSCHANGED:
        {
            WINDOWPOS *pos = (WINDOWPOS*)lparam;
            if (!(pos->flags & SWP_NOSIZE)) HandleWindowResize(hwnd, &platform.appScreenWidth, &platform.appScreenHeight);

            DefWindowProc(hwnd, msg, wparam, lparam);
        } break;
        case WM_GETDPISCALEDSIZE:
        {
            SIZE *inoutSize = (SIZE *)lparam;
            UINT newDpi = (UINT)wparam; // TODO: WARNING: Converting from WPARAM = UINT_PTR

            // For the following flag changes, a window resize event should be posted,
            // Keep the pending size when minimized, maximized, or borderless
            if (CORE.Window.flags & FLAG_WINDOW_MINIMIZED) return TRUE;
            if (CORE.Window.flags & FLAG_WINDOW_MAXIMIZED) return TRUE;
            if (CORE.Window.flags & FLAG_MASK_COVER) return TRUE;

            float dpiScale = ((float)newDpi)/96.0f;
            bool dpiScaling = CORE.Window.flags & FLAG_WINDOW_HIGHDPI;
            // Calculate client area in pixels from points
            SIZE desired = {
                .cx = dpiScaling? (int)roundf((float)platform.appScreenWidth*dpiScale) : platform.appScreenWidth,
                .cy = dpiScaling? (int)roundf((float)platform.appScreenHeight*dpiScale) : platform.appScreenHeight
            };
            // Calculate and set window size, including title and window border
            *inoutSize = CalcWindowSize(newDpi, desired, MakeWindowStyle(CORE.Window.flags));

            result = TRUE;
        } break;
        case WM_DPICHANGED:
        {
            // Get current dpi scale factor
            //float scalex = LOWORD(wparam)/96.0f;
            //float scaley = HIWORD(wparam)/96.0f;

            RECT *suggestedRect = (RECT *)lparam;

            // Never set the window size to anything other than the suggested rect here
            // Doing so can cause a window to stutter between monitors when transitioning between them
            if (platform.perMonitorDPIAwareV2)
            {
                BOOL succeeded = SetWindowPos(hwnd, NULL,
                    suggestedRect->left, suggestedRect->top,
                    suggestedRect->right - suggestedRect->left,
                    suggestedRect->bottom - suggestedRect->top,
                    SWP_NOZORDER | SWP_NOACTIVATE);

                if (succeeded == FALSE) TRACELOG(LOG_ERROR, "Failed to set window position [ERROR: %lu]", GetLastError());
            }
            else if ((CORE.Window.flags & FLAG_WINDOW_HIGHDPI) && !(CORE.Window.flags & (FLAG_WINDOW_MINIMIZED | FLAG_WINDOW_MAXIMIZED | FLAG_MASK_COVER)))
            {
                // This code only runs when per-monitor dpi awareness v1 is enabled. v1 doesn't send WM_GETDPISCALEDSIZE and the suggested RECT is subtly wrong (doesn't properly account for the title bar size)
                // The solution is to manually build a properly sized RECT (and ignore the advice of not using the suggested RECT). This makes a per-monitor v1 process match the behavior of per-monitor v2
                // To prevent the potential stutter when dragging the window between monitors, the manually sized rect keeps the same center as the suggested RECT
                RECT windowRect = { 0 };
                RECT clientRect = { 0 };
                GetWindowRect(hwnd, &windowRect);
                GetClientRect(hwnd, &clientRect);

                float dpiScale = ((float)LOWORD(wparam))/96.0f;
                int width = (int)roundf(platform.appScreenWidth*dpiScale) + (windowRect.right - windowRect.left) - clientRect.right;
                int height = (int)roundf(platform.appScreenHeight*dpiScale) + (windowRect.bottom - windowRect.top) - clientRect.bottom;
                int x = suggestedRect->left + ((suggestedRect->right - suggestedRect->left) - width)/2;
                int y = suggestedRect->top + ((suggestedRect->bottom - suggestedRect->top) - height)/2;

                BOOL succeeded = SetWindowPos(hwnd, NULL, x, y, width, height, SWP_NOZORDER | SWP_NOACTIVATE);

                if (succeeded == FALSE) TRACELOG(LOG_ERROR, "Failed to set window position [ERROR: %lu]", GetLastError());
            }

            HandleWindowResize(hwnd, &platform.appScreenWidth, &platform.appScreenHeight);

        } break;
        case WM_SETCURSOR:
        {
            // Called when mouse moves, enters/leaves window...
            if (LOWORD(lparam) == HTCLIENT)
            {
                SetCursor(CORE.Input.Mouse.cursorHidden? NULL : platform.cursor);
                return 0;
            }

            result = DefWindowProc(hwnd, msg, wparam, lparam);
        } break;
        case WM_PAINT:
        {
            if (rlGetVersion() == RL_OPENGL_SOFTWARE) // Using software renderer
            {
                PAINTSTRUCT ps = { 0 };
                HDC hdc = BeginPaint(hwnd, &ps);

                // Blit from memory DC to window DC
                BitBlt(hdc, 0, 0, CORE.Window.render.width, CORE.Window.render.height, platform.hdcmem, 0, 0, SRCCOPY);

                EndPaint(hwnd, &ps);
            }
            else DefWindowProc(hwnd, msg, wparam, lparam);
        } break;
        case WM_INPUT:
        {
            HandleRawInput(lparam);
            result = DefWindowProcW(hwnd, msg, wparam, lparam);
        } break;
        case WM_MOUSEMOVE:
        {
            if (!platform.mouseTracked)
            {
                // Request a notification when the mouse leaves the window
                TRACKMOUSEEVENT track = { 0 };
                track.cbSize = sizeof(track);
                track.dwFlags = TME_LEAVE;
                track.hwndTrack = hwnd;
                if (TrackMouseEvent(&track)) platform.mouseTracked = true;

                CORE.Input.Mouse.cursorOnScreen = true;
            }

            if (!CORE.Input.Mouse.cursorLocked)
            {
                CORE.Input.Mouse.currentPosition.x = (float)GET_X_LPARAM(lparam);
                CORE.Input.Mouse.currentPosition.y = (float)GET_Y_LPARAM(lparam);
                CORE.Input.Touch.position[0] = CORE.Input.Mouse.currentPosition;
            }
        } break;
        case WM_MOUSELEAVE:
        {
            platform.mouseTracked = false;
            CORE.Input.Mouse.cursorOnScreen = false;
        } break;
        case WM_CHAR:
        {
            WCHAR unit = (WCHAR)wparam;

            if ((unit >= 0xd800) && (unit <= 0xdbff)) platform.highSurrogate = unit;
            else
            {
                unsigned int codepoint = unit;

                // NOTE: Characters outside of the basic multilingual plane are sent in two messages
                if ((unit >= 0xdc00) && (unit <= 0xdfff))
                {
                    if (platform.highSurrogate != 0) codepoint = 0x10000 + (((unsigned int)platform.highSurrogate - 0xd800) << 10) + (unit - 0xdc00);
                    else codepoint = 0;
                }

                platform.highSurrogate = 0;
                HandleChar(codepoint);
            }
        } break;
        case WM_SYSCHAR: break; // Characters of ALT+key combinations are not text input, not processed to avoid the default sound
        case WM_UNICHAR:
        {
            // NOTE: Sent by other applications, TRUE reports that the message is supported
            if (wparam == UNICODE_NOCHAR) result = TRUE;
            else HandleChar((unsigned int)wparam);
        } break;
        case WM_POINTERDOWN:
        case WM_POINTERUPDATE:
        case WM_POINTERUP:
        case WM_POINTERCAPTURECHANGED:
        {
            // NOTE: Touch points are processed here, mouse and pen are processed by the default mouse messages
            UINT32 pointerId = (UINT32)LOWORD(wparam);
            DWORD pointerType = 0;

            if ((pGetPointerType != NULL) && pGetPointerType(pointerId, &pointerType) && (pointerType == POINTER_TYPE_TOUCH))
            {
                POINT position = { GET_X_LPARAM(lparam), GET_Y_LPARAM(lparam) };
                ScreenToClient(hwnd, &position);

                int touchAction = TOUCH_ACTION_UP;
                if (msg == WM_POINTERDOWN) touchAction = TOUCH_ACTION_DOWN;
                else if (msg == WM_POINTERUPDATE) touchAction = TOUCH_ACTION_MOVE;

                HandleTouch(touchAction, (int)pointerId, (float)position.x*CORE.Input.Mouse.scale.x, (float)position.y*CORE.Input.Mouse.scale.y);
            }
            else result = DefWindowProcW(hwnd, msg, wparam, lparam);
        } break;
        case WM_DROPFILES: HandleDroppedFiles((HDROP)wparam); break;
        case WM_IME_STARTCOMPOSITION:
        case WM_IME_COMPOSITION:
        {
            UpdateInputMethodPosition(hwnd);
            result = DefWindowProcW(hwnd, msg, wparam, lparam);
        } break;
        case WM_SYSKEYDOWN:
        {
            HandleKey(wparam, lparam, 1);
            // WM_SYSKEYDOWN events should still call DefWindowProc to allow things like alt+F4 to work
            result = DefWindowProcW(hwnd, msg, wparam, lparam);
        } break;
        case WM_SYSKEYUP:
        {
            HandleKey(wparam, lparam, 0);
            result = DefWindowProcW(hwnd, msg, wparam, lparam);
        } break;
        case WM_KEYDOWN: HandleKey(wparam, lparam, 1); break;
        case WM_KEYUP: HandleKey(wparam, lparam, 0); break;
        case WM_SYSCOMMAND:
        {
            if ((wparam & 0xFFF0) == SC_KEYMENU)
            {
                // Prevent the window from freezing when pressing Alt or F10
                result = 0;
            }
            else
            {
                result = DefWindowProcW(hwnd, msg, wparam, lparam);
            }
        } break;
        case WM_LBUTTONDOWN: HandleMouseButton(MOUSE_BUTTON_LEFT, 1); break;
        case WM_LBUTTONUP  : HandleMouseButton(MOUSE_BUTTON_LEFT, 0); break;
        case WM_RBUTTONDOWN: HandleMouseButton(MOUSE_BUTTON_RIGHT, 1); break;
        case WM_RBUTTONUP  : HandleMouseButton(MOUSE_BUTTON_RIGHT, 0); break;
        case WM_MBUTTONDOWN: HandleMouseButton(MOUSE_BUTTON_MIDDLE, 1); break;
        case WM_MBUTTONUP  : HandleMouseButton(MOUSE_BUTTON_MIDDLE, 0); break;
        case WM_XBUTTONDOWN:
        {
            switch (HIWORD(wparam))
            {
                case XBUTTON1: HandleMouseButton(MOUSE_BUTTON_SIDE, 1); break;
                case XBUTTON2: HandleMouseButton(MOUSE_BUTTON_EXTRA, 1); break;
                default: TRACELOG(LOG_WARNING, "TODO: Handle ex mouse button DOWN [wparam=%u]", HIWORD(wparam)); break;
            }
        } break;
        case WM_XBUTTONUP:
        {
            switch (HIWORD(wparam))
            {
                case XBUTTON1: HandleMouseButton(MOUSE_BUTTON_SIDE, 0); break;
                case XBUTTON2: HandleMouseButton(MOUSE_BUTTON_EXTRA, 0); break;
                default: TRACELOG(LOG_WARNING, "TODO: Handle ex mouse button UP [wparam=%u]", HIWORD(wparam)); break;
            }
        } break;
        case WM_MOUSEWHEEL: CORE.Input.Mouse.currentWheelMove.y += ((float)GET_WHEEL_DELTA_WPARAM(wparam))/WHEEL_DELTA; break;
        // Horizontal wheel movement is inverted to match GLFW
        case WM_MOUSEHWHEEL: CORE.Input.Mouse.currentWheelMove.x -= ((float)GET_WHEEL_DELTA_WPARAM(wparam))/WHEEL_DELTA; break;

        default: result = DefWindowProcW(hwnd, msg, wparam, lparam); // Message passed directly for execution (default behaviour)
    }
    //------------------------------------------------------------------------------------

    // Sanity check for flags
    if (platform.hwnd == hwnd) CheckFlags("After WndProc", hwnd, CORE.Window.flags, MakeWindowStyle(CORE.Window.flags), mask);

    // Operations to execute after the above check
    if (flagsOp.set & flagsOp.clear) TRACELOG(LOG_WARNING, "WIN32: FLAGS: Flags 0x%x were both set and cleared", flagsOp.set & flagsOp.clear);

    DWORD save = CORE.Window.flags;
    CORE.Window.flags |= flagsOp.set;
    CORE.Window.flags &= ~flagsOp.clear;
    if (save != CORE.Window.flags) TRACELOG(LOG_DEBUG, "WIN32: FLAGS: Current deferred flags: 0x%x > 0x%x (diff 0x%x)", save, CORE.Window.flags, save ^ CORE.Window.flags);

    return result;
}

// Handle keyboard input event
static void HandleKey(WPARAM wparam, LPARAM lparam, char state)
{
    bool isExtendedKey = ((HIWORD(lparam) & KF_EXTENDED) != 0);

    WORD scancode = LOBYTE(HIWORD(lparam));
    if (scancode == 0)
    {
        // Certain keys don't arrive with a scancode; fall back to virtual key mapping
        scancode = MapVirtualKeyW((UINT)wparam, MAPVK_VK_TO_VSC);
    }

    KeyboardKey key = KEY_NULL;
    if (isExtendedKey)
    {
        unsigned int index = scancode - EXTENDED_KEY_OFFSET;
        unsigned int arrayLength = sizeof(keyboardKeyFromExtendedScancode)/sizeof(keyboardKeyFromExtendedScancode[0]);
        if (index < arrayLength) key = keyboardKeyFromExtendedScancode[index];
    }
    else
    {
        unsigned int arrayLength = sizeof(keyboardKeyFromScancode)/sizeof(keyboardKeyFromScancode[0]);
        if (scancode < arrayLength) key = keyboardKeyFromScancode[scancode];
    }

    bool ignoreKey = false;
    if ((wparam == VK_CONTROL) && !isExtendedKey)
    {
        // Pressing AltGr causes Windows to send a fake Left Ctrl message
        // GLFW does something similar in win32_window.c
        MSG next = { 0 };
        LONG time = GetMessageTime();

        if (PeekMessageW(&next, NULL, 0, 0, PM_NOREMOVE))
        {
            // If the current message is Left Ctrl (VK_CONTROL) and the next message is AltGr (VK_MENU) with the exact same timestamp,
            // silently drop the fake control key event
            bool isNextExtendedKey = ((HIWORD(next.lParam) & KF_EXTENDED) != 0);
            if ((next.wParam == VK_MENU) && ((LONG)next.time == time) && isNextExtendedKey)
            {
                if ((next.message == WM_KEYDOWN) || (next.message == WM_SYSKEYDOWN) ||
                    (next.message == WM_KEYUP) || (next.message == WM_SYSKEYUP)) ignoreKey = true;
            }
        }
    }

    if (!ignoreKey)
    {
        if (key != KEY_NULL)
        {
            bool isKeyRepeatEvent = ((HIWORD(lparam) & KF_REPEAT) != 0);

            if (state == 1)
            {
                CORE.Input.Keyboard.currentKeyState[key] = 1;
                platform.keyReleasePending[key] = false;

                if (isKeyRepeatEvent) CORE.Input.Keyboard.keyRepeatInFrame[key] = 1;
                else
                {
                    platform.keyPressedInFrame[key] = true;

                    // Add key to the queue
                    if (CORE.Input.Keyboard.keyPressedQueueCount < MAX_KEY_PRESSED_QUEUE)
                    {
                        CORE.Input.Keyboard.keyPressedQueue[CORE.Input.Keyboard.keyPressedQueueCount] = key;
                        CORE.Input.Keyboard.keyPressedQueueCount++;
                    }

                    if (key == CORE.Input.Keyboard.exitKey) CORE.Window.shouldClose = true;
                }
            }
            else if (platform.keyPressedInFrame[key]) platform.keyReleasePending[key] = true; // Key press is kept for the current frame
            else CORE.Input.Keyboard.currentKeyState[key] = 0;
        }
        else if (wparam != VK_PACKET) TRACELOG(LOG_WARNING, "INPUT: Unknown (or currently unhandled) virtual keycode %d (0x%x)", wparam, wparam);
    }
}

// Register character input
static void HandleChar(unsigned int codepoint)
{
    // Check if there is space available in the queue, control characters are not registered
    if ((codepoint >= 32) && (codepoint != 127) && (CORE.Input.Keyboard.charPressedQueueCount < MAX_CHAR_PRESSED_QUEUE))
    {
        // Add character to the queue
        CORE.Input.Keyboard.charPressedQueue[CORE.Input.Keyboard.charPressedQueueCount] = (int)codepoint;
        CORE.Input.Keyboard.charPressedQueueCount++;
    }
}

// Register touch event
static void HandleTouch(int touchAction, int id, float x, float y)
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

    platform.touchIds[slot] = id;
    platform.touchPositions[slot] = (Vector2){ x, y };

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

    // Map primary touch point to mouse position (pixels, as registered by mouse messages)
    if (CORE.Input.Touch.pointCount > 0)
    {
        CORE.Input.Mouse.currentPosition.x = CORE.Input.Touch.position[0].x/CORE.Input.Mouse.scale.x;
        CORE.Input.Mouse.currentPosition.y = CORE.Input.Touch.position[0].y/CORE.Input.Mouse.scale.y;
    }
}

// Register dropped files
static void HandleDroppedFiles(HDROP drop)
{
    unsigned int count = DragQueryFileW(drop, 0xFFFFFFFF, NULL, 0);

    if (count > 0)
    {
        // In case previous dropped filepaths have not been freed, free them
        if (CORE.Window.dropFileCount > 0)
        {
            for (unsigned int i = 0; i < CORE.Window.dropFileCount; i++) RL_FREE(CORE.Window.dropFilepaths[i]);

            RL_FREE(CORE.Window.dropFilepaths);

            CORE.Window.dropFileCount = 0;
            CORE.Window.dropFilepaths = NULL;
        }

        // WARNING: Paths are freed by the system on DragFinish(), they need to be copied
        CORE.Window.dropFileCount = count;
        CORE.Window.dropFilepaths = (char **)RL_CALLOC(CORE.Window.dropFileCount, sizeof(char *));

        for (unsigned int i = 0; i < CORE.Window.dropFileCount; i++)
        {
            CORE.Window.dropFilepaths[i] = (char *)RL_CALLOC(MAX_FILEPATH_LENGTH, sizeof(char));

            unsigned int length = DragQueryFileW(drop, i, NULL, 0);
            WCHAR *path = (WCHAR *)RL_CALLOC(length + 1, sizeof(WCHAR));
            DragQueryFileW(drop, i, path, length + 1);
            WideCharToMultiByte(CP_UTF8, 0, path, -1, CORE.Window.dropFilepaths[i], MAX_FILEPATH_LENGTH - 1, NULL, NULL);
            RL_FREE(path);
        }
    }

    DragFinish(drop);
}

// Update gamepads state (XInput)
// NOTE: Not connected gamepads are checked once per second, checking them is slow
static void PollGamepads(void)
{
    if (pXInputGetState == NULL) return;

    // XInput buttons bits and their raylib gamepad buttons
    static const struct { WORD mask; int button; } buttons[] = {
        { 0x0001, GAMEPAD_BUTTON_LEFT_FACE_UP },
        { 0x0002, GAMEPAD_BUTTON_LEFT_FACE_DOWN },
        { 0x0004, GAMEPAD_BUTTON_LEFT_FACE_LEFT },
        { 0x0008, GAMEPAD_BUTTON_LEFT_FACE_RIGHT },
        { 0x0010, GAMEPAD_BUTTON_MIDDLE_RIGHT },
        { 0x0020, GAMEPAD_BUTTON_MIDDLE_LEFT },
        { 0x0040, GAMEPAD_BUTTON_LEFT_THUMB },
        { 0x0080, GAMEPAD_BUTTON_RIGHT_THUMB },
        { 0x0100, GAMEPAD_BUTTON_LEFT_TRIGGER_1 },
        { 0x0200, GAMEPAD_BUTTON_RIGHT_TRIGGER_1 },
        { 0x1000, GAMEPAD_BUTTON_RIGHT_FACE_DOWN },
        { 0x2000, GAMEPAD_BUTTON_RIGHT_FACE_RIGHT },
        { 0x4000, GAMEPAD_BUTTON_RIGHT_FACE_LEFT },
        { 0x8000, GAMEPAD_BUTTON_RIGHT_FACE_UP }
    };

    double time = GetTime();
    bool checkConnected = (time >= platform.gamepadCheckTime);
    if (checkConnected) platform.gamepadCheckTime = time + 1.0;

    for (int i = 0; (i < XINPUT_MAX_GAMEPADS) && (i < MAX_GAMEPADS); i++)
    {
        if (!CORE.Input.Gamepad.ready[i] && !checkConnected) continue;

        XInputState state = { 0 };
        if (pXInputGetState((DWORD)i, &state) != ERROR_SUCCESS)
        {
            if (CORE.Input.Gamepad.ready[i])
            {
                CORE.Input.Gamepad.ready[i] = false;
                CORE.Input.Gamepad.axisCount[i] = 0;
                platform.gamepadVibrationEnd[i] = 0.0;
                memset(CORE.Input.Gamepad.name[i], 0, MAX_GAMEPAD_NAME_LENGTH);
                memset(CORE.Input.Gamepad.currentButtonState[i], 0, sizeof(CORE.Input.Gamepad.currentButtonState[i]));
                memset(CORE.Input.Gamepad.previousButtonState[i], 0, sizeof(CORE.Input.Gamepad.previousButtonState[i]));
                memset(CORE.Input.Gamepad.axisState[i], 0, sizeof(CORE.Input.Gamepad.axisState[i]));
            }

            continue;
        }

        if (!CORE.Input.Gamepad.ready[i])
        {
            CORE.Input.Gamepad.ready[i] = true;
            CORE.Input.Gamepad.axisCount[i] = 6;
            snprintf(CORE.Input.Gamepad.name[i], MAX_GAMEPAD_NAME_LENGTH, "%s", XINPUT_GAMEPAD_NAME);
        }

        // Register previous gamepad states
        for (int k = 0; k < MAX_GAMEPAD_BUTTONS; k++) CORE.Input.Gamepad.previousButtonState[i][k] = CORE.Input.Gamepad.currentButtonState[i][k];

        for (unsigned int k = 0; k < sizeof(buttons)/sizeof(buttons[0]); k++)
        {
            if (state.gamepad.buttons & buttons[k].mask)
            {
                CORE.Input.Gamepad.currentButtonState[i][buttons[k].button] = 1;
                CORE.Input.Gamepad.lastButtonPressed = buttons[k].button;
            }
            else CORE.Input.Gamepad.currentButtonState[i][buttons[k].button] = 0;
        }

        // Axes in the range [-1, 1], sticks Y axis points down, triggers are -1 when released
        CORE.Input.Gamepad.axisState[i][GAMEPAD_AXIS_LEFT_X] = fmaxf((float)state.gamepad.thumbLX/32767.0f, -1.0f);
        CORE.Input.Gamepad.axisState[i][GAMEPAD_AXIS_LEFT_Y] = -fmaxf((float)state.gamepad.thumbLY/32767.0f, -1.0f);
        CORE.Input.Gamepad.axisState[i][GAMEPAD_AXIS_RIGHT_X] = fmaxf((float)state.gamepad.thumbRX/32767.0f, -1.0f);
        CORE.Input.Gamepad.axisState[i][GAMEPAD_AXIS_RIGHT_Y] = -fmaxf((float)state.gamepad.thumbRY/32767.0f, -1.0f);
        CORE.Input.Gamepad.axisState[i][GAMEPAD_AXIS_LEFT_TRIGGER] = ((float)state.gamepad.leftTrigger/255.0f)*2.0f - 1.0f;
        CORE.Input.Gamepad.axisState[i][GAMEPAD_AXIS_RIGHT_TRIGGER] = ((float)state.gamepad.rightTrigger/255.0f)*2.0f - 1.0f;

        // Register buttons for 2nd triggers (because they are axes)
        if (CORE.Input.Gamepad.axisState[i][GAMEPAD_AXIS_LEFT_TRIGGER] > 0.1f)
        {
            CORE.Input.Gamepad.currentButtonState[i][GAMEPAD_BUTTON_LEFT_TRIGGER_2] = 1;
            CORE.Input.Gamepad.lastButtonPressed = GAMEPAD_BUTTON_LEFT_TRIGGER_2;
        }
        else CORE.Input.Gamepad.currentButtonState[i][GAMEPAD_BUTTON_LEFT_TRIGGER_2] = 0;

        if (CORE.Input.Gamepad.axisState[i][GAMEPAD_AXIS_RIGHT_TRIGGER] > 0.1f)
        {
            CORE.Input.Gamepad.currentButtonState[i][GAMEPAD_BUTTON_RIGHT_TRIGGER_2] = 1;
            CORE.Input.Gamepad.lastButtonPressed = GAMEPAD_BUTTON_RIGHT_TRIGGER_2;
        }
        else CORE.Input.Gamepad.currentButtonState[i][GAMEPAD_BUTTON_RIGHT_TRIGGER_2] = 0;

        // Stop vibration when its time is over
        if ((platform.gamepadVibrationEnd[i] > 0.0) && (time >= platform.gamepadVibrationEnd[i]) && (pXInputSetState != NULL))
        {
            XInputVibration vibration = { 0 };
            pXInputSetState((DWORD)i, &vibration);
            platform.gamepadVibrationEnd[i] = 0.0;
        }
    }
}

// Confine the cursor to the window center while it is disabled and the window is focused, release it otherwise
static void UpdateCursorClip(void)
{
    if (CORE.Input.Mouse.cursorLocked && !(CORE.Window.flags & FLAG_WINDOW_UNFOCUSED) && (GetForegroundWindow() == platform.hwnd))
    {
        RECT clientRect = { 0 };
        POINT topLeft = { 0 };
        GetClientRect(platform.hwnd, &clientRect);
        ClientToScreen(platform.hwnd, &topLeft);

        LONG centerX = topLeft.x + clientRect.right/2;
        LONG centerY = topLeft.y + clientRect.bottom/2;
        RECT clipRect = { centerX, centerY, centerX + 1, centerY + 1 };
        if (!ClipCursor(&clipRect)) TRACELOG(LOG_WARNING, "WIN32: Failed to clip cursor [ERROR: %lu]", GetLastError());
    }
    else ClipCursor(NULL);
}

// Place input method windows next to the last clicked position
// NOTE: Text being composed is drawn by the input method, there is no text cursor
static void UpdateInputMethodPosition(HWND hwnd)
{
    if ((pImmGetContext == NULL) || (pImmReleaseContext == NULL)) return;

    HIMC context = pImmGetContext(hwnd);
    if (context == NULL) return;

    if (pImmSetCompositionWindow != NULL)
    {
        COMPOSITIONFORM composition = { 0 };
        composition.dwStyle = CFS_POINT;
        composition.ptCurrentPos = platform.textInputPosition;
        pImmSetCompositionWindow(context, &composition);
    }

    if (pImmSetCandidateWindow != NULL)
    {
        CANDIDATEFORM candidate = { 0 };
        candidate.dwIndex = 0;
        candidate.dwStyle = CFS_CANDIDATEPOS;
        candidate.ptCurrentPos = platform.textInputPosition;
        pImmSetCandidateWindow(context, &candidate);
    }

    pImmReleaseContext(hwnd, context);
}

// Handle mouse button input event
static void HandleMouseButton(int button, char state)
{
    if (state == 1)
    {
        // Capture the mouse on MOUSEDOWN events and no other buttons are pressed
        bool allMouseButtonsReleased = true;
        for (unsigned int i = 0; i < MAX_MOUSE_BUTTONS; i++)
        {
            if ((CORE.Input.Mouse.currentButtonState[i] != 0) && !platform.buttonReleasePending[i])
            {
                allMouseButtonsReleased = false;
                break;
            }
        }

        if (allMouseButtonsReleased) SetCapture(platform.hwnd);

        // Register current mouse button state
        CORE.Input.Mouse.currentButtonState[button] = 1;
        CORE.Input.Touch.currentTouchState[button] = 1;
        platform.buttonPressedInFrame[button] = true;
        platform.buttonReleasePending[button] = false;

        // Input method windows are placed next to the last clicked position
        POINT position = { 0 };
        if (GetCursorPos(&position) && ScreenToClient(platform.hwnd, &position)) platform.textInputPosition = position;
    }
    else
    {
        // NOTE: Button press is kept for the current frame if it is released in the same frame
        if (platform.buttonPressedInFrame[button]) platform.buttonReleasePending[button] = true;
        else
        {
            CORE.Input.Mouse.currentButtonState[button] = 0;
            CORE.Input.Touch.currentTouchState[button] = 0;
        }

        // Release mouse capture when all mouse buttons are released
        bool allMouseButtonsReleased = true;
        for (unsigned int i = 0; i < MAX_MOUSE_BUTTONS; i++)
        {
            if ((CORE.Input.Mouse.currentButtonState[i] != 0) && !platform.buttonReleasePending[i])
            {
                allMouseButtonsReleased = false;
                break;
            }
        }

        if (allMouseButtonsReleased) ReleaseCapture();
    }
}

// Handle raw input event
static void HandleRawInput(LPARAM lparam)
{
    RAWINPUT input = { 0 };
    UINT inputSize = 0;

    if (GetRawInputData((HRAWINPUT)lparam, RID_INPUT, NULL, &inputSize, sizeof(RAWINPUTHEADER)) != 0) return;
    if (inputSize > sizeof(input)) return;

    UINT size = GetRawInputData((HRAWINPUT)lparam, RID_INPUT, &input, &inputSize, sizeof(RAWINPUTHEADER));

    if (size == (UINT)-1) return;
    if (input.header.dwType != RIM_TYPEMOUSE) return;

    // NOTE: Raw mouse motion is only used while the cursor is disabled and the window is focused
    if (!CORE.Input.Mouse.cursorLocked || (CORE.Window.flags & FLAG_WINDOW_UNFOCUSED)) return;

    LONG deltaX = input.data.mouse.lLastX;
    LONG deltaY = input.data.mouse.lLastY;

    if (input.data.mouse.usFlags & MOUSE_MOVE_ABSOLUTE)
    {
        // Absolute pointing devices (remote desktop, tablets, virtual machines) report [0, 65535] over the primary monitor or the virtual desktop
        POINT position = { 0 };
        int width = GetSystemMetrics(SM_CXSCREEN);
        int height = GetSystemMetrics(SM_CYSCREEN);

        if (input.data.mouse.usFlags & MOUSE_VIRTUAL_DESKTOP)
        {
            position.x = GetSystemMetrics(SM_XVIRTUALSCREEN);
            position.y = GetSystemMetrics(SM_YVIRTUALSCREEN);
            width = GetSystemMetrics(SM_CXVIRTUALSCREEN);
            height = GetSystemMetrics(SM_CYVIRTUALSCREEN);
        }

        position.x += (LONG)(((float)input.data.mouse.lLastX/65535.0f)*(float)width);
        position.y += (LONG)(((float)input.data.mouse.lLastY/65535.0f)*(float)height);

        deltaX = platform.rawAbsoluteValid? (position.x - platform.rawAbsolutePosition.x) : 0;
        deltaY = platform.rawAbsoluteValid? (position.y - platform.rawAbsolutePosition.y) : 0;
        platform.rawAbsolutePosition = position;
        platform.rawAbsoluteValid = true;
    }

    CORE.Input.Mouse.currentPosition.x += (float)deltaX;
    CORE.Input.Mouse.currentPosition.y += (float)deltaY;
    CORE.Input.Touch.position[0] = CORE.Input.Mouse.currentPosition;
}

// Handle window resizing event
static void HandleWindowResize(HWND hwnd, int *width, int *height)
{
    if (CORE.Window.flags & FLAG_WINDOW_MINIMIZED) return;

    // Get client size (framebuffer inside the window)
    RECT rect = { 0 };
    GetClientRect(hwnd, &rect);
    SIZE clientSize = { rect.right, rect.bottom };

    if ((CORE.Window.render.width != clientSize.cx) || (CORE.Window.render.height != clientSize.cy)) CORE.Window.resizedLastFrame = true;
    CORE.Window.render.width = (int)clientSize.cx;
    CORE.Window.render.height = (int)clientSize.cy;
    CORE.Window.currentFbo.width = (int)clientSize.cx;
    CORE.Window.currentFbo.height = (int)clientSize.cy;

    float dpiScale = ((float)GetDpiForWindowCompat(hwnd))/96.0f;
    bool highdpi = !!(CORE.Window.flags & FLAG_WINDOW_HIGHDPI);
    unsigned int screenWidth = highdpi? (unsigned int)roundf(((float)clientSize.cx)/dpiScale) : clientSize.cx;
    unsigned int screenHeight = highdpi? (unsigned int)roundf(((float)clientSize.cy)/dpiScale) : clientSize.cy;
    CORE.Window.screen.width = screenWidth;
    CORE.Window.screen.height = screenHeight;

    if (AdoptWindowResize(CORE.Window.flags))
    {
        TRACELOG(LOG_DEBUG, "WIN32: WINDOW: Updating app size to [%ix%i] from window resize", screenWidth, screenHeight);
        *width = screenWidth;
        *height = screenHeight;
    }

    CORE.Window.screenScale = MatrixScale( (float)CORE.Window.render.width/CORE.Window.screen.width,
        (float)CORE.Window.render.height/CORE.Window.screen.height, 1.0f);

    if (CORE.Window.flags & FLAG_WINDOW_HIGHDPI)
    {
        float mouseScale = 1.0f/dpiScale;
        CORE.Input.Mouse.scale = (Vector2){ mouseScale, mouseScale };
    }
}

// Update window style
static void UpdateWindowStyle(HWND hwnd, unsigned desiredFlags)
{
    DWORD current = STYLE_MASK_WRITABLE & MakeWindowStyle(CORE.Window.flags);
    DWORD desired = STYLE_MASK_WRITABLE & MakeWindowStyle(desiredFlags);

    if (current != desired)
    {
        // Client area position and size, kept when the window frame changes
        RECT clientRect = { 0 };
        POINT clientOrigin = { 0 };
        GetClientRect(hwnd, &clientRect);
        ClientToScreen(hwnd, &clientOrigin);
        bool floating = (!IsIconic(hwnd) && !IsZoomed(hwnd));

        SetLastError(0);
        DWORD previous = STYLE_MASK_WRITABLE & SetWindowLongPtrW(hwnd, GWL_STYLE, desired);
        if (previous != current)
        {
            TRACELOG(LOG_ERROR, "WIN32: WINDOW: SetWindowLongPtr() returned writable flags 0x%x but expected 0x%x (diff=0x%x, error=%lu)",
                previous, current, previous ^ current, GetLastError());
        }

        CheckFlags("UpdateWindowStyle", hwnd, desiredFlags, desired, STYLE_MASK_WRITABLE);

        // NOTE: New window frame is only applied by SetWindowPos() with SWP_FRAMECHANGED
        if (floating)
        {
            RECT rect = { clientOrigin.x, clientOrigin.y, clientOrigin.x + clientRect.right, clientOrigin.y + clientRect.bottom };
            AdjustWindowRectExForDpiCompat(&rect, desired, FALSE, WINDOW_STYLE_EX, GetDpiForWindowCompat(hwnd));
            SetWindowPos(hwnd, NULL, rect.left, rect.top, rect.right - rect.left, rect.bottom - rect.top, SWP_NOZORDER | SWP_NOACTIVATE | SWP_FRAMECHANGED);
        }
        else SetWindowPos(hwnd, NULL, 0, 0, 0, 0, SWP_NOMOVE | SWP_NOSIZE | SWP_NOZORDER | SWP_NOACTIVATE | SWP_FRAMECHANGED);
    }

    // Minimized takes precedence over maximized
    Mized currentMized = MIZED_NONE;
    Mized desiredMized = MIZED_NONE;
    if (CORE.Window.flags & FLAG_WINDOW_MINIMIZED) currentMized = MIZED_MIN;
    else if (CORE.Window.flags & FLAG_WINDOW_MAXIMIZED) currentMized = MIZED_MAX;
    if (desiredFlags & FLAG_WINDOW_MINIMIZED) desiredMized = MIZED_MIN;
    else if (desiredFlags & FLAG_WINDOW_MAXIMIZED) desiredMized = MIZED_MAX;

    if (currentMized != desiredMized)
    {
        switch (desiredMized)
        {
            case MIZED_NONE: ShowWindow(hwnd, SW_RESTORE); break;
            case MIZED_MIN: ShowWindow(hwnd, SW_MINIMIZE); break;
            case MIZED_MAX: ShowWindow(hwnd, SW_MAXIMIZE); break;
        }
    }
}

// Sanitize flags
static unsigned SanitizeFlags(unsigned flags)
{
    // Fullscreen and borderless windowed modes exclude each other, window is undecorated in both
    if ((flags & FLAG_MASK_COVER) == FLAG_MASK_COVER) flags &= ~FLAG_BORDERLESS_WINDOWED_MODE;
    if (flags & FLAG_MASK_COVER) flags |= FLAG_WINDOW_UNDECORATED;

    if (flags & FLAG_WINDOW_MAXIMIZED)
    {
        if (flags & FLAG_MASK_COVER)
        {
            // NOTE: Maximized state is restored when leaving fullscreen and borderless windowed modes
            flags &= ~FLAG_WINDOW_MAXIMIZED;
        }

        if (~flags & FLAG_WINDOW_RESIZABLE)
        {
            if (!(CORE.Window.flags & FLAG_WINDOW_MAXIMIZED))
            {
                TRACELOG(LOG_WARNING, "WIN32: WINDOW: Cannot maximize a non-resizable window");
                flags &= ~FLAG_WINDOW_MAXIMIZED;
            }
            else if (CORE.Window.flags & FLAG_WINDOW_RESIZABLE)
            {
                TRACELOG(LOG_WARNING, "WIN32: WINDOW: Cannot set window as non-resizable when maximized");
                flags |= FLAG_WINDOW_RESIZABLE;
            }
        }
        else if (!(CORE.Window.flags & FLAG_WINDOW_MAXIMIZED))
        {
            if (CORE.Window.flags & FLAG_WINDOW_MINIMIZED)
            {
                // Window needs to be unminimized before it can be maximized since minimizing takes precedence
                flags &= ~FLAG_WINDOW_MINIMIZED;
            }
            else if ((flags & FLAG_WINDOW_MINIMIZED) && !(CORE.Window.flags & FLAG_WINDOW_MINIMIZED))
            {
                TRACELOG(LOG_WARNING, "WIN32: WINDOW: Cannot minimize and maximize a window in the same frame");
                flags &= ~FLAG_WINDOW_MINIMIZED;
                flags &= ~FLAG_WINDOW_MAXIMIZED;
            }
        }
    }

    // Don't change flags that can't be updated after window initialization
    flags = (flags & ~FLAG_MASK_NO_UPDATE) | (CORE.Window.flags & FLAG_MASK_NO_UPDATE);

    return flags;
}

// All window state changes from raylib flags go through this function. It performs
// whatever operations are needed to update the window state to match the desired flags
// In most cases this function should not update CORE.Window.flags directly, instead,
// the window itself should update CORE.Window.flags in response to actual state changes
// This means that CORE.Window.flags should always represent the actual state of the
// window. This function will continue to perform these update operations so long as
// the state continues to change
//
// This design takes care of many odd corner cases. For example, in case of restoring
// a window that was previously maximized AND minimized and those two flags need to be removed,
// ShowWindow with SW_RESTORE twice need to be actually calleed. Another example is
// wheen having a maximized window, if the undecorated flag is modified then the window style
// needs to be updated, but updating the style would mean the window size would change
// causing the window to lose its Maximized state which would mean the window size
// needs to be updated, followed by the update of window style, a second time, to restore that maximized
// state. This implementation is able to handle any/all of these special situations with a
// retry loop that continues until either the desired state is reached or the state stops changing
static void UpdateFlags(HWND hwnd, unsigned desiredFlags, int width, int height)
{
    // Flags that apply immediately without needing any operations
    CORE.Window.flags |= (desiredFlags & FLAG_MASK_NO_UPDATE);

    int vsync = (desiredFlags & FLAG_VSYNC_HINT)? 1 : 0;
    if (wglSwapIntervalEXT)
    {
        wglSwapIntervalEXT(vsync);
        if (vsync) CORE.Window.flags |= FLAG_VSYNC_HINT;
        else CORE.Window.flags &= ~FLAG_VSYNC_HINT;
    }

    // Window always on top
    if ((desiredFlags ^ CORE.Window.flags) & FLAG_WINDOW_TOPMOST)
    {
        bool topmost = ((desiredFlags & FLAG_WINDOW_TOPMOST) != 0);
        SetWindowPos(hwnd, topmost? HWND_TOPMOST : HWND_NOTOPMOST, 0, 0, 0, 0, SWP_NOMOVE | SWP_NOSIZE | SWP_NOACTIVATE);

        if (topmost) CORE.Window.flags |= FLAG_WINDOW_TOPMOST;
        else CORE.Window.flags &= ~FLAG_WINDOW_TOPMOST;
    }

    // Mouse events pass through the window
    if ((desiredFlags ^ CORE.Window.flags) & FLAG_WINDOW_MOUSE_PASSTHROUGH)
    {
        bool passthrough = ((desiredFlags & FLAG_WINDOW_MOUSE_PASSTHROUGH) != 0);
        UpdateLayeredStyle(hwnd, passthrough, platform.opacity);

        if (passthrough) CORE.Window.flags |= FLAG_WINDOW_MOUSE_PASSTHROUGH;
        else CORE.Window.flags &= ~FLAG_WINDOW_MOUSE_PASSTHROUGH;
    }

    // Keep running while minimized
    CORE.Window.flags = (CORE.Window.flags & ~FLAG_WINDOW_ALWAYS_RUN) | (desiredFlags & FLAG_WINDOW_ALWAYS_RUN);

    // Fullscreen and borderless windowed modes
    bool covered = ((CORE.Window.flags & FLAG_MASK_COVER) != 0);
    bool cover = ((desiredFlags & FLAG_MASK_COVER) != 0);

    if (covered && !cover) LeaveCoverMode(hwnd, desiredFlags);
    else if (covered && cover)
    {
        // Mode and resizable flag changes do not change the window
        unsigned keptFlags = FLAG_MASK_COVER | FLAG_WINDOW_RESIZABLE;
        CORE.Window.flags = (CORE.Window.flags & ~keptFlags) | (desiredFlags & keptFlags);
    }

    // NOTE: Window state is updated before covering the monitor
    unsigned stateFlags = desiredFlags;
    if (!covered && cover) stateFlags = (desiredFlags & ~(FLAG_MASK_COVER | FLAG_WINDOW_UNDECORATED)) | (CORE.Window.flags & FLAG_WINDOW_UNDECORATED);

    // TODO: Review flags update code
    DWORD previousStyle = 0;
    for (unsigned attempt = 1; ; attempt++)
    {
        CheckFlags("UpdateFlags", hwnd, CORE.Window.flags, MakeWindowStyle(CORE.Window.flags), STYLE_MASK_ALL);

        bool windowSizeUpdated = false;
        if (MakeWindowStyle(CORE.Window.flags) == MakeWindowStyle(stateFlags))
        {
            windowSizeUpdated = UpdateWindowSize(hwnd, width, height, stateFlags);
            if ((FLAG_MASK_REQUIRED & stateFlags) == (FLAG_MASK_REQUIRED & CORE.Window.flags)) break;
        }

        if ((attempt > 1) && (previousStyle == MakeWindowStyle(CORE.Window.flags)) && !windowSizeUpdated)
        {
            TRACELOG(LOG_ERROR, "WIN32: WINDOW: UpdateFlags() failed after %u attempt(s) wanted 0x%x but is 0x%x (diff=0x%x)",
                attempt, stateFlags, CORE.Window.flags, stateFlags ^ CORE.Window.flags);
            break;
        }

        previousStyle = MakeWindowStyle(CORE.Window.flags);
        UpdateWindowStyle(hwnd, stateFlags);
    }

    if (!covered && cover) EnterCoverMode(hwnd, desiredFlags);
}

// Update window layered style: mouse passthrough and opacity
static void UpdateLayeredStyle(HWND hwnd, bool passthrough, float opacity)
{
    if (hwnd == NULL) return;

    LONG_PTR exStyle = GetWindowLongPtrW(hwnd, GWL_EXSTYLE);
    bool layered = (passthrough || (opacity < 1.0f));

    if (layered) exStyle |= WS_EX_LAYERED;
    else exStyle &= ~WS_EX_LAYERED;

    if (passthrough) exStyle |= WS_EX_TRANSPARENT;
    else exStyle &= ~WS_EX_TRANSPARENT;

    SetWindowLongPtrW(hwnd, GWL_EXSTYLE, exStyle);

    if (layered) SetLayeredWindowAttributes(hwnd, 0, (BYTE)(255.0f*opacity), LWA_ALPHA);
}

// Cover the monitor with the window: fullscreen or borderless windowed mode
// NOTE: The window is undecorated and has the monitor size, the video mode is not changed
static void EnterCoverMode(HWND hwnd, unsigned desiredFlags)
{
    if (CORE.Window.flags & FLAG_WINDOW_MINIMIZED) return;

    GetWindowRect(hwnd, &platform.coverRestoreRect);

    MONITORINFO info = { 0 };
    info.cbSize = sizeof(info);
    if (!GetMonitorInfoW(MonitorFromWindow(hwnd, MONITOR_DEFAULTTONEAREST), &info))
    {
        TRACELOG(LOG_ERROR, "WIN32: Failed to get monitor info [ERROR: %lu]", GetLastError());
        return;
    }

    DWORD current = STYLE_MASK_WRITABLE & MakeWindowStyle(CORE.Window.flags);
    DWORD desired = STYLE_MASK_WRITABLE & MakeWindowStyle(desiredFlags);

    if (current != desired)
    {
        // Flags are updated by the window on style change
        platform.coverEnterFlags = (desiredFlags & FLAG_MASK_COVER);
        SetWindowLongPtrW(hwnd, GWL_STYLE, desired);
        platform.coverEnterFlags = 0;
    }
    else CORE.Window.flags |= (desiredFlags & FLAG_MASK_COVER);

    SetWindowPos(hwnd, HWND_TOP, info.rcMonitor.left, info.rcMonitor.top, info.rcMonitor.right - info.rcMonitor.left,
        info.rcMonitor.bottom - info.rcMonitor.top, SWP_NOOWNERZORDER | SWP_FRAMECHANGED);
}

// Restore the window from fullscreen or borderless windowed mode
static void LeaveCoverMode(HWND hwnd, unsigned desiredFlags)
{
    // NOTE: Resizable flag does not change the window while it covers the monitor
    CORE.Window.flags = (CORE.Window.flags & ~FLAG_WINDOW_RESIZABLE) | (desiredFlags & FLAG_WINDOW_RESIZABLE);

    DWORD current = STYLE_MASK_WRITABLE & MakeWindowStyle(CORE.Window.flags);
    DWORD desired = STYLE_MASK_WRITABLE & MakeWindowStyle(desiredFlags);

    if (current != desired)
    {
        // Flags are updated by the window on style change
        platform.coverLeave = true;
        SetWindowLongPtrW(hwnd, GWL_STYLE, desired);
        platform.coverLeave = false;
    }
    else CORE.Window.flags &= ~FLAG_MASK_COVER;

    SetWindowPos(hwnd, NULL, platform.coverRestoreRect.left, platform.coverRestoreRect.top,
        platform.coverRestoreRect.right - platform.coverRestoreRect.left, platform.coverRestoreRect.bottom - platform.coverRestoreRect.top,
        SWP_NOZORDER | SWP_NOOWNERZORDER | SWP_FRAMECHANGED);
}

// Check if OpenGL extension is available
static bool IsWglExtensionAvailable(HDC hdc, const char *extension)
{
    bool result = false;

    if (wglGetExtensionsStringARB != NULL)
    {
        const char *extList = wglGetExtensionsStringARB(hdc);
        if (extList != NULL)
        {
            // Simple substring search (could use strtok or strstr)
            if (strstr(extList, extension) != NULL) result = true;
        }
    }

    return result;
}
