package apihandler

import (
	"net/http"
	"github.com/emicklei/go-restful"
	"github.com/cloudware-controller/service"
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

}

func (controller *RoutingController) handleAddRoutings(request *restful.Request, response *restful.Response) {

}

func (controller *RoutingController) handleDeleteRouting(request *restful.Request, response *restful.Response) {

}

func (controller *RoutingController) handleDeleteRoutings(request *restful.Request, response *restful.Response) {

}
