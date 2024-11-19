# Makefile for livox-go-wrapper

# Paths
ROOT_DIR := $(CURDIR)
BUILD_DIR := $(ROOT_DIR)/build
LIVOX_SDK_DIR := $(ROOT_DIR)/Livox-SDK
LIVOX_SDK_BUILD_DIR := $(LIVOX_SDK_DIR)/build
GO_DIR := $(ROOT_DIR)/go
INCLUDE_DIR := $(ROOT_DIR)/include

# Commands
CMAKE := cmake
CMAKE_GENERATOR := "Visual Studio 15 2017 Win64"
GIT := git
GO := go
CONFIG := Release

# Check if we're running on Windows
ifeq ($(OS),Windows_NT)
    COPY := copy /Y
    MKDIR := mkdir
    RM := del /F /Q
    RMDIR := rmdir /S /Q
    NULL_DEVICE := NUL
else
    COPY := cp
    MKDIR := mkdir -p
    RM := rm -f
    RMDIR := rm -rf
    NULL_DEVICE := /dev/null
endif

.PHONY: all clean build_sdk build_wrapper build_go install init

# Default target
all: init build_sdk build_wrapper build_go

# Initialize project (clone Livox SDK if needed)
init:
	@if not exist "$(LIVOX_SDK_DIR)\.git" ( \
		echo Cloning Livox SDK... && \
		$(GIT) submodule add https://github.com/Livox-SDK/Livox-SDK.git 2>$(NULL_DEVICE) || \
		$(GIT) submodule update --init --recursive \
	)

# Build Livox SDK
build_sdk: init
	@if not exist "$(LIVOX_SDK_BUILD_DIR)" $(MKDIR) "$(LIVOX_SDK_BUILD_DIR)"
	cd $(LIVOX_SDK_BUILD_DIR) && $(CMAKE) .. -G $(CMAKE_GENERATOR)
	cd $(LIVOX_SDK_BUILD_DIR) && $(CMAKE) --build . --config $(CONFIG)

# Create build directory
$(BUILD_DIR):
	@if not exist "$(BUILD_DIR)" $(MKDIR) "$(BUILD_DIR)"

# Configure CMake for wrapper
configure: $(BUILD_DIR)
	cd $(BUILD_DIR) && $(CMAKE) .. -G $(CMAKE_GENERATOR)

# Build C++ wrapper
build_wrapper: configure
	cd $(BUILD_DIR) && $(CMAKE) --build . --config $(CONFIG)

# Build Go application
build_go: build_wrapper
	cd $(GO_DIR) && $(GO) build -o livox_app.exe main.go

# Clean build artifacts
clean:
	@if exist "$(BUILD_DIR)" $(RMDIR) "$(BUILD_DIR)"
	@if exist "$(LIVOX_SDK_BUILD_DIR)" $(RMDIR) "$(LIVOX_SDK_BUILD_DIR)"
	@if exist "$(GO_DIR)\livox_app.exe" $(RM) "$(GO_DIR)\livox_app.exe"

# Help target
help:
	@echo Available targets:
	@echo   all           : Build everything (default)
	@echo   init         : Initialize/update Livox SDK submodule
	@echo   build_sdk    : Build Livox SDK
	@echo   build_wrapper: Build C++ wrapper
	@echo   build_go     : Build Go application
	@echo   clean        : Clean build artifacts
	@echo   help         : Show this help message

# Build order dependencies
build_wrapper: build_sdk
build_go: build_wrapper