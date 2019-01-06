package Structs

import (
	"encoding/binary"
	"github.com/OpenSatelliteProject/SatHelperApp/XRIT/PacketData"
)

type RiceCompressionRecord struct {
	Type  byte
	Flags uint16
	Pixel byte
	Line  byte
}

func MakeRiceCompressionRecord(data []byte) *RiceCompressionRecord {
	v := RiceCompressionRecord{}

	v.Type = PacketData.RiceCompressionRecord

	v.Flags = binary.BigEndian.Uint16(data[0:2])
	v.Pixel = data[3]
	v.Line = data[4]

	return &v
}

func (rcr *RiceCompressionRecord) GetType() int {
	return int(rcr.Type)
}
