package disklayout

import (
	"unsafe"

	"github.com/asalih/go-ext/common"
)

// DirentNew represents the ext4 directory entry struct. This emulates Linux's
// ext4_dir_entry_2 struct. The FileName can not be more than 255 bytes so we
// only need 8 bits to store the NameLength. As a result, NameLength has been
// shortened and the other 8 bits are used to encode the file type. Use the
// FileTypeRaw field only if the SbDirentFileType feature is set.
//
// Note: This struct can be of variable size on disk. The one described below
// is of maximum size and the FileName beyond NameLength bytes might contain
// garbage.
//
// +marshal
type DirentNew struct {
	InodeNumber  uint32            `struc:"uint32,little"`
	RecordLength uint16            `struc:"uint16,little"`
	NameLength   uint8             `struc:"uint8,sizeof=FileNameRaw"`
	FileTypeRaw  uint8             `struc:"uint8"`
	FileNameRaw  [MaxFileName]byte `struc:"[]byte"`
}

// Compiles only if DirentNew implements Dirent.
var _ Dirent = (*DirentNew)(nil)

func (d *DirentNew) SizeBytes() int {
	return int(unsafe.Sizeof(*d))
}

func (d *DirentNew) UnmarshalBytes(src []byte) error {
	return common.UnmarshalBytes(d, src)
}

// Inode implements Dirent.Inode.
func (d *DirentNew) Inode() uint32 { return d.InodeNumber }

// RecordSize implements Dirent.RecordSize.
func (d *DirentNew) RecordSize() uint16 { return d.RecordLength }

// Name implements Dirent.FileName.
func (d *DirentNew) Name() string {
	return string(d.FileNameRaw[:d.NameLength])
}

// FileType implements Dirent.FileType.
func (d *DirentNew) FileType() (InodeType, bool) {
	if inodeType, ok := inodeTypeByFileType[d.FileTypeRaw]; ok {
		return inodeType, true
	}

	return InodeType(d.FileTypeRaw), false
}
