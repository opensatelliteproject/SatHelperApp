#!/bin/bash


REV_VAR="github.com/opensatelliteproject/SatHelperApp.RevString"
VERSION_VAR="github.com/opensatelliteproject/SatHelperApp.VersionString"
BUILD_DATE_VAR="github.com/opensatelliteproject/SatHelperApp.CompilationDate"
BUILD_TIME_VAR="github.com/opensatelliteproject/SatHelperApp.CompilationTime"
REPO_VERSION=$(git describe --always --dirty --tags)
REPO_REV=$(git rev-parse HEAD)
BUILD_DATE=$(date +"%d%m%Y")
BUILD_TIME=$(date +"%H%M%S")
GOBUILD_VERSION_ARGS="-ldflags \"-X '${REV_VAR}=${REPO_REV}' -X '${VERSION_VAR}=${REPO_VERSION}' -X '${BUILD_DATE_VAR}=${BUILD_DATE}' -X '${BUILD_TIME_VAR}=${BUILD_TIME}'\""


echo "REV_VAR=${REV_VAR}"
echo "VERSION_VAR=${VERSION_VAR}"
echo "BUILD_DATE_VAR=${BUILD_DATE_VAR}"
echo "BUILD_TIME_VAR=${BUILD_TIME_VAR}"
echo "REPO_VERSION=${REPO_VERSION}"
echo "REPO_REV=${REPO_REV}"
echo "BUILD_DATE=${BUILD_DATE}"
echo "BUILD_TIME=${BUILD_TIME}"
echo "GOBUILD_VERSION_ARGS=${GOBUILD_VERSION_ARGS}"


TAG=`git describe --exact-match --tags HEAD 2>/dev/null`
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
    echo go build ${GOBUILD_VERSION_ARGS} -o ../../bins/${i}
    bash -c "go build ${GOBUILD_VERSION_ARGS} -o ../../bins/${i}"
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
