package main

import (
	"path"
	"strings"
)

func getDeploymentName(imageTag string) string {
	image := strings.Split(imageTag, ":")[0]
	_, deployment := path.Split(image)
	return deployment
}

func failOnError(err error) {
	if err != nil {
		panic(err)
	}
}
