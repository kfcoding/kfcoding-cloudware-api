## 云件控制器

1. build

```
go build -o controller main.go

build using docker golang
docker run -it -v /Users/wsl/go/src:/go/src golang:1.10.3-alpine3.8 sh

cd src/github.com/kfcoding-cloudware-controller/ && go build -o controller main.go

scp controller root@worker:/home/cloudware-controller

cd /home/kfcoding-cloudware-controller && \
docker build -t daocloud.io/shaoling/kfcoding-cloudware-controller:v1.8 .

```