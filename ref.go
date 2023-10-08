package ext

type refInode struct {
	inode inode
}

// newRefInode is the inode ref
func newRefInode(args inodeArgs) (*refInode, error) {
	file := &refInode{}
	file.inode.init(args, file)
	return file, nil
}

func (in *inode) isRefInode() bool {
	_, ok := in.impl.(*refInode)
	return ok
}
