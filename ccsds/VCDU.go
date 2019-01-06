package ccsds

type VCDU struct {
	data []byte
}

func MakeVCDU(data []byte) *VCDU {
	return &VCDU{
		data: data,
	}
}

func (vcdu *VCDU) Version() int {
	return int(vcdu.data[0]&0xC0) >> 6
}

func (vcdu *VCDU) SCID() int {
	return int(vcdu.data[0]&0x3f)<<2 | int(vcdu.data[1]&0xc0)>>6
}

func (vcdu *VCDU) VCID() int {
	return int(vcdu.data[1]) & 0x3f
}

func (vcdu *VCDU) Counter() int {
	return (int(vcdu.data[2]) << 16) | (int(vcdu.data[3]) << 8) | int(vcdu.data[4])
}

func (vcdu *VCDU) Data() []byte {
	return vcdu.data[6:]
}

func (vcdu *VCDU) Replay() bool {
	return (vcdu.data[5] & 0x80) > 0
}
