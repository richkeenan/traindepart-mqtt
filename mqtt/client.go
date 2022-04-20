package mqtt

import (
	"fmt"
	paho "github.com/eclipse/paho.mqtt.golang"
	"log"
)

type Client struct {
	topicPrefix string
	pahoClient  paho.Client

	logger *log.Logger
}

func New(broker string, port int, username string, password string, topicPrefix string, logger *log.Logger) (*Client, error) {
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
		logger:      logger,
	}

	return c, nil
}

func (c *Client) Send(suffix, payload interface{}) {
	topic := fmt.Sprintf("%s/%s", c.topicPrefix, suffix)
	c.logger.Printf("Sending message on topic %s", topic)
	c.pahoClient.Publish(topic, 0, false, payload)
}
