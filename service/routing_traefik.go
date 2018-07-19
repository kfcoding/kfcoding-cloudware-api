package service

import (
	"github.com/kfcoding-cloudware-controller/etcd"
	"github.com/kfcoding-cloudware-controller/types"
	"log"
	"github.com/kfcoding-cloudware-controller/configs"
	"path"
	"context"
)

type RoutingTraefikService struct {
	etcdClient *etcd.MyEtcdClient
}

func GetRoutingTraefikService(etcdClient *etcd.MyEtcdClient) *RoutingTraefikService {
	return &RoutingTraefikService{
		etcdClient: etcdClient,
	}
}

func (service *RoutingTraefikService) AddRule(rule *types.RoutingBody) error {
	return service.addToEtcd(rule)
}

func (service *RoutingTraefikService) AddRules(rules []interface{}) error {
	var err error
	for _, v := range rules {
		body, _ := v.(map[string]interface{})
		err = service.addToEtcd(&types.RoutingBody{
			Name: body["Name"].(string),
			URL:  body["URL"].(string),
			Rule: body["Rule"].(string),
		})
		if err != nil {
			log.Print("AddRules error: ", err)
		}
	}
	return err
}

func (service *RoutingTraefikService) DeleteRule(rule *types.RoutingBody) error {
	return service.deleteFromEtcd(rule)
}

func (service *RoutingTraefikService) DeleteRules(rules []interface{}) error {
	var err error
	for _, v := range rules {
		body, _ := v.(map[string]interface{})
		err = service.deleteFromEtcd(&types.RoutingBody{
			Name: body["Name"].(string),
			URL:  body["URL"].(string),
			Rule: body["Rule"].(string),
		})
		if err != nil {
			log.Print("AddRules error: ", err)
		}
	}
	return err
}

func (service *RoutingTraefikService) addToEtcd(rule *types.RoutingBody) error {
	// set backend
	key := path.Join(configs.PrefixTraefik, "backends/", rule.Name, "/servers/1/url")
	value := rule.URL
	if _, err := service.etcdClient.EctdClientV3.Put(context.Background(), key, value); nil != err {
		log.Println("add rule to etcd put backend error: ", err)
		return err
	}

	//set frontend
	key = path.Join(configs.PrefixTraefik, "frontends/", rule.Name, "/backend")
	value = rule.Name
	if _, err := service.etcdClient.EctdClientV3.Put(context.Background(), key, value); nil != err {
		log.Println("add rule to etcd put frontend error: ", err)
		return err
	}
	key = path.Join(configs.PrefixTraefik, "frontends/", rule.Name, "/routes/1/rule")
	value = rule.Rule
	if _, err := service.etcdClient.EctdClientV3.Put(context.Background(), key, value); nil != err {
		log.Println("add rule to etcd put frontend error: ", err)
		return err
	}

	return nil
}

func (service *RoutingTraefikService) deleteFromEtcd(rule *types.RoutingBody) error {
	// delete backend
	key := path.Join(configs.PrefixTraefik, "backends/", rule.Name, "/servers/1/url")
	if _, err := service.etcdClient.EctdClientV3.Delete(context.Background(), key); nil != err {
		log.Println("delete rule to etcd delete backend error: ", err)
		return err
	}

	// delete frontend
	key = path.Join(configs.PrefixTraefik, "frontends/", rule.Name, "/backend")
	if _, err := service.etcdClient.EctdClientV3.Delete(context.Background(), key); nil != err {
		log.Println("delete rule to etcd delete frontend error: ", err)
		return err
	}
	key = path.Join(configs.PrefixTraefik, "frontends/", rule.Name, "/routes/1/rule")
	if _, err := service.etcdClient.EctdClientV3.Delete(context.Background(), key); nil != err {
		log.Println("delete rule to etcd delete frontend error: ", err)
		return err
	}
	return nil
}
