package main

import (
	"log"
	"os"
	"time"

	"github.com/getcarina/libcarina"
	"github.com/samalba/dockerclient"
)

func main() {
	var err error

	username := os.Getenv("CARINA_USERNAME")
	apiKey := os.Getenv("CARINA_APIKEY")
	clusterName := os.Args[1]

	// Connect to Carina
	cli, _ := libcarina.NewClusterClient(libcarina.BetaEndpoint, username, apiKey)
	cluster, err := cli.Get(clusterName)
	if err != nil {
		panic(err)
	}

	// Wait for it to come up...
	for cluster.Status == "new" || cluster.Status == "building" {
		time.Sleep(10 * time.Second)
		cluster, err = cli.Get(clusterName)
		if err != nil {
			break
		}
	}

	// Get the IP of the host and a *tls.Config
	host, tlsConfig, _ := cli.GetDockerConfig(clusterName)

	// Straight to Docker, do what you need
	docker, _ := dockerclient.NewDockerClient(host, tlsConfig)

	// Create a container
	containerConfig := &dockerclient.ContainerConfig{
		Image:       "whoa/tiny",
		AttachStdin: true,
		Tty:         true,
	}

	containerID, err := docker.CreateContainer(containerConfig, "foobar", nil)
	if err != nil {
		log.Fatal(err)
	}

	// Start the container
	hostConfig := &dockerclient.HostConfig{
		PublishAllPorts: true,
	}
	err = docker.StartContainer(containerID, hostConfig)
	if err != nil {
		log.Fatal(err)
	}

}
