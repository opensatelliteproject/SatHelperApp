package Demuxer

type BaseDemuxer interface {
	Init()
	Start()
	Stop()
	SendFrame([]byte)
	GetName() string
}
