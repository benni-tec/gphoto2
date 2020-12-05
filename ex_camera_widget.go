package gphoto2

import "time"

type WalkWidgetsFunc func(depth int, parent, widget *CameraWidget) error

func (w *CameraWidget) Walk(fn WalkWidgetsFunc, autodispose bool) error {
    if pw, err := w.Parent(); err != nil {
        return err
    } else {
        if autodispose {
            //defer pw.Close()
        }
        if err := w.walk(0, fn, pw, func() (*CameraWidget, error) {
            return w, nil
        }, autodispose); err != nil {
            return err
        } else {
            return nil
        }
    }
}

func (w *CameraWidget) walk(depth int, walkFn WalkWidgetsFunc, pw *CameraWidget, getFn func() (*CameraWidget, error), autodispose bool) error {
    w, err := getFn()
    if err != nil {
        panic(err)
    }

    if err := walkFn(depth, pw, w); err != nil {
        return err
    }

    if count, err := w.ChildrenCount(); err != nil {
        panic(err)
    } else {
        for n := 0; n < count; n++ {
            if err := w.walk(depth+1, walkFn, w, func() (*CameraWidget, error) {
                return w.Child(n)
            }, autodispose); err != nil {
                return err
            }
        }
    }

    return nil
}

func (w *CameraWidget) Value() (interface{}, error) {
    _, wvt, err := w.Type()
    if err != nil {
        return nil, err
    }

    switch wvt {
    case WidgetValueString:
        return w.ValueString()
    case WidgetValueInt:
        return w.ValueInt()
    case WidgetValueFloat:
        return w.ValueFloat()
    case WidgetValueDate:
        return w.ValueDate()
    case WidgetValueWeird:
        return nil, nil
    case WidgetValueInvalid:
        return nil, nil
    default:
        return nil, nil
    }
}

func (w *CameraWidget) SetValue(val interface{}) error {
    panic("not implemented")

    _, wvt, err := w.Type()
    if err != nil {
        return err
    }

    setString := func() error {
        return w.SetValueString("")
    }

    setInt := func() error {
        return w.SetValueInt(0)
    }

    setFloat := func() error {
        return w.SetValueFloat(0)
    }

    setDate := func() error {
        return w.SetValueDate(time.Time{})
    }

    switch wvt {
    case WidgetValueString:
        return setString()
    case WidgetValueInt:
        return setInt()
    case WidgetValueFloat:
        return setFloat()
    case WidgetValueDate:
        return setDate()
    case WidgetValueWeird:
        return nil
    case WidgetValueInvalid:
        return nil
    default:
        return nil
    }
}

