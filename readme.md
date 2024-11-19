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

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/livox-go-wrapper.git
cd livox-go-wrapper
```

2. Build using make (this will handle all steps including SDK setup):
```bash
make
```

Or build step by step:

```bash
# Initialize and build Livox SDK
make init
make build_sdk

# Build the wrapper
make build_wrapper

# Build Go application
make build_go
```


## Point Cloud Data Structure

The `Point` struct represents a single point in 3D space:

```go
type Point struct {
    X         float32  // X coordinate in meters
    Y         float32  // Y coordinate in meters
    Z         float32  // Z coordinate in meters
    Intensity float32  // Normalized intensity (0.0-1.0)
}
```

## API Reference

### Scanner

```go
// Create a new scanner
scanner := livox.NewScanner(bufferSize int)

// Start scanning
err := scanner.Start()

// Stop scanning
scanner.Stop()

// Get point cloud channel
pointChan := scanner.PointCloud()

// Get device events channel
deviceChan := scanner.DeviceEvents()

// Get list of connected devices
devices := scanner.GetDevices()
```


## Usage Examples

### Basic Connection Example
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
            // Process point cloud data here
        }
    }()

    // Wait for interrupt signal
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan
}
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

### Build Steps

1. Clone with submodules:
```bash
git clone https://github.com/nhaumann/livox-go-wrapper.git
cd livox-go-wrapper
```

2. Build everything:
```bash
make
```

The build process will:
1. Clone and build the Livox SDK
2. Build the C wrapper library
3. Build the Go package and example application

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
   - Verify all DLLs are in the correct locations
   - Check PATH includes Visual Studio build tools


## Acknowledgments

- [Livox SDK](https://github.com/Livox-SDK/Livox-SDK) - The underlying SDK this wrapper is built upon
