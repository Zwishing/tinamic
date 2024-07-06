// nsq/consumer.go
package nsq

import (
	"github.com/nsqio/go-nsq"
)

type NSQConsumer struct {
	Consumer *nsq.Consumer
}

type NSQConfig struct {
	Topic       string
	Channel     string
	Address     string
	MaxAttempts uint16
}

func NewNSQConfig() *NSQConfig {
	return &NSQConfig{
		Topic:       "add_data",
		Channel:     "channel1",
		Address:     "1.92.113.25:4161",
		MaxAttempts: 10,
	}
}

func NewNSQConsumer(config *NSQConfig, handler nsq.Handler) (*NSQConsumer, error) {
	consumer, err := nsq.NewConsumer(config.Topic, config.Channel, nsq.NewConfig())
	if err != nil {
		return nil, err
	}

	consumer.AddHandler(handler)

	if err := consumer.ConnectToNSQLookupd(config.Address); err != nil {
		return nil, err
	}

	return &NSQConsumer{Consumer: consumer}, nil
}

func (c *NSQConsumer) Stop() {

	c.Consumer.Stop()
}
