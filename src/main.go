package main

import (
	"os"
	"log"
	"gorage/src/server"
	"gorage/src/config"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"encoding/json"
)

const (
	version = 0.1
	buildVersion = 100
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
			deleteTarget(os.Args[:2])
		case "start":
			return true
		}
		return false
	} else {
		return true
	}
}

func printVersion() {
	fmt.Println("gorage version" , version)
	fmt.Println("build:" , buildVersion)
}

func printHelp() {
	fmt.Println("Help info")
}

func deleteTarget(target []string) {
	fmt.Println("delete target", target)
	if db, err := leveldb.OpenFile(config.GetDataBase(), nil); err != nil {
		log.Println("Open Database faild.")
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
				if err := os.Remove(fileDir + fileName); err != nil {
					log.Println("Remove file faild")
					log.Println("file:", fileDir + fileName)
					log.Println("key:", key)
				} else {
					log.Println("Remove file finished.")
					log.Println("file:", fileDir + fileName)
					log.Println("key:", key)
				}
			} else {
				log.Println("Not found value by key", key)
			}

			// delete data
			log.Println("delte key:", key)
			err = db.Delete([]byte(key), nil)
		}
		defer db.Close()
	}
}

func printItemsList() {
	if db, err := leveldb.OpenFile(config.GetDataBase(), nil); err != nil {
		log.Println("Open Database faild.")
	} else {
		item := db.NewIterator(nil, nil)
		for item.Next() {
			// Remember that the contents of the returned slice should not be modified, and
			// only valid until the next call to Next.
			value := item.Value()
			fmt.Println(string(value))
		}
		item.Release()
		err = item.Error()
		defer db.Close()
	}
}