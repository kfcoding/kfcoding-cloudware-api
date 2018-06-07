package main

import (
	"net/http"
	"github.com/kfcoding-ingress-controller/handler"
	"github.com/kfcoding-ingress-controller/kftype"
	"log"
	"github.com/kfcoding-ingress-controller/configs"
)

func main() {

	configs.InitEnv()

	channel := make(chan *kftype.Request, configs.QueueSize)

	go handler.StartK8sHandler(channel)

	http.Handle("/api/", handler.CreateHTTPAPIHandler(channel))

	log.Println("Start rest server")

	log.Fatal(http.ListenAndServe(configs.ServerAddress, nil))

}
