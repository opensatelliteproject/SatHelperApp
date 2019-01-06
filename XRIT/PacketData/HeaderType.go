package PacketData

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
