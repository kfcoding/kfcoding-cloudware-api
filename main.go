package main

import (
	"net/http"
	"log"
	"github.com/kfcoding-cloudware-controller/configs"
	"github.com/kfcoding-cloudware-controller/apihandler"
	"github.com/kfcoding-cloudware-controller/service"
	"github.com/kfcoding-cloudware-controller/etcd"
	"path"
)

func main() {

	configs.InitEnv()

	etcdClient := etcd.GetMyEtcdClient()
	keeperService := service.GetKeeperEtcdService(etcdClient)
	routingService := service.GetRoutingTraefikService(etcdClient)
	cloudwareService := service.GetCloudwareK8sService(keeperService, routingService)

	watcher := service.GetEtcdWatcher(etcdClient)
	go watcher.Watcher(path.Join(configs.KeeperPrefix, configs.Version), cloudwareService)

	http.Handle("/keep/", apihandler.CreateKeeperController(keeperService))
	http.Handle("/routing/", apihandler.CreateRoutingController(routingService))
	http.Handle("/cloudware/", apihandler.CreateCloudwareController(cloudwareService))

	log.Println("Start cloudware server")
	log.Fatal(http.ListenAndServe(configs.ServerAddress, nil))

}
