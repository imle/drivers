package nfc

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"tinygo.org/x/drivers/nfc/driver"
	"tinygo.org/x/drivers/nfc/tag"
)

type Device struct {
	driver driver.PN532Driver
}

func NewDevice(driver driver.PN532Driver) *Device {
	return &Device{
		driver: driver,
	}
}

type FirmwareVersion struct {
	ICVersion        uint8
	Version          uint8
	Revision         uint8
	ISO18092         bool
	ISOIEC14443TypeB bool
	ISOIEC14443TypeA bool
}

func (f FirmwareVersion) String() string {
	return fmt.Sprintf("PN5%X v%d.%d", f.ICVersion, f.Version, f.Revision)
}

func (d *Device) GetFirmwareVersion() (*FirmwareVersion, error) {
	resp, err := d.driver.SendCommand(driver.NewGetFirmwareVersionCommand(), 1*time.Second)
	if err != nil {
		return nil, err
	}

	version := uint32(0)
	err = binary.Read(bytes.NewBuffer(resp.Data), binary.BigEndian, &version)
	if err != nil {
		return nil, err
	}

	fv := &FirmwareVersion{
		ICVersion:        uint8(version >> 24),
		Version:          uint8(version >> 16),
		Revision:         uint8(version >> 8),
		ISO18092:         version&0x04 > 0,
		ISOIEC14443TypeB: version&0x02 > 0,
		ISOIEC14443TypeA: version&0x01 > 0,
	}

	return fv, nil
}

type SAMConfig struct {
	Mode              driver.SAMConfigurationMode
	TimeoutMultiplier uint8 // 50ms * TimeoutMultiplier = timeout
	DoNotUseIRQPin    bool
}

func (d *Device) SAMConfig(c *SAMConfig) error {
	if c == nil {
		c = &SAMConfig{
			Mode:              driver.SAMConfigurationModeNormalMode,
			TimeoutMultiplier: 0x14, // 0x14 * 50ms = 1 second
			DoNotUseIRQPin:    false,
		}
	}

	_, err := d.driver.SendCommand(driver.NewSAMConfigurationCommand(c.Mode, c.TimeoutMultiplier, !c.DoNotUseIRQPin), 1*time.Second)
	if err != nil {
		return err
	}
	return nil
}

func (d *Device) ReadRegister(addr uint16) (uint8, error) {
	resp, err := d.driver.SendCommand(driver.NewReadRegisterCommand(addr), 1*time.Second)
	if err != nil {
		return 0, err
	}

	return resp.Data[0], nil
}

func (d *Device) WriteRegister(addr uint16, value uint8) error {
	_, err := d.driver.SendCommand(driver.NewWriteRegisterCommand(addr, value), 1*time.Second)
	if err != nil {
		return err
	}

	return nil
}

func (d *Device) readPassiveTargetID(baudRate driver.CardBaudRate, timeout time.Duration) ([]byte, error) {
	resp, err := d.driver.SendCommand(driver.NewInListPassiveTargetCommand(baudRate), timeout)
	if err != nil {
		return nil, err
	}

	if resp.Data[0] != 1 {
		return nil, nil
	}

	return resp.Data, nil
}

func (d *Device) GetTagIfPresent(timeout time.Duration) (driver.Card, error) {
	resp, err := d.driver.SendCommand(driver.NewInListPassiveTargetCommand(driver.CardBaudRateMifareIso14443A), timeout)
	if err != nil {
		if err == driver.TimeoutError {
			return nil, nil
		}

		return nil, err
	}

	if resp.Data[0] == 0 {
		return nil, nil
	}

	if resp.Data[0] == 2 {
		return nil, errors.New("library not able to process 2 cards")
	}

	uidLength := resp.Data[5]

	if uidLength == tag.MifareClassicUIDLength {
		uid := [tag.MifareClassicUIDLength]byte{
			resp.Data[6], resp.Data[7], resp.Data[8], resp.Data[9],
		}
		return tag.NewMifareClassic(d.driver, uid, resp.Data[1]), nil
	} else if uidLength == tag.MifareUltralightUIDLength {
		//uid := [tag.MifareUltralightUIDLength]byte{
		//	resp.Data[6], resp.Data[7], resp.Data[8], resp.Data[9], resp.Data[10], resp.Data[11], resp.Data[12],
		//}
		return nil, errors.New("unimplemented card type")
	} else {
		return nil, errors.New("invalid uid length found")
	}
}

func (d *Device) Configure(c *SAMConfig) error {
	if err := d.driver.Configure(); err != nil {
		return err
	}
	if _, err := d.GetFirmwareVersion(); err != nil {
		return err
	}
	if err := d.SAMConfig(c); err != nil {
		return err
	}

	return nil
}
