package Structs

import (
	"encoding/binary"
	"github.com/OpenSatelliteProject/SatHelperApp/XRIT/PacketData"
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
	v.Columns = binary.BigEndian.Uint16(data[1:2])
	v.Lines = binary.BigEndian.Uint16(data[2:3])
	v.Compression = data[4]

	return &v
}

func (isr *ImageStructureRecord) GetType() int {
	return int(isr.Type)
}
