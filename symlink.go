package ext

import (
	"strings"
)

// symlink represents a symlink inode.
//
// +stateify savable
type symlink struct {
	inode  inode
	target string // immutable
}

// newSymlink is the symlink constructor. It reads out the symlink target from
// the inode (however it might have been stored).
func newSymlink(args inodeArgs) (*symlink, error) {
	var link []byte

	// If the symlink target is lesser than 60 bytes, its stores in inode.Data().
	// Otherwise either extents or block maps will be used to store the link.
	size := args.diskInode.Size()
	if size < 60 {
		link = args.diskInode.Data()[:size]
	} else {
		// Create a regular file out of this inode and read out the target.
		regFile, err := newRegularFile(args)
		if err != nil {
			return nil, err
		}

		link = make([]byte, size)
		if n, err := regFile.impl.ReadAt(link, 0); uint64(n) < size {
			return nil, err
		}
	}

	file := &symlink{target: string(link)}
	file.inode.init(args, file)
	return file, nil
}

func (in *inode) isSymlink() bool {
	_, ok := in.impl.(*symlink)
	return ok
}

func (f *symlink) ReadAt(p []byte, off int64) (n int, err error) {
	r := strings.NewReader(f.target)
	return r.ReadAt(p, off)
}
