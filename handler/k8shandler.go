package handler

import (
	"k8s.io/client/kubernetes/client"
	"k8s.io/client/kubernetes/config"
	"github.com/kfcoding-ingress-controller/kftype"
	"log"
	"context"
	"encoding/json"
)

type K8sHandler struct {
	k8sClient *client.APIClient
}

func StartK8sHandler(channel chan *kftype.Request) {

	log.Println("Start k8s handler")

	c, err := config.LoadKubeConfig()
	if err != nil {
		log.Fatal(err)
	}
	client := client.NewAPIClient(c)

	handler := K8sHandler{
		k8sClient: client,
	}

	for {
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
}

func (handler *K8sHandler) handleAddIngressRule(request *kftype.Request) {

	rule :=
		"{" +
			"\"host\": \"" + request.Pod + ".kfcoding.com\"," +
			"\"http\": {" +
			"\"paths\": [{" +
			"\"backend\": {" +
			"\"serviceName\": \"" + request.Pod + "svc" + "\"," +
			"\"servicePort\": 3000" +
			"}" +
			"}]" +
			"}" +
			"}"

	body := []byte("[{\"op\":\"add\", \"path\":\"/spec/rules/1\", \"value\":" + rule + "}]")

	log.Println("[{\"op\":\"add\", \"path\":\"/spec/rules/1\", \"value\":" + rule + "}]")

	var f interface{}
	err := json.Unmarshal(body, &f)
	result, resp, err := handler.k8sClient.ExtensionsV1beta1Api.PatchNamespacedIngress(context.Background(), request.Ingress, request.Namespace, f, nil)

	log.Println("result: ", result)
	log.Println("resp: ", resp)

	request.Done <- err
}

func (handler *K8sHandler) handleDeleteIngressRule(request *kftype.Request) {

	request.Done <- nil
}
