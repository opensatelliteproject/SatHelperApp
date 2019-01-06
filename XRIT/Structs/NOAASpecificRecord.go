package Structs

import (
	"encoding/binary"
	"github.com/OpenSatelliteProject/SatHelperApp/XRIT/PacketData"
	"github.com/OpenSatelliteProject/SatHelperApp/XRIT/Presets"
)

type NOAASpecificRecord struct {
	Type         byte
	Signature    string
	ProductID    uint16
	ProductSubID uint16
	Parameter    uint16
	Compression  byte
}

func MakeNOAASpecificRecord(data []byte) *NOAASpecificRecord {
	v := NOAASpecificRecord{}

	v.Type = PacketData.NOAASpecificHeader

	v.Signature = string(data[:4])
	v.ProductID = binary.BigEndian.Uint16(data[4:6])
	v.ProductSubID = binary.BigEndian.Uint16(data[6:8])
	v.ProductSubID = binary.BigEndian.Uint16(data[8:10])
	v.Compression = data[11]

	return &v
}

func (nsr *NOAASpecificRecord) Product() PacketData.NOAAProduct {
	return Presets.GetProductById(int(nsr.ProductID))
}

func (nsr *NOAASpecificRecord) GetType() int {
	return int(nsr.Type)
}
