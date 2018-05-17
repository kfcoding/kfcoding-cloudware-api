package handler

import (
	"k8s.io/client/kubernetes/client"
	"github.com/kfcoding-ingress-controller/kftype"
	"log"
	"time"
)

type K8sHandler struct {
	k8sClient *client.APIClient
}

func StartK8sHandler(channel chan *kftype.Request) {

	log.Println("Start k8s handler")

	//c, err := config.LoadKubeConfig()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//client := client.NewAPIClient(c)
	//
	//handler := K8sHandler{
	//k8sClient: client,
	//}

	handler := K8sHandler{
		//k8sClient: client,
	}

	select {
	case request := <-channel:

		log.Println("select request ", *request)

		if request.Option == kftype.INGRESS_RULE_ADD {
			handler.handleAddIngressRule(request)
		} else if request.Option == kftype.INGRESS_RULE_DELETE {
			handler.handleDeleteIngressRule(request)
		}
	}

}

func (handler *K8sHandler) handleAddIngressRule(request *kftype.Request) {

	time.After(5 * time.Second)

	request.Done <- "ok"
}

func (handler *K8sHandler) handleDeleteIngressRule(request *kftype.Request) {

	time.After(5 * time.Second)

	request.Done <- "ok"
}
