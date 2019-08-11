package XRIT

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/opensatelliteproject/SatHelperApp/XRIT/NOAAProductID"
	"github.com/opensatelliteproject/SatHelperApp/XRIT/PacketData"
	"github.com/opensatelliteproject/SatHelperApp/XRIT/ScannerSubProduct"
	"github.com/opensatelliteproject/SatHelperApp/XRIT/Structs"
	"strconv"
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

	AllHeaders []Structs.BaseRecord `json:"-"`

	// Non XRIT Headers - Auxiliary stuff
	TemperatureLUT        []float32
	ImageDataFunctionHash string
	IDFName               string
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
		fname = xh.DCSFilenameHeader.StringData
	} else if xh.AnnotationHeader != nil {
		fname = xh.AnnotationHeader.StringData
	}

	if xh.SegmentIdentificationHeader != nil && xh.SegmentIdentificationHeader.MaxSegments != 1 {
		fname += fmt.Sprintf("seg%03d", xh.SegmentIdentificationHeader.Sequence)
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
	return &Header{
		AllHeaders: make([]Structs.BaseRecord, 0),
	}
}

func MakeFromJSON(data string) (*Header, error) {
	var xh *Header
	err := json.Unmarshal([]byte(data), xh)
	if err != nil {
		return nil, err
	}
	return xh, nil
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
		xh.computeTemperatureLUT()
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

	xh.AllHeaders = append(xh.AllHeaders, record)
}

func (xh *Header) ToBaseNameString() string {
	if xh.SegmentIdentificationHeader != nil && xh.SegmentIdentificationHeader.COMS1 {
		// COMS-1 Image
		return fmt.Sprintf("%s- %s", xh.ImageNavigationHeader.ProjectionName, xh.IDFName)
	}
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
	}

	return baseName
}

func (xh *Header) ToNameString() string {
	baseName := xh.ToBaseNameString()

	if xh.SegmentIdentificationHeader != nil {
		baseName = fmt.Sprintf("%s (ID: %d Seg: %d/%d)", baseName, xh.SegmentIdentificationHeader.ImageID, xh.SegmentIdentificationHeader.Sequence, xh.SegmentIdentificationHeader.MaxSegments-1)
	}

	return baseName
}

func (xh *Header) IsFalseColorPiece() bool {
	if xh.Product().ID == NOAAProductID.GOES16_ABI || xh.Product().ID == NOAAProductID.GOES17_ABI {
		return xh.SubProduct().ID == 2 || xh.SubProduct().ID == 14
	}

	return false
}

func (xh *Header) ToJSON() string {
	data, _ := json.MarshalIndent(xh, "", "   ")
	return string(data)
}

func (xh *Header) computeTemperatureLUT() {
	if xh.ImageDataFunctionHeader == nil {
		return
	}

	d := map[string]string{}
	lines := strings.Split(xh.ImageDataFunctionHeader.StringData, "\n")
	for _, v := range lines {
		o := strings.Split(v, ":=")
		if len(o) == 2 {
			d[o[0]] = o[1]
		}
		if o[0] == "_NAME" {
			xh.IDFName = o[1]
		}
	}

	lut := make([]float32, 256)

	lastV := float32(0) // Use for filling non contiguous space (if any)
	for i := 0; i < 256; i++ {
		stri := strconv.FormatInt(int64(i), 10)
		if d[stri] != "" {
			v, err := strconv.ParseFloat(d[stri], 32)
			if err == nil {
				lastV = float32(v)
			}
		}
		lut[i] = lastV
	}

	xh.TemperatureLUT = lut

	h := sha256.New()
	_, _ = h.Write([]byte(xh.ImageDataFunctionHeader.StringData))

	xh.ImageDataFunctionHash = fmt.Sprintf("%x", h.Sum(nil))
}

// GetTemperatureLUT returns a LUT for converting pixel value (index) to degrees kelvin.
// Returns nil in case of no ImageDataFunctionHeader available
func (xh *Header) GetTemperatureLUT() []float32 {
	return xh.TemperatureLUT
}
