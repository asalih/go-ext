package ext

import (
	"io/fs"

	"github.com/asalih/go-ext/disklayout"
	"github.com/asalih/go-ext/linux"
	"github.com/asalih/go-ext/syserror"
)

// inode represents an ext inode.
//
// inode uses the same inheritance pattern that pkg/sentry/vfs structures use.
// This has been done to increase memory locality.
//
// Implementations:
//
//	inode --
//	       |-- dir
//	       |-- symlink
//	       |-- regular--
//	                   |-- extent file
//	                   |-- block map file
//
// +stateify savable
type inode struct {
	// name of entry
	name string

	// refs is a reference count. refs is accessed using atomic memory operations.
	refs int64

	// fsR is the containing filesystem.
	fsR *FileSystem

	// inodeNum is the inode number of this inode on disk. This is used to
	// identify inodes within the ext filesystem.
	inodeNum uint32

	// blkSize is the fs data block size. Same as filesystem.sb.BlockSize().
	blkSize uint64

	// diskInode gives us access to the inode struct on disk. Immutable.
	diskInode disklayout.Inode

	// This is immutable. The first field of the implementations must have inode
	// as the first field to ensure temporality.
	impl interface{}
}

var _ fs.DirEntry = (*inode)(nil)

type inodeArgs struct {
	fs        *FileSystem
	inodeNum  uint32
	blkSize   uint64
	diskInode disklayout.Inode
}

// newInode is the inode constructor. Reads the inode off disk. Identifies
// inodes based on the absolute inode number on disk.
func newInode(fsR *FileSystem, inodeNum uint32) (*inode, error) {
	if inodeNum == 0 {
		panic("inode number 0 on ext filesystems is not possible")
	}

	inodeRecordSize := fsR.sb.InodeSize()
	var diskInode disklayout.Inode
	if inodeRecordSize == disklayout.OldInodeSize {
		diskInode = &disklayout.InodeOld{}
	} else {
		diskInode = &disklayout.InodeNew{}
	}

	// Calculate where the inode is actually placed.
	inodesPerGrp := fsR.sb.InodesPerGroup()
	blkSize := fsR.sb.BlockSize()
	inodeTableOff := fsR.bgs[getBGNum(inodeNum, inodesPerGrp)].InodeTable() * blkSize
	inodeOff := inodeTableOff + uint64(uint32(inodeRecordSize)*getBGOff(inodeNum, inodesPerGrp))

	if err := readFromDisk(fsR.dev, int64(inodeOff), diskInode); err != nil {
		return nil, err
	}

	// Build the inode based on its type.
	args := inodeArgs{
		fs:        fsR,
		inodeNum:  inodeNum,
		blkSize:   blkSize,
		diskInode: diskInode,
	}

	switch diskInode.Mode().FileType() {
	case linux.ModeSymlink:
		f, err := newSymlink(args)
		if err != nil {
			return nil, err
		}
		return &f.inode, nil
	case linux.ModeRegular:
		f, err := newRegularFile(args)
		if err != nil {
			return nil, err
		}
		return &f.inode, nil
	case linux.ModeDirectory:
		f, err := newDirectory(args, fsR.sb.IncompatibleFeatures().DirentFileType)
		if err != nil {
			return nil, err
		}
		return &f.inode, nil
	default:
		// TODO(b/134676337): Return appropriate errors for sockets, pipes and devices.
		return nil, syserror.EINVAL
	}
}

func (in *inode) init(args inodeArgs, impl interface{}) {
	in.fsR = args.fs
	in.inodeNum = args.inodeNum
	in.blkSize = args.blkSize
	in.diskInode = args.diskInode
	in.impl = impl
}

func (in *inode) isDir() bool {
	_, ok := in.impl.(*directory)
	return ok
}

// statTo writes the statx fields to the output parameter.
func (in *inode) statTo(stat *Statx) {
	stat.Mask = linux.STATX_TYPE | linux.STATX_MODE | linux.STATX_NLINK |
		linux.STATX_UID | linux.STATX_GID | linux.STATX_INO | linux.STATX_SIZE |
		linux.STATX_ATIME | linux.STATX_CTIME | linux.STATX_MTIME
	stat.Blksize = uint32(in.blkSize)
	stat.Mode = uint16(in.diskInode.Mode())
	stat.Nlink = uint32(in.diskInode.LinksCount())
	stat.UID = uint32(in.diskInode.UID())
	stat.GID = uint32(in.diskInode.GID())
	stat.Ino = uint64(in.inodeNum)
	stat.Size = in.diskInode.Size()
	stat.Atime = in.diskInode.AccessTime()
	stat.Ctime = in.diskInode.ChangeTime()
	stat.Mtime = in.diskInode.ModificationTime()
	// stat.DevMajor = linux.UNNAMED_MAJOR
	// stat.DevMinor = in.fsR.devMinor
	// TODO(b/134676337): Set stat.Blocks which is the number of 512 byte blocks
	// (including metadata blocks) required to represent this file.
}

func (i *inode) Name() string {
	return i.name
}

func (i *inode) IsDir() bool {
	return i.isDir()
}

func (i *inode) Type() fs.FileMode {
	return i.diskInode.Mode().FSMode()
}

func (i *inode) Info() (fs.FileInfo, error) {
	return &fileInfo{inode: i}, nil
}

// getBGNum returns the block group number that a given inode belongs to.
func getBGNum(inodeNum uint32, inodesPerGrp uint32) uint32 {
	return (inodeNum - 1) / inodesPerGrp
}

// getBGOff returns the offset at which the given inode lives in the block
// group's inode table, i.e. the index of the inode in the inode table.
func getBGOff(inodeNum uint32, inodesPerGrp uint32) uint32 {
	return (inodeNum - 1) % inodesPerGrp
}
