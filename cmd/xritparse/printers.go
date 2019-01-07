package main

import (
	"encoding/hex"
	"fmt"
	"github.com/OpenSatelliteProject/SatHelperApp/XRIT/PacketData"
	"github.com/OpenSatelliteProject/SatHelperApp/XRIT/Structs"
	"strings"
)

func PrintAncillaryText(r *Structs.AncillaryText, printStructuredHeader, printImageDataRecord bool) {
	fmt.Println("Ancillary Record:")
	vs := strings.Split(r.StringData, ";")
	fmt.Println("    Data:")
	for _, v := range vs {
		if len(v) > 0 {
			fmt.Printf("        %s;\n", v)
		}
	}
}

func PrintAnnotationRecord(r *Structs.AnnotationRecord, printStructuredHeader, printImageDataRecord bool) {
	fmt.Println("Annotation Record:")
	fmt.Printf("    FileName: %s\n", r.StringData)
}

func PrintDCSFilenameRecord(r *Structs.DCSFilenameRecord, printStructuredHeader, printImageDataRecord bool) {
	fmt.Println("DCS FileName Record:")
	fmt.Printf("    FileName: %s\n", r.StringData)
}

func PrintHeaderStructuredRecord(r *Structs.HeaderStructuredRecord, printStructuredHeader, printImageDataRecord bool) {
	fmt.Println("Header Structured Record:")
	if printStructuredHeader {
		t := strings.Split(r.StringData, "UI")
		fmt.Println("   Data:")
		for _, v := range t {
			fmt.Printf("        %s\n", v)
		}
	} else {
		fmt.Println("   Data: {HIDDEN}")
	}
}

func PrintImageDataFunctionRecord(r *Structs.ImageDataFunctionRecord, printStructuredHeader, printImageDataRecord bool) {
	fmt.Println("Image Data Function Record:")
	if printImageDataRecord {
		fmt.Printf("   Data: %s\n", r.StringData)
	} else {
		fmt.Println("   Data: {HIDDEN}")
	}
}

func PrintImageNavigationRecord(r *Structs.ImageNavigationRecord, printStructuredHeader, printImageDataRecord bool) {
	fmt.Println("Image Navigation Record:")
	fmt.Printf("    Projection Name: %s\n", r.ProjectionName)
	fmt.Printf("    Column Scaling Factor: %d\n", r.ColumnScalingFactor)
	fmt.Printf("    Line Scaling Factor: %d\n", r.LineScalingFactor)
	fmt.Printf("    Column Offset: %d\n", r.ColumnOffset)
	fmt.Printf("    Line Offset: %d\n", r.LineOffset)
}

func PrintImageStructureRecord(r *Structs.ImageStructureRecord, printStructuredHeader, printImageDataRecord bool) {
	fmt.Println("Image Structure Header:")
	fmt.Printf("    Bits Per Pixel: %d\n", r.BitsPerPixel)
	fmt.Printf("    Columns: %d\n", r.Columns)
	fmt.Printf("    Lines: %d\n", r.Lines)
	fmt.Printf("    Compression: %s", PacketData.GetCompressionTypeString(int(r.Compression)))
}

func PrintNOAASpecificRecord(r *Structs.NOAASpecificRecord, printStructuredHeader, printImageDataRecord bool) {
	fmt.Println("NOAA Specific Header:")
	fmt.Printf("    Signature: %s\n", r.Signature)
	fmt.Printf("    Product ID: %s (%d)\n", r.Product().Name, r.ProductID)
	fmt.Printf("    Sub Product ID: %s (%d)\n", r.SubProduct().Name, r.ProductSubID)
	fmt.Printf("    Compression: %s\n", PacketData.GetCompressionTypeString(int(r.Compression)))
	fmt.Printf("    Parameter: %d\n", r.Parameter)
}

func PrintPrimaryRecord(r *Structs.PrimaryRecord, printStructuredHeader, printImageDataRecord bool) {
	fmt.Println("Primary Header:")
	fmt.Printf("    FileTypeCode: %s\n", PacketData.GetFileTypeCodeString(int(r.FileTypeCode)))
	fmt.Printf("    Header Length: %d\n", r.HeaderLength)
	fmt.Printf("    Data Length: %d\n", r.DataLength)
}

func PrintRiceCompressionRecord(r *Structs.RiceCompressionRecord, printStructuredHeader, printImageDataRecord bool) {
	fmt.Println("Rice Compression Record:")
	fmt.Printf("    Flags: %d\n", r.Flags)
	fmt.Printf("    Pixel: %d\n", r.Pixel)
	fmt.Printf("    Line: %d\n", r.Line)
}

func PrintSegmentIdentificationRecord(r *Structs.SegmentIdentificationRecord, printStructuredHeader, printImageDataRecord bool) {
	fmt.Println("Segment Identification Header:")
	fmt.Printf("    ImageID: %d\n", r.ImageID)
	fmt.Printf("    Sequence: %d\n", r.Sequence)
	fmt.Printf("    Start Column: %d\n", r.StartColumn)
	fmt.Printf("    Start Line: %d\n", r.StartLine)
	fmt.Printf("    Number of Segments: %d\n", r.MaxSegments)
	fmt.Printf("    Width: %d\n", r.MaxColumns)
	fmt.Printf("    Height: %d\n", r.MaxRows)
}
func PrintTimestampRecord(r *Structs.TimestampRecord, printStructuredHeader, printImageDataRecord bool) {
	fmt.Println("Timestamp Record:")
	fmt.Printf("    Days: %d\n", r.Days)
	fmt.Printf("    Milisseconds: %d\n", r.Milisseconds)
	fmt.Printf("    DateTime: %s\n", r.GetDateTime())
}

func PrintUnknownHeader(r *Structs.UnknownHeader, printStructuredHeader, printImageDataRecord bool) {
	fmt.Printf("Unknown Header (%d):", r.Type)
	v := strings.ToUpper(hex.EncodeToString(r.Data))
	fmt.Printf("    Hex Encoded Data: %s\n", v)
}
