//+build amd64

package gphoto2

import "C"

func cSize(sz int) C.ulong {
    return C.ulong(sz)
}
