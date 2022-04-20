package main

import (
	"encoding/json"
	"fmt"
	"github.com/traindepartmqtt/mqtt"
	"github.com/traindepartmqtt/rtt"
	"os"
)

func main() {
	cfg, err := getConfig()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	rttClient := rtt.New(cfg.Rtt.Username, cfg.Rtt.Password)
	result, err := rttClient.Search(cfg.Rtt.StationFrom, cfg.Rtt.StationTo)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	mqttClient, err := mqtt.New(cfg.Mqtt.Broker, cfg.Mqtt.Port, cfg.Mqtt.Username, cfg.Mqtt.Password, cfg.Mqtt.TopicPrefix)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	departures, err := toDepartures(result)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	payload, _ := json.Marshal(departures)
	mqttClient.Send("trainstatus", string(payload))
}
