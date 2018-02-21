// +build android

#include "_cgo_export.h"
#include <android/log.h>

void log_info(const char *msg) {
    __android_log_print(ANDROID_LOG_INFO, "raylib", msg);
}

void log_warn(const char *msg) {
    __android_log_print(ANDROID_LOG_WARN, "raylib", msg);
}

void log_error(const char *msg) {
    __android_log_print(ANDROID_LOG_ERROR, "raylib", msg);
}

void log_debug(const char *msg) {
    __android_log_print(ANDROID_LOG_DEBUG, "raylib", msg);
}

const char* get_internal_storage_path() {
    return internal_storage_path;
}
