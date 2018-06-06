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