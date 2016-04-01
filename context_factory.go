package main

import (
	"os"
	"path"

	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/project"
)

var APP_NAME = "wonka"

func CreateContext(args ...string) (string, *docker.Context, error) {
	prog, args := parseArgs(args...)

	if prog == "" {
		return "", nil, nil
	}

	configFile, err := findConfig(prog)
	if err != nil || configFile == "" {
		return "", nil, err
	}

	serviceFactory := &ServiceFactory{
		prog: prog,
		args: args,
	}
	context := &docker.Context{
		Context: project.Context{
			Rebuild:           true,
			ComposeFile:       configFile,
			ProjectName:       prog,
			ServiceFactory:    serviceFactory,
			EnvironmentLookup: &ConfigEnvironment{},
		},
	}

	context.Builder = NewBuilder(context)
	serviceFactory.Context = context

	return prog, context, nil
}

func findConfig(prog string) (string, error) {
	var err error
	configFile := GetConfigPath("apps", prog)
	if stat, err := os.Stat(configFile); os.IsNotExist(err) || stat.IsDir() {
		configFile += ".yml"
		if _, err = os.Stat(configFile); os.IsNotExist(err) {
			return "", nil
		}
	}

	return configFile, err

}

func parseArgs(args ...string) (string, []string) {
	prog := path.Base(args[0])
	if prog == APP_NAME {
		if len(args) < 2 {
			return "", nil
		}

		prog = path.Base(args[1])
		args = args[1:]
	}

	return prog, append([]string{prog}, args[1:]...)
}
