package main

import (
	"net/http"
	"log"
	"github.com/cloudware-controller/configs"
	"github.com/cloudware-controller/apihandler"
	"github.com/cloudware-controller/service"
	"github.com/cloudware-controller/etcd"
)

func main() {

	configs.InitEnv()

	etcdClient := etcd.GetMyEtcdClient()
	keeperService := service.GetKeeperEtcdService(etcdClient)
	routingService := service.GetRoutingTraefikService(etcdClient)
	cloudwareService := service.GetCloudwareK8sService(keeperService, routingService)

	watcher := service.GetEtcdWatcher(etcdClient)
	go watcher.Watcher(configs.KeeperPrefix, cloudwareService)

	http.Handle("/keep/", apihandler.CreateKeeperController(keeperService))
	http.Handle("/routing/", apihandler.CreateRoutingController(routingService))
	http.Handle("/cloudware/", apihandler.CreateCloudwareController(cloudwareService))

	log.Println("Start rest server")
	log.Fatal(http.ListenAndServe(configs.ServerAddress, nil))

}
