/*
 * LimeDevice.cpp
 */

#include "LimeDevice.h"
#include <SatHelper/exceptions/SatHelperException.h>

std::string LimeDevice::libraryVersion;
SoapySDR::Device* LimeDevice::device;
SoapySDR::Stream* LimeDevice::transfer;
SoapySDR::Kwargs LimeDevice::args;

uint32_t LimeDevice::GetCenterFrequency() {
	return device->getFrequency(SOAPY_SDR_RX, 0);
}

const std::string &LimeDevice::GetName() {
	return name;
}

uint32_t LimeDevice::GetSampleRate() {
	return device->getSampleRate(SOAPY_SDR_RX, 0);
}

void LimeDevice::SetSamplesAvailableCallback(GoDeviceCallback *cb) {
	this->cb = cb;
}

LimeDevice::~LimeDevice() {
	SoapySDR::Device::unmake(device);
}

bool LimeDevice::Init() {
	SoapySDR::setLogLevel(SOAPY_SDR_FATAL);
	args = SoapySDR::KwargsFromString("driver=lime,device=0");
    device = SoapySDR::Device::make(args);
	return (bool)device;
};

void LimeDevice::Destroy() {
	SoapySDR::Device::unmake(device);
};

void LimeDevice::Start() {
	transfer = device->setupStream(SOAPY_SDR_RX, SOAPY_SDR_CF32);
	if (device != NULL && transfer != NULL) {
		device->activateStream(transfer);
	} else {
		std::cerr << "Device not loaded!" << std::endl;
	}
}

void LimeDevice::Stop() {
	if (device != NULL && transfer != NULL) {
		device->deactivateStream(transfer);
		device->closeStream(transfer);
	} else {
		std::cerr << "Device not loaded!" << std::endl;
	}
}

void LimeDevice::SetLNAGain(uint8_t value) {
	device->setGain(SOAPY_SDR_RX, 0, "LNA", (double)value);
}

void LimeDevice::SetTIAGain(uint8_t value) {
	device->setGain(SOAPY_SDR_RX, 0, "TIA", (double)value);
}

void LimeDevice::SetPGAGain(uint8_t value) {
	device->setGain(SOAPY_SDR_RX, 0, "PGA", (double)value);
}

void LimeDevice::SetAntenna(std::string antenna) {
	device->setAntenna(SOAPY_SDR_RX, 0, antenna);
}

uint32_t LimeDevice::SetSampleRate(uint32_t sampleRate) {
	device->setSampleRate(SOAPY_SDR_RX, 0, (double)sampleRate);
	return device->getSampleRate(SOAPY_SDR_RX, 0);
}

uint32_t LimeDevice::SetCenterFrequency(uint32_t centerFrequency) {
	device->setFrequency(SOAPY_SDR_RX, 0, (double)centerFrequency);
	return device->getFrequency(SOAPY_SDR_RX, 0);
}

void LimeDevice::GetSamples(uint16_t samples) {
	if (cb != NULL) {
		int flags;
		long long timeNs;
		void* buffs[] = {buff};
		
		device->readStream(transfer, buffs, samples, flags, timeNs);
		cb->cbFloatIQ(buff, samples);
	}	
}