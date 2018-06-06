package server

import (
	"net/http"
	"log"
	"fmt"
	"os"
	"io"
	"gorage/src/utils"
	"strings"
	"github.com/satori/go.uuid"
	"gorage/src/config"
	"time"
	"gorage/src/data"
	"encoding/json"
	"github.com/syndtr/goleveldb/leveldb"
	"io/ioutil"
	"strconv"
	"github.com/fatih/color"
)

func StartServer(address string, port string) error {
	http.HandleFunc("/", index)
	http.HandleFunc("/upload", uploadHandle)
	http.HandleFunc("/delete", deleteHandle)
	http.HandleFunc("/list", listHandle)
	http.HandleFunc("/item", itemHandle)

	fs := http.FileServer(http.Dir(config.GetStorageDir()))
	http.Handle("/images/", http.StripPrefix("/images", fs))
	return http.ListenAndServe(address + ":" + port, nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/" || r.RequestURI == "/index.html" || r.RequestURI == "/index.htm" {
		fmt.Fprintf(w, "{\"code\": 200, \"msg\": \"Service running...\"}")
	} else {
		fmt.Fprintf(w, "{\"code\": 404, \"error\": \"404 Not Found.\"}")
	}
}

func itemHandle(w http.ResponseWriter, r *http.Request) {
	if !strings.EqualFold(r.Method, "get") {
		fmt.Fprintf(w, "{\"code\": 200, \"error\": \"Error Method.\"}")
		return
	}
	r.ParseForm()
	UUID := r.Form["UUID"][0]

	if db, err := leveldb.OpenFile(config.GetDataBase(), nil); err == nil {
		if value, err := db.Get([]byte(UUID), nil); err == nil {
			fmt.Fprintf(w, "{\"code\": 200, \"data\": %s}", value)
		}
	} else {
		fmt.Fprintf(w, "{\"code\": 200, \"error\": \"System exception.\"}")
	}
}

func listHandle(w http.ResponseWriter, r *http.Request) {
	if !strings.EqualFold(r.Method, "get") {
		fmt.Fprintf(w, "{\"code\": 200, \"error\": \"Error Method.\"}")
		return
	}
	r.ParseForm()
	if page, err := strconv.Atoi(string(r.Form["page"][0])); err == nil {
		start := (page - 1) * 10
		end := page * 10
		keys := utils.GetListWithStartAndEnd(start, end)

		str, _ := json.Marshal(keys)
		fmt.Fprintf(w, "{\"code\": 200, \"data\": %s}", string(str))
	} else {
		fmt.Fprintf(w, "{\"code\": 200, \"error\": \"Error parm.\"}")
	}
}

func uploadHandle(w http.ResponseWriter, r *http.Request) {
	if !strings.EqualFold(r.Method, "post") {
		fmt.Fprintf(w, "{\"code\": 200, \"error\": \"Error Method.\"}")
		return
	}
	// get file in form
	formFile, header, err := r.FormFile("file")

	if !utils.VerifyFileType(header.Filename) {
		fmt.Fprintf(w, "{\"code\": 200, \"error\": \"Invalid file type.\"}")
		return
	}

	if err != nil {
		log.Printf("Get form file failed: %s\n", err)
		fmt.Fprintf(w, "{\"code\": 200, \"error\": \"Get form file failed.\"}")
		return
	}
	defer formFile.Close()

	// storage directory
	fileDir := fmt.Sprintf("%d", time.Now().Year()) +
		"/" +
		fmt.Sprintf("%d", int(time.Now().Month())) +
		"/" +
		fmt.Sprintf("%d", int(time.Now().Day())) +
		"/" +
		utils.GetRandomString(16) +
		"/"

	if err := utils.CheckoutDir(config.GetStorageDir() + fileDir); err != nil {
		fmt.Fprintf(w, "{\"code\": 200, \"error\": \"Server error.\"}")
		return
	}
	filePath := config.GetStorageDir() + fileDir + header.Filename // 文件保存path

	gotFile, err := os.Create(filePath)
	if err != nil {
		log.Printf("Create failed: %s\n", err)
		fmt.Fprintf(w, "{\"code\": 200, \"error\": \"Create file failed.\"}")
		return
	}
	defer gotFile.Close()

	// 读取表单文件，写入保存文件
	_, err = io.Copy(gotFile, formFile)
	if err != nil {
		log.Printf("Write file failed: %s\n", err)
		fmt.Fprintf(w, "{\"code\": 200, \"error\": \"Write file failed.\"}")
		return
	}

	thisTime := time.Now()
	itemUUID := uuid.Must(uuid.NewV4()).String()
	// 记录上传数据
	item := data.UploadItem{
		UUID: itemUUID,
		FileName: header.Filename,
		Directory: fileDir,
		TagTime: strconv.FormatInt(thisTime.UnixNano()/1e6,10),
		UploadTime: thisTime.Format("2006-01-02 15:04:05"),
	}
	mData, err := json.Marshal(item)
	if err != nil{
		log.Fatalln(err)
	}

	if db, err := leveldb.OpenFile(config.GetDataBase(), nil); err != nil {
		log.Println("Open Database faild.")
	} else {
		err = db.Put([]byte(itemUUID), mData, nil)
		defer db.Close()
	}

	url := config.GetURL() + "content/" + fileDir + header.Filename
	url = strings.Replace(url, "//","/", -1)
	fmt.Fprintf(w, "{\"code\": 200, \"msg\": \"Upload finished.\", \"data\":%s, \"url\":\"%s\"}", string(mData), url)
}

func deleteHandle(w http.ResponseWriter, r *http.Request)  {
	if !strings.EqualFold(r.Method, "delete") {
		fmt.Fprintf(w, "{\"code\": 200, \"error\": \"Error Method.\"}")
		return
	}
	// Delete
	if db, err := leveldb.OpenFile(config.GetDataBase(), nil); err != nil {
		log.Println("Open Database faild.")
		log.Println(err)
	} else {
		if bodyContent, err := ioutil.ReadAll(r.Body); err == nil {
			var f interface{}
			json.Unmarshal(bodyContent, &f)
			bodyMap := f.(map[string]interface{})
			if key, ok := bodyMap["key"].(string); ok {
				// get value of the key
				value, err := db.Get([]byte(key), nil)
				if err == nil {
					// delete file
					var f interface{}
					json.Unmarshal(value, &f)
					valueMap := f.(map[string]interface{})
					fileDir := valueMap["Directory"].(string)
					fileName := valueMap["FileName"].(string)
					if !utils.CheckoutIfFileExists(fileDir + fileName) {
						log.Println("File not found, file:", fileDir + fileName)
					}
					if err := os.Remove(fileDir + fileName); err != nil {
						log.Println("Remove file faild, file:", fileDir + fileName)
					}
					if err := os.Remove(fileDir); err != nil {
						log.Println("Remove directory")
					} else {
						color.Red("Error in remove directory: %s", fileDir)
					}
				} else {
					log.Println("Not found value by key")
					fmt.Fprintf(w, "{\"code\": 200, \"error\": \"Not found value by key.\"}")
					defer db.Close()
					return
				}

				// delete data
				log.Println("delte key:", key)
				err = db.Delete([]byte(key), nil)
				if err != nil {
					log.Println(err)
					fmt.Fprintf(w, "{\"code\": 200, \"error\": \"Delete value faild in database.\"}")
					defer db.Close()
					return
				}
				fmt.Fprintf(w, "{\"code\": 200, \"error\": \"Delete finished.\"}")
			}
		}
		defer db.Close()
	}
}
