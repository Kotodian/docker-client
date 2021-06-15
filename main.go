package main

import (
	"flag"
	docker "github.com/fsouza/go-dockerclient"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"path/filepath"
	"strings"
)

var (
	imageTag      = flag.String("image", "", "image:tag")
	unixSock      = flag.String("sock", "unix://var/run/docker.sock", "docker sock")
	uat           = flag.Bool("uat", false, "need to deploy in uat")
	kubeConfig    = flag.String("kube_conf", filepath.Join(homedir.HomeDir(), ".kube", "config"), "kube config")
	uatKubeConfig = flag.String("uat_kube_conf", filepath.Join(homedir.HomeDir(), "uat", "config"), "uat kube config")
)

func init() {
	flag.Parse()
	var err error
	dockerClient, err = docker.NewClient(*unixSock)
	failOnError("failed to connect docker daemon error: ", err)
	devConfig, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)
	failOnError("failed to build dev kube config error: ", err)
	devKubeClient, err = kubernetes.NewForConfig(devConfig)
	failOnError("failed to connect dev k8s master error: ", err)

	if *uat {
		uatConfig, err := clientcmd.BuildConfigFromFlags("", *uatKubeConfig)
		failOnError("failed to build uat kube config error: ", err)
		uatKubeClient, err = kubernetes.NewForConfig(uatConfig)
		failOnError("failed to connect uat k8s master error: ", err)
	}
}

func main() {
	log.Println("start pulling image ", *imageTag)
	err := pullImage(*imageTag)
	failOnError("pull image error", err)

	devDeploymentName := getDeploymentName(*imageTag)
	log.Println("start updating dev k8s deployment ", devDeploymentName)
	err = updateDeployment(newDeployment(devDeploymentName, *imageTag))
	failOnError("update dev k8s deployment error: ", err)

	if *uat {
		newImage := strings.ReplaceAll(*imageTag, "beta", "rc")
		log.Printf("start tagging image %s to a new image %s\n", *imageTag, newImage)
		err = tagImage(*imageTag, strings.Split(newImage, ":")[1])
		failOnError("tagging image error: ", err)

		log.Printf("start pushing image %s\n", newImage)
		err = pushImage(newImage)
		failOnError("push image error: ", err)

		log.Printf("start deleting image %s\n", newImage)
		err = deleteImage(newImage)
		failOnError("delete image error: ", err)

		uatDeploymentName := getDeploymentName(newImage)
		log.Println("start updating uat k8s deployment ", uatDeploymentName)
		err = updateDeployment(newDeployment(uatDeploymentName, newImage))
		failOnError("update uat k8s deployment error: ", err)
	}
}
