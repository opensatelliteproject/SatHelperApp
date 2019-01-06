package Structs

type SegmentIdentificationRecord struct {
	Type byte
	Size uint16

	ImageID     uint16
	Sequence    uint16
	StartColumn uint16
	StartLine   uint16
	MaxSegments uint16
	MaxColumns  uint16
	MaxRows     uint16
}

func (sir *SegmentIdentificationRecord) GetType() int {
	return int(sir.Type)
}
