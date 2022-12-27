package dpdriver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

////////////////////////////////////////////////////////////////////////////////////////
type DriverConfig struct {
	HttpConf   HttpConfig   `yaml:"httpconf"`
}

type HttpConfig struct {
	Listen string `yaml:"listen"`
}

////////////////////////////////////////////////////////////////////////////////////////
// 数据上下文
type TestHttpSrv struct {
	listen     string
	demoDriver *DemoDPDriver
}

func NewTestHttpSrv(httpConf HttpConfig, sd *DemoDPDriver) *TestHttpSrv {
	httpSrv := &TestHttpSrv{
		listen:     httpConf.Listen,
		demoDriver: sd,
	}

	sd.logger.Infof("NewHttpSrv listen:%s done", httpSrv.listen)
	return httpSrv
}

func (srv *TestHttpSrv) Run() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	//api: 用来测试一些基本功能
	router.POST("/device/status", srv.postDeviceStatus)      // 设置子设备状态为在线或离线
	router.POST("/device/mod/name", srv.postModifyName)      // 修改子设备名
	router.POST("/device/mod/attr/:cid", srv.postModifyAttr) // 修改子设备扩展属性
	router.POST("/device/upload/:cid", srv.postUploadImage)  // 上传一张图片

	//router.POST("/device/postdp", srv.postDpValue) // 上报一个DP点

	//router.POST("/storage/v2/putdata", srv.putCustomStorageV2)
	//router.GET("/storage/v2/getdata", srv.getCustomStorageV2)
	//router.DELETE("/storage/v2/deletedata", srv.deleteCustomStorageV2)
	//router.GET("/storage/v2/querydata", srv.queryStroageV2)
	//router.GET("/storage/v2/getkeys", srv.getCustomStorageKeysV2)

	srv.demoDriver.logger.Infof("NewHttpSrv Run Run...")
	if err := router.Run(srv.listen); err != nil {
		srv.demoDriver.logger.Errorf("Run failed, listen:%s err:%s", srv.listen, err)
	}
}

// /////////////////////////////////////////////////////////////
func GinResult(message string, data interface{}) *gin.H {
	return &gin.H{
		"result": map[string]interface{}{
			"message": message,
			"data":    data,
		},
	}
}

//for test: 改变一个子设备的状态
type DeviceStatus struct {
	Cid    string `form:"cid"`
	Status int32  `form:"status"`
}

func (srv *TestHttpSrv) postDeviceStatus(c *gin.Context) {
	var devStatus DeviceStatus
	err := c.Bind(&devStatus)
	if err != nil {
		c.JSON(http.StatusBadRequest, GinResult("fail", err))
		return
	}

	cid := devStatus.Cid
	virDev, ok := srv.demoDriver.GetDeviceShadow(cid)
	if !ok {
		c.JSON(http.StatusBadRequest, GinResult("fail", "no cid"))
		return
	}

	//status:0 offline; status:1 online
	if devStatus.Status == 0 {
		virDev.CheckDevStatus(false)
	} else {
		virDev.CheckDevStatus(true)
	}

	srv.demoDriver.logger.Infof("postDeviceStatus cid:%s, Status:%v", cid, devStatus.Status)
	c.JSON(http.StatusOK, GinResult("ok", err))
}

func (srv *TestHttpSrv) postUploadImage(c *gin.Context) {
	cid := c.Param("cid")
	_, ok := srv.demoDriver.GetDeviceShadow(cid)
	if !ok {
		c.JSON(http.StatusBadRequest, GinResult("fail", "no cid"))
		return
	}

	err := TestUploadImage(srv.demoDriver, cid)
	srv.demoDriver.logger.Infof("postUploadImage cid:%s, err:%v", cid, err)
	c.JSON(http.StatusOK, GinResult("ok", err))
}

func (srv *TestHttpSrv) postModifyName(c *gin.Context) {
	cid := c.Query("cid")
	_, ok := srv.demoDriver.GetDeviceShadow(cid)
	if !ok {
		c.JSON(http.StatusBadRequest, GinResult("fail", "no cid"))
		return
	}

	name := c.Query("name")
	if len(cid) <= 0 {
		c.JSON(http.StatusBadRequest, GinResult("failed", "no cid"))
		return
	}

	err := TestModifyName(srv.demoDriver, cid, name)
	srv.demoDriver.logger.Infof("postModifyName ModifyName cid:%s, name:%s, err:%v", cid, name, err)
	c.JSON(http.StatusOK, GinResult("ok", err))
}

func (srv *TestHttpSrv) postModifyAttr(c *gin.Context) {
	cid := c.Param("cid")
	_, ok := srv.demoDriver.GetDeviceShadow(cid)
	if !ok {
		c.JSON(http.StatusBadRequest, GinResult("fail", "no cid"))
		return
	}

	err := TestModifyExt(srv.demoDriver, cid)
	srv.demoDriver.logger.Infof("postModifyAttr ModifyAttr cid:%s, err:%v", cid, err)
	c.JSON(http.StatusOK, GinResult("ok", err))
}

//func (ctx *TestHttpSrv) putCustomStorageV2(c *gin.Context) {
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
//	err = ctx.demoDriver.PutCustomStorageV2(req.Key, value)
//}
//
//func (ctx *TestHttpSrv) getCustomStorageV2(c *gin.Context) {
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
//	result, err = ctx.demoDriver.GetCustomStorageV2(keys)
//	if err != nil {
//		status = http.StatusInternalServerError
//	}
//
//	c.JSON(status, GinResult("ok", result))
//}
//
//func (ctx *TestHttpSrv) deleteCustomStorageV2(c *gin.Context) {
//	cids := []string{}
//	keysStr := c.Query("keys")
//	if keysStr != "" {
//		cids = strings.Split(keysStr, ",")
//	}
//	err := ctx.demoDriver.DeleteCustomStorageV2(cids)
//
//	c.JSON(http.StatusOK, GinResult("ok", err))
//}
//
//func (ctx *TestHttpSrv) queryStroageV2(c *gin.Context) {
//	prefix := c.Query("prefix")
//	fmt.Printf("prefix %s", prefix)
//	m, _ := ctx.demoDriver.QueryCustomStorageV2(prefix)
//	c.JSON(http.StatusOK, GinResult("ok", m))
//}
//
//func (ctx *TestHttpSrv) getCustomStorageKeysV2(c *gin.Context) {
//	var (
//		status = http.StatusOK
//	)
//	prefix := c.Query("prefix")
//	fmt.Printf("prefix %s", prefix)
//	strs, err := ctx.demoDriver.GetCustomStorageKeysV2(prefix)
//	if err != nil {
//		status = http.StatusInternalServerError
//	}
//	c.JSON(status, GinResult("ok", strs))
//}
