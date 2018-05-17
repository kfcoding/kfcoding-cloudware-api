package handler

import (
	"net/http"
	"github.com/emicklei/go-restful"
	"github.com/kfcoding-ingress-controller/kftype"
	"github.com/kubernetes/dashboard/src/app/backend/resource/ingress"
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
			To(apiHandler.handleAddIngressRule).
			Writes(ingress.IngressDetail{}))

	apiV1Ws.Route(
		apiV1Ws.DELETE("/ingress/{namespace}/{ingress}/{pod}").
			To(apiHandler.handleDeleteIngressRule).
			Writes(ingress.IngressDetail{}))

	wsContainer := restful.NewContainer()
	wsContainer.EnableContentEncoding(true)
	wsContainer.Add(apiV1Ws)

	return wsContainer
}

func (apiHandler *APIHandler) handleAddIngressRule(request *restful.Request, response *restful.Response) {

	namespace := request.PathParameter("namespace")
	ingress := request.PathParameter("ingress")
	pod := request.PathParameter("pod")

	if namespace == "" || ingress == "" || pod == "" {
		response.WriteError(http.StatusInternalServerError, errors.New("参数不完整"))
		return
	}

	done := make(chan string)

	req := &kftype.Request{
		Option:    kftype.INGRESS_RULE_ADD,
		Done:      done,
		Namespace: namespace,
		Ingress:   ingress,
		Pod:       pod,
	}

	apiHandler.channel <- req

	result := <-done
	response.Write([]byte(result))
}

func (apiHandler *APIHandler) handleDeleteIngressRule(request *restful.Request, response *restful.Response) {

	namespace := request.PathParameter("namespace")
	ingress := request.PathParameter("ingress")
	pod := request.PathParameter("pod")

	if namespace == "" || ingress == "" || pod == "" {
		response.WriteError(http.StatusInternalServerError, errors.New("参数不完整"))
		return
	}

	done := make(chan string)

	req := &kftype.Request{
		Option:    kftype.INGRESS_RULE_DELETE,
		Done:      done,
		Namespace: namespace,
		Ingress:   ingress,
		Pod:       pod,
	}

	apiHandler.channel <- req

	result := <-done
	response.Write([]byte(result))
}
