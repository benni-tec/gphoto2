package gphoto2

// #cgo pkg-config: libgphoto2
// #include <gphoto2.h>
import "C"

import "fmt"

type Error int

const (
	Err                   Error = C.GP_ERROR
	ErrBadParameters      Error = C.GP_ERROR_BAD_PARAMETERS
	ErrNoMemory           Error = C.GP_ERROR_NO_MEMORY
	ErrLibrary            Error = C.GP_ERROR_LIBRARY
	ErrUnknownPort        Error = C.GP_ERROR_UNKNOWN_PORT
	ErrNotSupported       Error = C.GP_ERROR_NOT_SUPPORTED
	ErrIO                 Error = C.GP_ERROR_IO
	ErrFixedLimitExceeded Error = C.GP_ERROR_FIXED_LIMIT_EXCEEDED
	ErrTimeout            Error = C.GP_ERROR_TIMEOUT
	ErrIOSupportedSerial  Error = C.GP_ERROR_IO_SUPPORTED_SERIAL
	ErrIOSupportedUSB     Error = C.GP_ERROR_IO_SUPPORTED_USB
	ErrIOInit             Error = C.GP_ERROR_IO_INIT
	ErrIORead             Error = C.GP_ERROR_IO_READ
	ErrIOWrite            Error = C.GP_ERROR_IO_WRITE
	ErrIOUpdate           Error = C.GP_ERROR_IO_UPDATE
	ErrIOSerialSpeed      Error = C.GP_ERROR_IO_SERIAL_SPEED
	ErrIOUSBClearHalt     Error = C.GP_ERROR_IO_USB_CLEAR_HALT
	ErrIOUSBFind          Error = C.GP_ERROR_IO_USB_FIND
	ErrIOUSBClaim         Error = C.GP_ERROR_IO_USB_CLAIM
	ErrIOLock             Error = C.GP_ERROR_IO_LOCK
	ErrHal                Error = C.GP_ERROR_HAL
	ErrCorruptedData      Error = C.GP_ERROR_CORRUPTED_DATA
	ErrFileExists         Error = C.GP_ERROR_FILE_EXISTS
	ErrModelNotFound      Error = C.GP_ERROR_MODEL_NOT_FOUND
	ErrDirectoryNotFound  Error = C.GP_ERROR_DIRECTORY_NOT_FOUND
	ErrFileNotFound       Error = C.GP_ERROR_FILE_NOT_FOUND
	ErrDirectoryExists    Error = C.GP_ERROR_DIRECTORY_EXISTS
	ErrCameraBusy         Error = C.GP_ERROR_CAMERA_BUSY
	ErrPathNotAbsolute    Error = C.GP_ERROR_PATH_NOT_ABSOLUTE
	ErrCancel             Error = C.GP_ERROR_CANCEL
	ErrCameraError        Error = C.GP_ERROR_CAMERA_ERROR
	ErrOsFailure          Error = C.GP_ERROR_OS_FAILURE
	ErrNoSpace            Error = C.GP_ERROR_NO_SPACE
)

func (err Error) Error() string {
	message := C.GoString(C.gp_result_as_string(C.int(err)))
	if message == "" {
		message = "libgphoto2: unknown error"
	}
	return fmt.Sprintf("libgphoto2: [%d] %s", err, message)
}

func toError(code C.int) error {
	if code >= 0 {
		return nil
	}

	return Error(code)
}
