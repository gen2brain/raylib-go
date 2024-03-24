//go:build android

package main

/*

#cgo  LDFLAGS: -landroid

#include <android/native_activity.h>
#include <stdlib.h>

extern struct ANativeActivity *GetANativeActivity(void); // if you need android_app (jni) put and call this

static char* getJNI(uintptr_t *vmp, uintptr_t* envp, uintptr_t* ctx, int* attachedp) {
	JNIEnv* env;
	struct ANativeActivity * activity = GetANativeActivity();
	*attachedp = 0;
        JavaVM* vm = activity->vm;
	switch ((*vm)->GetEnv(vm, (void**)&env, JNI_VERSION_1_6)) {
	case JNI_OK:
		break;
	case JNI_EDETACHED:
		if ((*vm)->AttachCurrentThread(vm, &env, 0) != 0) {
			return "cannot attach to JVM";
		}
		*attachedp = 1;
		break;
	case JNI_EVERSION:
		return "bad JNI version";
	default:
		return "unknown JNI error from GetEnv";
	}

	*envp = (uintptr_t)env;
	*vmp = (uintptr_t)vm;
	*ctx = (uintptr_t)(activity->clazz);
	return NULL;
}

static char* checkException(uintptr_t jnienv) {
	jthrowable exc;
	JNIEnv* env = (JNIEnv*)jnienv;

	if (!(*env)->ExceptionCheck(env)) {
		return NULL;
	}

	exc = (*env)->ExceptionOccurred(env);
	(*env)->ExceptionClear(env);

	jclass clazz = (*env)->FindClass(env, "java/lang/Throwable");
	jmethodID toString = (*env)->GetMethodID(env, clazz, "toString", "()Ljava/lang/String;");
	jobject msgStr = (*env)->CallObjectMethod(env, exc, toString);
	return (char*)(*env)->GetStringUTFChars(env, msgStr, 0);
}

static void unlockJNI(uintptr_t vmp) {
	JavaVM *vm = (JavaVM *)vmp;
	(*vm)->DetachCurrentThread(vm);
}

const char * getInternalStoragePathMy(){
	return GetANativeActivity()->internalDataPath;
}

const char * getExternalStoragePathMy(){
	return GetANativeActivity()->externalDataPath;
}
*/
import "C"
import (
	"errors"
	"fmt"
	"runtime"
	"unsafe"
)

func doSpecific() {
	fmt.Println("calling GetANativeActivity()->internalDataPath: ", C.GoString(C.getInternalStoragePathMy())) // this is only example to show how work GetANativeActivity()
	fmt.Println("calling GetANativeActivity()->externalDataPath: ", C.GoString(C.getExternalStoragePathMy()))
	// fmt.Println("calling rl.HomeDir(): ", rl.HomeDir()) // this work too
	err := RunOnJVM(func(vm, env, ctx uintptr) error {
		// do anything with jni
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}

func RunOnJVM(fn func(vm, env, ctx uintptr) error) error { // this is show how get to C pointer to JNIEnv*,JavaVM* and class of activity  
	if fn == nil {
		return errors.New("fn is nil")
	}
	errch := make(chan error)
	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		var envp, vmp, ctxp C.uintptr_t
		attached := C.int(0)
		if errStr := C.getJNI(&vmp, &envp, &ctxp, &attached); errStr != nil {
			errch <- errors.New(C.GoString(errStr))
			return
		}
		if attached != 0 {
			defer C.unlockJNI(vmp)
		}
		if err := fn(uintptr(vmp), uintptr(envp), uintptr(ctxp)); err != nil {
			errch <- err
			return
		}

		if exc := C.checkException(envp); exc != nil {
			errch <- errors.New(C.GoString(exc))
			C.free(unsafe.Pointer(exc))
			return
		}
		errch <- nil
	}()
	return <-errch
}
