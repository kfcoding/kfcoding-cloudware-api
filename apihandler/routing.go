package apihandler

import (
	"net/http"
	"github.com/emicklei/go-restful"
	"github.com/kfcoding-cloudware-controller/service"
	"github.com/go-simplejson"
	"github.com/kfcoding-cloudware-controller/types"
	"log"
	"strings"
	"github.com/kfcoding-cloudware-controller/configs"
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

	cors := restful.CrossOriginResourceSharing{
		AllowedMethods: []string{"POST", "OPTIONS", "GET"},
		AllowedHeaders: []string{"Authorization", "Content-Type", "Accept", "Token"},
		CookiesAllowed: false,
		Container:      wsContainer}
	wsContainer.Filter(cors.Filter)

	return wsContainer
}

func (controller *RoutingController) handleAddRouting(request *restful.Request, response *restful.Response) {
	if !controller.checkToken(request, response) {
		return
	}
	body := &types.RoutingBody{}
	if err := request.ReadEntity(body); nil != err {
		log.Printf("handleAddRouting error: %+v\n", err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, types.ResponseBody{Error: err.Error()})
		return
	}
	if body.Name == "" || body.URL == "" || body.Rule == "" {
		log.Printf("handleAddRouting error: 参数不正确")
		response.WriteHeaderAndEntity(http.StatusInternalServerError, types.ResponseBody{Error: "参数不正确"})
		return
	}
	log.Printf("handleAddRouting: %+v\n", body)

	err := controller.routingService.AddRule(body)

	if err != nil {
		response.WriteHeaderAndEntity(http.StatusOK, types.ResponseBody{Error: err.Error()})
	} else {
		response.WriteHeaderAndEntity(http.StatusOK, types.ResponseBody{Data: "ok"})
	}
}

func (controller *RoutingController) handleDeleteRouting(request *restful.Request, response *restful.Response) {
	if !controller.checkToken(request, response) {
		return
	}
	body := &types.RoutingBody{}
	if err := request.ReadEntity(body); nil != err {
		log.Printf("handleDeleteRouting error: %+v\n", err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, types.ResponseBody{Error: err.Error()})
		return
	}
	if body.Name == "" {
		log.Printf("handleDeleteRouting error: 参数不正确")
		response.WriteHeaderAndEntity(http.StatusInternalServerError, types.ResponseBody{Error: "参数不正确"})
		return
	}
	log.Print("handleDeleteRouting: ", body)

	err := controller.routingService.AddRule(body)

	if err != nil {
		response.WriteHeaderAndEntity(http.StatusOK, types.ResponseBody{Error: err.Error()})
	} else {
		response.WriteHeaderAndEntity(http.StatusOK, types.ResponseBody{Data: "ok"})
	}
}

func (controller *RoutingController) handleDeleteRoutings(request *restful.Request, response *restful.Response) {
	if !controller.checkToken(request, response) {
		return
	}
	// get body
	json := &simplejson.Json{}
	if err := request.ReadEntity(json); nil != err {
		log.Printf("handleDeleteRoutings error: %+v\n", err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, types.ResponseBody{Data: err.Error()})
		return
	}
	rules, err := json.Get("rules").Array()
	if nil != err {
		log.Printf("handleDeleteRoutings error: %+v\n", err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, types.ResponseBody{Data: err.Error()})
		return
	}
	log.Printf("handleDeleteRoutings: %+v\n", json)

	err = controller.routingService.DeleteRules(rules)

	// return
	if err != nil {
		response.WriteHeaderAndEntity(http.StatusInternalServerError, types.ResponseBody{Error: err.Error()})
	} else {
		response.WriteHeaderAndEntity(http.StatusOK, types.ResponseBody{Data: "ok"})
	}

}

func (controller *RoutingController) handleAddRoutings(request *restful.Request, response *restful.Response) {
	if !controller.checkToken(request, response) {
		return
	}
	// get body
	json := &simplejson.Json{}
	if err := request.ReadEntity(json); nil != err {
		log.Printf("handleAddRoutings error: %+v\n", err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, types.ResponseBody{Data: err.Error()})
		return
	}
	rules, err := json.Get("rules").Array()
	if nil != err {
		log.Printf("handleAddRoutings error: %+v\n", err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, types.ResponseBody{Data: err.Error()})
		return
	}
	log.Printf("handleAddRoutings: %+v\n", json)

	err = controller.routingService.AddRules(rules)

	if err != nil {
		response.WriteHeaderAndEntity(http.StatusOK, types.ResponseBody{Error: err.Error()})
	} else {
		response.WriteHeaderAndEntity(http.StatusOK, types.ResponseBody{Data: "ok"})
	}

}

func (controller *RoutingController) checkToken(request *restful.Request, response *restful.Response) bool {
	token := request.HeaderParameter("Token")
	if strings.Compare(token, configs.Token) != 0 {
		log.Print("认证失败")
		response.WriteHeaderAndEntity(http.StatusUnauthorized, types.ResponseBody{Error: "认证失败"})
		return false
	}
	return true
}
