package driver

import (
	"tinygo.org/x/drivers/nfc/ndef"
)

type CardBaudRate uint8

const (
	CardBaudRateMifareIso14443A CardBaudRate = 0x00
)

type Card interface {
	GetIdentifier() []byte
	WriteNDEFMessages(msg []*ndef.Message) error
	ReadNDEFMessages() ([]*ndef.Message, error)
}
