// +build android

#include "_cgo_export.h"

struct android_app *GetAndroidApp();

void android_init() {
    struct android_app *app;
    app = GetAndroidApp();
    asset_manager = app->activity->assetManager;
    internal_storage_path = app->activity->internalDataPath;
}
