#include "raylib.h"
#include "utils_callback.h"

const char * rayLoadFileDataCallback(const char *fileName, unsigned int *bytesRead) {
	return loadFileDataCallbackGo(fileName, strlen(fileName), bytesRead);
}

void setLoadFileDataCallbackWrapper(void) {
	SetLoadFileDataCallback(rayLoadFileDataCallback);
}
