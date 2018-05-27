package server

import (
	"net/http"
	"log"
	"strings"
	"fmt"
	"os"
	"io"
)

func StartServer(address string, port string) error {
	http.HandleFunc("/", index)
	http.HandleFunc("/upload", uploadHandle)
	return http.ListenAndServe(address + ":" + port, nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Host)
	fmt.Fprintf(w, "{\"code\": 200, \"msg\": \"Service running...\"}")
}

func uploadHandle(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" && r.Method != "post" {
		fmt.Fprintf(w, "{\"code\": 200, \"error\": \"Error Method.\"}")
		return
	}
	// 根据字段名获取表单文件
	formFile, header, err := r.FormFile("file")

	validFileType := false
	if strings.HasSuffix(header.Filename, ".png") {
		validFileType = true
	}
	if strings.HasSuffix(header.Filename, ".jpg") {
		validFileType = true
	}
	if strings.HasSuffix(header.Filename, ".jpeg") {
		validFileType = true
	}
	if strings.HasSuffix(header.Filename, ".bmp") {
		validFileType = true
	}
	if strings.HasSuffix(header.Filename, ".apng") {
		validFileType = true
	}
	if strings.HasSuffix(header.Filename, ".webp") {
		validFileType = true
	}
	if !validFileType {
		fmt.Fprintf(w, "{\"code\": 200, \"error\": \"Invalid file type.\"}")
		return
	}

	if err != nil {
		log.Printf("Get form file failed: %s\n", err)
		fmt.Fprintf(w, "{\"code\": 200, \"error\": \"Get form file failed.\"}")
		return
	}
	defer formFile.Close()

	// 创建保存文件
	gotFile, err := os.Create("." + r.URL.Path + "/" + header.Filename)
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
}
