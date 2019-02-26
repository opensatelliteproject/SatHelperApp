package XRIT

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/opensatelliteproject/SatHelperApp/XRIT/PacketData"
	"github.com/opensatelliteproject/SatHelperApp/XRIT/Structs"
	"io"
	"os"
)

type fpReader interface {
	Seek(offset int64, whence int) (int64, error)
	Read(p []byte) (n int, err error)
}

func MemoryParseFile(data []byte) (*Header, error) {
	breader := bytes.NewReader(data)
	return parseFileFromReader(breader)
}

func ParseFile(filename string) (*Header, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	return parseFileFromReader(f)
}

func parseFileFromReader(reader fpReader) (*Header, error) {

	headerType, _, data, err := readHeader(reader)

	if err != nil {
		return nil, err
	}

	if headerType != PacketData.PrimaryHeader {
		return nil, fmt.Errorf("first header is not primary (type == %d)", PacketData.PrimaryHeader)
	}

	primaryHeader := Structs.MakePrimaryRecord(data)

	pos, _ := reader.Seek(0, io.SeekCurrent)
	max := int64(primaryHeader.HeaderLength)

	xh := MakeXRITHeader()
	xh.SetHeader(primaryHeader)

	for pos < max {
		headerType, _, data, err = readHeader(reader)
		if err != nil {
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
		case PacketData.PrimaryHeader:
			header = Structs.MakePrimaryRecord(data)
		default:
			header = Structs.MakeUnknownHeader(headerType, data)
		}

		xh.SetHeader(header)

		pos, _ = reader.Seek(0, io.SeekCurrent)
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
	size = binary.BigEndian.Uint16(head[1:]) - 3 // Already read 3 bytes

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
