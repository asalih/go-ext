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

// directoryFD represents a directory file description. It implements
// vfs.FileDescriptionImpl.
//
// +stateify savable
type regularFileFD struct {
	fileDescription

	// off is the file offset. off is accessed using atomic memory operations.
	off int64

	// offMu serializes operations that may mutate off.
	// offMu sync.Mutex `state:"nosave"`
}

// // Release implements vfs.FileDescriptionImpl.Release.
// func (fd *regularFileFD) Release(context.Context) {}

// // PRead implements vfs.FileDescriptionImpl.PRead.
// func (fd *regularFileFD) PRead(ctx context.Context, dst usermem.IOSequence, offset int64, opts vfs.ReadOptions) (int64, error) {
// 	safeReader := safemem.FromIOReaderAt{
// 		ReaderAt: fd.inode().impl.(*regularFile).impl,
// 		Offset:   offset,
// 	}

// 	// Copies data from disk directly into usermem without any intermediate
// 	// allocations (if dst is converted into BlockSeq such that it does not need
// 	// safe copying).
// 	return dst.CopyOutFrom(ctx, safeReader)
// }

// // Read implements vfs.FileDescriptionImpl.Read.
// func (fd *regularFileFD) Read(ctx context.Context, dst usermem.IOSequence, opts vfs.ReadOptions) (int64, error) {
// 	n, err := fd.PRead(ctx, dst, fd.off, opts)
// 	fd.offMu.Lock()
// 	fd.off += n
// 	fd.offMu.Unlock()
// 	return n, err
// }

// // ConfigureMMap implements vfs.FileDescriptionImpl.ConfigureMMap.
// func (fd *regularFileFD) ConfigureMMap(ctx context.Context, opts *memmap.MMapOpts) error {
// 	// TODO(b/134676337): Implement mmap(2).
// 	return linuxerr.ENODEV
// }
