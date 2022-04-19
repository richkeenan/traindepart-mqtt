package main

import (
	"fmt"
	"os"
	"trainstatus/rtt"
)

func main() {
	cfg, err := getConfig()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	rttClient := rtt.New(cfg.Rtt.Username, cfg.Rtt.Password)
	result, err := rttClient.Search("BAL", "VIC")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%v\n", result)
}
