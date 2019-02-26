package main

import (
	"fmt"
	"github.com/opensatelliteproject/SatHelperApp"
	"github.com/opensatelliteproject/SatHelperApp/XRIT"
	"github.com/opensatelliteproject/SatHelperApp/XRIT/PacketData"
	"github.com/opensatelliteproject/SatHelperApp/XRIT/Structs"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

func printHeaders(header *XRIT.Header, printStructuredHeader, printImageDataRecord bool) {
	for _, v := range header.AllHeaders {
		switch v.GetType() {
		case PacketData.AncillaryTextRecord:
			PrintAncillaryText(v.(*Structs.AncillaryText), printStructuredHeader, printImageDataRecord)
		case PacketData.AnnotationRecord:
			PrintAnnotationRecord(v.(*Structs.AnnotationRecord), printStructuredHeader, printImageDataRecord)
		case PacketData.DCSFileNameRecord:
			PrintDCSFilenameRecord(v.(*Structs.DCSFilenameRecord), printStructuredHeader, printImageDataRecord)
		case PacketData.HeaderStructuredRecord:
			PrintHeaderStructuredRecord(v.(*Structs.HeaderStructuredRecord), printStructuredHeader, printImageDataRecord)
		case PacketData.ImageDataFunctionRecord:
			PrintImageDataFunctionRecord(v.(*Structs.ImageDataFunctionRecord), printStructuredHeader, printImageDataRecord)
		case PacketData.ImageNavigationRecord:
			PrintImageNavigationRecord(v.(*Structs.ImageNavigationRecord), printStructuredHeader, printImageDataRecord)
		case PacketData.ImageStructureRecord:
			PrintImageStructureRecord(v.(*Structs.ImageStructureRecord), printStructuredHeader, printImageDataRecord)
		case PacketData.NOAASpecificHeader:
			PrintNOAASpecificRecord(v.(*Structs.NOAASpecificRecord), printStructuredHeader, printImageDataRecord)
		case PacketData.PrimaryHeader:
			PrintPrimaryRecord(v.(*Structs.PrimaryRecord), printStructuredHeader, printImageDataRecord)
		case PacketData.RiceCompressionRecord:
			PrintRiceCompressionRecord(v.(*Structs.RiceCompressionRecord), printStructuredHeader, printImageDataRecord)
		case PacketData.SegmentIdentificationRecord:
			PrintSegmentIdentificationRecord(v.(*Structs.SegmentIdentificationRecord), printStructuredHeader, printImageDataRecord)
		case PacketData.TimestampRecord:
			PrintTimestampRecord(v.(*Structs.TimestampRecord), printStructuredHeader, printImageDataRecord)
		default:
			PrintUnknownHeader(v.(*Structs.UnknownHeader), printStructuredHeader, printImageDataRecord)
		}
		fmt.Println("")
	}
}

func parseFile(filename string, printStructuredHeader, printImageDataRecord bool) {
	xh, err := XRIT.ParseFile(filename)
	if err != nil {
		fmt.Printf("Error parsing file %s: %s\n", filename, err)
		os.Exit(1)
	}

	printHeaders(xh, printStructuredHeader, printImageDataRecord)
}

func main() {
	kingpin.Version(SatHelperApp.GetVersion())

	files := kingpin.Arg("filename", "File name to parse").Required().ExistingFiles()

	printStructuredHeader := kingpin.Flag("h", "Print Structured Header Record").Default("false").Bool()
	printImageDataRecord := kingpin.Flag("i", "Print Image Data Record").Default("false").Bool()

	kingpin.Parse()

	for _, v := range *files {
		parseFile(v, *printStructuredHeader, *printImageDataRecord)
	}
}
