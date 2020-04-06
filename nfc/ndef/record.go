package ndef

import (
	"bytes"
	"encoding/binary"
)

type Record struct {
	TypeNameFormat  TypeNameFormat
	idLengthPresent bool
	shortRecord     bool
	chunkFlag       bool
	messageEnd      bool
	messageBegin    bool

	RecordId []byte
	Payload  []byte

	typeLength    uint8
	payloadLength uint32
	idLength      uint8
	recordType    RecordTypeDefinition
}

func (r *Record) PayloadString() string {
	return string(r.Payload)
}

func (r *Record) String() string {
	return string(r.Payload)
}

func (r *Record) Marshal() ([]byte, error) {
	buffer := bytes.Buffer{}

	buffer.WriteByte(r.header())

	buffer.WriteByte(r.typeLength)

	if r.shortRecord {
		buffer.WriteByte(byte(r.payloadLength))
	} else {
		_ = binary.Write(&buffer, binary.BigEndian, r.payloadLength)
	}

	if r.idLengthPresent {
		buffer.WriteByte(r.idLength)
	}

	if r.typeLength > 0 {
		buffer.WriteByte(byte(r.TypeNameFormat))
	}

	if r.idLengthPresent && r.idLength > 0 {
		buffer.Write(r.RecordId)
	}

	buffer.Write(r.Payload)

	return buffer.Bytes(), nil
}

func (r Record) header() byte {
	var header byte = 0x00
	if r.messageBegin {
		header |= 0x80
	}
	if r.messageEnd {
		header |= 0x40
	}
	if r.chunkFlag {
		header |= 0x20
	}
	if len(r.Payload) < 0xFF {
		header |= 0x10
	}
	if len(r.RecordId) > 0 {
		header |= 0x08
	}
	header |= byte(r.TypeNameFormat)

	return header
}
