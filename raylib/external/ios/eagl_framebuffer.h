// EAGL has no default framebuffer: framebuffer name 0 is mapped to the framebuffer created by the platform layer
#ifndef EAGL_FRAMEBUFFER_H
#define EAGL_FRAMEBUFFER_H

extern unsigned int eaglDefaultFramebuffer;

// Resolve the default framebuffer and bind the resolved one, returns 0 if it is not multisampled
extern int eaglResolveDefaultFramebuffer(void);

static inline void glBindFramebufferEAGL(GLenum target, GLuint framebuffer)
{
    glBindFramebuffer(target, (framebuffer == 0)? eaglDefaultFramebuffer : framebuffer);
}

static inline void glGetIntegervEAGL(GLenum pname, GLint *data)
{
    glGetIntegerv(pname, data);

    // NOTE: GL_READ_FRAMEBUFFER_BINDING (0x8CAA) is not defined by the OpenGL ES 2.0 header
    if (((pname == GL_FRAMEBUFFER_BINDING) || (pname == 0x8CAA)) && (*data == (GLint)eaglDefaultFramebuffer)) *data = 0;
}

// NOTE: Multisampled framebuffers can not be read, the default one is resolved first
static inline void glReadPixelsEAGL(GLint x, GLint y, GLsizei width, GLsizei height, GLenum format, GLenum type, GLvoid *pixels)
{
    GLint framebuffer = 0;
    glGetIntegerv(GL_FRAMEBUFFER_BINDING, &framebuffer);

    if ((framebuffer == (GLint)eaglDefaultFramebuffer) && eaglResolveDefaultFramebuffer())
    {
        glReadPixels(x, y, width, height, format, type, pixels);
        glBindFramebuffer(GL_FRAMEBUFFER, eaglDefaultFramebuffer);
    }
    else glReadPixels(x, y, width, height, format, type, pixels);
}

#define glBindFramebuffer glBindFramebufferEAGL
#define glGetIntegerv glGetIntegervEAGL
#define glReadPixels glReadPixelsEAGL

#endif // EAGL_FRAMEBUFFER_H
