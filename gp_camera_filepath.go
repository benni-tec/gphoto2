package gphoto2

// #cgo pkg-config: libgphoto2
// #include <gphoto2.h>
// #include <stdlib.h>
import "C"

import (
	"io"
	"path"
	"unsafe"
)

type CameraFilePath struct {
	cam *Camera

	folder string
	file   string
}

func newCameraFilePath(cam *Camera, folder, file string) CameraFilePath {
	cf := CameraFilePath{
		cam:    cam,
		folder: folder,
		file:   file,
	}

	return cf
}

func (cf CameraFilePath) cFolderFile() (c_folder, c_file *C.char, free func()) {
	c_folder, c_folder_free := cString(cf.folder)
	c_file, c_file_free := cString(cf.file)

	return c_folder, c_file, func() {
		c_folder_free()
		c_file_free()
	}
}

func (cf CameraFilePath) Folder() string {
	return cf.folder
}

func (cf CameraFilePath) Name() string {
	return cf.file
}

func (cf CameraFilePath) Path() string {
	return path.Join(cf.Folder(), cf.Name())
}

func (cf CameraFilePath) Delete() error {
	c_folder, c_file, c_ff_free := cf.cFolderFile()
	defer c_ff_free()

	return toError(C.gp_camera_file_delete(cf.cam.c_ref, c_folder, c_file, cf.cam.ctx.c_ref))
}

type CameraFileInfo struct {
	File CameraFileInfoFile
}

type CameraFileInfoFile struct {
	Size          int64
	MTime         int64
	Width, Height int
}

func (cf CameraFilePath) GetInfo() (*CameraFileInfo, error) {
	c_info := new(C.CameraFileInfo)

	c_folder, c_file, c_ff_free := cf.cFolderFile()
	defer c_ff_free()

	if err := toError(C.gp_camera_file_get_info(
		cf.cam.c_ref,
		c_folder,
		c_file,
		c_info,
		cf.cam.ctx.c_ref,
	)); err != nil {
		return nil, err
	}

	return &CameraFileInfo{
		File: CameraFileInfoFile{
			Size:   int64(c_info.file.size),
			MTime:  int64(c_info.file.mtime),
			Width:  int(c_info.file.width),
			Height: int(c_info.file.height),
		},
	}, nil
}

func (cf CameraFilePath) ReadOffsetSize(offset, size int) ([]byte, error) {
	var (
		c_size   = cSize(size)
		c_offset = cSize(offset)
	)

	c_buf, c_buf_free := cMalloc(size)
	defer c_buf_free()

	c_folder, c_file, c_ff_free := cf.cFolderFile()
	defer c_ff_free()

	if err := toError(C.gp_camera_file_read(
		cf.cam.c_ref,
		c_folder,
		c_file,
		C.GP_FILE_TYPE_NORMAL,
		c_offset,
		c_buf,
		&c_size,
		cf.cam.ctx.c_ref,
	)); err != nil {
		return nil, err
	}

	var (
		buf = goBytes(c_buf, int(c_size))
		err error
	)
	if len(buf) < size {
		err = io.EOF
	}
	return buf, err
}

func (cf CameraFilePath) ReadOffset(offset int, p []byte) (int, error) {
	var (
		c_size   = cSize(len(p))
		c_offset = cSize(offset)
	)

	c_buf := (*C.char)(unsafe.Pointer(&p[0]))

	c_folder, c_file, c_ff_free := cf.cFolderFile()
	defer c_ff_free()

	if err := toError(C.gp_camera_file_read(
		cf.cam.c_ref,
		c_folder,
		c_file,
		C.GP_FILE_TYPE_NORMAL,
		c_offset,
		c_buf,
		&c_size,
		cf.cam.ctx.c_ref,
	)); err != nil {
		return 0, err
	}

	if int(c_size) < len(p) {
		return int(c_size), io.EOF
	} else {
		return int(c_size), nil
	}
}
