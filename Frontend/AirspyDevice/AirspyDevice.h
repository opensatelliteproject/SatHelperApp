/*
 * AirspyDevice.h
 *
 *  Created on: 24/12/2016
 *      Author: Lucas Teske
 */

#ifndef SRC_AIRSPYDEVICE_H_
#define SRC_AIRSPYDEVICE_H_

#include <cstdint>
#include <iostream>
#include <sstream>
#include <vector>
#include <string>
#include <functional>

extern "C" {
#include <libairspy/airspy.h>
}

#include "../DeviceParameters.h"
class AirspyDevice {
private:
	static std::string libraryVersion;

	GoDeviceCallback *cb;
	uint8_t boardId;
	std::string firmwareVersion;
	std::string partNumber;
	std::string serialNumber;
	std::vector<uint32_t> availableSampleRates;
	std::string name;
	airspy_device* device;

	uint32_t sampleRate;
	uint32_t centerFrequency;
	uint8_t lnaGain;
	uint8_t vgaGain;
	uint8_t mixerGain;

	int SamplesAvailableCallback(airspy_transfer *transfer);
public:
	AirspyDevice(GoDeviceCallback *cb);
	virtual ~AirspyDevice();

	static void Initialize();
	static void DeInitialize();

	uint32_t SetSampleRate(uint32_t sampleRate);
	uint32_t SetCenterFrequency(uint32_t centerFrequency);
	const std::vector<uint32_t>& GetAvailableSampleRates();
	void Start();
	void Stop();
	void SetAGC(bool agc);

	bool Init();
	void Destroy();

	void SetLNAGain(uint8_t value);
	void SetVGAGain(uint8_t value);
	void SetMixerGain(uint8_t value);
	void SetBiasT(uint8_t value);

	uint32_t GetCenterFrequency();

	const std::string &GetName();

	uint32_t GetSampleRate();

	void SetSamplesAvailableCallback(GoDeviceCallback *cb);

};

#endif /* SRC_AIRSPYDEVICE_H_ */
