package main

import (
	"os"
	"log"
	"gorage/src/server"
	"gorage/src/config"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"encoding/json"
	"github.com/fatih/color"
)

const (
	version = "0.1"
	buildVersion = "102"
)

func main() {
	if err := config.LoadConfig(); err != "" {
		log.Fatal(err)
		os.Exit(0)
	}
	if parseArgs() {
		log.Fatal(server.StartServer(config.GetHost(), config.GetPort()))
	}

}

func parseArgs() bool {
	if len(os.Args) > 1 {
		arg := os.Args[1]
		switch arg {
		case "help":
			printHelp()
		case "-h":
			printHelp()
		case "-help":
			printHelp()
		case "-v":
			printVersion()
		case "-V":
			printVersion()
		case "-version":
			printVersion()
		case "--version":
			printVersion()
		case "version":
			printVersion()
		case "list":
			printItemsList()
		case "delete":
			deleteTarget(os.Args[2:])
		case "start":
			return true
		}
		return false
	} else {
		return true
	}
}

func printVersion() {
	fmt.Print("version:")
	c := color.New(color.FgGreen)
	c.Println(version)

	fmt.Print("build:")
	c.Println(buildVersion)
}

func printHelp() {
	fmt.Println("Help info")
}

func deleteTarget(target []string) {
	if db, err := leveldb.OpenFile(config.GetDataBase(), nil); err != nil {
		color.Red("Error in Open Database.")
	} else {
		for _, key := range target {
			value, err := db.Get([]byte(key), nil)
			if err == nil {
				log.Println(string(value))
				// delete file
				var f interface{}
				json.Unmarshal(value, &f)
				valueMap := f.(map[string]interface{})
				fileDir := valueMap["Directory"].(string)
				fileName := valueMap["FileName"].(string)
				if err := os.Remove(config.GetStorageDir() + fileDir + fileName); err != nil {
					color.Red("Error in remove file: %s", config.GetStorageDir() + fileDir + fileName)
					log.Println("file:", fileDir + fileName)
					log.Println("key:", key)
				} else {
					log.Println("Remove file finished.")
					log.Println("file:", fileDir + fileName)
					log.Println("key:", key)
				}

				// delete data
				log.Println("delte value by key:", key)
				err = db.Delete([]byte(key), nil)
				if err != nil {
					color.Red("Error in delete value by key!")
				}
			} else {
				color.Red("Not found value by key:%s", key)
			}
		}
		defer db.Close()
	}
}

func printItemsList() {
	if db, err := leveldb.OpenFile(config.GetDataBase(), nil); err != nil {
		color.Red("Error in Open Database.")
	} else {
		item := db.NewIterator(nil, nil)
		color.Blue("--------------------")
		for item.Next() {
			value := item.Value()
			var f interface{}
			json.Unmarshal(value, &f)
			bodyMap := f.(map[string]interface{})
			color.Green("UUID:%s", bodyMap["UUID"].(string))
			color.White("FileName:%s", bodyMap["FileName"].(string))
			color.White("Directory:%s", bodyMap["Directory"].(string))
			color.White("URL:%s", config.GetURL() + "content/" + bodyMap["Directory"].(string) + bodyMap["FileName"].(string))

			color.Blue("--------------------")
		}
		item.Release()
		err = item.Error()
		if err != nil {
			color.Red("Parse error")
		}
		defer db.Close()
	}
}