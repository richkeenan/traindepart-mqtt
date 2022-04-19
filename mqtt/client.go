package mqtt

import (
	"fmt"
	paho "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	topicPrefix string
	pahoClient  paho.Client
}

func New(broker string, port int, username string, password string, topicPrefix string) (*Client, error) {
	opts := paho.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	if username != "" {
		opts.SetUsername(username)
		opts.SetPassword(password)
	}

	p := paho.NewClient(opts)
	if token := p.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("error connecting to MQTT broker, %v", token.Error())
	}

	c := &Client{
		topicPrefix: topicPrefix,
		pahoClient:  p,
	}

	return c, nil
}

func (c *Client) Send(suffix, payload interface{}) {
	c.pahoClient.Publish(fmt.Sprintf("%s/%s", c.topicPrefix, suffix), 0, false, payload)
}
