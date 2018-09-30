package Demuxer

type BaseDemuxer interface {
	Init()
	Start()
	Stop()
	SendData([]byte)
	GetName() string
}
