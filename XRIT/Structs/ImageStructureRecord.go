package Structs

type ImageStructureRecord struct {
	Type byte
	Size uint16

	BitsPerPixel byte
	Columns      uint16
	Lines        uint16
	Compression  byte
}

func (isr *ImageStructureRecord) GetType() int {
	return int(isr.Type)
}
