package main

import (
	"time"
	"trainstatus/rtt"
)

// TODO timezones?!?
func toDepartures(result *rtt.SearchResult) ([]Departure, error) {
	departures := make([]Departure, 0)
	now := time.Now()
	for _, service := range result.Services {
		plannedDepartureTime, err := time.Parse("1504", service.LocationDetail.GbttBookedDeparture)
		if err != nil {
			return nil, err
		}
		plannedDepartureTime = plannedDepartureTime.AddDate(now.Year(), int(now.Month()), now.Day())

		expectedDepartureTime, err := time.Parse("1504", service.LocationDetail.RealtimeDeparture)
		if err != nil {
			return nil, err
		}
		expectedDepartureTime = expectedDepartureTime.AddDate(now.Year(), int(now.Month()), now.Day())

		departures = append(departures, Departure{
			ScheduledTime: plannedDepartureTime,
			ExpectedTime:  expectedDepartureTime,
		})
	}

	return departures, nil
}

type Departure struct {
	ScheduledTime time.Time `json:"scheduledTime"`
	ExpectedTime  time.Time `json:"expectedTime"` // accounts for delays
}
