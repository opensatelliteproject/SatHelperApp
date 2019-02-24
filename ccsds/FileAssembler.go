package ccsds

import (
	"fmt"
	"github.com/OpenSatelliteProject/SatHelperApp/ImageProcessor"
	"github.com/OpenSatelliteProject/SatHelperApp/Logger"
	"github.com/OpenSatelliteProject/SatHelperApp/XRIT"
	"github.com/OpenSatelliteProject/SatHelperApp/XRIT/PacketData"
	"github.com/OpenSatelliteProject/goaec/szwrap"
	"github.com/mewkiz/pkg/osutil"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"sync"
)

type FileAssembler struct {
	sync.Mutex
	tmpFolder string
	outFolder string
	msduCache map[int]*MSDUInfo
	ip        *ImageProcessor.ImageProcessor
}

func MakeFileAssembler() *FileAssembler {
	return &FileAssembler{
		tmpFolder: "tmp",
		outFolder: "out",
		msduCache: make(map[int]*MSDUInfo),
		ip:        ImageProcessor.MakeImageProcessor(),
	}
}

func (fa *FileAssembler) SetTemporaryFolder(folder string) {
	fa.tmpFolder = folder
}

func (fa *FileAssembler) SetOutputFolder(folder string) {
	fa.outFolder = folder
}

func (fa *FileAssembler) PutMSDU(msdu *MSDU) {
	if msdu.Sequence == SequenceFirstSegment || msdu.Sequence == SequenceSingleData {
		if fa.msduCache[msdu.APID] != nil {
			minfo := fa.msduCache[msdu.APID]
			SLog.Warn("Received first segment to %03x but last data wasn't saved to disk yet! Forcing dump.", msdu.APID)
			ofilename := path.Join(fa.tmpFolder, strconv.FormatInt(int64(msdu.ChannelId), 10))
			ofilename = path.Join(ofilename, minfo.FileName)
			fa.handleFile(msdu.ChannelId, ofilename)
			fa.msduCache[msdu.APID] = nil
		}

		minfo := MakeMSDUInfo()
		minfo.APID = msdu.APID
		minfo.FileName = fmt.Sprintf("%03x.lrittmp", msdu.APID)
		minfo.LastPacketNumber = msdu.PacketNumber
		head, err := XRIT.MemoryParseFile(msdu.Data[10:])
		if err != nil {
			SLog.Error("Error parsing XRIT Header: %s", err.Error())
			minfo.Header = XRIT.MakeXRITHeader()
		}
		minfo.Header = head

		fa.msduCache[minfo.APID] = minfo
	} else if msdu.Sequence == SequenceLastSegment || msdu.Sequence == SequenceContinuedSegment {
		if fa.msduCache[msdu.APID] == nil {
			SLog.Debug("Orphan packet for APID %03x!", msdu.APID)
			return
		}
	}

	firstOrSinglePacket := msdu.Sequence == SequenceFirstSegment || msdu.Sequence == SequenceSingleData

	msduInfo := fa.msduCache[msdu.APID]
	msduInfo.Refresh()

	tmpPath := path.Join(fa.tmpFolder, strconv.FormatInt(int64(msdu.ChannelId), 10))

	if !osutil.Exists(tmpPath) {
		_ = os.MkdirAll(tmpPath, 0777)
	}

	filename := path.Join(tmpPath, msduInfo.FileName)

	dataToSave := msdu.Data

	if firstOrSinglePacket {
		dataToSave = dataToSave[10:]
	}

	missedPackets := msduInfo.LastPacketNumber - msdu.PacketNumber - 1

	if msduInfo.LastPacketNumber == 16383 && msdu.PacketNumber == 0 {
		missedPackets = 0
	}

	if msduInfo.Header.Compression() == PacketData.LRIT_RICE && !firstOrSinglePacket {
		if missedPackets > 0 {
			SLog.Warn("Missed %d packets on image. Filling with null bytes. Last Packet Number %d and current %d", missedPackets, msduInfo.LastPacketNumber, msdu.PacketNumber)
			fill := make([]byte, msduInfo.Header.ImageStructureHeader.Columns)
			for missedPackets > 0 {
				_ = ioutil.WriteFile(filename, fill, os.ModeAppend)
			}
		}

		if msduInfo.Header.RiceCompressionHeader == nil { // Fix bug for GOES-15 TX after GOES-16 Switch. Weird but let's try defaults
			d, err := szwrap.NOAADecompress(dataToSave, 8, 16, int(msduInfo.Header.ImageStructureHeader.Columns), szwrap.SZ_ALLOW_K13_OPTION_MASK|szwrap.SZ_MSB_OPTION_MASK|szwrap.SZ_NN_OPTION_MASK)
			if err != nil {
				SLog.Error("Error decompressing: %s", err.Error())
				return
			}
			dataToSave = d
		} else {
			d, err := szwrap.NOAADecompress(dataToSave, 8, int(msduInfo.Header.RiceCompressionHeader.Pixel), int(msduInfo.Header.ImageStructureHeader.Columns), int(msduInfo.Header.RiceCompressionHeader.Flags))
			if err != nil {
				SLog.Error("Error decompressing: %s", err.Error())
				return
			}
			dataToSave = d
		}
	}

	msduInfo.LastPacketNumber = msdu.PacketNumber

	mode := os.O_WRONLY

	if firstOrSinglePacket {
		mode |= os.O_CREATE
	} else {
		mode |= os.O_APPEND
	}

	defer func() {
		if msdu.Sequence == SequenceLastSegment || msdu.Sequence == SequenceSingleData {
			fa.handleFile(msdu.ChannelId, filename)
			fa.msduCache[msdu.APID] = nil
		}
	}()

	f, err := os.OpenFile(filename, mode, 0777)
	if err != nil {
		SLog.Error(err.Error())
		return
	}

	n, err := f.Write(dataToSave)
	if err != nil {
		SLog.Error(err.Error())
	}

	if n != len(dataToSave) {
		SLog.Error("Error saving all data. Expected %d bytes saved %d bytes", len(dataToSave), n)
	}

	_ = f.Close()
}

func (fa *FileAssembler) handleFile(vcid int, filename string) {
	//SLog.Debug("File to handle: %s", filename)
	xh, err := XRIT.ParseFile(filename)

	if err != nil {
		SLog.Error("FileAssembler::handleFile - Error parsing file %s: %s", filename, err)
		_ = os.Remove(filename)
		return
	}

	pathName := fmt.Sprintf("VCID-%d", vcid)
	n, ok := XRIT.VCID2Name[vcid]

	if ok {
		pathName = n
	}

	outBase := path.Join(fa.outFolder, pathName)

	if !osutil.Exists(outBase) {
		_ = os.MkdirAll(outBase, 0777)
	}

	if xh.SegmentIdentificationHeader != nil && xh.SegmentIdentificationHeader.MaxSegments != 1 {
		xh.Filename()
	}

	SLog.Info("New file (%s): %s", xh.ToNameString(), path.Join(pathName, xh.Filename()))
	newPath := path.Join(outBase, xh.Filename())
	//SLog.Debug("Moving %s to %s", filename, newPath)
	err = os.Rename(filename, newPath)
	if err != nil {
		SLog.Error("Error moving file %s to %s: %s", filename, newPath, err)
	}

	go PostHandleFile(newPath, outBase, vcid, fa.ip)
}
