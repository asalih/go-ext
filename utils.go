package ext

import (
	"io"

	"github.com/asalih/go-ext/common"
	"github.com/asalih/go-ext/disklayout"
	"github.com/asalih/go-ext/syserror"
)

// readFromDisk performs a binary read from disk into the given struct from
// the absolute offset provided.
func readFromDisk(dev io.ReaderAt, abOff int64, v common.Unmarshal) error {
	n := v.SizeBytes()
	buf := make([]byte, n)
	if read, _ := dev.ReadAt(buf, abOff); read < int(n) {
		return syserror.EIO
	}

	return v.UnmarshalBytes(buf)
}

// readSuperBlock reads the SuperBlock from block group 0 in the underlying
// device. There are three versions of the superblock. This function identifies
// and returns the correct version.
func readSuperBlock(dev io.ReaderAt) (disklayout.SuperBlock, error) {
	var sb disklayout.SuperBlock = &disklayout.SuperBlockOld{}
	if err := readFromDisk(dev, disklayout.SbOffset, sb); err != nil {
		return nil, err
	}
	if sb.Revision() == disklayout.OldRev {
		return sb, nil
	}

	sb = &disklayout.SuperBlock32Bit{}
	if err := readFromDisk(dev, disklayout.SbOffset, sb); err != nil {
		return nil, err
	}
	if !sb.IncompatibleFeatures().Is64Bit {
		return sb, nil
	}

	sb = &disklayout.SuperBlock64Bit{}
	if err := readFromDisk(dev, disklayout.SbOffset, sb); err != nil {
		return nil, err
	}
	return sb, nil
}

// blockGroupsCount returns the number of block groups in the ext fs.
func blockGroupsCount(sb disklayout.SuperBlock) uint64 {
	blocksCount := sb.BlocksCount()
	blocksPerGroup := uint64(sb.BlocksPerGroup())

	// Round up the result. float64 can compromise precision so do it manually.
	return (blocksCount + blocksPerGroup - 1) / blocksPerGroup
}

// readBlockGroups reads the block group descriptor table from block group 0 in
// the underlying device.
func readBlockGroups(dev io.ReaderAt, sb disklayout.SuperBlock) ([]disklayout.BlockGroup, error) {
	bgCount := blockGroupsCount(sb)
	bgdSize := uint64(sb.BgDescSize())
	is64Bit := sb.IncompatibleFeatures().Is64Bit
	bgds := make([]disklayout.BlockGroup, bgCount)

	for i, off := uint64(0), uint64(sb.FirstDataBlock()+1)*sb.BlockSize(); i < bgCount; i, off = i+1, off+bgdSize {
		if is64Bit {
			bgds[i] = &disklayout.BlockGroup64Bit{}
		} else {
			bgds[i] = &disklayout.BlockGroup32Bit{}
		}

		if err := readFromDisk(dev, int64(off), bgds[i]); err != nil {
			return nil, err
		}
	}
	return bgds, nil
}
