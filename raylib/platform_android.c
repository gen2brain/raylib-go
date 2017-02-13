// +build android

#include "_cgo_export.h"

void android_main(struct android_app *app) {
    androidMain(app);
}

void android_init(void *state) {
    struct android_app *app;
    app = (struct android_app *)state;
    asset_manager = app->activity->assetManager;
    internal_storage_path = app->activity->internalDataPath;
}
