package handler

import (
	"net/http"
	"github.com/emicklei/go-restful"
	"github.com/kfcoding-ingress-controller/kftype"
	"errors"
	"github.com/kfcoding-ingress-controller/etcd"
	"log"
)

type APIHandler struct {
	channel    chan *kftype.Request
	etcdClient *etcd.MyEtcdClient
}

func CreateHTTPAPIHandler(channel chan *kftype.Request) (http.Handler) {

	apiHandler := APIHandler{
		channel:    channel,
		etcdClient: etcd.GetMyEtcdClient(),
	}

	apiV1Ws := new(restful.WebService)

	apiV1Ws.Path("/api/").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	apiV1Ws.Route(
		apiV1Ws.PUT("/cloudware").
			To(apiHandler.handleAddIngressRule))

	apiV1Ws.Route(
		apiV1Ws.DELETE("/cloudware").
			To(apiHandler.handleDeleteIngressRule))

	apiV1Ws.Route(
		apiV1Ws.PUT("/cloudware/keepalive").
			To(apiHandler.handleKeepAlive))

	wsContainer := restful.NewContainer()
	wsContainer.EnableContentEncoding(true)
	wsContainer.Add(apiV1Ws)

	return wsContainer
}

func (apiHandler *APIHandler) handleKeepAlive(request *restful.Request, response *restful.Response) {

	body := &kftype.Request{}
	if err := request.ReadEntity(body); nil != err {
		response.WriteHeaderAndEntity(http.StatusInternalServerError, kftype.Response{Content: err.Error()})
		return
	}

	log.Print("handleKeepAlive: ", body.Pod)

	err := handleKeepCloudwareAlive(body, apiHandler.etcdClient)

	if err != nil {
		response.WriteHeaderAndEntity(http.StatusInternalServerError, kftype.Response{Content: err.Error()})
	} else {
		response.WriteHeaderAndEntity(http.StatusOK, kftype.Response{Content: ""})
	}

}

func (apiHandler *APIHandler) handleAddIngressRule(request *restful.Request, response *restful.Response) {

	if !apiHandler.checkCache(response) {
		return
	}

	body := &kftype.Request{}
	if err := request.ReadEntity(body); nil != err {
		response.WriteHeaderAndEntity(http.StatusInternalServerError, kftype.Response{Content: err.Error()})
		return
	}

	if body.Namespace == "" || body.Ingress == "" || body.Pod == "" {
		response.WriteHeaderAndEntity(http.StatusInternalServerError, kftype.Response{Content: "Incomplete parameters"})
		return
	}

	log.Print("handleAddIngressRule: ", body.Pod)

	done := make(chan error)
	body.Option = kftype.IngressRoleAdd
	body.Done = done

	apiHandler.channel <- body

	err := <-done

	close(done)

	if err != nil {
		response.WriteHeaderAndEntity(http.StatusInternalServerError, kftype.Response{Content: err.Error()})
	} else {
		response.WriteHeaderAndEntity(http.StatusOK, kftype.Response{Content: ""})
	}

}

func (apiHandler *APIHandler) handleDeleteIngressRule(request *restful.Request, response *restful.Response) {

	if !apiHandler.checkCache(response) {
		return
	}

	body := &kftype.Request{}
	if err := request.ReadEntity(body); nil != err {
		response.WriteHeaderAndEntity(http.StatusInternalServerError, kftype.Response{Content: err.Error()})
		return
	}

	if body.Namespace == "" || body.Ingress == "" || body.Pod == "" {
		response.WriteHeaderAndEntity(http.StatusInternalServerError, kftype.Response{Content: "Incomplete parameters"})
		return
	}

	log.Print("handleAddIngressRule: ", body.Pod)

	done := make(chan error)
	body.Option = kftype.IngressRoleDelete
	body.Done = done

	apiHandler.channel <- body

	err := <-done

	close(done)

	if err != nil {
		response.WriteHeaderAndEntity(http.StatusInternalServerError, kftype.Response{Content: err.Error()})
	} else {
		response.WriteHeaderAndEntity(http.StatusOK, kftype.Response{Content: ""})
	}

}

func (apiHandler *APIHandler) checkCache(response *restful.Response) bool {
	if cap(apiHandler.channel)-len(apiHandler.channel) <= 0 {
		response.WriteError(http.StatusInternalServerError, errors.New("The queue is full"))
		return false
	}
	return true
}
