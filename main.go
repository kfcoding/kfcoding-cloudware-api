package main

import (
	"net/http"
	"github.com/cloudware-controller/handler"
	"github.com/cloudware-controller/kftype"
	"log"
	"github.com/cloudware-controller/configs"
)

func main() {

	configs.InitEnv()

	channel := make(chan *kftype.Request, configs.QueueSize)

	go handler.StartK8sHandler(channel)

	http.Handle("/api/", handler.CreateHTTPAPIHandler(channel))

	log.Println("Start rest server")

	log.Fatal(http.ListenAndServe(configs.ServerAddress, nil))

}
