BINARY_NAME=pastee
MAIN_PACKAGE=./cmd/pastee
BUILD_DIR=bin

# Defaults
GOOS ?=
GOARCH ?=
EXE ?=

# Export so child processes (go build) see them on all shells
export GOOS
export GOARCH
export EXE

# OS/arch detection
ifeq ($(OS), Windows_NT)
	GOOS ?= windows
	GOARCH ?= amd64
	EXE := .exe
	SHELL := cmd
else
	UNAME_S := $(shell uname -s)
	UNAME_M := $(shell uname -m)
	ifeq ($(UNAME_S), Darwin)
		GOOS ?= darwin
		# Apple Silicon vs Intel
		ifeq ($(UNAME_M),arm64)
			GOARCH ?= arm64
		else
			GOARCH ?= amd64
		endif
	else
		GOOS ?= linux
		ifeq ($(UNAME_M),aarch64)
			GOARCH ?= arm64
		else
			GOARCH ?= amd64
		endif
	endif
endif

.PHONY: all build clean run

all: build

build:
	@echo "Building for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BUILD_DIR)/$(BINARY_NAME)$(EXE) $(MAIN_PACKAGE)
	@echo "Binary generated at $(BUILD_DIR)/$(BINARY_NAME)$(EXE)"

run:
	go run $(MAIN_PACKAGE)

clean:
	@echo "🧹 Cleaning..."
	@rm -rf $(BUILD_DIR)
