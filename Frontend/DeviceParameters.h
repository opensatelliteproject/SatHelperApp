/*
 * DeviceParameters.h
 *
 *  Created on: 31/01/2017
 *      Author: Lucas Teske
 */

#ifndef SRC_DEVICEPARAMETERS_H_
#define SRC_DEVICEPARAMETERS_H_


#define FRONTEND_SAMPLETYPE_FLOATIQ 0
#define FRONTEND_SAMPLETYPE_S16IQ 1
#define FRONTEND_SAMPLETYPE_S8IQ 2

class GoDeviceCallback {
public:
    virtual void cbFloatIQ(void *data, int length) {}
    virtual void cbS16IQ(int16_t *data, int length) {}
    virtual void cbS8IQ(int8_t *data, int length) {}
    virtual ~GoDeviceCallback() {}
};

#endif /* SRC_DEVICEPARAMETERS_H_ */