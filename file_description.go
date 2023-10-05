package ext

// fileDescription is embedded by ext implementations of
// vfs.FileDescriptionImpl.
type fileDescription struct {
	// vfsfd vfs.FileDescription
	// vfs.FileDescriptionDefaultImpl
	// vfs.LockFD
}

func (fd *fileDescription) inode() *inode {
	// return fd.vfsfd.Dentry().Impl().(*dentry).inode
	return nil
}

// Stat implements vfs.FileDescriptionImpl.Stat.
// func (fd *fileDescription) Stat(ctx context.Context, opts vfs.StatOptions) (linux.Statx, error) {
// 	// var stat linux.Statx
// 	// fd.inode().statTo(&stat)
// 	return nil, nil
// }

// SetStat implements vfs.FileDescriptionImpl.SetStat.
// func (fd *fileDescription) SetStat(ctx context.Context, opts vfs.SetStatOptions) error {
// 	if opts.Stat.Mask == 0 {
// 		return nil
// 	}
// 	return linuxerr.EPERM
// }

// SetStat implements vfs.FileDescriptionImpl.StatFS.
// func (fd *fileDescription) StatFS(ctx context.Context) (linux.Statfs, error) {
// 	// var stat linux.Statfs
// 	// fd.filesystem().statTo(&stat)
// 	return nil, nil
// }

// Sync implements vfs.FileDescriptionImpl.Sync.
// func (fd *fileDescription) Sync(ctx context.Context) error {
// 	return nil
// }
