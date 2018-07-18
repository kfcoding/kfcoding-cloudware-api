package apihandler
//
//import (
//	"net/http"
//	"github.com/emicklei/go-restful"
//	"github.com/cloudware-controller/types"
//	"log"
//	"strings"
//	"github.com/cloudware-controller/configs"
//	"github.com/go-simplejson"
//)
//
//type APIHandler struct {
//	etcdHandler *EtcdHandler
//}
//
//func CreateHTTPAPIHandler(etcd *EtcdHandler) (http.Handler) {
//
//	apiHandler := APIHandler{
//		etcdHandler: etcd,
//	}
//
//	apiV1Ws := new(restful.WebService)
//
//	apiV1Ws.Path("/api/").
//		Consumes(restful.MIME_JSON).
//		Produces(restful.MIME_JSON)
//
//	apiV1Ws.Route(
//		apiV1Ws.POST("/cloudware/routing").
//			To(apiHandler.handleAddCloudwareRouting))
//
//	apiV1Ws.Route(
//		apiV1Ws.DELETE("/cloudware/routing").
//			To(apiHandler.handleDelCloudwareRouting))
//
//	apiV1Ws.Route(
//		apiV1Ws.POST("/traefik/routings").
//			To(apiHandler.handleAddTraefikRoutings))
//
//	apiV1Ws.Route(
//		apiV1Ws.DELETE("/traefik/routings").
//			To(apiHandler.handleDelTraefikRoutings))
//
//	apiV1Ws.Route(
//		apiV1Ws.POST("/cloudware/keepalive").
//			To(apiHandler.handleKeepCloudwareAlive))
//
//	wsContainer := restful.NewContainer()
//	wsContainer.EnableContentEncoding(true)
//	wsContainer.Add(apiV1Ws)
//
//	return wsContainer
//}
//
//func (apiHandler *APIHandler) handleAddCloudwareRouting(request *restful.Request, response *restful.Response) {
//
//	// check token
//	if apiHandler.checkToken(request, response) == false {
//		return
//	}
//
//	// get body
//	body := &types.Request{}
//	if err := request.ReadEntity(body); nil != err {
//		response.WriteHeaderAndEntity(http.StatusInternalServerError, types.Response{Content: err.Error()})
//		return
//	}
//	if body.Name == "" || body.URL == "" || body.URL == "" {
//		response.WriteHeaderAndEntity(http.StatusInternalServerError, types.Response{Content: "Incomplete parameters"})
//		return
//	}
//
//	log.Print("handleAddRouting: ", body)
//
//	// add routing
//	if err := apiHandler.etcdHandler.handleAddRouting(body); err != nil {
//		response.WriteHeaderAndEntity(http.StatusInternalServerError, types.Response{Content: err.Error()})
//		return
//	}
//
//	// keep alive
//	if err := apiHandler.etcdHandler.handleKeepCloudwareAlive(body); err != nil {
//		response.WriteHeaderAndEntity(http.StatusInternalServerError, types.Response{Content: err.Error()})
//		return
//	}
//
//	// return
//	response.WriteHeaderAndEntity(http.StatusOK, types.Response{Content: ""})
//
//}
//
//func (apiHandler *APIHandler) handleDelCloudwareRouting(request *restful.Request, response *restful.Response) {
//
//	// check token
//	if apiHandler.checkToken(request, response) == false {
//		return
//	}
//
//	// get body
//	body := &types.Request{}
//	if err := request.ReadEntity(body); nil != err {
//		response.WriteHeaderAndEntity(http.StatusInternalServerError, types.Response{Content: err.Error()})
//		return
//	}
//	if body.Name == "" {
//		response.WriteHeaderAndEntity(http.StatusInternalServerError, types.Response{Content: "Incomplete parameters"})
//		return
//	}
//
//	log.Print("handleDeleteRouting: ", body)
//
//	// delete routing
//	if err := apiHandler.etcdHandler.handleDeleteRouting(body); err != nil {
//		response.WriteHeaderAndEntity(http.StatusInternalServerError, types.Response{Content: err.Error()})
//		return
//	}
//
//	// return
//	response.WriteHeaderAndEntity(http.StatusOK, types.Response{Content: ""})
//}
//
//func (apiHandler *APIHandler) handleAddTraefikRoutings(request *restful.Request, response *restful.Response) {
//
//	// check token
//	if apiHandler.checkToken(request, response) == false {
//		return
//	}
//
//	// get body
//	json := &simplejson.Json{}
//	if err := request.ReadEntity(json); nil != err {
//		response.WriteHeaderAndEntity(http.StatusInternalServerError, types.Response{Content: err.Error()})
//		return
//	}
//	log.Print("handleAddRoutings: ", json)
//
//	rules, err := json.Get("rules").Array()
//	if nil != err {
//		response.WriteHeaderAndEntity(http.StatusInternalServerError, types.Response{Content: err.Error()})
//		return
//	}
//
//	for _, v := range (rules) {
//		body, _ := v.(map[string]interface{})
//		err := apiHandler.etcdHandler.handleAddRouting(&types.Request{
//			Name: body["Name"].(string),
//			URL:  body["URL"].(string),
//			Rule: body["Rule"].(string),
//		})
//		if err != nil {
//			response.WriteHeaderAndEntity(http.StatusInternalServerError, types.Response{Content: err.Error()})
//			return
//		}
//	}
//
//	response.WriteHeaderAndEntity(http.StatusOK, types.Response{Content: ""})
//}
//
//func (apiHandler *APIHandler) handleDelTraefikRoutings(request *restful.Request, response *restful.Response) {
//
//	// check token
//	if apiHandler.checkToken(request, response) == false {
//		return
//	}
//
//	// get body
//	json := &simplejson.Json{}
//	if err := request.ReadEntity(json); nil != err {
//		response.WriteHeaderAndEntity(http.StatusInternalServerError, types.Response{Content: err.Error()})
//		return
//	}
//	log.Print("handleDeleteRoutings: ", json)
//
//	rules, err := json.Get("rules").Array()
//	if nil != err {
//		response.WriteHeaderAndEntity(http.StatusInternalServerError, types.Response{Content: err.Error()})
//		return
//	}
//
//	for _, v := range (rules) {
//		body, _ := v.(map[string]interface{})
//		err := apiHandler.etcdHandler.handleDeleteRouting(&types.Request{
//			Name: body["Name"].(string),
//		})
//		if err != nil {
//			response.WriteHeaderAndEntity(http.StatusInternalServerError, types.Response{Content: err.Error()})
//			return
//		}
//	}
//
//	response.WriteHeaderAndEntity(http.StatusOK, types.Response{Content: ""})
//}
//
//func (apiHandler *APIHandler) handleKeepCloudwareAlive(request *restful.Request, response *restful.Response) {
//
//	body := &types.Request{}
//	if err := request.ReadEntity(body); nil != err {
//		response.WriteHeaderAndEntity(http.StatusInternalServerError, types.Response{Content: err.Error()})
//		return
//	}
//
//	log.Print("handleKeepAlive: ", body)
//
//	if err := apiHandler.etcdHandler.handleKeepCloudwareAlive(body); err != nil {
//		response.WriteHeaderAndEntity(http.StatusInternalServerError, types.Response{Content: err.Error()})
//		return
//	}
//
//	response.WriteHeaderAndEntity(http.StatusOK, types.Response{Content: ""})
//
//}
//
//func (apiHandler *APIHandler) checkToken(request *restful.Request, response *restful.Response) bool {
//	token := request.HeaderParameter("Authorization")
//	if strings.Compare(token, configs.Token) != 0 {
//		response.WriteHeaderAndEntity(http.StatusUnauthorized, types.Response{Content: ""})
//		return false
//	}
//	return true
//}
