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

			if request.Option == kftype.IngressRoleAdd {
				handler.handleAddIngressRule(request)
			} else if request.Option == kftype.IngressRoleDelete {
				handler.handleDeleteIngressRule(request)
			}
		}
	}

}

func (handler *K8sHandler) handleAddIngressRule(request *kftype.Request) {

	rule := "{" +
		"\"host\": \"" + request.Pod + ".cloudware.kfcoding.com\"," +
		"\"http\": {" +
		"\"paths\": [{" +
		"\"backend\": {" +
		"\"serviceName\": \"" + request.Pod + "\"," +
		"\"servicePort\": 9800" +
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

	if nil != err {
		log.Print("handleAddIngressRule error", err)
	}

	request.Done <- err
}

func (handler *K8sHandler) handleDeleteIngressRule(request *kftype.Request) {

	list, _, err := handler.k8sClient.ExtensionsV1beta1Api.ListNamespacedIngress(context.Background(), request.Namespace, nil)

	if nil != err {
		log.Print("handleDeleteIngressRule, ListNamespacedIngress", err)
		request.Done <- err
		return
	}

	for _, v := range list.Items {
		if strings.Compare(v.Metadata.Name, request.Ingress) == 0 {
			for i, rule := range v.Spec.Rules {
				if strings.Compare(rule.Host, request.Pod+".cloudware.kfcoding.com") == 0 {
					body := []byte("[{\"op\":\"remove\", \"path\":\"/spec/rules/" + strconv.Itoa(i) + "\"}]")
					var f interface{}
					if err := json.Unmarshal(body, &f); err != nil {
						log.Print("handleDeleteIngressRule, Unmarshal error", err)
						request.Done <- err
						return
					}
					_, _, err := handler.k8sClient.ExtensionsV1beta1Api.PatchNamespacedIngress(context.Background(), request.Ingress, request.Namespace, f, nil)
					log.Print("handleDeleteIngressRule, PatchNamespacedIngress", err)
					request.Done <- err
					return
				}
			}
		}
	}

	log.Print("handleDeleteIngressRule", ", No ingress rule "+request.Pod)
	request.Done <- errors.New("No ingress rule " + request.Pod)
}
