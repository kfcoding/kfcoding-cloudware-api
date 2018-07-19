package service

import "github.com/kfcoding-cloudware-controller/types"

type CloudwareService interface {
	CreateCloudwareApi(body *types.CloudwareBody) (string, error)

	WatcherCallback(body *types.KeeperBody)
}
