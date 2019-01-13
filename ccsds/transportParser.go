package ccsds

import (
	"encoding/binary"
	"github.com/OpenSatelliteProject/SatHelperApp/Logger"
)

var skipChannels = []int{63, 2047}

type TransportParser struct {
	id          int
	cb          func(*MSDU)
	frameBuffer []byte
	msduSize    int
	tmpMSDU     *MSDU
}

func MakeTransportParser(channelId int, onMSDU func(*MSDU)) *TransportParser {
	return &TransportParser{
		id:          channelId,
		cb:          onMSDU,
		frameBuffer: make([]byte, 0),
		msduSize:    0,
	}
}

func SkipChannel(channelId int) bool {
	for _, v := range skipChannels {
		if v == channelId {
			return true
		}
	}

	return false
}

func (tp *TransportParser) closeFrame() {
	if tp.tmpMSDU != nil {
		tp.tmpMSDU.finalize()
		if !SkipChannel(tp.id) && !SkipChannel(tp.tmpMSDU.APID) { // Skip fill packets
			msdu := tp.tmpMSDU.Clone()
			tp.cb(msdu)
		}
	}
	tp.frameBuffer = make([]byte, 0)
	tp.tmpMSDU = nil
}

func (tp *TransportParser) parseFrame(data []byte) {
	if len(data) == 0 {
		/* Ignore */
		return
	}
	//SLog.Debug("TransportParser[%d]::parseFrame([%d]byte)", tp.id, len(data))
	if tp.tmpMSDU != nil {
		tp.closeFrame()
	}

	if len(data) > 6 {
		tp.tmpMSDU = MakeMSDUWithHeader(tp.id, data[:6])
		data = data[6:]
	} else {
		tp.frameBuffer = data
	}

	if tp.tmpMSDU != nil {
		data = tp.tmpMSDU.AddBytes(data)
		if tp.tmpMSDU.Closed() {
			tp.closeFrame()
		}
		if len(data) > 0 {
			tp.parseFrame(data)
		}
	}
}

func (tp *TransportParser) WriteChannelData(data *VCDU) {
	if data.VCID() != tp.id {
		SLog.Error("TransportParser: Wrong channel. Expected %d got %d", data.VCID(), tp.id)
		return
	}

	if len(data.Data()) != 886 {
		SLog.Error("TransportParser: Wrong frame size. Expected %d got %d", 886, len(data.data))
		return
	}

	if data.Replay() {
		SLog.Warn("Replay Packet: TODO: IMPLEMENT-ME")
		// TODO
		return
	}

	frame := data.Data()
	fhp := binary.BigEndian.Uint16(frame[:2]) & 0x7FF

	if fhp != 2047 && fhp > uint16(len(frame)) {
		SLog.Error("ERROR: FHP > FRAME")
		return
	}

	frame = frame[2:]

	if fhp != 2047 { // Has packet start
		// There is a header outside the start. Let's split it
		a := frame[:fhp]
		frame = frame[fhp:]

		if tp.tmpMSDU == nil && len(tp.frameBuffer) > 0 { // Not enough bytes to assemble a header last time
			c := append(tp.frameBuffer, a...)
			tp.parseFrame(c)
		} else if tp.tmpMSDU != nil {
			tp.tmpMSDU.AddBytes(a)
			tp.closeFrame()
		}

		tp.parseFrame(frame)
	} else {
		if len(tp.frameBuffer) > 0 && tp.tmpMSDU == nil { // Not enough bytes to assemble a header last time
			c := append(tp.frameBuffer, frame...)
			tp.parseFrame(c)
		} else if tp.tmpMSDU == nil {
			//SLog.Warn("EDGYEDGY CASE")
			//tp.tmpMSDU = MakeMSDUWithHeader(tp.id, frame)
		} else {
			tp.tmpMSDU.AddBytes(frame)
		}
	}

	if tp.tmpMSDU != nil && tp.tmpMSDU.Closed() {
		tp.closeFrame()
	}
}
