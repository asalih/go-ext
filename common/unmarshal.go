package common

import (
	"bytes"

	"github.com/lunixbochs/struc"
	"golang.org/x/xerrors"
)

type Unmarshal interface {
	SizeBytes() int

	// UnmarshalBytes deserializes a type from src.
	// Precondition: buf must be at least SizeBytes() in length.
	UnmarshalBytes(buf []byte) error
}

func UnmarshalBytes(v Unmarshal, buf []byte) error {
	br := bytes.NewBuffer(buf)
	if err := struc.Unpack(br, v); err != nil {
		return xerrors.Errorf("failed to binary read super block: %w", err)
	}

	return nil
}
