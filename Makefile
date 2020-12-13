PACKAGE := opensatelliteproject/SatHelperApp
REV_VAR := github.com/opensatelliteproject/SatHelperApp.RevString
VERSION_VAR := github.com/opensatelliteproject/SatHelperApp.VersionString
BUILD_DATE_VAR := github.com/opensatelliteproject/SatHelperApp.CompilationDate
BUILD_TIME_VAR := github.com/opensatelliteproject/SatHelperApp.CompilationTime
REPO_VERSION := $(shell git describe --always --dirty --tags)
REPO_REV := $(shell git rev-parse --sq HEAD)
BUILD_DATE := $(shell date +"%b %d %Y")
BUILD_TIME := $(shell date +"%H:%M:%S")

PATH := $(PATH):/usr/lib/go-1.11/bin

GOBIN := $(shell PATH=$PATH:/usr/lib/go-1.14/bin:/usr/local/Cellar/go/1.14/bin command -v go 2> /dev/null)
BASEDIR := $(CURDIR)
GOPATH := $(CURDIR)/.gopath
BASE := $(GOPATH)/src/$(PACKAGE)
DESTDIR?=/usr/local/bin
GOBUILD_VERSION_ARGS := -ldflags "-X $(REV_VAR)=$(REPO_REV) -X $(VERSION_VAR)=$(REPO_VERSION) -X \"$(BUILD_DATE_VAR)=$(BUILD_DATE)\" -X $(BUILD_TIME_VAR)=$(BUILD_TIME)"

INTSIZE := $(shell getconf LONG_BIT)

.PHONY: all build $(BASE)
.NOTPARALLEL: pre deps update

all: | $(BASE) pre deps update build

$(BASE):
	@echo Linking virtual GOPATH
	@mkdir -p "$(dir $@)"
	@sudo mount --bind $(CURDIR) "$(dir $@)"

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
	@cd $(BASE) && GO111MODULE=on GOPATH=$(GOPATH) $(GOBIN) get ./...
	@echo Deps for SatHelperApp
	@cd $(BASE)/cmd/SatHelperApp && GO111MODULE=on GOPATH=$(GOPATH) $(GOBIN) get ./...
	@echo Deps for DemuxReplay
	@cd $(BASE)/cmd/demuxReplay && GO111MODULE=on GOPATH=$(GOPATH) $(GOBIN) get ./...
	@echo Deps for xritparse
	@cd $(BASE)/cmd/xritparse && GO111MODULE=on GOPATH=$(GOPATH) $(GOBIN) get ./...
	@echo Deps for xritcat
	@cd $(BASE)/cmd/xritcat && GO111MODULE=on GOPATH=$(GOPATH) $(GOBIN) get ./...
	@echo Deps for xritimg
	@cd $(BASE)/cmd/xritimg && GO111MODULE=on GOPATH=$(GOPATH) $(GOBIN) get ./...
	@echo Deps for xritpdcs
	@cd $(BASE)/cmd/xritpdcs && GO111MODULE=on GOPATH=$(GOPATH) $(GOBIN) get ./...
	@echo Deps for MultiSegmentDump
	@cd $(BASE)/cmd/MultiSegmentDump && GO111MODULE=on GOPATH=$(GOPATH) $(GOBIN) get ./...
	@echo Deps for rpcClient
	@cd $(BASE)/cmd/rpcClient && GO111MODULE=on GOPATH=$(GOPATH) $(GOBIN) get ./...
	@echo Deps for SatUI
	@cd $(BASE)/cmd/SatUI && GO111MODULE=on GOPATH=$(GOPATH) $(GOBIN) get ./...

do-static: | $(BASE)
	@echo "Updating Code to have static libLimeSuite"
	@sed -i 's/-lLimeSuite/-l:libLimeSuite.a -l:libstdc++.a -static-libgcc -lm -lusb-1.0/g' $(GOPATH)/pkg/mod/github.com/myriadrf/limedrv*/limewrap/limewrap.go

update: | do-static $(BASE)
	@echo Updating AirspyDevice Wrapper
	@cd $(BASE) && swig -cgo -go -c++ -intgosize $(INTSIZE) Frontend/AirspyDevice/AirspyDevice.i

	@echo Updating RTLSDR Wrappper
	@cd $(BASE) && swig -cgo -go -c++ -intgosize $(INTSIZE) Frontend/RTLSDRDevice/RTLSDRDevice.i

build: | $(BASE)
	@echo Building SatHelperApp
	@cd $(BASE)/cmd/SatHelperApp && GO111MODULE=on GOPATH=$(GOPATH) $(GOBIN) build $(GOBUILD_VERSION_ARGS) -o $(BASEDIR)/SatHelperApp

	@echo Building DemuxReplay
	@cd $(BASE)/cmd/demuxReplay && GO111MODULE=on GOPATH=$(GOPATH) $(GOBIN) build $(GOBUILD_VERSION_ARGS) -o $(BASEDIR)/DemuxReplay

	@echo Building xritparse
	@cd $(BASE)/cmd/xritparse && GO111MODULE=on GOPATH=$(GOPATH) $(GOBIN) build $(GOBUILD_VERSION_ARGS) -o $(BASEDIR)/xritparse

	@echo Building xritcat
	@cd $(BASE)/cmd/xritcat && GO111MODULE=on GOPATH=$(GOPATH) $(GOBIN) build $(GOBUILD_VERSION_ARGS) -o $(BASEDIR)/xritcat

	@echo Building xritimg
	@cd $(BASE)/cmd/xritimg && GO111MODULE=on GOPATH=$(GOPATH) $(GOBIN) build $(GOBUILD_VERSION_ARGS) -o $(BASEDIR)/xritimg

	@echo Building xritpdcs
	@cd $(BASE)/cmd/xritpdcs && GO111MODULE=on GOPATH=$(GOPATH) $(GOBIN) build $(GOBUILD_VERSION_ARGS) -o $(BASEDIR)/xritpdcs

	@echo Building MultiSegmentDump
	@cd $(BASE)/cmd/MultiSegmentDump && GO111MODULE=on GOPATH=$(GOPATH) $(GOBIN) build $(GOBUILD_VERSION_ARGS) -o $(BASEDIR)/SatHelperDump

	@echo Building rpcClient
	@cd $(BASE)/cmd/rpcClient && GO111MODULE=on GOPATH=$(GOPATH) $(GOBIN) build $(GOBUILD_VERSION_ARGS) -o $(BASEDIR)/SatHelperClient

	@echo Building SatUI
	@GO111MODULE=off GOPATH=$(GOPATH) $(GOBIN) get github.com/asticode/go-astilectron-bundler/...
	@GO111MODULE=off GOPATH=$(GOPATH) $(GOBIN) install github.com/asticode/go-astilectron-bundler/astilectron-bundler
	@mkdir -p $(BASEDIR)/SatUI/
	@cd $(BASE)/cmd/SatUI && $(GOPATH)/bin/astilectron-bundler -l -o $(BASEDIR)/SatUI/

install: | $(BASE)
	@echo Installing
	@cd $(BASE) && cp $(BASEDIR)/SatHelperApp $(DESTDIR)/SatHelperApp
	@chmod +x $(DESTDIR)/SatHelperApp
	@cd $(BASE) && cp $(BASEDIR)/DemuxReplay $(DESTDIR)/DemuxReplay
	@chmod +x $(DESTDIR)/DemuxReplay
	@cd $(BASE) && cp $(BASEDIR)/xritparse $(DESTDIR)/xritparse
	@chmod +x $(DESTDIR)/xritparse
	@cd $(BASE) && cp $(BASEDIR)/xritcat $(DESTDIR)/xritcat
	@chmod +x $(DESTDIR)/xritcat
	@cd $(BASE) && cp $(BASEDIR)/xritimg $(DESTDIR)/xritimg
	@chmod +x $(DESTDIR)/xritimg
	@cd $(BASE) && cp $(BASEDIR)/xritpdcs $(DESTDIR)/xritpdcs
	@chmod +x $(DESTDIR)/xritpdcs
	@cd $(BASE) && cp $(BASEDIR)/SatHelperDump $(DESTDIR)/SatHelperDump
	@chmod +x $(DESTDIR)/SatHelperDump
	@cd $(BASE) && cp $(BASEDIR)/SatHelperClient $(DESTDIR)/SatHelperClient
	@chmod +x $(DESTDIR)/SatHelperClient
	@cd $(BASE) && cp $(BASEDIR)/SatHelperClient $(DESTDIR)/SatHelperClient
	@chmod +x $(DESTDIR)/SatHelperClient
	@cd $(BASE) && cp $(BASEDIR)/SatUI/$(shell go env GOOS)-$(shell go env GOARCH)/SatUI $(DESTDIR)/SatUI
	@chmod +x $(DESTDIR)/SatUI

test:
	go test -v -race $(shell go list ./... | grep -v /parts/ | grep -v /prime/ | grep -v /snap/ | grep -v /stage/ | grep -v /tmp/ | grep -v /librtlsdr/ )
