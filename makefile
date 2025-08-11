BINARY_NAME=pastee
MAIN_PACKAGE=./cmd/pastee
BUILD_DIR=bin

# Defaults
GOOS ?=
GOARCH ?=
EXE ?=

export GOOS
export GOARCH
export EXE

ifeq ($(OS),Windows_NT)
	PROC_ARCH := $(PROCESSOR_ARCHITECTURE)
	GOOS := windows
	EXE := .exe
	SHELL := cmd
	MKDIR := if not exist $(BUILD_DIR) mkdir $(BUILD_DIR)

	# Detect architecture
	ifeq ($(PROC_ARCH),ARM64)
		GOARCH := arm64
	else ifeq ($(PROC_ARCH),AMD64)
		GOARCH := amd64
	else ifeq ($(PROC_ARCH),x86)
		GOARCH := 386
	else
		GOARCH := amd64
	endif

	SETEVARS := set GOOS=$(GOOS) && set GOARCH=$(GOARCH) &&
else
	UNAME_S := $(shell uname -s)
	UNAME_M := $(shell uname -m)
	MKDIR := mkdir -p $(BUILD_DIR)
	SETEVARS := GOOS=$(GOOS) GOARCH=$(GOARCH)

	ifeq ($(UNAME_S),Darwin)
		GOOS := darwin
		ifeq ($(UNAME_M),arm64)
			GOARCH := arm64
		else
			GOARCH := amd64
		endif
	else
		GOOS := linux
		ifeq ($(UNAME_M),aarch64)
			GOARCH := arm64
		else
			GOARCH := amd64
		endif
	endif
endif

.PHONY: all build clean run

all: build

build:
	@echo Building for $(GOOS)/$(GOARCH)...
	@$(MKDIR)
	$(SETEVARS) go build -o $(BUILD_DIR)/$(BINARY_NAME)$(EXE) $(MAIN_PACKAGE)
	@echo Binary generated at $(BUILD_DIR)/$(BINARY_NAME)$(EXE)

run:
	go run $(MAIN_PACKAGE)

clean:
	@echo 🧹 Cleaning...
ifeq ($(OS),Windows_NT)
	if exist $(BUILD_DIR) rmdir /S /Q $(BUILD_DIR)
else
	rm -rf $(BUILD_DIR)
endif
