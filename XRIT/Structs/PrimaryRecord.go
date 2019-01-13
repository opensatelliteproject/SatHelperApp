package Structs

import (
	"encoding/binary"
	"github.com/OpenSatelliteProject/SatHelperApp/XRIT/PacketData"
)

type PrimaryRecord struct {
	Type         byte
	Size         uint16
	FileTypeCode byte
	HeaderLength uint32
	DataLength   uint64
}

func MakePrimaryRecord(data []byte) *PrimaryRecord {
	v := PrimaryRecord{}

	v.Type = PacketData.PrimaryHeader

	v.FileTypeCode = data[0]
	v.HeaderLength = binary.BigEndian.Uint32(data[1:5])
	v.DataLength = binary.BigEndian.Uint64(data[5:13])

	return &v
}

func (pr *PrimaryRecord) GetType() int {
	return int(pr.Type)
}
