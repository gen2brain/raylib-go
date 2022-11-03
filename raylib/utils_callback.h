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
