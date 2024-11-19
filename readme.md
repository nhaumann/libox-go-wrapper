# Livox Go Wrapper

A Go wrapper for the Livox LiDAR SDK that provides simple access to Livox devices through channels and Go-idiomatic interfaces.

## Prerequisites

- Windows 10 or later
- Git
- Visual Studio 2017
- CMake 3.10 or later
- Go 1.16 or later
- A Livox LiDAR device (tested with Mid-40/Mid-70/Horizon/Tele-15/HAP)

## Project Structure

```
livox-go-wrapper/
├── include/              # C wrapper header files
├── src/                 # C wrapper implementation
├── go/                  # Go package files
│   └── livox/          # Main package
├── build/              # Build output directory
├── Livox-SDK/          # Livox SDK submodule
└── CMakeLists.txt      # CMake configuration
```

## Building from Source

### Prerequisites

1. Install Visual Studio 2017
   - Ensure C++ desktop development is selected during installation

2. Install CMake
   - Download from https://cmake.org/download/
   - Add to system PATH

3. Install Go
   - Download from https://golang.org/dl/
   - Add to system PATH

### Building the Package

1. Clone the repository and SDK:
```bash
git clone https://github.com/nhaumann/livox-go-wrapper.git
cd livox-go-wrapper
git submodule add https://github.com/Livox-SDK/Livox-SDK.git
```

2. Build Livox SDK:
```bash
cd Livox-SDK
mkdir build
cd build
cmake .. -G "Visual Studio 15 2017 Win64"
cmake --build . --config Release
cd ../..
```

3. Build the wrapper:
```bash
mkdir build
cd build
cmake .. -G "Visual Studio 15 2017 Win64"
cmake --build . --config Release
cd ..
```

The wrapper DLL will be created at `build/bin/Release/livox_wrapper.dll`.

## Usage

### Import the Package

In your Go project:

```bash
go get github.com/nhaumann/livox-go-wrapper/go/livox
```

### Basic Example

```go
package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"
    
    "github.com/nhaumann/livox-go-wrapper/go/livox"
)

func main() {
    // Create new scanner with buffer for 100 point clouds
    scanner := livox.NewScanner(100)
    
    // Start the scanner
    if err := scanner.Start(); err != nil {
        log.Fatalf("Failed to start scanner: %v", err)
    }
    defer scanner.Stop()

    // Handle device events
    go func() {
        for event := range scanner.DeviceEvents() {
            if event.Connected {
                log.Printf("Device connected: %s (Handle: %d)", 
                    event.Device.BroadcastCode, event.Device.Handle)
            } else {
                log.Printf("Device disconnected: %s (Handle: %d)", 
                    event.Device.BroadcastCode, event.Device.Handle)
            }
        }
    }()

    // Handle point cloud data
    go func() {
        for cloud := range scanner.PointCloud() {
            log.Printf("Received point cloud from device %d: %d points", 
                cloud.DeviceHandle, len(cloud.Points))
            // Process points
            for _, point := range cloud.Points {
                _ = point // Process each point (X, Y, Z, Intensity)
            }
        }
    }()

    // Wait for interrupt signal
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan
}
```

### Data Types

```go
// Point represents a single point in 3D space
type Point struct {
    X         float32  // X coordinate in meters
    Y         float32  // Y coordinate in meters
    Z         float32  // Z coordinate in meters
    Intensity float32  // Normalized intensity (0.0-1.0)
}

// PointCloud represents a collection of points from a single scan
type PointCloud struct {
    DeviceHandle uint8
    Points       []Point
}

// Device represents a connected Livox device
type Device struct {
    Handle        uint8
    BroadcastCode string
}

// DeviceEvent represents a device connection/disconnection event
type DeviceEvent struct {
    Device    Device
    Connected bool
}
```

### API Reference

```go
// Create a new scanner with specified buffer size
scanner := livox.NewScanner(bufferSize int)

// Start scanning
err := scanner.Start()

// Stop scanning
scanner.Stop()

// Get point cloud channel (receive data)
pointChan := scanner.PointCloud()

// Get device events channel (monitor connections)
deviceChan := scanner.DeviceEvents()

// Get list of connected devices
devices := scanner.GetDevices()
```

## Troubleshooting

### Common Issues

1. **CMake cannot find Visual Studio**
   - Ensure Visual Studio 2017 is installed
   - Run `cmake --help` to see available generators
   - Use the exact generator name in the Makefile

2. **Livox SDK build fails**
   - Check Visual Studio installation
   - Ensure all Windows SDKs are installed
   - Try building the SDK manually following their documentation

3. **Go build fails**
   - Check that CGO is enabled (set CGO_ENABLED=1)
   - Verify wrapper DLL is in the correct location
   - Check PATH includes Visual Studio build tools

4. **Package import fails**
   - Ensure correct import path: `github.com/nhaumann/livox-go-wrapper/go/livox`
   - Run `go mod tidy` to update dependencies
   - Verify go.mod file exists in project root

## Acknowledgments

- [Livox SDK](https://github.com/Livox-SDK/Livox-SDK) - The underlying SDK this wrapper is built upon