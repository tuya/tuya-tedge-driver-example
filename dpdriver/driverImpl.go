package dpdriver

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
	"github.com/tuya/tuya-tedge-driver-sdk-go/dpmodel"
	"github.com/tuya/tuya-tedge-driver-sdk-go/service"
)

/*
1. DemoDpDriver 必须实现接口 `type DPModelDriver interface`
   接口定义：`tedge-driver-sdk-go/dpmodel/interface.go`
*/
type DemoDPDriver struct {
	mux        sync.RWMutex
	logger     commons.TedgeLogger
	dpmDS      *service.DPDriverService
	driverConf DriverConfig
	deviceMap  map[string]*DeviceShadow //子设备列表
}

func NewDemoDPDriver(dpDs *service.DPDriverService) *DemoDPDriver {
	return &DemoDPDriver{
		dpmDS:     dpDs,
		logger:    dpDs.GetLogger(),
		deviceMap: make(map[string]*DeviceShadow),
	}
}

//////////////////////////////////////////////////////////////////////////////////////////
//1.接收 Tedge/云端 下发的MQTT消息
//2.注意：不要在该接口做阻塞性操作
func (dd *DemoDPDriver) HandleCommands(ctx context.Context, cid string, msg dpmodel.CommandRequest, protocols map[string]commons.ProtocolProperties, dpExtend dpmodel.DPExtendInfo) error {
	lc := dd.logger
	msgStr, _ := json.Marshal(msg)
	lc.Infof("HandleCommands cid:%s, dpMessage:%s", cid, msgStr)

	virDev, ok := dd.GetDeviceShadow(cid)
	if !ok {
		dd.logger.Errorf("HandleCommands cid:%s, not found device, ignore!", cid)
		return nil
	}

	//don't block
	go virDev.ProcessDpMessage(msg)

	return nil
}

//1.在Tedge控制台页面，新增、激活、更新子设备属性、删除子设备时，回调该接口
//2.注意：不要在该接口做阻塞性操作
//3.如果接入的设备不需要在Tedge控制台页面手动新增子设备，则该接口实现为空即可
func (dd *DemoDPDriver) DeviceNotify(ctx context.Context, action commons.DeviceNotifyType, cid string, device commons.DeviceInfo) error {
	lc := dd.logger
	deviceStr, _ := json.Marshal(device)
	lc.Infof("DeviceNotify action:%s cid:%s: deviceStr:%s", action, cid, deviceStr)

	switch action {
	case commons.DeviceActiveNotify, commons.DeviceAddNotify, commons.DeviceUpdateNotify:
		if device.ActiveStatus != commons.DeviceActiveStatusActivated {
			lc.Warnf("DeviceAddNotify cid:%s AddDevice failed, not activated", cid)
			return nil
		}

		dd.DelDeviceShadow(cid)
		virDev := dd.SetDeviceShadow(device)
		go virDev.Run()
	case commons.DeviceDeleteNotify:
		dd.DelDeviceShadow(cid)
	}
	return nil
}

//1.ProductNotify 产品增删改通知,删除产品时product参数为空
//2.注意：不要在该接口做阻塞性操作
func (dd *DemoDPDriver) ProductNotify(ctx context.Context, t commons.ProductNotifyType, pid string, product dpmodel.DPModelProduct) error {
	dd.logger.Infof("ProductNotify pid:%s", pid)
	return nil
}

//在Tedge页面，更新驱动实例，或停止驱动实例时，回调该接口，驱动程序进行资源清理
func (dd *DemoDPDriver) Stop(ctx context.Context) error {
	dd.logger.Infof("Stop in...")
	return nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////
func (dd *DemoDPDriver) GetDeviceShadow(cid string) (*DeviceShadow, bool) {
	dd.mux.Lock()
	defer dd.mux.Unlock()

	virDev, ok := dd.deviceMap[cid]
	return virDev, ok
}

func (dd *DemoDPDriver) SetDeviceShadow(device commons.DeviceInfo) *DeviceShadow {
	dd.mux.Lock()
	defer dd.mux.Unlock()

	newVirDev := NewVirtualDevice(&device, dd)
	dd.deviceMap[device.Cid] = newVirDev
	return newVirDev
}

func (dd *DemoDPDriver) DelDeviceShadow(cid string) {
	dd.mux.Lock()
	defer dd.mux.Unlock()

	virDev, ok := dd.deviceMap[cid]
	if ok {
		virDev.Stop()
		delete(dd.deviceMap, cid)
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////////
func (dd *DemoDPDriver) Run() error {
	//1.获取驱动专家模式自定义配置
	driverConfigByte, _ := json.Marshal(dd.dpmDS.GetCustomConfig())
	dd.logger.Infof("DemoDPDriver Run driverConfigByte:%s", driverConfigByte)

	err := json.Unmarshal(driverConfigByte, &dd.driverConf)
	dd.logger.Infof("DemoDPDriver dd.driverConf:%+v, err:%v", dd.driverConf, err)

	//2.初始化子设备列表
	go InitSubDevices(dd)

	if len(dd.driverConf.HttpConf.Listen) != 0 {
		go NewTestHttpSrv(dd.driverConf.HttpConf, dd).Run()
	}

	return nil
}
