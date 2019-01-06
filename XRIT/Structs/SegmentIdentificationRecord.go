package Structs

import (
	"encoding/binary"
	"github.com/OpenSatelliteProject/SatHelperApp/XRIT/PacketData"
)

type SegmentIdentificationRecord struct {
	Type        byte
	ImageID     uint16
	Sequence    uint16
	StartColumn uint16
	StartLine   uint16
	MaxSegments uint16
	MaxColumns  uint16
	MaxRows     uint16
}

func MakeSegmentIdentificationRecord(data []byte) *SegmentIdentificationRecord {
	v := SegmentIdentificationRecord{}

	v.Type = PacketData.SegmentIdentificationRecord

	v.ImageID = binary.BigEndian.Uint16(data[0:2])
	v.Sequence = binary.BigEndian.Uint16(data[2:4])
	v.StartColumn = binary.BigEndian.Uint16(data[4:6])
	v.StartLine = binary.BigEndian.Uint16(data[6:8])
	v.MaxSegments = binary.BigEndian.Uint16(data[8:10])
	v.MaxColumns = binary.BigEndian.Uint16(data[10:12])
	v.MaxRows = binary.BigEndian.Uint16(data[12:14])

	return &v
}

func (sir *SegmentIdentificationRecord) GetType() int {
	return int(sir.Type)
}
