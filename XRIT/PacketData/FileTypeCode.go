package PacketData

import "fmt"

// FileTypeCode
const (
	UNKNOWN = -1

	// By LRIT/HRIT Standard
	// Section 4 of LRIT/HRIT Global Specification, CGMS 03, August 12, 1999
	IMAGE          = 0
	MESSAGES       = 1
	TEXT           = 2
	ENCRYPTION_KEY = 3
	RESERVED4      = 4

	METEOROLOGICAL_DATA = 128

	// NOAA
	DCS   = 130
	EMWIN = 214
)

var FileTypeCode = map[int]string{
	UNKNOWN:             "Unknown",
	IMAGE:               "Image",
	MESSAGES:            "Messages",
	TEXT:                "Text",
	ENCRYPTION_KEY:      "Encryption Key",
	RESERVED4:           "Reserved",
	METEOROLOGICAL_DATA: "Meteorological Data",
	DCS:                 "DCS",
	EMWIN:               "EMWIN",
}

func GetFileTypeCodeString(fileTypeCode int) string {
	v, ok := FileTypeCode[fileTypeCode]
	if ok {
		return v
	}

	return fmt.Sprintf("Unknown (%d)", fileTypeCode)
}
