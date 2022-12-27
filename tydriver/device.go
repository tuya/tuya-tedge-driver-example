package tydriver

import (
	"time"

	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
	"github.com/tuya/tuya-tedge-driver-sdk-go/thingmodel"
)

//1.每个子设备必须有cid，cid在该网关下必须唯一
//2.cid可以是设备SN、MAC等
type DeviceShadow struct {
	Cid    string
	Name   string
	driver *DemoTyDriver
	done   chan struct{}
}

func NewDeviceShadow(dev *commons.TMDeviceInfo, driver *DemoTyDriver) *DeviceShadow {
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
	lc.Infof("vd Run in...")

	tickerA := time.NewTicker(60 * time.Second)
	tickerB := time.NewTicker(60 * time.Second)

	for {
		select {
		case <-tickerA.C:
			//演示：定时检测设备在线/离线状态，并向云端上报
			vd.CheckDevStatus(true)
		case <-tickerB.C:
			//演示：向云端上报设备属性值；根据实际情况上报
			vd.ReportProperty()
			vd.ReportEvent()
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

	dpDs := vd.driver.tymDS
	dpDs.ReportDeviceStatus(devStatus)
}

//2.定时上报一个DP消息
func (vd *DeviceShadow) ReportProperty() {
	//示例:
	//1.假设设备的pid为：""，该 pid定义了4个属性点:
	// status(101): bool类型，表示开关状态；
	// temp(102): value数值类型，表示当前温度；
	// msg(103): string字符类型，表示自定义消息；
	// day(104): enum枚举类型，表示星期几

	//获取设备指定属性当前值，并向云端上报
	//1.调用设备提供的接口；2.或由设备推送过来；3.或其它途径
	//TODO: implement me

	//randI := rand.Int() % 2
	//bStatus := true
	//if randI == 0 {
	//	bStatus = false
	//}
	//vd.ReportStatus(bStatus)
	//
	//vd.ReportTemp(randI)
	//
	//vd.ReportMsg("This is a test msg!")
}

////for test
//func (vd *DeviceShadow) ReportStatus(bStatus bool) {
//	bValue := dpmodel.NewWithDPValue("101", commons.BoolType, bStatus)
//
//	var dPValues []*dpmodel.WithDPValue
//	dPValues = append(dPValues, bValue)
//	vd.driver.tymDS.ReportWithDPData(vd.Cid, dPValues)
//}
//
////for test
//func (vd *DeviceShadow) ReportTemp(iValue int) {
//	intValue := dpmodel.NewWithDPValue("102", commons.ValueType, iValue)
//
//	var dPValues []*dpmodel.WithDPValue
//	dPValues = append(dPValues, intValue)
//	vd.driver.tymDS.ReportWithDPData(vd.Cid, dPValues)
//}
//
////for test
//func (vd *DeviceShadow) ReportMsg(str1 string) {
//	strValue := dpmodel.NewWithDPValue("103", commons.StringType, str1)
//
//	var dPValues []*dpmodel.WithDPValue
//	dPValues = append(dPValues, strValue)
//	vd.driver.tymDS.ReportWithDPData(vd.Cid, dPValues)
//}
//
////for test
//func (vd *DeviceShadow) ReportDay(eValue string) {
//	//注：枚举值必须在iot平台pid中定义的保持一致
//	//枚举值: 1h, 2h,
//	if eValue != "1h" && eValue != "2h" {
//		return
//	}
//	enumValue := dpmodel.NewWithDPValue("104", commons.EnumType, eValue)
//
//	var dPValues []*dpmodel.WithDPValue
//	dPValues = append(dPValues, enumValue)
//	vd.driver.tymDS.ReportWithDPData(vd.Cid, dPValues)
//}

func (vd *DeviceShadow) ReportEvent() {

}

//3.处理云端或TEdge下发的属性设置消息
func (vd *DeviceShadow) ProcessPropertySet(msg thingmodel.PropertySet) {
	lc := vd.driver.logger
	lc.Infof("ProcessPropertySet cid:%s, MsgId:%s", vd.Cid, msg.MsgId)

	//for key, value := range msg.Data {
	//
	//}
}
