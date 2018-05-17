package handler

import (
	"k8s.io/client/kubernetes/client"
	"k8s.io/client/kubernetes/config"
	"github.com/kfcoding-ingress-controller/kftype"
	"log"
	"context"
	"encoding/json"
	"strings"
	"errors"
	"strconv"
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
	body := []byte("[{\"op\":\"add\", \"path\":\"/spec/rules/0\", \"value\":" + rule + "}]")

	var f interface{}
	if err := json.Unmarshal(body, &f); err != nil {
		request.Done <- err
		return
	}

	_, _, err := handler.k8sClient.ExtensionsV1beta1Api.PatchNamespacedIngress(context.Background(), request.Ingress, request.Namespace, f, nil)

	request.Done <- err
}

func (handler *K8sHandler) handleDeleteIngressRule(request *kftype.Request) {

	list, _, err := handler.k8sClient.ExtensionsV1beta1Api.ListNamespacedIngress(context.Background(), request.Namespace, nil)

	if nil != err {
		request.Done <- err
		return
	}

	for _, v := range list.Items {
		if strings.Compare(v.Metadata.Name, request.Ingress) == 0 {
			var i = 0
			for _, rule := range v.Spec.Rules {
				if strings.Compare(rule.Host, request.Pod+".kfcoding.com") == 0 {
					body := []byte("[{\"op\":\"remove\", \"path\":\"/spec/rules/" + strconv.Itoa(i) + "\"}]")
					var f interface{}
					if err := json.Unmarshal(body, &f); err != nil {
						request.Done <- err
						return
					}
					_, _, err := handler.k8sClient.ExtensionsV1beta1Api.PatchNamespacedIngress(context.Background(), request.Ingress, request.Namespace, f, nil)
					request.Done <- err
					return
				}
				i++

			}

		}
	}

	request.Done <- errors.New("No ingress rule " + request.Pod)
}
