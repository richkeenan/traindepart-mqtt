package rtt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	HttpClient *http.Client
	username   string
	password   string
}

func New(username, password string) *Client {
	client := http.Client{}
	return &Client{
		HttpClient: &client,
		username:   username,
		password:   password,
	}
}

// rttapi_richkeenan
//  os.Getenv("PASSWORD")

func (c *Client) DoIt(from, to string) (*SearchResult, error) {
	url := fmt.Sprintf("https://api.rtt.io/api/v1/json/search/%s/to/%s", from, to)
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(c.username, c.password)

	resp, err := c.HttpClient.Do(req)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var searchResult SearchResult
	err = json.Unmarshal(body, &searchResult)
	if err != nil {
		return nil, err
	}

	return &searchResult, nil
}
