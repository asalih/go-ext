package ext

import (
	"github.com/asalih/go-ext/disklayout"
)

type directoryEntries map[string]disklayout.Dirent

// directory represents a directory inode. It holds the childList in memory.
type directory struct {
	inode inode

	// childMap maps the child's filename to the dirent structure stored in
	// childList. This adds some data replication but helps in faster path
	// traversal. For consistency, key == childMap[key].diskDirent.FileName().
	// Immutable.
	childMap directoryEntries
}

// newDirectory is the directory constructor.
func newDirectory(args inodeArgs, newDirent bool) (*directory, error) {
	file := &directory{
		childMap: make(directoryEntries),
	}
	file.inode.init(args, file)

	// Initialize childList by reading dirents from the underlying file.
	if args.diskInode.Flags().Index {
		// TODO(b/134676337): Support hash tree directories. Currently only the '.'
		// and '..' entries are read in.

		// Users cannot navigate this hash tree directory yet.
		// log.Warningf("hash tree directory being used which is unsupported")
		return file, nil
	}

	// The dirents are organized in a linear array in the file data.
	// Extract the file data and decode the dirents.
	regFile, err := newRegularFile(args)
	if err != nil {
		return nil, err
	}

	// buf is used as scratch space for reading in dirents from disk and
	// unmarshalling them into dirent structs.
	buf := make([]byte, disklayout.DirentSize)
	size := args.diskInode.Size()
	for off, inc := uint64(0), uint64(0); off < size; off += inc {
		toRead := size - off
		if toRead > disklayout.DirentSize {
			toRead = disklayout.DirentSize
		}
		if n, err := regFile.impl.ReadAt(buf[:toRead], int64(off)); uint64(n) < toRead {
			return nil, err
		}

		var curDirent disklayout.Dirent
		if newDirent {
			curDirent = &disklayout.DirentNew{}
		} else {
			curDirent = &disklayout.DirentOld{}
		}
		curDirent.UnmarshalBytes(buf)

		if curDirent.Inode() != 0 && len(curDirent.Name()) != 0 {
			// Inode number and name length fields being set to 0 is used to indicate
			// an unused dirent.
			file.childMap[curDirent.Name()] = curDirent
		}

		// The next dirent is placed exactly after this dirent record on disk.
		inc = uint64(curDirent.RecordSize())
	}

	return file, nil
}
