package ext

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path"
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

func Check(r io.ReaderAt) error {
	sb, err := readSuperBlock(r)
	if err != nil {
		return xerrors.Errorf("failed to parse super block: %w", err)
	}

	if err := isCompatible(sb); err != nil {
		return err
	}

	return nil
}

// NewFS is created io/fs.FS for ext4 filesystem
func NewFS(r io.ReaderAt) (*FileSystem, error) {
	sb, err := readSuperBlock(r)
	if err != nil {
		return nil, xerrors.Errorf("failed to parse super block: %w", err)
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

	return fs, nil
}

func isCompatible(sb disklayout.SuperBlock) error {
	if sb.Magic() != common.EXT_SUPER_MAGIC {
		return syserror.EINVAL
	}

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

func (f *FileSystem) Open(name string) (fs.File, error) {
	name = strings.TrimPrefix(name, "/")
	if !fs.ValidPath(name) {
		return nil, fs.ErrInvalid
	}

	dirName, fileName := filepath.Split(name)
	dirEntries, err := f.readDirEntry(dirName)
	if err != nil {
		return nil, err
	}

	for _, entry := range dirEntries {
		if entry.IsDir() || entry.Name() != fileName {
			continue
		}

		inode, ok := entry.(*inode)
		if !ok {
			return nil, xerrors.Errorf("unspecified error, entry is not dir entry %+v", entry)
		}

		if inode.isRefInode() {
			return nil, errors.New("must be file or symlink")
		}

		return &file{
			info: &fileInfo{inode: inode},
		}, nil
	}
	return nil, fs.ErrNotExist
}

func (f *FileSystem) Stat(name string) (fs.FileInfo, error) {
	fi, err := f.Open(name)
	if err != nil {
		info, err := f.ReadDirInfo(name)
		if err != nil {
			return nil, xerrors.Errorf("failed to read dir info: %w", err)
		}
		return info, nil
	}
	info, err := fi.Stat()
	if err != nil {
		return nil, xerrors.Errorf("failed to stat file: %w", err)
	}
	return info, nil
}

func (f *FileSystem) ReadDirInfo(name string) (fs.FileInfo, error) {
	if name == "/" {
		inode, err := newInode(f, disklayout.RootDirInode)
		inode.name = "/"
		if err != nil {
			return nil, xerrors.Errorf("failed to parse root inode: %w", err)
		}
		return &fileInfo{
			inode: inode,
		}, nil
	}
	name = strings.TrimRight(name, "/")
	dirs, dir := path.Split(name)
	dirEntries, err := f.readDirEntry(dirs)
	if err != nil {
		return nil, xerrors.Errorf("failed to read dir entry: %w", err)
	}
	for _, entry := range dirEntries {
		if entry.Name() == strings.Trim(dir, "/") {
			return entry.Info()
		}
	}
	return nil, fs.ErrNotExist
}

func (f *FileSystem) readDirEntry(name string) ([]fs.DirEntry, error) {
	rootEntries, err := f.listEntries(disklayout.RootDirInode)
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

		rootEntries, err = f.listInoEntries(currentIno)
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

func (f *FileSystem) listEntries(ino uint32) ([]*inode, error) {
	in, err := newInode(f, ino)
	if err != nil {
		return nil, xerrors.Errorf("failed to get inode: %w", err)
	}

	return f.listInoEntries(in)
}

func (f *FileSystem) listInoEntries(in *inode) ([]*inode, error) {
	dir, ok := in.impl.(*directory)
	if !ok {
		return nil, xerrors.Errorf("inode is not dir: %d", in.inodeNum)
	}

	inodes := make([]*inode, 0)
	for name, d := range dir.childMap {
		if d.Inode() == 0 || name == "." || name == ".." {
			continue
		}

		entry, err := newInode(f, d.Inode())
		if err != nil {
			return nil, err
		}

		entry.name = name
		inodes = append(inodes, entry)
	}

	return inodes, nil
}
