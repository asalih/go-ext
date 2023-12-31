package disklayout

import (
	"fmt"
	"os"
)

// InodeType enumerates types of Inodes.
type InodeType int

const (
	// RegularFile is a regular file.
	RegularFile InodeType = iota

	// SpecialFile is a file that doesn't support SeekEnd. It is used for
	// things like proc files.
	SpecialFile

	// Directory is a directory.
	Directory

	// SpecialDirectory is a directory that *does* support SeekEnd. It's
	// the opposite of the SpecialFile scenario above. It similarly
	// supports proc files.
	SpecialDirectory

	// Symlink is a symbolic link.
	Symlink

	// Pipe is a pipe (named or regular).
	Pipe

	// Socket is a socket.
	Socket

	// CharacterDevice is a character device.
	CharacterDevice

	// BlockDevice is a block device.
	BlockDevice

	// Anonymous is an anonymous type when none of the above apply.
	// Epoll fds and event-driven fds fit this category.
	Anonymous
)

// String returns a human-readable representation of the InodeType.
func (n InodeType) String() string {
	switch n {
	case RegularFile, SpecialFile:
		return "file"
	case Directory, SpecialDirectory:
		return "directory"
	case Symlink:
		return "symlink"
	case Pipe:
		return "pipe"
	case Socket:
		return "socket"
	case CharacterDevice:
		return "character-device"
	case BlockDevice:
		return "block-device"
	case Anonymous:
		return "anonymous"
	default:
		return "unknown"
	}
}

// LinuxType returns the linux file type for this inode type.
func (n InodeType) LinuxType() uint32 {
	switch n {
	case RegularFile, SpecialFile:
		return ModeRegular
	case Directory, SpecialDirectory:
		return ModeDirectory
	case Symlink:
		return ModeSymlink
	case Pipe:
		return ModeNamedPipe
	case CharacterDevice:
		return ModeCharacterDevice
	case BlockDevice:
		return ModeBlockDevice
	case Socket:
		return ModeSocket
	default:
		return 0
	}
}

// ToDirentType converts an InodeType to a linux dirent type field.
func ToDirentType(nodeType InodeType) uint8 {
	switch nodeType {
	case RegularFile, SpecialFile:
		return DT_REG
	case Symlink:
		return DT_LNK
	case Directory, SpecialDirectory:
		return DT_DIR
	case Pipe:
		return DT_FIFO
	case CharacterDevice:
		return DT_CHR
	case BlockDevice:
		return DT_BLK
	case Socket:
		return DT_SOCK
	default:
		return DT_UNKNOWN
	}
}

// ToInodeType coverts a linux file type to InodeType.
func ToInodeType(linuxFileType FileMode) InodeType {
	switch linuxFileType {
	case ModeRegular:
		return RegularFile
	case ModeDirectory:
		return Directory
	case ModeSymlink:
		return Symlink
	case ModeNamedPipe:
		return Pipe
	case ModeCharacterDevice:
		return CharacterDevice
	case ModeBlockDevice:
		return BlockDevice
	case ModeSocket:
		return Socket
	default:
		panic(fmt.Sprintf("unknown file mode: %d", linuxFileType))
	}
}

// StableAttr contains Inode attributes that will be stable throughout the
// lifetime of the Inode.
//
// +stateify savable
type StableAttr struct {
	// Type is the InodeType of a InodeOperations.
	Type InodeType

	// DeviceID is the device on which a InodeOperations resides.
	DeviceID uint64

	// InodeID uniquely identifies InodeOperations on its device.
	InodeID uint64

	// BlockSize is the block size of data backing this InodeOperations.
	BlockSize int64

	// DeviceFileMajor is the major device number of this Node, if it is a
	// device file.
	DeviceFileMajor uint16

	// DeviceFileMinor is the minor device number of this Node, if it is a
	// device file.
	DeviceFileMinor uint32
}

// IsRegular returns true if StableAttr.Type matches a regular file.
func IsRegular(s StableAttr) bool {
	return s.Type == RegularFile
}

// IsFile returns true if StableAttr.Type matches any type of file.
func IsFile(s StableAttr) bool {
	return s.Type == RegularFile || s.Type == SpecialFile
}

// IsDir returns true if StableAttr.Type matches any type of directory.
func IsDir(s StableAttr) bool {
	return s.Type == Directory || s.Type == SpecialDirectory
}

// IsSymlink returns true if StableAttr.Type matches a symlink.
func IsSymlink(s StableAttr) bool {
	return s.Type == Symlink
}

// IsPipe returns true if StableAttr.Type matches any type of pipe.
func IsPipe(s StableAttr) bool {
	return s.Type == Pipe
}

// IsAnonymous returns true if StableAttr.Type matches any type of anonymous.
func IsAnonymous(s StableAttr) bool {
	return s.Type == Anonymous
}

// IsSocket returns true if StableAttr.Type matches any type of socket.
func IsSocket(s StableAttr) bool {
	return s.Type == Socket
}

// IsCharDevice returns true if StableAttr.Type matches a character device.
func IsCharDevice(s StableAttr) bool {
	return s.Type == CharacterDevice
}

// UnstableAttr contains Inode attributes that may change over the lifetime
// of the Inode.
//
// +stateify savable
type UnstableAttr struct {
	// Size is the file size in bytes.
	Size int64

	// Usage is the actual data usage in bytes.
	Usage int64

	// Perms is the protection (read/write/execute for user/group/other).
	Perms FilePermissions

	// Owner describes the ownership of this file.
	Owner FileOwner

	// AccessTime is the time of last access
	AccessTime uint32

	// ModificationTime is the time of last modification.
	ModificationTime uint32

	// StatusChangeTime is the time of last attribute modification.
	StatusChangeTime uint32

	// Links is the number of hard links.
	Links uint64
}

// SetOwner sets the owner and group if they are valid.
//
// This method is NOT thread-safe. Callers must prevent concurrent calls.
// func (ua *UnstableAttr) SetOwner(ctx context.Context, owner FileOwner) {
// 	if owner.UID.Ok() {
// 		ua.Owner.UID = owner.UID
// 	}
// 	if owner.GID.Ok() {
// 		ua.Owner.GID = owner.GID
// 	}
// 	ua.StatusChangeTime = ktime.NowFromContext(ctx)
// }

// SetPermissions sets the permissions.
//
// This method is NOT thread-safe. Callers must prevent concurrent calls.
// func (ua *UnstableAttr) SetPermissions(ctx context.Context, p FilePermissions) {
// 	ua.Perms = p
// 	ua.StatusChangeTime = ktime.NowFromContext(ctx)
// }

// // WithCurrentTime returns u with AccessTime == ModificationTime == current time.
// func WithCurrentTime(ctx context.Context, u UnstableAttr) UnstableAttr {
// 	t := ktime.NowFromContext(ctx)
// 	u.AccessTime = t
// 	u.ModificationTime = t
// 	u.StatusChangeTime = t
// 	return u
// }

// AttrMask contains fields to mask StableAttr and UnstableAttr.
//
// +stateify savable
type AttrMask struct {
	Type             bool
	DeviceID         bool
	InodeID          bool
	BlockSize        bool
	Size             bool
	Usage            bool
	Perms            bool
	UID              bool
	GID              bool
	AccessTime       bool
	ModificationTime bool
	StatusChangeTime bool
	Links            bool
}

// Empty returns true if all fields in AttrMask are false.
func (a AttrMask) Empty() bool {
	return a == AttrMask{}
}

// PermMask are file access permissions.
//
// +stateify savable
type PermMask struct {
	// Read indicates reading is permitted.
	Read bool

	// Write indicates writing is permitted.
	Write bool

	// Execute indicates execution is permitted.
	Execute bool
}

// OnlyRead returns true when only the read bit is set.
func (p PermMask) OnlyRead() bool {
	return p.Read && !p.Write && !p.Execute
}

// String implements the fmt.Stringer interface for PermMask.
func (p PermMask) String() string {
	return fmt.Sprintf("PermMask{Read: %v, Write: %v, Execute: %v}", p.Read, p.Write, p.Execute)
}

const (
	S_IROTH = 0x4
	S_IWOTH = 0x2
	S_IXOTH = 0x1
)

// Mode returns the system mode (unix.S_IXOTH, etc.) for these permissions
// in the "other" bits.
func (p PermMask) Mode() (mode os.FileMode) {
	if p.Read {
		mode |= S_IROTH
	}
	if p.Write {
		mode |= S_IWOTH
	}
	if p.Execute {
		mode |= S_IXOTH
	}
	return
}

// SupersetOf returns true iff the permissions in p are a superset of the
// permissions in other.
func (p PermMask) SupersetOf(other PermMask) bool {
	if !p.Read && other.Read {
		return false
	}
	if !p.Write && other.Write {
		return false
	}
	if !p.Execute && other.Execute {
		return false
	}
	return true
}

// FilePermissions represents the permissions of a file, with
// Read/Write/Execute bits for user, group, and other.
//
// +stateify savable
type FilePermissions struct {
	User  PermMask
	Group PermMask
	Other PermMask

	// Sticky, if set on directories, restricts renaming and deletion of
	// files in those directories to the directory owner, file owner, or
	// CAP_FOWNER. The sticky bit is ignored when set on other files.
	Sticky bool

	// SetUID executables can call UID-setting syscalls without CAP_SETUID.
	SetUID bool

	// SetGID executables can call GID-setting syscalls without CAP_SETGID.
	SetGID bool
}

// PermsFromMode takes the Other permissions (last 3 bits) of a FileMode and
// returns a set of PermMask.
func PermsFromMode(mode FileMode) (perms PermMask) {
	perms.Read = mode&ModeOtherRead != 0
	perms.Write = mode&ModeOtherWrite != 0
	perms.Execute = mode&ModeOtherExec != 0
	return
}

// FilePermsFromP9 converts a p9.FileMode to a FilePermissions struct.
// func FilePermsFromP9(mode FileMode) FilePermissions {
// 	return FilePermsFromMode(FileMode(mode))
// }

// FilePermsFromMode converts a system file mode to a FilePermissions struct.
func FilePermsFromMode(mode FileMode) (fp FilePermissions) {
	perm := mode.Permissions()
	fp.Other = PermsFromMode(perm)
	fp.Group = PermsFromMode(perm >> 3)
	fp.User = PermsFromMode(perm >> 6)
	fp.Sticky = mode&ModeSticky == ModeSticky
	fp.SetUID = mode&ModeSetUID == ModeSetUID
	fp.SetGID = mode&ModeSetGID == ModeSetGID
	return
}

// LinuxMode returns the linux mode_t representation of these permissions.
func (f FilePermissions) LinuxMode() FileMode {
	m := FileMode(f.User.Mode()<<6 | f.Group.Mode()<<3 | f.Other.Mode())
	if f.SetUID {
		m |= ModeSetUID
	}
	if f.SetGID {
		m |= ModeSetGID
	}
	if f.Sticky {
		m |= ModeSticky
	}
	return m
}

// OSMode returns the Go runtime's OS independent os.FileMode representation of
// these permissions.
func (f FilePermissions) OSMode() os.FileMode {
	m := os.FileMode(f.User.Mode()<<6 | f.Group.Mode()<<3 | f.Other.Mode())
	if f.SetUID {
		m |= os.ModeSetuid
	}
	if f.SetGID {
		m |= os.ModeSetgid
	}
	if f.Sticky {
		m |= os.ModeSticky
	}
	return m
}

// AnyExecute returns true if any of U/G/O have the execute bit set.
func (f FilePermissions) AnyExecute() bool {
	return f.User.Execute || f.Group.Execute || f.Other.Execute
}

// AnyWrite returns true if any of U/G/O have the write bit set.
func (f FilePermissions) AnyWrite() bool {
	return f.User.Write || f.Group.Write || f.Other.Write
}

// AnyRead returns true if any of U/G/O have the read bit set.
func (f FilePermissions) AnyRead() bool {
	return f.User.Read || f.Group.Read || f.Other.Read
}

// HasSetUIDOrGID returns true if either the setuid or setgid bit is set.
func (f FilePermissions) HasSetUIDOrGID() bool {
	return f.SetUID || f.SetGID
}

// DropSetUIDAndMaybeGID turns off setuid, and turns off setgid if f allows
// group execution.
func (f *FilePermissions) DropSetUIDAndMaybeGID() {
	f.SetUID = false
	if f.Group.Execute {
		f.SetGID = false
	}
}

// FileOwner represents ownership of a file.
//
// +stateify savable
type FileOwner struct {
	UID uint32
	GID uint32
}

// RootOwner corresponds to KUID/KGID 0/0.
var RootOwner = FileOwner{
	UID: 0,
	GID: 0,
}
