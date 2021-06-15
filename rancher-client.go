package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/rancher/norman/clientbase"
	"github.com/rancher/norman/types"
	rancherClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"os"
)

var devRancherClient *rancherClient.Client

func init() {
	var err error
	devRancherClient, err = rancherClient.NewClient(&clientbase.ClientOpts{
		URL:      os.Getenv("DEV_RANCHER_API"),
		TokenKey: os.Getenv("DEV_RANCHER_BEAR_TOKEN"),
	})
	if err != nil {
		panic(err)
	}
}

func clusters() (*rancherClient.ClusterCollection, error) {
	return devRancherClient.Cluster.List(&types.ListOpts{})
}

func projects() (*rancherClient.ProjectCollection, error) {
	return devRancherClient.Project.List(&types.ListOpts{})
}

func workload(client *rancherClient.ProjectCollection) {

}
