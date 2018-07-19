package main

import (
	"github.com/coreos/etcd/clientv3"
	"time"
	"github.com/cloudware-controller/configs"
	"context"
	"log"
)

func main() {

	ectdClientV3, err := clientv3.New(clientv3.Config{
		Endpoints:   configs.EtcdEndPoints,
		DialTimeout: 10 * time.Second,
		Username:    "root",
		Password:    "kfcoding",
	})

	if err != nil {
		panic(err)
	}

	//ectdClientV3.Delete(context.Background(), "/", clientv3.WithPrefix())

	resp, err := ectdClientV3.Get(context.Background(), "/", clientv3.WithPrefix())
	if err != nil {
		panic(err)
	}
	for _, kv := range resp.Kvs {
		log.Print(string(kv.Key), "=>", string(kv.Value))
	}
}
