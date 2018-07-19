package apihandler

import (
	"net/http"
	"github.com/emicklei/go-restful"
	"github.com/cloudware-controller/service"
	"github.com/go-simplejson"
	"github.com/cloudware-controller/types"
	"log"
)

type RoutingController struct {
	routingService service.RoutingService
}

func CreateRoutingController(routingService service.RoutingService) (http.Handler) {

	routingController := RoutingController{
		routingService: routingService,
	}

	apiV1Ws := new(restful.WebService)

	apiV1Ws.Path("/routing").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	apiV1Ws.Route(
		apiV1Ws.POST("/add").
			To(routingController.handleAddRouting))

	apiV1Ws.Route(
		apiV1Ws.POST("/adds").
			To(routingController.handleAddRoutings))

	apiV1Ws.Route(
		apiV1Ws.POST("/delete").
			To(routingController.handleDeleteRouting))

	apiV1Ws.Route(
		apiV1Ws.POST("/deletes").
			To(routingController.handleDeleteRoutings))

	wsContainer := restful.NewContainer()
	wsContainer.EnableContentEncoding(true)
	wsContainer.Add(apiV1Ws)

	return wsContainer
}

func (controller *RoutingController) handleAddRouting(request *restful.Request, response *restful.Response) {
	body := &types.RoutingBody{}
	if err := request.ReadEntity(body); nil != err {
		log.Print("handleAddRouting error: ", err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, types.ResponseBody{Error: err.Error()})
		return
	}
	log.Print("handleAddRouting: ", body)

	err := controller.routingService.AddRule(body)

	if err != nil {
		response.WriteHeaderAndEntity(http.StatusOK, types.ResponseBody{Error: err.Error()})
	} else {
		response.WriteHeaderAndEntity(http.StatusOK, types.ResponseBody{Data: "ok"})
	}
}

func (controller *RoutingController) handleAddRoutings(request *restful.Request, response *restful.Response) {
	// get body
	json := &simplejson.Json{}
	if err := request.ReadEntity(json); nil != err {
		log.Print("handleAddRoutings error: ", err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, types.ResponseBody{Data: err.Error()})
		return
	}
	rules, err := json.Get("rules").Array()
	if nil != err {
		log.Print("handleAddRoutings error: ", err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, types.ResponseBody{Data: err.Error()})
		return
	}
	log.Print("handleAddRoutings: ", json)

	err = controller.routingService.AddRules(rules)

	if err != nil {
		response.WriteHeaderAndEntity(http.StatusOK, types.ResponseBody{Error: err.Error()})
	} else {
		response.WriteHeaderAndEntity(http.StatusOK, types.ResponseBody{Data: "ok"})
	}

}

func (controller *RoutingController) handleDeleteRouting(request *restful.Request, response *restful.Response) {
	body := &types.RoutingBody{}
	if err := request.ReadEntity(body); nil != err {
		log.Print("handleAddRouting error: ", err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, types.ResponseBody{Error: err.Error()})
		return
	}
	log.Print("handleAddRouting: ", body)

	err := controller.routingService.AddRule(body)

	if err != nil {
		response.WriteHeaderAndEntity(http.StatusOK, types.ResponseBody{Error: err.Error()})
	} else {
		response.WriteHeaderAndEntity(http.StatusOK, types.ResponseBody{Data: "ok"})
	}
}

func (controller *RoutingController) handleDeleteRoutings(request *restful.Request, response *restful.Response) {
	// get body
	json := &simplejson.Json{}
	if err := request.ReadEntity(json); nil != err {
		log.Print("handleAddRoutings error: ", err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, types.ResponseBody{Data: err.Error()})
		return
	}
	rules, err := json.Get("rules").Array()
	if nil != err {
		log.Print("handleAddRoutings error: ", err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, types.ResponseBody{Data: err.Error()})
		return
	}
	log.Print("handleAddRoutings: ", json)

	err = controller.routingService.DeleteRules(rules)

	// return
	if err != nil {
		response.WriteHeaderAndEntity(http.StatusInternalServerError, types.ResponseBody{Error: err.Error()})
	} else {
		response.WriteHeaderAndEntity(http.StatusOK, types.ResponseBody{Data: "ok"})
	}

}
