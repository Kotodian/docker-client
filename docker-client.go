package main

import (
	docker "github.com/fsouza/go-dockerclient"
	"strings"
)

var dockerClient *docker.Client

func init() {
	var err error
	dockerClient, err = docker.NewClient(*unixSock)
	if err != nil {
		panic(err)
	}
}

func images() ([]docker.APIImages, error) {
	return dockerClient.ListImages(docker.ListImagesOptions{
		All:     true,
		Digests: false,
		Filter:  "",
	})
}

func pullImage(imageTag string) error {
	repositoryTag := strings.Split(imageTag, ":")
	repository, tag := repositoryTag[0], repositoryTag[1]
	return dockerClient.PullImage(docker.PullImageOptions{
		All:        false,
		Repository: repository,
		Tag:        tag,
	}, docker.AuthConfiguration{})
}

func tagImage(imageTag string, newTag string) error {
	return dockerClient.TagImage(imageTag,
		docker.TagImageOptions{Tag: newTag},
	)
}

func pushImage(imageTag string) error {
	nameTag := strings.Split(imageTag, ":")
	return dockerClient.PushImage(docker.PushImageOptions{
		Name: nameTag[0],
		Tag:  nameTag[1],
	}, docker.AuthConfiguration{})
}
