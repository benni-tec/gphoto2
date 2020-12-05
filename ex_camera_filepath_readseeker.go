package gphoto2

import (
	"errors"
	"io"
)

var _ io.ReadSeeker = (*fileReadSeeker)(nil)

type fileReadSeeker struct {
	f      CameraFilePath
	offset int
}

func (cf CameraFilePath) ReadSeeker() (io.ReadSeeker, error) {
	rc := &fileReadSeeker{f: cf, offset: 0}

	return rc, nil
}

func (rs *fileReadSeeker) Read(p []byte) (n int, err error) {
	n, err = rs.f.ReadOffset(rs.offset, p)
	rs.offset += n
	return n, err
}

// FIXME
func (rs *fileReadSeeker) Read2(p []byte) (n int, err error) {
	data, err := rs.f.ReadOffsetSize(rs.offset, len(p))
	rs.offset += len(data)
	copy(p, data)
	return n, err
}

func (rs *fileReadSeeker) Seek(offset int64, whence int) (int64, error) {
	n := int64(rs.offset)
	switch whence {
	case io.SeekStart:
		n = offset
	case io.SeekCurrent:
		n += offset
	case io.SeekEnd:
		if fi, err := rs.f.GetInfo(); err != nil {
			return int64(rs.offset), err
		} else {
			size := fi.File.Size
			n = size + offset
		}
	default:
		return int64(rs.offset), errors.New("gphoto.ReadSeeker.Seek: invalid whence")
	}

	if n < 0 {
		return int64(rs.offset), errors.New("gphoto.ReadSeeker.Seek: negative position")
	}

	rs.offset = int(n)
	return int64(rs.offset), nil

}
