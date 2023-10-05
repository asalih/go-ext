package disklayout

import (
	"unsafe"

	"github.com/asalih/go-ext/common"
)

// BlockGroup64Bit emulates struct ext4_group_desc in fs/ext4/ext4.h.
// It is the block group descriptor struct for 64-bit ext4 filesystems.
// It implements BlockGroup interface. It is an extension of the 32-bit
// version of BlockGroup.
//
// +marshal
type BlockGroup64Bit struct {
	// We embed the 32-bit struct here because 64-bit version is just an extension
	// of the 32-bit version.
	BlockGroup32Bit

	// 64-bit specific fields.
	BlockBitmapHi         uint32 `struc:"uint32,little"`
	InodeBitmapHi         uint32 `struc:"uint32,little"`
	InodeTableHi          uint32 `struc:"uint32,little"`
	FreeBlocksCountHi     uint16 `struc:"uint16,little"`
	FreeInodesCountHi     uint16 `struc:"uint16,little"`
	UsedDirsCountHi       uint16 `struc:"uint16,little"`
	ItableUnusedHi        uint16 `struc:"uint16,little"`
	ExcludeBitmapHi       uint32 `struc:"uint32,little"`
	BlockBitmapChecksumHi uint16 `struc:"uint16,little"`
	InodeBitmapChecksumHi uint16 `struc:"uint16,little"`
	Reserved              uint32 `struc:"uint32,little"`
}

// Compiles only if BlockGroup64Bit implements BlockGroup.
var _ BlockGroup = (*BlockGroup64Bit)(nil)

func (bg *BlockGroup64Bit) SizeBytes() int {
	return int(unsafe.Sizeof(*bg))
}

func (bg *BlockGroup64Bit) UnmarshalBytes(src []byte) error {
	return common.UnmarshalBytes(bg, src)
}

// Methods to override. Checksum() and Flags() are not overridden.

// InodeTable implements BlockGroup.InodeTable.
func (bg *BlockGroup64Bit) InodeTable() uint64 {
	return (uint64(bg.InodeTableHi) << 32) | uint64(bg.InodeTableLo)
}

// BlockBitmap implements BlockGroup.BlockBitmap.
func (bg *BlockGroup64Bit) BlockBitmap() uint64 {
	return (uint64(bg.BlockBitmapHi) << 32) | uint64(bg.BlockBitmapLo)
}

// InodeBitmap implements BlockGroup.InodeBitmap.
func (bg *BlockGroup64Bit) InodeBitmap() uint64 {
	return (uint64(bg.InodeBitmapHi) << 32) | uint64(bg.InodeBitmapLo)
}

// ExclusionBitmap implements BlockGroup.ExclusionBitmap.
func (bg *BlockGroup64Bit) ExclusionBitmap() uint64 {
	return (uint64(bg.ExcludeBitmapHi) << 32) | uint64(bg.ExcludeBitmapLo)
}

// FreeBlocksCount implements BlockGroup.FreeBlocksCount.
func (bg *BlockGroup64Bit) FreeBlocksCount() uint32 {
	return (uint32(bg.FreeBlocksCountHi) << 16) | uint32(bg.FreeBlocksCountLo)
}

// FreeInodesCount implements BlockGroup.FreeInodesCount.
func (bg *BlockGroup64Bit) FreeInodesCount() uint32 {
	return (uint32(bg.FreeInodesCountHi) << 16) | uint32(bg.FreeInodesCountLo)
}

// DirectoryCount implements BlockGroup.DirectoryCount.
func (bg *BlockGroup64Bit) DirectoryCount() uint32 {
	return (uint32(bg.UsedDirsCountHi) << 16) | uint32(bg.UsedDirsCountLo)
}

// UnusedInodeCount implements BlockGroup.UnusedInodeCount.
func (bg *BlockGroup64Bit) UnusedInodeCount() uint32 {
	return (uint32(bg.ItableUnusedHi) << 16) | uint32(bg.ItableUnusedLo)
}

// BlockBitmapChecksum implements BlockGroup.BlockBitmapChecksum.
func (bg *BlockGroup64Bit) BlockBitmapChecksum() uint32 {
	return (uint32(bg.BlockBitmapChecksumHi) << 16) | uint32(bg.BlockBitmapChecksumLo)
}

// InodeBitmapChecksum implements BlockGroup.InodeBitmapChecksum.
func (bg *BlockGroup64Bit) InodeBitmapChecksum() uint32 {
	return (uint32(bg.InodeBitmapChecksumHi) << 16) | uint32(bg.InodeBitmapChecksumLo)
}
