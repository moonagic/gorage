# gorage
An image hosting tool written by go.  
[![Go Report Card](https://goreportcard.com/badge/github.com/moonagic/gorage)](https://goreportcard.com/report/github.com/moonagic/gorage)

## How To Get Started
* [Download the release](https://github.com/moonagic/gorage/releases) and mv to somewhere, for example `/usr/local/bin/gorage`
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
Then just run `gorage start`

## Command line
* List all uploaded project
<img src="https://github.com/moonagic/gorage/blob/master/images/list.png" width="600">

* Delete uploaded project
<img src="https://github.com/moonagic/gorage/blob/master/images/delete.png" width="600">

