package ndef

type RecordTypeDefinition string

const (
	RecordTypeDefinitionGeneric     RecordTypeDefinition = ""
	RecordTypeDefinitionURI         RecordTypeDefinition = "U"
	RecordTypeDefinitionText        RecordTypeDefinition = "T"
	RecordTypeDefinitionSmartPoster RecordTypeDefinition = "Sp"
)

func (t RecordTypeDefinition) String() string {
	switch t {
	case RecordTypeDefinitionGeneric:
		return "Generic"
	case RecordTypeDefinitionURI:
		return "URI"
	case RecordTypeDefinitionText:
		return "Text"
	case RecordTypeDefinitionSmartPoster:
		return "Smart Poster"
	default:
		panic("unimplemented Type Name Format encountered")
	}
}
