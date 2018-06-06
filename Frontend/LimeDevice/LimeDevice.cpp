/*
 * LimeDevice.cpp
 */

#include "LimeDevice.h"
#include <SatHelper/exceptions/SatHelperException.h>

std::string LimeDevice::libraryVersion;
SoapySDR::Device* LimeDevice::device;
SoapySDR::Stream* LimeDevice::transfer;

uint32_t LimeDevice::GetCenterFrequency() {
	return centerFrequency;
}

const std::string &LimeDevice::GetName() {
	return name;
}

uint32_t LimeDevice::GetSampleRate() {
	return sampleRate;
}

void LimeDevice::SetSamplesAvailableCallback(GoDeviceCallback *cb) {
	this->cb = cb;
}

const std::vector<uint32_t>& LimeDevice::GetAvailableSampleRates() {
	return availableSampleRates;
}

LimeDevice::LimeDevice() {
	SoapySDR::Kwargs args = SoapySDR::KwargsFromString("driver=lime,device=0");
    device = SoapySDR::Device::make(args);

    if (device == NULL) {
        throw SatHelperException("There was an error initializing LimeSDR.");
    }

	if (device != NULL) {
		SetLNAGain(50);
		SetAntenna();
	}
}

LimeDevice::~LimeDevice() {
	SoapySDR::Device::unmake(device);
}

void LimeDevice::Initialize() {
	transfer = device->setupStream(SOAPY_SDR_RX, SOAPY_SDR_CF32);
}

void LimeDevice::DeInitialize() {
	SoapySDR::Device::unmake(device);
}

void LimeDevice::Start() {
	Initialize();
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
		DeInitialize();
	} else {
		std::cerr << "Device not loaded!" << std::endl;
	}
}

void LimeDevice::SetAGC(bool agc) {}

void LimeDevice::SetLNAGain(uint8_t value) {
	device->setGain(SOAPY_SDR_RX, 0, (double)value);
}

void LimeDevice::SetAntenna() {
	device->setAntenna(SOAPY_SDR_RX, 0, "LNAH");
}

uint32_t LimeDevice::SetSampleRate(uint32_t sampleRate) {
	device->setSampleRate(SOAPY_SDR_RX, 0, (double)sampleRate);
	return this->sampleRate = device->getSampleRate(SOAPY_SDR_RX, 0);
}

uint32_t LimeDevice::SetCenterFrequency(uint32_t centerFrequency) {
	device->setFrequency(SOAPY_SDR_RX, 0, (double)centerFrequency);
	return this->centerFrequency = device->getFrequency(SOAPY_SDR_RX, 0);
}