/*
 * RtlFrontend.h
 *
 *  Created on: 25/02/2017
 *      Author: Lucas Teske
 */

#ifndef SRC_RTLFRONTEND_H_
#define SRC_RTLFRONTEND_H_

#include <thread>
#include <cstdint>
#include <iostream>
#include <sstream>
#include <vector>
#include <string>
#include <functional>
#include <SatHelper/exceptions.h>

extern "C" {
#include <rtl-sdr.h>
}

#include "../DeviceParameters.h"

class RtlFrontend {
private:
  static std::vector<uint32_t> supportedSampleRates;

  uint32_t sampleRate;
  uint32_t centerFrequency;
  int deviceId;
  rtlsdr_dev_t *device;
  std::string deviceName;
  float lut[256];
  float alpha;
  float iavg;
  float qavg;
  std::thread *mainThread;
  uint8_t lnaGain;
  uint8_t vgaGain;
  uint8_t mixerGain;
  bool agc;
  GoDeviceCallback *cb;

  void threadWork();
  void refreshGains();
  void internalCallback(unsigned char *data, unsigned int length);
protected:
  static void rtlCallback(unsigned char *data, unsigned int length, void *ctx);

public:
  RtlFrontend(GoDeviceCallback *cb);
  virtual ~RtlFrontend();

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
};

#endif /* SRC_RTLFRONTEND_H_ */
