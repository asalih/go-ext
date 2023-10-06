package ext

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/asalih/go-ext/common"
	"github.com/asalih/go-ext/disklayout"
	"github.com/asalih/go-ext/syserror"
	"golang.org/x/xerrors"
)

// FileSystem is implemented io/fs interface
type FileSystem struct {
	dev io.ReaderAt

	sb  disklayout.SuperBlock
	bgs []disklayout.BlockGroup
}

// NewFS is created io/fs.FS for ext4 filesystem
func NewFS(r io.ReaderAt) (*FileSystem, error) {
	sb, err := readSuperBlock(r)
	if err != nil {
		return nil, xerrors.Errorf("failed to parse super block: %w", err)
	}

	if sb.Magic() != common.EXT_SUPER_MAGIC {
		return nil, syserror.EINVAL
	}

	if err := isCompatible(sb); err != nil {
		return nil, err
	}

	bgs, err := readBlockGroups(r, sb)
	if err != nil {
		return nil, err
	}

	fs := &FileSystem{
		dev: r,
		sb:  sb,
		bgs: bgs,
	}

	// rin, err := newInode(fs, disklayout.RootDirInode)
	// if err != nil {
	// 	return nil, err
	// }

	// fs.rootInode = rin

	return fs, nil
}

func isCompatible(sb disklayout.SuperBlock) error {
	// Please note that what is being checked is limited based on the fact that we
	// are mounting readonly and that we are not journaling. When mounting
	// read/write or with a journal, this must be reevaluated.
	incompatFeatures := sb.IncompatibleFeatures()
	if incompatFeatures.MetaBG {
		return errors.New("ext fs: meta block groups are not supported")
	}
	if incompatFeatures.MMP {
		return errors.New("ext fs: multiple mount protection is not supported")
	}
	if incompatFeatures.Encrypted {
		return errors.New("ext fs: encrypted inodes not supported")
	}
	if incompatFeatures.InlineData {
		return errors.New("ext fs: inline files not supported")
	}
	return nil
}

func (f *FileSystem) ReadDir(path string) ([]fs.DirEntry, error) {
	dirEntries, err := f.readDirEntry(path)
	if err != nil {
		return nil, err
	}

	return dirEntries, nil
}

func (ext4 *FileSystem) readDirEntry(name string) ([]fs.DirEntry, error) {
	rootEntries, err := ext4.listEntries(disklayout.RootDirInode)
	if err != nil {
		return nil, xerrors.Errorf("failed to list file infos: %w", err)
	}

	var currentIno *inode
	dirs := strings.Split(strings.Trim(filepath.Clean(name), string(os.PathSeparator)), string(os.PathSeparator))
	if len(dirs) == 1 && dirs[0] == "." || dirs[0] == "" {
		var dirEntries []fs.DirEntry
		for _, sino := range rootEntries {
			if sino.Name() == "." || sino.Name() == ".." {
				continue
			}
			dirEntries = append(dirEntries, sino)
		}
		return dirEntries, nil
	}

	for i, dir := range dirs {
		found := false
		for _, fileInfo := range rootEntries {
			if fileInfo.Name() != dir {
				continue
			}
			if !fileInfo.IsDir() {
				return nil, xerrors.Errorf("%s is file, directory: %w", fileInfo.Name(), fs.ErrNotExist)
			}
			found = true
			currentIno = fileInfo
			break
		}

		if !found {
			return nil, fs.ErrNotExist
		}

		rootEntries, err = ext4.listInoEntries(currentIno)
		if err != nil {
			return nil, xerrors.Errorf("failed to list directory entries inode(%d): %w", currentIno, err)
		}
		if i != len(dirs)-1 {
			continue
		}

		var dirEntries []fs.DirEntry
		for _, fileInfo := range rootEntries {
			// Skip current directory and parent directory
			// infinit loop in walkDir
			if fileInfo.Name() == "." || fileInfo.Name() == ".." {
				continue
			}

			dirEntries = append(dirEntries, fileInfo)
		}
		return dirEntries, nil
	}
	return nil, fs.ErrNotExist
}

func (ext4 *FileSystem) listEntries(ino uint32) ([]*inode, error) {
	in, err := newInode(ext4, ino)
	if err != nil {
		return nil, xerrors.Errorf("failed to get inode: %w", err)
	}

	return ext4.listInoEntries(in)
}

func (ext4 *FileSystem) listInoEntries(in *inode) ([]*inode, error) {
	dir, ok := in.impl.(*directory)
	if !ok {
		return nil, xerrors.Errorf("inode is not dir: %d", in.inodeNum)
	}

	inodes := make([]*inode, 0)
	for name, d := range dir.childMap {
		if d.Inode() == 0 || name == "." || name == ".." {
			continue
		}

		entry, err := newInode(ext4, d.Inode())
		if err != nil {
			return nil, err
		}

		entry.name = name
		inodes = append(inodes, entry)
	}

	return inodes, nil
}
