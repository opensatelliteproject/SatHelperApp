[![Snap Status](https://build.snapcraft.io/badge/opensatelliteproject/SatHelperApp.svg)](https://build.snapcraft.io/user/opensatelliteproject/SatHelperApp)


SatHelperApp
============

The OpenSatelliteProject Satellite Helper Application! This is currently a LRIT/HRIT Demodulator / Decoder program based on [libSatHelper](https://github.com/opensatelliteproject/libsathelper/) and [xritdemod](https://github.com/opensatelliteproject/xritdemod).

It is currently WIP and in Alpha State. Use with care.


Building
========

That's a standard go project. Make sure you have `libSatHelper` and `libairspy` installed and run:

```bash
go get github.com/OpenSatelliteProject/SatHelperApp
```

It will be installed into your `${GOPATH}/bin`. If you have it it on your path, just run `SatHelperApp` 

Have fun!


Ubuntu Instructions to get it running
=====================================

Base tools:
```bash
sudo apt install build-essential cmake swig
```

Quick Instructions to get GO 1.10 running:

```bash
sudo add-apt-repository ppa:gophers/archive
sudo apt-get update
sudo apt-get install golang-1.10-go
mkdir ~/go
export GOPATH=~/go
export GOROOT=/usr/lib/go-1.10
export PATH=$PATH:$GOPATH/bin:$GOROOT/bin
```

Install LibSatHelper:
```bash
git clone https://github.com/opensatelliteproject/libsathelper/
cd libsathelper
make libcorrect
sudo make libcorrect-install
make
sudo make install
```

Install libAirspy:
```bash
git clone https://github.com/airspy/airspyone_host/
cd airspyone_host
mkdir build
cd build
cmake .. -DINSTALL_UDEV_RULES=ON
make -j4
sudo make install
sudo ldconfig
```

Quick Instructions to get SatHelperApp running assuming you have Go 1.10, libSatHelper and libAirspy installed.
```bash
go get github.com/OpenSatelliteProject/SatHelperApp
SatHelperApp
```

This should create a `SatHelperApp.cfg` file in the same folder you ran `SatHelperApp`. You can edit it and tune for your needs.

Have fun!
