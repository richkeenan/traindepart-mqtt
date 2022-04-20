package main

import (
	"encoding/json"
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/traindepartmqtt/mqtt"
	"github.com/traindepartmqtt/rtt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	crontab    *cron.Cron
	logger     *log.Logger
	specParser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month)
)

func main() {
	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Println("Starting TrainDepartMQTT")

	logger.Println("Loading configuration")
	cfg, err := getConfig()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	rttClient := rtt.New(cfg.Rtt.Username, cfg.Rtt.Password)

	logger.Printf("Polling for departures from %s to %s", cfg.Rtt.StationFrom, cfg.Rtt.StationTo)

	logger.Println("Connecting to MQTT broker")
	mqttClient, err := mqtt.New(cfg.Mqtt.Broker, cfg.Mqtt.Port, cfg.Mqtt.Username, cfg.Mqtt.Password, cfg.Mqtt.TopicPrefix, logger)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	logger.Println("Connected to MQTT broker")

	crontab = cron.New()
	defer crontab.Stop()
	// every minute
	schedule, err := specParser.Parse("0 * * ? *")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	crontab.Schedule(schedule, cron.FuncJob(func() {
		err = sendDepartures(rttClient, cfg.Rtt.StationFrom, cfg.Rtt.StationTo, mqttClient)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}))
	crontab.Start()

	// Perform task immediately
	err = sendDepartures(rttClient, cfg.Rtt.StationFrom, cfg.Rtt.StationTo, mqttClient)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	wait()
}

func sendDepartures(rttClient *rtt.Client, stationFrom, stationTo string, mqttClient *mqtt.Client) error {
	logger.Println("Searching for departures")
	result, err := rttClient.Search(stationFrom, stationTo)
	if err != nil {
		return err
	}

	departures, err := toDepartures(result)
	if err != nil {
		return err
	}

	payload, _ := json.Marshal(departures)
	mqttClient.Send("trainstatus", string(payload))
	return nil
}

func wait() {
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
}
