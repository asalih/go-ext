package common

// Filesystem types used in statfs(2).
//
// See linux/magic.h.
const (
	ANON_INODE_FS_MAGIC   = 0x09041934
	CGROUP_SUPER_MAGIC    = 0x27e0eb
	DEVPTS_SUPER_MAGIC    = 0x00001cd1
	EXT_SUPER_MAGIC       = 0xef53
	FUSE_SUPER_MAGIC      = 0x65735546
	OVERLAYFS_SUPER_MAGIC = 0x794c7630
	PIPEFS_MAGIC          = 0x50495045
	PROC_SUPER_MAGIC      = 0x9fa0
	RAMFS_MAGIC           = 0x09041934
	SOCKFS_MAGIC          = 0x534F434B
	SYSFS_MAGIC           = 0x62656572
	TMPFS_MAGIC           = 0x01021994
	V9FS_MAGIC            = 0x01021997
)

// Filesystem path limits, from uapi/linux/limits.h.
const (
	NAME_MAX = 255
	PATH_MAX = 4096
)
