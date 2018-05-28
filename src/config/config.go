package config

import (
	"encoding/json"
	"io/ioutil"
	"runtime"
	"strings"
)

var (
	config map[string]string
)

// LoadConfig read config file
func LoadConfig() string {
	var configFile string
	switch runtime.GOOS {
	case "darwin":
		configFile = ".config"
		break
	case "linux":
		configFile = "/etc/gorage/config"
		break
	case "windows":
		configFile = "config-windows"
		break
	}

	var result []byte
	result, err := ioutil.ReadFile(configFile)
	if err == nil {
		var f interface{}
		json.Unmarshal(result, &f)
		m := f.(map[string]interface{})
		localURL, ok1 := m["url"].(string)
		if !strings.HasSuffix(localURL, "/") {
			localURL += "/"
		}
		localHost, ok2 := m["host"].(string)
		localPort, ok3 := m["port"].(string)
		localType, ok4 := m["fileType"].(string)
		localStorage, ok5 := m["storageDir"].(string)
		if !strings.HasSuffix(localStorage, "/") {
			localStorage += "/"
		}
		localDB, ok6 := m["db"].(string)
		if ok1 && ok2 && ok3 && ok4 && ok5 && ok6 {
			config = make(map[string]string)
			config["url"] = localURL
			config["host"] = localHost
			config["port"] = localPort
			config["fileType"] = localType
			config["storageDir"] = localStorage
			config["localDatabase"] = localDB
			return ""
		}
		return "Broken config."
	}
	return "Can not find config file...in \"/etc/gorage/config\""
}

// GetURL
func GetURL() string {
	return config["url"]
}
// GetHost
func GetHost() string {
	return config["host"]
}

// GetPort
func GetPort() string {
	return config["port"]
}

// GetTypes
func GetTypes() string {
	return config["fileType"]
}

// GetStorageDir
func GetStorageDir() string {
	return config["storageDir"]
}

// GetDataBase
func GetDataBase() string {
	return config["localDatabase"]
}