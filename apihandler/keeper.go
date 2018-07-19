package apihandler

import (
	"net/http"
	"github.com/emicklei/go-restful"
	"github.com/cloudware-controller/service"
	"github.com/cloudware-controller/types"
	"log"
)

type KeeperController struct {
	keeperService service.KeeperService
}

func CreateKeeperController(keeperService service.KeeperService) (http.Handler) {

	keeperController := KeeperController{
		keeperService: keeperService,
	}

	apiV1Ws := new(restful.WebService)
	apiV1Ws.Path("/keep/").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	apiV1Ws.Route(
		apiV1Ws.POST("/cloudware").
			To(keeperController.handleKeepAlive))

	wsContainer := restful.NewContainer()
	wsContainer.EnableContentEncoding(true)
	wsContainer.Add(apiV1Ws)

	return wsContainer
}

func (controller *KeeperController) handleKeepAlive(request *restful.Request, response *restful.Response) {
	body := &types.KeeperBody{}
	if err := request.ReadEntity(body); nil != err {
		log.Printf("handleKeepAlive error: %+v\n", err)
		response.WriteHeaderAndEntity(
			http.StatusInternalServerError,
			types.ResponseBody{Error: err.Error()})
		return
	}
	log.Printf("handleKeepAlive: %+v\n", body)

	if !controller.keeperService.Check(body) {
		log.Print("Keep ok, " + body.Name + " not exist")
		return
	}

	controller.keeperService.Keep(body)
	response.WriteHeaderAndEntity(http.StatusOK, types.ResponseBody{})
}
