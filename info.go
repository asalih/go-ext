package ext

import (
	"io/fs"
	"time"
)

type fileInfo struct {
	*inode
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
	var st Statx
	f.statTo(&st)
	return st
}
