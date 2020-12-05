package gphoto2

// #cgo pkg-config: libgphoto2
// #include <gphoto2.h>
// #include <stdlib.h>
import "C"

import (
	"fmt"
	"time"
	"unsafe"
)

type CameraWidget struct {
	c_ref *C.CameraWidget

	cam *Camera
}

func (cam *Camera) RootWidget() (*CameraWidget, func() error, error) {
	cw := &CameraWidget{
		cam: cam,
	}

	c_name, c_name_free := cString("")
	defer c_name_free()

	if err := toError(C.gp_widget_new(C.GP_WIDGET_WINDOW, c_name, &cw.c_ref)); err != nil {
		return nil, nil, err
	} else if err := toError(C.gp_camera_get_config(cam.c_ref, &cw.c_ref, cam.ctx.c_ref)); err != nil {
		cw.Close()
		return nil, nil, err
	}

	return cw, func() error {
		return toError(C.gp_camera_set_config(cam.c_ref, cw.c_ref, cam.ctx.c_ref))
	}, nil
}

func (cw *CameraWidget) GetID() (int, error) {
	var c_id C.int

	if err := toError(C.gp_widget_get_id(cw.c_ref, &c_id)); err != nil {
		return 0, err
	} else {
		return int(c_id), nil
	}
}

func (cw *CameraWidget) Parent() (*CameraWidget, error) {
	cwParent := &CameraWidget{
		cam: cw.cam,
	}

	if err := toError(C.gp_widget_get_parent(cw.c_ref, &cwParent.c_ref)); err != nil {
		return nil, err
	}

	return cw, nil
}

func (cw *CameraWidget) Close() {
	if err := toError(C.gp_widget_free(cw.c_ref)); err != nil {
		panic(err)
	}
}

func (cw *CameraWidget) Readonly() (bool, error) {
	var c_ro C.int

	if err := toError(C.gp_widget_get_readonly(cw.c_ref, &c_ro)); err != nil {
		return false, err
	}

	return c_ro == 1, nil
}

func (cw *CameraWidget) Name() (string, error) {
	var c_name *C.char

	if err := toError(C.gp_widget_get_name(cw.c_ref, &c_name)); err != nil {
		return "", err
	} else {
		//defer C.free(unsafe.Pointer(c_name))
	}

	return goString(c_name), nil
}

func (cw *CameraWidget) Label() (string, error) {
	var c_name *C.char

	if err := toError(C.gp_widget_get_label(cw.c_ref, &c_name)); err != nil {
		return "", err
	} else {
		//defer C.free(unsafe.Pointer(c_name))
	}

	return goString(c_name), nil
}

/***************************************************************
** Widget Children
**/

func (cw *CameraWidget) ChildByName(name string) (*CameraWidget, error) {
	cwChild := &CameraWidget{
		cam: cw.cam,
	}

	c_name, c_name_free := cString(name)
	defer c_name_free()

	if err := toError(C.gp_widget_get_child_by_name(cw.c_ref, c_name, &cwChild.c_ref)); err != nil {
		return nil, err
	}

	return cwChild, nil
}

func (cw *CameraWidget) ChildByLabel(label string) (*CameraWidget, error) {
	cwChild := &CameraWidget{
		cam: cw.cam,
	}

	c_label, c_label_free := cString(label)
	defer c_label_free()

	if err := toError(C.gp_widget_get_child_by_label(cw.c_ref, c_label, &cwChild.c_ref)); err != nil {
		return nil, err
	}

	return cwChild, nil
}

func (cw *CameraWidget) ChildrenCount() (int, error) {
	c_count := C.gp_widget_count_children(cw.c_ref)

	if err := toError(c_count); err != nil {
		return 0, err
	} else {
		return int(c_count), nil
	}
}

func (cw *CameraWidget) Child(n int) (*CameraWidget, error) {
	cwChild := &CameraWidget{
		cam: cw.cam,
	}

	if err := toError(C.gp_widget_get_child(cw.c_ref, (C.int)(n), &cwChild.c_ref)); err != nil {
		return nil, err
	}

	return cwChild, nil
}

func (cw *CameraWidget) ChildByID(id int) (*CameraWidget, error) {
	cwChild := &CameraWidget{
		cam: cw.cam,
	}

	if err := toError(C.gp_widget_get_child_by_id(cw.c_ref, (C.int)(id), &cwChild.c_ref)); err != nil {
		return nil, err
	}

	return cwChild, nil
}

/***************************************************************
** Widget Value
**/

func (cw *CameraWidget) value(val unsafe.Pointer) error {
	if err := toError(C.gp_widget_get_value(cw.c_ref, val)); err != nil {
		return err
	}
	return nil
}

func (cw *CameraWidget) ValueString() (string, error) {
	var c_val *C.char

	if err := cw.value(unsafe.Pointer(&c_val)); err != nil {
		return "", err
	}

	if c_val != nil {
	    // FIXME
		// defer C.free(unsafe.Pointer(c_val))
	}

	return goString(c_val), nil
}

func (cw *CameraWidget) ValueInt() (int, error) {
	var c_val C.int

	if err := cw.value(unsafe.Pointer(&c_val)); err != nil {
		return 0, err
	}

	return int(c_val), nil
}

func (cw *CameraWidget) ValueFloat() (float32, error) {
	var c_val C.float

	if err := cw.value(unsafe.Pointer(&c_val)); err != nil {
		return 0, err
	}

	return float32(c_val), nil
}

func (cw *CameraWidget) ValueDate() (time.Time, error) {
	var c_val C.int

	if err := cw.value(unsafe.Pointer(&c_val)); err != nil {
		return time.Time{}, err
	}

	return time.Unix(int64(c_val), 0), nil
}

func (cw *CameraWidget) ChoicesCount() (int, error) {
	c_count := C.gp_widget_count_choices(cw.c_ref)

	if err := toError(c_count); err != nil {
		return 0, err
	} else {
		return int(c_count), nil
	}
}

func (cw *CameraWidget) Choices() ([]string, error) {
	choicesCount, err := cw.ChoicesCount()
	if err != nil {
		return nil, err
	}

	choices := make([]string, choicesCount)
	for n := 0; n < choicesCount; n++ {
		var c_val *C.char
		if err := toError(C.gp_widget_get_choice(cw.c_ref, C.int(n), &c_val)); err != nil {
			return nil, err
		}
		choices[n] = goString(c_val)
	}

	return choices, nil
}

func (cw *CameraWidget) setValue(val unsafe.Pointer) error {
	if err := toError(C.gp_widget_set_value(cw.c_ref, val)); err != nil {
		return err
	}
	return nil
}

func (cw *CameraWidget) SetValueString(val string) error {
	c_val, c_val_free := cString(val)
	defer c_val_free()

	return cw.setValue(unsafe.Pointer(c_val))
}

func (cw *CameraWidget) SetValueInt(val int) error {
	return cw.setValue(unsafe.Pointer(uintptr(C.int(val))))
}

func (cw *CameraWidget) SetValueFloat(val float32) error {
	return cw.setValue(unsafe.Pointer(uintptr(C.float(val))))
}

func (cw *CameraWidget) SetValueDate(val time.Time) error {
	return cw.setValue(unsafe.Pointer(uintptr(C.int(val.Unix()))))
}

/***************************************************************
** Widget Type
**/

type WidgetValueType string

const (
	WidgetValueInvalid WidgetValueType = ""

	WidgetValueString WidgetValueType = "string"
	WidgetValueInt    WidgetValueType = "int"
	WidgetValueFloat  WidgetValueType = "float"
	WidgetValueDate   WidgetValueType = "date"
	WidgetValueWeird  WidgetValueType = "weird"
)

type WidgetType string

const (
	WidgetInvalid WidgetType = ""

	// Window widget This is the toplevel configuration widget. It should likely contain multiple widget seciton entries
	WidgetWindow WidgetType = "window"

	// Section widget (think Tab)
	WidgetSection WidgetType = "section"

	// Text widget
	WidgetText WidgetType = "text"

	// Slider widget
	WidgetRange WidgetType = "range"

	// Toggle widget (think check box)
	WidgetToggle WidgetType = "toggle"

	// Radio button widget
	WidgetRadio WidgetType = "radio"

	// Menu widget (same as RADIO)
	WidgetMenu WidgetType = "menu"

	// Button press widget
	WidgetButton WidgetType = "button"

	// Date entering widget
	WidgetDate WidgetType = "date"
)

func (cw *CameraWidget) Type() (WidgetType, WidgetValueType, error) {
	var c_type C.CameraWidgetType

	if err := toError(C.gp_widget_get_type(cw.c_ref, &c_type)); err != nil {
		return WidgetInvalid, WidgetValueInvalid, err
	}

	switch c_type {
	case C.GP_WIDGET_WINDOW:
		return WidgetWindow, WidgetValueWeird, nil
	case C.GP_WIDGET_SECTION:
		return WidgetSection, WidgetValueWeird, nil
	case C.GP_WIDGET_TEXT:
		return WidgetText, WidgetValueString, nil
	case C.GP_WIDGET_RANGE:
		return WidgetRange, WidgetValueFloat, nil
	case C.GP_WIDGET_TOGGLE:
		return WidgetToggle, WidgetValueInt, nil
	case C.GP_WIDGET_RADIO:
		return WidgetRadio, WidgetValueString, nil
	case C.GP_WIDGET_MENU:
		return WidgetMenu, WidgetValueInt, nil
	case C.GP_WIDGET_BUTTON:
		return WidgetButton, WidgetValueInt, nil
	case C.GP_WIDGET_DATE:
		return WidgetDate, WidgetValueDate, nil
	default:
		return WidgetInvalid, WidgetValueInvalid, fmt.Errorf("unhandled widget type value: %d", int(c_type))
	}
}
