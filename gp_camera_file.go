package gphoto2

// #cgo pkg-config: libgphoto2
// #include <gphoto2.h>
// #include <stdlib.h>
import "C"

type CameraFile struct {
	cam *Camera

	c_ref *C.CameraFile
}

func newCameraFileAllocate(cam *Camera) *CameraFile {
    cf := &CameraFile{
        cam: cam,
    }

    if err := toError(C.gp_file_new(&cf.c_ref)); err != nil {
        panic(err)
    }

    return cf
}

func newCameraFileOpen(cam *Camera, path string) (*CameraFile, error) {
    cf := newCameraFileAllocate(cam)

    c_path, c_path_free := cString(path)
    defer c_path_free()

    if err := toError(C.gp_file_open(cf.c_ref, c_path)); err != nil {
        cf.Close()
        return nil, err
    }

    return cf, nil
}


func newCameraFileGet(cam *Camera, cfp CameraFilePath) (*CameraFile, error) {
    cf := newCameraFileAllocate(cam)

    c_folder, c_file, c_ff_free := cfp.cFolderFile()
    defer c_ff_free()

    if err := toError(C.gp_camera_file_get(cam.c_ref, c_folder, c_file, C.GP_FILE_TYPE_NORMAL, cf.c_ref, cam.ctx.c_ref)); err != nil {
        cf.Close()
        return nil, err
    }

    return cf, nil
}

func (cf *CameraFile) Close() {
	if err := toError(C.gp_file_free(cf.c_ref)); err != nil {
	    panic(err)
    } else {
        cf.c_ref = nil
    }
}

func (cf *CameraFile) GetPath() (string, error) {
    var c_name *C.char

    if err := toError(C.gp_file_get_name(cf.c_ref, &c_name)); err != nil {
        return "", err
    }

    return goString(c_name), nil
}

func (cf *CameraFile) GetDataAndSize() ([]byte, error) {
    var c_size C.ulong
    var c_buf *C.char

    if err := toError(C.gp_file_get_data_and_size(cf.c_ref, &c_buf, &c_size)); err != nil {
        return nil, err
    }

    return goBytes(c_buf, int(c_size)), nil
}
