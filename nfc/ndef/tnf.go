package ndef

type TypeNameFormat byte

const (
	TypeNameFormatEmpty       TypeNameFormat = 0
	TypeNameFormatWellKnown   TypeNameFormat = 1
	TypeNameFormatMIMEMedia   TypeNameFormat = 2
	TypeNameFormatAbsoluteURI TypeNameFormat = 3
	TypeNameFormatExternal    TypeNameFormat = 4
	TypeNameFormatUnknown     TypeNameFormat = 5
	TypeNameFormatUnchanged   TypeNameFormat = 6
)

func (t TypeNameFormat) String() string {
	switch t {
	case TypeNameFormatEmpty:
		return "Empty"
	case TypeNameFormatWellKnown:
		return "Well-Known"
	case TypeNameFormatMIMEMedia:
		return "MIME Media"
	case TypeNameFormatAbsoluteURI:
		return "Absolute URI"
	case TypeNameFormatExternal:
		return "External"
	case TypeNameFormatUnknown:
		return "Unknown"
	case TypeNameFormatUnchanged:
		return "Unchanged"
	default:
		panic("unimplemented Type Name Format encountered")
	}
}
