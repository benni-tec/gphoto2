package gphoto2

// #cgo pkg-config: libgphoto2
// #include <gphoto2.h>
// #include <stdlib.h>
import "C"

type CameraAbilities struct {
	ID      string
	Model   string
	Library string
}

func (cam *Camera) GetAbilities() (*CameraAbilities, error) {
	var c_ca C.CameraAbilities
	if err := toError(C.gp_camera_get_abilities(cam.c_ref, &c_ca)); err != nil {
		return nil, err
	}

	return &CameraAbilities{
		ID:      goString(&c_ca.id[0]),
		Model:   goString(&c_ca.model[0]),
		Library: goString(&c_ca.library[0]),
	}, nil
}
