package etcd

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"strings"
	"strconv"
	"log"
	"sync"
	"github.com/kfcoding-cloudware-controller/configs"
)

type MyEtcdClient struct {
	EctdClientV3 *clientv3.Client
}

func (e *MyEtcdClient) Put(key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse) {
	ctx, cancel := context.WithTimeout(context.Background(), configs.RequestTimeout)
	resp, err := e.EctdClientV3.Put(ctx, key, val, opts...)
	cancel()
	if nil != err {
		log.Println("error put : ", err)
	}
	return resp
}

func (e *MyEtcdClient) Get(key string, opts ...clientv3.OpOption) (*clientv3.GetResponse) {
	ctx, cancel := context.WithTimeout(context.Background(), configs.RequestTimeout)
	resp, err := e.EctdClientV3.Get(ctx, key, opts ...)
	cancel()
	if nil != err {
		log.Println("error get : ", err)
	}
	return resp
}

func (e *MyEtcdClient) Delete(key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse) {
	ctx, cancel := context.WithTimeout(context.Background(), configs.RequestTimeout)
	resp, err := e.EctdClientV3.Delete(ctx, key, opts ...)
	cancel()
	if nil != err {
		log.Println("error delete : ", err)
	}
	return resp
}

func (e *MyEtcdClient) CheckExist(id string) (bool) {
	if e.Get(id).Count > 0 {
		return true
	}
	return false
}

func (e *MyEtcdClient) GetErrorType(err error) int {
	strs := strings.Split(err.Error(), ":")
	if len(strs) <= 0 {
		return -1
	}
	code, _ := strconv.Atoi(strings.TrimSpace(strs[0]))
	return code
}

var once sync.Once
var myEtcdClient *MyEtcdClient

func GetMyEtcdClient() *MyEtcdClient {
	once.Do(func() {

		var err error
		var config clientv3.Config
		if configs.EtcdUsername != "" {
			config = clientv3.Config{
				Endpoints:   configs.EtcdEndPoints,
				DialTimeout: configs.RequestTimeout,
				Username:    configs.EtcdUsername,
				Password:    configs.EtcdPassword,
			}
		} else {
			config = clientv3.Config{
				Endpoints:   configs.EtcdEndPoints,
				DialTimeout: configs.RequestTimeout,
			}
		}
		ectdClientV3, err := clientv3.New(config)
		if err != nil {
			log.Fatal("Error: new common client error:", err)
		}

		myEtcdClient = &MyEtcdClient{
			EctdClientV3: ectdClientV3,
		}
	})
	return myEtcdClient
}
