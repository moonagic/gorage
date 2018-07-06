# gorage
A simple image hosting tool written by go.  
[![Go Report Card](https://goreportcard.com/badge/github.com/moonagic/gorage)](https://goreportcard.com/report/github.com/moonagic/gorage)

## How To Get Started
* [Download the release](https://github.com/moonagic/gorage/releases) and move to somewhere, for example `/usr/local/bin/gorage`
* Create a config file in `/etc/gorage/config`
for example:
```
{
  "url": "https://example.com/",
  "host": "127.0.0.1",
  "port": "9909",
  "fileType": "png, jpg, jpeg, webp, bmp, apng, ttf, zip, sh",
  "storageDir": "/var/www/content/",
  "db": "/etc/gorage/.database"
}
```
Then just run `gorage -start`

## Command line
* List all uploaded project
<img src="https://github.com/moonagic/gorage/blob/master/images/list.png">

* Delete uploaded project
<img src="https://github.com/moonagic/gorage/blob/master/images/delete.png">

## Usage
* Upload file  
```bash
curl -F "file=@example.png" http://127.0.0.1:9909/upload
{"code": 200, "msg": "Upload finished.", "data":{"UUID":"09c742d7-a8b4-4923-ace7-199aa0e2d169","FileName":"example.png","Directory":"2018/6/13/wfxyoyyxqu43bap7/","TagTime":"1528882517983","UploadTime":"2018-06-13 17:35:17"}, "url":"http://example.com/content/2018/6/13/wfxyoyyxqu43bap7/example.png"}
```

* List all uploaded project  
```bash
curl -X get http://127.0.0.1:9909/list?page=1
{"code": 200, "data": [{"Index":0,"UUID":"572ccde8-a42c-4c28-9260-7e030d4fb8e5","TagTime":"1528186905334"},{"Index":1,"UUID":"c8368d8b-8ca2-4ebe-9667-e018625aa8e1","TagTime":"1528267017421"},{"Index":2,"UUID":"38bbe867-02a8-4e2a-8203-74830afd9454","TagTime":"1528267854439"}]}
```

* Delete uploaded project  
```bash
curl -X delete http://127.0.0.1:9909/delete -d '{"key":"572ccde8-a42c-4c28-9260-7e030d4fb8e5"}'
{"code": 200, "msg": "Delete finished."}
```
