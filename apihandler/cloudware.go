package apihandler

import (
	"net/http"
	"github.com/emicklei/go-restful"
	"github.com/kfcoding-cloudware-controller/service"
	"github.com/kfcoding-cloudware-controller/types"
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
	apiV1Ws.Path("/cloudware").
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
		log.Print("handleCreateCloudware error: ", err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, types.ResponseBody{Error: err.Error()})
		return
	}
	if body.Image == "" {
		log.Print("handleCreateCloudware error: Image 不能为空")
		response.WriteHeaderAndEntity(http.StatusInternalServerError, types.ResponseBody{Error: "Image 不能为空"})
		return
	}
	log.Printf("handleCreateCloudware: %+v\n", body)

	data, err := controller.cloudwareService.CreateCloudwareApi(body)

	if err == nil {
		response.WriteHeaderAndEntity(http.StatusOK, types.ResponseBody{Data: data})
	} else {
		response.WriteHeaderAndEntity(http.StatusInternalServerError, types.ResponseBody{Error: err.Error()})
	}
}
