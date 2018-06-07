/*
 * SpyserverDevice.h
 *
 *  Created on: 07/06/2018
 *      Author: Lucas Teske
 *      Based on Youssef Touil (youssef@live.com) C# implementation.
 */

#ifndef SRC_SPYSERVERFRONTEND_H_
#define SRC_SPYSERVERFRONTEND_H_

#define NOMINMAX
#include <SatHelper/sathelper.h>
#include <atomic>
#include "SpyserverProtocol.h"
#include "../DeviceParameters.h"


#define Q(x) #x
#define QUOTE(x) Q(x)

#ifndef MAJOR_VERSION
#define MAJOR_VERSION unk
#endif
#ifndef MINOR_VERSION
#define MINOR_VERSION unk
#endif
#ifndef MAINT_VERSION
#define MAINT_VERSION unk
#endif
#ifndef GIT_SHA1
#define GIT_SHA1 unk
#endif

enum ParserPhase {
  AcquiringHeader,
  ReadingData
};

#define SAMPLE_BUFFER_SIZE 256 * 1024

class SpyserverDevice {
private:
  static constexpr unsigned int BufferSize = 64 * 1024;
  const int DefaultDisplayPixels = 2000;
  const int DefaultFFTRange = 127;
  const uint32_t ProtocolVersion = SPYSERVER_PROTOCOL_VERSION;
  const std::string SoftwareID = std::string("SatHelperApp " QUOTE(MAJOR_VERSION) "."  QUOTE(MINOR_VERSION) "." QUOTE(MAINT_VERSION));
  const std::string NameNoDevice = std::string("SpyServer - No Device");
  const std::string NameAirspyOne = std::string("SpyServer - Airspy One");
  const std::string NameAirspyHF = std::string("SpyServer - Airspy HF+");
  const std::string NameRTLSDR = std::string("SpyServer - RTLSDR");
  const std::string NameUnknown = std::string("SpyServer - Unknown Device");

  GoDeviceCallback *cb;
  SatHelper::TcpClient client;

  std::atomic_bool terminated;
  std::atomic_bool streaming;
  std::atomic_bool gotDeviceInfo;
  std::atomic_bool gotSyncInfo;
  std::atomic_bool canControl;
  std::atomic_bool isConnected;

  uint8_t *headerData;
  uint8_t *bodyBuffer;
  uint64_t bodyBufferLength;
  uint32_t parserPosition;
  uint32_t lastSequenceNumber;

  std::thread *receiverThread;

  SatHelperException error;
  std::atomic_bool hasError;
  std::string hostname;
  int port;

  DeviceInfo deviceInfo;
  MessageHeader header;

  uint32_t streamingMode;
  uint32_t parserPhase;

  uint32_t droppedBuffers;
  std::atomic<int64_t> down_stream_bytes;


  uint32_t minimumTunableFrequency;
  uint32_t maximumTunableFrequency;
  uint32_t deviceCenterFrequency;
  uint32_t channelCenterFrequency;
  uint32_t channelDecimationStageCount;
  int32_t gain;

  std::vector<uint32_t> availableSampleRates;
  SatHelper::CircularBuffer<uint8_t> dataS8Queue;

  // Not the best way, I know
  float fBuffer[SAMPLE_BUFFER_SIZE];
  int16_t s16Buffer[SAMPLE_BUFFER_SIZE];

  bool SayHello();
  void Cleanup();
  void OnConnect();
  bool SetSetting(uint32_t settingType, std::vector<uint32_t> params);
  bool SendCommand(uint32_t cmd, std::vector<uint8_t> args);
  void ParseMessage(char *buffer, uint32_t len);
  int ParseHeader(char *buffer, uint32_t len);
  int ParseBody(char *buffer, uint32_t len);
  void ProcessDeviceInfo();
  void ProcessClientSync();
  void ProcessUInt8Samples();
  void ProcessInt16Samples();
  void ProcessFloatSamples();
  void ProcessUInt8FFT();
  void HandleNewMessage();
  void SetStreamState();
  void threadLoop();

public:
  SpyserverDevice(std::string hostname, int port);
  virtual ~SpyserverDevice();

  uint32_t SetSampleRate(uint32_t sampleRate);
  uint32_t SetCenterFrequency(uint32_t centerFrequency);
  const std::vector<uint32_t>& GetAvailableSampleRates();
  void Start();
  void Stop();
  void SetAGC(bool agc);

  void SetLNAGain(uint8_t value);
  void SetVGAGain(uint8_t value);
  void SetMixerGain(uint8_t value);
  void SetBiasT(uint8_t value);

  uint32_t GetCenterFrequency();

  const std::string &GetName();

  uint32_t GetSampleRate();

  void SetSamplesAvailableCallback(GoDeviceCallback *cb);

  bool Init();
  void Destroy();

  void Connect();
  void Disconnect();

};

#endif /* SRC_SPYSERVERFRONTEND_H_ */
