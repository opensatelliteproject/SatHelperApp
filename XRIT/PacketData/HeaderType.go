package PacketData

import "fmt"

// HeaderType
const (
	Unknown                 = -1
	PrimaryHeader           = 0
	ImageStructureRecord    = 1
	ImageNavigationRecord   = 2
	ImageDataFunctionRecord = 3
	AnnotationRecord        = 4
	TimestampRecord         = 5
	AncillaryTextRecord     = 6
	KeyRecord               = 7
	Head9                   = 9 // Weird

	SegmentIdentificationRecord = 128
	NOAASpecificHeader          = 129
	HeaderStructuredRecord      = 130
	RiceCompressionRecord       = 131
	DCSFileNameRecord           = 132
)

var HeaderType = map[int]string{
	Unknown:                 "Unknown",
	PrimaryHeader:           "Primary Header",
	ImageStructureRecord:    "Image Structure Record",
	ImageNavigationRecord:   "Image Navigation Record",
	ImageDataFunctionRecord: "Image Data Function Record",
	AnnotationRecord:        "Annotation Record",
	TimestampRecord:         "Timestamp Record",
	AncillaryTextRecord:     "Ancillary Text Record",
	KeyRecord:               "Key Record",
	Head9:                   "Unknown",
	SegmentIdentificationRecord: "Segment Identification Record",
	NOAASpecificHeader:          "NOAA Specific Header",
	HeaderStructuredRecord:      "Header Structured Record",
	RiceCompressionRecord:       "Rice Compression Record",
	DCSFileNameRecord:           "DCS StringData Record",
}

func GetHeaderTypeString(headerType int) string {
	v, ok := HeaderType[headerType]
	if ok {
		return v
	}

	return fmt.Sprintf("Unknown (%d)", headerType)
}
