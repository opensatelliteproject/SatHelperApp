package main

import (
	"github.com/OpenSatelliteProject/libsathelper"
	"github.com/foize/go.fifo"
	"time"
	"log"
	. "github.com/logrusorgru/aurora"
	"github.com/OpenSatelliteProject/SatHelperApp/Models"
)

func initDecoder() {
	if CurrentConfig.Decoder.UseLastFrameData {
		viterbiData = make([]byte, CODEDFRAMESIZE + LASTFRAMEDATABITS)
		decodedData = make([]byte, FRAMESIZE + LASTFRAMEDATA)
		lastFrameEnd = make([]byte, LASTFRAMEDATABITS)

		viterbi = SatHelper.NewViterbi27(FRAMEBITS + LASTFRAMEDATABITS)

		for i := 0; i < LASTFRAMEDATABITS; i++ {
			lastFrameEnd[i] = 128
		}
	} else {
		viterbiData = make([]byte, CODEDFRAMESIZE)
		decodedData = make([]byte, FRAMESIZE)

		viterbi = SatHelper.NewViterbi27(FRAMEBITS)
	}

	codedData = make([]byte, CODEDFRAMESIZE)
	rsCorrectedData = make([]byte, FRAMESIZE)
	rsWorkBuffer = make([]byte, 255)

	reedSolomon = SatHelper.NewReedSolomon()
	correlator = SatHelper.NewCorrelator()
	packetFixer = SatHelper.NewPacketFixer()

	syncWord = make([]byte, 4)

	reedSolomon.SetCopyParityToOutput(true)

	if CurrentConfig.Base.Mode == "lrit" {
		correlator.AddWord(LritUw0)
		correlator.AddWord(LritUw2)
	} else {
		correlator.AddWord(HritUw0)
		correlator.AddWord(HritUw2)
	}

	symbolsFifo = fifo.NewQueue()

	log.Printf(Cyan("Use Last Frame Data: %t").String(), Bold(Green(CurrentConfig.Decoder.UseLastFrameData)))
}

func decoderLoop() {
	isCorrupted := false
	lastFrameOk := false

	var localStats Models.Statistics
	var averageRSCorrections float32 = 0.0
	var averageVitCorrections float32 = 0.0
	var lostPacketsPerChannel [256]int64
	var lastPacketCount [256]int64
	var receivedPacketsPerChannel [256]int64
	var flywheelCount = 0

	for running {
		if symbolsFifo.Len() >= CODEDFRAMESIZE {
			if localStats.TotalPackets % AverageLastNSamples == 0 {
				averageRSCorrections = 0
				averageVitCorrections = 0
			}
			for i := 0; i < CODEDFRAMESIZE; i++ {
				codedData[i] = symbolsFifo.Next().(byte)
			}

			if flywheelCount == DefaultFlywheelRecheck {
				lastFrameOk = false
				flywheelCount = 0
			}

			// This reduces CPU Usage
			if !lastFrameOk {
				correlator.Correlate(&codedData[0], CODEDFRAMESIZE)
			} else {
				// If we got a good lock before, let's just check if the sync is in correct pos.

				correlator.Correlate(&codedData[0], CODEDFRAMESIZE / 16)
				if correlator.GetHighestCorrelationPosition() != 0 {
					// Oh no, that means something happened :/
					correlator.Correlate(&codedData[0], CODEDFRAMESIZE);
					lastFrameOk = false
					flywheelCount = 0
				}
			}
			flywheelCount++

			word := correlator.GetCorrelationWordNumber()
			pos := correlator.GetHighestCorrelationPosition()
			corr := correlator.GetHighestCorrelation()
			phaseShift := SatHelper.DEG_0
			if word == 1 {
				phaseShift = SatHelper.DEG_180
			}

			if corr < MINCORRELATIONBITS {
				log.Printf(Red("Correlation didn't match criteria of %d bits. Got %d\n").String(), Bold(MINCORRELATIONBITS), Bold(corr))
			}

			if pos != 0 {
				// Sync frame
				shiftWithConstantSize(&codedData, int(pos), CODEDFRAMESIZE)
				for symbolsFifo.Len() < int(pos) {
					// Wait enough data
					time.Sleep(time.Microsecond)
				}
				offset := CODEDFRAMESIZE - pos
				for i := offset; i < CODEDFRAMESIZE; i++ {
					codedData[i] = symbolsFifo.Next().(byte)
				}
			}

			if CurrentConfig.Base.Mode == "lrit" {
				packetFixer.FixPacket(&codedData[0], CODEDFRAMESIZE, phaseShift, false)
			}


			if CurrentConfig.Decoder.UseLastFrameData {
				for i := 0; i < LASTFRAMEDATABITS; i++ {
					viterbiData[i] = lastFrameEnd[i]
				}
				for i := LASTFRAMEDATABITS; i < CODEDFRAMESIZE + LASTFRAMEDATABITS; i++ {
					viterbiData[i] = codedData[i-LASTFRAMEDATABITS]
				}
			} else {
				for i := 0; i < CODEDFRAMESIZE; i++ {
					viterbiData[i] = codedData[i]
				}
			}

			viterbi.Decode(&viterbiData[0], &decodedData[0])

			if CurrentConfig.Base.Mode == "hrit" {
				nrzmDecodeSize := FRAMESIZE
				if CurrentConfig.Decoder.UseLastFrameData {
					nrzmDecodeSize += LASTFRAMEDATA
				}

				SatHelper.DifferentialEncodingNrzmDecode(&decodedData[0], nrzmDecodeSize)
			}

			signalErrors := float32(viterbi.GetPercentBER())
			signalErrors = 100 - (signalErrors * 10)
			signalQuality := uint8(signalErrors)

			averageVitCorrections += float32(viterbi.GetBER())

			if CurrentConfig.Decoder.UseLastFrameData {
				shiftWithConstantSize(&decodedData, LASTFRAMEDATA / 2, FRAMESIZE + LASTFRAMEDATA / 2)
				for i := 0; i < LASTFRAMEDATABITS; i++ {
					lastFrameEnd[i] = viterbiData[CODEDFRAMESIZE + i]
				}
			}

			for i:=0; i<4; i++ {
				syncWord[i] = decodedData[i]
				localStats.SyncWord[i] = decodedData[i]
			}

			shiftWithConstantSize(&decodedData, 4, FRAMESIZE - 4)

			localStats.AverageVitCorrections += uint16(viterbi.GetBER())
			localStats.TotalPackets += 1

			SatHelper.DeRandomizerDeRandomize(&decodedData[0], FRAMESIZE-4)

			derrors := make([]int32, RSBLOCKS)

			for i := 0; i < RSBLOCKS; i++ {
				reedSolomon.Deinterleave(&decodedData[0], &rsWorkBuffer[0], byte(i), RSBLOCKS)
				derrors[i] = int32(int8(reedSolomon.Decode_ccsds(&rsWorkBuffer[0])))
				reedSolomon.Interleave(&rsWorkBuffer[0], &rsCorrectedData[0], byte(i), RSBLOCKS)
				if derrors[i] != -1 {
					averageRSCorrections += float32(derrors[i])
				}
				localStats.RsErrors[i] = derrors[i]
			}

			if derrors[0] == -1 && derrors[1] == -1 && derrors[2] == -1 && derrors[3] == -1 {
				isCorrupted = true
				lastFrameOk = false
				localStats.DroppedPackets += 1
			} else {
				isCorrupted = false
				lastFrameOk = true
			}


			scid := ((rsCorrectedData[0] & 0x3F) << 2) | (rsCorrectedData[1] & 0xC0) >> 6
			vcid := rsCorrectedData[1] & 0x3F
			counter := uint(rsCorrectedData[2])
			counter = SatHelper.ToolsSwapEndianess(counter)
			counter &= 0xFFFFFF00
			counter >>= 8

			if ! isCorrupted {
				if lastPacketCount[vcid] + 1 != int64(counter) && lastPacketCount[vcid] > -1 {
					lostCount := int(int64(counter) - lastPacketCount[vcid] - 1)
					localStats.LostPackets += uint64(lostCount)
					lostPacketsPerChannel[vcid] += int64(lostCount)
				}
				lastPacketCount[vcid] = int64(counter)
				if receivedPacketsPerChannel[vcid] == -1 {
					receivedPacketsPerChannel[vcid] = 1
				} else {
					receivedPacketsPerChannel[vcid] = receivedPacketsPerChannel[vcid] + 1
				}

				localStats.SCID = scid
				localStats.VCID = vcid

				localStats.PacketNumber = uint64(counter)
				localStats.VitErrors = uint16(viterbi.GetBER())
				localStats.FrameBits = FRAMEBITS
				localStats.SignalQuality = signalQuality
				localStats.SyncCorrelation = uint8(corr)
				switch phaseShift {
					case SatHelper.DEG_0: localStats.PhaseCorrection = 0; break
					case SatHelper.DEG_90: localStats.PhaseCorrection = 1; break
					case SatHelper.DEG_180: localStats.PhaseCorrection = 2; break
					case SatHelper.DEG_270: localStats.PhaseCorrection = 3; break
				}

				if localStats.TotalPackets % AverageLastNSamples == 0 {
					localStats.AverageRSCorrections = uint8(averageRSCorrections / 4)
					localStats.AverageVitCorrections = uint16(averageVitCorrections)
				} else {
					localStats.AverageRSCorrections = uint8(averageRSCorrections / float32(4*(localStats.TotalPackets % AverageLastNSamples)))
					localStats.AverageVitCorrections = uint16(averageVitCorrections / float32(localStats.TotalPackets % AverageLastNSamples))
				}
				localStats.FrameLock = 1
				localStats.DecoderFifoUsage = uint8(100 * float32(symbolsFifo.Len()) / float32(FifoSize))
				localStats.DemodulatorFifoUsage = demodFifoUsage

				if demuxer != nil {
					demuxer.SendFrame(rsCorrectedData)
				}

			} else {
				localStats.FrameLock = 0
			}

			for i := 0; i < 256; i++ {
				localStats.ReceivedPacketsPerChannel[i] = receivedPacketsPerChannel[i]
				localStats.LostPacketsPerChannel[i] = lostPacketsPerChannel[i]
			}
			SetStats(localStats)
		} else {
			time.Sleep(time.Microsecond)
		}
	}
}