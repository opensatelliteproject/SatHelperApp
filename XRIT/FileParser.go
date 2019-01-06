package XRIT

import (
	"encoding/binary"
	"fmt"
	"github.com/OpenSatelliteProject/SatHelperApp/XRIT/PacketData"
	"github.com/OpenSatelliteProject/SatHelperApp/XRIT/Structs"
	"io"
	"os"
)

func ParseFile(filename string) (*Header, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	headerType, _, data, err := readHeader(f)

	if headerType != PacketData.PrimaryHeader {
		return nil, fmt.Errorf("first header is not primary (type == %d)", PacketData.PrimaryHeader)
	}

	primaryHeader := Structs.MakePrimaryRecord(data)

	pos, _ := f.Seek(0, io.SeekCurrent)
	max := int64(primaryHeader.HeaderLength)

	xh := MakeXRITHeader()
	xh.SetHeader(primaryHeader)

	for pos < max {
		headerType, _, data, err = readHeader(f)
		if headerType != PacketData.PrimaryHeader {
			return xh, err
		}

		var header Structs.BaseRecord

		switch headerType {
		case PacketData.AncillaryTextRecord:
			header = Structs.MakeAncillaryText(data)
		case PacketData.AnnotationRecord:
			header = Structs.MakeAnnotationRecord(data)
		case PacketData.DCSFileNameRecord:
			header = Structs.MakeDCSFilenameRecord(data)
		case PacketData.HeaderStructuredRecord:
			header = Structs.MakeHeaderStructuredRecord(data)
		case PacketData.ImageDataFunctionRecord:
			header = Structs.MakeImageDataFunctionRecord(data)
		case PacketData.ImageNavigationRecord:
			header = Structs.MakeImageNavigationRecord(data)
		case PacketData.ImageStructureRecord:
			header = Structs.MakeImageStructureRecord(data)
		case PacketData.NOAASpecificHeader:
			header = Structs.MakeNOAASpecificRecord(data)
		case PacketData.RiceCompressionRecord:
			header = Structs.MakeRiceCompressionRecord(data)
		case PacketData.SegmentIdentificationRecord:
			header = Structs.MakeSegmentIdentificationRecord(data)
		case PacketData.TimestampRecord:
			header = Structs.MakeTimestampRecord(data)
		default:
			header = Structs.MakeUnknownHeader(headerType, data)
		}

		xh.SetHeader(header)

		pos, _ = f.Seek(0, io.SeekCurrent)
	}

	return xh, nil
}

// readHeader reads a single header from XRIT File
func readHeader(reader io.Reader) (headerType byte, size uint16, data []byte, err error) {
	head := make([]byte, 3)
	n, err := reader.Read(head)

	if err != nil {
		return 0, 0, nil, err
	}

	if n != 3 {
		return 0, 0, nil, fmt.Errorf("EOF")
	}

	headerType = head[0]
	size = binary.BigEndian.Uint16(head[1:])

	data = make([]byte, size)
	n, err = reader.Read(data)

	if err != nil {
		return 0, 0, nil, err
	}

	if n != int(size) {
		return 0, 0, nil, fmt.Errorf("EOF")
	}

	return headerType, size, data, nil
}
