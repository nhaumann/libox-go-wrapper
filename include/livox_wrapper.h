// livox_wrapper.h
#ifndef LIVOX_WRAPPER_H_
#define LIVOX_WRAPPER_H_

#ifdef __cplusplus
extern "C"
{
#endif

#ifdef _WIN32
#define LIVOX_API __declspec(dllexport)
#else
#define LIVOX_API
#endif

#include <stdint.h>
#include "livox_sdk.h"

    // Function pointer types for callbacks
    typedef void (*PointCloudCb)(uint8_t handle, uint8_t *data, uint32_t data_num, uint8_t data_type);
    typedef void (*DeviceInfoCb)(uint8_t handle, char *broadcast_code, uint8_t connected);

    // Core functions
    LIVOX_API void RegisterPointCloudCallback(PointCloudCb cb);
    LIVOX_API void RegisterDeviceInfoCallback(DeviceInfoCb cb);
    LIVOX_API uint32_t InitSdk(void);
    LIVOX_API uint32_t StopSdk(void);

#ifdef __cplusplus
}
#endif

#endif // LIVOX_WRAPPER_H_