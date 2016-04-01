package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/Sirupsen/logrus"
	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/project"
)

func main() {
	if err := run(); err != nil {
		logrus.Fatal(err)
	}
}

func run() error {
	null := make(chan project.ProjectEvent)
	go func() {
		for {
			<-null
		}
	}()

	prog, context, err := CreateContext(os.Args...)
	if context == nil {
		return fmt.Errorf("Failed to find program for %s", os.Args)
	}

	if err != nil {
		logrus.Fatal(err)
	}

	p, err := docker.NewProject(context)
	if err != nil {
		return err
	}

	p.AddListener(null)

	service, err := p.CreateService(strings.ToLower(prog))
	if err != nil {
		return err
	}

	if err := p.Create(); err != nil {
		return err
	}

	containers, err := service.Containers()
	if err != nil {
		return err
	}

	if len(containers) == 0 {
		return fmt.Errorf("Failed to find containers for service: %s", prog)
	}

	id, err := containers[0].Id()
	if err != nil {
		return err
	}

	services := []string{}

	for k, _ := range p.Configs {
		if k != prog {
			services = append(services, k)
		}
	}

	if len(services) > 0 {
		if err := p.Up(services...); err != nil {
			return err
		}
	}

	// Because attaching is rocket science....
	return syscall.Exec("/usr/bin/docker", []string{"/usr/bin/docker", "start", "-ai", id}, os.Environ())
}
