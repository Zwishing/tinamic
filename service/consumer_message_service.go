package service

import (
	"fmt"
	"github.com/nsqio/go-nsq"
	message "tinamic/pkg/nsq"
)

type ConsumerMinioAdd struct {
}

func NewConsumerMinioAdd() *ConsumerMinioAdd {
	return &ConsumerMinioAdd{}
}

// HandleMessage 是需要实现的处理消息的方法
func (m *ConsumerMinioAdd) HandleMessage(msg *nsq.Message) (err error) {
	fmt.Println(msg.ID, msg.NSQDAddress, string(msg.Body))
	//m.repo.SaveDataSource()
	fmt.Println(11111)
	return nil
}

func (m *ConsumerMinioAdd) ConsumerMinioAdd() {
	c := &message.NSQConfig{
		Topic:       "add_data",
		Channel:     "first",
		Address:     "1.92.113.25:4161",
		MaxAttempts: 5,
	}
	consumer, err := message.NewNSQConsumer(c, m.HandleMessage)
	if err != nil {
		return
	}
	//consumer.Stop()
	fmt.Println(consumer)
}
