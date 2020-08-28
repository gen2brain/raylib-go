#if defined(__cplusplus)
extern "C" {            // Prevents name mangling of functions
#endif


void setLogCallbackWrapper(void);                 // enable the call-back
void internalTraceLogCallbackGo(int, void*, int); // Go function that will get called

#if defined(__cplusplus)
}
#endif
