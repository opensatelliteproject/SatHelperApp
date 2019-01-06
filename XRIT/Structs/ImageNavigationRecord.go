package Structs

type ImageNavigationRecord struct {
	Type byte
	Size uint16

	ProjectionName      string
	ColumnScalingFactor uint32
	LineScalingFactor   uint32
	ColumnOffset        int32
	LineOffset          int32
}

func (imr *ImageNavigationRecord) GetType() int {
	return int(imr.Type)
}
