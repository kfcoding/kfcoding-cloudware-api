package service

import "github.com/cloudware-controller/types"

type CloudwareService interface {
	CreateCloudwareApi(body *types.CloudwareBody) (string, error)

	WatcherCallback(body *types.KeeperBody)
}
