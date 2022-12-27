package tydriver

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

////////////////////////////////////////////////////////////////////////////////////////
type DriverConfig struct {
	HttpConf   HttpConfig   `yaml:"httpconf"`
	DeviceConf DeviceConfig `yaml:"deviceconf"`
}

type HttpConfig struct {
	Listen string `yaml:"listen"`
}

type DeviceConfig struct {
	TestKey1 string `yaml:"testkey1"`
	TestKey2 string `yaml:"testkey2"`
}

// 数据上下文
type HttpSrv struct {
	listen string
	tyd    *DemoTyDriver
}

func NewTestHttpSrv(httpConf HttpConfig, sd *DemoTyDriver) *HttpSrv {
	httpSrv := &HttpSrv{
		listen: httpConf.Listen,
		tyd:    sd,
	}
	sd.logger.Infof("NewHttpSrv listen:%s done", httpSrv.listen)

	return httpSrv
}

func (ctx *HttpSrv) Run() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	//router.POST("/device/status", ctx.postDeviceStatus)                // 置指定子设备在线或离线
	//router.POST("/device/alert", ctx.httpSdkAlert)                     // 告警
	//router.POST("/device/report/event", ctx.postReportEvent)           // 上报事件
	//router.POST("/device/report/property", ctx.postPropertyReport)     // 上报属性
	//router.POST("/device/report/coordinate", ctx.postReportCoordinate) // 上报坐标，结构型

	//router.POST("/storage/v2/putdata", ctx.putCustomStorageV2)
	//router.GET("/storage/v2/getdata", ctx.getCustomStorageV2)
	//router.DELETE("/storage/v2/deletedata", ctx.deleteCustomStorageV2)
	//router.GET("/storage/v2/querydata", ctx.queryStroageV2)
	//router.GET("/storage/v2/getkeys", ctx.queryStorageKeysV2)

	ctx.tyd.logger.Infof("NewHttpSrv Run Run...")
	if err := router.Run(ctx.listen); err != nil {
		panic(fmt.Sprintf("Run err:%s", err))
	}
}

//////////////////////////////////////////////////////////////////////////////////////
func GinResult(message string, data interface{}) *gin.H {
	return &gin.H{
		"result": map[string]interface{}{
			"message": message,
			"data":    data,
		},
	}
}

// Test: 手动设置一个子设备的状态
type DeviceStatus struct {
	Cid    string `form:"cid"`
	Status int32  `form:"status"`
}

//func (ctx *HttpSrv) postDeviceStatus(c *gin.Context) {
//	var devStatus DeviceStatus
//	err := c.Bind(&devStatus)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, GinResult("fail", err))
//		return
//	}
//
//	if len(devStatus.Cid) == 0 {
//		c.JSON(http.StatusBadRequest, GinResult("fail", "param no cid"))
//		return
//	}
//
//	//status:0 offline; status:1 online
//	if devStatus.Status == 0 {
//		ctx.tyd.TestDeviceOffline(devStatus.Cid)
//	} else {
//		ctx.tyd.TestDeviceOnline(devStatus.Cid)
//	}
//
//	ctx.tyd.logger.Infof("postDeviceStatus cid:%s, Status:%v", devStatus.Cid, devStatus.Status)
//	c.JSON(http.StatusOK, GinResult("ok", err))
//}

// 触发一个驱动告警
type AlertMessage struct {
	Level string `form:"level"`
	Msg   string `form:"msg"`
}

//func (ctx *HttpSrv) httpSdkAlert(c *gin.Context) {
//	var alertMsg AlertMessage
//	err := c.Bind(&alertMsg)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, GinResult("fail", err))
//		return
//	}
//
//	err = ctx.tyd.ReportAlert(alertMsg.Level, alertMsg.Msg)
//	ctx.tyd.logger.Debugf("1.httpSdkAlarm ReportAlarm err:%v", err)
//
//	c.JSON(http.StatusOK, GinResult("ok", err))
//}

// 触发一个事件上报
type EventMessage struct {
	Cid string `form:"cid"`
	Msg string `form:"msg"`
}

//func (ctx *HttpSrv) postReportEvent(c *gin.Context) {
//	var eventMsg EventMessage
//	err := c.Bind(&eventMsg)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, GinResult("fail", err))
//		return
//	}
//
//	ctx.tyd.ReportEvent(eventMsg.Cid, eventMsg.Msg)
//	ctx.tyd.logger.Debugf("1.postReportEvent ReportEvent done.")
//
//	c.JSON(http.StatusOK, GinResult("ok", err))
//}

////PropertyReport
//func (ctx *HttpSrv) postPropertyReport(c *gin.Context) {
//	var propertyMsg EventMessage
//	err := c.Bind(&propertyMsg)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, GinResult("fail", err))
//		return
//	}
//
//	ctx.tyd.PropertyReport(propertyMsg.Cid, propertyMsg.Msg)
//	ctx.tyd.logger.Debugf("1.postPropertyReport PropertyReport done.")
//
//	c.JSON(http.StatusOK, GinResult("ok", err))
//}
//
////PropertyReportCoordinate
//func (ctx *HttpSrv) postReportCoordinate(c *gin.Context) {
//	cid := c.Query("cid")
//	if len(cid) <= 0 {
//		c.JSON(http.StatusBadRequest, GinResult("failed", "no cid"))
//		return
//	}
//
//	err := ctx.tyd.PropertyReportCoordinate(cid)
//	ctx.tyd.logger.Debugf("1.postReportCoordinate PropertyReportCoordinate done.")
//
//	c.JSON(http.StatusOK, GinResult("ok", err))
//}

//func (ctx *HttpSrv) putCustomStorageV2(c *gin.Context) {
//	type PutReq struct {
//		Key   string      `json:"key"`
//		Value interface{} `json:"value"`
//	}
//
//	var (
//		err    error
//		status = http.StatusOK
//	)
//	defer func() {
//		c.JSON(status, GinResult("ok", err))
//	}()
//
//	req := &PutReq{}
//	if err = c.BindJSON(req); err != nil {
//		status = http.StatusBadRequest
//		return
//	}
//
//	if req.Key == "" || req.Value == nil {
//		status = http.StatusBadRequest
//		return
//	}
//
//	value, _ := json.Marshal(req.Value)
//
//	err = ctx.tyd.PutCustomStorageV2(req.Key, value)
//}
//
//func (ctx *HttpSrv) getCustomStorageV2(c *gin.Context) {
//	var (
//		err    error
//		status = http.StatusOK
//		result map[string]interface{}
//	)
//
//	keysStr := c.Query("keys")
//
//	keys := strings.Split(keysStr, ",")
//	if len(keys) == 1 && keys[0] == "" {
//		keys = []string{}
//	}
//	fmt.Printf("keys %v, len %d\n", keys, len(keys))
//	result, err = ctx.tyd.GetCustomStorageV2(keys)
//	if err != nil {
//		status = http.StatusInternalServerError
//	}
//
//	c.JSON(status, GinResult("ok", result))
//}
//
//func (ctx *HttpSrv) deleteCustomStorageV2(c *gin.Context) {
//	cids := []string{}
//	keysStr := c.Query("keys")
//	if keysStr != "" {
//		cids = strings.Split(keysStr, ",")
//	}
//	err := ctx.tyd.DeleteCustomStorageV2(cids)
//
//	c.JSON(http.StatusOK, GinResult("ok", err))
//}
//
//func (ctx *HttpSrv) queryStroageV2(c *gin.Context) {
//	prefix := c.Query("prefix")
//	fmt.Printf("prefix %s", prefix)
//	m, _ := ctx.tyd.QueryCustomStorageV2(prefix)
//	c.JSON(http.StatusOK, GinResult("ok", m))
//}
//
//func (ctx *HttpSrv) queryStorageKeysV2(c *gin.Context) {
//	prefix := c.Query("prefix")
//	fmt.Printf("prefix %s", prefix)
//	m, _ := ctx.tyd.QueryCustomStorageKeys(prefix)
//	c.JSON(http.StatusOK, GinResult("ok", m))
//}
