package Models

type Statistics struct {
	SCID uint8
	VCID uint8
	PacketNumber uint64
	VitErrors uint16
	FrameBits uint16
	RsErrors [4]int32
	SignalQuality uint8
	SyncCorrelation uint8
	PhaseCorrection uint8
	LostPackets uint64
	AverageVitCorrections uint16
	AverageRSCorrections uint8
	DroppedPackets uint64
	ReceivedPacketsPerChannel [256]int64
	LostPacketsPerChannel [256]int64
	TotalPackets uint64
	StartTime uint32
	SyncWord [4]uint8
	FrameLock uint8
	DemodulatorFifoUsage uint8
	DecoderFifoUsage uint8
}
