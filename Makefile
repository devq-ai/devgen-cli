BINARY_NAME := devgen
BUILD_DIR := build

.PHONY: build clean help install

help:
	@echo "DevGen CLI Build System"

build:
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) main.go
	@echo "Build complete"

clean:
	@rm -rf $(BUILD_DIR)

install: build
	@cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/ || echo "Try: sudo make install"
