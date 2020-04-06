package driver

import (
	"time"
)

const _debug = false

const (
	PN532Preamble   byte = 0x00
	PN532StartCode1 byte = 0x00
	PN532StartCode2 byte = 0xFF
	PN532Postamble  byte = 0x00

	PN532HostToPN532 byte = 0xD4
	PN532PN532ToHost byte = 0xD5

	PN532AckWaitTime = 10 * time.Millisecond // ms, timeout of waiting for ACK
)

type PN532Response uint8

const (
	ResponseInDataExchange      PN532Response = 0x41
	ResponseInListPassiveTarget PN532Response = 0x4B

	GpioValidationBit = 0x80
	GpioP30           = 0
	GpioP31           = 1
	GpioP32           = 2
	GpioP33           = 3
	GpioP34           = 4
	GpioP35           = 5
)

type Response struct {
	Data []byte
}

type PN532Driver interface {
	Configure() error
	WakeUp()
	SendCommand(command Command, timeout time.Duration) (*Response, error)
}
