package apihandler

import (
	"net/http"
	"github.com/emicklei/go-restful"
	"github.com/cloudware-controller/service"
	"github.com/cloudware-controller/types"
	"log"
)

type CloudwareController struct {
	cloudwareService service.CloudwareService
}

func CreateCloudwareController(cloudwareService service.CloudwareService) (http.Handler) {

	cloudwareController := CloudwareController{
		cloudwareService: cloudwareService,
	}

	apiV1Ws := new(restful.WebService)

	apiV1Ws.Path("/cloudware/").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	apiV1Ws.Route(
		apiV1Ws.POST("/").
			To(cloudwareController.handleCreateCloudware))

	wsContainer := restful.NewContainer()
	wsContainer.EnableContentEncoding(true)
	wsContainer.Add(apiV1Ws)

	return wsContainer
}

func (controller *CloudwareController) handleCreateCloudware(request *restful.Request, response *restful.Response) {
	body := &types.CloudwareBody{}
	if err := request.ReadEntity(body); nil != err {
		log.Print("handleCreateCloudware: ", err)
		response.WriteHeaderAndEntity(
			http.StatusInternalServerError,
			types.ResponseBody{Error: err.Error()})
		return
	}
	log.Print("handleCreateCloudware: ", body)

	data, err := controller.cloudwareService.CreateCloudwareApi(body)

	if err == nil {
		response.WriteHeaderAndEntity(http.StatusOK, types.ResponseBody{Data: data})
	} else {
		log.Print("handleCreateCloudware: ", err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, types.ResponseBody{Error: err.Error()})
	}
}
