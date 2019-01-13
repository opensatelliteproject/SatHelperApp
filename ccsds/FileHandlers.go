package ccsds

import (
	"archive/zip"
	"github.com/OpenSatelliteProject/SatHelperApp/Logger"
	"github.com/OpenSatelliteProject/SatHelperApp/XRIT"
	"github.com/OpenSatelliteProject/SatHelperApp/XRIT/PacketData"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func PostHandleFile(filename, outBase string) {
	xh, err := XRIT.ParseFile(filename)

	if err != nil {
		SLog.Error("Error parsing file %s: %s", filename, err)
		_ = os.Remove(filename)
		return
	}

	switch xh.Compression() {
	case PacketData.ZIP:
		SLog.Debug("File %s is a zip. Decompressing it.", filename)
		stripFileHeader(filename, int64(xh.PrimaryHeader.HeaderLength))
		handleZipFile(filename, outBase)
	case PacketData.GIF:
		handleRawFileStrip(filename, xh)
	case PacketData.JPEG:
		handleRawFileStrip(filename, xh)
	}
}

func handleRawFileStrip(filename string, xh *XRIT.Header) {
	SLog.Debug("File %s is a %s.", filename, PacketData.GetCompressionTypeString(xh.Compression()))
	newName := strings.Replace(filename, ".lrit", PacketData.GetCompressionTypeExtension(xh.Compression()), 1)
	stripFileHeader(filename, int64(xh.PrimaryHeader.HeaderLength))
	err := os.Rename(filename, newName)
	if err != nil {
		SLog.Error("Error moving file %s to %s: %s", filename, newName, err)
	}
}

func handleZipFile(filename, outBase string) {
	r, err := zip.OpenReader(filename)
	if err != nil {
		SLog.Error("Error opening %s as zip file: %s", filename, err)
		return
	}

	for _, zipFile := range r.File {
		SLog.Info("New file decompressed: %s", zipFile.Name)
		outFile := path.Join(outBase, zipFile.Name)
		f, err := os.Create(outFile)
		if err != nil {
			SLog.Error("Error creating %s: %s", outFile, err)
			continue
		}

		r, err := zipFile.Open()

		if err != nil {
			SLog.Error("Internal ZIP Error: %s", err)
			_ = f.Close()
			continue
		}

		_, err = io.Copy(f, r)

		if err != nil {
			SLog.Error("Error writing data: %s", err)
		}

		_ = f.Close()
		_ = r.Close()
	}

	_ = r.Close()
	SLog.Debug("Removing file %s", filename)
	_ = os.Remove(filename)
}

func stripFileHeader(filename string, offset int64) {
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		SLog.Error("Error reading file %s: %s", filename, err)
		return
	}

	if len(data) < int(offset) {
		SLog.Error("File %s smaller than offset.", filename)
		return
	}

	data = data[offset:]

	err = ioutil.WriteFile(filename, data, os.ModePerm)

	if err != nil {
		SLog.Error("Error writing file %s: %s", filename, err)
	}
}
