package ext

import (
	"io"
)

// regularFile represents a regular file's inode. This too follows the
// inheritance pattern prevelant in the vfs layer described in
// pkg/sentry/vfs/README.md.
//
// +stateify savable
type regularFile struct {
	inode inode

	// This is immutable. The first field of fileReader implementations must be
	// regularFile to ensure temporality.
	// io.ReaderAt is more strict than io.Reader in the sense that a partial read
	// is always accompanied by an error. If a read spans past the end of file, a
	// partial read (within file range) is done and io.EOF is returned.
	impl io.ReaderAt
}

// newRegularFile is the regularFile constructor. It figures out what kind of
// file this is and initializes the fileReader.
func newRegularFile(args inodeArgs) (*regularFile, error) {
	if args.diskInode.Flags().Extents {
		file, err := newExtentFile(args)
		if err != nil {
			return nil, err
		}
		return &file.regFile, nil
	}

	file, err := newBlockMapFile(args)
	if err != nil {
		return nil, err
	}
	return &file.regFile, nil
}

func (in *inode) isRegular() bool {
	_, ok := in.impl.(*regularFile)
	return ok
}
