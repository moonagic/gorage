SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o gorage-linux-amd64 src/main.go
tar cfJ gorage-linux-amd64.tar.xz gorage-linux-amd64

