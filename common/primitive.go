package common

import "encoding/binary"

// Int8 is a marshal.Marshallable implementation for int8.
type Int8 int8

// Uint8 is a marshal.Marshallable implementation for uint8.
type Uint8 uint8

// Int16 is a marshal.Marshallable implementation for int16.
type Int16 int16

// Uint16 is a marshal.Marshallable implementation for uint16.
type Uint16 uint16

// Int32 is a marshal.Marshallable implementation for int32.
type Int32 int32

// Uint32 is a marshal.Marshallable implementation for uint32.
type Uint32 uint32

// Int64 is a marshal.Marshallable implementation for int64.
type Int64 int64

// Uint64 is a marshal.Marshallable implementation for uint64.
type Uint64 uint64

// SizeBytes implements marshal.Marshallable.SizeBytes.
//
//go:nosplit
func (i *Int16) SizeBytes() int {
	return 2
}

// UnmarshalBytes implements marshal.Marshallable.UnmarshalBytes.
func (i *Int16) UnmarshalBytes(src []byte) error {
	*i = Int16(int16(binary.LittleEndian.Uint16(src[:2])))
	return nil
}

// SizeBytes implements marshal.Marshallable.SizeBytes.
//
//go:nosplit
func (i *Int32) SizeBytes() int {
	return 4
}

// UnmarshalBytes implements marshal.Marshallable.UnmarshalBytes.
func (i *Int32) UnmarshalBytes(src []byte) error {
	*i = Int32(int32(binary.LittleEndian.Uint32(src[:4])))
	return nil
}

// SizeBytes implements marshal.Marshallable.SizeBytes.
//
//go:nosplit
func (i *Int64) SizeBytes() int {
	return 8
}

// UnmarshalBytes implements marshal.Marshallable.UnmarshalBytes.
func (i *Int64) UnmarshalBytes(src []byte) error {
	*i = Int64(int64(binary.LittleEndian.Uint64(src[:8])))
	return nil
}

// SizeBytes implements marshal.Marshallable.SizeBytes.
//
//go:nosplit
func (i *Int8) SizeBytes() int {
	return 1
}

// UnmarshalBytes implements marshal.Marshallable.UnmarshalBytes.
func (i *Int8) UnmarshalBytes(src []byte) error {
	*i = Int8(int8(src[0]))
	return nil
}

// SizeBytes implements marshal.Marshallable.SizeBytes.
//
//go:nosplit
func (u *Uint16) SizeBytes() int {
	return 2
}

// UnmarshalBytes implements marshal.Marshallable.UnmarshalBytes.
func (u *Uint16) UnmarshalBytes(src []byte) error {
	*u = Uint16(uint16(binary.LittleEndian.Uint16(src[:2])))
	return nil
}

// SizeBytes implements marshal.Marshallable.SizeBytes.
//
//go:nosplit
func (u *Uint32) SizeBytes() int {
	return 4
}

// UnmarshalBytes implements marshal.Marshallable.UnmarshalBytes.
func (u *Uint32) UnmarshalBytes(src []byte) error {
	*u = Uint32(uint32(binary.LittleEndian.Uint32(src[:4])))
	return nil
}

// SizeBytes implements marshal.Marshallable.SizeBytes.
//
//go:nosplit
func (u *Uint64) SizeBytes() int {
	return 8
}

// UnmarshalBytes implements marshal.Marshallable.UnmarshalBytes.
func (u *Uint64) UnmarshalBytes(src []byte) error {
	*u = Uint64(uint64(binary.LittleEndian.Uint64(src[:8])))
	return nil
}

// SizeBytes implements marshal.Marshallable.SizeBytes.
//
//go:nosplit
func (u *Uint8) SizeBytes() int {
	return 1
}

// UnmarshalBytes implements marshal.Marshallable.UnmarshalBytes.
func (u *Uint8) UnmarshalBytes(src []byte) error {
	*u = Uint8(uint8(src[0]))
	return nil
}
