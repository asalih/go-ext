package disklayout

import (
	"time"
	"unsafe"

	"github.com/asalih/go-ext/common"
)

// InodeNew represents ext4 inode structure which can be bigger than
// OldInodeSize. The actual size of this struct should be determined using
// inode.ExtraInodeSize. Accessing any field here should be verified with the
// actual size. The extra space between the end of the inode struct and end of
// the inode record can be used to store extended attr.
//
// If the TimeExtra fields are in scope, the lower 2 bits of those are used
// to extend their counter part to be 34 bits wide; the rest (upper) 30 bits
// are used to provide nanoscond precision. Hence, these timestamps will now
// overflow in May 2446.
// See https://www.kernel.org/doc/html/latest/filesystems/ext4/dynamic.html#inode-timestamps.
//
// +marshal
type InodeNew struct {
	InodeOld

	ExtraInodeSize        uint16 `struc:"uint16,little"`
	ChecksumHi            uint16 `struc:"uint16,little"`
	ChangeTimeExtra       uint32 `struc:"uint32,little"`
	ModificationTimeExtra uint32 `struc:"uint32,little"`
	AccessTimeExtra       uint32 `struc:"uint32,little"`
	CreationTime          uint32 `struc:"uint32,little"`
	CreationTimeExtra     uint32 `struc:"uint32,little"`
	VersionHi             uint32 `struc:"uint32,little"`
	ProjectID             uint32 `struc:"uint32,little"`
}

// Compiles only if InodeNew implements Inode.
var _ Inode = (*InodeNew)(nil)

func (sb *InodeNew) SizeBytes() int {
	return int(unsafe.Sizeof(*sb))
}

func (in *InodeNew) UnmarshalBytes(src []byte) error {
	return common.UnmarshalBytes(in, src)
}

// fromExtraTime decodes the extra time and constructs the kernel time struct
// with nanosecond precision.
func fromExtraTime(lo int32, extra uint32) time.Time {
	// See description above InodeNew for format.
	seconds := (int64(extra&0x3) << 32) + int64(lo)
	nanoseconds := int64(extra >> 2)
	return time.Unix(seconds, nanoseconds)
}

// Only override methods which change due to ext4 specific fields.

// Size implements Inode.Size.
func (in *InodeNew) Size() uint64 {
	return (uint64(in.SizeHi) << 32) | uint64(in.SizeLo)
}

// InodeSize implements Inode.InodeSize.
func (in *InodeNew) InodeSize() uint16 {
	return OldInodeSize + in.ExtraInodeSize
}

// ChangeTime implements Inode.ChangeTime.
func (in *InodeNew) ChangeTime() time.Time {
	// Apply new timestamp logic if inode.ChangeTimeExtra is in scope.
	if in.ExtraInodeSize >= 8 {
		return fromExtraTime(in.ChangeTimeRaw, in.ChangeTimeExtra)
	}

	return in.InodeOld.ChangeTime()
}

// ModificationTime implements Inode.ModificationTime.
func (in *InodeNew) ModificationTime() time.Time {
	// Apply new timestamp logic if inode.ModificationTimeExtra is in scope.
	if in.ExtraInodeSize >= 12 {
		return fromExtraTime(in.ModificationTimeRaw, in.ModificationTimeExtra)
	}

	return in.InodeOld.ModificationTime()
}

// AccessTime implements Inode.AccessTime.
func (in *InodeNew) AccessTime() time.Time {
	// Apply new timestamp logic if inode.AccessTimeExtra is in scope.
	if in.ExtraInodeSize >= 16 {
		return fromExtraTime(in.AccessTimeRaw, in.AccessTimeExtra)
	}

	return in.InodeOld.AccessTime()
}
