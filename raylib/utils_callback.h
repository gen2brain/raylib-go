#if defined(SUPPORT_TRACELOG)
    #define TRACELOG(level, ...) TraceLog(level, __VA_ARGS__)

    #if defined(SUPPORT_TRACELOG_DEBUG)
        #define TRACELOGD(...) TraceLog(LOG_DEBUG, __VA_ARGS__)
    #else
        #define TRACELOGD(...) (void)0
    #endif
#else
    #define TRACELOG(level, ...) (void)0
    #define TRACELOGD(...) (void)0
#endif

#if defined(__cplusplus)
extern "C" {            // Prevents name mangling of functions
#endif

typedef unsigned char * cchar_t;


void setLoadFileDataCallbackWrapper(void);
void unsetLoadFileDataCallbackWrapper(void);
void setLoadFileTextCallbackWrapper(void);
void unsetLoadFileTextCallbackWrapper(void);

#if defined(__cplusplus)
}
#endif
