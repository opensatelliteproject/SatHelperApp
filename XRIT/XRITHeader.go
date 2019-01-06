package XRIT

import (
	"fmt"
	"github.com/OpenSatelliteProject/SatHelperApp/XRIT/NOAAProductID"
	"github.com/OpenSatelliteProject/SatHelperApp/XRIT/PacketData"
	"github.com/OpenSatelliteProject/SatHelperApp/XRIT/ScannerSubProduct"
	"github.com/OpenSatelliteProject/SatHelperApp/XRIT/Structs"
	"strings"
)

type Header struct {
	AncillaryHeader             *Structs.AncillaryText
	AnnotationHeader            *Structs.AnnotationRecord
	DCSFilenameHeader           *Structs.DCSFilenameRecord
	HeaderStructuredHeader      *Structs.HeaderStructuredRecord
	ImageDataFunctionHeader     *Structs.ImageDataFunctionRecord
	ImageNavigationHeader       *Structs.ImageNavigationRecord
	ImageStructureHeader        *Structs.ImageStructureRecord
	NOAASpecificHeader          *Structs.NOAASpecificRecord
	PrimaryHeader               *Structs.PrimaryRecord
	RiceCompressionHeader       *Structs.RiceCompressionRecord
	SegmentIdentificationHeader *Structs.SegmentIdentificationRecord
	TimestampHeader             *Structs.TimestampRecord
	UnknownHeaders              []*Structs.UnknownHeader
}

func (xh *Header) Product() *PacketData.NOAAProduct {
	if xh.NOAASpecificHeader != nil {
		v := xh.NOAASpecificHeader.Product()
		return &v
	}

	v := PacketData.MakeNOAAProduct(-1)

	return &v
}

func (xh *Header) SubProduct() *PacketData.NOAASubProduct {
	prod := xh.Product()

	if xh.NOAASpecificHeader != nil {
		v := prod.GetSubProduct(int(xh.NOAASpecificHeader.ProductSubID))
		return &v
	}

	v := PacketData.MakeSubProduct(-1, "Unknown")
	return &v
}

func (xh *Header) IsFullDisk() bool {
	return xh.PrimaryHeader.FileTypeCode == PacketData.IMAGE &&
		(xh.Product().ID == NOAAProductID.GOES13_ABI ||
			xh.Product().ID == NOAAProductID.GOES15_ABI ||
			xh.Product().ID == NOAAProductID.GOES16_ABI ||
			xh.Product().ID == NOAAProductID.GOES17_ABI ||
			xh.Product().ID == NOAAProductID.HIMAWARI8_ABI) &&

		(xh.SubProduct().ID == ScannerSubProduct.INFRARED_FULLDISK ||
			xh.SubProduct().ID == ScannerSubProduct.VISIBLE_FULLDISK ||
			xh.SubProduct().ID == ScannerSubProduct.WATERVAPOUR_FULLDISK)
}

func (xh *Header) Filename() string {
	fname := ""
	if xh.DCSFilenameHeader != nil {
		fname = xh.DCSFilenameHeader.Filename
	} else if xh.AnnotationHeader != nil {
		fname = xh.AnnotationHeader.Filename
	}

	if fname != "" {
		if len(strings.Split(fname, ".")) == 1 {
			fname = fname + ".lrit"
		}
	}

	return fname
}

func (xh *Header) Compression() int {
	if xh.NOAASpecificHeader != nil {
		return int(xh.NOAASpecificHeader.Compression)
	} else if xh.ImageStructureHeader != nil {
		return int(xh.ImageStructureHeader.Compression)
	}

	return PacketData.NO_COMPRESSION
}

func (xh *Header) IsCompressed() bool {
	return xh.Compression() != PacketData.NO_COMPRESSION
}

func MakeXRITHeader() *Header {
	return &Header{}
}

func MakeXRITHeaderWithHeaders(records []Structs.BaseRecord) *Header {
	xh := MakeXRITHeader()

	for _, v := range records {
		xh.SetHeader(v)
	}

	return xh
}

func (xh *Header) SetHeader(record Structs.BaseRecord) {
	switch record.GetType() {
	case PacketData.AncillaryTextRecord:
		xh.AncillaryHeader = record.(*Structs.AncillaryText)
	case PacketData.AnnotationRecord:
		xh.AnnotationHeader = record.(*Structs.AnnotationRecord)
	case PacketData.DCSFileNameRecord:
		xh.DCSFilenameHeader = record.(*Structs.DCSFilenameRecord)
	case PacketData.HeaderStructuredRecord:
		xh.HeaderStructuredHeader = record.(*Structs.HeaderStructuredRecord)
	case PacketData.ImageDataFunctionRecord:
		xh.ImageDataFunctionHeader = record.(*Structs.ImageDataFunctionRecord)
	case PacketData.ImageNavigationRecord:
		xh.ImageNavigationHeader = record.(*Structs.ImageNavigationRecord)
	case PacketData.ImageStructureRecord:
		xh.ImageStructureHeader = record.(*Structs.ImageStructureRecord)
	case PacketData.NOAASpecificHeader:
		xh.NOAASpecificHeader = record.(*Structs.NOAASpecificRecord)
	case PacketData.PrimaryHeader:
		xh.PrimaryHeader = record.(*Structs.PrimaryRecord)
	case PacketData.RiceCompressionRecord:
		xh.RiceCompressionHeader = record.(*Structs.RiceCompressionRecord)
	case PacketData.SegmentIdentificationRecord:
		xh.SegmentIdentificationHeader = record.(*Structs.SegmentIdentificationRecord)
	case PacketData.TimestampRecord:
		xh.TimestampHeader = record.(*Structs.TimestampRecord)
	default:
		xh.UnknownHeaders = append(xh.UnknownHeaders, record.(*Structs.UnknownHeader))
	}
}

func (xh *Header) ToNameString() string {
	baseName := xh.Product().Name

	if xh.SubProduct().Name != "Unknown" {
		baseName = xh.Product().Name + " - " + xh.SubProduct().Name
	}

	if xh.Product().ID == NOAAProductID.GOES13_ABI ||
		xh.Product().ID == NOAAProductID.GOES15_ABI ||
		xh.Product().ID == NOAAProductID.GOES16_ABI ||
		xh.Product().ID == NOAAProductID.GOES17_ABI ||
		xh.Product().ID == NOAAProductID.HIMAWARI8_ABI {
		baseName = xh.Product().Name + " - " + xh.SubProduct().Name

		if xh.SegmentIdentificationHeader != nil {
			baseName = fmt.Sprintf("%s (ID: %s Seg: %d/%d)", baseName, xh.SegmentIdentificationHeader.ImageID, xh.SegmentIdentificationHeader.Sequence, xh.SegmentIdentificationHeader.MaxSegments)
		}
	}

	return baseName
}
