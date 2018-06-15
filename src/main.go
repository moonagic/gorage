package main

import (
	"encoding/json"
	"fmt"
	"gorage/src/config"
	"gorage/src/server"
	"gorage/src/utils"
	"log"
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	version = "0.1"
)

func main() {

	if err := config.LoadConfig(); err != "" {
		log.Fatal(err)
		os.Exit(0)
	}
	config.LoadKeyCache()
	if parseArgs() {
		log.Fatal(server.StartServer(config.GetHost(), config.GetPort()))
	}

}

func parseArgs() bool {
	if len(os.Args) > 1 {
		arg := os.Args[1]
		switch arg {
		case "-h", "-help":
			printHelp()
		case "-v", "-version":
			printVersion()
		case "-l", "-list":
			if len(os.Args) > 2 {
				if page, err := strconv.Atoi(os.Args[2]); err == nil {
					printItemsListWithPage(page)
				}
				return false
			}
			printItemsList()
		case "-d":
			deleteTarget(os.Args[2:])
		case "-start":
			return true
		}
		return false
	}
	return true
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

	fmt.Println("         -l / -list  [page(option)]")
	fmt.Println("             list all uploaded items.")
	fmt.Println("")

	fmt.Println("         -v / -version")
	fmt.Println("             Print softwear version.")
	fmt.Println("")

	fmt.Println("         -d [target UUID(s)]")
	fmt.Println("             Delete item with UUIS(s).")
	fmt.Println("")

	fmt.Println("         -start")
	fmt.Println("             Launch softwear.")
	fmt.Println("")

	fmt.Println("         -h / -help")
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
					color.Red("Error in remove file: %s", config.GetStorageDir()+fileDir+fileName)
					log.Println("file:", fileDir+fileName)
					log.Println("key:", key)
				} else {
					log.Println("Remove file finished.")
					log.Println("file:", fileDir+fileName)
					log.Println("key:", key)
				}
				if err := os.Remove(config.GetStorageDir() + fileDir); err != nil {
					log.Println("Remove directory")
				} else {
					color.Red("Error in remove directory: %s", config.GetStorageDir()+fileDir)
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

func printItemsListWithPage(page int) {
	start := (page - 1) * 10
	end := page * 10
	keys := utils.GetListWithStartAndEnd(start, end)
	if db, err := leveldb.OpenFile(config.GetDataBase(), nil); err != nil {
		color.Red("Error in Open Database.")
	} else {
		color.Blue("--------------------")
		for i := 0; i < len(keys); i++ {
			if value, err := db.Get([]byte(keys[i].UUID), nil); err == nil {
				color.Green("Index:%d", keys[i].Index)
				showItemValue(value)
				color.Blue("--------------------")
			} else {
				// get error
				color.Red("Error in get value with key.")
			}
		}
		defer db.Close()
	}
}

func printItemsList() {
	printItemsListWithPage(1)
	if len(config.KeyCacheArray) > 10 {
		color.Green("Too many data, only the top 10 datas displayed. You can use pagination queries.")
	}
}

func showItemValue(value []byte) {
	var f interface{}
	json.Unmarshal(value, &f)
	bodyMap := f.(map[string]interface{})
	color.Green("UUID:%s", bodyMap["UUID"].(string))
	color.White("FileName:%s", bodyMap["FileName"].(string))
	color.White("Directory:%s", bodyMap["Directory"].(string))
	color.White("URL:%s", config.GetURL()+"content/"+bodyMap["Directory"].(string)+bodyMap["FileName"].(string))
	color.White("UploadTime:%s", bodyMap["UploadTime"].(string))
	color.White("TagTime:%s", bodyMap["TagTime"].(string))
}
