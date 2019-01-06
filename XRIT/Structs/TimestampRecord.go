package Structs

import "time"

type TimestampRecord struct {
	Type byte
	Size uint16

	Days         uint16
	Milisseconds uint32
}

func (tr *TimestampRecord) GetDateTime() time.Time {
	d := time.Date(1958, 1, 1, 0, 0, 0, 0, time.UTC)
	d.Add(time.Duration(tr.Days*3600*24) * time.Second)
	d.Add(time.Duration(tr.Milisseconds) * time.Millisecond)

	return d
}

func (tr *TimestampRecord) GetType() int {
	return int(tr.Type)
}
