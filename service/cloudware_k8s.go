package service

import (
	"github.com/kfcoding-cloudware-controller/types"
	"log"
	"github.com/kfcoding-cloudware-controller/template"
	"github.com/kfcoding-cloudware-controller/configs"
	"github.com/satori/go.uuid"
	"encoding/json"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/typed/core/v1"
	v12 "k8s.io/api/core/v1"
	v13 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CloudwareK8sService struct {
	Keeper           KeeperService
	Routing          RoutingService
	PodInterface     v1.PodInterface
	ServiceInterface v1.ServiceInterface
	K8sClient        *kubernetes.Clientset
}

func GetCloudwareK8sService(keeper KeeperService, routing RoutingService) *CloudwareK8sService {

	cfg, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal("Could not init in cluster config: ", err.Error())
	}
	K8sClient, err := kubernetes.NewForConfig(cfg)
	PodInterface := K8sClient.CoreV1().Pods(configs.Namespace)
	ServiceInterface := K8sClient.CoreV1().Services(configs.Namespace)

	return &CloudwareK8sService{
		Keeper:           keeper,
		Routing:          routing,
		K8sClient:        K8sClient,
		PodInterface:     PodInterface,
		ServiceInterface: ServiceInterface,
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

func (service *CloudwareK8sService) CreateCloudwareService(name string) (*v12.Service, error) {
	var serviceBody v12.Service
	err := json.Unmarshal([]byte(template.CloudwareService), &serviceBody)
	if nil != err {
		log.Print("CreateCloudwareService error: ", err)
		return nil, err
	}

	serviceBody.Name = name
	serviceBody.Namespace = configs.Namespace
	serviceBody.Labels["app"] = name
	serviceBody.Spec.Selector["app"] = name

	v1Service, err := service.ServiceInterface.Create(&serviceBody)

	if nil != err {
		log.Print("CreateCloudwareService error: ", err)
		return nil, err
	} else {
		log.Printf("CreateCloudwareService ok")
		return v1Service, nil
	}
}

func (service *CloudwareK8sService) CreateCloudwarePod(body *types.CloudwareBody) (string, error) {
	var podBody v12.Pod
	err := json.Unmarshal([]byte(template.CloudwarePod), &podBody)
	if nil != err {
		log.Print("CreateCloudwarePod error: ", err)
		return "", err
	}
	var name = "cloudware-" + uuid.Must(uuid.NewV4()).String()
	podBody.Name = name
	podBody.Namespace = configs.Namespace
	podBody.Labels["app"] = name
	podBody.Spec.Containers[1].Image = body.Image

	_, err = service.PodInterface.Create(&podBody)

	if nil != err {
		log.Print("CreateCloudwarePod error: ", err)
		return "", err
	} else {
		log.Printf("CreateCloudwarePod ok")
		return name, nil
	}

}

func (service *CloudwareK8sService) DeleteCloudwareService(name string) (string, error) {

	racePeriodSeconds := int64(0)
	var propagationPolicy v13.DeletionPropagation
	propagationPolicy = "Background"

	options := &v13.DeleteOptions{
		TypeMeta: v13.TypeMeta{
			Kind:       "DeleteOptions",
			APIVersion: "v1",
		},
		GracePeriodSeconds: &racePeriodSeconds,
		PropagationPolicy:  &propagationPolicy,
	}
	err := service.ServiceInterface.Delete(name, options)

	if nil != err {
		log.Print("DeleteCloudwareService error: ", err)
		return "", err
	} else {
		log.Printf("DeleteCloudwareService ok")
		return "", nil
	}
}

func (service *CloudwareK8sService) DeleteCloudwarePod(name string) (string, error) {

	racePeriodSeconds := int64(0)
	var propagationPolicy v13.DeletionPropagation
	propagationPolicy = "Background"

	options := &v13.DeleteOptions{
		TypeMeta: v13.TypeMeta{
			Kind:       "DeleteOptions",
			APIVersion: "v1",
		},
		GracePeriodSeconds: &racePeriodSeconds,
		PropagationPolicy:  &propagationPolicy,
	}
	err := service.PodInterface.Delete(name, options)

	if nil != err {
		log.Print("DeleteCloudwarePod error: ", err)
		return "", err
	} else {
		log.Printf("DeleteCloudwarePod ok")
		return "", nil
	}
}
