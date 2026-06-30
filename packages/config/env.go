package config

import (
	"os"
	"strings"
)

type EnvProvider struct{}

func NewEnvProvider() *EnvProvider {
	return &EnvProvider{}
}

func (e *EnvProvider) Name() string { return "env" }

func (e *EnvProvider) Load() (map[string]interface{}, error) {
	data := make(map[string]interface{})
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			data[parts[0]] = parts[1]
		}
	}
	return data, nil
}
