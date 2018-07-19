package service

import "github.com/kfcoding-cloudware-controller/types"

type RoutingService interface {
	AddRule(*types.RoutingBody) error

	AddRules([]interface{}) error

	DeleteRule(*types.RoutingBody) error

	DeleteRules([]interface{}) error
}
