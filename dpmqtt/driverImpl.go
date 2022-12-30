package dpmqtt

import (
	"context"
	"encoding/json"

	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
	"github.com/tuya/tuya-tedge-driver-sdk-go/dpmodel"
	"github.com/tuya/tuya-tedge-driver-sdk-go/service"
)

/*
1. DemoMqttDPDriver 必须实现接口 `type DPModelDriver interface`
   接口定义：`tedge-driver-sdk-go/dpmodel/interface.go`
*/
type DemoMqttDPDriver struct {
	logger     commons.TedgeLogger
	dpmDS      *service.DPDriverService
	driverConf DriverConfig
	mqttItf    *MqttItf                 //处理来自设备的MQTT消息
}

func NewDemoMqttDPDriver(dpDs *service.DPDriverService) *DemoMqttDPDriver {
	mqttDp := &DemoMqttDPDriver{
		dpmDS:     dpDs,
		logger:    dpDs.GetLogger(),
	}

	mqttDp.mqttItf = NewMqttItf(mqttDp)
	return mqttDp
}

//////////////////////////////////////////////////////////////////////////////////////////
//1.接收 Tedge/云端 下发的MQTT消息
//2.注意：不要在该接口做阻塞性操作
func (dd *DemoMqttDPDriver) HandleCommands(ctx context.Context, cid string, msg dpmodel.CommandRequest, protocols map[string]commons.ProtocolProperties, dpExtend dpmodel.DPExtendInfo) error {
	lc := dd.logger
	msgStr, _ := json.Marshal(msg)
	lc.Infof("HandleCommands cid:%s, dpMessage:%s", cid, msgStr)

	//TODO: implement me
	//...

	return nil
}

//1.在Tedge控制台页面，新增、激活、更新子设备属性、删除子设备时，回调该接口
//2.注意：不要在该接口做阻塞性操作
//3.如果接入的设备不需要在Tedge控制台页面手动新增子设备，则该接口实现为空即可
func (dd *DemoMqttDPDriver) DeviceNotify(ctx context.Context, action commons.DeviceNotifyType, cid string, device commons.DeviceInfo) error {
	lc := dd.logger
	deviceStr, _ := json.Marshal(device)
	lc.Infof("DeviceNotify action:%s cid:%s: deviceStr:%s", action, cid, deviceStr)

	//TODO: implement me
	//...

	return nil
}

//1.ProductNotify 产品增删改通知,删除产品时product参数为空
//2.注意：不要在该接口做阻塞性操作
func (dd *DemoMqttDPDriver) ProductNotify(ctx context.Context, t commons.ProductNotifyType, pid string, product dpmodel.DPModelProduct) error {
	dd.logger.Infof("ProductNotify pid:%s", pid)
	return nil
}

//在Tedge页面，更新驱动实例，或停止驱动实例时，回调该接口，驱动程序进行资源清理
func (dd *DemoMqttDPDriver) Stop(ctx context.Context) error {
	dd.logger.Infof("Stop in...")
	return nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////
func (dd *DemoMqttDPDriver) Run() error {
	//1.获取驱动专家模式自定义配置
	driverConfigByte, _ := json.Marshal(dd.dpmDS.GetCustomConfig())
	dd.logger.Infof("DemoMqttDPDriver Run driverConfigByte:%s", driverConfigByte)

	err := json.Unmarshal(driverConfigByte, &dd.driverConf)
	dd.logger.Infof("DemoMqttDPDriver dd.driverConf:%+v, err:%v", dd.driverConf, err)

	//TODO: implement me
	//...

	return nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////
func (dd *DemoMqttDPDriver) GetMqtt() *MqttItf {
	return dd.mqttItf
}

func (dd *DemoMqttDPDriver) GetMqttUsername() string {

	//TODO: implement me
	//return dd.driverConf.MqttConf.Username

	return "username"
}
