package ndef

type UriPrefix uint8

// Prefixes for NDEF Records (to identify record type)
const (
	UriPrefixNone        UriPrefix = 0x00
	UriPrefixHttpWwwDot  UriPrefix = 0x01
	UriPrefixHttpsWwwDot UriPrefix = 0x02
	UriPrefixHttp        UriPrefix = 0x03
	UriPrefixHttps       UriPrefix = 0x04
	UriPrefixTel         UriPrefix = 0x05
	UriPrefixMailto      UriPrefix = 0x06
	UriPrefixFtpAnonAt   UriPrefix = 0x07
	UriPrefixFtpFtpDot   UriPrefix = 0x08
	UriPrefixFtps        UriPrefix = 0x09
	UriPrefixSftp        UriPrefix = 0x0A
	UriPrefixSmb         UriPrefix = 0x0B
	UriPrefixNfs         UriPrefix = 0x0C
	UriPrefixFtp         UriPrefix = 0x0D
	UriPrefixDav         UriPrefix = 0x0E
	UriPrefixNews        UriPrefix = 0x0F
	UriPrefixTelnet      UriPrefix = 0x10
	UriPrefixImap        UriPrefix = 0x11
	UriPrefixRtsp        UriPrefix = 0x12
	UriPrefixUrn         UriPrefix = 0x13
	UriPrefixPop         UriPrefix = 0x14
	UriPrefixSip         UriPrefix = 0x15
	UriPrefixSips        UriPrefix = 0x16
	UriPrefixTftp        UriPrefix = 0x17
	UriPrefixBtspp       UriPrefix = 0x18
	UriPrefixBtl2cap     UriPrefix = 0x19
	UriPrefixBtgoep      UriPrefix = 0x1A
	UriPrefixTcpobex     UriPrefix = 0x1B
	UriPrefixIrdaobex    UriPrefix = 0x1C
	UriPrefixFile        UriPrefix = 0x1D
	UriPrefixUrnEpcId    UriPrefix = 0x1E
	UriPrefixUrnEpcTag   UriPrefix = 0x1F
	UriPrefixUrnEpcPat   UriPrefix = 0x20
	UriPrefixUrnEpcRaw   UriPrefix = 0x21
	UriPrefixUrnEpc      UriPrefix = 0x22
	UriPrefixUrnNfc      UriPrefix = 0x23
)

var (
	UriPrefixToString = map[UriPrefix]string{
		UriPrefixNone:        "",
		UriPrefixHttpWwwDot:  "http://www.",
		UriPrefixHttpsWwwDot: "https://www.",
		UriPrefixHttp:        "http://",
		UriPrefixHttps:       "https://",
		UriPrefixTel:         "tel:",
		UriPrefixMailto:      "mailto:",
		UriPrefixFtpAnonAt:   "ftp://anonymous:anonymous@",
		UriPrefixFtpFtpDot:   "ftp://ftp.",
		UriPrefixFtps:        "ftps://",
		UriPrefixSftp:        "sftp://",
		UriPrefixSmb:         "smb://",
		UriPrefixNfs:         "nfs://",
		UriPrefixFtp:         "ftp://",
		UriPrefixDav:         "dav://",
		UriPrefixNews:        "news:",
		UriPrefixTelnet:      "telnet://",
		UriPrefixImap:        "imap:",
		UriPrefixRtsp:        "rtsp://",
		UriPrefixUrn:         "urn:",
		UriPrefixPop:         "pop:",
		UriPrefixSip:         "sip:",
		UriPrefixSips:        "sips:",
		UriPrefixTftp:        "tftp:",
		UriPrefixBtspp:       "btspp://",
		UriPrefixBtl2cap:     "btl2cap://",
		UriPrefixBtgoep:      "btgoep://",
		UriPrefixTcpobex:     "tcpobex://",
		UriPrefixIrdaobex:    "irdaobex://",
		UriPrefixFile:        "file://",
		UriPrefixUrnEpcId:    "urn:epc:id:",
		UriPrefixUrnEpcTag:   "urn:epc:tag:",
		UriPrefixUrnEpcPat:   "urn:epc:pat:",
		UriPrefixUrnEpcRaw:   "urn:epc:raw:",
		UriPrefixUrnEpc:      "urn:epc:",
		UriPrefixUrnNfc:      "urn:nfc:",
	}

	StringToUriPrefix = map[string]UriPrefix{
		"":                           UriPrefixNone,
		"http://www.":                UriPrefixHttpWwwDot,
		"https://www.":               UriPrefixHttpsWwwDot,
		"http://":                    UriPrefixHttp,
		"https://":                   UriPrefixHttps,
		"tel:":                       UriPrefixTel,
		"mailto:":                    UriPrefixMailto,
		"ftp://anonymous:anonymous@": UriPrefixFtpAnonAt,
		"ftp://ftp.":                 UriPrefixFtpFtpDot,
		"ftps://":                    UriPrefixFtps,
		"sftp://":                    UriPrefixSftp,
		"smb://":                     UriPrefixSmb,
		"nfs://":                     UriPrefixNfs,
		"ftp://":                     UriPrefixFtp,
		"dav://":                     UriPrefixDav,
		"news:":                      UriPrefixNews,
		"telnet://":                  UriPrefixTelnet,
		"imap:":                      UriPrefixImap,
		"rtsp://":                    UriPrefixRtsp,
		"urn:":                       UriPrefixUrn,
		"pop:":                       UriPrefixPop,
		"sip:":                       UriPrefixSip,
		"sips:":                      UriPrefixSips,
		"tftp:":                      UriPrefixTftp,
		"btspp://":                   UriPrefixBtspp,
		"btl2cap://":                 UriPrefixBtl2cap,
		"btgoep://":                  UriPrefixBtgoep,
		"tcpobex://":                 UriPrefixTcpobex,
		"irdaobex://":                UriPrefixIrdaobex,
		"file://":                    UriPrefixFile,
		"urn:epc:id:":                UriPrefixUrnEpcId,
		"urn:epc:tag:":               UriPrefixUrnEpcTag,
		"urn:epc:pat:":               UriPrefixUrnEpcPat,
		"urn:epc:raw:":               UriPrefixUrnEpcRaw,
		"urn:epc:":                   UriPrefixUrnEpc,
		"urn:nfc:":                   UriPrefixUrnNfc,
	}
)

func (n UriPrefix) Prefix() string {
	return UriPrefixToString[n]
}

func LookupPrefix(s string) UriPrefix {
	return StringToUriPrefix[s]
}
