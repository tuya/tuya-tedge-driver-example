[English](README.md) | [中文版](README_CN.md)
# driver-example: TEdge Driver Demo

## Build
* `go mod vendor && make build`

## Release image

1. Prepare the environment: `docker buildx create --name mybuilder`, `docker buildx use mybuilder`
2. Package and release image: `./script/dockerbuild.sh v1.0.0`

**Note**: Both x86 and ARMv7 images are released by the script`dockerbuild.sh`.

## Description
* TEdge supports two running modes: DP model and things data model (TuyaLink)
* The drivers for both the DP model and TuyaLink model are described in the following samples.
* Unless otherwise specified, TEdge runs with the DP model.
* The following driver features are described in the sample:
   - Standard driver programming
   - Implementation of the driver service interface
   - Sub-device activation
   - Sub-device status update
   - Sub-device DP reporting
   - Sub-device command processing

### DP model
* Sample: `dpdriver`

### Things data model
* Sample: `tydriver`

## Technical support
Tuya IoT Developer Platform: https://developer.tuya.com/en/

Tuya Developer Help Center: https://support.tuya.com/en/help

Tuya Service Ticket System: https://service.console.tuya.com/