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
		if arg == "help" || arg == "-h" || arg == "-help" {
			printHelp()
		}
		if arg == "-v" || arg == "-V" || arg == "-version" || arg == "version" {
			printVersion()
		}
		if arg == "-l" || arg == "list" || arg == "-list" {
			printItemsList()
		}
		if arg == "delete" || arg == "-d" {
			deleteTarget(os.Args[2:])
		}
		if arg == "start" || arg =="-start" {
			return true
		}
		return false
	} else {
		return true
	}
}

func printVersion() {
	c := color.New(color.FgGreen)
	c.Print("version:")
	fmt.Println(version)
}

func printHelp() {
	fmt.Println("gorage", version)

	fmt.Println("usage:")
	fmt.Println("")
	fmt.Println("     gorage")
	fmt.Println("")

	fmt.Println("         -l / -list / list")
	fmt.Println("             list all uploaded items.")
	fmt.Println("")

	fmt.Println("         -v / -V / -version / version")
	fmt.Println("             Print softwear version.")
	fmt.Println("")

	fmt.Println("         -d / delete [target UUID(s)]")
	fmt.Println("             Delete item with UUIS(s).")
	fmt.Println("")

	fmt.Println("         -start / start")
	fmt.Println("             Launch softwear.")
	fmt.Println("")

	fmt.Println("         -h / -help / help")
	fmt.Println("             Print this message.")
	fmt.Println("")
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