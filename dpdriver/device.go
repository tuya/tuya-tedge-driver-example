package dpdriver

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
	"github.com/tuya/tuya-tedge-driver-sdk-go/dpmodel"
)

//1.每个子设备必须有cid，cid在该网关下必须唯一
//2.cid可以是设备SN、MAC等
type DeviceShadow struct {
	Cid    string
	Name   string
	driver *DemoDPDriver
	done   chan struct{}
}

func NewVirtualDevice(dev *commons.DeviceInfo, driver *DemoDPDriver) *DeviceShadow {
	virtualDev := &DeviceShadow{
		Cid:    dev.Cid,
		Name:   dev.BaseAttr.Name,
		driver: driver,
		done:   make(chan struct{}),
	}

	return virtualDev
}

func (vd *DeviceShadow) Run() {
	lc := vd.driver.logger
	lc.Infof("vd Run in, cid:%s", vd.Cid)
	vd.CheckDevStatus(true)

	tickerA := time.NewTicker(60 * time.Second)
	tickerB := time.NewTicker(60 * time.Second)

	for {
		select {
		case <-tickerA.C:
			//演示：定时检测设备在线/离线状态，并向云端上报
			vd.CheckDevStatus(true)
		case <-tickerB.C:
			//演示：向云端上报设备属性值；根据实际情况上报
			vd.ReportDp()
		case <-vd.done:
			tickerA.Stop()
			tickerB.Stop()
			lc.Infof("vd Run cid:%s exit", vd.Cid)
			return
		}
	}
}

func (vd *DeviceShadow) Stop() {
	close(vd.done)
}

//1.定时检测并更新设备状态
func (vd *DeviceShadow) CheckDevStatus(online bool) {
	//检测设备在线状态, 设备状态由驱动负责维护;
	//1.调用设备提供的接口；2.或者通过其它方式判断设备状态; 3.驱动定时向Tedge上报设备状态
	//TODO: implement me
	isOnline := online

	//b.调用sdk api，向云端更新设备状态
	devStatus := &commons.DeviceStatus{}
	if isOnline {
		// 设备在线
		devStatus.Online = append(devStatus.Online, vd.Cid)
	} else {
		// 设备离线
		devStatus.Offline = append(devStatus.Offline, vd.Cid)
	}

	dpDs := vd.driver.dpmDS
	dpDs.ReportDeviceStatus(devStatus)
}

//2.定时上报一个DP消息
func (vd *DeviceShadow) ReportDp() {
	//示例:
	//1.假设设备的pid为：""，该 pid定义了4个DP点:
	// switch(16): bool类型，表示"开关状态"；
	// electric_total(46): value数值类型，表示当前"总电量"；
	// device_alarm(31): string字符类型，表示"设备告警"；
	// state(65): enum枚举类型(online, offline)，表示"在线状态"

	//获取设备指定属性当前值，并向云端上报
	//1.调用设备提供的接口；2.或由设备推送过来；3.或其它途径
	//TODO: implement me

	randI := rand.Int() % 2
	bStatus := true
	if randI == 0 {
		bStatus = false
	}
	vd.ReportDp16(bStatus)

	vd.ReportDp46(int64(randI))

	vd.ReportDp31("This is a test warning!")

	vd.ReportDp65("online")
}

//for test
func (vd *DeviceShadow) ReportDp16(bStatus bool) {
	bValue := dpmodel.NewWithDPValue("16", commons.BoolType, bStatus)

	var dPValues []*dpmodel.WithDPValue
	dPValues = append(dPValues, bValue)
	vd.driver.dpmDS.ReportWithDPData(vd.Cid, dPValues)
}

//for test
func (vd *DeviceShadow) ReportDp46(iValue int64) {
	intValue := dpmodel.NewWithDPValue("46", commons.ValueType, iValue)

	var dPValues []*dpmodel.WithDPValue
	dPValues = append(dPValues, intValue)
	vd.driver.dpmDS.ReportWithDPData(vd.Cid, dPValues)
}

//for test
func (vd *DeviceShadow) ReportDp31(str1 string) {
	strValue := dpmodel.NewWithDPValue("31", commons.StringType, str1)

	var dPValues []*dpmodel.WithDPValue
	dPValues = append(dPValues, strValue)
	vd.driver.dpmDS.ReportWithDPData(vd.Cid, dPValues)
}

//for test
func (vd *DeviceShadow) ReportDp65(eValue string) {
	//注：枚举值必须在iot平台pid中定义的保持一致
	//枚举值: online, offline,
	if eValue != "online" && eValue != "offline" {
		return
	}
	enumValue := dpmodel.NewWithDPValue("65", commons.EnumType, eValue)

	var dPValues []*dpmodel.WithDPValue
	dPValues = append(dPValues, enumValue)
	vd.driver.dpmDS.ReportWithDPData(vd.Cid, dPValues)
}

//3.处理云端或Tedge下发的DP消息
func (vd *DeviceShadow) ProcessDpMessage(msg dpmodel.CommandRequest) {
	cid := vd.Cid
	lc := vd.driver.logger

	dps, ok := msg.Data["dps"]
	if !ok {
		lc.Errorf("ProcessDpMessage cid:%s, msg.Data['dps'] not exist, can't happened", cid)
		return
	}

	dpsMap, ok := dps.(map[string]interface{})
	if !ok {
		lc.Errorf("ProcessDpMessage cid:%s, dps must be `map[string]interface{}`", cid)
		return
	}

	dpsMaps, _ := json.Marshal(dpsMap)
	lc.Infof("ProcessDpMessage cid:%s, dpsMaps:%s", cid, dpsMaps)

	//1.调用设备提供的接口，将云端下发的指令发送到设备
	//TODO: implement me
}
