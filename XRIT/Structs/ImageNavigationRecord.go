package Structs

import (
	"encoding/binary"
	"github.com/OpenSatelliteProject/SatHelperApp/XRIT/PacketData"
)

type ImageNavigationRecord struct {
	Type byte

	ProjectionName      string
	ColumnScalingFactor uint32
	LineScalingFactor   uint32
	ColumnOffset        int32
	LineOffset          int32
}

func MakeImageNavigationRecord(data []byte) *ImageNavigationRecord {
	inr := ImageNavigationRecord{}

	inr.Type = PacketData.ImageNavigationRecord

	inr.ProjectionName = string(data[:32])
	inr.ColumnScalingFactor = binary.BigEndian.Uint32(data[32:36])
	inr.LineScalingFactor = binary.BigEndian.Uint32(data[36:40])
	inr.ColumnOffset = int32(binary.BigEndian.Uint32(data[40:44]))
	inr.LineOffset = int32(binary.BigEndian.Uint32(data[44:48]))

	return &inr
}

func (imr *ImageNavigationRecord) GetType() int {
	return int(imr.Type)
}
