package gphoto2

// #include <stdlib.h>
import "C"

import "unsafe"

type freeFunc func()

func goString(charPtr *C.char) string {
	return C.GoString((*C.char)(charPtr))
}

func cString(str string) (*C.char, freeFunc) {
	c_str := C.CString(str)
	return c_str, func() {
		C.free(unsafe.Pointer(c_str))
	}
}

func goBytes(charPtr *C.char, size int) []byte {
	return C.GoBytes(unsafe.Pointer(charPtr), C.int(size))
}

func cBytes(bytes []byte) (*C.char, freeFunc) {
	c_bytes := C.CBytes(bytes)
	return (*C.char)(c_bytes), func() {
		C.free(unsafe.Pointer(c_bytes))
	}
}

func cMalloc(size int) (*C.char, freeFunc) {
	c_buf := (*C.char)(C.malloc(C.size_t(size)))
	return c_buf, func() {
		C.free(unsafe.Pointer(c_buf))
	}
}
