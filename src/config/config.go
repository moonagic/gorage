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
	var configFile string
	switch runtime.GOOS {
	case "darwin":
		configFile = ".config"
		break
	case "linux":
		configFile = "/etc/imagesStorage/config"
		break
	case "windows":
		configFile = "D:\\goWorkspace\\src\\imagesStorage\\config"
		break
	}

	var result []byte
	result, err := ioutil.ReadFile(configFile)
	if err == nil {
		var f interface{}
		json.Unmarshal(result, &f)
		m := f.(map[string]interface{})
		localHost, ok2 := m["host"].(string)
		localPort, ok3 := m["port"].(string)
		localType, ok4 := m["fileType"].(string)
		if ok2 && ok3 && ok4 {
			config = make(map[string]string)
			config["host"] = localHost
			config["port"] = localPort
			config["fileType"] = localType
			return ""
		}
		return "Broken config."
	}
	return "Can not find config file...in \"/etc/imagesStorage/config\""
}

// GetHost 获取监听Host
func GetHost() string {
	return config["host"]
}

// GetPort 获取监听port
func GetPort() string {
	return config["port"]
}

// GetTypeps 获取可用文件后缀
func GetTypeps() string {
	return config["fileType"];
}