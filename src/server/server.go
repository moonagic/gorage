package server

import (
	"net/http"
	"log"
	"fmt"
	"os"
	"io"
	"imagesStorage/src/utils"
	"strings"
	"github.com/satori/go.uuid"
	"imagesStorage/src/config"
	"time"
)

func StartServer(address string, port string) error {
	http.HandleFunc("/", index)
	http.HandleFunc("/upload", uploadHandle)
	http.HandleFunc("/delete", deleteHandle)
	return http.ListenAndServe(address + ":" + port, nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Host)
	fmt.Fprintf(w, "{\"code\": 200, \"msg\": \"Service running...\"}")
}

func uploadHandle(w http.ResponseWriter, r *http.Request) {
	if !strings.EqualFold(r.Method, "post") {
		fmt.Fprintf(w, "{\"code\": 200, \"error\": \"Error Method.\"}")
		return
	}
	// 根据字段名获取表单文件
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

	fileDir := config.GetStorageDir() + fmt.Sprintf("%d", time.Now().Year()) + "/" + fmt.Sprintf("%d", int(time.Now().Month())) + "/" + fmt.Sprintf("%d", int(time.Now().Day())) + "/" // 文件保存dir
	log.Println(fileDir)
	log.Println(time.Now().Year())
	log.Println(int(time.Now().Month()))
	log.Println(time.Now().Day())
	if err := utils.CheckoutDir(fileDir); err != nil {
		fmt.Fprintf(w, "{\"code\": 200, \"error\": \"Server error.\"}")
		return
	}
	filePath := fileDir + header.Filename // 文件保存path
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
	fmt.Fprintf(w, "{\"code\": 200, \"msg\": \"Upload finished.\"}")

	// 记录上传数据
}

func deleteHandle(w http.ResponseWriter, r *http.Request)  {
	if !strings.EqualFold(r.Method, "delete") {
		fmt.Fprintf(w, "{\"code\": 200, \"error\": \"Error Method.\"}")
		return
	}
	// 删除处理
	log.Println(uuid.Must(uuid.NewV4()))
}