/**********************************************************************************************
*
*   rcore_ios - Functions to manage window, graphics device and inputs
*
*   PLATFORM: IOS
*       - iOS, iPadOS (arm64)
*
*   LIMITATIONS:
*       - Window is always fullscreen, window management functions are not available
*       - SetGamepadMappings() and SetGamepadVibration() not available (gamepads are mapped by the system)
*       - No drawing while the application is not active (OpenGL ES is not allowed in background)
*       - FLAG_WINDOW_HIGHDPI not available: screen size is always the framebuffer size in pixels
*       - Hardware keyboard keys are not repeated by the system
*
*   POSSIBLE IMPROVEMENTS:
*       - Pointer devices (mouse, trackpad)
*       - Gamepad vibration (CoreHaptics)
*
*   ADDITIONAL NOTES:
*       - The UIKit application is started first, the program runs from the main run loop once the window exists
*       - Requested screen size is letterboxed to the display aspect ratio and scaled by the system
*
*   CONFIGURATION:
*       #define RCORE_PLATFORM_CUSTOM_FLAG
*           Custom flag for rcore on target platform -not used-
*
*   DEPENDENCIES:
*       - UIKit, QuartzCore: Application, window, events
*       - OpenGLES framework: EAGL context (deprecated by Apple, still available)
*       - CoreGraphics: Software renderer framebuffer presentation
*       - GameController: Gamepads
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

#import <UIKit/UIKit.h>             // UIKit and Foundation
#import <QuartzCore/QuartzCore.h>   // Required for: CAEAGLLayer, CADisplayLink
#import <GameController/GameController.h> // Required for: GCController, GCExtendedGamepad
#if !defined(GRAPHICS_API_OPENGL_SOFTWARE)
#import <OpenGLES/EAGL.h>           // Required for: EAGLContext
#import <OpenGLES/EAGLDrawable.h>   // Required for: renderbufferStorage:fromDrawable:
#endif

#undef Point                        // rcore symbol collision
#undef Size                         // rcore symbol collision

#include <os/log.h>                 // Required for: os_log_with_type()
#include <dlfcn.h>                  // Required for: dlsym()
#include <time.h>                   // Required for: clock_gettime()
#include <math.h>                   // Required for: round()

// Platform framebuffers use the real names, name 0 is only mapped for rlgl
#undef glBindFramebuffer
#undef glGetIntegerv
#undef glReadPixels

//----------------------------------------------------------------------------------
// Defines and Macros
//----------------------------------------------------------------------------------
#define GL_DEPTH24_STENCIL8_PACKED  0x88F0  // GL_DEPTH24_STENCIL8 (OpenGL ES 3.0), GL_DEPTH24_STENCIL8_OES
#define KEYMAP_SIZE                 232     // Size of the HID keyboard usages to raylib keys map

//----------------------------------------------------------------------------------
// Types and Structures Definition
//----------------------------------------------------------------------------------
@interface RaylibApplicationDelegate : UIResponder <UIApplicationDelegate>
@end

@interface RaylibSceneDelegate : UIResponder <UIWindowSceneDelegate>
@end

@interface RaylibViewController : UIViewController
@end

@interface RaylibView : UIView
@end

@interface RaylibDisplayLinkTarget : NSObject
@end

typedef struct {
    UIWindow *window;                   // Application window
    RaylibViewController *viewController; // Window root view controller
    RaylibView *view;                   // Window content view (touch events, OpenGL ES drawable)
    CADisplayLink *displayLink;         // Display refresh timer, paces frames with V-Sync
    RaylibDisplayLinkTarget *displayLinkTarget; // Display refresh timer target

#if defined(GRAPHICS_API_OPENGL_SOFTWARE)
    unsigned char *softwarePixels;      // Software renderer framebuffer copy (BGRA)
    int softwareWidth;                  // Software renderer framebuffer width
    int softwareHeight;                 // Software renderer framebuffer height
    CGColorSpaceRef colorSpace;         // Color space for the framebuffer image
#else
    EAGLContext *context;               // OpenGL ES context
    unsigned int framebuffer;           // Framebuffer presented on screen
    unsigned int colorRenderbuffer;     // Color renderbuffer, storage provided by the view layer
    unsigned int depthRenderbuffer;     // Depth and stencil renderbuffer
    unsigned int msaaFramebuffer;       // Multisampled framebuffer, resolved into the presented one (FLAG_MSAA_4X_HINT)
    unsigned int msaaColorRenderbuffer; // Multisampled color renderbuffer
#endif
    int drawableWidth;                  // Framebuffer width in pixels
    int drawableHeight;                 // Framebuffer height in pixels

    bool initialized;                   // Platform initialized, framebuffer can be resized
    bool appEnabled;                    // Application is active (foreground, receiving events)
    bool displayRefreshed;              // Display refresh happened since the last buffer swap
    bool frameSwapped;                  // Frame presented, V-Sync waits for the next display refresh
    bool closeRequested;                // Application termination requested
    bool eventReceived;                 // Input or application event received (event waiting)
    bool displaySize;                   // Screen size follows the display size, requested with InitWindow(0, 0)
    UIInterfaceOrientationMask orientations; // Interface orientations allowed for the requested screen size

    UITouch *touches[MAX_TOUCH_POINTS]; // Active touches, the index is the touch point id
    GCController *gamepads[MAX_GAMEPADS]; // Connected game controllers (extended gamepad profile), the index is the gamepad id

    char *clipboardText;                // Clipboard text, returned by GetClipboardText()
    char storagePath[MAX_FILEPATH_LENGTH]; // Application documents directory, base path for storage
} PlatformData;

//----------------------------------------------------------------------------------
// Global Variables Definition
//----------------------------------------------------------------------------------
extern CoreData CORE;                   // Global CORE state context

static PlatformData platform = { 0 };   // Platform specific data

#if !defined(GRAPHICS_API_OPENGL_SOFTWARE)
unsigned int eaglDefaultFramebuffer = 0; // Framebuffer used by rlgl as default framebuffer (name 0)
#endif

// Map HID keyboard usages (UIKeyboardHIDUsage) to raylib KeyboardKey values
static const int hidToRaylibMap[KEYMAP_SIZE] = {
    [0x04] = KEY_A, [0x05] = KEY_B, [0x06] = KEY_C, [0x07] = KEY_D, [0x08] = KEY_E, [0x09] = KEY_F, [0x0A] = KEY_G,
    [0x0B] = KEY_H, [0x0C] = KEY_I, [0x0D] = KEY_J, [0x0E] = KEY_K, [0x0F] = KEY_L, [0x10] = KEY_M, [0x11] = KEY_N,
    [0x12] = KEY_O, [0x13] = KEY_P, [0x14] = KEY_Q, [0x15] = KEY_R, [0x16] = KEY_S, [0x17] = KEY_T, [0x18] = KEY_U,
    [0x19] = KEY_V, [0x1A] = KEY_W, [0x1B] = KEY_X, [0x1C] = KEY_Y, [0x1D] = KEY_Z, [0x1E] = KEY_ONE, [0x1F] = KEY_TWO,
    [0x20] = KEY_THREE, [0x21] = KEY_FOUR, [0x22] = KEY_FIVE, [0x23] = KEY_SIX, [0x24] = KEY_SEVEN, [0x25] = KEY_EIGHT,
    [0x26] = KEY_NINE, [0x27] = KEY_ZERO, [0x28] = KEY_ENTER, [0x29] = KEY_ESCAPE, [0x2A] = KEY_BACKSPACE,
    [0x2B] = KEY_TAB, [0x2C] = KEY_SPACE, [0x2D] = KEY_MINUS, [0x2E] = KEY_EQUAL, [0x2F] = KEY_LEFT_BRACKET,
    [0x30] = KEY_RIGHT_BRACKET, [0x31] = KEY_BACKSLASH, [0x33] = KEY_SEMICOLON, [0x34] = KEY_APOSTROPHE,
    [0x35] = KEY_GRAVE, [0x36] = KEY_COMMA, [0x37] = KEY_PERIOD, [0x38] = KEY_SLASH, [0x39] = KEY_CAPS_LOCK,
    [0x3A] = KEY_F1, [0x3B] = KEY_F2, [0x3C] = KEY_F3, [0x3D] = KEY_F4, [0x3E] = KEY_F5, [0x3F] = KEY_F6,
    [0x40] = KEY_F7, [0x41] = KEY_F8, [0x42] = KEY_F9, [0x43] = KEY_F10, [0x44] = KEY_F11, [0x45] = KEY_F12,
    [0x46] = KEY_PRINT_SCREEN, [0x47] = KEY_SCROLL_LOCK, [0x48] = KEY_PAUSE, [0x49] = KEY_INSERT, [0x4A] = KEY_HOME,
    [0x4B] = KEY_PAGE_UP, [0x4C] = KEY_DELETE, [0x4D] = KEY_END, [0x4E] = KEY_PAGE_DOWN, [0x4F] = KEY_RIGHT,
    [0x50] = KEY_LEFT, [0x51] = KEY_DOWN, [0x52] = KEY_UP, [0x53] = KEY_NUM_LOCK, [0x54] = KEY_KP_DIVIDE,
    [0x55] = KEY_KP_MULTIPLY, [0x56] = KEY_KP_SUBTRACT, [0x57] = KEY_KP_ADD, [0x58] = KEY_KP_ENTER, [0x59] = KEY_KP_1,
    [0x5A] = KEY_KP_2, [0x5B] = KEY_KP_3, [0x5C] = KEY_KP_4, [0x5D] = KEY_KP_5, [0x5E] = KEY_KP_6, [0x5F] = KEY_KP_7,
    [0x60] = KEY_KP_8, [0x61] = KEY_KP_9, [0x62] = KEY_KP_0, [0x63] = KEY_KP_DECIMAL, [0x65] = KEY_KB_MENU,
    [0x67] = KEY_KP_EQUAL, [0xE0] = KEY_LEFT_CONTROL, [0xE1] = KEY_LEFT_SHIFT, [0xE2] = KEY_LEFT_ALT,
    [0xE3] = KEY_LEFT_SUPER, [0xE4] = KEY_RIGHT_CONTROL, [0xE5] = KEY_RIGHT_SHIFT, [0xE6] = KEY_RIGHT_ALT,
    [0xE7] = KEY_RIGHT_SUPER,
};

//----------------------------------------------------------------------------------
// Module Internal Functions Declaration
//----------------------------------------------------------------------------------
int InitPlatform(void);          // Initialize platform (graphics, inputs and more)
void ClosePlatform(void);        // Close platform

static bool InitGraphicsContext(void);                          // Initialize OpenGL ES context or software framebuffer
static void SetupOrientation(void);                             // Request the interface orientation for the requested screen size
static bool LoadFramebuffer(void);                              // Load framebuffer for the current view size
#if !defined(GRAPHICS_API_OPENGL_SOFTWARE)
int eaglResolveDefaultFramebuffer(void);                        // Resolve multisampled default framebuffer into the presented one
#endif
static void UnloadFramebuffer(void);                            // Unload framebuffer
static void SetupFramebuffer(int width, int height);            // Setup framebuffer and screen sizes for a display size
static void UpdateFramebufferSize(void);                        // Update framebuffer after a view size change (rotation)
static void ProcessRunLoop(bool wait);                          // Process pending run loop events, optionally waiting for one
static void HandleTouches(NSSet *touches, int action);          // Register touch points change
static void HandlePresses(NSSet *presses, bool pressed);        // Register hardware keyboard keys change
static void PollGamepads(void);                                 // Register gamepads connection and state
static void TraceLogIOS(int logType, const char *text, va_list args); // Send trace log messages to the system log
#if !defined(GRAPHICS_API_OPENGL_SOFTWARE)
static void *GetProcAddressIOS(const char *name);               // Get OpenGL ES procedure address
#endif

// Program entry point, called from the main run loop when the application window is ready
extern void ios_run(void);

//----------------------------------------------------------------------------------
// Module Functions Declaration
//----------------------------------------------------------------------------------
// NOTE: Functions declaration is provided by raylib.h

//----------------------------------------------------------------------------------
// Module Functions Definition: Application
//----------------------------------------------------------------------------------

// Application entry point
// NOTE: UIApplicationMain() never returns, the program runs inside the main run loop
__attribute__((weak)) int main(int argc, char *argv[])
{
    @autoreleasepool
    {
        SetTraceLogCallback(TraceLogIOS);

        return UIApplicationMain(argc, argv, nil, NSStringFromClass([RaylibApplicationDelegate class]));
    }
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
void ToggleFullscreen(void)
{
    TRACELOG(LOG_WARNING, "ToggleFullscreen() not available on target platform");
}

// Toggle borderless windowed mode
void ToggleBorderlessWindowed(void)
{
    TRACELOG(LOG_WARNING, "ToggleBorderlessWindowed() not available on target platform");
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

    // State change: FLAG_VSYNC_HINT
    if (!FLAG_IS_SET(CORE.Window.flags, FLAG_VSYNC_HINT) && FLAG_IS_SET(flags, FLAG_VSYNC_HINT)) FLAG_SET(CORE.Window.flags, FLAG_VSYNC_HINT);

    // State change: FLAG_WINDOW_ALWAYS_RUN
    if (!FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_ALWAYS_RUN) && FLAG_IS_SET(flags, FLAG_WINDOW_ALWAYS_RUN)) FLAG_SET(CORE.Window.flags, FLAG_WINDOW_ALWAYS_RUN);
}

// Clear window configuration state flags
void ClearWindowState(unsigned int flags)
{
    // State change: FLAG_VSYNC_HINT
    if (FLAG_IS_SET(CORE.Window.flags, FLAG_VSYNC_HINT) && FLAG_IS_SET(flags, FLAG_VSYNC_HINT)) FLAG_CLEAR(CORE.Window.flags, FLAG_VSYNC_HINT);

    // State change: FLAG_WINDOW_ALWAYS_RUN
    if (FLAG_IS_SET(CORE.Window.flags, FLAG_WINDOW_ALWAYS_RUN) && FLAG_IS_SET(flags, FLAG_WINDOW_ALWAYS_RUN)) FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_ALWAYS_RUN);
}

// Set icon for window
void SetWindowIcon(Image image)
{
    TRACELOG(LOG_WARNING, "SetWindowIcon() not available on target platform");
}

// Set icon for window, multiple images
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
    return (void *)platform.window; // Type: UIWindow*
}

// Get number of monitors
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
    return CORE.Window.display.width;
}

// Get selected monitor height (currently used by monitor)
int GetMonitorHeight(int monitor)
{
    return CORE.Window.display.height;
}

// Get selected monitor physical width in millimetres
int GetMonitorPhysicalWidth(int monitor)
{
    TRACELOG(LOG_WARNING, "GetMonitorPhysicalWidth() not available on target platform");
    return 0;
}

// Get selected monitor physical height in millimetres
int GetMonitorPhysicalHeight(int monitor)
{
    TRACELOG(LOG_WARNING, "GetMonitorPhysicalHeight() not available on target platform");
    return 0;
}

// Get selected monitor refresh rate
int GetMonitorRefreshRate(int monitor)
{
    @autoreleasepool
    {
        return (int)[[UIScreen mainScreen] maximumFramesPerSecond];
    }
}

// Get the human-readable, UTF-8 encoded name of the selected monitor
const char *GetMonitorName(int monitor)
{
    TRACELOG(LOG_WARNING, "GetMonitorName() not available on target platform");
    return "";
}

// Get window position XY on monitor
Vector2 GetWindowPosition(void)
{
    return (Vector2){ 0, 0 };
}

// Get window scale DPI factor for current monitor
Vector2 GetWindowScaleDPI(void)
{
    @autoreleasepool
    {
        float scale = (float)[[UIScreen mainScreen] nativeScale];

        return (Vector2){ scale, scale };
    }
}

// Set clipboard text content
void SetClipboardText(const char *text)
{
    @autoreleasepool
    {
        [[UIPasteboard generalPasteboard] setString:((text != NULL)? [NSString stringWithUTF8String:text] : @"")];
    }
}

// Get clipboard text content
// NOTE: returned string is allocated and freed by the platform
// NOTE: The system asks the user before reading content copied by another application
const char *GetClipboardText(void)
{
    RL_FREE(platform.clipboardText);
    platform.clipboardText = NULL;

    @autoreleasepool
    {
        UIPasteboard *pasteboard = [UIPasteboard generalPasteboard];
        NSString *text = [pasteboard hasStrings]? [pasteboard string] : nil;

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
        UIPasteboard *pasteboard = [UIPasteboard generalPasteboard];
        UIImage *pasteboardImage = [pasteboard hasImages]? [pasteboard image] : nil;
        NSData *data = (pasteboardImage != nil)? UIImagePNGRepresentation(pasteboardImage) : nil;

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
    CORE.Input.Mouse.cursorHidden = false;
}

// Hides mouse cursor
void HideCursor(void)
{
    CORE.Input.Mouse.cursorHidden = true;
}

// Enables cursor (unlock cursor)
void EnableCursor(void)
{
    // Set cursor position in the middle
    SetMousePosition(CORE.Window.screen.width/2, CORE.Window.screen.height/2);

    CORE.Input.Mouse.cursorLocked = false;
}

// Disables cursor (lock cursor)
void DisableCursor(void)
{
    // Set cursor position in the middle
    SetMousePosition(CORE.Window.screen.width/2, CORE.Window.screen.height/2);

    CORE.Input.Mouse.cursorLocked = true;
}

// Swap back buffer with front buffer (screen drawing)
void SwapScreenBuffer(void)
{
    @autoreleasepool
    {
#if defined(GRAPHICS_API_OPENGL_SOFTWARE)
        int width = platform.drawableWidth;
        int height = platform.drawableHeight;

        // Resize software renderer framebuffer to the drawable size, the frame is not presented
        if ((platform.softwareWidth != width) || (platform.softwareHeight != height))
        {
            platform.softwareWidth = width;
            platform.softwareHeight = height;
            RL_FREE(platform.softwarePixels);
            platform.softwarePixels = (unsigned char *)RL_MALLOC(width*height*4);
            rlResizeFramebuffer(width, height);
            return;
        }

        // NOTE: Software renderer framebuffer copy is in BGRA order
        rlCopyFramebuffer(0, 0, width, height, PIXELFORMAT_UNCOMPRESSED_R8G8B8A8, platform.softwarePixels);

        CFDataRef data = CFDataCreate(NULL, platform.softwarePixels, width*height*4);
        CGDataProviderRef provider = CGDataProviderCreateWithCFData(data);
        CGBitmapInfo info = kCGBitmapByteOrder32Little | kCGImageAlphaNoneSkipFirst;
        CGImageRef image = CGImageCreate(width, height, 8, 32, width*4, platform.colorSpace, info, provider, NULL, false, kCGRenderingIntentDefault);

        [CATransaction begin];
        [CATransaction setDisableActions:YES];
        [[platform.view layer] setContents:(id)image];
        [CATransaction commit];

        CGImageRelease(image);
        CGDataProviderRelease(provider);
        CFRelease(data);
#else
        if (platform.context == nil) return;

        eaglResolveDefaultFramebuffer();

        glBindRenderbuffer(GL_RENDERBUFFER, platform.colorRenderbuffer);
        [platform.context presentRenderbuffer:GL_RENDERBUFFER];

        glBindFramebuffer(GL_FRAMEBUFFER, eaglDefaultFramebuffer);
#endif

        // NOTE: V-Sync wait happens in PollInputEvents(), input events are processed while waiting
        platform.frameSwapped = true;
        platform.displayRefreshed = false;
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
            if (nsurl != nil) [[UIApplication sharedApplication] openURL:nsurl options:@{} completionHandler:nil];
            else TRACELOG(LOG_WARNING, "OpenURL() failed to open the URL");
        }
    }
}

//----------------------------------------------------------------------------------
// Module Functions Definition: Inputs
//----------------------------------------------------------------------------------

// Set internal gamepad mappings
int SetGamepadMappings(const char *mappings)
{
    TRACELOG(LOG_WARNING, "SetGamepadMappings() not available on target platform");
    return 0;
}

// Set gamepad vibration
void SetGamepadVibration(int gamepad, float leftMotor, float rightMotor, float duration)
{
    TRACELOG(LOG_WARNING, "SetGamepadVibration() not available on target platform");
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
    TRACELOG(LOG_WARNING, "SetMouseCursor() not available on target platform");
}

// Get physical key name
const char *GetKeyName(int key)
{
    TRACELOG(LOG_WARNING, "GetKeyName() not implemented on target platform");
    return "";
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

    // Register previous gamepad states
    for (int i = 0; i < MAX_GAMEPADS; i++)
    {
        if (CORE.Input.Gamepad.ready[i])
        {
            for (int k = 0; k < MAX_GAMEPAD_BUTTONS; k++) CORE.Input.Gamepad.previousButtonState[i][k] = CORE.Input.Gamepad.currentButtonState[i][k];
        }
    }

    CORE.Window.resizedLastFrame = false;

    // Process application and input events
    // NOTE: Program is paused while the application is not active, drawing is not allowed in background
    @autoreleasepool
    {
        ProcessRunLoop(false);

        if (CORE.Window.eventWaiting)
        {
            // Wait for an input or application event, display refresh is paused meanwhile
            [platform.displayLink setPaused:YES];
            while (!platform.eventReceived && !platform.closeRequested) ProcessRunLoop(true);
            [platform.displayLink setPaused:!platform.appEnabled];

            CORE.Time.previous = GetTime();
        }
        else if (platform.frameSwapped && FLAG_IS_SET(CORE.Window.flags, FLAG_VSYNC_HINT))
        {
            // Wait for the next display refresh after a frame was presented (V-Sync)
            while (!platform.displayRefreshed && platform.appEnabled) ProcessRunLoop(true);
        }

        platform.frameSwapped = false;
        platform.eventReceived = false;

        if (!platform.appEnabled && !platform.closeRequested)
        {
            while (!platform.appEnabled && !platform.closeRequested) ProcessRunLoop(true);

            CORE.Time.previous = GetTime();
        }

        PollGamepads();
    }

    if (platform.closeRequested) CORE.Window.shouldClose = true;
}

//----------------------------------------------------------------------------------
// Module Internal Functions Definition
//----------------------------------------------------------------------------------

// Initialize platform: graphics, inputs and more
int InitPlatform(void)
{
    @autoreleasepool
    {
        if (platform.view == nil)
        {
            TRACELOG(LOG_WARNING, "IOS: Application is not running, the program must be started by the platform main() (c-archive build)");
            return -1;
        }

        SetupOrientation();

        // Initialize graphics context
        if (!InitGraphicsContext())
        {
            TRACELOG(LOG_FATAL, "PLATFORM: Failed to initialize graphics device");
            ClosePlatform();
            return -1;
        }

        // NOTE: Screen size is the framebuffer size in pixels, content is not scaled
        FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_HIGHDPI);
        FLAG_SET(CORE.Window.flags, FLAG_FULLSCREEN_MODE);
        if (!platform.appEnabled) FLAG_SET(CORE.Window.flags, FLAG_WINDOW_UNFOCUSED);

        CORE.Window.ready = true;
    }

    TRACELOG(LOG_INFO, "DISPLAY: Device initialized successfully");
    TRACELOG(LOG_INFO, "    > Display size: %i x %i", CORE.Window.display.width, CORE.Window.display.height);
    TRACELOG(LOG_INFO, "    > Screen size:  %i x %i", CORE.Window.screen.width, CORE.Window.screen.height);
    TRACELOG(LOG_INFO, "    > Render size:  %i x %i", CORE.Window.render.width, CORE.Window.render.height);
    TRACELOG(LOG_INFO, "    > Viewport offsets: %i, %i", CORE.Window.renderOffset.x, CORE.Window.renderOffset.y);

    // Load OpenGL extensions
    // NOTE: GL procedures address loader is required to load extensions
    //----------------------------------------------------------------------------
#if !defined(GRAPHICS_API_OPENGL_SOFTWARE)
    rlLoadExtensions(GetProcAddressIOS);
#endif
    //----------------------------------------------------------------------------

    // Initialize timing system
    //----------------------------------------------------------------------------
    InitTimer();
    //----------------------------------------------------------------------------

    // Initialize storage system
    //----------------------------------------------------------------------------
    @autoreleasepool
    {
        NSString *documents = [NSSearchPathForDirectoriesInDomains(NSDocumentDirectory, NSUserDomainMask, YES) firstObject];
        if (documents != nil) snprintf(platform.storagePath, MAX_FILEPATH_LENGTH, "%s", [documents fileSystemRepresentation]);
    }

    CORE.Storage.basePath = platform.storagePath;
    //----------------------------------------------------------------------------

    platform.initialized = true;

#if defined(GRAPHICS_API_OPENGL_SOFTWARE)
    TRACELOG(LOG_INFO, "PLATFORM: IOS (Software): Initialized successfully");
#else
    TRACELOG(LOG_INFO, "PLATFORM: IOS: Initialized successfully");
#endif

    return 0;
}

// Close platform
void ClosePlatform(void)
{
    @autoreleasepool
    {
        platform.initialized = false;
        platform.displaySize = false;

#if defined(GRAPHICS_API_OPENGL_SOFTWARE)
        RL_FREE(platform.softwarePixels);
        if (platform.colorSpace != NULL) CGColorSpaceRelease(platform.colorSpace);
        platform.softwarePixels = NULL;
        platform.softwareWidth = 0;
        platform.softwareHeight = 0;
        platform.colorSpace = NULL;
#else
        if (platform.context != nil)
        {
            [EAGLContext setCurrentContext:platform.context];
            UnloadFramebuffer();
            [EAGLContext setCurrentContext:nil];
            [platform.context release];
            platform.context = nil;
        }
#endif

        for (int i = 0; i < MAX_GAMEPADS; i++)
        {
            [platform.gamepads[i] release];
            platform.gamepads[i] = nil;
            CORE.Input.Gamepad.ready[i] = false;
        }

        RL_FREE(platform.clipboardText);
        platform.clipboardText = NULL;
    }
}

// Initialize OpenGL ES context or software framebuffer
static bool InitGraphicsContext(void)
{
#if defined(GRAPHICS_API_OPENGL_SOFTWARE)
    [[platform.view layer] setContentsGravity:kCAGravityResize];
    [[platform.view layer] setOpaque:YES];
    platform.colorSpace = CGColorSpaceCreateDeviceRGB();
#else
    EAGLRenderingAPI api = (rlGetVersion() == RL_OPENGL_ES_30)? kEAGLRenderingAPIOpenGLES3 : kEAGLRenderingAPIOpenGLES2;

    platform.context = [[EAGLContext alloc] initWithAPI:api];
    if (platform.context == nil)
    {
        TRACELOG(LOG_WARNING, "EAGL: Failed to create OpenGL ES context");
        return false;
    }

    if (![EAGLContext setCurrentContext:platform.context])
    {
        TRACELOG(LOG_WARNING, "EAGL: Failed to make OpenGL ES context current");
        return false;
    }

    CAEAGLLayer *layer = (CAEAGLLayer *)[platform.view layer];
    [layer setOpaque:YES];
    [layer setDrawableProperties:@{ kEAGLDrawablePropertyRetainedBacking: @NO, kEAGLDrawablePropertyColorFormat: kEAGLColorFormatRGBA8 }];
#endif

    return LoadFramebuffer();
}

// Request the interface orientation for the requested screen size
// NOTE: Landscape for a screen wider than tall, portrait for a taller one, like the Android platform
static void SetupOrientation(void)
{
    platform.orientations = UIInterfaceOrientationMaskAll;

    if (CORE.Window.screen.width > CORE.Window.screen.height) platform.orientations = UIInterfaceOrientationMaskLandscape;
    else if (CORE.Window.screen.width < CORE.Window.screen.height) platform.orientations = UIInterfaceOrientationMaskPortrait;

    CGSize size = [platform.view bounds].size;
    bool landscape = (size.width > size.height);
    bool rotate = ((landscape && (platform.orientations == UIInterfaceOrientationMaskPortrait)) ||
        (!landscape && (platform.orientations == UIInterfaceOrientationMaskLandscape)));

    if ([[NSProcessInfo processInfo] operatingSystemVersion].majorVersion >= 16)
    {
#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wunguarded-availability-new"
        [platform.viewController setNeedsUpdateOfSupportedInterfaceOrientations];

        if (rotate)
        {
            UIWindowSceneGeometryPreferencesIOS *preferences = [[[UIWindowSceneGeometryPreferencesIOS alloc] initWithInterfaceOrientations:platform.orientations] autorelease];
            [[platform.window windowScene] requestGeometryUpdateWithPreferences:preferences errorHandler:nil];
        }
#pragma clang diagnostic pop
    }
    else [UIViewController attemptRotationToDeviceOrientation];

    // Wait for the rotation, the framebuffer is created for the final view size
    if (rotate)
    {
        double timeout = CFAbsoluteTimeGetCurrent() + 1.0;

        while (CFAbsoluteTimeGetCurrent() < timeout)
        {
            ProcessRunLoop(true);

            size = [platform.view bounds].size;
            if ((size.width > size.height) != landscape) break;
        }
    }
}

// Load framebuffer for the current view size
// NOTE: The layer drawable size is the view size in points multiplied by the layer contents scale
static bool LoadFramebuffer(void)
{
    CGSize size = [platform.view bounds].size;
    CGFloat nativeScale = [[UIScreen mainScreen] nativeScale];

    CORE.Window.display.width = (int)round(size.width*nativeScale);
    CORE.Window.display.height = (int)round(size.height*nativeScale);

    SetupFramebuffer(CORE.Window.display.width, CORE.Window.display.height);

    int width = CORE.Window.screen.width + CORE.Window.renderOffset.x;

#if defined(GRAPHICS_API_OPENGL_SOFTWARE)
    platform.drawableWidth = width;
    platform.drawableHeight = CORE.Window.screen.height + CORE.Window.renderOffset.y;
#else
    CAEAGLLayer *layer = (CAEAGLLayer *)[platform.view layer];
    [platform.view setContentScaleFactor:((CGFloat)width + 0.5)/size.width];

    glGenFramebuffers(1, &platform.framebuffer);
    glBindFramebuffer(GL_FRAMEBUFFER, platform.framebuffer);

    glGenRenderbuffers(1, &platform.colorRenderbuffer);
    glBindRenderbuffer(GL_RENDERBUFFER, platform.colorRenderbuffer);

    if (![platform.context renderbufferStorage:GL_RENDERBUFFER fromDrawable:layer])
    {
        TRACELOG(LOG_WARNING, "EAGL: Failed to create framebuffer storage from the view layer");
        return false;
    }

    glFramebufferRenderbuffer(GL_FRAMEBUFFER, GL_COLOR_ATTACHMENT0, GL_RENDERBUFFER, platform.colorRenderbuffer);
    glGetRenderbufferParameteriv(GL_RENDERBUFFER, GL_RENDERBUFFER_WIDTH, &platform.drawableWidth);
    glGetRenderbufferParameteriv(GL_RENDERBUFFER, GL_RENDERBUFFER_HEIGHT, &platform.drawableHeight);

    // NOTE: Depth and stencil buffers are attached to the framebuffer used for drawing
    bool msaa = FLAG_IS_SET(CORE.Window.flags, FLAG_MSAA_4X_HINT);

    if (msaa)
    {
        glGenFramebuffers(1, &platform.msaaFramebuffer);
        glBindFramebuffer(GL_FRAMEBUFFER, platform.msaaFramebuffer);

        glGenRenderbuffers(1, &platform.msaaColorRenderbuffer);
        glBindRenderbuffer(GL_RENDERBUFFER, platform.msaaColorRenderbuffer);
    #if defined(GRAPHICS_API_OPENGL_ES3)
        glRenderbufferStorageMultisample(GL_RENDERBUFFER, 4, GL_RGBA8, platform.drawableWidth, platform.drawableHeight);
    #else
        glRenderbufferStorageMultisampleAPPLE(GL_RENDERBUFFER, 4, GL_RGBA8_OES, platform.drawableWidth, platform.drawableHeight);
    #endif
        glFramebufferRenderbuffer(GL_FRAMEBUFFER, GL_COLOR_ATTACHMENT0, GL_RENDERBUFFER, platform.msaaColorRenderbuffer);
    }

    glGenRenderbuffers(1, &platform.depthRenderbuffer);
    glBindRenderbuffer(GL_RENDERBUFFER, platform.depthRenderbuffer);

    if (!msaa) glRenderbufferStorage(GL_RENDERBUFFER, GL_DEPTH24_STENCIL8_PACKED, platform.drawableWidth, platform.drawableHeight);
    #if defined(GRAPHICS_API_OPENGL_ES3)
    else glRenderbufferStorageMultisample(GL_RENDERBUFFER, 4, GL_DEPTH24_STENCIL8_PACKED, platform.drawableWidth, platform.drawableHeight);
    #else
    else glRenderbufferStorageMultisampleAPPLE(GL_RENDERBUFFER, 4, GL_DEPTH24_STENCIL8_PACKED, platform.drawableWidth, platform.drawableHeight);
    #endif

    glFramebufferRenderbuffer(GL_FRAMEBUFFER, GL_DEPTH_ATTACHMENT, GL_RENDERBUFFER, platform.depthRenderbuffer);
    glFramebufferRenderbuffer(GL_FRAMEBUFFER, GL_STENCIL_ATTACHMENT, GL_RENDERBUFFER, platform.depthRenderbuffer);

    if (glCheckFramebufferStatus(GL_FRAMEBUFFER) != GL_FRAMEBUFFER_COMPLETE)
    {
        TRACELOG(LOG_WARNING, "EAGL: Framebuffer is not complete [STATUS: 0x%x]", glCheckFramebufferStatus(GL_FRAMEBUFFER));
        return false;
    }

    eaglDefaultFramebuffer = msaa? platform.msaaFramebuffer : platform.framebuffer;
#endif

    // NOTE: Layer drawable size is rounded by the system, offsets follow the real size
    CORE.Window.renderOffset.x = platform.drawableWidth - CORE.Window.screen.width;
    CORE.Window.renderOffset.y = platform.drawableHeight - CORE.Window.screen.height;
    CORE.Window.render.width = CORE.Window.screen.width;
    CORE.Window.render.height = CORE.Window.screen.height;
    CORE.Window.currentFbo.width = CORE.Window.render.width;
    CORE.Window.currentFbo.height = CORE.Window.render.height;

    return true;
}

#if !defined(GRAPHICS_API_OPENGL_SOFTWARE)
// Resolve multisampled default framebuffer into the presented one, bound on return
// NOTE: Returns 0 if the default framebuffer is not multisampled
int eaglResolveDefaultFramebuffer(void)
{
    if (platform.msaaFramebuffer == 0) return 0;

    #if defined(GRAPHICS_API_OPENGL_ES3)
    glBindFramebuffer(GL_READ_FRAMEBUFFER, platform.msaaFramebuffer);
    glBindFramebuffer(GL_DRAW_FRAMEBUFFER, platform.framebuffer);
    glBlitFramebuffer(0, 0, platform.drawableWidth, platform.drawableHeight, 0, 0, platform.drawableWidth, platform.drawableHeight, GL_COLOR_BUFFER_BIT, GL_NEAREST);
    #else
    glBindFramebuffer(GL_READ_FRAMEBUFFER_APPLE, platform.msaaFramebuffer);
    glBindFramebuffer(GL_DRAW_FRAMEBUFFER_APPLE, platform.framebuffer);
    glResolveMultisampleFramebufferAPPLE();
    #endif

    glBindFramebuffer(GL_FRAMEBUFFER, platform.framebuffer);

    return 1;
}
#endif

// Unload framebuffer
static void UnloadFramebuffer(void)
{
#if !defined(GRAPHICS_API_OPENGL_SOFTWARE)
    if (platform.framebuffer != 0) glDeleteFramebuffers(1, &platform.framebuffer);
    if (platform.colorRenderbuffer != 0) glDeleteRenderbuffers(1, &platform.colorRenderbuffer);
    if (platform.depthRenderbuffer != 0) glDeleteRenderbuffers(1, &platform.depthRenderbuffer);
    if (platform.msaaFramebuffer != 0) glDeleteFramebuffers(1, &platform.msaaFramebuffer);
    if (platform.msaaColorRenderbuffer != 0) glDeleteRenderbuffers(1, &platform.msaaColorRenderbuffer);

    platform.framebuffer = 0;
    platform.colorRenderbuffer = 0;
    platform.depthRenderbuffer = 0;
    platform.msaaFramebuffer = 0;
    platform.msaaColorRenderbuffer = 0;
    eaglDefaultFramebuffer = 0;
#endif
}

// Setup framebuffer and screen sizes for a display size
// NOTE: Requested screen size is extended to the display aspect ratio (border bars),
// the framebuffer is scaled to the display by the system
static void SetupFramebuffer(int width, int height)
{
    if ((CORE.Window.screen.width == 0) || (CORE.Window.screen.height == 0)) platform.displaySize = true;

    if (platform.displaySize)
    {
        CORE.Window.screen.width = width;
        CORE.Window.screen.height = height;
    }

    float displayRatio = (float)width/(float)height;
    float screenRatio = (float)CORE.Window.screen.width/(float)CORE.Window.screen.height;

    if (displayRatio <= screenRatio)
    {
        CORE.Window.renderOffset.x = 0;
        CORE.Window.renderOffset.y = (int)round((float)CORE.Window.screen.width/displayRatio) - CORE.Window.screen.height;
    }
    else
    {
        CORE.Window.renderOffset.x = (int)round((float)CORE.Window.screen.height*displayRatio) - CORE.Window.screen.width;
        CORE.Window.renderOffset.y = 0;
    }

    CORE.Window.render.width = CORE.Window.screen.width;
    CORE.Window.render.height = CORE.Window.screen.height;
}

// Update framebuffer after a view size change (rotation)
static void UpdateFramebufferSize(void)
{
    if (!platform.initialized) return;

#if !defined(GRAPHICS_API_OPENGL_SOFTWARE)
    [EAGLContext setCurrentContext:platform.context];
#endif

    UnloadFramebuffer();

    if (!LoadFramebuffer())
    {
        TRACELOG(LOG_WARNING, "EAGL: Failed to resize framebuffer");
        return;
    }

    // Reset viewport and projection matrix for new size
    SetupViewport(CORE.Window.render.width, CORE.Window.render.height);

    CORE.Window.resizedLastFrame = true;
    platform.eventReceived = true;
}

// Process pending run loop events, optionally waiting for one
static void ProcessRunLoop(bool wait)
{
    if (wait) CFRunLoopRunInMode(kCFRunLoopDefaultMode, 1.0, true);

    while (CFRunLoopRunInMode(kCFRunLoopDefaultMode, 0.0, true) == kCFRunLoopRunHandledSource) { }
}

// Register touch points change
static void HandleTouches(NSSet *touches, int action)
{
    platform.eventReceived = true;

    CGSize size = [platform.view bounds].size;
    float scaleX = (float)platform.drawableWidth/(float)size.width;
    float scaleY = (float)platform.drawableHeight/(float)size.height;

    if (action == TOUCH_ACTION_DOWN)
    {
        for (UITouch *touch in touches)
        {
            for (int i = 0; i < MAX_TOUCH_POINTS; i++)
            {
                if (platform.touches[i] == nil)
                {
                    platform.touches[i] = touch;
                    break;
                }
            }
        }
    }

#if SUPPORT_GESTURES_SYSTEM
    // NOTE: Released touch points are still reported to the gestures system
    GestureEvent gestureEvent = { 0 };
    gestureEvent.touchAction = action;

    for (int i = 0; i < MAX_TOUCH_POINTS; i++)
    {
        if (platform.touches[i] == nil) continue;

        CGPoint point = [platform.touches[i] locationInView:platform.view];

        gestureEvent.pointId[gestureEvent.pointCount] = i;
        gestureEvent.position[gestureEvent.pointCount].x = ((float)point.x*scaleX - (float)CORE.Window.renderOffset.x/2)/(float)GetScreenWidth();
        gestureEvent.position[gestureEvent.pointCount].y = ((float)point.y*scaleY - (float)CORE.Window.renderOffset.y/2)/(float)GetScreenHeight();
        gestureEvent.pointCount++;
    }

    // Gesture data is sent to gestures-system for processing
    ProcessGestureEvent(gestureEvent);
#endif

    if ((action == TOUCH_ACTION_UP) || (action == TOUCH_ACTION_CANCEL))
    {
        for (UITouch *touch in touches)
        {
            for (int i = 0; i < MAX_TOUCH_POINTS; i++)
            {
                if (platform.touches[i] == touch) platform.touches[i] = nil;
            }
        }
    }

    // Register active touch points
    int count = 0;
    for (int i = 0; i < MAX_TOUCH_POINTS; i++)
    {
        if (platform.touches[i] == nil) continue;

        CGPoint point = [platform.touches[i] locationInView:platform.view];

        CORE.Input.Touch.pointId[count] = i;
        CORE.Input.Touch.position[count].x = (float)point.x*scaleX - (float)CORE.Window.renderOffset.x/2;
        CORE.Input.Touch.position[count].y = (float)point.y*scaleY - (float)CORE.Window.renderOffset.y/2;
        count++;
    }

    CORE.Input.Touch.pointCount = count;

    // NOTE: Touch is mapped to the mouse left button, first touch point to the mouse position
    CORE.Input.Touch.currentTouchState[MOUSE_BUTTON_LEFT] = (count > 0);

    if (count > 0) CORE.Input.Mouse.currentPosition = CORE.Input.Touch.position[0];
    if (action == TOUCH_ACTION_DOWN) CORE.Input.Mouse.previousPosition = CORE.Input.Mouse.currentPosition;
}

// Register hardware keyboard keys change
// NOTE: Held keys are not repeated by the system
static void HandlePresses(NSSet *presses, bool pressed)
{
    platform.eventReceived = true;

    for (UIPress *press in presses)
    {
        UIKey *pressKey = [press key];
        if (pressKey == nil) continue;

        long usage = (long)[pressKey keyCode];
        int key = ((usage >= 0) && (usage < KEYMAP_SIZE))? hidToRaylibMap[usage] : (int)KEY_NULL;

        if ((key > KEY_NULL) && (key < MAX_KEYBOARD_KEYS))
        {
            CORE.Input.Keyboard.currentKeyState[key] = pressed;

            if (pressed)
            {
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

        if (!pressed || ([pressKey modifierFlags] & UIKeyModifierCommand)) continue;

        // Characters input
        // NOTE: Non-printable keys report control characters or a name, like "UIKeyInputUpArrow"
        NSString *characters = [pressKey characters];
        NSUInteger length = [characters length];
        if ((length == 0) || (length > 2)) continue;

        unsigned int codepoint = [characters characterAtIndex:0];

        // Combine UTF-16 surrogate pairs, any other two units string is not a character
        if (length == 2)
        {
            unsigned int low = [characters characterAtIndex:1];
            if ((codepoint < 0xD800) || (codepoint > 0xDBFF) || (low < 0xDC00) || (low > 0xDFFF)) continue;

            codepoint = 0x10000 + ((codepoint - 0xD800) << 10) + (low - 0xDC00);
        }

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

// Register gamepads connection and state
// NOTE: Controllers are mapped by the system to the extended gamepad profile (Xbox layout),
// the state is updated by the main run loop and it is only received while the application is active
static void PollGamepads(void)
{
    NSArray<GCController *> *controllers = [GCController controllers];

    // Unregister disconnected gamepads
    for (int i = 0; i < MAX_GAMEPADS; i++)
    {
        if ((platform.gamepads[i] != nil) && ![controllers containsObject:platform.gamepads[i]])
        {
            [platform.gamepads[i] release];
            platform.gamepads[i] = nil;

            CORE.Input.Gamepad.ready[i] = false;
            memset(CORE.Input.Gamepad.name[i], 0, MAX_GAMEPAD_NAME_LENGTH);
            memset(CORE.Input.Gamepad.currentButtonState[i], 0, sizeof(CORE.Input.Gamepad.currentButtonState[i]));
            memset(CORE.Input.Gamepad.axisState[i], 0, sizeof(CORE.Input.Gamepad.axisState[i]));
        }
    }

    // Register connected gamepads, in the first free slot
    for (GCController *controller in controllers)
    {
        if ([controller extendedGamepad] == nil) continue;

        bool registered = false;
        for (int i = 0; i < MAX_GAMEPADS; i++)
        {
            if (platform.gamepads[i] == controller) registered = true;
        }

        for (int i = 0; !registered && (i < MAX_GAMEPADS); i++)
        {
            if (platform.gamepads[i] != nil) continue;

            platform.gamepads[i] = [controller retain];
            registered = true;

            const char *name = [[controller vendorName] UTF8String];
            snprintf(CORE.Input.Gamepad.name[i], MAX_GAMEPAD_NAME_LENGTH, "%s", (name != NULL)? name : "Gamepad");

            CORE.Input.Gamepad.ready[i] = true;
            CORE.Input.Gamepad.axisCount[i] = 6;
        }
    }

    // Register gamepads state
    for (int i = 0; i < MAX_GAMEPADS; i++)
    {
        if (platform.gamepads[i] == nil) continue;

        GCExtendedGamepad *gamepad = [platform.gamepads[i] extendedGamepad];

        // NOTE: Options, home and thumbstick buttons are not available on all controllers (nil inputs are not pressed),
        // trigger buttons are registered from the triggers axis
        GCControllerButtonInput *inputs[MAX_GAMEPAD_BUTTONS] = { nil };
        inputs[GAMEPAD_BUTTON_LEFT_FACE_UP] = [[gamepad dpad] up];
        inputs[GAMEPAD_BUTTON_LEFT_FACE_RIGHT] = [[gamepad dpad] right];
        inputs[GAMEPAD_BUTTON_LEFT_FACE_DOWN] = [[gamepad dpad] down];
        inputs[GAMEPAD_BUTTON_LEFT_FACE_LEFT] = [[gamepad dpad] left];
        inputs[GAMEPAD_BUTTON_RIGHT_FACE_UP] = [gamepad buttonY];
        inputs[GAMEPAD_BUTTON_RIGHT_FACE_RIGHT] = [gamepad buttonB];
        inputs[GAMEPAD_BUTTON_RIGHT_FACE_DOWN] = [gamepad buttonA];
        inputs[GAMEPAD_BUTTON_RIGHT_FACE_LEFT] = [gamepad buttonX];
        inputs[GAMEPAD_BUTTON_LEFT_TRIGGER_1] = [gamepad leftShoulder];
        inputs[GAMEPAD_BUTTON_RIGHT_TRIGGER_1] = [gamepad rightShoulder];
        inputs[GAMEPAD_BUTTON_MIDDLE_LEFT] = [gamepad buttonOptions];
        inputs[GAMEPAD_BUTTON_MIDDLE] = [gamepad buttonHome];
        inputs[GAMEPAD_BUTTON_MIDDLE_RIGHT] = [gamepad buttonMenu];
        inputs[GAMEPAD_BUTTON_LEFT_THUMB] = [gamepad leftThumbstickButton];
        inputs[GAMEPAD_BUTTON_RIGHT_THUMB] = [gamepad rightThumbstickButton];

        for (int k = 0; k < MAX_GAMEPAD_BUTTONS; k++)
        {
            if ([inputs[k] isPressed])
            {
                CORE.Input.Gamepad.currentButtonState[i][k] = 1;
                CORE.Input.Gamepad.lastButtonPressed = k;
            }
            else CORE.Input.Gamepad.currentButtonState[i][k] = 0;
        }

        // NOTE: Thumbsticks vertical axis is positive up, raylib expects it positive down
        // Triggers range is [0..1], raylib expects [-1..1]
        CORE.Input.Gamepad.axisState[i][GAMEPAD_AXIS_LEFT_X] = [[[gamepad leftThumbstick] xAxis] value];
        CORE.Input.Gamepad.axisState[i][GAMEPAD_AXIS_LEFT_Y] = -[[[gamepad leftThumbstick] yAxis] value];
        CORE.Input.Gamepad.axisState[i][GAMEPAD_AXIS_RIGHT_X] = [[[gamepad rightThumbstick] xAxis] value];
        CORE.Input.Gamepad.axisState[i][GAMEPAD_AXIS_RIGHT_Y] = -[[[gamepad rightThumbstick] yAxis] value];
        CORE.Input.Gamepad.axisState[i][GAMEPAD_AXIS_LEFT_TRIGGER] = [[gamepad leftTrigger] value]*2.0f - 1.0f;
        CORE.Input.Gamepad.axisState[i][GAMEPAD_AXIS_RIGHT_TRIGGER] = [[gamepad rightTrigger] value]*2.0f - 1.0f;

        // Register buttons for pressure triggers, same threshold as the GLFW platform
        for (int k = GAMEPAD_AXIS_LEFT_TRIGGER; k <= GAMEPAD_AXIS_RIGHT_TRIGGER; k++)
        {
            int button = (k == GAMEPAD_AXIS_LEFT_TRIGGER)? GAMEPAD_BUTTON_LEFT_TRIGGER_2 : GAMEPAD_BUTTON_RIGHT_TRIGGER_2;

            if (CORE.Input.Gamepad.axisState[i][k] > 0.1f)
            {
                CORE.Input.Gamepad.currentButtonState[i][button] = 1;
                CORE.Input.Gamepad.lastButtonPressed = button;
            }
            else CORE.Input.Gamepad.currentButtonState[i][button] = 0;
        }
    }
}

// Send trace log messages to the system log
// NOTE: Standard output is not visible on devices
static void TraceLogIOS(int logType, const char *text, va_list args)
{
    char buffer[MAX_TRACELOG_MSG_LENGTH] = { 0 };
    const char *level = "";
    os_log_type_t type = OS_LOG_TYPE_DEFAULT;

    switch (logType)
    {
        case LOG_TRACE: level = "TRACE: "; break;
        case LOG_DEBUG: level = "DEBUG: "; break;
        case LOG_INFO: level = "INFO: "; break;
        case LOG_WARNING: level = "WARNING: "; type = OS_LOG_TYPE_ERROR; break;
        case LOG_ERROR: level = "ERROR: "; type = OS_LOG_TYPE_ERROR; break;
        case LOG_FATAL: level = "FATAL: "; type = OS_LOG_TYPE_FAULT; break;
        default: break;
    }

    vsnprintf(buffer, sizeof(buffer), text, args);
    os_log_with_type(OS_LOG_DEFAULT, type, "[raylib] %{public}s%{public}s", level, buffer);

    if (logType == LOG_FATAL) exit(EXIT_FAILURE);
}

#if !defined(GRAPHICS_API_OPENGL_SOFTWARE)
// Get OpenGL ES procedure address
static void *GetProcAddressIOS(const char *name)
{
    return dlsym(RTLD_DEFAULT, name);
}
#endif

//----------------------------------------------------------------------------------
// Module Internal Classes Implementation
//----------------------------------------------------------------------------------
@implementation RaylibApplicationDelegate

- (UISceneConfiguration *)application:(UIApplication *)application configurationForConnectingSceneSession:(UISceneSession *)session options:(UISceneConnectionOptions *)options
{
    UISceneConfiguration *configuration = [[[UISceneConfiguration alloc] initWithName:@"Default Configuration" sessionRole:[session role]] autorelease];
    [configuration setDelegateClass:[RaylibSceneDelegate class]];

    return configuration;
}

- (void)applicationWillTerminate:(UIApplication *)application
{
    platform.closeRequested = true;
}

@end

@implementation RaylibSceneDelegate

- (void)scene:(UIScene *)scene willConnectToSession:(UISceneSession *)session options:(UISceneConnectionOptions *)options
{
    if (platform.window != nil) return;

    platform.window = [[UIWindow alloc] initWithWindowScene:(UIWindowScene *)scene];
    platform.viewController = [[RaylibViewController alloc] init];
    platform.view = [[RaylibView alloc] initWithFrame:[platform.window bounds]];

    [platform.view setAutoresizingMask:(UIViewAutoresizingFlexibleWidth | UIViewAutoresizingFlexibleHeight)];
    [platform.view setMultipleTouchEnabled:YES];
    [platform.view setBackgroundColor:[UIColor blackColor]];
    [platform.viewController setView:platform.view];
    [platform.window setRootViewController:platform.viewController];
    [platform.window makeKeyAndVisible];
    [platform.view becomeFirstResponder];

    platform.displayLinkTarget = [[RaylibDisplayLinkTarget alloc] init];
    platform.displayLink = [[CADisplayLink displayLinkWithTarget:platform.displayLinkTarget selector:@selector(displayRefresh:)] retain];
    [platform.displayLink addToRunLoop:[NSRunLoop mainRunLoop] forMode:NSRunLoopCommonModes];

    // Run the program from the main run loop, after the application launch is completed
    [self performSelector:@selector(runProgram) withObject:nil afterDelay:0.0];
}

- (void)runProgram
{
    ios_run();

    exit(0);
}

- (void)sceneDidBecomeActive:(UIScene *)scene
{
    platform.appEnabled = true;
    platform.eventReceived = true;
    [platform.displayLink setPaused:NO];
    FLAG_CLEAR(CORE.Window.flags, FLAG_WINDOW_UNFOCUSED);
}

- (void)sceneWillResignActive:(UIScene *)scene
{
    platform.appEnabled = false;
    [platform.displayLink setPaused:YES];
    FLAG_SET(CORE.Window.flags, FLAG_WINDOW_UNFOCUSED);

    // NOTE: Touches are cancelled by the system when the application is not active, key releases are not reported
    memset(CORE.Input.Keyboard.currentKeyState, 0, sizeof(CORE.Input.Keyboard.currentKeyState));
    memset(CORE.Input.Touch.currentTouchState, 0, sizeof(CORE.Input.Touch.currentTouchState));
    memset(platform.touches, 0, sizeof(platform.touches));
    CORE.Input.Touch.pointCount = 0;

#if !defined(GRAPHICS_API_OPENGL_SOFTWARE)
    // Finish pending drawing commands, OpenGL ES is not allowed in background
    if (platform.context != nil) glFinish();
#endif
}

@end

@implementation RaylibViewController

- (BOOL)prefersStatusBarHidden
{
    return YES;
}

- (BOOL)prefersHomeIndicatorAutoHidden
{
    return YES;
}

- (UIRectEdge)preferredScreenEdgesDeferringSystemGestures
{
    return UIRectEdgeAll;
}

- (UIInterfaceOrientationMask)supportedInterfaceOrientations
{
    return (platform.orientations != 0)? platform.orientations : UIInterfaceOrientationMaskAll;
}

@end

@implementation RaylibView

#if !defined(GRAPHICS_API_OPENGL_SOFTWARE)
+ (Class)layerClass
{
    return [CAEAGLLayer class];
}
#endif

- (void)layoutSubviews
{
    [super layoutSubviews];

    CGSize size = [self bounds].size;
    CGFloat nativeScale = [[UIScreen mainScreen] nativeScale];

    if (((int)round(size.width*nativeScale) != (int)CORE.Window.display.width) || ((int)round(size.height*nativeScale) != (int)CORE.Window.display.height)) UpdateFramebufferSize();
}

- (BOOL)canBecomeFirstResponder
{
    return YES;
}

- (void)pressesBegan:(NSSet<UIPress *> *)presses withEvent:(UIPressesEvent *)event
{
    HandlePresses(presses, true);
}

- (void)pressesEnded:(NSSet<UIPress *> *)presses withEvent:(UIPressesEvent *)event
{
    HandlePresses(presses, false);
}

- (void)pressesCancelled:(NSSet<UIPress *> *)presses withEvent:(UIPressesEvent *)event
{
    HandlePresses(presses, false);
}

- (void)touchesBegan:(NSSet<UITouch *> *)touches withEvent:(UIEvent *)event
{
    HandleTouches(touches, TOUCH_ACTION_DOWN);
}

- (void)touchesMoved:(NSSet<UITouch *> *)touches withEvent:(UIEvent *)event
{
    HandleTouches(touches, TOUCH_ACTION_MOVE);
}

- (void)touchesEnded:(NSSet<UITouch *> *)touches withEvent:(UIEvent *)event
{
    HandleTouches(touches, TOUCH_ACTION_UP);
}

- (void)touchesCancelled:(NSSet<UITouch *> *)touches withEvent:(UIEvent *)event
{
    HandleTouches(touches, TOUCH_ACTION_CANCEL);
}

@end

@implementation RaylibDisplayLinkTarget

- (void)displayRefresh:(CADisplayLink *)link
{
    platform.displayRefreshed = true;
}

@end
