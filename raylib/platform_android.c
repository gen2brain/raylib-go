// +build android

#include "_cgo_export.h"

void android_main(struct android_app *app) {
    androidMain(app);
}

void init_asset_manager(void *state) {
    struct android_app *app;
	app = (struct android_app *)state;
	asset_manager = app->activity->assetManager;
}
