package main

import (
	"net/http"
	"log"
	"github.com/cloudware-controller/configs"
	"github.com/cloudware-controller/handler"
)

func main() {

	configs.InitEnv()

	etcd := handler.GetEtcdHandler()
	go etcd.Watcher()
	http.Handle("/api/", handler.CreateHTTPAPIHandler(etcd))

	http.Handle("/api/", handler.CreateHTTPAPIHandler(nil))

	log.Println("Start rest server")
	log.Fatal(http.ListenAndServe(configs.ServerAddress, nil))

}
