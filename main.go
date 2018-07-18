package main

import (
	"net/http"
	"log"
	"github.com/cloudware-controller/configs"
	"github.com/cloudware-controller/apihandler"
	"github.com/cloudware-controller/service"
)

func main() {

	configs.InitEnv()

	keeperService := service.GetKeeperEtcdService()
	http.Handle("/keep/", apihandler.CreateKeeperController(keeperService))

	routingService := service.GetRoutingTraefikService()
	http.Handle("/routing/", apihandler.CreateRoutingController(routingService))

	cloudwareService := service.GetCloudwareK8sService(keeperService, routingService)
	http.Handle("/cloudware/", apihandler.CreateCloudwareController(cloudwareService))

	log.Println("Start rest server")
	log.Fatal(http.ListenAndServe(configs.ServerAddress, nil))

}
