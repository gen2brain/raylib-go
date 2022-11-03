#include "raylib.h"
#include "utils_callback.h"
#include <string.h>                     // Required for: strcpy(), strcat()
#include <stdlib.h>

unsigned char * rayLoadFileDataCallback(const char *fileName, unsigned int *bytesRead) {
	unsigned char ** ref = (unsigned char **)malloc(sizeof(unsigned char *));
	loadFileDataCallbackGo(fileName, strlen(fileName), bytesRead, ref);
	unsigned char * p = (unsigned char *)malloc(*bytesRead);
	memcpy(p, *ref, *bytesRead);
	free(ref);
	return p;
}

void setLoadFileDataCallbackWrapper(void) {
	SetLoadFileDataCallback(rayLoadFileDataCallback);
}

char * rayLoadFileTextCallback(const char *fileName) {
	char ** ref = (char **)malloc(sizeof(char *));
	int sz = 0;
	loadFileTextCallbackGo(fileName, strlen(fileName), &sz, ref);
	char * refref = *ref;
	// unsigned char * p = (unsigned char *)malloc(*bytesRead);
	// memcpy(p, *ref, *bytesRead);
	free(ref);
	return refref;
	// return p;
}

void setLoadFileTextCallbackWrapper(void) {
	SetLoadFileTextCallback(rayLoadFileTextCallback);
}