package Structs

import (
	"encoding/binary"
	"github.com/opensatelliteproject/SatHelperApp/XRIT/PacketData"
)

type ImageStructureRecord struct {
	Type byte

	BitsPerPixel byte
	Columns      uint16
	Lines        uint16
	Compression  byte
}

func MakeImageStructureRecord(data []byte) *ImageStructureRecord {
	v := ImageStructureRecord{}

	v.Type = PacketData.ImageStructureRecord

	v.BitsPerPixel = data[0]
	v.Columns = binary.BigEndian.Uint16(data[1:3])
	v.Lines = binary.BigEndian.Uint16(data[3:5])
	v.Compression = data[5]

	return &v
}

func (isr *ImageStructureRecord) GetType() int {
	return int(isr.Type)
}
