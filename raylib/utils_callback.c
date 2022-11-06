#include "raylib.h"
#include "utils_callback.h"
#include <string.h>                     // Required for: strcpy(), strcat()
#include <stdlib.h>

unsigned char * rayLoadFileDataCallback(const char *fileName, unsigned int *bytesRead) {
	unsigned char ** ref = (unsigned char **)malloc(sizeof(unsigned char *));
	loadFileDataCallbackGo(fileName, strlen(fileName), bytesRead, ref);
	if (*bytesRead == -1) {
		TRACELOG(LOG_WARNING, "FILE: [%s] failed to load binary file", fileName);
		return NULL;
	}
	unsigned char * p = (unsigned char *)malloc(*bytesRead);
	memcpy(p, *ref, *bytesRead);
	free(ref);
	return p;
}

void setLoadFileDataCallbackWrapper(void) {
	SetLoadFileDataCallback(rayLoadFileDataCallback);
}

void unsetLoadFileDataCallbackWrapper(void) {
	SetLoadFileDataCallback(NULL);
}

char * rayLoadFileTextCallback(const char *fileName) {
	char ** ref = (char **)malloc(sizeof(char *));
	int sz = 0;
	loadFileTextCallbackGo(fileName, strlen(fileName), &sz, ref);
	if (sz == -1) {
		TRACELOG(LOG_WARNING, "FILE: [%s] failed to load text file", fileName);
		return NULL;
	}
	char * refref = *ref;
	free(ref);
	return refref;
}

void setLoadFileTextCallbackWrapper(void) {
	SetLoadFileTextCallback(rayLoadFileTextCallback);
}

void unsetLoadFileTextCallbackWrapper(void) {
	SetLoadFileTextCallback(NULL);
}