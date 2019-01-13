package Structs

import (
	"strconv"
	"time"
)

const DCSFrameMark = "\x02\x02\x18"

type DCSPacket struct {
	Address         string
	DateTime        time.Time
	Status          string
	Signal          int
	FrequencyOffset int
	ModIndexNormal  string
	DataQualNominal string
	Channel         string
	SourceCode      string
}

func DCSDateToTime(date string) time.Time {
	// %y%j%H%M%S
	year := "20" + date[:2]
	dayOfYear, _ := strconv.ParseInt(date[2:5], 10, 32)
	hour := date[5:7]
	minute := date[7:9]
	second := date[9:11]
	dt, _ := time.Parse("2006150405", year+hour+minute+second)
	return dt.Add(time.Duration(dayOfYear*3600*24) * time.Second)
}

func MakeDCSPacket(data []byte) *DCSPacket {
	signal, _ := strconv.ParseInt(string(data[21:23]), 10, 32)
	freqOffset, _ := strconv.ParseInt(string(data[23:25]), 10, 32)

	return &DCSPacket{
		Address:         string(data[:9]),
		DateTime:        DCSDateToTime(string(data[9:20])),
		Status:          string(data[20]),
		Signal:          int(signal),
		FrequencyOffset: int(freqOffset),
		ModIndexNormal:  string(data[25]),
		DataQualNominal: string(data[26]),
		Channel:         string(data[27:31]),
		SourceCode:      string(data[31:33]),
	}
}

type DCSData struct {
	Header  string
	Packets []*DCSPacket
}

func ParseDCS(data []byte) *DCSData {
	baseHeader := string(data[:64])
	content := data[64:]

	packets := make([]*DCSPacket, 0)
	lastStart := 0

	for i := 0; i < len(content)-len(DCSFrameMark); i++ {
		if string(content[i:i+3]) == DCSFrameMark {
			if i-1 > 0 {
				packetData := content[lastStart : i-1]
				packets = append(packets, MakeDCSPacket(packetData))
			}
			i += 3
			lastStart = i
		}
	}

	return &DCSData{
		Header:  baseHeader,
		Packets: packets,
	}
}
