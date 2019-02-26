package main

import (
	"fmt"
	"github.com/opensatelliteproject/SatHelperApp"
	"github.com/opensatelliteproject/SatHelperApp/XRIT"
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

	_, _ = os.Stdout.Write(data)
}

func main() {
	kingpin.Version(SatHelperApp.GetVersion())

	files := kingpin.Arg("filename", "File name to print content").Required().ExistingFiles()

	kingpin.Parse()

	for _, v := range *files {
		catFile(v)
	}
}
