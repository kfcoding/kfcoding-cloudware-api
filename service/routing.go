package service

import "github.com/cloudware-controller/types"

type RoutingService interface {
	AddRule(*types.RoutingBody) error

	AddRules([]interface{}) error

	DeleteRule(*types.RoutingBody) error

	DeleteRules([]interface{}) error
}
