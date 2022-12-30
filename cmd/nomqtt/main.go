package main

import (
	"fmt"

	"github.com/tuya/tuya-tedge-driver-example/dpdriver"
	"github.com/tuya/tuya-tedge-driver-example/tydriver"
	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
	"github.com/tuya/tuya-tedge-driver-sdk-go/service"
)

const (
	serviceName string = "driver-example"
)

//TEdge 驱动开发示例：DP模型、tyLink 模型二合一示例
//TEdge 边缘网关运行模式：1.DP模型(commons.DPModel); 2.tyLink模型(commons.ThingModel)
//驱动程序根据实际情况只实现任一模式即可
func main() {
	sdkLog := commons.DefaultLogger(commons.DebugLevel, serviceName)
	base := service.NewBaseService(sdkLog)

	// 1. DP 模型驱动
	tEdgeModel := base.GetTEdgeModel()
	if tEdgeModel == commons.DPModel {
		//step1: 创建驱动服务
		ds := service.NewDPServiceWithBase(base)
		sdkLog.Infof("TEdge run in dp model")

		//step2: 实现驱动接口 `DPModelDriver`
		dDpDriver := dpdriver.NewDemoDPDriver(ds)
		go dDpDriver.Run()

		//step3: 启动驱动服务
		//Start: blocked
		//注：simpleDriver 必须实现接口 `type DPModelDriver interface`
		err := ds.Start(dDpDriver)
		if err != nil {
			sdkLog.Errorf("dp driver:%s start error:%s", serviceName, err)
			panic(fmt.Sprintf("dp driver:%s start error:%s", serviceName, err))
		}
		return
	}

	if tEdgeModel == commons.ThingModel { // 2. tyLink 模型驱动
		ds := service.NewTyServiceWithBase(base)
		sdkLog.Infof("TEdge run in tyLink model")

		dTyDriver := tydriver.NewDemoTyDriver(ds)
		go dTyDriver.Run()

		//Start: blocked
		err := ds.Start(dTyDriver)
		if err != nil {
			sdkLog.Errorf("tyLink driver:%s start error:%s", serviceName, err)
			panic(fmt.Sprintf("tylink driver:%s start error:%s", serviceName, err))
		}
		return
	}

	// should never reach here
	sdkLog.Errorf("tEdgeModel:%s, driver not support, panic", tEdgeModel)
	panic(fmt.Sprintf("tEdgeModel:%s, driver not support, panic", tEdgeModel))
}
