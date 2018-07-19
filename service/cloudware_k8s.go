package service

import (
	"github.com/cloudware-controller/types"
	"k8s.io/client/kubernetes/config"
	"k8s.io/client/kubernetes/client"
	"log"
	"github.com/cloudware-controller/template"
	"github.com/cloudware-controller/configs"
	"github.com/satori/go.uuid"
	"encoding/json"
	"context"
)

type CloudwareK8sService struct {
	Keeper    KeeperService
	Routing   RoutingService
	K8sClient *client.APIClient
}

func GetCloudwareK8sService(keeper KeeperService, routing RoutingService) *CloudwareK8sService {

	c, err := config.LoadKubeConfig()
	if err != nil {
		panic(err.Error())
	}

	return &CloudwareK8sService{
		Keeper:    keeper,
		Routing:   routing,
		K8sClient: client.NewAPIClient(c),
	}
}

func (service *CloudwareK8sService) WatcherCallback(body *types.KeeperBody) {
	log.Printf("WatcherCallback: %+v\n", body)

	service.DeleteCloudwareService(body.Name)

	service.DeleteCloudwarePod(body.Name)

	service.Routing.DeleteRule(&types.RoutingBody{Name: body.Name})

	log.Printf("WatcherCallback ok: %+v\n", body)
}

func (service *CloudwareK8sService) CreateCloudwareApi(body *types.CloudwareBody) (string, error) {
	name, err := service.CreateCloudwarePod(body)
	if nil != err {
		return "", err
	}
	v1Service, err := service.CreateCloudwareService(name)
	if nil != err {
		return "", err
	}
	service.Keeper.Keep(&types.KeeperBody{Name: name})
	wsAddr := name + "." + configs.WsAddrSuffix
	service.Routing.AddRule(&types.RoutingBody{
		Name: name,
		URL:  "http://" + v1Service.Spec.ClusterIP + ":9800",
		Rule: "Host: " + wsAddr,
	})
	return wsAddr, nil
}

func (service *CloudwareK8sService) CreateCloudwareService(name string) (*client.V1Service, error) {
	var serviceBody client.V1Service
	err := json.Unmarshal([]byte(template.CloudwareService), &serviceBody)
	if nil != err {
		log.Print("CreateCloudwareService error: ", err)
		return nil, err
	}
	serviceBody.Metadata.Name = name
	serviceBody.Metadata.Namespace = configs.Namespace
	serviceBody.Metadata.Labels["app"] = name
	serviceBody.Spec.Selector["app"] = name

	v1Service, _, err := service.K8sClient.CoreV1Api.CreateNamespacedService(
		context.Background(), configs.Namespace,
		serviceBody, nil)

	if nil != err {
		log.Print("CreateCloudwareService error: ", err)
		return nil, err
	} else {
		log.Printf("CreateCloudwareService ok: %+v\n", v1Service)
		return &v1Service, nil
	}
}

func (service *CloudwareK8sService) CreateCloudwarePod(body *types.CloudwareBody) (string, error) {
	var podBody client.V1Pod
	err := json.Unmarshal([]byte(template.CloudwarePod), &podBody)
	if nil != err {
		log.Print("CreateCloudwarePod error: ", err)
		return "", err
	}
	var name = "cloudware-" + uuid.Must(uuid.NewV4()).String()
	podBody.Metadata.Name = name
	podBody.Metadata.Namespace = configs.Namespace
	podBody.Metadata.Labels["app"] = name
	podBody.Spec.Containers[1].Image = body.Image

	pod, _, err := service.K8sClient.CoreV1Api.CreateNamespacedPod(
		context.Background(), configs.Namespace,
		podBody, nil)

	if nil != err {
		log.Print("CreateCloudwarePod error: ", err)
		return "", err
	} else {
		log.Printf("CreateCloudwarePod ok: %+v\n", pod)
		return name, nil
	}

}

func (service *CloudwareK8sService) DeleteCloudwareService(serviceName string) (string, error) {
	res, _, err := service.K8sClient.CoreV1Api.DeleteNamespacedService(
		context.Background(),
		serviceName,
		configs.Namespace,
		client.V1DeleteOptions{
			ApiVersion:         "v1",
			Kind:               "DeleteOptions",
			GracePeriodSeconds: 0,
			OrphanDependents:   false,
			PropagationPolicy:  "Background",
		},
		nil)

	if nil != err {
		log.Print("DeleteCloudwareService error: ", err)
		return "", err
	} else {
		log.Printf("DeleteCloudwareService ok: %+v\n", res)
		return "", nil
	}
}

func (service *CloudwareK8sService) DeleteCloudwarePod(podName string) (string, error) {
	res, _, err := service.K8sClient.CoreV1Api.DeleteNamespacedPod(
		context.Background(),
		podName,
		configs.Namespace,
		client.V1DeleteOptions{
			ApiVersion:         "v1",
			Kind:               "DeleteOptions",
			GracePeriodSeconds: 0,
			OrphanDependents:   false,
			PropagationPolicy:  "Background",
		},
		nil)

	log.Printf("DeleteCloudwarePod ok: %+v %+v\n", res, err)
	return "", err

}
