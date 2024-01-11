package disklayout

import (
	"unsafe"

	"github.com/asalih/go-ext/common"
)

// SuperBlockOld implements SuperBlock and represents the old version of the
// superblock struct. Should be used only if RevLevel = OldRev.
//
// +marshal
type SuperBlockOld struct {
	InodesCountRaw      uint32 `struc:"uint32,little"`
	BlocksCountLo       uint32 `struc:"uint32,little"`
	ReservedBlocksCount uint32 `struc:"uint32,little"`
	FreeBlocksCountLo   uint32 `struc:"uint32,little"`
	FreeInodesCountRaw  uint32 `struc:"uint32,little"`
	FirstDataBlockRaw   uint32 `struc:"uint32,little"`
	LogBlockSize        uint32 `struc:"uint32,little"`
	LogClusterSize      uint32 `struc:"uint32,little"`
	BlocksPerGroupRaw   uint32 `struc:"uint32,little"`
	ClustersPerGroupRaw uint32 `struc:"uint32,little"`
	InodesPerGroupRaw   uint32 `struc:"uint32,little"`
	Mtime               uint32 `struc:"uint32,little"`
	Wtime               uint32 `struc:"uint32,little"`
	MountCountRaw       uint16 `struc:"uint16,little"`
	MaxMountCountRaw    uint16 `struc:"uint16,little"`
	MagicRaw            uint16 `struc:"uint16,little"`
	State               uint16 `struc:"uint16,little"`
	Errors              uint16 `struc:"uint16,little"`
	MinorRevLevel       uint16 `struc:"uint16,little"`
	LastCheck           uint32 `struc:"uint32,little"`
	CheckInterval       uint32 `struc:"uint32,little"`
	CreatorOS           uint32 `struc:"uint32,little"`
	RevLevel            uint32 `struc:"uint32,little"`
	DefResUID           uint16 `struc:"uint16,little"`
	DefResGID           uint16 `struc:"uint16,little"`
}

// Compiles only if SuperBlockOld implements SuperBlock.
var _ SuperBlock = (*SuperBlockOld)(nil)

func (sb *SuperBlockOld) SizeBytes() int {
	return int(unsafe.Sizeof(*sb))
}

func (sb *SuperBlockOld) UnmarshalBytes(src []byte) error {
	return common.UnmarshalBytes(sb, src)
}

// ExtType implements SuperBlock.ExtType
func (sb *SuperBlockOld) ExtType() ExtType { return 0 }

// InodesCount implements SuperBlock.InodesCount.
func (sb *SuperBlockOld) InodesCount() uint32 { return sb.InodesCountRaw }

// BlocksCount implements SuperBlock.BlocksCount.
func (sb *SuperBlockOld) BlocksCount() uint64 { return uint64(sb.BlocksCountLo) }

// FreeBlocksCount implements SuperBlock.FreeBlocksCount.
func (sb *SuperBlockOld) FreeBlocksCount() uint64 { return uint64(sb.FreeBlocksCountLo) }

// FreeInodesCount implements SuperBlock.FreeInodesCount.
func (sb *SuperBlockOld) FreeInodesCount() uint32 { return sb.FreeInodesCountRaw }

// MountCount implements SuperBlock.MountCount.
func (sb *SuperBlockOld) MountCount() uint16 { return sb.MountCountRaw }

// MaxMountCount implements SuperBlock.MaxMountCount.
func (sb *SuperBlockOld) MaxMountCount() uint16 { return sb.MaxMountCountRaw }

// FirstDataBlock implements SuperBlock.FirstDataBlock.
func (sb *SuperBlockOld) FirstDataBlock() uint32 { return sb.FirstDataBlockRaw }

// BlockSize implements SuperBlock.BlockSize.
func (sb *SuperBlockOld) BlockSize() uint64 { return 1 << (10 + sb.LogBlockSize) }

// BlocksPerGroup implements SuperBlock.BlocksPerGroup.
func (sb *SuperBlockOld) BlocksPerGroup() uint32 { return sb.BlocksPerGroupRaw }

// ClusterSize implements SuperBlock.ClusterSize.
func (sb *SuperBlockOld) ClusterSize() uint64 { return 1 << (10 + sb.LogClusterSize) }

// ClustersPerGroup implements SuperBlock.ClustersPerGroup.
func (sb *SuperBlockOld) ClustersPerGroup() uint32 { return sb.ClustersPerGroupRaw }

// InodeSize implements SuperBlock.InodeSize.
func (sb *SuperBlockOld) InodeSize() uint16 { return OldInodeSize }

// InodesPerGroup implements SuperBlock.InodesPerGroup.
func (sb *SuperBlockOld) InodesPerGroup() uint32 { return sb.InodesPerGroupRaw }

// BgDescSize implements SuperBlock.BgDescSize.
func (sb *SuperBlockOld) BgDescSize() uint16 { return 32 }

// CompatibleFeatures implements SuperBlock.CompatibleFeatures.
func (sb *SuperBlockOld) CompatibleFeatures() CompatFeatures { return CompatFeatures{} }

// IncompatibleFeatures implements SuperBlock.IncompatibleFeatures.
func (sb *SuperBlockOld) IncompatibleFeatures() IncompatFeatures { return IncompatFeatures{} }

// ReadOnlyCompatibleFeatures implements SuperBlock.ReadOnlyCompatibleFeatures.
func (sb *SuperBlockOld) ReadOnlyCompatibleFeatures() RoCompatFeatures { return RoCompatFeatures{} }

// Magic implements SuperBlock.Magic.
func (sb *SuperBlockOld) Magic() uint16 { return sb.MagicRaw }

// Revision implements SuperBlock.Revision.
func (sb *SuperBlockOld) Revision() SbRevision { return SbRevision(sb.RevLevel) }
