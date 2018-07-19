package service

import "github.com/kfcoding-cloudware-controller/types"

type KeeperService interface {
	Keep(body *types.KeeperBody)

	Remove(body *types.KeeperBody)

	Check(body *types.KeeperBody) bool
}
