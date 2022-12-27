# driver-example: Tedge 驱动程序Demo

## 编译
* `go mod vendor && make build`

## 发布镜像
* 环境准备：`docker buildx create --name mybuilder`，`docker buildx use mybuilder`
* 打包并发布：`./script/dockerbuild.sh v1.0.0`
* 镜像发布说明：`dockerbuild.sh` 该脚本会同时发布 x86、armv7 两种架构镜像

## 说明
* Tedge有两种运行模式：DP模型、物模型(TyLink)
* 该驱动程序Demo同时演示了物模型、DP模型两种驱动
* 除非特别说明，Tedge默认情况都运行在DP模型下
* 驱动程序功能介绍：
    - 标准的驱动程序开发范式
    - 驱动服务接口实现示例
    - 子设备激活示例
    - 子设备状态更新示例
    - 子设备DP上报示例
    - 子设备指令处理示例

### DP模型
* 代码实现示例：`dpdriver`

### 物模型
* 代码实现示例：`tydriver`

## 技术支持
Tuya IoT Developer Platform: https://developer.tuya.com/en/

Tuya Developer Help Center: https://support.tuya.com/en/help

Tuya Work Order System: https://service.console.tuya.com/
