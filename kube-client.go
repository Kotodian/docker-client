package main

import (
	"context"
	"errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
)

var devKubeClient *kubernetes.Clientset
var uatKubeClient *kubernetes.Clientset

type deployment struct {
	// name deployment name like core-gw || ac-ocpp || admin
	name string
	// image is the docker image
	image string
	// uat is uat environment?
	uat bool
}

func newDeployment(name string, image string) *deployment {
	uat := false
	tag := strings.Split(image, ":")[1]
	if strings.Split(tag, "_")[1] == "rc" {
		uat = true
	}
	return &deployment{
		name:  name,
		image: image,
		uat:   uat,
	}
}

func updateDeployment(deploy *deployment) error {
	var client *kubernetes.Clientset
	if deploy.uat {
		if uatKubeClient == nil {
			return errors.New("no uat kube client")
		}
		client = uatKubeClient
	} else {
		client = devKubeClient
	}
	ctx := context.Background()
	// get the old deployment
	oldDeployment, err := client.AppsV1().Deployments("default").Get(ctx, deploy.name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	// deepcopy a new one
	newDeployment := oldDeployment.DeepCopy()
	// update the deployment
	oldImage := newDeployment.Spec.Template.Spec.Containers[0].Image
	if oldImage != deploy.image {
		if !strings.Contains(oldImage, deploy.name) {
			return errors.New("this image doesn't belong to this deployment")
		}
		newDeployment.Spec.Template.Spec.Containers[0].Image = deploy.image
		_, err = client.AppsV1().Deployments("default").Update(
			context.Background(),
			newDeployment, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}
