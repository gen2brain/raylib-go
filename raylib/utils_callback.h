#if defined(__cplusplus)
extern "C" {            // Prevents name mangling of functions
#endif

typedef unsigned char * cchar_t;


// unsigned char * rayLoadFileDataCallback(const char *, unsigned int *);                 // enable the call-back

void setLoadFileDataCallbackWrapper(void);

#if defined(__cplusplus)
}
#endif
