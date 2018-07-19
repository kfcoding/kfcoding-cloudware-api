package service

import "github.com/cloudware-controller/types"

type KeeperService interface {
	Keep(body *types.KeeperBody)

	Remove(body *types.KeeperBody)

	Check(body *types.KeeperBody) bool
}
