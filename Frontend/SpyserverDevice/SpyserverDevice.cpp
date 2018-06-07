/*
 * SpyserverDevice.cpp
 *
 *  Created on: 21/04/2017
 *      Author: Lucas Teske
 *      Based on Youssef Touil (youssef@live.com) C# implementation.
 */

#include "SpyserverDevice.h"
#include <cstring>
#include <algorithm>

SpyserverDevice::SpyserverDevice(std::string hostname, int port) :
    client(hostname, port), terminated(false), streaming(false), gotDeviceInfo(false),
    gotSyncInfo(false), canControl(false), isConnected(false), headerData(new uint8_t[sizeof(MessageHeader)]), bodyBuffer(NULL),
    bodyBufferLength(0), parserPosition(0), lastSequenceNumber(0),
    receiverThread(NULL), hasError(false), hostname(hostname), port(port), streamingMode(STREAM_MODE_IQ_ONLY),
    parserPhase(0), droppedBuffers(0), down_stream_bytes(0), minimumTunableFrequency(0), maximumTunableFrequency(0),
    deviceCenterFrequency(0), channelCenterFrequency(0), channelDecimationStageCount(0), gain(0),
    dataS8Queue(SAMPLE_BUFFER_SIZE), cb(NULL) {

  // std::cout << "SpyserverDevice(" << hostname << ", " << port << ")" << std::endl;
}

SpyserverDevice::~SpyserverDevice() {
  Disconnect();
  delete[] headerData;
}

void SpyserverDevice::SetBiasT(uint8_t value) {
  // // std::cerr << "Not supported by OSP" << std::endl;
}


bool SpyserverDevice::Init() {
  try {
    Connect();
    return true;
  } catch (std::exception &e) {
    // std::cerr << e.what() << std::endl;
    return false;
  }
}

void SpyserverDevice::Destroy() {
  Disconnect();
}

void SpyserverDevice::Connect() {
  if (receiverThread != NULL) {
    return;
  }
  // std::cout << "SpyServer: Trying to connect" << std::endl;
  client.Connect();
  isConnected = true;
  // std::cout << "SpyServer: Connected" << std::endl;

  SayHello();
  Cleanup();

  terminated = false;
  hasError = false;
  gotSyncInfo = false;
  gotDeviceInfo = false;

  receiverThread  = new std::thread(&SpyserverDevice::threadLoop, this);

  for (int i=0; i<1000 && !hasError; i++) {
    if (gotDeviceInfo) {
      if (deviceInfo.DeviceType == DEVICE_INVALID) {
        error = SatHelperException("Server is up but no device is available");
        hasError = true;
        break;
      }

      if (gotSyncInfo) {
        // std::cout << "Got sync Info" << std::endl;
        OnConnect();
        return;
      }
    }
    std::this_thread::sleep_for(std::chrono::milliseconds(1));
  }
  Disconnect();
  if (hasError) {
    hasError = false;
    throw error;
  }

  throw SatHelperException("Server didn't send the device capability and synchronization info.");
}

void SpyserverDevice::Disconnect() {
  // std::cerr << "DISCONNECTED" << std::endl;
  terminated = true;
  if (isConnected) {
    client.Close();
  }

  if (receiverThread != NULL) {
    receiverThread->join();
    receiverThread = NULL;
  }

  Cleanup();
}

void SpyserverDevice::OnConnect() {
  SetSetting(SETTING_STREAMING_MODE, { streamingMode });
  SetSetting(SETTING_IQ_FORMAT, { STREAM_FORMAT_INT16 });
  SetSetting(SETTING_FFT_FORMAT, { STREAM_FORMAT_UINT8 });
  //SetSetting(SETTING_FFT_DISPLAY_PIXELS, { displayPixels });
  //SetSetting(SETTING_FFT_DB_OFFSET, { fftOffset });
  //SetSetting(SETTING_FFT_DB_RANGE, { fftRange });
  //deviceInfo.MaximumSampleRate
  //availableSampleRates
  for (unsigned int i=0; i<=deviceInfo.DecimationStageCount; i++) {
    availableSampleRates.push_back(deviceInfo.MaximumSampleRate / (double)(1 << i));
  }
}

bool SpyserverDevice::SetSetting(uint32_t settingType, std::vector<uint32_t> params) {
  std::vector<uint8_t> argBytes;
  if (params.size() > 0) {
    argBytes = std::vector<uint8_t>(sizeof(SettingType) + params.size() * sizeof(uint32_t));
    uint8_t *settingBytes = (uint8_t *) &settingType;
    for (unsigned int i=0; i<sizeof(uint32_t); i++) {
      argBytes[i] = settingBytes[i];
    }

    std::memcpy(&argBytes[0]+sizeof(uint32_t), &params[0], sizeof(uint32_t) * params.size());
  } else {
    argBytes = std::vector<uint8_t>();
  }

  return SendCommand(CMD_SET_SETTING, argBytes);
}

bool SpyserverDevice::SayHello() {
  const uint8_t *protocolVersionBytes = (const uint8_t *) &ProtocolVersion;
  const uint8_t *softwareVersionBytes = (const uint8_t *) SoftwareID.c_str();
  std::vector<uint8_t> args = std::vector<uint8_t>(sizeof(ProtocolVersion) + SoftwareID.size());

  std::memcpy(&args[0], protocolVersionBytes, sizeof(ProtocolVersion));
  std::memcpy(&args[0] + sizeof(ProtocolVersion), softwareVersionBytes, SoftwareID.size());

  return SendCommand(CMD_HELLO, args);
}

void SpyserverDevice::Cleanup() {
    deviceInfo.DeviceType = 0;
    deviceInfo.DeviceSerial = 0;
    deviceInfo.DecimationStageCount = 0;
    deviceInfo.GainStageCount = 0;
    deviceInfo.MaximumSampleRate = 0;
    deviceInfo.MaximumBandwidth = 0;
    deviceInfo.MaximumGainIndex = 0;
    deviceInfo.MinimumFrequency = 0;
    deviceInfo.MaximumFrequency = 0;

    gain = 0;
    //displayCenterFrequency = 0;
    //deviceCenterFrequency = 0;
    //displayDecimationStageCount = 0;
    //channelDecimationStageCount = 0;
    //minimumTunableFrequency = 0;
    //maximumTunableFrequency = 0;
    canControl = false;
    gotDeviceInfo = false;
    gotSyncInfo = false;

    lastSequenceNumber = ((uint32_t)-1);
    droppedBuffers = 0;
    down_stream_bytes = 0;

    parserPhase = AcquiringHeader;
    parserPosition = 0;

    streaming = false;
    terminated = true;
}

void SpyserverDevice::threadLoop() {
  parserPhase = AcquiringHeader;
  parserPosition = 0;

  char buffer[BufferSize];
  try {
    while(!terminated) {
      if (terminated) {
        break;
      }
      uint32_t availableData = client.AvailableData();
      if (availableData > 0) {
        availableData = availableData > BufferSize ? BufferSize : availableData;
        client.Receive(buffer, availableData);
        ParseMessage(buffer, availableData);
      }
    }
  } catch (SatHelperException &e) {
    error = e;
    // std::cerr << e.what() << std::endl;
  }
  if (bodyBuffer != NULL) {
    delete[] bodyBuffer;
    bodyBuffer = NULL;
  }

  Cleanup();
}

void SpyserverDevice::ParseMessage(char *buffer, uint32_t len) {
  down_stream_bytes++;

  int consumed;
  while (len > 0 && !terminated) {
    if (parserPhase == AcquiringHeader) {
      while (parserPhase == AcquiringHeader && len > 0) {
        consumed = ParseHeader(buffer, len);
        buffer += consumed;
        len -= consumed;
      }

      if (parserPhase == ReadingData) {
        uint8_t client_major = (SPYSERVER_PROTOCOL_VERSION >> 24) & 0xFF;
        uint8_t client_minor = (SPYSERVER_PROTOCOL_VERSION >> 16) & 0xFF;

        uint8_t server_major = (header.ProtocolID >> 24) & 0xFF;
        uint8_t server_minor = (header.ProtocolID >> 16) & 0xFF;
        //uint16_t server_build = (header.ProtocolID & 0xFFFF);

        if (client_major != server_major || client_minor != server_minor) {
          throw SatHelperException("Server is running an unsupported protocol version.");
        }

        if (header.BodySize > SPYSERVER_MAX_MESSAGE_BODY_SIZE) {
          throw SatHelperException("The server is probably buggy.");
        }

        if (bodyBuffer == NULL || bodyBufferLength < header.BodySize) {
          if (bodyBuffer != NULL) {
            delete[] bodyBuffer;
          }

          bodyBuffer = new uint8_t[header.BodySize];
        }
      }
    }

    if (parserPhase == ReadingData) {
      consumed = ParseBody(buffer, len);
      buffer += consumed;
      len -= consumed;

      if (parserPhase == AcquiringHeader) {
        if (header.MessageType != MSG_TYPE_DEVICE_INFO && header.MessageType != MSG_TYPE_CLIENT_SYNC) {
          int32_t gap = header.SequenceNumber - lastSequenceNumber - 1;
          lastSequenceNumber = header.SequenceNumber;
          droppedBuffers += gap;
          if (gap > 0) {
            // std::cerr << "Lost " << gap << " frames from SpyServer!" << std::endl;
          }
        }
        HandleNewMessage();
      }
    }
  }
}

int SpyserverDevice::ParseHeader(char *buffer, uint32_t length) {
  auto consumed = 0;

  while (length > 0) {
    int to_write = std::min((uint32_t)(sizeof(MessageHeader) - parserPosition), length);
    std::memcpy(&header + parserPosition, buffer, to_write);
    length -= to_write;
    buffer += to_write;
    parserPosition += to_write;
    consumed += to_write;
    if (parserPosition == sizeof(MessageHeader)) {
      parserPosition = 0;
      if (header.BodySize > 0) {
        parserPhase = ReadingData;
      }

      return consumed;
    }
  }

  return consumed;
}

int SpyserverDevice::ParseBody(char* buffer, uint32_t length) {
  auto consumed = 0;

  while (length > 0) {
    int to_write = std::min((int) header.BodySize - parserPosition, length);
    std::memcpy(bodyBuffer + parserPosition, buffer, to_write);
    length -= to_write;
    buffer += to_write;
    parserPosition += to_write;
    consumed += to_write;

    if (parserPosition == header.BodySize) {
      parserPosition = 0;
      parserPhase = AcquiringHeader;
      return consumed;
    }
  }

  return consumed;
}

bool SpyserverDevice::SendCommand(uint32_t cmd, std::vector<uint8_t> args) {
  if (!isConnected) {
    return false;
  }

  bool result;
  uint32_t headerLen = sizeof(CommandHeader);
  uint16_t argLen = args.size();
  uint8_t *buffer = new uint8_t[headerLen + argLen];

  CommandHeader header;
  header.CommandType = cmd;
  header.BodySize = argLen;

  for (uint32_t i=0; i<sizeof(CommandHeader); i++) {
    buffer[i] = ((uint8_t *)(&header))[i];
  }

  if (argLen > 0) {
    for (uint16_t i=0; i<argLen; i++) {
      buffer[i+headerLen] = args[i];
    }
  }

  try {
    client.Send((char *)buffer, headerLen+argLen);
    result = true;
  } catch (SatHelperException &e) {
    result = false;
  }

  delete[] buffer;
  return result;
}

void SpyserverDevice::HandleNewMessage() {
  if (terminated) {
    return;
  }

  switch (header.MessageType) {
    case MSG_TYPE_DEVICE_INFO:
      ProcessDeviceInfo();
      break;
    case MSG_TYPE_CLIENT_SYNC:
      ProcessClientSync();
      break;
    case MSG_TYPE_UINT8_IQ:
      ProcessUInt8Samples();
      break;
    case MSG_TYPE_INT16_IQ:
      ProcessInt16Samples();
      break;
    case MSG_TYPE_FLOAT_IQ:
      ProcessFloatSamples();
      break;
    case MSG_TYPE_UINT8_FFT:
      ProcessUInt8FFT();
      break;
    default:
      break;
  }
}

void SpyserverDevice::ProcessDeviceInfo() {
  std::memcpy(&deviceInfo, bodyBuffer, sizeof(DeviceInfo));
  minimumTunableFrequency = deviceInfo.MinimumFrequency;
  maximumTunableFrequency = deviceInfo.MaximumFrequency;
  gotDeviceInfo = true;
}

void SpyserverDevice::ProcessClientSync() {
  ClientSync sync;
  std::memcpy(&sync, bodyBuffer, sizeof(ClientSync));

  canControl = sync.CanControl != 0;
  gain = (int) sync.Gain;
  deviceCenterFrequency = sync.DeviceCenterFrequency;
  channelCenterFrequency = sync.IQCenterFrequency;
  //displayCenterFrequency = sync.FFTCenterFrequency;

  switch (streamingMode) {
  case STREAM_MODE_FFT_ONLY:
  case STREAM_MODE_FFT_IQ:
    minimumTunableFrequency = sync.MinimumFFTCenterFrequency;
    maximumTunableFrequency = sync.MaximumFFTCenterFrequency;
    break;
  case STREAM_MODE_IQ_ONLY:
    minimumTunableFrequency = sync.MinimumIQCenterFrequency;
    maximumTunableFrequency = sync.MaximumIQCenterFrequency;
    break;
  }

  gotSyncInfo = true;
}

void SpyserverDevice::ProcessUInt8Samples() {
  // Spy Server sends out uint8_t that is the sample shifted by 128.
  // So we store as uint8, but sends to the callback converted as float.
  uint32_t numSamples = header.BodySize;
  if (dataS8Queue.size() + numSamples >= SAMPLE_BUFFER_SIZE) {
    uint32_t samplesToAdd = SAMPLE_BUFFER_SIZE - dataS8Queue.size();
    dataS8Queue.addSamples(bodyBuffer, samplesToAdd);
    numSamples -= samplesToAdd;

    // CircularBuffer is now full, so copy to output buffer
    dataS8Queue.unsafe_lockMutex();
    for (int i=0; i<SAMPLE_BUFFER_SIZE; i++) {
      fBuffer[i] = (dataS8Queue.unsafe_takeSample() - 128) / 128.f;
    }
    dataS8Queue.unsafe_unlockMutex();

    // Write output to callback
    if (cb != NULL) {
      cb->cbFloatIQ(fBuffer, SAMPLE_BUFFER_SIZE / 2);
    }

    // Add Remaining Samples.
    if (numSamples > 0) {
      dataS8Queue.addSamples((bodyBuffer) + samplesToAdd, numSamples);
    }
  } else {
    // Circular Buffer has enough space, so just add.
    dataS8Queue.addSamples(bodyBuffer, numSamples);
  }
}

void SpyserverDevice::ProcessInt16Samples() {
  uint32_t numSamples = header.BodySize / 2;
  if (cb != NULL) {
    cb->cbS16IQ((int16_t *)bodyBuffer, numSamples / 2);
  }
}

void SpyserverDevice::ProcessFloatSamples() {
  uint32_t numSamples = header.BodySize / 4;
  if (cb != NULL) {
    cb->cbFloatIQ((float *)bodyBuffer, numSamples / 2);
  }
}

void SpyserverDevice::ProcessUInt8FFT() {
  // TODO
  // // std::cerr << "UInt8 FFT Samples processing not implemented!!!" << std::endl;
}

void SpyserverDevice::SetStreamState() {
  SetSetting(SETTING_STREAMING_ENABLED, {(unsigned int)(streaming ? 1 : 0)});
}

uint32_t SpyserverDevice::SetSampleRate(uint32_t sampleRate) {
  for (unsigned int i=0; i<availableSampleRates.size(); i++) {
    if (availableSampleRates[i] == sampleRate) {
            channelDecimationStageCount = i;
            SetSetting(SETTING_IQ_DECIMATION, {channelDecimationStageCount});
            return GetSampleRate();
    }
  }
  std::cerr << "Sample rate not supported: " << sampleRate << std::endl;
  std::cerr << "Supported Sample Rates: " << std::endl;
  for (uint32_t sr: availableSampleRates) {
    std::cout << "  " << sr << std::endl;
  }
  return GetSampleRate();
}

uint32_t SpyserverDevice::SetCenterFrequency(uint32_t centerFrequency) {
    channelCenterFrequency = centerFrequency;
    SetSetting(SETTING_IQ_FREQUENCY, {channelCenterFrequency});
    return centerFrequency;
}

const std::vector<uint32_t>& SpyserverDevice::GetAvailableSampleRates() {
  return availableSampleRates;
}

void SpyserverDevice::Start() {
    if (!streaming) {
        streaming = true;
        down_stream_bytes = 0;
        SetStreamState();
    }
}

void SpyserverDevice::Stop() {
    if (streaming) {
        streaming = false;
        down_stream_bytes = 0;
        SetStreamState();
    }
}

void SpyserverDevice::SetAGC(bool agc) {
  // std::cerr << "AGC Not Supported" << std::endl;
}

void SpyserverDevice::SetLNAGain(uint8_t value) {
    gain = value;
    SetSetting(SETTING_GAIN, {(uint32_t)gain});
}

void SpyserverDevice::SetVGAGain(uint8_t value) {
  // std::cerr << "VGA Gain Not Supported" << std::endl;
}

void SpyserverDevice::SetMixerGain(uint8_t value) {
  // std::cerr << "Mixer Gain Not Supported" << std::endl;
}

uint32_t SpyserverDevice::GetCenterFrequency() {
  return channelCenterFrequency;
}

const std::string &SpyserverDevice::GetName() {
  switch (deviceInfo.DeviceType) {
  case DEVICE_INVALID:
    return SpyserverDevice::NameNoDevice;
  case DEVICE_AIRSPY_ONE:
    return SpyserverDevice::NameAirspyOne;
  case DEVICE_AIRSPY_HF:
    return SpyserverDevice::NameAirspyHF;
  case DEVICE_RTLSDR:
    return SpyserverDevice::NameRTLSDR;
  default:
    return SpyserverDevice::NameUnknown;
  }
}

uint32_t SpyserverDevice::GetSampleRate() {
  return (int) (deviceInfo.MaximumSampleRate / (double) (1 << channelDecimationStageCount));
}
void SpyserverDevice::SetSamplesAvailableCallback(GoDeviceCallback *cb) {
  this->cb = cb;
}

