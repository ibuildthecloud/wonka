package main

import (
	"fmt"
	"os"
	"strings"
)

var CONFIG = "${HOME}/.wonka/apps/%s"

func getHome() string {
	home := os.Getenv("WONKA_HOME")
	if home == "" {
		home = os.ExpandEnv("${HOME}/.wonka")
	}
	return home
}

func GetConfigPath(dir, name string) string {
	return fmt.Sprintf("%s/%s/%s", getHome(), dir, strings.ToLower(name))
}
