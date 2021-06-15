package main

import (
	"flag"
	"fmt"
)

var (
	imageTag = flag.String("image", "", "image:tag")
	unixSock = flag.String("sock", "unix://var/run/docker.sock", "docker sock")
	uat      = flag.Bool("uat", false, "need to deploy in uat")
)

func main() {
	flag.Parse()
	collection, err := clusters()
	if err != nil {
		panic(err)
	}
	fmt.Println(collection.Data[0].ID)
}
