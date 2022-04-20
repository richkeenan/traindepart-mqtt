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
		Username    string `env:"RTT_USERNAME" cli:"rtt-username" yaml:"username"`
		Password    string `env:"RTT_PASSWORD" cli:"rtt-password" yaml:"password"`
		StationFrom string `env:"RTT_FROM" cli:"rtt-from" yaml:"station_from"`
		StationTo   string `env:"RTT_TO" cli:"rtt-to" yaml:"station_to"`
	}
	Mqtt struct {
		Broker      string `env:"MQTT_BROKER" cli:"mqtt-broker" yaml:"broker"`
		Port        int    `env:"MQTT_PORT" cli:"mqtt-port" yaml:"port"`
		Username    string `env:"MQTT_USERNAME" cli:"mqtt-username" yaml:"username"`
		Password    string `env:"MQTT_PASSWORD" cli:"mqtt-password" yaml:"password"`
		TopicPrefix string `env:"MQTT_TOPIC_PREFIX" cli:"mqtt-topic-prefix" yaml:"topic_prefix"`
	}
}

func getConfig() (*Config, error) {
	var cfg Config
	cfg.Mqtt.TopicPrefix = "traindepart"
	cfg.Mqtt.Port = 1883

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
	if cfg.Rtt.StationFrom == "" {
		return nil, fmt.Errorf("rtt station from not set")
	}
	if cfg.Rtt.StationTo == "" {
		return nil, fmt.Errorf("rtt station to not set")
	}

	return &cfg, nil
}
