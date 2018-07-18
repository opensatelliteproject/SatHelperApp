/*
 * LimeDevice.h
 *
 *  Created on: 24/12/2016
 *      Author: Lucas Teske
 */

#ifndef SRC_LIMEDEVICE_H_
#define SRC_LIMEDEVICE_H_

#include <cstdint>
#include <iostream>
#include <sstream>
#include <vector>
#include <string>
#include <functional>

#include <lime/LimeSuite.h>
#include <SatHelper/exceptions/SatHelperException.h>

#include "../DeviceParameters.h"

class LimeDevice {
private:
	static std::string libraryVersion;
	static lms_device_t* device;
	static lms_stream_t transfer;

	GoDeviceCallback *cb;
	std::string name = "LimeSDR (LMS7)";
	float buff[16384];

public:
	LimeDevice(GoDeviceCallback *cb) : cb(cb) {};
	virtual ~LimeDevice();

	const std::string GetName();
	
	uint32_t SetSampleRate(uint32_t wishedSampleRate);
	uint32_t SetCenterFrequency(uint32_t wishedFrequency);

	uint32_t GetCenterFrequency();
	uint32_t GetSampleRate();

	void SetLNAGain(uint8_t value);
	void SetTIAGain(uint8_t value);
	void SetPGAGain(uint8_t value);
	void SetAntenna(std::string antenna);
	void SetSamplesAvailableCallback(GoDeviceCallback *cb);

	bool Init();
	void Destroy();

	void Start();
	void Stop();
	void GetSamples(uint16_t samples);
};

#endif /* SRC_LIMEDEVICE_H_ */
