package gphoto2

// #cgo pkg-config: libgphoto2
// #include <gphoto2.h>
// #include <stdlib.h>
import "C"

import "unsafe"

type CameraList struct {
	c_ref *C.CameraList
}

func newCameraList() CameraList {
	cl := CameraList{}

	if err := toError(C.gp_list_new(&cl.c_ref)); err != nil {
		panic(err)
	}

	return cl
}

func (cl CameraList) Close() {
    C.free(unsafe.Pointer(cl.c_ref))

    // FIXME: Why this fails? "pointer being freed was not allocated"
	//if err := toError(C.gp_list_free(cl.c_ref)); err != nil {
	//	panic(err)
	//}
}

func (cl CameraList) each(fn func(key, value string)) {
	var (
		size = int(C.gp_list_count(cl.c_ref))
	)

	if size < 0 {
		return
	}

	for i := 0; i < size; i++ {
		(func() {
			var c_key *C.char
			var c_val *C.char

			C.gp_list_get_name(cl.c_ref, C.int(i), &c_key)
			C.gp_list_get_value(cl.c_ref, C.int(i), &c_val)
			defer C.free(unsafe.Pointer(c_key))
			defer C.free(unsafe.Pointer(c_val))
			key := C.GoString(c_key)
			val := C.GoString(c_val)

			fn(key, val)
		})()
	}
}

func (cl CameraList) ToMap() map[string]string {
	vals := map[string]string{}

	cl.each(func(key, value string) {
		vals[key] = value
	})

	return vals
}

func (cl CameraList) Keys() []string {
	var keys []string

	cl.each(func(key, value string) {
		keys = append(keys, key)
	})

	return keys
}
