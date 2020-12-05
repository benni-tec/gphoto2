package gphoto2

// #cgo pkg-config: libgphoto2
// #include <gphoto2.h>
// #include <stdlib.h>
import "C"

import "unsafe"

type CameraEventType int

const (
	EventUnknown   CameraEventType = C.GP_EVENT_UNKNOWN
	EventTimeout   CameraEventType = C.GP_EVENT_TIMEOUT
	EventFileAdded CameraEventType = C.GP_EVENT_FILE_ADDED
)

type CameraEvent struct {
	Type   CameraEventType
	Folder string
	File   string
}

func goCameraEvent(voidPtr unsafe.Pointer, eventType C.CameraEventType) *CameraEvent {
	ce := new(CameraEvent)
	ce.Type = CameraEventType(eventType)

	if ce.Type == EventFileAdded {
		cameraFilePath := (*C.CameraFilePath)(voidPtr)
		ce.File = C.GoString((*C.char)(&cameraFilePath.name[0]))
		ce.Folder = C.GoString((*C.char)(&cameraFilePath.folder[0]))
	}

	return ce
}

func (cam *Camera) WaitForEvent(timeoutMS int) (*CameraEvent, error) {
	var eventType C.CameraEventType
	var vp unsafe.Pointer

	if err := toError(C.gp_camera_wait_for_event(
		cam.c_ref, C.int(timeoutMS), &eventType, &vp, cam.ctx.c_ref,
	)); err != nil {
		return nil, err
	} else {
        defer C.free(vp)

		return goCameraEvent(vp, eventType), nil
	}
}
