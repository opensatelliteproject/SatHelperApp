package Structs

type PrimaryRecord struct {
	Type         byte
	Size         uint16
	FileTypeCode byte
	HeaderLength uint32
	DataLength   uint64
}

func (pr *PrimaryRecord) GetType() int {
	return int(pr.Type)
}
