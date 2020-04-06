package tag

import (
	"tinygo.org/x/drivers/nfc/ndef"
)

type FeliCaCmd uint8

const (
	FeliCaReadMaxServiceNum    = 16
	FeliCaReadMaxBlockNum      = 12 // for typical FeliCa card
	FeliCaWriteMaxServiceNum   = 16
	FeliCaWriteMaxBlockNum     = 10 // for typical FeliCa card
	FeliCaReqServiceMaxNodeNum = 32

	FeliCaCmdPolling                FeliCaCmd = 0x00
	FeliCaCmdRequestService         FeliCaCmd = 0x02
	FeliCaCmdRequestResponse        FeliCaCmd = 0x04
	FeliCaCmdReadWithoutEncryption  FeliCaCmd = 0x06
	FeliCaCmdWriteWithoutEncryption FeliCaCmd = 0x08
	FeliCaCmdRequestSystemCode      FeliCaCmd = 0x0C
)

type FeliCa struct {
}

func (c FeliCa) WriteNDEFMessages(msg []*ndef.Message) error {
	panic("implement me")
}

func (c FeliCa) ReadNDEFMessages() ([]*ndef.Message, error) {
	panic("implement me")
}
