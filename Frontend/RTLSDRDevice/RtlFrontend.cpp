/*
 * RtlFrontend.cpp
 *
 *  Created on: 25/02/2017
 *      Author: Lucas Teske
 */

#include "RtlFrontend.h"
#include <cmath>

std::vector<uint32_t> RtlFrontend::supportedSampleRates = {
  240000, 300000, 960000, 1152000, 1200000, 1440000, 1600000, 1800000, 1920000, 2400000, 2560000, 2880000, 3200000
};

RtlFrontend::RtlFrontend(GoDeviceCallback *cb) :
    sampleRate(2560000), centerFrequency(106300000), deviceId(0), alpha(0), iavg(0), qavg(0), lnaGain(0), vgaGain(0), mixerGain(0),
     agc(false) {
  this->cb = cb;
  for (int i = 0; i < 256; i++) {
    lut[i] = (i - 128) * (1.f / 127.f);
  }
  mainThread = NULL;
}

RtlFrontend::~RtlFrontend() {
  rtlsdr_close(device);
}

bool RtlFrontend::Init() {
  if (rtlsdr_open(&device, deviceId) != 0) {
	Log(cb).Get(logERROR) << "Failed to open rtlsdr with id " << deviceId << std::endl;
    return false;
  }

  deviceName = std::string(rtlsdr_get_device_name(deviceId));

  return true;
}

void RtlFrontend::SetBiasT(uint8_t value) {
  std::cerr << "BiasT in RtlSdr is not supported by OSP" << std::endl;
}

uint32_t RtlFrontend::SetSampleRate(uint32_t sampleRate) {
  this->sampleRate = sampleRate;
  rtlsdr_set_sample_rate(device, sampleRate);
  return sampleRate;
}

uint32_t RtlFrontend::SetCenterFrequency(uint32_t centerFrequency) {
  this->centerFrequency = centerFrequency;
  rtlsdr_set_center_freq(device, centerFrequency);
  return centerFrequency;
}

const std::vector<uint32_t>& RtlFrontend::GetAvailableSampleRates() {
  return RtlFrontend::supportedSampleRates;
}

void RtlFrontend::Start() {
  alpha = 1.f - exp(-1.0 / (sampleRate * 0.05f));
  iavg = 0;
  qavg = 0;

  if (mainThread != NULL) {
    throw SatHelperException("The worker is already running!");
  }
  tuner = rtlsdr_get_tuner_type(device);

  if (rtlsdr_set_sample_rate(device, sampleRate) != 0) {
    std::cerr << "Cannot set sample rate to " << sampleRate << std::endl;
    throw SatHelperException("Cannot set sample rate.");
  }

  if (rtlsdr_set_center_freq(device, centerFrequency) != 0) {
    std::cerr << "Cannot set center frequency to " << centerFrequency
        << std::endl;
    throw SatHelperException("Cannot set center frequency.");
  }

  if (rtlsdr_set_tuner_gain_mode(device, !agc) != 0) {
    std::cerr << "Cannot enable / disable Tuner AGC" << std::endl;
    throw SatHelperException("Cannot set Tuner AGC");
  }

  refreshGains();

  if (rtlsdr_reset_buffer(device) != 0) {
    throw SatHelperException("Cannot reset device buffer");
  }
  mainThread = new std::thread(&RtlFrontend::threadWork, this);
}

void RtlFrontend::rtlCallback(unsigned char *data, unsigned int length, void *ctx) {
  RtlFrontend *frontend = (RtlFrontend *)ctx;
  frontend->internalCallback(data, length);
}

void RtlFrontend::threadWork() {
  rtlsdr_read_async(device, rtlCallback, this, 0, 16384);
}

void RtlFrontend::internalCallback(unsigned char *data, unsigned int length) {
  float *iq = new float[length];

  for (unsigned int i=0; i<length; i++) {
    iq[i] = lut[data[i]];
    if (i%1) {
      qavg += alpha * (iq[i] - qavg);
      iq[i] -= qavg;
    } else {
      iavg += alpha * (iq[i] - iavg);
      iq[i] -= iavg;
    }
  }

  if (cb != NULL) {
    this->cb->cbFloatIQ(iq, length/2);
  }

  delete[] iq;
}

void RtlFrontend::refreshGains() {
  switch (tuner) {
    case RTLSDR_TUNER_R820T: rtlsdr_set_tuner_gain_ext(device, lnaGain, mixerGain, vgaGain); break;
    case RTLSDR_TUNER_E4000:
      rtlsdr_set_tuner_gain(device, lnaGain);
//      rtlsdr_set_tuner_if_gain(device, 0, mixerGain);
//      rtlsdr_set_tuner_if_gain(device, 1, vgaGain);
//      rtlsdr_set_tuner_if_gain(device, 2, lnaGain);
      break;
    default:
      rtlsdr_set_tuner_gain(device, lnaGain);

  }
}

void RtlFrontend::Stop() {
  rtlsdr_cancel_async(device);
  if (mainThread != NULL && mainThread->joinable()) {
    mainThread->join();
  }
  mainThread = NULL;
}

void RtlFrontend::SetAGC(bool agc) {
  rtlsdr_set_tuner_gain_mode(device, !agc);
  this->agc = agc;
}

void RtlFrontend::SetLNAGain(uint8_t value) {
  lnaGain = value;
  refreshGains();
}

void RtlFrontend::SetVGAGain(uint8_t value) {
  vgaGain = value;
  refreshGains();
}

void RtlFrontend::SetMixerGain(uint8_t value) {
  mixerGain = value;
  refreshGains();
}

uint32_t RtlFrontend::GetCenterFrequency() {
  return centerFrequency;
}

const std::string &RtlFrontend::GetName() {
  return deviceName;
}

uint32_t RtlFrontend::GetSampleRate() {
  return sampleRate;
}

void RtlFrontend::SetSamplesAvailableCallback(GoDeviceCallback *cb) {
  this->cb = cb;
}
