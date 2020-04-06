package driver

type PN532Command uint8

const (
	Diagnose              PN532Command = 0x00
	GetFirmwareVersion    PN532Command = 0x02
	GetGeneralStatus      PN532Command = 0x04
	ReadRegister          PN532Command = 0x06
	WriteRegister         PN532Command = 0x08
	ReadGpio              PN532Command = 0x0C
	WriteGpio             PN532Command = 0x0E
	SetSerialBaudRate     PN532Command = 0x10
	SetParameters         PN532Command = 0x12
	SAMConfiguration      PN532Command = 0x14
	PowerDown             PN532Command = 0x16
	RFConfiguration       PN532Command = 0x32
	RFRegulationTest      PN532Command = 0x58
	InJumpForDep          PN532Command = 0x56
	InJumpForPSL          PN532Command = 0x46
	InListPassiveTarget   PN532Command = 0x4A
	InAtr                 PN532Command = 0x50
	InPSL                 PN532Command = 0x4E
	InDataExchange        PN532Command = 0x40
	InCommunicateThru     PN532Command = 0x42
	InDeselect            PN532Command = 0x44
	InRelease             PN532Command = 0x52
	InSelect              PN532Command = 0x54
	InAutoPoll            PN532Command = 0x60
	TgInitAsTarget        PN532Command = 0x8C
	TgSetGeneralBytes     PN532Command = 0x92
	TgGetData             PN532Command = 0x86
	TgSetData             PN532Command = 0x8E
	TgSetMetaData         PN532Command = 0x94
	TgGetInitiatorCommand PN532Command = 0x88
	TgResponseToInitiator PN532Command = 0x90
	TgGetTargetStatus     PN532Command = 0x8A
)

type SAMConfigurationMode uint8

const (
	SAMConfigurationModeNormalMode  SAMConfigurationMode = 0x01 // SAM is not used; this is the default mode.
	SAMConfigurationModeVirtualCard SAMConfigurationMode = 0x02 // couple PN532+SAM is seen as only one contactless SAM card from the external world.
	SAMConfigurationModeWiredCard   SAMConfigurationMode = 0x03 // host controller can access to the SAM with standard PCD commands (InListPassiveTarget. InDataExchange, …),
	SAMConfigurationModeDualCard    SAMConfigurationMode = 0x04 // both the PN532 and the SAM are visible from the external world as two separated targets.
)

const maxCardsAtOnce byte = 1

type Command struct {
	header []byte
	body   []byte
}

func NewDiagnoseCommand() Command {
	return Command{
		header: []byte{byte(Diagnose)},
	}
}

func NewGetFirmwareVersionCommand() Command {
	return Command{
		header: []byte{byte(GetFirmwareVersion)},
	}
}

func NewGetGeneralStatusCommand() Command {
	return Command{
		header: []byte{byte(GetGeneralStatus)},
	}
}

func NewReadRegisterCommand(addr uint16) Command {
	return Command{
		header: []byte{byte(ReadRegister), byte(addr >> 8), byte(addr)},
	}
}

func NewWriteRegisterCommand(addr uint16, value uint8) Command {
	return Command{
		header: []byte{byte(WriteRegister), byte(addr >> 8), byte(addr), value},
	}
}

func NewReadGpioCommand() Command {
	return Command{
		header: []byte{byte(ReadGpio)},
	}
}

func NewWriteGpioCommand() Command {
	return Command{
		header: []byte{byte(WriteGpio)},
	}
}

func NewSetSerialBaudRateCommand() Command {
	return Command{
		header: []byte{byte(SetSerialBaudRate)},
	}
}

func NewSetParametersCommand() Command {
	return Command{
		header: []byte{byte(SetParameters)},
	}
}

func NewSAMConfigurationCommand(mode SAMConfigurationMode, timeoutMultiplier uint8, useIRQ bool) Command {
	var irq byte = 0x00
	if useIRQ {
		irq = 0x01
	}
	return Command{
		header: []byte{byte(SAMConfiguration), byte(mode), timeoutMultiplier, irq},
	}
}

func NewPowerDownCommand() Command {
	return Command{
		header: []byte{byte(PowerDown)},
	}
}

func NewRFConfigurationCommand() Command {
	return Command{
		header: []byte{byte(RFConfiguration)},
	}
}

func NewRFRegulationTestCommand() Command {
	return Command{
		header: []byte{byte(RFRegulationTest)},
	}
}

func NewInJumpForDepCommand() Command {
	return Command{
		header: []byte{byte(InJumpForDep)},
	}
}

func NewInJumpForPSLCommand() Command {
	return Command{
		header: []byte{byte(InJumpForPSL)},
	}
}

func NewInListPassiveTargetCommand(baudRate CardBaudRate) Command {
	return Command{
		header: []byte{byte(InListPassiveTarget), maxCardsAtOnce, byte(baudRate)},
	}
}

func NewInAtrCommand() Command {
	return Command{
		header: []byte{byte(InAtr)},
	}
}

func NewInPSLCommand() Command {
	return Command{
		header: []byte{byte(InPSL)},
	}
}

func NewInDataExchangeCommand(data []byte) Command {
	return Command{
		header: []byte{byte(InDataExchange)},
		body:   data,
	}
}

func NewInCommunicateThruCommand() Command {
	return Command{
		header: []byte{byte(InCommunicateThru)},
	}
}

func NewInDeselectCommand() Command {
	return Command{
		header: []byte{byte(InDeselect)},
	}
}

func NewInReleaseCommand() Command {
	return Command{
		header: []byte{byte(InRelease)},
	}
}

func NewInSelectCommand() Command {
	return Command{
		header: []byte{byte(InSelect)},
	}
}

func NewInAutoPollCommand() Command {
	return Command{
		header: []byte{byte(InAutoPoll)},
	}
}

func NewTgInitAsTargetCommand() Command {
	return Command{
		header: []byte{byte(TgInitAsTarget)},
	}
}

func NewTgSetGeneralBytesCommand() Command {
	return Command{
		header: []byte{byte(TgSetGeneralBytes)},
	}
}

func NewTgGetDataCommand() Command {
	return Command{
		header: []byte{byte(TgGetData)},
	}
}

func NewTgSetDataCommand() Command {
	return Command{
		header: []byte{byte(TgSetData)},
	}
}

func NewTgSetMetaDataCommand() Command {
	return Command{
		header: []byte{byte(TgSetMetaData)},
	}
}

func NewTgGetInitiatorCommandCommand() Command {
	return Command{
		header: []byte{byte(TgGetInitiatorCommand)},
	}
}

func NewTgResponseToInitiatorCommand() Command {
	return Command{
		header: []byte{byte(TgResponseToInitiator)},
	}
}

func NewTgGetTargetStatusCommand() Command {
	return Command{
		header: []byte{byte(TgGetTargetStatus)},
	}
}
