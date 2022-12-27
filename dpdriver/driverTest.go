package dpdriver

import (
	"context"
	"fmt"

	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
	"github.com/tuya/tuya-tedge-driver-sdk-go/dpmodel"
)

const (
	TestPid1     = "en1agoqpbynmgb5l" //演示产品pid：霍尼韦尔能耗系统
	TestDevCid1  = "test_cid_0001"    //演示子设备cid
	TestDevName1 = "子设备名1"            //演示子设备设备名
)

////////////////////////////////////////////////////////////////////////////////////////////
func InitSubDevices(ddp *DemoDPDriver) {
	//Test: 新增一个pid
	TestAddProduct(ddp)

	//Test: 新增并激活一个子设备
	TestAddDevice(ddp)

	//1.获取真实子设备列表，并依次激活; 若子设备已经激活过，则不要再次激活
	//必填字段：cid: 子设备id，网关下唯一，cid可以是设备SN、MAC等
	//必填字段：pid: 每个子设备必须绑定一个pid，在iot平台创建
	//必填字段：name: 子设备名
	//调用设备提供的接口；第三方系统推送；或者其它任何方式
	//TODO: implement me
	//......

	//2.运行子设备
	//驱动每次重启后，可以从sdk中获取所有已经激活过的子设备
	devices := ddp.dpmDS.GetActiveDevices()
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
func TestAddProduct(ddp *DemoDPDriver) {
	//新增"产品(pid)"有两种方式，pid必须在iot平台上创建
	//方式1：在Tedge控制台页面，手动新增
	//方式2：在驱动中调用sdk接口新增

	dpp1 := dpmodel.DPModelProductAddInfo{
		Id:          TestPid1,
		Name:        "演示设备",
		Description: "This a test product",
	}
	ddp.dpmDS.AddProduct(dpp1)
}

//for test: 新增并激活一个子设备
func TestAddDevice(ddp *DemoDPDriver) error {
	//设备自定义扩展属性：json
	extendData := make(map[string]interface{})
	extendData["testExt1"] = "test ext value1"

	deviceInfo := &commons.DeviceMeta{
		Cid:       TestDevCid1, //子设备cid，必填，网关下必须唯一
		ProductId: TestPid1,    //子设备pid：必填
		BaseAttr: commons.BaseProperty{
			Name: TestDevName1, //子设备名，必填
		},
		ExtendedAttr: commons.ExtendedProperty{
			InstallLocation: "设备安装地址",   //可为空
			ExtendData:      extendData, //子设备扩展字段，根据实际情况填写，可为空
		},
	}

	//新增并激活一个子设备
	err := ddp.dpmDS.ActiveDevice(*deviceInfo)
	if err != nil {
		//激活子设备失败，向TEdge上报一条告警
		ddp.dpmDS.ReportAlert(context.Background(), commons.WARN, fmt.Sprintf("active device cid:%s fail", deviceInfo.Cid))
		ddp.logger.Warnf("TestAddDevice ActiveDevice cid:%s, err:%s", deviceInfo.Cid, err)
		return err
	}

	ddp.logger.Infof("TestAddDevice success, cid:%s, productId:%s, deviceName:%s", deviceInfo.Cid, deviceInfo.ProductId, deviceInfo.BaseAttr.Name)
	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
//for test: 上传一张图片
func TestUploadImage(ddp *DemoDPDriver, cid string) error {
	//上传一张图片到OSS
	content := []byte("XXXX") //Read content from picture

	resp, err := ddp.dpmDS.UploadFile(content, cid, "", 5)
	ddp.logger.Debugf("TestUploadImage UploadFile resp:%s, err:%v", resp, err)

	return nil
}

//for test: 更新设备名；设备名变化时
func TestModifyName(ddp *DemoDPDriver, cid, name string) error {
	baseAttr := commons.BaseProperty{
		Name: name,
	}

	err := ddp.dpmDS.SetDeviceBaseAttr(cid, baseAttr)
	ddp.logger.Infof("ModifyDeviceBaseAttr cid:%s, name:%s, err:%v", cid, name, err)

	return err
}

//for test: 更新设备扩展属性
func TestModifyExt(ddp *DemoDPDriver, cid string) error {
	extendData := make(map[string]interface{})
	extendData["testExt1"] = "modify test ext ext1"
	extendData["testExt2"] = "modify test ext ext2"

	property := commons.ExtendedProperty{
		InstallLocation: "新的安装地址",
		ExtendData:      extendData,
	}

	err := ddp.dpmDS.SetDeviceExtendProperty(cid, property)
	return err
}

///////////////////////////////////////////////////////////////////////////////////////////////
////测试添加配置V2
//func (demoDriver *DemoDPDriver) PutCustomStorageV2(key string, value []byte) error {
//	return demoDriver.dpmDS.PutCustomStorageV2(map[string][]byte{
//		key: value,
//	})
//}
//
////测获取配置V2
//func (demoDriver *DemoDPDriver) GetCustomStorageV2(keys []string) (map[string]interface{}, error) {
//	storageMap, err := demoDriver.dpmDS.GetCustomStorageV2(keys)
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
//
////删除存储信息
//func (demoDriver *DemoDPDriver) DeleteCustomStorageV2(keys []string) error {
//	return demoDriver.dpmDS.DeleteCustomStorageV2(keys)
//}
//
////获取所有存储信息
//func (demoDriver *DemoDPDriver) QueryCustomStorageV2(prefix string) (map[string]interface{}, error) {
//	storageMap, err := demoDriver.dpmDS.QueryCustomStorageV2(prefix)
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
//
////获取keys
//func (demoDriver *DemoDPDriver) GetCustomStorageKeysV2(prefix string) ([]string, error) {
//	return demoDriver.dpmDS.GetCustomStorageKeysV2(prefix)
//}
