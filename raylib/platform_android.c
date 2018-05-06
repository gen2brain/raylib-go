// +build android

#include "_cgo_export.h"

void android_init() {
    struct android_app *app;
    app = GetAndroidApp();
    asset_manager = app->activity->assetManager;
    internal_storage_path = app->activity->internalDataPath;
}
