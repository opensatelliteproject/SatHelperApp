package PacketData

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
