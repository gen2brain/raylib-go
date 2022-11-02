#if defined(__cplusplus)
extern "C" {            // Prevents name mangling of functions
#endif

typedef const char cchar_t;


const char * rayLoadFileDataCallback(const char *, unsigned int *);                 // enable the call-back

#if defined(__cplusplus)
}
#endif
