package gphoto2

// #cgo pkg-config: libgphoto2
// #include <gphoto2.h>
// #include <stdlib.h>
import "C"

type Context struct {
    c_ref *C.GPContext
}

func NewContext() *Context {
    return &Context{
        c_ref: C.gp_context_new(),
    }
}

func (c *Context) Cancel() {
    // TODO: Handle result
    C.gp_context_cancel(c.c_ref)
    //if err := toError(C.gp_context_cancel(c.c_ref)); err != nil {
    //    panic(err)
    //}
}
