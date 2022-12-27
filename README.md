[English](README.md) | [中文版](README_CN.md)
# driver-example: Tedge Driver Demo

## build
* `go mod vendor && make build`

## Publish Docker Image
* Prepare environment: `docker buildx create --name mybuilder` && `docker buildx use mybuilder`
* Package and Publish: `./script/dockerbuild.sh v1.0.0`
* Script Release Notes：`dockerbuild.sh` This script will build x86 and armv7 architecture images at the same time.

## Driver Notes
* Tedge supports two running modes: DP mode and Thing mode (TyLink).
* This Driver Demo contains code samples of both DP model and Thing model.
* Unless otherwise specified, Tedge runs in DP model by default.
* Driver samples introduction：
    - Standard driver development paradigm
    - Example of Driver service interface implementation
    - Example of sub device activation
    - Example of sub device status update
    - Example of sub device DP message reporting
    - Example of sub device command processing

### DP mode
* Dp code samples: `dpdriver`

### TyLink mode
* TyLink code samples: `tydriver`

## Technical Support
Tuya IoT Developer Platform: https://developer.tuya.com/en/

Tuya Developer Help Center: https://support.tuya.com/en/help

Tuya Work Order System: https://service.console.tuya.com/
