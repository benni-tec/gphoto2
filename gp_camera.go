package gphoto2

// #cgo pkg-config: libgphoto2
// #include <gphoto2.h>
// #include <stdlib.h>
import "C"

import (
	"path"
	"strings"
)

type CaptureType int

const (
	CaptureImage CaptureType = C.GP_CAPTURE_IMAGE
	CaptureMovie CaptureType = C.GP_CAPTURE_MOVIE
	CaptureSound CaptureType = C.GP_CAPTURE_SOUND
)

type Camera struct {
	c_ref *C.Camera

	ctx *Context
}

func NewCamera(ctx *Context) *Camera {
	cam := &Camera{
		ctx: ctx,
	}

	if err := toError(C.gp_camera_new(&cam.c_ref)); err != nil {
		panic(err)
	}

	return cam
}

func (cam *Camera) Init() error {
	if err := toError(C.gp_camera_init(cam.c_ref, cam.ctx.c_ref)); err != nil {
		return err
	}

	return nil
}

func (cam *Camera) TriggerCapture() error {
	return toError(C.gp_camera_trigger_capture(cam.c_ref, cam.ctx.c_ref))
}

func (cam *Camera) Capture() (*CameraFilePath, error) {
	var c_path C.CameraFilePath

	if err := toError(C.gp_camera_capture(cam.c_ref, C.CameraCaptureType(CaptureImage), &c_path, cam.ctx.c_ref)); err != nil {
		return nil, err
	}

	cf := newCameraFilePath(cam,
		C.GoString(&c_path.folder[0]),
		C.GoString(&c_path.name[0]),
	)

	return &cf, nil
}

func (cam *Camera) CapturePreview() (*CameraFile, error) {
	cf := newCameraFileAllocate(cam)

	if err := toError(C.gp_camera_capture_preview(cam.c_ref, cf.c_ref, cam.ctx.c_ref)); err != nil {
		cf.Close()
		return nil, err
	}

	return cf, nil
}

func (cam *Camera) Summary() (string, error) {
	var text C.CameraText
	if err := toError(C.gp_camera_get_summary(cam.c_ref, &text, cam.ctx.c_ref)); err != nil {
		return "", err
	}
	return goString(&text.text[0]), nil
}

func (cam *Camera) Exit() error {
	return toError(C.gp_camera_exit(cam.c_ref, cam.ctx.c_ref))
}

type listCFunc func(*C.Camera, *C.char, *C.CameraList, *C.GPContext) C.int

func (cam *Camera) list(path string, lfn listCFunc) ([]string, error) {
	if path == "" {
		path = "/"
	}

	cl := newCameraList()
	defer cl.Close()

	c_folder, c_folder_free := cString(path)
	defer c_folder_free()

	if err := toError(lfn(cam.c_ref, c_folder, cl.c_ref, cam.ctx.c_ref)); err != nil {
		return nil, err
	}

	return cl.Keys(), nil
}

func (cam *Camera) ListFolders(path string) ([]string, error) {
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	return cam.list(path, func(camera *C.Camera, char *C.char, list *C.CameraList, ctx *C.GPContext) C.int {
		return C.gp_camera_folder_list_folders(camera, char, list, ctx)
	})
}

func (cam *Camera) ListFiles(path string) ([]string, error) {
	return cam.list(path, func(camera *C.Camera, char *C.char, list *C.CameraList, ctx *C.GPContext) C.int {
		return C.gp_camera_folder_list_files(camera, char, list, ctx)
	})
}

func (cam *Camera) File2(folder, file string) CameraFilePath {
	return newCameraFilePath(cam, folder, file)
}

func (cam *Camera) File(p string) CameraFilePath {
	folder, file := path.Split(p)
	return newCameraFilePath(cam, folder, file)
}
