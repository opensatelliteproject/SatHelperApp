%module(directors="1") SpyserverDevice
%{
#include "SpyserverDevice.h"
%}

%insert(cgo_comment_typedefs) %{
#cgo CXXFLAGS: -std=c++11 -O0
#cgo LDFLAGS: -lSatHelper
%}

%include "stdint.i"
%include "stl.i"
%include "std_vector.i"

%feature("director") GoDeviceCallback;
%rename("SpyserverDeviceCallback") GoDeviceCallback;

%include "../DeviceParameters.h"

%template(Vector32u) std::vector<uint32_t>;
%template(Vector32f) std::vector<float>;
%template(Vector16i) std::vector<int16_t>;
%template(Vector8i) std::vector<int8_t>;

%include "./SpyserverDevice.h"