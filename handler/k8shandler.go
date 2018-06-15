package handler

import (
	"github.com/cloudware-controller/kftype"
	"log"
	"context"
	"encoding/json"
	"strings"
	"errors"
	"strconv"
	"github.com/cloudware-controller/etcd"
	"github.com/coreos/etcd/clientv3"
	"github.com/cloudware-controller/configs"
	"fmt"
	"net/http"
	"io/ioutil"

	"k8s.io/client-go/kubernetes"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/rest"
)

type K8sHandler struct {
	k8sClient  *kubernetes.Clientset
	etcdClient *etcd.MyEtcdClient
}

func getK8sClientOutOfCluster() *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

func getK8sClientInCluster() *kubernetes.Clientset {

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	return clientset
}

func StartK8sHandler(channel chan *kftype.Request) {

	client := getK8sClientInCluster()

	//client := getK8sClientOutOfCluster()

	handler := K8sHandler{
		k8sClient:  client,
		etcdClient: etcd.GetMyEtcdClient(),
	}

	go handler.Watcher()

	for {
		select {
		case request := <-channel:

			log.Println("select request ", *request)

			if request.Option == kftype.IngressRoleAdd {
				handler.handleAddIngressRule(request)
			} else if request.Option == kftype.IngressRoleDelete {
				request.Done <- handler.handleDeleteIngressRule(request)
			}
		}
	}

}

func (handler *K8sHandler) Watcher() {
	log.Println("Start Watcher")

	rch := handler.etcdClient.EctdClientV3.Watch(context.Background(), configs.PrefixAlive, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			switch ev.Type {
			case 1: //DELETE
				strs := strings.Split(string(ev.Kv.Key), "/")
				body := &kftype.Request{}
				body.Namespace = configs.Namespace
				body.Pod = strs[len(strs)-1]
				body.Ingress = configs.IngressName

				handler.handleDeleteIngressRule(body)

				url := configs.ApiServerAddress + "/cloudware/deleteContainer?podName=" + body.Pod + "&type=0"
				req, _ := http.NewRequest("DELETE", url, nil)
				res, _ := http.DefaultClient.Do(req)
				defer res.Body.Close()
				body1, _ := ioutil.ReadAll(res.Body)
				fmt.Println(string(body1))

			}
		}
	}

}

func handleKeepCloudwareAlive(request *kftype.Request, etcdClient *etcd.MyEtcdClient) error {
	var resp *clientv3.LeaseGrantResponse
	var err error
	if resp, err = etcdClient.EctdClientV3.Grant(context.TODO(), configs.CloudWareTTL); err != nil {
		log.Println("handleKeepCloudwareAlive: ", err)
		return err
	}
	key := configs.PrefixAlive + request.Pod
	if err != nil {
		log.Print(err)
		return err
	}
	if _, err = etcdClient.EctdClientV3.Put(context.TODO(), key, "", clientv3.WithLease(resp.ID)); nil != err {
		log.Println("handleKeepCloudwareAlive: ", err)
		return err
	}
	return nil
}

func (handler *K8sHandler) handleAddIngressRule(request *kftype.Request) {

	rule := "{" +
		"\"host\": \"" + request.Pod + configs.WsAddrSuffix + "\"," +
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

	log.Print(string(body))

	ingress := handler.k8sClient.ExtensionsV1beta1().Ingresses(configs.Namespace)
	result, err := ingress.Patch(configs.IngressName, types.JSONPatchType, []byte(body))
	log.Print(result)

	if nil != err {
		log.Print("handleAddIngressRule error: ", err)
	}

	err = handleKeepCloudwareAlive(request, handler.etcdClient)

	request.Done <- err
}

func (handler *K8sHandler) handleDeleteIngressRule(request *kftype.Request) (error) {
	ingress := handler.k8sClient.ExtensionsV1beta1().Ingresses(configs.Namespace)
	result, err := ingress.Get(configs.IngressName, v1.GetOptions{})

	if nil != err {
		return err
	}

	log.Print(result)

	for k, v := range result.Spec.Rules {
		if strings.Compare(strings.Split(v.Host, ".")[0], request.Pod) == 0 {
			body := []byte("[{\"op\":\"remove\", \"path\":\"/spec/rules/" + strconv.Itoa(k) + "\"}]")
			var f interface{}

			if err := json.Unmarshal(body, &f); err != nil {
				log.Print("handleDeleteIngressRule, Unmarshal error", err)
				return err
			}

			log.Print(string(body))

			result, err = ingress.Patch(configs.IngressName, types.JSONPatchType, body)
			log.Print(result)

			if nil != err {
				log.Print("handleDeleteIngressRule, PatchNamespacedIngress", err)
			}

			return err
		}
	}

	log.Print("handleDeleteIngressRule", ", No ingress rule "+request.Pod)
	return errors.New("No ingress rule " + request.Pod)
}
