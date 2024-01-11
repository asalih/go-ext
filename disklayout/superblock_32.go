package disklayout

import (
	"unsafe"

	"github.com/asalih/go-ext/common"
)

// SuperBlock32Bit implements SuperBlock and represents the 32-bit version of
// the ext4_super_block struct in fs/ext4/ext4.h. Should be used only if
// RevLevel = DynamicRev and 64-bit feature is disabled.
//
// +marshal
type SuperBlock32Bit struct {
	// We embed the old superblock struct here because the 32-bit version is just
	// an extension of the old version.
	SuperBlockOld

	FirstInode         uint32     `struc:"uint32,little"`
	InodeSizeRaw       uint16     `struc:"uint16,little"`
	BlockGroupNumber   uint16     `struc:"uint16,little"`
	FeatureCompat      uint32     `struc:"uint32,little"`
	FeatureIncompat    uint32     `struc:"uint32,little"`
	FeatureRoCompat    uint32     `struc:"uint32,little"`
	UUID               [16]byte   `struc:"[16]byte"`
	VolumeName         [16]byte   `struc:"[16]byte"`
	LastMounted        [64]byte   `struc:"[64]byte"`
	AlgoUsageBitmap    uint32     `struc:"uint32,little"`
	PreallocBlocks     byte       `struc:"byte,little"`
	PreallocDirBlocks  byte       `struc:"byte,little"`
	ReservedGdtBlocks  uint16     `struc:"uint16,little"`
	JournalUUID        [16]byte   `struc:"[16]byte"`
	JournalInum        uint32     `struc:"uint32,little"`
	JournalDev         uint32     `struc:"uint32,little"`
	LastOrphan         uint32     `struc:"uint32,little"`
	HashSeed           [4]uint32  `struc:"[4]uint32,little"`
	DefaultHashVersion byte       `struc:"byte"`
	JnlBackupType      byte       `struc:"byte"`
	BgDescSizeRaw      uint16     `struc:"uint16,little"`
	DefaultMountOpts   uint32     `struc:"uint32,little"`
	FirstMetaBg        uint32     `struc:"uint32,little"`
	MkfsTime           uint32     `struc:"uint32,little"`
	JnlBlocks          [17]uint32 `struc:"[17]uint32,little"`
}

// Compiles only if SuperBlock32Bit implements SuperBlock.
var _ SuperBlock = (*SuperBlock32Bit)(nil)

// Only override methods which change based on the additional fields above.
// Not overriding SuperBlock.BgDescSize because it would still return 32 here.

func (sb *SuperBlock32Bit) SizeBytes() int {
	return int(unsafe.Sizeof(*sb))
}

func (sb *SuperBlock32Bit) UnmarshalBytes(src []byte) error {
	return common.UnmarshalBytes(sb, src)
}

// InodeSize implements SuperBlock.InodeSize.
func (sb *SuperBlock32Bit) InodeSize() uint16 {
	return sb.InodeSizeRaw
}

// CompatibleFeatures implements SuperBlock.CompatibleFeatures.
func (sb *SuperBlock32Bit) CompatibleFeatures() CompatFeatures {
	return CompatFeaturesFromInt(sb.FeatureCompat)
}

// IncompatibleFeatures implements SuperBlock.IncompatibleFeatures.
func (sb *SuperBlock32Bit) IncompatibleFeatures() IncompatFeatures {
	return IncompatFeaturesFromInt(sb.FeatureIncompat)
}

// ReadOnlyCompatibleFeatures implements SuperBlock.ReadOnlyCompatibleFeatures.
func (sb *SuperBlock32Bit) ReadOnlyCompatibleFeatures() RoCompatFeatures {
	return RoCompatFeaturesFromInt(sb.FeatureRoCompat)
}

// ExtType implements SuperBlock.ExtType
func (sb *SuperBlock32Bit) ExtType() ExtType {
	return getExtType(sb.FeatureCompat, sb.FeatureIncompat)
}
