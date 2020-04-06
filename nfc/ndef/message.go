package ndef

import (
	"bytes"
	"encoding/binary"
	"io"
)

type Message struct {
	Records []*Record
}

func (m *Message) Marshal() ([]byte, error) {
	buf := &bytes.Buffer{}

	_, err := m.WriteTo(buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (m *Message) WriteTo(w io.Writer) (n int64, err error) {
	panic("implement me")
}

type Reader interface {
	io.Reader
	io.ByteReader
}

func MustDecode(reader Reader) *Message {
	decode, err := Decode(reader)
	if err != nil {
		panic(err)
	}
	return decode
}

func Decode(reader Reader) (*Message, error) {
	msg := &Message{
		Records: []*Record{},
	}

	recordIndex := 0
	for {
		header, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}

		msg.Records = append(msg.Records, &Record{
			TypeNameFormat:  TypeNameFormat(header & 0x07),
			idLengthPresent: header&0x08 == 0x08,
			shortRecord:     header&0x10 == 0x10,
			chunkFlag:       header&0x20 == 0x20,
			messageEnd:      header&0x40 == 0x40,
			messageBegin:    header&0x80 == 0x80,
		})

		record := msg.Records[recordIndex]

		record.typeLength, _ = reader.ReadByte()

		if record.shortRecord {
			length, err := reader.ReadByte()
			if err != nil {
				return nil, err
			}

			record.payloadLength = uint32(length)
		} else {
			err = binary.Read(reader, binary.BigEndian, &record.payloadLength)
			if err != nil {
				return nil, err
			}
		}

		if record.idLengthPresent {
			record.idLength, _ = reader.ReadByte()
		}

		if record.typeLength > 0 {
			recordType := make([]byte, record.typeLength)
			_, err = reader.Read(recordType)
			if err != nil {
				return nil, err
			}
			record.recordType = RecordTypeDefinition(recordType)
		}

		if record.idLength > 0 {
			record.RecordId = make([]byte, record.idLength)
			_, err = reader.Read(record.RecordId)
			if err != nil {
				return nil, err
			}
		}

		if record.idLength > 0 {
			record.RecordId = make([]byte, record.idLength)
			_, err = reader.Read(record.RecordId)
			if err != nil {
				return nil, err
			}
		}

		record.Payload = make([]byte, record.payloadLength)
		_, err = reader.Read(record.Payload)
		if err != nil {
			return nil, err
		}

		if record.messageEnd {
			break
		}
	}

	return msg, nil
}
