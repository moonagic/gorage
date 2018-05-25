package config

import (
	"encoding/json"
	"io/ioutil"
	"runtime"
)

var (
	config map[string]string
)

// LoadConfig 读取配置文件
func LoadConfig() string {
	// 测试环境
	if runtime.GOOS == "darwin" {
		config = make(map[string]string)
		config["host"] = "127.0.0.1"
		config["port"] = "9090"
		return ""
	}

	result, err := ioutil.ReadFile("/etc/imagesStorage/config")
	if err == nil {
		var f interface{}
		json.Unmarshal(result, &f)
		m := f.(map[string]interface{})
		localHost, ok2 := m["host"].(string)
		localPort, ok3 := m["port"].(string)
		if ok2 && ok3 {
			config = make(map[string]string)
			config["host"] = localHost
			config["port"] = localPort
			return ""
		}
		return "Broken config."
	}
	return "Can not find config file...in \"/etc/gowebhook/config\""
}

// GetHost 获取监听Host
func GetHost() string {
	return config["host"]
}

// GetPort 获取监听port
func GetPort() string {
	return config["port"]
}
