package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"

	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/project"
)

type Builder struct {
	*docker.DaemonBuilder
	context *docker.Context
}

func NewBuilder(context *docker.Context) *Builder {
	return &Builder{
		DaemonBuilder: docker.NewDaemonBuilder(context),
		context:       context,
	}
}

func (b *Builder) Build(p *project.Project, service project.Service) (string, error) {
	if service.Config().Build == "" {
		return service.Config().Image, nil
	}

	digest := sha256.New()
	tar, err := docker.CreateTar(p, service.Name())
	if err != nil {
		return "", err
	}
	defer tar.Close()

	_, err = io.Copy(digest, tar)
	if err != nil {
		return "", err
	}

	hexString := hex.EncodeToString(digest.Sum([]byte{}))
	cacheFile := GetConfigPath("cache", hexString)

	imageId := ""
	imageIdBytes, err := ioutil.ReadFile(cacheFile)
	if err == nil {
		client := b.context.ClientFactory.Create(service)
		if image, err := client.InspectImage(string(imageIdBytes)); err == nil && image != nil {
			imageId = string(imageIdBytes)
		}
	}

	if imageId == "" {
		imageId, err = b.DaemonBuilder.Build(p, service)
		if err != nil {
			return "", err
		}

		cacheDir := GetConfigPath("cache", "")
		if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
			if err := os.MkdirAll(cacheDir, 0755); err != nil {
				return "", err
			}
		}

		if err := ioutil.WriteFile(cacheFile, []byte(imageId), 0644); err != nil {
			return "", err
		}
	}

	return imageId, nil
}
