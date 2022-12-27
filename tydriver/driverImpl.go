package tydriver

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
	"github.com/tuya/tuya-tedge-driver-sdk-go/service"
	"github.com/tuya/tuya-tedge-driver-sdk-go/thingmodel"
)

/*
1. DemoTyDriver 必须实现接口 `type ThingModelDriver interface`
   接口定义：`tedge-driver-sdk-go/thingsmodel/interface.go`
*/
type DemoTyDriver struct {
	mux       sync.RWMutex
	logger    commons.TedgeLogger
	tymDS     *service.TyDriverService
	driConfig map[string]interface{}   //驱动自定义配置，yaml格式(专家模式)
	deviceMap map[string]*DeviceShadow //子设备列表
}

func NewDemoTyDriver(tymDS *service.TyDriverService) *DemoTyDriver {
	return &DemoTyDriver{
		tymDS:     tymDS,
		logger:    tymDS.GetLogger(),
		deviceMap: make(map[string]*DeviceShadow),
	}
}

//////////////////////////////////////////////////////////////////////////////////////////
//1.在 Tedge 控制台页面，新增、激活、更新子设备属性、删除子设备时，回调该接口
//2.注意：不要在该接口做阻塞性操作
//3.如果接入的设备不需要在 Tedge 控制台页面手动新增子设备，则该接口实现为空即可
func (tyd *DemoTyDriver) DeviceNotify(ctx context.Context, action commons.DeviceNotifyType, cid string, device commons.TMDeviceInfo) error {
	lc := tyd.logger
	deviceStr, _ := json.Marshal(device)
	lc.Infof("DeviceNotify action:%s cid:%s: deviceStr:%s", action, cid, deviceStr)

	if device.ActiveStatus != commons.DeviceActiveStatusActivated {
		lc.Warnf("DeviceAddNotify cid:%s AddDevice failed, not activated", cid)
		return nil
	}

	switch action {
	case commons.DeviceActiveNotify, commons.DeviceAddNotify, commons.DeviceUpdateNotify:
		tyd.DelDeviceShadow(cid)

		virDev := tyd.SetDeviceShadow(device)
		go virDev.Run()
	case commons.DeviceDeleteNotify:
		tyd.DelDeviceShadow(cid)
	}
	return nil
}

//注意：不要在该接口做阻塞性操作
func (tyd *DemoTyDriver) ProductNotify(ctx context.Context, t commons.ProductNotifyType, pid string, product thingmodel.ThingModelProduct) error {
	productStr, _ := json.Marshal(product)
	tyd.logger.Debugf("ProductNotify pid:%s action:%s: productStr:%s", pid, t, productStr)

	return nil
}

//在Tedge页面，更新驱动实例，或停止驱动实例时，回调该接口，驱动程序进行资源清理
func (tyd *DemoTyDriver) Stop(ctx context.Context) error {
	tyd.logger.Infof("Stop in...")
	return nil
}

// HandlePropertySet 云端向设备发送属性设置消息
// 注意：不要在该接口做阻塞性操作
func (tyd *DemoTyDriver) HandlePropertySet(ctx context.Context, cid string, msg thingmodel.PropertySet, protocols map[string]commons.ProtocolProperties) error {
	msgStr, _ := json.Marshal(msg)
	tyd.logger.Debugf("HandlePropertySet msg: cid: %s, propertySet:%s", cid, msgStr)

	virDev, ok := tyd.GetDeviceShadow(cid)
	if !ok {
		tyd.logger.Errorf("HandlePropertySet cid:%s, not found device, ignore!", cid)
		return nil
	}

	//don't block
	go virDev.ProcessPropertySet(msg)

	return nil
}

// HandlePropertyGet 云端主动向设备发起属性查询的请求
//注意：不要在该接口做阻塞性操作
func (tyd *DemoTyDriver) HandlePropertyGet(ctx context.Context, cid string, data thingmodel.PropertyGet, protocols map[string]commons.ProtocolProperties) error {
	dataStr, _ := json.Marshal(data)
	tyd.logger.Debugf("HandlePropertyGet msg: cid: %s, propertyGet:%s", cid, dataStr)

	return nil
}

// HandleActionExecute 云端向设备发送动作执行消息
//注意：不要在该接口做阻塞性操作
func (tyd *DemoTyDriver) HandleActionExecute(ctx context.Context, cid string, data thingmodel.ActionExecuteRequest, protocols map[string]commons.ProtocolProperties) error {
	dataStr, _ := json.Marshal(data)
	tyd.logger.Debugf("HandlePropertyGet msg: cid: %s, actionCmd:%s", cid, dataStr)

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////
func (tyd *DemoTyDriver) GetDeviceShadow(cid string) (*DeviceShadow, bool) {
	tyd.mux.Lock()
	defer tyd.mux.Unlock()

	virDev, ok := tyd.deviceMap[cid]
	return virDev, ok
}

func (tyd *DemoTyDriver) SetDeviceShadow(device commons.TMDeviceInfo) *DeviceShadow {
	tyd.mux.Lock()
	defer tyd.mux.Unlock()

	newVirDev := NewDeviceShadow(&device, tyd)
	tyd.deviceMap[device.Cid] = newVirDev
	return newVirDev
}

func (tyd *DemoTyDriver) DelDeviceShadow(cid string) {
	tyd.mux.Lock()
	defer tyd.mux.Unlock()

	virDev, ok := tyd.deviceMap[cid]
	if ok {
		virDev.Stop()
		delete(tyd.deviceMap, cid)
	}
}

////////////////////////////////////////////////////////////////////////////////////////
func (tyd *DemoTyDriver) Run() error {
	lc := tyd.logger
	//1.获取驱动专家模式自定义配置
	tyd.driConfig = tyd.tymDS.GetCustomConfig()
	driverConfigByte, _ := json.Marshal(tyd.driConfig)
	lc.Infof("Run driverConfigByte:%s", driverConfigByte)

	//2.初始化子设备列表
	go InitSubDevices(tyd)

	//3.
	var driverConf DriverConfig
	err := json.Unmarshal(driverConfigByte, &driverConf)
	lc.Infof("Initialize templateConf:%+v, err:%v", driverConf, err)

	if len(driverConf.HttpConf.Listen) != 0 {
		go NewTestHttpSrv(driverConf.HttpConf, tyd).Run()
	}

	return nil
}
