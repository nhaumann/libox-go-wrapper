cmake_minimum_required(VERSION 3.10)

project(livox_go_wrapper)

# Specify the updated paths to the Livox SDK relative to project root
set(LIVOX_SDK_PATH "${CMAKE_CURRENT_SOURCE_DIR}/Livox-SDK")
set(LIVOX_SDK_INCLUDE "${LIVOX_SDK_PATH}/sdk_core/include")
set(LIVOX_SDK_LIB "${LIVOX_SDK_PATH}/build/sdk_core/Release")

# Verify SDK paths exist
if(NOT EXISTS ${LIVOX_SDK_INCLUDE})
    message(FATAL_ERROR "Livox SDK include path not found: ${LIVOX_SDK_INCLUDE}")
endif()

if(NOT EXISTS ${LIVOX_SDK_LIB})
    message(FATAL_ERROR "Livox SDK library path not found: ${LIVOX_SDK_LIB}. Please build the SDK first.")
endif()

# Include directories
include_directories(${LIVOX_SDK_INCLUDE})
include_directories(${CMAKE_SOURCE_DIR}/include)

# Link directories
link_directories(${LIVOX_SDK_LIB})

# Add the wrapper source
add_library(livox_wrapper SHARED 
    src/livox_wrapper.c
)

# Link Livox SDK to the wrapper
target_link_libraries(livox_wrapper 
    ${LIVOX_SDK_LIB}/livox_sdk_static.lib
)

# Ensure symbols are exported
set_target_properties(livox_wrapper PROPERTIES 
    WINDOWS_EXPORT_ALL_SYMBOLS ON
    OUTPUT_NAME "livox_wrapper"
)

# Set output directories
set_target_properties(livox_wrapper PROPERTIES
    RUNTIME_OUTPUT_DIRECTORY "${CMAKE_BINARY_DIR}/bin"
    LIBRARY_OUTPUT_DIRECTORY "${CMAKE_BINARY_DIR}/lib"
    ARCHIVE_OUTPUT_DIRECTORY "${CMAKE_BINARY_DIR}/lib"
)

# Add compile definitions if needed
target_compile_definitions(livox_wrapper PRIVATE
    WIN32_LEAN_AND_MEAN
    _CRT_SECURE_NO_WARNINGS
)

# Print configuration info for debugging
message(STATUS "Livox SDK Include: ${LIVOX_SDK_INCLUDE}")
message(STATUS "Livox SDK Lib: ${LIVOX_SDK_LIB}")
message(STATUS "Project Source Dir: ${CMAKE_SOURCE_DIR}")
message(STATUS "Build Output Dir: ${CMAKE_BINARY_DIR}")
