package nsq

import (
	"github.com/nsqio/go-nsq"
	"strconv"
)

// NsqConfig
type NsqConfig struct {
	addr string
	*nsq.Config
}

type NsqProducer struct {
	*nsq.Producer
}

type NsqMessage struct {
	*nsq.Message
}

func NewNsqConfig(host string, port int) *NsqConfig {
	return &NsqConfig{
		addr:   host + ":" + strconv.Itoa(port),
		Config: nsq.NewConfig(),
	}
}

func NewNsqProducer(config *NsqConfig) (*NsqProducer, error) {
	producer, err := nsq.NewProducer(config.addr, config.Config)
	return &NsqProducer{
		producer,
	}, err
}

func NewNsqConsumer(topic string, channel string, handler nsq.Handler, config *NsqConfig) error {
	consumer, err := nsq.NewConsumer(topic, channel, config.Config)
	if err != nil {
		return err
	}
	consumer.AddHandler(handler)
	if err := consumer.ConnectToNSQLookupd(config.addr); err != nil { // 通过lookupd查询
		return err
	}
	return nil
}
