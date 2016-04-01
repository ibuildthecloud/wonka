package main

import (
	"os"
	"strconv"

	"github.com/docker/docker/pkg/term"
	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/project"
)

type ServiceFactory struct {
	docker.ServiceFactory

	prog string
	args []string
}

func (s *ServiceFactory) Create(p *project.Project, name string, serviceConfig *project.ServiceConfig) (project.Service, error) {
	config := *serviceConfig
	if config.User == "${USER}" {
		config.User = strconv.Itoa(os.Getuid())
	}
	config.StdinOpen = true
	if _, isTerm := term.GetFdInfo(os.Stdin); isTerm {
		config.Tty = true
	}
	config.Command = project.NewCommand(s.args...)

	return s.ServiceFactory.Create(p, name, &config)
}
