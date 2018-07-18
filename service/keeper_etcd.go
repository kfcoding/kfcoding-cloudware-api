package service

import (
	"github.com/cloudware-controller/etcd"
	"github.com/cloudware-controller/types"
)

type KeeperEtcdService struct {
	etcdClient *etcd.MyEtcdClient
}

func GetKeeperEtcdService() *KeeperEtcdService {
	return &KeeperEtcdService{
		etcdClient: etcd.GetMyEtcdClient(),
	}
}

func (keeper *KeeperEtcdService) Keep(body *types.KeeperBody) {

}

func (keeper *KeeperEtcdService) Remove(body *types.KeeperBody) {

}

func (keeper *KeeperEtcdService) Check(body *types.KeeperBody) {

}
