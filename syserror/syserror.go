package syserror

import (
	"errors"
	"syscall"
)

// The following variables have the same meaning as their syscall equivalent.
var (
	EIDRM    = error(syscall.Errno(0x2b))
	EINTR    = error(syscall.Errno(0x4))
	EIO      = error(syscall.Errno(0x5))
	EISDIR   = error(syscall.Errno(0x15))
	ENOENT   = error(syscall.Errno(0x2))
	ENOEXEC  = error(syscall.Errno(0x8))
	ENOMEM   = error(syscall.Errno(0xc))
	ENOTSOCK = error(syscall.Errno(0x58))
	ENOSPC   = error(syscall.Errno(0x1c))
	ENOSYS   = error(syscall.Errno(0x26))
)

var (
	// ErrWouldBlock is an internal error used to indicate that an operation
	// cannot be satisfied immediately, and should be retried at a later
	// time, possibly when the caller has received a notification that the
	// operation may be able to complete. It is used by implementations of
	// the kio.File interface.
	ErrWouldBlock = errors.New("request would block")

	// ErrInterrupted is returned if a request is interrupted before it can
	// complete.
	ErrInterrupted = errors.New("request was interrupted")

	// ErrExceedsFileSizeLimit is returned if a request would exceed the
	// file's size limit.
	ErrExceedsFileSizeLimit = errors.New("exceeds file size limit")

	// EINVAL is returned for invalid argument
	EINVAL = errors.New("invalid argument")
)
