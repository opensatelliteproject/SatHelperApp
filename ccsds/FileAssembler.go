package ccsds

import (
	"fmt"
	"github.com/OpenSatelliteProject/SatHelperApp/Logger"
	"github.com/mewkiz/pkg/osutil"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type FileAssembler struct {
	tmpFolder string
	outFolder string
	msduCache map[int]*MSDUInfo
}

func MakeFileAssembler() *FileAssembler {
	return &FileAssembler{
		tmpFolder: "tmp",
		outFolder: "out",
		msduCache: make(map[int]*MSDUInfo),
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
			fa.handleFile(ofilename)
			fa.msduCache[msdu.APID] = nil
		}

		minfo := MakeMSDUInfo()
		minfo.APID = msdu.APID
		minfo.FileName = fmt.Sprintf("%03x.lrittmp", msdu.APID)
		minfo.LastPacketNumber = msdu.PacketNumber

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

	// TODO Handle Compression and other stuff

	mode := os.ModeAppend

	if firstOrSinglePacket {
		f, err := os.Create(filename)
		if err != nil {
			SLog.Error("Error creating file: %s", err)
			return
		}
		_ = f.Close()
	}

	err := ioutil.WriteFile(filename, dataToSave, mode)
	if err != nil {
		SLog.Error(err.Error())
	}

	if msdu.Sequence == SequenceLastSegment || msdu.Sequence == SequenceSingleData {
		fa.handleFile(filename)
		fa.msduCache[msdu.APID] = nil
	}
}

func (fa *FileAssembler) handleFile(filename string) {
	SLog.Debug("File to handle: %s", filename)
}
