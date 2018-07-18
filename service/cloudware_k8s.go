package service

import "github.com/cloudware-controller/types"

type CloudwareK8sService struct {
	Keeper  KeeperService
	Routing RoutingService
}

func GetCloudwareK8sService(keeper KeeperService, routing RoutingService) *CloudwareK8sService {
	return &CloudwareK8sService{
		Keeper:  keeper,
		Routing: routing,
	}
}

func (service *CloudwareK8sService) CreateCloudwareApi(body *types.CloudwareBody) (string, error) {
	return "ok", nil
}

func (service *CloudwareK8sService) WatcherCallback(body *types.KeeperBody) {

}
