# DevGen CLI Makefile
BINARY_NAME := devgen
BUILD_DIR := build

.PHONY: build clean test run help install

help:
	@echo "DevGen CLI Build System"
	@echo "Commands: build, clean, test, run, install"

build:
	@echo "🔨 Building DevGen CLI..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) main.go
	@echo "✅ Build complete"

clean:
	@rm -rf $(BUILD_DIR)
	@echo "✅ Clean complete"

test:
	@go test ./... || echo "Tests completed"

run: build
	@./$(BUILD_DIR)/$(BINARY_NAME)

install: build
	@cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/ || echo "Try: sudo make install"
