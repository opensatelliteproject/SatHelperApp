/*
 * LimeDevice.h
 *
 *  Created on: 24/12/2016
 *      Author: Lucas Teske
 */

#ifndef SRC_LIMEDEVICE_H_
#define SRC_LIMEDEVICE_H_

#include <SoapySDR/Device.hpp>
#include <SoapySDR/Formats.hpp>
#include <SoapySDR/Errors.hpp>
#include <cstdint>
#include <iostream>
#include <sstream>
#include <vector>
#include <string>
#include <functional>

extern "C" {
//#include <SoapySDR/Device.h>
//#include <SoapySDR/Formats.h>
}

#include "../DeviceParameters.h"

class LimeDevice {
private:
	static std::string libraryVersion;
	static SoapySDR::Device* device;
	static SoapySDR::Stream* transfer;

	GoDeviceCallback *cb;
	uint8_t boardId;
	std::string firmwareVersion;
	std::string partNumber;
	std::string serialNumber;
	std::string name = "LimeSDR-Mini";

	uint32_t sampleRate;
	uint32_t centerFrequency;
	uint8_t gain;

	std::vector<uint32_t> availableSampleRates;
	int SamplesAvailableCallback();

public:
	LimeDevice();
	virtual ~LimeDevice();

	static void Initialize();
	static void DeInitialize();

	uint32_t SetSampleRate(uint32_t sampleRate);
	uint32_t SetCenterFrequency(uint32_t centerFrequency);
	const std::vector<uint32_t>& GetAvailableSampleRates();
	void Start();
	void Stop();
	void SetAGC(bool agc);
	void SetLNAGain(uint8_t value);
	void SetAntenna();
	uint32_t GetCenterFrequency();

	const std::string &GetName();

	uint32_t GetSampleRate();

	void SetSamplesAvailableCallback(GoDeviceCallback *cb);
};

#endif /* SRC_LIMEDEVICE_H_ */
