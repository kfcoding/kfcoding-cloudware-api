package service

import (
	"github.com/kfcoding-cloudware-controller/etcd"
	"log"
	"github.com/coreos/etcd/clientv3"
	"strings"
	"context"
	"github.com/kfcoding-cloudware-controller/types"
)

type EtcdWatcher struct {
	etcdClient *etcd.MyEtcdClient
}

func GetEtcdWatcher(etcdClient *etcd.MyEtcdClient) *EtcdWatcher {
	return &EtcdWatcher{
		etcdClient: etcdClient,
	}
}

func (watcher *EtcdWatcher) Watcher(prefix string, service CloudwareService) {
	log.Println("Start Etcd Watcher")

	rch := watcher.etcdClient.EctdClientV3.Watch(context.Background(), prefix, clientv3.WithPrefix())

	for wresp := range rch {
		for _, ev := range wresp.Events {
			//fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			// /kfcoding/v1/1/cloduware-321dsnaknfkdsnjf9afndks
			switch ev.Type {
			case 1: //DELETE
				log.Print("listen etcd delete: ", string(ev.Kv.Key))
				strs := strings.Split(string(ev.Kv.Key), "/")
				body := &types.KeeperBody{Name: strs[len(strs)-1]}
				service.WatcherCallback(body)
			}
		}
	}

}
