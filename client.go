package main

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/project"
	fDockerClient "github.com/fsouza/go-dockerclient"
	"github.com/samalba/dockerclient"
)

func GetClient(context *docker.Context, service project.Service) (*fDockerClient.Client, error) {
	samsClient := context.ClientFactory.Create(service)
	if impl, ok := samsClient.(*dockerclient.DockerClient); ok {
		client, err := fDockerClient.NewClient("unix:///var/run/docker.sock")
		if err != nil {
			return nil, err
		}

		client.TLSConfig = impl.TLSConfig
		return client, err
	} else {
		return nil, fmt.Errorf("Invalid docker client")
	}
}

func Attach(id string, client *fDockerClient.Client) (<-chan bool, error) {
	result := make(chan bool, 1)
	container, err := client.InspectContainer(id)
	if err != nil {
		return nil, err
	}

	opts := fDockerClient.AttachToContainerOptions{
		Container:    id,
		Stream:       true,
		InputStream:  os.Stdin,
		OutputStream: os.Stdout,
		ErrorStream:  os.Stderr,
		Stdin:        true,
		Stdout:       true,
		Stderr:       true,
	}

	if container.Config.Tty {
		opts.RawTerminal = true
	} else {
		opts.Logs = true
	}

	go func() {
		if err := client.AttachToContainer(opts); err != nil {
			logrus.Fatal(err)
		}
		result <- true
	}()

	return result, nil
}
