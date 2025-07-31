BINARY_NAME=pastee

MAIN_PACKAGE=./cmd/pastee

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# Carpeta de salida (opcional)
BUILD_DIR=bin

.PHONY: all build clean run

# Compilar la aplicación
all: build

build:
	@echo "🛠️  Compilando para $(GOOS)/$(GOARCH)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "✅ Binario generado en $(BUILD_DIR)/$(BINARY_NAME)"

# Ejecutar directamente
run:
	go run $(MAIN_PACKAGE)

# Limpiar archivos compilados
clean:
	@echo "🧹 Cleaning..."
	@rm -rf $(BUILD_DIR)
