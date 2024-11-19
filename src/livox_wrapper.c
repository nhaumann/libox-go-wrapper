// livox_wrapper.c
#include "livox_wrapper.h"
#include <stdio.h>
#include <string.h>

// Static storage for device information
static void (*point_cloud_callback)(uint8_t, uint8_t *, uint32_t, uint8_t) = NULL;
static void (*device_info_callback)(uint8_t, char *, uint8_t) = NULL;

// Callback to handle point cloud data
static void OnLidarDataCallback(uint8_t handle, LivoxEthPacket *data, uint32_t data_num, void *client_data)
{
    if (data && point_cloud_callback)
    {
        point_cloud_callback(
            handle,
            (uint8_t *)data->data,
            data_num,
            data->data_type);
    }
}

// Device state change callback
static void OnDeviceInfoChange(const DeviceInfo *info, DeviceEvent type)
{
    if (info == NULL || info->handle >= kMaxLidarCount || !device_info_callback)
    {
        return;
    }

    uint8_t handle = info->handle;
    uint8_t is_connected = 0;

    switch (type)
    {
    case kEventConnect:
        is_connected = 1;
        printf("[C] Lidar connected: %s\n", info->broadcast_code);
        break;

    case kEventDisconnect:
        is_connected = 0;
        printf("[C] Lidar disconnected: %s\n", info->broadcast_code);
        break;

    case kEventStateChange:
        is_connected = 1;
        printf("[C] Lidar state changed: %s\n", info->broadcast_code);
        break;
    }

    // Create a modifiable copy of the broadcast code
    char broadcast_code[kBroadcastCodeSize];
    strncpy(broadcast_code, info->broadcast_code, kBroadcastCodeSize - 1);
    broadcast_code[kBroadcastCodeSize - 1] = '\0';

    device_info_callback(handle, broadcast_code, is_connected);

    // If device is connected and in normal state, start sampling
    if (is_connected && info->state == kLidarStateNormal)
    {
        LidarStartSampling(handle, NULL, NULL);
    }
}

// Broadcast callback
static void OnDeviceBroadcast(const BroadcastDeviceInfo *info)
{
    if (info == NULL || info->dev_type == kDeviceTypeHub)
    {
        return;
    }

    printf("[C] Found broadcast device: %s\n", info->broadcast_code);

    uint8_t handle = 0;
    if (AddLidarToConnect(info->broadcast_code, &handle) == kStatusSuccess)
    {
        SetDataCallback(handle, OnLidarDataCallback, NULL);
    }
}

LIVOX_API void RegisterPointCloudCallback(void (*callback)(uint8_t, uint8_t *, uint32_t, uint8_t))
{
    point_cloud_callback = callback;
}

LIVOX_API void RegisterDeviceInfoCallback(void (*callback)(uint8_t, char *, uint8_t))
{
    device_info_callback = callback;
}

LIVOX_API uint32_t InitSdk()
{
    // Initialize SDK
    if (!Init())
    {
        printf("[C] Failed to initialize Livox SDK\n");
        return 1;
    }

    // Set callbacks
    SetBroadcastCallback(OnDeviceBroadcast);
    SetDeviceStateUpdateCallback(OnDeviceInfoChange);

    // Start device discovery
    if (!Start())
    {
        printf("[C] Failed to start Livox SDK\n");
        Uninit();
        return 2;
    }

    printf("[C] Livox SDK initialized successfully\n");
    return 0;
}

LIVOX_API uint32_t StopSdk()
{
    Uninit();
    printf("[C] Livox SDK stopped\n");
    return 0;
}