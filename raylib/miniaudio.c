// The implementation of mini_al needs to #include windows.h which means it needs to go into
// it's own translation unit. Not doing this will cause conflicts with CloseWindow(), etc.
#define MAL_IMPLEMENTATION
#include "external/miniaudio.h"
