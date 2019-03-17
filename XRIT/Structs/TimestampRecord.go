package Structs

import (
	"encoding/binary"
	"github.com/opensatelliteproject/SatHelperApp/XRIT/PacketData"
	"time"
)

type TimestampRecord struct {
	Type         byte
	Days         uint16
	Milisseconds uint32
}

func MakeTimestampRecord(data []byte) *TimestampRecord {
	v := TimestampRecord{}

	v.Type = PacketData.TimestampRecord
	v.Days = binary.BigEndian.Uint16(data[1:3])
	v.Milisseconds = binary.BigEndian.Uint32(data[3:7])

	return &v
}

func (tr *TimestampRecord) GetDateTime() time.Time {
	d := time.Date(1958, 1, 1, 0, 0, 0, 0, time.UTC)
	d = d.Add(time.Duration(int64(tr.Days)*3600*24) * time.Second)
	d = d.Add(time.Duration(tr.Milisseconds) * time.Millisecond)

	return d
}

func (tr *TimestampRecord) GetType() int {
	return int(tr.Type)
}
