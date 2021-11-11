# Makefile for logging project
#
# Targets:
#   all: builds the code
#   build: build the linux code
#   release: build the darwin code
#   test: runs the tests

.PHONY:test build release all

# Package path
PACKAGE_NAME = logging
PACKAGE_LIB_NAME = $(PACKAGE_NAME).a
PACKAGE_IMPORT_PATH = zonst/$(PACKAGE_NAME)
PACKAGE_STORE_FILENAME = zonst/$(PACKAGE_LIB_NAME)

RELEASE_PACKAGE_PATH = $(HOME)/go-pkg
dst_root_path = $(RELEASE_PACKAGE_PATH)/$(PACKAGE_NAME)
dst_package_path = $(RELEASE_PACKAGE_PATH)/$(PACKAGE_NAME)/pkg
dst_src_path = $(RELEASE_PACKAGE_PATH)/$(PACKAGE_NAME)/src/$(PACKAGE_IMPORT_PATH)

GOVERSION := $(shell go version |awk '{print $$3}')

all:test build release
	
test:
	go test -v

build:
	GOOS=linux go build -o pkg/pkg/$(GOVERSION)/linux_amd64/$(PACKAGE_STORE_FILENAME)
	GOOS=darwin go build -o pkg/pkg/$(GOVERSION)/darwin_amd64/$(PACKAGE_STORE_FILENAME)
	GOOS=windows go build -o pkg/pkg/$(GOVERSION)/windows_amd64/$(PACKAGE_STORE_FILENAME)
	@echo successful installed!

install:
	GOOS=linux go install
	GOOS=darwin go install
	GOOS=windows go install
	@echo successful installed!

release:
	@echo copy makefile
	@[ -d $(dst_root_path) ] || mkdir -p $(dst_root_path)
	cp -rf pkg/Makefile $(dst_root_path)
	@echo copy package
	@[ -d $(dst_package_path) ] || mkdir -p $(dst_package_path)
	cp -rf pkg/pkg/* $(dst_package_path)
	@echo copy src
	@[ -d $(dst_src_path) ] || mkdir -p $(dst_src_path)
	cp -rf pkg/docs.go $(dst_src_path)
	@echo successful copied!

