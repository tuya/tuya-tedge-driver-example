package tydriver

import (
	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
	"github.com/tuya/tuya-tedge-driver-sdk-go/thingmodel"
)

const (
	TestTyPid1     = "8drwkxvkvxv6wwou" //演示产品pid
	TestTyDevCid1  = "test_tycid_0001"  //演示子设备cid
	TestTyDevName1 = "Ty子设备名1"        //演示子设备设备名
)

////////////////////////////////////////////////////////////////////////////////////////////
func InitSubDevices(ddp *DemoTyDriver) {
	//Test: 新增一个pid
	TestAddProduct(ddp)

	//Test: 新增并激活一个子设备
	//注：要激活一个 tylink 子设备，必须在iot平台获取/购买pid对应的授权码，一个授权码可激活一个子设备
	TestAddDevice(ddp)

	//1.获取子设备列表，并激活; 若子设备已经激活，则不要再次激活
	//调用设备提供的接口；第三方系统推送；或者其它任何方式
	//TODO: implement me
	//......

	//2.运行子设备
	//驱动每次重启后，可以从sdk中获取所有已经激活过的子设备
	devices := ddp.tymDS.GetActiveDevices()
	for cid, dev := range devices {
		_, ok := ddp.GetDeviceShadow(cid)
		if ok {
			continue
		}

		virDev := ddp.SetDeviceShadow(dev)
		go virDev.Run()
	}
}

//for test: 新增一个产品(pid)
func TestAddProduct(ddp *DemoTyDriver) {
	//新增"产品(pid)"有两种方式，pid必须在iot平台上创建
	//1.在Tedge控制台页面，手动新增
	//2.在驱动中调用sdk接口新增

	dpp1 := thingmodel.AddProductReq{
		Id:          TestTyPid1,
		Name:        "演示设备",
		Description: "This a test product",
	}
	ddp.tymDS.AddProduct(dpp1)
}

//for test: 新增并激活一个子设备，每个子设备必须绑定一个pid
func TestAddDevice(ddp *DemoTyDriver) error {
	extendData := make(map[string]interface{})
	extendData["testExt1"] = "test ext value1"

	deviceInfo := &commons.TMDeviceMeta{
		Cid:       TestTyDevCid1, //子设备cid，必填，网关下必须唯一
		ProductId: TestTyPid1,    //子设备pid：必填
		BaseAttr: commons.BaseProperty{
			Name: TestTyDevName1, //子设备名，必填
		},
		ExtendedAttr: commons.ExtendedProperty{
			InstallLocation: "设备安装地址",   //可为空
			ExtendData:      extendData, //子设备扩展字段，根据实际情况填写，可为空
		},
	}

	//新增并激活一个子设备
	err := ddp.tymDS.ActiveDevice(*deviceInfo)
	if err != nil {
		ddp.logger.Warnf("TestAddDevice ActiveDevice err:%s", err)
		return err
	}

	ddp.logger.Infof("TestAddDevice success, cid:%s, productId:%s, deviceName:%s", deviceInfo.Cid, deviceInfo.ProductId, deviceInfo.BaseAttr.Name)
	return nil
}

//////////////////////////////////////////////////////////////////////////////////////
////测试添加配置V2
//func (tyd *DemoTyDriver) PutCustomStorageV2(key string, value []byte) error {
//	return tyd.dpSdk.PutCustomStorage(map[string][]byte{
//		key: value,
//	})
//}
//
////测获取配置V2
//func (tyd *DemoTyDriver) GetCustomStorageV2(keys []string) (map[string]interface{}, error) {
//	storageMap, err := tyd.dpSdk.GetCustomStorage(keys)
//	if err != nil {
//		return nil, err
//	}
//
//	result := map[string]interface{}{}
//	for key, value := range storageMap {
//		result[key] = string(value)
//	}
//
//	return result, nil
//}

////删除存储信息
//func (tyd *DemoTyDriver) DeleteCustomStorageV2(keys []string) error {
//	return tyd.dpSdk.DeleteCustomStorage(keys)
//}
//
////查询存储信息
//func (tyd *DemoTyDriver) QueryCustomStorageV2(prefix string) (map[string][]byte, error) {
//	return tyd.dpSdk.QueryCustomStorage(prefix)
//}
//
////查询存储的key
//func (tyd *DemoTyDriver) QueryCustomStorageKeys(prefix string) ([]string, error) {
//	return tyd.dpSdk.GetCustomStorageKeys(prefix)
//}
