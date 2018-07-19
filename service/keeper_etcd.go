package service

import (
	"github.com/cloudware-controller/etcd"
	"github.com/cloudware-controller/types"
	"github.com/coreos/etcd/clientv3"
	"github.com/cloudware-controller/configs"
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
	log.Print("Keep: ", body)

	var resp *clientv3.LeaseGrantResponse
	var err error
	if resp, err = keeper.etcdClient.EctdClientV3.Grant(
		context.TODO(), configs.KeeperTTL); err != nil {
		log.Println("Keep error: ", err)
		return
	}

	key := path.Join(configs.KeeperPrefix, configs.Version, configs.TypeCloudware, body.Name)
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
	key := path.Join(configs.KeeperPrefix, configs.Version, configs.TypeCloudware, body.Name)
	return keeper.etcdClient.CheckExist(key)
}
