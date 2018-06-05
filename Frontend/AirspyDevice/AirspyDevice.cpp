/*
 * AirspyDevice.cpp
 *
 *  Created on: 24/12/2016
 *      Author: Lucas Teske
 */

#include "AirspyDevice.h"
#include <SatHelper/exceptions/SatHelperException.h>

#define AIRSPY_ERROR airspy_error_name((enum airspy_error)result)

std::string AirspyDevice::libraryVersion;

uint32_t AirspyDevice::GetCenterFrequency() {
	return centerFrequency;
}

const std::string &AirspyDevice::GetName() {
	return name;
}

uint32_t AirspyDevice::GetSampleRate() {
	return sampleRate;
}

void AirspyDevice::SetSamplesAvailableCallback(GoDeviceCallback *cb) {
	this->cb = cb;
}


void AirspyDevice::SetBiasT(uint8_t value) {
	int result = airspy_set_rf_bias(device, value);
	if (result != AIRSPY_SUCCESS) {
		std::cerr << "Error setting BiasT to " << value << ": " << AIRSPY_ERROR << std::endl;
	}
}

AirspyDevice::AirspyDevice() {
	int result = airspy_open(&device);
	if (result != AIRSPY_SUCCESS) {
		std::cerr << "Cannot open device: " << AIRSPY_ERROR << std::endl;
		throw SatHelperException("There was an error initializing AirSpy.");
	}

	if (device != NULL) {
		// Get Board ID
		result = airspy_board_id_read(device, &boardId);
		if (result != AIRSPY_SUCCESS) {
			std::cerr << "Cannot get BoardId: " << AIRSPY_ERROR << std::endl;
			throw SatHelperException("There was an error initializing AirSpy.");
		}
		// std::cout << "  Device Board ID: " << (int) boardId << std::endl;

		// Get Firmware Version
		char version[256];
		result = airspy_version_string_read(device, version, 255);
		if (result != AIRSPY_SUCCESS) {
			std::cerr << "Cannot get Firmware Version: " << AIRSPY_ERROR
					<< std::endl;
			throw SatHelperException("There was an error initializing AirSpy.");
		}
		firmwareVersion = std::string(version);
		// std::cout << "  Firmware Version: " << firmwareVersion << std::endl;

		// Get Serial Number and Part Number
		airspy_read_partid_serialno_t ser;
		result = airspy_board_partid_serialno_read(device, &ser);
		if (result != AIRSPY_SUCCESS) {
			std::cerr << "Cannot get Serial Data: " << AIRSPY_ERROR
					<< std::endl;
			throw SatHelperException("There was an error initializing AirSpy.");
		}

		std::stringstream ss;
		ss << "0x" << std::hex << ser.part_id[0] << " 0x" << std::hex
				<< ser.part_id[1];
		partNumber = ss.str();
		// std::cout << "  Part Number: " << partNumber << std::endl;

		ss.str( std::string() );
		ss.clear();

		ss << "0x" << std::hex << ser.serial_no[2] << std::hex
				<< ser.serial_no[3];
		serialNumber = ss.str();
		// std::cout << "  Serial Number: " << serialNumber << std::endl;

		uint32_t sampRateCount;
		uint32_t *sampleRates;
		airspy_get_samplerates(device, &sampRateCount, 0);

		sampleRates = new uint32_t[sampRateCount];
		airspy_get_samplerates(device, sampleRates, sampRateCount);
		for (unsigned i = 0; i < sampRateCount; i++) {
			if (i == 0) {
				// std::cout << "Setting default sample rate as " << (sampleRates[i] * 1e-6) << " MSPS" << std::endl;
				sampleRate = sampleRates[i];
			}
			availableSampleRates.push_back(sampleRates[i]);
		}
		delete sampleRates;
		ss.str( std::string() );
		ss.clear();
		ss << "AirSpy(" << (int) boardId << ") - " << ser.serial_no[2]
				<< ser.serial_no[3];
		name = ss.str();
		// std::cout << "Device Name: " << name << std::endl;

		SetCenterFrequency(106300000);
		result = airspy_set_sample_type(device, AIRSPY_SAMPLE_FLOAT32_IQ);
		if (result != AIRSPY_SUCCESS) {
			std::cerr << "Cannot set Sample Type: " << AIRSPY_ERROR
					<< std::endl;
			throw SatHelperException("There was an error initializing AirSpy.");
		}

		SetSampleRate(availableSampleRates[0]);

		SetLNAGain(8);
		SetMixerGain(5);
		SetVGAGain(5);
	}
}

AirspyDevice::~AirspyDevice() {
	int result = airspy_close(device);
	if (result != AIRSPY_SUCCESS) {
		std::cerr << "Cannot close device: " << AIRSPY_ERROR << std::endl;
	}
}

void AirspyDevice::Initialize() {
	int result = airspy_init();
	if (result != AIRSPY_SUCCESS) {
		std::cerr << "Error initializing Airspy Library: " << AIRSPY_ERROR
				<< std::endl;
		return;
	}
	airspy_lib_version_t lib_version;
	airspy_lib_version(&lib_version);
	std::stringstream ss;
	ss << "Airspy Library " << lib_version.major_version << "."
			<< lib_version.minor_version << "." << lib_version.minor_version;
	AirspyDevice::libraryVersion = ss.str();
	// std::cout << "AirSpy Library version: " << libraryVersion << std::endl;
}

void AirspyDevice::DeInitialize() {
	airspy_exit();
}

void AirspyDevice::SetAGC(bool agc) {
	if (agc) {
		airspy_set_lna_agc(device, 1);
		airspy_set_mixer_agc(device, 1);
	} else {
		airspy_set_lna_agc(device, 0);
		airspy_set_mixer_agc(device, 0);
		SetLNAGain(lnaGain);
		SetMixerGain(mixerGain);
	}
}

void AirspyDevice::SetLNAGain(uint8_t value) {
	int result = airspy_set_lna_gain(device, value) ;
	if (result != AIRSPY_SUCCESS) {
		std::cout << "Error setting LNA Gain to " << (int) value << ": " << AIRSPY_ERROR << std::endl;
	}
	lnaGain = value;
}

void AirspyDevice::SetVGAGain(uint8_t value) {
	int result = airspy_set_vga_gain(device, value);
	if (result != AIRSPY_SUCCESS) {
		std::cout << "Error setting VGA Gain to " << (int) value << ": " << AIRSPY_ERROR << std::endl;
	}
	vgaGain = value;
}

void AirspyDevice::SetMixerGain(uint8_t value) {
	int result = airspy_set_mixer_gain(device, value);
	if (result != AIRSPY_SUCCESS) {
		std::cout << "Error setting Mixer Gain to " << (int) value << ": " << AIRSPY_ERROR << std::endl;
	}
	mixerGain = value;
}

const std::vector<uint32_t>& AirspyDevice::GetAvailableSampleRates() {
	return availableSampleRates;
}

void AirspyDevice::Start() {
	if (device != NULL) {
		// std::cout << "Starting Airspy Frontend" << std::endl;
		int result = airspy_start_rx(device, [](airspy_transfer* transfer) {
			return ((AirspyDevice*) transfer->ctx)->SamplesAvailableCallback(transfer);
		}, this);
		if (result != AIRSPY_SUCCESS) {
			std::cerr << "There was a problem starting AirSpy Streaming." << std::endl;
		}
	} else {
		std::cerr << "Device not loaded!" << std::endl;
	}
}

void AirspyDevice::Stop() {
	if (device != NULL && airspy_is_streaming(device) == AIRSPY_TRUE) {
		airspy_stop_rx(device);
	} else {
		std::cerr << "Device not loaded!" << std::endl;
	}
}

uint32_t AirspyDevice::SetSampleRate(uint32_t sampleRate) {
	if (this->sampleRate != sampleRate) {
		if (airspy_is_streaming(device) == AIRSPY_TRUE) {
			Stop();
			int result = airspy_set_samplerate(device, sampleRate);
			if (result != AIRSPY_SUCCESS) {
				std::cerr << "Cannot change device sample rate: "
						<< AIRSPY_ERROR << std::endl;
			} else {
				this->sampleRate = sampleRate;
			}
			Start();
		} else {
			int result = airspy_set_samplerate(device, sampleRate);
			if (result != AIRSPY_SUCCESS) {
				std::cerr << "Cannot change device sample rate: "
						<< AIRSPY_ERROR << std::endl;
			} else {
				this->sampleRate = sampleRate;
			}
		}
	}
	return this->sampleRate;
}

uint32_t AirspyDevice::SetCenterFrequency(uint32_t centerFrequency) {
	centerFrequency = centerFrequency < 24000000 ? 24000000 : centerFrequency;
	centerFrequency = centerFrequency > 1750000000 ? 1750000000 : centerFrequency;
	int result = airspy_set_freq(device, centerFrequency);
	if (result != AIRSPY_SUCCESS) {
		std::cerr << "Cannot change device center frequency: " << AIRSPY_ERROR << std::endl;
	} else {
		this->centerFrequency = centerFrequency;
	}
	return this->centerFrequency;
}

int AirspyDevice::SamplesAvailableCallback(airspy_transfer * transfer) {
	if (transfer->dropped_samples > 0) {
		// std::cerr << "Warning! " << transfer->dropped_samples << " samples has been dropped!" << std::endl;
	}
	if (cb != NULL) {
	    cb->cbFloatIQ(transfer->samples, transfer->sample_count);
	}
	return 0;
}