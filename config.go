package main

import (
	"errors"
	"fmt"
	"github.com/gravitational/configure"
	"io/ioutil"
	"os"
)

type Config struct {
	Rtt struct {
		Username string `env:"RTT_USERNAME" cli:"rtt-username" yaml:"username"`
		Password string `env:"RTT_PASSWORD" cli:"rtt-password" yaml:"password"`
	}
}

func getConfig() (*Config, error) {
	var cfg Config
	err := configure.ParseEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading environment variables, %v", err)
	}

	configFile := "/config/config.yaml" // When running in Docker
	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		configFile = "./config.yaml" // When running locally
	}
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("error reading YAML file, %v", err)
	}

	err = configure.ParseYAML(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("error parsing YAML file, %v", err)
	}

	err = configure.ParseCommandLine(&cfg, os.Args[1:])
	if err != nil {
		return nil, fmt.Errorf("error parsing command line, %v", err)
	}

	if cfg.Rtt.Username == "" {
		return nil, fmt.Errorf("rtt username not set")
	}
	if cfg.Rtt.Password == "" {
		return nil, fmt.Errorf("rtt password not set")
	}

	return &cfg, nil
}
