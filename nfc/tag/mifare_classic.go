package tag

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"time"

	"tinygo.org/x/drivers/nfc/driver"
	"tinygo.org/x/drivers/nfc/ndef"
)

type MifareCmd uint8

const (
	MifareCmdAuthA           MifareCmd = 0x60
	MifareCmdAuthB           MifareCmd = 0x61
	MifareCmdRead            MifareCmd = 0x30
	MifareCmdWrite           MifareCmd = 0xA0
	MifareCmdWriteUltralight MifareCmd = 0xA2
	MifareCmdTransfer        MifareCmd = 0xB0
	MifareCmdDecrement       MifareCmd = 0xC0
	MifareCmdIncrement       MifareCmd = 0xC1
	MifareCmdStore           MifareCmd = 0xC2
)

type MifareTLVBlockType uint8

const (
	MifareTLVBlockTypeNULL        MifareTLVBlockType = 0x00 // These blocks should be ignored
	MifareTLVBlockTypeNDEFMessage MifareTLVBlockType = 0x03 // Block contains an NDEF message
	MifareTLVBlockTypeProprietary MifareTLVBlockType = 0xFD // Block contains proprietary information
	MifareTLVBlockTypeTerminator  MifareTLVBlockType = 0xFE // Last TLV block in the data area
)

type MifareTLVSize uint8

const (
	MifareTLVSizeShort MifareTLVSize = 2
	MifareTLVSizeLong  MifareTLVSize = 4
)

type MifareAuth uint8

const (
	MifareAuthA = MifareAuth(MifareCmdAuthA)
	MifareAuthB = MifareAuth(MifareCmdAuthB)
)

var madKey = [6]byte{0xA0, 0xA1, 0xA2, 0xA3, 0xA4, 0xA5}
var ncfKey = [6]byte{0xD3, 0xF7, 0xD3, 0xF7, 0xD3, 0xF7}

type MifareApplicationDirectory struct {
}

const MifareClassicUIDLength = 4

type MifareClassic struct {
	shield   driver.PN532Driver
	uid      [4]byte
	targetId uint8
}

func NewMifareClassic(shield driver.PN532Driver, uid [4]byte, targetId uint8) *MifareClassic {
	return &MifareClassic{shield: shield, uid: uid, targetId: targetId}
}

func (c MifareClassic) GetIdentifier() []byte {
	return c.uid[:]
}

func (c MifareClassic) WriteNDEFMessages(msg []*ndef.Message) error {
	panic("implement me")
}

const (
	BlockSize            byte = 16
	BlocksPerSmallSector byte = 4
	BlocksPerLargeSector byte = 16
)

var blockBuffer = make([]byte, BlockSize)

func (c MifareClassic) ReadNDEFMessages() (msg []*ndef.Message, errOut error) {
	defer func() {
		if r := recover(); r != nil {
			errOut = r.(error)
		}
	}()

	reader, writer := io.Pipe()

	go func() {
		var block byte = 4

		if c.isFirstBlockInSector(block) {
			err := c.authenticateBlock(MifareAuthB, ncfKey, block)
			if err != nil {
				panic(err)
			}
		}

		data, err := c.readBlock(block)
		if err != nil {
			panic(err)
		}

		_, err = writer.Write(data)
		if err != nil {
			panic(err)
		}

		block++

		if c.isLastBlockInSector(block) {
			block++
		}
	}()

	var messages []*ndef.Message

	singleByte := make([]byte, 1)

	for i := 0; ; {
		_, err := reader.Read(singleByte)
		if err != nil {
			return nil, err
		}
		i += 1

		blockType := MifareTLVBlockType(singleByte[0])

		if blockType == MifareTLVBlockTypeTerminator {
			// Done processing, no more data
			break
		}

		if blockType == MifareTLVBlockTypeNULL {
			// Get length byte (should be 0) and continue to next record
			_, err := reader.Read(singleByte)
			if err != nil {
				return nil, err
			}
			i += 1
			continue
		}

		if blockType != MifareTLVBlockTypeProprietary && blockType != MifareTLVBlockTypeNDEFMessage {
			// If record is not one of the known types then the data is corrupted
			return nil, fmt.Errorf("invalid block type found: 0x%02X", blockType)
		}

		// Get length byte
		_, err = reader.Read(singleByte)
		if err != nil {
			return nil, err
		}
		i += 1

		messageLength := uint16(singleByte[0])
		// check for 3 byte size
		if singleByte[0] == 0xFF {
			err = binary.Read(reader, binary.BigEndian, &messageLength)
			if err != nil {
				return nil, err
			}
			i += 2
		}

		data := make([]byte, messageLength)

		_, err = reader.Read(data)
		if err != nil {
			return nil, err
		}
		i += int(messageLength)

		if blockType == MifareTLVBlockTypeProprietary {
			// Read all the proprietary data and go to next record
			continue
		}

		msg, err := ndef.Decode(bytes.NewBuffer(data))
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func (c MifareClassic) authenticateBlock(auth MifareAuth, key [6]byte, block byte) error {
	data := []byte{
		c.targetId,                                     // Target ID
		byte(auth),                                     // AuthA or AuthB
		key[0], key[1], key[2], key[3], key[4], key[5], // Auth Key
		block,                                  // Block number
		c.uid[0], c.uid[1], c.uid[2], c.uid[3], // UID
	}

	resp, err := c.shield.SendCommand(driver.NewInDataExchangeCommand(data), 1*time.Second)
	if err != nil {
		return err
	}

	if resp.Data[0] != 0x00 {
		return fmt.Errorf("unable to authenticate block: 0x%02X", block)
	}

	return nil
}

func (c MifareClassic) readBlock(block byte) ([]byte, error) {
	data := []byte{
		c.targetId,          // Target ID (relevant if we allowed more than one card to be read
		byte(MifareCmdRead), // Read command
		block,               // Block number
	}

	resp, err := c.shield.SendCommand(driver.NewInDataExchangeCommand(data), 1*time.Second)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}

func (c MifareClassic) isFirstBlockInSector(block uint8) bool {
	if block < 0x80 {
		return block%BlocksPerSmallSector == 0
	} else {
		return block%BlocksPerLargeSector == 0
	}
}

func (c MifareClassic) isLastBlockInSector(block uint8) bool {
	if block < 0x80 {
		return (block+1)%BlocksPerSmallSector == 0
	} else {
		return (block+1)%BlocksPerLargeSector == 0
	}
}
