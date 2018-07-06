mkdir gorage-linux-amd64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gorage-linux-amd64/gorage-linux-amd64 src/main.go
tar cfJ gorage-linux-amd64.tar.xz gorage-linux-amd64/gorage-linux-amd64
rm -rf gorage-linux-amd64/
