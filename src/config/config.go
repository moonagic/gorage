package config

import (
	"encoding/json"
	"io/ioutil"
)

var (
	config map[string]string
)

// LoadConfig 读取配置文件
func LoadConfig() string {
	result, err := ioutil.ReadFile("/etc/imagesStorage/config")
	if err == nil {
		var f interface{}
		json.Unmarshal(result, &f)
		m := f.(map[string]interface{})
		localSecret, ok1 := m["secret"].(string)
		localHost, ok2 := m["host"].(string)
		localPort, ok3 := m["port"].(string)
		if ok1 && ok2 && ok3 {
			config = make(map[string]string)
			config["secret"] = localSecret
			config["host"] = localHost
			config["port"] = localPort
			return ""
		}
		return "Broken config."
	}
	return "Can not find config file...in \"/etc/gowebhook/config\""
}

// GetSecret 获取secret
func GetSecret() string {
	return config["secret"]
}

// GetHost 获取监听Host
func GetHost() string {
	return config["host"]
}

// GetPort 获取监听port
func GetPort() string {
	return config["port"]
}
