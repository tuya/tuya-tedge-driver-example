package dpmqtt

////////////////////////////////////////////////////////////////////////////////////////
type DriverConfig struct {
	HttpConf HttpConfig `yaml:"httpconf"`
	MqttConf MqttConfig `yaml:"mqttconf"`
}

type HttpConfig struct {
	Listen string `yaml:"listen"`
}

type MqttConfig struct {
	Username string `yaml:"username"`
}
