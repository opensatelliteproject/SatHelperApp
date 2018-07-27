/*
 * LimeDevice.cpp
 */

#include "LimeDevice.h"

lms_device_t* LimeDevice::device;
lms_stream_t LimeDevice::transfer;

uint32_t LimeDevice::GetCenterFrequency() {
	double center_frequency;
	LMS_GetLOFrequency(device, LMS_CH_RX, 0, &center_frequency);
	return (uint32_t)center_frequency;
}

const std::string LimeDevice::GetName() {
	return name;
}

uint32_t LimeDevice::GetSampleRate() {
	double sample_rate;
	LMS_GetSampleRate(device, LMS_CH_RX, 0, &sample_rate, NULL);
	return (uint32_t)sample_rate;
}

void LimeDevice::SetSamplesAvailableCallback(GoDeviceCallback *cb) {
	this->cb = cb;
}

LimeDevice::~LimeDevice() {
	LMS_Close(device);
}

void LimeDevice::Destroy() {
	LMS_Close(device);
};

bool LimeDevice::Init() {
    lms_info_str_t list[8];
	
    if (LMS_GetDeviceList(list) < 0) {
		Log(cb).Get(logERROR) << "Error getting device list.";
		return false;
	}
        
    if (LMS_Open(&device, list[0], NULL)) {
		Log(cb).Get(logERROR) << "Error opening device. Check permissions.";
		return false;
	}

    if (LMS_Init(device) != 0) {
		Log(cb).Get(logERROR) << "Error initializing device.";
		return false;
	}

	if (LMS_SetLPFBW(device, LMS_CH_RX, 0, 3e6) != 0) {
		Log(cb).Get(logERROR) << "Error setting analog filter.";
		return false;
	}

	if (LMS_Calibrate(device, LMS_CH_RX, 0, 3e6, 0) != 0) {
		Log(cb).Get(logERROR) << "Error calibrating device. Maybe the port is saturated.";
		return false;
	}
    	
	return (bool)device;
};

void LimeDevice::Start() {
    transfer.channel = 0;
    transfer.fifoSize = 1024 * 1024;
    transfer.throughputVsLatency = 1.0;
    transfer.isTx = false;
    transfer.dataFmt = lms_stream_t::LMS_FMT_F32;
    
    if (LMS_SetupStream(device, &transfer) != 0) {
		Log(cb).Get(logERROR) << "Error starting streaming.";
		return;
	}

	LMS_StartStream(&transfer);
}

void LimeDevice::Stop() {
	if (device != NULL && &transfer != NULL) {
		LMS_StopStream(&transfer);
    	LMS_DestroyStream(device, &transfer);
	} else {
		Log(cb).Get(logERROR) << "Device not loaded!";
	}
}

void LimeDevice::SetLNAGain(uint8_t value) {
	if (LMS_SetNormalizedGain(device, LMS_CH_RX, 0, (double)value/50) != 0)	
		Log(cb).Get(logERROR) << "Error setting gain.";
}

void LimeDevice::SetTIAGain(uint8_t value) {
	Log(cb).Get(logWARN) << "Unfortunately TIA Gain isn't supported by the LMS7 API.";
}

void LimeDevice::SetPGAGain(uint8_t value) {
	Log(cb).Get(logWARN) << "Unfortunately PGA Gain isn't supported by the LMS7 API.";
}

void LimeDevice::SetAntenna(std::string wishedAntenna) {
	int antennaCode = LMS_PATH_LNAH;

	switch ((char)wishedAntenna[3]) {
		case 'H': antennaCode = LMS_PATH_LNAH; break;
		case 'W': antennaCode = LMS_PATH_LNAW; break;
		case 'L': antennaCode = LMS_PATH_LNAL; break;
		default:
			Log(cb).Get(logWARN) << "Invalid Antenna: Falling back to LNAH. Check your settings!";
			antennaCode = LMS_PATH_LNAH;
			break;
	}

	if (LMS_SetAntenna(device, LMS_CH_RX, 0, antennaCode) != 0) {
		Log(cb).Get(logERROR) << "Error setting antenna!";
	}
}

uint32_t LimeDevice::SetSampleRate(uint32_t wishedSampleRate) {
	double actualSampleRate;

	if (LMS_SetSampleRate(device, (double)wishedSampleRate, 4) != 0) {
		Log(cb).Get(logERROR) << "Error setting sample rate.";
	}
	
	LMS_GetSampleRate(device, LMS_CH_RX, 0, &actualSampleRate, NULL);
	return (uint32_t)actualSampleRate;
}

uint32_t LimeDevice::SetCenterFrequency(uint32_t wishedFrequency) {
	double actualFrequency;

	if (LMS_SetLOFrequency(device, LMS_CH_RX, 0, wishedFrequency) != 0) {
		Log(cb).Get(logERROR) << "Error setting frequency!";
	}

	LMS_GetLOFrequency(device, LMS_CH_RX, 0, &actualFrequency);
	return (uint32_t)actualFrequency;
}

void LimeDevice::GetSamples(uint16_t samples) {
	if (cb != NULL) {
		uint16_t samplesRead = LMS_RecvStream(&transfer, buff, samples, NULL, 1000);
		cb->cbFloatIQ(buff, samplesRead);
	}
}
