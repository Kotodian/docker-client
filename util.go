package main

import (
	"errors"
	"path"
	"strings"
)

func getDeploymentName(imageTag string) string {
	image := strings.Split(imageTag, ":")[0]
	_, deployment := path.Split(image)
	return deployment
}

func failOnError(message string, err error) {
	if err != nil {
		panic(errors.New(message + err.Error()))
	}
}
