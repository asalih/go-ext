package disklayout

import (
	"unsafe"

	"github.com/asalih/go-ext/common"
)

// Extents were introduced in ext4 and provide huge performance gains in terms
// data locality and reduced metadata block usage. Extents are organized in
// extent trees. The root node is contained in inode.BlocksRaw.
//
// Terminology:
//   - Physical Block:
//       Filesystem data block which is addressed normally wrt the entire
//       filesystem (addressed with 48 bits).
//
//   - File Block:
//       Data block containing *only* file data and addressed wrt to the file
//       with only 32 bits. The (i)th file block contains file data from
//       byte (i * sb.BlockSize()) to ((i+1) * sb.BlockSize()).

const (
	// ExtentHeaderSize is the size of the header of an extent tree node.
	ExtentHeaderSize = 12

	// ExtentEntrySize is the size of an entry in an extent tree node.
	// This size is the same for both leaf and internal nodes.
	ExtentEntrySize = 12

	// ExtentMagic is the magic number which must be present in the header.
	ExtentMagic = 0xf30a
)

// ExtentEntryPair couples an in-memory ExtendNode with the ExtentEntry that
// points to it. We want to cache these structs in memory to avoid repeated
// disk reads.
//
// Note: This struct itself does not represent an on-disk struct.
type ExtentEntryPair struct {
	// Entry points to the child node on disk.
	Entry ExtentEntry
	// Node points to child node in memory. Is nil if the current node is a leaf.
	Node *ExtentNode
}

// ExtentNode represents an extent tree node. For internal nodes, all Entries
// will be ExtendIdxs. For leaf nodes, they will all be Extents.
//
// Note: This struct itself does not represent an on-disk struct.
type ExtentNode struct {
	Header  ExtentHeader
	Entries []ExtentEntryPair
}

// ExtentEntry represents an extent tree node entry. The entry can either be
// an ExtentIdx or Extent itself. This exists to simplify navigation logic.
type ExtentEntry interface {
	common.Unmarshal

	// FileBlock returns the first file block number covered by this entry.
	FileBlock() uint32

	// PhysicalBlock returns the child physical block that this entry points to.
	PhysicalBlock() uint64
}

// ExtentHeader emulates the ext4_extent_header struct in ext4. Each extent
// tree node begins with this and is followed by `NumEntries` number of:
//   - Extent         if `Depth` == 0
//   - ExtentIdx      otherwise
//
// +marshal
type ExtentHeader struct {
	// Magic in the extent magic number, must be 0xf30a.
	Magic uint16 `struc:"uint16,little"`

	// NumEntries indicates the number of valid entries following the header.
	NumEntries uint16 `struc:"uint16,little"`

	// MaxEntries that could follow the header. Used while adding entries.
	MaxEntries uint16 `struc:"uint16,little"`

	// Height represents the distance of this node from the farthest leaf. Please
	// note that Linux incorrectly calls this `Depth` (which means the distance
	// of the node from the root).
	Height     uint16 `struc:"uint16,little"`
	Generation uint32 `struc:"uint32,little"`
}

func (ex *ExtentHeader) SizeBytes() int {
	return int(unsafe.Sizeof(*ex))
}

func (ex *ExtentHeader) UnmarshalBytes(src []byte) error {
	return common.UnmarshalBytes(ex, src)
}

// ExtentIdx emulates the ext4_extent_idx struct in ext4. Only present in
// internal nodes. Sorted in ascending order based on FirstFileBlock since
// Linux does a binary search on this. This points to a block containing the
// child node.
//
// +marshal
type ExtentIdx struct {
	FirstFileBlock uint32 `struc:"uint32,little"`
	ChildBlockLo   uint32 `struc:"uint32,little"`
	ChildBlockHi   uint16 `struc:"uint16,little"`
	Reserved       uint16 `struc:"uint16,little"`
}

// Compiles only if ExtentIdx implements ExtentEntry.
var _ ExtentEntry = (*ExtentIdx)(nil)

// FileBlock implements ExtentEntry.FileBlock.
func (ei *ExtentIdx) FileBlock() uint32 {
	return ei.FirstFileBlock
}

// PhysicalBlock implements ExtentEntry.PhysicalBlock. It returns the
// physical block number of the child block.
func (ei *ExtentIdx) PhysicalBlock() uint64 {
	return (uint64(ei.ChildBlockHi) << 32) | uint64(ei.ChildBlockLo)
}

func (ei *ExtentIdx) SizeBytes() int {
	return int(unsafe.Sizeof(*ei))
}

func (ei *ExtentIdx) UnmarshalBytes(src []byte) error {
	return common.UnmarshalBytes(ei, src)
}

// Extent represents the ext4_extent struct in ext4. Only present in leaf
// nodes. Sorted in ascending order based on FirstFileBlock since Linux does a
// binary search on this. This points to an array of data blocks containing the
// file data. It covers `Length` data blocks starting from `StartBlock`.
//
// +marshal
type Extent struct {
	FirstFileBlock uint32 `struc:"uint32,little"`
	Length         uint16 `struc:"uint16,little"`
	StartBlockHi   uint16 `struc:"uint16,little"`
	StartBlockLo   uint32 `struc:"uint32,little"`
}

// Compiles only if Extent implements ExtentEntry.
var _ ExtentEntry = (*Extent)(nil)

func (sb *Extent) SizeBytes() int {
	return int(unsafe.Sizeof(*sb))
}

func (e *Extent) UnmarshalBytes(src []byte) error {
	return common.UnmarshalBytes(e, src)
}

// FileBlock implements ExtentEntry.FileBlock.
func (e *Extent) FileBlock() uint32 {
	return e.FirstFileBlock
}

// PhysicalBlock implements ExtentEntry.PhysicalBlock. It returns the
// physical block number of the first data block this extent covers.
func (e *Extent) PhysicalBlock() uint64 {
	return (uint64(e.StartBlockHi) << 32) | uint64(e.StartBlockLo)
}
