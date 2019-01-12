package main

import (
	"fmt"
	"github.com/OpenSatelliteProject/SatHelperApp"
	"github.com/OpenSatelliteProject/SatHelperApp/XRIT"
	"github.com/OpenSatelliteProject/SatHelperApp/XRIT/Structs"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"io/ioutil"
	"os"
)

func catFile(filename string) {
	xh, err := XRIT.ParseFile(filename)
	if err != nil {
		fmt.Printf("Error parsing file %s: %s\n", filename, err)
		os.Exit(1)
	}
	offset := xh.PrimaryHeader.HeaderLength

	f, err := os.Open(filename)

	if err != nil {
		fmt.Printf("Error parsing file %s: %s\n", filename, err)
		os.Exit(1)
	}

	_, err = f.Seek(int64(offset), io.SeekStart)
	if err != nil {
		fmt.Printf("Error parsing file %s: %s\n", filename, err)
		os.Exit(1)
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Printf("Error parsing file %s: %s\n", filename, err)
		os.Exit(1)
	}

	d := Structs.ParseDCS(data)
	fmt.Printf("Header: %s\n", d.Header)
	fmt.Println("   # Address             Date / Time           Status  Signal  Frequency Offset  MIN  DQN  Channel  Source  ")
	for i, v := range d.Packets {
		fmt.Printf("%4d %8s  %29s    %1s      %2d dB         %2d          %1s    %1s    %4s      %2s    \n", i, v.Address, v.DateTime, v.Status, v.Signal, v.FrequencyOffset, v.ModIndexNormal, v.DataQualNominal, v.Channel, v.SourceCode)
	}
}

func main() {
	kingpin.Version(SatHelperApp.GetVersion())

	files := kingpin.Arg("filename", "DCS File name to print content").Required().ExistingFiles()

	kingpin.Parse()

	for _, v := range *files {
		catFile(v)
	}
}
