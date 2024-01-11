package disklayout

import (
	"unsafe"

	"github.com/asalih/go-ext/common"
)

// SuperBlock64Bit implements SuperBlock and represents the 64-bit version of
// the ext4_super_block struct in fs/ext4/ext4.h. This sums up to be exactly
// 1024 bytes (smallest possible block size) and hence the superblock always
// fits in no more than one data block. Should only be used when the 64-bit
// feature is set.
//
// +marshal
type SuperBlock64Bit struct {
	// We embed the 32-bit struct here because 64-bit version is just an extension
	// of the 32-bit version.
	SuperBlock32Bit

	BlocksCountHi           uint32     `struc:"uint32,little"`
	ReservedBlocksCountHi   uint32     `struc:"uint32,little"`
	FreeBlocksCountHi       uint32     `struc:"uint32,little"`
	MinInodeSize            uint16     `struc:"uint16,little"`
	WantInodeSize           uint16     `struc:"uint16,little"`
	Flags                   uint32     `struc:"uint32,little"`
	RaidStride              uint16     `struc:"uint16,little"`
	MmpInterval             uint16     `struc:"uint16,little"`
	MmpBlock                uint64     `struc:"uint64,little"`
	RaidStripeWidth         uint32     `struc:"uint32,little"`
	LogGroupsPerFlex        byte       `struc:"byte"`
	ChecksumType            byte       `struc:"byte"`
	EncryptionLevel         byte       `struc:"byte"`
	ReservedPad             byte       `struc:"byte"`
	KbytesWritten           uint64     `struc:"uint64,little"`
	SnapshotInum            uint32     `struc:"uint32,little"`
	SnapshotID              uint32     `struc:"uint32,little"`
	SnapshotRsrvBlocksCount uint64     `struc:"uint64,little"`
	SnapshotList            uint32     `struc:"uint32,little"`
	ErrorCount              uint32     `struc:"uint32,little"`
	FirstErrorTime          uint32     `struc:"uint32,little"`
	FirstErrorInode         uint32     `struc:"uint32,little"`
	FirstErrorBlock         uint64     `struc:"uint64,little"`
	FirstErrorFunction      [32]byte   `struc:"[32]pad"`
	FirstErrorLine          uint32     `struc:"uint32,little"`
	LastErrorTime           uint32     `struc:"uint32,little"`
	LastErrorInode          uint32     `struc:"uint32,little"`
	LastErrorLine           uint32     `struc:"uint32,little"`
	LastErrorBlock          uint64     `struc:"uint64,little"`
	LastErrorFunction       [32]byte   `struc:"[32]pad"`
	MountOpts               [64]byte   `struc:"[64]pad"`
	UserQuotaInum           uint32     `struc:"uint32,little"`
	GroupQuotaInum          uint32     `struc:"uint32,little"`
	OverheadBlocks          uint32     `struc:"uint32,little"`
	BackupBgs               [2]uint32  `struc:"[2]uint32,little"`
	EncryptAlgos            [4]byte    `struc:"[4]pad"`
	EncryptPwSalt           [16]byte   `struc:"[16]pad"`
	LostFoundInode          uint32     `struc:"uint32,little"`
	ProjectQuotaInode       uint32     `struc:"uint32,little"`
	ChecksumSeed            uint32     `struc:"uint32,little"`
	WtimeHi                 byte       `struc:"byte"`
	MtimeHi                 byte       `struc:"byte"`
	MkfsTimeHi              byte       `struc:"byte"`
	LastCheckHi             byte       `struc:"byte"`
	FirstErrorTimeHi        byte       `struc:"byte"`
	LastErrorTimeHi         byte       `struc:"byte"`
	Reserved1               [2]byte    `struc:"[2]pad"`
	Encoding                uint16     `struc:"uint16,little"`
	EncodingFlags           uint16     `struc:"uint16,little"`
	Reserved2               [95]uint32 `struc:"[95]uint32,little"`
	Checksum                uint32     `struc:"uint32,little"`
}

// Compiles only if SuperBlock64Bit implements SuperBlock.
var _ SuperBlock = (*SuperBlock64Bit)(nil)

func (sb *SuperBlock64Bit) SizeBytes() int {
	return int(unsafe.Sizeof(*sb))
}

func (sb *SuperBlock64Bit) UnmarshalBytes(src []byte) error {
	return common.UnmarshalBytes(sb, src)
}

// Only override methods which change based on the 64-bit feature.

// BlocksCount implements SuperBlock.BlocksCount.
func (sb *SuperBlock64Bit) BlocksCount() uint64 {
	return (uint64(sb.BlocksCountHi) << 32) | uint64(sb.BlocksCountLo)
}

// FreeBlocksCount implements SuperBlock.FreeBlocksCount.
func (sb *SuperBlock64Bit) FreeBlocksCount() uint64 {
	return (uint64(sb.FreeBlocksCountHi) << 32) | uint64(sb.FreeBlocksCountLo)
}

// BgDescSize implements SuperBlock.BgDescSize.
func (sb *SuperBlock64Bit) BgDescSize() uint16 { return sb.BgDescSizeRaw }

// ExtType implements SuperBlock.ExtType
func (sb *SuperBlock64Bit) ExtType() ExtType {
	return getExtType(sb.FeatureCompat, sb.FeatureIncompat)
}
