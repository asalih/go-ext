package disklayout

import (
	"unsafe"

	"github.com/asalih/go-ext/common"
)

// DirentOld represents the old directory entry struct which does not contain
// the file type. This emulates Linux's ext4_dir_entry struct.
//
// Note: This struct can be of variable size on disk. The one described below
// is of maximum size and the FileName beyond NameLength bytes might contain
// garbage.
//
// +marshal
type DirentOld struct {
	InodeNumber  uint32            `struc:"uint32,little"`
	RecordLength uint16            `struc:"uint16,little"`
	NameLength   uint16            `struc:"uint16,little"`
	FileNameRaw  [MaxFileName]byte `struc:"[]byte"`
}

// Compiles only if DirentOld implements Dirent.
var _ Dirent = (*DirentOld)(nil)

func (d *DirentOld) SizeBytes() int {
	return int(unsafe.Sizeof(*d))
}

func (d *DirentOld) UnmarshalBytes(src []byte) error {
	return common.UnmarshalBytes(d, src)
}

// Inode implements Dirent.Inode.
func (d *DirentOld) Inode() uint32 { return d.InodeNumber }

// RecordSize implements Dirent.RecordSize.
func (d *DirentOld) RecordSize() uint16 { return d.RecordLength }

// FileName implements Dirent.FileName.
func (d *DirentOld) Name() string {
	return string(d.FileNameRaw[:d.NameLength])
}

// FileType implements Dirent.FileType.
func (d *DirentOld) FileType() (InodeType, bool) {
	return Anonymous, false
}
