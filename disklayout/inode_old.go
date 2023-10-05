package disklayout

import (
	"time"
	"unsafe"

	"github.com/asalih/go-ext/common"
	"github.com/asalih/go-ext/linux"
)

const (
	// OldInodeSize is the inode size in ext2/ext3.
	OldInodeSize = 128
)

// InodeOld implements Inode interface. It emulates ext2/ext3 inode struct.
// Inode struct size and record size are both 128 bytes for this.
//
// All fields representing time are in seconds since the epoch. Which means that
// they will overflow in January 2038.
//
// +marshal
type InodeOld struct {
	ModeRaw uint16 `struc:"uint16,little"`
	UIDLo   uint16 `struc:"uint16,little"`
	SizeLo  uint32 `struc:"uint32,little"`

	// The time fields are signed integers because they could be negative to
	// represent time before the epoch.
	AccessTimeRaw       int32 `struc:"uint32,little"`
	ChangeTimeRaw       int32 `struc:"uint32,little"`
	ModificationTimeRaw int32 `struc:"uint32,little"`
	DeletionTimeRaw     int32 `struc:"uint32,little"`

	GIDLo         uint16   `struc:"uint16,little"`
	LinksCountRaw uint16   `struc:"uint16,little"`
	BlocksCountLo uint32   `struc:"uint32,little"`
	FlagsRaw      uint32   `struc:"uint32,little"`
	VersionLo     uint32   `struc:"uint32,little"` // This is OS dependent.
	DataRaw       [60]byte `struc:"[60]byte,little"`
	Generation    uint32   `struc:"uint32,little"`
	FileACLLo     uint32   `struc:"uint32,little"`
	SizeHi        uint32   `struc:"uint32,little"`
	ObsoFaddr     uint32   `struc:"uint32,little"`

	// OS dependent fields have been inlined here.
	BlocksCountHi uint16 `struc:"uint16,little"`
	FileACLHi     uint16 `struc:"uint16,little"`
	UIDHi         uint16 `struc:"uint16,little"`
	GIDHi         uint16 `struc:"uint16,little"`
	ChecksumLo    uint16 `struc:"uint16,little"`
	Unused        uint16 `struc:"uint16,little"`
}

// Compiles only if InodeOld implements Inode.
var _ Inode = (*InodeOld)(nil)

func (in *InodeOld) SizeBytes() int {
	return int(unsafe.Sizeof(*in))
}

func (in *InodeOld) UnmarshalBytes(src []byte) error {
	return common.UnmarshalBytes(in, src)
}

// Mode implements Inode.Mode.
func (in *InodeOld) Mode() linux.FileMode { return linux.FileMode(in.ModeRaw) }

// UID implements Inode.UID.
func (in *InodeOld) UID() uint32 {
	return (uint32(in.UIDHi) << 16) | uint32(in.UIDLo)
}

// GID implements Inode.GID.
func (in *InodeOld) GID() uint32 {
	return (uint32(in.GIDHi) << 16) | uint32(in.GIDLo)
}

// Size implements Inode.Size.
func (in *InodeOld) Size() uint64 {
	// In ext2/ext3, in.SizeHi did not exist, it was instead named in.DirACL.
	return uint64(in.SizeLo)
}

// InodeSize implements Inode.InodeSize.
func (in *InodeOld) InodeSize() uint16 { return OldInodeSize }

// AccessTime implements Inode.AccessTime.
func (in *InodeOld) AccessTime() time.Time {
	return time.Unix(int64(in.AccessTimeRaw), 0)
}

// ChangeTime implements Inode.ChangeTime.
func (in *InodeOld) ChangeTime() time.Time {
	return time.Unix(int64(in.ChangeTimeRaw), 0)
}

// ModificationTime implements Inode.ModificationTime.
func (in *InodeOld) ModificationTime() time.Time {
	return time.Unix(int64(in.ModificationTimeRaw), 0)
}

// DeletionTime implements Inode.DeletionTime.
func (in *InodeOld) DeletionTime() time.Time {
	return time.Unix(int64(in.DeletionTimeRaw), 0)
}

// LinksCount implements Inode.LinksCount.
func (in *InodeOld) LinksCount() uint16 { return in.LinksCountRaw }

// Flags implements Inode.Flags.
func (in *InodeOld) Flags() InodeFlags { return InodeFlagsFromInt(in.FlagsRaw) }

// Data implements Inode.Data.
func (in *InodeOld) Data() []byte { return in.DataRaw[:] }
