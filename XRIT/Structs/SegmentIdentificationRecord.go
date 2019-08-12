package Structs

import (
	"encoding/binary"
	"github.com/opensatelliteproject/SatHelperApp/XRIT/PacketData"
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
	COMS1       bool
}

func MakeSegmentIdentificationRecord(data []byte) *SegmentIdentificationRecord {
	v := SegmentIdentificationRecord{}

	v.Type = PacketData.SegmentIdentificationRecord

	if len(data) == 4 {
		v.StartColumn = 0
		v.Sequence = uint16(data[0])
		v.MaxSegments = uint16(data[1])
		v.StartLine = binary.BigEndian.Uint16(data[2:4])
		v.COMS1 = true
	} else {
		v.ImageID = binary.BigEndian.Uint16(data[0:2])
		v.Sequence = binary.BigEndian.Uint16(data[2:4])
		v.StartColumn = binary.BigEndian.Uint16(data[4:6])
		v.StartLine = binary.BigEndian.Uint16(data[6:8])
		v.MaxSegments = binary.BigEndian.Uint16(data[8:10])
		v.MaxColumns = binary.BigEndian.Uint16(data[10:12])
		v.MaxRows = binary.BigEndian.Uint16(data[12:14])
		v.COMS1 = false
	}

	return &v
}

func (sir *SegmentIdentificationRecord) GetType() int {
	return int(sir.Type)
}
