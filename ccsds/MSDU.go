package ccsds

import (
	"encoding/binary"
	"github.com/OpenSatelliteProject/SatHelperApp/Logger"
)

type MSDU struct {
	ChannelId        int
	Version          int
	APID             int
	Priority         int
	Type             int
	HasSecondHeader  bool
	PrimaryHeader    []byte
	Sequence         int
	PacketNumber     int
	PacketLength     int
	FullPacketLength int
	Data             []byte

	CRC           uint16
	CalculatedCRC uint16

	lengthValid        bool
	fullData           []byte
	currentFullDataPos int
	headerParsed       bool
	closed             bool
}

func MakeMSDUWithHeader(channelId int, header []byte) *MSDU {
	m := MSDU{
		ChannelId:     channelId,
		PrimaryHeader: make([]byte, 6),
		fullData:      make([]byte, 0, 8192+6), // Max Length is 8192 + 6 bytes header
		headerParsed:  false,
	}
	m.fullData = m.fullData[:len(header)]
	copy(m.fullData, header)
	m.currentFullDataPos += len(header)
	m.parseHeader()

	return &m
}

func (msdu *MSDU) parseHeader() {
	if len(msdu.fullData) < 6 || msdu.headerParsed {
		return
	}

	data := msdu.fullData

	o := binary.BigEndian.Uint16(data[:2])

	msdu.Version = int((o & 0xE000) >> 13)

	msdu.Type = int((o & 0x1000) >> 12)
	msdu.HasSecondHeader = ((o & 0x800) >> 11) > 0
	msdu.APID = int(o & 0x7FF)
	msdu.Priority = msdu.APID / 32

	o = binary.BigEndian.Uint16(data[2:4])

	msdu.Sequence = int((o & 0xC000) >> 14)
	msdu.PacketNumber = int(o & 0x3FFF)
	msdu.PacketLength = SizeFromMSDUHeader(data) + 1
	msdu.FullPacketLength = msdu.PacketLength + 6

	if msdu.FullPacketLength > cap(msdu.fullData) {
		SLog.Warn("Received a packet that is reporting to be bigger than %d (got %d). Skipping parse...", cap(msdu.fullData), msdu.FullPacketLength)
		msdu.FullPacketLength = 0
		msdu.finalize()
		return
	}

	msdu.fullData = msdu.fullData[:msdu.FullPacketLength] // Already 8192 reserved, so we can do that.
	msdu.headerParsed = true
}

func (msdu *MSDU) AddBytes(data []byte) []byte {
	if msdu.headerParsed {
		bytesToAdd := msdu.FullPacketLength - msdu.currentFullDataPos
		if len(data) < bytesToAdd {
			bytesToAdd = len(data)
		}

		a := data[:bytesToAdd]
		data = data[bytesToAdd:]

		copy(msdu.fullData[msdu.currentFullDataPos:], a)

		msdu.currentFullDataPos += bytesToAdd

		if msdu.currentFullDataPos == msdu.FullPacketLength {
			msdu.finalize()
		}

		return data
	}

	remainingBytes := len(msdu.fullData) - msdu.currentFullDataPos
	if remainingBytes < len(data) {
		missingBytes := len(data) - remainingBytes
		msdu.fullData = msdu.fullData[:len(msdu.fullData)+missingBytes] // Already 8192 reserved, so we can do that.
	}

	copy(msdu.fullData[msdu.currentFullDataPos:], data)
	msdu.currentFullDataPos += len(data)

	if msdu.currentFullDataPos > 6 {
		msdu.parseHeader()
	}

	return data[0:0]
}

func (msdu *MSDU) finalize() {
	if !msdu.closed {
		data := msdu.fullData
		data = data[6:]

		msdu.lengthValid = msdu.PacketLength == len(data)
		msdu.CRC = binary.BigEndian.Uint16(data[len(data)-2:])

		msdu.CalculatedCRC = CRC(data[:len(data)-2])

		msdu.Data = data
		msdu.closed = true
	}
}

func (msdu *MSDU) Closed() bool {
	return msdu.closed
}

func (msdu *MSDU) Clone() *MSDU {
	v := *msdu
	return &v
}

func (msdu *MSDU) Valid() bool {
	if !msdu.lengthValid {
		SLog.Debug("Not valid because length. Expected %d got %d", msdu.PacketLength, len(msdu.Data))
		return false
	}

	if msdu.APID == 2047 {
		return true // Don't check CRC for Fill Frames
	}

	if msdu.CRC != msdu.CalculatedCRC {
		SLog.Debug("Not valid because CRC. Expected 0x%04x got 0x%04x", msdu.CRC, msdu.CalculatedCRC)
		return false
	}

	return true
}

func SizeFromMSDUHeader(data []byte) int {
	if len(data) < 6 {
		panic("Not enough data for sizing MSDU!!")
	}

	l := binary.BigEndian.Uint16(data[4:6])
	return int(l)
}

func CRC(data []byte) uint16 {
	lsb := byte(0xFF)
	msb := byte(0xFF)

	for _, v := range data {
		x := v ^ msb
		x ^= x >> 4
		msb = lsb ^ (x >> 3) ^ (x << 4)
		lsb = x ^ (x << 5)
	}

	return uint16(msb)<<8 + uint16(lsb)
}
