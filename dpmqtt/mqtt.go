package dpmqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/tuya/tuya-tedge-driver-sdk-go/commons"
)

var _ commons.MqttDriver = &MqttItf{}

type MqttItf struct {
	driver *DemoMqttDPDriver
	log    commons.TedgeLogger
}

func NewMqttItf(dpDriver *DemoMqttDPDriver) *MqttItf {
	return &MqttItf{
		driver: dpDriver,
		log:    dpDriver.dpmDS.GetLogger(),
	}
}

//
func (m *MqttItf) Auth(clientId, username, password string) (bool, error) {
	m.log.Infof("mqtt auth, clientId:%s%s username:%s password:%s", clientId, username, password)

	//TODO: implement me
	//......

	return true, nil
}

func (m *MqttItf) Sub(clientId, username, topic string, qos byte) (bool, error) {
	m.log.Infof("mqtt subscribe, clientId:%s%s username:%s topic:%s qos:%v", clientId, username, topic, qos)

	//TODO: implement me
	//......

	return true, nil
}

func (m *MqttItf) Pub(clientId, username, topic string, qos byte, retained bool) (bool, error) {
	m.log.Infof("mqtt publish: clientId:%s%s username:%s topic:%s qos:%v, retained:%v", clientId, username, topic, qos, retained)

	//TODO: implement me
	//......

	return true, nil
}

func (m *MqttItf) UnSub(clientId, username string, topics []string) {
	m.log.Infof("mqtt unsubscribe clientId:%s%s username:%s topic:%s", clientId, username, topics)

	//TODO: implement me
	//......

	return
}

func (m *MqttItf) Connected(clientId, username, ip, port string) {
	m.log.Infof("mqtt connect clientId:%s%s username:%s ip:%s port:%s", clientId, username, ip, port)

	//TODO: implement me
	//......
}

func (m *MqttItf) Closed(clientId, username string) {
	m.log.Infof("mqtt Closed clientId:%s username:%s", clientId, username)

	//TODO: implement me
	//......
}

func (m *MqttItf) ConnectHandler(client mqtt.Client) {
	//TODO: implement me
	//......
}
