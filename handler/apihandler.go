package handler

import (
	"net/http"
	"github.com/emicklei/go-restful"
	"github.com/cloudware-controller/kftype"
	"log"
	"strings"
	"github.com/cloudware-controller/configs"
)

type APIHandler struct {
	etcdHandler *EtcdHandler
}

func CreateHTTPAPIHandler(etcd *EtcdHandler) (http.Handler) {

	apiHandler := APIHandler{
		etcdHandler: etcd,
	}

	apiV1Ws := new(restful.WebService)

	apiV1Ws.Path("/api/").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	apiV1Ws.Route(
		apiV1Ws.POST("/cloudware/routing").
			To(apiHandler.handleAddRouting))

	apiV1Ws.Route(
		apiV1Ws.DELETE("/cloudware/routing").
			To(apiHandler.handleDeleteRouting))

	apiV1Ws.Route(
		apiV1Ws.POST("/cloudware/keepalive").
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

	log.Print("handleKeepAlive: ", body)

	err := apiHandler.etcdHandler.handleKeepCloudwareAlive(body)

	if err != nil {
		response.WriteHeaderAndEntity(http.StatusInternalServerError, kftype.Response{Content: err.Error()})
	} else {
		response.WriteHeaderAndEntity(http.StatusOK, kftype.Response{Content: ""})
	}

}

func (apiHandler *APIHandler) handleAddRouting(request *restful.Request, response *restful.Response) {

	// check token
	if apiHandler.checkToken(request, response) == false {
		return
	}

	// get body
	body := &kftype.Request{}
	if err := request.ReadEntity(body); nil != err {
		response.WriteHeaderAndEntity(http.StatusInternalServerError, kftype.Response{Content: err.Error()})
		return
	}
	if body.Name == "" || body.URL == "" || body.URL == "" {
		response.WriteHeaderAndEntity(http.StatusInternalServerError, kftype.Response{Content: "Incomplete parameters"})
		return
	}

	log.Print("handleAddRouting: ", body)

	// add routing
	err := apiHandler.etcdHandler.handleAddRouting(body)

	// return
	if err != nil {
		response.WriteHeaderAndEntity(http.StatusInternalServerError, kftype.Response{Content: err.Error()})
	} else {
		response.WriteHeaderAndEntity(http.StatusOK, kftype.Response{Content: ""})
	}

}

func (apiHandler *APIHandler) handleDeleteRouting(request *restful.Request, response *restful.Response) {

	// check token
	if apiHandler.checkToken(request, response) == false {
		return
	}

	// get body
	body := &kftype.Request{}
	if err := request.ReadEntity(body); nil != err {
		response.WriteHeaderAndEntity(http.StatusInternalServerError, kftype.Response{Content: err.Error()})
		return
	}
	if body.Name == "" {
		response.WriteHeaderAndEntity(http.StatusInternalServerError, kftype.Response{Content: "Incomplete parameters"})
		return
	}

	log.Print("handleDeleteRouting: ", body)

	// delete routing
	err := apiHandler.etcdHandler.handleDeleteRouting(body)

	// return
	if err != nil {
		response.WriteHeaderAndEntity(http.StatusInternalServerError, kftype.Response{Content: err.Error()})
	} else {
		response.WriteHeaderAndEntity(http.StatusOK, kftype.Response{Content: ""})
	}
}

func (apiHandler *APIHandler) checkToken(request *restful.Request, response *restful.Response) bool {
	token := request.HeaderParameter("Authorization")
	if strings.Compare(token, configs.Token) != 0 {
		response.WriteHeaderAndEntity(http.StatusUnauthorized, kftype.Response{Content: ""})
		return false
	}
	return true
}
