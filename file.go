package ext

import (
	"errors"
	"io"
	"io/fs"
	"time"
)

type file struct {
	info *fileInfo

	position int64
}

type fileInfo struct {
	*inode
}

type Statx struct {
	Mask      uint32
	Blksize   uint32
	Nlink     uint32
	UID       uint32
	GID       uint32
	Mode      uint16
	Ino       uint64
	Size      uint64
	Blocks    uint64
	Atime     time.Time
	Ctime     time.Time
	Mtime     time.Time
	RdevMajor uint32
	RdevMinor uint32
	DevMajor  uint32
	DevMinor  uint32
}

func (f *fileInfo) Name() string {
	return f.name
}

func (f *fileInfo) Size() int64 {
	return int64(f.diskInode.Size())
}

func (f *fileInfo) Mode() fs.FileMode {
	return f.Type()
}

func (f *fileInfo) ModTime() time.Time {
	return f.diskInode.ChangeTime()
}

func (f *fileInfo) IsDir() bool {
	return f.isDir()
}

func (f *fileInfo) Sys() any {
	st := new(Statx)
	f.statTo(st)
	return st
}

func (f *file) Stat() (fs.FileInfo, error) {
	return f.info, nil
}

func (f *file) Read(b []byte) (int, error) {
	var rdr io.ReaderAt
	if f.info.inode.isSymlink() {
		sl, ok := f.info.inode.impl.(*symlink)
		if !ok {
			return 0, fs.ErrInvalid
		}
		rdr = sl

	} else if f.info.inode.isRegular() {
		rf, ok := f.info.inode.impl.(*regularFile)
		if !ok {
			return 0, fs.ErrInvalid
		}
		rdr = rf.impl
	} else {
		return 0, fs.ErrInvalid
	}

	sz := f.info.Size()
	toRead := len(b)
	if f.position+int64(toRead) > sz {
		toRead = int(sz - f.position)
	}

	n, err := rdr.ReadAt(b[:toRead], f.position)
	if err != nil {
		return n, err
	}
	if toRead != len(b) {
		err = io.EOF
	}

	f.position += int64(n)
	return n, err
}

func (f *file) ReadAt(p []byte, off int64) (n int, err error) {
	var rdr io.ReaderAt
	if f.info.inode.isSymlink() {
		sl, ok := f.info.inode.impl.(*symlink)
		if !ok {
			return 0, fs.ErrInvalid
		}
		rdr = sl

	} else if f.info.inode.isRegular() {
		rf, ok := f.info.inode.impl.(*regularFile)
		if !ok {
			return 0, fs.ErrInvalid
		}
		rdr = rf.impl
	} else {
		return 0, fs.ErrInvalid
	}

	n, err = rdr.ReadAt(p, off)
	if err != nil {
		return n, err
	}

	f.position = off
	return n, err
}

func (f *file) Close() error {
	return nil
}

// Seek implements vfs.FileDescriptionImpl.Seek.
func (f *file) Seek(offset int64, whence int) (ret int64, err error) {
	var newPos int64

	switch whence {
	case io.SeekStart:
		newPos = offset
	case io.SeekCurrent:
		newPos = f.position + offset
	case io.SeekEnd:
		newPos = f.info.Size() + offset
	default:
		return 0, errors.New("invalid whence value")
	}

	if newPos < 0 {
		return 0, errors.New("negative position")
	}

	f.position = newPos
	return newPos, nil
}
