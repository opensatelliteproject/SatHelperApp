PACKAGE := OpenSatelliteProject/SatHelperApp
REV_VAR := main.RevString
VERSION_VAR := main.VersionString
BUILD_DATE_VAR := main.CompilationDate
BUILD_TIME_VAR := main.CompilationTime
REPO_VERSION := $(shell git describe --always --dirty --tags)
REPO_REV := $(shell git rev-parse --sq HEAD)
BUILD_DATE := $(shell date +"%b %d %Y")
BUILD_TIME := $(shell date +"%H:%M:%S")

BASEDIR := $(CURDIR)
GOPATH := $(CURDIR)/.gopath
BASE := $(GOPATH)/src/$(PACKAGE)
DESTDIR?=/usr/local/bin 
GOBUILD_VERSION_ARGS := -ldflags "-X $(REV_VAR)=$(REPO_REV) -X $(VERSION_VAR)=$(REPO_VERSION) -X \"$(BUILD_DATE_VAR)=$(BUILD_DATE)\" -X $(BUILD_TIME_VAR)=$(BUILD_TIME)"

.PHONY: all
.NOTPARALLEL: deps update

all: | $(BASE) deps update build

$(BASE):
	@echo Linking virtual GOPATH
	@mkdir -p $(dir $@)
	@ln -sf $(CURDIR) $@

clean:
	@echo Cleaning virtual GOPATH
	@rm -fr .gopath

deps: | $(BASE)
	@echo Downloading dependencies
	@cd $(BASE) && go get

update: | $(BASE)
	@echo Updating AirspyDevice Wrapper
	@cd $(BASE) && swig -cgo -go -c++ -intgosize 64 Frontend/AirspyDevice/AirspyDevice.i
	@cd $(BASE) && swig -cgo -go -c++ -intgosize 64 Frontend/SpyserverDevice/SpyserverDevice.i

	@echo Updating LimeDevice Wrapper
	@cd $(BASE) && swig -cgo -go -c++ -intgosize 64 Frontend/LimeDevice/LimeDevice.i

build: | $(BASE)
	@echo Building SatHelperApp
	@cd $(BASE) && go build $(GOBUILD_VERSION_ARGS) -o $(BASEDIR)/SatHelperApp

install: | $(BASE)
	@echo Installing
	@cd $(BASE) && cp $(BASEDIR)/SatHelperApp $(DESTDIR)/SatHelperApp
	@chmod +x $(DESTDIR)/SatHelperApp
