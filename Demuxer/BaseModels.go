package Demuxer

type BaseDemuxer interface {
	Init()
	Start()
	Stop()
	ProcessFrame([]byte)
	GetName() string
}
