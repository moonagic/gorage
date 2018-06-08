package config

import (
	"encoding/json"
	"gorage/src/data"
	"io/ioutil"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/syndtr/goleveldb/leveldb"
)

var (
	config        map[string]string
	KeyCacheArray keyCache
)

type keyCache []data.KeyMap

func (k keyCache) Len() int { return len(k) }
func (k keyCache) Less(i, j int) bool {
	int640, _ := strconv.ParseInt(k[i].TagTime, 10, 64)
	int641, _ := strconv.ParseInt(k[j].TagTime, 10, 64)
	return int640 < int641
}
func (k keyCache) Swap(i, j int) { k[i], k[j] = k[j], k[i] }

// LoadConfig read config file
func LoadConfig() string {
	var configFile string
	switch runtime.GOOS {
	case "darwin":
		configFile = "./config-macintosh"
		break
	case "linux":
		configFile = "/etc/gorage/config"
		break
	case "windows":
		configFile = "../config-windows"
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
	return "Can not find config file...in \"" + configFile + "\""
}

// LoadKeyCache
func LoadKeyCache() {
	if db, err := leveldb.OpenFile(GetDataBase(), nil); err == nil {
		item := db.NewIterator(nil, nil)
		for item.Next() {
			value := item.Value()
			var f interface{}
			json.Unmarshal(value, &f)
			bodyMap := f.(map[string]interface{})
			keyModel := data.KeyMap{
				UUID:    bodyMap["UUID"].(string),
				TagTime: bodyMap["TagTime"].(string),
			}
			KeyCacheArray = append(KeyCacheArray, keyModel)
		}
		item.Release()
		err = item.Error()
		if err != nil {
			color.Red("Parse error")
		}
		defer db.Close()
	}
	sort.Sort(KeyCacheArray)

	for i := 0; i < len(KeyCacheArray); i++ {
		KeyCacheArray[i].Index = i
	}
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
