package main

import (
	"k8s.io/client/kubernetes/client"
	"k8s.io/client/kubernetes/config"
	"github.com/cloudware-controller/configs"
	"context"
	"log"
)

func main() {
	c, err := config.LoadKubeConfig()
	if err != nil {
		panic(err.Error())
	}
	K8sClient := client.NewAPIClient(c)

	body := client.V1DeleteOptions{
		ApiVersion:         "v1",
		Kind:               "DeleteOptions",
		GracePeriodSeconds: 0,
		OrphanDependents:   false,
		PropagationPolicy:  "Background",
	}
	res, resp, err := K8sClient.CoreV1Api.DeleteNamespacedPod(
		context.Background(),
		"cloudware-1-7987f857f9-jqf26",
		configs.Namespace,
		body, nil)

	log.Print("Kind ", res.Kind)
	log.Print("ApiVersion ", res.ApiVersion)
	log.Print("Status ", res.Status)
	log.Print("Metadata ", res.Metadata)
	log.Print("Code ", res.Code)
	log.Print("Message ", res.Message)
	log.Print("Reason ", res.Reason)
	log.Print("Details ", res.Details)
	log.Print(res)
	log.Print(resp)
	log.Print(err)
	//var podBody client.V1Pod
	//err = json.Unmarshal([]byte(template.CloudwarePod), &podBody)
	//if nil != err {
	//	log.Print("CreateCloudwareService: ", err)
	//}
	//var name = "cloudware-12345"
	//podBody.Metadata.Name = name
	//podBody.Metadata.Namespace = configs.Namespace
	//podBody.Metadata.Labels["app"] = name
	//podBody.Spec.Containers[1].Image = "daocloud.io/shaoling/kfcoding-rstudio-latest:master"
	//
	//log.Print(podBody)
	//log.Print(podBody.Spec)
	//log.Print(podBody.Metadata)
	//log.Print(podBody.Kind)
	//log.Print(podBody.ApiVersion)
	//pod, res, err1 := K8sClient.CoreV1Api.CreateNamespacedPod(
	//	context.Background(),
	//	"kfcoding-alpha",
	//	podBody, nil)
	//
	//log.Print(pod)
	//log.Print(err1)
	//log.Print(res)
}
