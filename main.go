package main

import (
	"net/http"
	"github.com/kfcoding-ingress-controller/handler"
	"github.com/kfcoding-ingress-controller/kftype"
	"log"
)

func main() {

	const QUEUE_SIZE = 1000

	channel := make(chan *kftype.Request, QUEUE_SIZE)

	go handler.StartK8sHandler(channel)

	http.Handle("/apis/", handler.CreateHTTPAPIHandler(channel))

	log.Println("Start rest server")

	log.Fatal(http.ListenAndServe("0.0.0.0:9090", nil))

}
