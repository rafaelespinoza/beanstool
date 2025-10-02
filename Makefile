# Environment
BASE_PATH := $(shell pwd)
BIN_DIR := $(BASE_PATH)/bin
DIST_PATH := $(BASE_PATH)/dist

# Go parameters
GOCMD = go

.PHONY: all build test dependencies vet install clean release

# Rules
all: test build

build: dependencies
	mkdir -p $(BIN_DIR) && $(GOCMD) build -o $(BIN_DIR)/beanstool .

test: dependencies
	cd $(BASE_PATH); $(GOCMD) test -v ./...

dependencies:
	$(GOCMD) mod tidy

vet:
	$(GOCMD) vet ./...

install:
	cp -f $(BIN_DIR)/beanstool /usr/bin/

clean:
	rm -rf $(DIST_PATH)
	$(GOCMD) clean .

# Automates binary building on many platforms. This Makefile won't install the
# tool for you. Read more at https://goreleaser.com.
#
# This target is set up to keep effects local by default, via the --snapshot
# flag. Use the Makefile variable, ARGS, to override that:
# 	$ make release ARGS='--snapshot=false'
GORELEASER ?= goreleaser
release:
	$(GORELEASER) release --clean --snapshot $(ARGS)
