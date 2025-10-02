# Package configuration
PROJECT = beanstool
COMMANDS = beanstool

# Environment
BASE_PATH := $(shell pwd)
BUILD_PATH := $(BASE_PATH)/build
BIN_DIR := $(BASE_PATH)/bin
VERSION ?= $(shell git branch 2> /dev/null | sed -e '/^[^*]/d' -e 's/* \(.*\)/\1/' | tr '/' '_')
BUILD ?= $(shell date --utc +%FT%T%z | tr ' ' 'T')
ASSETS := static

# PACKAGES
PKG_OS = darwin linux
PKG_ARCH = amd64
PKG_CONTENT = LICENSE

# Go parameters
GOCMD = go

.PHONY: all build test dependencies install packages clean

# Rules
all: test build

build: dependencies
	mkdir -p $(BIN_DIR)
	for cmd in $(COMMANDS); do \
		$(GOCMD) build -o $(BIN_DIR)/ -ldflags "-X main.version=$(VERSION) -X main.build=$(BUILD)" $${cmd}.go; \
	done

test: dependencies
	cd $(BASE_PATH); $(GOCMD) test -v ./...

dependencies:
	$(GOCMD) mod tidy

vet:
	$(GOCMD) vet ./...

install:
	for cmd in $(COMMANDS); do \
		cp -rf $${cmd} /usr/bin/; \
	done

packages: clean
	for os in $(PKG_OS); do \
		for arch in $(PKG_ARCH); do \
			cd $(BASE_PATH); \
			mkdir -p $(BUILD_PATH)/$(PROJECT)_$(VERSION)_$${os}_$${arch}; \
			for cmd in $(COMMANDS); do \
				GOOS=$${os} GOARCH=$${arch} $(GOCMD) build -ldflags "-X main.version=$(VERSION) -X main.build=\"$(BUILD)\"" -o $(BUILD_PATH)/$(PROJECT)_$(VERSION)_$${os}_$${arch}/$${cmd} $${cmd}.go ; \
			done; \
			for content in $(PKG_CONTENT); do \
				cp -rf $${content} $(BUILD_PATH)/$(PROJECT)_$(VERSION)_$${os}_$${arch}/; \
			done; \
			cd  $(BUILD_PATH) && tar -cvzf $(BUILD_PATH)/$(PROJECT)_$(VERSION)_$${os}_$${arch}.tar.gz $(PROJECT)_$(VERSION)_$${os}_$${arch}/; \
		done; \
	done;

clean:
	echo $(VERSION)
	rm -rf $(BUILD_PATH)

	$(GOCMD) clean .
