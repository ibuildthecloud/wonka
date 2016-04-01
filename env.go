package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/docker/libcompose/project"
)

type ConfigEnvironment struct {
}

func appendEnv(array []string, key, value string) []string {
	parts := strings.SplitN(key, "/", 2)
	if len(parts) == 2 {
		key = parts[1]
	}

	return append(array, fmt.Sprintf("%s=%s", key, value))
}

func lookupKeys(env map[string]string, keys ...string) []string {
	for _, key := range keys {
		if strings.HasSuffix(key, "*") {
			result := []string{}
			for envKey, envValue := range env {
				keyPrefix := key[:len(key)-1]
				if strings.HasPrefix(envKey, keyPrefix) {
					result = appendEnv(result, envKey, envValue)
				}
			}

			if len(result) > 0 {
				return result
			}
		} else if value, ok := env[key]; ok {
			return appendEnv([]string{}, key, value)
		}
	}

	return []string{}
}

func (c *ConfigEnvironment) Lookup(key, serviceName string, serviceConfig *project.ServiceConfig) []string {
	env := map[string]string{}
	for _, val := range os.Environ() {
		parts := strings.SplitN(val, "=", 2)
		if len(parts) == 2 {
			env[parts[0]] = parts[1]
		}
	}
	result := lookupKeys(env, key)
	sort.Strings(result)
	return result
}
