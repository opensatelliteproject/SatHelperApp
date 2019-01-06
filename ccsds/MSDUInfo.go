package ccsds

import (
	"github.com/OpenSatelliteProject/SatHelperApp/XRIT"
	"time"
)

const MSDUTimeout = 15 * 60 // 15 minutes

type MSDUInfo struct {
	APID             int
	ReceivedTime     time.Time
	FileName         string
	LastPacketNumber int
	Header           *XRIT.Header
}

func MakeMSDUInfo() *MSDUInfo {
	return &MSDUInfo{
		ReceivedTime: time.Now(),
	}
}

func (mi *MSDUInfo) Expired() bool {
	return time.Since(mi.ReceivedTime).Seconds() > MSDUTimeout
}

func (mi *MSDUInfo) Refresh() {
	mi.ReceivedTime = time.Now()
}
