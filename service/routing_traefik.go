package service

import "github.com/cloudware-controller/etcd"

type RoutingTraefikService struct {
	etcdClient *etcd.MyEtcdClient
}

func GetRoutingTraefikService(etcdClient *etcd.MyEtcdClient) *RoutingTraefikService {
	return &RoutingTraefikService{
		etcdClient: etcdClient,
	}
}

func (service *RoutingTraefikService) AddRule() {

}

func (service *RoutingTraefikService) AddRules() {

}

func (service *RoutingTraefikService) DeleteRule() {

}

func (service *RoutingTraefikService) DeleteRules() {

}
