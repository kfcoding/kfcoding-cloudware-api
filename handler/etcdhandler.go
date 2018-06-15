package handler

import (
	"github.com/cloudware-controller/kftype"
	"log"
	"context"
	"github.com/cloudware-controller/etcd"
	"github.com/coreos/etcd/clientv3"
	"github.com/cloudware-controller/configs"
	"fmt"
	"github.com/kfcoding-cloudware-traefik/config"
	"net/http"
	"io/ioutil"
)

type EtcdHandler struct {
	etcdClient *etcd.MyEtcdClient
}

func GetEtcdHandler() *EtcdHandler {
	return &EtcdHandler{
		etcdClient: etcd.GetMyEtcdClient(),
	}
}

func (handler *EtcdHandler) Watcher() {
	log.Println("Start Watcher")

	rch := handler.etcdClient.EctdClientV3.Watch(context.Background(), configs.PrefixAlive, clientv3.WithPrefix())

	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			switch ev.Type {
			case 1: //DELETE
				body := &kftype.Request{Pod: string(ev.Kv.Value)}
				handler.handleDeleteRouting(body)
				handler.handleDeleteContaienr(body)
			}
		}
	}

}

func (handler *EtcdHandler) handleKeepCloudwareAlive(request *kftype.Request) error {

	var resp *clientv3.LeaseGrantResponse
	var err error
	if resp, err = handler.etcdClient.EctdClientV3.Grant(context.TODO(), configs.CloudWareTTL); err != nil {
		log.Println("handleKeepCloudwareAlive error: ", err)
		return err
	}

	key := configs.PrefixAlive + request.Pod
	if _, err = handler.etcdClient.EctdClientV3.Put(context.TODO(), key, request.Pod, clientv3.WithLease(resp.ID)); nil != err {
		log.Println("handleKeepCloudwareAlive error: ", err)
		return err
	}

	return nil
}

func (handler *EtcdHandler) handleAddRouting(request *kftype.Request) (error) {

	// set backend
	key := configs.PrefixTraefik + "backends/" + request.Pod + "/servers/1/url"
	value := request.URL
	if _, err := handler.etcdClient.EctdClientV3.Put(context.Background(), key, value); nil != err {
		log.Println("handleAddRouting put backend error: ", err)
		return err
	}

	//set frontend
	key = config.PrefixTraefik + "frontends/" + request.Pod + "/backend"
	value = request.Pod
	if _, err := handler.etcdClient.EctdClientV3.Put(context.Background(), key, value); nil != err {
		log.Println("handleAddRouting put frontend error: ", err)
		return err
	}
	key = config.PrefixTraefik + "frontends/" + request.Pod + "/routes/1/rule"
	value = "Path:/" + request.Pod
	if _, err := handler.etcdClient.EctdClientV3.Put(context.Background(), key, value); nil != err {
		log.Println("handleAddRouting put frontend error: ", err)
		return err
	}

	// set keepalive
	if err := handler.handleKeepCloudwareAlive(request); err != nil {
		return err
	}

	return nil
}

func (handler *EtcdHandler) handleDeleteRouting(request *kftype.Request) (error) {
	// delete backend
	key := configs.PrefixTraefik + "backends/" + request.Pod + "/servers/1/url"
	if _, err := handler.etcdClient.EctdClientV3.Delete(context.Background(), key); nil != err {
		log.Println("handleDeleteRouting delete backend error: ", err)
		return err
	}

	// delete frontend
	key = config.PrefixTraefik + "frontends/" + request.Pod + "/backend"
	if _, err := handler.etcdClient.EctdClientV3.Delete(context.Background(), key); nil != err {
		log.Println("handleAddRouting put frontend error: ", err)
		return err
	}
	key = config.PrefixTraefik + "frontends/" + request.Pod + "/routes/1/rule"
	if _, err := handler.etcdClient.EctdClientV3.Delete(context.Background(), key); nil != err {
		log.Println("handleAddRouting put frontend error: ", err)
		return err
	}
	return nil
}

func (handler *EtcdHandler) handleDeleteContaienr(body *kftype.Request) {
	url := configs.ApiServerAddress + "/cloudware/deleteContainer?podName=" + body.Pod + "&type=0"
	req, _ := http.NewRequest("DELETE", url, nil)
	res, _ := http.DefaultClient.Do(req)
	resp, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(resp))
	res.Body.Close()
}
