#include "raylib.h"
#include "utils_log.h"
#include <stdio.h>                      // Required for: vprintf()
#include <string.h>                     // Required for: strcpy(), strcat()

#define MAX_TRACELOG_BUFFER_SIZE   128  // As defined in utils.c from raylib

void rayLogWrapperCallback(int logType, const char *text, va_list args) {
	char buffer[MAX_TRACELOG_BUFFER_SIZE] = { 0 };

	vsprintf(buffer, text, args);

	internalTraceLogCallbackGo(logType, buffer, strlen(buffer));
}

void setLogCallbackWrapper(void) {
	SetTraceLogCallback(rayLogWrapperCallback);
}
