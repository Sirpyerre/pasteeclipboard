BINARY_NAME=pastee

MAIN_PACKAGE=./cmd/pastee

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# Output folder
BUILD_DIR=bin

.PHONY: all build clean run

all: build

build:
	@echo "🛠️  Compilando para $(GOOS)/$(GOARCH)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "✅ Binario generado en $(BUILD_DIR)/$(BINARY_NAME)"

run:
	go run $(MAIN_PACKAGE)

clean:
	@echo "🧹 Cleaning..."
	@rm -rf $(BUILD_DIR)
