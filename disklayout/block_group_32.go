package disklayout

import (
	"unsafe"

	"github.com/asalih/go-ext/common"
)

// BlockGroup32Bit emulates the first half of struct ext4_group_desc in
// fs/ext4/ext4.h. It is the block group descriptor struct for ext2, ext3 and
// 32-bit ext4 filesystems. It implements BlockGroup interface.
//
// +marshal
type BlockGroup32Bit struct {
	BlockBitmapLo         uint32 `struc:"uint32,little"`
	InodeBitmapLo         uint32 `struc:"uint32,little"`
	InodeTableLo          uint32 `struc:"uint32,little"`
	FreeBlocksCountLo     uint16 `struc:"uint16,little"`
	FreeInodesCountLo     uint16 `struc:"uint16,little"`
	UsedDirsCountLo       uint16 `struc:"uint16,little"`
	FlagsRaw              uint16 `struc:"uint16,little"`
	ExcludeBitmapLo       uint32 `struc:"uint32,little"`
	BlockBitmapChecksumLo uint16 `struc:"uint16,little"`
	InodeBitmapChecksumLo uint16 `struc:"uint16,little"`
	ItableUnusedLo        uint16 `struc:"uint16,little"`
	ChecksumRaw           uint16 `struc:"uint16,little"`
}

// Compiles only if BlockGroup32Bit implements BlockGroup.
var _ BlockGroup = (*BlockGroup32Bit)(nil)

func (bg *BlockGroup32Bit) SizeBytes() int {
	return int(unsafe.Sizeof(*bg))
}

func (bg *BlockGroup32Bit) UnmarshalBytes(src []byte) error {
	return common.UnmarshalBytes(bg, src)
}

// InodeTable implements BlockGroup.InodeTable.
func (bg *BlockGroup32Bit) InodeTable() uint64 { return uint64(bg.InodeTableLo) }

// BlockBitmap implements BlockGroup.BlockBitmap.
func (bg *BlockGroup32Bit) BlockBitmap() uint64 { return uint64(bg.BlockBitmapLo) }

// InodeBitmap implements BlockGroup.InodeBitmap.
func (bg *BlockGroup32Bit) InodeBitmap() uint64 { return uint64(bg.InodeBitmapLo) }

// ExclusionBitmap implements BlockGroup.ExclusionBitmap.
func (bg *BlockGroup32Bit) ExclusionBitmap() uint64 { return uint64(bg.ExcludeBitmapLo) }

// FreeBlocksCount implements BlockGroup.FreeBlocksCount.
func (bg *BlockGroup32Bit) FreeBlocksCount() uint32 { return uint32(bg.FreeBlocksCountLo) }

// FreeInodesCount implements BlockGroup.FreeInodesCount.
func (bg *BlockGroup32Bit) FreeInodesCount() uint32 { return uint32(bg.FreeInodesCountLo) }

// DirectoryCount implements BlockGroup.DirectoryCount.
func (bg *BlockGroup32Bit) DirectoryCount() uint32 { return uint32(bg.UsedDirsCountLo) }

// UnusedInodeCount implements BlockGroup.UnusedInodeCount.
func (bg *BlockGroup32Bit) UnusedInodeCount() uint32 { return uint32(bg.ItableUnusedLo) }

// BlockBitmapChecksum implements BlockGroup.BlockBitmapChecksum.
func (bg *BlockGroup32Bit) BlockBitmapChecksum() uint32 { return uint32(bg.BlockBitmapChecksumLo) }

// InodeBitmapChecksum implements BlockGroup.InodeBitmapChecksum.
func (bg *BlockGroup32Bit) InodeBitmapChecksum() uint32 { return uint32(bg.InodeBitmapChecksumLo) }

// Checksum implements BlockGroup.Checksum.
func (bg *BlockGroup32Bit) Checksum() uint16 { return bg.ChecksumRaw }

// Flags implements BlockGroup.Flags.
func (bg *BlockGroup32Bit) Flags() BGFlags { return BGFlagsFromInt(bg.FlagsRaw) }
