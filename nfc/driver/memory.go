package driver

import (
	"bytes"
)

type PN532MemoryDriver struct {
	buffer *bytes.Buffer
}

func NewMemoryPN532(buffer *bytes.Buffer) *PN532MemoryDriver {
	return &PN532MemoryDriver{buffer: buffer}
}

func NewMemoryPN532FromBytes(b []byte) *PN532MemoryDriver {
	return NewMemoryPN532(bytes.NewBuffer(b))
}

func (d PN532MemoryDriver) Configure() {
	panic("implement me")
}
