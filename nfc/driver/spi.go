package driver

import (
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"machine"
)

const (
	DataWrite  byte = 0x01
	StatusRead byte = 0x02
	DataRead   byte = 0x03
)

var (
	AckTimeoutError   = errors.New("timeout checking ack")
	TimeoutError      = errors.New("timeout waiting for response")
	InvalidFrameError = errors.New("invalid frame")
)

type PN532SPIDriver struct {
	spi machine.SPI
	ss  machine.Pin
}

func NewPN532SPIDriver(spi machine.SPI, ss machine.Pin) *PN532SPIDriver {
	return &PN532SPIDriver{
		spi: spi,
		ss:  ss,
	}
}

type PN532SPIDriverConfig struct {
	SS machine.Pin
	machine.SPIConfig
}

func (d PN532SPIDriver) isReady() (bool, error) {
	d.ss.Low()
	defer d.ss.High()

	if _, err := d.spi.Transfer(StatusRead); err != nil {
		return false, err
	}
	status, err := d.spi.Transfer(0x00)
	if err != nil {
		return false, err
	}

	return status != 0, nil
}

func (d PN532SPIDriver) Configure() error {
	d.ss.Configure(machine.PinConfig{Mode: machine.PinOutput})

	err := d.spi.Configure(machine.SPIConfig{
		LSBFirst: true,
	})
	if err != nil {
		return err
	}

	d.WakeUp()
	return nil
}

func (d PN532SPIDriver) WakeUp() {
	d.ss.High()
	d.ss.Low()

	time.Sleep(400 * time.Millisecond)
	d.ss.High()
	time.Sleep(10 * time.Millisecond)
}

func (d PN532SPIDriver) isReadyTimeout(timeout time.Duration) (ready bool, err error) {
	start := time.Now()

	for ms := 0 * time.Millisecond; !ready; ms = time.Now().Sub(start) {
		if timeout > 0 && ms > timeout {
			return false, nil
		}

		ready, err = d.isReady()
		if err != nil {
			return false, err
		}

		time.Sleep(2 * time.Millisecond)
	}

	return ready, nil
}

func (d PN532SPIDriver) SendCommand(command Command, timeout time.Duration) (*Response, error) {
	// Write command
	err := func() error {
		d.ss.Low()
		defer d.ss.High()
		time.Sleep(2 * time.Millisecond)

		length := byte(len(command.header) + len(command.body))

		packet := make([]byte, 9+length)
		pointer := 0

		packet[pointer] = DataWrite
		pointer++

		packet[pointer] = PN532Preamble
		pointer++
		packet[pointer] = PN532StartCode1
		pointer++
		packet[pointer] = PN532StartCode2
		pointer++
		checksum := PN532Preamble + PN532StartCode1 + PN532StartCode2

		packet[pointer] = length + 1
		pointer++
		packet[pointer] = ^length
		pointer++

		packet[pointer] = PN532HostToPN532
		pointer++
		checksum += PN532HostToPN532

		for i := range command.header {
			packet[pointer] = command.header[i]
			pointer++
			checksum += command.header[i]
		}

		for i := range command.body {
			packet[pointer] = command.body[i]
			pointer++
			checksum += command.body[i]
		}

		packet[pointer] = ^checksum
		pointer++
		packet[pointer] = PN532Postamble
		pointer++

		fmt.Println(hex.EncodeToString(packet))

		return d.spi.Tx(packet, nil)
	}()
	if err != nil {
		return nil, err
	}

	// Wait for ready to ack
	ready, err := d.isReadyTimeout(PN532AckWaitTime)
	if err != nil {
		return nil, err
	}
	if !ready {
		return nil, AckTimeoutError
	}

	ack, err := d.readAckFrame()
	if err != nil {
		return nil, err
	}
	if !ack {
		return nil, errors.New("timeout waiting for ack")
	}

	// Read Response
	return func() (*Response, error) {
		ready, err := d.isReadyTimeout(timeout)
		if err != nil {
			return nil, err
		}
		if !ready {
			return nil, TimeoutError
		}

		d.ss.Low()
		defer d.ss.High()
		time.Sleep(1 * time.Millisecond)

		if _, err = d.spi.Transfer(DataRead); err != nil {
			return nil, err
		}

		header := make([]byte, 7)
		if err = d.spi.Tx(nil, header); err != nil {
			return nil, err
		}

		fmt.Print(hex.EncodeToString(header))

		if header[0] != PN532Preamble {
			return nil, InvalidFrameError
		}
		if header[1] != PN532StartCode1 {
			return nil, InvalidFrameError
		}
		if header[2] != PN532StartCode2 {
			return nil, InvalidFrameError
		}

		length := header[3]
		lengthChecksum := header[4]

		if length+lengthChecksum != 0 {
			return nil, InvalidFrameError
		}

		if header[5] != PN532PN532ToHost {
			return nil, InvalidFrameError
		}
		if header[6] != command.header[0]+1 {
			return nil, InvalidFrameError
		}
		sum := PN532PN532ToHost + header[6]

		data := make([]byte, length)
		if err = d.spi.Tx(nil, data); err != nil {
			return nil, err
		}
		fmt.Print(hex.EncodeToString(data))

		for i := 0; i < len(data); i++ {
			sum += data[i]
		}

		dataChecksum := data[length-1]
		fmt.Printf(" (0x%02X 0x%02X 0x%02X)\n", sum, dataChecksum, sum+dataChecksum)
		if sum+dataChecksum != 0 {
			return nil, InvalidFrameError
		}
		if data[length-1] != PN532Postamble {
			return nil, InvalidFrameError
		}

		return &Response{Data: data[0 : length-2]}, nil
	}()
}

var pn532AckFrame = []byte{0x00, 0x00, 0xFF, 0x00, 0xFF, 0x00}

func (d PN532SPIDriver) readAckFrame() (bool, error) {
	buf := make([]byte, len(pn532AckFrame))

	d.ss.Low()
	defer d.ss.High()
	time.Sleep(1 * time.Millisecond)
	_, err := d.spi.Transfer(DataRead)
	if err != nil {
		return false, err
	}

	for i := 0; i < len(buf); i++ {
		buf[i], err = d.spi.Transfer(0x00)
		if err != nil {
			return false, err
		}
	}

	for i := range buf {
		if buf[i] != pn532AckFrame[i] {
			return false, nil
		}
	}

	return true, nil
}
