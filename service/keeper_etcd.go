package service

import (
	"github.com/kfcoding-cloudware-controller/etcd"
	"github.com/kfcoding-cloudware-controller/types"
	"github.com/coreos/etcd/clientv3"
	"github.com/kfcoding-cloudware-controller/configs"
	"log"
	"context"
	"path"
)

type KeeperEtcdService struct {
	etcdClient *etcd.MyEtcdClient
}

func GetKeeperEtcdService(etcdClient *etcd.MyEtcdClient) *KeeperEtcdService {
	return &KeeperEtcdService{
		etcdClient: etcdClient,
	}
}

func (keeper *KeeperEtcdService) Keep(body *types.KeeperBody) {
	log.Printf("Keep: %+v\n", body)

	var resp *clientv3.LeaseGrantResponse
	var err error
	if resp, err = keeper.etcdClient.EctdClientV3.Grant(
		context.TODO(), int64(configs.KeeperTTL)); err != nil {
		log.Println("Keep error: ", err)
		return
	}

	key := path.Join(configs.KeeperPrefix, configs.Version, body.Name)
	if _, err = keeper.etcdClient.EctdClientV3.Put(
		context.TODO(), key, "",
		clientv3.WithLease(resp.ID)); nil != err {
		log.Println("Keep error: ", err)
		return
	}
	log.Print("Keep ok")
}

func (keeper *KeeperEtcdService) Remove(body *types.KeeperBody) {

}

func (keeper *KeeperEtcdService) Check(body *types.KeeperBody) bool {
	key := path.Join(configs.KeeperPrefix, configs.Version, body.Name)
	return keeper.etcdClient.CheckExist(key)
}
