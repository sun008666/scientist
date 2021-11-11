# Makefile for logging project
#
# Targets:
#   all: builds the code
#   build: build the linux code
#   release: build the darwin code
#   test: runs the tests

.PHONY:install all

# Package path
PACKAGE_NAME := logging
PACKAGE_BASE_PATH := zonst
PACKAGE_LIB_NAME := $(PACKAGE_NAME).a
PACKAGE_IMPORT_PATH := $(PACKAGE_BASE_PATH)/$(PACKAGE_NAME)
PACKAGE_STORE_FILENAME := $(PACKAGE_BASE_PATH)/$(PACKAGE_LIB_NAME)
GOVERSION := $(shell go version |awk '{print $$3}')

GOENVPATH := $$GOPATH

linux_pkg := $(GOENVPATH)/pkg/linux_amd64/$(PACKAGE_BASE_PATH)
darwin_pkg := $(GOENVPATH)/pkg/darwin_amd64/$(PACKAGE_BASE_PATH)
windows_pkg := $(GOENVPATH)/pkg/windows_amd64/$(PACKAGE_BASE_PATH)
linux_src_pkg := pkg/$(GOVERSION)/linux_amd64/$(PACKAGE_STORE_FILENAME)
darwin_src_pkg := pkg/$(GOVERSION)/darwin_amd64/$(PACKAGE_STORE_FILENAME)
windows_src_pkg := pkg/$(GOVERSION)/windows_amd64/$(PACKAGE_STORE_FILENAME)

all:binary source
	
binary:
	@[ -d $(linux_pkg) ] || mkdir -p $(linux_pkg)
	cp -rf $(linux_src_pkg) $(linux_pkg)
	@[ -d $(darwin_pkg) ] || mkdir -p $(darwin_pkg)
	cp -rf $(darwin_src_pkg) $(darwin_pkg)
	@[ -d $(windows_pkg) ] || mkdir -p $(windows_pkg)
	cp -rf $(windows_src_pkg) $(windows_pkg)
	@echo cp -rf $(linux_src_pkg) $(linux_pkg)
	@echo cp -rf $(darwin_src_pkg) $(darwin_pkg)
	@echo cp -rf $(windows_src_pkg) $(windows_pkg)
	@echo 'pkg installed'

source:
	@echo 'GOPATH='$(GOENVPATH)
	@[ -d $(GOENVPATH)/src/$(PACKAGE_IMPORT_PATH) ] || mkdir -p $(GOENVPATH)/src/$(PACKAGE_IMPORT_PATH)
	cp -rf src/$(PACKAGE_IMPORT_PATH)/docs.go $(GOENVPATH)/src/$(PACKAGE_IMPORT_PATH)
	@echo cp -rf src/$(PACKAGE_IMPORT_PATH)/docs.go $(GOENVPATH)/src/$(PACKAGE_IMPORT_PATH)
	@echo 'src installed'
