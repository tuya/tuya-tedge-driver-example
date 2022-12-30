package main

import (
	"fmt"

	"github.com/tuya/tuya-tedge-driver-example/dpmqtt"
	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
	"github.com/tuya/tuya-tedge-driver-sdk-go/service"
)

const (
	serviceName string = "driver-example"
)

//TEdge 驱动开发示例：DP模型，对接MQTT协议设备
func main() {
	sdkLog := commons.DefaultLogger(commons.DebugLevel, serviceName)
	//step1: 创建驱动服务
	ds := service.NewDPService(sdkLog)

	// DP 模型驱动
	tEdgeModel := ds.GetTEdgeModel()
	if tEdgeModel == commons.DPModel {
		sdkLog.Infof("TEdge run in dp model")

		//step2: 实现驱动接口 `DPModelDriver`
		dDpDriver := dpmqtt.NewDemoMqttDPDriver(ds)
		go dDpDriver.Run()

		//step3: 启动驱动服务, WithMqtt: 对接MQTT协议设备
		//Start: blocked
		//注：simpleDriver 必须实现接口 `type DPModelDriver interface

		mqttOptions := service.WithMqtt(dDpDriver.GetMqtt(), dDpDriver.GetMqttUsername(), dDpDriver.GetMqtt().ConnectHandler)
		err := ds.Start(dDpDriver, mqttOptions)
		if err != nil {
			sdkLog.Errorf("dp driver:%s start error:%s", serviceName, err)
			panic(fmt.Sprintf("dp driver:%s start error:%s", serviceName, err))
		}
		return
	}

	// should never reach here
	sdkLog.Errorf("tEdgeModel:%s, driver not support, panic", tEdgeModel)
	panic(fmt.Sprintf("tEdgeModel:%s, driver not support, panic", tEdgeModel))
}
