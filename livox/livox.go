package livox

/*
#include <stdint.h>
#include <stdlib.h>
#include "livox_sdk.h"

#cgo CFLAGS: -I${SRCDIR}/../include
#cgo LDFLAGS: -L${SRCDIR}/../include -llivox_wrapper -llivox_sdk_static

#include "livox_wrapper.h"

extern void pointCloudCallback(uint8_t handle, uint8_t *data, uint32_t data_num, uint8_t data_type);
extern void deviceInfoCallback(uint8_t handle, char* broadcast_code, uint8_t connected);
*/
import "C"
import (
	"fmt"
	"log"
	"sync"
	"unsafe"
)

// Point represents a single point in 3D space with intensity
type Point struct {
	X         float32
	Y         float32
	Z         float32
	Intensity float32
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

// Scanner manages the Livox SDK and provides channels for point cloud and device data
type Scanner struct {
	devicesMutex sync.RWMutex
	devices      map[uint8]string

	// Channels for communication
	pointChan  chan PointCloud
	deviceChan chan DeviceEvent

	// Control channels
	stopChan chan struct{}
	doneChan chan struct{}
}

var (
	// Global scanner instance for C callback handling
	globalScanner *Scanner
	once          sync.Once
)

// NewScanner creates a new Scanner instance
func NewScanner(pointCloudBufferSize int) *Scanner {
	s := &Scanner{
		devices:    make(map[uint8]string),
		pointChan:  make(chan PointCloud, pointCloudBufferSize),
		deviceChan: make(chan DeviceEvent, 10),
		stopChan:   make(chan struct{}),
		doneChan:   make(chan struct{}),
	}

	// Set global scanner instance for C callbacks
	once.Do(func() {
		globalScanner = s
	})

	return s
}

// Start initializes the Livox SDK and begins scanning
func (s *Scanner) Start() error {
	if ret := C.InitSdk(); ret != 0 {
		return fmt.Errorf("failed to initialize SDK: %d", ret)
	}

	// Register callbacks
	C.RegisterPointCloudCallback((C.PointCloudCb)(C.pointCloudCallback))
	C.RegisterDeviceInfoCallback((C.DeviceInfoCb)(C.deviceInfoCallback))

	return nil
}

// Stop gracefully shuts down the scanner and SDK
func (s *Scanner) Stop() {
	close(s.stopChan)
	C.StopSdk()
	close(s.pointChan)
	close(s.deviceChan)
	<-s.doneChan
}

// PointCloud returns the channel for receiving point cloud data
func (s *Scanner) PointCloud() <-chan PointCloud {
	return s.pointChan
}

// DeviceEvents returns the channel for receiving device connection events
func (s *Scanner) DeviceEvents() <-chan DeviceEvent {
	return s.deviceChan
}

// GetDevices returns a slice of currently connected devices
func (s *Scanner) GetDevices() []Device {
	s.devicesMutex.RLock()
	defer s.devicesMutex.RUnlock()

	devices := make([]Device, 0, len(s.devices))
	for handle, code := range s.devices {
		devices = append(devices, Device{
			Handle:        handle,
			BroadcastCode: code,
		})
	}
	return devices
}

//export pointCloudCallback
func pointCloudCallback(handle C.uint8_t, data *C.uint8_t, dataNum C.uint32_t, dataType C.uint8_t) {
	if globalScanner == nil {
		return
	}

	// Convert C array to Go slice
	pointData := (*[1 << 30]C.LivoxRawPoint)(unsafe.Pointer(data))[:dataNum:dataNum]

	// Create point cloud with initial capacity
	cloud := PointCloud{
		DeviceHandle: uint8(handle),
		Points:       make([]Point, 0, dataNum),
	}

	// Convert points
	for _, p := range pointData {
		intensity := float32(p.reflectivity) / 255.0 // Normalize intensity to 0-1 range
		cloud.Points = append(cloud.Points, Point{
			X:         float32(p.x) / 1000.0, // Convert to meters
			Y:         float32(p.y) / 1000.0,
			Z:         float32(p.z) / 1000.0,
			Intensity: intensity,
		})
	}

	// Try to send point cloud data, skip if channel is full
	select {
	case globalScanner.pointChan <- cloud:
	default:
		log.Println("Warning: Point cloud channel full, dropping frame")
	}
}

//export deviceInfoCallback
func deviceInfoCallback(handle C.uint8_t, broadcastCode *C.char, connected C.uint8_t) {
	if globalScanner == nil {
		return
	}

	code := C.GoString(broadcastCode)
	device := Device{
		Handle:        uint8(handle),
		BroadcastCode: code,
	}

	globalScanner.devicesMutex.Lock()
	if connected == 0 {
		delete(globalScanner.devices, uint8(handle))
	} else {
		globalScanner.devices[uint8(handle)] = code
	}
	globalScanner.devicesMutex.Unlock()

	// Send device event
	select {
	case globalScanner.deviceChan <- DeviceEvent{
		Device:    device,
		Connected: connected != 0,
	}:
	default:
		log.Println("Warning: Device event channel full, dropping event")
	}
}
