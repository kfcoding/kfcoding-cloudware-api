package handler

import (
	"net/http"
	"github.com/emicklei/go-restful"
	"github.com/kfcoding-ingress-controller/kftype"
	"errors"
)

type APIHandler struct {
	channel chan *kftype.Request
}

func CreateHTTPAPIHandler(channel chan *kftype.Request) (http.Handler) {

	apiHandler := APIHandler{
		channel: channel,
	}

	apiV1Ws := new(restful.WebService)

	apiV1Ws.Path("/apis/extensions/v1beta1").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	apiV1Ws.Route(
		apiV1Ws.PUT("/ingress/{namespace}/{ingress}/{pod}").
			To(apiHandler.handleAddIngressRule))

	apiV1Ws.Route(
		apiV1Ws.DELETE("/ingress/{namespace}/{ingress}/{pod}").
			To(apiHandler.handleDeleteIngressRule))

	wsContainer := restful.NewContainer()
	wsContainer.EnableContentEncoding(true)
	wsContainer.Add(apiV1Ws)

	return wsContainer
}

func (apiHandler *APIHandler) handleAddIngressRule(request *restful.Request, response *restful.Response) {

	if !apiHandler.checkCache(response) {
		return
	}

	namespace := request.PathParameter("namespace")
	ingress := request.PathParameter("ingress")
	pod := request.PathParameter("pod")

	if namespace == "" || ingress == "" || pod == "" {
		response.WriteError(http.StatusInternalServerError, errors.New("Incomplete parameters"))
		return
	}

	done := make(chan error)

	req := &kftype.Request{
		Option:    kftype.INGRESS_RULE_ADD,
		Done:      done,
		Namespace: namespace,
		Ingress:   ingress,
		Pod:       pod,
	}

	apiHandler.channel <- req

	err := <-done

	close(done)

	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
	} else {
		response.WriteEntity(nil)
	}

}

func (apiHandler *APIHandler) handleDeleteIngressRule(request *restful.Request, response *restful.Response) {

	if !apiHandler.checkCache(response) {
		return
	}

	namespace := request.PathParameter("namespace")
	ingress := request.PathParameter("ingress")
	pod := request.PathParameter("pod")

	if namespace == "" || ingress == "" || pod == "" {
		response.WriteError(http.StatusInternalServerError, errors.New("Incomplete parameters"))
		return
	}

	done := make(chan error)

	req := &kftype.Request{
		Option:    kftype.INGRESS_RULE_DELETE,
		Done:      done,
		Namespace: namespace,
		Ingress:   ingress,
		Pod:       pod,
	}

	apiHandler.channel <- req

	err := <-done

	close(done)

	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
	} else {
		response.WriteEntity(nil)
	}

}

func (apiHandler *APIHandler) checkCache(response *restful.Response) bool {
	if cap(apiHandler.channel)-len(apiHandler.channel) <= 0 {
		response.WriteError(http.StatusInternalServerError, errors.New("The queue is full"))
		return false
	}
	return true
}
