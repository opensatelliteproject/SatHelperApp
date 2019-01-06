package Structs

import (
	"github.com/OpenSatelliteProject/SatHelperApp/XRIT/PacketData"
	"github.com/OpenSatelliteProject/SatHelperApp/XRIT/Presets"
)

type NOAASpecificRecord struct {
	Type byte
	Size uint16

	Signature    string
	ProductID    uint16
	ProductSubID uint16
	Parameter    uint16
	Compression  byte
}

func (nsr *NOAASpecificRecord) Product() PacketData.NOAAProduct {
	return Presets.GetProductById(int(nsr.ProductID))
}

func (nsr *NOAASpecificRecord) GetType() int {
	return int(nsr.Type)
}
