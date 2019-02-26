PACKAGE := OpenSatelliteProject/SatHelperApp
REV_VAR := main.RevString
VERSION_VAR := main.VersionString
BUILD_DATE_VAR := main.CompilationDate
BUILD_TIME_VAR := main.CompilationTime
REPO_VERSION := $(shell git describe --always --dirty --tags)
REPO_REV := $(shell git rev-parse --sq HEAD)
BUILD_DATE := $(shell date +"%b %d %Y")
BUILD_TIME := $(shell date +"%H:%M:%S")

PATH := $(PATH):/usr/lib/go-1.10/bin

GOBIN := $(shell PATH=$PATH:/usr/lib/go-1.10/bin:/usr/local/Cellar/go/1.10.2/bin command -v go 2> /dev/null)
BASEDIR := $(CURDIR)
GOPATH := $(CURDIR)/.gopath
BASE := $(GOPATH)/src/$(PACKAGE)
DESTDIR?=/usr/local/bin
GOBUILD_VERSION_ARGS := -ldflags "-X $(REV_VAR)=$(REPO_REV) -X $(VERSION_VAR)=$(REPO_VERSION) -X \"$(BUILD_DATE_VAR)=$(BUILD_DATE)\" -X $(BUILD_TIME_VAR)=$(BUILD_TIME)"

INTSIZE := $(shell getconf LONG_BIT)

.PHONY: all build
.NOTPARALLEL: pre deps update

all: | $(BASE) pre deps update build

$(BASE):
	@echo Linking virtual GOPATH
	@mkdir -p $(dir $@)
	@ln -sf $(CURDIR) $@

pre:
	@echo Prechecking
ifndef GOBIN
	$(error "GO executable not found")
endif

clean:
	@echo Cleaning virtual GOPATH
	@rm -fr .gopath

deps: | $(BASE)
	@echo Downloading dependencies
	@cd $(BASE) && GOPATH=$(GOPATH) $(GOBIN) get
	@cd $(BASE)/cmd/SatHelperApp && GOPATH=$(GOPATH) $(GOBIN) get
	@cd $(BASE)/cmd/demuxReplay && GOPATH=$(GOPATH) $(GOBIN) get
	@cd $(BASE)/cmd/xritparse && GOPATH=$(GOPATH) $(GOBIN) get

update: | $(BASE)
	@echo Updating AirspyDevice Wrapper
	@cd $(BASE) && swig -cgo -go -c++ -intgosize $(INTSIZE) Frontend/AirspyDevice/AirspyDevice.i
	@cd $(BASE) && swig -cgo -go -c++ -intgosize $(INTSIZE) Frontend/SpyserverDevice/SpyserverDevice.i

	@echo Updating LimeDevice Wrapper
	@cd $(BASE) && swig -cgo -go -c++ -intgosize $(INTSIZE) Frontend/LimeDevice/LimeDevice.i

	@echo Updating RTLSDR Wrappper
	@cd $(BASE) && swig -cgo -go -c++ -intgosize $(INTSIZE) Frontend/RTLSDRDevice/RTLSDRDevice.i

build: | $(BASE)
	@echo Building SatHelperApp
	@cd $(BASE)/cmd/SatHelperApp && GOPATH=$(GOPATH) $(GOBIN) build $(GOBUILD_VERSION_ARGS) -o $(BASEDIR)/SatHelperApp
	@echo Building DemuxReplay
	@cd $(BASE)/cmd/demuxReplay && GOPATH=$(GOPATH) $(GOBIN) build $(GOBUILD_VERSION_ARGS) -o $(BASEDIR)/DemuxReplay
	@echo Building xritparse
	@cd $(BASE)/cmd/xritparse && GOPATH=$(GOPATH) $(GOBIN) build $(GOBUILD_VERSION_ARGS) -o $(BASEDIR)/xritparse
	@echo Building xritcat
	@cd $(BASE)/cmd/xritcat && GOPATH=$(GOPATH) $(GOBIN) build $(GOBUILD_VERSION_ARGS) -o $(BASEDIR)/xritcat
	@echo Building xritimg
	@cd $(BASE)/cmd/xritimg && GOPATH=$(GOPATH) $(GOBIN) build $(GOBUILD_VERSION_ARGS) -o $(BASEDIR)/xritimg
	@echo Building xritpdcs
	@cd $(BASE)/cmd/xritpdcs && GOPATH=$(GOPATH) $(GOBIN) build $(GOBUILD_VERSION_ARGS) -o $(BASEDIR)/xritpdcs


install: | $(BASE)
	@echo Installing
	@cd $(BASE) && cp $(BASEDIR)/SatHelperApp $(DESTDIR)/SatHelperApp
	@chmod +x $(DESTDIR)/SatHelperApp
	@cd $(BASE) && cp $(BASEDIR)/SatHelperApp $(DESTDIR)/DemuxReplay
	@chmod +x $(DESTDIR)/DemuxReplay
	@cd $(BASE) && cp $(BASEDIR)/xritparse $(DESTDIR)/xritparse
	@chmod +x $(DESTDIR)/xritparse
	@cd $(BASE) && cp $(BASEDIR)/xritcat $(DESTDIR)/xritcat
	@chmod +x $(DESTDIR)/xritcat
	@cd $(BASE) && cp $(BASEDIR)/xritimg $(DESTDIR)/xritimg
	@chmod +x $(DESTDIR)/xritimg
	@cd $(BASE) && cp $(BASEDIR)/xritpdcs $(DESTDIR)/xritpdcs
	@chmod +x $(DESTDIR)/xritpdcs

test:
	go test -v -race $(shell go list ./... | grep -v /parts/ | grep -v /prime/ | grep -v /snap/ | grep -v /stage/ | grep -v /tmp/ | grep -v /librtlsdr/ )
