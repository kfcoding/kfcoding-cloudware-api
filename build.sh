#!/bin/bash
go build -o CloudwareController main.go
docker build -t kfcoding.com/cloudware-controller:v1.0 .
docker save kfcoding.com/cloudware-controller:v1.0> cloudware-controller.tar.gz
scp cloudware-controller.tar.gz root@node2:/root
scp cloudware-controller.tar.gz root@node3:/root