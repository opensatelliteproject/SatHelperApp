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

enum TLogLevel {logERROR, logWARN, logINFO, logDEBUG};

class GoDeviceCallback {
public:
    virtual void cbFloatIQ(void *data, int length) {}
    virtual void cbS16IQ(void *data, int length) {}
    virtual void cbS8IQ(void *data, int length) {}
    virtual void Info(std::string) {}
    virtual void Error(std::string) {}
    virtual void Warn(std::string) {}
    virtual void Debug(std::string) {}
    virtual ~GoDeviceCallback() {}
};


class Log {
private:
  GoDeviceCallback *cb;
  TLogLevel level;
protected:
  std::ostringstream os;
public:
  Log(GoDeviceCallback *cb) : cb(cb) {}
  std::ostringstream& Get(TLogLevel level = logINFO) {
    this->level = level;
    return os;
  }

  ~Log() {
    switch (this->level) {
      case logERROR: cb->Error(os.str()); break;
      case logDEBUG: cb->Debug(os.str()); break;
      case logWARN: cb->Warn(os.str()); break;
      default: cb->Info(os.str()); break;
    }
  }
};

#endif /* SRC_DEVICEPARAMETERS_H_ */