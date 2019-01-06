package Structs

type RiceCompressionRecord struct {
	Type byte
	Size uint16

	Flags uint16
	Pixel byte
	Line  byte
}

func (rcr *RiceCompressionRecord) GetType() int {
	return int(rcr.Type)
}
