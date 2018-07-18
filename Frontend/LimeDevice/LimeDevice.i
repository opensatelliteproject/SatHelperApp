%module(directors="1") LimeDevice
%{
#include "LimeDevice.h"
%}

%insert(cgo_comment_typedefs) %{
#cgo CXXFLAGS: -std=c++11 -O3
#cgo LDFLAGS: -lLimeSuite
%}

%include "stdint.i"
%include "stl.i"
%include "std_vector.i"

%feature("director") GoDeviceCallback;
%rename(LimeDeviceCallback) GoDeviceCallback;

%include "../DeviceParameters.h"

%template(Vector32u) std::vector<uint32_t>;
%template(Vector32f) std::vector<float>;
%template(Vector16i) std::vector<int16_t>;
%template(Vector8i) std::vector<int8_t>;

%include "./LimeDevice.h"
