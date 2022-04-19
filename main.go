package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"trainstatus/rtt"
)

func main() {
	url := "https://api.rtt.io/api/v1/json/search/BAL/to/VIC"
	client := http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth("rttapi_richkeenan", os.Getenv("PASSWORD"))

	resp, err := client.Do(req)

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var searchResult rtt.SearchResult
	err = json.Unmarshal(body, &searchResult)
	if err != nil {
		panic(err.Error())
	}
}
