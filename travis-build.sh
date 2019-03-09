#!/bin/bash

TAG=`git describe --exact-match --tags HEAD`

if [[ $? -eq 0 ]];
then
  echo "Releasing for tag ${TAG}"
  ORIGINAL_FOLDER="`pwd`"
  echo "I'm in `pwd`"
  mkdir -p bins
  mkdir -p zips

  echo "Building RTLSDR"
  git clone https://github.com/librtlsdr/librtlsdr.git
  cd librtlsdr
  mkdir -p build && cd build
  cmake ..
  make -j10
  sudo make install
  sudo ldconfig
  cd ..

  echo "Going back to $ORIGINAL_FOLDER"
  cd "$ORIGINAL_FOLDER"

  echo "Building Static LimeSuite"
  git clone https://github.com/myriadrf/LimeSuite.git
  cd LimeSuite
  git checkout stable
  mkdir -p builddir && cd builddir
  cmake ../ -DBUILD_SHARED_LIBS=OFF
  make -j10
  sudo make install
  sudo ldconfig
  cd ..

  echo "Going back to $ORIGINAL_FOLDER"
  cd "$ORIGINAL_FOLDER"

  echo "Updating Code to have static libLimeSuite"
  sed -i 's/-lLimeSuite/-l:libLimeSuite.a -l:libstdc++.a -lm -lusb-1.0/g' Frontend/LimeDevice/LimeDevice.go
  sed -i 's/-lLimeSuite/-l:libLimeSuite.a -l:libstdc++.a -lm -lusb-1.0/g' ../../myriadrf/limedrv/limewrap/limewrap.go

  echo "Building"
  cd cmd
  for i in *
  do
    echo "Building $i"
    cd ${i}
    echo go build -o ../../bins/${i}
    go build -o ../../bins/${i}
    echo "Zipping ${i}-${TAG}-linux-amd64.zip"
    zip -r "../../zips/${i}-${TAG}-linux-amd64.zip" ../../bins/$i
    cd ..
  done
  cd ..
  echo "Binaries: "
  ls -la bins
  echo "Zip Files: "
  ls -la zips
else
  echo "No tags for current commit. Skipping releases."
fi
